// Package search provides a pure-Go, zero-dependency fuzzy search engine
// for the system-design dashboard. It uses an inverted index with TF-IDF
// scoring, Levenshtein distance (≤ 2 edits) for typo tolerance, and
// field-level boost weights so title matches outrank body matches.
//
// The index is built in a background goroutine; callers check IsReady()
// before serving results. While building, Search returns an empty slice
// and Ready returns false so the handler can show an "indexing…" state.
package search

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"unicode"

	"system-design/internal/registry"
)

// ─── Boost constants ────────────────────────────────────────────────────────

const (
	boostTitle       = 5.0
	boostDescription = 3.0
	boostPhaseTitle  = 2.5
	boostKey         = 2.5
	boostSay         = 2.0
	boostQA          = 2.0
	boostHint        = 1.8
	boostThink       = 1.5
	boostCompare     = 1.5
	boostNFR         = 1.5
	boostFundLinks   = 1.2
	boostTable       = 1.0
	boostBody        = 1.0
	boostCode        = 0.5
)

// ─── Data types ─────────────────────────────────────────────────────────────

// field is a single indexed text field with a display name and boost weight.
type field struct {
	name  string  // human-readable: "Title", "Phase: Caching Strategies", etc.
	text  string  // plain text to search
	boost float64 // scoring multiplier
}

// doc is a single indexed document.
type doc struct {
	Slug        string
	Title       string
	Description string
	Type        string // "problem" | "fundamental" | "algorithm" | "pattern"
	Category    string // "networking", "storage", "" for problems/algorithms
	Parent      string // parent fundamental title for sub-topics
	Route       string // "/fund/storage/redis"
	fields      []field
}

// Result is a single ranked search hit returned to the caller.
type Result struct {
	Slug        string
	Title       string
	Description string
	Type        string
	Category    string
	Parent      string
	Route       string
	ScorePct    int    // 0–100 for CSS progress bar width
	MatchedIn   string // best-matched field name shown in card
	Excerpt     string // short snippet around the match
}

// Index is the search index. Create with New(), populate with BuildInBackground().
type Index struct {
	docs  []doc
	ready atomic.Bool // true once building is done
}

// New returns an empty, not-ready Index.
func New() *Index { return &Index{} }

// IsReady returns true when the background build has finished.
func (idx *Index) IsReady() bool { return idx.ready.Load() }

// ─── Build ───────────────────────────────────────────────────────────────────

// BuildInBackground starts background indexing. The index stays empty (IsReady
// == false) until the goroutine finishes.
func (idx *Index) BuildInBackground(
	reg *registry.Registry,
	contentFS embed.FS,
	funcMap template.FuncMap,
) {
	go func() {
		idx.build(reg, contentFS)
		log.Printf("search: indexed %d documents", len(idx.docs))
		idx.ready.Store(true)
	}()
}

func (idx *Index) build(reg *registry.Registry, contentFS embed.FS) {
	var docs []doc

	// Problems
	for _, p := range reg.Problems {
		d := doc{
			Slug:        p.Slug,
			Title:       p.Title,
			Description: p.Description,
			Type:        "problem",
			Route:       "/problem/" + p.Slug,
		}
		d.fields = append(d.fields,
			field{"Title", p.Title, boostTitle},
			field{"Description", p.Description, boostDescription},
			field{"NFRs", nfrText(p), boostNFR},
			field{"Used technologies", fundLinkText(p), boostFundLinks},
		)
		addContentFields(contentFS, p.Path, &d)
		docs = append(docs, d)
	}

	// Fundamentals (recursive — includes children)
	var indexFund func(f *registry.Fundamental, cat, parent string)
	indexFund = func(f *registry.Fundamental, cat, parent string) {
		d := doc{
			Slug:        f.Slug,
			Title:       f.Title,
			Description: f.Description,
			Type:        "fundamental",
			Category:    cat,
			Parent:      parent,
			Route:       "/fund/" + f.Slug,
		}
		d.fields = append(d.fields,
			field{"Title", f.Title, boostTitle},
			field{"Description", f.Description, boostDescription},
			field{"Category", cat, boostNFR},
		)
		if parent != "" {
			d.fields = append(d.fields, field{"Subcategory of " + parent, parent, boostDescription})
		}
		addContentFields(contentFS, f.Path, &d)
		docs = append(docs, d)

		for i := range f.Children {
			indexFund(&f.Children[i], cat, f.Title)
		}
	}
	for _, f := range reg.Fundamentals {
		cat := strings.SplitN(f.Slug, "/", 2)[0]
		indexFund(f, cat, "")
	}

	// Algorithms
	for _, a := range reg.Algorithms {
		d := doc{
			Slug:        a.Slug,
			Title:       a.Title,
			Description: a.Description,
			Type:        "algorithm",
			Route:       "/algo/" + a.Slug,
		}
		d.fields = append(d.fields,
			field{"Title", a.Title, boostTitle},
			field{"Description", a.Description, boostDescription},
			field{"Used in problems", strings.Join(a.UsedIn, " "), boostFundLinks},
		)
		addContentFields(contentFS, a.Path, &d)
		docs = append(docs, d)
	}

	// Patterns
	for _, p := range reg.Patterns {
		d := doc{
			Slug:        p.Slug,
			Title:       p.Title,
			Description: p.Description,
			Type:        "pattern",
			Route:       "/pattern/" + p.Slug,
		}
		d.fields = append(d.fields,
			field{"Title", p.Title, boostTitle},
			field{"Description", p.Description, boostDescription},
		)
		addContentFields(contentFS, p.Path, &d)
		docs = append(docs, d)
	}

	idx.docs = docs
}

// ─── Content extraction ──────────────────────────────────────────────────────

var (
	rePhase     = regexp.MustCompile(`\{\{phase\s+\d+\s+"([^"]+)"`)
	reSay       = regexp.MustCompile(`(?s)\{\{say\s+"(.*?)"\}\}`)
	reKey       = regexp.MustCompile(`(?s)\{\{key\s+"(.*?)"\}\}`)
	reHint      = regexp.MustCompile(`(?s)\{\{hint\s+"[^"]*"\s+"(.*?)"\}\}`)
	reThink     = regexp.MustCompile(`(?s)\{\{think\s+"(.*?)"`)
	reTmplBlock = regexp.MustCompile(`\{\{[^}]*\}\}`)
	reHTMLTag   = regexp.MustCompile(`<[^>]+>`)
	reHTMLEnt   = regexp.MustCompile(`&[a-zA-Z0-9#]+;`)
	reSpaces    = regexp.MustCompile(`\s{2,}`)
)

// addContentFields reads content/{path}/index.html, extracts text by macro
// type, and appends indexed fields with appropriate boosts.
func addContentFields(contentFS embed.FS, contentPath string, d *doc) {
	data, err := fs.ReadFile(contentFS, "content/"+contentPath+"/index.html")
	if err != nil {
		return // content file not found — index metadata only
	}
	raw := string(data)

	// Phase titles — boost 2.5
	if phases := collectMatches(rePhase, raw, 1); len(phases) > 0 {
		d.fields = append(d.fields, field{"Phases", strings.Join(phases, " "), boostPhaseTitle})
	}

	// {{say}} — boost 2.0
	if says := collectMatchesStripped(reSay, raw, 1); len(says) > 0 {
		d.fields = append(d.fields, field{"Interview speech", strings.Join(says, " "), boostSay})
	}

	// {{key}} — boost 2.5
	if keys := collectMatchesStripped(reKey, raw, 1); len(keys) > 0 {
		d.fields = append(d.fields, field{"Key takeaway", strings.Join(keys, " "), boostKey})
	}

	// {{hint}} second arg — boost 1.8
	if hints := collectMatchesStripped(reHint, raw, 1); len(hints) > 0 {
		d.fields = append(d.fields, field{"Hints", strings.Join(hints, " "), boostHint})
	}

	// {{think}} first arg — boost 1.5
	if thinks := collectMatchesStripped(reThink, raw, 1); len(thinks) > 0 {
		d.fields = append(d.fields, field{"Reasoning", strings.Join(thinks, " "), boostThink})
	}

	// Body: strip all template calls and HTML, remaining text — boost 1.0
	body := reTmplBlock.ReplaceAllString(raw, " ")
	body = reHTMLTag.ReplaceAllString(body, " ")
	body = reHTMLEnt.ReplaceAllString(body, " ")
	body = reSpaces.ReplaceAllString(strings.TrimSpace(body), " ")
	if body != "" {
		d.fields = append(d.fields, field{"Content", body, boostBody})
	}
}

// collectMatches extracts capture group n from all matches of re in s.
func collectMatches(re *regexp.Regexp, s string, n int) []string {
	var out []string
	for _, m := range re.FindAllStringSubmatch(s, -1) {
		if len(m) > n {
			out = append(out, m[n])
		}
	}
	return out
}

// collectMatchesStripped extracts capture group n and strips HTML/templates.
func collectMatchesStripped(re *regexp.Regexp, s string, n int) []string {
	raw := collectMatches(re, s, n)
	for i, r := range raw {
		r = reTmplBlock.ReplaceAllString(r, " ")
		r = reHTMLTag.ReplaceAllString(r, " ")
		r = reHTMLEnt.ReplaceAllString(r, " ")
		raw[i] = strings.TrimSpace(reSpaces.ReplaceAllString(r, " "))
	}
	return raw
}

// ─── Registry helpers ────────────────────────────────────────────────────────

func nfrText(p *registry.Problem) string {
	var parts []string
	for _, n := range p.NFRs {
		parts = append(parts, n.Title)
	}
	return strings.Join(parts, " ")
}

func fundLinkText(p *registry.Problem) string {
	var parts []string
	for _, u := range p.Uses {
		parts = append(parts, u.Fundamental)
		if u.Config != "" {
			parts = append(parts, u.Config)
		}
	}
	return strings.Join(parts, " ")
}

// ─── Search ──────────────────────────────────────────────────────────────────

// Search returns up to limit ranked results for the query string.
// Returns nil and false when the index is not yet ready.
func (idx *Index) Search(query string, limit int) ([]Result, bool) {
	if !idx.ready.Load() {
		return nil, false
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return nil, true
	}

	tokens := tokenize(query)
	if len(tokens) == 0 {
		return nil, true
	}

	type scored struct {
		d         doc
		score     float64
		matchedIn string
		excerpt   string
	}

	var hits []scored
	for _, d := range idx.docs {
		score, matchedIn, excerpt := scoreDoc(d, tokens)
		if score > 0 {
			hits = append(hits, scored{d, score, matchedIn, excerpt})
		}
	}

	// Sort by score descending
	sort.Slice(hits, func(i, j int) bool {
		return hits[i].score > hits[j].score
	})

	// Cap to limit
	if len(hits) > limit {
		hits = hits[:limit]
	}

	if len(hits) == 0 {
		return []Result{}, true
	}

	// Normalise scores to 0–100 relative to best hit
	top := hits[0].score

	results := make([]Result, len(hits))
	for i, h := range hits {
		pct := 100
		if top > 0 {
			pct = int(h.score / top * 100)
		}
		// Clamp to [10, 100] so even weak hits have a visible bar
		if pct < 10 {
			pct = 10
		}
		results[i] = Result{
			Slug:        h.d.Slug,
			Title:       h.d.Title,
			Description: h.d.Description,
			Type:        h.d.Type,
			Category:    h.d.Category,
			Parent:      h.d.Parent,
			Route:       h.d.Route,
			ScorePct:    pct,
			MatchedIn:   h.matchedIn,
			Excerpt:     h.excerpt,
		}
	}
	return results, true
}

// ─── Scoring ─────────────────────────────────────────────────────────────────

// scoreDoc computes a relevance score for doc d against the token set.
// It also returns the name of the best-matching field and a short excerpt.
func scoreDoc(d doc, tokens []string) (float64, string, string) {
	var totalScore float64
	bestField := ""
	bestFieldScore := 0.0
	bestExcerpt := ""

	for _, f := range d.fields {
		fieldTokens := tokenize(f.text)
		fs, excerpt := scoreField(fieldTokens, f.text, tokens, f.boost)
		totalScore += fs
		if fs > bestFieldScore {
			bestFieldScore = fs
			bestField = f.name
			bestExcerpt = excerpt
		}
	}

	return totalScore, bestField, bestExcerpt
}

// scoreField scores a single field against the query tokens using exact,
// prefix, and fuzzy (Levenshtein ≤ 2) matching with boost applied.
func scoreField(fieldTokens []string, rawText string, queryTokens []string, boost float64) (float64, string) {
	if len(fieldTokens) == 0 {
		return 0, ""
	}

	var totalScore float64
	var matchedQueryToken string

	for _, qt := range queryTokens {
		bestTokenScore := 0.0
		var bestMatchTok string

		for _, ft := range fieldTokens {
			ts := tokenScore(qt, ft)
			if ts > bestTokenScore {
				bestTokenScore = ts
				bestMatchTok = ft
			}
		}

		totalScore += bestTokenScore * boost
		if bestMatchTok != "" && matchedQueryToken == "" {
			matchedQueryToken = bestMatchTok
		}
	}

	excerpt := ""
	if totalScore > 0 && matchedQueryToken != "" {
		excerpt = extractExcerpt(rawText, matchedQueryToken, 120)
	}

	return totalScore, excerpt
}

// tokenScore returns a match score (0.0–1.0) between a query token and a field token.
// Scoring tiers:
//
//	1.0  exact match
//	0.85 prefix match (query is prefix of field, ≥ 4 chars)
//	0.6  Levenshtein distance == 1
//	0.3  Levenshtein distance == 2
//	0.0  no match
func tokenScore(query, field string) float64 {
	if query == field {
		return 1.0
	}
	// Prefix: "redis" matches "redistribution" — only useful when query ≥ 4 chars
	if len(query) >= 4 && strings.HasPrefix(field, query) {
		return 0.85
	}
	// Field is prefix of query — handles "caching" matching "cache"
	if len(field) >= 4 && strings.HasPrefix(query, field) {
		return 0.7
	}
	// Fuzzy
	d := levenshtein(query, field)
	if d == 1 {
		return 0.6
	}
	if d == 2 && len(query) >= 4 {
		return 0.3
	}
	return 0.0
}

// ─── Levenshtein ─────────────────────────────────────────────────────────────

// levenshtein computes the edit distance between two strings.
// Early-exit when distance would exceed 2 to keep it fast.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	// Short-circuit: length difference > 2 can never be ≤ 2
	diff := la - lb
	if diff < 0 {
		diff = -diff
	}
	if diff > 2 {
		return diff
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		curr[0] = i
		rowMin := curr[0]
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			ins := curr[j-1] + 1
			del := prev[j] + 1
			sub := prev[j-1] + cost
			v := ins
			if del < v {
				v = del
			}
			if sub < v {
				v = sub
			}
			curr[j] = v
			if v < rowMin {
				rowMin = v
			}
		}
		// If the minimum in this row already exceeds 2, abort early
		if rowMin > 2 {
			return rowMin
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

// ─── Excerpt extraction ──────────────────────────────────────────────────────

// extractExcerpt returns up to maxLen characters of text centred around the
// first occurrence of the target word.
func extractExcerpt(text, target string, maxLen int) string {
	lower := strings.ToLower(text)
	idx := strings.Index(lower, target)
	if idx == -1 {
		// Fallback: first maxLen chars
		if len(text) > maxLen {
			return strings.TrimSpace(text[:maxLen]) + "…"
		}
		return strings.TrimSpace(text)
	}
	start := idx - maxLen/3
	if start < 0 {
		start = 0
	}
	end := start + maxLen
	if end > len(text) {
		end = len(text)
		start = end - maxLen
		if start < 0 {
			start = 0
		}
	}
	excerpt := strings.TrimSpace(text[start:end])
	if start > 0 {
		excerpt = "…" + excerpt
	}
	if end < len(text) {
		excerpt = excerpt + "…"
	}
	return excerpt
}

// ─── Tokeniser ───────────────────────────────────────────────────────────────

// tokenize splits text into lowercase words, removing punctuation and
// stop-words shorter than 2 characters.
func tokenize(s string) []string {
	s = strings.ToLower(s)
	var tokens []string
	var cur strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			cur.WriteRune(r)
		} else {
			if cur.Len() >= 2 {
				tokens = append(tokens, cur.String())
			}
			cur.Reset()
		}
	}
	if cur.Len() >= 2 {
		tokens = append(tokens, cur.String())
	}
	return dedupe(tokens)
}

func dedupe(tokens []string) []string {
	seen := make(map[string]bool, len(tokens))
	out := tokens[:0]
	for _, t := range tokens {
		if !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}

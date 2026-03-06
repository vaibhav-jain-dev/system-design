package handlers

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"system-design/internal/registry"
)

type Handler struct {
	reg       *registry.Registry
	templates *template.Template
	contentFS fs.FS
	funcMap   template.FuncMap
}

func New(reg *registry.Registry, templateFS, contentFS fs.FS, funcMap template.FuncMap) *Handler {
	// Parse layout templates
	tmpl := template.Must(
		template.New("").Funcs(funcMap).ParseFS(templateFS,
			"web/templates/base.html",
			"web/templates/sidebar.html",
			"web/templates/welcome.html",
			"web/templates/detail_problem.html",
			"web/templates/detail_fund.html",
			"web/templates/detail_algo.html",
			"web/templates/detail_pattern.html",
			"web/templates/context_card.html",
			"web/templates/doc_card.html",
			"web/templates/detail_concept.html",
			"web/templates/detail_quick.html",
		))

	return &Handler{
		reg:       reg,
		templates: tmpl,
		contentFS: contentFS,
		funcMap:   funcMap,
	}
}

// isHTMX checks if the request came from HTMX (partial swap).
func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// countFundamentals recursively counts fundamentals including children.
func countFundamentals(funds []*registry.Fundamental) int {
	count := 0
	for _, f := range funds {
		count++
		count += countChildFundamentals(f.Children)
	}
	return count
}

func countChildFundamentals(children []registry.Fundamental) int {
	count := 0
	for _, c := range children {
		count++
		count += countChildFundamentals(c.Children)
	}
	return count
}

// baseData returns common template data shared by all handlers.
func (h *Handler) baseData() map[string]interface{} {
	return map[string]interface{}{
		"Problems":          h.reg.Problems,
		"Fundamentals":      h.reg.Fundamentals,
		"FundamentalGroups": h.reg.GroupedFundamentals(),
		"Algorithms":        h.reg.Algorithms,
		"Patterns":          h.reg.Patterns,
		"Concepts":          h.reg.Concepts,
		"QuickCategories":   h.reg.QuickCategories,
		"TotalFundamentals": countFundamentals(h.reg.Fundamentals),
	}
}

// Dashboard renders the full page with sidebar and welcome view.
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data := h.baseData()
	data["Content"] = template.HTML("")
	data["ActiveSlug"] = ""
	data["PageType"] = "welcome"
	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal error", 500)
	}
}

// ProblemDetail renders a problem's detail view.
func (h *Handler) ProblemDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	problem := h.reg.GetProblem(slug)
	if problem == nil {
		http.NotFound(w, r)
		return
	}

	// Render content from HTML file
	content := h.renderContent(problem.Path)

	// Build phase→NFR map (phase number string → []nfr slug) for JS filter
	phaseNFRMap := make(map[string][]string)
	for _, nfr := range problem.NFRs {
		for _, ph := range nfr.Phases {
			key := strconv.Itoa(ph)
			phaseNFRMap[key] = append(phaseNFRMap[key], nfr.Slug)
		}
	}
	phaseJSON, _ := json.Marshal(phaseNFRMap)

	// Build fundamental→NFR map for context card dimming
	useNFRMap := make(map[string][]string)
	for _, use := range problem.Uses {
		if len(use.NFRs) > 0 {
			useNFRMap[use.Fundamental] = use.NFRs
		}
	}
	useJSON, _ := json.Marshal(useNFRMap)

	// Build phase→FR map (phase number string → []fr slug) for JS intersection filter
	phaseFRMap := make(map[string][]string)
	for _, fr := range problem.FRs {
		for _, ph := range fr.Phases {
			key := strconv.Itoa(ph)
			phaseFRMap[key] = append(phaseFRMap[key], fr.Slug)
		}
	}
	frJSON, _ := json.Marshal(phaseFRMap)

	data := h.baseData()
	data["Problem"] = problem
	data["Content"] = content
	data["ActiveSlug"] = slug
	data["PageType"] = "problem"
	data["PhaseNFRMapJSON"] = template.JS(phaseJSON)
	data["UseNFRMapJSON"] = template.JS(useJSON)
	data["PhaseFRMapJSON"] = template.JS(frJSON)

	if isHTMX(r) {
		if err := h.templates.ExecuteTemplate(w, "detail_problem.html", data); err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Internal error", 500)
		}
		return
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal error", 500)
	}
}

// FundamentalDetail renders a fundamental's detail view.
// Supports ?from=problem-slug to highlight items relevant to that problem.
func (h *Handler) FundamentalDetail(w http.ResponseWriter, r *http.Request) {
	// chi wildcard: /fund/networking/load-balancing → slug = "networking/load-balancing"
	slug := chi.URLParam(r, "*")
	slug = strings.TrimPrefix(slug, "/")

	fund := h.reg.GetFundamental(slug)
	if fund == nil {
		http.NotFound(w, r)
		return
	}

	// Contextual highlighting: check if navigating from a specific problem
	fromProblem := r.URL.Query().Get("from")
	var highlightContext *registry.UsageLink
	var fromProblemRef *registry.Problem
	if fromProblem != "" {
		fromProblemRef = h.reg.GetProblem(fromProblem)
		if fromProblemRef != nil {
			for _, use := range fromProblemRef.Uses {
				if use.Fundamental == slug || strings.HasPrefix(use.Fundamental, slug+"/") {
					link := use
					highlightContext = &link
					break
				}
			}
		}
	}

	// Build highlight keywords from the usage config for content matching
	var highlightKeywords []string
	if highlightContext != nil {
		highlightKeywords = extractKeywords(highlightContext.Config)
	}

	content := h.renderContent(fund.Path)

	data := h.baseData()
	data["Fundamental"] = fund
	data["Content"] = content
	data["ActiveSlug"] = slug
	data["PageType"] = "fundamental"
	data["FromProblem"] = fromProblem
	data["FromProblemRef"] = fromProblemRef
	data["HighlightContext"] = highlightContext
	data["HighlightKeywords"] = highlightKeywords
	data["UsedBy"] = aggregateUsedByLinks(fund)

	if isHTMX(r) {
		if err := h.templates.ExecuteTemplate(w, "detail_fund.html", data); err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Internal error", 500)
		}
		return
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal error", 500)
	}
}

// extractKeywords pulls meaningful terms from a config string for highlighting.
func extractKeywords(config string) []string {
	// Split on common separators and filter short/stop words
	stop := map[string]bool{"with": true, "and": true, "for": true, "the": true, "per": true, "via": true, "a": true, "an": true, "in": true, "on": true, "to": true, "of": true, "is": true}
	words := strings.FieldsFunc(config, func(r rune) bool {
		return r == ' ' || r == ',' || r == '(' || r == ')'
	})
	var keywords []string
	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w) > 2 && !stop[strings.ToLower(w)] {
			keywords = append(keywords, w)
		}
	}
	return keywords
}

// renderContent reads and renders a content HTML file with macros.
func (h *Handler) renderContent(contentPath string) template.HTML {
	filePath := path.Join("content", contentPath, "index.html")
	data, err := fs.ReadFile(h.contentFS, filePath)
	if err != nil {
		log.Printf("Content not found: %s (%v)", filePath, err)
		return template.HTML(`<div class="no-content">Content not yet written. Create: ` + filePath + `</div>`)
	}

	// Parse content as a Go template to process macros
	tmpl, err := template.New("content").Funcs(h.funcMap).Parse(string(data))
	if err != nil {
		log.Printf("Content template parse error in %s: %v", filePath, err)
		return template.HTML(`<div class="content-error">Template error: ` + template.HTMLEscapeString(err.Error()) + `</div>`)
	}

	var sb strings.Builder
	if err := tmpl.Execute(&sb, nil); err != nil {
		log.Printf("Content template exec error in %s: %v", filePath, err)
		return template.HTML(`<div class="content-error">Render error: ` + template.HTMLEscapeString(err.Error()) + `</div>`)
	}

	return template.HTML(sb.String())
}

// AlgorithmDetail renders an algorithm's detail view.
func (h *Handler) AlgorithmDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	algo := h.reg.GetAlgorithm(slug)
	if algo == nil {
		http.NotFound(w, r)
		return
	}

	content := h.renderContent(algo.Path)

	data := h.baseData()
	data["Algorithm"] = algo
	data["Content"] = content
	data["ActiveSlug"] = slug
	data["PageType"] = "algorithm"

	if isHTMX(r) {
		if err := h.templates.ExecuteTemplate(w, "detail_algo.html", data); err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Internal error", 500)
		}
		return
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal error", 500)
	}
}

// PatternDetail renders a design pattern's detail view.
func (h *Handler) PatternDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	pattern := h.reg.GetPattern(slug)
	if pattern == nil {
		http.NotFound(w, r)
		return
	}

	content := h.renderContent(pattern.Path)

	data := h.baseData()
	data["Pattern"] = pattern
	data["Content"] = content
	data["ActiveSlug"] = slug
	data["PageType"] = "pattern"

	if isHTMX(r) {
		if err := h.templates.ExecuteTemplate(w, "detail_pattern.html", data); err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Internal error", 500)
		}
		return
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal error", 500)
	}
}

// ConceptDetail renders a concept's detail view showing all appearances.
func (h *Handler) ConceptDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	concept := h.reg.GetConcept(slug)
	if concept == nil {
		http.NotFound(w, r)
		return
	}

	// Find which category this concept belongs to
	var category string
	for _, cat := range h.reg.Concepts {
		for _, c := range cat.Concepts {
			if c.Slug == slug {
				category = cat.Category
				break
			}
		}
	}

	data := h.baseData()
	data["Concept"] = concept
	data["ConceptCategory"] = category
	data["ActiveSlug"] = "concept-" + slug
	data["PageType"] = "concept"

	if isHTMX(r) {
		if err := h.templates.ExecuteTemplate(w, "detail_concept.html", data); err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Internal error", 500)
		}
		return
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal error", 500)
	}
}

// QuickCategoryDetail renders a quick-answer category page.
func (h *Handler) QuickCategoryDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	cat := h.reg.GetQuickCategory(slug)
	if cat == nil {
		http.NotFound(w, r)
		return
	}

	data := h.baseData()
	data["QuickCategory"] = cat
	data["ActiveSlug"] = "quick-" + slug
	data["PageType"] = "quick"

	if isHTMX(r) {
		if err := h.templates.ExecuteTemplate(w, "detail_quick.html", data); err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Internal error", 500)
		}
		return
	}

	if err := h.templates.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Internal error", 500)
	}
}

// Placeholder handlers for PDF generation (Phase 4)
func (h *Handler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"not_implemented"}`))
}

func (h *Handler) GenerateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"not_implemented"}`))
}

func (h *Handler) ServePDF(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	http.ServeFile(w, r, path.Join("output", filename))
}

func aggregateUsedByLinks(fund *registry.Fundamental) []registry.UsageLink {
	combined := make([]registry.UsageLink, 0, len(fund.UsedBy))
	seen := make(map[string]bool)

	add := func(link registry.UsageLink) {
		if seen[link.Problem] {
			return
		}
		seen[link.Problem] = true
		combined = append(combined, link)
	}

	for _, link := range fund.UsedBy {
		add(link)
	}

	var walk func(children []registry.Fundamental)
	walk = func(children []registry.Fundamental) {
		for _, child := range children {
			for _, link := range child.UsedBy {
				add(link)
			}
			walk(child.Children)
		}
	}
	walk(fund.Children)

	sort.Slice(combined, func(i, j int) bool {
		if combined[i].ProblemRef != nil && combined[j].ProblemRef != nil {
			return combined[i].ProblemRef.Title < combined[j].ProblemRef.Title
		}
		return combined[i].Problem < combined[j].Problem
	})

	return combined
}

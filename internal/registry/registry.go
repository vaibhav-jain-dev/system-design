package registry

import (
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// NFRDef defines the display properties of a standard non-functional requirement.
type NFRDef struct {
	Title string
	Color string // CSS color value
}

// StandardNFRs is the canonical set of NFRs with display metadata.
// Problems reference these by slug in their nfrs[] list.
var StandardNFRs = map[string]NFRDef{
	"scalability":   {Title: "Scalability", Color: "#6366F1"},
	"performance":   {Title: "Performance", Color: "#2563EB"},
	"availability":  {Title: "Availability", Color: "#059669"},
	"consistency":   {Title: "Consistency", Color: "#D97706"},
	"durability":    {Title: "Durability", Color: "#7C3AED"},
	"security":      {Title: "Security", Color: "#DC2626"},
	"cost":          {Title: "Cost", Color: "#64748B"},
	"observability": {Title: "Observability", Color: "#0891B2"},
}

// ProblemNFR tags a problem with a standard NFR and the phase numbers where it's addressed.
type ProblemNFR struct {
	Slug   string `yaml:"slug"`
	Phases []int  `yaml:"phases"`

	// Resolved from StandardNFRs (populated after load)
	Title string `yaml:"-"`
	Color string `yaml:"-"`
}

// FunctionalRequirement is a problem-specific feature that can be selected
// to filter content to phases that implement it.
type FunctionalRequirement struct {
	Slug   string `yaml:"slug"`
	Title  string `yaml:"title"`
	Phases []int  `yaml:"phases"`
}

// UsageLink represents a bidirectional link between a Problem and a Fundamental.
type UsageLink struct {
	// Forward: which fundamental is used
	Fundamental string `yaml:"fundamental"`
	// Reverse: which problem uses it (auto-filled)
	Problem string `yaml:"-"`
	// Context fields
	Config  string `yaml:"config"`
	Why     string `yaml:"why"`
	NotThis string `yaml:"not_this"`
	Risk    string `yaml:"risk"`
	Caveats string `yaml:"caveats"`
	// NFR tags: which non-functional requirements this use addresses
	NFRs []string `yaml:"nfrs"`

	// Resolved references (populated after load)
	FundamentalRef *Fundamental `yaml:"-"`
	ProblemRef     *Problem     `yaml:"-"`
}

type DocMeta struct {
	Type   string `yaml:"type"`
	Script string `yaml:"script"`
	Output string `yaml:"output"`
}

type Problem struct {
	Slug        string                  `yaml:"slug"`
	Title       string                  `yaml:"title"`
	Description string                  `yaml:"description"`
	Path        string                  `yaml:"path"`
	Category    string                  `yaml:"category"` // e.g. "distributed" — used to group in sidebar
	Docs        []DocMeta               `yaml:"docs"`
	NFRs        []ProblemNFR            `yaml:"nfrs"`
	FRs         []FunctionalRequirement `yaml:"functional_requirements"`
	Uses        []UsageLink             `yaml:"uses"`

	// Algorithms used in this problem (auto-derived from Algorithm.UsedIn reverse)
	Algorithms []*Algorithm `yaml:"-"`
}

type Fundamental struct {
	Slug        string        `yaml:"slug"`
	Title       string        `yaml:"title"`
	Description string        `yaml:"description"`
	Path        string        `yaml:"path"`
	Children    []Fundamental `yaml:"children"`
	UsedBy      []UsageLink   `yaml:"-"`

	// Optional cross-reference to an algorithm that implements this concept
	RelatedAlgorithm    string     `yaml:"related_algorithm"`
	RelatedAlgorithmRef *Algorithm `yaml:"-"`
}

type Algorithm struct {
	Slug        string   `yaml:"slug"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Path        string   `yaml:"path"`
	UsedIn      []string `yaml:"used_in"`
	// Resolved references (populated after load)
	UsedInProblems []*Problem `yaml:"-"`

	// Optional cross-reference to a fundamental that covers this concept in depth
	RelatedFundamental    string       `yaml:"related_fundamental"`
	RelatedFundamentalRef *Fundamental `yaml:"-"`
}

type Pattern struct {
	Slug        string `yaml:"slug"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
}

// QuickQuestion is a direct interview Q&A with full reasoning, tradeoffs and caveats.
type QuickQuestion struct {
	Number   int    `yaml:"number"`
	Question string `yaml:"question"`
	Needed   string `yaml:"needed"`   // One-liner: what to say first
	Answer   string `yaml:"answer"`   // Full explanation with reasoning + justification
	Tradeoff string `yaml:"tradeoff"` // Alternatives and why they are worse/different
	Caveat   string `yaml:"caveat"`   // Common mistakes, gotchas, what NOT to say
	Watch    string `yaml:"watch"`    // The thing that makes the answer stand out

	// Optional inline SVG/HTML diagram showing the flow described in the answer
	DiagramHTML string `yaml:"diagram_html"`

	// Optional cross-references (slugs)
	RelatedProblems     []string `yaml:"related_problems"`
	RelatedFundamentals []string `yaml:"related_fundamentals"`

	// Resolved references (populated after load)
	RelatedProblemRefs     []*Problem     `yaml:"-"`
	RelatedFundamentalRefs []*Fundamental `yaml:"-"`
}

// QuickCategory groups related quick questions under a topic.
type QuickCategory struct {
	Slug        string          `yaml:"slug"`
	Title       string          `yaml:"title"`
	Description string          `yaml:"description"`
	NumberRange string          `yaml:"number_range"`
	Questions   []QuickQuestion `yaml:"questions"`
}

// ConceptAppearance records where a concept appears (section-level granularity).
type ConceptAppearance struct {
	Type    string `yaml:"type"`    // "problem", "fundamental", "algorithm", "pattern"
	Slug    string `yaml:"slug"`    // e.g. "url-shortener", "storage/redis"
	Section string `yaml:"section"` // e.g. "Caching Deep Dive"
	Phase   int    `yaml:"phase"`   // phase number (0 if not applicable)

	// Resolved reference (populated after load)
	Title string `yaml:"-"` // resolved title of the target
	URL   string `yaml:"-"` // resolved URL path
}

// Concept is a cross-cutting topic that appears across multiple content types.
type Concept struct {
	Slug        string              `yaml:"slug"`
	Title       string              `yaml:"title"`
	Description string              `yaml:"description"`
	AppearsIn   []ConceptAppearance `yaml:"appears_in"`
}

// ConceptCategory groups related concepts under a named category.
type ConceptCategory struct {
	Category string    `yaml:"category"`
	Concepts []Concept `yaml:"concepts"`
}

type registryFile struct {
	Problems        []Problem        `yaml:"problems"`
	Fundamentals    []Fundamental    `yaml:"fundamentals"`
	Algorithms      []Algorithm      `yaml:"algorithms"`
	Patterns        []Pattern        `yaml:"patterns"`
	Concepts        []ConceptCategory `yaml:"concepts"`
	QuickCategories []QuickCategory  `yaml:"quick_categories"`
}

// Registry holds the loaded knowledge graph.
type Registry struct {
	Problems        []*Problem
	Fundamentals    []*Fundamental
	Algorithms      []*Algorithm
	Patterns        []*Pattern
	Concepts        []*ConceptCategory
	QuickCategories []*QuickCategory

	problemsBySlug        map[string]*Problem
	fundamentalsBySlug    map[string]*Fundamental
	fundamentalAncestors  map[string][]*Fundamental
	algorithmsBySlug      map[string]*Algorithm
	patternsBySlug        map[string]*Pattern
	conceptsBySlug        map[string]*Concept
	quickCategoriesBySlug map[string]*QuickCategory
}

// Load parses _registry.yaml and builds the knowledge graph with reverse links.
func Load(fsys fs.FS, path string) (*Registry, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read registry: %w", err)
	}

	var raw registryFile
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse registry: %w", err)
	}

	reg := &Registry{
		problemsBySlug:        make(map[string]*Problem),
		fundamentalsBySlug:    make(map[string]*Fundamental),
		fundamentalAncestors:  make(map[string][]*Fundamental),
		algorithmsBySlug:      make(map[string]*Algorithm),
		patternsBySlug:        make(map[string]*Pattern),
		conceptsBySlug:        make(map[string]*Concept),
		quickCategoriesBySlug: make(map[string]*QuickCategory),
	}

	// Index fundamentals (including children)
	for i := range raw.Fundamentals {
		f := &raw.Fundamentals[i]
		reg.Fundamentals = append(reg.Fundamentals, f)
		reg.indexFundamental(f, nil)
	}

	// Index problems and build forward/reverse links
	for i := range raw.Problems {
		p := &raw.Problems[i]
		reg.Problems = append(reg.Problems, p)
		reg.problemsBySlug[p.Slug] = p

		// Resolve NFR display metadata from StandardNFRs
		for j := range p.NFRs {
			if def, ok := StandardNFRs[p.NFRs[j].Slug]; ok {
				p.NFRs[j].Title = def.Title
				p.NFRs[j].Color = def.Color
			} else {
				log.Printf("WARNING: problem %q references unknown NFR slug %q", p.Slug, p.NFRs[j].Slug)
			}
		}

		for j := range p.Uses {
			link := &p.Uses[j]
			link.Problem = p.Slug
			link.ProblemRef = p
			if f, ok := reg.fundamentalsBySlug[link.Fundamental]; ok {
				link.FundamentalRef = f
				reverseLink := *link
				reverseLink.ProblemRef = p
				appendUsedByUnique(f, reverseLink)

				// If a problem links to a child fundamental (e.g. networking/cdn/cloudfront),
				// also surface it on parent fundamentals (e.g. networking/cdn) so the
				// reverse links stay complete at every level.
				for _, ancestor := range reg.fundamentalAncestors[link.Fundamental] {
					appendUsedByUnique(ancestor, reverseLink)
				}
			} else {
				log.Printf("WARNING: problem %q references non-existent fundamental %q", p.Slug, link.Fundamental)
			}
		}
	}

	// Index algorithms and resolve problem references
	for i := range raw.Algorithms {
		a := &raw.Algorithms[i]
		reg.Algorithms = append(reg.Algorithms, a)
		reg.algorithmsBySlug[a.Slug] = a
		for _, problemSlug := range a.UsedIn {
			if p, ok := reg.problemsBySlug[problemSlug]; ok {
				a.UsedInProblems = append(a.UsedInProblems, p)
				// Reverse link: problem knows which algorithms it uses
				p.Algorithms = append(p.Algorithms, a)
			} else {
				log.Printf("WARNING: algorithm %q references non-existent problem %q", a.Slug, problemSlug)
			}
		}
	}

	// Resolve algorithm ↔ fundamental cross-references
	for _, a := range reg.Algorithms {
		if a.RelatedFundamental != "" {
			a.RelatedFundamentalRef = reg.fundamentalsBySlug[a.RelatedFundamental]
			if a.RelatedFundamentalRef == nil {
				log.Printf("WARNING: algorithm %q references non-existent fundamental %q", a.Slug, a.RelatedFundamental)
			}
		}
	}
	for _, f := range reg.fundamentalsBySlug {
		if f.RelatedAlgorithm != "" {
			f.RelatedAlgorithmRef = reg.algorithmsBySlug[f.RelatedAlgorithm]
			if f.RelatedAlgorithmRef == nil {
				log.Printf("WARNING: fundamental %q references non-existent algorithm %q", f.Slug, f.RelatedAlgorithm)
			}
		}
	}

	// Index patterns
	for i := range raw.Patterns {
		pt := &raw.Patterns[i]
		reg.Patterns = append(reg.Patterns, pt)
		reg.patternsBySlug[pt.Slug] = pt
	}

	// Index concepts and resolve references
	for i := range raw.Concepts {
		cat := &raw.Concepts[i]
		reg.Concepts = append(reg.Concepts, cat)
		for j := range cat.Concepts {
			c := &cat.Concepts[j]
			reg.conceptsBySlug[c.Slug] = c
			for k := range c.AppearsIn {
				a := &c.AppearsIn[k]
				switch a.Type {
				case "problem":
					if p := reg.problemsBySlug[a.Slug]; p != nil {
						a.Title = p.Title
						a.URL = "/problem/" + a.Slug
					}
				case "fundamental":
					if f := reg.fundamentalsBySlug[a.Slug]; f != nil {
						a.Title = f.Title
						a.URL = "/fund/" + a.Slug
					}
				case "algorithm":
					if al := reg.algorithmsBySlug[a.Slug]; al != nil {
						a.Title = al.Title
						a.URL = "/algo/" + a.Slug
					}
				case "pattern":
					if pt := reg.patternsBySlug[a.Slug]; pt != nil {
						a.Title = pt.Title
						a.URL = "/pattern/" + a.Slug
					}
				}
			}
		}
	}

	// Index quick categories and resolve cross-references
	for i := range raw.QuickCategories {
		cat := &raw.QuickCategories[i]
		reg.QuickCategories = append(reg.QuickCategories, cat)
		reg.quickCategoriesBySlug[cat.Slug] = cat
		for j := range cat.Questions {
			q := &cat.Questions[j]
			for _, pSlug := range q.RelatedProblems {
				if p := reg.problemsBySlug[pSlug]; p != nil {
					q.RelatedProblemRefs = append(q.RelatedProblemRefs, p)
				}
			}
			for _, fSlug := range q.RelatedFundamentals {
				if f := reg.fundamentalsBySlug[fSlug]; f != nil {
					q.RelatedFundamentalRefs = append(q.RelatedFundamentalRefs, f)
				}
			}
		}
	}

	return reg, nil
}

func (r *Registry) indexFundamental(f *Fundamental, ancestors []*Fundamental) {
	r.fundamentalsBySlug[f.Slug] = f
	if len(ancestors) > 0 {
		r.fundamentalAncestors[f.Slug] = append([]*Fundamental(nil), ancestors...)
	}

	nextAncestors := append(append([]*Fundamental(nil), ancestors...), f)
	for i := range f.Children {
		child := &f.Children[i]
		r.indexFundamental(child, nextAncestors)
	}
}

func appendUsedByUnique(f *Fundamental, link UsageLink) {
	for _, existing := range f.UsedBy {
		if existing.Problem == link.Problem {
			return
		}
	}
	f.UsedBy = append(f.UsedBy, link)
}

// GetProblem returns a problem by slug.
func (r *Registry) GetProblem(slug string) *Problem {
	return r.problemsBySlug[slug]
}

// GetFundamental returns a fundamental by slug path (e.g., "networking/load-balancing").
func (r *Registry) GetFundamental(slug string) *Fundamental {
	return r.fundamentalsBySlug[slug]
}

// AllFundamentals returns all fundamentals (flat, including children).
func (r *Registry) AllFundamentals() []*Fundamental {
	var all []*Fundamental
	for _, f := range r.fundamentalsBySlug {
		all = append(all, f)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Slug < all[j].Slug
	})
	return all
}

// GetAlgorithm returns an algorithm by slug.
func (r *Registry) GetAlgorithm(slug string) *Algorithm {
	return r.algorithmsBySlug[slug]
}

// GetPattern returns a pattern by slug.
func (r *Registry) GetPattern(slug string) *Pattern {
	return r.patternsBySlug[slug]
}

// GetConcept returns a concept by slug.
func (r *Registry) GetConcept(slug string) *Concept {
	return r.conceptsBySlug[slug]
}

// GetQuickCategory returns a quick category by slug.
func (r *Registry) GetQuickCategory(slug string) *QuickCategory {
	return r.quickCategoriesBySlug[slug]
}

// FundamentalGroup groups top-level fundamentals under a category name
// derived from the first path segment of their slug (e.g. "networking", "storage").
type FundamentalGroup struct {
	Category string
	Items    []*Fundamental
}

// GroupedFundamentals returns fundamentals organised by their top-level category.
// The order follows the order in which categories first appear in the registry.
func (r *Registry) GroupedFundamentals() []FundamentalGroup {
	seen := make(map[string]int) // category → index in groups slice
	var groups []FundamentalGroup

	for _, f := range r.Fundamentals {
		cat := categoryOf(f.Slug)
		idx, ok := seen[cat]
		if !ok {
			idx = len(groups)
			seen[cat] = idx
			groups = append(groups, FundamentalGroup{Category: cat})
		}
		groups[idx].Items = append(groups[idx].Items, f)
	}
	return groups
}

// categoryOf extracts the human-readable category from a fundamental slug.
// E.g. "networking/load-balancing" → "Networking", "storage/redis" → "Storage".
func categoryOf(slug string) string {
	var prefix string
	if idx := strings.Index(slug, "/"); idx >= 0 {
		prefix = slug[:idx]
	} else {
		prefix = slug
	}
	if len(prefix) == 0 {
		return "Other"
	}
	return strings.ToUpper(prefix[:1]) + prefix[1:]
}

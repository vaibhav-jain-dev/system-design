package registry

import (
	"fmt"
	"io/fs"
	"log"
	"sort"

	"gopkg.in/yaml.v3"
)

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
	Slug        string      `yaml:"slug"`
	Title       string      `yaml:"title"`
	Description string      `yaml:"description"`
	Path        string      `yaml:"path"`
	Docs        []DocMeta   `yaml:"docs"`
	Uses        []UsageLink `yaml:"uses"`
}

type Fundamental struct {
	Slug        string        `yaml:"slug"`
	Title       string        `yaml:"title"`
	Description string        `yaml:"description"`
	Path        string        `yaml:"path"`
	Children    []Fundamental `yaml:"children"`
	UsedBy      []UsageLink   `yaml:"-"`
}

type Algorithm struct {
	Slug        string   `yaml:"slug"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Path        string   `yaml:"path"`
	UsedIn      []string `yaml:"used_in"`
	// Resolved references (populated after load)
	UsedInProblems []*Problem `yaml:"-"`
}

type Pattern struct {
	Slug        string `yaml:"slug"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
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
	Problems     []Problem         `yaml:"problems"`
	Fundamentals []Fundamental     `yaml:"fundamentals"`
	Algorithms   []Algorithm       `yaml:"algorithms"`
	Patterns     []Pattern         `yaml:"patterns"`
	Concepts     []ConceptCategory `yaml:"concepts"`
}

// Registry holds the loaded knowledge graph.
type Registry struct {
	Problems     []*Problem
	Fundamentals []*Fundamental
	Algorithms   []*Algorithm
	Patterns     []*Pattern
	Concepts     []*ConceptCategory

	problemsBySlug     map[string]*Problem
	fundamentalsBySlug map[string]*Fundamental
	algorithmsBySlug   map[string]*Algorithm
	patternsBySlug     map[string]*Pattern
	conceptsBySlug     map[string]*Concept
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
		problemsBySlug:     make(map[string]*Problem),
		fundamentalsBySlug: make(map[string]*Fundamental),
		algorithmsBySlug:   make(map[string]*Algorithm),
		patternsBySlug:     make(map[string]*Pattern),
		conceptsBySlug:     make(map[string]*Concept),
	}

	// Index fundamentals (including children)
	for i := range raw.Fundamentals {
		f := &raw.Fundamentals[i]
		reg.Fundamentals = append(reg.Fundamentals, f)
		reg.indexFundamental(f)
	}

	// Index problems and build forward/reverse links
	for i := range raw.Problems {
		p := &raw.Problems[i]
		reg.Problems = append(reg.Problems, p)
		reg.problemsBySlug[p.Slug] = p

		for j := range p.Uses {
			link := &p.Uses[j]
			link.Problem = p.Slug
			link.ProblemRef = p
			if f, ok := reg.fundamentalsBySlug[link.Fundamental]; ok {
				link.FundamentalRef = f
				reverseLink := *link
				reverseLink.ProblemRef = p
				f.UsedBy = append(f.UsedBy, reverseLink)
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
			} else {
				log.Printf("WARNING: algorithm %q references non-existent problem %q", a.Slug, problemSlug)
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

	return reg, nil
}

func (r *Registry) indexFundamental(f *Fundamental) {
	r.fundamentalsBySlug[f.Slug] = f
	for i := range f.Children {
		child := &f.Children[i]
		r.indexFundamental(child)
	}
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

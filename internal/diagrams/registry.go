// Package diagrams provides a centralized registry of all visual diagrams
// used across the system design content. Each diagram is identified by a
// unique slug and includes metadata for discovery and maintenance.
//
// Usage in templates:
//
//	{{diagram "rate-limiter-architecture"}}          — lookup by slug
//	{{diagram "Title" `<div>inline HTML</div>`}}     — inline mode (backward compat)
package diagrams

// Type represents the rendering mode of a diagram.
type Type string

const (
	TypeHTML  Type = "html"  // Rendered from HTML/CSS classes
	TypeImage Type = "image" // Rendered from an image file (PNG/SVG)
)

// Diagram holds the content and metadata for a single visual diagram.
type Diagram struct {
	Slug        string // Unique identifier, e.g. "rate-limiter-architecture"
	Title       string // Display title shown above diagram
	Description string // Short description of what the diagram shows
	ContentFile string // Which content file uses this diagram, e.g. "problems/rate-limiter"
	Type        Type   // html or image
	HTML        string // Raw HTML content (for TypeHTML)
	ImagePath   string // Path relative to /static/img/diagrams/ (for TypeImage)
}

// Registry is a map of slug → Diagram for fast lookup.
type Registry struct {
	diagrams map[string]*Diagram
	all      []*Diagram
}

// New creates an empty diagram registry.
func New() *Registry {
	return &Registry{
		diagrams: make(map[string]*Diagram),
	}
}

// Register adds a diagram to the registry. Panics on duplicate slug.
func (r *Registry) Register(d *Diagram) {
	if _, exists := r.diagrams[d.Slug]; exists {
		panic("duplicate diagram slug: " + d.Slug)
	}
	r.diagrams[d.Slug] = d
	r.all = append(r.all, d)
}

// Get returns a diagram by slug, or nil if not found.
func (r *Registry) Get(slug string) *Diagram {
	return r.diagrams[slug]
}

// All returns all registered diagrams.
func (r *Registry) All() []*Diagram {
	return r.all
}

// ByContentFile returns all diagrams for a given content file path.
func (r *Registry) ByContentFile(contentFile string) []*Diagram {
	var result []*Diagram
	for _, d := range r.all {
		if d.ContentFile == contentFile {
			result = append(result, d)
		}
	}
	return result
}

// Count returns the total number of registered diagrams.
func (r *Registry) Count() int {
	return len(r.all)
}

// BuildDefault creates and populates the registry with all diagrams.
func BuildDefault() *Registry {
	r := New()
	registerRateLimiter(r)
	registerInstagram(r)
	registerURLShortener(r)
	registerAlgorithms(r)
	registerFundamentals(r)
	registerPatterns(r)
	return r
}

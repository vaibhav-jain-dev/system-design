package handlers

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"

	"system-design/internal/registry"
)

type Handler struct {
	reg        *registry.Registry
	templates  *template.Template
	contentFS  fs.FS
	funcMap    template.FuncMap
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
			"web/templates/context_card.html",
			"web/templates/doc_card.html",
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

// Dashboard renders the full page with sidebar and welcome view.
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Problems":     h.reg.Problems,
		"Fundamentals": h.reg.Fundamentals,
		"Content":      template.HTML(""),
		"ActiveSlug":   "",
		"PageType":     "welcome",
	}
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

	data := map[string]interface{}{
		"Problem":      problem,
		"Content":      content,
		"Problems":     h.reg.Problems,
		"Fundamentals": h.reg.Fundamentals,
		"ActiveSlug":   slug,
		"PageType":     "problem",
	}

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
				if use.Fundamental == slug {
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

	data := map[string]interface{}{
		"Fundamental":       fund,
		"Content":           content,
		"Problems":          h.reg.Problems,
		"Fundamentals":      h.reg.Fundamentals,
		"ActiveSlug":        slug,
		"PageType":          "fundamental",
		"FromProblem":       fromProblem,
		"FromProblemRef":    fromProblemRef,
		"HighlightContext":  highlightContext,
		"HighlightKeywords": highlightKeywords,
	}

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
		return r == ' ' || r == ',' || r == '(' || r == ')' || r == '-'
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
		return template.HTML(`<div class="content-error">Template error: ` + err.Error() + `</div>`)
	}

	var sb strings.Builder
	if err := tmpl.Execute(&sb, nil); err != nil {
		log.Printf("Content template exec error in %s: %v", filePath, err)
		return template.HTML(`<div class="content-error">Render error: ` + err.Error() + `</div>`)
	}

	return template.HTML(sb.String())
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

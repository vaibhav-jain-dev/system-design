package macros

import (
	"fmt"
	"html/template"
	"strings"
)

// FuncMap returns the template function map for content macros.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"say":       say,
		"thought":   thought,
		"avoid":     avoid,
		"key":       key,
		"phase":     phase,
		"code":      code,
		"qa":        qa,
		"followup":  followup,
		"checklist": checklist,
		"compare":   compare,
		"table":     tableMacro,
		"info":      info,
		"diagram":   diagram,

		// Helpers for building structured data in templates
		"options": optionsList,
		"best":    optBest,
		"alt":     optAlt,
		"nofit":   optNofit,
		"rows":    rowsList,
		"row":     rowItem,

		// Utility functions for templates
		"map":      makeMap,
		"multiply": multiply,
		"slugIcon": slugIcon,
		"contains": strings.Contains,
	}
}

// makeMap creates a map from alternating key-value pairs.
func makeMap(pairs ...interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i+1 < len(pairs); i += 2 {
		if key, ok := pairs[i].(string); ok {
			m[key] = pairs[i+1]
		}
	}
	return m
}

// multiply returns a * b (for CSS calc in templates).
func multiply(a, b int) int {
	return a * b
}

// slugIcon returns an inline SVG icon for a given slug.
func slugIcon(slug string) template.HTML {
	icons := map[string]string{
		"icon-url":           `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>`,
		"icon-rate-limit":    `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>`,
		"icon-chat":          `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/></svg>`,
		"icon-load-balancer": `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="5" r="3"/><circle cx="5" cy="19" r="3"/><circle cx="19" cy="19" r="3"/><line x1="12" y1="8" x2="5" y2="16"/><line x1="12" y1="8" x2="19" y2="16"/></svg>`,
		"icon-cdn":           `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>`,
		"icon-redis":         `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>`,
		"icon-database":      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>`,
		"icon-networking":    `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="6" height="6" rx="1"/><rect x="16" y="2" width="6" height="6" rx="1"/><rect x="9" y="16" width="6" height="6" rx="1"/><path d="M5 8v3a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8"/><line x1="12" y1="13" x2="12" y2="16"/></svg>`,
		"icon-storage":       `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>`,
		"icon-compute":       `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="4" width="16" height="16" rx="2"/><path d="M9 9h6v6H9z"/><path d="M9 1v3M15 1v3M9 20v3M15 20v3M20 9h3M20 14h3M1 9h3M1 14h3"/></svg>`,
		"icon-messaging":     `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>`,
		"icon-fundamental":   `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 2 7 12 12 22 7 12 2"/><polyline points="2 17 12 22 22 17"/><polyline points="2 12 12 17 22 12"/></svg>`,
		"icon-instagram":     `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="20" rx="5"/><circle cx="12" cy="12" r="5"/><circle cx="17.5" cy="6.5" r="1.5"/></svg>`,
		"icon-problem":       `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><circle cx="12" cy="12" r="10"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`,
	}

	slugToIcon := map[string]string{
		"url-shortener":              "icon-url",
		"rate-limiter":               "icon-rate-limit",
		"chat-system":                "icon-chat",
		"networking/load-balancing":   "icon-load-balancer",
		"networking/load-balancing/alb": "icon-load-balancer",
		"networking/load-balancing/nlb": "icon-load-balancer",
		"networking/cdn":              "icon-cdn",
		"networking/cdn/cloudfront":   "icon-cdn",
		"storage/redis":              "icon-redis",
		"storage/dynamodb":           "icon-database",
		"storage/postgres":           "icon-database",
		"instagram":                  "icon-instagram",
		"problem":                    "icon-problem",
		"fundamental":                "icon-fundamental",
	}

	iconName := ""
	if name, ok := slugToIcon[slug]; ok {
		iconName = name
	} else if strings.HasPrefix(slug, "networking") {
		iconName = "icon-networking"
	} else if strings.HasPrefix(slug, "storage") {
		iconName = "icon-storage"
	} else if strings.HasPrefix(slug, "compute") {
		iconName = "icon-compute"
	} else if strings.HasPrefix(slug, "messaging") {
		iconName = "icon-messaging"
	} else {
		iconName = "icon-fundamental"
	}

	if svg, ok := icons[iconName]; ok {
		return template.HTML(svg)
	}
	return template.HTML(icons["icon-fundamental"])
}

// say renders an interview say-box.
func say(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="say-box"><span class="say-label">Say:</span> <em>"%s"</em></div>`, text))
}

// thought renders a thought cloud for reasoning/math/failure scenarios.
func thought(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="thought-cloud">%s</div>`, text))
}

// avoid renders an avoid box for common mistakes.
func avoid(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="avoid-box"><span class="avoid-label">Avoid:</span> %s</div>`, text))
}

// key renders a key takeaway box.
func key(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="key-takeaway"><strong>%s</strong></div>`, text))
}

// phase renders a phase/section header.
func phase(num int, title, time string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="phase-header">
			<span class="phase-number">%d</span>
			<span class="phase-title">%s</span>
			<span class="phase-time">%s</span>
		</div>`, num, title, time))
}

// code renders a syntax-highlighted code block.
func code(lang, content string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="code-block" data-lang="%s"><pre><code>%s</code></pre></div>`,
		lang, template.HTMLEscapeString(content)))
}

// qa renders an interviewer Q&A card.
func qa(args ...string) template.HTML {
	if len(args) < 2 {
		return ""
	}
	q, a := args[0], args[1]
	html := fmt.Sprintf(
		`<div class="qa-card">
			<p class="qa-interviewer"><strong>Interviewer:</strong> "%s"</p>
			<p class="qa-you"><strong>You:</strong> "%s"</p>`, q, a)

	if len(args) >= 4 {
		fq, fa := args[2], args[3]
		html += fmt.Sprintf(
			`<div class="qa-followup">
				<p class="qa-interviewer"><strong>Follow-up:</strong> "%s"</p>
				<p class="qa-you"><strong>You:</strong> "%s"</p>
			</div>`, fq, fa)
	}
	html += `</div>`
	return template.HTML(html)
}

// followup renders a follow-up question card.
func followup(question, answer string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="followup-box">
			<p class="followup-q"><strong>If they ask:</strong> "%s"</p>
			<p class="followup-a"><strong>You:</strong> "%s"</p>
		</div>`, question, answer))
}

// checklist renders a green checklist.
func checklist(items ...string) template.HTML {
	var sb strings.Builder
	sb.WriteString(`<div class="checklist"><ul>`)
	for _, item := range items {
		sb.WriteString(fmt.Sprintf(`<li>%s</li>`, item))
	}
	sb.WriteString(`</ul></div>`)
	return template.HTML(sb.String())
}

// info renders an info tooltip for term definitions.
func info(term, definition string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<span class="info-term">%s<span class="info-icon" title="%s">ℹ</span><span class="info-tooltip">%s</span></span>`,
		term, definition, definition))
}

// Option types for compare macro
type CompareOption struct {
	Kind    string // "best", "alt", "nofit"
	Name    string
	Reason  string
}

func optBest(name, reason string) CompareOption  { return CompareOption{"best", name, reason} }
func optAlt(name, reason string) CompareOption   { return CompareOption{"alt", name, reason} }
func optNofit(name, reason string) CompareOption { return CompareOption{"nofit", name, reason} }

func optionsList(opts ...CompareOption) []CompareOption { return opts }

// compare renders a color-coded comparison card.
func compare(title string, opts []CompareOption) template.HTML {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<div class="compare-card"><div class="compare-title">%s</div>`, title))
	for _, opt := range opts {
		icon := map[string]string{"best": "✓", "alt": "~", "nofit": "✗"}[opt.Kind]
		sb.WriteString(fmt.Sprintf(
			`<div class="compare-option compare-%s"><span class="compare-icon">%s</span><strong>%s</strong> — %s</div>`,
			opt.Kind, icon, opt.Name, opt.Reason))
	}
	sb.WriteString(`</div>`)
	return template.HTML(sb.String())
}

// Row type for table macro
type TableRow struct {
	Cells []string
}

func rowItem(cells ...string) TableRow  { return TableRow{cells} }
func rowsList(rows ...TableRow) []TableRow { return rows }

// tableMacro renders a styled table.
func tableMacro(title string, rows []TableRow) template.HTML {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<div class="styled-table"><div class="table-title">%s</div><table>`, title))
	if len(rows) > 0 {
		sb.WriteString(`<thead><tr>`)
		for _, cell := range rows[0].Cells {
			sb.WriteString(fmt.Sprintf(`<th>%s</th>`, cell))
		}
		sb.WriteString(`</tr></thead><tbody>`)
		for _, row := range rows[1:] {
			sb.WriteString(`<tr>`)
			for _, cell := range row.Cells {
				sb.WriteString(fmt.Sprintf(`<td>%s</td>`, cell))
			}
			sb.WriteString(`</tr>`)
		}
		sb.WriteString(`</tbody>`)
	}
	sb.WriteString(`</table></div>`)
	return template.HTML(sb.String())
}

// diagram renders an SVG-based diagram from node/edge definitions.
// For now returns a placeholder — will be expanded with full SVG rendering.
func diagram(title string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="diagram-container">
			<div class="diagram-title">%s</div>
			<div class="diagram-placeholder">Diagram: %s</div>
		</div>`, title, title))
}

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

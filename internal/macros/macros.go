package macros

import (
	"fmt"
	"html/template"
	"strings"

	"system-design/internal/diagrams"
)

// FuncMap returns the template function map for content macros.
// The diagram registry enables slug-based diagram lookup.
func FuncMap(diagramReg *diagrams.Registry) template.FuncMap {
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
		"diagram":   makeDiagramFunc(diagramReg),
		"details":   details,
		"stageNav":  stageNav,
		"anchor":    anchor,
		"deepQA":    deepQA,

		// Thought process macros
		"hint":      hint,
		"think":     think,
		"triggerQs": triggerQs,

		// Interview insight macros
		"mustKnow":       mustKnow,
		"goodToKnow":     goodToKnow,
		"caveat":         caveat,
		"collapseSection": collapseSection,

		// Helpers for building structured data in templates
		"options": optionsList,
		"best":    optBest,
		"alt":     optAlt,
		"nofit":   optNofit,
		"rows":    rowsList,
		"row":     rowItem,
		"whyNot":  whyNotItem,
		"whatIf":  whatIfItem,
		"how":     howItem,

		// Utility functions for templates
		"map":          makeMap,
		"multiply":     multiply,
		"slugIcon":     slugIcon,
		"contains":     strings.Contains,
		"slugImagePath": slugImagePath,
		"safe":         func(s string) template.HTML { return template.HTML(s) },
	}
}

// slugImagePath maps a content slug to a tech image path under /static/img/tech/.
// Fundamental slugs contain "/" (e.g. "storage/redis") which would break URL paths,
// so we use an explicit whitelist of known images.
func slugImagePath(kind, slug string) string {
	knownImages := map[string]string{
		"problem:url-shortener":               "/static/img/tech/problem-url-shortener.svg",
		"problem:rate-limiter":                "/static/img/tech/problem-rate-limiter.svg",
		"problem:instagram":                   "/static/img/tech/problem-instagram.svg",
		"fundamental:storage/redis":           "/static/img/tech/fund-redis.svg",
		"fundamental:storage/dynamodb":        "/static/img/tech/fund-dynamodb.svg",
		"fundamental:networking/load-balancing": "/static/img/tech/fund-load-balancer.svg",
		"fundamental:networking/cdn":          "/static/img/tech/fund-cdn.svg",
		"algorithm:base62-encoding":           "/static/img/tech/algo-base62-encoding.svg",
		"algorithm:consistent-hashing":        "/static/img/tech/algo-consistent-hashing.svg",
		"algorithm:token-bucket":              "/static/img/tech/algo-token-bucket.svg",
		"algorithm:bloom-filter":              "/static/img/tech/algo-bloom-filter.svg",
		"algorithm:snowflake-id":              "/static/img/tech/algo-snowflake-id.svg",
		"algorithm:trie":                      "/static/img/tech/algo-trie.svg",
		"algorithm:geohash":                   "/static/img/tech/algo-geohash.svg",
		"pattern:rag":                         "/static/img/tech/pat-rag.svg",
		"pattern:agent-tools":                 "/static/img/tech/pat-agent-tools.svg",
		"pattern:prompt-chaining":             "/static/img/tech/pat-prompt-chaining.svg",
		"pattern:guardrails":                  "/static/img/tech/pat-guardrails.svg",
		"pattern:embeddings-vector-search":    "/static/img/tech/pat-embeddings-vector-search.svg",
	}
	if path, ok := knownImages[kind+":"+slug]; ok {
		return path
	}
	return ""
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
		"icon-url":             `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>`,
		"icon-rate-limit":      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>`,
		"icon-chat":            `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/></svg>`,
		"icon-load-balancer":   `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="5" r="3"/><circle cx="5" cy="19" r="3"/><circle cx="19" cy="19" r="3"/><line x1="12" y1="8" x2="5" y2="16"/><line x1="12" y1="8" x2="19" y2="16"/></svg>`,
		"icon-cdn":             `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>`,
		"icon-redis":           `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>`,
		"icon-database":        `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>`,
		"icon-networking":      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="6" height="6" rx="1"/><rect x="16" y="2" width="6" height="6" rx="1"/><rect x="9" y="16" width="6" height="6" rx="1"/><path d="M5 8v3a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8"/><line x1="12" y1="13" x2="12" y2="16"/></svg>`,
		"icon-storage":         `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>`,
		"icon-compute":         `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="4" width="16" height="16" rx="2"/><path d="M9 9h6v6H9z"/><path d="M9 1v3M15 1v3M9 20v3M15 20v3M20 9h3M20 14h3M1 9h3M1 14h3"/></svg>`,
		"icon-messaging":       `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>`,
		"icon-fundamental":     `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 2 7 12 12 22 7 12 2"/><polyline points="2 17 12 22 22 17"/><polyline points="2 12 12 17 22 12"/></svg>`,
		"icon-instagram":       `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="20" rx="5"/><circle cx="12" cy="12" r="5"/><circle cx="17.5" cy="6.5" r="1.5"/></svg>`,
		"icon-problem":         `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><circle cx="12" cy="12" r="10"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`,
		"icon-algorithm":       `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/><path d="M8 7h8M8 11h6M8 15h4"/></svg>`,
		"icon-base62":          `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/><line x1="14" y1="4" x2="10" y2="20"/></svg>`,
		"icon-token-bucket":    `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><circle cx="12" cy="14" r="3"/></svg>`,
		"icon-consistent-hash": `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><circle cx="12" cy="5" r="1.5" fill="currentColor"/><circle cx="18" cy="10" r="1.5" fill="currentColor"/><circle cx="16" cy="17" r="1.5" fill="currentColor"/><circle cx="7" cy="16" r="1.5" fill="currentColor"/><circle cx="5" cy="9" r="1.5" fill="currentColor"/></svg>`,
		"icon-bloom-filter":    `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="6" width="20" height="12" rx="2"/><line x1="6" y1="6" x2="6" y2="18"/><line x1="10" y1="6" x2="10" y2="18"/><line x1="14" y1="6" x2="14" y2="18"/><line x1="18" y1="6" x2="18" y2="18"/><circle cx="6" cy="10" r="1" fill="currentColor"/><circle cx="14" cy="10" r="1" fill="currentColor"/><circle cx="10" cy="14" r="1" fill="currentColor"/><circle cx="18" cy="14" r="1" fill="currentColor"/></svg>`,
		"icon-snowflake-id":    `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="2" x2="12" y2="22"/><line x1="2" y1="12" x2="22" y2="12"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/><line x1="19.07" y1="4.93" x2="4.93" y2="19.07"/></svg>`,
		"icon-pattern":         `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2a4 4 0 0 1 4 4c0 1.95-1.4 3.57-3.25 3.92L12 22"/><path d="M12 2a4 4 0 0 0-4 4c0 1.95 1.4 3.57 3.25 3.92"/><circle cx="12" cy="14" r="2"/><line x1="8" y1="18" x2="16" y2="18"/></svg>`,
		"icon-rag":             `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/><line x1="8" y1="8" x2="14" y2="8"/><line x1="8" y1="11" x2="14" y2="11"/><line x1="8" y1="14" x2="12" y2="14"/></svg>`,
		"icon-agent-tools":     `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>`,
		"icon-prompt-chain":    `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>`,
		"icon-guardrails":      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="M9 12l2 2 4-4"/></svg>`,
		"icon-embeddings":      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="6" cy="6" r="3"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="18" r="3"/><line x1="9" y1="6" x2="15" y2="6"/><line x1="6" y1="9" x2="6" y2="15"/><line x1="18" y1="9" x2="18" y2="15"/><line x1="9" y1="18" x2="15" y2="18"/><line x1="9" y1="8" x2="15" y2="16"/></svg>`,
		// Category icons (used for sidebar category labels)
		"icon-cat-networking":  `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 0 0-2.91-.09z"/><path d="m12 15-3-3a22 22 0 0 1 2-3.95A12.88 12.88 0 0 1 22 2c0 2.72-.78 7.5-6 11a22.35 22.35 0 0 1-4 2z"/><path d="M9 12H4s.55-3.03 2-4c1.62-1.08 5 0 5 0"/><path d="M12 15v5s3.03-.55 4-2c1.08-1.62 0-5 0-5"/></svg>`,
		"icon-cat-storage":     `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>`,
		"icon-cat-distributed": `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="5" r="3"/><circle cx="5" cy="19" r="3"/><circle cx="19" cy="19" r="3"/><path d="M12 8v13M5 16l7-5M19 16l-7-5"/></svg>`,
		"icon-cat-messaging":   `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 2H8l-4 5h16l-4-5z"/><line x1="12" y1="12" x2="12" y2="16"/><line x1="10" y1="14" x2="14" y2="14"/></svg>`,
		"icon-cat-caching":     `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>`,
		// New fundamental icons
		"icon-sharding":        `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/><line x1="12" y1="8" x2="12" y2="22"/><path d="M3 11l9 3 9-3"/></svg>`,
		"icon-geospatial":      `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="10" r="3"/><path d="M12 2a8 8 0 0 0-8 8c0 5 8 14 8 14s8-9 8-14a8 8 0 0 0-8-8z"/></svg>`,
		"icon-cap-theorem":     `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 22 20 2 20"/><line x1="12" y1="8" x2="12" y2="14"/><line x1="8" y1="14" x2="16" y2="14"/><text x="4" y="19" font-size="4">C</text><text x="12" y="19" font-size="4">A</text><text x="19" y="19" font-size="4">P</text></svg>`,
		"icon-dist-rate-limit": `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/><line x1="6" y1="6" x2="18" y2="18"/></svg>`,
		"icon-saga":            `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="8" width="5" height="8" rx="1"/><rect x="9" y="8" width="5" height="8" rx="1"/><rect x="16" y="8" width="5" height="8" rx="1"/><path d="M7 12h2M14 12h2"/><path d="M4 4l16 16"/></svg>`,
		"icon-circuit-breaker": `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8h1a4 4 0 0 1 0 8h-1"/><path d="M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z"/><line x1="6" y1="1" x2="6" y2="4"/><line x1="10" y1="1" x2="10" y2="4"/><line x1="14" y1="1" x2="14" y2="4"/></svg>`,
		"icon-websocket":       `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14"/><path d="M5 12l4-4m-4 4 4 4"/><path d="M19 12l-4-4m4 4-4 4"/></svg>`,
		"icon-s3":              `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 8l-9-5-9 5v8l9 5 9-5V8z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>`,
		"icon-elasticsearch":   `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/><line x1="8" y1="8" x2="14" y2="8"/><line x1="8" y1="11" x2="14" y2="11"/><line x1="8" y1="14" x2="12" y2="14"/></svg>`,
		"icon-trie":            `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="3" r="2"/><circle cx="5" cy="15" r="2"/><circle cx="19" cy="15" r="2"/><circle cx="9" cy="21" r="2"/><circle cx="15" cy="21" r="2"/><line x1="12" y1="5" x2="5" y2="13"/><line x1="12" y1="5" x2="19" y2="13"/><line x1="5" y1="17" x2="9" y2="19"/><line x1="19" y1="17" x2="15" y2="19"/></svg>`,
		"icon-geohash":         `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="12" y1="3" x2="12" y2="21"/><circle cx="16" cy="16" r="2" fill="currentColor"/></svg>`,
		// New pattern icons
		"icon-auth":            `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/><circle cx="12" cy="16" r="1" fill="currentColor"/></svg>`,
		"icon-payments":        `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="1" y="4" width="22" height="16" rx="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>`,
		"icon-coupons":         `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M20 7H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/><line x1="5" y1="12" x2="5.01" y2="12"/><line x1="19" y1="12" x2="19.01" y2="12"/></svg>`,
		"icon-notifications":   `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>`,
		"icon-content":         `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="20" rx="3"/><line x1="7" y1="8" x2="17" y2="8"/><line x1="7" y1="12" x2="17" y2="12"/><line x1="7" y1="16" x2="13" y2="16"/></svg>`,
		"icon-search":          `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/><line x1="8" y1="8" x2="14" y2="8"/><line x1="8" y1="11" x2="14" y2="11"/></svg>`,
		"icon-reliability":     `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="M9 12l2 2 4-4"/></svg>`,
		// Quick byte icons
		"icon-qb":              `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M9 18h6"/><path d="M10 22h4"/><path d="M12 2a7 7 0 0 1 7 7c0 2.38-1.19 4.47-3 5.74V17a1 1 0 0 1-1 1H9a1 1 0 0 1-1-1v-2.26C6.19 13.47 5 11.38 5 9a7 7 0 0 1 7-7z"/></svg>`,
	}

	slugToIcon := map[string]string{
		// Problems
		"url-shortener":                 "icon-url",
		"rate-limiter":                  "icon-rate-limit",
		"chat-system":                   "icon-chat",
		"instagram":                     "icon-instagram",
		"ticket-booking":                "icon-problem",
		"food-delivery":                 "icon-problem",
		"ride-hailing":                  "icon-problem",
		"search-autocomplete":           "icon-search",
		"twitter-feed":                  "icon-instagram",
		"google-calendar":               "icon-problem",
		"id-generator":                  "icon-snowflake-id",
		"distributed-cache":             "icon-redis",
		"notification-system":           "icon-notifications",
		"file-storage":                  "icon-s3",
		"collaborative-editing":         "icon-websocket",
		"payment-system":                "icon-payments",
		"logging-system":                "icon-elasticsearch",
		"recommendation-system":         "icon-rag",
		// Fundamentals — networking
		"networking/load-balancing":     "icon-load-balancer",
		"networking/load-balancing/alb": "icon-load-balancer",
		"networking/load-balancing/nlb": "icon-load-balancer",
		"networking/cdn":                "icon-cdn",
		"networking/cdn/cloudfront":     "icon-cdn",
		"networking/websockets":         "icon-websocket",
		// Fundamentals — storage
		"storage/redis":                 "icon-redis",
		"storage/dynamodb":              "icon-database",
		"storage/postgres":              "icon-database",
		"storage/sharding":              "icon-sharding",
		"storage/geospatial":            "icon-geospatial",
		"storage/s3":                    "icon-s3",
		"storage/elasticsearch":         "icon-elasticsearch",
		// Fundamentals — distributed
		"distributed/consistent-hashing": "icon-consistent-hash",
		"distributed/cap-theorem":        "icon-cap-theorem",
		"distributed/rate-limiting":      "icon-dist-rate-limit",
		"distributed/saga-pattern":       "icon-saga",
		"distributed/circuit-breaker":    "icon-circuit-breaker",
		// Type defaults
		"problem":                        "icon-problem",
		"fundamental":                    "icon-fundamental",
		"algorithm":                      "icon-algorithm",
		// Algorithms
		"base62-encoding":               "icon-base62",
		"token-bucket":                  "icon-token-bucket",
		"consistent-hashing":            "icon-consistent-hash",
		"bloom-filter":                  "icon-bloom-filter",
		"snowflake-id":                  "icon-snowflake-id",
		"trie":                          "icon-trie",
		"geohash":                       "icon-geohash",
		// Patterns — LLM
		"pattern":                       "icon-pattern",
		"rag":                           "icon-rag",
		"agent-tools":                   "icon-agent-tools",
		"prompt-chaining":               "icon-prompt-chain",
		"guardrails":                    "icon-guardrails",
		"embeddings-vector-search":      "icon-embeddings",
		// Patterns — system
		"auth":                          "icon-auth",
		"payments":                      "icon-payments",
		"coupons":                       "icon-coupons",
		"messaging":                     "icon-messaging",
		"content":                       "icon-content",
		"search":                        "icon-search",
		"reliability":                   "icon-reliability",
		// Quick byte patterns
		"qb-distributed-counters":       "icon-qb",
		"qb-rate-limiting":              "icon-qb",
		"qb-id-generation":              "icon-qb",
		"qb-caching":                    "icon-qb",
		"qb-messaging-kafka":            "icon-qb",
		// Category names (used in welcome page)
		"Networking":          "icon-cat-networking",
		"Storage":             "icon-cat-storage",
		"Distributed":         "icon-cat-distributed",
		"Distributed Systems": "icon-cat-distributed",
		"Messaging":           "icon-cat-messaging",
		"Caching":             "icon-cat-caching",
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

// say renders an interview say-box with speech icon.
func say(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="say-box"><span class="say-label"><svg class="box-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg> Say</span> <em>"%s"</em></div>`, text))
}

// thought renders a thought cloud for reasoning/math/failure scenarios.
func thought(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="thought-cloud">%s</div>`, text))
}

// avoid renders an avoid box for common mistakes.
func avoid(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="avoid-box"><span class="avoid-label"><svg class="box-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg> Avoid</span> %s</div>`, text))
}

// key renders a key takeaway box.
func key(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="key-takeaway"><span class="key-label"><svg class="box-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg> Key Insight</span> %s</div>`, text))
}

// details renders a collapsible details/summary block (for SQL, code, etc.).
func details(summary, lang, content string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<details class="collapsible-code"><summary>%s</summary><div class="code-block" data-lang="%s"><pre><code>%s</code></pre></div></details>`,
		summary, lang, template.HTMLEscapeString(content)))
}

// phase renders a phase/section header with an anchor ID for navigation.
// The data-phase attribute is used by the NFR filter system to map phases to NFRs.
func phase(num int, title, time string) template.HTML {
	slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	return template.HTML(fmt.Sprintf(
		`<div class="phase-header" id="phase-%d-%s" data-phase="%d">
			<span class="phase-number">%d</span>
			<span class="phase-title">%s</span>
			<span class="phase-time">%s</span>
		</div>`, num, slug, num, num, title, time))
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
	escapedDef := template.HTMLEscapeString(definition)
	return template.HTML(fmt.Sprintf(
		`<span class="info-term">%s<span class="info-icon" title="%s">ℹ</span><span class="info-tooltip">%s</span></span>`,
		term, escapedDef, definition))
}

// Option types for compare macro
type CompareOption struct {
	Kind   string // "best", "alt", "nofit"
	Name   string
	Reason string
}

func optBest(name, reason string) CompareOption  { return CompareOption{"best", name, reason} }
func optAlt(name, reason string) CompareOption   { return CompareOption{"alt", name, reason} }
func optNofit(name, reason string) CompareOption { return CompareOption{"nofit", name, reason} }

func optionsList(opts ...CompareOption) []CompareOption { return opts }

// compare renders a color-coded comparison card.
func compare(title string, opts []CompareOption) template.HTML {
	icons := map[string]string{
		"best":  `<svg class="compare-svg-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>`,
		"alt":   `<svg class="compare-svg-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="5" y1="12" x2="19" y2="12"/><polyline points="12 5 19 12 12 19"/></svg>`,
		"nofit": `<svg class="compare-svg-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>`,
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<div class="compare-card"><div class="compare-title">%s</div>`, title))
	for _, opt := range opts {
		icon := icons[opt.Kind]
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

func rowItem(cells ...string) TableRow     { return TableRow{cells} }
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

// makeDiagramFunc creates a diagram rendering function with registry access.
//
// Three modes:
//
//	{{diagram "slug"}}                              — lookup from registry by slug
//	{{diagram "Title" `<div>HTML diagram</div>`}}   — inline HTML mode (backward compat)
//	{{diagram "Title" `ascii art here`}}            — legacy ASCII mode
func makeDiagramFunc(reg *diagrams.Registry) func(string, ...string) template.HTML {
	return func(titleOrSlug string, args ...string) template.HTML {
		// Mode 1: Slug-based lookup (no second argument)
		if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
			if d := reg.Get(titleOrSlug); d != nil {
				return renderDiagram(d)
			}
			// No match — render placeholder
			return template.HTML(fmt.Sprintf(
				`<div class="diagram-container" data-slug="%s">
				<div class="diagram-title">%s</div>
				<div class="diagram-placeholder">Diagram: %s</div>
			</div>`, titleOrSlug, titleOrSlug, titleOrSlug))
		}

		// Mode 2/3: Inline content (backward compat)
		content := strings.TrimSpace(args[0])
		if strings.HasPrefix(content, "<") {
			return template.HTML(fmt.Sprintf(
				`<div class="diagram-container">
				<div class="diagram-title">%s</div>
				%s
			</div>`, titleOrSlug, content))
		}
		// Legacy ASCII art mode
		art := template.HTMLEscapeString(args[0])
		return template.HTML(fmt.Sprintf(
			`<div class="diagram-container">
				<div class="diagram-title">%s</div>
				<pre class="diagram-art">%s</pre>
			</div>`, titleOrSlug, art))
	}
}

// renderDiagram renders a Diagram from the registry with interactive info icon,
// fullscreen toggle, and zoom controls.
func renderDiagram(d *diagrams.Diagram) template.HTML {
	const quickLegend = `<div class="diagram-quick-legend" aria-label="Diagram reading guide">
		<span><strong>Flow:</strong> follow arrows left-to-right or top-to-bottom</span>
		<span><strong>Colors:</strong> green = recommended, amber = caution, red = risk</span>
	</div>`

	// Build info icon with tooltip if description exists
	infoHTML := ""
	if d.Description != "" {
		infoHTML = fmt.Sprintf(
			`<span class="diagram-info" onclick="event.stopPropagation(); this.classList.toggle('active')">
				<svg class="diagram-info-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
				<span class="diagram-tooltip">%s</span>
			</span>`, template.HTMLEscapeString(d.Description))
	}

	// Fullscreen toggle button
	fullscreenBtn := `<button class="diagram-fullscreen-btn" title="Toggle fullscreen" onclick="event.stopPropagation(); var d=this.closest('.diagram-container'); d.classList.toggle('fullscreen'); if(d.classList.contains('fullscreen')){document.body.style.overflow='hidden'}else{document.body.style.overflow=''}">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3"/></svg>
	</button>`

	// Zoom controls
	zoomControls := `<div class="diagram-zoom-controls">
		<button class="diagram-zoom-btn" title="Zoom in" onclick="event.stopPropagation(); var c=this.closest('.diagram-container'); var s=parseFloat(c.dataset.zoom||'1'); s=Math.min(s+0.15,2); c.dataset.zoom=s; c.querySelector('.diagram-body').style.transform='scale('+s+')'; c.querySelector('.diagram-body').style.transformOrigin='top center'">+</button>
		<button class="diagram-zoom-btn" title="Zoom out" onclick="event.stopPropagation(); var c=this.closest('.diagram-container'); var s=parseFloat(c.dataset.zoom||'1'); s=Math.max(s-0.15,0.5); c.dataset.zoom=s; c.querySelector('.diagram-body').style.transform='scale('+s+')'; c.querySelector('.diagram-body').style.transformOrigin='top center'">−</button>
		<button class="diagram-zoom-btn" title="Reset zoom" onclick="event.stopPropagation(); var c=this.closest('.diagram-container'); c.dataset.zoom='1'; c.querySelector('.diagram-body').style.transform='scale(1)'">⊙</button>
	</div>`

	switch d.Type {
	case diagrams.TypeImage:
		return template.HTML(fmt.Sprintf(
			`<div class="diagram-container diagram-interactive" data-slug="%s" data-zoom="1">
				<div class="diagram-header">
					<div class="diagram-title">%s</div>
					%s
					%s
				</div>
				%s
				<div class="diagram-body">
					<div class="diagram-body-overlay"></div>
					<img src="/static/img/diagrams/%s" alt="%s" class="diagram-img" loading="lazy">
				</div>
				%s
			</div>`, d.Slug, d.Title, infoHTML, fullscreenBtn, quickLegend, d.ImagePath, d.Title, zoomControls))
	default: // TypeHTML
		return template.HTML(fmt.Sprintf(
			`<div class="diagram-container diagram-interactive" data-slug="%s" data-zoom="1">
				<div class="diagram-header">
					<div class="diagram-title">%s</div>
					%s
					%s
				</div>
				%s
				<div class="diagram-body">
					<div class="diagram-body-overlay"></div>
					%s
				</div>
				%s
			</div>`, d.Slug, d.Title, infoHTML, fullscreenBtn, quickLegend, d.HTML, zoomControls))
	}
}

// StageItem represents a navigation item for the stage navigator.
type StageItem struct {
	Num   int
	Title string
	ID    string
}

// stageNav renders a sticky navigation bar for problem stages/phases.
// Usage: {{stageNav "Requirements" 1 "MVP" 2 "Upload Flow" 3 ...}}
func stageNav(args ...interface{}) template.HTML {
	var items []StageItem
	for i := 0; i+1 < len(args); i += 2 {
		title := fmt.Sprintf("%v", args[i])
		num := 0
		switch v := args[i+1].(type) {
		case int:
			num = v
		case float64:
			num = int(v)
		}
		slug := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
		slug = strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
				return r
			}
			return -1
		}, slug)
		items = append(items, StageItem{num, title, fmt.Sprintf("phase-%d-%s", num, slug)})
	}

	var sb strings.Builder
	sb.WriteString(`<nav class="stage-nav"><button class="stage-nav-arrow" aria-label="Scroll left" onclick="var i=this.nextElementSibling;i.scrollTo({left:Math.max(0,i.scrollLeft-200),behavior:'smooth'})">&#8249;</button><div class="stage-nav-inner">`)
	for _, item := range items {
		sb.WriteString(fmt.Sprintf(
			`<a href="#%s" class="stage-nav-item" onclick="event.preventDefault(); document.getElementById('%s').scrollIntoView({behavior:'smooth', block:'start'})"><span class="stage-nav-num">%d</span><span class="stage-nav-label">%s</span></a>`,
			item.ID, item.ID, item.Num, item.Title))
	}
	sb.WriteString(`</div></nav>`)
	return template.HTML(sb.String())
}

// anchor creates a named anchor point for navigation.
func anchor(id string) template.HTML {
	return template.HTML(fmt.Sprintf(`<div id="%s" class="anchor-point"></div>`, id))
}

// deepQA renders a nested Q&A section.
// It takes raw HTML content that should use the .dqa-* CSS classes.
// Usage: {{deepQA "Section Title" `<div class="dqa-item">...</div>`}}
func deepQA(title, content string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="dqa-section">
			<div class="dqa-header">
				<svg class="dqa-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
				<span>%s</span>
			</div>
			%s
		</div>`, title, content))
}

// ── Interview Insight Macros ───────────────────────────────────

// mustKnow renders a "Must Know" box with critical interview knowledge.
// Usage: {{mustKnow "point1" "point2" "point3"}}
func mustKnow(items ...string) template.HTML {
	var sb strings.Builder
	sb.WriteString(`<div class="must-know-box"><div class="must-know-header"><svg class="box-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg> Must Know</div><ul>`)
	for _, item := range items {
		sb.WriteString(fmt.Sprintf(`<li>%s</li>`, item))
	}
	sb.WriteString(`</ul></div>`)
	return template.HTML(sb.String())
}

// goodToKnow renders a "Good to Know" box with supplementary knowledge.
// Usage: {{goodToKnow "point1" "point2" "point3"}}
func goodToKnow(items ...string) template.HTML {
	var sb strings.Builder
	sb.WriteString(`<div class="good-to-know-box"><div class="good-to-know-header"><svg class="box-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg> Good to Know</div><ul>`)
	for _, item := range items {
		sb.WriteString(fmt.Sprintf(`<li>%s</li>`, item))
	}
	sb.WriteString(`</ul></div>`)
	return template.HTML(sb.String())
}

// caveat renders a caveat/warning callout box.
// Usage: {{caveat "Important caveat text here"}}
func caveat(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="caveat-box"><span class="caveat-label"><svg class="box-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg> Caveat</span> %s</div>`, text))
}

// collapseSection renders a collapsible section (closed by default) for content
// that's supplementary or less relevant for core interview prep.
// Usage: {{collapseSection "Section Title" "HTML content here"}}
func collapseSection(title, content string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<details class="collapse-section"><summary class="collapse-section-summary"><svg class="collapse-section-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 18l6-6-6-6"/></svg> %s</summary><div class="collapse-section-body">%s</div></details>`,
		title, content))
}

// ── Thought Process Macros ─────────────────────────────────────

// hint renders an inline cloud icon that shows a popup with the thought process.
// Short hint is visible on the icon, click opens a modal with the full explanation.
// Usage: {{hint "short hint" "detailed explanation of why you thought this"}}
func hint(short, detail string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<span class="hint-trigger" onclick="this.querySelector('.hint-popup').classList.toggle('show')">
			<svg class="hint-cloud" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 10h-1.26A8 8 0 1 0 9 20h9a5 5 0 0 0 0-10z"/></svg>
			<span class="hint-label">%s</span>
			<span class="hint-popup" onclick="event.stopPropagation()">
				<span class="hint-popup-header">
					<svg class="hint-popup-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M18 10h-1.26A8 8 0 1 0 9 20h9a5 5 0 0 0 0-10z"/></svg>
					Why this?
					<button class="hint-close" onclick="event.stopPropagation(); this.closest('.hint-popup').classList.remove('show')">&times;</button>
				</span>
				<span class="hint-popup-body">%s</span>
			</span>
		</span>`, short, detail))
}

// ThinkChain represents a single step in a nested thought chain.
type ThinkChain struct {
	Kind    string // "whyNot", "whatIf", "how"
	Title   string
	Content string
}

func whyNotItem(title, content string) ThinkChain { return ThinkChain{"whyNot", title, content} }
func whatIfItem(title, content string) ThinkChain { return ThinkChain{"whatIf", title, content} }
func howItem(title, content string) ThinkChain    { return ThinkChain{"how", title, content} }

// think renders an enhanced thought cloud with nested reasoning chains.
// Usage: {{think "main reasoning" (whyNot "alternative" "reason") (whatIf "scenario" "response") ...}}
func think(mainThought string, chains ...ThinkChain) template.HTML {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<div class="think-block"><div class="think-main">%s</div>`, mainThought))

	if len(chains) > 0 {
		sb.WriteString(`<div class="think-chains">`)
		for _, c := range chains {
			icon := ""
			label := ""
			cssClass := ""
			switch c.Kind {
			case "whyNot":
				icon = `<svg class="think-chain-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>`
				label = "Why not"
				cssClass = "think-why-not"
			case "whatIf":
				icon = `<svg class="think-chain-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`
				label = "What if"
				cssClass = "think-what-if"
			case "how":
				icon = `<svg class="think-chain-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`
				label = "How"
				cssClass = "think-how"
			}
			sb.WriteString(fmt.Sprintf(
				`<details class="think-chain %s">
					<summary class="think-chain-summary">%s <span class="think-chain-label">%s:</span> %s</summary>
					<div class="think-chain-body">%s</div>
				</details>`, cssClass, icon, label, c.Title, c.Content))
		}
		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</div>`)
	return template.HTML(sb.String())
}

// triggerQs renders a collapsible section with a bulb icon showing potential
// interviewer questions that the current section could trigger.
// Usage: {{triggerQs "What might they ask?" "Q1" "A1" "Q2" "A2" ...}}
func triggerQs(title string, qaPairs ...string) template.HTML {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`<details class="trigger-qs">
			<summary class="trigger-qs-summary">
				<svg class="trigger-qs-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
					<path d="M9 18h6"/><path d="M10 22h4"/>
					<path d="M12 2a7 7 0 0 0-4 12.7V17h8v-2.3A7 7 0 0 0 12 2z"/>
				</svg>
				<span>%s</span>
				<span class="trigger-qs-count">%d</span>
			</summary>
			<div class="trigger-qs-body">`, title, len(qaPairs)/2))

	for i := 0; i+1 < len(qaPairs); i += 2 {
		q, a := qaPairs[i], qaPairs[i+1]
		sb.WriteString(fmt.Sprintf(
			`<div class="trigger-qa">
				<div class="trigger-q">
					<svg class="trigger-q-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
					%s
				</div>
				<div class="trigger-a">%s</div>
			</div>`, q, a))
	}

	sb.WriteString(`</div></details>`)
	return template.HTML(sb.String())
}

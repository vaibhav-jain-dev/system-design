# AI Project Guide — System Design Interview Prep Dashboard

> **If you are an AI agent working on this project, read this file completely before doing anything.** This is a self-contained guide that covers architecture, data flow, conventions, and every file you need to know about. It replaces the need to read README.md for your purposes.

---

## Table of Contents

1. [What This Project Does](#what-this-project-does)
2. [How to Run It](#how-to-run-it)
3. [Architecture at a Glance](#architecture-at-a-glance)
4. [Request Lifecycle](#request-lifecycle)
5. [File Map — Every File and What It Does](#file-map)
6. [The Knowledge Graph (Registry)](#the-knowledge-graph)
7. [Content Authoring System](#content-authoring-system)
8. [Macro Reference](#macro-reference)
9. [Diagram System](#diagram-system)
10. [Template System](#template-system)
11. [NFR Filter System](#nfr-filter-system)
12. [CSS Class Conventions](#css-class-conventions)
13. [How to Add a New Problem (Step-by-Step)](#how-to-add-a-new-problem)
14. [How to Add a New Fundamental](#how-to-add-a-new-fundamental)
15. [How to Add a New Algorithm](#how-to-add-a-new-algorithm)
16. [How to Add/Modify a Diagram](#how-to-addmodify-a-diagram)
17. [Content Quality Rules](#content-quality-rules)
18. [Common Pitfalls and Gotchas](#common-pitfalls-and-gotchas)
19. [Testing Your Changes](#testing-your-changes)

---

## What This Project Does

A Go web application that serves **rich, interactive HTML content** for system design interview preparation. It has:

- **Problems** — full system design walkthroughs (URL Shortener, Rate Limiter, Instagram, Chat System, etc.) with 12 phases each
- **Fundamentals** — deep dives into infrastructure topics (Redis, DynamoDB, Load Balancing, CDN, etc.) with 8 phases each
- **Algorithms** — implementation-focused content (Base62, Bloom Filter, Consistent Hashing, etc.) with 8 phases
- **Patterns** — AI/ML design patterns (RAG, Agent Tools, Prompt Chaining, etc.)
- **Concepts** — cross-cutting topics that link across all content types
- **Quick Questions** — rapid-fire interview Q&A organized by category

All content is **cross-linked** via a YAML knowledge graph. Problems reference fundamentals they use, algorithms reference problems they appear in, and reverse links are auto-derived at startup.

---

## How to Run It

```bash
go run main.go
# Server starts at http://localhost:8080
```

Dependencies: Go 1.24+, chi router, google/uuid, gopkg.in/yaml.v3. All vendored or in go.sum.

---

## Architecture at a Glance

```
┌─────────────────────────────────────────────────────────┐
│                     Browser                              │
│  HTMX partial swaps + Alpine.js sidebar collapse         │
└────────────┬────────────────────────────────────────────┘
             │ HTTP
┌────────────▼────────────────────────────────────────────┐
│  main.go — chi router                                    │
│                                                          │
│  Routes:                                                 │
│    GET /                    → Dashboard (welcome page)   │
│    GET /problem/{slug}      → Problem detail             │
│    GET /fund/*              → Fundamental detail          │
│    GET /algo/{slug}         → Algorithm detail            │
│    GET /pattern/{slug}      → Pattern detail              │
│    GET /concept/{slug}      → Concept detail              │
│    GET /quick/{slug}        → Quick category detail       │
│    GET /practice            → Practice page               │
│    GET /highlights          → Highlights dashboard        │
│    POST /api/generate/{slug}→ PDF generation (stub)       │
│    GET /static/*            → Static files (CSS/JS)       │
└────────────┬────────────────────────────────────────────┘
             │
   ┌─────────┼──────────┐
   ▼         ▼          ▼
┌──────┐ ┌───────┐ ┌────────┐
│Registry│ │Macros │ │Diagrams│
│(YAML) │ │(FuncMap)│ │(Registry)│
└──┬───┘ └───┬───┘ └────┬───┘
   │         │          │
   │    ┌────▼──────────▼───┐
   │    │  Go html/template  │
   │    │  Content parsing   │
   │    └────────┬──────────┘
   │             │
   ▼             ▼
┌────────────────────────┐
│  content/**/ index.html │ ← HTML fragments with Go template macros
│  web/templates/*.html   │ ← Layout templates (base, sidebar, detail_*)
└────────────────────────┘
```

**Key insight**: Content files are NOT raw HTML. They are **Go template fragments** that use macros like `{{say "text"}}`, `{{diagram "slug"}}`, `{{phase 1 "Title" "5 min"}}`. The Go server parses each content file as a Go template, executes the macros (which produce HTML), then wraps the result in a layout template.

---

## Request Lifecycle

Here's exactly what happens when a user visits `/problem/rate-limiter`:

1. **Router** (`main.go:56`) matches the route, calls `h.ProblemDetail`
2. **Handler** (`handlers.go:155-216`):
   - Looks up `rate-limiter` in `registry.problemsBySlug`
   - Calls `renderContent("problems/rate-limiter")` which:
     a. Reads `content/problems/rate-limiter/index.html` from the embedded FS
     b. Runs `preprocessContent()` — converts backtick strings in `{{actions}}` to double-quoted strings (Go templates don't support backticks inside `{{ }}`)
     c. Parses the file as a Go template with the macro FuncMap
     d. Executes the template → produces HTML string
   - Builds NFR/FR JSON maps for the JS filter system
   - Checks `HX-Request` header:
     - **If HTMX request** → renders only `detail_problem.html` (partial swap)
     - **If full page load** → renders `base.html` which includes sidebar + detail area
3. **Browser** receives HTML, HTMX handles subsequent navigation as partial swaps

---

## File Map

### Go Source Files

| File | Purpose | Key Functions |
|------|---------|---------------|
| `main.go` | Entry point. Embeds static/template/content FS, loads registry, builds diagram registry, creates FuncMap, wires chi routes | `main()` |
| `internal/registry/registry.go` | Parses `_registry.yaml`, builds knowledge graph with forward/reverse links, resolves cross-references | `Load()`, `GetProblem()`, `GetFundamental()`, `GetAlgorithm()`, `GetPattern()`, `GetConcept()`, `GroupedFundamentals()` |
| `internal/handlers/handlers.go` | All HTTP handlers. Renders content through templates, handles HTMX partial vs full page | `Dashboard()`, `ProblemDetail()`, `FundamentalDetail()`, `AlgorithmDetail()`, `PatternDetail()`, `ConceptDetail()`, `renderContent()`, `preprocessContent()` |
| `internal/macros/macros.go` | Go template FuncMap. Every content macro (`say`, `think`, `hint`, `phase`, `diagram`, `code`, `table`, `compare`, etc.) | `FuncMap()`, `say()`, `think()`, `hint()`, `phase()`, `makeDiagramFunc()`, etc. |
| `internal/diagrams/registry.go` | Diagram struct, Registry type, `BuildDefault()` which registers all diagram domain files | `BuildDefault()`, `Register()`, `Get()` |
| `internal/diagrams/rate_limiter.go` | Rate limiter diagrams (slug prefix: `rl-`) | `registerRateLimiter()` |
| `internal/diagrams/instagram.go` | Instagram diagrams (slug prefix: `ig-`) | `registerInstagram()` |
| `internal/diagrams/url_shortener.go` | URL shortener diagrams (slug prefix: `url-`) | `registerURLShortener()` |
| `internal/diagrams/algorithms.go` | Algorithm diagrams (slug prefix: `algo-`) | `registerAlgorithms()` |
| `internal/diagrams/fundamentals.go` | Fundamental diagrams (slug prefix: `fund-`) | `registerFundamentals()` |
| `internal/diagrams/patterns.go` | Pattern diagrams (slug prefix: `pat-`) | `registerPatterns()` |
| `internal/diagrams/chat_system.go` | Chat system diagrams (slug prefix: `cs-`) | `registerChatSystem()` |
| `internal/diagrams/food_delivery.go` | Food delivery diagrams | `registerFoodDelivery()` |
| `internal/diagrams/ticket_booking.go` | Ticket booking diagrams | `registerTicketBooking()` |
| `internal/diagrams/ride_hailing.go` | Ride hailing diagrams | `registerRideHailing()` |
| `internal/diagrams/search_autocomplete.go` | Search autocomplete diagrams | `registerSearchAutocomplete()` |
| `internal/diagrams/twitter_feed.go` | Twitter feed diagrams | `registerTwitterFeed()` |
| `internal/diagrams/google_calendar.go` | Google Calendar diagrams | `registerGoogleCalendar()` |
| `internal/diagrams/google_docs.go` | Google Docs diagrams | `registerGoogleDocs()` |
| `internal/diagrams/payment_system.go` | Payment system diagrams | `registerPaymentSystem()` |
| `internal/diagrams/id_generator.go` | ID generator diagrams | `registerIDGenerator()` |
| `internal/diagrams/distributed_cache.go` | Distributed cache diagrams | `registerDistributedCache()` |
| `internal/diagrams/notification_system.go` | Notification system diagrams | `registerNotificationSystem()` |
| `internal/diagrams/file_storage.go` | File storage diagrams | `registerFileStorage()` |
| `internal/diagrams/logging_system.go` | Logging system diagrams | `registerLoggingSystem()` |
| `internal/diagrams/recommendation.go` | Recommendation system diagrams | `registerRecommendationSystem()` |

### Template Files (`web/templates/`)

| File | Purpose |
|------|---------|
| `base.html` | Layout shell — sidebar + detail area + book mode JS + stage nav JS + NFR filter JS |
| `sidebar.html` | Collapsible tree navigation (Problems, Fundamentals grouped by category, Algorithms, Patterns) |
| `welcome.html` | Default dashboard view with stats and quick links |
| `detail_problem.html` | Problem detail — renders content + NFR panel + FR panel + context cards (uses fundamentals) + algorithm pills |
| `detail_fund.html` | Fundamental detail — content + sub-topic pills + "Used In" reverse context cards + contextual highlighting |
| `detail_algo.html` | Algorithm detail — content + "Used in" problem pills |
| `detail_pattern.html` | Pattern detail — content only |
| `detail_concept.html` | Concept detail — shows all appearances as clickable links |
| `detail_quick.html` | Quick category Q&A page |
| `detail_quick_all.html` | All quick categories on one page |
| `detail_practice.html` | Interactive practice/solution-checker page |
| `detail_highlights.html` | Highlights review dashboard |
| `context_card.html` | Reusable cross-link card showing config/why/not_this/risk/caveats |
| `doc_card.html` | PDF generation status card (placeholder) |

### Content Files (`content/`)

| Path Pattern | Type | Count |
|-------------|------|-------|
| `content/_registry.yaml` | Knowledge graph definition | 1 |
| `content/problems/{slug}/index.html` | Problem content (12 phases) | 18 |
| `content/fundamentals/{category}/{slug}/index.html` | Fundamental content (8 phases) | ~15 |
| `content/algorithms/{slug}/index.html` | Algorithm content (8 phases) | 7 |
| `content/patterns/{slug}/index.html` | Pattern content | 5 |

### Static Files (`web/static/`)

| File | Purpose |
|------|---------|
| `css/style.css` | All styling — layout, sidebar, diagrams, macros, book mode, mobile responsive, NFR filter |
| `js/htmx.min.js` | Vendored HTMX (~14KB) for partial page swaps |
| `js/alpine.min.js` | Vendored Alpine.js (~15KB) for sidebar collapse/expand |

---

## The Knowledge Graph

The YAML registry (`content/_registry.yaml`) defines four content types and their relationships:

### Data Types (defined in `internal/registry/registry.go`)

```
Problem
  ├── Slug, Title, Description, Path
  ├── Category ("distributed" or empty for core)
  ├── NFRs[] → ProblemNFR{Slug, Phases[], Title, Color}
  ├── FRs[] → FunctionalRequirement{Slug, Title, Phases[]}
  ├── Uses[] → UsageLink (forward links to fundamentals)
  └── Algorithms[] → *Algorithm (auto-derived reverse links)

Fundamental
  ├── Slug (path-like: "storage/redis", "networking/cdn/cloudfront")
  ├── Title, Description, Path
  ├── Children[] → Fundamental (hierarchical sub-topics)
  ├── UsedBy[] → UsageLink (auto-derived reverse links from problems)
  └── RelatedAlgorithm → *Algorithm

Algorithm
  ├── Slug, Title, Description, Path
  ├── UsedIn[] → string (problem slugs from YAML)
  ├── UsedInProblems[] → *Problem (resolved pointers)
  └── RelatedFundamental → *Fundamental

Pattern
  └── Slug, Title, Description, Path

UsageLink (bidirectional problem↔fundamental link)
  ├── Fundamental (slug), Problem (slug)
  ├── Config, Why, NotThis, Risk, Caveats (rich context)
  ├── NFRs[] (which NFRs this use addresses)
  └── FundamentalRef, ProblemRef (resolved pointers)
```

### How Reverse Links Work

At startup, `registry.Load()`:
1. Indexes all fundamentals (including nested children) into `fundamentalsBySlug`
2. For each problem's `Uses[]` entries, creates reverse `UsageLink` on the fundamental's `UsedBy[]`
3. Propagates reverse links up to ancestor fundamentals (e.g., if a problem uses `networking/cdn/cloudfront`, the link also appears on `networking/cdn`)
4. For each algorithm's `UsedIn[]`, resolves problem pointers and adds reverse `Algorithms[]` to the problem

**You never need to manually specify reverse links.** Just add `uses:` entries on problems and `used_in:` on algorithms — everything else is auto-derived.

---

## Content Authoring System

Content files are Go template fragments — no `<html>`, no `<head>`, no `<body>`. Just the body content using macros.

### Problem File Structure (12 phases, required)

```
{{stageNav "Phase 1 Title" 1 "Phase 2 Title" 2 ... "Interview Deep-Dive Q&A" 12}}

{{phase 1 "Requirements & Problem Scope" "5 min"}}
{{say "Opening statement..."}}
{{diagram "prefix-requirements"}}
{{think "Reasoning..." (whyNot "X" "reason") (whatIf "Y" "response")}}
{{hint "short" "detailed explanation"}}
{{triggerQs "Title" "Q1" "A1" "Q2" "A2"}}

... phases 2-11 ...

{{phase 12 "Interview Deep-Dive Q&A" "10 min"}}
{{deepQA "Title" `<div class="dqa-item">...</div>`}}

{{key "One-sentence takeaway."}}
```

### Fundamental / Algorithm File Structure (8 phases, no stageNav, no deepQA)

```
{{phase 1 "What Is {Topic}?" "3 min"}}
{{say "Opening explanation..."}}
{{diagram "fund-topic-overview"}}
{{think "..." (whyNot "..." "...") (how "..." "...")}}
{{hint "..." "..."}}

... phases 2-8 ...

{{checklist "Point 1" "Point 2"}}
{{key "Key takeaway."}}
```

---

## Macro Reference

Every macro is defined in `internal/macros/macros.go` and available via the FuncMap.

| Macro | Arguments | Renders As |
|-------|-----------|------------|
| `{{say "text"}}` | string (supports HTML inside) | Green-border italic quote box with speech icon |
| `{{think "main" (whyNot "X" "R") (whatIf "Y" "R") (how "Q" "A")}}` | string + variadic ThinkChain | Expandable thought block with red/amber/blue sub-chains |
| `{{hint "short" "detail"}}` | two strings | Inline cloud icon, click opens popup |
| `{{triggerQs "Title" "Q1" "A1" "Q2" "A2" ...}}` | string + variadic Q/A pairs | Collapsible bulb section with Q&A |
| `{{thought "text"}}` | string | Gray cloud box (legacy — prefer `think`) |
| `{{avoid "text"}}` | string | Red-border warning box |
| `{{key "text"}}` | string | Blue-border bold key insight box |
| `{{phase N "Title" "Time"}}` | int, string, string | Phase header with anchor ID `phase-N-slug` and `data-phase="N"` |
| `{{code "lang" "content"}}` | string, string | Dark code block |
| `{{diagram "slug"}}` | string | Looks up diagram from registry, renders with fullscreen/zoom controls |
| `{{diagram "Title" "<html>"}}` | string, string | Inline HTML diagram (backward compat) |
| `{{table "Title" (rows (row "a" "b") (row "c" "d"))}}` | string + []TableRow | Styled table (first row = header) |
| `{{compare "Title" (options (best "X" "R") (alt "Y" "R") (nofit "Z" "R"))}}` | string + []CompareOption | Color-coded card (green/amber/red) |
| `{{qa "Q" "A"}}` or `{{qa "Q" "A" "FQ" "FA"}}` | variadic strings | Q&A card with optional follow-up |
| `{{followup "Q" "A"}}` | two strings | Amber-border follow-up card |
| `{{checklist "item1" "item2" ...}}` | variadic strings | Green checklist with check marks |
| `{{info "term" "definition"}}` | two strings | Inline term with hover tooltip |
| `{{details "summary" "lang" "code"}}` | three strings | Collapsible `<details>` with code |
| `{{stageNav "Title1" 1 "Title2" 2 ...}}` | variadic (string, int pairs) | Sticky horizontal phase nav bar |
| `{{anchor "id"}}` | string | Invisible anchor div |
| `{{deepQA "Title" "<html>"}}` | two strings | Deep Q&A section using `.dqa-*` classes |
| `{{mustKnow "p1" "p2" ...}}` | variadic strings | "Must Know" box with bullet list |
| `{{goodToKnow "p1" "p2" ...}}` | variadic strings | "Good to Know" info box |
| `{{caveat "text"}}` | string | Caveat/warning callout |
| `{{collapseSection "Title" "HTML"}}` | two strings | Collapsible section (closed by default) |

### Helper Functions (used inside macros)

| Function | Purpose |
|----------|---------|
| `options(...)` | Wraps CompareOption items into a slice for `compare` |
| `best("name", "reason")` | Green option for `compare` |
| `alt("name", "reason")` | Amber option for `compare` |
| `nofit("name", "reason")` | Red option for `compare` |
| `rows(...)` | Wraps TableRow items into a slice for `table` |
| `row("a", "b", ...)` | Creates a table row |
| `whyNot("title", "content")` | Red chain for `think` |
| `whatIf("title", "content")` | Amber chain for `think` |
| `how("title", "content")` | Blue chain for `think` |

---

## Diagram System

Diagrams live in Go files under `internal/diagrams/`, NOT inline in content HTML.

### How It Works

1. Each domain has a Go file (e.g., `rate_limiter.go`) that registers diagrams via `r.Register(&Diagram{...})`
2. `BuildDefault()` in `registry.go` calls all domain registration functions
3. The `FuncMap` gets a `diagram` function that looks up slugs from the registry
4. Content files reference diagrams: `{{diagram "rl-architecture"}}`

### Diagram Struct

```go
type Diagram struct {
    Slug        string  // "rl-architecture"
    Title       string  // "Architecture"
    Description string  // Tooltip text
    ContentFile string  // "problems/rate-limiter"
    Type        Type    // TypeHTML or TypeImage
    HTML        string  // Raw HTML using CSS diagram classes
    ImagePath   string  // For TypeImage: filename in /static/img/diagrams/
}
```

### Slug Prefix Convention

| Prefix | Domain | File |
|--------|--------|------|
| `rl-` | Rate Limiter | `rate_limiter.go` |
| `ig-` | Instagram | `instagram.go` |
| `url-` | URL Shortener | `url_shortener.go` |
| `algo-` | Algorithms | `algorithms.go` |
| `fund-` | Fundamentals | `fundamentals.go` |
| `pat-` | Patterns | `patterns.go` |
| `cs-` | Chat System | `chat_system.go` |
| `fd-` | Food Delivery | `food_delivery.go` |
| `tb-` | Ticket Booking | `ticket_booking.go` |
| `rh-` | Ride Hailing | `ride_hailing.go` |
| `sa-` | Search Autocomplete | `search_autocomplete.go` |
| `tf-` | Twitter Feed | `twitter_feed.go` |
| `gc-` | Google Calendar | `google_calendar.go` |
| `gd-` | Google Docs | `google_docs.go` |
| `pay-` | Payment System | `payment_system.go` |
| `idg-` | ID Generator | `id_generator.go` |
| `dc-` | Distributed Cache | `distributed_cache.go` |
| `ns-` | Notification System | `notification_system.go` |
| `fs-` | File Storage | `file_storage.go` |
| `log-` | Logging System | `logging_system.go` |
| `rec-` | Recommendation | `recommendation.go` |

### Diagram CSS Classes

**Layout**: `.d-cols` (CSS grid), `.d-col`, `.d-flow` (horizontal flex), `.d-flow-v` (vertical flex), `.d-row`, `.d-branch` / `.d-branch-arm`

**Boxes**: `.d-box` + color: `.blue`, `.green`, `.purple`, `.amber`, `.red`, `.gray`, `.indigo`

**Grouping**: `.d-group` (dashed border), `.d-group-title`

**Arrows**: `.d-arrow` (horizontal `→`), `.d-arrow-down` (vertical `↓`)

**Entity/DB**: `.d-entity`, `.d-entity-header [color]`, `.d-entity-body`, `.pk` (PK badge), `.fk` (FK badge), `.idx` + `.idx-btree`/`.idx-hash`/`.idx-gin`

**Specialized**: `.d-bitfield` / `.d-bitfield-segment`, `.d-ring` / `.d-ring-node`, `.d-subproblem [color]`

---

## Template System

### HTMX Partial Swaps

When a user clicks a sidebar link, HTMX sends a request with `HX-Request: true` header. The handler detects this and renders ONLY the detail template (not the full page with sidebar). HTMX swaps just the detail area.

For direct URL access (no HTMX header), the handler renders the full `base.html` layout.

### Template Data

Every handler passes a `map[string]interface{}` with:
- `Problems`, `CoreProblems`, `DistributedProblems` — for sidebar
- `Fundamentals`, `FundamentalGroups` — for sidebar
- `Algorithms`, `Patterns`, `Concepts`, `QuickCategories` — for sidebar
- `ActiveSlug` — highlights current item in sidebar
- `PageType` — determines which detail template to render
- Page-specific: `Problem`, `Fundamental`, `Algorithm`, `Pattern`, `Content` (rendered HTML)

### The `preprocessContent` Function

**Critical gotcha**: Go's `html/template` does NOT support backtick raw strings inside `{{ }}` actions. But content files use backticks heavily for multi-line strings in macros like `{{say \`...\`}}` and `{{deepQA "title" \`...\`}}`.

The `preprocessContent()` function in `handlers.go` scans the source and converts backtick strings inside template actions to properly escaped double-quoted strings BEFORE parsing. This is why backticks work in content files even though Go templates don't natively support them.

---

## NFR Filter System

Each problem has Non-Functional Requirements (NFRs) and Functional Requirements (FRs) mapped to phases.

### 8 Standard NFRs

| Slug | Color | What It Covers |
|------|-------|---------------|
| `scalability` | Indigo (#6366F1) | Horizontal scaling, sharding, capacity |
| `performance` | Blue (#2563EB) | Latency, caching, CDN, hot-path |
| `availability` | Green (#059669) | Fault tolerance, failover, multi-AZ |
| `consistency` | Amber (#D97706) | CAP tradeoffs, strong/eventual |
| `durability` | Purple (#7C3AED) | Persistence, replication, WAL |
| `security` | Red (#DC2626) | Auth, encryption, abuse prevention |
| `cost` | Gray (#64748B) | Cost analysis, capacity planning |
| `observability` | Cyan (#0891B2) | Monitoring, tracing, alerting |

### How It Works

1. YAML defines `nfrs:` with phase mappings per problem
2. Handler builds `PhaseNFRMapJSON` (phase → NFR slugs) and `UseNFRMapJSON` (fundamental → NFR slugs)
3. Template renders NFR chips in a panel
4. JavaScript (in `base.html`) wires click handlers that dim/restore phase sections and context cards based on selected NFRs

### NFR Tagging Rules

- **Never tag Phase 1** (requirements) unless the problem is fundamentally about that NFR
- **Never tag Phase 12** (Interview Q&A) — it covers everything
- Tag a phase ONLY if its PRIMARY content is about that NFR
- Max 2 NFRs per `uses:` entry (pick the most dominant)
- If any single NFR is mapped to 7+ phases, you're over-tagging

---

## CSS Class Conventions

### Content Highlighting

| Class | Color | Use For |
|-------|-------|---------|
| `<span class="hl">` | Amber | Key numbers, thresholds |
| `<span class="hl-blue">` | Blue | Concepts, algorithm names |
| `<span class="hl-red">` | Red | Warnings, anti-patterns |

### Text Formatting

- `<strong>` — bold key terms on first use and important numbers
- `<small>` — caveats, cost notes, side details
- `<ul><li>` — bullet lists (prefer over paragraphs for 3+ ideas)

### API Documentation Component (`.api-doc`)

Used in problems with API design phases. Two tabs:
- **Quick** — monospace table (`api-quick`, `api-quick-row`, `api-method`, `api-quick-path`, `api-quick-comment`)
- **Detailed** — cards with request tables and response JSON (`api-card`, `api-req-table`, `api-responses`, `api-response-body`)

JSON token classes: `.jk` (key, blue), `.jv` (string value, green), `.jn` (number, orange), `.jb` (boolean, purple), `.jp` (punctuation, gray), `.jc` (comment, muted)

---

## How to Add a New Problem

### Step 1: Create the diagram file

Create `internal/diagrams/{domain}.go`:

```go
package diagrams

func registerMyDomain(r *Registry) {
    r.Register(&Diagram{
        Slug:        "md-requirements",
        Title:       "Requirements Overview",
        Description: "Functional and non-functional requirements",
        ContentFile: "problems/my-domain",
        Type:        TypeHTML,
        HTML:        `<div class="d-flow">...</div>`,
    })
    // Register more diagrams...
}
```

### Step 2: Register in BuildDefault

Edit `internal/diagrams/registry.go`, add to `BuildDefault()`:

```go
registerMyDomain(r)
```

### Step 3: Add to _registry.yaml

```yaml
problems:
  - slug: my-domain
    title: "My Domain"
    description: "One-line description"
    path: problems/my-domain
    category: ""  # or "distributed"
    nfrs:
      - slug: scalability
        phases: [3, 5, 7]
      - slug: performance
        phases: [2, 4, 5]
      # ... more NFRs
    functional_requirements:
      - slug: core-feature
        title: "Core Feature"
        phases: [2, 3]
      # ... more FRs
    uses:
      - fundamental: storage/redis
        config: "Specific config"
        why: "Why chosen"
        not_this: "Rejected alternative"
        risk: "Key risk"
        caveats: "Cost/operational note"
        nfrs: [performance, availability]
      # ... more fundamentals
```

### Step 4: Create content file

Create `content/problems/my-domain/index.html` following the 12-phase structure with `stageNav`, all macros, and Phase 12 `deepQA`.

### Step 5: Verify

Run `go run main.go` and check:
- No startup errors (missing diagrams, broken YAML, etc.)
- Content renders at `/problem/my-domain`
- Sidebar shows the new problem
- Context cards appear with correct links
- NFR filter dims/shows correct phases

---

## How to Add a New Fundamental

### Step 1: Add to _registry.yaml

```yaml
fundamentals:
  - slug: category/my-topic
    title: "My Topic"
    description: "One-line description"
    path: fundamentals/category/my-topic
    children:  # optional
      - slug: category/my-topic/subtopic
        title: "Subtopic"
        path: fundamentals/category/my-topic/subtopic
```

### Step 2: Create content file

Create `content/fundamentals/category/my-topic/index.html` — 8 phases, no stageNav, no deepQA.

### Step 3: Add uses entries

If any existing problem uses this fundamental, add a `uses:` entry in that problem's YAML definition.

### Step 4: Add diagrams

Register diagrams in the appropriate Go file (usually `fundamentals.go` with `fund-` prefix).

---

## How to Add a New Algorithm

Same as fundamentals but:
1. Add under `algorithms:` in YAML with `used_in: [problem-slug-1, problem-slug-2]`
2. Content goes in `content/algorithms/{slug}/index.html`
3. Diagrams use `algo-` prefix in `algorithms.go`

---

## How to Add/Modify a Diagram

1. Find which Go file owns the slug prefix (see table above)
2. Add/edit a `Register()` call:

```go
r.Register(&Diagram{
    Slug:        "prefix-name",
    Title:       "Human Title",
    Description: "What this shows",
    ContentFile: "problems/slug-or-fundamentals/path",
    Type:        TypeHTML,
    HTML:        `<div class="d-flow">
        <div class="d-box blue">A</div>
        <div class="d-arrow">→</div>
        <div class="d-box green">B</div>
    </div>`,
})
```

3. Reference in content: `{{diagram "prefix-name"}}`
4. Restart server

**Duplicate slug = panic at startup.** Always use unique slugs.

---

## Content Quality Rules

### Writing Style

- **Bullet points over paragraphs** — break 3+ ideas into `<ul><li>`
- **Bold key terms** with `<strong>` on first use
- **Highlight numbers** with `<span class="hl">`
- **No filler** — no "great question", no "let me think", start with the answer
- **No subjective ratings** — replace "excellent" with actual numbers
- **Q&A answers: fact first** — open with the technical answer, then justify
- **`{{say}}` blocks**: break long ones into bullet lists for scannable delivery
- **`{{hint}}` max 2-3 sentences** — use `{{think}}` for longer reasoning

### Structural Requirements

| Type | Phases | stageNav | deepQA (Phase 12) | hints/phase | triggerQs |
|------|--------|----------|--------------------| ------------|-----------|
| Problem | 12 | Required | Required (5 items, 3 nesting levels) | 2-3 | 5-6 phases |
| Fundamental | 8 | No | No | 1-2 | 3-4 phases |
| Algorithm | 8 | No | No | 1-2 | 3-4 phases |

### Every problem phase opens with `{{say "..."}}` — what to literally say in the interview.

---

## Common Pitfalls and Gotchas

### 1. Backtick Preprocessing

Content files can use backticks inside `{{ }}` actions because `preprocessContent()` converts them to double-quoted strings. But this means:
- **Backticks inside backtick strings don't work** — you can't nest them
- **The conversion escapes `\`, `"`, `\n`, `\t`** — if your content has these, they'll be escaped
- If you see weird rendering, check that your backtick strings are properly closed

### 2. Embedded Filesystem

All files are embedded via `//go:embed`. This means:
- Changes to content/template/static files require **restarting the server**
- The server reads from the binary, not from disk at runtime
- New files must be under the embedded directories

### 3. Fundamental Slug Format

Fundamental slugs are **path-like**: `storage/redis`, `networking/cdn/cloudfront`. This affects:
- Route matching: `/fund/*` uses chi wildcard, not `{slug}`
- YAML references: problems reference fundamentals by their full path slug
- Image paths: `slugImagePath()` uses an explicit whitelist because `/` breaks URL patterns

### 4. Diagram Registry Panics

Registering a diagram with a duplicate slug causes a **panic at startup**. Always check existing slugs before adding new ones.

### 5. HTMX Partial Swap

If your template changes affect the sidebar, they won't show up during HTMX navigation — only on full page reload. The sidebar is only rendered with `base.html`.

### 6. NFR Phase Tagging

Over-tagging is the most common mistake. Phase 1 (requirements) and Phase 12 (Q&A) should almost never be tagged. A phase should be tagged ONLY if its PRIMARY content teaches that NFR.

### 7. Context Card Highlighting

When a user navigates from a problem to a fundamental (`?from=problem-slug`), the handler extracts keywords from the `config` field and highlights matching content. Make sure your `config` field contains meaningful keywords that actually appear in the fundamental's content.

### 8. Category Field on Problems

The `category` field on problems is used for sidebar grouping:
- Empty or missing → appears under "Core Problems"
- `"distributed"` → appears under "Distributed Systems"

---

## Testing Your Changes

```bash
# Build and run — catches template parse errors, YAML issues, duplicate diagram slugs
go run main.go

# Run tests (if any exist)
go test ./...

# Check for compile errors
go build ./...
```

**What to verify manually:**
1. Server starts without errors in the log
2. Navigate to your new/modified content page
3. All diagrams render (no "Diagram: slug" placeholder text)
4. Sidebar shows correct hierarchy
5. Cross-links work (context cards, algorithm pills, "Used In" sections)
6. NFR filter dims correct phases
7. HTMX navigation works (click sidebar items, check partial swap)

---

## Quick Reference: Where Things Live

| I want to... | Look at... |
|--------------|-----------|
| Add a route | `main.go` (chi router setup) |
| Add a handler | `internal/handlers/handlers.go` |
| Add a macro | `internal/macros/macros.go` |
| Add a diagram | `internal/diagrams/{domain}.go` + `registry.go:BuildDefault()` |
| Add content | `content/{type}/{slug}/index.html` |
| Add to knowledge graph | `content/_registry.yaml` |
| Change layout/sidebar | `web/templates/base.html` or `sidebar.html` |
| Change detail rendering | `web/templates/detail_{type}.html` |
| Change styling | `web/static/css/style.css` |
| Add a new NFR type | `internal/registry/registry.go:StandardNFRs` |
| Understand data types | `internal/registry/registry.go` (all structs at top) |

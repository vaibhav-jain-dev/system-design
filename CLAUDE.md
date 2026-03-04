# System Design Interview Prep Dashboard

## What This Project Is

A Go web application for system design interview preparation. It serves **rich HTML content** for system design problems (URL Shortener, Rate Limiter, Instagram) and infrastructure fundamentals (Load Balancing, Redis, DynamoDB, CDN), with deep cross-linking between them. A Python subprocess handles PDF export.

**Run it:** `go run main.go` → `http://localhost:8080`

---

## Architecture Overview

```
User Browser
    │
    ├── GET /                          → Dashboard (sidebar + welcome)
    ├── GET /problem/{slug}            → Problem detail page
    ├── GET /fund/{slug...}            → Fundamental detail (slug is path-like: "networking/load-balancing")
    ├── GET /algo/{slug}               → Algorithm detail
    ├── GET /pattern/{slug}            → Design pattern detail
    ├── POST /api/generate/{slug}      → Generate PDF (Go spawns Python subprocess)
    ├── GET /api/status/{taskID}       → PDF generation status (SSE)
    └── GET /pdf/{filename}            → Serve generated PDF
```

**Request flow:** Go server (chi router) → parses YAML registry → renders content HTML files through Go templates (applying macros like `{{say}}`, `{{diagram}}`) → wraps in layout template → serves to browser. HTMX handles partial page swaps for SPA-like navigation.

---

## Project Structure

```
system-design/
├── main.go                              # Entry point: loads registry, builds diagram registry, wires routes
├── go.mod / go.sum
│
├── internal/
│   ├── registry/registry.go             # Parses _registry.yaml → builds knowledge graph with reverse links
│   ├── handlers/handlers.go             # All HTTP handlers: Dashboard, ProblemDetail, FundamentalDetail, etc.
│   ├── macros/macros.go                 # Go template FuncMap: say, thought, avoid, key, phase, code, qa, diagram, etc.
│   └── diagrams/                        # Centralized diagram registry (102 diagrams total)
│       ├── registry.go                  # Diagram struct, Registry type, BuildDefault()
│       ├── rate_limiter.go              # 19 diagrams (slug prefix: rl-)
│       ├── instagram.go                 # 24 diagrams (slug prefix: ig-)
│       ├── url_shortener.go             # 21 diagrams (slug prefix: url-)
│       ├── algorithms.go                # 13 diagrams (slug prefix: algo-)
│       ├── fundamentals.go              # 15 diagrams (slug prefix: fund-)
│       ├── patterns.go                  # 10 diagrams (slug prefix: pat-)
│       └── README.md                    # Diagram library documentation
│
├── web/
│   ├── templates/
│   │   ├── base.html                    # Layout shell: sidebar + detail area + book mode JS + stage nav JS
│   │   ├── sidebar.html                 # Collapsible tree: Problems, Fundamentals (hierarchical), Algorithms, Patterns
│   │   ├── welcome.html                 # Default view with stats and quick links
│   │   ├── detail_problem.html          # Problem: content + "Uses Fundamentals" context cards
│   │   ├── detail_fund.html             # Fundamental: content + sub-topic pills + "Used In" reverse cards + highlight
│   │   ├── detail_algo.html             # Algorithm: content + "Used in" problem pills
│   │   ├── detail_pattern.html          # Pattern: content only
│   │   ├── context_card.html            # Reusable cross-link card (config, why, not_this, risk, caveats)
│   │   └── doc_card.html                # PDF status card (placeholder)
│   └── static/
│       ├── css/style.css                # All styling: layout, sidebar, diagrams, macros, book mode, mobile responsive
│       └── js/
│           ├── htmx.min.js              # Vendored HTMX for partial page swaps
│           └── alpine.min.js            # Vendored Alpine.js for sidebar collapse/expand
│
├── content/                             # HTML content files (Go template fragments, no <html>/<head>)
│   ├── _registry.yaml                   # Knowledge graph: problems, fundamentals, algorithms, patterns + relationships
│   ├── problems/
│   │   ├── url-shortener/index.html     # ~441 lines, uses macros + diagram slugs
│   │   ├── instagram/index.html         # ~603 lines
│   │   └── rate-limiter/index.html      # ~424 lines
│   ├── fundamentals/
│   │   ├── networking/
│   │   │   ├── load-balancing/index.html
│   │   │   ├── load-balancing/alb/index.html
│   │   │   ├── load-balancing/nlb/index.html
│   │   │   ├── cdn/index.html
│   │   │   └── cdn/cloudfront/index.html
│   │   └── storage/
│   │       ├── redis/index.html
│   │       └── dynamodb/index.html
│   ├── algorithms/
│   │   ├── base62-encoding/index.html
│   │   ├── bloom-filter/index.html
│   │   ├── consistent-hashing/index.html
│   │   ├── snowflake-id/index.html
│   │   └── token-bucket/index.html
│   └── patterns/
│       ├── rag/index.html
│       ├── agent-tools/index.html
│       ├── prompt-chaining/index.html
│       ├── guardrails/index.html
│       └── embeddings-vector-search/index.html
│
├── engine/
│   ├── generate_pdf.py                  # Python: Playwright HTML→PDF (A4, page numbers, appendix sections)
│   └── requirements.txt
│
└── output/                              # Generated PDFs (gitignored)
```

---

## Knowledge Graph (content/_registry.yaml)

The YAML registry defines **four content types** and their relationships:

### Content Types

| Type | Slug Pattern | Route | Example |
|------|-------------|-------|---------|
| **Problem** | simple slug | `/problem/{slug}` | `url-shortener`, `rate-limiter`, `instagram` |
| **Fundamental** | path-like slug | `/fund/{slug...}` | `networking/load-balancing`, `storage/redis` |
| **Algorithm** | simple slug | `/algo/{slug}` | `base62-encoding`, `consistent-hashing` |
| **Pattern** | simple slug | `/pattern/{slug}` | `rag`, `agent-tools` |

### Relationships

**Problem → Fundamental** (via `uses` array, with rich context):
```yaml
problems:
  - slug: url-shortener
    uses:
      - fundamental: storage/redis          # FK to fundamental slug
        config: "ElastiCache read-through, TTL 24h"   # How it's configured
        why: "Hot URLs (top 20%) account for 80% reads" # Why this choice
        not_this: "Memcached — no persistence"          # What NOT to use
        risk: "Thundering herd on cache miss"            # Risk to mention
        caveats: "r6g.large = ~$92/mo"                   # Cost/operational note
```

**Fundamental → Problem** (auto-derived reverse links): At startup, `registry.Load()` scans all `problem.Uses` entries and populates each `Fundamental.UsedBy[]` with reverse `UsageLink` structs. No manual reverse linking needed.

**Fundamental → Children** (hierarchical): Fundamentals can have `children[]` for sub-topics (e.g., Load Balancing → ALB, NLB). Children are full fundamentals with their own content pages.

**Algorithm → Problem** (via `used_in` array): Simple slug references, resolved to `*Problem` pointers at startup.

### Go Data Types (internal/registry/registry.go)

```go
type Problem struct {
    Slug, Title, Description, Path string
    Uses []UsageLink                        // Forward links to fundamentals
}

type Fundamental struct {
    Slug, Title, Description, Path string
    Children []Fundamental                   // Sub-topics (hierarchical)
    UsedBy   []UsageLink                     // Reverse links from problems (auto-derived)
}

type Algorithm struct {
    Slug, Title, Description, Path string
    UsedIn         []string                  // Problem slugs (from YAML)
    UsedInProblems []*Problem                // Resolved pointers (auto-derived)
}

type Pattern struct {
    Slug, Title, Description, Path string
}

type UsageLink struct {
    Fundamental    string                    // Forward: fundamental slug
    Problem        string                    // Reverse: problem slug (auto-filled)
    Config, Why, NotThis, Risk, Caveats string  // Rich context fields
    FundamentalRef *Fundamental              // Resolved pointer
    ProblemRef     *Problem                  // Resolved pointer
}
```

---

## Content Authoring System

Content files are **Go template fragments** — no `<html>` or `<head>`, just body content using macros. The Go server parses each content file as a Go template, executes macros, then wraps in the layout.

### Available Macros (internal/macros/macros.go)

| Macro | Purpose | Renders As |
|-------|---------|------------|
| `{{say "text"}}` | What to literally say in interview | Green-border italic quote box |
| `{{thought "text"}}` | Side reasoning, math, failure scenarios (legacy) | Gray cloud box |
| `{{think "main" (whyNot "X" "reason") (whatIf "Y" "resp") (how "Q" "A")}}` | Enhanced thought with nested reasoning chains | Expandable thought block with red/amber/blue sub-chains |
| `{{hint "short" "detail"}}` | Inline thought-process trigger (cloud icon) | Click to show popup explaining "why this decision" |
| `{{triggerQs "Title" "Q1" "A1" "Q2" "A2" ...}}` | Potential interviewer questions per section | Collapsible bulb section with Q&A pairs |
| `{{avoid "text"}}` | Common mistakes, "never say X" | Red-border warning box |
| `{{key "text"}}` | Key takeaway, one-liner | Blue-border bold box |
| `{{phase N "Title" "Time"}}` | Section header with number + time badge | Phase header with anchor ID `phase-N-slug` |
| `{{code "lang" "content"}}` | Code block | Dark code block with syntax highlighting |
| `{{diagram "slug"}}` | Diagram from registry | Diagram container with title + HTML/image content |
| `{{diagram "Title" "<html>"}}` | Inline diagram (backward compat) | Same container, inline HTML |
| `{{table "Title" (rows (row "a" "b") (row "c" "d"))}}` | Data table | Styled table (first row = header) |
| `{{compare "Title" (options (best "X" "reason") (alt "Y" "reason") (nofit "Z" "reason"))}}` | Decision comparison | Color-coded card (green/amber/red) |
| `{{qa "Q" "A" "FQ" "FA"}}` | Interviewer Q&A | Q&A card with optional follow-up |
| `{{followup "Q" "A"}}` | Likely follow-up question | Amber-border card |
| `{{checklist "item1" "item2"}}` | Section summary checklist | Green checklist with check marks |
| `{{info "term" "definition"}}` | Inline tooltip | Term with hover tooltip |
| `{{details "summary" "lang" "code"}}` | Collapsible code block | `<details>` with syntax-highlighted code |
| `{{stageNav "Title1" 1 "Title2" 2 ...}}` | Sticky phase navigation bar | Horizontal nav linking to phase anchors |
| `{{anchor "id"}}` | Named anchor point | Invisible div with ID |
| `{{deepQA "Title" "<html>"}}` | Deep Q&A section | Section with custom HTML using `.dqa-*` CSS classes |

### Content File Example

```html
{{stageNav "Requirements" 1 "Architecture" 2 "Data Model" 3}}

{{phase 1 "Requirements" "5 min"}}

{{say "Let me clarify the scope..."}}

{{diagram "rl-requirements"}}

{{think "Key decision: fail-open vs fail-closed. Fail-open is correct for most cases."
  (whyNot "fail-closed" "A full outage is worse than a few seconds of unlimited traffic.")
  (whatIf "Redis is down for 5 minutes" "Local in-memory counters take over with 80% accuracy.")
}}

{{hint "why sorted set?" "ZADD is O(log N) and ZRANGEBYSCORE counts in a window atomically."}}

{{phase 2 "Architecture" "10 min"}}

{{diagram "rl-architecture"}}

{{code "python" `class TokenBucket:
    def is_allowed(self, key, capacity, rate):
        ...`}}

{{compare "Database Choice" (options
  (best "Redis" "Sub-ms latency, atomic ops")
  (alt "DynamoDB" "Works but slower")
  (nofit "Postgres" "Too slow for per-request checks")
)}}

{{qa "Why Redis over DynamoDB?" "Redis gives O(1) atomic operations..."}}
```

---

## Diagram Registry (internal/diagrams/)

Diagrams are stored as **Go structs** in domain-specific files, not inline in content HTML. Content files reference diagrams by slug: `{{diagram "rl-architecture"}}`.

### Diagram Struct

```go
type Diagram struct {
    Slug        string  // Unique ID: "rl-architecture"
    Title       string  // Display title: "Architecture"
    Description string  // What this shows (for discovery)
    ContentFile string  // Which content uses it: "problems/rate-limiter"
    Type        Type    // TypeHTML or TypeImage
    HTML        string  // Raw HTML (for TypeHTML) — uses CSS diagram classes
    ImagePath   string  // Path in /static/img/diagrams/ (for TypeImage)
}
```

### Slug Prefixes

| Prefix | Domain | File | Count |
|--------|--------|------|-------|
| `rl-` | Rate Limiter | `rate_limiter.go` | 19 |
| `ig-` | Instagram | `instagram.go` | 24 |
| `url-` | URL Shortener | `url_shortener.go` | 21 |
| `algo-` | Algorithms | `algorithms.go` | 13 |
| `fund-` | Fundamentals | `fundamentals.go` | 15 |
| `pat-` | Patterns | `patterns.go` | 10 |

### Diagram CSS Classes

**Layout:** `.d-cols` (CSS grid), `.d-col`, `.d-flow` (horizontal flex), `.d-flow-v` (vertical flex), `.d-row`, `.d-branch` / `.d-branch-arm`

**Boxes:** `.d-box` + color: `.blue`, `.green`, `.purple`, `.amber`, `.red`, `.gray`, `.indigo`

**Grouping:** `.d-group` (dashed border), `.d-group-title`

**Arrows:** `.d-arrow` (horizontal), `.d-arrow-down` (vertical)

**Entity/DB:** `.d-entity`, `.d-entity-header [color]`, `.d-entity-body`, `.pk` (PK badge), `.fk` (FK badge), `.idx` + `.idx-btree`/`.idx-hash`/`.idx-gin`/etc.

**Specialized:** `.d-bitfield` / `.d-bitfield-segment`, `.d-ring` / `.d-ring-node`, `.d-subproblem [color]`

### Adding/Modifying a Diagram

1. Edit the appropriate Go file (e.g., `internal/diagrams/rate_limiter.go`)
2. Add or modify a `Register()` call with the diagram HTML
3. Reference in content: `{{diagram "rl-new-slug"}}`
4. Rebuild: `go run main.go`

---

## UI Architecture

### Layout (base.html)

Two-panel layout: **sidebar** (collapsible tree) + **detail area** (content view).

- HTMX partial swaps: clicking sidebar items sends `HX-Request: true`, server returns only the detail template (not full page). Direct URL access returns full page.
- Alpine.js: sidebar tree expand/collapse (client-side only, no server round-trip)
- Book mode: CSS multi-column layout for wide screens (>1600px), JS calculates column count
- Stage nav: sticky horizontal bar linking to `{{phase}}` anchors, auto-highlights via IntersectionObserver

### Template Dispatch

```
base.html
  ├── sidebar.html (always rendered)
  └── detail area (conditional):
      ├── PageType="problem"     → detail_problem.html
      ├── PageType="fundamental" → detail_fund.html
      ├── PageType="algorithm"   → detail_algo.html
      ├── PageType="pattern"     → detail_pattern.html
      └── PageType="welcome"     → welcome.html
```

### Cross-linking Flow

**Problem → Fundamental:** `detail_problem.html` renders `context_card.html` for each `problem.Uses[]` entry. Card shows config/why/not_this/risk/caveats. Links to `/fund/{slug}?from={problem-slug}`.

**Fundamental → Problem:** `detail_fund.html` renders `context_card.html` for each `fundamental.UsedBy[]` entry (auto-derived reverse links).

**Contextual highlighting:** When navigating to a fundamental from a problem (`?from=problem-slug`), the handler extracts keywords from the UsageLink's config field and passes them to the template. Client-side JS highlights matching table cells, compare options, and checklist items.

---

## Handler Logic (internal/handlers/handlers.go)

`renderContent(path)`: Reads `content/{path}/index.html` from embedded FS → parses as Go template with FuncMap (macros) → executes → returns `template.HTML`.

`isHTMX(r)`: Checks `HX-Request` header. If true, renders only the detail template. If false, renders full `base.html` with sidebar.

Template data always includes: `Problems`, `Fundamentals`, `Algorithms`, `Patterns` (for sidebar), plus page-specific data (`Problem`, `Fundamental`, etc.).

---

## PDF Generation

**Flow:** User clicks Generate → `POST /api/generate/{slug}` → Go handler collects problem content + all linked fundamental content paths → spawns `python3 engine/generate_pdf.py --config config.json` → Python uses Playwright (headless Chromium) to render combined HTML as A4 PDF.

**PDF includes:** Problem content + lettered appendix sections for each linked fundamental.

**Layout:** A4 portrait, 2cm left margin, Inter + JetBrains Mono fonts, page numbers bottom-right.

**Status:** PDF generation is currently placeholder (`not_implemented`). The Python engine and config format are ready.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Server | Go + chi router |
| Templates | Go `html/template` |
| Interactivity | HTMX (vendored, ~14KB) |
| Sidebar collapse | Alpine.js (vendored, ~15KB) |
| Styling | Single CSS file, CSS Grid + Flexbox, CSS variables |
| PDF | Python + Playwright (subprocess) |
| Content | HTML fragments in `content/` with Go template macros |
| Registry | YAML knowledge graph |

---

## Startup Sequence (main.go)

1. Load `content/_registry.yaml` → build `Registry` (problems, fundamentals with reverse links, algorithms, patterns)
2. Build `diagrams.Registry` via `diagrams.BuildDefault()` (registers all 102 diagrams)
3. Create `template.FuncMap` via `macros.FuncMap(diagramReg)` (macros get diagram registry for slug lookup)
4. Create `Handler` with registry, template FS, content FS, and FuncMap
5. Parse all layout templates with FuncMap
6. Set up chi routes and start server on :8080

---

## Gold Standard: Writing Content

### Problem File Structure (content/problems/{slug}/index.html)

Every problem file follows this exact structure with 12 phases:

```html
{{stageNav "Requirements & Problem Scope" 1 "Phase Title" 2 ... "Interview Deep-Dive Q&A" 12}}

{{phase 1 "Requirements & Problem Scope" "5 min"}}

{{say "Opening statement clarifying scope and assumptions..."}}

{{diagram "prefix-requirements"}}

{{think "Core design decision reasoning..."
  (whyNot "rejected alternative" "specific reason with numbers")
  (whatIf "failure scenario" "concrete response with fallback strategy")
}}

{{hint "why this choice?" "Detailed explanation of thought process — include numbers, latency, cost."}}

{{triggerQs "Questions this section might trigger"
  "Specific question an interviewer would ask?"
  "Direct answer. No filler. Start with the technical fact. Include numbers. 2-4 sentences."
  "Second question?"
  "Direct answer with specifics."
}}

{{phase 2 "Next Phase" "X min"}}
... (repeat pattern for phases 2-11)

{{phase 12 "Interview Deep-Dive Q&A" "10 min"}}

{{deepQA "Critical Interview Questions" `
<div class="dqa-item">
  <div class="dqa-q">Top-level question?</div>
  <div class="dqa-a">Direct answer with specifics.</div>
  <div class="dqa-sub">
    <div class="dqa-sub-q">Follow-up question?</div>
    <div class="dqa-sub-a">Answer with <span class="dqa-key">highlighted terms</span>.</div>
    <div class="dqa-deep">
      <div class="dqa-deep-q">Deep follow-up?</div>
      <div class="dqa-deep-a">Deepest answer.</div>
    </div>
  </div>
</div>
... (exactly 5 dqa-items, each with 3-level nesting)
`}}

{{key "One-sentence key takeaway for the entire problem."}}
```

**Required elements per problem file:**
- `{{stageNav}}` with exactly 12 phases
- Phase 12 is always "Interview Deep-Dive Q&A" with `{{deepQA}}` (5 items, 3-level nesting)
- Every phase opens with `{{say "..."}}` (what to literally say in interview)
- `{{hint}}` on every design decision (2-3 per phase minimum)
- `{{think}}` replaces `{{thought}}` — always include at least 1 `whyNot` or `whatIf` chain
- `{{triggerQs}}` on 5-6 key phases with 2-3 Q&A pairs each
- `{{diagram "slug"}}` for every visual (never inline HTML)
- `{{code "lang" "..."}}` for implementations
- `{{compare}}` for technology/algorithm selection decisions
- `{{table}}` for structured data (requirements, estimates, etc.)
- One `{{checklist}}` for section summaries
- One `{{key}}` as closing takeaway
- Q&A answers are ALWAYS direct — no "great question", no "this is the hardest", start with the answer

**Phase timing convention:**
- Requirements: 5 min
- Core algorithm/architecture phases: 8-10 min
- Supporting phases: 3-5 min
- Interview Q&A: 10 min

### Fundamental File Structure (content/fundamentals/{path}/index.html)

Fundamentals use 8 phases (NOT 12), no `{{stageNav}}`, no `{{deepQA}}`:

```html
{{phase 1 "What Is {Topic}?" "3 min"}}

{{say "Opening explanation..."}}

{{diagram "fund-topic-overview"}}

{{think "Key concept reasoning..."
  (whyNot "common misconception" "why it's wrong")
  (how "core mechanism" "how it actually works")
}}

{{hint "why this matters?" "Interview relevance explanation."}}

{{phase 2 "Core Concepts" "5 min"}}

{{table "Comparison" (rows
  (row "Feature" "Option A" "Option B")
  (row "Latency" "1ms" "10ms")
)}}

{{triggerQs "Interview Questions"
  "Question?" "Direct answer."
}}

... (8 phases total)

{{checklist "Summary point 1" "Summary point 2"}}

{{key "Key takeaway."}}
```

**Required elements per fundamental file:**
- 8 phases (not 12)
- `{{hint}}` 1-2 per phase
- `{{think}}` with chains (replacing all `{{thought}}`)
- `{{triggerQs}}` on 3-4 phases
- `{{qa}}` inline throughout (not deferred to end)
- `{{table}}` for structured comparisons (heavy use)
- `{{diagram "slug"}}` for visuals
- One `{{checklist}}` and one `{{key}}` at end

### Algorithm File Structure (content/algorithms/{slug}/index.html)

Same as fundamentals (8 phases, no stageNav, no deepQA) but with emphasis on:
- Full implementation in `{{code}}` (the centerpiece)
- Complexity analysis in `{{table}}`
- `{{compare}}` for algorithm alternatives
- `{{hint}}` on data structure and parameter choices

### Diagram Registration (internal/diagrams/{domain}.go)

```go
func registerNewDomain(r *Registry) {
    r.Register(&Diagram{
        Slug:        "prefix-diagram-name",     // prefix matches domain
        Title:       "Human Readable Title",
        Description: "What this diagram shows",
        ContentFile: "problems/slug",            // or "fundamentals/path"
        Type:        TypeHTML,
        HTML:        `<div class="d-flow">
            <div class="d-box blue">Component A</div>
            <div class="d-arrow">→</div>
            <div class="d-box green">Component B</div>
        </div>`,
    })
}
```

**Diagram slug conventions:**
| Prefix | Domain | File |
|--------|--------|------|
| `rl-` | Rate Limiter | `rate_limiter.go` |
| `ig-` | Instagram | `instagram.go` |
| `url-` | URL Shortener | `url_shortener.go` |
| `algo-` | Algorithms | `algorithms.go` |
| `fund-` | Fundamentals | `fundamentals.go` |
| `pat-` | Patterns | `patterns.go` |

**For a new problem**, create a new file `internal/diagrams/{domain}.go` with a `registerDomain(r *Registry)` function, and add it to `BuildDefault()` in `registry.go`.

### Registry Entry (_registry.yaml)

```yaml
# New problem
problems:
  - slug: chat-system
    title: "Chat System"
    description: "Design WhatsApp: real-time messaging, presence, group chats"
    path: problems/chat-system
    uses:
      - fundamental: storage/redis
        config: "Pub/Sub for real-time message delivery"
        why: "Sub-ms publish to connected WebSocket sessions"
        not_this: "Polling — adds 1-5 second latency to message delivery"
        risk: "Redis Pub/Sub is fire-and-forget — message loss on disconnect"
        caveats: "Pub/Sub does not persist. Use Kafka for message durability."

# New fundamental
fundamentals:
  - slug: messaging/kafka
    title: "Apache Kafka"
    description: "Distributed event streaming, partitions, consumer groups"
    path: fundamentals/messaging/kafka
```

### Writing Guidelines for Q&A

**DO:**
- "Redis sorted sets store members with float64 scores. ZADD is O(log N). For 1M members, that's ~20 comparisons per insert."
- "Fan-out-on-write for < 100K followers costs at most 100K Redis writes per post. At 1ms each = 100 seconds of background work."

**DON'T:**
- "This is the hardest DynamoDB challenge." (just explain why it's hard)
- "Great question!" (remove entirely)
- "Let me think about this..." (just give the answer)
- "There are several approaches..." (name the best one first)

## Common Tasks

### Add a new problem
1. Add entry to `content/_registry.yaml` under `problems:` with slug, title, description, path, and `uses:` links
2. Create `content/problems/{slug}/index.html` following the Gold Standard above (12 phases, stageNav, deepQA)
3. Create `internal/diagrams/{domain}.go` with diagram registrations, add to `BuildDefault()` in `registry.go`
4. Restart server

### Add a new fundamental
1. Add entry to `content/_registry.yaml` under `fundamentals:` (optionally with `children:`)
2. Create `content/fundamentals/{path}/index.html` following the Gold Standard (8 phases, no stageNav)
3. Add `uses:` entries in any problems that reference it
4. Restart server — reverse links are auto-derived

### Modify a diagram
1. Find the diagram slug in content: `{{diagram "rl-architecture"}}`
2. Find the Go file: `internal/diagrams/rate_limiter.go` (rl- prefix)
3. Edit the HTML in the `Register()` call
4. Restart server

### Add an image diagram
1. Place image in `web/static/img/diagrams/`
2. Register in appropriate Go file with `Type: TypeImage, ImagePath: "filename.png"`
3. Reference in content: `{{diagram "slug"}}`

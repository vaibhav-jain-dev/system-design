# System Design Interview Prep — Knowledge Dashboard

A web-based knowledge dashboard for system design interview prep. Two types of content deeply cross-linked: **Problems** (system design questions) and **Fundamentals** (infrastructure concepts), with ADHD-friendly visual formatting.

## Quick Start

```bash
# Run locally
go run main.go
# → http://localhost:8080

# Run with Docker
docker compose up --build
```

## Architecture

- **Go web server** (chi router + html/template + HTMX) — serves the dashboard
- **Python subprocess** (Playwright) — PDF export engine
- **Content as Go template fragments** — HTML files using helper macros
- **YAML registry** — defines the knowledge graph (problems → fundamentals links)
- **localStorage cache** — remembers reading position, scroll, sidebar state

## Project Structure

```
├── main.go                        # HTTP server entry point
├── internal/
│   ├── registry/registry.go       # YAML parsing, reverse index
│   ├── macros/macros.go           # Template helper functions (say, thought, qa, etc.)
│   └── handlers/handlers.go       # Route handlers
├── web/templates/                 # Go HTML templates
├── web/static/                    # CSS, HTMX, Alpine.js
├── content/
│   ├── _registry.yaml             # Knowledge graph definition
│   ├── problems/                  # Problem content (HTML fragments)
│   └── fundamentals/              # Fundamental content (deep hierarchy)
├── engine/generate_pdf.py         # PDF export engine
├── Dockerfile
└── docker-compose.yml
```

---

## Instructions for LLMs Working on This Project

### When Adding a New Problem

1. **Create the content file** at `content/problems/{slug}/index.html` using Go template macros
2. **Add to `_registry.yaml`** under `problems:` with all fields (slug, title, description, path, uses)
3. **For each fundamental in `uses:`**, check if it exists:
   - If **new fundamental**: Create its content file at `content/fundamentals/{category}/{slug}/index.html` with COMPLETE deep knowledge. Add it to `_registry.yaml` under `fundamentals:`.
   - If **existing fundamental**: Verify its content covers everything needed for this problem. If not, update the fundamental's content.
4. **Always include rich `uses:` metadata**: config, why, not_this, risk, caveats. Never just a slug link.

### When Adding a New Fundamental

1. **Create the content file** with complete deep knowledge using macros
2. **Add to `_registry.yaml`** under `fundamentals:` (with `children:` if it has sub-topics)
3. **Scan existing problems** — if any problem uses this fundamental but doesn't have it in `uses:`, add it

### Content Rules (CRITICAL)

1. **Diagrams over text** — Use HTML/SVG diagrams, NEVER ASCII art or Mermaid. Every architecture should be visual.
2. **ADHD-friendly** — Max 2-3 sentences per concept. Alternate visual components (say-box, table, thought-cloud, code-block). Never 3+ tables back-to-back.
3. **Mobile-friendly** — Content must work on iPhone 15 (393px). Book mode auto-activates on wide screens (>1600px).
4. **Info tooltips** — Every technical term should have `{{info "term" "definition"}}` on first use.
5. **Both perspectives** — Cover from building POV (how to implement) and interview POV (what to say, follow-up answers).
6. **Cost caveats** — Never say "serverless = free". Always include per-request pricing and monthly estimates.
7. **Failure scenarios** — Every component needs "what if this goes down?" in a thought-cloud.

### Contextual Highlighting

When a user navigates from a problem to a fundamental (via context card), items relevant to that problem are auto-highlighted. The mechanism:
- Context cards include `?from=problem-slug` in the URL
- The handler extracts keywords from the problem's `config` field for that fundamental
- Client-side JS highlights matching table rows, compare options, and checklist items
- Non-matching items are dimmed (50% opacity)
- "Reset Highlight" button removes the context

**For PDF export**: The highlight context should be preserved — when generating a PDF from a problem, the appendix fundamentals should have the same highlighting applied (relevant items marked).

### Available Content Macros

| Macro | Purpose |
|-------|---------|
| `{{say "..."}}` | What to literally say in interview (green box) |
| `{{thought "..."}}` | Side reasoning, math, failure scenarios (gray cloud) |
| `{{avoid "..."}}` | Common mistakes (red box) |
| `{{key "..."}}` | Key takeaway, one-liner (blue box) |
| `{{phase N "Title" "Time"}}` | Section header with number + time |
| `{{code "lang" "..."}}` | Code block |
| `{{qa "Q" "A" "FQ" "FA"}}` | Interviewer Q&A with optional follow-up |
| `{{followup "Q" "A"}}` | Follow-up question card (amber) |
| `{{checklist "..." "..."}}` | Green checklist with ✓ |
| `{{compare "Title" (options (best ...) (alt ...) (nofit ...))}}` | Comparison card |
| `{{table "Title" (rows (row ...) ...)}}` | Styled table |
| `{{info "term" "definition"}}` | Info tooltip (ℹ) for term definitions |
| `{{diagram "title"}}` | Diagram placeholder (to be expanded with SVG) |

### Registry YAML Format

```yaml
problems:
  - slug: my-problem
    title: "Problem Title"
    description: "One-line description"
    path: problems/my-problem
    uses:
      - fundamental: category/slug
        config: "Specific config/algorithm used"
        why: "Why this was selected"
        not_this: "Rejected alternative + reason"
        risk: "Key risk or failure mode"
        caveats: "Cost, limits, gotchas"

fundamentals:
  - slug: category/slug
    title: "Fundamental Title"
    description: "One-line description"
    path: fundamentals/category/slug
    children:  # optional sub-topics
      - slug: category/slug/child
        title: "Child Topic"
        path: fundamentals/category/slug/child
```

---

## Content Quality Standards

All content must meet these standards before merging. These ensure consistency, accuracy, and interview-readiness across the knowledge base.

### Structure Requirements

| Content Type | Phases | stageNav | deepQA (Phase 12) |
|-------------|--------|----------|--------------------|
| **Problem** | 12 | Required | Required (5 items, 3-level nesting) |
| **Fundamental** | 8 | Not used | Not used |
| **Algorithm** | 8 | Not used | Not used |

### Required Macros Per Content Type

**Problems (12 phases):**
- `{{stageNav}}` linking all 12 phases
- Every phase opens with `{{say "..."}}` (interview speech)
- `{{hint}}` on every design decision (2-3 per phase minimum)
- `{{think}}` with at least 1 `whyNot` or `whatIf` chain per usage
- `{{triggerQs}}` on 5-6 key phases with 2-3 Q&A pairs each
- `{{diagram "slug"}}` for every visual (never inline HTML diagrams)
- `{{compare}}` for technology/algorithm selection decisions
- `{{table}}` for structured data (requirements, estimates)
- `{{checklist}}` for section summaries
- `{{key}}` as closing takeaway
- Phase 12: `{{deepQA}}` with exactly 5 items, each with 3-level nesting

**Fundamentals / Algorithms (8 phases):**
- `{{hint}}` 1-2 per phase
- `{{think}}` with reasoning chains
- `{{triggerQs}}` on 3-4 phases
- `{{qa}}` inline throughout
- `{{table}}` for structured comparisons
- `{{diagram "slug"}}` for visuals
- `{{checklist}}` and `{{key}}` at end

### Language and Formatting Rules

1. **Bullet points over paragraphs** — any paragraph with 3+ connected ideas becomes a bullet list with sub-bullets
2. **Bold key terms** with `<strong>` on first use and for important numbers
3. **Highlight classes:**
   - `<span class="hl">` (amber) — key numbers, thresholds, capacities
   - `<span class="hl-blue">` — concepts, algorithm names
   - `<span class="hl-red">` — warnings, anti-patterns, failure modes
4. **`<small>` for caveats** — cost notes, side details
5. **No filler** — no "great question", no "let me think", no "there are several approaches". Start with the answer.
6. **No subjective ratings** — replace "fair", "excellent" with actual numbers
7. **Q&A answers: fact first** — open with the technical answer, then justify

### Factual Accuracy

- All cloud service metrics (PoP counts, pricing, limits) must be verified against current documentation
- Cost estimates must specify pricing model (on-demand vs reserved)
- Performance numbers must cite the source or measurement conditions
- When a claim cannot be verified, mark it with `<small>(approximate)</small>`

### Diagram Standards

- All diagrams registered in `internal/diagrams/` Go files — never inline HTML in content
- Slug prefix matches domain (`rl-`, `ig-`, `url-`, `algo-`, `fund-`, `pat-`)
- Each problem should have diagrams for: requirements, architecture, data model, and key algorithms
- Each fundamental should have at least 2 diagrams

---

### PDF Layout Rules

| Property | Value |
|----------|-------|
| Page size | A4 portrait |
| Margin left | 2 cm |
| Margin bottom | 0.5 cm |
| Margin top | 0 cm |
| Margin right | 0 cm |
| Page number | Bottom-right corner |
| Font | Inter (body), JetBrains Mono (code) |

PDF for a problem includes: problem content + all linked fundamentals as lettered appendix sections (auto-ordered by relevance).

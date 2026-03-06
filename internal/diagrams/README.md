# Diagram Library

Centralized registry of all visual diagrams used in the system design content.

## Architecture

```
internal/diagrams/
  registry.go       — Diagram struct, Registry type, BuildDefault()
  rate_limiter.go   — Rate Limiter problem diagrams (rl-*)
  instagram.go      — Instagram problem diagrams (ig-*)
  url_shortener.go  — URL Shortener problem diagrams (url-*)
  algorithms.go     — Algorithm diagrams (algo-*)
  fundamentals.go   — Fundamentals diagrams (fund-*)
  patterns.go       — AI/ML pattern diagrams (pat-*)
```

## Usage in Templates

### Slug-based lookup (recommended)
```html
{{diagram "rl-architecture"}}
```
Looks up the diagram by slug from the registry. If not found, renders a placeholder.

### Inline HTML (backward compatible)
```html
{{diagram "Title" `<div class="d-cols">...</div>`}}
```
Renders inline HTML directly. Still works for quick prototyping.

## Diagram Struct

```go
type Diagram struct {
    Slug        string // Unique ID: "rl-architecture", "ig-data-model"
    Title       string // Display title shown above diagram
    Description string // What this diagram shows (for discovery)
    ContentFile string // Which content uses it: "problems/rate-limiter"
    Type        Type   // TypeHTML or TypeImage
    HTML        string // Raw HTML (for TypeHTML)
    ImagePath   string // Path in /static/img/diagrams/ (for TypeImage)
}
```

## Slug Conventions

| Prefix | Content Area | Example |
|--------|-------------|---------|
| `rl-`  | Rate Limiter | `rl-architecture`, `rl-token-bucket` |
| `ig-`  | Instagram | `ig-data-model`, `ig-feed-strategy` |
| `url-` | URL Shortener | `url-api-design`, `url-base62` |
| `algo-`| Algorithms | `algo-consistent-hash-ring` |
| `fund-`| Fundamentals | `fund-cdn-request-flow` |
| `pat-` | AI/ML Patterns | `pat-rag-pipeline` |

## Adding a New Diagram

1. Add to the appropriate Go file (e.g., `rate_limiter.go`):
```go
r.Register(&Diagram{
    Slug:        "rl-new-diagram",
    Title:       "New Diagram Title",
    Description: "What this diagram shows",
    ContentFile: "problems/rate-limiter",
    Type:        TypeHTML,
    HTML:        `<div class="d-flow-v">...</div>`,
})
```

2. Reference in the content template:
```html
{{diagram "rl-new-diagram"}}
```

## Adding an Image Diagram

1. Place the image in `web/static/img/diagrams/`
2. Register it:
```go
r.Register(&Diagram{
    Slug:        "rl-sequence-diagram",
    Title:       "Rate Limit Sequence Diagram",
    Description: "Sequence diagram of rate limit check flow",
    ContentFile: "problems/rate-limiter",
    Type:        TypeImage,
    ImagePath:   "rl-sequence-diagram.png",
})
```

## Available CSS Classes for HTML Diagrams

### Layout
- `.d-cols` — CSS grid (auto-fit columns, stacks on mobile)
- `.d-col` — Column container
- `.d-flow` — Horizontal flex (stacks vertical on mobile)
- `.d-flow-v` — Vertical flex
- `.d-row` — Horizontal inline layout (wraps)
- `.d-branch` / `.d-branch-arm` — Branching flows

### Boxes
- `.d-box` — Base styled box
- Colors: `.d-box.blue`, `.green`, `.purple`, `.amber`, `.red`, `.gray`, `.indigo`

### Grouping
- `.d-group` — Dashed border group with background
- `.d-group-title` — Uppercase label for groups

### Arrows & Labels
- `.d-arrow` — Horizontal arrow
- `.d-arrow-down` — Vertical arrow
- `.d-label` — Small italic annotation

### Entity/Database Diagrams
- `.d-entity` — Table container
- `.d-entity-header [color]` — Table name header
- `.d-entity-body` — Fields section
- `.pk` — Primary key field (auto "PK" badge)
- `.fk` — Foreign key field (auto "FK" badge)
- `.idx` — Indexed field with type badges:
  `.idx-btree`, `.idx-hash`, `.idx-gin`, `.idx-gsi`, `.idx-unique`, `.idx-composite`

### Relationships
- `.d-er-lines` — Container for relationship lines
- `.d-er-connector` — Individual connector
- `.d-er-from`, `.d-er-to`, `.d-er-type` — Cardinality labels

### Specialized
- `.d-bitfield` / `.d-bitfield-segment` — Bit layout visualization
- `.d-ring` / `.d-ring-node` — Consistent hashing ring
- `.d-subproblem [color]` — Sub-problem cards with icon

## Diagram Index

### Rate Limiter (rl-*)
| Slug | Title | Description |
|------|-------|-------------|
| `rl-requirements` | Requirements | Functional & non-functional requirements overview |
| `rl-estimates` | Estimates | Back-of-envelope capacity estimates |
| `rl-headers` | Rate Limit Headers | HTTP response headers for 200 and 429 |
| `rl-rules` | Rate Limit Rules | Configuration by tier, endpoint, identity |
| `rl-algorithm-comparison` | Algorithm Comparison | Token Bucket vs Sliding Window vs Fixed Window |
| `rl-token-bucket` | Token Bucket | How Token Bucket algorithm works |
| `rl-sliding-window` | Sliding Window Counter | How Sliding Window works |
| `rl-fixed-window-burst` | Fixed Window Burst | Boundary burst problem visualization |
| `rl-architecture` | Architecture | Full request path with rate limiting |
| `rl-hop-by-hop` | Hop-by-Hop | Detailed request flow with annotations |
| `rl-distributed` | Distributed Challenges | Multi-server and sharding solutions |
| `rl-multi-region` | Multi-Region | Cross-region rate limiting strategies |
| `rl-data-model` | Data Model | Redis key schema for each algorithm |
| `rl-lua-flow` | Lua Script Flow | Atomic Redis execution diagram |
| `rl-elasticache` | ElastiCache Topology | 3-shard cluster layout |
| `rl-tradeoffs` | Trade-offs | Key decision matrix |
| `rl-edge-cases` | Edge Cases | DDoS, hot key, clock skew mitigations |
| `rl-cost` | Cost Breakdown | Infrastructure cost estimate |
| `rl-sub-problems` | Sub-Problems | Building blocks overview |

### Instagram (ig-*)
See `instagram.go` for full list (24 diagrams).

### URL Shortener (url-*)
See `url_shortener.go` for full list (21 diagrams).

### Algorithms (algo-*)
See `algorithms.go` for full list (13 diagrams).

### Fundamentals (fund-*)
See `fundamentals.go` for full list.

### Patterns (pat-*)
See `patterns.go` for full list (10 diagrams).

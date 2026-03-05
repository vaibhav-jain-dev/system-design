# Diagram Clarity Improvement Plan

## Context

All 102 diagrams across 6 domain files were audited for self-explanatory clarity. The goal: every diagram should be fully understandable without requiring external context. No box, label, or arrow should raise a question it doesn't answer. This plan catalogs every ambiguity found and specifies the fix.

---

## Common Issue Patterns

| Pattern | Count | Description |
|---------|-------|-------------|
| **Missing failure path** | ~18 | Happy path shown, no explanation of what happens on error/timeout/crash |
| **Vague action label** | ~14 | Labels like "process", "handle", "cached?", "Maybe" that don't explain what actually happens |
| **Black-box component** | ~12 | Component shown as box with no internals or explanation of mechanism |
| **Missing decision logic** | ~10 | Branch shown but no explanation of what triggers which path |
| **Incomplete flow** | ~8 | Flow starts but doesn't show end state or return path |
| **Undefined term/acronym** | ~6 | ARC, CRR, singleflight mentioned without definition |

---

## File: `internal/diagrams/rate_limiter.go` (19 diagrams)

### rl-requirements
- **Issue:** "Fail-open for most APIs, fail-closed for payments" stated but no explanation of consequences
- **Fix:** Add a sub-label: "Fail-open = allows traffic during Redis outage (brief over-limit). Fail-closed = rejects all traffic (zero revenue risk for payments)"

### rl-capacity-estimation
- **Issue:** "2-3 Redis ops per check" — which ops? Reader doesn't know if it's INCR, GET+SET, or EXPIRE
- **Fix:** Add label: "1 EVAL (Lua script) = internally: GET current count + INCR + EXPIRE = 2-3 logical ops"
- **Issue:** "50M keys x 100B = 5 GB" — 100B per key not broken down
- **Fix:** Add: "Key ~40B (rl:{user}:{endpoint}:{window}) + Value ~20B (counter+TTL) + Redis overhead ~40B = ~100B"

### rl-rules-config
- **Issue:** How system determines auth tier not shown (JWT? API key lookup?)
- **Fix:** Add box at top of flow: "API Gateway extracts tier from JWT claims (plan: free/pro/enterprise)"
- **Issue:** "POST /login: 5 req/15min + CAPTCHA" — unclear if CAPTCHA is part of rate limiter
- **Fix:** Clarify label: "Rate limiter returns 429 + CAPTCHA challenge URL (CAPTCHA is a separate service)"
- **Issue:** How three dimensions (tier, endpoint, identity) combine is missing
- **Fix:** Add label: "Evaluation order: identity key first → endpoint rule match → tier multiplier applied"

### rl-algorithm-comparison
- **Issue:** Sliding Window Counter has star but no explanation of why it's best
- **Fix:** Add sub-label: "Best overall: smooth accuracy + O(1) memory + O(1) time"
- **Issue:** Sliding Window Log "Memory: O(N) per key" — N definition unclear
- **Fix:** Clarify: "N = number of individual request timestamps stored in the window (e.g., 100 req/min = up to 100 entries)"

### rl-token-bucket
- **Issue:** "tokens >= 1?" — why not "tokens >= request_cost"?
- **Fix:** Change label to: "tokens >= cost? (cost=1 for most APIs, cost=N for batch endpoints)"
- **Issue:** Float tokens between 0-1 not addressed
- **Fix:** Add label: "Tokens are float64. Partial tokens accumulate via refill. Request needs >= 1.0 whole token."

### rl-sliding-window
- **Issue:** 0.75 weight not derived — reader can't reconstruct it
- **Fix:** Add calculation box: "Weight = overlap / window_size = 45min / 60min = 0.75 (we're 15min into current window, so 45min of previous window overlaps)"
- **Issue:** What happens at exact window boundary not shown
- **Fix:** Add note: "At minute 60: weight=0.0, previous window fully expired, only current window counts"

### rl-architecture
- **Issue:** "Check rules" — where are rules stored? Local cache or Redis?
- **Fix:** Change to: "Check rules (in-memory, loaded at startup, hot-reloaded via config change event)"
- **Issue:** "Extract key" is vague
- **Fix:** Change to: "Extract key: Authorization header → JWT → user_id, else API-Key header, else client IP"
- **Issue:** Flow ends at ALLOW/REJECT without showing response headers
- **Fix:** Add final box: "Attach X-RateLimit-Limit/Remaining/Reset headers to response (both ALLOW and REJECT)"

### rl-hop-by-hop
- **Issue:** JWT → API-Key → IP fallback priority unclear
- **Fix:** Rewrite label: "Priority chain: (1) JWT user_id if valid → (2) API-Key header if present → (3) client IP as last resort. Invalid JWT = reject 401, don't fall through"
- **Issue:** "3 separate commands" not specified
- **Fix:** Change to: "Single EVAL replaces 3 commands: GET counter + INCR counter + EXPIRE TTL"

### rl-distributed-challenges
- **Issue:** Hash tag `{user:123}` purpose not explained
- **Fix:** Add label: "Hash tags {user:123} force Redis Cluster to route all keys for same user to same shard (CRC16 hashes only the tag content)"
- **Issue:** What happens if target shard is down not shown
- **Fix:** Add amber box: "Shard down → fail-open: allow request, increment local counter, sync when shard recovers"

### rl-multi-region
- **Issue:** "Medium accuracy" undefined in percentage terms
- **Fix:** Change to: "~95% accuracy (local counters may miss cross-region requests within sync interval)"
- **Issue:** No decision guidance on which approach to pick
- **Fix:** Add a d-label at bottom: "Decision: Use local counters + gossip sync (recommended). Global Redis only if exact global limit is required (e.g., billing caps)"

### rl-data-model
- **Issue:** Token refill calculation not shown
- **Fix:** Add formula box: "current_tokens = min(capacity, stored_tokens + (now - last_refill) * refill_rate)"
- **Issue:** window_ts not defined
- **Fix:** Change label to: "window_ts = floor(now_unix / window_seconds) — e.g., for 60s window at 12:15:37 → 1234567860"
- **Issue:** TTL 2x not explained
- **Fix:** Add: "TTL = 2 x window_sec to keep previous window for weighted calculation (sliding window needs both current + previous)"
- **Issue:** Rules → Token Bucket 1:N relationship stated in text but not shown visually
- **Fix:** Add a d-arrow between Rules entity and Token Bucket key showing "1 rule → N keys (one per user)"

### rl-lua-execution
- **Issue:** Which algorithm this Lua script implements is unspecified
- **Fix:** Add title annotation: "Token Bucket Lua Script (for sliding window, see separate script)"
- **Issue:** "Return {allowed, remaining}" — remaining of what?
- **Fix:** Change to: "Return {allowed: bool, remaining_tokens: int, retry_after_ms: int}"

### rl-elasticache-topology
- **Issue:** Failover behavior not explained
- **Fix:** Add label: "Primary fails → ElastiCache auto-promotes replica within 15-30s. Replica serves reads during failover. Writes queue or fail-open."
- **Issue:** 39 GB memory calculation not shown
- **Fix:** Add: "3 shards x 2 nodes x 6.5 GB (r6g.large) = 39 GB total. Usable ~26 GB (3 primaries x 6.5 GB, replicas for failover)"

### rl-tradeoffs
- **Issue:** Redis "down" duration not specified — different strategies for different durations
- **Fix:** Add: "Redis down <30s: fail-open, local counters absorb. Down >30s: alert, consider circuit breaker to shed load"
- **Issue:** Cost of centralized vs sidecar not compared
- **Fix:** Add sub-labels: "Centralized: 1 service to manage, single point of failure. Sidecar: per-pod overhead (~5MB RAM each), no single point of failure"

### rl-edge-cases
- **Issue:** DDoS — how WAF + geo-blocking protect rate limiter not explained
- **Fix:** Add: "WAF blocks known attack patterns before traffic reaches rate limiter. Geo-blocking drops traffic from non-served regions. Rate limiter only handles legitimate-looking traffic."
- **Issue:** Hot Key — local cache consistency unclear
- **Fix:** Add: "Local cache (10ms TTL) = app server caches REJECT decisions for 10ms. Stale by at most 10ms. On cache miss → Redis check. Prevents one user from consuming all Redis ops."
- **Issue:** API Key Sharing — "anomaly detection" undefined
- **Fix:** Change to: "Per-key limits + anomaly detection (flag if >5 distinct IPs use same key in 1 minute → alert + temporary per-IP sub-limits)"
- **Issue:** Webhook Retry Storms — "separate tier" scope unclear
- **Fix:** Change to: "Webhook endpoints get a dedicated rate limit tier (e.g., 100 req/s vs 1000 req/s for user APIs). Rate limiter checks endpoint path to select tier."
- **Issue:** Redis Failover — how counters rebuild not explained
- **Fix:** Add: "Counters rebuild organically: each new request increments from 0. Within one window period, counters are accurate again. Brief over-limit window is accepted (fail-open)."

### rl-cost-breakdown
- **Issue:** Sidecar vs middleware distinction unclear
- **Fix:** Add clarification: "Rate limiting runs as middleware in your existing ECS tasks (not a separate sidecar container). Zero additional container cost."
- **Issue:** Multi-region cost not shown
- **Fix:** Add row: "Multi-region: x2-3 ElastiCache cost ($1100-$1650/mo) + cross-region data transfer (~$50/mo)"

---

## File: `internal/diagrams/instagram.go` (24 diagrams)

### ig-mvp-architecture
- **Issue:** "Pre-signed URL upload" — who generates the URL not specified
- **Fix:** Add: "API server generates pre-signed S3 URL (valid 15 min) → client uploads directly to S3"

### ig-upload-read-flow
- **Issue:** CDN URL generation at upload time unclear
- **Fix:** Add: "CDN URL = deterministic: https://cdn.example.com/media/{post_id}/{size}.jpg — constructed from post_id, not stored separately"

### ig-write-path
- **Issue:** "S3 503 → client retries" — no fallback if S3 permanently fails
- **Fix:** Add: "Client retries 3x with backoff. After 3 failures → return 503 to user with 'try again later'. Upload is idempotent via pre-signed URL."

### ig-read-path
- **Issue:** "Failure: slow query, timeout 5s" — what happens after timeout undefined
- **Fix:** Add: "Timeout 5s → return 504 to client. Feed service falls back to cached/stale feed if available, else trending posts."

### ig-cdn-media-pipeline
- **Issue:** "Accept header → WebP vs JPEG negotiation" — who negotiates unclear
- **Fix:** Change to: "CloudFront checks Accept header. If 'image/webp' present → serve WebP variant from S3. Else → serve JPEG. All variants pre-generated by Lambda."

### ig-hybrid-fanout
- **Issue:** Celebrity feed — pre-computed or queried live unclear
- **Fix:** Add label: "Celebrity posts queried at read time: SELECT FROM posts WHERE user_id IN (followed_celebrities) ORDER BY created_at LIMIT 20. Merged with pre-computed fan-out feed."

### ig-stage2-architecture
- **Issue:** S3 → Lambda trigger mechanism not explained
- **Fix:** Change to: "S3 Event Notification (ObjectCreated) → triggers Lambda → generates 4 image sizes (150px, 360px, 720px, 1080px)"

### ig-notification-flow
- **Issue:** When deduplication happens and batching strategy unclear
- **Fix:** Add: "Dedup window: 5 minutes. Batch similar events: 'User A and 3 others liked your post' instead of 4 separate notifications."

### ig-peak-traffic
- **Issue:** "Fallback to trending (pre-computed)" — computation freshness unknown
- **Fix:** Add: "Trending feed recomputed every 5 minutes by background job. Serves as fallback when personalized feed times out."

### ig-caching-strategy
- **Issue:** Cache stampede protection — which layer and interaction unclear
- **Fix:** Split into two labels: "Singleflight: applied at app server layer (Go sync.Once per cache key). Probabilistic early expiry: applied at Redis layer (each key expires at TTL x random(0.8-1.0) to stagger rebuilds)."

### ig-multi-region
- **Issue:** Consistency model and conflict resolution not mentioned
- **Fix:** Add: "DynamoDB Global Tables: last-writer-wins (LWW) conflict resolution using item timestamps. Eventually consistent reads across regions (~1s lag)."

### ig-stories-lifecycle
- **Issue:** 25h instead of 24h — unexplained
- **Fix:** Change label to: "S3 lifecycle deletes after 25h (1h buffer for timezone edge cases — user posts at 11:59 PM, story should last until 11:59 PM next day in their timezone)"

### ig-explore-search
- **Issue:** NSFW/spam filter — flagged content handling incomplete
- **Fix:** Add: "Flagged content hidden from Explore immediately. Human review SLA: 4 hours. If approved, content re-enters Explore. If rejected, removed permanently + user warned."

### ig-distributed-counter
- **Issue:** Shard index computation from post_id not shown
- **Fix:** Add: "Shard index = CRC32(post_id) % num_shards. Read: query all shards in parallel, sum results. Aggregate cached 30s."

### ig-scaling-timeline
- **Issue:** Who monitors triggers and how often not explained
- **Fix:** Add: "CloudWatch alarms monitor these metrics. Auto-scaling triggers within 2 minutes of threshold breach. Architecture changes (add region, shard DB) are manual decisions."

### ig-content-moderation
- **Issue:** Human review SLA and queue backup handling missing
- **Fix:** Add: "Human review SLA: 4 hours. Queue > 10K items → auto-reject scores 0.7-0.9, only escalate 0.5-0.7. Team: 50 moderators, ~200 reviews/hour each."

### ig-monitoring-slos
- **Issue:** "Celebrity merge" step unexplained
- **Fix:** Change to: "Celebrity merge (15ms): fan-out-on-read queries for celebrity posts, merged into pre-computed feed by timestamp"

### ig-analytics-pipeline
- **Issue:** Flink window parameters not specified
- **Fix:** Add: "Flink: 1-minute tumbling windows for real-time counts (likes, views). 5-minute sliding windows for trending score computation."

### ig-security-layers
- **Issue:** What happens to toxic content not explained
- **Fix:** Add: "Toxic caption → reject post creation with error message. Toxic comment → shadow-hide (author sees it, others don't) + queue for review."

---

## File: `internal/diagrams/url_shortener.go` (21 diagrams)

### url-bit-layout
- **Issue:** "seconds mod 65536" — wraparound after 18.2 hours not addressed
- **Fix:** Add: "Wraps every 18.2h — acceptable because uniqueness comes from the full ID (timestamp + machine + sequence), not timestamp alone"

### url-collision-resolution
- **Issue:** "Collision probability ~0.14% per attempt" — is 5 retries enough?
- **Fix:** Add: "5 retries: probability of all 5 colliding = (0.0014)^5 ≈ 0 (5.4 x 10^-15). Effectively guaranteed to succeed."

### url-architecture
- **Issue:** "async analytics" — when/how fires ambiguous
- **Fix:** Change to: "At redirect: fire-and-forget analytics event to SQS (non-blocking, <1ms overhead). Worker processes events for click counts."

### url-write-read-paths
- **Issue:** Safe Browsing API provider and failure handling missing
- **Fix:** Add: "Google Safe Browsing API v4 (~50ms). On API timeout → allow URL creation (fail-open) + queue for async re-check within 5 minutes."

### url-scaling-strategy
- **Issue:** Multi-AZ redundancy headroom not accounted for
- **Fix:** Add: "23 tasks x 1.5 (50% headroom for AZ failure) = 35 tasks across 3 AZs. If 1 AZ fails, remaining 2 AZs handle peak load."

### url-multi-az
- **Issue:** Failover behavior for ALB and Redis not explained
- **Fix:** Add: "ALB auto-removes unhealthy AZ targets in ~30s. Redis: ElastiCache auto-promotes replica in different AZ within 15-30s. ECS replaces failed tasks in ~60s."

### url-reliability-patterns
- **Issue:** Retry backoff parameters not specified
- **Fix:** Change to: "Retry with exponential backoff (base=100ms, max=5s, jitter=±50ms). Max 3 retries. Circuit breaker opens after 5 consecutive failures."

### url-edge-cases
- **Issue:** "Same long URL → two different short codes (by design)" — why by design not explained
- **Fix:** Add: "By design: each short URL tracks its own analytics (click count, referrer). Deduplication would merge analytics across different campaigns sharing the same destination."

### url-multi-region
- **Issue:** DynamoDB Global Tables consistency model not explained
- **Fix:** Add: "Eventually consistent across regions (<1s). Reads within same region are strongly consistent. Conflict resolution: last-writer-wins by timestamp."

### url-security-layers
- **Issue:** CAPTCHA trigger and failure mode incomplete
- **Fix:** Add: "CAPTCHA triggered by rate limiter (>5 creates from anonymous IP). Served by reCAPTCHA v3 (invisible). Failed CAPTCHA → 403 with retry link."

### url-bot-detection
- **Issue:** "Chrome UA + curl TLS = bot" — unclear if this is one example or the trigger
- **Fix:** Add: "One of multiple signals. Bot score = weighted sum of: TLS fingerprint mismatch (0.4), behavioral anomaly (0.3), IP reputation (0.2), header analysis (0.1). Score > 0.7 → block."

### url-monitoring-slos
- **Issue:** Tracing "breakdown" not explained
- **Fix:** Add: "X-Ray shows per-service latency waterfall: which component is slow? ALB routing (1ms) → ECS handler (2ms) → Redis lookup (0.5ms) → DynamoDB fallback (5ms)"

### url-db-access-patterns
- **Issue:** "Adaptive capacity redistributes within minutes" — automatic or manual unclear
- **Fix:** Change to: "DynamoDB adaptive capacity is fully automatic. Detects hot partitions and redistributes throughput within 5-30 minutes. No manual intervention needed."

### url-base62-encoding
- Minor: Decimal place values shown but pedagogical purpose unclear
- **Fix:** Add label: "Base62 digits have positional value like decimal — rightmost = 62^0, next = 62^1, etc."

### url-qr-code
- **Issue:** Lambda@Edge QR generation is a black box
- **Fix:** Add: "Lambda@Edge uses `qrcode` library. Input: short URL. Output: PNG (200x200, ~2KB). Cached at edge for 24h via Cache-Control header."

### url-alias-validation
- **Issue:** Homoglyph check scope unclear
- **Fix:** Add: "Unicode confusables list (ICU). Normalizes aliases to ASCII skeleton, checks against reserved brand names (paypal, google, etc.)"

### url-analytics-schema
- **Issue:** IP → Geo lookup — real-time or cached unclear
- **Fix:** Change to: "MaxMind GeoIP2 local database (updated weekly). In-memory lookup, <0.1ms. No external API call."

---

## File: `internal/diagrams/algorithms.go` (13 diagrams)

### algo-base62-url-shortener
- **Issue:** "CDN (cached?)" and "Redis (cached?)" — question marks suggest uncertainty
- **Fix:** Remove question marks. Change to: "CDN (302 not cached)" and "Redis (cache-aside, TTL 24h)"
- **Issue:** Read path fallback condition not explained
- **Fix:** Add: "Redis miss → DynamoDB lookup → write result back to Redis (cache-aside pattern)"

### algo-bloom-filter-cache-layer
- **Issue:** "Maybe here" is vague
- **Fix:** Change to: "Possibly exists (check DB to confirm)"
- **Issue:** Cache negative results not addressed
- **Fix:** Add: "DB confirms not found → cache null result with short TTL (60s) to prevent repeated misses"

### algo-consistent-hash-ring
- **Issue:** Hash function not specified
- **Fix:** Add label: "Hash function: CRC32 or MD5 (consistent across all nodes). Maps keys and nodes to same 0-2^32 ring."

### algo-consistent-hash-vnodes
- **Issue:** How virtual nodes are assigned not explained
- **Fix:** Add: "Each physical node gets 100-200 virtual nodes at deterministic positions: hash(node_id + '-' + vnode_index). Uniform distribution guaranteed by hash function."

### algo-consistent-hash-rebalance
- **Issue:** How system knows which keys to move not explained
- **Fix:** Add: "New node D inserted between A and B on ring. Keys in range (A, D] now belong to D instead of B. Only keys in this range are migrated — other nodes unaffected."

### algo-snowflake-btree-insert
- **Issue:** "???" in UUID column is confusing
- **Fix:** Change to: "random page" with label: "UUID = random → inserts scatter across B-tree → constant page splits and cache misses"

### algo-snowflake-distributed-arch
- **Issue:** How machine_id prevents collisions not stated
- **Fix:** Add: "Each server assigned unique machine_id (0-1023) at deployment. Same timestamp + same sequence can occur on different servers, but machine_id makes the full ID unique."

### algo-token-vs-leaky-bucket
- **Issue:** Why outputs differ not explained with concrete numbers
- **Fix:** Add concrete example: "Both receive 10 requests at once. Token bucket (capacity=5): 5 pass immediately, 5 rejected. Leaky bucket (rate=1/s): 1 passes per second, 9 queued for up to 9 seconds."

### algo-token-bucket-api-arch
- **Issue:** "atomic HMGET + HMSET" — what Lua script checks not explained
- **Fix:** Add: "Lua script: (1) HMGET tokens, last_refill → (2) calculate refilled tokens → (3) if tokens >= 1: HMSET tokens-1, update last_refill → return ALLOW. Else return REJECT."
- **Issue:** Response headers returned on ALLOW only or both unclear
- **Fix:** Add: "Headers returned on BOTH allow (200) and reject (429). 429 adds Retry-After header."

---

## File: `internal/diagrams/fundamentals.go` (15 diagrams)

### fund-cloudfront-event-hooks
- **Issue:** "Viewer Request" vs "Origin Request" — same request or different concepts unclear
- **Fix:** Add: "Same HTTP request at different stages: Viewer Request = at edge (before cache check). Origin Request = at origin (only on cache miss)."
- **Issue:** CF Functions vs Lambda@Edge — when to use which not explained
- **Fix:** Add: "CF Functions: <1ms, JS only, simple transforms (headers, redirects). Lambda@Edge: <5s, any runtime, complex logic (auth, A/B testing, image resize)."

### fund-cloudfront-origin-shield
- **Issue:** "Reduces origin load by 90%" — context-dependent claim
- **Fix:** Add: "90% reduction assumes typical web traffic (80% cacheable content, 24h TTL). Benefits decrease with low-TTL or uncacheable content."

### fund-lb-multi-az
- **Issue:** Route 53 → ALB health check direction unclear
- **Fix:** Add: "Route 53 health checks ALB endpoint directly (HTTP 200 on /health). If ALB fails health check → Route 53 stops routing DNS to that region."
- **Issue:** "zonal shift via ARC" — ARC undefined
- **Fix:** Change to: "Zonal shift via ARC (AWS Application Recovery Controller) — manually or automatically shift traffic away from impaired AZ"

### fund-alb-path-routing
- **Issue:** WAF rule evaluation location unclear
- **Fix:** Add: "WAF is attached to ALB. Rules evaluated at ALB before listener rules. Block → 403. Allow → proceed to path-based routing."

### fund-alb-security-layers
- **Issue:** "SQL injection, XSS, rate limit" — separate types or examples unclear
- **Fix:** Change to: "WAF Rule Groups: (1) SQL injection detection, (2) XSS detection, (3) rate-based rules. Each evaluated independently, first match wins."
- **Issue:** AWS Shield position in stack unclear
- **Fix:** Add: "AWS Shield operates at L3/L4 (network layer) — absorbs DDoS volumetric attacks before they reach ALB. Transparent, always-on, no ALB configuration needed."

### fund-nlb-architecture
- **Issue:** NLB routing algorithm not specified
- **Fix:** Add: "NLB uses flow hash algorithm: hash(src_ip, src_port, dst_ip, dst_port, protocol) → deterministic target selection. Same client connection always reaches same target."

### fund-nlb-privatelink
- **Issue:** PrivateLink mechanism is a black box
- **Fix:** Add: "PrivateLink creates ENI (Elastic Network Interface) with private IP in consumer VPC. Traffic stays on AWS backbone, never traverses public internet. TCP/UDP only."

### fund-dynamodb-key-structure
- **Issue:** SK ordering basis not explained
- **Fix:** Add: "SK sorts lexicographically. Prefix convention (ORDER#, PROFILE#) groups item types. Within a PK, query SK with BEGINS_WITH('ORDER#') returns all orders sorted by date."

### fund-dynamodb-single-table
- **Issue:** How to query "all orders for user" not shown
- **Fix:** Add: "Query: PK = 'USER#U123' AND SK BEGINS_WITH('ORDER#') → returns all orders sorted by date. This is the core single-table access pattern."

### fund-dynamodb-write-sharding
- **Issue:** Read path for sharded writes not shown — how to read back?
- **Fix:** Add: "Write: random shard suffix (0-9). Read: parallel QUERY on all 10 shard keys → aggregate results. Trade-off: 10x read queries for even write distribution."

### fund-redis-cache-aside
- **Issue:** TTL and null caching not addressed
- **Fix:** Add: "SET with TTL (e.g., 300s). Cache null/empty results with short TTL (60s) to prevent cache stampede on missing keys."

### fund-redis-cluster
- **Issue:** Key-to-slot mapping not explained
- **Fix:** Add: "Slot = CRC16(key) mod 16384. Adding a shard → migrate slot ranges (e.g., slots 0-5460 split to give new shard 4096-5460). Redis Cluster handles migration automatically."

### fund-redis-failover
- **Issue:** What happens to writes during 10-30s detection window not addressed
- **Fix:** Add: "During detection: writes to failed primary are lost (Redis async replication). App should retry writes after failover completes. Use WAIT command for critical writes requiring sync replication."
- **Issue:** "Clients Reconnect Automatically" — depends on client library
- **Fix:** Change to: "Clients reconnect automatically (Lettuce/Jedis/ioredis support auto-reconnect). ElastiCache DNS endpoint resolves to new primary."

---

## File: `internal/diagrams/patterns.go` (10 diagrams)

### pat-react-agent-loop
- **Issue:** Loop termination logic not explained
- **Fix:** Add: "Loop exits when: (1) LLM outputs final_answer tool call, or (2) max iterations reached (typically 10), or (3) LLM explicitly states task is complete."
- **Issue:** Tool failure handling not shown
- **Fix:** Add amber box on OBSERVE: "Tool error → error message passed back to LLM as observation. LLM decides: retry with different params, try different tool, or report failure."

### pat-agent-architecture-patterns
- **Issue:** When to use which architecture not explained
- **Fix:** Add labels: "Simple (1 LLM): <3 tools, single domain. Router: multiple domains, need specialization. Multi-Agent: complex workflows, need verification between steps."
- **Issue:** Router decision mechanism not shown
- **Fix:** Add: "Router = LLM call with system prompt listing sub-agents and their capabilities. Routes based on intent classification."

### pat-embedding-space
- **Issue:** "Meaning axis 1/2" too abstract
- **Fix:** Change to concrete example: "Semantic dimension (e.g., topic)" and "Semantic dimension (e.g., sentiment)" with note: "Actual dimensions are learned, not human-interpretable. 2D shown for illustration; real embeddings use 768-3072 dimensions."

### pat-hnsw-index-structure
- **Issue:** Layer assignment rule not explained
- **Fix:** Add: "Each node assigned to layers probabilistically: P(layer L) = 1/M^L. All nodes in Layer 0. ~1/M in Layer 1. ~1/M^2 in Layer 2. Search starts at top, drills down."

### pat-embedding-pipeline
- **Issue:** OpenAI vs self-host tradeoff not guided
- **Fix:** Add: "OpenAI: simpler, $0.0001/1K tokens, requires internet. Self-host: higher upfront cost, no data leaves your infra, lower latency (<10ms vs ~100ms)."
- **Issue:** "upsert" not defined
- **Fix:** Change to: "pgvector UPSERT (INSERT ... ON CONFLICT UPDATE) — re-index if document content changed"

### pat-guardrail-pipeline
- **Issue:** Error vs "ask to rephrase" decision logic missing
- **Fix:** Add: "PII detected → ask to rephrase (recoverable). Injection detected → return error (block). Topic filter → return error (policy violation)."
- **Issue:** Output guardrail failure actions undefined
- **Fix:** Add: "Hallucination check fail → retry with stricter prompt (max 2 retries). Toxicity fail → return fallback response. Format fail → retry with format instructions."

### pat-prompt-chain-pipeline
- **Issue:** Whether steps are LLM calls or post-processing unclear
- **Fix:** Add: "Each step = separate LLM call with specialized prompt. Output of step N becomes input context for step N+1."
- **Issue:** Why only step 3 retries not explained
- **Fix:** Add: "Steps 1-2 are extraction/classification (deterministic enough). Step 3 is generation (creative, may need retry). Step 4 validates — if validation fails, retry step 3."

### pat-gate-pattern-flow
- **Issue:** "Gate Validator" — LLM or regex?
- **Fix:** Change to: "Gate Validator (LLM call with rubric: 'Does this output meet criteria X, Y, Z? Answer PASS or FAIL with reason')"
- **Issue:** Abort vs fallback decision not guided
- **Fix:** Add: "Max retries exhausted → abort if quality is critical (legal, medical). Use fallback (generic response) if partial answer is acceptable."

### pat-rag-pipeline
- **Issue:** When reranking is needed vs optional not explained
- **Fix:** Add: "Rerank when: retrieval returns >20 candidates, or precision matters more than latency (+50-100ms). Skip when: <10 candidates or latency budget <200ms."
- **Issue:** Chunk count and formatting not specified
- **Fix:** Add: "Top-k chunks (k=3-5 typical). Format: each chunk wrapped with source metadata. Total context must fit model's window."

### pat-two-stage-retrieval-reranking
- **Issue:** "dense + BM25" — merged or separate unclear
- **Fix:** Change to: "Hybrid retrieval: run dense (vector) and BM25 (keyword) in parallel → merge results using Reciprocal Rank Fusion (RRF). Take union of top-50 from each."
- **Issue:** Cross-encoder cost not addressed
- **Fix:** Add: "Cross-encoder scores each (query, doc) pair independently. 50 candidates x ~2ms each = ~100ms. More accurate than bi-encoder but O(N) per query."

---

## Implementation Approach

### Execution order
1. **High-severity first:** rl-rules-config, rl-architecture, rl-data-model, rl-lua-execution, rl-edge-cases, rl-multi-region (these are diagrams interviewers would probe most)
2. **Then medium-severity** across all files
3. **Then low-severity** polish

### How to apply fixes
Each fix above modifies the `HTML` field of a `Diagram` struct in the corresponding Go file. Changes are:
- Adding `<div class="d-label">...</div>` sub-labels to existing boxes
- Changing existing label text to be more precise
- Adding new boxes (typically `d-box amber` or `d-box gray`) for failure paths
- Adding `<div class="d-note">` annotations where supported

### Files to modify
- `internal/diagrams/rate_limiter.go` — 18 diagrams with fixes
- `internal/diagrams/instagram.go` — 18 diagrams with fixes
- `internal/diagrams/url_shortener.go` — 17 diagrams with fixes
- `internal/diagrams/algorithms.go` — 9 diagrams with fixes
- `internal/diagrams/fundamentals.go` — 14 diagrams with fixes
- `internal/diagrams/patterns.go` — 10 diagrams with fixes

### Verification
1. `go build ./...` — ensure all Go files compile
2. `go run main.go` — start server
3. Visit each diagram in browser and verify: (a) no rendering errors, (b) added labels are visible and readable, (c) diagram doesn't overflow its container
4. Spot-check 5-10 diagrams across different files for clarity improvement

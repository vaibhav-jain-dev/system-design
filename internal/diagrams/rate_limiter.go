package diagrams

func registerRateLimiter(r *Registry) {
	r.Register(&Diagram{
		Slug:        "rl-requirements",
		Title:       "Functional & Non-Functional Requirements",
		Description: "Prioritized functional requirements and non-functional targets for a distributed rate limiter",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P0 &#8212; Core (Must Have)</div>
      <div class="d-flow-v">
        <div class="d-box green">Rate limit by user/API key</div>
        <div class="d-box green">Rate limit by IP (anonymous)</div>
        <div class="d-box green">Configurable rules per endpoint/tier</div>
        <div class="d-box green">Return X-RateLimit-* headers</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P1 &#8212; Important</div>
      <div class="d-flow-v">
        <div class="d-box blue">Distributed counting across servers</div>
        <div class="d-box blue">Fail-open when Redis is down</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P2 &#8212; Nice to Have</div>
      <div class="d-flow-v">
        <div class="d-box gray">Multi-region rate limiting</div>
        <div class="d-box gray">Admin dashboard for rules</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">Latency: &lt; 1ms overhead per request</div>
        <div class="d-box purple">Availability: 99.999% (in hot path)</div>
        <div class="d-box purple">Throughput: 1M+ decisions/second</div>
        <div class="d-box amber">Accuracy: Allow 5% over-limit tolerance</div>
        <div class="d-box amber">Fail-open: If limiter fails, allow traffic</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Key Decisions</div>
      <div class="d-flow-v">
        <div class="d-box red">Fail-open vs Fail-closed?</div>
        <div class="d-label">Fail-open for most APIs, fail-closed for payments</div>
        <div class="d-box red">Exact vs Approximate counting?</div>
        <div class="d-label">Lua atomic scripts for single-region accuracy</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-capacity-estimation",
		Title:       "Back-of-Envelope Estimates",
		Description: "Traffic, storage, and Redis capacity estimates for rate limiting at scale",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Traffic</div>
      <div class="d-flow-v">
        <div class="d-box blue">100K RPS across all services</div>
        <div class="d-box blue">Peak (5x) = 500K checks/sec</div>
        <div class="d-box purple">2-3 Redis ops per check</div>
        <div class="d-box purple">= 300K&#8211;1.5M Redis ops/sec peak</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Storage</div>
      <div class="d-flow-v">
        <div class="d-box amber">~100 bytes per counter</div>
        <div class="d-box amber">10M active users &#215; 5 endpoints</div>
        <div class="d-box amber">= 50M keys &#215; 100B = 5 GB</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Redis Capacity</div>
      <div class="d-flow-v">
        <div class="d-box green">100K+ ops/sec per node</div>
        <div class="d-box green">3-node cluster = 300K ops/sec</div>
        <div class="d-box green">6 nodes for 500K peak</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-headers",
		Title:       "Rate Limit Headers (Every Response)",
		Description: "HTTP response headers for allowed (200) and rejected (429) rate-limited requests",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">200 OK (Allowed)</div>
        <div class="d-flow-v">
          <div class="d-box green" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">HTTP/1.1 200 OK
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 67
X-RateLimit-Reset: 1704067260
X-RateLimit-Policy: "100;w=60"</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">429 Too Many Requests</div>
        <div class="d-flow-v">
          <div class="d-box red" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1704067260
Retry-After: 45</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-rules-config",
		Title:       "Rate Limit Rules Configuration",
		Description: "Rate limit rules organized by auth tier, endpoint type, and identity",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">By Auth Tier</div>
      <div class="d-flow-v">
        <div class="d-box gray">Free: 100 req/min</div>
        <div class="d-box blue">Pro: 1,000 req/min</div>
        <div class="d-box green">Enterprise: 10,000 req/min</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">By Endpoint Type</div>
      <div class="d-flow-v">
        <div class="d-box green">GET /api/*: 500 req/min</div>
        <div class="d-box amber">POST /api/*: 100 req/min</div>
        <div class="d-box red">POST /login: 5 req/15min + CAPTCHA</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">By Identity</div>
      <div class="d-flow-v">
        <div class="d-box purple">By API Key: per-key limits</div>
        <div class="d-box amber">By IP: 50 req/min (anonymous)</div>
        <div class="d-box red">By IP + Path: login brute-force</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-algorithm-comparison",
		Title:       "Algorithm Comparison Overview",
		Description: "Side-by-side comparison of four rate limiting algorithms with trade-offs",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Token Bucket</div>
      <div class="d-flow-v">
        <div class="d-box green">Memory: O(1) per key</div>
        <div class="d-box green">Allows controlled bursts</div>
        <div class="d-box blue">Used by: AWS API Gateway</div>
        <div class="d-label">Best for: APIs that tolerate bursts</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Sliding Window Counter &#9733;</div>
      <div class="d-flow-v">
        <div class="d-box green">Memory: O(1) per key</div>
        <div class="d-box green">No boundary burst problem</div>
        <div class="d-box blue">Used by: Cloudflare, Stripe</div>
        <div class="d-label">Best for: Most use cases (recommended)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fixed Window Counter</div>
      <div class="d-flow-v">
        <div class="d-box amber">Memory: O(1) per key</div>
        <div class="d-box red">2x burst at boundaries!</div>
        <div class="d-box gray">Simplest to implement</div>
        <div class="d-label">Best for: MVP only</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Sliding Window Log</div>
      <div class="d-flow-v">
        <div class="d-box green">Highest accuracy</div>
        <div class="d-box red">Memory: O(N) per key!</div>
        <div class="d-box gray">N = requests in window</div>
        <div class="d-label">Best for: Low-rate limits (login)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-token-bucket",
		Title:       "Token Bucket — How It Works",
		Description: "Visual flow of the token bucket algorithm showing refill, check, and allow/reject logic",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box indigo" style="min-width:140px;text-align:center">
      <strong>Bucket</strong><br>
      capacity: 10<br>
      tokens: 7
    </div>
    <div class="d-arrow">&#8592;</div>
    <div class="d-box green" style="min-width:120px;text-align:center">
      <strong>Refill</strong><br>
      +2 tokens/sec
    </div>
  </div>
  <div class="d-arrow-down">&#8595; request arrives</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box green">tokens &#8805; 1?</div>
        <div class="d-arrow-down">&#8595; YES</div>
        <div class="d-box green">tokens-- &#8594; ALLOW</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box red">tokens &lt; 1?</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">REJECT (429)</div>
      </div>
    </div>
  </div>
  <div class="d-label">Burst: can send up to 'capacity' requests instantly, then rate-limited to refill_rate/sec</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-sliding-window",
		Title:       "Sliding Window Counter — How It Works",
		Description: "Weighted count calculation across previous and current windows to prevent boundary bursts",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Previous Window</div>
      <div class="d-box amber" style="text-align:center">
        count = 84<br>
        <small>11:00 &#8212; 12:00</small>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Current Window</div>
      <div class="d-box blue" style="text-align:center">
        count = 36<br>
        <small>12:00 &#8212; 13:00</small>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" style="text-align:center">
    <strong>Weighted Count</strong><br>
    now = 12:15 &#8594; 75% of prev window overlaps<br>
    weighted = 84 &#215; 0.75 + 36 = <strong>99</strong><br>
    limit = 100 &#8594; <strong>ALLOW</strong> (1 remaining)
  </div>
  <div class="d-label">Eliminates boundary burst: smoothly transitions between windows using time-weighted average</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-fixed-window-burst",
		Title:       "Fixed Window — The Boundary Burst Problem",
		Description: "Demonstrates how fixed window counters allow 2x burst at window boundaries",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Window 1 (11:00&#8211;12:00)</div>
      <div class="d-flow-v">
        <div class="d-box gray" style="text-align:center">
          &#8230; quiet &#8230;<br>
          <small>0 requests until 11:59</small>
        </div>
        <div class="d-box red" style="text-align:center">
          <strong>11:59:30 &#8594; 100 requests!</strong><br>
          <small>All at end of window</small>
        </div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Window 2 (12:00&#8211;13:00)</div>
      <div class="d-flow-v">
        <div class="d-box red" style="text-align:center">
          <strong>12:00:01 &#8594; 100 requests!</strong><br>
          <small>Counter reset to 0</small>
        </div>
        <div class="d-box gray" style="text-align:center">
          &#8230; rest of window &#8230;
        </div>
      </div>
    </div>
  </div>
  <div class="d-box red" style="text-align:center">
    <strong>Result: 200 requests in ~30 seconds!</strong><br>
    Limit was 100/min but user got 2x at the boundary
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-architecture",
		Title:       "Rate Limiter Architecture — Full Request Path",
		Description: "End-to-end request flow from client through ALB, API server, rate limit middleware, and Redis",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>Client Request</strong></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" style="text-align:center">Route 53 (DNS)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" style="text-align:center">ALB (Load Balancer)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" style="text-align:center"><strong>API Server (ECS)</strong></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" style="text-align:center;border:2px solid var(--amber)">
    <strong>Rate Limit Middleware</strong><br>
    Extract key (API key or IP) &#8594; Check rules &#8594; Query Redis
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box red" style="text-align:center">
          <strong>Redis</strong><br>
          Lua script (atomic)<br>
          INCR + EXPIRE
        </div>
        <div class="d-arrow-down">&#8595; response</div>
      </div>
    </div>
  </div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box green" style="text-align:center">
          <strong>ALLOW</strong><br>
          Continue to app logic<br>
          Add X-RateLimit-* headers
        </div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box red" style="text-align:center">
          <strong>REJECT</strong><br>
          429 Too Many Requests<br>
          + Retry-After header
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-hop-by-hop",
		Title:       "Request Flow (Hop-by-Hop Detail)",
		Description: "Detailed six-hop request flow from client through DNS, ALB, middleware, Redis, and back",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">H1: Client sends API request with API key or JWT</div>
  <div class="d-label">HTTPS &#8594; TLS termination at ALB</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray">H2: ALB routes to healthy ECS container</div>
  <div class="d-label">Round-robin or least-connections &#8226; Health check: /health</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">H3: Middleware extracts rate limit key</div>
  <div class="d-label">Parse JWT &#8594; user_id OR API-Key header &#8594; key OR IP fallback</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber">H4: Lua script executes atomically in Redis</div>
  <div class="d-label">&lt; 1ms &#8226; Single eval replaces 3 separate commands &#8226; No race conditions</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green">H5a: ALLOWED &#8594; Continue to application logic</div>
  <div class="d-box red">H5b: REJECTED &#8594; Return 429 with Retry-After</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box blue">H6: Response includes X-RateLimit-* headers (always)</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-distributed-challenges",
		Title:       "Distributed Challenges & Solutions",
		Description: "How multiple app servers share Redis and how hash tags ensure correct sharding",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Challenge: Multiple App Servers</div>
        <div class="d-flow-v">
          <div class="d-flow">
            <div class="d-box blue">Server A</div>
            <div class="d-box blue">Server B</div>
            <div class="d-box blue">Server C</div>
          </div>
          <div class="d-arrow-down">&#8595; all share</div>
          <div class="d-box red" style="text-align:center"><strong>Same Redis Cluster</strong><br>Lua scripts &#8594; atomic counting</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Challenge: Redis Sharding</div>
        <div class="d-flow-v">
          <div class="d-box amber" style="text-align:center">Key: <code>rl:{user:123}:api</code></div>
          <div class="d-label">Hash tags {user:123} ensure same slot</div>
          <div class="d-flow">
            <div class="d-box gray">Shard 1</div>
            <div class="d-box green">Shard 2 &#10003;</div>
            <div class="d-box gray">Shard 3</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-multi-region",
		Title:       "Multi-Region Rate Limiting Strategies",
		Description: "Three approaches to multi-region rate limiting: gossip sync, per-region limits, and global Redis",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">&#9733; Async Gossip Sync (Recommended)</div>
      <div class="d-flow-v">
        <div class="d-flow">
          <div class="d-box blue" style="text-align:center">US-East<br>Redis (local)</div>
          <div class="d-box blue" style="text-align:center">EU-West<br>Redis (local)</div>
        </div>
        <div class="d-flow">
          <div class="d-arrow">&#8596; sync every 1-5s</div>
        </div>
        <div class="d-box green" style="text-align:center">Zero latency impact &#8226; Medium accuracy &#8226; Used by Cloudflare</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Per-Region Limits (Simplest)</div>
      <div class="d-flow-v">
        <div class="d-flow">
          <div class="d-box amber" style="text-align:center">US: 100/min</div>
          <div class="d-box amber" style="text-align:center">EU: 100/min</div>
        </div>
        <div class="d-box red" style="text-align:center">Total possible = 200/min &#8226; Acceptable for most APIs</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Global Redis (Avoid)</div>
      <div class="d-flow-v">
        <div class="d-box red" style="text-align:center">+50-100ms per request! &#8226; Cross-region latency &#8226; Unacceptable</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-data-model",
		Title:       "Data Model (Redis Key Schema)",
		Description: "Redis key schemas for token bucket, sliding window, fixed window, and rate limit rules",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header red">Token Bucket (HASH)</div>
      <div class="d-entity-body">
        <div class="pk">KEY: rl:token:{user}:{endpoint}</div>
        <div class="idx idx-hash">tokens FLOAT</div>
        <div class="idx idx-hash">last_refill FLOAT (unix ts)</div>
        <div>TTL: ceil(capacity / refill_rate) + 1</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header blue">Sliding Window (STRING)</div>
      <div class="d-entity-body">
        <div class="pk">KEY: rl:sw:{user}:{endpoint}:{window_ts}</div>
        <div class="idx idx-hash">value INTEGER (counter)</div>
        <div>TTL: window_sec &#215; 2</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header green">Rate Limit Rules (JSON/Config)</div>
      <div class="d-entity-body">
        <div class="pk">endpoint VARCHAR</div>
        <div class="idx">tier ENUM (free|pro|enterprise)</div>
        <div>limit INTEGER</div>
        <div>window_sec INTEGER</div>
        <div>algorithm ENUM (token|sliding|fixed)</div>
        <div>fail_mode ENUM (open|closed)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header amber">Fixed Window (STRING)</div>
      <div class="d-entity-body">
        <div class="pk">KEY: rl:fw:{user}:{endpoint}:{window_start}</div>
        <div class="idx idx-hash">value INTEGER (counter)</div>
        <div>TTL: window_sec</div>
      </div>
    </div>
  </div>
</div>
<div class="d-er-lines">
  <div class="d-er-connector">
    <span class="d-er-from">Rate Limit Rules</span>
    <span class="d-er-type">1:N</span>
    <span class="d-er-to">Token Bucket / Sliding Window keys</span>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-lua-execution",
		Title:       "Redis Lua Script Execution Flow",
		Description: "How a Lua script executes atomically inside Redis in a single TCP round-trip",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>App Server</strong><br>EVAL lua_script 1 key limit window now</div>
  <div class="d-arrow-down">&#8595; single TCP round-trip</div>
  <div class="d-box red" style="text-align:center;border:2px solid var(--red)">
    <strong>Redis (Atomic Execution)</strong><br>
    No other commands execute during Lua script
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-box amber">1. HMGET key tokens last_refill</div>
        <div class="d-box amber">2. Calculate refilled tokens</div>
        <div class="d-box amber">3. Check if tokens &#8805; 1</div>
        <div class="d-box amber">4. HMSET key new_tokens now</div>
        <div class="d-box amber">5. EXPIRE key ttl</div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center">
          <strong>Why Lua?</strong><br>
          &#10003; Atomic (no race conditions)<br>
          &#10003; 1 round-trip vs 3-5<br>
          &#10003; Server-side execution<br>
          &#10003; No distributed locks needed
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box blue" style="text-align:center">Return {allowed, remaining}</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-elasticache-topology",
		Title:       "ElastiCache Cluster Topology",
		Description: "Three-shard Redis cluster layout with primary and replica nodes per shard",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 1 (slots 0&#8211;5460)</div>
      <div class="d-flow-v">
        <div class="d-box red" style="text-align:center"><strong>Primary</strong><br>r6g.large</div>
        <div class="d-box gray" style="text-align:center">Replica<br>r6g.large</div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 2 (slots 5461&#8211;10922)</div>
      <div class="d-flow-v">
        <div class="d-box red" style="text-align:center"><strong>Primary</strong><br>r6g.large</div>
        <div class="d-box gray" style="text-align:center">Replica<br>r6g.large</div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 3 (slots 10923&#8211;16383)</div>
      <div class="d-flow-v">
        <div class="d-box red" style="text-align:center"><strong>Primary</strong><br>r6g.large</div>
        <div class="d-box gray" style="text-align:center">Replica<br>r6g.large</div>
      </div>
    </div>
  </div>
  <div class="d-label">3 shards &#215; 2 nodes = 6 nodes total &#8226; 300K+ ops/sec &#8226; 39 GB memory</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-tradeoffs",
		Title:       "Key Trade-offs Matrix",
		Description: "Trade-off comparison for fail-open vs fail-closed, centralized vs sidecar, and exact vs approximate",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fail-open vs Fail-closed</div>
      <div class="d-flow-v">
        <div class="d-box green">Fail-open: allow traffic if Redis down</div>
        <div class="d-label">&#10003; Most APIs: availability &gt; protection</div>
        <div class="d-box red">Fail-closed: block traffic if Redis down</div>
        <div class="d-label">&#10003; Payment APIs: prevent fraud</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Centralized vs Sidecar</div>
      <div class="d-flow-v">
        <div class="d-box blue">Centralized rate limit service</div>
        <div class="d-label">&#10003; Monolith: single place to manage</div>
        <div class="d-box purple">Sidecar (Envoy/Istio filter)</div>
        <div class="d-label">&#10003; Microservices: per-service limits</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Exact vs Approximate</div>
      <div class="d-flow-v">
        <div class="d-box green">Lua atomic scripts (exact)</div>
        <div class="d-label">&#10003; Single-region accuracy</div>
        <div class="d-box amber">Local counters + sync (approximate)</div>
        <div class="d-label">&#10003; Multi-region: accept 5% error</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-edge-cases",
		Title:       "Edge Cases & Mitigations",
		Description: "Six common edge cases (DDoS, hot keys, clock skew, key sharing, retry storms, failover) with mitigations",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">DDoS (Millions of IPs)</div>
        <div class="d-flow-v">
          <div class="d-box red">IP rate limiting alone won't work</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">AWS WAF + CloudFront geo-blocking</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Hot Key Problem</div>
        <div class="d-flow-v">
          <div class="d-box red">One user = 50% of checks on 1 shard</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Local cache (10ms TTL) + hash tags</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Clock Skew</div>
        <div class="d-flow-v">
          <div class="d-box red">Servers disagree on time</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Use Redis TIME command, not local clock</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">API Key Sharing</div>
        <div class="d-flow-v">
          <div class="d-box red">Multiple users share one key</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Per-key limits + anomaly detection</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Webhook Retry Storms</div>
        <div class="d-flow-v">
          <div class="d-box red">External service retries 1000x</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Separate tier + exponential backoff</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Redis Failover</div>
        <div class="d-flow-v">
          <div class="d-box red">Primary fails, counter reset</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Fail-open + counters rebuild</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-cost-breakdown",
		Title:       "Infrastructure Cost Breakdown",
		Description: "Monthly infrastructure costs for ElastiCache, ECS, and CloudWatch totaling ~$580/mo",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Monthly Infrastructure</div>
      <div class="d-flow-v">
        <div class="d-box red" style="text-align:center"><strong>ElastiCache Redis</strong><br>3 shards &#215; 2 nodes (r6g.large)<br><strong>$550/mo</strong></div>
        <div class="d-box gray" style="text-align:center"><strong>ECS Fargate (sidecar)</strong><br>CPU overhead minimal<br><strong>$0</strong> (existing containers)</div>
        <div class="d-box blue" style="text-align:center"><strong>CloudWatch</strong><br>Custom metrics for hits<br><strong>$30/mo</strong></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Total</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center;font-size:1.1rem"><strong>~$580/month</strong><br>incremental cost to add rate limiting</div>
        <div class="d-label">Rate limiting is cheap. The Redis cluster is likely already present for caching. Engineering cost in algorithms and edge cases is the real investment.</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-sub-problems",
		Title:       "Rate Limiter Sub-Problems & Building Blocks",
		Description: "Six key sub-problems: algorithm selection, atomic counting, multi-region sync, failure handling, rule config, and observability",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-subproblem green">
      <div class="d-subproblem-icon">&#9889;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Algorithm Selection</div>
        <div class="d-subproblem-desc">Token Bucket vs Sliding Window vs Fixed Window &#8212; trade-offs in burst tolerance, memory, accuracy</div>
      </div>
    </div>
    <div class="d-subproblem blue">
      <div class="d-subproblem-icon">&#128274;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Atomic Counting</div>
        <div class="d-subproblem-desc">Redis Lua scripts for race-free increment + check. Single round-trip replaces 3-5 commands</div>
      </div>
    </div>
    <div class="d-subproblem purple">
      <div class="d-subproblem-icon">&#127760;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Multi-Region Sync</div>
        <div class="d-subproblem-desc">Gossip-based async sync vs per-region limits vs global Redis. Latency vs accuracy trade-off</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-subproblem amber">
      <div class="d-subproblem-icon">&#9888;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Failure Handling</div>
        <div class="d-subproblem-desc">Fail-open vs fail-closed. Circuit breaker on Redis connection. Local cache fallback</div>
      </div>
    </div>
    <div class="d-subproblem red">
      <div class="d-subproblem-icon">&#128736;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Rule Configuration</div>
        <div class="d-subproblem-desc">Per-endpoint, per-tier, per-identity limits. Hot-reload without restart. A/B testing of limits</div>
      </div>
    </div>
    <div class="d-subproblem indigo">
      <div class="d-subproblem-icon">&#128200;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Observability</div>
        <div class="d-subproblem-desc">Hit rate metrics, false positive tracking, per-key analytics dashboard, alerting on anomalies</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

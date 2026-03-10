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
        <div class="d-box green" data-tip="Primary identity: API key header (X-API-Key). Fallback: JWT user_id claim. Enables per-customer quotas."><span class="d-step">1</span>Rate limit by user/API key <div class="d-tag green">&#10003; must have</div></div>
        <div class="d-box green" data-tip="Unauthenticated traffic identified by X-Forwarded-For IP. Use CIDR grouping to handle NAT/proxies."><span class="d-step">2</span>Rate limit by IP (anonymous)</div>
        <div class="d-box green" data-tip="Rules stored in config service (Redis or S3). Hot-reload without restart. Supports regex endpoint matching."><span class="d-step">3</span>Configurable rules per endpoint/tier</div>
        <div class="d-box green" data-tip="X-RateLimit-Limit (quota), X-RateLimit-Remaining (left), X-RateLimit-Reset (Unix epoch when window resets). Required by RFC 6585."><span class="d-step">4</span>Return X-RateLimit-* headers</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P1 &#8212; Important</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="All app server instances share one Redis cluster. Lua scripts ensure atomic read-modify-write with no race conditions across nodes.">Distributed counting across servers <div class="d-tag blue">Redis</div></div>
        <div class="d-box blue" data-tip="If Redis is unreachable, middleware catches the error and returns allow=true. Prevents a limiter outage from taking down your entire API.">Fail-open when Redis is down <div class="d-tag amber">circuit breaker</div></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P2 &#8212; Nice to Have</div>
      <div class="d-flow-v">
        <div class="d-box gray" data-tip="Each region tracks local counters and gossip-syncs every 1-5s. Accepts ~5% inaccuracy in exchange for zero cross-region latency.">Multi-region rate limiting</div>
        <div class="d-box gray" data-tip="UI for ops team to update limits, view top offenders, and preview rule changes before promotion to prod.">Admin dashboard for rules</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="Redis call (0.3-0.5ms avg) + key extraction (0.05ms) + header write (0.05ms). Total stays under 1ms p99.">Latency: &lt; 1ms overhead per request <span class="d-metric latency">&lt;1ms</span></div>
        <div class="d-box purple" data-tip="Rate limiter is in the hot path. Multi-AZ ElastiCache with automatic failover achieves 99.999% (5 nines). Fail-open as last resort.">Availability: 99.999% (in hot path) <span class="d-metric throughput">5 nines</span></div>
        <div class="d-box purple" data-tip="3-node Redis cluster at 100K ops/sec each = 300K. Scale to 6 nodes for peak 500K. Each check needs 2-3 Redis ops.">Throughput: 1M+ decisions/second <span class="d-metric throughput">1M+ /sec</span></div>
        <div class="d-box amber" data-tip="Multi-region gossip sync intentionally allows slight over-limit. Single-region Lua scripts are exact. Choose based on use case.">Accuracy: Allow 5% over-limit tolerance</div>
        <div class="d-box amber" data-tip="Default policy. The alternative (fail-closed) is appropriate only for payment or fraud-sensitive flows.">Fail-open: If limiter fails, allow traffic <div class="d-tag green">default policy</div></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Key Decisions</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Fail-open keeps your API up when Redis has a brief outage. Fail-closed is correct only when false positives (allowing over-limit requests) have financial or security consequences.">Fail-open vs Fail-closed? <div class="d-tag red">discuss in interview</div></div>
        <div class="d-label">Fail-open for most APIs, fail-closed for payments</div>
        <div class="d-box red" data-tip="Lua EVAL runs atomically inside Redis — no other commands execute concurrently. One round-trip replaces 3-5 separate INCR/GET/EXPIRE calls.">Exact vs Approximate counting? <div class="d-tag red">discuss in interview</div></div>
        <div class="d-label">Lua atomic scripts for single-region accuracy</div>
      </div>
    </div>
  </div>
</div>
<div class="d-legend">
  <div class="d-legend-item"><div class="d-legend-color green"></div>P0 must-have</div>
  <div class="d-legend-item"><div class="d-legend-color blue"></div>P1 important</div>
  <div class="d-legend-item"><div class="d-legend-color gray"></div>P2 nice-to-have</div>
  <div class="d-legend-item"><div class="d-legend-color purple"></div>Non-functional</div>
  <div class="d-legend-item"><div class="d-legend-color red"></div>Key decision</div>
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
        <div class="d-box blue" data-tip="Baseline steady-state. Each request triggers one rate limit check (one Lua EVAL call).">100K RPS across all services <span class="d-metric throughput">100K RPS</span></div>
        <div class="d-box blue" data-tip="Marketing events, flash sales, or viral traffic spikes. Cluster must handle this without degradation.">Peak (5x) = 500K checks/sec <span class="d-metric throughput">500K peak</span></div>
        <div class="d-box purple" data-tip="Sliding window counter: INCR (current window) + GET (previous window) = 2 ops. Token bucket: HMGET (read tokens + timestamp) + HMSET (write back) = 2 ops. Sliding window log (sorted set): ZREMRANGEBYSCORE + ZCARD + ZADD = 3 ops. Lua script bundles all into 1 round-trip — no intermediate state is exposed.">2–3 Redis ops per check (algorithm-dependent)</div>
        <div class="d-box purple" data-tip="Range spans baseline to peak: lower = 100K baseline RPS × 3 ops = 300K ops/sec (steady state); upper = 500K peak RPS × 3 ops = 1.5M ops/sec (5x traffic spike). With a 10ms local in-process cache, repeat keys for the same user bypass Redis entirely — effective Redis load drops to 20–40% of these numbers.">= 300K (100K×3, baseline) to 1.5M (500K×3, peak) Redis ops/sec <span class="d-metric throughput">1.5M peak</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Storage</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Redis HASH key overhead (~64B) + field names (~20B) + values (~16B). Actual measured ~80-120B per token bucket entry.">~100 bytes per counter <span class="d-metric size">~100B</span></div>
        <div class="d-box amber" data-tip="10M daily active users each tracked on 5 distinct rate-limited endpoints (login, search, upload, API, admin).">10M active users &#215; 5 endpoints</div>
        <div class="d-box amber" data-tip="Step by step: 10M users × 5 endpoints = 50M keys. Each key: ~100B overhead = 50M × 100B = 5,000,000,000 B = 5 GB. r6g.large has 13 GB RAM — well within capacity. With TTL expiry on inactive users, active memory is typically 40–60% of max: ~2–3 GB live.">= 50M keys × 100B = 5,000M B = 5 GB <span class="d-metric size">5 GB</span></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Redis Capacity</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="r6g.large benchmarked at 100-150K ops/sec for simple string/hash ops. Lua scripts are slightly slower (~80K/sec).">100K+ ops/sec per node <span class="d-metric throughput">100K/node</span></div>
        <div class="d-box green" data-tip="Handles baseline 100K RPS comfortably. Add replica per shard for read scaling and failover.">3-node cluster = 300K ops/sec <span class="d-metric throughput">300K ops/sec</span></div>
        <div class="d-box green" data-tip="Need to handle 1.5M ops/sec peak. r6g.large does ~100K ops/sec → 1.5M / 100K = 15 nodes min. But 6 nodes with ~250K ops/sec each (Lua scripts ~80–100K/sec with pipelining) — use 6 as baseline, add replicas for read HA. AWS ElastiCache online resharding allows shard expansion without downtime.">6 nodes for 500K peak (1.5M ops / 250K per-node capacity) <div class="d-tag green">&#10003; recommended</div></div>
      </div>
    </div>
  </div>
</div>
<div class="d-cols" style="margin-top:8px">
  <div class="d-col">
    <div class="d-number"><div class="d-number-value">100K</div><div class="d-number-label">Baseline RPS</div></div>
  </div>
  <div class="d-col">
    <div class="d-number"><div class="d-number-value">5 GB</div><div class="d-number-label">Redis Memory</div></div>
  </div>
  <div class="d-col">
    <div class="d-number"><div class="d-number-value">6</div><div class="d-number-label">Redis Nodes (peak)</div></div>
  </div>
  <div class="d-col">
    <div class="d-number"><div class="d-number-value">&lt;1ms</div><div class="d-number-label">Overhead/Request</div></div>
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
        <div class="d-group-title"><span class="d-status active"></span> 200 OK (Allowed)</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="X-RateLimit-Limit: the configured quota. X-RateLimit-Remaining: tokens left in this window. X-RateLimit-Reset: Unix epoch timestamp when the window resets. X-RateLimit-Policy: IETF draft header (limit;w=window_seconds)." style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">HTTP/1.1 200 OK
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 67
X-RateLimit-Reset: 1704067260
X-RateLimit-Policy: "100;w=60"</div>
        </div>
        <div class="d-caption">Always present on every response, even non-rate-limited paths.</div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title"><span class="d-status error"></span> 429 Too Many Requests</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="Retry-After: seconds until the client can retry. Clients MUST respect this and implement exponential backoff. Never retry immediately on 429." style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">HTTP/1.1 429 Too Many Requests
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1704067260
Retry-After: 45</div>
        </div>
        <div class="d-caption">RFC 6585 status code. Body should include JSON error with message.</div>
      </div>
    </div>
  </div>
  <div class="d-legend">
    <div class="d-legend-item"><div class="d-legend-color green"></div>Request allowed — headers inform client of remaining quota</div>
    <div class="d-legend-item"><div class="d-legend-color red"></div>Request rejected — Retry-After tells client when to retry</div>
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
        <div class="d-box gray" data-tip="Identify via JWT claim or API key metadata lookup. Free tier limit enforced to encourage upgrades.">Free: 100 req/min <span class="d-metric throughput">100/min</span></div>
        <div class="d-box blue" data-tip="Pro customers get 10x free tier. Stored in account metadata, cached in Redis for O(1) lookup.">Pro: 1,000 req/min <span class="d-metric throughput">1K/min</span></div>
        <div class="d-box green" data-tip="Custom limits negotiated per contract. Enterprise keys bypass shared-tier checks and hit a dedicated rule set.">Enterprise: 10,000 req/min <span class="d-metric throughput">10K/min</span> <div class="d-tag green">&#10003; custom SLA</div></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">By Endpoint Type</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Read endpoints are cheap. Higher limit reflects low server cost per request. Cached at CDN edge for popular resources.">GET /api/*: 500 req/min <span class="d-metric throughput">500/min</span></div>
        <div class="d-box amber" data-tip="Writes trigger DB mutations, validation, and side-effects. 5x lower than GET to protect backend.">POST /api/*: 100 req/min <span class="d-metric throughput">100/min</span></div>
        <div class="d-box red" data-tip="Brute-force protection. After 5 failures in 15 min, serve CAPTCHA. After 10, block IP for 1 hour. Store attempts in Redis ZSET.">POST /login: 5 req/15min + CAPTCHA <span class="d-metric throughput">5/15min</span> <div class="d-tag red">security</div></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">By Identity</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="Key format: rl:{api_key}:{endpoint}. Limits defined in key metadata, cached in Redis. Allows per-customer negotiated limits.">By API Key: per-key limits <div class="d-tag blue">Redis</div></div>
        <div class="d-box amber" data-tip="Use X-Forwarded-For (first untrusted IP). Must use CIDR grouping for IPv6. AWS WAF handles volumetric attacks before this layer.">By IP: 50 req/min (anonymous) <span class="d-metric throughput">50/min</span></div>
        <div class="d-box red" data-tip="Composite key: rl:{ip}:{path}. Detects credential stuffing attacks targeting login from many IPs.">By IP + Path: login brute-force <div class="d-tag red">security</div></div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Rules evaluated in priority order: API Key &rarr; JWT user_id &rarr; IP. Most specific match wins.</div>`,
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
        <div class="d-box green" data-tip="2 fields in a Redis HASH: tokens (float) and last_refill (unix ts). Fixed size regardless of request volume.">Memory: O(1) per key <span class="d-metric size">O(1)</span></div>
        <div class="d-box green" data-tip="Users can spend saved-up tokens instantly — up to capacity. Good for bursty mobile clients that batch requests.">Allows controlled bursts</div>
        <div class="d-box blue" data-tip="AWS API Gateway and Lambda use token bucket natively. Capacity = burst limit, refill_rate = rate limit.">Used by: AWS API Gateway <div class="d-tag blue">AWS</div></div>
        <div class="d-label">Best for: APIs that tolerate bursts</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">⭐ Sliding Window Counter <div class="d-tag green">recommended</div></div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Stores only 2 integers: current window count and previous window count. Same O(1) as fixed window but without the boundary burst.">Memory: O(1) per key <span class="d-metric size">O(1)</span></div>
        <div class="d-box green" data-tip="Weighted formula: weighted = prev_count × (1 - elapsed/window) + curr_count. Smoothly interpolates between windows.">No boundary burst problem <div class="d-tag green">&#10003; better than fixed</div></div>
        <div class="d-box blue" data-tip="Cloudflare uses this in their global rate limiting product (announced 2019). Stripe uses it for API key quotas.">Used by: Cloudflare, Stripe <div class="d-tag blue">production proven</div></div>
        <div class="d-label">Best for: Most use cases (recommended)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fixed Window Counter</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Single INCR key with EXPIRE. Dead simple. One Redis command per request.">Memory: O(1) per key <span class="d-metric size">O(1)</span></div>
        <div class="d-box red" data-tip="User sends 100 req at 11:59:59 and 100 req at 12:00:01. Both windows reset independently. 200 requests pass in 2 seconds.">2x burst at boundaries! <div class="d-tag red">known flaw</div></div>
        <div class="d-box gray" data-tip="INCR key; if result == 1 then EXPIRE key window_sec end. Can be done without Lua.">Simplest to implement</div>
        <div class="d-label">Best for: MVP only</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Sliding Window Log</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Stores exact timestamp of every request. No interpolation needed — just count entries in [now - window, now].">Highest accuracy</div>
        <div class="d-box red" data-tip="Redis ZSET with one entry per request. At 1000 req/min per user × 10M users = 10 billion entries in worst case.">Memory: O(N) per key! <span class="d-metric size">O(N)</span> <div class="d-tag red">memory risk</div></div>
        <div class="d-box gray" data-tip="N = max requests allowed per window. Login at 5/15min = 5 entries max. Acceptable. 1000 req/min = 1000 entries. Not acceptable.">N = requests in window</div>
        <div class="d-label">Best for: Low-rate limits (login)</div>
      </div>
    </div>
  </div>
</div>
<div class="d-legend">
  <div class="d-legend-item"><div class="d-legend-color green"></div>Recommended / strength</div>
  <div class="d-legend-item"><div class="d-legend-color amber"></div>Neutral / trade-off</div>
  <div class="d-legend-item"><div class="d-legend-color red"></div>Weakness / avoid</div>
  <div class="d-legend-item"><div class="d-legend-color blue"></div>Real-world usage</div>
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
    <div class="d-box indigo" data-tip="Two fields stored in Redis HASH: tokens (float) and last_refill (unix timestamp). O(1) memory per key." style="min-width:140px;text-align:center">
      <strong>Bucket</strong> <span class="d-metric size">O(1)</span><br>
      capacity: 10<br>
      tokens: 7
    </div>
    <div class="d-arrow">&#8592;</div>
    <div class="d-box green" data-tip="Tokens refill continuously. On each check, calculate elapsed time since last_refill and add (elapsed × rate). Cap at capacity." style="min-width:120px;text-align:center">
      <strong>Refill</strong><br>
      +2 tokens/sec
    </div>
  </div>
  <div class="d-arrow-down">&#8595; request arrives</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box green" data-tip="Decrement tokens by 1 atomically. Refill first, then check.">tokens &#8805; 1?</div>
        <div class="d-arrow-down">&#8595; YES</div>
        <div class="d-box green"><span class="d-status active"></span>tokens-- &#8594; ALLOW</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box red" data-tip="Return 429 with Retry-After = (1 - tokens) / refill_rate seconds.">tokens &lt; 1?</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red"><span class="d-status error"></span>REJECT (429)</div>
      </div>
    </div>
  </div>
  <div class="d-caption">Burst: can send up to <strong>capacity</strong> requests instantly, then rate-limited to refill_rate/sec. Used by AWS API Gateway.</div>
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
      <div class="d-box amber" data-tip="Key: rl:sw:{user}:{endpoint}:1704063600 (11:00 epoch). Counter = 84 requests were made during 11:00-12:00. TTL = window_sec × 2 keeps this key alive for the weighted calculation." style="text-align:center">
        count = 84<br>
        <small>11:00 &#8212; 12:00</small>
        <span class="d-metric size">84 req</span>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Current Window</div>
      <div class="d-box blue" data-tip="Key: rl:sw:{user}:{endpoint}:1704067200 (12:00 epoch). 15 minutes into this window, 36 requests have been made so far." style="text-align:center">
        count = 36<br>
        <small>12:00 &#8212; 13:00</small>
        <span class="d-metric size">36 req</span>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; Lua script reads both keys atomically</div>
  <div class="d-box purple" data-tip="Formula: overlap = (window_end - now) / window_sec = (12:00 - 11:45) / 60min = 0.25... wait, now=12:15 means 75% of prev window is within [now-60min, now]. weighted = prev×overlap + curr." style="text-align:center">
    <strong>Weighted Count</strong> <div class="d-tag green">⭐ key insight</div><br>
    now = 12:15 &#8594; 75% of prev window overlaps<br>
    weighted = 84 &#215; 0.75 + 36 = <strong>99</strong><br>
    limit = 100 &#8594; <span class="d-status active"></span><strong>ALLOW</strong> (1 remaining)
  </div>
  <div class="d-legend">
    <div class="d-legend-item"><div class="d-legend-color amber"></div>Previous window count (weighted)</div>
    <div class="d-legend-item"><div class="d-legend-color blue"></div>Current window count (full weight)</div>
    <div class="d-legend-item"><div class="d-legend-color purple"></div>Weighted total (must be &lt; limit)</div>
  </div>
  <div class="d-caption">Eliminates boundary burst: smoothly transitions between windows using time-weighted average. <strong>Used by Cloudflare and Stripe.</strong> O(1) memory — only 2 counters per key.</div>
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
        <div class="d-box gray" data-tip="User is quiet for 59 minutes. Fixed window counter stays at 0 the whole time. Nothing interesting happens." style="text-align:center">
          &#8230; quiet &#8230;<br>
          <small>0 requests until 11:59</small>
        </div>
        <div class="d-box red" data-tip="User fires 100 requests in the last 30 seconds of the window. Counter hits exactly 100 = limit. All 100 pass. Window 1 counter = 100 / 100 limit." style="text-align:center">
          <span class="d-status error"></span><strong>11:59:30 &#8594; 100 requests!</strong> <span class="d-metric throughput">100 req</span><br>
          <small>All at end of window — counter = 100/100</small>
          <div class="d-tag red">limit reached</div>
        </div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Window 2 (12:00&#8211;13:00)</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Clock strikes 12:00:00. Fixed window counter resets to 0. User fires another 100 requests immediately. All 100 pass. The system just allowed 200 req in 30 seconds." style="text-align:center">
          <span class="d-status error"></span><strong>12:00:01 &#8594; 100 requests!</strong> <span class="d-metric throughput">100 req</span><br>
          <small>Counter reset to 0 — all 100 pass again</small>
          <div class="d-tag red">counter reset!</div>
        </div>
        <div class="d-box gray" style="text-align:center">
          &#8230; rest of window &#8230;
        </div>
      </div>
    </div>
  </div>
  <div class="d-box red" data-tip="This is the classic fixed window attack. A knowledgeable user can always double their rate limit by timing requests around window boundaries. Use Sliding Window Counter to eliminate this." style="text-align:center">
    <span class="d-status error"></span><strong>Result: 200 requests in ~30 seconds!</strong> <span class="d-metric throughput">2x limit</span><br>
    Limit was 100/min but user got 2× at the boundary
  </div>
  <div class="d-caption">Fix: use Sliding Window Counter. weighted = prev × overlap + curr. The boundary burst attack is mathematically impossible with weighted counts.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-architecture",
		Title:       "Rate Limiter Architecture — Full Request Path",
		Description: "End-to-end request flow from client through ALB, API server, rate limit middleware, and Redis",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="HTTPS request with API key in header or JWT Bearer token. Client-side should implement exponential backoff on 429." style="text-align:center"><span class="d-step">1</span><strong>Client Request</strong></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" data-tip="DNS resolution adds ~1-5ms. Use Route 53 latency-based routing for multi-region." style="text-align:center">Route 53 (DNS) <span class="d-metric latency">&lt;5ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="Application Load Balancer distributes across ECS containers. TLS termination happens here." style="text-align:center">ALB (Load Balancer) <span class="d-metric latency">~1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="ECS Fargate containers. Rate limit middleware runs as first middleware in chain, before auth." style="text-align:center"><span class="d-step">2</span><strong>API Server (ECS)</strong> <span class="d-metric throughput">100K RPS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="Middleware extracts key (API key &gt; JWT user_id &gt; IP fallback), loads tier rules, then calls Redis Lua script. Total overhead &lt;1ms." style="text-align:center;border:2px solid var(--amber)">
    <span class="d-step">3</span><strong>Rate Limit Middleware</strong><br>
    Extract key (API key or IP) &#8594; Check rules &#8594; Query Redis
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box red" data-tip="Single EVAL command runs atomically. No race conditions. 1 TCP round-trip replaces 3-5 separate commands." style="text-align:center">
          <span class="d-step">4</span><strong>Redis</strong> <span class="d-metric latency">&lt;1ms</span><br>
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
        <div class="d-box green" data-tip="Headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset (Unix epoch)" style="text-align:center">
          <span class="d-status active"></span><strong>ALLOW</strong><br>
          Continue to app logic<br>
          Add X-RateLimit-* headers
        </div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box red" data-tip="Return 429 with Retry-After header (seconds). Client should back off exponentially." style="text-align:center">
          <span class="d-status error"></span><strong>REJECT</strong><br>
          429 Too Many Requests<br>
          + Retry-After header
        </div>
      </div>
    </div>
  </div>
  <div class="d-legend">
    <div class="d-legend-item"><div class="d-legend-color green"></div>Success path</div>
    <div class="d-legend-item"><div class="d-legend-color red"></div>Rejection path</div>
    <div class="d-legend-item"><div class="d-legend-color amber"></div>Decision point</div>
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
  <div class="d-box blue" data-tip="Client includes API key in X-API-Key header or JWT in Authorization: Bearer. HTTPS/TLS 1.3."><span class="d-step">1</span>Client sends API request with API key or JWT</div>
  <div class="d-label">HTTPS &#8594; TLS termination at ALB <span class="d-metric latency">~10ms TLS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" data-tip="ALB distributes via round-robin or least-connections. Health check on /health every 30s."><span class="d-step">2</span>ALB routes to healthy ECS container <span class="d-metric latency">~1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="Priority: API-Key header → JWT user_id → X-Forwarded-For IP → socket IP. Key format: rl:{identity}:{endpoint}"><span class="d-step">3</span>Middleware extracts rate limit key <span class="d-metric latency">&lt;0.1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="EVAL sha1 1 key limit window now → returns {allowed: bool, remaining: int, reset: epoch}"><span class="d-step">4</span>Lua script executes atomically in Redis <span class="d-metric latency">&lt;1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green"><span class="d-step">5a</span><span class="d-status active"></span>ALLOWED &#8594; Continue to application logic</div>
  <div class="d-box red"><span class="d-step">5b</span><span class="d-status error"></span>REJECTED &#8594; Return 429 with Retry-After</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box blue" data-tip="Always present regardless of allow/reject: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, X-RateLimit-Policy"><span class="d-step">6</span>Response includes X-RateLimit-* headers (always)</div>
  <div class="d-caption">Total overhead: <strong>&lt;2ms</strong> added to every request. The Redis call dominates (~0.5ms avg).</div>
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
            <div class="d-box blue" data-tip="Each ECS container runs independently. Without shared state, each would maintain its own counter — user could exceed the limit N×containers times.">Server A</div>
            <div class="d-box blue" data-tip="Server B also sends all rate limit checks to the shared Redis cluster. Lua atomicity ensures no double-counting.">Server B</div>
            <div class="d-box blue" data-tip="Horizontal scaling does not multiply the quota. Adding more servers never grants users more requests.">Server C</div>
          </div>
          <div class="d-arrow-down">&#8595; all share</div>
          <div class="d-box red" data-tip="Lua EVAL is atomic: HMGET + compute + HMSET run as one indivisible unit. No other Redis command executes in between." style="text-align:center"><strong>Same Redis Cluster</strong> <div class="d-tag blue">Redis</div><br>Lua scripts &#8594; atomic counting <span class="d-metric latency">&lt;1ms</span></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Challenge: Redis Sharding</div>
        <div class="d-flow-v">
          <div class="d-box amber" data-tip="Without hash tags, rl:user:123:login and rl:user:123:api could land on different shards. Lua EVAL cannot span shards — it would error with CROSSSLOT." style="text-align:center">Key: <code>rl:{user:123}:api</code> <div class="d-tag amber">hash tag</div></div>
          <div class="d-label">Hash tags <code>{user:123}</code> ensure all keys for one user land on the same shard</div>
          <div class="d-flow">
            <div class="d-box gray" data-tip="Slot range 0-5460. CRC16('user:123') mod 16384 maps to a slot outside this range.">Shard 1</div>
            <div class="d-box green" data-tip="All rl:{user:123}:* keys share the same hash slot and are guaranteed to land here. Lua script runs atomically on this shard." style="text-align:center"><span class="d-status active"></span>Shard 2 &#10003; <div class="d-tag green">user:123 here</div></div>
            <div class="d-box gray" data-tip="Slot range 10923-16383. user:123 maps below this range.">Shard 3</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-caption">Key insight: Lua scripts are atomic but cannot span shards. Hash tags <code>{}</code> pin all related keys to one shard, enabling safe multi-key Lua execution.</div>
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
      <div class="d-group-title">⭐ Async Gossip Sync <div class="d-tag green">recommended</div></div>
      <div class="d-flow-v">
        <div class="d-flow">
          <div class="d-box blue" data-tip="US-East counts locally. Periodically ships delta to EU-West (e.g. 'user 123 has used 47 more tokens'). Sub-ms local reads." style="text-align:center">US-East<br>Redis (local) <span class="d-metric latency">&lt;1ms local</span></div>
          <div class="d-box blue" data-tip="EU-West similarly counts locally. Merges incoming deltas using a CRDT-like max(local, remote) strategy." style="text-align:center">EU-West<br>Redis (local) <span class="d-metric latency">&lt;1ms local</span></div>
        </div>
        <div class="d-flow">
          <div class="d-arrow">&#8596; sync every 1-5s <span class="d-metric latency">async</span></div>
        </div>
        <div class="d-box green" data-tip="~5% over-limit tolerance is acceptable for most APIs. A user globally limited to 1000/min might see up to 1050 in edge cases." style="text-align:center"><span class="d-status active"></span>Zero latency impact &#8226; ~5% accuracy loss &#8226; Used by Cloudflare <div class="d-tag blue">Cloudflare</div></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Per-Region Limits (Simplest)</div>
      <div class="d-flow-v">
        <div class="d-flow">
          <div class="d-box amber" data-tip="Each region enforces the full quota independently. No cross-region communication needed." style="text-align:center">US: 100/min <span class="d-metric throughput">100/min</span></div>
          <div class="d-box amber" data-tip="EU users get a fresh 100/min limit. A user routing through both regions can effectively get 200/min globally." style="text-align:center">EU: 100/min <span class="d-metric throughput">100/min</span></div>
        </div>
        <div class="d-box red" data-tip="A determined user with VPN can switch regions to double their quota. For most use cases this over-provisioning is acceptable." style="text-align:center">Total possible = 200/min &#8226; Acceptable for most APIs <div class="d-tag amber">trade-off</div></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Global Redis <div class="d-tag red">avoid</div></div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Cross-region Redis call from EU-West to us-east-1: ~80ms round-trip. This adds 80ms to every API request. Never acceptable in a hot path." style="text-align:center"><span class="d-status error"></span>+50-100ms per request! <span class="d-metric latency">+80ms</span><br>Cross-region latency &#8226; Unacceptable</div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Interview answer: Start with per-region limits (simplest). Upgrade to gossip sync when global accuracy becomes a product requirement.</div>`,
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
      <div class="d-entity-header red">Token Bucket (HASH) <div class="d-tag blue">Redis HASH</div></div>
      <div class="d-entity-body">
        <div class="pk" data-tip="Hash tag {user} ensures all endpoints for one user land on same shard. Enables multi-key Lua scripts.">KEY: rl:token:{user}:{endpoint}</div>
        <div class="idx idx-hash" data-tip="Current token count as float. Supports fractional tokens for sub-second refill rates.">tokens FLOAT <span class="d-metric size">8B</span></div>
        <div class="idx idx-hash" data-tip="Unix timestamp (float) of last refill. Used to calculate elapsed time: now - last_refill.">last_refill FLOAT (unix ts) <span class="d-metric size">8B</span></div>
        <div data-tip="Auto-expire when user is inactive. Buffer of +1 ensures key lives long enough to catch late-arriving requests.">TTL: ceil(capacity / refill_rate) + 1</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header blue">⭐ Sliding Window (STRING) <div class="d-tag green">recommended</div></div>
      <div class="d-entity-body">
        <div class="pk" data-tip="window_ts = floor(now / window_sec). Two keys exist simultaneously: current and previous window.">KEY: rl:sw:{user}:{endpoint}:{window_ts}</div>
        <div class="idx idx-hash" data-tip="Simple integer counter. INCR is O(1) and atomic without Lua. Lua needed only for the weighted read.">value INTEGER (counter) <span class="d-metric size">8B</span></div>
        <div data-tip="Keep previous window alive for the weighted calculation. Multiply by 2 to ensure both current and prior window coexist.">TTL: window_sec &#215; 2</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header green">Rate Limit Rules (JSON/Config) <div class="d-tag amber">config store</div></div>
      <div class="d-entity-body">
        <div class="pk" data-tip="Matched via regex or exact path. More specific rules take precedence. e.g. /login overrides /api/*.">endpoint VARCHAR</div>
        <div class="idx" data-tip="Determines base quota. Resolved from API key metadata or JWT claims at auth time.">tier ENUM (free|pro|enterprise)</div>
        <div data-tip="Requests per window_sec. Combined: limit=100, window_sec=60 means 100 req/min.">limit INTEGER</div>
        <div>window_sec INTEGER</div>
        <div data-tip="Algorithm is per-endpoint. Login uses sliding_log (exact, low N). APIs use sliding_counter (O(1)).">algorithm ENUM (token|sliding|fixed)</div>
        <div data-tip="open: allow on Redis error. closed: reject on Redis error. Default open except for payment endpoints.">fail_mode ENUM (open|closed)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header amber">Fixed Window (STRING)</div>
      <div class="d-entity-body">
        <div class="pk" data-tip="window_start = floor(now / window_sec) * window_sec. Resets sharply at each boundary — source of the burst problem.">KEY: rl:fw:{user}:{endpoint}:{window_start}</div>
        <div class="idx idx-hash">value INTEGER (counter)</div>
        <div data-tip="Key expires exactly at window end. New window starts fresh at 0, enabling the 2x boundary burst attack.">TTL: window_sec</div>
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
</div>
<div class="d-caption">All keys use hash tags <code>{user}</code> to ensure same-shard locality. TTLs auto-clean inactive user keys — no manual eviction needed.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rl-lua-execution",
		Title:       "Redis Lua Script Execution Flow",
		Description: "How a Lua script executes atomically inside Redis in a single TCP round-trip",
		ContentFile: "problems/rate-limiter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="EVALSHA uses SHA1 of the cached script — avoids re-sending full Lua source on every call. Falls back to EVAL on cache miss (first call only)." style="text-align:center"><span class="d-step">1</span><strong>App Server</strong><br>EVALSHA &lt;sha1&gt; 1 key limit window now <div class="d-tag blue">1 TCP round-trip</div></div>
  <div class="d-arrow-down">&#8595; single TCP round-trip <span class="d-metric latency">~0.3ms network</span></div>
  <div class="d-box red" data-tip="Redis is single-threaded. Lua script blocks ALL other commands during execution. Keep scripts under 1ms — never do I/O or loops over large data inside Lua." style="text-align:center;border:2px solid var(--red)">
    <span class="d-step">2</span><strong>Redis (Atomic Execution)</strong> <span class="d-metric latency">&lt;0.2ms CPU</span><br>
    No other commands execute during Lua script <div class="d-tag amber">single-threaded lock</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Read current state: remaining tokens (float) and timestamp of last refill (float). HMGET returns both fields in one call."><span class="d-step">2a</span> HMGET key tokens last_refill</div>
        <div class="d-box amber" data-tip="new_tokens = min(capacity, old_tokens + (now - last_refill) × refill_rate). Capped at capacity to prevent over-accumulation."><span class="d-step">2b</span> Calculate refilled tokens</div>
        <div class="d-box amber" data-tip="If new_tokens &ge; 1: allowed=true, new_tokens -= 1. Else: allowed=false, compute retry_after = (1 - new_tokens) / refill_rate."><span class="d-step">2c</span> Check if tokens &#8805; 1 &rarr; allow/deny</div>
        <div class="d-box amber" data-tip="Persist updated token count and current timestamp. HMSET replaces both fields atomically inside the same script."><span class="d-step">2d</span> HMSET key new_tokens now</div>
        <div class="d-box amber" data-tip="TTL = ceil(capacity / refill_rate) + 1 second. Key auto-expires when user is inactive. No separate cleanup job needed."><span class="d-step">2e</span> EXPIRE key ttl <span class="d-metric latency">auto-cleanup</span></div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Alternative without Lua: WATCH + MULTI/EXEC (optimistic locking). Requires retry loop on contention. Much higher latency under load." style="text-align:center">
          <strong>Why Lua?</strong> <div class="d-tag green">&#10003; recommended</div><br>
          &#10003; Atomic — no race conditions<br>
          &#10003; 1 round-trip vs 3-5 commands<br>
          &#10003; Server-side — no network for math<br>
          &#10003; No WATCH/MULTI/EXEC retry loop<br>
          &#10003; No distributed locks needed
        </div>
        <div class="d-number"><div class="d-number-value">1</div><div class="d-number-label">TCP Round-Trip</div></div>
        <div class="d-number"><div class="d-number-value">&lt;0.5ms</div><div class="d-number-label">Total Redis Time</div></div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box blue" data-tip="Go struct unmarshalled from Redis response. allowed drives the middleware decision. remaining written to X-RateLimit-Remaining header. reset_at written to X-RateLimit-Reset." style="text-align:center"><span class="d-step">3</span>Return {allowed: bool, remaining: int, reset_at: epoch}</div>
  <div class="d-caption">Total Redis time: <strong>&lt;0.5ms</strong> per check. Script SHA1 is cached after first EVAL — all subsequent calls use EVALSHA (faster).</div>
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
        <div class="d-box red" data-tip="Primary handles all writes (EVAL/EVALSHA). r6g.large: 2 vCPU, 13 GB RAM, Graviton2. ~100K ops/sec for Lua workloads." style="text-align:center"><span class="d-status active"></span><strong>Primary</strong><br>r6g.large <span class="d-metric throughput">100K ops/s</span></div>
        <div class="d-box gray" data-tip="Replica receives async replication from primary. Promotes automatically on primary failure (~10-30s failover). Can serve read-only queries." style="text-align:center">Replica<br>r6g.large <div class="d-tag gray">standby</div></div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 2 (slots 5461&#8211;10922)</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Handles ~33% of keyspace. Hash tag routing ensures related keys (same user) land on the same shard automatically." style="text-align:center"><span class="d-status active"></span><strong>Primary</strong><br>r6g.large <span class="d-metric throughput">100K ops/s</span></div>
        <div class="d-box gray" data-tip="Multi-AZ: replica in different Availability Zone than primary. Survives AZ-level failure." style="text-align:center">Replica<br>r6g.large <div class="d-tag gray">multi-AZ</div></div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 3 (slots 10923&#8211;16383)</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Third shard rounds out the 16384-slot keyspace. AWS ElastiCache supports online resharding to add shards without downtime." style="text-align:center"><span class="d-status active"></span><strong>Primary</strong><br>r6g.large <span class="d-metric throughput">100K ops/s</span></div>
        <div class="d-box gray" style="text-align:center">Replica<br>r6g.large <div class="d-tag gray">multi-AZ</div></div>
      </div>
    </div>
  </div>
  <div class="d-cols" style="margin-top:8px">
    <div class="d-col">
      <div class="d-number"><div class="d-number-value">6</div><div class="d-number-label">Total Nodes</div></div>
    </div>
    <div class="d-col">
      <div class="d-number"><div class="d-number-value">300K+</div><div class="d-number-label">ops/sec baseline</div></div>
    </div>
    <div class="d-col">
      <div class="d-number"><div class="d-number-value">78 GB</div><div class="d-number-label">Total Memory</div></div>
    </div>
    <div class="d-col">
      <div class="d-number"><div class="d-number-value">~$550</div><div class="d-number-label">per month</div></div>
    </div>
  </div>
  <div class="d-caption">3 shards × 2 nodes = 6 nodes total <span class="d-metric throughput">300K+ ops/sec</span> <span class="d-metric size">78 GB memory</span> <span class="d-metric cost">~$550/mo</span>. Scale to 6 shards for 500K peak ops/sec.</div>
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
        <div class="d-box green" data-tip="Circuit breaker catches Redis connection errors and returns allow=true. Brief Redis outage does not cascade to a full API outage. Risk: a few seconds of unlimited traffic.">⭐ Fail-open: allow traffic if Redis down <div class="d-tag green">default</div></div>
        <div class="d-label">&#10003; Most APIs: availability &gt; protection</div>
        <div class="d-box red" data-tip="On Redis error, return 503 Service Unavailable. Prevents any request from slipping through. Risk: Redis hiccup takes down your entire API.">Fail-closed: block traffic if Redis down <div class="d-tag red">use sparingly</div></div>
        <div class="d-label">&#10003; Payment APIs: prevent fraud</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Centralized vs Sidecar</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Shared middleware library or dedicated rate limit microservice. All services call the same Redis cluster. Easy to monitor in one place.">Centralized rate limit service <div class="d-tag blue">monolith</div></div>
        <div class="d-label">&#10003; Single place to manage, update rules</div>
        <div class="d-box purple" data-tip="Envoy sidecar or Istio EnvoyFilter applies rate limiting at the service mesh layer. Each service gets independent limits without code changes.">Sidecar (Envoy/Istio filter) <div class="d-tag purple">microservices</div></div>
        <div class="d-label">&#10003; Microservices: per-service limits, no code change</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Exact vs Approximate</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Lua EVAL is atomic. No race condition between reading the counter and incrementing it. Perfect accuracy within a single region.">⭐ Lua atomic scripts (exact) <div class="d-tag green">single-region</div></div>
        <div class="d-label">&#10003; Single-region: 100% accuracy guaranteed</div>
        <div class="d-box amber" data-tip="Each app server keeps local counters (Go sync.Map or similar) and periodically ships deltas to Redis. Accepts ~5% over-limit in exchange for zero Redis calls on most requests.">Local counters + sync (approximate) <div class="d-tag amber">multi-region</div></div>
        <div class="d-label">&#10003; Multi-region: accept 5% error, save ~80% Redis ops</div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Interview tip: Always state your default (fail-open, centralized, exact) and explain when you would deviate (payments, microservices, multi-region).</div>`,
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
          <div class="d-box red" data-tip="1M attacker IPs each sending 1 req/sec = 1M RPS. Your Redis cluster cannot store or check 1M unique keys fast enough. Per-IP rate limiting fails at DDoS scale."><span class="d-status error"></span>IP rate limiting alone won't work <div class="d-tag red">scale failure</div></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="AWS WAF Layer 7 rules block known bad IPs at the edge. CloudFront geo-blocking drops traffic from high-risk regions before it reaches your origin."><span class="d-status active"></span>AWS WAF + CloudFront geo-blocking <div class="d-tag green">edge layer</div></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Hot Key Problem</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="Enterprise customer making 10K req/min = one Redis shard receiving 10K ops/sec for a single key hash slot. Other users on that shard suffer high latency."><span class="d-status error"></span>One user = 50% of checks on 1 shard <div class="d-tag red">hotspot</div></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="App server caches the rate limit decision locally for 10ms. 10K req/min = 167/sec; at 10ms TTL, only ~17 actually hit Redis. 90% cache hit rate for hot keys."><span class="d-status active"></span>Local cache (10ms TTL) + hash tags <div class="d-tag green">90% cache hit</div></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Clock Skew</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="Two app servers may have clocks 100ms apart. If the token bucket uses local time for refill calculation, a fast clock refills tokens faster — user gets more quota than intended."><span class="d-status error"></span>Servers disagree on time <div class="d-tag red">accuracy bug</div></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="Inside the Lua script, call redis.call('TIME') to get the authoritative Redis server time. Consistent across all app servers regardless of NTP drift."><span class="d-status active"></span>Use Redis TIME command, not local clock <div class="d-tag green">single source of truth</div></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">API Key Sharing</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="Team at a startup shares one API key. 10 developers all hitting the Pro limit (1000/min) together = 10K req/min through one key. Legitimate but potentially abusive."><span class="d-status error"></span>Multiple users share one key <div class="d-tag amber">gray area</div></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="Track unique source IPs per API key. If key A suddenly comes from 50 unique IPs vs historical 3, flag for review. Can also enforce IP allowlist per key."><span class="d-status active"></span>Per-key limits + anomaly detection <div class="d-tag green">ML signal</div></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Webhook Retry Storms</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="External payment provider or CI system hits your webhook endpoint, gets a 5xx, and retries with no backoff. 1000 retries in 10 seconds. Your 429 responses make them retry harder."><span class="d-status error"></span>External service retries 1000x <div class="d-tag red">retry storm</div></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="Separate rate limit tier for known webhook senders (by IP or API key). Return 429 with Retry-After: 60. Document in your API that webhooks must use exponential backoff."><span class="d-status active"></span>Separate tier + exponential backoff <div class="d-tag green">isolate</div></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Redis Failover</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="Primary fails. Replica promotes (~10-30s window). During promotion, writes fail. After promotion, the new primary has the replica's data — may be slightly behind (async replication)."><span class="d-status error"></span>Primary fails, counter reset risk <div class="d-tag amber">~10-30s window</div></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="During failover window (~30s), fail-open allows all traffic. After promotion, new primary starts with slightly stale counters — resets to 0 in worst case. Acceptable for 30s."><span class="d-status active"></span>Fail-open + counters rebuild <div class="d-tag green">graceful degradation</div></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-caption">Most edge cases reduce to: (1) shed load at the edge (WAF/CDN) before it reaches Redis, and (2) fail-open gracefully rather than cascading.</div>
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
        <div class="d-box red" data-tip="r6g.large: 2 vCPUs, 13 GB RAM, Graviton2 processor. ~$150/mo on-demand per node × 6 nodes = $900/mo. Reserved pricing (~40% discount) brings it to ~$550/mo." style="text-align:center"><strong>ElastiCache Redis</strong> <div class="d-tag blue">Redis</div><br>3 shards &#215; 2 nodes (r6g.large)<br><strong>$550/mo</strong> <span class="d-metric cost">94% of total</span></div>
        <div class="d-box gray" data-tip="Rate limiting middleware runs inside existing ECS Fargate containers. Adds ~0.5% CPU overhead per container. No new containers needed." style="text-align:center"><strong>ECS Fargate (middleware)</strong><br>CPU overhead minimal (~0.5%)<br><strong>$0</strong> incremental <span class="d-metric cost">0%</span></div>
        <div class="d-box blue" data-tip="4 custom metrics × $0.30/metric/mo = $1.20. Add dimensions (endpoint, tier) = ~100 metric streams × $0.30 = $30/mo. Includes alarms for spike detection." style="text-align:center"><strong>CloudWatch Metrics</strong><br>rate_limit_hit, miss, error per tier<br><strong>$30/mo</strong> <span class="d-metric cost">5%</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Total Cost</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="This is the incremental cost assuming no existing Redis cluster. If Redis is already present for app caching, the marginal cost is only the extra memory (5 GB = ~$50/mo more)." style="text-align:center;font-size:1.1rem"><span class="d-status active"></span><strong>~$580/month</strong><br>incremental cost to add rate limiting</div>
        <div class="d-label">Rate limiting is cheap. The Redis cluster is likely already present for caching. Engineering cost in algorithms and edge cases is the real investment.</div>
      </div>
    </div>
    <div class="d-cols" style="margin-top:8px">
      <div class="d-col">
        <div class="d-number"><div class="d-number-value">$550</div><div class="d-number-label">Redis (6 nodes)</div></div>
      </div>
      <div class="d-col">
        <div class="d-number"><div class="d-number-value">$30</div><div class="d-number-label">Observability</div></div>
      </div>
      <div class="d-col">
        <div class="d-number"><div class="d-number-value">$580</div><div class="d-number-label">Total / Month</div></div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Cost per request at 100K RPS = $580 / (100K × 86400 × 30) = <strong>$0.0000000022/request</strong>. Essentially free per-request.</div>`,
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
        <div class="d-subproblem-title">Algorithm Selection <div class="d-tag green">core decision</div></div>
        <div class="d-subproblem-desc">Token Bucket vs Sliding Window vs Fixed Window &#8212; trade-offs in burst tolerance, memory O(1) vs O(N), accuracy. <strong>Sliding Window Counter recommended</strong> for most APIs.</div>
      </div>
    </div>
    <div class="d-subproblem blue">
      <div class="d-subproblem-icon">&#128274;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Atomic Counting <div class="d-tag blue">Redis Lua</div></div>
        <div class="d-subproblem-desc">Redis Lua EVAL scripts for race-free increment + check. Single round-trip &lt;0.5ms replaces 3-5 separate commands. Hash tags ensure same-shard locality.</div>
      </div>
    </div>
    <div class="d-subproblem purple">
      <div class="d-subproblem-icon">&#127760;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Multi-Region Sync <div class="d-tag amber">advanced</div></div>
        <div class="d-subproblem-desc">Gossip-based async sync (1-5s) vs per-region limits (simplest) vs global Redis (avoid: +80ms). Accept ~5% accuracy loss for zero latency impact.</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-subproblem amber">
      <div class="d-subproblem-icon">&#9888;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Failure Handling <div class="d-tag red">interview signal</div></div>
        <div class="d-subproblem-desc">Fail-open (default) vs fail-closed (payments). Circuit breaker catches Redis errors. Local 10ms cache as last-resort fallback. Prevents limiter outage cascading to API outage.</div>
      </div>
    </div>
    <div class="d-subproblem red">
      <div class="d-subproblem-icon">&#128736;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Rule Configuration <div class="d-tag blue">operational</div></div>
        <div class="d-subproblem-desc">Per-endpoint, per-tier, per-identity limits stored in config service. Hot-reload via Redis pub/sub without restart. Supports A/B testing of limit values.</div>
      </div>
    </div>
    <div class="d-subproblem indigo">
      <div class="d-subproblem-icon">&#128200;</div>
      <div class="d-subproblem-text">
        <div class="d-subproblem-title">Observability <div class="d-tag blue">CloudWatch</div></div>
        <div class="d-subproblem-desc">Emit rate_limit_hit, rate_limit_miss, rate_limit_error metrics per endpoint/tier. False-positive rate (legitimate users blocked) is the key SLO to track. Alert on sudden spike in 429 rate.</div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">The hardest sub-problem is <strong>Atomic Counting</strong> — correctly implementing the Lua script with hash tags, TTL, and clock sync. Algorithm selection is easier than it looks once you rule out Sliding Log for high-limit APIs.</div>`,
	})
}

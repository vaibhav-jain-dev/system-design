package diagrams

func registerURLShortener(r *Registry) {
	r.Register(&Diagram{
		Slug:        "url-api-design",
		Title:       "API Design",
		Description: "Prioritized API endpoints for the URL shortener service.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P0 — Core (Must Have)</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Returns 201 + short_code. Idempotent if same URL shortened twice — each call creates a new code."><div class="d-tag green">Core</div>POST /api/v1/urls &#8594; shorten URL, returns short_code</div>
        <div class="d-box green" data-tip="301 Moved Permanently. Browser caches redirect for 24h. Zero latency on repeat visits."><div class="d-tag green">Core</div>GET /{short_code} &#8594; 301 redirect to original URL</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P1 — Important</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="User-supplied alias up to 30 chars. Validated against reserved words, profanity filter, and homoglyph check.">PUT /api/v1/urls?alias=my-link &#8594; custom alias</div>
        <div class="d-box blue" data-tip="DynamoDB TTL attribute. Default 5 years (1825 days). Background scanner deletes within 48h of expiry at no cost.">TTL per URL (default 5 years) &#8594; auto-expiration</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P2 — Nice to Have</div>
      <div class="d-flow-v">
        <div class="d-box gray" data-tip="Click count + geo breakdown from click_analytics table. Served async — never affects redirect latency.">GET /api/v1/urls/{code}/stats &#8594; click analytics</div>
        <div class="d-box gray" data-tip="Sets expires_at = now() for soft delete; DynamoDB TTL cleans up within 48h. Cached 301s still work until they expire.">DELETE /api/v1/urls/{code} &#8594; remove short URL</div>
      </div>
    </div>
  </div>
</div>
<div class="d-legend">
  <div class="d-legend-item"><div class="d-legend-color green"></div>P0 must-have</div>
  <div class="d-legend-item"><div class="d-legend-color blue"></div>P1 important</div>
  <div class="d-legend-item"><div class="d-legend-color gray"></div>P2 nice-to-have</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-redirects",
		Title:       "301 vs 302 Redirect Comparison",
		Description: "Trade-offs between permanent and temporary HTTP redirects.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">301 Moved Permanently</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Browser stores redirect in local cache; subsequent clicks never leave the browser. Zero latency.">Browser caches redirect locally</div>
        <div class="d-box green" data-tip="CloudFront serves cached 301 from the nearest edge PoP. Eliminates round-trip to origin.">CDN caches at edge PoPs</div>
        <div class="d-box green" data-tip="No origin hit = no server cost. Scales to millions of clicks at near-zero marginal cost.">Repeat visits never hit origin</div>
        <div class="d-box amber" data-tip="Cache-Control: max-age=86400. Once cached, destination is locked for up to 24h in browsers.">Cannot change destination after caching</div>
        <div class="d-box amber" data-tip="No server request = no click event. Use 302 or server-side redirect to preserve analytics.">Lose analytics on cached clicks</div>
      </div>
      <div class="d-flow-v">
        <div class="d-label">Traffic flow (2nd click by same user):</div>
        <div class="d-box blue" data-tip="User's browser that previously received and cached the 301 response.">Browser</div>
        <div class="d-arrow-down">&#8595; cached locally</div>
        <div class="d-box green"><span class="d-status active"></span>&#10003; Direct to destination <span class="d-metric latency">0ms</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">302 Found (Temporary)</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="HTTP spec says 302 must not be cached unless response includes explicit Cache-Control or Expires.">Browser does NOT cache</div>
        <div class="d-box red" data-tip="Each redirect request travels to origin: DNS → CDN → ALB → ECS. Adds latency and cost.">Every click hits your servers</div>
        <div class="d-box red" data-tip="At 57K peak RPS, 302 sends all traffic to origin vs 301 absorbing 60%+ at CDN edge."><span class="d-metric throughput">2&#8212;5&#215;</span> more origin traffic</div>
        <div class="d-box green" data-tip="Every request hits the server, so every click can be logged with geo, device, referrer.">Full analytics on every click</div>
        <div class="d-box green" data-tip="Update original_url in DynamoDB; next request picks up new destination immediately.">Can change destination anytime</div>
      </div>
      <div class="d-flow-v">
        <div class="d-label">Traffic flow (2nd click by same user):</div>
        <div class="d-box blue" data-tip="User's browser. Has no cached 302 response — must always go back to origin.">Browser</div>
        <div class="d-arrow-down">&#8595; no cache</div>
        <div class="d-box purple" data-tip="ECS API server processes the GET /{code} request, looks up DynamoDB/Redis, issues 302 Location header.">Your server <span class="d-metric latency">5&#8212;50ms</span></div>
        <div class="d-arrow-down">&#8595; redirect</div>
        <div class="d-box amber" data-tip="Final destination the user reaches after the 302 redirect chain.">Destination</div>
      </div>
    </div>
  </div>
</div>
<div class="d-flow-v">
  <div class="d-box indigo" data-tip="Twitter uses 301. Bitly uses 302 (analytics). Choose based on whether you need click counting or speed.">Best practice: 301 by default for performance &#8212; switch to 302 only for links needing click tracking or mutable destinations</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-nfr-estimates",
		Title:       "Non-Functional Requirements & Back-of-Envelope Estimates",
		Description: "Availability targets, latency goals, and scale math for the URL shortener.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">NFR Targets</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="52 minutes downtime per year. Multi-AZ deployment, auto-scaling, health checks."><span class="d-step">1</span>Availability: 99.99% <span class="d-metric throughput">52 min/yr</span> <div class="d-tag green">4 nines</div></div>
        <div class="d-box green" data-tip="p50 < 5ms (cache hit), p99 < 10ms (cache miss + DB). CDN hits are < 2ms."><span class="d-step">2</span>Redirect latency: &lt;10ms p99 <span class="d-metric latency">&lt;10ms p99</span></div>
        <div class="d-box blue" data-tip="DynamoDB 11 9's durability. Point-in-time recovery enabled. Cross-region replication for DR."><span class="d-step">3</span>Durability: never lose a URL mapping <div class="d-tag blue">11 nines</div></div>
        <div class="d-box amber" data-tip="Cache TTL 1h means updates take up to 1h to propagate. Acceptable for URL shortener — URLs rarely change."><span class="d-step">4</span>Consistency: eventual OK (cache + DB sync)</div>
        <div class="d-box red" data-tip="Fail-open: if Redis is down, serve from DynamoDB only. Redirects still work at higher latency. Never return 503 for existing URLs.">Fail-open: Redis down &#8594; serve from DB only <div class="d-tag red">discuss in interview</div></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Math</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="100M / 86,400s = 1,157 avg. 5x safety factor for peak bursts during viral events.">100M new URLs/day = <span class="d-metric throughput">1,157 write QPS</span> (5x peak = <span class="d-metric throughput">5,785</span>)</div>
        <div class="d-box purple" data-tip="10:1 read-to-write ratio is typical for link shorteners. Most traffic is redirects, not URL creation.">10:1 read:write = <span class="d-metric throughput">11,570 read QPS</span> (5x peak = <span class="d-metric throughput">57,850</span>)</div>
        <div class="d-box purple" data-tip="5-year horizon is the default TTL. After 5 years URLs auto-expire unless explicitly renewed.">5 years: 100M &#215; 365 &#215; 5 = <span class="d-metric size">182.5B</span> total URLs</div>
        <div class="d-box amber" data-tip="500 bytes per row: short_code (7B) + original_url (200B) + metadata (293B). DynamoDB max item 400KB — no concern.">Storage: 182.5B &#215; 500B = <span class="d-metric size">~91 TB</span> over 5 years</div>
      </div>
      <div class="d-flow">
        <div class="d-number"><div class="d-number-value">100M</div><div class="d-number-label">URLs/day</div></div>
        <div class="d-number"><div class="d-number-value">57K</div><div class="d-number-label">peak read RPS</div></div>
        <div class="d-number"><div class="d-number-value">91 TB</div><div class="d-number-label">5-yr storage</div></div>
      </div>
    </div>
  </div>
</div>
<div class="d-legend">
  <div class="d-legend-item"><div class="d-legend-color green"></div>Availability / latency</div>
  <div class="d-legend-item"><div class="d-legend-color blue"></div>Durability</div>
  <div class="d-legend-item"><div class="d-legend-color amber"></div>Consistency trade-off</div>
  <div class="d-legend-item"><div class="d-legend-color purple"></div>Scale math</div>
  <div class="d-legend-item"><div class="d-legend-color red"></div>Key decision</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-bandwidth-infra",
		Title:       "Bandwidth, Caching & Infrastructure Estimates",
		Description: "Bandwidth, cache sizing, and infrastructure estimates for each tier.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Bandwidth Estimation</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="500B per redirect response: short_code lookup + 301 Location header. Average, not peak.">Read: 11,570 RPS &#215; 500B = <span class="d-metric throughput">5.8 MB/s avg</span></div>
        <div class="d-box blue" data-tip="5x safety factor applied to average RPS. CloudFront absorbs most of this at edge.">Peak read bandwidth: 57,850 &#215; 500B = <span class="d-metric throughput">29 MB/s</span></div>
        <div class="d-box blue" data-tip="1KB request body: URL string + metadata. Write path is 10x less traffic than reads.">Write: 1,157 RPS &#215; 1KB (req body) = <span class="d-metric throughput">1.2 MB/s</span></div>
        <div class="d-box amber" data-tip="At $0.085/GB CloudFront egress: ~$42.50/day or ~$1,275/mo just for egress. Budget for this.">Daily egress: 5.8 MB/s &#215; 86,400 = <span class="d-metric size">~500 GB/day</span></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Caching Estimation (80/20 Rule)</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Pareto principle: top 20% of URLs (viral content, popular links) account for 80% of redirect traffic.">20% hot URLs generate 80% traffic</div>
        <div class="d-box green" data-tip="Daily active set resets each day. Cumulative storage grows but working set stays ~20M URLs.">Daily URLs to cache: 100M &#215; 0.2 = <span class="d-metric size">20M URLs</span></div>
        <div class="d-box green" data-tip="r6g.large = 13 GB RAM, ~$92/mo. Comfortably holds 20M × 500B = 10 GB working set.">Cache memory: 20M &#215; 500B = <span class="d-metric size">10 GB</span> (fits 1 Redis node)</div>
        <div class="d-box purple" data-tip="95% cache hit means 5% of 57K peak = ~2.9K RPS actually reach DynamoDB. Well within on-demand limits.">At 90% cache hit: only <span class="d-metric throughput">1,157 RPS</span> reach DB</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Infrastructure Sizing</div>
      <div class="d-flow-v">
        <div class="d-box indigo" data-tip="ECS Fargate tasks, 2 vCPU / 4 GB each. Auto-scale on CPU 60%. Stateless — scale horizontally without coordination.">API Servers: <span class="d-metric throughput">4-6 instances</span> (each handles 2-3K RPS)</div>
        <div class="d-box indigo" data-tip="On-demand billing: $1.25/M reads, $1.25/M writes. Switch to provisioned + reserved at 10K+ sustained RPS to save ~70%.">Database: DynamoDB on-demand (auto-scales)</div>
        <div class="d-box indigo" data-tip="r6g.large: 13 GB RAM, ~$92/mo. Cluster mode: 3 primary + 2 replicas each = 9 nodes for HA and read scaling.">Cache: 1x r6g.large (<span class="d-metric size">13 GB</span>) &#8594; 3-node cluster at scale</div>
        <div class="d-box indigo" data-tip="CloudFront absorbs 60%+ of reads at the edge. Each cache hit saves a DynamoDB read (~$0.00000125).">CDN: CloudFront 400+ PoPs (absorbs <span class="d-metric throughput">60%+</span> reads)</div>
        <div class="d-box indigo" data-tip="KGS pre-generates key batches of 1000. Lambda invocations cost ~$0.0000002 each. 2 instances for HA.">KGS: 2 Lambda instances + DynamoDB table</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-data-model",
		Title:       "Data Model",
		Description: "DynamoDB table schema, GSI design, and access patterns.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">DynamoDB: urls table</div>
      <div class="d-entity">
        <div class="d-entity-header indigo">urls</div>
        <div class="d-entity-body">
          <div class="pk">short_code STRING (Partition Key)</div>
          <div>original_url STRING</div>
          <div class="idx idx-gsi">user_id STRING</div>
          <div>created_at NUMBER (epoch)</div>
          <div>expires_at NUMBER (TTL auto-delete)</div>
          <div>click_count NUMBER (atomic ADD)</div>
        </div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">GSI: user-urls-index</div>
      <div class="d-entity">
        <div class="d-entity-header purple">user-urls-index</div>
        <div class="d-entity-body">
          <div class="pk">user_id STRING</div>
          <div class="pk">created_at NUMBER (SK)</div>
          <div>short_code, original_url, click_count (projected)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Access Patterns</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="GetItem by partition key is always O(1) regardless of table size. Single-digit ms.">Redirect: GetItem(PK=short_code) <span class="d-metric latency">&lt;5ms</span></div>
        <div class="d-box green" data-tip="attribute_not_exists prevents overwriting existing URLs. Idempotent with condition.">Create: PutItem(condition: attribute_not_exists)</div>
        <div class="d-box blue" data-tip="GSI query returns user's URLs sorted by created_at. Cursor-based pagination with ExclusiveStartKey.">User URLs: Query(GSI, user_id) &#8594; paginated</div>
        <div class="d-box blue" data-tip="Soft delete: set expires_at to now. Hard delete: DeleteItem. Cache invalidation via CloudFront API.">Delete: DeleteItem(PK=short_code)</div>
        <div class="d-box amber" data-tip="DynamoDB TTL runs a background scanner. Items deleted within 48h of expiry. No cost for TTL deletes.">Expiration: TTL auto-deletes expired URLs</div>
        <div class="d-box amber" data-tip="UpdateItem with SET click_count = click_count + :one. Atomic, no read-modify-write race.">Analytics: atomic ADD on click_count</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-base62-encoding",
		Title:       "Base62 Encoding: How a Number Becomes a Short Code",
		Description: "Step-by-step visualization of Base62 encoding from decimal to short code.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Step 1: Counter Value (decimal 123,456,789)</div>
    <div class="d-bitfield">
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">1</div>
        <div class="d-bitfield-name">10^8</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">2</div>
        <div class="d-bitfield-name">10^7</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">3</div>
        <div class="d-bitfield-name">10^6</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">4</div>
        <div class="d-bitfield-name">10^5</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">5</div>
        <div class="d-bitfield-name">10^4</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">6</div>
        <div class="d-bitfield-name">10^3</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">7</div>
        <div class="d-bitfield-name">10^2</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">8</div>
        <div class="d-bitfield-name">10^1</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">9</div>
        <div class="d-bitfield-name">10^0</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; divide repeatedly by 62</div>
  <div class="d-group">
    <div class="d-group-title">Step 2: Convert to Base62 Digits (remainders)</div>
    <div class="d-bitfield">
      <div class="d-bitfield-segment" style="background: var(--green-bg, #dcfce7); border-color: var(--green-border, #86efac);">
        <div class="d-bitfield-bits">8</div>
        <div class="d-bitfield-name">62^4</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--green-bg, #dcfce7); border-color: var(--green-border, #86efac);">
        <div class="d-bitfield-bits">M</div>
        <div class="d-bitfield-name">62^3</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--green-bg, #dcfce7); border-color: var(--green-border, #86efac);">
        <div class="d-bitfield-bits">0</div>
        <div class="d-bitfield-name">62^2</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--green-bg, #dcfce7); border-color: var(--green-border, #86efac);">
        <div class="d-bitfield-bits">k</div>
        <div class="d-bitfield-name">62^1</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--green-bg, #dcfce7); border-color: var(--green-border, #86efac);">
        <div class="d-bitfield-bits">X</div>
        <div class="d-bitfield-name">62^0</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; map to alphabet</div>
  <div class="d-group">
    <div class="d-group-title">Step 3: Base62 Alphabet (62 characters)</div>
    <div class="d-bitfield">
      <div class="d-bitfield-segment" style="background: var(--purple-bg, #f3e8ff); border-color: var(--purple-border, #d8b4fe);">
        <div class="d-bitfield-bits">a&#8212;z</div>
        <div class="d-bitfield-name">0&#8212;25</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--amber-bg, #fef3c7); border-color: var(--amber-border, #fcd34d);">
        <div class="d-bitfield-bits">A&#8212;Z</div>
        <div class="d-bitfield-name">26&#8212;51</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--indigo-bg, #e0e7ff); border-color: var(--indigo-border, #a5b4fc);">
        <div class="d-bitfield-bits">0&#8212;9</div>
        <div class="d-bitfield-name">52&#8212;61</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; result</div>
  <div class="d-group">
    <div class="d-group-title">Output: Short Code</div>
    <div class="d-row">
      <div class="d-box green">123,456,789 &#8594; "8M0kX"</div>
      <div class="d-box blue">7 chars = 62^7 = 3.52 trillion unique codes</div>
      <div class="d-box amber">At 100M/day &#8594; lasts 96 years</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-bit-layout",
		Title:       "URL ID Bit Layout (43 bits for 7-char Base62)",
		Description: "Three options for structuring 43-bit URL IDs: counter, snowflake, and KGS.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Option A: Simple Counter ID (43 bits)</div>
    <div class="d-bitfield">
      <div class="d-bitfield-segment" style="background: var(--green-bg, #dcfce7); border-color: var(--green-border, #86efac);">
        <div class="d-bitfield-bits">bits 42&#8212;0</div>
        <div class="d-bitfield-name">Counter value (43 bits = 8.8 trillion)</div>
      </div>
    </div>
    <div class="d-label">Simple but sequential &#8212; guessable URLs</div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Option B: Snowflake-style ID (43 bits)</div>
    <div class="d-bitfield">
      <div class="d-bitfield-segment" style="background: var(--blue-bg, #dbeafe); border-color: var(--blue-border, #93c5fd);">
        <div class="d-bitfield-bits">bits 42&#8212;27</div>
        <div class="d-bitfield-name">Timestamp (16 bits &#8212; seconds mod 65536)</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--purple-bg, #f3e8ff); border-color: var(--purple-border, #d8b4fe);">
        <div class="d-bitfield-bits">bits 26&#8212;17</div>
        <div class="d-bitfield-name">Shard/Worker ID (10 bits &#8212; 1024 workers)</div>
      </div>
      <div class="d-bitfield-segment" style="background: var(--amber-bg, #fef3c7); border-color: var(--amber-border, #fcd34d);">
        <div class="d-bitfield-bits">bits 16&#8212;0</div>
        <div class="d-bitfield-name">Sequence (17 bits &#8212; 131K/sec/worker)</div>
      </div>
    </div>
    <div class="d-label">Non-sequential, embeds metadata &#8212; but more complex</div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Option C: KGS Pre-generated (recommended)</div>
    <div class="d-bitfield">
      <div class="d-bitfield-segment" style="background: var(--indigo-bg, #e0e7ff); border-color: var(--indigo-border, #a5b4fc);">
        <div class="d-bitfield-bits">bits 42&#8212;0</div>
        <div class="d-bitfield-name">Random unique value (43 bits, pre-generated by KGS)</div>
      </div>
    </div>
    <div class="d-label">Non-guessable, zero collisions, no coordination at runtime</div>
    <div class="d-tag green">recommended</div>
  </div>
  <div class="d-row">
    <div class="d-box green" data-tip="Base62 alphabet (a-z, A-Z, 0-9) maps any integer to a URL-safe 7-char code with no special characters.">All options &#8594; Base62 encode &#8594; 7-character short code</div>
    <div class="d-box purple" data-tip="43-bit space (8.8T) exceeds 62^7 = 3.5T, ensuring safe padding. At 100M/day the space lasts ~96 years.">43 bits &#8776; 8.8 trillion values &#8811; 3.5T (62^7)</div>
  </div>
  <div class="d-caption">Option A (counter) is simplest but produces guessable sequential URLs. Option C (KGS) is recommended for production: zero collision risk and no runtime coordination needed.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-keygen-approaches",
		Title:       "Key Generation Approaches Compared",
		Description: "Comparison of counter, hash-truncate, and KGS approaches for key generation.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Approach 1: Counter + Base62</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Single-writer bottleneck. Zookeeper or DynamoDB atomic counter. Max ~10K TPS before contention."><span class="d-step">1</span>Global atomic counter (DB or Zookeeper)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="O(log₆₂ N) divisions. 7-char code = enough for 3.5T unique URLs. Deterministic, no randomness."><span class="d-step">2</span>counter++ &#8594; Base62 encode &#8594; "Ab3xK9"</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Conditional PutItem to guarantee uniqueness. Sequential IDs are guessable — enumerate competitor URLs."><span class="d-step">3</span>Store mapping in DynamoDB</div>
      </div>
      <div class="d-tag amber">sequential — guessable</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Approach 2: MD5/SHA + Truncate</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="MD5 produces 128-bit hash. Same URL always produces same hash = deduplication possible, but same short code for identical long URLs."><span class="d-step">1</span>Hash long URL: MD5("https://...")</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Take the least-significant 43 bits from the 128-bit MD5 digest. Birthday paradox: 0.14% collision at 100B URLs."><span class="d-step">2</span>Take first 43 bits &#8594; Base62(7 chars)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="GetItem to check if code exists. Retry with hash(url + salt++) if collision. Adds DB read per write."><span class="d-step">3</span>Check collision &#8594; retry if exists</div>
      </div>
      <div class="d-tag amber">collision retry needed</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Approach 3: Key Generation Service</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Lambda + DynamoDB. Pre-generates and marks keys 'used'. Runs offline. Zero runtime coordination."><span class="d-step">1</span>KGS pre-generates batches of unique keys</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Each ECS task fetches 1000 keys at startup and refills when 200 remain. O(1) local assignment, no network call per URL."><span class="d-step">2</span>App server requests batch of <span class="d-metric size">1,000</span> keys</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Pop from local in-memory queue. Sub-microsecond. Lost keys on crash are acceptable (tiny fraction of keyspace)."><span class="d-step">3</span>Assign next key from local batch <span class="d-metric latency">&lt;0.1ms</span></div>
      </div>
      <div class="d-tag green">recommended</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-collision-resolution",
		Title:       "Collision Resolution: MD5/SHA Truncation Approach",
		Description: "Flow of MD5 hash truncation with collision detection and retry strategy.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="Any long URL submitted by client. MD5 is deterministic: same URL always maps to same hash."><span class="d-step">1</span>Input: long_url = "https://example.com/very/long/path"</div>
  <div class="d-arrow-down">&#8595; MD5 hash</div>
  <div class="d-box purple" data-tip="128-bit MD5 digest. Cryptographically weak but fast (&lt;1µs). Collision resistance not needed here — we check DB anyway."><span class="d-step">2</span>MD5 &#8594; "e4d909c290d0fb1ca068ffaddf22cbd0" <span class="d-metric size">128 bits</span></div>
  <div class="d-arrow-down">&#8595; take first 43 bits</div>
  <div class="d-box green" data-tip="43 bits → 8.8T unique values. Base62 encodes 43 bits into exactly 7 URL-safe characters."><span class="d-step">3</span>Truncate &#8594; Base62 encode &#8594; "Ab3xK9p" (7 chars)</div>
  <div class="d-arrow-down">&#8595; check DB</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label"><span class="d-status active"></span>No collision</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Conditional PutItem: attribute_not_exists(short_code). Atomic — no race condition between check and write."><span class="d-step">4</span>&#10003; Store mapping in DB</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><span class="d-step">5</span>Return 201 + short URL to client</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label"><span class="d-status error"></span>Collision detected</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="PutItem conditional check failed: another URL owns this short_code. Extremely rare (&lt;0.14% at 100B URLs).">&#215; short_code already exists</div>
        <div class="d-arrow-down">&#8595; retry strategy</div>
        <div class="d-box amber" data-tip="Append integer salt to URL string before hashing. Each retry produces a different hash, hence different truncation.">Append counter: hash(url + "1")</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="New 43-bit truncation from new hash. Statistically independent from previous attempt.">Re-truncate &#8594; Base62 &#8594; new code</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Max 5 retries. P(5 consecutive collisions) ≈ (0.0014)^5 ≈ negligible. Return 503 only if all fail.">Check DB again (max 5 retries)</div>
      </div>
    </div>
  </div>
  <div class="d-caption">Birthday paradox: at 100B URLs, collision probability ≈ 0.14% per attempt. KGS approach eliminates this entirely.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-architecture",
		Title:       "System Architecture",
		Description: "End-to-end system architecture from clients through CDN, API, cache, and storage.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="Browser, mobile apps, and API consumers. 100M new URLs/day, 10:1 read:write ratio.">Clients (Browser / Mobile / API consumers) <span class="d-metric throughput">57K peak RPS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Latency-based routing sends users to the nearest region. Weighted routing for blue/green deploys.">Route 53 (DNS) &#8594; latency-based routing <span class="d-metric latency">&lt;5ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="400+ edge PoPs worldwide. Cache 301 redirects with TTL 24h. Absorbs 60%+ of read traffic."><span class="d-status active"></span>CloudFront (CDN) &#8212; 301 cached at edge <span class="d-metric throughput">60% absorbed</span></div>
  <div class="d-arrow-down">&#8595; cache miss</div>
  <div class="d-box indigo" data-tip="Application Load Balancer. TLS termination, /health check every 30s. Cross-AZ distribution.">ALB (Load Balancer) &#8212; TLS termination <span class="d-metric latency">~1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box green" data-tip="Stateless ECS Fargate tasks. Auto-scale on CPU 60%. Each handles 2-3K RPS.">API Server 1</div>
    <div class="d-box green">API Server 2</div>
    <div class="d-box green">API Server N</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box red" data-tip="Cache-aside pattern. 10 GB covers 20M hot URLs. 95% hit rate. TTL 1h.">ElastiCache Redis <span class="d-metric latency">&lt;1ms</span></div>
    <div class="d-box amber" data-tip="Source of truth. On-demand capacity. GetItem O(1) by short_code PK.">DynamoDB <span class="d-metric latency">~5ms</span></div>
    <div class="d-box purple" data-tip="Pre-generates batches of unique keys. Lambda + DynamoDB. Each API server holds a local batch of 1000 keys.">KGS</div>
  </div>
  <div class="d-arrow-down">&#8595; async analytics</div>
  <div class="d-row">
    <div class="d-box gray" data-tip="Click events streamed asynchronously. No impact on redirect latency.">Kinesis (click stream)</div>
    <div class="d-box gray">Lambda (aggregate)</div>
    <div class="d-box gray">S3 (archive)</div>
  </div>
  <div class="d-legend">
    <div class="d-legend-item"><div class="d-legend-color purple"></div>Network edge</div>
    <div class="d-legend-item"><div class="d-legend-color green"></div>Compute</div>
    <div class="d-legend-item"><div class="d-legend-color red"></div>Cache</div>
    <div class="d-legend-item"><div class="d-legend-color amber"></div>Storage</div>
    <div class="d-legend-item"><div class="d-legend-color gray"></div>Analytics (async)</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-write-read-paths",
		Title:       "Write & Read Paths (Hop by Hop)",
		Description: "Detailed write and read request paths through every system component.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">WRITE PATH (Create Short URL)</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="d-step">1</span>Client: POST /api/v1/urls</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Round-robin to healthy ECS container. TLS terminated here."><span class="d-step">2</span>ALB &#8594; healthy ECS task <span class="d-metric latency">~1ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="URL format validation + Google Safe Browsing API check (async, non-blocking). Reject malicious URLs."><span class="d-step">3</span>Validate URL + Safe Browsing <span class="d-metric latency">~50ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="Each server holds 1000 pre-generated keys. Refill when 200 remain. No contention, O(1) assignment."><span class="d-step">4</span>Get next key from local KGS batch <span class="d-metric latency">&lt;0.1ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Conditional PutItem: attribute_not_exists(short_code). Guarantees uniqueness. On-demand billing."><span class="d-step">5</span>DynamoDB PutItem <span class="d-metric latency">~5ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray"><span class="d-step">6</span>Return 201 + short URL</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">READ PATH (Redirect)</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="d-step">1</span>Client: GET /Ab3xK9</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="301 cached with Cache-Control: max-age=86400. 60% of traffic served here."><span class="d-step">2</span>CloudFront: cache HIT? &#8594; 301 <span class="d-metric latency">&lt;5ms</span></div>
        <div class="d-arrow-down">&#8595; MISS</div>
        <div class="d-box indigo"><span class="d-step">3</span>ALB &#8594; ECS task</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="Cache key: url:{short_code}. Value: original URL. Hit rate ~95% for warm codes."><span class="d-step">4</span>ElastiCache: GET url:Ab3xK9 <span class="d-metric latency">&lt;1ms</span></div>
        <div class="d-label"><span class="d-status active"></span>HIT (95%)? &#8594; return 301</div>
        <div class="d-arrow-down">&#8595; MISS (5%)</div>
        <div class="d-box amber" data-tip="Single-digit ms. GetItem by PK is O(1) regardless of table size."><span class="d-step">5</span>DynamoDB: GetItem <span class="d-metric latency">~5ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red"><span class="d-step">6</span>Write to ElastiCache (TTL 1h)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray"><span class="d-step">7</span>Return 301 + Location header</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-caching-layers",
		Title:       "Four-Layer Caching Architecture",
		Description: "Four cache layers from browser to DynamoDB with hit rates and latencies.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="HTTP 301 sets Cache-Control: max-age=86400 (24h). Browser caches locally — subsequent clicks are instant. Trade-off: lose analytics for cached clicks.">L1: Browser Cache (301 Cache-Control: max-age=86400) <span class="d-metric latency">0ms</span> <span class="d-metric throughput">~30% hit</span></div>
  <div class="d-arrow-down">&#8595; miss</div>
  <div class="d-box purple" data-tip="400+ edge PoPs globally. TTL 24h. Cache key is the URL path. Invalidation via CreateInvalidation API if URL destination changes.">L2: CloudFront Edge (400+ PoPs, TTL 24h) <span class="d-metric latency">&lt;5ms</span> <span class="d-metric throughput">~60% hit</span></div>
  <div class="d-arrow-down">&#8595; miss</div>
  <div class="d-box red" data-tip="Cache-aside pattern. Key: url:{short_code}. 20M hot URLs × 500B = 10 GB. r6g.large = 13 GB. TTL 1h to balance freshness vs hit rate.">L3: ElastiCache Redis (regional, TTL 1h) <span class="d-metric latency">&lt;1ms</span> <span class="d-metric throughput">~95% hit</span></div>
  <div class="d-arrow-down">&#8595; miss (5%)</div>
  <div class="d-box amber" data-tip="Source of truth. On-demand capacity scales automatically. GetItem by PK is always O(1). TTL attribute auto-deletes expired URLs.">L4: DynamoDB (source of truth) <span class="d-metric latency">~5ms</span> <span class="d-metric throughput">~2.9K RPS</span></div>
  <div class="d-caption">Combined hit rate across all layers: <strong>99.7%</strong> of redirects served from cache. Only <strong>~170 RPS</strong> reach DynamoDB at peak.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-redirect-lookup",
		Title:       "Example: Redirect Lookup",
		Description: "Sample redirect lookup results showing cache hits, misses, and error cases.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box green" data-tip="Cache key: url:Ab3xK9. Redis returns original_url in O(1). ECS issues 301 Location header immediately. No DynamoDB call.">
    <span class="d-status active"></span><span class="d-step">1</span>GET /Ab3xK9 <span class="d-metric latency">0.3ms</span>
    <div class="d-label">ElastiCache HIT &#8594; 301 to https://example.com/very/long/path?utm=...</div>
  </div>
  <div class="d-box amber" data-tip="Cache key url:Xz9mQ2 not found (cold or evicted). Falls through to DynamoDB GetItem. Writes result back to Redis with TTL 1h.">
    <span class="d-step">2</span>GET /Xz9mQ2 <span class="d-metric latency">4.8ms</span>
    <div class="d-label">Cache MISS &#8594; DynamoDB GetItem &#8594; 301 to https://docs.google.com/spreadsheets/d/...</div>
  </div>
  <div class="d-box red" data-tip="DynamoDB TTL scanner deleted the item. GetItem returns no item but expires_at was in the past. Serve 410 Gone (not 404) — signals permanent removal to crawlers.">
    <span class="d-status error"></span><span class="d-step">3</span>GET /expired1 <span class="d-metric latency">5.1ms</span>
    <div class="d-label">DynamoDB TTL expired &#8594; HTTP 410 Gone (permanent removal signal)</div>
  </div>
  <div class="d-box red" data-tip="GetItem returns no item and no TTL attribute — code was never issued. Return 404 Not Found. Do NOT reveal whether a code existed-then-expired vs never existed.">
    <span class="d-status error"></span><span class="d-step">4</span>GET /notfound <span class="d-metric latency">4.2ms</span>
    <div class="d-label">DynamoDB no item &#8594; HTTP 404 Not Found</div>
  </div>
</div>
<div class="d-legend">
  <div class="d-legend-item"><div class="d-legend-color green"></div>Cache HIT (&lt;1ms)</div>
  <div class="d-legend-item"><div class="d-legend-color amber"></div>Cache MISS (~5ms)</div>
  <div class="d-legend-item"><div class="d-legend-color red"></div>Error (410/404)</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-scaling-strategy",
		Title:       "Scaling Strategy by Component",
		Description: "Horizontal scaling approach for API and cache layers with sizing math.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">API Layer</div>
      <div class="d-flow">
        <div class="d-number"><div class="d-number-value">23</div><div class="d-number-label">tasks at peak</div></div>
        <div class="d-number"><div class="d-number-value">2.5K</div><div class="d-number-label">RPS per task</div></div>
      </div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="CloudWatch metric: ECSServiceAverageCPUUtilization. Target 60% CPU. Scale-out cooldown 60s. Scale-in cooldown 300s.">ECS auto-scaling on CPU (target <span class="d-metric throughput">60%</span>)</div>
        <div class="d-box green" data-tip="2 vCPU / 4 GB Fargate task. Handles ~2-3K RPS of redirect traffic, limited by network I/O to Redis/DynamoDB.">Each task handles <span class="d-metric throughput">2-3K RPS</span></div>
        <div class="d-box green" data-tip="No shared state: each task operates independently. Scale to 100+ tasks without coordination overhead.">Stateless: scale horizontally to 100+ tasks</div>
        <div class="d-label">Peak <span class="d-metric throughput">57K RPS</span> ÷ 2.5K/task = ~23 tasks</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Cache Layer</div>
      <div class="d-flow">
        <div class="d-number"><div class="d-number-value">9</div><div class="d-number-label">Redis nodes</div></div>
        <div class="d-number"><div class="d-number-value">39 GB</div><div class="d-number-label">total cache</div></div>
      </div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Cluster mode distributes keyspace across 3 primary nodes. 16,384 hash slots divided equally (~5,461 per shard).">ElastiCache cluster mode: 3 shards</div>
        <div class="d-box red" data-tip="Consistent hashing distributes keys. url:Ab3xK9 maps to exactly one shard via CRC16(key) mod 16384."><span class="d-metric size">16,384</span> hash slots across shards</div>
        <div class="d-box red" data-tip="2 replicas per primary = 6 additional nodes. Each replica handles reads. 9 total = 13 GB × 3 primaries = 39 GB total capacity.">2 read replicas per shard = 9 nodes</div>
        <div class="d-label">Total: <span class="d-metric size">39 GB</span> cache, <span class="d-metric throughput">300K+</span> reads/sec</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-multi-az",
		Title:       "Multi-AZ Deployment",
		Description: "Three availability zone deployment layout with ECS tasks, Redis, and DynamoDB.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box purple" data-tip="AWS ALB automatically distributes traffic across healthy targets in all AZs. Cross-AZ LB adds ~0.01ms but ensures even distribution.">ALB (cross-AZ load balancing)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">AZ-1a</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="Stateless Fargate task. 2 vCPU / 4 GB. Handles ~2-3K RPS. Auto-replaces on health check failure in &lt;30s."><span class="d-status active"></span>ECS Task 1</div>
          <div class="d-box green" data-tip="Stateless Fargate task. Scales independently of other AZs."><span class="d-status active"></span>ECS Task 2</div>
          <div class="d-box red" data-tip="Redis primary node: accepts all writes and replicates to replicas in AZ-1b and AZ-1c. Sync replication."><span class="d-status active"></span>Redis Primary</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">AZ-1b</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="Independent failure domain. AZ-1b outage leaves AZ-1a and AZ-1c serving 4 of 6 tasks."><span class="d-status active"></span>ECS Task 3</div>
          <div class="d-box green"><span class="d-status active"></span>ECS Task 4</div>
          <div class="d-box red" data-tip="Read replica in AZ-1b. Promotes to primary in 10-30s if AZ-1a Redis fails. Async replication lag &lt;1ms.">Redis Replica</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">AZ-1c</div>
        <div class="d-flow-v">
          <div class="d-box green"><span class="d-status active"></span>ECS Task 5</div>
          <div class="d-box green"><span class="d-status active"></span>ECS Task 6</div>
          <div class="d-box red" data-tip="Second read replica. Provides an additional failover target and spreads read load across 3 AZs.">Redis Replica</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="DynamoDB replicates all writes to 3 AZs before acknowledging. AZ failure is transparent to the application.">DynamoDB (multi-AZ by default, 3 replicas)</div>
  <div class="d-legend">
    <div class="d-legend-item"><div class="d-legend-color green"></div>Compute (stateless)</div>
    <div class="d-legend-item"><div class="d-legend-color red"></div>Cache (Redis)</div>
    <div class="d-legend-item"><div class="d-legend-color amber"></div>Storage (DynamoDB)</div>
    <div class="d-legend-item"><div class="d-legend-color purple"></div>Network (ALB)</div>
  </div>
  <div class="d-caption">Any single AZ failure leaves 4 ECS tasks + Redis replica serving traffic. Recovery is automatic: ALB stops routing to failed AZ within 30s.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-reliability-patterns",
		Title:       "Reliability Patterns Overview",
		Description: "Redundancy, resilience, and protection patterns for high availability.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Redundancy</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="6 ECS tasks across 3 AZs. ALB stops routing to failed AZ within 30s via health checks.">Multi-AZ: ECS across 3 AZs</div>
        <div class="d-label"><span class="d-metric latency">&lt;1 min</span> recovery</div>
        <div class="d-box green" data-tip="DynamoDB writes replicated synchronously to 3 AZs before ACK. AZ failure is completely transparent.">DB replication: DynamoDB 3 replicas</div>
        <div class="d-label"><span class="d-metric latency">&lt;1 sec</span> (transparent)</div>
        <div class="d-box green" data-tip="ElastiCache Multi-AZ: automatic failover promotes read replica to primary. DNS update propagates in ~30s.">Cache failover: Redis Multi-AZ</div>
        <div class="d-label"><span class="d-metric latency">10&#8212;30 sec</span> (promote replica)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Resilience</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Hystrix-style circuit breaker: open after 5 consecutive failures, half-open after 10s cooldown. Redirects still work from Redis during DB outage.">Circuit breaker on DynamoDB calls</div>
        <div class="d-label">Immediate fallback to cache (95% reads still served)</div>
        <div class="d-box blue" data-tip="Base delay 100ms, max 5s, multiplied by 2^attempt. Jitter ±20% prevents synchronized retries from all ECS tasks.">Retry with exponential backoff + jitter</div>
        <div class="d-label">Prevents thundering herd on recovery</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Protection</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Token bucket per user_id (auth) or IP (anon). 100 req/s per user. WAF rule at ALB level — never reaches ECS.">Rate limiting: <span class="d-metric throughput">100 req/s</span> per user</div>
        <div class="d-label">Immediate (HTTP 429)</div>
        <div class="d-box amber" data-tip="If DynamoDB is degraded, serve redirects from Redis only. Writes return 503. Users can still follow existing short links.">Graceful degradation: cache-only mode</div>
        <div class="d-label"><span class="d-metric throughput">95%</span> reads still work during DB outage</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-edge-cases",
		Title:       "Edge Cases & How We Handle Them",
		Description: "URL lifecycle edge cases and security layers for the shortener service.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">URL Lifecycle Edge Cases</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="410 signals permanent removal vs 404 temporary. Tells crawlers and browsers to remove from cache/index. DynamoDB TTL auto-deletes within 48h."><span class="d-status error"></span>Expired URL &#8594; HTTP 410 Gone (not 404)</div>
        <div class="d-box amber" data-tip="GetItem returns no item. Return 404 with JSON error body. Do NOT reveal whether code was valid-but-expired."><span class="d-status error"></span>Non-existent code &#8594; HTTP 404 Not Found</div>
        <div class="d-box amber" data-tip="PutItem conditional expression fails. Return 409 with suggestion to pick a different alias."><span class="d-status error"></span>Custom alias conflict &#8594; HTTP 409 Conflict</div>
        <div class="d-box red" data-tip="Google Safe Browsing API v4: synchronous check &lt;50ms. Async URIBL check follows. Reject with 422 and error message."><span class="d-status error"></span>Malicious URL &#8594; Safe Browsing API check &#8594; reject</div>
        <div class="d-box green" data-tip="By design: same long URL gets two different short codes. Prevents enumeration of who else shortened the same URL. Deduplication optional."><span class="d-status active"></span>Same long URL shortened twice &#8594; two different short codes (by design)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Security Layers</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Token bucket in Redis: 100 tokens/sec per IP. WAF rule enforced at ALB — abuse never reaches ECS or DynamoDB.">Rate limiting: <span class="d-metric throughput">100 req/s</span> per IP (token bucket)</div>
        <div class="d-box red" data-tip="Regex: must start with http/https, max 2048 chars, RFC 3986 compliant. Reject immediately before any DB call.">Input validation: URL format regex + length limit</div>
        <div class="d-box red" data-tip="TLS 1.3 terminated at ALB. HSTS header on all responses. Internal traffic (ALB → ECS) is also TLS via ACM cert.">HTTPS everywhere (TLS 1.3 via ALB)</div>
        <div class="d-box red" data-tip="AWS WAF managed rule: SQLi, XSS, known bad IPs. Applied to both /api/v1/urls writes and /{code} reads.">WAF: block SQL injection, XSS in custom aliases</div>
        <div class="d-box red" data-tip="Safe Browsing API v4. Synchronous check on URL creation (&lt;50ms). Nightly re-scan of all URLs against updated threat feeds.">Google Safe Browsing API: reject malicious URLs</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-analytics-pipeline",
		Title:       "Async Analytics Pipeline",
		Description: "Asynchronous click event pipeline from redirect through Kinesis to storage.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="301 redirect is returned to client first. Analytics call is async goroutine — zero impact on redirect latency."><span class="d-step">1</span>Redirect happens (GET /{code} &#8594; 301)</div>
  <div class="d-arrow-down">&#8595; fire-and-forget async</div>
  <div class="d-box green" data-tip="Non-blocking PutRecord call. If Kinesis is down, click data is lost — acceptable trade-off. Redirect latency is never affected."><span class="d-step">2</span>Kinesis Data Streams <span class="d-metric latency">&lt;1ms async</span> <span class="d-metric throughput">up to 1MB/s/shard</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Real-time</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="Lambda triggered every 60s. Aggregates click_count increments and batches UpdateItem calls to reduce DynamoDB WCU."><span class="d-step">3a</span>Lambda (batch aggregate per minute)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="UpdateItem with ADD expression: atomic increment. No race condition even at 100K clicks/sec on viral URL."><span class="d-step">4a</span>DynamoDB: atomic ADD click_count</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Archive</div>
      <div class="d-flow-v">
        <div class="d-box gray" data-tip="Kinesis Firehose buffers 5 min or 128 MB, then writes Snappy-compressed Parquet to S3. Cost: ~$0.029/GB."><span class="d-step">3b</span>Kinesis Firehose &#8594; S3 (parquet)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray" data-tip="SQL on S3 via Athena. Pay per query ($5/TB scanned). Parquet columnar format reduces scan cost by ~10x vs JSON."><span class="d-step">4b</span>Athena: ad-hoc analytics</div>
      </div>
    </div>
  </div>
  <div class="d-caption">Analytics is fully decoupled from the redirect path. A Kinesis outage loses click data but never affects redirect availability or latency.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-multi-region",
		Title:       "Multi-Region Architecture",
		Description: "Geo-distributed multi-region deployment with DynamoDB Global Tables.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box purple" data-tip="Route 53 latency-based routing resolves DNS to the region with lowest RTT for the user. Updates take ~60s to propagate.">Route 53 (Latency-based DNS &#8594; nearest region)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">us-east-1 (Americas)</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="400+ edge PoPs globally. us-east-1 CloudFront serves North and South America.">CloudFront</div>
          <div class="d-box green" data-tip="ECS Fargate tasks + r6g.large ElastiCache Redis. Stateless compute, regional cache for hot URLs.">ECS + Redis</div>
          <div class="d-box amber" data-tip="Primary replica. All writes from any region eventually replicate here via Global Tables. &lt;1s lag."><span class="d-status active"></span>DynamoDB (primary)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">eu-west-1 (Europe)</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="European edge PoPs serve EU traffic locally, complying with data residency if needed.">CloudFront</div>
          <div class="d-box green" data-tip="Same stack as us-east-1. Writes from EU users route to local DynamoDB replica then replicate globally.">ECS + Redis</div>
          <div class="d-box amber" data-tip="Active replica — accepts writes. Conflict resolution: last-writer-wins (DynamoDB Global Tables default).">DynamoDB (replica)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">ap-south-1 (Asia)</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="APAC edge PoPs. Covers India, Southeast Asia, Australia with low-latency redirects.">CloudFront</div>
          <div class="d-box green" data-tip="Handles APAC write traffic locally. Cache warm-up takes ~1h after regional deployment.">ECS + Redis</div>
          <div class="d-box amber" data-tip="Active replica. APAC users get &lt;10ms redirect latency from regional stack instead of 200ms+ cross-region.">DynamoDB (replica)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-legend">
    <div class="d-legend-item"><div class="d-legend-color purple"></div>CDN / DNS</div>
    <div class="d-legend-item"><div class="d-legend-color green"></div>Compute + Cache</div>
    <div class="d-legend-item"><div class="d-legend-color amber"></div>Database</div>
  </div>
  <div class="d-caption">DynamoDB Global Tables: multi-region multi-active with &lt;1s replication. A full regional failure triggers Route 53 health-check failover in ~60s.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-cost-scaling",
		Title:       "Cost Scaling Overview",
		Description: "Monthly infrastructure cost estimates at MVP, growth, and scale tiers.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">MVP (1K QPS)</div>
      <div class="d-flow">
        <div class="d-number"><div class="d-number-value">$400</div><div class="d-number-label">per month</div></div>
        <div class="d-number"><div class="d-number-value">1K</div><div class="d-number-label">peak QPS</div></div>
      </div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="2 Fargate tasks (~$60/mo) + t4g.micro Redis (~$15/mo) + DynamoDB on-demand (~$50/mo) + CloudFront 1TB (~$85/mo).">&#8776; $400/mo total</div>
        <div class="d-box gray" data-tip="2 Fargate tasks: 2 vCPU / 4 GB each, ~$30/mo per task. 1 Redis cache.r7g.small for hot-URL working set.">2 ECS tasks + 1 Redis node</div>
        <div class="d-box gray" data-tip="On-demand billing: $1.25/M reads + $1.25/M writes. No capacity planning needed at MVP scale.">DynamoDB on-demand</div>
        <div class="d-box gray" data-tip="1 TB egress at $0.085/GB = $85/mo. CloudFront serves bulk of reads. Origin bandwidth minimal.">CloudFront 1 TB</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Growth (10K QPS)</div>
      <div class="d-flow">
        <div class="d-number"><div class="d-number-value">$2,400</div><div class="d-number-label">per month</div></div>
        <div class="d-number"><div class="d-number-value">10K</div><div class="d-number-label">peak QPS</div></div>
      </div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="8 Fargate tasks (~$240/mo) + 3-node Redis (~$280/mo) + DynamoDB provisioned (~$400/mo) + CloudFront 10TB (~$850/mo).">&#8776; $2,400/mo total</div>
        <div class="d-box gray" data-tip="8 tasks × $30/mo each. 3-node Redis cluster: 1 primary + 2 replicas, r7g.large × 3 = ~$280/mo.">8 ECS tasks + 3-node Redis</div>
        <div class="d-box gray" data-tip="Switch from on-demand to provisioned capacity at 10K QPS to save ~40%. Use auto-scaling with 20% buffer.">DynamoDB provisioned</div>
        <div class="d-box gray" data-tip="10 TB egress at blended $0.085/GB = ~$850/mo. CloudFront traffic grows linearly with QPS.">CloudFront 10 TB</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale (100K QPS)</div>
      <div class="d-flow">
        <div class="d-number"><div class="d-number-value">$12K</div><div class="d-number-label">per month</div></div>
        <div class="d-number"><div class="d-number-value">100K</div><div class="d-number-label">peak QPS</div></div>
      </div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="40 Fargate tasks (~$1,200) + Redis cluster 9 nodes (~$2,500) + DynamoDB reserved (~$4,200) + CloudFront 50TB (~$4,200).">&#8776; $12,100/mo total</div>
        <div class="d-box gray" data-tip="40 tasks × $30/mo. Redis cluster mode: 3 shards × (1 primary + 2 replicas) = 9 nodes, 16,384 hash slots.">40 ECS tasks + Redis cluster</div>
        <div class="d-box gray" data-tip="1-year reserved capacity saves ~30% vs provisioned. At 100K QPS, write path alone needs ~10K WCU.">DynamoDB reserved</div>
        <div class="d-box amber" data-tip="50 TB × $0.085/GB = ~$4,250/mo. CDN egress becomes the dominant cost at scale. Consider private pricing.">CloudFront 50 TB <span class="d-metric throughput">35% of cost</span></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-security-layers",
		Title:       "Enhanced Security Layers",
		Description: "Five-layer defense: preview, bot detection, domain reputation, alias validation, rate limiting.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="All requests enter through AWS WAF at the ALB/CloudFront boundary before reaching any application code.">Incoming Request: POST /api/v1/urls or GET /{code}</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Write Path Security (URL Creation)</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="Token bucket per user_id/IP in Redis. Authenticated: 100/min. Anonymous: 10/min. Exceeding returns HTTP 429 immediately."><span class="d-step">1</span>Rate Limiting (<span class="d-metric throughput">100/min</span> auth, <span class="d-metric throughput">10/min</span> anon)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red" data-tip="Google reCAPTCHA v3 score. Anonymous users challenged after 5 URL creations. Score &lt;0.5 requires explicit checkbox CAPTCHA."><span class="d-step">2</span>CAPTCHA (anonymous after 5 creates)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber" data-tip="Regex ^[a-zA-Z0-9_-]{3,30}$, reserved word list (~500 terms), Aho-Corasick profanity filter, Unicode homoglyph normalisation."><span class="d-step">3</span>Custom Alias Validation (regex + reserved + profanity + homoglyph)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple" data-tip="Sync: Google Safe Browsing API v4 (&lt;50ms). Async: URIBL + WHOIS age check. New domains (&lt;7 days) flagged for review."><span class="d-step">4</span>Domain Reputation (Safe Browsing sync + URIBL/WHOIS async)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="attribute_not_exists(short_code) condition on PutItem. Final atomicity guarantee — even if two requests pass all checks simultaneously."><span class="d-step">5</span>DynamoDB Conditional Write</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Read Path Security (Redirect)</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="AWS WAF Bot Control managed rule group. Blocks known bot signatures, high-rate IPs, and Tor exit nodes at CloudFront/ALB."><span class="d-step">1</span>AWS WAF Bot Control at ALB</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber" data-tip="JA3/JA4 TLS fingerprint cross-checked with User-Agent. Chrome UA + curl TLS fingerprint = mismatched = bot flag."><span class="d-step">2</span>TLS Fingerprint + UA Cross-check</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple" data-tip="Composite risk score: bot score + URL domain reputation + click velocity. Score &gt;0.7 shows interstitial warning page."><span class="d-step">3</span>Risk Score Check (high risk &#8594; preview page)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="Normal 301/302 redirect served. Latency target p99 &lt;10ms. Analytics event fired asynchronously."><span class="d-step">4</span>Normal redirect (301/302) <span class="d-metric latency">&lt;10ms p99</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-flow-v">
    <div class="d-box gray" data-tip="Scheduled Lambda + Step Functions at 02:00 UTC. Re-checks all URLs against latest Safe Browsing + PhishTank feeds. Disables flagged URLs with 410 response.">Nightly Batch: Re-scan all URLs against updated threat feeds &#8594; disable flagged URLs (410 Gone)</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-bot-detection",
		Title:       "Bot Detection Pipeline",
		Description: "Multi-signal bot detection from TLS fingerprinting to behavioral analysis.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="All read and write requests are scored before processing. Score computation adds &lt;1ms latency via in-memory lookup.">Incoming HTTP Request</div>
  <div class="d-arrow-down">&#8595; extract signals</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Signal 1: TLS Fingerprint</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="JA3 hashes the TLS ClientHello fields (cipher suites, extensions, elliptic curves) into 32-char MD5. Unique per client library.">JA3/JA4 hash of TLS handshake</div>
          <div class="d-box gray" data-tip="AWS WAF maintains updated lists of known bot TLS fingerprints: Puppeteer, Playwright, Selenium WebDriver, curl, wget.">Known bot fingerprints: Puppeteer, Selenium, curl, wget</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Signal 2: UA Consistency</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="Real Chrome browser has a specific TLS fingerprint that differs from curl. Spoofing UA string alone doesn't change TLS fingerprint.">Cross-check UA string with TLS fingerprint</div>
          <div class="d-box gray" data-tip="Chrome/120 UA + curl/7.x TLS fingerprint = identity mismatch = confidence bot. Score incremented by 0.4.">Chrome UA + curl TLS = bot (mismatch)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Signal 3: Behavioral</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="Sliding window rate tracked in Redis. Uniform inter-request timing (robots are precise; humans are not) raises score.">Creation velocity + timing patterns</div>
          <div class="d-box gray" data-tip="&gt;10 URL creates/min with &lt;5% timing variance + no mouse events = high-confidence bot.">&#62;10 URLs/min with identical patterns = bot</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; compute</div>
  <div class="d-box indigo" data-tip="Weighted sum of all signals. Weights tuned to minimize false positives (&lt;0.1% of legitimate users challenged).">Bot Score: 0.0 (human) &#8594; 1.0 (bot)</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Score &lt; 0.3</div>
      <div class="d-box green" data-tip="Normal processing path. Majority of legitimate traffic falls here. No additional latency."><span class="d-status active"></span>Allow (normal flow)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Score 0.3 &#8212; 0.7</div>
      <div class="d-box amber" data-tip="reCAPTCHA v3 challenge injected. If score confirmed &gt;0.5, serve checkbox CAPTCHA. Legitimate users pass; bots fail.">Challenge (CAPTCHA)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Score &gt; 0.7</div>
      <div class="d-box red" data-tip="HTTP 429 Too Many Requests with Retry-After header. IP added to block list for 1h in WAF."><span class="d-status error"></span>Block (HTTP 429)</div>
    </div>
  </div>
  <div class="d-caption">Bot scoring adds &lt;1ms per request via Redis lookup. False positive rate target: &lt;0.1% of legitimate users challenged.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-alias-validation",
		Title:       "Custom Alias Validation Flow",
		Description: "Five-step validation pipeline for custom short URL aliases.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="User-supplied alias from PUT /api/v1/urls?alias=my-link. Checked before any DB write.">Input: custom_alias = "my-link"</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="^[a-zA-Z0-9_-]{3,30}$ — alphanumeric, hyphen, underscore only. Min 3 chars to prevent trivial aliases. Max 30 to fit in URL."><span class="d-step">1</span>Regex: ^[a-zA-Z0-9_-]{3,30}$ &#8594; PASS</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="~500 reserved words: admin, api, www, cdn, health, status, login, logout, etc. Hash-set O(1) lookup."><span class="d-step">2</span>Reserved Words: admin, api, www, cdn... &#8594; PASS</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="Aho-Corasick multi-pattern matching: O(N) scan across ~5K profanity terms regardless of alias length. Sub-millisecond."><span class="d-step">3</span>Profanity Filter (Aho-Corasick, <span class="d-metric size">~5K terms</span>) &#8594; PASS</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="Unicode normalisation + visual lookalike map. Converts 'paypa1' to 'paypal', 'g00gle' to 'google'. Prevents brand impersonation."><span class="d-step">4</span>Homoglyph Check (paypa1 &#8776; paypal?) &#8594; PASS</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="attribute_not_exists(short_code) condition. Atomic — no race between two users claiming the same alias simultaneously."><span class="d-step">5</span>DynamoDB: PutItem(condition: attribute_not_exists) &#8594; PASS or 409</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label"><span class="d-status active"></span>All pass</div>
      <div class="d-box green" data-tip="201 Created with Location header pointing to the new short URL. Alias reserved in DynamoDB atomically.">201 Created: short.ly/my-link</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label"><span class="d-status error"></span>Regex fails</div>
      <div class="d-box red" data-tip="400 Bad Request with detailed error: 'Alias contains invalid characters. Use letters, numbers, hyphens, underscores only.'">400 Bad Request</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label"><span class="d-status error"></span>Reserved/profanity</div>
      <div class="d-box red" data-tip="422 Unprocessable Entity. Generic message to avoid revealing the reserved word list to attackers.">422 Unprocessable</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label"><span class="d-status error"></span>Already taken</div>
      <div class="d-box red" data-tip="409 Conflict with suggestion: 'Try my-link-2 or my-link-2026'. Never reveal who owns the alias.">409 Conflict</div>
    </div>
  </div>
  <div class="d-caption">All validation steps run in-memory except Step 5 (DynamoDB). Total latency for a valid alias: &lt;5ms.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-qr-code",
		Title:       "QR Code Generation Architecture",
		Description: "Lambda@Edge QR code generation with CloudFront caching.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="QR endpoint accepts size (px) and format (png/svg). Cache key includes all query params so size variants are cached independently."><span class="d-step">1</span>Client: GET /api/v1/urls/Ab3xK9/qr?size=300&amp;format=png</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Request routed to nearest CloudFront edge PoP. Cache key: {short_code}+{size}+{format}. TTL 24h — QR content never changes."><span class="d-step">2</span>CloudFront Edge PoP</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label"><span class="d-status active"></span>Cache HIT (24h TTL)</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Binary PNG/SVG served from edge cache. No Lambda invocation needed. Scales to millions of QR views for viral URLs.">Return cached QR image <span class="d-metric latency">&lt;5ms</span></div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Cache MISS</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Lambda@Edge runs in the nearest edge PoP, not us-east-1. Reduces cold-start latency from ~200ms to ~50ms. Node.js qrcode library."><span class="d-step">3</span>Lambda@Edge: generate QR code <span class="d-metric latency">~50ms cold</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray" data-tip="QR encodes the short URL itself (e.g., https://short.ly/Ab3xK9), not the destination. Scanning always resolves through the shortener."><span class="d-step">4</span>Encode: https://short.ly/Ab3xK9</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Response includes Cache-Control: max-age=86400. CloudFront caches at edge for 24h. Subsequent requests served instantly."><span class="d-step">5</span>Return PNG/SVG + cache at edge</div>
      </div>
    </div>
  </div>
  <div class="d-row">
    <div class="d-box indigo" data-tip="150×150px PNG. Suitable for business cards and receipts. QR version 3, ~7% error correction.">150px: business cards</div>
    <div class="d-box indigo" data-tip="300×300px PNG. Restaurant menus, flyers, product labels. QR version 5, 15% error correction.">300px: posters, menus</div>
    <div class="d-box indigo" data-tip="600×600px SVG (vector). Billboards and large prints. SVG scales infinitely. 30% error correction for damaged surfaces.">600px: billboards</div>
  </div>
  <div class="d-caption">Cost: $0.00006 per 10K Lambda@Edge requests + CloudFront transfer. At 1M QR views/day, &lt;$10/mo after warm cache is established.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-analytics-schema",
		Title:       "Enhanced Analytics Pipeline & Schema",
		Description: "Click event enrichment and dimensional analytics storage in DynamoDB.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="301 redirect returned to client first. Analytics goroutine fires concurrently — redirect latency is never affected."><span class="d-step">1</span>Redirect: GET /Ab3xK9 &#8594; 301</div>
  <div class="d-arrow-down">&#8595; fire-and-forget</div>
  <div class="d-box green" data-tip="Raw click event: {short_code, timestamp_ms, ip, user_agent, referrer, accept_language}. ~300 bytes per event. 2 Kinesis shards = 2MB/s capacity."><span class="d-step">2</span>Kinesis Data Streams (raw click event) <span class="d-metric size">~300B/event</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Lambda trigger: batch size 100 records, max window 60s. Parallelism factor 1 per shard. Enriches raw events before aggregation."><span class="d-step">3</span>Lambda Enrichment (per batch, every 60s)</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-box gray" data-tip="MaxMind GeoIP2 Precision City database. IP → country, region, city, lat/lng. ~2ms lookup. GDPR: truncate last octet in EU.">IP &#8594; Geo (MaxMind GeoIP2)</div>
        <div class="d-box gray" data-tip="ua-parser library: User-Agent string → {browser, browser_version, os, device_type}. Cached — same UA seen millions of times.">UA &#8594; Device/Browser parsing</div>
        <div class="d-box gray" data-tip="Extract eTLD+1 from Referer header. Maps 'https://t.co/xyz' → 'twitter.com'. Identifies traffic sources.">Referrer &#8594; domain extraction</div>
        <div class="d-box gray" data-tip="Reuse bot score from write path if already computed. Otherwise compute from IP + UA signals. Score &gt;0.7 = exclude from analytics.">Bot score computation</div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Drop events where bot_score &gt; 0.7 before aggregation. Prevents inflated click counts from scrapers and monitoring tools."><span class="d-step">4</span>Filter bots (score &gt; 0.7)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Group by {short_code, date, dimension_type, dimension_value}. Sum click_count per group. Reduces N events to K dimension rows."><span class="d-step">5</span>Aggregate by dimensions</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="BatchWriteItem: up to 25 UpdateItem operations per call. ADD expression: atomic increment. No read-modify-write needed."><span class="d-step">6</span>Batch write to DynamoDB</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-entity">
    <div class="d-entity-header indigo">click_analytics table</div>
    <div class="d-entity-body">
      <div class="pk">PK: short_code#date (e.g., Ab3xK9#2024-03-05)</div>
      <div class="pk">SK: dimension#value (e.g., country#IN, device#mobile)</div>
      <div>click_count NUMBER (atomic ADD)</div>
      <div>unique_ips NUMBER (HyperLogLog approximation)</div>
    </div>
  </div>
  <div class="d-caption">Composite PK isolates analytics writes per URL per day: a viral URL's analytics never contend with other URLs. Dimension SK allows querying "clicks by country" or "clicks by device" with a single Query call.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-monitoring-slos",
		Title:       "Monitoring & SLO Dashboard",
		Description: "Service level objectives, key metrics, and alert thresholds.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">SLOs</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Multi-AZ + multi-region deployment. Error budget: 52 min/yr. Alert when monthly burn rate exceeds 5%.">Availability: <span class="d-metric throughput">99.99%</span> (<span class="d-metric latency">52 min/yr</span>)</div>
        <div class="d-box green" data-tip="p50 &lt;2ms (cache hit), p95 &lt;5ms, p99 &lt;10ms (DB hit). Measured at ALB, excludes client network.">Redirect p99: <span class="d-metric latency">&lt;10ms</span></div>
        <div class="d-box green" data-tip="Write SLO is looser than read: URL creation can retry. 0.1% failure budget = 100K failed creates/day at 100M/day scale.">Write Success: <span class="d-metric throughput">99.9%</span></div>
        <div class="d-box green" data-tip="CloudFront alone absorbs 60%. Combined with Redis: &gt;99% of all traffic served from cache. DB sees &lt;1% of total RPS.">Cache Hit: <span class="d-metric throughput">&gt;60%</span> (CloudFront)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Real-time Dashboard</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="CloudWatch metric: RequestCount on ALB. Alert if drops &gt;20% from rolling 1h average (possible outage or traffic shift).">RPS (requests per second)</div>
        <div class="d-box blue" data-tip="CloudWatch: TargetResponseTime histogram. P99 SLO alert: &gt;10ms for 5 consecutive minutes triggers PagerDuty.">Latency percentiles (p50/p95/p99)</div>
        <div class="d-box blue" data-tip="4xx alert: &gt;5% rate (client errors spike = possible attack). 5xx alert: &gt;0.1% rate (server errors = immediate page).">Error rate (4xx, 5xx)</div>
        <div class="d-box blue" data-tip="ElastiCache: CacheHits/(CacheHits+CacheMisses). Alert if falls below 50% — may indicate cache eviction or cold start.">Cache hit ratio</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Business Dashboard</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="Kinesis analytics → DynamoDB. Trailing 7-day moving average. Alert if drops &gt;30% (possible API abuse or market event).">New URLs/day</div>
        <div class="d-box purple" data-tip="Aggregate from click_analytics table via Athena query. Top 100 domains by click volume. Updated daily.">Top domains shortened</div>
        <div class="d-box purple" data-tip="MaxMind GeoIP data aggregated by Kinesis Lambda. Shows heatmap of click origins. Used for capacity planning.">Geo distribution</div>
        <div class="d-box purple" data-tip="DAU from unique user_id values in DynamoDB. Anonymous users counted by hashed IP. 30-day retention window.">Active users</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Cost Dashboard</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="AWS Cost Explorer API. Per-service daily spend: ECS, ElastiCache, DynamoDB, CloudFront, Kinesis. Alert on 20% day-over-day spike.">Daily spend by service</div>
        <div class="d-box amber" data-tip="DynamoDB consumed RCU/WCU vs provisioned. Alert if consumed &gt;80% of provisioned — time to scale or switch billing mode.">DynamoDB RCU/WCU usage</div>
        <div class="d-box amber" data-tip="CloudFront DataTransfer-Out-Bytes metric. At scale, egress is 35% of total cost. Alert if exceeds monthly budget.">CloudFront transfer (GB)</div>
        <div class="d-box amber" data-tip="AWS Budgets alert at 80% and 100% of monthly budget. MoM trend chart for forecasting Reserved Instance purchases.">Month-over-month trend</div>
      </div>
    </div>
  </div>
</div>
<div class="d-flow-v">
  <div class="d-box red" data-tip="X-Ray trace: segments for each hop (ALB → ECS → Redis → DynamoDB). Identifies which layer adds latency. Sample rate: 5% in production to control cost.">Distributed Tracing: AWS X-Ray &#8212; ALB &#8594; ECS &#8594; Redis &#8594; DynamoDB (end-to-end latency breakdown)</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-db-access-patterns",
		Title:       "DynamoDB Access Patterns & Hot Partition Handling",
		Description: "Access patterns, GSI design, and hot partition mitigation for the URL table.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Primary Table: urls</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Single-item lookup by partition key. O(1) regardless of table size. &lt;1ms in same AZ, ~5ms cross-AZ."><span class="d-step">R</span>GetItem(PK=short_code) &#8594; O(1) <span class="d-metric latency">&lt;5ms</span></div>
        <div class="d-box green" data-tip="Condition: attribute_not_exists(short_code). Atomic — prevents race between two concurrent creates of same code."><span class="d-step">W</span>PutItem(condition: attribute_not_exists) &#8594; O(1)</div>
        <div class="d-box blue" data-tip="Hard delete by PK. For soft delete: set expires_at=now() and let TTL clean up. Invalidate CloudFront cache via API."><span class="d-step">D</span>DeleteItem(PK=short_code) &#8594; O(1)</div>
        <div class="d-box amber" data-tip="DynamoDB TTL scanner runs every 48h. Deletes items with expires_at &lt; now(). Zero RCU/WCU consumed for TTL deletions.">TTL auto-deletes expired URLs &#8594; zero cost</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">GSI-1: user-urls-index</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="GSI partition key: user_id. Sort key: created_at (epoch seconds). Enables time-sorted listing of a user's URLs.">PK: user_id | SK: created_at</div>
        <div class="d-box purple" data-tip="ALL projection copies every attribute to the GSI. Avoids extra GetItem calls. ~2× storage cost vs KEYS_ONLY — acceptable at URL sizes.">Projection: ALL (full item copy)</div>
        <div class="d-box purple" data-tip="Query(KeyConditionExpression=user_id=:uid). ExclusiveStartKey for cursor pagination. ScanIndexForward=false for newest-first.">Query: "List my URLs" &#8594; paginated</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Hot Partition Handling</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="10M clicks/day = ~116 RPS on a single short_code partition key. Without caching, would hit DynamoDB partition directly.">Viral URL: <span class="d-metric throughput">10M clicks/day</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="95% cache hit rate means only 5% of 116 RPS = ~6 RPS reach DynamoDB. Well within single-partition limit.">Cache hit rate 95%+ &#8594; only <span class="d-metric throughput">~500K DB reads/day</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Each DynamoDB partition supports 3,000 RCU/s and 1,000 WCU/s. 6 RPS << 3,000 RCU. No throttling risk.">DynamoDB: <span class="d-metric throughput">3,000 RCU</span> per partition = <span class="d-metric throughput">3M reads/sec</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue" data-tip="DynamoDB adaptive capacity automatically boosts throughput for hot partitions, redistributing capacity from cold partitions within minutes.">Adaptive capacity redistributes within minutes</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Analytics: Separate Table</div>
      <div class="d-flow-v">
        <div class="d-box gray" data-tip="click_analytics table has PK=short_code#date, separate from urls table. Viral URL analytics writes isolated to their own partition.">PK: short_code#date &#8594; isolates analytics writes</div>
        <div class="d-box gray" data-tip="Without isolation, 10M click_count ADDs/day on the urls table would throttle the redirect path by consuming partition WCU budget.">Prevents viral URL analytics from affecting redirects</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-sub-problems",
		Title:       "URL Shortener Sub-Problems & Building Blocks",
		Description: "Key sub-problems and reusable building blocks referenced by the URL shortener.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-subproblem">
    <div class="d-subproblem-icon indigo">&#128273;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Key Generation (Base62 Encoding)</div>
      <div class="d-subproblem-desc">62^7 = 3.5T unique codes. KGS pre-generates batches, zero collisions.</div>
    </div>
    <div class="d-subproblem-link">&#8594; Base62 Encoding</div>
  </div>
  <div class="d-subproblem">
    <div class="d-subproblem-icon red">&#9889;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Rate Limiting API Abuse</div>
      <div class="d-subproblem-desc">Token bucket at ALB: 100 req/s per user. Prevents URL spam and DDoS.</div>
    </div>
    <div class="d-subproblem-link">&#8594; Token Bucket / Rate Limiter</div>
  </div>
  <div class="d-subproblem">
    <div class="d-subproblem-icon purple">&#127760;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">CDN Caching (CloudFront)</div>
      <div class="d-subproblem-desc">301 cached at 400+ edge PoPs. 60%+ reads never reach origin.</div>
    </div>
    <div class="d-subproblem-link">&#8594; CDN / CloudFront</div>
  </div>
  <div class="d-subproblem">
    <div class="d-subproblem-icon green">&#128203;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Consistent Hashing (if SQL sharding)</div>
      <div class="d-subproblem-desc">Hash ring with virtual nodes. Add/remove shards moves only K/N keys.</div>
    </div>
    <div class="d-subproblem-link">&#8594; Consistent Hashing</div>
  </div>
  <div class="d-subproblem">
    <div class="d-subproblem-icon amber">&#128200;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Bloom Filter (Optional Collision Check)</div>
      <div class="d-subproblem-desc">O(1) probabilistic check before DB write. Reduces unnecessary lookups.</div>
    </div>
    <div class="d-subproblem-link">&#8594; Bloom Filter</div>
  </div>
</div>`,
	})
}

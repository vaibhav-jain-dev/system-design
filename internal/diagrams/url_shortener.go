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
        <div class="d-box blue">Browser</div>
        <div class="d-arrow-down">&#8595; cached locally</div>
        <div class="d-box green"><span class="d-status active"></span>&#10003; Direct to destination <span class="d-metric latency">0ms</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">302 Found (Temporary)</div>
      <div class="d-flow-v">
        <div class="d-box red">Browser does NOT cache</div>
        <div class="d-box red">Every click hits your servers</div>
        <div class="d-box red">2&#8212;5&#215; more origin traffic</div>
        <div class="d-box green">Full analytics on every click</div>
        <div class="d-box green">Can change destination anytime</div>
      </div>
      <div class="d-flow-v">
        <div class="d-label">Traffic flow (2nd click by same user):</div>
        <div class="d-box blue">Browser</div>
        <div class="d-arrow-down">&#8595; no cache</div>
        <div class="d-box purple">Your server (5&#8212;50ms)</div>
        <div class="d-arrow-down">&#8595; redirect</div>
        <div class="d-box amber">Destination</div>
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
        <div class="d-box green" data-tip="52 minutes downtime per year. Multi-AZ deployment, auto-scaling, health checks.">Availability: 99.99% <span class="d-metric throughput">52 min/yr</span></div>
        <div class="d-box green" data-tip="p50 < 5ms (cache hit), p99 < 10ms (cache miss + DB). CDN hits are < 2ms.">Redirect latency: &lt; 10ms p99</div>
        <div class="d-box blue" data-tip="DynamoDB 11 9's durability. Point-in-time recovery enabled. Cross-region replication for DR.">Durability: never lose a URL mapping</div>
        <div class="d-box amber" data-tip="Cache TTL 1h means updates take up to 1h to propagate. Acceptable for URL shortener — URLs rarely change.">Consistency: eventual OK (cache + DB sync)</div>
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
  </div>
  <div class="d-row">
    <div class="d-box green">All options &#8594; Base62 encode &#8594; 7-character short code</div>
    <div class="d-box purple">43 bits &#8776; 8.8 trillion values &#8811; 3.5T (62^7)</div>
  </div>
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
        <div class="d-box blue">Global atomic counter (DB or Zookeeper)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">counter++ &#8594; Base62 encode &#8594; "Ab3xK9"</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Store mapping in DynamoDB</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Approach 2: MD5/SHA + Truncate</div>
      <div class="d-flow-v">
        <div class="d-box blue">Hash long URL: MD5("https://...")</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Take first 43 bits &#8594; Base62(7 chars)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Check collision &#8594; retry if exists</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Approach 3: Key Generation Service</div>
      <div class="d-flow-v">
        <div class="d-box blue">KGS pre-generates batches of unique keys</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">App server requests batch of 1000 keys</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Assign next key from local batch</div>
      </div>
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
  <div class="d-box blue">Input: long_url = "https://example.com/very/long/path"</div>
  <div class="d-arrow-down">&#8595; MD5 hash</div>
  <div class="d-box purple">MD5 &#8594; "e4d909c290d0fb1ca068ffaddf22cbd0" (128 bits)</div>
  <div class="d-arrow-down">&#8595; take first 43 bits</div>
  <div class="d-box green">Truncate &#8594; Base62 encode &#8594; "Ab3xK9p" (7 chars)</div>
  <div class="d-arrow-down">&#8595; check DB</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">No collision</div>
      <div class="d-flow-v">
        <div class="d-box green">&#10003; Store mapping in DB</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Return short URL to client</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Collision detected</div>
      <div class="d-flow-v">
        <div class="d-box red">&#215; short_code already exists</div>
        <div class="d-arrow-down">&#8595; retry strategy</div>
        <div class="d-box amber">Append counter: hash(url + "1")</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Re-truncate &#8594; Base62 &#8594; new code</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Check DB again (max 5 retries)</div>
      </div>
    </div>
  </div>
  <div class="d-label">Birthday paradox: at 100B URLs, collision probability &#8776; 0.14% per attempt</div>
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
		HTML: `<div class="d-query-result">
  <table>
    <tr><th>Operation</th><th>Key</th><th>Result</th><th>Latency</th><th>Source</th></tr>
    <tr><td>GET /Ab3xK9</td><td>url:Ab3xK9</td><td>https://example.com/very/long/path?utm=...</td><td>0.3ms</td><td>ElastiCache HIT</td></tr>
    <tr><td>GET /Xz9mQ2</td><td>url:Xz9mQ2</td><td>https://docs.google.com/spreadsheets/d/...</td><td>4.8ms</td><td>DynamoDB (cache miss)</td></tr>
    <tr><td>GET /expired1</td><td>url:expired1</td><td>HTTP 410 Gone</td><td>5.1ms</td><td>DynamoDB (TTL expired)</td></tr>
    <tr><td>GET /notfound</td><td>url:notfound</td><td>HTTP 404 Not Found</td><td>4.2ms</td><td>DynamoDB (no item)</td></tr>
  </table>
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
      <div class="d-flow-v">
        <div class="d-box green">ECS auto-scaling on CPU (target 60%)</div>
        <div class="d-box green">Each task handles 2-3K RPS</div>
        <div class="d-box green">Stateless: scale horizontally to 100+ tasks</div>
        <div class="d-label">Peak 57K RPS &#247; 2.5K/task = ~23 tasks</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Cache Layer</div>
      <div class="d-flow-v">
        <div class="d-box red">ElastiCache cluster mode: 3 shards</div>
        <div class="d-box red">16,384 hash slots across shards</div>
        <div class="d-box red">2 read replicas per shard = 9 nodes</div>
        <div class="d-label">Total: 39 GB cache, 300K+ reads/sec</div>
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
  <div class="d-box purple">ALB (cross-AZ load balancing)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">AZ-1a</div>
        <div class="d-flow-v">
          <div class="d-box green">ECS Task 1</div>
          <div class="d-box green">ECS Task 2</div>
          <div class="d-box red">Redis Primary</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">AZ-1b</div>
        <div class="d-flow-v">
          <div class="d-box green">ECS Task 3</div>
          <div class="d-box green">ECS Task 4</div>
          <div class="d-box red">Redis Replica</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">AZ-1c</div>
        <div class="d-flow-v">
          <div class="d-box green">ECS Task 5</div>
          <div class="d-box green">ECS Task 6</div>
          <div class="d-box red">Redis Replica</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber">DynamoDB (multi-AZ by default, 3 replicas)</div>
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
        <div class="d-box green">Multi-AZ: ECS across 3 AZs</div>
        <div class="d-label">&lt; 1 min recovery</div>
        <div class="d-box green">DB replication: DynamoDB 3 replicas</div>
        <div class="d-label">&lt; 1 sec (transparent)</div>
        <div class="d-box green">Cache failover: Redis Multi-AZ</div>
        <div class="d-label">10&#8212;30 sec (promote replica)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Resilience</div>
      <div class="d-flow-v">
        <div class="d-box blue">Circuit breaker on DynamoDB calls</div>
        <div class="d-label">Immediate fallback to cache</div>
        <div class="d-box blue">Retry with exponential backoff + jitter</div>
        <div class="d-label">Prevents thundering herd</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Protection</div>
      <div class="d-flow-v">
        <div class="d-box red">Rate limiting: 100 req/s per user</div>
        <div class="d-label">Immediate (HTTP 429)</div>
        <div class="d-box amber">Graceful degradation: cache-only mode</div>
        <div class="d-label">95% reads still work during DB outage</div>
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
        <div class="d-box amber">Expired URL &#8594; HTTP 410 Gone (not 404)</div>
        <div class="d-box amber">Non-existent code &#8594; HTTP 404 Not Found</div>
        <div class="d-box amber">Custom alias conflict &#8594; HTTP 409 Conflict</div>
        <div class="d-box amber">Malicious URL &#8594; Safe Browsing API check &#8594; reject</div>
        <div class="d-box amber">Same long URL shortened twice &#8594; two different short codes (by design)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Security Layers</div>
      <div class="d-flow-v">
        <div class="d-box red">Rate limiting: 100 req/s per IP (token bucket)</div>
        <div class="d-box red">Input validation: URL format regex + length limit</div>
        <div class="d-box red">HTTPS everywhere (TLS 1.3 via ALB)</div>
        <div class="d-box red">WAF: block SQL injection, XSS in custom aliases</div>
        <div class="d-box red">Google Safe Browsing API: reject malicious URLs</div>
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
  <div class="d-box blue"><span class="d-step">1</span>Redirect happens (GET /{code} &#8594; 301)</div>
  <div class="d-arrow-down">&#8595; fire-and-forget async</div>
  <div class="d-box green" data-tip="Non-blocking PutRecord call. If Kinesis is down, click data is lost (acceptable trade-off for redirect latency)."><span class="d-step">2</span>Kinesis Data Streams <span class="d-metric latency">&lt;1ms async</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Real-time</div>
      <div class="d-flow-v">
        <div class="d-box purple">Lambda (batch aggregate per minute)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">DynamoDB: atomic ADD click_count</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Archive</div>
      <div class="d-flow-v">
        <div class="d-box gray">Kinesis Firehose &#8594; S3 (parquet)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray">Athena: ad-hoc analytics</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-multi-region",
		Title:       "Multi-Region Architecture",
		Description: "Geo-distributed multi-region deployment with DynamoDB Global Tables.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box purple">Route 53 (Latency-based DNS &#8594; nearest region)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">us-east-1 (Americas)</div>
        <div class="d-flow-v">
          <div class="d-box purple">CloudFront</div>
          <div class="d-box green">ECS + Redis</div>
          <div class="d-box amber">DynamoDB (primary)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">eu-west-1 (Europe)</div>
        <div class="d-flow-v">
          <div class="d-box purple">CloudFront</div>
          <div class="d-box green">ECS + Redis</div>
          <div class="d-box amber">DynamoDB (replica)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">ap-south-1 (Asia)</div>
        <div class="d-flow-v">
          <div class="d-box purple">CloudFront</div>
          <div class="d-box green">ECS + Redis</div>
          <div class="d-box amber">DynamoDB (replica)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-label">DynamoDB Global Tables: multi-region, multi-active, &lt;1s replication</div>
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
      <div class="d-flow-v">
        <div class="d-box green">&#8776; $400/mo total</div>
        <div class="d-box gray">2 ECS tasks + 1 Redis node</div>
        <div class="d-box gray">DynamoDB on-demand</div>
        <div class="d-box gray">CloudFront 1 TB</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Growth (10K QPS)</div>
      <div class="d-flow-v">
        <div class="d-box blue">&#8776; $2,400/mo total</div>
        <div class="d-box gray">8 ECS tasks + 3-node Redis</div>
        <div class="d-box gray">DynamoDB provisioned</div>
        <div class="d-box gray">CloudFront 10 TB</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale (100K QPS)</div>
      <div class="d-flow-v">
        <div class="d-box purple">&#8776; $12,100/mo total</div>
        <div class="d-box gray">40 ECS tasks + Redis cluster</div>
        <div class="d-box gray">DynamoDB reserved</div>
        <div class="d-box amber">CloudFront 50 TB (35% of cost)</div>
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
  <div class="d-box blue">Incoming Request: POST /api/v1/urls or GET /{code}</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Write Path Security (URL Creation)</div>
        <div class="d-flow-v">
          <div class="d-box red">Layer 1: Rate Limiting (100/min auth, 10/min anon)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red">Layer 2: CAPTCHA (anonymous after 5 creates)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Layer 3: Custom Alias Validation (regex + reserved + profanity + homoglyph)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Layer 4: Domain Reputation (Safe Browsing sync + URIBL/WHOIS async)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Layer 5: DynamoDB Conditional Write</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Read Path Security (Redirect)</div>
        <div class="d-flow-v">
          <div class="d-box red">AWS WAF Bot Control at ALB</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">TLS Fingerprint + UA Cross-check</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Risk Score Check (high risk &#8594; preview page)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Normal redirect (301/302)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-flow-v">
    <div class="d-box gray">Nightly Batch: Re-scan all URLs against updated threat feeds &#8594; disable flagged URLs (410 Gone)</div>
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
  <div class="d-box blue">Incoming HTTP Request</div>
  <div class="d-arrow-down">&#8595; extract signals</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Signal 1: TLS Fingerprint</div>
        <div class="d-flow-v">
          <div class="d-box purple">JA3/JA4 hash of TLS handshake</div>
          <div class="d-box gray">Known bot fingerprints: Puppeteer, Selenium, curl, wget</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Signal 2: UA Consistency</div>
        <div class="d-flow-v">
          <div class="d-box purple">Cross-check UA string with TLS fingerprint</div>
          <div class="d-box gray">Chrome UA + curl TLS = bot (mismatch)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Signal 3: Behavioral</div>
        <div class="d-flow-v">
          <div class="d-box purple">Creation velocity + timing patterns</div>
          <div class="d-box gray">&gt;10 URLs/min with identical patterns = bot</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; compute</div>
  <div class="d-box indigo">Bot Score: 0.0 (human) &#8594; 1.0 (bot)</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Score &lt; 0.3</div>
      <div class="d-box green">Allow (normal flow)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Score 0.3 &#8212; 0.7</div>
      <div class="d-box amber">Challenge (CAPTCHA)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Score &gt; 0.7</div>
      <div class="d-box red">Block (HTTP 429)</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-alias-validation",
		Title:       "Custom Alias Validation Flow",
		Description: "Five-step validation pipeline for custom short URL aliases.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Input: custom_alias = "my-link"</div>
  <div class="d-arrow-down">&#8595; Step 1</div>
  <div class="d-box green">Regex: ^[a-zA-Z0-9_-]{3,30}$ &#8594; PASS</div>
  <div class="d-arrow-down">&#8595; Step 2</div>
  <div class="d-box green">Reserved Words: admin, api, www, cdn... &#8594; PASS</div>
  <div class="d-arrow-down">&#8595; Step 3</div>
  <div class="d-box green">Profanity Filter (Aho-Corasick, ~5K terms) &#8594; PASS</div>
  <div class="d-arrow-down">&#8595; Step 4</div>
  <div class="d-box green">Homoglyph Check (paypa1 &#8776; paypal?) &#8594; PASS</div>
  <div class="d-arrow-down">&#8595; Step 5</div>
  <div class="d-box amber">DynamoDB: PutItem(condition: attribute_not_exists) &#8594; PASS or 409</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">All pass</div>
      <div class="d-box green">201 Created: short.ly/my-link</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Regex fails</div>
      <div class="d-box red">400 Bad Request</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Reserved/profanity</div>
      <div class="d-box red">422 Unprocessable</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Already taken</div>
      <div class="d-box red">409 Conflict</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-qr-code",
		Title:       "QR Code Generation Architecture",
		Description: "Lambda@Edge QR code generation with CloudFront caching.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Client: GET /api/v1/urls/Ab3xK9/qr?size=300&format=png</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">CloudFront Edge PoP</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Cache HIT (24h TTL)</div>
      <div class="d-flow-v">
        <div class="d-box green">Return cached QR image (&lt;5ms)</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Cache MISS</div>
      <div class="d-flow-v">
        <div class="d-box amber">Lambda@Edge: generate QR code</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray">Encode: https://short.ly/Ab3xK9</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Return PNG/SVG + cache at edge</div>
      </div>
    </div>
  </div>
  <div class="d-row">
    <div class="d-box indigo">150px: business cards</div>
    <div class="d-box indigo">300px: posters, menus</div>
    <div class="d-box indigo">600px: billboards</div>
  </div>
  <div class="d-label">Cost: $0.00006 per 10K requests (Lambda@Edge) + CloudFront transfer</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "url-analytics-schema",
		Title:       "Enhanced Analytics Pipeline & Schema",
		Description: "Click event enrichment and dimensional analytics storage in DynamoDB.",
		ContentFile: "problems/url-shortener",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Redirect: GET /Ab3xK9 &#8594; 301</div>
  <div class="d-arrow-down">&#8595; fire-and-forget</div>
  <div class="d-box green">Kinesis Data Streams (raw click event)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Lambda Enrichment (per batch, every 60s)</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-box gray">IP &#8594; Geo (MaxMind GeoIP2)</div>
        <div class="d-box gray">UA &#8594; Device/Browser parsing</div>
        <div class="d-box gray">Referrer &#8594; domain extraction</div>
        <div class="d-box gray">Bot score computation</div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-box amber">Filter bots (score &gt; 0.7)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Aggregate by dimensions</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Batch write to DynamoDB</div>
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
        <div class="d-box green">Availability: 99.99% (52 min/yr)</div>
        <div class="d-box green">Redirect p99: &lt;10ms</div>
        <div class="d-box green">Write Success: 99.9%</div>
        <div class="d-box green">Cache Hit: &gt;60%</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Real-time Dashboard</div>
      <div class="d-flow-v">
        <div class="d-box blue">RPS (requests per second)</div>
        <div class="d-box blue">Latency percentiles (p50/p95/p99)</div>
        <div class="d-box blue">Error rate (4xx, 5xx)</div>
        <div class="d-box blue">Cache hit ratio</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Business Dashboard</div>
      <div class="d-flow-v">
        <div class="d-box purple">New URLs/day</div>
        <div class="d-box purple">Top domains shortened</div>
        <div class="d-box purple">Geo distribution</div>
        <div class="d-box purple">Active users</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Cost Dashboard</div>
      <div class="d-flow-v">
        <div class="d-box amber">Daily spend by service</div>
        <div class="d-box amber">DynamoDB RCU/WCU usage</div>
        <div class="d-box amber">CloudFront transfer (GB)</div>
        <div class="d-box amber">Month-over-month trend</div>
      </div>
    </div>
  </div>
</div>
<div class="d-flow-v">
  <div class="d-box red">Distributed Tracing: AWS X-Ray &#8212; ALB &#8594; ECS &#8594; Redis &#8594; DynamoDB (end-to-end latency breakdown)</div>
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
        <div class="d-box green">GetItem(PK=short_code) &#8594; O(1), &lt;5ms</div>
        <div class="d-box green">PutItem(condition: attribute_not_exists) &#8594; O(1)</div>
        <div class="d-box blue">DeleteItem(PK=short_code) &#8594; O(1)</div>
        <div class="d-box amber">TTL auto-deletes expired URLs &#8594; zero cost</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">GSI-1: user-urls-index</div>
      <div class="d-flow-v">
        <div class="d-box purple">PK: user_id | SK: created_at</div>
        <div class="d-box purple">Projection: ALL (full item copy)</div>
        <div class="d-box purple">Query: "List my URLs" &#8594; paginated</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Hot Partition Handling</div>
      <div class="d-flow-v">
        <div class="d-box red">Viral URL: 10M clicks/day</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Cache hit rate 95%+ &#8594; only 500K DB reads/day</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">DynamoDB: 3,000 RCU per partition = 3M reads/sec</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Adaptive capacity redistributes within minutes</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Analytics: Separate Table</div>
      <div class="d-flow-v">
        <div class="d-box gray">PK: short_code#date &#8594; isolates analytics writes</div>
        <div class="d-box gray">Prevents viral URL analytics from affecting redirects</div>
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

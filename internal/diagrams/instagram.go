package diagrams

func registerInstagram(r *Registry) {
	r.Register(&Diagram{
		Slug:        "ig-api-design",
		Title:       "API Design",
		Description: "API endpoint prioritization for core, important, and nice-to-have features",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P0 — Core (Must Have)</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Multipart upload: client gets pre-signed S3 URL, uploads directly, then POSTs metadata. Max 50MB. Returns media_id for CDN URL construction.">POST /api/v1/media &#8594; upload photo, returns media_id <div class="d-tag green">recommended</div></div>
        <div class="d-box green" data-tip="Cursor-based pagination (cursor=last_post_id). Limit 20 items. Sorted by score DESC. Avoids OFFSET N which degrades with depth.">GET /api/v1/feed &#8594; paginated feed (cursor-based)</div>
        <div class="d-box green" data-tip="Idempotent toggle: follow if not following, unfollow if following. Returns 200 with new state. Writes to follows table and triggers fan-out.">POST /api/v1/follow/{user_id} &#8594; follow/unfollow</div>
        <div class="d-box green" data-tip="Idempotent: second like is a no-op (409 if duplicate). Composite PK (user_id, post_id) enforces dedup. Counter incremented via sharded Redis INCR.">POST /api/v1/like/{post_id} &#8594; toggle like</div>
        <div class="d-box green" data-tip="Writes to comments table with (post_id, created_at) index. Max 2200 chars. Async toxicity check — published immediately, removed async if flagged.">POST /api/v1/comment/{post_id} &#8594; add comment</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P1 — Important</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="DynamoDB TTL field set to now+86400s. S3 lifecycle deletes media at 25h. Redis sorted set holds story order. View tracking via bitmap.">POST /api/v1/stories &#8594; 24hr ephemeral content</div>
        <div class="d-box blue" data-tip="Pre-computed explore feed from ML ranking model. Batch Spark job runs hourly. Personalized per interest cluster. Cached in Redis, refreshed on scroll.">GET /api/v1/explore &#8594; trending + personalized</div>
        <div class="d-box blue" data-tip="SSE preferred for unidirectional push. EventSource reconnects automatically. Falls back to long-polling for older clients. Max 100K concurrent SSE connections per gateway.">GET /api/v1/notifications &#8594; SSE/WebSocket</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P2 — Nice to Have</div>
      <div class="d-flow-v">
        <div class="d-box gray" data-tip="E2E encrypted DMs require separate key exchange infra. Separate WebSocket cluster. Out of scope for initial design — adds significant complexity.">POST /api/v1/messages/{user_id} &#8594; DMs</div>
        <div class="d-box gray" data-tip="Video transcoding pipeline (FFmpeg) adds significant cost and latency. Separate upload flow from photos. Out of scope for initial design.">POST /api/v1/reels &#8594; short video</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-nfr-estimates",
		Title:       "Non-Functional Requirements & Back-of-Envelope Estimates",
		Description: "Non-functional requirements and back-of-envelope capacity estimates",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">NFR Targets</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="99.99% = 52 min downtime/yr. Achieved via multi-AZ deployment, ALB health checks, auto-scaling, and circuit breakers.">Availability: 99.99% (52 min downtime/yr)</div>
        <div class="d-box green" data-tip="p99 &lt;200ms for feed. Achieved by: Redis cache (1-2ms), CDN edge (&lt;50ms for images), read replicas. Measure with CloudWatch p99 latency alarm.">Latency: Feed load &lt; 200ms p99 <span class="d-metric latency">&lt;200ms p99</span></div>
        <div class="d-box blue" data-tip="2B MAU = monthly active. 500M DAU = daily active. 25% DAU/MAU ratio. Peak during events: 10x baseline.">Scale: 2B MAU, 500M DAU <span class="d-metric throughput">500M DAU</span></div>
        <div class="d-box blue" data-tip="10 reads per write means caching is extremely effective. Redis hit rate target: 95%+. Cache-first architecture justified.">Read:Write ratio: ~10:1</div>
        <div class="d-box amber" data-tip="Feed is eventual: stale by up to 5min (Redis TTL). Likes are strong: composite PK in Postgres enforces dedup. Auth is always strong consistency.">Consistency: Eventual for feed, strong for likes</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Back-of-Envelope Math</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="500M &#215; 10 = 5B/day ÷ 86400s = 57,870 RPS. Round to 58K. Peak 10x = 580K RPS. This drives the caching and fan-out architecture.">500M DAU &#215; 10 loads/day = 5B feed req/day &#8776; <span class="d-metric throughput">58K RPS</span></div>
        <div class="d-box purple" data-tip="100M ÷ 86400 = 1,157 uploads/sec baseline. 5x peak = 5,785/sec. Pre-signed URLs offload this traffic from app servers to S3 directly.">100M photos/day &#8776; <span class="d-metric throughput">1,150 uploads/sec</span> (5x peak = 5,750/sec)</div>
        <div class="d-box purple" data-tip="Likes + comments are ~10x photo count. Each like is a DB write (DynamoDB) and a Redis INCR. Total write RPS ~6K sustained, 60K peak.">Likes/comments add 10x writes &#8594; total <span class="d-metric throughput">~6K write RPS</span></div>
        <div class="d-box amber" data-tip="100M photos &#215; 2MB = 200TB raw. After 4-size resize pipeline: thumbnail(15KB)+small(40KB)+medium(120KB)+full(300KB) = 475KB total per photo. Actual: ~47TB/day.">100M/day &#215; 2MB avg = <span class="d-metric size">200 TB/day new media</span></div>
        <div class="d-box amber" data-tip="200TB &#215; 365 = 73PB/year. S3 Standard: $23/TB = $1.7M/mo. S3 Intelligent-Tiering moves cold data to cheaper tiers automatically.">Year 1 storage: <span class="d-metric size">~73 PB</span> (200TB &#215; 365)</div>
      </div>
    </div>
  </div>
</div>
<div class="d-flow">
  <div class="d-number"><div class="d-number-value">58K</div><div class="d-number-label">Feed RPS</div></div>
  <div class="d-number"><div class="d-number-value">1,150</div><div class="d-number-label">Upload/sec</div></div>
  <div class="d-number"><div class="d-number-value">200 TB</div><div class="d-number-label">Media/day</div></div>
  <div class="d-number"><div class="d-number-value">73 PB</div><div class="d-number-label">Year-1 Storage</div></div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-mvp-architecture",
		Title:       "MVP Architecture",
		Description: "MVP architecture with monolith, Postgres, S3, and CloudFront",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-cols">
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box blue" data-tip="Mobile-first. iOS and Android clients use HTTP/2. Web uses React SPA. Client prefetches feed on app launch to hide latency."><span class="d-step">1</span> Client (iOS / Android / Web)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box purple" data-tip="Caches static assets (JS, CSS, images) at 400+ PoPs. Cache-Control: max-age=31536000 for immutable assets. Reduces origin load by 90%."><span class="d-step">2</span> CloudFront (CDN) <span class="d-metric latency">&lt;50ms</span></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box purple" data-tip="Layer 7 load balancer. Round-robin across ECS tasks. Health checks every 10s. Connection draining 30s on deploys. Terminates TLS."><span class="d-step">3</span> ALB (Load Balancer) <span class="d-metric latency">1ms</span></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box green" data-tip="Single Django/FastAPI process handles all routes. 2&#215; t3.large at MVP. Stateless — any instance handles any request. Scale horizontally as load grows."><span class="d-step">4</span> ECS (Django / FastAPI Monolith) <div class="d-tag green">start here</div></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-row">
        <div class="d-box indigo" data-tip="Single RDS Postgres db.r6g.large. ACID for all writes. B-tree indexes on (user_id, created_at). Upgrade to read replicas when p99 &gt; 50ms.">Postgres (RDS) <span class="d-metric latency">5-20ms</span></div>
        <div class="d-box red" data-tip="Single ElastiCache node for session storage and feed cache. 5-min TTL on feed, 1-hr on profiles. Upgrade to cluster when memory &gt; 80%.">Redis (ElastiCache) <span class="d-metric latency">&lt;1ms</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">MEDIA STORAGE</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Client gets pre-signed URL (15-min expiry), uploads directly to S3. Bypasses app server — avoids 200TB/day saturating ECS. 11 nines durability."><span class="d-step">1</span> S3 Bucket (Photos) <div class="d-tag blue">S3</div></div>
        <div class="d-label">Pre-signed URL upload (bypasses app server)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="S3 origin + CloudFront distribution. Cache-Control: public, max-age=31536000. WebP served to supporting clients. Reduces S3 GET costs by 10x."><span class="d-step">2</span> CloudFront (CDN) <span class="d-metric latency">&lt;50ms global</span></div>
        <div class="d-label">Edge-cached images</div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">MVP runs as a single monolith on 2 ECS tasks. Start here and extract services only when a specific bottleneck is measured. Total cost: ~$500/mo for &lt;1M users.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-mvp-tech-stack",
		Title:       "MVP Tech Stack (Visual)",
		Description: "Visual breakdown of MVP tech stack components and monthly costs",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Compute</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="2 Fargate tasks, each 2 vCPU 4GB RAM. Auto-scales to 8 tasks at CPU>70%. Stateless — sessions in Redis. $120/mo baseline.">ECS Fargate: 2&#215; t3.large &#8212; <span class="d-metric cost">$120/mo</span></div>
        <div class="d-box purple" data-tip="Layer 7 load balancer with TLS termination. Health checks every 10s. Connection draining 30s. Cross-zone enabled.">ALB (Load Balancer) &#8212; <span class="d-metric cost">$25/mo</span> <div class="d-tag blue">recommended</div></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Storage</div>
      <div class="d-flow-v">
        <div class="d-box indigo" data-tip="db.r6g.large: 2 vCPU, 16GB RAM. gp3 500GB SSD. Multi-AZ standby for failover. Upgrade to r6g.2xlarge when connections exceed 1K/sec.">RDS Postgres db.r6g.large &#8212; <span class="d-metric cost">$200/mo</span></div>
        <div class="d-box red" data-tip="t4g.medium: 2 vCPU, 4GB RAM. Sufficient for &lt;1M users. Upgrade to cluster when memory exceeds 80% or latency &gt;5ms.">ElastiCache Redis t4g.medium &#8212; <span class="d-metric cost">$50/mo</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Media &amp; Delivery</div>
      <div class="d-flow-v">
        <div class="d-box amber">S3 Standard &#8212; $23/TB/mo</div>
        <div class="d-box purple">CloudFront CDN &#8212; $85/10TB</div>
      </div>
    </div>
  </div>
</div>
<div class="d-box blue">Total: &#8776;$500/mo for &lt;1M users</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-data-model",
		Title:       "Data Model (Entity Relationships)",
		Description: "Entity relationship diagram for users, posts, follows, likes, and comments",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity" data-tip="Sharded by user_id at 50M+ users. Each shard holds ~10M users. Postgres BIGSERIAL gives 9.2 quintillion IDs.">
      <div class="d-entity-header blue">users <span class="d-metric size">~2B rows</span></div>
      <div class="d-entity-body">
        <div class="pk">id BIGSERIAL</div>
        <div class="idx idx-unique">username VARCHAR(30)</div>
        <div class="idx idx-unique">email VARCHAR(255)</div>
        <div>bio TEXT</div>
        <div>avatar_url TEXT</div>
        <div>created_at TIMESTAMP</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity" data-tip="Hottest table by write volume. Composite index on (user_id, created_at DESC) is critical for profile page queries. Migrated to DynamoDB at scale for single-digit ms reads.">
      <div class="d-entity-header green">posts <span class="d-metric throughput">1,150 writes/sec</span></div>
      <div class="d-entity-body">
        <div class="pk">id BIGSERIAL</div>
        <div class="fk">user_id BIGINT &#8594; users.id</div>
        <div>caption TEXT</div>
        <div>media_url TEXT NOT NULL</div>
        <div>media_type VARCHAR(10)</div>
        <div class="idx idx-btree">created_at TIMESTAMP</div>
      </div>
    </div>
    <div style="font-size:0.68rem; color:var(--text-muted); margin-top:4px; text-align:center;">idx: (user_id, created_at DESC)</div>
  </div>
  <div class="d-col">
    <div class="d-entity" data-tip="Composite PK (follower_id, followee_id) prevents duplicates. Reverse index on followee_id enables 'who follows me?' in O(log N). At scale, this is the largest table by row count.">
      <div class="d-entity-header purple">follows <span class="d-metric size">~400B rows</span></div>
      <div class="d-entity-body">
        <div class="pk fk">follower_id BIGINT &#8594; users.id</div>
        <div class="pk fk">followee_id BIGINT &#8594; users.id</div>
        <div>created_at TIMESTAMP</div>
      </div>
    </div>
    <div style="font-size:0.68rem; color:var(--text-muted); margin-top:4px; text-align:center;">idx: followee_id (reverse lookup)</div>
    <div class="d-entity" style="margin-top: 0.75rem;" data-tip="Composite PK deduplicates likes. At viral scale, sharded counters aggregate across 100 Redis shards every 5s. DynamoDB at scale for 100K+ WPS.">
      <div class="d-entity-header amber">likes <span class="d-metric throughput">100K+ WPS peak</span></div>
      <div class="d-entity-body">
        <div class="pk fk">user_id BIGINT &#8594; users.id</div>
        <div class="pk fk">post_id BIGINT &#8594; posts.id</div>
        <div>created_at TIMESTAMP</div>
      </div>
    </div>
    <div style="font-size:0.68rem; color:var(--text-muted); margin-top:4px; text-align:center;">idx: post_id (count query)</div>
  </div>
  <div class="d-col">
    <div class="d-entity" data-tip="Ordered by (post_id, created_at DESC) for threaded display. Text indexed with GIN for spam detection. At scale, partitioned by post_id range.">
      <div class="d-entity-header red">comments <span class="d-metric throughput">~6K WPS</span></div>
      <div class="d-entity-body">
        <div class="pk">id BIGSERIAL</div>
        <div class="fk">user_id BIGINT &#8594; users.id</div>
        <div class="fk">post_id BIGINT &#8594; posts.id</div>
        <div>text TEXT NOT NULL</div>
        <div class="idx idx-btree">created_at TIMESTAMP</div>
      </div>
    </div>
    <div style="font-size:0.68rem; color:var(--text-muted); margin-top:4px; text-align:center;">idx: (post_id, created_at DESC)</div>
  </div>
</div>
<div class="d-er-lines">
  <div class="d-er-connector"><span class="d-er-from">users</span> <span class="d-er-type">1:N</span> <span class="d-er-to">posts</span></div>
  <div class="d-er-connector"><span class="d-er-from">users</span> <span class="d-er-type">M:N</span> <span class="d-er-to">users</span> (via follows)</div>
  <div class="d-er-connector"><span class="d-er-from">users</span> <span class="d-er-type">M:N</span> <span class="d-er-to">posts</span> (via likes)</div>
  <div class="d-er-connector"><span class="d-er-from">posts</span> <span class="d-er-type">1:N</span> <span class="d-er-to">comments</span></div>
</div>
<div class="d-legend">
  <span class="d-legend-item"><span class="d-legend-color blue"></span>Identity</span>
  <span class="d-legend-item"><span class="d-legend-color green"></span>Content</span>
  <span class="d-legend-item"><span class="d-legend-color purple"></span>Relationships</span>
  <span class="d-legend-item"><span class="d-legend-color amber"></span>Engagement (high write)</span>
  <span class="d-legend-item"><span class="d-legend-color red"></span>User-generated text</span>
</div>
<div class="d-caption">Five tables handle 99% of Instagram's core data. The follows table is the largest by row count (~400B), while likes sees the highest peak write throughput during viral events.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-feed-query-result",
		Title:       "Example: Feed Query Result",
		Description: "Example feed query result showing post data with engagement counts",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-query-result">
  <table>
    <tr><th>post_id</th><th>username</th><th>caption</th><th>like_count</th><th>comment_count</th><th>created_at</th></tr>
    <tr><td>4521</td><td>@foodie_jane</td><td>Best ramen in Tokyo...</td><td>342</td><td>28</td><td>2 min ago</td></tr>
    <tr><td>4519</td><td>@travel_mike</td><td>Sunset in Santorini</td><td>1,205</td><td>89</td><td>15 min ago</td></tr>
    <tr><td>4515</td><td>@dev_sarah</td><td>My setup for 2024</td><td>567</td><td>43</td><td>1 hr ago</td></tr>
  </table>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-upload-read-flow",
		Title:       "Upload & Read Flow",
		Description: "Side-by-side write path for photo upload and read path for feed loading",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">WRITE PATH (Photo Upload)</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="d-step">1</span> Client</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Direct upload to S3 via pre-signed URL bypasses app server. Max file size 50MB. URL expires in 15 minutes."><span class="d-step">2</span> S3 (Pre-signed URL Upload) <span class="d-metric latency">200ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue"><span class="d-step">3</span> POST /media (with S3 key)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Round-robin across ECS tasks. Health checks every 10s. Drains unhealthy targets in 30s."><span class="d-step">4</span> ALB <span class="d-metric latency">1ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Validates S3 key exists, writes metadata to DB. Triggers async fan-out via Kafka event."><span class="d-step">5</span> ECS (App Server) <span class="d-metric latency">15ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo" data-tip="ACID write for post record. At scale, migrated to DynamoDB for single-digit ms. Postgres failover to standby in 30s."><span class="d-step">6</span> Postgres (Write post record) <span class="d-metric latency">5ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="DEL key invalidation, not update. Prevents stale cache serving old post count on profile."><span class="d-step">7</span> Redis (Invalidate cache) <span class="d-metric latency">&lt;1ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray"><span class="d-step">8</span> Response (post_id + CDN URL) <span class="d-status active"></span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">READ PATH (Feed Load) <span class="d-metric throughput">58K RPS</span></div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="d-step">1</span> Client</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple"><span class="d-step">2</span> ALB <span class="d-metric latency">1ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><span class="d-step">3</span> ECS (App Server)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="90%+ hit rate for active users. ZREVRANGE returns pre-sorted post IDs. Miss triggers DB query."><span class="d-step">4</span> Redis (Cache check) <span class="d-metric latency">1-2ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo" data-tip="Only hit on cache miss (~10%). Fan-out-on-read query joins followees + recent posts. Timeout at 5s."><span class="d-step">5</span> Postgres (Cache miss) <span class="d-metric latency">20-50ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red"><span class="d-step">6</span> Redis (Cache feed, TTL 5m)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="400+ edge PoPs. Cache-Control: max-age=31536000 for immutable images. Origin Shield reduces origin hits by 90%."><span class="d-step">7</span> CloudFront (Serve images) <span class="d-metric latency">&lt;50ms global</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray"><span class="d-step">8</span> Response (Feed JSON + CDN URLs) <span class="d-status active"></span></div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Write path total: ~220ms (dominated by S3 upload). Read path total: &lt;55ms on cache hit (90%+ of requests), ~100ms on cache miss. The 100:1 read/write ratio makes caching highly effective.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-write-path",
		Title:       "Write Path: Photo Upload (Hop by Hop)",
		Description: "Hop-by-hop write path for photo upload from client to database",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box blue">W1: Client uploads photo to S3 via pre-signed URL</div>
  <div class="d-label">S3 | Direct upload, no app server bottleneck | Failure: S3 503 &#8594; client retries</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box blue">W2: Client calls POST /media with S3 key</div>
  <div class="d-label">ALB &#8594; ECS | Metadata only, photo already in S3</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">W3: App server writes post record to Postgres</div>
  <div class="d-label">RDS Postgres | ACID for post creation | Failure: DB failover to standby (30s)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red">W4: Invalidate user profile cache</div>
  <div class="d-label">ElastiCache Redis | Profile shows latest post count | Failure: cache miss, serve stale</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green">W5: Return post_id + CDN URL to client</div>
  <div class="d-label">Response | Client can immediately share the post</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-read-path",
		Title:       "Read Path: Feed Load (Hop by Hop)",
		Description: "Hop-by-hop read path for feed loading with cache layers",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box blue">R1: Client requests GET /api/v1/feed</div>
  <div class="d-label">ALB &#8594; ECS | App server handles feed assembly | Failure: ALB drains unhealthy node</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red">R2: Check Redis for cached feed</div>
  <div class="d-label">ElastiCache Redis | 90%+ cache hit for active users | Failure: cache miss, fall through</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">R3: Query Postgres: followees then recent posts</div>
  <div class="d-label">RDS Postgres | Fan-out-on-read: assemble feed at query time | Failure: slow query, timeout 5s</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red">R4: Cache assembled feed in Redis (TTL 5 min)</div>
  <div class="d-label">ElastiCache Redis | Next request hits cache | Failure: Redis full, evict LRU</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">R5: Client loads images via CDN URLs</div>
  <div class="d-label">CloudFront &#8594; S3 | Edge-cached images, &lt;50ms globally | Failure: cache miss, origin 200ms</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-cdn-media-pipeline",
		Title:       "CDN & Media Storage Pipeline",
		Description: "CDN and media storage pipeline with image resizing and edge distribution",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box blue">Client (iOS / Android / Web)</div>
  <div class="d-arrow-down">&#8595; pre-signed URL upload</div>
  <div class="d-box amber">S3 Ingest Bucket (raw uploads)</div>
  <div class="d-arrow-down">&#8595; S3 Event Notification</div>
  <div class="d-box green">Lambda: Image Processing Pipeline</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Thumbnail</div>
      <div class="d-box purple">150&#215;150</div>
      <div class="d-label">~15KB</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Small</div>
      <div class="d-box purple">320&#215;320</div>
      <div class="d-label">~40KB</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Medium</div>
      <div class="d-box purple">640&#215;640</div>
      <div class="d-label">~120KB</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Full</div>
      <div class="d-box purple">1080&#215;1080</div>
      <div class="d-label">~300KB</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; all sizes saved</div>
  <div class="d-box amber">S3 Media Bucket (processed images, WebP + JPEG)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">CDN Distribution</div>
        <div class="d-flow-v">
          <div class="d-box purple">CloudFront Origin Shield (regional cache)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">CloudFront Edge PoPs (400+ locations)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">User device (&lt;50ms globally)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">CDN Config</div>
        <div class="d-flow-v">
          <div class="d-box green">Cache-Control: public, max-age=31536000</div>
          <div class="d-box green">Accept header &#8594; WebP vs JPEG negotiation</div>
          <div class="d-box green">Origin Shield &#8594; single origin fetch per region</div>
          <div class="d-box amber">S3 CRR to eu-west-1, ap-south-1 for global reads</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-hybrid-fanout",
		Title:       "Feed Strategy: Hybrid Fan-out Architecture",
		Description: "Hybrid fan-out architecture comparing write-based and read-based strategies",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fan-out-on-Write (Normal Users &lt; 100K followers)</div>
      <div class="d-flow-v">
        <div class="d-box green"><span class="d-step">1</span> User creates post <span class="d-status active"></span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue" data-tip="Reads from follows table or Redis set. Avg user has ~200 followers. For 1K followers, fan-out completes in &lt;50ms."><span class="d-step">2</span> Feed Service reads followers list <span class="d-metric latency">2ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="ZADD is O(log N) per write. For 1K followers = 1K ZADD ops. Pipelined in batches of 100 for ~10ms total. Score = Unix timestamp for time-ordering."><span class="d-step">3</span> Write post_id to each follower's Redis sorted set <span class="d-metric throughput">100K ZADD/sec</span></div>
        <div class="d-label">ZADD feed:{user_id} {timestamp} {post_id}</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Prevents unbounded memory growth. 500 posts &#215; 16 bytes = 8KB per user feed. For 500M users = ~4TB Redis."><span class="d-step">4</span> Trim to latest 500: ZREMRANGEBYRANK feed:{uid} 0 -501 <span class="d-metric size">8KB/user</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><span class="d-step">5</span> &#10003; Feed pre-computed. Read = O(1) Redis lookup <span class="d-metric latency">&lt;2ms read</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fan-out-on-Read (Celebrities &gt; 100K followers)</div>
      <div class="d-flow-v">
        <div class="d-box green"><span class="d-step">1</span> Celebrity creates post <span class="d-status active"></span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="A celebrity with 200M followers would need 200M Redis writes. At 100K writes/sec = 33 minutes of fan-out. Unacceptable latency."><span class="d-step">2</span> Skip fan-out (would be 200M writes) <span class="d-status error"></span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue"><span class="d-step">3</span> Mark as celebrity_post in Redis set <span class="d-metric latency">&lt;1ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="At read time, fetch pre-computed feed + query latest N posts from each followed celebrity. Typically 5-10 celebrities per user. Merge sort in memory."><span class="d-step">4</span> On feed read: merge pre-computed + celebrity posts <span class="d-metric latency">15-30ms</span></div>
        <div class="d-label">DB query: SELECT FROM posts WHERE user_id IN (celeb_ids)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber"><span class="d-step">5</span> Sort by timestamp, return top N</div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Hybrid approach: fan-out-on-write for 99% of users (fast reads), fan-out-on-read for top 1% celebrities (avoids 200M write storms). The threshold of 100K followers balances write cost vs read latency.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-stage2-architecture",
		Title:       "Stage 2 Architecture (Visual)",
		Description: "Stage 2 architecture with feed service, read replicas, and async processing",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box blue">Client (iOS / Android / Web)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">CloudFront + Origin Shield</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">ALB</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Monolith (ECS &#215;8)</div>
        <div class="d-flow-v">
          <div class="d-box green">API Server (Django/FastAPI)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box indigo">Postgres Primary</div>
            <div class="d-box indigo">Read Replica &#215;2</div>
          </div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Feed Service (ECS &#215;4)</div>
        <div class="d-flow-v">
          <div class="d-box green">Fan-out Worker</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red">Redis Cluster (3-node)</div>
          <div class="d-label">Sorted sets per user</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Async Processing</div>
        <div class="d-flow-v">
          <div class="d-box amber">SQS (likes/comments queue)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Lambda (batch writer)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">S3 &#8594; Lambda &#8594; 4 image sizes</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-service-db-ownership",
		Title:       "Microservices: Service &#8594; Database Ownership",
		Description: "Microservices mapped to their owned databases and storage technologies",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-cols">
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green">User Service</div>
      <div class="d-label">Profiles, auth, follows</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box indigo">Postgres (sharded)</div>
      <div class="d-label">ACID for auth, B-tree for follows</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green">Post Service</div>
      <div class="d-label">Posts, media metadata</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">DynamoDB</div>
      <div class="d-label">PK=user_id, SK=timestamp</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green">Feed Service</div>
      <div class="d-label">Feed computation + cache</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box red">Redis Cluster</div>
      <div class="d-label">Sorted sets, sub-ms reads</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green">Media Service</div>
      <div class="d-label">Upload, process, serve</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">S3 + Lambda</div>
      <div class="d-label">11 nines durability</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green">Engagement Service</div>
      <div class="d-label">Likes, comments, shares</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">DynamoDB</div>
      <div class="d-label">Sharded counters, 100K+ WPS</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green">Notification Svc</div>
      <div class="d-label">Push, in-app, email</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">SQS + DynamoDB</div>
      <div class="d-label">Async, TTL auto-cleanup</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green">Search/Explore</div>
      <div class="d-label">Hashtags, user search</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box purple">Elasticsearch</div>
      <div class="d-label">Full-text, autocomplete</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-microservices-architecture",
		Title:       "Microservices Architecture",
		Description: "Full microservices architecture with API gateway and Kafka event bus",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box blue" data-tip="500M DAU across iOS, Android, and Web. Mobile clients use HTTP/2 multiplexing. Average session: 30 minutes, 10 feed loads.">Client (iOS / Android / Web) <span class="d-metric throughput">58K RPS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="400+ edge PoPs globally. Caches static assets (images, JS, CSS) with 1-year TTL. Origin Shield reduces origin fetches by 90%.">CloudFront (CDN) <span class="d-metric latency">&lt;50ms edge</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Path-based routing to microservices. WAF blocks malicious requests. Rate limiting at 100 req/min per user. gRPC internally between services.">ALB (API Gateway) <span class="d-metric latency">1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box green" data-tip="Profiles, auth, follow graph. Sharded Postgres by user_id. ACID for auth operations. gRPC for internal calls.">User Svc <span class="d-status active"></span></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box indigo" data-tip="Sharded by user_id. B-tree indexes on follows for O(log N) lookups. Read replicas for profile reads.">Postgres <span class="d-metric latency">5ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green" data-tip="Post CRUD and media metadata. PK=user_id, SK=timestamp for time-ordered queries. Single-digit ms at any scale.">Post Svc <span class="d-status active"></span></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="On-demand capacity. PK=user_id, SK=created_at. DAX cache for hot posts. Global Tables for multi-region.">DynamoDB <span class="d-metric latency">3ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green" data-tip="Hybrid fan-out engine. Writes to Redis sorted sets for normal users. Merges celebrity posts at read time.">Feed Svc <span class="d-status active"></span></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box red" data-tip="50-node cluster. Sorted sets per user feed. 8KB per user &#215; 500M users = ~4TB. Sub-ms ZREVRANGE for feed reads.">Redis Cluster <span class="d-metric latency">&lt;2ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green" data-tip="Pre-signed URL upload, Lambda resize to 4 sizes (150px to 1080px), WebP + JPEG. 11 nines durability.">Media Svc <span class="d-status active"></span></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="200TB/day new media. S3 Intelligent-Tiering for cost optimization. CRR to 3 regions.">S3 <span class="d-metric size">200 TB/day</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green" data-tip="Likes, comments, shares. Sharded counters for viral posts. Deduplication via composite PK.">Engagement Svc <span class="d-status active"></span></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="Sharded counters: 100 shards per viral post. Aggregate every 5s. 100K+ writes/sec sustained.">DynamoDB <span class="d-metric throughput">100K+ WPS</span></div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" data-tip="Decouples all services. 3 brokers, 50 partitions per topic. Retention 7 days. Enables event sourcing and async processing.">Kafka (MSK) &#8212; Event Bus: post_created, user_followed, post_liked <span class="d-metric throughput">500K events/sec</span></div>
</div>
<div class="d-legend">
  <span class="d-legend-item"><span class="d-legend-color blue"></span>Client / External</span>
  <span class="d-legend-item"><span class="d-legend-color purple"></span>Networking / CDN</span>
  <span class="d-legend-item"><span class="d-legend-color green"></span>Application Services</span>
  <span class="d-legend-item"><span class="d-legend-color indigo"></span>Relational DB</span>
  <span class="d-legend-item"><span class="d-legend-color amber"></span>NoSQL / Object Storage</span>
  <span class="d-legend-item"><span class="d-legend-color red"></span>Cache (Redis)</span>
  <span class="d-legend-item"><span class="d-legend-color gray"></span>Event Streaming</span>
</div>
<div class="d-caption">Each microservice owns its database. Kafka event bus decouples all services and enables eventual consistency. Total request path: client &#8594; CDN &#8594; ALB &#8594; service &#8594; DB in &lt;100ms for reads.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-notification-flow",
		Title:       "Notification System Flow",
		Description: "Notification system flow from event producers through Kafka to push and in-app delivery",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Event Producers</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="Fires on every new post. Triggers fan-out to all followers' notification feeds."><span class="d-step">1</span> Post Svc: post_created <span class="d-status active"></span></div>
          <div class="d-box green" data-tip="Highest volume event source. Viral posts generate 1M+ like events. Batched to avoid notification spam."><span class="d-step">1</span> Engagement Svc: post_liked, comment_added <span class="d-metric throughput">100K events/sec</span></div>
          <div class="d-box green"><span class="d-step">1</span> User Svc: user_followed</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-label">all events flow into</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="Partitioned by target_user_id for ordering guarantees. 50 partitions, 7-day retention. Consumer lag monitored via CloudWatch."><span class="d-step">2</span> Kafka (MSK) &#8212; notification topic <span class="d-metric throughput">500K msg/sec</span></div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="Deduplicates using (event_type, target_user, source_entity) key in Redis with 1h TTL. Batches likes into 'alice and 12 others liked your post' style messages."><span class="d-step">3</span> Notification Service (ECS consumers) <span class="d-metric latency">50-200ms</span></div>
  <div class="d-label">Dedup by (event_type, target_user, source_entity) &#8212; batch similar events</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">Mobile Push</div>
        <div class="d-flow-v">
          <div class="d-box amber" data-tip="Fan-out to platform-specific endpoints. Handles token rotation and delivery receipts. $1 per 1M publishes."><span class="d-step">4a</span> SNS Platform App <span class="d-metric cost">$1/1M msgs</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box blue">APNs (iOS) <span class="d-metric latency">100-500ms</span></div>
            <div class="d-box green">FCM (Android) <span class="d-metric latency">100-500ms</span></div>
          </div>
        </div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">Web Real-time</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="SSE preferred over WebSocket for unidirectional notifications. 100K concurrent connections per gateway instance. Redis Pub/Sub for cross-instance delivery."><span class="d-step">4b</span> SSE / WebSocket Gateway <span class="d-metric latency">&lt;50ms</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">Browser client <span class="d-status active"></span></div>
        </div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">In-App Badge</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="Atomic INCR for badge count. Reset to 0 on notification tab open. O(1) operation, sub-ms latency."><span class="d-step">4c</span> Redis INCR badge:{user_id} <span class="d-metric latency">&lt;1ms</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber" data-tip="TTL 30 days auto-deletes old notifications. PK=user_id, SK=timestamp. On-demand capacity handles traffic spikes."><span class="d-step">5</span> DynamoDB (notification history, TTL 30d) <span class="d-metric latency">3ms</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-label">Rate limit: max 50 push notifications/hr per user &#8212; batch likes: "alice and 12 others liked your post"</div>
</div>
<div class="d-caption">End-to-end notification latency: &lt;500ms for push, &lt;50ms for web real-time. Like batching reduces push volume by 90% during viral events. Rate limiting prevents notification fatigue.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-peak-traffic",
		Title:       "Peak Traffic Defense Layers",
		Description: "Three-layer defense strategy for handling 10x peak traffic spikes",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box red">&#9889; Peak Traffic: 580K RPS (10&#215; normal)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Layer 1: Absorb</div>
        <div class="d-flow-v">
          <div class="d-box purple">CDN pre-warming (push to all PoPs)</div>
          <div class="d-box purple">Rate limiter (100 req/min per user)</div>
          <div class="d-box purple">Auto-scaling ECS (CPU/memory triggers)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Layer 2: Buffer</div>
        <div class="d-flow-v">
          <div class="d-box amber">SQS write queue (likes/comments)</div>
          <div class="d-box amber">Sharded counters (100 shards)</div>
          <div class="d-box amber">Increased cache TTL during peaks</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Layer 3: Degrade</div>
        <div class="d-flow-v">
          <div class="d-box gray">Circuit breakers (fail fast)</div>
          <div class="d-box gray">Serve cached/stale feed</div>
          <div class="d-box gray">Fallback to trending (pre-computed)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-label">Normal: 58K RPS &#8594; Peak: 580K RPS &#8594; Super Bowl: 1M likes/sec on one post</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-caching-strategy",
		Title:       "Multi-Layer Caching Strategy",
		Description: "Multi-layer caching strategy with L1 local, L2 Redis, and L3 database tiers",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box blue" data-tip="58K RPS total. Feed requests dominate at 80%. Profile and post detail make up the remaining 20%.">Incoming Request (GET /feed, GET /post, GET /profile) <span class="d-metric throughput">58K RPS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">L1 &#8212; Local In-Memory Cache</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="Process-local cache eliminates network hops entirely. Short 30s TTL prevents stale data. Each of 8 ECS tasks has its own 256MB cache. LRU eviction when full."><span class="d-step">1</span> Caffeine / Guava LRU (per ECS task) <span class="d-metric latency">&lt;1ms</span></div>
          <div class="d-label">TTL: 30s | Size: 256MB per instance</div>
          <div class="d-label">Hit rate: ~60% | No network hop, fastest layer</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">L2 &#8212; Redis Cluster (Distributed)</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="50-node cluster across 3 AZs. r6g.xlarge instances. Handles 500K ops/sec. Sorted sets for feeds, hashes for profiles. Cluster mode enabled for auto-sharding."><span class="d-step">2</span> ElastiCache Redis Cluster (50 nodes) <span class="d-metric latency">1-2ms</span> <span class="d-metric cost">~$18K/mo</span></div>
          <div class="d-label">TTL: 5min (feed), 1hr (profiles)</div>
          <div class="d-label">Hit rate: ~95% | Shared across all app instances</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">L3 &#8212; Database (Source of Truth)</div>
        <div class="d-flow-v">
          <div class="d-box indigo" data-tip="Only 5% of requests reach the database. Postgres for relational data (users, follows), DynamoDB for high-write engagement data. Read replicas absorb read spikes."><span class="d-step">3</span> Postgres / DynamoDB <span class="d-metric latency">5-50ms</span></div>
          <div class="d-label">Always consistent, highest cost</div>
          <div class="d-label">Only hit on L1 + L2 miss (~5% of requests)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; Invalidation Strategy</div>
  <div class="d-row">
    <div class="d-box amber" data-tip="Delete-on-write prevents serving stale data. L1 short TTL means max 30s staleness. No explicit L1 invalidation needed for most cases.">Write-through: DB write &#8594; delete L2 key &#8594; L1 expires via TTL</div>
    <div class="d-box amber" data-tip="Kafka consumer on each ECS task listens for invalidation events. Ensures L1 consistency across all instances within 100ms of write.">Kafka event &#8594; all instances invalidate L1 (pub/sub) <span class="d-metric latency">&lt;100ms</span></div>
  </div>
  <div class="d-label">Cache stampede protection: singleflight pattern + probabilistic early expiration (TTL &#215; random(0.8, 1.0))</div>
</div>
<div class="d-legend">
  <span class="d-legend-item"><span class="d-legend-color green"></span>L1: Local (sub-ms)</span>
  <span class="d-legend-item"><span class="d-legend-color red"></span>L2: Redis (1-2ms)</span>
  <span class="d-legend-item"><span class="d-legend-color indigo"></span>L3: Database (5-50ms)</span>
  <span class="d-legend-item"><span class="d-legend-color amber"></span>Invalidation path</span>
</div>
<div class="d-caption">Three-tier caching absorbs 95% of reads before hitting the database. Effective latency: 60% at &lt;1ms (L1), 35% at 1-2ms (L2), 5% at 5-50ms (L3). Total cache infrastructure cost: ~$20K/mo saves ~$200K/mo in database costs.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-multi-region",
		Title:       "Multi-Region Architecture",
		Description: "Multi-region architecture with Route 53 latency-based routing across three regions",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box purple">Route 53 (Global DNS &#8212; Latency-based routing)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">us-east-1 (Americas)</div>
        <div class="d-flow-v">
          <div class="d-box purple">CloudFront</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">ALB</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">ECS Services</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box indigo">Postgres</div>
            <div class="d-box amber">DynamoDB</div>
            <div class="d-box red">Redis</div>
          </div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">eu-west-1 (Europe)</div>
        <div class="d-flow-v">
          <div class="d-box purple">CloudFront</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">ALB</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">ECS Services</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box indigo">Postgres</div>
            <div class="d-box amber">DynamoDB</div>
            <div class="d-box red">Redis</div>
          </div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">ap-south-1 (India)</div>
        <div class="d-flow-v">
          <div class="d-box purple">CloudFront</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">ALB</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">ECS Services</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box indigo">Postgres</div>
            <div class="d-box amber">DynamoDB</div>
            <div class="d-box red">Redis</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-label">DynamoDB Global Tables + S3 CRR replicate across all regions</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-stories-lifecycle",
		Title:       "Stories Lifecycle",
		Description: "Stories lifecycle with 24-hour TTL auto-deletion and view tracking",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-row">
    <div class="d-box blue">User uploads story</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber">S3 (media file)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">DynamoDB (TTL = 24h)</div>
  </div>
  <div class="d-label">S3 lifecycle policy deletes media after 25h (timezone buffer)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box red">Redis sorted set (per-user story feed)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple">Redis bitmap (view tracking per story_id)</div>
  </div>
  <div class="d-label">Stories are time-ordered, not ranked &#8212; bitmap gives O(1) "has user seen this?" check</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box blue">Client renders story ring</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box gray">Unseen stories first, then by recency</div>
  </div>
  <div class="d-arrow-down">&#8595; after 24 hours</div>
  <div class="d-box red">DynamoDB TTL auto-deletes &#8594; Redis keys expire &#8594; Story gone</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-sub-problems",
		Title:       "Instagram Sub-Problems & Building Blocks",
		Description: "Key sub-problems and building blocks for Instagram at scale",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-subproblem">
    <div class="d-subproblem-icon red">&#9829;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Like Counter at Scale</div>
      <div class="d-subproblem-desc">Sharded counters + deduplication for 1M likes/sec on viral posts</div>
    </div>
    <div class="d-subproblem-link">&#8594; Consistent Hashing</div>
  </div>
  <div class="d-subproblem">
    <div class="d-subproblem-icon blue">&#8635;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Feed Ranking & Fan-out</div>
      <div class="d-subproblem-desc">Hybrid fan-out-on-write + fan-out-on-read with Redis sorted sets</div>
    </div>
    <div class="d-subproblem-link">&#8594; Redis Sorted Sets</div>
  </div>
  <div class="d-subproblem">
    <div class="d-subproblem-icon purple">&#128247;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Media Processing Pipeline</div>
      <div class="d-subproblem-desc">S3 event &#8594; Lambda &#8594; resize to 4 sizes &#8594; CDN distribution</div>
    </div>
    <div class="d-subproblem-link">&#8594; CDN / CloudFront</div>
  </div>
  <div class="d-subproblem">
    <div class="d-subproblem-icon amber">&#9889;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Real-time Notifications</div>
      <div class="d-subproblem-desc">Kafka events &#8594; push via FCM/APNs + SSE for web</div>
    </div>
    <div class="d-subproblem-link">&#8594; Rate Limiter</div>
  </div>
  <div class="d-subproblem">
    <div class="d-subproblem-icon green">&#128336;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Ephemeral Stories (24h TTL)</div>
      <div class="d-subproblem-desc">DynamoDB TTL auto-delete + Redis bitmap for view tracking</div>
    </div>
    <div class="d-subproblem-link">&#8594; DynamoDB TTL</div>
  </div>
  <div class="d-subproblem">
    <div class="d-subproblem-icon indigo">&#128256;</div>
    <div class="d-subproblem-text">
      <div class="d-subproblem-title">Distributed ID Generation</div>
      <div class="d-subproblem-desc">Snowflake IDs: 64-bit, time-sortable, zero coordination</div>
    </div>
    <div class="d-subproblem-link">&#8594; Snowflake ID</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-explore-search",
		Title:       "Explore & Search Pipeline",
		Description: "Explore and search pipeline with content moderation and ML ranking",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Data Ingestion</div>
        <div class="d-flow-v">
          <div class="d-box green">Kafka events: post_created, post_liked, comment_added</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">ML Content Classifier</div>
          <div class="d-label">NSFW/spam filter &#8594; reject or flag for review</div>
          <div class="d-arrow-down">&#8595; approved content</div>
          <div class="d-box amber">Elasticsearch Index (posts, hashtags, users)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Explore Feed Generation</div>
        <div class="d-flow-v">
          <div class="d-box indigo">Spark Batch Job (hourly)</div>
          <div class="d-label">engagement_rate = (likes + comments) / impressions</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">ML Ranking Model</div>
          <div class="d-label">Collaborative filtering: users who liked X also liked Y</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red">Redis (pre-computed explore feeds per interest cluster)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box blue">GET /explore &#8594; personalized trending grid</div>
    <div class="d-box blue">GET /search?q=sunset &#8594; Elasticsearch autocomplete</div>
    <div class="d-box blue">GET /tags/travel &#8594; hashtag feed (ES query)</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-distributed-counter",
		Title:       "Distributed Counter Architecture",
		Description: "Sharded counter architecture for handling 1M likes per second",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box blue">1M concurrent like requests</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Deduplicate: check likes table PK (user_id, post_id)</div>
  <div class="d-arrow-down">&#8595; new likes only</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Shard 0</div>
      <div class="d-box amber">INCRBY 1</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Shard 1</div>
      <div class="d-box amber">INCRBY 1</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">...</div>
      <div class="d-box gray">N shards</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Shard 99</div>
      <div class="d-box amber">INCRBY 1</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; aggregate every 5s</div>
  <div class="d-box green">Total = SUM(all shards) &#8594; cached in Redis (TTL 5s)</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-scaling-timeline",
		Title:       "Scaling Timeline: Architecture Evolution",
		Description: "Architecture evolution timeline from MVP through global scale",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">MVP (0&#8212;1M users)</div>
        <div class="d-flow-v">
          <div class="d-box green">Monolith (Django/FastAPI)</div>
          <div class="d-box indigo">Single Postgres (RDS)</div>
          <div class="d-box red">Single Redis</div>
          <div class="d-box amber">S3 + CloudFront</div>
          <div class="d-label">~$500/mo | 1 team</div>
          <div class="d-label">Pull-based feed (SQL query)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Growth (1M&#8212;50M users)</div>
        <div class="d-flow-v">
          <div class="d-box green">Monolith + Feed Service</div>
          <div class="d-box indigo">Postgres + 2 read replicas</div>
          <div class="d-box red">Redis cluster (3 nodes)</div>
          <div class="d-box amber">S3 &#8594; Lambda &#8594; 4 sizes</div>
          <div class="d-label">~$12.5K/mo | 3 teams</div>
          <div class="d-label">Hybrid fan-out + SQS buffering</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Scale (50M&#8212;500M users)</div>
        <div class="d-flow-v">
          <div class="d-box green">7 Microservices (ECS)</div>
          <div class="d-box indigo">Sharded Postgres + DynamoDB</div>
          <div class="d-box red">Redis cluster (50 nodes)</div>
          <div class="d-box gray">Kafka (MSK) event bus</div>
          <div class="d-label">~$150K/mo | 10 teams</div>
          <div class="d-label">Event-driven, gRPC, sharded counters</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Global (500M+ users)</div>
        <div class="d-flow-v">
          <div class="d-box green">Multi-region microservices</div>
          <div class="d-box indigo">DynamoDB Global Tables</div>
          <div class="d-box red">Per-region Redis clusters</div>
          <div class="d-box purple">Route 53 latency routing</div>
          <div class="d-label">~$1M+/mo | 30+ teams</div>
          <div class="d-label">CRDTs, S3 CRR, multi-region Kafka</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-flow">
    <div class="d-box blue">0&#8212;1M</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">1M&#8212;50M</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">50M&#8212;500M</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">500M+</div>
  </div>
  <div class="d-label">Key triggers: feed latency &gt; 200ms &#8594; add fan-out | write RPS &gt; 10K &#8594; shard DB | user latency &gt; 150ms &#8594; add region</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-distributed-system",
		Title:       "Complete Distributed System Vision",
		Description: "Complete distributed system vision at 500M+ DAU with all services",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box blue" data-tip="500M DAU generating 58K RPS baseline, 580K RPS at 10x peak. HTTP/2 with gzip. Mobile clients prefetch feed on app launch."><span class="d-step">1</span> Clients (iOS / Android / Web) &#8212; 500M DAU <span class="d-metric throughput">58K RPS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Latency-based routing directs users to nearest of 3 regions (us-east-1, eu-west-1, ap-south-1). Health checks failover in 30s. TTL 60s."><span class="d-step">2</span> Route 53 (Latency-based DNS) &#8594; nearest region <span class="d-metric latency">1-5ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="400+ edge PoPs. Origin Shield per region. Cache-Control: immutable for media. 95%+ cache hit rate for images. Saves ~$200K/mo in origin bandwidth."><span class="d-step">3</span> CloudFront (400+ edge PoPs) &#8212; static + API caching <span class="d-metric latency">&lt;50ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="WAF blocks SQL injection and XSS. Rate limiter: 100 req/min per user via token bucket. ALB path-based routing to microservices. mTLS between services."><span class="d-step">4</span> API Gateway (ALB + WAF + rate limiter) <span class="d-metric latency">1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box green" data-tip="Auth, profiles, follow graph. Sharded Postgres. 10 ECS tasks.">User Svc <span class="d-status active"></span></div>
    <div class="d-box green" data-tip="Post CRUD. DynamoDB PK=user_id, SK=timestamp. 8 ECS tasks.">Post Svc <span class="d-status active"></span></div>
    <div class="d-box green" data-tip="Hybrid fan-out. Redis sorted sets. Most latency-critical service.">Feed Svc <span class="d-status active"></span></div>
    <div class="d-box green" data-tip="Pre-signed upload, Lambda resize. 200TB/day new media.">Media Svc <span class="d-status active"></span></div>
    <div class="d-box green" data-tip="Sharded counters, dedup. 100K+ WPS at peak.">Engagement Svc <span class="d-status active"></span></div>
    <div class="d-box green" data-tip="Kafka consumer, SNS push, SSE web. Rate limited 50/hr per user.">Notification Svc <span class="d-status active"></span></div>
    <div class="d-box green" data-tip="Elasticsearch full-text search, autocomplete. Explore feed via ML ranking.">Search Svc <span class="d-status active"></span></div>
  </div>
  <div class="d-arrow-down">&#8595; all services produce events &#8595;</div>
  <div class="d-box red" data-tip="3 brokers per region, 50 partitions per topic. Multi-region Kafka (MirrorMaker 2) for cross-region event replication. 7-day retention. Exactly-once semantics."><span class="d-step">5</span> Kafka (MSK) &#8212; Event Bus <span class="d-metric throughput">500K events/sec</span></div>
  <div class="d-label">Events: post_created, user_followed, post_liked, comment_added, story_expired</div>
  <div class="d-arrow-down">&#8595; consumed by downstream services &#8595;</div>
  <div class="d-row">
    <div class="d-box amber" data-tip="Global Tables replicate across 3 regions with &lt;1s lag. On-demand capacity auto-scales. Used for posts, likes, notifications.">DynamoDB Global Tables <span class="d-metric latency">3ms</span></div>
    <div class="d-box indigo" data-tip="Sharded by user_id. Primary in us-east-1 with cross-region read replicas. Used for users, follows, auth.">Postgres (sharded) <span class="d-metric latency">5ms</span></div>
    <div class="d-box red" data-tip="50-node cluster per region. Sorted sets for feeds, hashes for profiles. 4TB total across all users.">Redis Cluster (feed cache) <span class="d-metric latency">&lt;2ms</span></div>
    <div class="d-box amber" data-tip="73 PB/year growth. S3 Intelligent-Tiering. CRR to 3 regions. 11 nines durability.">S3 (media) <span class="d-metric size">200 TB/day</span></div>
    <div class="d-box purple" data-tip="3-node cluster. Inverted index for hashtags, user search. Autocomplete via edge n-grams.">Elasticsearch <span class="d-metric latency">10ms</span></div>
  </div>
</div>
<div class="d-legend">
  <span class="d-legend-item"><span class="d-legend-color blue"></span>Client</span>
  <span class="d-legend-item"><span class="d-legend-color purple"></span>Network / CDN / Search</span>
  <span class="d-legend-item"><span class="d-legend-color indigo"></span>Gateway / Relational DB</span>
  <span class="d-legend-item"><span class="d-legend-color green"></span>Microservices</span>
  <span class="d-legend-item"><span class="d-legend-color red"></span>Cache / Event Bus</span>
  <span class="d-legend-item"><span class="d-legend-color amber"></span>NoSQL / Object Storage</span>
</div>
<div class="d-caption">Full distributed system at 500M+ DAU: 7 microservices, 5 storage technologies, 3 regions. Total infra cost: ~$1M+/mo. Request path from client to response: &lt;100ms for cached reads, &lt;300ms for writes.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-security-layers",
		Title:       "Security & Content Moderation Layers",
		Description: "Multi-layer security architecture for Instagram at scale.",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Upload Path Security</div>
      <div class="d-flow-v">
        <div class="d-box red">Rate Limiting (100 uploads/day per user)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">PhotoDNA Hash Check (CSAM, &lt;100ms)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">ML Image Classifier (NSFW/violence, ~200ms async)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">NLP Caption/Comment Toxicity (~50ms)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Publish to Feed + Explore Index</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Account Security</div>
      <div class="d-flow-v">
        <div class="d-box purple">2FA (TOTP authenticator)</div>
        <div class="d-box purple">Login anomaly detection (device + geo)</div>
        <div class="d-box purple">Session: JWT 15-min access + 7-day refresh</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Bot Prevention</div>
      <div class="d-flow-v">
        <div class="d-box indigo">Follow/unfollow churn detection</div>
        <div class="d-box indigo">Comment spam (duplicate text detection)</div>
        <div class="d-box indigo">Behavioral bot scoring (ML model)</div>
        <div class="d-box indigo">Shadowban for suspicious accounts</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-content-moderation",
		Title:       "Content Moderation Pipeline",
		Description: "Multi-stage content moderation from upload to human review.",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">New Post / Comment Uploaded</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red">Stage 1: PhotoDNA Hash Match (known bad content, &lt;100ms)</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Match</div>
      <div class="d-box red">Block + Report to NCMEC (legally required)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">No match</div>
      <div class="d-box green">Continue &#8595;</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber">Stage 2: ML Classification (NSFW/violence/spam, ~200ms async)</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Score &gt; 0.9</div>
      <div class="d-box red">Auto-remove + notify creator</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Score 0.5 &#8212; 0.9</div>
      <div class="d-box amber">Flag for human review (SQS queue)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Score &lt; 0.5</div>
      <div class="d-box green">Approved &#8594; publish to feed</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Stage 3: User Reports &#8594; Human Review Queue (4-24h SLA)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray">Stage 4: Appeal Process &#8594; Second Reviewer (24-48h SLA)</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-analytics-pipeline",
		Title:       "Analytics & Engagement Pipeline",
		Description: "Event-driven analytics from user actions to real-time dashboards and batch reports.",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-row">
    <div class="d-box blue">post_viewed</div>
    <div class="d-box blue">post_liked</div>
    <div class="d-box blue">post_commented</div>
    <div class="d-box blue">story_viewed</div>
    <div class="d-box blue">follow_action</div>
  </div>
  <div class="d-arrow-down">&#8595; Kafka (MSK)</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Real-time (Flink)</div>
      <div class="d-flow-v">
        <div class="d-box purple">Apache Flink (sliding window aggregation)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">Redis counters (engagement rate, trending score)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Grafana dashboards (&lt;2 min lag)</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Batch (Firehose)</div>
      <div class="d-flow-v">
        <div class="d-box amber">Kinesis Firehose &#8594; S3 (Parquet)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray">Athena / Spark (ad-hoc analytics)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray">Retention cohorts, A/B test results</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">ML Training</div>
      <div class="d-flow-v">
        <div class="d-box indigo">Feature store (engagement signals)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo">Feed ranking model training (daily)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo">Explore personalization model</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-monitoring-slos",
		Title:       "Monitoring & SLO Dashboard",
		Description: "Service level objectives and monitoring metrics for Instagram at scale.",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">SLOs by Feature</div>
      <div class="d-flow-v">
        <div class="d-box green">Feed: 99.99% available, p99 &lt;300ms</div>
        <div class="d-box green">Upload: 99.9% success rate</div>
        <div class="d-box green">Image Load: p99 &lt;200ms (CDN)</div>
        <div class="d-box blue">Notifications: 99.5% within 5s</div>
        <div class="d-box blue">Search: p99 &lt;100ms</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Key Alerts</div>
      <div class="d-flow-v">
        <div class="d-box red">Feed 5xx &gt; 0.1% &#8594; page on-call</div>
        <div class="d-box red">Kafka lag &gt; 100K &#8594; scale consumers</div>
        <div class="d-box amber">Redis &gt; 80% memory &#8594; scale cluster</div>
        <div class="d-box amber">CDN hit ratio &lt; 80% &#8594; check TTL</div>
        <div class="d-box amber">DynamoDB throttles &#8594; auto-scale</div>
      </div>
    </div>
  </div>
</div>
<div class="d-flow-v">
  <div class="d-box indigo">Distributed Tracing (OpenTelemetry): Feed request &#8594; Feed Svc (5ms) &#8594; Redis (0.5ms) &#8594; Post Svc gRPC (10ms) &#8594; User Svc gRPC (5ms) &#8594; Celebrity merge (15ms) = 36ms total</div>
</div>`,
	})
}

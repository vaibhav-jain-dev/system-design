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
        <div class="d-box amber" data-tip="S3 Standard: 11 nines durability. Pre-signed URL upload bypasses app server. Use Intelligent-Tiering at &gt;10TB to move cold objects to cheaper tiers.">S3 Standard &#8212; <span class="d-metric cost">$23/TB/mo</span> <div class="d-tag blue">S3</div></div>
        <div class="d-box purple" data-tip="$85/10TB outbound from edge. 400+ PoPs. Origin Shield adds ~$10/TB but cuts origin load by 90%. Cache-Control: max-age=31536000 for immutable assets.">CloudFront CDN &#8212; <span class="d-metric cost">$85/10TB</span></div>
      </div>
    </div>
  </div>
</div>
<div class="d-box blue" data-tip="$500/mo total: Compute $145 + Storage $250 + CDN $85 + misc. Scales linearly to ~$5K/mo at 10M users before optimization.">Total: &#8776;<span class="d-metric cost">$500/mo</span> for &lt;1M users <div class="d-tag green">MVP budget</div></div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-data-model",
		Title:       "Data Model (Entity Relationships)",
		Description: "Entity relationship diagram for users, posts, follows, likes, and comments",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow" style="justify-content:center; margin-bottom:12px;">
  <div class="d-number"><div class="d-number-value">5</div><div class="d-number-label">core tables</div></div>
  <div class="d-number"><div class="d-number-value">400B</div><div class="d-number-label">follows rows</div></div>
  <div class="d-number"><div class="d-number-value">100K+</div><div class="d-number-label">peak WPS (likes)</div></div>
  <div class="d-number"><div class="d-number-value">1,150</div><div class="d-number-label">post writes/sec</div></div>
</div>
<!-- Primary chain: users → posts → comments with visual arrows -->
<div style="display:flex; align-items:center; gap:0; overflow-x:auto; padding:4px 0; width:100%;">
  <div style="flex-shrink:0; min-width:175px;">
    <div class="d-entity" data-tip="Sharded by user_id at 50M+ users. Postgres BIGSERIAL gives 9.2 quintillion IDs.">
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
  <div class="d-rel-card" style="flex-shrink:0; padding:0 10px; align-self:center;">
    <div class="d-rel-card-label">1:N</div>
    <div class="d-rel-line-arrow" style="min-width:44px;"></div>
  </div>
  <div style="flex-shrink:0; min-width:185px;">
    <div class="d-entity" data-tip="Hottest table by write volume. Composite index (user_id, created_at DESC) critical for profile page queries.">
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
    <div style="font-size:0.65rem; color:var(--text-muted); margin-top:3px; text-align:center;">idx: (user_id, created_at DESC)</div>
  </div>
  <div class="d-rel-card" style="flex-shrink:0; padding:0 10px; align-self:center;">
    <div class="d-rel-card-label">1:N</div>
    <div class="d-rel-line-arrow" style="min-width:44px;"></div>
  </div>
  <div style="flex-shrink:0; min-width:185px;">
    <div class="d-entity" data-tip="Ordered by (post_id, created_at DESC) for threaded display. GIN index on text for spam detection.">
      <div class="d-entity-header red">comments <span class="d-metric throughput">~6K WPS</span></div>
      <div class="d-entity-body">
        <div class="pk">id BIGSERIAL</div>
        <div class="fk">user_id BIGINT &#8594; users.id</div>
        <div class="fk">post_id BIGINT &#8594; posts.id</div>
        <div>text TEXT NOT NULL</div>
        <div class="idx idx-btree">created_at TIMESTAMP</div>
      </div>
    </div>
    <div style="font-size:0.65rem; color:var(--text-muted); margin-top:3px; text-align:center;">idx: (post_id, created_at DESC)</div>
  </div>
</div>
<!-- M:N junction tables below with relationship labels -->
<div style="display:flex; gap:20px; margin-top:14px; justify-content:center; flex-wrap:wrap;">
  <div style="min-width:195px; text-align:center;">
    <div style="display:flex; align-items:center; justify-content:center; gap:6px; margin-bottom:6px; font-size:0.68rem; color:var(--text-muted); font-weight:600;">
      <span style="color:var(--blue); font-weight:700;">users</span>
      <div class="d-rel-line-arrow" style="min-width:18px;"></div>
      <span style="font-family:var(--font-mono); color:var(--indigo); font-size:0.7rem; font-weight:700;">M:N</span>
      <div class="d-rel-line-arrow" style="min-width:18px; transform:rotate(180deg);"></div>
      <span style="color:var(--blue); font-weight:700;">users</span>
    </div>
    <div class="d-entity" data-tip="Composite PK (follower_id, followee_id) prevents duplicates. Reverse index on followee_id enables 'who follows me?' in O(log N). At 2B users, largest table by row count.">
      <div class="d-entity-header purple">follows <span class="d-metric size">~400B rows</span></div>
      <div class="d-entity-body">
        <div class="pk fk">follower_id BIGINT &#8594; users.id</div>
        <div class="pk fk">followee_id BIGINT &#8594; users.id</div>
        <div>created_at TIMESTAMP</div>
      </div>
    </div>
    <div style="font-size:0.65rem; color:var(--text-muted); margin-top:3px;">idx: followee_id (reverse lookup)</div>
  </div>
  <div style="min-width:195px; text-align:center;">
    <div style="display:flex; align-items:center; justify-content:center; gap:6px; margin-bottom:6px; font-size:0.68rem; color:var(--text-muted); font-weight:600;">
      <span style="color:var(--blue); font-weight:700;">users</span>
      <div class="d-rel-line-arrow" style="min-width:18px;"></div>
      <span style="font-family:var(--font-mono); color:var(--indigo); font-size:0.7rem; font-weight:700;">M:N</span>
      <div class="d-rel-line-arrow" style="min-width:18px;"></div>
      <span style="color:var(--green); font-weight:700;">posts</span>
    </div>
    <div class="d-entity" data-tip="Composite PK deduplicates likes. At viral scale, sharded Redis counters aggregate across 100 shards every 5s before writing to DynamoDB.">
      <div class="d-entity-header amber">likes <span class="d-metric throughput">100K+ WPS peak</span></div>
      <div class="d-entity-body">
        <div class="pk fk">user_id BIGINT &#8594; users.id</div>
        <div class="pk fk">post_id BIGINT &#8594; posts.id</div>
        <div>created_at TIMESTAMP</div>
      </div>
    </div>
    <div style="font-size:0.65rem; color:var(--text-muted); margin-top:3px;">idx: post_id (count query)</div>
  </div>
</div>
<div class="d-legend" style="margin-top:12px;">
  <span class="d-legend-item"><span class="d-legend-color blue"></span>Identity</span>
  <span class="d-legend-item"><span class="d-legend-color green"></span>Content</span>
  <span class="d-legend-item"><span class="d-legend-color purple"></span>M:N self-join</span>
  <span class="d-legend-item"><span class="d-legend-color amber"></span>Engagement (high write)</span>
  <span class="d-legend-item"><span class="d-legend-color red"></span>User-generated text</span>
</div>
<div class="d-caption">Five tables handle 99% of Instagram's core data. Arrows show FK relationships. The follows table is the largest (~400B rows); likes sees the highest peak write throughput during viral events.</div>`,
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
  <div class="d-box blue" data-tip="Pre-signed URL expires in 15 min. Client uploads directly to S3 — bypasses all app servers. Max 50MB. S3 returns ETag on success."><span class="d-step">1</span> W1: Client uploads photo to S3 via pre-signed URL <span class="d-metric latency">~200ms</span> <span class="d-status active"></span></div>
  <div class="d-label">S3 | Direct upload, no app server bottleneck | Failure: S3 503 &#8594; client retries</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box blue" data-tip="Only metadata (S3 key, caption, media_type) sent to app server. Photo already in S3. ALB routes to any healthy ECS task."><span class="d-step">2</span> W2: Client calls POST /media with S3 key <span class="d-metric latency">1ms ALB</span></div>
  <div class="d-label">ALB &#8594; ECS | Metadata only, photo already in S3</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="ACID INSERT into posts table. Returns new post_id (BIGSERIAL). Triggers async Kafka event post_created for fan-out workers."><span class="d-step">3</span> W3: App server writes post record to Postgres <span class="d-metric latency">5-15ms</span> <span class="d-status active"></span></div>
  <div class="d-label">RDS Postgres | ACID for post creation | Failure: DB failover to standby (30s)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red" data-tip="DEL profile:{user_id} key. Prevents stale post count on profile page. Sub-ms. If Redis is down, cache miss on next read — safe fallback."><span class="d-step">4</span> W4: Invalidate user profile cache <span class="d-metric latency">&lt;1ms</span></div>
  <div class="d-label">ElastiCache Redis | Profile shows latest post count | Failure: cache miss, serve stale</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="CDN URL format: https://cdn.instagram.com/{post_id}/{size}.jpg. Client can immediately display the photo and share the link."><span class="d-step">5</span> W5: Return post_id + CDN URL to client <span class="d-status active"></span></div>
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
  <div class="d-box blue" data-tip="GET /api/v1/feed?cursor=last_post_id&amp;limit=20. Cursor-based pagination avoids OFFSET degradation at depth. ALB routes to any ECS task."><span class="d-step">1</span> R1: Client requests GET /api/v1/feed <span class="d-status active"></span></div>
  <div class="d-label">ALB &#8594; ECS | App server handles feed assembly | Failure: ALB drains unhealthy node</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red" data-tip="ZREVRANGE feed:{user_id} 0 19. Returns 20 post_ids in timestamp order. 90%+ hit rate for active users. 1-2ms. Miss triggers DB path."><span class="d-step">2</span> R2: Check Redis for cached feed <span class="d-metric latency">1-2ms</span> <span class="d-status active"></span></div>
  <div class="d-label">ElastiCache Redis | 90%+ cache hit for active users | Failure: cache miss, fall through</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="Only on 10% cache miss. JOIN follows + posts WHERE followee_id IN (...) ORDER BY created_at DESC LIMIT 20. Index on (user_id, created_at). Timeout 5s."><span class="d-step">3</span> R3: Query Postgres: followees then recent posts <span class="d-metric latency">20-50ms</span></div>
  <div class="d-label">RDS Postgres | Fan-out-on-read: assemble feed at query time | Failure: slow query, timeout 5s</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red" data-tip="SET feed:{user_id} {post_ids} EX 300. 5-min TTL balances freshness vs DB load. If Redis full, LRU evicts least-recently-accessed feeds."><span class="d-step">4</span> R4: Cache assembled feed in Redis (TTL 5 min) <span class="d-metric latency">&lt;1ms</span></div>
  <div class="d-label">ElastiCache Redis | Next request hits cache | Failure: Redis full, evict LRU</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="CDN URL per image. CloudFront checks edge cache first. On miss, fetches from Origin Shield (one per region), then S3 origin. Hit rate: 95%+."><span class="d-step">5</span> R5: Client loads images via CDN URLs <span class="d-metric latency">&lt;50ms global</span> <span class="d-status active"></span></div>
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
  <div class="d-box blue" data-tip="Client receives pre-signed URL from API server (15-min expiry). Uploads directly to S3 over HTTPS. Bypasses all app server bandwidth constraints."><span class="d-step">1</span> Client (iOS / Android / Web) <span class="d-metric latency">~200ms upload</span></div>
  <div class="d-arrow-down">&#8595; pre-signed URL upload</div>
  <div class="d-box amber" data-tip="Raw uploads bucket. Lifecycle policy: delete after 48h once processing completes. Versioning disabled to save cost. Server-side encryption (SSE-S3)."><span class="d-step">2</span> S3 Ingest Bucket (raw uploads) <div class="d-tag blue">S3</div></div>
  <div class="d-arrow-down">&#8595; S3 Event Notification</div>
  <div class="d-box green" data-tip="Lambda triggered by S3 ObjectCreated event. Uses Pillow/Sharp for resizing. Generates 4 sizes + WebP variant per size. Timeout: 30s. Memory: 1GB."><span class="d-step">3</span> Lambda: Image Processing Pipeline <span class="d-metric latency">2-5s processing</span></div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Thumbnail</div>
      <div class="d-box purple" data-tip="Used for grid views and feed previews. Served to slow connections. WebP saves ~30% vs JPEG.">150&#215;150 <span class="d-metric size">~15KB</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Small</div>
      <div class="d-box purple" data-tip="Mobile feed default. Good balance of quality and bandwidth. WebP served when Accept header includes image/webp.">320&#215;320 <span class="d-metric size">~40KB</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Medium</div>
      <div class="d-box purple" data-tip="High-DPI mobile screens and tablet views. Most common served size. Quality 85 JPEG, 80 WebP.">640&#215;640 <span class="d-metric size">~120KB</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Full</div>
      <div class="d-box purple" data-tip="Desktop and zoom views. Only served on explicit user action. Quality 90 JPEG to preserve photo detail.">1080&#215;1080 <span class="d-metric size">~300KB</span></div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; all sizes saved</div>
  <div class="d-box amber" data-tip="Permanent media bucket. S3 Intelligent-Tiering moves objects to Infrequent Access after 30 days. SSE-S3 encryption. Public read disabled — served only through CloudFront."><span class="d-step">4</span> S3 Media Bucket (processed images, WebP + JPEG) <div class="d-tag blue">S3</div> <span class="d-metric size">200 TB/day</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">CDN Distribution</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="One Origin Shield per region collapses all edge PoP misses into a single origin fetch. Reduces S3 GET requests by 90%. Costs ~$10/TB but saves far more."><span class="d-step">5</span> CloudFront Origin Shield (regional cache) <span class="d-metric latency">5-10ms</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple" data-tip="400+ PoPs worldwide. Cache-Control: max-age=31536000 means images are cached at edge indefinitely. TTL reset only on CDN invalidation."><span class="d-step">6</span> CloudFront Edge PoPs (400+ locations) <span class="d-metric latency">&lt;20ms</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue" data-tip="Browser receives image in under 50ms globally for cached content. HTTP/2 multiplexing loads multiple images simultaneously."><span class="d-step">7</span> User device <span class="d-metric latency">&lt;50ms globally</span> <span class="d-status active"></span></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">CDN Config</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="Immutable assets — URLs include content hash. 1-year TTL means no origin hits after first request per PoP.">Cache-Control: public, max-age=31536000 <div class="d-tag green">recommended</div></div>
          <div class="d-box green" data-tip="CloudFront Lambda@Edge checks Accept header. Serves WebP to Chrome/Firefox (30% smaller), JPEG to Safari/older clients.">Accept header &#8594; WebP vs JPEG negotiation</div>
          <div class="d-box green" data-tip="Without Origin Shield, each of 400 PoPs would fetch separately from S3 on first miss — 400x origin load. Shield collapses this to 1 fetch per region.">Origin Shield &#8594; single origin fetch per region</div>
          <div class="d-box amber" data-tip="CRR to eu-west-1 and ap-south-1 means European and Asian users get sub-50ms reads from local S3, not cross-ocean S3 fetches.">S3 CRR to eu-west-1, ap-south-1 for global reads</div>
        </div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Upload to CDN-ready: ~5s (Lambda resize). Once at CDN edge: &lt;50ms globally. Pre-signed upload bypasses all app servers — 200TB/day never touches ECS.</div>`,
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
  <div class="d-box blue" data-tip="1M–50M users. Mobile-first HTTP/2. Client prefetches 20 feed items on app launch to hide latency.">Client (iOS / Android / Web) <span class="d-metric throughput">~5K RPS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Origin Shield reduces origin load by 90%. Cache-Control: max-age=31536000 for immutable assets. Handles all static media without touching app servers.">CloudFront + Origin Shield <span class="d-metric latency">&lt;50ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Layer 7 path-based routing. /feed &#8594; Feed Service, /media &#8594; Monolith, /api/v1/* &#8594; Monolith. Health checks every 10s.">ALB <span class="d-metric latency">1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Monolith (ECS &#215;8)</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="8 ECS tasks, each 4 vCPU 8GB RAM. Handles all API routes except feed. Auto-scales to 16 tasks at CPU&gt;70%. Stateless — sessions in Redis.">API Server (Django/FastAPI) <div class="d-tag green">stateless</div></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box indigo" data-tip="All writes go to primary. db.r6g.2xlarge at this scale. Trigger: add replicas when primary CPU &gt;60%.">Postgres Primary <span class="d-status active"></span></div>
            <div class="d-box indigo" data-tip="Read replicas serve all SELECT queries: feeds, profiles, posts. Lag typically &lt;100ms from primary. Failover to replica if primary fails.">Read Replica &#215;2 <span class="d-metric latency">5ms reads</span></div>
          </div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Feed Service (ECS &#215;4)</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="Consumes Kafka post_created events. For each post, reads follower list, writes post_id to each follower's Redis sorted set in batches of 100.">Fan-out Worker <span class="d-metric throughput">100K ZADD/sec</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red" data-tip="3-node cluster for redundancy (1 primary + 2 replicas). Each user's feed is a sorted set of post_ids. ZREVRANGE 0 19 returns latest 20 posts in &lt;2ms.">Redis Cluster (3-node) <span class="d-metric latency">&lt;2ms</span></div>
          <div class="d-label">Sorted sets per user</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Async Processing</div>
        <div class="d-flow-v">
          <div class="d-box amber" data-tip="Buffers like and comment writes during peak traffic. SQS standard queue — at-least-once delivery. Batch size 10 per Lambda invocation.">SQS (likes/comments queue) <span class="d-metric throughput">10K msg/sec</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="Batch writes 10 likes at once to DynamoDB. Reduces write units by 10x vs individual writes. Triggered by SQS event.">Lambda (batch writer) <span class="d-metric latency">50ms batch</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber" data-tip="S3 ObjectCreated event triggers Lambda. Generates 4 sizes (150, 320, 640, 1080px) + WebP variants. Stores all 8 files to media bucket.">S3 &#8594; Lambda &#8594; 4 image sizes <span class="d-metric latency">2-5s</span></div>
        </div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Stage 2 separates the feed fan-out into a dedicated service while keeping the monolith for all other routes. Cost: ~$12.5K/mo for 1M–50M users.</div>`,
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
      <div class="d-box green" data-tip="Owns auth, profiles, follow graph. gRPC API. Shards Postgres by user_id at 50M+ users. ~10 ECS tasks.">User Service <span class="d-status active"></span></div>
      <div class="d-label">Profiles, auth, follows</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box indigo" data-tip="Sharded by user_id. B-tree indexes on (follower_id, followee_id) and reverse (followee_id). ACID for auth writes. Read replicas for follow graph queries.">Postgres (sharded) <span class="d-metric latency">5ms</span></div>
      <div class="d-label">ACID for auth, B-tree for follows</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green" data-tip="Owns post CRUD and media metadata. DynamoDB chosen over Postgres for infinite horizontal scale and single-digit ms at any throughput.">Post Service <span class="d-status active"></span></div>
      <div class="d-label">Posts, media metadata</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="PK=user_id enables O(1) lookup of all posts by a user. SK=timestamp gives time-ordered results without an index. DAX cache for hot posts.">DynamoDB <span class="d-metric latency">3ms</span> <div class="d-tag blue">on-demand</div></div>
      <div class="d-label">PK=user_id, SK=timestamp</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green" data-tip="Hybrid fan-out engine. Consumes Kafka post_created events. Writes to sorted sets for normal users (&lt;100K followers), merges celebrity posts at read time.">Feed Service <span class="d-status active"></span></div>
      <div class="d-label">Feed computation + cache</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box red" data-tip="50-node cluster across 3 AZs. One sorted set per user: ZADD feed:{uid} {ts} {post_id}. ZREVRANGE for reads. 8KB per user &#215; 500M = 4TB total.">Redis Cluster <span class="d-metric latency">&lt;2ms</span> <span class="d-metric size">4TB total</span></div>
      <div class="d-label">Sorted sets, sub-ms reads</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green" data-tip="Issues pre-signed S3 URLs for direct upload. Triggers Lambda resize pipeline on upload. Manages CDN URLs. 200TB/day never touches app servers.">Media Service <span class="d-status active"></span></div>
      <div class="d-label">Upload, process, serve</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="S3: 11 nines durability, Intelligent-Tiering cuts cost 40%. Lambda resize: 4 sizes + WebP, ~3s end-to-end. S3 CRR to 3 regions.">S3 + Lambda <div class="d-tag blue">S3</div> <span class="d-metric size">200 TB/day</span></div>
      <div class="d-label">11 nines durability</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green" data-tip="Handles likes (dedup via composite PK), comments, and shares. Sharded counters for viral posts — 100 Redis shards aggregate every 5s.">Engagement Service <span class="d-status active"></span></div>
      <div class="d-label">Likes, comments, shares</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="Composite PK (user_id, post_id) prevents duplicate likes. On-demand capacity auto-scales to handle viral bursts. Aggregated counts written every 5s.">DynamoDB <span class="d-metric throughput">100K+ WPS</span></div>
      <div class="d-label">Sharded counters, 100K+ WPS</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green" data-tip="Consumes Kafka events, deduplicates (1h window), batches likes into 'alice and 12 others liked'. Sends via SNS (APNs/FCM) and SSE gateway.">Notification Svc <span class="d-status active"></span></div>
      <div class="d-label">Push, in-app, email</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="SQS buffers notification jobs for retry on failure. DynamoDB stores history with TTL=30 days. PK=user_id, SK=timestamp for ordered retrieval.">SQS + DynamoDB <span class="d-metric latency">&lt;500ms push</span></div>
      <div class="d-label">Async, TTL auto-cleanup</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box green" data-tip="Kafka consumer indexes approved posts into Elasticsearch. Explore feed generated by hourly Spark job using engagement_rate = (likes+comments)/impressions.">Search/Explore <span class="d-status active"></span></div>
      <div class="d-label">Hashtags, user search</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box purple" data-tip="3-node cluster. Inverted index for hashtag and caption full-text search. Edge n-gram analyzer for autocomplete. Query latency &lt;10ms at p99.">Elasticsearch <span class="d-metric latency">10ms</span></div>
      <div class="d-label">Full-text, autocomplete</div>
    </div>
  </div>
</div>
<div class="d-legend">
  <span class="d-legend-item"><span class="d-legend-color green"></span>Microservice (owns its DB)</span>
  <span class="d-legend-item"><span class="d-legend-color indigo"></span>Relational (ACID, B-tree)</span>
  <span class="d-legend-item"><span class="d-legend-color amber"></span>NoSQL / Object Storage</span>
  <span class="d-legend-item"><span class="d-legend-color red"></span>Cache (Redis)</span>
  <span class="d-legend-item"><span class="d-legend-color purple"></span>Search</span>
</div>
<div class="d-caption">Each service owns exactly one primary datastore — no shared databases. DB choice is driven by access pattern: ACID for auth, DynamoDB for high-write engagement, Redis for sub-ms feed reads.</div>`,
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
		HTML:        `<div class="d-flow" style="justify-content:center; margin-bottom:12px;">
  <div class="d-number"><div class="d-number-value">500K</div><div class="d-number-label">events/sec (Kafka)</div></div>
  <div class="d-number"><div class="d-number-value">&lt;500ms</div><div class="d-number-label">push delivery</div></div>
  <div class="d-number"><div class="d-number-value">&lt;50ms</div><div class="d-number-label">web real-time</div></div>
  <div class="d-number"><div class="d-number-value">90%</div><div class="d-number-label">like batching reduction</div></div>
</div>
<div class="d-flow-v">
  <div class="d-group" style="width:100%">
    <div class="d-group-title">Step 1 — Event Producers</div>
    <div class="d-row" style="flex-wrap:wrap; gap:8px;">
      <div class="d-box green" data-tip="Fires on every new post. Triggers fan-out to all followers' notification feeds."><span class="d-step">1a</span> Post Svc: post_created <span class="d-status active"></span></div>
      <div class="d-box green" data-tip="Highest volume event source. Viral posts generate 1M+ like events. Batched to avoid notification spam."><span class="d-step">1b</span> Engagement Svc: post_liked, comment_added <span class="d-metric throughput">100K events/sec</span></div>
      <div class="d-box green" data-tip="Fired when user A follows user B. Lower volume than likes but fan-out is large for celebrity accounts."><span class="d-step">1c</span> User Svc: user_followed</div>
    </div>
  </div>
  <div class="d-arrow-down" data-tip="All event types published to single Kafka topic. Partitioned by target_user_id so all events for a user go to the same partition.">&#8595; all events publish to Kafka</div>
  <div class="d-box red" style="width:100%; text-align:center;" data-tip="Partitioned by target_user_id for ordering guarantees. 50 partitions, 7-day retention. Consumer lag monitored via CloudWatch."><span class="d-step">2</span> <strong>Kafka (MSK)</strong> — notification topic <span class="d-metric throughput">500K msg/sec</span> <span class="d-tag red">50 partitions</span></div>
  <div class="d-arrow-down">&#8595; ECS consumer group</div>
  <div class="d-box indigo" style="width:100%; text-align:center;" data-tip="Deduplicates using (event_type, target_user, source_entity) key in Redis with 1h TTL. Batches likes into 'alice and 12 others liked your post' style messages."><span class="d-step">3</span> <strong>Notification Service</strong> (ECS consumers) <span class="d-metric latency">50-200ms</span> — dedup + batch similar events</div>
  <div class="d-arrow-down">&#8595; fan-out to delivery channels</div>
  <div class="d-branch" style="width:100%;">
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">Mobile Push</div>
        <div class="d-flow-v">
          <div class="d-box amber" data-tip="Fan-out to platform-specific endpoints. Handles token rotation and delivery receipts. $1 per 1M publishes."><span class="d-step">4a</span> SNS Platform App <span class="d-metric cost">$1/1M msgs</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box blue" data-tip="Apple Push Notification Service. Requires device token, certificate auth. 100-500ms delivery.">APNs (iOS) <span class="d-metric latency">100-500ms</span></div>
            <div class="d-box green" data-tip="Firebase Cloud Messaging. Handles Android + web push. Higher throughput ceiling than APNs.">FCM (Android) <span class="d-metric latency">100-500ms</span></div>
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
          <div class="d-box blue" data-tip="EventSource API reconnects automatically on drop. No custom client code needed.">Browser client <span class="d-status active"></span></div>
        </div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">In-App Badge</div>
        <div class="d-flow-v">
          <div class="d-box red" data-tip="Atomic INCR for badge count. Reset to 0 on notification tab open. O(1) operation, sub-ms latency."><span class="d-step">4c</span> Redis INCR badge:{user_id} <span class="d-metric latency">&lt;1ms</span></div>
          <div class="d-arrow-down">&#8595; persist history</div>
          <div class="d-box amber" data-tip="TTL 30 days auto-deletes old notifications. PK=user_id, SK=timestamp. On-demand capacity handles traffic spikes."><span class="d-step">5</span> DynamoDB (notification history, TTL 30d) <span class="d-metric latency">3ms</span></div>
        </div>
      </div>
    </div>
  </div>
</div>
<div class="d-legend">
  <div class="d-legend-item"><div class="d-legend-color green"></div>Event producers</div>
  <div class="d-legend-item"><div class="d-legend-color red"></div>Kafka stream</div>
  <div class="d-legend-item"><div class="d-legend-color indigo"></div>Processing (dedup/batch)</div>
  <div class="d-legend-item"><div class="d-legend-color amber"></div>Mobile push</div>
  <div class="d-legend-item"><div class="d-legend-color purple"></div>Web real-time</div>
</div>
<div class="d-caption">Rate limit: max 50 push/hr per user. Like batching reduces push volume 90% during viral events. End-to-end: &lt;500ms push, &lt;50ms web.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-peak-traffic",
		Title:       "Peak Traffic Defense Layers",
		Description: "Three-layer defense strategy for handling 10x peak traffic spikes",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box red" data-tip="10x traffic spike happens during live events (Super Bowl, World Cup). 580K RPS sustained for ~2 hours. Viral posts generate 1M likes/sec on a single post_id.">&#9889; Peak Traffic: <span class="d-metric throughput">580K RPS</span> (10&#215; normal) <span class="d-status error"></span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Layer 1: Absorb</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="Pre-warm CDN before known events (Super Bowl). Push popular content to all 400+ PoPs 30 min before kickoff. Eliminates cold-cache miss spike."><span class="d-step">1</span> CDN pre-warming (push to all PoPs) <div class="d-tag green">proactive</div></div>
          <div class="d-box purple" data-tip="Token bucket: 100 req/min per user_id. Implemented at ALB WAF level. Returns 429 with Retry-After header. Protects origin from burst."><span class="d-step">2</span> Rate limiter (100 req/min per user) <span class="d-metric throughput">blocks 80% of burst</span></div>
          <div class="d-box purple" data-tip="ECS auto-scales on CPU&gt;70%. Scale-out completes in ~3 min. Pre-scale for known events: set min tasks to 3x normal 30 min before."><span class="d-step">3</span> Auto-scaling ECS (CPU/memory triggers) <span class="d-metric latency">3 min scale-out</span></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Layer 2: Buffer</div>
        <div class="d-flow-v">
          <div class="d-box amber" data-tip="SQS absorbs write spikes. At 1M likes/sec, SQS queues messages without dropping. Lambda batch consumer writes 10 at a time to DynamoDB."><span class="d-step">1</span> SQS write queue (likes/comments) <span class="d-metric throughput">1M msg/sec</span></div>
          <div class="d-box amber" data-tip="100 Redis shards per post. Each shard handles 10K INCR/sec. Aggregate every 5s. Prevents hot-key bottleneck on viral post_id."><span class="d-step">2</span> Sharded counters (100 shards) <span class="d-metric throughput">1M INCR/sec</span></div>
          <div class="d-box amber" data-tip="Raise feed TTL from 5 min to 15 min during peaks. Reduces DB fan-out reads by 3x. Slightly staler feed is acceptable vs total unavailability."><span class="d-step">3</span> Increased cache TTL during peaks <span class="d-metric latency">5 min &#8594; 15 min TTL</span></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Layer 3: Degrade</div>
        <div class="d-flow-v">
          <div class="d-box gray" data-tip="Hystrix/Resilience4j circuit breaker. If error rate &gt;50% for 10s, open circuit, return fallback. Prevents cascading failure across services."><span class="d-step">1</span> Circuit breakers (fail fast) <span class="d-status error"></span></div>
          <div class="d-box gray" data-tip="Serve Redis-cached feed even if 30 min stale. Better to show slightly old feed than 503 error. Users rarely notice 30 min staleness."><span class="d-step">2</span> Serve cached/stale feed <span class="d-status active"></span></div>
          <div class="d-box gray" data-tip="If personalized feed fails, fall back to pre-computed trending feed (top 100 posts globally). Always available from Redis, no user-specific computation."><span class="d-step">3</span> Fallback to trending (pre-computed) <span class="d-status active"></span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-flow">
    <div class="d-number"><div class="d-number-value">58K</div><div class="d-number-label">Normal RPS</div></div>
    <div class="d-number"><div class="d-number-value">580K</div><div class="d-number-label">Peak RPS</div></div>
    <div class="d-number"><div class="d-number-value">1M</div><div class="d-number-label">Likes/sec (viral)</div></div>
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
  <div class="d-box purple" data-tip="Latency-based routing: directs each user to nearest region. Health check polling every 30s. Failover in &lt;60s if region goes down. TTL=60s for client DNS cache.">Route 53 (Global DNS &#8212; Latency-based routing) <span class="d-metric latency">1-5ms DNS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">us-east-1 (Americas) <div class="d-tag blue">primary</div></div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="400+ PoPs in Americas. All US media requests served from CloudFront edge — never hits S3 origin for cached content.">CloudFront <span class="d-metric latency">&lt;20ms US</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple" data-tip="WAF rules, rate limiting, path-based routing to microservices. TLS termination.">ALB <span class="d-metric latency">1ms</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="Full set of 7 microservices. ~100 ECS tasks total. Auto-scales independently per service.">ECS Services <span class="d-metric throughput">~30K RPS</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box indigo" data-tip="Postgres primary for user/auth data. Cross-region read replicas in eu-west-1 and ap-south-1.">Postgres <span class="d-metric latency">5ms</span></div>
            <div class="d-box amber" data-tip="DynamoDB Global Tables: us-east-1 is the write primary. Replicates to other regions in &lt;1s.">DynamoDB <span class="d-metric latency">3ms</span></div>
            <div class="d-box red" data-tip="50-node Redis cluster. Not replicated cross-region — feeds are regional. 4TB per cluster.">Redis <span class="d-metric latency">&lt;2ms</span></div>
          </div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">eu-west-1 (Europe) <div class="d-tag blue">GDPR</div></div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="EU PoPs. GDPR requires EU user data stays in eu-west-1. Geolocation override in Route 53 forces EU users to this region regardless of latency.">CloudFront <span class="d-metric latency">&lt;20ms EU</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple" data-tip="Same ALB config as us-east-1. GDPR data boundary enforced at ALB layer — no cross-region user data forwarding.">ALB <span class="d-metric latency">1ms</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="Full microservices stack. EU user writes go to eu-west-1 Postgres primary — not replicated to US for GDPR compliance.">ECS Services <span class="d-metric throughput">~15K RPS</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box indigo" data-tip="Postgres primary for EU users. Isolated from us-east-1 Postgres for GDPR. Postgres logical replication for non-PII data only.">Postgres <span class="d-metric latency">5ms</span></div>
            <div class="d-box amber" data-tip="DynamoDB Global Tables receive replication from us-east-1. On-demand capacity scales independently.">DynamoDB <span class="d-metric latency">3ms</span></div>
            <div class="d-box red" data-tip="EU Redis cluster. EU user feeds pre-computed and cached here. Not cross-region — EU users always read from this cluster.">Redis <span class="d-metric latency">&lt;2ms</span></div>
          </div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">ap-south-1 (India/Asia)</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="APAC PoPs. India has 200M+ Instagram users — dedicated region reduces cross-ocean latency from 200ms to 20ms.">CloudFront <span class="d-metric latency">&lt;20ms APAC</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple" data-tip="Same ALB config. Higher mobile traffic mix — connection pooling tuned for LTE network patterns (higher reconnect rate).">ALB <span class="d-metric latency">1ms</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" data-tip="Full microservices. ap-south-1 added when India user latency exceeded 150ms from us-east-1.">ECS Services <span class="d-metric throughput">~13K RPS</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box indigo" data-tip="Postgres cross-region read replica from us-east-1. Read replica lag: typically &lt;100ms. Failover promotes replica to primary in 60s.">Postgres <span class="d-metric latency">5ms</span></div>
            <div class="d-box amber" data-tip="DynamoDB Global Tables third replica. S3 CRR also mirrors media bucket here for edge-origin proximity.">DynamoDB <span class="d-metric latency">3ms</span></div>
            <div class="d-box red" data-tip="APAC Redis cluster. Indian user feeds cached here. 4TB for regional user base.">Redis <span class="d-metric latency">&lt;2ms</span></div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-label">DynamoDB Global Tables + S3 CRR replicate across all regions</div>
</div>
<div class="d-legend">
  <span class="d-legend-item"><span class="d-legend-color purple"></span>Network / CDN / DNS</span>
  <span class="d-legend-item"><span class="d-legend-color green"></span>Application Services</span>
  <span class="d-legend-item"><span class="d-legend-color indigo"></span>Relational DB (Postgres)</span>
  <span class="d-legend-item"><span class="d-legend-color amber"></span>NoSQL (DynamoDB Global Tables)</span>
  <span class="d-legend-item"><span class="d-legend-color red"></span>Cache (regional Redis)</span>
</div>
<div class="d-caption">Three active-active regions each handle full traffic independently. DynamoDB Global Tables and S3 CRR replicate data. Redis and Postgres are regional — no cross-region cache or write replication for these layers.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-stories-lifecycle",
		Title:       "Stories Lifecycle",
		Description: "Stories lifecycle with 24-hour TTL auto-deletion and view tracking",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-row">
    <div class="d-box blue" data-tip="Client uploads story media directly to S3 via pre-signed URL (same as photos). Then POSTs metadata to /api/v1/stories with media_s3_key."><span class="d-step">1</span> User uploads story <span class="d-metric latency">~500ms</span></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber" data-tip="Story media stored in S3 with a separate lifecycle policy: delete at 25h (1h buffer for timezone edge cases). Objects stored as story/{story_id}/media.jpg."><span class="d-step">2</span> S3 (media file) <div class="d-tag blue">S3</div> TTL 25h</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green" data-tip="DynamoDB TTL field set to now+86400 (Unix timestamp). DynamoDB background job deletes expired items within ~48h of TTL. No application code needed for cleanup."><span class="d-step">3</span> DynamoDB (TTL = 24h) <span class="d-metric latency">3ms write</span></div>
  </div>
  <div class="d-label">S3 lifecycle policy deletes media after 25h (timezone buffer)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box red" data-tip="ZADD stories:{user_id} {timestamp} {story_id}. ZREVRANGE returns stories in chronological order. ZEXPIRE set to 24h. Each story_id is 8 bytes."><span class="d-step">4</span> Redis sorted set (per-user story feed) <span class="d-metric latency">&lt;1ms</span></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple" data-tip="SETBIT views:{story_id} {viewer_user_id} 1. O(1) read and write. 500M users = 62.5MB per story bitmap. BITCOUNT gives exact view count."><span class="d-step">4</span> Redis bitmap (view tracking per story_id) <span class="d-metric size">~62MB/story</span></div>
  </div>
  <div class="d-label">Stories are time-ordered, not ranked &#8212; bitmap gives O(1) "has user seen this?" check</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box blue" data-tip="Client renders circular avatar rings. Red ring = unseen stories. Stories sorted: unseen first (bitmap check), then by recency within each bucket."><span class="d-step">5</span> Client renders story ring <span class="d-status active"></span></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box gray" data-tip="GETBIT views:{story_id} {current_user_id} returns 0 (unseen) or 1 (seen). All unseen stories from all followed users shown before any seen story."><span class="d-step">5</span> Unseen stories first, then by recency</div>
  </div>
  <div class="d-arrow-down">&#8595; after 24 hours</div>
  <div class="d-box red" data-tip="DynamoDB TTL fires at exactly 24h. Redis sorted sets expire via separate TTL on the key. S3 lifecycle fires at 25h. All three cleanup systems are independent."><span class="d-step">6</span> DynamoDB TTL auto-deletes &#8594; Redis keys expire &#8594; Story gone <span class="d-status error"></span></div>
<div class="d-caption">Stories are ephemeral by design: three independent TTL mechanisms (DynamoDB TTL, Redis key expiry, S3 lifecycle) each handle cleanup. No single-point-of-failure in deletion path.</div>
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
          <div class="d-box green" data-tip="All three event types feed the search index. post_liked events update engagement_score for ranking. Kafka consumer group: search-indexer, lag monitored."><span class="d-step">1</span> Kafka events: post_created, post_liked, comment_added <span class="d-metric throughput">500K events/sec</span></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple" data-tip="Two-stage filter: PhotoDNA hash check (&lt;100ms) then ML classifier (&lt;200ms). Score &gt;0.9 auto-rejected, 0.5-0.9 goes to human queue, &lt;0.5 approved."><span class="d-step">2</span> ML Content Classifier <span class="d-metric latency">&lt;200ms</span></div>
          <div class="d-label">NSFW/spam filter &#8594; reject or flag for review</div>
          <div class="d-arrow-down">&#8595; approved content</div>
          <div class="d-box amber" data-tip="3-node Elasticsearch cluster. Index per content type: posts (caption, hashtags, media_type), users (username, bio). Edge n-gram for autocomplete. Refresh 1s."><span class="d-step">3</span> Elasticsearch Index (posts, hashtags, users) <span class="d-metric latency">10ms query</span></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Explore Feed Generation</div>
        <div class="d-flow-v">
          <div class="d-box indigo" data-tip="Reads 7 days of Kafka events from S3 (Parquet via Firehose). Computes engagement_rate per post per interest cluster. Runs hourly on EMR cluster."><span class="d-step">1</span> Spark Batch Job (hourly) <span class="d-metric latency">~30 min runtime</span></div>
          <div class="d-label">engagement_rate = (likes + comments) / impressions</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue" data-tip="Collaborative filtering: matrix factorization on user-post interaction matrix. Users assigned to interest clusters (travel, food, tech, etc). Top-100 posts per cluster."><span class="d-step">2</span> ML Ranking Model <div class="d-tag blue">collaborative filtering</div></div>
          <div class="d-label">Collaborative filtering: users who liked X also liked Y</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red" data-tip="Pre-computed explore feed per interest cluster (100 clusters). LRANGE explore:{cluster_id} 0 49 returns top 50 posts in &lt;2ms. Refreshed hourly."><span class="d-step">3</span> Redis (pre-computed explore feeds per interest cluster) <span class="d-metric latency">&lt;2ms read</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box blue" data-tip="Returns personalized grid from user's interest cluster Redis key. Falls back to global trending if cluster not yet computed.">GET /explore &#8594; personalized trending grid <span class="d-status active"></span></div>
    <div class="d-box blue" data-tip="Elasticsearch prefix query with edge n-gram. Returns top 10 suggestions in &lt;10ms. Results ranked by post engagement_score.">GET /search?q=sunset &#8594; Elasticsearch autocomplete <span class="d-metric latency">&lt;10ms</span></div>
    <div class="d-box blue" data-tip="Term query on hashtags field. Results sorted by engagement_rate DESC. Paginated with search_after for deep pagination without OFFSET degradation.">GET /tags/travel &#8594; hashtag feed (ES query) <span class="d-metric latency">&lt;10ms</span></div>
  </div>
</div>
<div class="d-caption">Explore is intentionally batch-computed (hourly) — freshness matters less than relevance. Search is near-real-time via Kafka streaming into Elasticsearch with 1s index refresh.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-distributed-counter",
		Title:       "Distributed Counter Architecture",
		Description: "Sharded counter architecture for handling 1M likes per second",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML:        `<div class="d-flow-v">
  <div class="d-box blue" data-tip="Super Bowl / World Cup goal: 1M likes/sec on a single celebrity post. A single Redis key would hit ~100K ops/sec limit — need sharding."><span class="d-step">1</span> 1M concurrent like requests <span class="d-metric throughput">1M req/sec</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="DynamoDB conditional write: put_item if (user_id, post_id) not exists. Returns ConditionalCheckFailedException for duplicates. Dedup cost: ~3ms per check."><span class="d-step">2</span> Deduplicate: check likes table PK (user_id, post_id) <span class="d-metric latency">3ms</span></div>
  <div class="d-arrow-down">&#8595; new likes only</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Shard 0</div>
      <div class="d-box amber" data-tip="Key: like_count:{post_id}:shard:0. Each shard handles ~10K INCR/sec. Spread across different Redis nodes."><span class="d-step">3</span> INCRBY 1 <span class="d-metric latency">&lt;1ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Shard 1</div>
      <div class="d-box amber" data-tip="Key: like_count:{post_id}:shard:1. Distributed across Redis cluster nodes by key hash."><span class="d-step">3</span> INCRBY 1 <span class="d-metric latency">&lt;1ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">...</div>
      <div class="d-box gray" data-tip="100 shards total. Shard selected by: hash(user_id) % 100. Ensures even distribution across shards.">N shards <span class="d-metric throughput">10K/sec each</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Shard 99</div>
      <div class="d-box amber" data-tip="Key: like_count:{post_id}:shard:99. Last shard completes the 100-shard fan-out."><span class="d-step">3</span> INCRBY 1 <span class="d-metric latency">&lt;1ms</span></div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; aggregate every 5s</div>
  <div class="d-box green" data-tip="Background job: MGET all 100 shard keys, SUM values, SET like_count:{post_id} {total} EX 5. Runs every 5s. Displayed count lags by max 5s."><span class="d-step">4</span> Total = SUM(all shards) &#8594; cached in Redis (TTL 5s) <span class="d-metric latency">5s max lag</span> <span class="d-status active"></span></div>
<div class="d-flow">
  <div class="d-number"><div class="d-number-value">100</div><div class="d-number-label">Shards</div></div>
  <div class="d-number"><div class="d-number-value">10K</div><div class="d-number-label">INCR/sec/shard</div></div>
  <div class="d-number"><div class="d-number-value">1M</div><div class="d-number-label">Total likes/sec</div></div>
  <div class="d-number"><div class="d-number-value">5s</div><div class="d-number-label">Display lag</div></div>
</div>
<div class="d-caption">100 shards handle 1M likes/sec on a single post. Deduplication in DynamoDB prevents double-counting. Displayed count lags by up to 5s — acceptable UX tradeoff for 1000x throughput gain.</div>
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
          <div class="d-box green" data-tip="2 ECS tasks, all routes on one codebase. Scale to 8 tasks before splitting. Django/FastAPI handles 500 req/sec per task.">Monolith (Django/FastAPI) <div class="d-tag green">start here</div></div>
          <div class="d-box indigo" data-tip="db.r6g.large: 2 vCPU, 16GB RAM. Split to read replicas when primary CPU &gt;60% sustained.">Single Postgres (RDS) <span class="d-metric latency">5-20ms</span></div>
          <div class="d-box red" data-tip="t4g.medium: 4GB RAM. Session storage and feed cache. Upgrade to cluster when memory hits 80%.">Single Redis <span class="d-metric latency">&lt;1ms</span></div>
          <div class="d-box amber" data-tip="S3 for uploads, CloudFront for delivery. No Lambda resizing yet — store only original. Add resize pipeline at 10K uploads/day.">S3 + CloudFront</div>
          <div class="d-label"><span class="d-metric cost">~$500/mo</span> | 1 team</div>
          <div class="d-label">Pull-based feed (SQL query)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Growth (1M&#8212;50M users)</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="Extract Feed Service when feed p99 &gt;200ms. Monolith handles everything else. Trigger: feed latency alarm on CloudWatch.">Monolith + Feed Service <div class="d-tag blue">first split</div></div>
          <div class="d-box indigo" data-tip="2 read replicas absorb 10x read load. Primary handles only writes (~1K WPS). Replicas serve feeds, profiles, post details.">Postgres + 2 read replicas <span class="d-metric latency">5ms reads</span></div>
          <div class="d-box red" data-tip="3-node cluster: 1 primary per shard + replicas. Hash-slot sharding. Handles 100K ops/sec total.">Redis cluster (3 nodes) <span class="d-metric throughput">100K ops/sec</span></div>
          <div class="d-box amber" data-tip="Lambda triggered on S3 ObjectCreated. Generates 4 sizes + WebP. Total 8 files per upload. 2-5s processing. Cost: ~$0.20 per 1K uploads.">S3 &#8594; Lambda &#8594; 4 sizes <span class="d-metric latency">2-5s resize</span></div>
          <div class="d-label"><span class="d-metric cost">~$12.5K/mo</span> | 3 teams</div>
          <div class="d-label">Hybrid fan-out + SQS buffering</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Scale (50M&#8212;500M users)</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="User, Post, Feed, Media, Engagement, Notification, Search. Each owns its DB. gRPC inter-service. Strangler fig pattern from monolith.">7 Microservices (ECS) <div class="d-tag blue">event-driven</div></div>
          <div class="d-box indigo" data-tip="Postgres sharded by user_id (10 shards). DynamoDB replaces Postgres for posts and likes — handles 100K WPS without sharding complexity.">Sharded Postgres + DynamoDB <span class="d-metric latency">3-5ms</span></div>
          <div class="d-box red" data-tip="50 nodes across 3 AZs. 4TB total for feeds. r6g.xlarge instances. 500K ops/sec capacity.">Redis cluster (50 nodes) <span class="d-metric size">4TB feeds</span></div>
          <div class="d-box gray" data-tip="Kafka MSK: 3 brokers, 50 partitions per topic. Decouples all services. 500K events/sec. Enables event sourcing and replay.">Kafka (MSK) event bus <span class="d-metric throughput">500K events/sec</span></div>
          <div class="d-label"><span class="d-metric cost">~$150K/mo</span> | 10 teams</div>
          <div class="d-label">Event-driven, gRPC, sharded counters</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Global (500M+ users)</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="3 regions: us-east-1, eu-west-1, ap-south-1. Each region is fully independent. Cross-region traffic only for writes that need global consistency.">Multi-region microservices <div class="d-tag blue">active-active</div></div>
          <div class="d-box indigo" data-tip="DynamoDB Global Tables replicate across all 3 regions with &lt;1s lag. Last-writer-wins conflict resolution. On-demand capacity per region.">DynamoDB Global Tables <span class="d-metric latency">&lt;1s replication</span></div>
          <div class="d-box red" data-tip="Separate 50-node cluster per region. Feeds are not cross-region replicated — users get their regional feed. No cross-region Redis traffic.">Per-region Redis clusters <span class="d-metric size">4TB per region</span></div>
          <div class="d-box purple" data-tip="Latency-based routing directs each user to nearest region. Health check failover in 30s. TTL 60s. Geolocation override for GDPR (EU users stay in eu-west-1).">Route 53 latency routing <span class="d-metric latency">1-5ms DNS</span></div>
          <div class="d-label"><span class="d-metric cost">~$1M+/mo</span> | 30+ teams</div>
          <div class="d-label">CRDTs, S3 CRR, multi-region Kafka</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-flow">
    <div class="d-box blue" data-tip="$500/mo, 1 team, monolith. Measure everything. Don't optimize prematurely.">0&#8212;1M <span class="d-metric cost">$500/mo</span></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue" data-tip="$12.5K/mo, 3 teams. Extract Feed Service when p99 &gt; 200ms. Add read replicas.">1M&#8212;50M <span class="d-metric cost">$12.5K/mo</span></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue" data-tip="$150K/mo, 10 teams. Full microservices via strangler fig. Kafka event bus.">50M&#8212;500M <span class="d-metric cost">$150K/mo</span></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue" data-tip="$1M+/mo, 30+ teams. Multi-region active-active. DynamoDB Global Tables.">500M+ <span class="d-metric cost">$1M+/mo</span></div>
  </div>
  <div class="d-flow">
    <div class="d-number"><div class="d-number-value">4</div><div class="d-number-label">Scale Phases</div></div>
    <div class="d-number"><div class="d-number-value">2,000x</div><div class="d-number-label">Cost Growth</div></div>
    <div class="d-number"><div class="d-number-value">30x</div><div class="d-number-label">Team Growth</div></div>
  </div>
  <div class="d-label">Key triggers: feed latency &gt; 200ms &#8594; add fan-out | write RPS &gt; 10K &#8594; shard DB | user latency &gt; 150ms &#8594; add region</div>
</div>
<div class="d-caption">Scale each layer only when a measured bottleneck demands it. The trigger metrics are more important than the architecture stages — measure first, then migrate.</div>`,
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
        <div class="d-box red" data-tip="Token bucket per user_id. 100 uploads/day prevents spam accounts from flooding the platform. Stored in Redis. Returns 429 with Retry-After when exceeded."><span class="d-step">1</span> Rate Limiting (100 uploads/day per user) <span class="d-metric throughput">blocks spam</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="PhotoDNA perceptual hash database from Microsoft. Compares upload hash to NCMEC database of known CSAM. Match = instant block + mandatory NCMEC report. Legally required."><span class="d-step">2</span> PhotoDNA Hash Check (CSAM, &lt;100ms) <span class="d-metric latency">&lt;100ms</span> <span class="d-status error"></span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Async ML classifier runs after upload is visible. Score 0-1: &gt;0.9 auto-remove, 0.5-0.9 human review queue, &lt;0.5 approved. 98.5% precision to minimize false positives."><span class="d-step">3</span> ML Image Classifier (NSFW/violence, ~200ms async) <span class="d-metric latency">200ms async</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="BERT-based NLP model. Hate speech, harassment, spam patterns. Checks caption + first 50 comments. 50ms inline for captions, async for comments."><span class="d-step">4</span> NLP Caption/Comment Toxicity (~50ms) <span class="d-metric latency">50ms inline</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Content approved for distribution. Fanout to follower feeds via Kafka. Indexed in Elasticsearch for search and explore."><span class="d-step">5</span> Publish to Feed + Explore Index <span class="d-status active"></span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Account Security</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="TOTP (Google Authenticator / Authy). 6-digit code, 30s window. Stored as encrypted TOTP secret in user table. Required for API access after trusted device expiry.">2FA (TOTP authenticator) <div class="d-tag green">recommended</div></div>
        <div class="d-box purple" data-tip="Flags login from new device + new country. Risk score from device fingerprint + IP geolocation. Score &gt;0.7 triggers SMS challenge before access granted.">Login anomaly detection (device + geo)</div>
        <div class="d-box purple" data-tip="Access token: JWT, 15-min TTL, signed with RS256. Refresh token: opaque 256-bit token in DynamoDB, 7-day TTL. Rotation on every refresh. Revoke list in Redis.">Session: JWT 15-min access + 7-day refresh <span class="d-metric latency">&lt;1ms verify</span></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Bot Prevention</div>
      <div class="d-flow-v">
        <div class="d-box indigo" data-tip="Sliding window: &gt;50 follow/unfollow actions in 1 hour = suspicious. Redis ZADD + ZCOUNT pattern. Immediate flag, re-evaluate after 24h.">Follow/unfollow churn detection <span class="d-metric throughput">50/hr limit</span></div>
        <div class="d-box indigo" data-tip="MinHash LSH for near-duplicate comment detection. Same user posting &gt;80% similar comments triggers spam review. Bloom filter for exact dedup.">Comment spam (duplicate text detection) <span class="d-metric latency">5ms check</span></div>
        <div class="d-box indigo" data-tip="Behavioral signals: typing speed, scroll patterns, tap coordinates (humans have variance, bots don't). Score 0-100. &gt;70 = bot, requires CAPTCHA.">Behavioral bot scoring (ML model)</div>
        <div class="d-box indigo" data-tip="Shadowban: content visible to the account but hidden from others. Account doesn't know they're banned — prevents evasion by account-switching. Lifted after 30 days with no violations.">Shadowban for suspicious accounts <span class="d-metric latency">30 day TTL</span></div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Defense-in-depth: rate limiting stops volume attacks, PhotoDNA stops known CSAM, ML stops novel NSFW, NLP stops harassment, bot detection stops coordinated abuse. Each layer is independent.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "ig-content-moderation",
		Title:       "Content Moderation Pipeline",
		Description: "Multi-stage content moderation from upload to human review.",
		ContentFile: "problems/instagram",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="Entry point for all new content: photos, videos, captions, and comments all go through the same moderation pipeline. Volume: ~1,150 uploads/sec + 6K comments/sec."><span class="d-step">1</span> New Post / Comment Uploaded <span class="d-metric throughput">~7K items/sec</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red" data-tip="Microsoft PhotoDNA: perceptual hash of image compared to NCMEC database of known CSAM. &lt;100ms synchronous check. Block is instant — content never becomes visible."><span class="d-step">2</span> Stage 1: PhotoDNA Hash Match (known bad content, &lt;100ms) <span class="d-metric latency">&lt;100ms</span></div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Match</div>
      <div class="d-box red" data-tip="Mandatory NCMEC report under 18 U.S.C. § 2258A. Content blocked instantly. Account flagged for review. Zero tolerance — no appeal process.">Block + Report to NCMEC (legally required) <span class="d-status error"></span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">No match</div>
      <div class="d-box green" data-tip="99.99%+ of content reaches this point. Proceed to ML classifier while content is published optimistically.">Continue &#8595; <span class="d-status active"></span></div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="Custom ResNet-based model. Inputs: image pixels + caption text. Outputs: score 0-1 per category (NSFW, graphic violence, spam, hate symbols). Runs async — content already published."><span class="d-step">3</span> Stage 2: ML Classification (NSFW/violence/spam, ~200ms async) <span class="d-metric latency">200ms async</span></div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Score &gt; 0.9</div>
      <div class="d-box red" data-tip="Automated removal within 200ms of upload. Push notification to creator: 'Your post was removed for violating Community Guidelines.' Can appeal within 30 days.">Auto-remove + notify creator <span class="d-status error"></span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Score 0.5 &#8212; 0.9</div>
      <div class="d-box amber" data-tip="SQS queue to human review team. Priority: higher score = higher priority. SLA: 4h for score &gt;0.7, 24h for 0.5-0.7. Content stays visible during review.">Flag for human review (SQS queue) <span class="d-metric latency">4-24h SLA</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Score &lt; 0.5</div>
      <div class="d-box green" data-tip="Content fully approved. Remains in feed. Future user reports can still trigger Stage 3 re-review.">Approved &#8594; publish to feed <span class="d-status active"></span></div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Crowdsourced reporting. 3+ reports from different accounts triggers priority review queue. Reviewer sees full context: account history, post details, report reasons."><span class="d-step">4</span> Stage 3: User Reports &#8594; Human Review Queue (4-24h SLA) <span class="d-metric latency">4-24h SLA</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" data-tip="Creator appeals automated removal. Second reviewer is more senior. Reviews full context + appeal statement. Decision is final. 24-48h SLA. ~15% of appeals succeed."><span class="d-step">5</span> Stage 4: Appeal Process &#8594; Second Reviewer (24-48h SLA) <span class="d-metric latency">24-48h SLA</span></div>
<div class="d-caption">Four-stage pipeline balances safety and false-positive rate. Automated stages handle 99.9% of cases. Human review is reserved for edge cases and appeals — can't be automated away.</div>
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
    <div class="d-box blue" data-tip="Impression + dwell time event. Fired when post visible for &gt;300ms. Key metric for engagement_rate denominator.">post_viewed <span class="d-status active"></span></div>
    <div class="d-box blue" data-tip="Like toggle event. Includes user_id, post_id, timestamp, action (like/unlike). Highest volume engagement signal.">post_liked <span class="d-status active"></span></div>
    <div class="d-box blue" data-tip="Comment created event. Includes sentiment score from inline NLP. Used for engagement_rate numerator.">post_commented <span class="d-status active"></span></div>
    <div class="d-box blue" data-tip="Story view with dwell time and completion rate. Completion rate (watched to end) is stronger signal than view count.">story_viewed <span class="d-status active"></span></div>
    <div class="d-box blue" data-tip="Follow and unfollow. Net follower growth rate is a key content quality signal. Unfollow after view = negative signal for explore ranking.">follow_action <span class="d-status active"></span></div>
  </div>
  <div class="d-arrow-down">&#8595; Kafka (MSK) <span class="d-metric throughput">500K events/sec</span></div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">Real-time (Flink)</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="5-minute sliding window aggregation. Computes trending_score = likes_last_5min / impressions_last_5min. Runs on 10-node Flink cluster on ECS."><span class="d-step">1</span> Apache Flink (sliding window aggregation) <span class="d-metric latency">&lt;2 min lag</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="ZADD trending:{category} {score} {post_id} for top trending posts. INCR engagement:{post_id} for real-time counters. Drives explore feed and trending topics."><span class="d-step">2</span> Redis counters (engagement rate, trending score) <span class="d-metric latency">&lt;1ms write</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Grafana boards: RPS, error rates, cache hit ratios, feed p99 latency. Refreshes every 30s. PagerDuty alerts for SLO breaches."><span class="d-step">3</span> Grafana dashboards <span class="d-metric latency">&lt;2 min lag</span></div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Batch (Firehose)</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Firehose buffers 5 min or 128MB, then flushes to S3 as Parquet. Partitioned by date/hour/event_type for efficient Athena queries."><span class="d-step">1</span> Kinesis Firehose &#8594; S3 (Parquet) <span class="d-metric latency">5 min buffer</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray" data-tip="Athena: serverless SQL over S3 Parquet, $5/TB scanned. Spark on EMR for complex ML feature computation. Both read the same S3 Parquet files."><span class="d-step">2</span> Athena / Spark (ad-hoc analytics)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray" data-tip="7-day and 30-day retention cohorts track user engagement over time. A/B test results compared using Mann-Whitney U test for statistical significance."><span class="d-step">3</span> Retention cohorts, A/B test results</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">ML Training</div>
      <div class="d-flow-v">
        <div class="d-box indigo" data-tip="Tecton/Feast feature store. Pre-computed user interest vectors, post embeddings, engagement history. Served at &lt;5ms for online inference."><span class="d-step">1</span> Feature store (engagement signals) <span class="d-metric latency">&lt;5ms serve</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo" data-tip="Daily retraining on last 30 days of events. Two-tower model: user tower + post tower. Dot product = relevance score. Deployed via SageMaker A/B rollout."><span class="d-step">2</span> Feed ranking model training (daily) <span class="d-latency">daily retrain</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo" data-tip="Hourly batch job assigns users to interest clusters (100 clusters: travel, food, tech, etc.). Top 100 posts per cluster precomputed into Redis."><span class="d-step">3</span> Explore personalization model <span class="d-metric latency">hourly batch</span></div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Three parallel consumers from Kafka: Flink for &lt;2min real-time dashboards, Firehose for durable batch analytics, and ML training pipeline for daily model retraining. All from the same event stream.</div>`,
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
        <div class="d-box green" data-tip="99.99% = 52 min downtime/year. Error budget: 5.26 min/month. Measured as: (total_minutes - error_minutes) / total_minutes. p99 &lt;300ms across all regions.">Feed: 99.99% available, p99 &lt;300ms <span class="d-metric latency">p99 &lt;300ms</span> <span class="d-status active"></span></div>
        <div class="d-box green" data-tip="Upload success = S3 pre-signed URL issued + client upload completes + metadata written to DB. 99.9% = 8.7 hours downtime/year. 0.1% acceptable for transient S3 issues.">Upload: 99.9% success rate <span class="d-status active"></span></div>
        <div class="d-box green" data-tip="Image load measured at CDN edge. 95%+ cache hit rate means most images served in &lt;20ms. p99 &lt;200ms includes the 5% CDN misses that hit origin.">Image Load: p99 &lt;200ms (CDN) <span class="d-metric latency">p99 &lt;200ms</span> <span class="d-status active"></span></div>
        <div class="d-box blue" data-tip="99.5% within 5s = 1.82 days downtime/year. Looser SLO because push notifications are best-effort (APNs/FCM don't guarantee delivery). SSE web notifications are stricter.">Notifications: 99.5% within 5s <span class="d-metric latency">&lt;5s push</span></div>
        <div class="d-box blue" data-tip="Elasticsearch p99 &lt;100ms. Autocomplete is critical UX — users expect instant results. Achieved via: 3-node cluster, index warm-up, query caching.">Search: p99 &lt;100ms <span class="d-metric latency">p99 &lt;100ms</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Key Alerts</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="CloudWatch alarm: error rate over 5-min window. &gt;0.1% = PagerDuty page to on-call. Runbook: check ALB 5xx, ECS task health, DB connections. Auto-rollback on deploy.">Feed 5xx &gt; 0.1% &#8594; page on-call <span class="d-status error"></span></div>
        <div class="d-box red" data-tip="Kafka consumer lag &gt;100K = fan-out falling behind. Action: scale ECS consumer tasks from 4 to 8. Root cause: viral post with 1M followers.">Kafka lag &gt; 100K &#8594; scale consumers <span class="d-status error"></span></div>
        <div class="d-box amber" data-tip="Redis memory &gt;80% = risk of OOM evictions. Action: scale from r6g.xlarge to r6g.2xlarge (doubles memory). Lead time: 5 min for ElastiCache scale-up.">Redis &gt; 80% memory &#8594; scale cluster <span class="d-metric size">80% threshold</span></div>
        <div class="d-box amber" data-tip="CDN hit ratio &lt;80% means too many origin requests. Root causes: TTL too short, cache-busting parameters in URLs, new content surge. Fix: increase TTL or pre-warm CDN.">CDN hit ratio &lt; 80% &#8594; check TTL</div>
        <div class="d-box amber" data-tip="DynamoDB throttles = provisioned capacity exceeded. DynamoDB auto-scaling responds in 2 min. Immediate fix: switch to on-demand mode for burst tolerance.">DynamoDB throttles &#8594; auto-scale</div>
      </div>
    </div>
  </div>
</div>
<div class="d-flow">
  <div class="d-number"><div class="d-number-value">99.99%</div><div class="d-number-label">Feed SLO</div></div>
  <div class="d-number"><div class="d-number-value">52 min</div><div class="d-number-label">Downtime Budget/yr</div></div>
  <div class="d-number"><div class="d-number-value">5</div><div class="d-number-label">Key SLOs</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box indigo" data-tip="OpenTelemetry traces propagated via W3C Trace-Context headers. Sampled at 1% for normal traffic, 100% for errors. Stored in Jaeger / AWS X-Ray. 15-day retention.">Distributed Tracing (OpenTelemetry): Feed request &#8594; Feed Svc (5ms) &#8594; Redis (0.5ms) &#8594; Post Svc gRPC (10ms) &#8594; User Svc gRPC (5ms) &#8594; Celebrity merge (15ms) = 36ms total</div>
</div>
<div class="d-caption">SLOs define the reliability contract. Error budgets enable teams to balance feature velocity vs reliability. When budget is consumed, freeze feature releases until SLO is met.</div>`,
	})
}

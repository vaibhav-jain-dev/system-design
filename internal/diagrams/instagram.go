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
        <div class="d-box green">POST /api/v1/media &#8594; upload photo, returns media_id</div>
        <div class="d-box green">GET /api/v1/feed &#8594; paginated feed (cursor-based)</div>
        <div class="d-box green">POST /api/v1/follow/{user_id} &#8594; follow/unfollow</div>
        <div class="d-box green">POST /api/v1/like/{post_id} &#8594; toggle like</div>
        <div class="d-box green">POST /api/v1/comment/{post_id} &#8594; add comment</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P1 — Important</div>
      <div class="d-flow-v">
        <div class="d-box blue">POST /api/v1/stories &#8594; 24hr ephemeral content</div>
        <div class="d-box blue">GET /api/v1/explore &#8594; trending + personalized</div>
        <div class="d-box blue">GET /api/v1/notifications &#8594; SSE/WebSocket</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P2 — Nice to Have</div>
      <div class="d-flow-v">
        <div class="d-box gray">POST /api/v1/messages/{user_id} &#8594; DMs</div>
        <div class="d-box gray">POST /api/v1/reels &#8594; short video</div>
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
        <div class="d-box green">Availability: 99.99% (52 min downtime/yr)</div>
        <div class="d-box green">Latency: Feed load &lt; 200ms p99</div>
        <div class="d-box blue">Scale: 2B MAU, 500M DAU</div>
        <div class="d-box blue">Read:Write ratio: ~10:1</div>
        <div class="d-box amber">Consistency: Eventual for feed, strong for likes</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Back-of-Envelope Math</div>
      <div class="d-flow-v">
        <div class="d-box purple">500M DAU &#215; 10 loads/day = 5B feed req/day &#8776; 58K RPS</div>
        <div class="d-box purple">100M photos/day &#8776; 1,150 uploads/sec (5x peak = 5,750/sec)</div>
        <div class="d-box purple">Likes/comments add 10x writes &#8594; total ~6K write RPS</div>
        <div class="d-box amber">100M/day &#215; 2MB avg = 200 TB/day new media</div>
        <div class="d-box amber">Year 1 storage: ~73 PB (200TB &#215; 365)</div>
      </div>
    </div>
  </div>
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
      <div class="d-box blue">Client (iOS / Android / Web)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box purple">CloudFront (CDN)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box purple">ALB (Load Balancer)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box green">ECS (Django / FastAPI Monolith)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-row">
        <div class="d-box indigo">Postgres (RDS)</div>
        <div class="d-box red">Redis (ElastiCache)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">MEDIA STORAGE</div>
      <div class="d-flow-v">
        <div class="d-box amber">S3 Bucket (Photos)</div>
        <div class="d-label">Pre-signed URL upload</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">CloudFront (CDN)</div>
        <div class="d-label">Edge-cached images</div>
      </div>
    </div>
  </div>
</div>`,
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
        <div class="d-box green">ECS Fargate: 2&#215; t3.large &#8212; $120/mo</div>
        <div class="d-box purple">ALB (Load Balancer) &#8212; $25/mo</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Storage</div>
      <div class="d-flow-v">
        <div class="d-box indigo">RDS Postgres db.r6g.large &#8212; $200/mo</div>
        <div class="d-box red">ElastiCache Redis t4g.medium &#8212; $50/mo</div>
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
    <div class="d-entity">
      <div class="d-entity-header blue">users</div>
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
    <div class="d-entity">
      <div class="d-entity-header green">posts</div>
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
    <div class="d-entity">
      <div class="d-entity-header purple">follows</div>
      <div class="d-entity-body">
        <div class="pk fk">follower_id BIGINT &#8594; users.id</div>
        <div class="pk fk">followee_id BIGINT &#8594; users.id</div>
        <div>created_at TIMESTAMP</div>
      </div>
    </div>
    <div style="font-size:0.68rem; color:var(--text-muted); margin-top:4px; text-align:center;">idx: followee_id (reverse lookup)</div>
    <div class="d-entity" style="margin-top: 0.75rem;">
      <div class="d-entity-header amber">likes</div>
      <div class="d-entity-body">
        <div class="pk fk">user_id BIGINT &#8594; users.id</div>
        <div class="pk fk">post_id BIGINT &#8594; posts.id</div>
        <div>created_at TIMESTAMP</div>
      </div>
    </div>
    <div style="font-size:0.68rem; color:var(--text-muted); margin-top:4px; text-align:center;">idx: post_id (count query)</div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header red">comments</div>
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
</div>`,
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
        <div class="d-box blue">Client</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">S3 (Pre-signed URL Upload)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">POST /media (with S3 key)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">ALB</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">ECS (App Server)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo">Postgres (Write post record)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">Redis (Invalidate cache)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray">Response (post_id + CDN URL)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">READ PATH (Feed Load)</div>
      <div class="d-flow-v">
        <div class="d-box blue">Client</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">ALB</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">ECS (App Server)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">Redis (Cache check)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo">Postgres (Cache miss)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">Redis (Cache feed, TTL 5m)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">CloudFront (Serve images)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray">Response (Feed JSON + CDN URLs)</div>
      </div>
    </div>
  </div>
</div>`,
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
        <div class="d-box green">User creates post</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Feed Service reads followers list</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Write post_id to each follower's Redis sorted set</div>
        <div class="d-label">ZADD feed:{user_id} {timestamp} {post_id}</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Trim to latest 500: ZREMRANGEBYRANK feed:{uid} 0 -501</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">&#10003; Feed pre-computed. Read = O(1) Redis lookup</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fan-out-on-Read (Celebrities &gt; 100K followers)</div>
      <div class="d-flow-v">
        <div class="d-box green">Celebrity creates post</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">Skip fan-out (would be 200M writes)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Mark as celebrity_post in Redis set</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">On feed read: merge pre-computed + celebrity posts</div>
        <div class="d-label">DB query: SELECT FROM posts WHERE user_id IN (celeb_ids)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Sort by timestamp, return top N</div>
      </div>
    </div>
  </div>
</div>`,
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
  <div class="d-box blue">Client (iOS / Android / Web)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">CloudFront (CDN)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">ALB (API Gateway)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box green">User Svc</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box indigo">Postgres</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green">Post Svc</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">DynamoDB</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green">Feed Svc</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box red">Redis Cluster</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green">Media Svc</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">S3</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green">Engagement Svc</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">DynamoDB</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray">Kafka (MSK) &#8212; Event Bus: post_created, user_followed, post_liked</div>
</div>`,
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
          <div class="d-box green">Post Svc: post_created</div>
          <div class="d-box green">Engagement Svc: post_liked, comment_added</div>
          <div class="d-box green">User Svc: user_followed</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-flow-v">
        <div class="d-label">all events flow into</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">Kafka (MSK) &#8212; notification topic</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">Notification Service (ECS consumers)</div>
  <div class="d-label">Dedup by (event_type, target_user, source_entity) &#8212; batch similar events</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">Mobile Push</div>
        <div class="d-flow-v">
          <div class="d-box amber">SNS Platform App</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-row">
            <div class="d-box blue">APNs (iOS)</div>
            <div class="d-box green">FCM (Android)</div>
          </div>
        </div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">Web Real-time</div>
        <div class="d-flow-v">
          <div class="d-box purple">SSE / WebSocket Gateway</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">Browser client</div>
        </div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">In-App Badge</div>
        <div class="d-flow-v">
          <div class="d-box red">Redis INCR badge:{user_id}</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">DynamoDB (notification history, TTL 30d)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-label">Rate limit: max 50 push notifications/hr per user &#8212; batch likes: "alice and 12 others liked your post"</div>
</div>`,
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
  <div class="d-box blue">Incoming Request (GET /feed, GET /post, GET /profile)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">L1 &#8212; Local In-Memory Cache</div>
        <div class="d-flow-v">
          <div class="d-box green">Caffeine / Guava LRU (per ECS task)</div>
          <div class="d-label">TTL: 30s | Size: 256MB per instance</div>
          <div class="d-label">Latency: &lt;1ms | Hit rate: ~60%</div>
          <div class="d-label">No network hop, fastest layer</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">L2 &#8212; Redis Cluster (Distributed)</div>
        <div class="d-flow-v">
          <div class="d-box red">ElastiCache Redis Cluster (50 nodes)</div>
          <div class="d-label">TTL: 5min (feed), 1hr (profiles)</div>
          <div class="d-label">Latency: 1-2ms | Hit rate: ~95%</div>
          <div class="d-label">Shared across all app instances</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">L3 &#8212; Database (Source of Truth)</div>
        <div class="d-flow-v">
          <div class="d-box indigo">Postgres / DynamoDB</div>
          <div class="d-label">Latency: 5-50ms</div>
          <div class="d-label">Always consistent, highest cost</div>
          <div class="d-label">Only hit on L1 + L2 miss (~5%)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; Invalidation Strategy</div>
  <div class="d-row">
    <div class="d-box amber">Write-through: DB write &#8594; delete L2 key &#8594; L1 expires via TTL</div>
    <div class="d-box amber">Kafka event &#8594; all instances invalidate L1 (pub/sub)</div>
  </div>
  <div class="d-label">Cache stampede protection: singleflight pattern + probabilistic early expiration (TTL &#215; random(0.8, 1.0))</div>
</div>`,
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
  <div class="d-box blue">Clients (iOS / Android / Web) &#8212; 500M DAU</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Route 53 (Latency-based DNS) &#8594; nearest region</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">CloudFront (400+ edge PoPs) &#8212; static + API caching</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">API Gateway (ALB + WAF + rate limiter)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box green">User Svc</div>
    <div class="d-box green">Post Svc</div>
    <div class="d-box green">Feed Svc</div>
    <div class="d-box green">Media Svc</div>
    <div class="d-box green">Engagement Svc</div>
    <div class="d-box green">Notification Svc</div>
    <div class="d-box green">Search Svc</div>
  </div>
  <div class="d-arrow-down">&#8595; all services produce events &#8595;</div>
  <div class="d-box red">Kafka (MSK) &#8212; Event Bus</div>
  <div class="d-label">Events: post_created, user_followed, post_liked, comment_added, story_expired</div>
  <div class="d-arrow-down">&#8595; consumed by downstream services &#8595;</div>
  <div class="d-row">
    <div class="d-box amber">DynamoDB Global Tables</div>
    <div class="d-box indigo">Postgres (sharded)</div>
    <div class="d-box red">Redis Cluster (feed cache)</div>
    <div class="d-box amber">S3 (media)</div>
    <div class="d-box purple">Elasticsearch</div>
  </div>
</div>`,
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

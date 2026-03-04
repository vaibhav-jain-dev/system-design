package diagrams

func registerFundamentals(r *Registry) {
	// -------------------------------------------------------
	// CDN
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-cdn-request-flow",
		Title:       "CDN Request Flow",
		Description: "Shows how a CDN request flows from users through edge PoPs, regional cache (Origin Shield), and origin server with latency at each layer",
		ContentFile: "fundamentals/networking/cdn",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-row">
    <div class="d-box blue">User (Sydney)</div>
    <div class="d-box blue">User (London)</div>
  </div>
  <div class="d-arrow-down">↓</div>
  <div class="d-row">
    <div class="d-box amber">Edge PoP - Sydney</div>
    <div class="d-box amber">Edge PoP - London</div>
  </div>
  <div class="d-label">HIT = Return cached response (~5-20ms)</div>
  <div class="d-arrow-down">↓ MISS</div>
  <div class="d-box green">Regional Cache (Origin Shield)</div>
  <div class="d-label">~30-60ms</div>
  <div class="d-arrow-down">↓ MISS</div>
  <div class="d-box purple">Origin Server (S3 / ALB)</div>
  <div class="d-label">~100-300ms</div>
</div>`,
	})

	// -------------------------------------------------------
	// CloudFront
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-cloudfront-multi-origin",
		Title:       "CloudFront Multi-Origin Architecture",
		Description: "CloudFront distribution with path-based behavior routing to S3, ALB, and MediaStore origins",
		ContentFile: "fundamentals/networking/cdn/cloudfront",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Users (global)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">CloudFront (400+ edge PoPs)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">/static/*</div>
      <div class="d-box green">S3 Origin</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">/api/*</div>
      <div class="d-box purple">ALB Origin</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">/media/*</div>
      <div class="d-box amber">MediaStore</div>
    </div>
  </div>
  <div class="d-label">Path-based behavior routing</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-cloudfront-event-hooks",
		Title:       "CloudFront Event Hooks",
		Description: "Request lifecycle through CloudFront showing viewer request/response and origin request/response event hooks for CF Functions and Lambda@Edge",
		ContentFile: "fundamentals/networking/cdn/cloudfront",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-box blue">Viewer Request</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber">CF Function</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box gray">Edge Cache</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue">Origin Request</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple">Lambda@Edge</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box green">Origin</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue">Origin Response</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple">Lambda@Edge</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box gray">Edge Cache</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue">Viewer Response</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber">CF Function</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box indigo">Client</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-cloudfront-origin-shield",
		Title:       "Origin Shield Architecture",
		Description: "Origin Shield as an additional caching layer between edge PoPs and origin, reducing origin load by 90%",
		ContentFile: "fundamentals/networking/cdn/cloudfront",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-row">
    <div class="d-box blue">Edge PoP (Sydney)</div>
    <div class="d-box blue">Edge PoP (London)</div>
    <div class="d-box blue">Edge PoP (Tokyo)</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Origin Shield (us-east-1)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green">Origin (S3/ALB)</div>
  <div class="d-label">Single cache layer reduces origin load by 90%</div>
</div>`,
	})

	// -------------------------------------------------------
	// Load Balancing
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-lb-multi-az",
		Title:       "Multi-AZ Load Balancing Architecture",
		Description: "Load balancer spanning multiple availability zones with Route 53 DNS and target groups in each AZ",
		ContentFile: "fundamentals/networking/load-balancing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Internet</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-box gray">Route 53 (DNS)</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-box green">ALB / NLB (spans all AZs)</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-row">
    <div class="d-group">
      <div class="d-group-title">AZ-1a</div>
      <div class="d-box indigo">Target Group EC2 / ECS</div>
    </div>
    <div class="d-group">
      <div class="d-group-title">AZ-1b</div>
      <div class="d-box indigo">Target Group EC2 / ECS</div>
    </div>
    <div class="d-group">
      <div class="d-group-title">AZ-1c</div>
      <div class="d-box indigo">Target Group EC2 / ECS</div>
    </div>
  </div>
  <div class="d-label">Health checks: ALB pings /health every 30s. Unhealthy target removed in ~60s. Entire AZ down: zonal shift via ARC.</div>
</div>`,
	})

	// -------------------------------------------------------
	// ALB
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-alb-path-routing",
		Title:       "ALB Path-Based Routing",
		Description: "ALB listener routing requests to different target groups based on URL path pattern with WAF rules evaluated before routing",
		ContentFile: "fundamentals/networking/load-balancing/alb",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Internet</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">ALB (Listener :443)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">/api/*</div>
      <div class="d-box green">API Target Group (ECS)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">/ws/*</div>
      <div class="d-box purple">WebSocket Target Group</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Default</div>
      <div class="d-box amber">Frontend Target Group (S3/ECS)</div>
    </div>
  </div>
  <div class="d-label">WAF rules evaluated before routing</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-alb-security-layers",
		Title:       "ALB Security Layers",
		Description: "Defense in depth showing request flow through AWS Shield, WAF rules, Cognito auth, ALB, and target group",
		ContentFile: "fundamentals/networking/load-balancing/alb",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Client</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red">AWS Shield (DDoS)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber">WAF Rules (SQL injection, XSS, rate limit)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Cognito Auth</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">ALB</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green">Target Group</div>
</div>`,
	})

	// -------------------------------------------------------
	// NLB
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-nlb-architecture",
		Title:       "NLB Architecture",
		Description: "NLB Layer 4 load balancing with static IP, distributing TCP/UDP traffic across targets in multiple AZs with source IP preserved",
		ContentFile: "fundamentals/networking/load-balancing/nlb",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Client (Static IP: 1.2.3.4)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">NLB (Layer 4 — TCP/UDP)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">AZ-1a</div>
        <div class="d-box green">Target 1</div>
        <div class="d-box green">Target 2</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">AZ-1b</div>
        <div class="d-box green">Target 3</div>
        <div class="d-box green">Target 4</div>
      </div>
    </div>
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">AZ-1c</div>
        <div class="d-box green">Target 5</div>
        <div class="d-box green">Target 6</div>
      </div>
    </div>
  </div>
  <div class="d-label">Source IP preserved</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-nlb-privatelink",
		Title:       "PrivateLink Architecture",
		Description: "NLB-based PrivateLink architecture connecting provider and consumer VPCs over AWS private network without traversing the public internet",
		ContentFile: "fundamentals/networking/load-balancing/nlb",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">PROVIDER VPC</div>
      <div class="d-flow-v">
        <div class="d-box green">Service</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo">NLB</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">VPC Endpoint Service</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-label">AWS PrivateLink (private network)</div>
      <div class="d-arrow">&#8594;</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">CONSUMER VPC</div>
      <div class="d-flow-v">
        <div class="d-box amber">VPC Endpoint</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Application</div>
      </div>
    </div>
  </div>
</div>`,
	})

	// -------------------------------------------------------
	// DynamoDB
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-dynamodb-key-structure",
		Title:       "DynamoDB Key Structure",
		Description: "Partition key distributes items across partitions, sort key orders items within each partition using the PK#value SK#value convention",
		ContentFile: "fundamentals/storage/dynamodb",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Table: ECommerce</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Partition Key (PK) distributes items</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">PARTITION A</div>
        <div class="d-flow-v">
          <div class="d-box green">PK = USER#U123</div>
          <div class="d-label">Sort Key orders items:</div>
          <div class="d-box gray">SK = ORDER#2024-01</div>
          <div class="d-box gray">SK = ORDER#2024-02</div>
          <div class="d-box gray">SK = PROFILE</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">PARTITION B</div>
        <div class="d-flow-v">
          <div class="d-box green">PK = USER#U456</div>
          <div class="d-label">Sort Key orders items:</div>
          <div class="d-box gray">SK = ORDER#2024-03</div>
          <div class="d-box gray">SK = PROFILE</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">PARTITION C</div>
        <div class="d-flow-v">
          <div class="d-box green">PK = PRODUCT#P789</div>
          <div class="d-label">Sort Key orders items:</div>
          <div class="d-box gray">SK = DETAILS</div>
          <div class="d-box gray">SK = REVIEW#U123</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-dynamodb-single-table",
		Title:       "Single-Table Design Layout",
		Description: "Single-table design showing multiple entity types (user profiles, orders, products, reviews) coexisting in one DynamoDB table with composite keys",
		ContentFile: "fundamentals/storage/dynamodb",
		Type:        TypeHTML,
		HTML: `<div class="d-group">
  <div class="d-group-title">TABLE: ECommerce (Single Table)</div>
  <div class="d-flow-v">
    <div class="d-group">
      <div class="d-group-title">PARTITION: USER#U123</div>
      <div class="d-row">
        <div class="d-box blue">SK = PROFILE<br>name, email</div>
        <div class="d-box green">SK = ORDER#2024-01-15#O456<br>total, status</div>
        <div class="d-box green">SK = ORDER#2024-02-20#O789<br>total, status</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">PARTITION: PRODUCT#P789</div>
      <div class="d-row">
        <div class="d-box purple">SK = DETAILS<br>name, price, category</div>
        <div class="d-box amber">SK = REVIEW#U123<br>rating, text</div>
        <div class="d-box amber">SK = REVIEW#U456<br>rating, text</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-dynamodb-write-sharding",
		Title:       "Write Sharding Pattern",
		Description: "Comparison of single partition key (throttled) vs write sharding across multiple partitions to solve the hot partition problem",
		ContentFile: "fundamentals/storage/dynamodb",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">WITHOUT SHARDING</div>
      <div class="d-flow-v">
        <div class="d-box blue">All Writes</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">POST#123 (single PK)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">Single Partition<br>THROTTLED!</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">WITH SHARDING</div>
      <div class="d-flow-v">
        <div class="d-box blue">All Writes</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Random Shard (0-9)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-row">
          <div class="d-box green">POST#123#0</div>
          <div class="d-box green">POST#123#1</div>
          <div class="d-box green">...</div>
          <div class="d-box green">POST#123#9</div>
        </div>
        <div class="d-label">Distributed across 10 partitions</div>
      </div>
    </div>
  </div>
</div>`,
	})

	// -------------------------------------------------------
	// Redis
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-redis-cache-aside",
		Title:       "Cache-Aside Pattern Flow",
		Description: "Cache-aside (lazy load) pattern showing application checking Redis first, then querying database on miss and populating cache",
		ContentFile: "fundamentals/storage/redis",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Application Request</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Check Redis Cache</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">HIT</div>
      <div class="d-box green">Return Cached Data</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">MISS</div>
      <div class="d-flow-v">
        <div class="d-box amber">Query Database</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Write Result to Redis</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Return Data</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-redis-cluster",
		Title:       "ElastiCache Cluster Architecture",
		Description: "Redis cluster mode enabled with 16,384 hash slots distributed across shards, each with primary and replica nodes",
		ContentFile: "fundamentals/storage/redis",
		Type:        TypeHTML,
		HTML: `<div class="d-label">Cluster Mode Enabled (16,384 hash slots)</div>
<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">SHARD 1 (slots 0-5460)</div>
      <div class="d-flow-v">
        <div class="d-box blue">Primary</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-row">
          <div class="d-box green">Replica 1</div>
          <div class="d-box green">Replica 2</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">SHARD 2 (slots 5461-10922)</div>
      <div class="d-flow-v">
        <div class="d-box blue">Primary</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-row">
          <div class="d-box green">Replica 1</div>
          <div class="d-box green">Replica 2</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">SHARD 3 (slots 10923-16383)</div>
      <div class="d-flow-v">
        <div class="d-box blue">Primary</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-row">
          <div class="d-box green">Replica 1</div>
          <div class="d-box green">Replica 2</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-redis-failover",
		Title:       "Redis Failover Flow",
		Description: "Redis primary node failure and automatic failover sequence through Sentinel/ElastiCache detection, replica promotion, DNS update, and client reconnection",
		ContentFile: "fundamentals/storage/redis",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box red">Primary Node Crashes</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber">Sentinel / ElastiCache Detects Failure (10-30s)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box blue">Promote Replica to Primary</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Update DNS Endpoint</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green">Clients Reconnect Automatically</div>
</div>`,
	})
}

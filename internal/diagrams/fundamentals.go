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
    <div class="d-box blue" data-tip="End users from global locations. CDN routes to nearest edge PoP via anycast DNS.">User (Sydney)</div>
    <div class="d-box blue" data-tip="End users from global locations. CDN routes to nearest edge PoP via anycast DNS.">User (London)</div>
  </div>
  <div class="d-arrow-down"><span class="d-step">1</span> ↓</div>
  <div class="d-row">
    <div class="d-box amber" data-tip="400+ global PoPs. Serves cached content with sub-20ms latency. Cache capacity ~1-10 TB per PoP.">Edge PoP - Sydney <span class="d-metric latency">5-20ms</span></div>
    <div class="d-box amber" data-tip="400+ global PoPs. Serves cached content with sub-20ms latency. Cache capacity ~1-10 TB per PoP.">Edge PoP - London <span class="d-metric latency">5-20ms</span></div>
  </div>
  <div class="d-label">HIT = Return cached response</div>
  <div class="d-arrow-down"><span class="d-step">2</span> ↓ MISS</div>
  <div class="d-box green" data-tip="Single regional cache layer between edge PoPs and origin. Collapses duplicate origin requests. Reduces origin load by up to 90%.">Regional Cache (Origin Shield) <span class="d-metric latency">30-60ms</span></div>
  <div class="d-arrow-down"><span class="d-step">3</span> ↓ MISS</div>
  <div class="d-box purple" data-tip="S3 for static assets, ALB for dynamic content. Full round-trip includes TLS negotiation + compute time.">Origin Server (S3 / ALB) <span class="d-metric latency">100-300ms</span></div>
  <div class="d-caption">Each cache layer reduces origin load exponentially: edge PoPs handle ~95% of requests, Origin Shield catches ~4%, only ~1% reaches origin.</div>
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
  <div class="d-box blue" data-tip="Global user base. CloudFront uses anycast to route each user to the lowest-latency edge PoP.">Users (global)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="400+ edge locations across 90+ cities. Supports HTTP/3, TLS 1.3, WebSocket. Up to 250K RPS per distribution.">CloudFront (400+ edge PoPs) <span class="d-metric throughput">250K RPS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">/static/*</div>
      <div class="d-box green" data-tip="Best for images, CSS, JS, fonts. S3 Transfer Acceleration enabled. Cost: $0.023/GB stored. Use OAC (Origin Access Control) for security.">S3 Origin <span class="d-metric latency">5-50ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">/api/*</div>
      <div class="d-box purple" data-tip="Dynamic API responses. Short TTL (0-60s) or no-cache. ALB health checks route to healthy targets. Supports sticky sessions.">ALB Origin <span class="d-metric latency">50-200ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">/media/*</div>
      <div class="d-box amber" data-tip="Optimized for video streaming. Supports HLS/DASH. Auto-scales to handle live events. Low-latency chunked transfer.">MediaStore <span class="d-metric latency">10-100ms</span></div>
    </div>
  </div>
  <div class="d-label">Path-based behavior routing</div>
  <div class="d-legend"><span class="d-box green" style="display:inline;padding:2px 6px;">Static</span> <span class="d-box purple" style="display:inline;padding:2px 6px;">Dynamic</span> <span class="d-box amber" style="display:inline;padding:2px 6px;">Streaming</span></div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-cloudfront-event-hooks",
		Title:       "CloudFront Event Hooks",
		Description: "Request lifecycle through CloudFront showing viewer request/response and origin request/response event hooks for CF Functions and Lambda@Edge",
		ContentFile: "fundamentals/networking/cdn/cloudfront",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-box blue" data-tip="Incoming viewer request. CF Functions run here for URL rewrites, header manipulation, A/B testing. Sub-ms execution, JS only."><span class="d-step">1</span> Viewer Request</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber" data-tip="CloudFront Functions: sub-ms latency, 10KB code limit, JS only. 10M RPS max. Use for lightweight transforms: URL rewrites, header normalization, redirects.">CF Function <span class="d-metric latency">&lt;1ms</span></div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box gray" data-tip="Edge cache lookup. HIT skips origin entirely. TTL controlled by Cache-Control headers or cache policy."><span class="d-step">2</span> Edge Cache</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue" data-tip="Only fires on cache MISS. Lambda@Edge can modify origin selection, add auth headers, normalize cache keys."><span class="d-step">3</span> Origin Request</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple" data-tip="Lambda@Edge: 5s timeout (viewer) / 30s (origin). Node.js or Python. 1MB code limit. Use for auth, dynamic origin selection, image resize.">Lambda@Edge <span class="d-metric latency">5-50ms</span></div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box green" data-tip="Origin fetches response. S3, ALB, API Gateway, or custom origin."><span class="d-step">4</span> Origin</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue" data-tip="Modify origin response before caching. Add security headers, transform content, set custom cache TTLs."><span class="d-step">5</span> Origin Response</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple" data-tip="Lambda@Edge can modify response headers, add CORS, transform body before caching at edge.">Lambda@Edge</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box gray"><span class="d-step">6</span> Edge Cache</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue" data-tip="Final response to viewer. CF Function can add security headers, modify response for A/B testing."><span class="d-step">7</span> Viewer Response</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber" data-tip="Lightweight response transforms: add security headers (CSP, HSTS), customize error pages.">CF Function</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box indigo"><span class="d-step">8</span> Client</div>
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
    <div class="d-box blue" data-tip="Each edge PoP independently caches content. Without Origin Shield, N PoPs = N cache misses to origin for the same object.">Edge PoP (Sydney)</div>
    <div class="d-box blue" data-tip="Each edge PoP independently caches content. Without Origin Shield, N PoPs = N cache misses to origin for the same object.">Edge PoP (London)</div>
    <div class="d-box blue" data-tip="Each edge PoP independently caches content. Without Origin Shield, N PoPs = N cache misses to origin for the same object.">Edge PoP (Tokyo)</div>
  </div>
  <div class="d-arrow-down"><span class="d-step">1</span> &#8595; MISS</div>
  <div class="d-box purple" data-tip="Choose the region closest to your origin. Acts as a single funnel: collapses multiple edge misses into one origin fetch. Incremental cost: ~$0.0075/10K requests.">Origin Shield (us-east-1) <span class="d-metric throughput">90% origin reduction</span></div>
  <div class="d-arrow-down"><span class="d-step">2</span> &#8595; MISS (only ~10%)</div>
  <div class="d-box green" data-tip="Origin receives only ~10% of total requests. Best paired with S3 for static or ALB with auto-scaling for dynamic.">Origin (S3/ALB) <span class="d-metric throughput">~10% of requests</span></div>
  <div class="d-caption">Without Origin Shield, 400 PoPs each miss independently = 400 origin fetches. With Origin Shield, those collapse to 1 fetch.</div>
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
  <div class="d-box blue" data-tip="Public internet traffic. Multiple ISP paths for redundancy."><span class="d-step">1</span> Internet</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-box gray" data-tip="Latency-based or weighted routing. Health-check failover with 10s interval. TTL 60s. Supports alias records for ALB/NLB."><span class="d-step">2</span> Route 53 (DNS) <span class="d-metric latency">1-50ms</span></div>
  <div class="d-arrow-down">↓</div>
  <div class="d-box green" data-tip="Cross-zone load balancing enabled by default on ALB. ALB: Layer 7 (HTTP), NLB: Layer 4 (TCP/UDP). Auto-scales to millions of RPS."><span class="d-step">3</span> ALB / NLB (spans all AZs) <span class="d-metric throughput">millions RPS</span></div>
  <div class="d-arrow-down">↓</div>
  <div class="d-row">
    <div class="d-group">
      <div class="d-group-title">AZ-1a</div>
      <div class="d-box indigo" data-tip="Min 2 targets per AZ for high availability. Auto-scaling group maintains desired count. Deregistration delay: 300s default.">Target Group EC2 / ECS</div>
    </div>
    <div class="d-group">
      <div class="d-group-title">AZ-1b</div>
      <div class="d-box indigo" data-tip="Cross-zone balancing distributes evenly regardless of target count per AZ.">Target Group EC2 / ECS</div>
    </div>
    <div class="d-group">
      <div class="d-group-title">AZ-1c</div>
      <div class="d-box indigo" data-tip="If entire AZ fails, Route 53 zonal shift via ARC removes it in ~60s. Remaining AZs absorb traffic.">Target Group EC2 / ECS</div>
    </div>
  </div>
  <div class="d-label">Health checks: ALB pings /health every 30s. Unhealthy target removed in ~60s. Entire AZ down: zonal shift via ARC.</div>
  <div class="d-caption">Minimum 3 AZs for production: tolerates 1 AZ failure while maintaining 66% capacity. Over-provision each AZ by 50%.</div>
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
  <div class="d-box blue" data-tip="HTTPS traffic on port 443. ALB terminates TLS — offloads CPU-intensive encryption from targets.">Internet</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="Listener supports up to 100 rules. Rules evaluated in priority order (1-50000). Supports host-based, path-based, HTTP header, query string, and source IP conditions.">ALB (Listener :443) <span class="d-metric throughput">100 rules max</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">/api/*</div>
      <div class="d-box green" data-tip="ECS Fargate tasks. Supports gRPC and HTTP/2. Slow start mode: ramps new targets from 0% to full over 30-900s. Health check: HTTP 200 on /health.">API Target Group (ECS) <span class="d-metric latency">50-200ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">/ws/*</div>
      <div class="d-box purple" data-tip="Sticky sessions via app cookie or ALB-generated cookie. WebSocket idle timeout: 4000s max. Connection draining on deregistration.">WebSocket Target Group <span class="d-metric latency">1-5ms per frame</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Default</div>
      <div class="d-box amber" data-tip="Static assets from S3 or server-rendered pages from ECS. Can return fixed response (maintenance page) or redirect (HTTP 301/302).">Frontend Target Group (S3/ECS)</div>
    </div>
  </div>
  <div class="d-label">WAF rules evaluated before routing</div>
  <div class="d-caption">Path-based routing eliminates the need for separate load balancers per service. One ALB handles all microservices with up to 100 routing rules.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-alb-security-layers",
		Title:       "ALB Security Layers",
		Description: "Defense in depth showing request flow through AWS Shield, WAF rules, Cognito auth, ALB, and target group",
		ContentFile: "fundamentals/networking/load-balancing/alb",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue"><span class="d-step">1</span> Client</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red" data-tip="Shield Standard: free, protects against L3/L4 DDoS (SYN flood, UDP reflection). Shield Advanced: $3K/mo, L7 protection, DRT team response, cost protection on scaling."><span class="d-step">2</span> AWS Shield (DDoS) <span class="d-metric throughput">Tbps mitigation</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="Up to 1500 rules per Web ACL. Rate-based rules: block IPs exceeding threshold (e.g., 2000 req/5min). Managed rule groups for OWASP Top 10. $5/mo per Web ACL + $1/mo per rule."><span class="d-step">3</span> WAF Rules (SQL injection, XSS, rate limit) <span class="d-metric throughput">1500 rules max</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="ALB-native Cognito integration: no auth code needed in app. Supports OIDC, SAML, social IdPs. JWT tokens validated at ALB layer. Unauthenticated users redirected to login."><span class="d-step">4</span> Cognito Auth <span class="d-metric latency">10-50ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="TLS termination, mTLS support. Access logs to S3 for audit. Security groups restrict inbound to 443/80 only."><span class="d-step">5</span> ALB</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="Security groups allow traffic only from ALB security group. No direct internet access. Private subnets for targets."><span class="d-step">6</span> Target Group</div>
  <div class="d-caption">Defense in depth: each layer independently blocks threats. A request surviving all 6 layers is authenticated, rate-limited, and sanitized.</div>
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
  <div class="d-box blue" data-tip="Client connects to NLB static IP or Elastic IP. Static IPs allow firewall whitelisting — not possible with ALB.">Client (Static IP: 1.2.3.4)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="Layer 4: TCP/UDP/TLS. Sub-ms latency (hardware flow-based). Millions of RPS. Static IP per AZ. Preserves source IP. No security groups (use NACLs). TLS passthrough or termination.">NLB (Layer 4 — TCP/UDP) <span class="d-metric latency">&lt;1ms</span> <span class="d-metric throughput">millions RPS</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">AZ-1a</div>
        <div class="d-box green" data-tip="Hash-based routing: 5-tuple (src IP, dst IP, src port, dst port, protocol). Same client always hits same target.">Target 1</div>
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
  <div class="d-legend"><span class="d-box indigo" style="display:inline;padding:2px 6px;">NLB</span> Layer 4 only — no HTTP routing, no WAF, no sticky sessions. Use ALB if you need L7 features.</div>
  <div class="d-caption">Choose NLB for: static IPs, extreme low latency, non-HTTP protocols (gRPC, MQTT, gaming), or PrivateLink endpoints.</div>
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
        <div class="d-box green" data-tip="Provider's internal service. Can be in any VPC, any account. Exposed only through NLB + Endpoint Service — no public IP needed.">Service</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo" data-tip="NLB required for PrivateLink. Routes traffic to service targets. One NLB can back multiple endpoint services.">NLB <span class="d-metric latency">&lt;1ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Accept/reject connection requests per consumer account. Supports cross-account and cross-region. Cost: $0.01/hr + $0.01/GB processed.">VPC Endpoint Service</div>
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
        <div class="d-box amber" data-tip="Creates an ENI in the consumer VPC. Gets a private IP. DNS resolves the endpoint to this private IP. No VPC peering or Transit Gateway needed.">VPC Endpoint <span class="d-metric size">ENI per AZ</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue" data-tip="Connects to endpoint using private DNS. Traffic never leaves AWS backbone network. Appears as a local resource.">Application</div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">PrivateLink eliminates the need for VPC peering, NAT gateways, or internet gateways. Traffic stays on AWS backbone — no public internet exposure.</div>`,
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
  <div class="d-box blue" data-tip="One DynamoDB table can handle 10M+ items. Max item size: 400KB. On-demand mode: pay per request. Provisioned mode: specify RCU/WCU.">Table: ECommerce</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="PK hashed to determine physical partition. High-cardinality keys distribute evenly. Each partition: 3000 RCU, 1000 WCU, 10GB max. Avoid low-cardinality PKs (e.g., status, date).">Partition Key (PK) distributes items <span class="d-metric throughput">1000 WCU/partition</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">PARTITION A</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="Composite key pattern: ENTITY#ID. Enables Query on PK to fetch all items for an entity. GetItem requires both PK + SK.">PK = USER#U123</div>
          <div class="d-label">Sort Key orders items:</div>
          <div class="d-box gray" data-tip="SK enables range queries: begins_with, between, >, <. Items sorted by SK in B-tree within partition. Query returns items in SK order.">SK = ORDER#2024-01 <span class="d-metric latency">&lt;10ms query</span></div>
          <div class="d-box gray">SK = ORDER#2024-02</div>
          <div class="d-box gray" data-tip="Singleton item pattern: use a fixed SK like PROFILE for 1:1 entity data alongside 1:N collections.">SK = PROFILE</div>
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
          <div class="d-box green" data-tip="Different entity types can coexist in the same table. PK prefix (USER#, PRODUCT#) acts as a type discriminator.">PK = PRODUCT#P789</div>
          <div class="d-label">Sort Key orders items:</div>
          <div class="d-box gray">SK = DETAILS</div>
          <div class="d-box gray">SK = REVIEW#U123</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-caption">Key design rule: model access patterns first, then design keys. Every Query must specify PK. SK enables flexible range queries within a partition.</div>
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
        <div class="d-box blue" data-tip="1:1 entity. GetItem(PK=USER#U123, SK=PROFILE) — single-digit ms. Max 400KB per item. Store denormalized data to avoid joins.">SK = PROFILE<br>name, email <span class="d-metric latency">&lt;5ms</span></div>
        <div class="d-box green" data-tip="1:N collection. Query(PK=USER#U123, SK begins_with ORDER#) returns all orders sorted by date. Paginate with LastEvaluatedKey.">SK = ORDER#2024-01-15#O456<br>total, status</div>
        <div class="d-box green" data-tip="Date-prefixed SK enables time-range queries: SK between ORDER#2024-01 and ORDER#2024-12 for yearly orders.">SK = ORDER#2024-02-20#O789<br>total, status</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">PARTITION: PRODUCT#P789</div>
      <div class="d-row">
        <div class="d-box purple" data-tip="Product details as singleton. GSI on category+price enables browse-by-category queries without scanning.">SK = DETAILS<br>name, price, category</div>
        <div class="d-box amber" data-tip="Reviews stored under product PK. Inverted GSI (PK=USER#U123, SK=REVIEW#P789) enables 'my reviews' query. This is the power of single-table design.">SK = REVIEW#U123<br>rating, text</div>
        <div class="d-box amber" data-tip="Multiple reviews per product. Query with SK begins_with REVIEW# returns all reviews. Add GSI for sorting by rating.">SK = REVIEW#U456<br>rating, text</div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Single-table design trades write simplicity for read efficiency: one Query fetches related entities that would require JOINs in SQL. Use GSIs for alternate access patterns.</div>`,
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
        <div class="d-box blue">All Writes <span class="d-metric throughput">10K WPS</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Single PK funnels all writes to one partition. DynamoDB partition limit: 1000 WCU. Above that = ProvisionedThroughputExceededException.">POST#123 (single PK) <span class="d-metric throughput">1000 WCU max</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="Hot partition problem. Even with on-demand mode, a single partition tops out at 1000 WCU. Adaptive capacity helps but does not fully solve it.">Single Partition<br>THROTTLED! <span class="d-status error"></span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">WITH SHARDING</div>
      <div class="d-flow-v">
        <div class="d-box blue">All Writes <span class="d-metric throughput">10K WPS</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Application appends random suffix (0-9) to PK. Distributes writes across 10 partitions. Read requires Scatter-Gather across all shards.">Random Shard (0-9) <span class="d-metric throughput">10x capacity</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-row">
          <div class="d-box green" data-tip="Each shard handles 1000 WCU independently. 10 shards = 10,000 WCU effective capacity.">POST#123#0 <span class="d-metric throughput">1K WCU</span></div>
          <div class="d-box green">POST#123#1 <span class="d-metric throughput">1K WCU</span></div>
          <div class="d-box green">...</div>
          <div class="d-box green">POST#123#9 <span class="d-metric throughput">1K WCU</span></div>
        </div>
        <div class="d-label">Distributed across 10 partitions <span class="d-metric throughput">10K WCU total</span></div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Trade-off: write sharding multiplies write capacity linearly but complicates reads — reads must scatter-gather across all shards and aggregate client-side.</div>`,
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
  <div class="d-box blue" data-tip="Application always checks cache first. If key exists and TTL not expired, return immediately. Cache-aside is the most common caching pattern."><span class="d-step">1</span> Application Request</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="GET key — O(1) operation. Sub-ms over ElastiCache in same AZ. Network hop adds ~0.1-0.5ms. Pipeline multiple GETs for batch lookups."><span class="d-step">2</span> Check Redis Cache <span class="d-metric latency">&lt;1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">HIT (~95%)</div>
      <div class="d-box green" data-tip="Cache hit: return immediately. No database query. Typical hit rate: 95-99% for read-heavy workloads with proper TTL tuning."><span class="d-step">3a</span> Return Cached Data <span class="d-metric latency">&lt;1ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">MISS (~5%)</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Cache miss triggers DB query. Protect against thundering herd with request coalescing (singleflight) or mutex lock on cache key."><span class="d-step">3b</span> Query Database <span class="d-metric latency">5-50ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="SET key value EX ttl — write result to cache with TTL. TTL prevents stale data. Typical TTL: 5min-24h depending on data freshness requirements."><span class="d-step">4</span> Write Result to Redis <span class="d-metric latency">&lt;1ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><span class="d-step">5</span> Return Data</div>
      </div>
    </div>
  </div>
  <div class="d-caption">Cache-aside is lazy: only caches data on first read. Cold start requires warming. Alternative: write-through caches on every write but wastes memory on unread data.</div>
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
        <div class="d-box blue" data-tip="Primary handles all writes for its slot range. Max memory: r6g.xlarge = 26GB. Writes replicated async to replicas. If primary fails, replica promoted in 10-30s.">Primary <span class="d-metric size">26 GB</span> <span class="d-metric throughput">100K ops/s</span></div>
        <div class="d-arrow-down">&#8595; async repl</div>
        <div class="d-row">
          <div class="d-box green" data-tip="Replica serves reads to scale read throughput. Replica lag typically &lt;1ms. Reader endpoint auto-distributes across replicas.">Replica 1 <span class="d-metric latency">&lt;1ms lag</span></div>
          <div class="d-box green" data-tip="Second replica for HA — survives primary + 1 replica failure. Cross-AZ replicas for AZ-level resilience.">Replica 2</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">SHARD 2 (slots 5461-10922)</div>
      <div class="d-flow-v">
        <div class="d-box blue">Primary <span class="d-metric size">26 GB</span> <span class="d-metric throughput">100K ops/s</span></div>
        <div class="d-arrow-down">&#8595; async repl</div>
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
        <div class="d-box blue">Primary <span class="d-metric size">26 GB</span> <span class="d-metric throughput">100K ops/s</span></div>
        <div class="d-arrow-down">&#8595; async repl</div>
        <div class="d-row">
          <div class="d-box green">Replica 1</div>
          <div class="d-box green">Replica 2</div>
        </div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">3 shards x 26GB = 78GB total. 3 shards x 100K ops/s = 300K ops/s. Scale horizontally by adding shards — online resharding redistributes hash slots with zero downtime.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-redis-failover",
		Title:       "Redis Failover Flow",
		Description: "Redis primary node failure and automatic failover sequence through Sentinel/ElastiCache detection, replica promotion, DNS update, and client reconnection",
		ContentFile: "fundamentals/storage/redis",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box red" data-tip="Primary crash: OOM kill, hardware failure, or network partition. All in-flight writes lost unless wait/WAIT command used for synchronous replication."><span class="d-step">1</span> Primary Node Crashes <span class="d-status error"></span> <span class="d-metric latency">T+0s</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="ElastiCache monitors via heartbeat every 1s. After N missed heartbeats (default 5), node marked as failing. Sentinel quorum (majority) must agree before failover."><span class="d-step">2</span> Sentinel / ElastiCache Detects Failure <span class="d-status error"></span> <span class="d-metric latency">10-30s</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box blue" data-tip="Replica with least replication lag selected. SLAVEOF NO ONE issued. Replica loads its dataset and begins accepting writes. Other replicas re-pointed to new primary."><span class="d-step">3</span> Promote Replica to Primary <span class="d-status active"></span> <span class="d-metric latency">5-10s</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="ElastiCache primary endpoint DNS updated to point to new primary IP. DNS TTL is 5s. Clients using primary endpoint auto-resolve to new primary."><span class="d-step">4</span> Update DNS Endpoint <span class="d-status active"></span> <span class="d-metric latency">5-15s</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="Clients reconnect using primary endpoint. Total downtime: 20-60s. Use retry with exponential backoff. Application should handle ReadOnlyError during transition."><span class="d-step">5</span> Clients Reconnect Automatically <span class="d-status active"></span> <span class="d-metric latency">Total: 20-60s</span></div>
  <div class="d-caption">Total failover window: 20-60 seconds. During this window, writes fail and reads may return stale data. Design applications to handle transient Redis unavailability gracefully.</div>
</div>`,
	})
}

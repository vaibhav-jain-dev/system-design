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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number" data-tip="95% edge hit rate = only 5% of requests reach origin. For 100K RPS: 95K served at &lt;20ms from edge, 5K reach origin. Without CDN all 100K would hit origin. This is why CDN is non-negotiable for media-heavy apps."><div class="d-number-value">95%</div><div class="d-number-label">Requests served from edge (5% reach origin)</div></div>
  <div class="d-number"><div class="d-number-value">400+</div><div class="d-number-label">Global PoPs</div></div>
</div>
<div class="d-flow-v">
  <div class="d-row">
    <div class="d-box blue" data-tip="End users from global locations. CDN routes to nearest edge PoP via anycast DNS.">User (Sydney)</div>
    <div class="d-box blue" data-tip="End users from global locations. CDN routes to nearest edge PoP via anycast DNS.">User (London)</div>
  </div>
  <div class="d-arrow-down"><span class="d-step">1</span> ↓</div>
  <div class="d-row">
    <div class="d-box amber" data-tip="400+ global PoPs. Cache capacity: ~1–10 TB per PoP SSD (CloudFront) or 50+ TB (Akamai Tier 1). 95% edge hit rate means: for 58K RPS, 55,100 RPS served from edge at &lt;20ms, only 2,900 RPS reach origin. Each edge hit saves ~100–300ms vs origin round trip.">Edge PoP - Sydney <span class="d-metric latency">5-20ms</span> <div class="d-tag green">~95% hit rate</div></div>
    <div class="d-box amber" data-tip="400+ global PoPs. Serves cached content with sub-20ms latency. Cache capacity ~1-10 TB per PoP.">Edge PoP - London <span class="d-metric latency">5-20ms</span></div>
  </div>
  <div class="d-label">HIT = Return cached response <span class="d-status active"></span></div>
  <div class="d-arrow-down"><span class="d-step">2</span> ↓ MISS <span class="d-status error"></span></div>
  <div class="d-box green" data-tip="Single regional cache layer between edge PoPs and origin. Collapses duplicate origin requests. Reduces origin load by up to 90%.">Regional Cache (Origin Shield) <span class="d-metric latency">30-60ms</span> <div class="d-tag green">90% origin reduction</div></div>
  <div class="d-arrow-down"><span class="d-step">3</span> ↓ MISS (only ~1%)</div>
  <div class="d-box purple" data-tip="S3 for static assets, ALB for dynamic content. Full round-trip includes TLS negotiation + compute time.">Origin Server (S3 / ALB) <span class="d-metric latency">100-300ms</span></div>
  <div class="d-caption">Cache hit math example: 58K RPS total. Edge PoPs (95% hit): 55,100 served at edge, 2,900 miss. Origin Shield (80% of misses): 2,320 served from regional cache, 580 miss. Origin sees only 580 RPS (1% of total). Cost savings: each edge hit avoids origin compute + DB read (≈$0.00001/req) = $0.58/sec = $50K/day saved on origin infrastructure.</div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">250K</div><div class="d-number-label">RPS per distribution</div></div>
  <div class="d-number"><div class="d-number-value">400+</div><div class="d-number-label">Edge locations</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box blue" data-tip="Global user base. CloudFront uses anycast to route each user to the lowest-latency edge PoP.">Users (global)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="400+ edge locations across 90+ cities. Supports HTTP/3, TLS 1.3, WebSocket. Default limit: 250K RPS per distribution (can be raised to unlimited via AWS support). Each PoP is an Anycast IP — DNS automatically routes users to nearest PoP (usually within 50ms). A request from Sydney goes to Sydney PoP, not US-EAST-1.">CloudFront (400+ edge PoPs) <span class="d-metric throughput">250K RPS</span> <div class="d-tag indigo">HTTP/3 + TLS 1.3</div></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">/static/*</div>
      <div class="d-box green" data-tip="Best for images, CSS, JS, fonts. S3 Transfer Acceleration enabled. Cost: $0.023/GB stored. Use OAC (Origin Access Control) for security.">S3 Origin <span class="d-metric latency">5-50ms</span> <div class="d-tag green">long TTL</div></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">/api/*</div>
      <div class="d-box purple" data-tip="Dynamic API responses. Short TTL (0-60s) or no-cache. ALB health checks route to healthy targets. Supports sticky sessions.">ALB Origin <span class="d-metric latency">50-200ms</span> <div class="d-tag amber">short TTL</div></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">/media/*</div>
      <div class="d-box amber" data-tip="Optimized for video streaming. Supports HLS/DASH. Auto-scales to handle live events. Low-latency chunked transfer.">MediaStore <span class="d-metric latency">10-100ms</span> <div class="d-tag indigo">HLS/DASH</div></div>
    </div>
  </div>
  <div class="d-label">Path-based behavior routing</div>
  <div class="d-legend"><span class="d-box green" style="display:inline;padding:2px 6px;">Static</span> long TTL &nbsp; <span class="d-box purple" style="display:inline;padding:2px 6px;">Dynamic</span> short TTL &nbsp; <span class="d-box amber" style="display:inline;padding:2px 6px;">Streaming</span> adaptive bitrate</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-cloudfront-event-hooks",
		Title:       "CloudFront Event Hooks",
		Description: "Request lifecycle through CloudFront showing viewer request/response and origin request/response event hooks for CF Functions and Lambda@Edge",
		ContentFile: "fundamentals/networking/cdn/cloudfront",
		Type:        TypeHTML,
		HTML: `<div class="d-flow" style="flex-wrap: wrap; gap: 0.25rem;">
  <div class="d-box blue" data-tip="Incoming viewer request. CF Functions run here for URL rewrites, header manipulation, A/B testing. Sub-ms execution, JS only."><span class="d-step">1</span> Viewer Request</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber" data-tip="CloudFront Functions: sub-ms latency, 10KB code limit, JS only. 10M RPS max. Use for lightweight transforms: URL rewrites, header normalization, redirects.">CF Function <span class="d-metric latency">&lt;1ms</span> <div class="d-tag green">10M RPS</div></div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box gray" data-tip="Edge cache lookup. HIT skips origin entirely. TTL controlled by Cache-Control headers or cache policy."><span class="d-step">2</span> Edge Cache</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue" data-tip="Only fires on cache MISS. Lambda@Edge can modify origin selection, add auth headers, normalize cache keys."><span class="d-step">3</span> Origin Request</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple" data-tip="Lambda@Edge: 5s timeout (viewer) / 30s (origin). Node.js or Python. 1MB code limit. Use for auth, dynamic origin selection, image resize.">Lambda@Edge <span class="d-metric latency">5-50ms</span> <div class="d-tag indigo">30s timeout</div></div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box green" data-tip="Origin fetches response. S3, ALB, API Gateway, or custom origin."><span class="d-step">4</span> Origin <span class="d-status active"></span></div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue" data-tip="Modify origin response before caching. Add security headers, transform content, set custom cache TTLs."><span class="d-step">5</span> Origin Response</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple" data-tip="Lambda@Edge can modify response headers, add CORS, transform body before caching at edge.">Lambda@Edge</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box gray" data-tip="Response stored in edge cache for future requests with matching cache key."><span class="d-step">6</span> Edge Cache</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue" data-tip="Final response to viewer. CF Function can add security headers, modify response for A/B testing."><span class="d-step">7</span> Viewer Response</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber" data-tip="Lightweight response transforms: add security headers (CSP, HSTS), customize error pages.">CF Function</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box indigo" data-tip="Final validated response delivered to client. Total latency for a cache hit: viewer request CF Function + cache lookup + viewer response CF Function = ~2ms."><span class="d-step">8</span> Client <span class="d-status active"></span></div>
</div>
<div class="d-legend">
  <span class="d-box amber" style="display:inline-block;padding:2px 8px;font-size:0.75rem;">CF Function</span> &lt;1ms, JS only &nbsp;
  <span class="d-box purple" style="display:inline-block;padding:2px 8px;font-size:0.75rem;">Lambda@Edge</span> 5-30s, Node/Python &nbsp;
  <span class="d-box gray" style="display:inline-block;padding:2px 8px;font-size:0.75rem;">Cache</span> TTL-based
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-cloudfront-origin-shield",
		Title:       "Origin Shield Architecture",
		Description: "Origin Shield as an additional caching layer between edge PoPs and origin, reducing origin load by 90%",
		ContentFile: "fundamentals/networking/cdn/cloudfront",
		Type:        TypeHTML,
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">400→1</div><div class="d-number-label">PoP misses collapse to 1 origin fetch</div></div>
  <div class="d-number"><div class="d-number-value">90%</div><div class="d-number-label">Origin load reduction</div></div>
</div>
<div class="d-flow-v">
  <div class="d-row">
    <div class="d-box blue" data-tip="Each edge PoP independently caches content. Without Origin Shield, N PoPs = N cache misses to origin for the same object.">Edge PoP (Sydney)</div>
    <div class="d-box blue" data-tip="Each edge PoP independently caches content. Without Origin Shield, N PoPs = N cache misses to origin for the same object.">Edge PoP (London)</div>
    <div class="d-box blue" data-tip="Each edge PoP independently caches content. Without Origin Shield, N PoPs = N cache misses to origin for the same object.">Edge PoP (Tokyo)</div>
  </div>
  <div class="d-arrow-down"><span class="d-step">1</span> &#8595; MISS (all 3 PoPs)</div>
  <div class="d-box purple" data-tip="Choose the region closest to your origin. Acts as a single funnel: collapses multiple edge misses into one origin fetch. Incremental cost: ~$0.0075/10K requests.">Origin Shield (us-east-1) <span class="d-metric throughput">90% origin reduction</span> <div class="d-tag green">recommended</div></div>
  <div class="d-arrow-down"><span class="d-step">2</span> &#8595; MISS (only ~10%)</div>
  <div class="d-box green" data-tip="Origin receives only ~10% of total requests. Best paired with S3 for static or ALB with auto-scaling for dynamic.">Origin (S3/ALB) <span class="d-metric throughput">~10% of requests</span></div>
  <div class="d-caption">Without Origin Shield: 400 PoPs, each with independent cache. If an object expires, up to 400 concurrent misses hit origin simultaneously (thundering herd). With Origin Shield: all 400 PoP misses are held and deduplicated — only 1 origin fetch happens. Cost: ~$0.0075/10K requests to Origin Shield, but saves 399 origin fetches per cache miss.</div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">3</div><div class="d-number-label">Minimum AZs for production</div></div>
  <div class="d-number"><div class="d-number-value">66%</div><div class="d-number-label">Capacity after 1 AZ failure</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box blue" data-tip="Public internet traffic. Multiple ISP paths for redundancy."><span class="d-step">1</span> Internet</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-box gray" data-tip="Latency-based or weighted routing. Health-check failover with 10s interval. TTL 60s. Supports alias records for ALB/NLB."><span class="d-step">2</span> Route 53 (DNS) <span class="d-metric latency">1-50ms</span> <div class="d-tag indigo">anycast</div></div>
  <div class="d-arrow-down">↓</div>
  <div class="d-box green" data-tip="Cross-zone load balancing enabled by default on ALB. ALB: Layer 7 (HTTP), NLB: Layer 4 (TCP/UDP). Auto-scales to millions of RPS."><span class="d-step">3</span> ALB / NLB (spans all AZs) <span class="d-metric throughput">millions RPS</span></div>
  <div class="d-arrow-down">↓</div>
  <div class="d-row">
    <div class="d-group">
      <div class="d-group-title">AZ-1a <span class="d-status active"></span></div>
      <div class="d-box indigo" data-tip="Min 2 targets per AZ for high availability. Auto-scaling group maintains desired count. Deregistration delay: 300s default.">Target Group EC2 / ECS</div>
    </div>
    <div class="d-group">
      <div class="d-group-title">AZ-1b <span class="d-status active"></span></div>
      <div class="d-box indigo" data-tip="Cross-zone balancing distributes evenly regardless of target count per AZ.">Target Group EC2 / ECS</div>
    </div>
    <div class="d-group">
      <div class="d-group-title">AZ-1c <span class="d-status active"></span></div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">100</div><div class="d-number-label">Max routing rules per listener</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box blue" data-tip="HTTPS traffic on port 443. ALB terminates TLS — offloads CPU-intensive encryption from targets.">Internet</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="Listener supports up to 100 rules. Rules evaluated in priority order (1-50000). Supports host-based, path-based, HTTP header, query string, and source IP conditions.">ALB (Listener :443) <span class="d-metric throughput">100 rules max</span> <div class="d-tag indigo">TLS termination</div></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">/api/* <span class="d-step">1</span></div>
      <div class="d-box green" data-tip="ECS Fargate tasks. Supports gRPC and HTTP/2. Slow start mode: ramps new targets from 0% to full over 30-900s. Health check: HTTP 200 on /health.">API Target Group (ECS) <span class="d-metric latency">50-200ms</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">/ws/* <span class="d-step">2</span></div>
      <div class="d-box purple" data-tip="Sticky sessions via app cookie or ALB-generated cookie. WebSocket idle timeout: 4000s max. Connection draining on deregistration.">WebSocket Target Group <span class="d-metric latency">1-5ms per frame</span> <div class="d-tag amber">sticky sessions</div></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">Default <span class="d-step">3</span></div>
      <div class="d-box amber" data-tip="Static assets from S3 or server-rendered pages from ECS. Can return fixed response (maintenance page) or redirect (HTTP 301/302).">Frontend Target Group (S3/ECS)</div>
    </div>
  </div>
  <div class="d-label">WAF rules evaluated before routing <div class="d-tag green">cheapest-first</div></div>
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
  <div class="d-box blue" data-tip="Raw internet request — untrusted. Must pass all 6 layers before reaching application logic."><span class="d-step">1</span> Client</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red" data-tip="Shield Standard: free, protects against L3/L4 DDoS (SYN flood, UDP reflection). Shield Advanced: $3K/mo, L7 protection, DRT team response, cost protection on scaling."><span class="d-step">2</span> AWS Shield (DDoS) <span class="d-metric throughput">Tbps mitigation</span> <div class="d-tag green">free tier available</div></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="Up to 1500 rules per Web ACL. Rate-based rules: block IPs exceeding threshold (e.g., 2000 req/5min). Managed rule groups for OWASP Top 10. $5/mo per Web ACL + $1/mo per rule."><span class="d-step">3</span> WAF Rules (SQL injection, XSS, rate limit) <span class="d-metric throughput">1500 rules max</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="ALB-native Cognito integration: no auth code needed in app. Supports OIDC, SAML, social IdPs. JWT tokens validated at ALB layer. Unauthenticated users redirected to login."><span class="d-step">4</span> Cognito Auth <span class="d-metric latency">10-50ms</span> <div class="d-tag indigo">OIDC/SAML</div></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="TLS termination, mTLS support. Access logs to S3 for audit. Security groups restrict inbound to 443/80 only."><span class="d-step">5</span> ALB</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="Security groups allow traffic only from ALB security group. No direct internet access. Private subnets for targets."><span class="d-step">6</span> Target Group <span class="d-status active"></span></div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">&lt;1ms</div><div class="d-number-label">NLB latency (Layer 4)</div></div>
  <div class="d-number"><div class="d-number-value">static</div><div class="d-number-label">IP per AZ (firewall-friendly)</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box blue" data-tip="Client connects to NLB static IP or Elastic IP. Static IPs allow firewall whitelisting — not possible with ALB.">Client (Static IP: 1.2.3.4)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="Layer 4: TCP/UDP/TLS. Sub-ms latency (hardware flow-based). Millions of RPS. Static IP per AZ. Preserves source IP. No security groups (use NACLs). TLS passthrough or termination.">NLB (Layer 4 — TCP/UDP) <span class="d-metric latency">&lt;1ms</span> <span class="d-metric throughput">millions RPS</span> <div class="d-tag indigo">source IP preserved</div></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-group">
        <div class="d-group-title">AZ-1a</div>
        <div class="d-box green" data-tip="Hash-based routing: 5-tuple (src IP, dst IP, src port, dst port, protocol). Same client always hits same target.">Target 1 <span class="d-status active"></span></div>
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
  <div class="d-label">Source IP preserved <div class="d-tag green">no X-Forwarded-For needed</div></div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">0</div><div class="d-number-label">Public internet exposure</div></div>
</div>
<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">PROVIDER VPC</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Provider's internal service. Can be in any VPC, any account. Exposed only through NLB + Endpoint Service — no public IP needed.">Service <span class="d-status active"></span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo" data-tip="NLB required for PrivateLink. Routes traffic to service targets. One NLB can back multiple endpoint services.">NLB <span class="d-metric latency">&lt;1ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Accept/reject connection requests per consumer account. Supports cross-account and cross-region. Cost: $0.01/hr + $0.01/GB processed.">VPC Endpoint Service <div class="d-tag indigo">cross-account</div></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-label">AWS PrivateLink (private network) <div class="d-tag green">never public internet</div></div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">400 KB</div><div class="d-number-label">Max item size</div></div>
  <div class="d-number"><div class="d-number-value">1,000</div><div class="d-number-label">WCU per partition max</div></div>
  <div class="d-number"><div class="d-number-value">&lt;10ms</div><div class="d-number-label">Query latency (any scale)</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box blue" data-tip="One DynamoDB table can handle 10M+ items. Max item size: 400KB. On-demand mode: pay per request. Provisioned mode: specify RCU/WCU.">Table: ECommerce</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="PK hashed to determine physical partition. High-cardinality keys distribute evenly. Each partition: 3000 RCU, 1000 WCU, 10GB max. Avoid low-cardinality PKs (e.g., status, date).">Partition Key (PK) distributes items <span class="d-metric throughput">1000 WCU/partition</span> <div class="d-tag amber">avoid low-cardinality PKs</div></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">PARTITION A</div>
        <div class="d-flow-v">
          <div class="d-box green" data-tip="Composite key pattern: ENTITY#ID. Enables Query on PK to fetch all items for an entity. GetItem requires both PK + SK.">PK = USER#U123 <div class="d-tag indigo">ENTITY#ID pattern</div></div>
          <div class="d-label">Sort Key orders items:</div>
          <div class="d-box gray" data-tip="SK enables range queries: begins_with, between, >, <. Items sorted by SK in B-tree within partition. Query returns items in SK order.">SK = ORDER#2024-01 <span class="d-metric latency">&lt;10ms query</span></div>
          <div class="d-box gray">SK = ORDER#2024-02</div>
          <div class="d-box gray" data-tip="Singleton item pattern: use a fixed SK like PROFILE for 1:1 entity data alongside 1:N collections.">SK = PROFILE <div class="d-tag gray">singleton</div></div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">1</div><div class="d-number-label">Table for all entity types</div></div>
  <div class="d-number"><div class="d-number-value">0</div><div class="d-number-label">JOINs needed</div></div>
</div>
<div class="d-group">
  <div class="d-group-title">TABLE: ECommerce (Single Table) <div class="d-tag indigo">single-table design</div></div>
  <div class="d-flow-v">
    <div class="d-group">
      <div class="d-group-title">PARTITION: USER#U123</div>
      <div class="d-row">
        <div class="d-box blue" data-tip="1:1 entity. GetItem(PK=USER#U123, SK=PROFILE) — single-digit ms. Max 400KB per item. Store denormalized data to avoid joins.">SK = PROFILE<br>name, email <span class="d-metric latency">&lt;5ms</span> <div class="d-tag gray">1:1</div></div>
        <div class="d-box green" data-tip="1:N collection. Query(PK=USER#U123, SK begins_with ORDER#) returns all orders sorted by date. Paginate with LastEvaluatedKey.">SK = ORDER#2024-01-15#O456<br>total, status <div class="d-tag green">1:N</div></div>
        <div class="d-box green" data-tip="Date-prefixed SK enables time-range queries: SK between ORDER#2024-01 and ORDER#2024-12 for yearly orders.">SK = ORDER#2024-02-20#O789<br>total, status</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">PARTITION: PRODUCT#P789</div>
      <div class="d-row">
        <div class="d-box purple" data-tip="Product details as singleton. GSI on category+price enables browse-by-category queries without scanning.">SK = DETAILS<br>name, price, category</div>
        <div class="d-box amber" data-tip="Reviews stored under product PK. Inverted GSI (PK=USER#U123, SK=REVIEW#P789) enables 'my reviews' query. This is the power of single-table design.">SK = REVIEW#U123<br>rating, text <div class="d-tag amber">GSI for inverted query</div></div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">10x</div><div class="d-number-label">Capacity increase with 10 shards</div></div>
</div>
<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">WITHOUT SHARDING <div class="d-tag amber">hot partition problem</div></div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="10K writes per second all targeting the same partition key.">All Writes <span class="d-metric throughput">10K WPS</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Single PK funnels all writes to one partition. DynamoDB partition limit: 1000 WCU. Above that = ProvisionedThroughputExceededException.">POST#123 (single PK) <span class="d-metric throughput">1000 WCU max</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" data-tip="Hot partition problem. Even with on-demand mode, a single partition tops out at 1000 WCU. Adaptive capacity helps but does not fully solve it.">Single Partition<br>THROTTLED! <span class="d-status error"></span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">WITH SHARDING <div class="d-tag green">recommended</div></div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Same 10K writes, now distributed across 10 shard partitions.">All Writes <span class="d-metric throughput">10K WPS</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Application appends random suffix (0-9) to PK. Distributes writes across 10 partitions. Read requires Scatter-Gather across all shards.">Random Shard (0-9) <span class="d-metric throughput">10x capacity</span> <div class="d-tag indigo">rand() % 10</div></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-row">
          <div class="d-box green" data-tip="Each shard handles 1000 WCU independently. 10 shards = 10,000 WCU effective capacity.">POST#123#0 <span class="d-metric throughput">1K WCU</span> <span class="d-status active"></span></div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">~95%</div><div class="d-number-label">Cache hit rate (well-tuned)</div></div>
  <div class="d-number"><div class="d-number-value">&lt;1ms</div><div class="d-number-label">Redis GET latency</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box blue" data-tip="Application always checks cache first. If key exists and TTL not expired, return immediately. Cache-aside is the most common caching pattern."><span class="d-step">1</span> Application Request</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="GET key — O(1) operation. Sub-ms over ElastiCache in same AZ. Network hop adds ~0.1-0.5ms. Pipeline multiple GETs for batch lookups."><span class="d-step">2</span> Check Redis Cache <span class="d-metric latency">&lt;1ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-label">HIT (~95%)</div>
      <div class="d-box green" data-tip="Cache hit: return immediately. No database query. Typical hit rate: 95-99% for read-heavy workloads with proper TTL tuning."><span class="d-step">3a</span> Return Cached Data <span class="d-metric latency">&lt;1ms</span> <span class="d-status active"></span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-label">MISS (~5%)</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Cache miss triggers DB query. Protect against thundering herd with request coalescing (singleflight) or mutex lock on cache key."><span class="d-step">3b</span> Query Database <span class="d-metric latency">5-50ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="SET key value EX ttl — write result to cache with TTL. TTL prevents stale data. Typical TTL: 5min-24h depending on data freshness requirements."><span class="d-step">4</span> Write Result to Redis <span class="d-metric latency">&lt;1ms</span> <div class="d-tag indigo">SET EX ttl</div></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><span class="d-step">5</span> Return Data <span class="d-status active"></span></div>
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
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">16,384</div><div class="d-number-label">Hash slots total</div></div>
  <div class="d-number"><div class="d-number-value">78 GB</div><div class="d-number-label">Total memory (3 shards × 26 GB)</div></div>
  <div class="d-number"><div class="d-number-value">300K</div><div class="d-number-label">ops/s cluster-wide</div></div>
</div>
<div class="d-label">Cluster Mode Enabled (16,384 hash slots) <div class="d-tag green">zero-downtime resharding</div></div>
<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">SHARD 1 (slots 0-5460)</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Primary handles all writes for its slot range. Max memory: r6g.xlarge = 26GB. Writes replicated async to replicas. If primary fails, replica promoted in 10-30s.">Primary <span class="d-metric size">26 GB</span> <span class="d-metric throughput">100K ops/s</span> <span class="d-status active"></span></div>
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
        <div class="d-box blue" data-tip="Independent primary for this slot range.">Primary <span class="d-metric size">26 GB</span> <span class="d-metric throughput">100K ops/s</span> <span class="d-status active"></span></div>
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
        <div class="d-box blue">Primary <span class="d-metric size">26 GB</span> <span class="d-metric throughput">100K ops/s</span> <span class="d-status active"></span></div>
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

	// -------------------------------------------------------
	// CAP Theorem
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-cap-overview",
		Title:       "CAP Theorem Overview",
		Description: "CP vs AP decision triangle showing the tradeoff between consistency and availability during network partitions",
		ContentFile: "fundamentals/distributed/cap-theorem",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box gray">Network Partition (always possible in distributed systems)</div>
  <div class="d-arrow-down">↓ When it happens, choose:</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">CP — Consistency + Partition Tolerance</div>
        <div class="d-flow-v">
          <div class="d-box blue">Return error or wait</div>
          <div class="d-box blue">Never return stale data</div>
          <div class="d-box blue">Examples: HBase, Zookeeper, etcd</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">AP — Availability + Partition Tolerance</div>
        <div class="d-flow-v">
          <div class="d-box green">Always return a response</div>
          <div class="d-box green">May return stale data</div>
          <div class="d-box green">Examples: Cassandra, DynamoDB, CouchDB</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-caption">CA (no partition tolerance) is impossible in a network — partitions always happen. The real choice is CP vs AP.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-cap-partition-behavior",
		Title:       "System Behavior During a Network Partition",
		Description: "Shows what happens step-by-step when a network partition occurs in CP vs AP systems",
		ContentFile: "fundamentals/distributed/cap-theorem",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-row">
    <div class="d-box gray" data-tip="Node A and Node B can exchange messages normally.">Node A &#8596; Node B <span class="d-status active"></span></div>
    <div class="d-arrow">&#8594; partition &#8594;</div>
    <div class="d-box red" data-tip="Network split: nodes cannot communicate. Both are still running.">Node A &#10007; Node B <span class="d-status error"></span></div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">CP System Response</div>
        <div class="d-flow-v" style="gap:0.5rem">
          <div class="d-box blue" data-tip="Only the majority partition (quorum) accepts writes."><span class="d-step">1</span> Majority partition keeps serving</div>
          <div class="d-box red" data-tip="Minority side returns errors or times out."><span class="d-step">2</span> Minority partition rejects requests</div>
          <div class="d-box blue" data-tip="When network heals, minority catches up from majority."><span class="d-step">3</span> Partition heals &#8594; minority syncs</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">AP System Response</div>
        <div class="d-flow-v" style="gap:0.5rem">
          <div class="d-box green" data-tip="Both sides continue accepting reads and writes."><span class="d-step">1</span> Both partitions keep serving</div>
          <div class="d-box amber" data-tip="Each side may have different data versions."><span class="d-step">2</span> Data diverges between sides</div>
          <div class="d-box green" data-tip="Conflict resolution: LWW, vector clocks, or CRDTs."><span class="d-step">3</span> Partition heals &#8594; reconcile conflicts</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-caption">CP sacrifices availability during partition. AP sacrifices consistency. Both recover when partition heals.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-lb-algorithm-comparison",
		Title:       "Load Balancing Algorithm Comparison",
		Description: "Visual comparison of Round Robin, Least Connections, and Weighted algorithms",
		ContentFile: "fundamentals/networking/load-balancing",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Round Robin</div>
      <div class="d-flow-v" style="gap:0.4rem">
        <div class="d-box blue" data-tip="Requests distributed 1-2-3-1-2-3 in order.">Server 1 &#8592; req 1, 4, 7</div>
        <div class="d-box blue" data-tip="Equal distribution regardless of server load.">Server 2 &#8592; req 2, 5, 8</div>
        <div class="d-box blue" data-tip="Simple but ignores server health and capacity.">Server 3 &#8592; req 3, 6, 9</div>
      </div>
      <div class="d-label">Simple, equal distribution. Ignores server load.</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Least Connections</div>
      <div class="d-flow-v" style="gap:0.4rem">
        <div class="d-box green" data-tip="Currently handling the fewest active connections.">Server 1: 3 active &#8592; next req</div>
        <div class="d-box gray" data-tip="Busy server receives fewer new requests.">Server 2: 8 active</div>
        <div class="d-box gray" data-tip="Moderate load.">Server 3: 5 active</div>
      </div>
      <div class="d-label">Routes to least busy server. Best for varied request durations.</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Weighted Round Robin</div>
      <div class="d-flow-v" style="gap:0.4rem">
        <div class="d-box indigo" data-tip="Large instance gets proportionally more traffic.">Server 1 (w=5): 50% traffic</div>
        <div class="d-box indigo" data-tip="Medium instance.">Server 2 (w=3): 30% traffic</div>
        <div class="d-box indigo" data-tip="Small instance.">Server 3 (w=2): 20% traffic</div>
      </div>
      <div class="d-label">For heterogeneous fleets (different instance sizes).</div>
    </div>
  </div>
</div>`,
	})

	// -------------------------------------------------------
	// Circuit Breaker
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-cb-overview",
		Title:       "Circuit Breaker Overview",
		Description: "How a circuit breaker wraps a dependency call to prevent cascade failures",
		ContentFile: "fundamentals/distributed/circuit-breaker",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Service A (caller)</div>
  <div class="d-arrow-down">↓ all calls go through</div>
  <div class="d-box purple">Circuit Breaker (proxy layer)<br><small>tracks success/failure counts, manages state</small></div>
  <div class="d-arrow-down">↓ forwards when CLOSED</div>
  <div class="d-box green">Service B (dependency)</div>
  <div class="d-caption">When B fails repeatedly, the breaker OPENS and returns errors immediately — no waiting for timeouts. Service A degrades gracefully instead of hanging.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-cb-state-machine",
		Title:       "Circuit Breaker State Machine",
		Description: "CLOSED → OPEN → HALF_OPEN state transitions based on failure thresholds",
		ContentFile: "fundamentals/distributed/circuit-breaker",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-group">
    <div class="d-group-title">CLOSED</div>
    <div class="d-box green">Requests pass through<br>Counting failures</div>
    <div class="d-label">failure rate &gt; threshold →</div>
  </div>
  <div class="d-arrow">→</div>
  <div class="d-group">
    <div class="d-group-title">OPEN</div>
    <div class="d-box red">Fail fast<br>No calls to dependency<br>Return fallback immediately</div>
    <div class="d-label">→ after timeout period</div>
  </div>
  <div class="d-arrow">→</div>
  <div class="d-group">
    <div class="d-group-title">HALF_OPEN</div>
    <div class="d-box amber">Allow 1 probe request<br>Success → CLOSED<br>Failure → OPEN</div>
  </div>
</div>`,
	})

	// -------------------------------------------------------
	// Saga Pattern
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-saga-overview",
		Title:       "Saga Pattern Overview",
		Description: "Sequence of local transactions with compensating transactions on failure",
		ContentFile: "fundamentals/distributed/saga-pattern",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">T1: Reserve Inventory</div>
    <div class="d-arrow">→</div>
    <div class="d-box blue">T2: Charge Customer</div>
    <div class="d-arrow">→</div>
    <div class="d-box blue">T3: Ship Order</div>
    <div class="d-arrow">→</div>
    <div class="d-box green">Done ✓</div>
  </div>
  <div class="d-label">If T3 fails → run compensating transactions backwards:</div>
  <div class="d-flow">
    <div class="d-box red">T3 FAILS</div>
    <div class="d-arrow">→</div>
    <div class="d-box amber">C2: Refund Customer</div>
    <div class="d-arrow">→</div>
    <div class="d-box amber">C1: Release Inventory</div>
    <div class="d-arrow">→</div>
    <div class="d-box gray">Rolled Back</div>
  </div>
  <div class="d-caption">Each step is a local transaction. No distributed lock. Compensating transactions must be idempotent.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-saga-choreography",
		Title:       "Choreography-Based Saga",
		Description: "Event-driven service chain with no central orchestrator",
		ContentFile: "fundamentals/distributed/saga-pattern",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-box blue">Order Service<br><small>emits OrderPlaced</small></div>
  <div class="d-arrow">→</div>
  <div class="d-box green">Inventory Service<br><small>listens OrderPlaced<br>emits InventoryReserved</small></div>
  <div class="d-arrow">→</div>
  <div class="d-box purple">Payment Service<br><small>listens InventoryReserved<br>emits PaymentCharged</small></div>
  <div class="d-arrow">→</div>
  <div class="d-box amber">Shipping Service<br><small>listens PaymentCharged<br>emits OrderShipped</small></div>
</div>
<div class="d-caption">No central coordinator. Each service reacts to events and emits new ones. Simple to add new services, but hard to trace overall saga state.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-saga-orchestration",
		Title:       "Orchestration-Based Saga",
		Description: "Central orchestrator directing services step by step",
		ContentFile: "fundamentals/distributed/saga-pattern",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box indigo">Saga Orchestrator<br><small>owns the saga_log, drives each step</small></div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-arrow">→ step 1</div>
    <div class="d-box blue">Inventory Service</div>
    <div class="d-arrow">← ok/fail</div>
  </div>
  <div class="d-flow">
    <div class="d-arrow">→ step 2</div>
    <div class="d-box green">Payment Service</div>
    <div class="d-arrow">← ok/fail</div>
  </div>
  <div class="d-flow">
    <div class="d-arrow">→ step 3</div>
    <div class="d-box amber">Shipping Service</div>
    <div class="d-arrow">← ok/fail</div>
  </div>
  <div class="d-caption">Orchestrator knows the full saga state. Easy to observe and debug. Single point of coordination — but not a single point of failure if the log is durable.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-saga-compensation",
		Title:       "Saga Compensation Flow",
		Description: "Forward steps and their corresponding compensating transactions",
		ContentFile: "fundamentals/distributed/saga-pattern",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Forward Steps</div>
      <div class="d-box blue">T1: reserve_inventory(order_id)</div>
      <div class="d-box blue">T2: charge_customer(order_id)</div>
      <div class="d-box blue">T3: create_shipment(order_id)</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Compensating Transactions</div>
      <div class="d-box amber">C1: release_inventory(order_id)</div>
      <div class="d-box amber">C2: refund_customer(order_id)</div>
      <div class="d-box amber">C3: cancel_shipment(order_id)</div>
    </div>
  </div>
</div>
<div class="d-caption">Every forward step must have a compensating transaction. Compensations must be idempotent — they may be called more than once if the orchestrator retries.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-saga-recovery",
		Title:       "Saga Crash Recovery",
		Description: "Orchestrator recovers from crash by reading saga_log and replaying from last successful step",
		ContentFile: "fundamentals/distributed/saga-pattern",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box red">Orchestrator Crashes mid-saga</div>
  <div class="d-arrow-down">↓ on restart</div>
  <div class="d-box indigo">Read saga_log from durable store (Postgres / DynamoDB)</div>
  <div class="d-arrow-down">↓ find last committed step</div>
  <div class="d-box amber">Determine: last completed step = T2 (charge_customer)</div>
  <div class="d-arrow-down">↓ replay forward idempotently</div>
  <div class="d-box green">Execute T3: create_shipment (safe — idempotent)</div>
  <div class="d-caption">Each step's idempotency key (order_id) ensures replaying a completed step has no side effects. The log is the source of truth for saga state.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-saga-vs-2pc",
		Title:       "Saga vs 2-Phase Commit",
		Description: "Comparison of 2PC (distributed locks) vs Saga (local commits) for distributed transactions",
		ContentFile: "fundamentals/distributed/saga-pattern",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">2-Phase Commit (2PC)</div>
      <div class="d-box red">Phase 1: coordinator sends PREPARE<br>all participants hold locks</div>
      <div class="d-box red">Phase 2: coordinator sends COMMIT<br>participants release locks</div>
      <div class="d-label">If coordinator crashes in phase 2 → locks held forever</div>
      <div class="d-box amber">Blocking, coordinator SPOF, poor availability</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Saga</div>
      <div class="d-box green">Each step is a local transaction<br>commits immediately, no locks held</div>
      <div class="d-box green">Failures trigger compensating transactions<br>eventual consistency</div>
      <div class="d-label">No distributed lock, no coordinator SPOF</div>
      <div class="d-box blue">Non-blocking, high availability, eventual consistency</div>
    </div>
  </div>
</div>`,
	})

	// -------------------------------------------------------
	// Rate Limiting
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-rl-token-bucket",
		Title:       "Token Bucket Algorithm",
		Description: "Token bucket fill and consume cycle for rate limiting",
		ContentFile: "fundamentals/distributed/rate-limiting",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Refill Thread<br><small>adds tokens at rate R/sec</small></div>
    <div class="d-arrow">→</div>
    <div class="d-box purple">Token Bucket<br><small>capacity C tokens<br>current: 7/10</small></div>
    <div class="d-arrow">→</div>
    <div class="d-box green">Incoming Request<br><small>consumes 1 token</small></div>
  </div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-box green">Token available → allow request</div>
    <div class="d-arrow">vs</div>
    <div class="d-box red">No token → reject with 429</div>
  </div>
  <div class="d-caption">Allows bursts up to capacity C. Sustained rate limited to R tokens/sec. Redis INCR + EX implements a leaky bucket variant in one atomic call.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-rl-distributed",
		Title:       "Distributed Rate Limiter",
		Description: "Multi-server rate limiter with shared Redis counter for consistent limits across the fleet",
		ContentFile: "fundamentals/distributed/rate-limiting",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">App Server 1</div>
    <div class="d-box blue">App Server 2</div>
    <div class="d-box blue">App Server 3</div>
  </div>
  <div class="d-arrow-down">↓ all servers check the same key</div>
  <div class="d-box purple">Redis Cluster<br><small>key: rate:user:12345:2024010312<br>INCR + EXPIRE (atomic Lua script)</small></div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-box green">count ≤ 1000 → allow</div>
    <div class="d-arrow">vs</div>
    <div class="d-box red">count &gt; 1000 → 429 Too Many Requests</div>
  </div>
  <div class="d-caption">Shared Redis counter ensures the limit applies across the entire fleet. Lua script makes INCR+EXPIRE atomic — no race conditions between servers.</div>
</div>`,
	})

	// -------------------------------------------------------
	// Message Queues
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-message-queue-overview",
		Title:       "Message Queue Overview",
		Description: "Producer → Queue → Consumer basic flow with decoupling benefits",
		ContentFile: "fundamentals/messaging/message-queues",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Producer A<br><small>order-service</small></div>
    <div class="d-box blue">Producer B<br><small>payment-service</small></div>
  </div>
  <div class="d-arrow-down">↓ publish messages</div>
  <div class="d-box purple">Message Queue / Topic<br><small>durable, ordered, replicated</small></div>
  <div class="d-arrow-down">↓ consume at own pace</div>
  <div class="d-flow">
    <div class="d-box green">Consumer A<br><small>notification-service</small></div>
    <div class="d-box green">Consumer B<br><small>analytics-service</small></div>
    <div class="d-box green">Consumer C<br><small>audit-service</small></div>
  </div>
  <div class="d-caption">Producers and consumers are decoupled in time and scale. Queue absorbs bursts. Consumers can be added or removed without changing producers.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-kafka-architecture",
		Title:       "Kafka Architecture",
		Description: "Topics, partitions, and consumer groups showing how Kafka scales throughput",
		ContentFile: "fundamentals/messaging/message-queues",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Producers (any service)</div>
  <div class="d-arrow-down">↓ write to partition by key hash</div>
  <div class="d-group">
    <div class="d-group-title">Topic: order-events (3 partitions, replication factor 3)</div>
    <div class="d-flow">
      <div class="d-box purple">Partition 0<br><small>broker 1 (leader)</small></div>
      <div class="d-box purple">Partition 1<br><small>broker 2 (leader)</small></div>
      <div class="d-box purple">Partition 2<br><small>broker 3 (leader)</small></div>
    </div>
  </div>
  <div class="d-arrow-down">↓ each consumer group reads independently</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Consumer Group: notifications</div>
        <div class="d-box green">Consumer 1 → P0</div>
        <div class="d-box green">Consumer 2 → P1</div>
        <div class="d-box green">Consumer 3 → P2</div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Consumer Group: analytics</div>
        <div class="d-box amber">Consumer 1 → P0+P1</div>
        <div class="d-box amber">Consumer 2 → P2</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-exactly-once-flow",
		Title:       "Exactly-Once Delivery Flow",
		Description: "Idempotent producer and transactional consumer achieving exactly-once semantics",
		ContentFile: "fundamentals/messaging/message-queues",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Producer<br><small>idempotent producer ID + sequence number</small></div>
  <div class="d-arrow-down">↓ broker deduplicates using (producer_id, seq)</div>
  <div class="d-box purple">Kafka Broker<br><small>stores message exactly once even on retry</small></div>
  <div class="d-arrow-down">↓ consumer reads + processes atomically</div>
  <div class="d-box green">Consumer<br><small>commits offset + side effect in same DB transaction</small></div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-box gray">BEGIN TRANSACTION</div>
    <div class="d-arrow">→</div>
    <div class="d-box gray">write result to DB</div>
    <div class="d-arrow">→</div>
    <div class="d-box gray">commit offset</div>
    <div class="d-arrow">→</div>
    <div class="d-box gray">COMMIT</div>
  </div>
  <div class="d-caption">Idempotent producer prevents duplicates at broker. Transactional consumer prevents duplicates at sink. Together: exactly-once end-to-end.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-backpressure-handling",
		Title:       "Backpressure Handling",
		Description: "Consumer-driven flow control preventing producer from overwhelming slow consumers",
		ContentFile: "fundamentals/messaging/message-queues",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Producer<br><small>1000 msg/s</small></div>
  <div class="d-arrow-down">↓</div>
  <div class="d-box purple">Queue / Buffer<br><small>growing backlog → signal to slow down</small></div>
  <div class="d-arrow-down">↓ consumer pulls at its own rate</div>
  <div class="d-box green">Consumer<br><small>200 msg/s (slow)</small></div>
  <div class="d-cols" style="margin-top:0.5rem;">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Backpressure Strategies</div>
        <div class="d-box amber">Drop (shed load) — reject new messages when full</div>
        <div class="d-box amber">Block producer — pause until consumer catches up</div>
        <div class="d-box amber">Scale consumers — add more consumer instances</div>
      </div>
    </div>
  </div>
  <div class="d-caption">Backpressure is essential to prevent OOM crashes. Kafka handles it naturally — consumers pull at their own rate, so producers cannot overwhelm them.</div>
</div>`,
	})

	// -------------------------------------------------------
	// WebSockets
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-websocket-overview",
		Title:       "WebSocket Connection Overview",
		Description: "HTTP upgrade handshake establishing a persistent bidirectional WebSocket connection",
		ContentFile: "fundamentals/networking/websockets",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Browser</div>
    <div class="d-arrow">→ HTTP GET /ws<br>Upgrade: websocket</div>
    <div class="d-box green">Server</div>
  </div>
  <div class="d-flow">
    <div class="d-box blue">Browser</div>
    <div class="d-arrow">← 101 Switching Protocols</div>
    <div class="d-box green">Server</div>
  </div>
  <div class="d-arrow-down">↓ TCP connection stays open</div>
  <div class="d-flow">
    <div class="d-box blue">Browser</div>
    <div class="d-arrow">⟷ bidirectional frames<br>(no HTTP overhead)</div>
    <div class="d-box green">Server</div>
  </div>
  <div class="d-caption">After the upgrade, both sides can send frames at any time without HTTP request/response overhead. Latency: ~1ms per frame (vs ~50ms for HTTP polling).</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-websocket-heartbeat",
		Title:       "WebSocket Ping/Pong Heartbeat",
		Description: "Ping/pong keepalive mechanism detecting dead connections",
		ContentFile: "fundamentals/networking/websockets",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box green">Server</div>
    <div class="d-arrow">→ PING frame (every 30s)</div>
    <div class="d-box blue">Client</div>
  </div>
  <div class="d-flow">
    <div class="d-box green">Server</div>
    <div class="d-arrow">← PONG frame</div>
    <div class="d-box blue">Client</div>
  </div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-box red">No PONG within 60s</div>
    <div class="d-arrow">→</div>
    <div class="d-box amber">Server closes connection<br>removes from session map</div>
  </div>
  <div class="d-caption">Heartbeats detect silently disconnected clients (mobile switching networks, NAT timeouts). Without heartbeats, zombie connections accumulate and exhaust file descriptors.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-websocket-scaling",
		Title:       "WebSocket Horizontal Scaling",
		Description: "Multiple WebSocket servers coordinated via Redis pub/sub for cross-server message delivery",
		ContentFile: "fundamentals/networking/websockets",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Client A (on WS Server 1)</div>
  <div class="d-label">sends message to Client B (on WS Server 2)</div>
  <div class="d-flow">
    <div class="d-box green">WS Server 1<br><small>has Client A connection</small></div>
    <div class="d-arrow">→ publish to Redis channel</div>
    <div class="d-box purple">Redis Pub/Sub</div>
    <div class="d-arrow">→ fan out to all servers</div>
    <div class="d-box green">WS Server 2<br><small>has Client B connection<br>delivers to Client B</small></div>
  </div>
  <div class="d-caption">Sticky sessions (via load balancer) bind a client to one WS server. Redis pub/sub routes cross-server messages. Each server only holds connections for its clients.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-redis-pubsub-fanout",
		Title:       "Redis Pub/Sub Fanout",
		Description: "Redis pub/sub message broadcast to multiple WebSocket server subscribers",
		ContentFile: "fundamentals/networking/websockets",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Publisher (any service)<br><small>PUBLISH chat:room:42 "Hello"</small></div>
  <div class="d-arrow-down">↓ Redis fan-out</div>
  <div class="d-box purple">Redis Pub/Sub<br><small>channel: chat:room:42</small></div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-box green">WS Server 1<br><small>subscribed to channel<br>3 connected clients in room 42</small></div>
    <div class="d-box green">WS Server 2<br><small>subscribed to channel<br>5 connected clients in room 42</small></div>
    <div class="d-box green">WS Server 3<br><small>subscribed to channel<br>2 connected clients in room 42</small></div>
  </div>
  <div class="d-caption">Redis pub/sub is fire-and-forget — no persistence. If a server is down when the message is published, it misses it. Use Kafka for durability + pub/sub fan-out.</div>
</div>`,
	})

	// -------------------------------------------------------
	// Sharding
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-shard-range",
		Title:       "Range-Based Sharding",
		Description: "Range-based partition assignment mapping key ranges to specific shards",
		ContentFile: "fundamentals/storage/sharding",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Shard Key: user_id</div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-box green">Shard 0<br><small>user_id 1 – 1,000,000</small></div>
    <div class="d-box green">Shard 1<br><small>user_id 1M – 2,000,000</small></div>
    <div class="d-box green">Shard 2<br><small>user_id 2M – 3,000,000</small></div>
    <div class="d-box amber">Shard 3 (hot!)<br><small>user_id 3M – 4,000,000<br>all new signups land here</small></div>
  </div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-box blue">Range query user_id 500K–600K</div>
    <div class="d-arrow">→ hits only Shard 0</div>
    <div class="d-box green">Shard 0</div>
  </div>
  <div class="d-caption">Range sharding preserves locality — range queries hit one shard. But monotonic keys (auto-increment, timestamp) always write to the latest shard: hotspot risk.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-shard-hash",
		Title:       "Hash-Based Sharding",
		Description: "Hash-based partition assignment distributing keys uniformly across shards",
		ContentFile: "fundamentals/storage/sharding",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">shard_id = hash(user_id) mod N</div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-box blue">user_id = 1001<br><small>hash → shard 2</small></div>
    <div class="d-box blue">user_id = 1002<br><small>hash → shard 0</small></div>
    <div class="d-box blue">user_id = 1003<br><small>hash → shard 1</small></div>
  </div>
  <div class="d-flow" style="margin-top:0.5rem;">
    <div class="d-box green">Shard 0<br><small>~equal load</small></div>
    <div class="d-box green">Shard 1<br><small>~equal load</small></div>
    <div class="d-box green">Shard 2<br><small>~equal load</small></div>
  </div>
  <div class="d-caption">Even distribution — no hotspots. Tradeoff: range queries must scatter to all shards. Adding a shard with simple mod moves (N-1)/N ≈ 80% of keys. Use consistent hashing instead.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-shard-directory",
		Title:       "Directory-Based Sharding",
		Description: "Lookup-table based routing mapping each key to a specific shard",
		ContentFile: "fundamentals/storage/sharding",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Incoming Request<br><small>key: user:12345</small></div>
  <div class="d-arrow-down">↓ lookup (cached locally)</div>
  <div class="d-box purple">Directory Service / Lookup Table<br><small>user:12345 → shard3<br>user:67890 → shard1</small></div>
  <div class="d-arrow-down">↓ route to correct shard</div>
  <div class="d-flow">
    <div class="d-box gray">Shard 1</div>
    <div class="d-box gray">Shard 2</div>
    <div class="d-box green">Shard 3 ← user:12345</div>
    <div class="d-box gray">Shard 4</div>
  </div>
  <div class="d-caption">Most flexible — any key can move to any shard without math. Cache the lookup table locally in each app server (4MB for 1M users). Downside: lookup service is a potential SPOF.</div>
</div>`,
	})

	// -------------------------------------------------------
	// Geospatial
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-geospatial-overview",
		Title:       "Geospatial Indexing Overview",
		Description: "2D coordinates mapped to 1D index key enabling efficient proximity queries",
		ContentFile: "fundamentals/storage/geospatial",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Without Spatial Index</div>
        <div class="d-box red">WHERE lat BETWEEN 40.7 AND 40.8<br>AND lng BETWEEN -74.0 AND -73.9</div>
        <div class="d-label">→ Full table scan on 2 columns</div>
        <div class="d-box red">O(N) — unusable at scale</div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">With Spatial Index</div>
        <div class="d-box green">Map (lat, lng) → 1D key<br><small>GeoHash string / S2 cell ID / H3 cell ID</small></div>
        <div class="d-label">→ Range scan on single key</div>
        <div class="d-box green">O(log N) — fast at any scale</div>
      </div>
    </div>
  </div>
  <div class="d-caption">Key insight: nearby points in 2D produce nearby keys in 1D (locality-preserving mapping). This converts a 2D proximity query into an efficient 1D range scan.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-geospatial-geohash",
		Title:       "GeoHash Encoding",
		Description: "World grid with base32 encoding hierarchy showing prefix-based proximity",
		ContentFile: "fundamentals/storage/geospatial",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">(40.7484, -73.9857)<br><small>Times Square, NYC</small></div>
    <div class="d-arrow">→ encode</div>
    <div class="d-box purple">dr5ru<br><small>5-char GeoHash ≈ 4.9km cell</small></div>
    <div class="d-arrow">→ longer</div>
    <div class="d-box green">dr5ruj<br><small>6-char ≈ 1.2km cell</small></div>
  </div>
  <div class="d-label">Nearby points share prefix:</div>
  <div class="d-flow">
    <div class="d-box green">dr5ruj = Times Square</div>
    <div class="d-box green">dr5rum = 500m away</div>
    <div class="d-box amber">dr5q = different neighborhood<br><small>no common prefix at 5 chars</small></div>
  </div>
  <div class="d-caption">Shared prefix = geographic proximity. But check 8 neighboring cells for boundary safety — two points 1m apart can be on opposite sides of a cell edge.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-geospatial-s2",
		Title:       "S2 Cell Hierarchy",
		Description: "S2 sphere cell hierarchy from cube face to centimeter-scale cells",
		ContentFile: "fundamentals/storage/geospatial",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Level 0<br><small>6 cube faces<br>~85M km² each</small></div>
    <div class="d-arrow">→ subdivide</div>
    <div class="d-box green">Level 10<br><small>~100 km²<br>city district</small></div>
    <div class="d-arrow">→ subdivide</div>
    <div class="d-box purple">Level 14<br><small>~0.6 km²<br>driver matching</small></div>
    <div class="d-arrow">→ subdivide</div>
    <div class="d-box amber">Level 30<br><small>~1 cm²<br>precise point</small></div>
  </div>
  <div class="d-label">Each cell has a unique 64-bit integer ID — store as BIGINT in any DB:</div>
  <div class="d-flow">
    <div class="d-box blue">Point (lat, lng)</div>
    <div class="d-arrow">→</div>
    <div class="d-box purple">S2 Cell ID (uint64)<br><small>3845968832901234688</small></div>
    <div class="d-arrow">→</div>
    <div class="d-box green">B-tree range scan<br><small>containment query</small></div>
  </div>
  <div class="d-caption">Used by Google Maps, Pokémon GO for polygon containment. Sphere-aware — no latitude distortion unlike GeoHash rectangles.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-geospatial-h3",
		Title:       "H3 Hexagonal Grid",
		Description: "Hexagonal H3 grid cells for equal-area spatial analytics",
		ContentFile: "fundamentals/storage/geospatial",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Resolution 5<br><small>252 km² per cell<br>city district</small></div>
    <div class="d-arrow">→ finer</div>
    <div class="d-box purple">Resolution 8<br><small>0.74 km² per cell<br>Uber surge zones</small></div>
    <div class="d-arrow">→ finer</div>
    <div class="d-box green">Resolution 11<br><small>0.0003 km² per cell<br>building scale</small></div>
  </div>
  <div class="d-label">Equal-area: every cell at the same resolution covers the same area:</div>
  <div class="d-flow">
    <div class="d-box amber">H3 cell at equator<br>0.74 km²</div>
    <div class="d-arrow">=</div>
    <div class="d-box amber">H3 cell near pole<br>0.74 km²</div>
    <div class="d-arrow">→</div>
    <div class="d-box green">density ratio = supply/demand<br>comparable everywhere</div>
  </div>
  <div class="d-caption">Uber uses H3 resolution 8 for surge pricing. Each hexagon has exactly 6 equidistant neighbors — smoother spatial aggregation than rectangles.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-geospatial-redis-geo",
		Title:       "Redis GEO Sorted Set",
		Description: "Redis GEO sorted set with 52-bit GeoHash score for real-time proximity queries",
		ContentFile: "fundamentals/storage/geospatial",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">GEOADD drivers:active<br><small>(-73.9857, 40.7484, "driver:1001")</small></div>
    <div class="d-arrow">→</div>
    <div class="d-box purple">Sorted Set<br><small>score = 52-bit GeoHash<br>member = "driver:1001"</small></div>
  </div>
  <div class="d-label">GEOSEARCH (replaces deprecated GEORADIUS):</div>
  <div class="d-flow">
    <div class="d-box green">GEOSEARCH drivers:active<br><small>FROMLONLAT -73.99 40.75<br>BYRADIUS 5 km ASC COUNT 10</small></div>
    <div class="d-arrow">→</div>
    <div class="d-box green">Returns 10 nearest drivers<br><small>O(N+log M) — sub-ms</small></div>
  </div>
  <div class="d-caption">No per-member TTL. Use a parallel Redis hash (drivers:last_seen) + background job to remove stale drivers older than 60s. Accuracy: ~0.6 meters.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-geospatial-postgis",
		Title:       "PostGIS Spatial Index Query Flow",
		Description: "PostGIS GIST index query flow with bounding-box pre-filter and exact distance computation",
		ContentFile: "fundamentals/storage/geospatial",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">ST_DWithin(location, point, 5000)<br><small>find all drivers within 5km</small></div>
  <div class="d-arrow-down">↓ GIST index: bounding box pre-filter</div>
  <div class="d-box amber">Bounding box scan<br><small>~1000 candidate rows from index</small></div>
  <div class="d-arrow-down">↓ exact distance computation</div>
  <div class="d-box green">Haversine distance filter<br><small>~50 rows within exact 5km radius</small></div>
  <div class="d-arrow-down">↓ ORDER BY dist_meters LIMIT 10</div>
  <div class="d-box green">Top 10 nearest drivers<br><small>total: ~10ms at 1M rows</small></div>
  <div class="d-caption">GIST index reduces O(N) full scan to O(log N) bounding-box lookup. Exact distance only computed on the small candidate set. 10x slower than Redis GEO but durable.</div>
</div>`,
	})

	// -------------------------------------------------------
	// Consistent Hashing (Fundamental)
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-consistent-hashing-overview",
		Title:       "Consistent Hashing Ring",
		Description: "Hash ring with node positions and key assignment to nearest clockwise node",
		ContentFile: "fundamentals/distributed/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-label">Hash ring [0, 2³²): nodes and keys mapped to positions</div>
  <div class="d-flow">
    <div class="d-box blue">Node A<br><small>position 10</small></div>
    <div class="d-box green">Node B<br><small>position 40</small></div>
    <div class="d-box purple">Node C<br><small>position 70</small></div>
  </div>
  <div class="d-label">Keys assigned to first node clockwise from their hash position:</div>
  <div class="d-flow">
    <div class="d-box amber">Key "user:1"<br><small>hash=15 → Node B (pos 40)</small></div>
    <div class="d-box amber">Key "user:2"<br><small>hash=55 → Node C (pos 70)</small></div>
    <div class="d-box amber">Key "user:3"<br><small>hash=80 → Node A (pos 10, wraps)</small></div>
  </div>
  <div class="d-caption">Adding a node only affects keys between the new node and its predecessor. O(K/N) keys move — not all keys like with modulo hashing.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-consistent-hashing-vnodes",
		Title:       "Virtual Nodes",
		Description: "Virtual nodes spread across ring for even load distribution",
		ContentFile: "fundamentals/distributed/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Without vnodes (uneven)</div>
        <div class="d-box red">Node A: handles 60% of keys</div>
        <div class="d-box amber">Node B: handles 25% of keys</div>
        <div class="d-box amber">Node C: handles 15% of keys</div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">With vnodes (even)</div>
        <div class="d-box blue">Node A: 150 positions on ring<br><small>handles ~33% of keys</small></div>
        <div class="d-box green">Node B: 150 positions on ring<br><small>handles ~33% of keys</small></div>
        <div class="d-box purple">Node C: 150 positions on ring<br><small>handles ~33% of keys</small></div>
      </div>
    </div>
  </div>
  <div class="d-caption">Cassandra uses 256 vnodes per node by default. More positions = better statistical balance. Adding a node steals a few vnodes from each existing node — balanced migration.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-consistent-hashing-rebalance",
		Title:       "Rebalancing on Node Add/Remove",
		Description: "Key movement when adding or removing a node from the consistent hash ring",
		ContentFile: "fundamentals/distributed/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Add Node D (pos 55)</div>
      <div class="d-box blue">Node A (pos 10)</div>
      <div class="d-box green">Node B (pos 40)</div>
      <div class="d-box purple">Node D NEW (pos 55)<br><small>steals keys 40–55 from Node C</small></div>
      <div class="d-box purple">Node C (pos 70)<br><small>loses keys 40–55 only</small></div>
      <div class="d-label">Only 1/N keys move → minimal disruption</div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Remove Node B (pos 40)</div>
      <div class="d-box blue">Node A (pos 10)</div>
      <div class="d-box red">Node B REMOVED (pos 40)</div>
      <div class="d-box purple">Node C (pos 70)<br><small>inherits keys 10–40 from B</small></div>
      <div class="d-label">Only B's keys move → rest unaffected</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-consistent-hashing-hotspot",
		Title:       "Hotspot Detection and Mitigation",
		Description: "Identifying and mitigating hot keys on the consistent hash ring",
		ContentFile: "fundamentals/distributed/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box red">Hot key: "celebrity:user:12345"<br><small>1M requests/sec all land on Node A</small></div>
  <div class="d-arrow-down">↓ mitigation strategies</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Write Sharding</div>
        <div class="d-box amber">Split into N virtual keys<br>user:12345#0...user:12345#9</div>
        <div class="d-box amber">Writes distributed across 10 nodes</div>
        <div class="d-box amber">Reads scatter-gather + merge</div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Local Cache</div>
        <div class="d-box green">Cache hot key in each app server</div>
        <div class="d-box green">Skip ring lookup for cached hits</div>
        <div class="d-box green">TTL invalidation (10s)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-consistent-hashing-failure",
		Title:       "Node Failure and Key Rerouting",
		Description: "Consistent hash ring behavior when a node fails and keys are rerouted",
		ContentFile: "fundamentals/distributed/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Node A (pos 10)</div>
    <div class="d-box red">Node B FAILED (pos 40)</div>
    <div class="d-box purple">Node C (pos 70)</div>
  </div>
  <div class="d-arrow-down">↓ Node B removed from ring</div>
  <div class="d-flow">
    <div class="d-box blue">Node A (pos 10)</div>
    <div class="d-box purple">Node C (pos 70)<br><small>inherits all of B's keys (pos 10–40)</small></div>
  </div>
  <div class="d-label">With replication factor R=3: replicas on next 2 clockwise nodes</div>
  <div class="d-flow">
    <div class="d-box amber">Key at pos 25 → primary: B (failed)</div>
    <div class="d-arrow">→</div>
    <div class="d-box green">Promoted replica: C<br><small>no data loss with R=3</small></div>
  </div>
  <div class="d-caption">With replication factor 3, node failure causes zero data loss. Keys served from surviving replicas within seconds of failure detection.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-consistent-hashing-summary",
		Title:       "Consistent Hashing Summary",
		Description: "Full ring with replication and key assignment walkthrough",
		ContentFile: "fundamentals/distributed/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-label">Full system: 4 nodes, 150 vnodes each, replication factor 3</div>
  <div class="d-flow">
    <div class="d-box blue">Node A<br><small>150 positions</small></div>
    <div class="d-box green">Node B<br><small>150 positions</small></div>
    <div class="d-box purple">Node C<br><small>150 positions</small></div>
    <div class="d-box amber">Node D<br><small>150 positions</small></div>
  </div>
  <div class="d-label">Key "user:999" flow:</div>
  <div class="d-flow">
    <div class="d-box gray">hash("user:999") → pos 234</div>
    <div class="d-arrow">→</div>
    <div class="d-box green">Primary: Node B (nearest clockwise vnode)</div>
  </div>
  <div class="d-flow">
    <div class="d-box purple">Replica 1: Node C</div>
    <div class="d-box amber">Replica 2: Node D</div>
    <div class="d-label">(next 2 clockwise nodes)</div>
  </div>
  <div class="d-caption">3 copies survive any single node failure. Adding Node E moves only 1/5 of keys. DynamoDB, Cassandra, and Redis Cluster all use consistent hashing variants.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-redis-failover",
		Title:       "Redis Failover Flow",
		Description: "Redis primary node failure and automatic failover sequence through Sentinel/ElastiCache detection, replica promotion, DNS update, and client reconnection",
		ContentFile: "fundamentals/storage/redis",
		Type:        TypeHTML,
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">20-60s</div><div class="d-number-label">Total failover window</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box red" data-tip="Primary crash: OOM kill, hardware failure, or network partition. All in-flight writes lost unless wait/WAIT command used for synchronous replication."><span class="d-step">1</span> Primary Node Crashes <span class="d-status error"></span> <span class="d-metric latency">T+0s</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="ElastiCache monitors via heartbeat every 1s. After N missed heartbeats (default 5), node marked as failing. Sentinel quorum (majority) must agree before failover."><span class="d-step">2</span> Sentinel / ElastiCache Detects Failure <span class="d-status error"></span> <span class="d-metric latency">10-30s</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box blue" data-tip="Replica with least replication lag selected. SLAVEOF NO ONE issued. Replica loads its dataset and begins accepting writes. Other replicas re-pointed to new primary."><span class="d-step">3</span> Promote Replica to Primary <span class="d-status active"></span> <span class="d-metric latency">5-10s</span> <div class="d-tag green">least-lag replica wins</div></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="ElastiCache primary endpoint DNS updated to point to new primary IP. DNS TTL is 5s. Clients using primary endpoint auto-resolve to new primary."><span class="d-step">4</span> Update DNS Endpoint <span class="d-status active"></span> <span class="d-metric latency">5-15s</span> <div class="d-tag indigo">TTL 5s</div></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="Clients reconnect using primary endpoint. Total downtime: 20-60s. Use retry with exponential backoff. Application should handle ReadOnlyError during transition."><span class="d-step">5</span> Clients Reconnect Automatically <span class="d-status active"></span> <span class="d-metric latency">Total: 20-60s</span></div>
  <div class="d-caption">Total failover window: 20-60 seconds. During this window, writes fail and reads may return stale data. Design applications to handle transient Redis unavailability gracefully.</div>
</div>`,
	})

	// -------------------------------------------------------
	// S3 / Object Storage
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-s3-overview",
		Title:       "S3 Object Storage Overview",
		Description: "Shows S3 bucket structure, durability via multi-AZ replication, and the three core operations (PUT/GET/DELETE)",
		ContentFile: "fundamentals/storage/s3",
		Type:        TypeHTML,
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">11 9s</div><div class="d-number-label">Durability (per object/year)</div></div>
  <div class="d-number"><div class="d-number-value">$0.023</div><div class="d-number-label">Per GB/month (Standard)</div></div>
  <div class="d-number"><div class="d-number-value">5 TB</div><div class="d-number-label">Max object size</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box blue" data-tip="Client (browser, app server, CDN). Three operations only: PUT (write), GET (read), DELETE. No in-place edit — must write a new object.">Client</div>
  <div class="d-arrow-down">&#8595; PUT / GET / DELETE</div>
  <div class="d-box indigo" data-tip="Globally unique name across all AWS accounts. Flat namespace — no real directories, just key prefixes. Versioning optional. One region per bucket.">S3 Bucket <span class="d-metric">flat key namespace</span></div>
  <div class="d-arrow-down">&#8595; automatic replication</div>
  <div class="d-flow">
    <div class="d-box green" data-tip="Each AZ is a separate data center with independent power and networking. Multi-AZ replication is automatic, synchronous, invisible to the caller.">AZ-1 replica</div>
    <div class="d-box green" data-tip="S3 replicates to at least 3 AZs within the region. All copies are identical. Loss of any 2 AZs still serves data from the third.">AZ-2 replica</div>
    <div class="d-box green" data-tip="Third copy ensures 11-nines durability. Probability of losing all 3 simultaneously is vanishingly small (~10^-11 per year per object).">AZ-3 replica</div>
  </div>
  <div class="d-caption">S3 objects are immutable: no partial write, no in-place update. Write a new object at the same key to 'update'. Three AZ replicas give 11-nines durability. At 1M objects, expected loss: 0.00001 objects/year.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-s3-presigned",
		Title:       "Pre-Signed URL Flow",
		Description: "Client upload and download flow using pre-signed URLs — server generates signed URL, client talks directly to S3",
		ContentFile: "fundamentals/storage/s3",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Upload Flow (Pre-Signed PUT)</div>
    <div class="d-flow">
      <div class="d-box blue" data-tip="Client requests permission to upload — does not send the file yet.">Client <span class="d-step">1</span> POST /upload-url</div>
      <div class="d-arrow">→</div>
      <div class="d-box indigo" data-tip="Generates pre-signed PUT URL in &lt;1ms (pure crypto, no network call to S3). Includes: bucket, key, expiry, HMAC-SHA256 signature.">App Server <span class="d-step">2</span> sign URL</div>
    </div>
    <div class="d-flow">
      <div class="d-box blue" data-tip="Client receives the signed URL and uploads directly — app server is NOT in the data path.">Client <span class="d-step">3</span> PUT to S3 URL</div>
      <div class="d-arrow">→</div>
      <div class="d-box green" data-tip="S3 validates the HMAC signature. If valid and not expired, stores the object. No AWS credentials exposed to client.">S3 <span class="d-tag green">direct upload</span> <span class="d-metric latency">~10 Gbps</span></div>
    </div>
  </div>
  <div class="d-group" style="margin-top: 0.75rem;">
    <div class="d-group-title">Download Flow (CloudFront + OAC)</div>
    <div class="d-flow">
      <div class="d-box blue">User</div>
      <div class="d-arrow">→</div>
      <div class="d-box amber" data-tip="CloudFront edge PoP. Checks cache. HIT: serves from edge (~10ms). MISS: fetches from S3 origin, caches, serves.">CloudFront Edge <span class="d-metric latency">~10ms HIT</span></div>
      <div class="d-arrow">→ MISS only</div>
      <div class="d-box green" data-tip="S3 bucket is NOT public. OAC (Origin Access Control) allows only CloudFront service principal to GET objects. Direct S3 URL returns 403.">S3 (private bucket) <span class="d-tag indigo">OAC</span></div>
    </div>
  </div>
  <div class="d-caption">Pre-signed URL removes app server from data path — at 1000 concurrent 5MB uploads, your server is not the bottleneck. OAC ensures S3 is never publicly accessible.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-s3-lifecycle",
		Title:       "S3 Lifecycle Tiering",
		Description: "Automatic cost optimization: objects move from Standard to IA to Glacier as they age",
		ContentFile: "fundamentals/storage/s3",
		Type:        TypeHTML,
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">$0.023</div><div class="d-number-label">Standard /GB/mo</div></div>
  <div class="d-number"><div class="d-number-value">$0.00099</div><div class="d-number-label">Deep Archive /GB/mo</div></div>
  <div class="d-number"><div class="d-number-value">23x</div><div class="d-number-label">Cost reduction to archive</div></div>
</div>
<div class="d-flow">
  <div class="d-box blue" data-tip="$0.023/GB/mo. Immediate access. Best for recently uploaded user media accessed frequently in first 30 days.">Standard <div class="d-tag blue">Day 0-30</div> <div class="d-metric">$0.023/GB</div></div>
  <div class="d-arrow">→ 30 days</div>
  <div class="d-box amber" data-tip="$0.0125/GB/mo. Immediate access but $0.01/GB retrieval fee. Good for media accessed monthly — thumbnails, user profile pics from 1-3 months ago.">Standard-IA <div class="d-tag amber">Day 30-90</div> <div class="d-metric">$0.0125/GB</div></div>
  <div class="d-arrow">→ 90 days</div>
  <div class="d-box purple" data-tip="$0.004/GB/mo. Millisecond retrieval but higher per-GB-retrieved cost. Good for compliance scans, quarterly archive access.">Glacier Instant <div class="d-tag purple">Day 90-365</div> <div class="d-metric">$0.004/GB</div></div>
  <div class="d-arrow">→ 1 year</div>
  <div class="d-box gray" data-tip="$0.00099/GB/mo. 12-48h retrieval. Best for regulatory archives (7-year retention for financial records). 23x cheaper than Standard.">Deep Archive <div class="d-tag gray">Year 1+</div> <div class="d-metric">$0.00099/GB</div></div>
</div>
<div class="d-caption">Lifecycle policies execute automatically — no code changes. 1 PB of media: Standard forever = $23K/mo. With lifecycle to Deep Archive after 1 year = ~$4K/mo. $19K/month savings at Instagram scale.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-s3-events",
		Title:       "S3 Event Notification Pipeline",
		Description: "S3 upload triggers Lambda/SQS/SNS for async processing: image resize, video transcode, content moderation",
		ContentFile: "fundamentals/storage/s3",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="User uploads file directly to S3 via pre-signed PUT URL. Upload completes — HTTP 200 returned to user immediately.">User uploads to S3 (pre-signed PUT) <span class="d-step">1</span></div>
  <div class="d-arrow-down">&#8595; s3:ObjectCreated event</div>
  <div class="d-box amber" data-tip="S3 emits an event notification within milliseconds of the PUT completing. Targets: Lambda (sync invoke), SQS (queue), SNS (fan-out to multiple subscribers).">S3 Event Notification <span class="d-step">2</span> <span class="d-metric latency">&lt;1s</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box green" data-tip="Direct Lambda invocation. Good for simple transforms: resize image to 3 sizes. Max 15 min execution. At-least-once — Lambda is idempotent.">Lambda <div class="d-tag green">image resize</div></div>
    <div class="d-box indigo" data-tip="SQS buffers events for batch processing. Lambda polls SQS (up to 10 messages per batch). Handles spiky upload bursts without throttling. DLQ for failures.">SQS → Lambda Workers <div class="d-tag indigo">video transcode</div></div>
    <div class="d-box red" data-tip="SNS fans out to multiple subscribers: moderation queue, analytics queue, email notification. One S3 event reaches N downstream consumers.">SNS → Moderation + Analytics <div class="d-tag red">content safety</div></div>
  </div>
  <div class="d-caption">S3 events are at-least-once delivery — design Lambda/worker functions to be idempotent. Same image may be processed twice on retry. Check output existence before reprocessing.</div>
</div>`,
	})

	// -------------------------------------------------------
	// Elasticsearch / OpenSearch
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-es-overview",
		Title:       "Elasticsearch Inverted Index",
		Description: "How Elasticsearch's inverted index enables full-text search: tokenization, indexing, and query execution",
		ContentFile: "fundamentals/storage/elasticsearch",
		Type:        TypeHTML,
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">&lt;50ms</div><div class="d-number-label">Full-text search latency</div></div>
  <div class="d-number"><div class="d-number-value">O(1)</div><div class="d-number-label">Inverted index lookup</div></div>
</div>
<div class="d-cols" style="gap: 1rem;">
  <div class="d-flow-v">
    <div class="d-label" style="font-weight: 600;">Documents (source)</div>
    <div class="d-box gray" style="font-size: 0.8rem;" data-tip="Doc 1: contains words 'redis', 'cache', 'fast'">Doc 1: "Redis cache is fast"</div>
    <div class="d-box gray" style="font-size: 0.8rem;" data-tip="Doc 2: contains words 'redis', 'sorted', 'sets'">Doc 2: "Redis sorted sets"</div>
    <div class="d-box gray" style="font-size: 0.8rem;" data-tip="Doc 3: contains words 'cache', 'invalidation'">Doc 3: "Cache invalidation"</div>
  </div>
  <div class="d-arrow">→ index</div>
  <div class="d-flow-v">
    <div class="d-label" style="font-weight: 600;">Inverted Index</div>
    <div class="d-box blue" style="font-size: 0.8rem;" data-tip="Term 'redis' appears in Doc 1 (pos 1) and Doc 2 (pos 1). O(1) lookup — no scan needed.">"redis" → [doc1, doc2]</div>
    <div class="d-box blue" style="font-size: 0.8rem;" data-tip="Term 'cache' appears in Doc 1 and Doc 3.">"cache" → [doc1, doc3]</div>
    <div class="d-box blue" style="font-size: 0.8rem;">"fast"  → [doc1]</div>
    <div class="d-box blue" style="font-size: 0.8rem;">"sorted" → [doc2]</div>
  </div>
  <div class="d-arrow">→ query "redis cache"</div>
  <div class="d-flow-v">
    <div class="d-label" style="font-weight: 600;">Result (BM25 ranked)</div>
    <div class="d-box green" data-tip="Doc 1 contains both 'redis' AND 'cache' — highest BM25 score.">#1 Doc 1 <span class="d-tag green">both terms</span></div>
    <div class="d-box amber" data-tip="Doc 2 has 'redis', Doc 3 has 'cache' — single term match, lower score.">#2 Doc 2, Doc 3</div>
  </div>
</div>
<div class="d-caption">Inverted index: during indexing, each term maps to the list of documents containing it. At query time: look up each term → intersect posting lists → score by BM25. O(1) per term lookup vs O(N) full scan in SQL.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-es-sharding",
		Title:       "Elasticsearch Shard Architecture",
		Description: "Index split into primary shards across nodes, with replica shards for fault tolerance and read scaling",
		ContentFile: "fundamentals/storage/elasticsearch",
		Type:        TypeHTML,
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">30-50 GB</div><div class="d-number-label">Target per shard</div></div>
  <div class="d-number"><div class="d-number-value">3</div><div class="d-number-label">Master-eligible nodes (odd)</div></div>
</div>
<div class="d-flow-v">
  <div class="d-box blue" data-tip="Client query is sent to any node (coordinating node). It fans out the query to all shards in parallel, collects results, merges, ranks, and returns top-N.">Client Query</div>
  <div class="d-arrow-down">&#8595; coordinating node fans out</div>
  <div class="d-flow">
    <div class="d-flow-v">
      <div class="d-label">Node 1</div>
      <div class="d-box green" data-tip="Primary shard 0 — stores 1/3 of the index data. Handles writes and reads.">P0 (primary)</div>
      <div class="d-box amber" data-tip="Replica of shard 1 — promoted to primary if Node 2 fails. Handles read requests.">R1 (replica)</div>
    </div>
    <div class="d-flow-v">
      <div class="d-label">Node 2</div>
      <div class="d-box green" data-tip="Primary shard 1 — stores another 1/3 of index data.">P1 (primary)</div>
      <div class="d-box amber" data-tip="Replica of shard 2 — fault tolerance for Node 3 failure.">R2 (replica)</div>
    </div>
    <div class="d-flow-v">
      <div class="d-label">Node 3</div>
      <div class="d-box green" data-tip="Primary shard 2 — final 1/3 of index data.">P2 (primary)</div>
      <div class="d-box amber" data-tip="Replica of shard 0 — fault tolerance for Node 1 failure.">R0 (replica)</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; merge + rank results</div>
  <div class="d-box blue">Top-N results returned to client</div>
  <div class="d-caption">3 primaries × 1 replica = 6 total shards across 3 nodes. Each primary + its replica are never on the same node. Losing any 1 node = no data loss. Query fans out to all 3 primaries in parallel — O(1/N) latency vs single shard.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-es-indexing-pipeline",
		Title:       "Elasticsearch Indexing Pipeline",
		Description: "Primary DB as source of truth → CDC → Kafka → sync worker → Elasticsearch derived index",
		ContentFile: "fundamentals/storage/elasticsearch",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="PostgreSQL or DynamoDB. Source of truth. All writes go here first. Never write to ES directly as the primary.">Primary Database <span class="d-tag blue">source of truth</span></div>
  <div class="d-arrow-down">&#8595; CDC (Debezium / DynamoDB Streams)</div>
  <div class="d-box amber" data-tip="Change Data Capture captures every INSERT/UPDATE/DELETE as an event. Debezium for Postgres (reads WAL). DynamoDB Streams for DynamoDB. At-least-once delivery.">Change Events Stream <span class="d-metric latency">&lt;100ms lag</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" data-tip="Kafka buffers CDC events. Consumer group 'es-sync' processes events. Dead letter queue for failed indexing attempts. Partition by entity_id for ordering.">Kafka / Kinesis <span class="d-tag indigo">buffer + replay</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" data-tip="Sync worker: transform DB row to ES document format, call ES index API. Batch index (bulk API) for throughput. Handle ES backpressure with Kafka consumer pause.">ES Sync Worker <span class="d-metric latency">~1s total lag</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Derived read index. Eventually consistent — 1s lag vs DB. Never the source of truth. Re-index from DB snapshot if ES index corrupts.">Elasticsearch Index <span class="d-tag purple">derived, eventually consistent</span></div>
  <div class="d-caption">1-second end-to-end indexing lag (DB write → searchable). Design applications to tolerate this: 'your tweet is being indexed...' is acceptable. Use force_refresh only for critical paths (costs ~50ms per refresh).</div>
</div>`,
	})

	// -------------------------------------------------------
	// HTTPS & TLS
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-https-tls-handshake",
		Title:       "TLS 1.3 Handshake",
		Description: "Full TLS 1.3 handshake showing 1-RTT connection establishment vs TLS 1.2's 2-RTT",
		ContentFile: "fundamentals/networking/https",
		Type:        TypeHTML,
		HTML: `<div class="d-cols" style="grid-template-columns: 1fr 1fr; gap: 2rem;">
  <div class="d-flow-v">
    <div class="d-group-title">TLS 1.3 — 1 Round Trip</div>
    <div class="d-box blue" data-tip="Client sends ClientHello with supported cipher suites, key_share (Diffie-Hellman public key), and supported versions.">Client → ClientHello<br><small>key_share + cipher list</small></div>
    <div class="d-arrow-down"><span class="d-step">1 RTT</span></div>
    <div class="d-box green" data-tip="Server picks cipher, responds with its key_share. Both sides now derive the session key. Server immediately sends Certificate + Finished.">Server → ServerHello<br>+ Certificate + Finished<br><small>derives session key here</small></div>
    <div class="d-arrow-down"></div>
    <div class="d-box amber" data-tip="Client verifies certificate, sends Finished. Encrypted application data can now flow.">Client → Finished<br><span class="d-tag green">Encrypted data starts</span></div>
    <div class="d-caption">1-RTT: ~50ms faster than TLS 1.2. 0-RTT resumption skips the round trip entirely for known servers.</div>
  </div>
  <div class="d-flow-v">
    <div class="d-group-title">TLS 1.2 — 2 Round Trips (legacy)</div>
    <div class="d-box gray" data-tip="Client sends ClientHello with cipher suites and random bytes.">Client → ClientHello</div>
    <div class="d-arrow-down"><span class="d-step">RTT 1</span></div>
    <div class="d-box gray" data-tip="Server picks cipher, sends its public certificate and ServerHelloDone.">Server → ServerHello<br>+ Certificate + Done</div>
    <div class="d-arrow-down"></div>
    <div class="d-box gray" data-tip="Client verifies certificate, encrypts a pre-master secret with the server's public key, sends ChangeCipherSpec.">Client → KeyExchange<br>+ ChangeCipherSpec</div>
    <div class="d-arrow-down"><span class="d-step">RTT 2</span></div>
    <div class="d-box gray" data-tip="Server decrypts pre-master secret, both derive session key, server sends ChangeCipherSpec + Finished.">Server → Finished</div>
    <div class="d-arrow-down"></div>
    <div class="d-box gray">Client → Finished<br><span class="d-tag amber">Encrypted data starts</span></div>
    <div class="d-caption">2-RTT: adds 100-200ms on top of TCP. TLS 1.2 is deprecated — migrate to TLS 1.3.</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-https-certificate-chain",
		Title:       "Certificate Chain of Trust",
		Description: "Root CA → Intermediate CA → Leaf certificate — how browsers verify server identity",
		ContentFile: "fundamentals/networking/https",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
    <div class="d-number"><div class="d-number-value">~200</div><div class="d-number-label">Trusted Root CAs in browsers</div></div>
    <div class="d-number"><div class="d-number-value">90 days</div><div class="d-number-label">Let's Encrypt cert lifetime</div></div>
    <div class="d-number"><div class="d-number-value">2048-bit</div><div class="d-number-label">Minimum RSA key size</div></div>
  </div>
  <div class="d-box purple" data-tip="~200 Root CAs are pre-installed in OS/browser trust stores. They never issue leaf certs directly — they sign Intermediate CAs offline and keep their private keys in air-gapped HSMs.">Root CA <span class="d-tag purple">self-signed · stored in browser/OS</span><br><small>DigiCert, Let's Encrypt ISRG Root, Comodo, GlobalSign</small></div>
  <div class="d-arrow-down">&#8595; signs (offline, air-gapped)</div>
  <div class="d-box indigo" data-tip="Intermediate CAs are signed by Root CAs and used online to issue leaf certificates. If compromised, only the Intermediate can be revoked without touching the Root.">Intermediate CA <span class="d-tag indigo">online · revocable</span><br><small>Let's Encrypt R3, DigiCert TLS RSA SHA256</small></div>
  <div class="d-arrow-down">&#8595; issues (ACME / manual)</div>
  <div class="d-box green" data-tip="The leaf certificate is what your server presents. Contains: Subject (domain), Public Key, Validity period, SAN (Subject Alternative Names), Signature by Intermediate CA.">Leaf Certificate <span class="d-tag green">your server's cert</span><br><small>CN=api.example.com · Valid: 90 days · SAN: *.example.com</small></div>
  <div class="d-arrow-down">&#8595; presented during TLS handshake</div>
  <div class="d-box blue" data-tip="Browser verifies: (1) cert signature is valid (signed by trusted Intermediate), (2) Intermediate cert is signed by trusted Root, (3) hostname matches CN/SAN, (4) cert not expired, (5) cert not revoked (OCSP/CRL).">Browser validates chain<br><small>signature → issuer → hostname → expiry → revocation</small></div>
  <div class="d-caption">Full chain must be served: leaf + intermediates. Missing intermediates cause 'UNABLE_TO_VERIFY_LEAF_SIGNATURE' errors in some clients even if the Root is trusted.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-https-mtls",
		Title:       "mTLS: Mutual Authentication",
		Description: "Standard TLS vs mTLS — both client and server present certificates; zero-trust microservices pattern",
		ContentFile: "fundamentals/networking/https",
		Type:        TypeHTML,
		HTML: `<div class="d-cols" style="grid-template-columns: 1fr 1fr; gap: 2rem;">
  <div class="d-flow-v">
    <div class="d-group-title">Standard TLS (one-way)</div>
    <div class="d-box blue">Client <small>(browser, mobile app)</small></div>
    <div class="d-arrow-down">→ ClientHello</div>
    <div class="d-box green">Server presents certificate</div>
    <div class="d-arrow-down">← Client verifies server identity</div>
    <div class="d-box amber">Session established<br><span class="d-tag green">Server authenticated ✓</span><br><span class="d-tag red">Client anonymous ✗</span></div>
    <div class="d-caption">Used for: public HTTPS. Client is not authenticated at TLS layer — use JWT/OAuth at application layer instead.</div>
  </div>
  <div class="d-flow-v">
    <div class="d-group-title">mTLS (zero-trust, both sides)</div>
    <div class="d-box blue">Service A <small>(e.g., payment-svc)</small><br><span class="d-tag blue">has client cert from internal CA</span></div>
    <div class="d-arrow-down">→ ClientHello + client cert</div>
    <div class="d-box green">Service B <small>(e.g., fraud-svc)</small><br><span class="d-tag green">verifies client cert</span></div>
    <div class="d-arrow-down">← Server cert + CertificateRequest</div>
    <div class="d-box amber">Session established<br><span class="d-tag green">Server authenticated ✓</span><br><span class="d-tag green">Client authenticated ✓</span></div>
    <div class="d-caption">Used for: microservices (Istio, Linkerd auto-provision), IoT devices, zero-trust internal APIs. Short-lived certs (24h) — Istio rotates automatically.</div>
  </div>
</div>`,
	})

	// -------------------------------------------------------
	// gRPC
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-grpc-overview",
		Title:       "gRPC Architecture",
		Description: "Client stub → HTTP/2 streams → server handler: how gRPC multiplexes calls over one connection",
		ContentFile: "fundamentals/networking/grpc",
		Type:        TypeHTML,
		HTML: `<div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
  <div class="d-number"><div class="d-number-value">10x</div><div class="d-number-label">Throughput vs REST (large payloads)</div></div>
  <div class="d-number"><div class="d-number-value">3x</div><div class="d-number-label">Smaller payload vs JSON</div></div>
  <div class="d-number"><div class="d-number-value">~20</div><div class="d-number-label">Languages with generated stubs</div></div>
</div>
<div class="d-flow-v">
  <div class="d-cols" style="grid-template-columns: 1fr 1fr; gap: 2rem;">
    <div class="d-flow-v">
      <div class="d-group-title">Client side</div>
      <div class="d-box blue" data-tip="Your application code calls methods on a generated stub as if they were local functions. The stub handles serialization, framing, and HTTP/2 transport transparently.">Application code<br><small>calls stub methods</small></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box indigo" data-tip="Auto-generated from .proto file. Serializes arguments to Protobuf binary. Creates HTTP/2 stream. Handles retries and deadlines.">Generated Stub<br><small>protobuf serialize → HTTP/2 stream</small></div>
    </div>
    <div class="d-flow-v">
      <div class="d-group-title">Server side</div>
      <div class="d-box green" data-tip="Your handler implements the interface generated from .proto. Receives deserialized Protobuf objects. Returns typed response.">Service Handler<br><small>implements .proto interface</small></div>
      <div class="d-arrow-down">&#8593;</div>
      <div class="d-box indigo" data-tip="Deserializes Protobuf binary from HTTP/2 stream. Dispatches to correct handler. Sends response as Protobuf on same stream.">gRPC Server Runtime<br><small>HTTP/2 stream → deserialize → dispatch</small></div>
    </div>
  </div>
  <div class="d-arrow-down">single TCP connection → multiple concurrent streams (HTTP/2 multiplexing)</div>
  <div class="d-box amber" data-tip="HTTP/2 multiplexes many gRPC calls over a single TCP connection. No head-of-line blocking per stream. Header compression (HPACK) reduces overhead. Bidirectional streaming native.">HTTP/2 Connection<br><span class="d-tag amber">stream 1: GetUser</span> <span class="d-tag blue">stream 3: CreateOrder</span> <span class="d-tag green">stream 5: StreamFeed</span><br><small>all multiplexed — no connection overhead per call</small></div>
  <div class="d-caption">One TCP connection handles dozens of concurrent gRPC calls. REST requires connection-per-request or HTTP keep-alive pool — gRPC's HTTP/2 base is more efficient at scale.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-grpc-streaming-modes",
		Title:       "gRPC 4 Streaming Modes",
		Description: "Unary, server streaming, client streaming, bidirectional streaming — when to use each",
		ContentFile: "fundamentals/networking/grpc",
		Type:        TypeHTML,
		HTML: `<div class="d-cols" style="grid-template-columns: 1fr 1fr; gap: 1.5rem;">
  <div class="d-flow-v">
    <div class="d-group-title">1. Unary RPC</div>
    <div class="d-box blue">Client → Request</div>
    <div class="d-arrow-down">&#8595;</div>
    <div class="d-box green">Server → Response</div>
    <div class="d-caption">Same as REST. Use for: CRUD operations, auth, simple lookups. Most common mode.</div>
  </div>
  <div class="d-flow-v">
    <div class="d-group-title">2. Server Streaming</div>
    <div class="d-box blue">Client → Request (once)</div>
    <div class="d-arrow-down">&#8595;</div>
    <div class="d-box green">Server → stream of responses<br><small>msg 1 → msg 2 → msg 3 → done</small></div>
    <div class="d-caption">Use for: live stock prices, log streaming, large dataset pagination, ML model inference results.</div>
  </div>
  <div class="d-flow-v">
    <div class="d-group-title">3. Client Streaming</div>
    <div class="d-box blue">Client → stream of requests<br><small>chunk 1 → chunk 2 → done</small></div>
    <div class="d-arrow-down">&#8595;</div>
    <div class="d-box green">Server → Response (once)</div>
    <div class="d-caption">Use for: file upload in chunks, bulk data ingestion, IoT sensor readings aggregation.</div>
  </div>
  <div class="d-flow-v">
    <div class="d-group-title">4. Bidirectional Streaming</div>
    <div class="d-box blue">Client → stream</div>
    <div class="d-arrow" style="align-self: center;">⇄ simultaneously</div>
    <div class="d-box green">Server → stream</div>
    <div class="d-caption">Use for: real-time chat between services, collaborative editing sync, game state updates, trading order book.</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-grpc-vs-rest",
		Title:       "gRPC vs REST Decision Guide",
		Description: "When to choose gRPC vs REST for service communication — performance, compatibility, use case matrix",
		ContentFile: "fundamentals/networking/grpc",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow" style="gap: 1.5rem; margin-bottom: 0.75rem; flex-wrap: wrap;">
    <div class="d-number"><div class="d-number-value">25-30%</div><div class="d-number-label">gRPC throughput gain (medium load)</div></div>
    <div class="d-number"><div class="d-number-value">10x</div><div class="d-number-label">gRPC throughput (large payloads)</div></div>
    <div class="d-number"><div class="d-number-value">~1/3</div><div class="d-number-label">Protobuf payload size vs JSON</div></div>
  </div>
  <div class="d-cols" style="grid-template-columns: 1fr 1fr; gap: 2rem; margin-top: 1rem;">
    <div class="d-group">
      <div class="d-group-title">Use gRPC when</div>
      <div class="d-flow-v" style="gap: 0.5rem;">
        <div class="d-box green" data-tip="gRPC generates type-safe stubs in 20+ languages from .proto files. One schema definition serves Go, Java, Python, Node.js simultaneously."><strong>Internal microservice-to-microservice</strong><br><small>type-safe contracts, multi-language stubs</small></div>
        <div class="d-box green" data-tip="Server streaming is first-class in gRPC. REST requires SSE or WebSocket for streaming — gRPC streaming is built into the protocol."><strong>Streaming data</strong><br><small>4 streaming modes native to protocol</small></div>
        <div class="d-box green" data-tip="Protobuf binary is ~3x smaller than JSON. HTTP/2 multiplexing means less latency per call. Critical for high-frequency calls between services."><strong>High-throughput, low-latency paths</strong><br><small>payment processing, fraud detection, real-time feeds</small></div>
        <div class="d-box green" data-tip="Google, Netflix, Square, CoreOS, Docker use gRPC for internal service communication."><strong>Polyglot microservices</strong><br><small>Go services calling Python ML models, Java backends</small></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Use REST when</div>
      <div class="d-flow-v" style="gap: 0.5rem;">
        <div class="d-box amber" data-tip="Browsers can't call gRPC directly — need grpc-web proxy (Envoy). REST works natively in every browser with fetch() or XMLHttpRequest."><strong>Public APIs / browser clients</strong><br><small>third-party developers, mobile SDKs, web frontends</small></div>
        <div class="d-box amber" data-tip="curl, Postman, HTTPie all work with REST out of the box. gRPC requires grpcurl or generated clients."><strong>Human-readable APIs</strong><br><small>debugging, partner integrations, webhooks</small></div>
        <div class="d-box amber" data-tip="REST responses can be cached by CDN, browser, or proxy. gRPC POST requests are not HTTP-cacheable."><strong>Cacheable responses</strong><br><small>CDN-cacheable GET endpoints, read-heavy public data</small></div>
        <div class="d-box amber" data-tip="Teams familiar with REST don't need to learn .proto files, code generation, and gRPC ecosystem. REST has lower initial overhead."><strong>Simple CRUD, team REST expertise</strong><br><small>no build step, no protoc toolchain</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	// -------------------------------------------------------
	// GraphQL
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "fund-graphql-overview",
		Title:       "GraphQL Architecture",
		Description: "Single endpoint, schema, resolvers, and how a client query maps to resolver execution",
		ContentFile: "fundamentals/networking/graphql",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols" style="grid-template-columns: 1fr auto 1fr; gap: 1rem; align-items: start;">
    <div class="d-flow-v">
      <div class="d-group-title">Client Query</div>
      <div class="d-box blue" style="font-family: monospace; font-size: 0.8rem; white-space: pre; text-align: left;" data-tip="Client specifies exactly which fields it needs. No over-fetching — only requested fields are returned. Nested queries traverse the schema graph.">query {
  user(id: "123") {
    name
    email
    posts(first: 5) {
      title
      likes
    }
  }
}</div>
    </div>
    <div class="d-arrow" style="align-self: center;">→</div>
    <div class="d-flow-v">
      <div class="d-group-title">GraphQL Server</div>
      <div class="d-box green" data-tip="GraphQL parses the query, validates against the schema, then executes each field resolver in parallel where possible.">Parse → Validate → Execute</div>
      <div class="d-arrow-down">&#8595; resolver chain</div>
      <div class="d-box amber" data-tip="userResolver: SELECT * FROM users WHERE id = 123 (1 query). postsResolver: DataLoader batches all post lookups into one query.">userResolver → postsResolver<br><small>DataLoader batches N+1 → 1 query</small></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box indigo" data-tip="Response contains exactly the requested fields — nothing more, nothing less. Nested structure mirrors query structure.">Response: exactly what was asked<br><small>{ user: { name, email, posts: [...] } }</small></div>
    </div>
  </div>
  <div class="d-arrow-down">single endpoint: POST /graphql</div>
  <div class="d-box purple" data-tip="GraphQL uses a single endpoint for all operations. The query body determines what data is returned. Unlike REST's N endpoints, GraphQL has 1 endpoint with a query language.">POST /graphql <span class="d-tag purple">one endpoint for all operations</span><br><small>Content-Type: application/json · body: { query, variables, operationName }</small></div>
  <div class="d-caption">GraphQL eliminates over-fetching (REST returns full objects) and under-fetching (REST requires N requests for N related resources). One query = one network round trip.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-graphql-n1-problem",
		Title:       "N+1 Problem & DataLoader Solution",
		Description: "Without DataLoader: 1 query for posts + N queries for authors. With DataLoader: 2 queries total",
		ContentFile: "fundamentals/networking/graphql",
		Type:        TypeHTML,
		HTML: `<div class="d-cols" style="grid-template-columns: 1fr 1fr; gap: 2rem;">
  <div class="d-flow-v">
    <div class="d-group-title" style="color: var(--red);">Without DataLoader — N+1 queries</div>
    <div class="d-box gray">Client requests: 100 posts + their authors</div>
    <div class="d-arrow-down">&#8595;</div>
    <div class="d-box red" data-tip="postsResolver fetches all 100 posts in 1 query. Then authorResolver is called once per post — 100 separate SELECT queries.">SELECT * FROM posts LIMIT 100<br><span class="d-tag red">1 query</span></div>
    <div class="d-arrow-down">&#8595; for each of 100 posts</div>
    <div class="d-box red" data-tip="Each post triggers its own author lookup. 100 posts = 100 separate DB round trips. At 1ms each = 100ms latency just for author lookups.">SELECT * FROM users WHERE id = 1<br>SELECT * FROM users WHERE id = 2<br>...<br>SELECT * FROM users WHERE id = 100<br><span class="d-tag red">100 queries</span></div>
    <div class="d-caption" style="color: var(--red);">Total: 101 DB queries. At 1ms each = 101ms just for DB. N becomes unbounded as client queries deeper.</div>
  </div>
  <div class="d-flow-v">
    <div class="d-group-title" style="color: var(--green);">With DataLoader — 2 queries total</div>
    <div class="d-box gray">Client requests: 100 posts + their authors</div>
    <div class="d-arrow-down">&#8595;</div>
    <div class="d-box green" data-tip="postsResolver fetches all 100 posts in 1 query. Each post's authorResolver calls dataloader.load(author_id) — not a DB query yet.">SELECT * FROM posts LIMIT 100<br><span class="d-tag green">1 query</span></div>
    <div class="d-arrow-down">&#8595; DataLoader collects all author_ids</div>
    <div class="d-box green" data-tip="DataLoader batches all 100 .load(id) calls into a single IN query. Results cached for the request lifetime — same author fetched multiple times = 1 DB hit.">SELECT * FROM users<br>WHERE id IN (1, 2, ..., 100)<br><span class="d-tag green">1 batched query</span></div>
    <div class="d-caption" style="color: var(--green);">Total: 2 DB queries regardless of result size. Per-request cache means duplicate author IDs don't re-query. This is the standard fix — use DataLoader in every resolver.</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fund-graphql-schema",
		Title:       "GraphQL Schema: Operations & Types",
		Description: "Query, Mutation, Subscription operations and how the type system enables compile-time validation",
		ContentFile: "fundamentals/networking/graphql",
		Type:        TypeHTML,
		HTML: `<div class="d-cols" style="grid-template-columns: 1fr 1fr; gap: 2rem;">
  <div class="d-flow-v">
    <div class="d-group-title">Three Operation Types</div>
    <div class="d-box blue" data-tip="Query is read-only. Resolvers run in parallel. Multiple queries can be batched in one network request. Idempotent — safe to retry.">Query (read)<br><small>query { user(id: 123) { name } }</small></div>
    <div class="d-box green" data-tip="Mutation is write. Resolvers run sequentially (to preserve order). Each mutation gets its own transaction context.">Mutation (write)<br><small>mutation { createPost(input: {...}) { id } }</small></div>
    <div class="d-box purple" data-tip="Subscription uses WebSocket (or SSE). Client subscribes to events; server pushes updates. Requires persistent connection — different infra than Query/Mutation.">Subscription (real-time)<br><small>subscription { onNewMessage { text sender } }</small></div>
    <div class="d-caption">Subscriptions require WebSocket infrastructure — separate from HTTP query/mutation servers. Scale differently.</div>
  </div>
  <div class="d-flow-v">
    <div class="d-group-title">Type System</div>
    <div class="d-box indigo" style="font-family: monospace; font-size: 0.78rem; white-space: pre; text-align: left;" data-tip="Every field has an explicit type. ! means non-nullable. [Post!]! means a non-nullable list of non-nullable Posts. Type system enforced at parse time — invalid queries rejected before execution.">type User {
  id: ID!
  name: String!
  email: String!
  posts: [Post!]!
  createdAt: DateTime!
}

type Post {
  id: ID!
  title: String!
  author: User!
  likes: Int!
}</div>
    <div class="d-caption">Strong type system: invalid queries are rejected at parse time (before any DB call). Client tooling (introspection) auto-generates TypeScript types.</div>
  </div>
</div>`,
	})
}

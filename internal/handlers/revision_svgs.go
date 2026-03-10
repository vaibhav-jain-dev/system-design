package handlers

// revisionSVGs maps problem slugs to hand-crafted SVG revision diagrams.
// Eye-friendly palette: warm white bg, soft pastel boxes, clear connectors.
//
// Color scheme:
//   Blue boxes (#e8f0fe / #4285f4)   = Architecture decisions
//   Green boxes (#e6f4ea / #34a853)  = Chosen technology / approach
//   Red boxes (#fce8e6 / #ea4335)    = Rejected alternative
//   Amber boxes (#fef7e0 / #f9ab00)  = Key numbers / metrics
//   Purple boxes (#f3e8fd / #9334e6) = Data stores
//   Teal boxes (#e0f7fa / #00897b)   = Caching / CDN layer
//   Gray lines (#9aa0a6)             = Connectors / arrows
//   Dark text (#202124)              = Primary text
//   Gray text (#5f6368)              = Labels on connectors

var revisionSVGs = map[string]string{

// ─────────────────────────────────────────────────────────────────────
// URL SHORTENER
// ─────────────────────────────────────────────────────────────────────
"url-shortener": `<svg viewBox="0 0 1100 720" xmlns="http://www.w3.org/2000/svg" font-family="Inter,system-ui,sans-serif">
  <defs>
    <marker id="us-arrow" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto"><path d="M0,0 L8,3 L0,6" fill="#9aa0a6"/></marker>
    <marker id="us-arrow-grn" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto"><path d="M0,0 L8,3 L0,6" fill="#34a853"/></marker>
    <marker id="us-arrow-red" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto"><path d="M0,0 L8,3 L0,6" fill="#ea4335"/></marker>
    <filter id="us-shadow"><feDropShadow dx="0" dy="1" stdDeviation="2" flood-opacity="0.08"/></filter>
  </defs>
  <rect width="1100" height="720" rx="12" fill="#fafbfd"/>

  <!-- Title -->
  <text x="550" y="32" text-anchor="middle" font-size="15" font-weight="700" fill="#1e293b">URL Shortener — Architecture Decision Map</text>

  <!-- ═══ WRITE PATH (left) ═══ -->
  <text x="280" y="62" text-anchor="middle" font-size="12" font-weight="600" fill="#4285f4" letter-spacing="1">WRITE PATH</text>

  <!-- Client -->
  <rect x="30" y="80" width="100" height="40" rx="6" fill="#f1f3f4" stroke="#dadce0" stroke-width="1" filter="url(#us-shadow)"/>
  <text x="80" y="105" text-anchor="middle" font-size="11" font-weight="500" fill="#202124">Client</text>

  <!-- Arrow -->
  <line x1="130" y1="100" x2="170" y2="100" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#us-arrow)"/>
  <text x="150" y="93" text-anchor="middle" font-size="8" fill="#5f6368">POST /shorten</text>

  <!-- API Gateway -->
  <rect x="172" y="80" width="110" height="40" rx="6" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="227" y="100" text-anchor="middle" font-size="10" font-weight="600" fill="#1a73e8">API Gateway</text>
  <text x="227" y="112" text-anchor="middle" font-size="8" fill="#5f6368">rate limit + auth</text>

  <!-- Arrow -->
  <line x1="282" y1="100" x2="322" y2="100" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#us-arrow)"/>

  <!-- Snowflake ID -->
  <rect x="324" y="76" width="130" height="48" rx="6" fill="#e6f4ea" stroke="#34a853" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="389" y="96" text-anchor="middle" font-size="10" font-weight="600" fill="#137333">Snowflake ID</text>
  <text x="389" y="108" text-anchor="middle" font-size="8" fill="#5f6368">64-bit, sortable, no collision</text>
  <text x="389" y="118" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">10K IDs/sec/node</text>

  <!-- Arrow -->
  <line x1="454" y1="100" x2="494" y2="100" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#us-arrow)"/>
  <text x="474" y="93" text-anchor="middle" font-size="8" fill="#5f6368">encode</text>

  <!-- Base62 -->
  <rect x="496" y="80" width="100" height="40" rx="6" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="546" y="98" text-anchor="middle" font-size="10" font-weight="600" fill="#1a73e8">Base62</text>
  <text x="546" y="110" text-anchor="middle" font-size="8" fill="#5f6368">7 chars = 3.5T</text>

  <!-- Arrow down to DynamoDB -->
  <line x1="546" y1="120" x2="546" y2="155" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#us-arrow)"/>
  <text x="555" y="140" font-size="8" fill="#5f6368">store</text>

  <!-- DynamoDB (shared) -->
  <rect x="476" y="157" width="140" height="48" rx="6" fill="#f3e8fd" stroke="#9334e6" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="546" y="176" text-anchor="middle" font-size="10" font-weight="600" fill="#7627bb">DynamoDB</text>
  <text x="546" y="188" text-anchor="middle" font-size="8" fill="#5f6368">PK=short_code, O(1) lookup</text>
  <text x="546" y="199" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">auto-scales writes</text>

  <!-- ═══ READ PATH (right) ═══ -->
  <text x="850" y="62" text-anchor="middle" font-size="12" font-weight="600" fill="#00897b" letter-spacing="1">READ PATH</text>

  <!-- Client read -->
  <rect x="670" y="80" width="100" height="40" rx="6" fill="#f1f3f4" stroke="#dadce0" stroke-width="1" filter="url(#us-shadow)"/>
  <text x="720" y="100" text-anchor="middle" font-size="11" font-weight="500" fill="#202124">Client</text>
  <text x="720" y="112" text-anchor="middle" font-size="8" fill="#5f6368">GET /abc1234</text>

  <!-- Arrow -->
  <line x1="770" y1="100" x2="820" y2="100" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#us-arrow)"/>

  <!-- CDN -->
  <rect x="822" y="76" width="120" height="48" rx="6" fill="#e0f7fa" stroke="#00897b" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="882" y="96" text-anchor="middle" font-size="10" font-weight="600" fill="#00695c">CloudFront CDN</text>
  <text x="882" y="108" text-anchor="middle" font-size="8" fill="#5f6368">edge cache, TTL 24h</text>
  <text x="882" y="118" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">~5ms latency</text>

  <!-- Arrow down on cache miss -->
  <line x1="882" y1="124" x2="882" y2="155" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#us-arrow)" stroke-dasharray="4,3"/>
  <text x="910" y="142" font-size="8" fill="#ea4335">miss</text>

  <!-- Redis -->
  <rect x="822" y="157" width="120" height="48" rx="6" fill="#e0f7fa" stroke="#00897b" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="882" y="176" text-anchor="middle" font-size="10" font-weight="600" fill="#00695c">Redis Cache</text>
  <text x="882" y="188" text-anchor="middle" font-size="8" fill="#5f6368">top 20% URLs</text>
  <text x="882" y="199" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">80% read traffic</text>

  <!-- Arrow from Redis miss to DynamoDB -->
  <line x1="822" y1="181" x2="616" y2="181" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#us-arrow)" stroke-dasharray="4,3"/>
  <text x="719" y="174" text-anchor="middle" font-size="8" fill="#ea4335">cache miss → DB</text>

  <!-- 302 Redirect result -->
  <rect x="960" y="157" width="110" height="48" rx="6" fill="#e6f4ea" stroke="#34a853" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="1015" y="176" text-anchor="middle" font-size="10" font-weight="600" fill="#137333">302 Redirect</text>
  <text x="1015" y="188" text-anchor="middle" font-size="8" fill="#5f6368">not 301</text>
  <text x="1015" y="199" text-anchor="middle" font-size="8" fill="#5f6368">analytics + mutable</text>

  <!-- Arrow from CDN/Redis hit to 302 -->
  <line x1="942" y1="100" x2="960" y2="170" stroke="#34a853" stroke-width="1.2" marker-end="url(#us-arrow-grn)" stroke-dasharray="4,3"/>
  <text x="965" y="135" font-size="8" fill="#34a853">hit</text>

  <!-- ═══ KEY DECISIONS (middle section) ═══ -->
  <line x1="30" y1="240" x2="1070" y2="240" stroke="#e2e8f0" stroke-width="1"/>
  <text x="550" y="265" text-anchor="middle" font-size="12" font-weight="600" fill="#1e293b" letter-spacing="1">KEY DECISIONS &amp; WHY</text>

  <!-- Decision 1: Why Snowflake not UUID -->
  <rect x="30" y="285" width="240" height="64" rx="8" fill="#e6f4ea" stroke="#34a853" stroke-width="1.5" filter="url(#us-shadow)"/>
  <text x="150" y="305" text-anchor="middle" font-size="10" font-weight="700" fill="#137333">Snowflake ID (chosen)</text>
  <text x="150" y="318" text-anchor="middle" font-size="9" fill="#202124">64-bit → 8-char Base62</text>
  <text x="150" y="330" text-anchor="middle" font-size="9" fill="#202124">time-sortable, zero coordination</text>
  <text x="150" y="342" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">10K IDs/sec per worker</text>

  <!-- Rejected: UUID -->
  <rect x="30" y="360" width="240" height="42" rx="8" fill="#fce8e6" stroke="#ea4335" stroke-width="1.2"/>
  <text x="45" y="378" font-size="9" fill="#c5221f" font-weight="600">✗ UUID</text>
  <text x="45" y="392" font-size="8" fill="#5f6368">128-bit = 22 chars Base64 — too long for URL</text>

  <!-- Connection line: decision 1 → write path -->
  <line x1="150" y1="285" x2="389" y2="124" stroke="#34a853" stroke-width="1" stroke-dasharray="4,3" opacity="0.5"/>

  <!-- Decision 2: Why DynamoDB -->
  <rect x="310" y="285" width="240" height="64" rx="8" fill="#e6f4ea" stroke="#34a853" stroke-width="1.5" filter="url(#us-shadow)"/>
  <text x="430" y="305" text-anchor="middle" font-size="10" font-weight="700" fill="#137333">DynamoDB (chosen)</text>
  <text x="430" y="318" text-anchor="middle" font-size="9" fill="#202124">single-digit ms PK lookup</text>
  <text x="430" y="330" text-anchor="middle" font-size="9" fill="#202124">auto-scales, no sharding needed</text>
  <text x="430" y="342" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">$1.25/M writes</text>

  <!-- Rejected: Postgres at scale -->
  <rect x="310" y="360" width="240" height="42" rx="8" fill="#fce8e6" stroke="#ea4335" stroke-width="1.2"/>
  <text x="325" y="378" font-size="9" fill="#c5221f" font-weight="600">✗ Postgres at scale</text>
  <text x="325" y="392" font-size="8" fill="#5f6368">manual sharding, connection pooling, ops burden</text>

  <!-- Decision 3: Why 302 -->
  <rect x="590" y="285" width="240" height="64" rx="8" fill="#e6f4ea" stroke="#34a853" stroke-width="1.5" filter="url(#us-shadow)"/>
  <text x="710" y="305" text-anchor="middle" font-size="10" font-weight="700" fill="#137333">302 Temporary (chosen)</text>
  <text x="710" y="318" text-anchor="middle" font-size="9" fill="#202124">every request hits server → analytics</text>
  <text x="710" y="330" text-anchor="middle" font-size="9" fill="#202124">target URL can be updated</text>
  <text x="710" y="342" text-anchor="middle" font-size="8" fill="#5f6368">CDN absorbs load anyway</text>

  <!-- Rejected: 301 -->
  <rect x="590" y="360" width="240" height="42" rx="8" fill="#fce8e6" stroke="#ea4335" stroke-width="1.2"/>
  <text x="605" y="378" font-size="9" fill="#c5221f" font-weight="600">✗ 301 Permanent</text>
  <text x="605" y="392" font-size="8" fill="#5f6368">browsers cache forever — no click tracking, no URL change</text>

  <!-- Decision 4: Why Redis -->
  <rect x="870" y="285" width="200" height="64" rx="8" fill="#e6f4ea" stroke="#34a853" stroke-width="1.5" filter="url(#us-shadow)"/>
  <text x="970" y="305" text-anchor="middle" font-size="10" font-weight="700" fill="#137333">Redis Cache (chosen)</text>
  <text x="970" y="318" text-anchor="middle" font-size="9" fill="#202124">80/20 rule: 20% URLs</text>
  <text x="970" y="330" text-anchor="middle" font-size="9" fill="#202124">serve 80% of all reads</text>
  <text x="970" y="342" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">sub-ms, TTL 24h</text>

  <!-- ═══ SCALE EVOLUTION ═══ -->
  <line x1="30" y1="430" x2="1070" y2="430" stroke="#e2e8f0" stroke-width="1"/>
  <text x="550" y="455" text-anchor="middle" font-size="12" font-weight="600" fill="#1e293b" letter-spacing="1">SCALE EVOLUTION</text>

  <!-- Stage 1 -->
  <rect x="30" y="475" width="310" height="80" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="185" y="495" text-anchor="middle" font-size="11" font-weight="700" fill="#1a73e8">Stage 1: 0 → 100K users</text>
  <text x="185" y="510" text-anchor="middle" font-size="9" fill="#202124">Postgres + single app server</text>
  <text x="185" y="523" text-anchor="middle" font-size="9" fill="#202124">No Redis, no CDN needed yet</text>
  <text x="185" y="536" text-anchor="middle" font-size="8" fill="#5f6368">Cost: ~$200/mo</text>
  <text x="185" y="548" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">handles 1K writes/sec easily</text>

  <!-- Arrow -->
  <line x1="340" y1="515" x2="380" y2="515" stroke="#4285f4" stroke-width="2" marker-end="url(#us-arrow)"/>

  <!-- Stage 2 -->
  <rect x="382" y="475" width="310" height="80" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="537" y="495" text-anchor="middle" font-size="11" font-weight="700" fill="#1a73e8">Stage 2: 100K → 10M</text>
  <text x="537" y="510" text-anchor="middle" font-size="9" fill="#202124">DynamoDB + Redis + CloudFront</text>
  <text x="537" y="523" text-anchor="middle" font-size="9" fill="#202124">Read replicas, cache-aside pattern</text>
  <text x="537" y="536" text-anchor="middle" font-size="8" fill="#5f6368">Cost: ~$2,500/mo</text>
  <text x="537" y="548" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">58K reads/sec, 6K writes/sec</text>

  <!-- Arrow -->
  <line x1="692" y1="515" x2="732" y2="515" stroke="#4285f4" stroke-width="2" marker-end="url(#us-arrow)"/>

  <!-- Stage 3 -->
  <rect x="734" y="475" width="310" height="80" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#us-shadow)"/>
  <text x="889" y="495" text-anchor="middle" font-size="11" font-weight="700" fill="#1a73e8">Stage 3: 10M+ global</text>
  <text x="889" y="510" text-anchor="middle" font-size="9" fill="#202124">Multi-region, DynamoDB Global Tables</text>
  <text x="889" y="523" text-anchor="middle" font-size="9" fill="#202124">CRDTs for counters, geo-routing</text>
  <text x="889" y="536" text-anchor="middle" font-size="8" fill="#5f6368">3x infra cost for global coverage</text>
  <text x="889" y="548" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">EU 150ms → 20ms latency</text>

  <!-- ═══ KEY NUMBERS ═══ -->
  <line x1="30" y1="580" x2="1070" y2="580" stroke="#e2e8f0" stroke-width="1"/>
  <text x="550" y="605" text-anchor="middle" font-size="12" font-weight="600" fill="#1e293b" letter-spacing="1">KEY NUMBERS TO REMEMBER</text>

  <rect x="30" y="620" width="150" height="50" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="105" y="640" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">7 chars</text>
  <text x="105" y="655" text-anchor="middle" font-size="8" fill="#5f6368">3.5 trillion combos</text>

  <rect x="200" y="620" width="150" height="50" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="275" y="640" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">500M URLs/day</text>
  <text x="275" y="655" text-anchor="middle" font-size="8" fill="#5f6368">~6K writes/sec</text>

  <rect x="370" y="620" width="150" height="50" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="445" y="640" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">10:1 read:write</text>
  <text x="445" y="655" text-anchor="middle" font-size="8" fill="#5f6368">~58K reads/sec</text>

  <rect x="540" y="620" width="150" height="50" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="615" y="640" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">5ms p99</text>
  <text x="615" y="655" text-anchor="middle" font-size="8" fill="#5f6368">CDN redirect latency</text>

  <rect x="710" y="620" width="150" height="50" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="785" y="640" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">TTL 24h</text>
  <text x="785" y="655" text-anchor="middle" font-size="8" fill="#5f6368">Redis + CDN cache</text>

  <rect x="880" y="620" width="190" height="50" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="975" y="640" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">$1.25/M writes</text>
  <text x="975" y="655" text-anchor="middle" font-size="8" fill="#5f6368">DynamoDB on-demand</text>

  <!-- Bottom key takeaway -->
  <rect x="250" y="690" width="600" height="24" rx="12" fill="#e8f0fe" stroke="#4285f4" stroke-width="1"/>
  <text x="550" y="706" text-anchor="middle" font-size="10" font-weight="600" fill="#1a73e8">Core insight: Snowflake ID eliminates collision checking — the only hard problem becomes caching at scale</text>
</svg>`,

// ─────────────────────────────────────────────────────────────────────
// RATE LIMITER
// ─────────────────────────────────────────────────────────────────────
"rate-limiter": `<svg viewBox="0 0 1100 750" xmlns="http://www.w3.org/2000/svg" font-family="Inter,system-ui,sans-serif">
  <defs>
    <marker id="rl-arr" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto"><path d="M0,0 L8,3 L0,6" fill="#9aa0a6"/></marker>
    <marker id="rl-arr-g" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto"><path d="M0,0 L8,3 L0,6" fill="#34a853"/></marker>
    <marker id="rl-arr-r" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto"><path d="M0,0 L8,3 L0,6" fill="#ea4335"/></marker>
    <filter id="rl-sh"><feDropShadow dx="0" dy="1" stdDeviation="2" flood-opacity="0.08"/></filter>
  </defs>
  <rect width="1100" height="750" rx="12" fill="#fafbfd"/>
  <text x="550" y="32" text-anchor="middle" font-size="15" font-weight="700" fill="#1e293b">Rate Limiter — Architecture Decision Map</text>

  <!-- ═══ REQUEST FLOW ═══ -->
  <text x="550" y="60" text-anchor="middle" font-size="12" font-weight="600" fill="#4285f4" letter-spacing="1">REQUEST FLOW (MIDDLEWARE PATTERN)</text>

  <!-- Client -->
  <rect x="30" y="80" width="90" height="40" rx="6" fill="#f1f3f4" stroke="#dadce0" stroke-width="1" filter="url(#rl-sh)"/>
  <text x="75" y="104" text-anchor="middle" font-size="11" font-weight="500" fill="#202124">Client</text>

  <line x1="120" y1="100" x2="155" y2="100" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#rl-arr)"/>

  <!-- API Gateway -->
  <rect x="157" y="78" width="110" height="44" rx="6" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#rl-sh)"/>
  <text x="212" y="98" text-anchor="middle" font-size="10" font-weight="600" fill="#1a73e8">API Gateway</text>
  <text x="212" y="112" text-anchor="middle" font-size="8" fill="#5f6368">extracts user key</text>

  <line x1="267" y1="100" x2="302" y2="100" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#rl-arr)"/>

  <!-- Rate Limiter Middleware -->
  <rect x="304" y="74" width="160" height="52" rx="6" fill="#e6f4ea" stroke="#34a853" stroke-width="1.5" filter="url(#rl-sh)"/>
  <text x="384" y="94" text-anchor="middle" font-size="10" font-weight="700" fill="#137333">Rate Limiter</text>
  <text x="384" y="106" text-anchor="middle" font-size="9" fill="#202124">Token Bucket / Sliding Window</text>
  <text x="384" y="118" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">check + decrement atomically</text>

  <!-- Allowed path -->
  <line x1="464" y1="90" x2="520" y2="90" stroke="#34a853" stroke-width="1.5" marker-end="url(#rl-arr-g)"/>
  <text x="492" y="83" text-anchor="middle" font-size="8" fill="#34a853" font-weight="600">200 OK</text>

  <rect x="522" y="76" width="100" height="36" rx="6" fill="#e6f4ea" stroke="#34a853" stroke-width="1"/>
  <text x="572" y="98" text-anchor="middle" font-size="10" font-weight="500" fill="#137333">App Server</text>

  <!-- Rejected path -->
  <line x1="464" y1="110" x2="520" y2="130" stroke="#ea4335" stroke-width="1.5" marker-end="url(#rl-arr-r)"/>
  <text x="500" y="128" font-size="8" fill="#ea4335" font-weight="600">429</text>

  <rect x="522" y="116" width="100" height="36" rx="6" fill="#fce8e6" stroke="#ea4335" stroke-width="1"/>
  <text x="572" y="131" text-anchor="middle" font-size="9" font-weight="500" fill="#c5221f">Too Many</text>
  <text x="572" y="144" text-anchor="middle" font-size="8" fill="#5f6368">Retry-After: N</text>

  <!-- Redis connection from middleware -->
  <line x1="384" y1="126" x2="384" y2="165" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#rl-arr)"/>
  <text x="395" y="148" font-size="8" fill="#5f6368">Lua script</text>

  <!-- Redis -->
  <rect x="314" y="167" width="140" height="48" rx="6" fill="#e0f7fa" stroke="#00897b" stroke-width="1.2" filter="url(#rl-sh)"/>
  <text x="384" y="186" text-anchor="middle" font-size="10" font-weight="600" fill="#00695c">Redis Cluster</text>
  <text x="384" y="198" text-anchor="middle" font-size="8" fill="#5f6368">sorted set per user key</text>
  <text x="384" y="209" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">sub-ms, atomic via Lua</text>

  <!-- Fail-open decision -->
  <rect x="670" y="76" width="200" height="52" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2" filter="url(#rl-sh)"/>
  <text x="770" y="94" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">FAIL-OPEN (critical)</text>
  <text x="770" y="107" text-anchor="middle" font-size="9" fill="#202124">Redis down → allow all traffic</text>
  <text x="770" y="120" text-anchor="middle" font-size="8" fill="#5f6368">outage > brief unlimited traffic</text>

  <line x1="622" y1="94" x2="668" y2="94" stroke="#f9ab00" stroke-width="1.2" stroke-dasharray="4,3"/>

  <!-- Local fallback -->
  <rect x="890" y="76" width="180" height="52" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1"/>
  <text x="980" y="94" text-anchor="middle" font-size="9" font-weight="600" fill="#1a73e8">Local In-Memory Fallback</text>
  <text x="980" y="107" text-anchor="middle" font-size="8" fill="#202124">per-node token bucket</text>
  <text x="980" y="120" text-anchor="middle" font-size="8" fill="#5f6368">~80% accuracy, no coordination</text>

  <line x1="870" y1="102" x2="888" y2="102" stroke="#4285f4" stroke-width="1" stroke-dasharray="4,3"/>

  <!-- ═══ ALGORITHM COMPARISON ═══ -->
  <line x1="30" y1="240" x2="1070" y2="240" stroke="#e2e8f0" stroke-width="1"/>
  <text x="550" y="265" text-anchor="middle" font-size="12" font-weight="600" fill="#1e293b" letter-spacing="1">ALGORITHM COMPARISON</text>

  <!-- Token Bucket -->
  <rect x="30" y="280" width="240" height="90" rx="8" fill="#e6f4ea" stroke="#34a853" stroke-width="1.5" filter="url(#rl-sh)"/>
  <text x="150" y="300" text-anchor="middle" font-size="10" font-weight="700" fill="#137333">Token Bucket (best for APIs)</text>
  <text x="150" y="315" text-anchor="middle" font-size="9" fill="#202124">refills at constant rate</text>
  <text x="150" y="328" text-anchor="middle" font-size="9" fill="#202124">allows bursts up to bucket size</text>
  <text x="150" y="341" text-anchor="middle" font-size="9" fill="#202124">2 Redis keys: tokens + timestamp</text>
  <text x="150" y="356" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">O(1) per request</text>

  <!-- Sliding Window Log -->
  <rect x="290" y="280" width="240" height="90" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#rl-sh)"/>
  <text x="410" y="300" text-anchor="middle" font-size="10" font-weight="700" fill="#1a73e8">Sliding Window Log</text>
  <text x="410" y="315" text-anchor="middle" font-size="9" fill="#202124">sorted set of timestamps</text>
  <text x="410" y="328" text-anchor="middle" font-size="9" fill="#202124">exact count, no approximation</text>
  <text x="410" y="341" text-anchor="middle" font-size="9" fill="#202124">ZREMRANGEBYSCORE + ZCARD + ZADD</text>
  <text x="410" y="356" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">O(log N) — Lua for atomicity</text>

  <!-- Fixed Window -->
  <rect x="550" y="280" width="240" height="90" rx="8" fill="#f1f3f4" stroke="#dadce0" stroke-width="1.2"/>
  <text x="670" y="300" text-anchor="middle" font-size="10" font-weight="700" fill="#5f6368">Fixed Window Counter</text>
  <text x="670" y="315" text-anchor="middle" font-size="9" fill="#202124">one counter per time window</text>
  <text x="670" y="328" text-anchor="middle" font-size="9" fill="#202124">simple but boundary burst problem</text>
  <text x="670" y="341" text-anchor="middle" font-size="9" fill="#ea4335">2x burst at window boundary</text>
  <text x="670" y="356" text-anchor="middle" font-size="8" fill="#5f6368">O(1) — INCR + EXPIRE</text>

  <!-- Leaky Bucket -->
  <rect x="810" y="280" width="240" height="90" rx="8" fill="#f1f3f4" stroke="#dadce0" stroke-width="1.2"/>
  <text x="930" y="300" text-anchor="middle" font-size="10" font-weight="700" fill="#5f6368">Leaky Bucket</text>
  <text x="930" y="315" text-anchor="middle" font-size="9" fill="#202124">constant drain rate</text>
  <text x="930" y="328" text-anchor="middle" font-size="9" fill="#202124">smooths output, no bursts</text>
  <text x="930" y="341" text-anchor="middle" font-size="9" fill="#ea4335">can't handle legitimate spikes</text>
  <text x="930" y="356" text-anchor="middle" font-size="8" fill="#5f6368">O(1) — queue-based</text>

  <!-- ═══ DISTRIBUTED CHALLENGES ═══ -->
  <line x1="30" y1="395" x2="1070" y2="395" stroke="#e2e8f0" stroke-width="1"/>
  <text x="550" y="420" text-anchor="middle" font-size="12" font-weight="600" fill="#1e293b" letter-spacing="1">DISTRIBUTED CHALLENGES</text>

  <!-- Race condition -->
  <rect x="30" y="440" width="330" height="70" rx="8" fill="#fce8e6" stroke="#ea4335" stroke-width="1.2" filter="url(#rl-sh)"/>
  <text x="195" y="460" text-anchor="middle" font-size="10" font-weight="700" fill="#c5221f">Race Condition: GET-then-SET</text>
  <text x="195" y="475" text-anchor="middle" font-size="9" fill="#202124">Two requests read same count → both decrement</text>
  <text x="195" y="488" text-anchor="middle" font-size="9" fill="#202124">→ limit exceeded by up to N concurrent requests</text>
  <text x="195" y="501" text-anchor="middle" font-size="8" fill="#5f6368">At 500K req/s, this happens constantly</text>

  <!-- Arrow to solution -->
  <line x1="360" y1="475" x2="400" y2="475" stroke="#34a853" stroke-width="1.5" marker-end="url(#rl-arr-g)"/>
  <text x="380" y="468" text-anchor="middle" font-size="8" fill="#34a853">fix</text>

  <!-- Solution -->
  <rect x="402" y="440" width="330" height="70" rx="8" fill="#e6f4ea" stroke="#34a853" stroke-width="1.5" filter="url(#rl-sh)"/>
  <text x="567" y="460" text-anchor="middle" font-size="10" font-weight="700" fill="#137333">Solution: Lua Script Atomicity</text>
  <text x="567" y="475" text-anchor="middle" font-size="9" fill="#202124">Redis executes Lua as single atomic operation</text>
  <text x="567" y="488" text-anchor="middle" font-size="9" fill="#202124">ZREMRANGEBYSCORE + ZCARD + ZADD in one call</text>
  <text x="567" y="501" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">0.1ms instead of 3 × 0.1ms round trips</text>

  <!-- Multi-region sync -->
  <rect x="752" y="440" width="320" height="70" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#rl-sh)"/>
  <text x="912" y="460" text-anchor="middle" font-size="10" font-weight="700" fill="#1a73e8">Multi-Region: Eventual Sync</text>
  <text x="912" y="475" text-anchor="middle" font-size="9" fill="#202124">each region has local Redis counter</text>
  <text x="912" y="488" text-anchor="middle" font-size="9" fill="#202124">sync every 5s → global counter ±5%</text>
  <text x="912" y="501" text-anchor="middle" font-size="8" fill="#5f6368">accuracy trade-off for latency</text>

  <!-- ═══ SCALE EVOLUTION ═══ -->
  <line x1="30" y1="535" x2="1070" y2="535" stroke="#e2e8f0" stroke-width="1"/>
  <text x="550" y="560" text-anchor="middle" font-size="12" font-weight="600" fill="#1e293b" letter-spacing="1">SCALE EVOLUTION</text>

  <rect x="30" y="575" width="330" height="65" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#rl-sh)"/>
  <text x="195" y="595" text-anchor="middle" font-size="11" font-weight="700" fill="#1a73e8">Stage 1: Single Service</text>
  <text x="195" y="610" text-anchor="middle" font-size="9" fill="#202124">in-memory hash map, single node</text>
  <text x="195" y="625" text-anchor="middle" font-size="8" fill="#5f6368">works up to ~10K req/sec on one box</text>

  <line x1="360" y1="607" x2="400" y2="607" stroke="#4285f4" stroke-width="2" marker-end="url(#rl-arr)"/>

  <rect x="402" y="575" width="330" height="65" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#rl-sh)"/>
  <text x="567" y="595" text-anchor="middle" font-size="11" font-weight="700" fill="#1a73e8">Stage 2: Centralized Redis</text>
  <text x="567" y="610" text-anchor="middle" font-size="9" fill="#202124">Redis Cluster, Lua scripts, fail-open</text>
  <text x="567" y="625" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">handles 500K+ req/sec</text>

  <line x1="732" y1="607" x2="772" y2="607" stroke="#4285f4" stroke-width="2" marker-end="url(#rl-arr)"/>

  <rect x="774" y="575" width="296" height="65" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#rl-sh)"/>
  <text x="922" y="595" text-anchor="middle" font-size="11" font-weight="700" fill="#1a73e8">Stage 3: Multi-Region</text>
  <text x="922" y="610" text-anchor="middle" font-size="9" fill="#202124">local Redis per region, async sync</text>
  <text x="922" y="625" text-anchor="middle" font-size="8" fill="#5f6368">±5% accuracy, sub-ms local latency</text>

  <!-- ═══ KEY NUMBERS ═══ -->
  <line x1="30" y1="660" x2="1070" y2="660" stroke="#e2e8f0" stroke-width="1"/>

  <rect x="30" y="675" width="160" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="110" y="694" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">500K req/sec</text>
  <text x="110" y="708" text-anchor="middle" font-size="8" fill="#5f6368">peak rate limit checks</text>

  <rect x="210" y="675" width="160" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="290" y="694" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">0.1ms</text>
  <text x="290" y="708" text-anchor="middle" font-size="8" fill="#5f6368">Lua script latency</text>

  <rect x="390" y="675" width="160" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="470" y="694" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">100 rules/user</text>
  <text x="470" y="708" text-anchor="middle" font-size="8" fill="#5f6368">IP + API key + endpoint</text>

  <rect x="570" y="675" width="160" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="650" y="694" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">429 + headers</text>
  <text x="650" y="708" text-anchor="middle" font-size="8" fill="#5f6368">Retry-After, X-RateLimit-*</text>

  <rect x="750" y="675" width="160" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="830" y="694" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">fail-open</text>
  <text x="830" y="708" text-anchor="middle" font-size="8" fill="#5f6368">Redis down = allow all</text>

  <rect x="930" y="675" width="140" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="1000" y="694" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">~$92/mo</text>
  <text x="1000" y="708" text-anchor="middle" font-size="8" fill="#5f6368">ElastiCache r6g.large</text>

  <!-- Bottom takeaway -->
  <rect x="200" y="730" width="700" height="18" rx="9" fill="#e8f0fe" stroke="#4285f4" stroke-width="1"/>
  <text x="550" y="743" text-anchor="middle" font-size="9" font-weight="600" fill="#1a73e8">Core insight: The hard problem isn't the algorithm — it's atomicity under concurrency and fail-open resilience</text>
</svg>`,

// ─────────────────────────────────────────────────────────────────────
// INSTAGRAM
// ─────────────────────────────────────────────────────────────────────
"instagram": `<svg viewBox="0 0 1100 820" xmlns="http://www.w3.org/2000/svg" font-family="Inter,system-ui,sans-serif">
  <defs>
    <marker id="ig-arr" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto"><path d="M0,0 L8,3 L0,6" fill="#9aa0a6"/></marker>
    <marker id="ig-arr-g" markerWidth="8" markerHeight="6" refX="8" refY="3" orient="auto"><path d="M0,0 L8,3 L0,6" fill="#34a853"/></marker>
    <filter id="ig-sh"><feDropShadow dx="0" dy="1" stdDeviation="2" flood-opacity="0.08"/></filter>
  </defs>
  <rect width="1100" height="820" rx="12" fill="#fafbfd"/>
  <text x="550" y="32" text-anchor="middle" font-size="15" font-weight="700" fill="#1e293b">Instagram — Architecture Decision Map</text>

  <!-- ═══ THE HARD PROBLEM ═══ -->
  <rect x="350" y="50" width="400" height="44" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.5" filter="url(#ig-sh)"/>
  <text x="550" y="70" text-anchor="middle" font-size="11" font-weight="700" fill="#e37400">THE HARD PROBLEM: News Feed at Scale</text>
  <text x="550" y="85" text-anchor="middle" font-size="9" fill="#202124">500M DAU × 10 feed loads/day × 200 followees = 1T queries/day if pull-based</text>

  <!-- ═══ HYBRID FAN-OUT (the core decision) ═══ -->
  <text x="550" y="115" text-anchor="middle" font-size="12" font-weight="600" fill="#4285f4" letter-spacing="1">CORE DECISION: HYBRID FAN-OUT</text>

  <!-- Fan-out-on-write for normal -->
  <rect x="30" y="130" width="320" height="80" rx="8" fill="#e6f4ea" stroke="#34a853" stroke-width="1.5" filter="url(#ig-sh)"/>
  <text x="190" y="150" text-anchor="middle" font-size="10" font-weight="700" fill="#137333">Fan-out-on-WRITE (normal users)</text>
  <text x="190" y="165" text-anchor="middle" font-size="9" fill="#202124">post → push to all follower feed caches</text>
  <text x="190" y="178" text-anchor="middle" font-size="9" fill="#202124">Redis ZADD to each follower's sorted set</text>
  <text x="190" y="191" text-anchor="middle" font-size="8" fill="#5f6368">1 post × 1K followers = 1,001 writes</text>
  <text x="190" y="203" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">feed read = single Redis ZREVRANGE → O(1)</text>

  <!-- Threshold -->
  <rect x="380" y="145" width="140" height="50" rx="20" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.5"/>
  <text x="450" y="166" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">100K followers</text>
  <text x="450" y="182" text-anchor="middle" font-size="9" fill="#5f6368">threshold</text>

  <!-- Lines from threshold to both -->
  <line x1="380" y1="170" x2="350" y2="170" stroke="#f9ab00" stroke-width="1.2"/>
  <line x1="520" y1="170" x2="560" y2="170" stroke="#f9ab00" stroke-width="1.2"/>

  <!-- Fan-out-on-read for celebrities -->
  <rect x="562" y="130" width="320" height="80" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="722" y="150" text-anchor="middle" font-size="10" font-weight="700" fill="#1a73e8">Fan-out-on-READ (celebrities)</text>
  <text x="722" y="165" text-anchor="middle" font-size="9" fill="#202124">celebrity posts stored centrally</text>
  <text x="722" y="178" text-anchor="middle" font-size="9" fill="#202124">merged at read time with cached feed</text>
  <text x="722" y="191" text-anchor="middle" font-size="8" fill="#5f6368">avoids 10M Redis writes per celebrity post</text>
  <text x="722" y="203" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">feed read adds ~15ms for merge</text>

  <!-- Rejected: pure fan-out-on-write -->
  <rect x="910" y="130" width="170" height="80" rx="8" fill="#fce8e6" stroke="#ea4335" stroke-width="1.2"/>
  <text x="995" y="150" text-anchor="middle" font-size="9" font-weight="600" fill="#c5221f">✗ Pure Write</text>
  <text x="995" y="165" text-anchor="middle" font-size="8" fill="#5f6368">10M follower post</text>
  <text x="995" y="178" text-anchor="middle" font-size="8" fill="#5f6368">= 10M cache writes</text>
  <text x="995" y="191" text-anchor="middle" font-size="8" fill="#5f6368">= 500 sec Redis time</text>
  <text x="995" y="203" text-anchor="middle" font-size="8" fill="#ea4335" font-weight="600">unacceptable</text>

  <!-- ═══ UPLOAD FLOW ═══ -->
  <line x1="30" y1="230" x2="1070" y2="230" stroke="#e2e8f0" stroke-width="1"/>
  <text x="280" y="255" text-anchor="middle" font-size="12" font-weight="600" fill="#00897b" letter-spacing="1">UPLOAD FLOW (4 HOPS)</text>

  <rect x="30" y="270" width="100" height="40" rx="6" fill="#f1f3f4" stroke="#dadce0" stroke-width="1" filter="url(#ig-sh)"/>
  <text x="80" y="294" text-anchor="middle" font-size="10" font-weight="500" fill="#202124">Client</text>

  <line x1="130" y1="290" x2="165" y2="290" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#ig-arr)"/>
  <text x="147" y="283" text-anchor="middle" font-size="7" fill="#5f6368">init</text>

  <rect x="167" y="268" width="130" height="44" rx="6" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="232" y="288" text-anchor="middle" font-size="9" font-weight="600" fill="#1a73e8">API → presigned URL</text>
  <text x="232" y="302" text-anchor="middle" font-size="8" fill="#5f6368">max 50MB, 15min expiry</text>

  <line x1="297" y1="290" x2="332" y2="290" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#ig-arr)"/>
  <text x="314" y="283" text-anchor="middle" font-size="7" fill="#5f6368">direct</text>

  <rect x="334" y="268" width="130" height="44" rx="6" fill="#f3e8fd" stroke="#9334e6" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="399" y="288" text-anchor="middle" font-size="9" font-weight="600" fill="#7627bb">S3 Upload</text>
  <text x="399" y="302" text-anchor="middle" font-size="8" fill="#5f6368">bypasses API servers</text>

  <line x1="464" y1="290" x2="499" y2="290" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#ig-arr)"/>
  <text x="481" y="283" text-anchor="middle" font-size="7" fill="#5f6368">event</text>

  <rect x="501" y="268" width="150" height="44" rx="6" fill="#e6f4ea" stroke="#34a853" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="576" y="288" text-anchor="middle" font-size="9" font-weight="600" fill="#137333">Lambda (resize)</text>
  <text x="576" y="302" text-anchor="middle" font-size="8" fill="#5f6368">4 variants + WebP + EXIF strip</text>

  <line x1="651" y1="290" x2="686" y2="290" stroke="#9aa0a6" stroke-width="1.5" marker-end="url(#ig-arr)"/>

  <rect x="688" y="268" width="120" height="44" rx="6" fill="#e0f7fa" stroke="#00897b" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="748" y="288" text-anchor="middle" font-size="9" font-weight="600" fill="#00695c">CloudFront</text>
  <text x="748" y="302" text-anchor="middle" font-size="8" fill="#5f6368">edge delivery</text>

  <!-- Why presigned -->
  <rect x="830" y="268" width="240" height="44" rx="6" fill="#fef7e0" stroke="#f9ab00" stroke-width="1"/>
  <text x="950" y="286" text-anchor="middle" font-size="9" font-weight="600" fill="#e37400">Why presigned URL?</text>
  <text x="950" y="300" text-anchor="middle" font-size="8" fill="#202124">5MB photo = 2-10 sec. Can't block API threads.</text>
  <line x1="830" y1="290" x2="808" y2="290" stroke="#f9ab00" stroke-width="1" stroke-dasharray="3,3"/>

  <!-- ═══ SERVICE SPLIT AT SCALE ═══ -->
  <line x1="30" y1="330" x2="1070" y2="330" stroke="#e2e8f0" stroke-width="1"/>
  <text x="550" y="355" text-anchor="middle" font-size="12" font-weight="600" fill="#1e293b" letter-spacing="1">MICROSERVICES AT 50M+ USERS (7 services, each owns its DB)</text>

  <!-- Service boxes -->
  <rect x="30" y="370" width="135" height="55" rx="6" fill="#e8f0fe" stroke="#4285f4" stroke-width="1" filter="url(#ig-sh)"/>
  <text x="97" y="390" text-anchor="middle" font-size="9" font-weight="600" fill="#1a73e8">User Service</text>
  <text x="97" y="403" text-anchor="middle" font-size="8" fill="#5f6368">Postgres (ACID)</text>
  <text x="97" y="416" text-anchor="middle" font-size="8" fill="#5f6368">sharded by user_id</text>

  <rect x="180" y="370" width="135" height="55" rx="6" fill="#f3e8fd" stroke="#9334e6" stroke-width="1" filter="url(#ig-sh)"/>
  <text x="247" y="390" text-anchor="middle" font-size="9" font-weight="600" fill="#7627bb">Post Service</text>
  <text x="247" y="403" text-anchor="middle" font-size="8" fill="#5f6368">DynamoDB</text>
  <text x="247" y="416" text-anchor="middle" font-size="8" fill="#5f6368">PK=user, SK=time</text>

  <rect x="330" y="370" width="135" height="55" rx="6" fill="#e0f7fa" stroke="#00897b" stroke-width="1" filter="url(#ig-sh)"/>
  <text x="397" y="390" text-anchor="middle" font-size="9" font-weight="600" fill="#00695c">Feed Service</text>
  <text x="397" y="403" text-anchor="middle" font-size="8" fill="#5f6368">Redis sorted sets</text>
  <text x="397" y="416" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">sub-ms reads</text>

  <rect x="480" y="370" width="135" height="55" rx="6" fill="#e6f4ea" stroke="#34a853" stroke-width="1" filter="url(#ig-sh)"/>
  <text x="547" y="390" text-anchor="middle" font-size="9" font-weight="600" fill="#137333">Media Service</text>
  <text x="547" y="403" text-anchor="middle" font-size="8" fill="#5f6368">S3 + Lambda</text>
  <text x="547" y="416" text-anchor="middle" font-size="8" fill="#5f6368">event-driven</text>

  <rect x="630" y="370" width="135" height="55" rx="6" fill="#fef7e0" stroke="#f9ab00" stroke-width="1" filter="url(#ig-sh)"/>
  <text x="697" y="390" text-anchor="middle" font-size="9" font-weight="600" fill="#e37400">Engagement Svc</text>
  <text x="697" y="403" text-anchor="middle" font-size="8" fill="#5f6368">DynamoDB sharded</text>
  <text x="697" y="416" text-anchor="middle" font-size="8" fill="#5f6368">100 counter shards</text>

  <rect x="780" y="370" width="135" height="55" rx="6" fill="#f1f3f4" stroke="#dadce0" stroke-width="1" filter="url(#ig-sh)"/>
  <text x="847" y="390" text-anchor="middle" font-size="9" font-weight="600" fill="#5f6368">Notification Svc</text>
  <text x="847" y="403" text-anchor="middle" font-size="8" fill="#5f6368">SQS + DynamoDB</text>
  <text x="847" y="416" text-anchor="middle" font-size="8" fill="#5f6368">5s delay OK</text>

  <rect x="930" y="370" width="140" height="55" rx="6" fill="#f1f3f4" stroke="#dadce0" stroke-width="1" filter="url(#ig-sh)"/>
  <text x="1000" y="390" text-anchor="middle" font-size="9" font-weight="600" fill="#5f6368">Search/Explore</text>
  <text x="1000" y="403" text-anchor="middle" font-size="8" fill="#5f6368">Elasticsearch</text>
  <text x="1000" y="416" text-anchor="middle" font-size="8" fill="#5f6368">autocomplete</text>

  <!-- Kafka bus connecting services -->
  <rect x="200" y="435" width="700" height="24" rx="12" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2"/>
  <text x="550" y="451" text-anchor="middle" font-size="9" font-weight="600" fill="#1a73e8">Kafka Event Bus — reads: sync gRPC (&lt;10ms) | writes: async Kafka events (never mix)</text>

  <!-- ═══ VIRAL POST HANDLING ═══ -->
  <line x1="30" y1="475" x2="1070" y2="475" stroke="#e2e8f0" stroke-width="1"/>
  <text x="550" y="500" text-anchor="middle" font-size="12" font-weight="600" fill="#1e293b" letter-spacing="1">VIRAL POST: 1M LIKES/SECOND</text>

  <!-- Problem -->
  <rect x="30" y="515" width="250" height="65" rx="8" fill="#fce8e6" stroke="#ea4335" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="155" y="535" text-anchor="middle" font-size="10" font-weight="700" fill="#c5221f">Problem: Single Counter</text>
  <text x="155" y="550" text-anchor="middle" font-size="9" fill="#202124">one Redis key = 100K ops/sec max</text>
  <text x="155" y="565" text-anchor="middle" font-size="9" fill="#202124">1M likes/sec → 10x over capacity</text>
  <text x="155" y="575" text-anchor="middle" font-size="8" fill="#ea4335">hot key bottleneck</text>

  <line x1="280" y1="547" x2="330" y2="547" stroke="#34a853" stroke-width="1.5" marker-end="url(#ig-arr-g)"/>

  <!-- Solution -->
  <rect x="332" y="515" width="280" height="65" rx="8" fill="#e6f4ea" stroke="#34a853" stroke-width="1.5" filter="url(#ig-sh)"/>
  <text x="472" y="535" text-anchor="middle" font-size="10" font-weight="700" fill="#137333">Solution: 100 Sharded Counters</text>
  <text x="472" y="550" text-anchor="middle" font-size="9" fill="#202124">random shard selection, lazy aggregate q5s</text>
  <text x="472" y="565" text-anchor="middle" font-size="9" fill="#202124">dedup: DynamoDB PutItem attribute_not_exists</text>
  <text x="472" y="575" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">100 shards × 100K = 10M ops/sec capacity</text>

  <!-- Graceful degradation -->
  <rect x="640" y="515" width="430" height="65" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="855" y="535" text-anchor="middle" font-size="10" font-weight="700" fill="#e37400">Graceful Degradation Order</text>
  <text x="855" y="550" text-anchor="middle" font-size="9" fill="#202124">1. stop explore batch → 2. buffer likes 10s → 3. skip celebrity merge → 4. disable story tracking</text>
  <text x="855" y="565" text-anchor="middle" font-size="9" fill="#202124">Feed serving + uploads degrade LAST</text>
  <text x="855" y="575" text-anchor="middle" font-size="8" fill="#5f6368">auto-scale + SQS buffering: $15K/mo vs $125K/mo over-provision</text>

  <!-- ═══ SCALE EVOLUTION ═══ -->
  <line x1="30" y1="600" x2="1070" y2="600" stroke="#e2e8f0" stroke-width="1"/>
  <text x="550" y="625" text-anchor="middle" font-size="12" font-weight="600" fill="#1e293b" letter-spacing="1">SCALE EVOLUTION</text>

  <rect x="30" y="640" width="200" height="68" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="130" y="658" text-anchor="middle" font-size="10" font-weight="700" fill="#1a73e8">0-1M: Monolith</text>
  <text x="130" y="672" text-anchor="middle" font-size="9" fill="#202124">Postgres, pull-based feed</text>
  <text x="130" y="685" text-anchor="middle" font-size="9" fill="#202124">IN-clause on posts table</text>
  <text x="130" y="698" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">~$500/mo</text>

  <line x1="230" y1="674" x2="260" y2="674" stroke="#4285f4" stroke-width="2" marker-end="url(#ig-arr)"/>

  <rect x="262" y="640" width="230" height="68" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="377" y="658" text-anchor="middle" font-size="10" font-weight="700" fill="#1a73e8">1M-50M: Hybrid Fan-out</text>
  <text x="377" y="672" text-anchor="middle" font-size="9" fill="#202124">Redis ZSET feeds, SQS async</text>
  <text x="377" y="685" text-anchor="middle" font-size="9" fill="#202124">skip inactive users (&gt;7d)</text>
  <text x="377" y="698" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">~$12.5K/mo</text>

  <line x1="492" y1="674" x2="522" y2="674" stroke="#4285f4" stroke-width="2" marker-end="url(#ig-arr)"/>

  <rect x="524" y="640" width="240" height="68" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="644" y="658" text-anchor="middle" font-size="10" font-weight="700" fill="#1a73e8">50M-500M: Microservices</text>
  <text x="644" y="672" text-anchor="middle" font-size="9" fill="#202124">7 services, Kafka, DynamoDB</text>
  <text x="644" y="685" text-anchor="middle" font-size="9" fill="#202124">gRPC reads, async writes</text>
  <text x="644" y="698" text-anchor="middle" font-size="8" fill="#5f6368">feed load: ~36ms total</text>

  <line x1="764" y1="674" x2="794" y2="674" stroke="#4285f4" stroke-width="2" marker-end="url(#ig-arr)"/>

  <rect x="796" y="640" width="274" height="68" rx="8" fill="#e8f0fe" stroke="#4285f4" stroke-width="1.2" filter="url(#ig-sh)"/>
  <text x="933" y="658" text-anchor="middle" font-size="10" font-weight="700" fill="#1a73e8">500M+: Multi-Region</text>
  <text x="933" y="672" text-anchor="middle" font-size="9" fill="#202124">3 regions, CRDTs for counters</text>
  <text x="933" y="685" text-anchor="middle" font-size="9" fill="#202124">DynamoDB Global Tables &lt;1s</text>
  <text x="933" y="698" text-anchor="middle" font-size="8" font-weight="600" fill="#f9ab00">EU: 150ms → 20ms</text>

  <!-- ═══ KEY NUMBERS ═══ -->
  <line x1="30" y1="725" x2="1070" y2="725" stroke="#e2e8f0" stroke-width="1"/>

  <rect x="30" y="740" width="140" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="100" y="758" text-anchor="middle" font-size="9" font-weight="700" fill="#e37400">200TB/day</text>
  <text x="100" y="772" text-anchor="middle" font-size="8" fill="#5f6368">media uploads</text>

  <rect x="185" y="740" width="140" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="255" y="758" text-anchor="middle" font-size="9" font-weight="700" fill="#e37400">10:1 R:W</text>
  <text x="255" y="772" text-anchor="middle" font-size="8" fill="#5f6368">58K read, 6K write</text>

  <rect x="340" y="740" width="140" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="410" y="758" text-anchor="middle" font-size="9" font-weight="700" fill="#e37400">$3.8M/mo</text>
  <text x="410" y="772" text-anchor="middle" font-size="8" fill="#5f6368">CDN at 500M DAU</text>

  <rect x="495" y="740" width="140" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="565" y="758" text-anchor="middle" font-size="9" font-weight="700" fill="#e37400">~36ms</text>
  <text x="565" y="772" text-anchor="middle" font-size="8" fill="#5f6368">feed load latency</text>

  <rect x="650" y="740" width="140" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="720" y="758" text-anchor="middle" font-size="9" font-weight="700" fill="#e37400">100K threshold</text>
  <text x="720" y="772" text-anchor="middle" font-size="8" fill="#5f6368">fan-out strategy</text>

  <rect x="805" y="740" width="140" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="875" y="758" text-anchor="middle" font-size="9" font-weight="700" fill="#e37400">73PB/year</text>
  <text x="875" y="772" text-anchor="middle" font-size="8" fill="#5f6368">storage growth</text>

  <rect x="960" y="740" width="110" height="42" rx="8" fill="#fef7e0" stroke="#f9ab00" stroke-width="1.2"/>
  <text x="1015" y="758" text-anchor="middle" font-size="9" font-weight="700" fill="#e37400">AP feed</text>
  <text x="1015" y="772" text-anchor="middle" font-size="8" fill="#5f6368">CP for likes</text>

  <!-- Bottom takeaway -->
  <rect x="200" y="798" width="700" height="18" rx="9" fill="#e8f0fe" stroke="#4285f4" stroke-width="1"/>
  <text x="550" y="811" text-anchor="middle" font-size="9" font-weight="600" fill="#1a73e8">Core insight: Hybrid fan-out (write for normal, read for celebrities) is THE differentiating answer</text>
</svg>`,
}

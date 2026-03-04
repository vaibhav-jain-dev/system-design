package diagrams

func registerSearchAutocomplete(r *Registry) {
	r.Register(&Diagram{
		Slug:        "sa-requirements",
		Title:       "Scale Estimates & Requirements",
		Description: "Functional and non-functional requirements with scale estimates for search autocomplete",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Estimates</div>
      <div class="d-flow-v">
        <div class="d-box blue">5B searches/day &#8594; 60K QPS</div>
        <div class="d-box blue">Avg 4 keystrokes before selection</div>
        <div class="d-box blue">60K &#215; 4 = 240K suggestion QPS</div>
        <div class="d-box purple">Peak (3x) = 720K QPS</div>
        <div class="d-box amber">100M unique query prefixes</div>
        <div class="d-box amber">Avg prefix 15 chars &#215; 100M = 1.5 GB raw</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green">Return top 10 suggestions per prefix</div>
        <div class="d-box green">Rank by popularity + personalization</div>
        <div class="d-box green">Support fuzzy matching (typo tolerance)</div>
        <div class="d-box blue">Filter inappropriate/harmful content</div>
        <div class="d-box blue">Geo-aware suggestions</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">Latency: &lt; 100ms p99</div>
        <div class="d-box purple">Availability: 99.99%</div>
        <div class="d-box amber">Freshness: trending within 15 min</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-api-design",
		Title:       "API Endpoints",
		Description: "REST API design for autocomplete suggestions and trending queries",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Core Suggestion API</div>
        <div class="d-flow-v">
          <div class="d-box green" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">GET /v1/suggestions?q={prefix}
    &amp;limit=10
    &amp;lang=en
    &amp;lat=37.77&amp;lon=-122.42
    &amp;user_id=abc123

200 OK
{
  "suggestions": [
    {"text": "search engine design", "score": 0.95},
    {"text": "search algorithms", "score": 0.87},
    ...
  ],
  "took_ms": 12
}</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Query Logging (Async)</div>
        <div class="d-flow-v">
          <div class="d-box blue" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">POST /v1/queries/log
{
  "query": "search engine design",
  "user_id": "abc123",
  "selected_rank": 1,
  "timestamp": 1704067200
}</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Trending API</div>
        <div class="d-flow-v">
          <div class="d-box amber" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">GET /v1/trending?region=US&amp;limit=10

200 OK
{"trending": ["breaking news", ...]}</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-trie-structure",
		Title:       "Trie Data Structure Visualization",
		Description: "How prefix trie stores suggestions with top-K results cached at each node",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box indigo" style="text-align:center"><strong>Root</strong></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box blue" style="text-align:center"><strong>s</strong></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue" style="text-align:center"><strong>se</strong></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-flow">
          <div class="d-branch">
            <div class="d-branch-arm">
              <div class="d-box green" style="text-align:center"><strong>sea</strong><br><small>top-10 cached</small></div>
              <div class="d-arrow-down">&#8595;</div>
              <div class="d-box green" style="text-align:center"><strong>sear</strong><br><small>search: 50K</small><br><small>seattle: 12K</small></div>
            </div>
            <div class="d-branch-arm">
              <div class="d-box amber" style="text-align:center"><strong>sel</strong><br><small>top-10 cached</small></div>
              <div class="d-arrow-down">&#8595;</div>
              <div class="d-box amber" style="text-align:center"><strong>self</strong><br><small>self care: 8K</small></div>
            </div>
          </div>
        </div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box purple" style="text-align:center"><strong>t</strong></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" style="text-align:center"><strong>te</strong></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" style="text-align:center"><strong>tes</strong><br><small>tesla: 30K</small><br><small>test: 25K</small></div>
      </div>
    </div>
  </div>
  <div class="d-label">Each node caches top-10 suggestions &#8212; eliminates DFS at query time. O(L) lookup where L = prefix length.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end architecture from client through CDN, API servers, Trie service, and data pipeline",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>Client (Browser/App)</strong><br>Debounce 200ms &#8594; GET /v1/suggestions?q=sea</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" style="text-align:center"><strong>CDN (CloudFront)</strong><br>Cache popular prefixes &#8226; TTL 5 min &#8226; ~60% hit rate</div>
  <div class="d-arrow-down">&#8595; cache miss</div>
  <div class="d-box green" style="text-align:center"><strong>API Gateway + Load Balancer</strong></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo" style="text-align:center"><strong>Suggestion Service (Stateless)</strong><br>Prefix lookup &#8594; Trie Service &#8594; Rank &#8594; Filter &#8594; Return</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box red" style="text-align:center"><strong>Trie Service</strong><br>In-memory Trie<br>Sharded by prefix range</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box amber" style="text-align:center"><strong>Redis Cache</strong><br>Hot prefix results<br>TTL 60s</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box purple" style="text-align:center"><strong>Personalization</strong><br>User history from<br>DynamoDB</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; offline pipeline</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box gray" style="text-align:center"><strong>Kafka</strong><br>Query logs stream</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box gray" style="text-align:center"><strong>Flink/Spark</strong><br>Aggregate counts<br>Build new Trie</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box gray" style="text-align:center"><strong>S3</strong><br>Trie snapshots<br>for recovery</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-ranking-algorithm",
		Title:       "Ranking Algorithm &#8212; Popularity + Personalization",
		Description: "Scoring formula combining global popularity, freshness decay, and user personalization signals",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box indigo" style="text-align:center"><strong>Final Score = w1&#183;Popularity + w2&#183;Freshness + w3&#183;Personal + w4&#183;Geo</strong></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Popularity (w1 = 0.5)</div>
        <div class="d-flow-v">
          <div class="d-box green">log10(query_count_30d)</div>
          <div class="d-label">Logarithmic scale prevents mega-queries from dominating</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Freshness (w2 = 0.2)</div>
        <div class="d-flow-v">
          <div class="d-box amber">decay = e^(-&#955; &#183; age_hours)</div>
          <div class="d-label">&#955; = 0.01 &#8594; half-life ~70 hours. Trending topics boosted</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Personalization (w3 = 0.2)</div>
        <div class="d-flow-v">
          <div class="d-box purple">User search history overlap</div>
          <div class="d-box purple">Clicked results affinity</div>
          <div class="d-label">Fall back to demographic cohort if new user</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Geo Boost (w4 = 0.1)</div>
        <div class="d-flow-v">
          <div class="d-box blue">Location relevance score</div>
          <div class="d-label">"pizza" in NYC &#8594; boost "pizza near me"</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-fuzzy-matching",
		Title:       "Fuzzy Matching &#8212; Typo Tolerance",
		Description: "Edit distance and phonetic matching approaches for handling misspelled queries",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Edit Distance (Levenshtein)</div>
      <div class="d-flow-v">
        <div class="d-box blue" style="text-align:center">User types: "amzon"</div>
        <div class="d-arrow-down">&#8595; edit distance = 1</div>
        <div class="d-box green" style="text-align:center">Match: "amazon" (insert 'a')</div>
        <div class="d-label">Allow edit distance &#8804; 2 for queries &gt; 4 chars</div>
        <div class="d-box amber" style="text-align:center">
          <strong>Implementation</strong><br>
          BK-Tree: O(log N) lookup<br>
          SymSpell: O(1) with precomputed deletes
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Phonetic Matching (Soundex/Metaphone)</div>
      <div class="d-flow-v">
        <div class="d-box blue" style="text-align:center">User types: "fone"</div>
        <div class="d-arrow-down">&#8595; same phonetic code</div>
        <div class="d-box green" style="text-align:center">Match: "phone"</div>
        <div class="d-label">Double Metaphone handles multilingual names</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Combined Pipeline</div>
      <div class="d-flow-v">
        <div class="d-box purple" style="text-align:center">
          1. Exact prefix match (Trie)<br>
          2. Edit distance &#8804; 1 corrections<br>
          3. Phonetic fallback<br>
          4. Merge + re-rank by score
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-trending-update",
		Title:       "Real-Time Trending Pipeline",
		Description: "Streaming pipeline from query logs through Kafka and Flink to update Trie with trending queries",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>User Search Queries</strong><br>240K QPS &#8594; query log events</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" style="text-align:center"><strong>Kafka (query-logs topic)</strong><br>Partitioned by prefix[0:2] &#8226; 7-day retention</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box amber" style="text-align:center"><strong>Flink &#8212; Real-Time</strong><br>Sliding window: 15 min<br>Count-Min Sketch for freq<br>Detect &#8805; 3x spike</div>
        <div class="d-arrow-down">&#8595; trending alert</div>
        <div class="d-box red" style="text-align:center"><strong>Hot Update</strong><br>Push trending to Trie nodes<br>via gRPC &#8226; &lt; 1 min delay</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box green" style="text-align:center"><strong>Spark &#8212; Batch (Daily)</strong><br>Full recount 30-day window<br>Rebuild Trie from scratch<br>Snapshot to S3</div>
        <div class="d-arrow-down">&#8595; scheduled swap</div>
        <div class="d-box green" style="text-align:center"><strong>Blue/Green Deploy</strong><br>New Trie replicas load snapshot<br>Traffic shifts atomically</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-data-model",
		Title:       "Data Model &#8212; Tables & Storage",
		Description: "Database schema for queries, suggestions, user history, and content filter tables",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">query_aggregates (DynamoDB)</div>
      <div class="d-entity-body">
        <div class="pk">PK: prefix (first 3 chars)</div>
        <div class="fk">SK: full_query</div>
        <div class="idx idx-hash">count_30d INT</div>
        <div class="idx idx-hash">count_7d INT</div>
        <div>last_seen TIMESTAMP</div>
        <div>language VARCHAR(5)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header green">trie_snapshots (S3)</div>
      <div class="d-entity-body">
        <div class="pk">KEY: trie/{shard_id}/{timestamp}.bin</div>
        <div>Serialized Trie (protobuf)</div>
        <div>~500 MB per shard</div>
        <div>Daily rebuild + incremental</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header purple">user_search_history (DynamoDB)</div>
      <div class="d-entity-body">
        <div class="pk">PK: user_id</div>
        <div class="fk">SK: timestamp#query</div>
        <div>query VARCHAR(200)</div>
        <div>selected_suggestion VARCHAR(200)</div>
        <div>TTL: 90 days</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header red">blocked_terms (DynamoDB)</div>
      <div class="d-entity-body">
        <div class="pk">PK: term</div>
        <div>category ENUM (hate|violence|adult|spam)</div>
        <div>severity ENUM (hard_block|soft_block)</div>
        <div>added_by VARCHAR</div>
        <div>added_at TIMESTAMP</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-caching-strategy",
		Title:       "Multi-Layer Caching Strategy",
		Description: "Three-tier caching: CDN for popular prefixes, Redis for warm results, in-memory Trie for all lookups",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>Client</strong><br>Local cache: last 50 queries &#8226; instant for repeated prefixes</div>
  <div class="d-arrow-down">&#8595; cache miss</div>
  <div class="d-box green" style="text-align:center"><strong>CDN (CloudFront)</strong><br>Cache key: prefix + region + lang<br>TTL: 5 min &#8226; Hit rate: ~60% for top 1000 prefixes<br>Cost: $0.01/10K requests</div>
  <div class="d-arrow-down">&#8595; cache miss</div>
  <div class="d-box amber" style="text-align:center"><strong>Redis (ElastiCache)</strong><br>Full suggestion lists for warm prefixes<br>TTL: 60s &#8226; Hit rate: ~85% of remaining<br>Key: sug:{prefix}:{lang}:{region}</div>
  <div class="d-arrow-down">&#8595; cache miss</div>
  <div class="d-box red" style="text-align:center"><strong>Trie Service (In-Memory)</strong><br>Full Trie traversal &#8594; rank &#8594; filter &#8594; return top 10<br>Latency: 5-15ms &#8226; Always available &#8226; Source of truth</div>
  <div class="d-label">Overall: ~94% of requests served from CDN + Redis. Only 6% hit Trie directly.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-scaling",
		Title:       "Trie Sharding by Prefix Range",
		Description: "Sharding strategy splitting the Trie across multiple servers by prefix character ranges",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box indigo" style="text-align:center"><strong>Prefix Router</strong><br>Hash first 2 chars &#8594; shard assignment</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 1: a-f</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center"><strong>Primary</strong><br>~25M prefixes<br>3 GB memory</div>
        <div class="d-box gray" style="text-align:center">Replica 1</div>
        <div class="d-box gray" style="text-align:center">Replica 2</div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 2: g-m</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center"><strong>Primary</strong><br>~30M prefixes<br>3.6 GB memory</div>
        <div class="d-box gray" style="text-align:center">Replica 1</div>
        <div class="d-box gray" style="text-align:center">Replica 2</div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 3: n-s</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center"><strong>Primary</strong><br>~28M prefixes<br>3.4 GB memory</div>
        <div class="d-box gray" style="text-align:center">Replica 1</div>
        <div class="d-box gray" style="text-align:center">Replica 2</div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 4: t-z</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center"><strong>Primary</strong><br>~17M prefixes<br>2 GB memory</div>
        <div class="d-box gray" style="text-align:center">Replica 1</div>
        <div class="d-box gray" style="text-align:center">Replica 2</div>
      </div>
    </div>
  </div>
  <div class="d-label">4 shards &#215; 3 replicas = 12 nodes &#8226; Uneven distribution: shard by 2-char prefix hash for better balance</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-filtering",
		Title:       "Content Filtering Pipeline",
		Description: "Multi-stage pipeline for filtering inappropriate, harmful, and low-quality suggestions",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>Raw Suggestions (Top 20 candidates)</strong></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red" style="text-align:center"><strong>Stage 1: Blocklist Filter</strong><br>Exact match against 50K blocked terms &#8226; O(1) hash lookup<br>Hard block: hate speech, violence, illegal content</div>
  <div class="d-arrow-down">&#8595; pass</div>
  <div class="d-box amber" style="text-align:center"><strong>Stage 2: Regex Pattern Filter</strong><br>Phone numbers, SSNs, email patterns<br>PII detection &#8226; ~200 regex rules</div>
  <div class="d-arrow-down">&#8595; pass</div>
  <div class="d-box purple" style="text-align:center"><strong>Stage 3: ML Classifier (Async)</strong><br>BERT-tiny model &#8226; Toxicity score 0-1<br>Block if score &gt; 0.8 &#8226; 2ms inference</div>
  <div class="d-arrow-down">&#8595; pass</div>
  <div class="d-box green" style="text-align:center"><strong>Stage 4: Quality Filter</strong><br>Remove duplicates (normalized), too-short (&lt; 2 chars),<br>too-long (&gt; 60 chars), gibberish (entropy check)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" style="text-align:center"><strong>Final: Top 10 Clean Suggestions</strong></div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-geo-suggestions",
		Title:       "Geo-Aware Suggestions",
		Description: "Location-based suggestion boosting with regional Trie overlays and geo-hash bucketing",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Architecture</div>
      <div class="d-flow-v">
        <div class="d-box blue" style="text-align:center"><strong>User Query</strong><br>"pizza" + lat/lon</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" style="text-align:center"><strong>Geo Hash</strong><br>lat/lon &#8594; geohash6<br>e.g., "9q8yyk" (SF)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-flow">
          <div class="d-branch">
            <div class="d-branch-arm">
              <div class="d-box green" style="text-align:center"><strong>Global Trie</strong><br>Base suggestions</div>
            </div>
            <div class="d-branch-arm">
              <div class="d-box purple" style="text-align:center"><strong>Regional Overlay</strong><br>Geo-boosted results</div>
            </div>
          </div>
        </div>
        <div class="d-arrow-down">&#8595; merge + re-rank</div>
        <div class="d-box green" style="text-align:center"><strong>Final Results</strong></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Example: "pizza" in San Francisco</div>
      <div class="d-flow-v">
        <div class="d-box green">1. "pizza near me" (geo-boosted &#215;1.5)</div>
        <div class="d-box green">2. "pizza delivery sf" (geo-boosted &#215;1.3)</div>
        <div class="d-box blue">3. "pizza dough recipe" (global)</div>
        <div class="d-box blue">4. "pizza hut" (global)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Regional Trie Storage</div>
      <div class="d-flow-v">
        <div class="d-box gray" style="text-align:center">~200 metro regions<br>Each: 1-5M local queries<br>Total: ~2 GB additional memory</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

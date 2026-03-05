package diagrams

func registerSearchAutocomplete(r *Registry) {
	r.Register(&Diagram{
		Slug:        "sa-requirements",
		Title:       "Functional & Non-Functional Requirements",
		Description: "Scale targets and requirements for search autocomplete: 5B queries/day, 500M users, <100ms p99",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P0 &#8212; Core (Must Have)</div>
      <div class="d-flow-v">
        <div class="d-box green">Return top-10 suggestions for prefix</div>
        <div class="d-box green">&#8804; 100ms p99 latency end-to-end</div>
        <div class="d-box green">Support 500M daily active users</div>
        <div class="d-box green">Handle 5B queries/day (&#8776; 58K QPS avg)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P1 &#8212; Important</div>
      <div class="d-flow-v">
        <div class="d-box blue">Trending queries surface within 15 min</div>
        <div class="d-box blue">Personalized suggestions per user</div>
        <div class="d-box blue">Multi-language / Unicode support</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P2 &#8212; Nice to Have</div>
      <div class="d-flow-v">
        <div class="d-box gray">Spell correction / fuzzy matching</div>
        <div class="d-box gray">Query categorization (web, images, news)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Estimates</div>
      <div class="d-flow-v">
        <div class="d-box purple">5B queries/day &#8776; 58K QPS average</div>
        <div class="d-box purple">Peak (5x) &#8776; 290K QPS</div>
        <div class="d-box purple">Avg 4 keystrokes per query &#8594; 230K suggestion QPS avg</div>
        <div class="d-box purple">Peak suggestions &#8776; 1.15M QPS</div>
        <div class="d-box amber">10 suggestions &#215; 50 bytes avg &#8776; 500B per response</div>
        <div class="d-box amber">Storage: 1B unique queries &#215; 30B avg &#8776; 30 GB trie</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Key Constraints</div>
      <div class="d-flow-v">
        <div class="d-box red">p99 latency &lt; 100ms (including network)</div>
        <div class="d-box red">Availability: 99.99% &#8212; search is critical path</div>
        <div class="d-box red">Freshness: trending queries within 15 min</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-api-design",
		Title:       "API Design",
		Description: "REST endpoint for suggestions, WebSocket streaming alternative, and analytics endpoint",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Primary &#8212; REST Suggestion Endpoint</div>
      <div class="d-flow-v">
        <div class="d-box green" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">GET /v1/suggest?q={prefix}&amp;limit=10&amp;lang=en&amp;user_id={uid}

Response 200:
{
  "query": "how to",
  "suggestions": [
    {"text": "how to tie a tie", "score": 0.95, "type": "trending"},
    {"text": "how to cook rice", "score": 0.89, "type": "popular"},
    ...
  ],
  "request_id": "abc-123",
  "latency_ms": 12
}

Headers:
  Cache-Control: public, max-age=60
  X-Request-Id: abc-123</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Alternative &#8212; WebSocket Streaming</div>
      <div class="d-flow-v">
        <div class="d-box blue" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">WS /v1/suggest/stream

Client &#8594; {"type":"keystroke","q":"ho"}
Server &#8594; {"type":"suggestions","q":"ho",
              "suggestions":[...]}

Client &#8594; {"type":"keystroke","q":"how"}
Server &#8594; {"type":"suggestions","q":"how",
              "suggestions":[...]}</div>
        <div class="d-box amber">Tradeoff: lower latency per keystroke, but higher server cost (persistent connections)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Analytics Endpoint</div>
      <div class="d-flow-v">
        <div class="d-box purple" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">POST /v1/analytics/query
{
  "query": "how to tie a tie",
  "selected_suggestion": 2,
  "session_id": "sess-456",
  "timestamp": 1704067200
}</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-trie-structure",
		Title:       "Trie Data Structure with Top-K Caching",
		Description: "Prefix trie with node-level counts and precomputed top-K suggestions at each node",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Trie with Precomputed Top-K at Each Node</div>
    <div class="d-flow-v">
      <div class="d-flow">
        <div class="d-box gray">(root) &#8212; 1B queries</div>
      </div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-flow">
        <div class="d-box blue">'h' &#8212; 120M queries<br/>top-K: [hello, how, hotel, ...]</div>
        <div class="d-box gray">'a' &#8212; 95M</div>
        <div class="d-box gray">'t' &#8212; 88M</div>
        <div class="d-box gray">... 23 more</div>
      </div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-flow">
        <div class="d-box green">'ho' &#8212; 45M queries<br/>top-K: [how to, hotel, home, ...]</div>
        <div class="d-box gray">'he' &#8212; 38M</div>
        <div class="d-box gray">'ha' &#8212; 22M</div>
      </div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-flow">
        <div class="d-box purple">'how' &#8212; 28M queries<br/>top-K: [how to, how much, how many, ...]</div>
        <div class="d-box gray">'hot' &#8212; 9M</div>
        <div class="d-box gray">'hom' &#8212; 8M</div>
      </div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-flow">
        <div class="d-box indigo">'how ' &#8212; 25M queries<br/>top-K: [how to tie a tie, how to cook, ...]</div>
      </div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Node Structure</div>
        <div class="d-flow-v">
          <div class="d-box amber" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">type TrieNode struct {
  children   [26+]*TrieNode  // a-z + special
  topK       []Suggestion    // precomputed top-10
  isTerminal bool
  frequency  int64
}</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Lookup: O(L) where L = prefix length</div>
        <div class="d-flow-v">
          <div class="d-box green">1. Walk trie: root &#8594; h &#8594; o &#8594; w</div>
          <div class="d-box green">2. Return node.topK (precomputed)</div>
          <div class="d-box green">3. No DFS needed &#8212; O(1) after walk</div>
          <div class="d-box blue">Memory: &#8776; 30 GB for 1B unique queries</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end architecture: CDN edge cache, ALB, Suggestion Service, Trie, Redis, DynamoDB",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Client<br/>(browser / mobile)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">CDN Edge Cache<br/>(CloudFront)<br/>TTL 60s for popular prefixes</div>
    <div class="d-arrow">&#8594; cache miss</div>
    <div class="d-box purple">ALB<br/>(route by prefix hash)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box indigo">Suggestion Service<br/>(stateful, trie in memory)</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Hot Path (read)</div>
        <div class="d-flow-v">
          <div class="d-box indigo">Suggestion Service</div>
          <div class="d-arrow-down">&#8595; O(L) lookup</div>
          <div class="d-box amber">In-Memory Trie<br/>&#8776; 30 GB per node<br/>precomputed top-K</div>
          <div class="d-arrow-down">&#8595; cache miss on personalization</div>
          <div class="d-box red">Redis Cluster<br/>user history + trending cache<br/>TTL 5 min</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Warm Path (update)</div>
        <div class="d-flow-v">
          <div class="d-box gray">Query Log (Kafka)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Spark Aggregation<br/>(hourly job)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Trie Builder<br/>(MapReduce)</div>
          <div class="d-arrow-down">&#8595; atomic swap</div>
          <div class="d-box amber">In-Memory Trie</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Cold Path (persist)</div>
        <div class="d-flow-v">
          <div class="d-box green">DynamoDB<br/>query_aggregates table<br/>persistent source of truth</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">S3<br/>serialized trie snapshots<br/>versioned, immutable</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-data-collection",
		Title:       "Query Logging &amp; Aggregation Pipeline",
		Description: "Data collection pipeline: search queries through Kafka to Spark aggregation to trie rebuild",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">User Search<br/>"how to cook"</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">Search Service<br/>(logs query + metadata)</div>
    <div class="d-arrow">&#8594; async</div>
    <div class="d-box green">Kafka<br/>topic: search-queries<br/>partitioned by prefix[0:2]</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Real-Time Path (&#8804; 15 min)</div>
        <div class="d-flow-v">
          <div class="d-box amber">Flink Streaming Job</div>
          <div class="d-arrow-down">&#8595; 5-min tumbling window</div>
          <div class="d-box amber">Trending Detector<br/>z-score &gt; 3&#963; &#8594; trending</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red">Redis Sorted Set<br/>trending:global<br/>ZADD score=frequency</div>
          <div class="d-arrow-down">&#8595; merge</div>
          <div class="d-box indigo">Suggestion Service<br/>(hot-patch trending into trie)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Batch Path (hourly)</div>
        <div class="d-flow-v">
          <div class="d-box purple">Spark Aggregation<br/>GROUP BY query, date<br/>SUM(count), decay old</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">DynamoDB<br/>query_aggregates<br/>updated frequency + recency</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Trie Builder (MapReduce)<br/>build new trie from aggregates</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">S3 Snapshot<br/>serialized trie v{N+1}</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Log Record Schema</div>
        <div class="d-flow-v">
          <div class="d-box gray" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">{
  "query": "how to cook rice",
  "user_id": "u-12345",
  "timestamp": 1704067200,
  "session_id": "sess-789",
  "locale": "en-US",
  "device": "mobile",
  "result_clicked": true,
  "position": 2
}</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-ranking-algorithm",
		Title:       "Suggestion Ranking Algorithm",
		Description: "Scoring formula: frequency, recency, personalization, and trending boost factors",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Final Score = frequency &#215; recency &#215; personalization &#215; trending</div>
    <div class="d-flow">
      <div class="d-box blue" style="min-width:160px">
        <strong>Frequency Score</strong><br/>
        log10(total_searches)<br/>
        &#8212; dampens popularity bias<br/>
        Range: 0 &#8594; 10
      </div>
      <div class="d-box green" style="min-width:160px">
        <strong>Recency Decay</strong><br/>
        e^(-&#955; &#215; age_hours)<br/>
        &#955; = 0.01 (24h half-life)<br/>
        Range: 0.0 &#8594; 1.0
      </div>
      <div class="d-box purple" style="min-width:160px">
        <strong>Personalization</strong><br/>
        1.0 + (0.5 &#215; user_affinity)<br/>
        user_affinity from history<br/>
        Range: 1.0 &#8594; 1.5
      </div>
      <div class="d-box red" style="min-width:160px">
        <strong>Trending Boost</strong><br/>
        1.0 + min(2.0, z_score / 3)<br/>
        z_score = (current - avg) / stddev<br/>
        Range: 1.0 &#8594; 3.0
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Example: "how to cook"</div>
        <div class="d-flow-v">
          <div class="d-box amber">frequency = log10(2,800,000) = 6.45</div>
          <div class="d-box amber">recency = e^(-0.01 &#215; 2) = 0.98</div>
          <div class="d-box amber">personalization = 1.0 (no history)</div>
          <div class="d-box amber">trending = 1.0 (not trending)</div>
          <div class="d-box green"><strong>final = 6.45 &#215; 0.98 &#215; 1.0 &#215; 1.0 = 6.32</strong></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Example: "super bowl score" (trending)</div>
        <div class="d-flow-v">
          <div class="d-box amber">frequency = log10(50,000) = 4.70</div>
          <div class="d-box amber">recency = e^(-0.01 &#215; 0.5) = 0.995</div>
          <div class="d-box amber">personalization = 1.3 (sports fan)</div>
          <div class="d-box amber">trending = 1.0 + min(2.0, 8.5/3) = 3.0</div>
          <div class="d-box green"><strong>final = 4.70 &#215; 0.995 &#215; 1.3 &#215; 3.0 = 18.24</strong></div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-trie-update",
		Title:       "Offline Trie Rebuild &amp; Atomic Swap",
		Description: "MapReduce trie rebuild pipeline with serialization, versioning, and zero-downtime atomic swap",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box purple">DynamoDB<br/>query_aggregates<br/>(source of truth)</div>
    <div class="d-arrow">&#8594; full scan</div>
    <div class="d-box blue">MapReduce Job<br/>Phase 1: Sort by prefix<br/>Phase 2: Build trie partitions</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">Merge &amp; Compute Top-K<br/>at every node (heap of size K)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber">Serialize to Binary<br/>trie-v{N+1}.bin<br/>&#8776; 30 GB compressed</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box amber">Upload to S3<br/>s3://tries/trie-v{N+1}.bin</div>
    <div class="d-arrow">&#8594; notify</div>
    <div class="d-box indigo">Coordinator Service<br/>(orchestrates rollout)</div>
    <div class="d-arrow">&#8594; rolling update</div>
    <div class="d-box indigo">Suggestion Nodes<br/>(download + deserialize)</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Atomic Swap on Each Node</div>
    <div class="d-cols">
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box green">1. Download trie-v{N+1} from S3</div>
          <div class="d-box green">2. Deserialize into new Trie object</div>
          <div class="d-box green">3. Warm up: prefetch hot prefixes</div>
          <div class="d-box green">4. Atomic pointer swap: trie = newTrie</div>
          <div class="d-box green">5. Old trie garbage collected</div>
        </div>
      </div>
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box gray">Zero downtime &#8212; reads hit old trie until swap</div>
          <div class="d-box gray">Rollback: swap back to trie-v{N} pointer</div>
          <div class="d-box gray">Canary: update 5% of nodes first, monitor p99</div>
          <div class="d-box gray">Rebuild frequency: hourly (batch), instant (trending hot-patch)</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-data-model",
		Title:       "Data Model",
		Description: "Core tables: search_queries, query_aggregates, trie_snapshots, user_search_history, trending_queries",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">search_queries</div>
      <div class="d-entity-body">
        <div><span class="pk">query_id</span> UUID</div>
        <div>query_text VARCHAR(255)</div>
        <div>user_id VARCHAR(64)</div>
        <div>session_id VARCHAR(64)</div>
        <div>timestamp BIGINT</div>
        <div>locale VARCHAR(10)</div>
        <div>device_type ENUM(desktop, mobile, tablet)</div>
        <div>result_clicked BOOLEAN</div>
        <div>click_position INT</div>
        <div><span class="idx idx-btree">idx_query_text_ts</span></div>
        <div><span class="idx idx-hash">idx_user_id</span></div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header blue">query_aggregates</div>
      <div class="d-entity-body">
        <div><span class="pk">query_text</span> VARCHAR(255)</div>
        <div><span class="pk">date</span> DATE</div>
        <div>daily_count BIGINT</div>
        <div>weekly_count BIGINT</div>
        <div>monthly_count BIGINT</div>
        <div>all_time_count BIGINT</div>
        <div>avg_click_position DECIMAL(3,1)</div>
        <div>last_searched_at BIGINT</div>
        <div><span class="idx idx-btree">idx_alltime_count</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header green">trie_snapshots</div>
      <div class="d-entity-body">
        <div><span class="pk">version</span> INT</div>
        <div>s3_path VARCHAR(512)</div>
        <div>size_bytes BIGINT</div>
        <div>node_count BIGINT</div>
        <div>query_count BIGINT</div>
        <div>built_at TIMESTAMP</div>
        <div>status ENUM(building, ready, active, archived)</div>
        <div>checksum_sha256 VARCHAR(64)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header green">user_search_history</div>
      <div class="d-entity-body">
        <div><span class="pk">user_id</span> VARCHAR(64)</div>
        <div><span class="pk">query_text</span> VARCHAR(255)</div>
        <div>search_count INT</div>
        <div>last_searched_at BIGINT</div>
        <div>click_through_rate DECIMAL(3,2)</div>
        <div><span class="fk">user_id</span> &#8594; users.id</div>
        <div><span class="idx idx-btree">idx_user_last_searched</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header amber">trending_queries</div>
      <div class="d-entity-body">
        <div><span class="pk">query_text</span> VARCHAR(255)</div>
        <div><span class="pk">window_start</span> BIGINT</div>
        <div>current_count BIGINT</div>
        <div>baseline_avg BIGINT</div>
        <div>baseline_stddev DECIMAL(10,2)</div>
        <div>z_score DECIMAL(5,2)</div>
        <div>is_trending BOOLEAN</div>
        <div>detected_at TIMESTAMP</div>
        <div><span class="idx idx-btree">idx_z_score</span></div>
        <div><span class="idx idx-btree">idx_trending_detected</span></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-personalization",
		Title:       "Personalization &amp; A/B Testing",
		Description: "User history weighting, collaborative filtering signals, and A/B test framework for ranking",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Layer 1: User History Weighting</div>
        <div class="d-flow-v">
          <div class="d-box blue">Recent searches (last 30 days)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">Compute user_affinity per category<br/>affinity = &#931;(recency_weight &#215; frequency)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Redis Hash<br/>user:{uid}:affinities<br/>{"sports": 0.8, "cooking": 0.6, ...}<br/>TTL 24h</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Layer 2: Collaborative Filtering</div>
        <div class="d-flow-v">
          <div class="d-box purple">Users who searched X also searched Y</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Item-item similarity matrix<br/>(precomputed offline, Spark ALS)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Top-50 similar queries per query<br/>stored in DynamoDB</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Boost co-occurring suggestions<br/>+0.2 score for collaborative match</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Layer 3: A/B Test Framework</div>
        <div class="d-flow-v">
          <div class="d-box amber">Experiment Config<br/>{"id": "exp-42", "param": "trending_boost",<br/> "control": 1.0, "treatment": 2.0,<br/> "traffic": 10%}</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Hash(user_id + exp_id) % 100<br/>&lt; 10 &#8594; treatment, else control</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Metrics: suggestion CTR, search success rate, time-to-click</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Significance test after 7 days<br/>p &lt; 0.05 &#8594; ship or revert</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Personalization Merge at Query Time</div>
    <div class="d-flow">
      <div class="d-box indigo">Trie Top-K<br/>(global)</div>
      <div class="d-arrow">&#8594; merge</div>
      <div class="d-box indigo">User History<br/>(Redis lookup)</div>
      <div class="d-arrow">&#8594; re-rank</div>
      <div class="d-box indigo">Collaborative Boost<br/>(if cache hit)</div>
      <div class="d-arrow">&#8594; apply</div>
      <div class="d-box indigo">A/B Params<br/>(experiment config)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Final Top-10<br/>Suggestions</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-scaling",
		Title:       "Trie Sharding &amp; Replication Strategy",
		Description: "Sharding trie by prefix range with replication and multi-layer caching",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Prefix-Range Sharding (4 shards)</div>
    <div class="d-flow">
      <div class="d-box blue">Shard 1<br/>prefixes a&#8211;f<br/>&#8776; 7.5 GB<br/>3 replicas</div>
      <div class="d-box green">Shard 2<br/>prefixes g&#8211;m<br/>&#8776; 8.2 GB<br/>3 replicas</div>
      <div class="d-box purple">Shard 3<br/>prefixes n&#8211;s<br/>&#8776; 7.8 GB<br/>3 replicas</div>
      <div class="d-box amber">Shard 4<br/>prefixes t&#8211;z<br/>&#8776; 6.5 GB<br/>3 replicas</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Routing Layer</div>
        <div class="d-flow-v">
          <div class="d-box indigo">ALB + prefix-based routing<br/>prefix[0] &#8594; shard assignment</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Each shard: 3 nodes (1 primary + 2 replicas)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Read from any replica (round-robin)</div>
          <div class="d-box gray">Total: 4 shards &#215; 3 replicas = 12 nodes<br/>Each node: r5.4xlarge (128 GB RAM, 16 vCPU)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Multi-Layer Cache</div>
        <div class="d-flow-v">
          <div class="d-box green">L1: Browser localStorage<br/>recent queries, TTL 1h<br/>hit rate &#8776; 15%</div>
          <div class="d-arrow-down">&#8595; miss</div>
          <div class="d-box green">L2: CDN Edge (CloudFront)<br/>popular 2-char prefixes<br/>hit rate &#8776; 40%</div>
          <div class="d-arrow-down">&#8595; miss</div>
          <div class="d-box green">L3: Redis Cluster<br/>hot 3-char prefixes<br/>hit rate &#8776; 30%</div>
          <div class="d-arrow-down">&#8595; miss (&#8776; 15% of requests)</div>
          <div class="d-box amber">L4: In-Memory Trie<br/>always available</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Rebalancing</div>
        <div class="d-flow-v">
          <div class="d-box red">Hot shard detected<br/>(e.g., shard 4 during "Taylor Swift" trending)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Option 1: Add replicas to hot shard<br/>3 &#8594; 5 replicas (auto-scale)</div>
          <div class="d-box gray">Option 2: Split shard<br/>t&#8211;z &#8594; (t&#8211;v) + (w&#8211;z)</div>
          <div class="d-box gray">Option 3: CDN cache longer TTL<br/>for trending prefix range</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-failure-modes",
		Title:       "Failure Modes &amp; Recovery",
		Description: "Trie corruption recovery, stale suggestions handling, and offensive content filtering",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Trie Corruption / Bad Deploy</div>
      <div class="d-flow-v">
        <div class="d-box red">Symptom: p99 spike or empty suggestions</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Detection: health check queries<br/>"the", "how", "what" must return &gt; 0 results</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Recovery: atomic rollback to trie-v{N-1}<br/>pointer swap in &lt; 1 second</div>
        <div class="d-box green">Prevention: canary deploy to 5% of nodes<br/>monitor for 10 min before full rollout</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Stale Suggestions</div>
      <div class="d-flow-v">
        <div class="d-box red">Symptom: trending event not reflected<br/>(e.g., breaking news missing)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Root cause: batch pipeline delay &gt; 1h</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Mitigation: real-time trending path<br/>Flink &#8594; Redis &#8594; hot-patch into trie<br/>bypasses batch, &#8804; 15 min latency</div>
        <div class="d-box green">Fallback: if Flink down, trending<br/>queries served from Redis cache (stale OK)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Offensive Content Filtering</div>
      <div class="d-flow-v">
        <div class="d-box red">Risk: autocomplete suggests harmful/offensive queries</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Layer 1: Static blocklist<br/>&#8776; 500K terms, checked at trie build time<br/>blocked queries never enter trie</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Layer 2: ML classifier<br/>BERT-tiny model, &lt; 2ms inference<br/>score &gt; 0.8 &#8594; suppress suggestion</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Layer 3: Human review queue<br/>flagged borderline cases (score 0.5&#8211;0.8)<br/>reviewed within 4 hours</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Emergency: Redis blocklist<br/>instant removal, propagates in &lt; 1 min</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "sa-typeahead-flow",
		Title:       "End-to-End Typeahead Flow",
		Description: "Complete keystroke flow: debounce, cache check, service call, rank, and render",
		ContentFile: "problems/search-autocomplete",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Keystroke &#8594; Suggestions (target: &lt; 100ms e2e)</div>
    <div class="d-flow-v">
      <div class="d-flow">
        <div class="d-box blue">User types 'h'</div>
        <div class="d-arrow">&#8594;</div>
        <div class="d-box blue">Debounce 50ms<br/>(skip if typing fast)</div>
        <div class="d-arrow">&#8594;</div>
        <div class="d-box green">Check localStorage<br/>key: "suggest:h"<br/>TTL 1h</div>
      </div>
      <div class="d-arrow-down">&#8595; cache miss</div>
      <div class="d-flow">
        <div class="d-box green">CDN Edge<br/>(CloudFront POP)<br/>&#8776; 5ms</div>
        <div class="d-arrow">&#8594; miss</div>
        <div class="d-box purple">ALB<br/>route: prefix[0] &#8594; shard<br/>&#8776; 1ms</div>
        <div class="d-arrow">&#8594;</div>
        <div class="d-box indigo">Suggestion Service<br/>trie.lookup("h")<br/>&#8776; 0.5ms</div>
      </div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-flow">
        <div class="d-box indigo">Personalize<br/>Redis: user affinities<br/>&#8776; 1ms</div>
        <div class="d-arrow">&#8594;</div>
        <div class="d-box indigo">Re-rank Top-K<br/>apply scoring formula<br/>&#8776; 0.2ms</div>
        <div class="d-arrow">&#8594;</div>
        <div class="d-box indigo">Filter Offensive<br/>blocklist check<br/>&#8776; 0.1ms</div>
        <div class="d-arrow">&#8594;</div>
        <div class="d-box green">Return JSON<br/>10 suggestions<br/>&#8776; 500 bytes</div>
      </div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-flow">
        <div class="d-box blue">Browser Render<br/>update dropdown<br/>&#8776; 5ms</div>
        <div class="d-arrow">&#8594;</div>
        <div class="d-box blue">Cache in localStorage<br/>for repeat prefix</div>
        <div class="d-arrow">&#8594;</div>
        <div class="d-box gray">Log impression<br/>async to Kafka</div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Latency Budget Breakdown</div>
    <div class="d-flow">
      <div class="d-box amber">Debounce<br/>50ms</div>
      <div class="d-box amber">Network<br/>&#8776; 20ms</div>
      <div class="d-box amber">CDN check<br/>&#8776; 5ms</div>
      <div class="d-box amber">Trie lookup<br/>&#8776; 0.5ms</div>
      <div class="d-box amber">Personalize<br/>&#8776; 1ms</div>
      <div class="d-box amber">Re-rank<br/>&#8776; 0.3ms</div>
      <div class="d-box amber">Network back<br/>&#8776; 20ms</div>
      <div class="d-box amber">Render<br/>&#8776; 5ms</div>
      <div class="d-box green"><strong>Total<br/>&#8776; 102ms</strong></div>
    </div>
  </div>
</div>`,
	})
}

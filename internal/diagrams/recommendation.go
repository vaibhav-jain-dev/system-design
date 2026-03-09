package diagrams

func registerRecommendationSystem(r *Registry) {
	r.Register(&Diagram{
		Slug:        "rec-requirements",
		Title:       "Requirements &amp; Scale",
		Description: "Scale targets, latency budget, and freshness requirements for a recommendation system",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">1B registered users</div>
        <div class="d-box purple">100M items (videos / products / articles)</div>
        <div class="d-box purple">10K recommendation requests/sec (peak)</div>
        <div class="d-box purple">Model training: daily batch (offline)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Quality Goals</div>
      <div class="d-flow-v">
        <div class="d-box green">CTR improvement: &gt;20% vs non-personalized baseline</div>
        <div class="d-box green">Real-time signal freshness: &lt;5 min</div>
        <div class="d-box green">Cold start: usable recs after 50 user interactions</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Serving Latency Budget</div>
      <div class="d-flow-v">
        <div class="d-box amber">P50 end-to-end: &lt;100ms</div>
        <div class="d-box amber">P99 end-to-end: &lt;200ms</div>
        <div class="d-box blue">Candidate generation: 30ms</div>
        <div class="d-box blue">Ranking: 50ms</div>
        <div class="d-box blue">Filtering + serving: 5ms</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Key Non-Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box gray">Diversity: not all recs from same category</div>
        <div class="d-box gray">Freshness: newly uploaded items surfaced within 1h</div>
        <div class="d-box gray">Explainability: &#8220;because you watched X&#8221; reasons</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-two-stage-architecture",
		Title:       "Two-Stage Architecture",
		Description: "Candidate generation, ranking, filtering pipeline with per-stage latency budgets",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Recommendation Pipeline</div>
    <div class="d-flow">
      <div class="d-box gray">100M items in corpus</div>
      <div class="d-arrow">&#8594; 30ms &#8594;</div>
      <div class="d-box blue">
        <strong>Stage 1: Candidate Generation</strong><br/>
        Retrieve 100&#8211;1,000 candidates<br/>
        ANN search + collaborative filtering<br/>
        Recall &gt; 90%
      </div>
      <div class="d-arrow">&#8594; 50ms &#8594;</div>
      <div class="d-box purple">
        <strong>Stage 2: Ranking</strong><br/>
        Score all candidates<br/>
        Deep learning model<br/>
        Predicted CTR + watch time
      </div>
      <div class="d-arrow">&#8594; 5ms &#8594;</div>
      <div class="d-box amber">
        <strong>Stage 3: Filtering</strong><br/>
        Remove watched/purchased<br/>
        Apply business rules<br/>
        Enforce diversity
      </div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">
        <strong>Serve</strong><br/>
        Top 10 to user<br/>
        &lt;100ms total
      </div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Why Two Stages?</div>
        <div class="d-flow-v">
          <div class="d-box green">Ranking all 100M items with deep model = impossible in &lt;100ms</div>
          <div class="d-box green">Candidate gen uses fast approximate methods to reduce to 1,000</div>
          <div class="d-box green">Ranker uses expensive model only on the 1,000 candidates</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Latency Trade-off</div>
        <div class="d-flow-v">
          <div class="d-box amber">Stage 1 (fast, lower precision): ANN with HNSW, recall ~90%</div>
          <div class="d-box amber">Stage 2 (slow, high precision): full deep model, 1,000 items</div>
          <div class="d-label">Recall loss at stage 1 is acceptable &#8212; we miss ~10% of truly best items but still serve highly relevant results.</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-candidate-generation",
		Title:       "Candidate Generation Methods",
		Description: "Four candidate retrieval strategies: collaborative filtering, content-based, trending, and personalized popularity",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">1. Collaborative Filtering</div>
        <div class="d-flow-v">
          <div class="d-box blue">&#8220;Users like you watched X&#8221;</div>
          <div class="d-box blue">Item&#8211;item similarity matrix from co-interaction</div>
          <div class="d-label">If user watched A and B, and others who watched A+B also watched C &#8594; recommend C</div>
          <div class="d-box blue">Matrix factorization (ALS) for implicit feedback</div>
          <div class="d-box green">Generates ~250 candidates</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">2. Content-Based Filtering</div>
        <div class="d-flow-v">
          <div class="d-box purple">Item embedding similarity</div>
          <div class="d-label">&#8220;You liked action movies, here are more action movies&#8221;</div>
          <div class="d-box purple">Item features: genre, tags, actors, description embedding</div>
          <div class="d-box purple">User profile: weighted average of liked item embeddings</div>
          <div class="d-box green">Generates ~250 candidates</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">3. Trending (Cold Start Fallback)</div>
        <div class="d-flow-v">
          <div class="d-box amber">Globally popular items by region</div>
          <div class="d-label">Sorted by engagement velocity (interactions in last 24h / age)</div>
          <div class="d-box amber">Refreshed every 15 minutes</div>
          <div class="d-box amber">Used when user has &lt;50 interactions</div>
          <div class="d-box green">Generates ~250 candidates</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">4. Personalized Popularity</div>
        <div class="d-flow-v">
          <div class="d-box indigo">Trending filtered by user&#8217;s category affinity</div>
          <div class="d-label">Top trending in categories user has engaged with in last 30 days</div>
          <div class="d-box indigo">Blends novelty (trending) with relevance (affinity)</div>
          <div class="d-box green">Generates ~250 candidates</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Union &amp; Dedup</div>
    <div class="d-flow">
      <div class="d-box gray">250 + 250 + 250 + 250 = 1,000 candidates (after dedup by item_id)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Pass all 1,000 to ranking stage</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-ml-pipeline",
		Title:       "ML Training Pipeline",
		Description: "End-to-end model training: feature engineering, feature store, training, registry, and A/B promotion",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Training Pipeline (Daily Batch)</div>
    <div class="d-flow">
      <div class="d-box gray">Raw user events (Kafka / S3)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Feature Engineering (Spark)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Feature Store (offline: S3, online: Redis)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Model Training (TensorFlow / PyTorch)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Model Registry</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box indigo">A/B Test &#8594; Production Serving</div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Offline Model (30-day history)</div>
        <div class="d-flow-v">
          <div class="d-box blue">Training data: last 30 days of click/watch events</div>
          <div class="d-box blue">Labels: positive = watched &gt;30% of video</div>
          <div class="d-box blue">Training time: 2&#8211;4 hours on GPU cluster</div>
          <div class="d-box blue">Deployed daily via model registry promotion</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Online Fine-Tuning (24h signals)</div>
        <div class="d-flow-v">
          <div class="d-box amber">Incremental training on last 24h data</div>
          <div class="d-box amber">Adapts to trending topics and seasonal shifts</div>
          <div class="d-box amber">Runs every 6 hours, updates model weights</div>
          <div class="d-box amber">Falls back to offline model if fine-tune diverges</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Model Registry</div>
        <div class="d-flow-v">
          <div class="d-box indigo">Versioned model artifacts in S3</div>
          <div class="d-box indigo">Metadata: train date, eval metrics, dataset hash</div>
          <div class="d-box indigo">Promotion requires: offline AUC &gt; 0.85 + online A/B CTR lift</div>
          <div class="d-box indigo">Rollback in &lt; 5 minutes via blue/green swap</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-feature-store",
		Title:       "Feature Store Architecture",
		Description: "Offline (S3+Spark) and online (Redis) feature stores with freshness guarantees at serving time",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Offline Feature Store (S3 + Spark)</div>
      <div class="d-flow-v">
        <div class="d-box blue">Batch features &#8212; computed daily by Spark jobs</div>
        <div class="d-flow-v" style="margin-top:4px">
          <div class="d-box gray">user_lifetime_watch_time (minutes)</div>
          <div class="d-box gray">category_preferences_30d (vector of category weights)</div>
          <div class="d-box gray">item_watch_count_30d (popularity signal)</div>
          <div class="d-box gray">user_embedding_128d (dense representation)</div>
          <div class="d-box gray">item_embedding_128d (dense representation)</div>
        </div>
        <div class="d-box blue">Freshness: updated once daily</div>
        <div class="d-label">Storage: Parquet on S3, loaded into Spark for training.</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Online Feature Store (Redis)</div>
      <div class="d-flow-v">
        <div class="d-box red">Real-time features &#8212; updated continuously</div>
        <div class="d-flow-v" style="margin-top:4px">
          <div class="d-box gray">user_session_clicks (last 30 min)</div>
          <div class="d-box gray">trending_now (top 1K items, 15-min TTL)</div>
          <div class="d-box gray">item_clicks_1h (per-item recency signal)</div>
          <div class="d-box gray">user_category_affinity_delta (session drift)</div>
        </div>
        <div class="d-box red">Freshness: &lt;5 minutes</div>
        <div class="d-label">Written by Flink stream processor consuming Kafka user-events.</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Feature Join at Serving Time</div>
      <div class="d-flow-v">
        <div class="d-box purple">Serving request arrives with user_id + item_ids</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Lookup offline features &#8212; Redis hash: offline:{user_id}</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">Lookup online features &#8212; Redis hash: online:{user_id}</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Concatenate into full feature vector &#8594; ranking model</div>
      </div>
      <div class="d-label">Both lookups in parallel. Total feature fetch: &lt;5ms (Redis pipeline).</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-embedding-retrieval",
		Title:       "Embedding-Based Retrieval (ANN)",
		Description: "Item and user embedding generation, FAISS ANN search, and HNSW index structure",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Embedding Generation</div>
    <div class="d-flow">
      <div class="d-box gray">Item metadata (title, tags, category)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Item embedding model (two-tower)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">128-dim item vector</div>
    </div>
    <div class="d-flow" style="margin-top:8px">
      <div class="d-box gray">User interaction history</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">User embedding model (two-tower)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">128-dim user vector</div>
    </div>
    <div class="d-label">Two-tower: separate encoders for user and item, trained to maximize dot product for positive interactions.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">ANN Search (FAISS / Elasticsearch kNN)</div>
        <div class="d-flow-v">
          <div class="d-box amber">Find 100 nearest item vectors to user vector</div>
          <div class="d-box amber">Similarity: cosine distance (dot product after L2 normalize)</div>
          <div class="d-box amber">FAISS IVF index: cluster 100M items into 10K centroids</div>
          <div class="d-label">At query time: find top-100 nearest centroids, search only those clusters &#8212; scans ~1% of corpus.</div>
          <div class="d-box green">Recall: ~90% of true top-100 found</div>
          <div class="d-box green">Latency: &lt;50ms at 100M items</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">HNSW Index Structure</div>
        <div class="d-flow-v">
          <div class="d-box indigo">Hierarchical Navigable Small World graph</div>
          <div class="d-box indigo">Multi-layer graph: top layer sparse, bottom dense</div>
          <div class="d-box indigo">Search: start at top layer, greedily descend</div>
          <div class="d-box indigo">Complexity: O(log N) per query vs O(N) brute force</div>
          <div class="d-label">Build time: O(N log N). Memory: ~80 bytes/vector &#8215; 100M = 8GB for index.</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-ranking-model",
		Title:       "Ranking Model (Deep Learning)",
		Description: "Wide and deep network architecture, feature inputs, output scores, and batch candidate scoring",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Input Features</div>
        <div class="d-flow-v">
          <div class="d-box blue">User features: embedding_128d, age_bucket, region, device</div>
          <div class="d-box blue">Item features: embedding_128d, category, age_hours, popularity_score</div>
          <div class="d-box blue">Context features: time_of_day, day_of_week, session_length_minutes</div>
          <div class="d-box blue">Interaction features: user&#215;item cross features, historical CTR for user+category</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Wide &amp; Deep Architecture</div>
        <div class="d-flow-v">
          <div class="d-box purple"><strong>Wide part:</strong> linear model on crossed features (memorization &#8212; learns specific user&#215;item patterns)</div>
          <div class="d-box purple"><strong>Deep part:</strong> 3-layer MLP on embeddings (generalization &#8212; learns abstract patterns)</div>
          <div class="d-box purple">Outputs combined: sigmoid(wide + deep) &#8594; predicted CTR</div>
          <div class="d-label">Wide part handles long-tail patterns. Deep part generalizes to unseen combinations.</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Model Outputs &amp; Final Score</div>
    <div class="d-flow">
      <div class="d-box amber">Predicted CTR (click-through rate)</div>
      <div class="d-box gray">&#215; 0.3</div>
      <div class="d-box amber">+ Predicted watch_time (seconds)</div>
      <div class="d-box gray">&#215; 0.7</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green"><strong>Final score</strong><br/>0.3&#215;CTR + 0.7&#215;watch_time</div>
    </div>
    <div class="d-label">Watch time weighted higher &#8212; optimizing for engagement depth, not just clicks. Business objective encoded in weights.</div>
  </div>
  <div class="d-group" style="margin-top:8px">
    <div class="d-group-title">Batch Scoring All 1,000 Candidates</div>
    <div class="d-flow">
      <div class="d-box indigo">1,000 candidates &#8594; batched into single model forward pass (GPU)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">All 1,000 scored in ~50ms</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Sort by final score &#8594; top 10 returned</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-cold-start",
		Title:       "Cold Start Problem",
		Description: "Strategies for new users (no history) and new items (no engagement), with resolution thresholds",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">New User Cold Start</div>
      <div class="d-flow-v">
        <div class="d-box red"><strong>Problem:</strong> No interaction history &#8594; no user embedding &#8594; collaborative filtering useless</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue"><strong>Step 1:</strong> Show trending items in user&#8217;s region (geolocation from signup)</div>
        <div class="d-box blue"><strong>Step 2:</strong> Onboarding preference prompt &#8212; 5 category selections</div>
        <div class="d-label">&#8220;What are you interested in? Comedy / Action / Tech / Sports / Food&#8221;</div>
        <div class="d-box blue"><strong>Step 3:</strong> Demographic-based recs from age_bucket + region cluster</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><strong>Resolved after:</strong> 50 user interactions &#8594; sufficient for embedding, switch to personalized recs</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">New Item Cold Start</div>
      <div class="d-flow-v">
        <div class="d-box red"><strong>Problem:</strong> No engagement data &#8594; no item popularity signal &#8594; won&#8217;t appear in collaborative filtering</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber"><strong>Step 1:</strong> Show to random 1% of users whose category affinity matches item&#8217;s category</div>
        <div class="d-label">Exploration: targeted random exposure to collect initial signals.</div>
        <div class="d-box amber"><strong>Step 2:</strong> Measure initial CTR over first 100 impressions</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><strong>If CTR &gt;5%:</strong> add item to regular candidate pool &#8594; collaborative filtering picks it up</div>
        <div class="d-box gray"><strong>If CTR &lt;1%:</strong> deprioritize, flag for content review</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px">
      <div class="d-group-title">Resolution Thresholds</div>
      <div class="d-flow-v">
        <div class="d-box green">New user: resolved after <strong>50 interactions</strong></div>
        <div class="d-box green">New item: resolved after <strong>100 impressions</strong></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-real-time-signals",
		Title:       "Real-Time Signal Processing",
		Description: "Kafka event stream, Flink consumer, sliding window feature updates, and session-level recommendation adaptation",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Event Ingestion</div>
    <div class="d-flow">
      <div class="d-box gray">User watches 45s of video</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Kafka event: &#123; user_id, item_id, event_type=watch, watch_time_seconds=45, timestamp &#125;</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Flink consumer (partitioned by user_id)</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Flink Stream Processing</div>
    <div class="d-flow">
      <div class="d-box blue">Flink sliding window job (1h window, 5min slide)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-flow-v">
        <div class="d-box green">Update: user_session_clicks (INCR)</div>
        <div class="d-box green">Update: user_watch_time_1h (SUM)</div>
        <div class="d-box green">Update: category_affinity_delta (weighted INCR by category)</div>
        <div class="d-box green">Update: item_clicks_1h (INCR per item_id)</div>
      </div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">Write to Redis online feature store</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Impact on Recommendations</div>
    <div class="d-flow">
      <div class="d-box amber">Next recommendation request (within 5 min)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Feature fetch includes updated session features</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Recs adapt to current session intent</div>
    </div>
    <div class="d-label">Example: user watches 3 cooking videos &#8594; category_affinity_delta[cooking] increases &#8594; next rec page shows more cooking content within 5 minutes.</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-ab-testing",
		Title:       "A/B Testing Framework",
		Description: "Deterministic user assignment, metric collection, statistical significance, and winner promotion",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Experiment Configuration</div>
    <div class="d-flow">
      <div class="d-box purple">Experiment: model_v2_vs_v1</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">50% users &#8594; Model A (control: current production)</div>
      <div class="d-box green">50% users &#8594; Model B (treatment: new model)</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Deterministic Assignment</div>
        <div class="d-flow-v">
          <div class="d-box indigo">bucket = hash(user_id + experiment_id) mod 100</div>
          <div class="d-box indigo">bucket 0&#8211;49 &#8594; control (Model A)</div>
          <div class="d-box indigo">bucket 50&#8211;99 &#8594; treatment (Model B)</div>
          <div class="d-label">Same user always gets same model &#8212; consistent experience, no cross-contamination. Hash ensures uniform distribution.</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Metrics &amp; Decision</div>
        <div class="d-flow-v">
          <div class="d-box amber">Primary: CTR (click-through rate)</div>
          <div class="d-box amber">Primary: avg watch_time per session</div>
          <div class="d-box amber">Secondary: session_length, user_retention_7d</div>
          <div class="d-box amber">Guardrails: latency P99 must not increase &gt;10ms</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green"><strong>Winner criteria:</strong> p-value &lt;0.05 on primary metric</div>
          <div class="d-box green"><strong>Minimum duration:</strong> 2 weeks (captures weekly seasonality)</div>
          <div class="d-box green"><strong>Min sample size:</strong> 10K users per bucket for 95% confidence</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Promotion to 100% Traffic</div>
    <div class="d-flow">
      <div class="d-box green">Experiment declares winner</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Promote Model B to model registry (stable tag)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Gradual rollout: 50% &#8594; 75% &#8594; 100% over 24h</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Monitor guardrail metrics during ramp</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-serving-optimization",
		Title:       "Serving Optimization",
		Description: "Pre-computed recommendations for top active users, on-demand for long tail, with cache invalidation strategy",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Pre-Computed Recs (Top 10M Active Users)</div>
        <div class="d-flow-v">
          <div class="d-box green">Offline job runs daily after model training</div>
          <div class="d-box green">Compute top-K recs for each of 10M most-active users</div>
          <div class="d-box green">Store in Redis: ZADD recs:{user_id} {score} {item_id}</div>
          <div class="d-label">K=50 pre-computed, serve top 10, rotate as user scrolls.</div>
          <div class="d-box green">Serve at request time: ZREVRANGE recs:{user_id} 0 9 &#8594; &lt;1ms</div>
          <div class="d-box blue">Covers <strong>80% of traffic</strong> (top 1% of users drive 80% of requests)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">On-Demand (Long-Tail Users)</div>
        <div class="d-flow-v">
          <div class="d-box amber">No pre-computed recs for user &#8594; compute on demand</div>
          <div class="d-box amber">Candidate gen + ranking + filtering pipeline</div>
          <div class="d-box amber">Target: &lt;100ms end-to-end</div>
          <div class="d-box amber">Covers remaining <strong>20% of traffic</strong></div>
          <div class="d-label">Long-tail users are infrequent &#8212; on-demand cost is acceptable and avoids wasted pre-compute for rarely-active users.</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Cache Invalidation Strategy</div>
    <div class="d-flow">
      <div class="d-box red">Invalidate pre-computed recs when:</div>
      <div class="d-flow-v">
        <div class="d-box gray">Model updated (daily) &#8594; recompute all 10M users over 2h</div>
        <div class="d-box gray">User has 5+ new interactions &#8594; mark stale &#8594; recompute async</div>
        <div class="d-box gray">TTL: 25 hours (auto-expire if not refreshed)</div>
      </div>
    </div>
    <div class="d-label">Recompute rate: 10M users / 7,200 sec (2h window) = 1,389 recomputes/sec. Single offline job at low priority.</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rec-monitoring",
		Title:       "Monitoring &amp; Business Metrics",
		Description: "Technical and business KPIs, alert thresholds, and A/B test dashboard for recommendation health",
		ContentFile: "problems/recommendation-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Technical Metrics</div>
      <div class="d-flow-v">
        <div class="d-box blue">Latency P50 / P99 per stage (gen / rank / filter / total)</div>
        <div class="d-box blue">Cache hit rate for pre-computed recs (target &gt;75%)</div>
        <div class="d-box blue">Model staleness: hours since last training run</div>
        <div class="d-box blue">Feature store freshness: age of online features (target &lt;5min)</div>
        <div class="d-box blue">Candidate generation recall (sampled offline evaluation)</div>
        <div class="d-box blue">Kafka consumer lag for real-time signal processing</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Business Metrics</div>
      <div class="d-flow-v">
        <div class="d-box green">CTR &#8212; click-through rate (target: &gt;8%)</div>
        <div class="d-box green">Watch time per session (minutes)</div>
        <div class="d-box green">Recommendation diversity &#8212; % unique categories in top-10</div>
        <div class="d-box green">Filter rate &#8212; % of candidates removed by post-ranking filters</div>
        <div class="d-label">Filter rate &gt;50% means candidate gen is poor quality &#8212; tuning signal.</div>
        <div class="d-box green">User retention D7 / D30 (weekly/monthly active ratio)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Alerts &amp; A/B Dashboard</div>
      <div class="d-flow-v">
        <div class="d-box red"><strong>Alert:</strong> CTR drops &gt;10% week-over-week &#8594; PagerDuty (model degradation)</div>
        <div class="d-box red"><strong>Alert:</strong> P99 latency &gt;300ms for 5 min &#8594; Slack warning</div>
        <div class="d-box red"><strong>Alert:</strong> Feature freshness &gt;15 min &#8594; Flink lag investigation</div>
        <div class="d-box amber"><strong>A/B dashboard:</strong> Live CTR lift per experiment bucket</div>
        <div class="d-box amber">Statistical significance indicator (p-value live)</div>
        <div class="d-box amber">Guardrail metric tracker (latency, diversity, filter rate)</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

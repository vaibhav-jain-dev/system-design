package diagrams

func registerTwitterFeed(r *Registry) {
	r.Register(&Diagram{
		Slug:        "tf-requirements",
		Title:       "Functional & Non-Functional Requirements",
		Description: "Scale targets and feature requirements for a Twitter/X feed system",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P0 &#8212; Core (Must Have)</div>
      <div class="d-flow-v">
        <div class="d-box green">Post tweets (text, images, links)</div>
        <div class="d-box green">Home timeline &#8212; aggregated feed from followed users</div>
        <div class="d-box green">Follow / unfollow users</div>
        <div class="d-box green">User timeline &#8212; single user&#8217;s tweets</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P1 &#8212; Important</div>
      <div class="d-flow-v">
        <div class="d-box blue">Like, retweet, reply</div>
        <div class="d-box blue">Hashtag search</div>
        <div class="d-box blue">Trending topics</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P2 &#8212; Nice to Have</div>
      <div class="d-flow-v">
        <div class="d-box gray">Notifications (mentions, likes)</div>
        <div class="d-box gray">Media upload (video, GIF)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">500M DAU</div>
        <div class="d-box purple">500M tweets/day &#8776; 5,800 writes/sec</div>
        <div class="d-box purple">300K timeline reads/sec (peak)</div>
        <div class="d-box purple">Average user follows 200 accounts</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box amber">Timeline latency: &lt; 200ms p99</div>
        <div class="d-box amber">Tweet publish: &lt; 5s to appear in follower feeds</div>
        <div class="d-box amber">Availability: 99.99%</div>
        <div class="d-box amber">Timeline freshness: &lt; 30s for non-celebrity tweets</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Key Decisions</div>
      <div class="d-flow-v">
        <div class="d-box red">Fan-out-on-write vs fan-out-on-read?</div>
        <div class="d-label">Hybrid: write for normal users, read for celebrities</div>
        <div class="d-box red">Chronological vs ranked feed?</div>
        <div class="d-label">Ranked with recency bias &#8212; engagement + freshness score</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-api-design",
		Title:       "API Design",
		Description: "Core REST API endpoints for tweet creation, timelines, and social graph",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Tweet Operations</div>
      <div class="d-flow-v">
        <div class="d-box green"><strong>POST /v1/tweets</strong><br/>Body: {text, media_ids[], reply_to?}<br/>Returns: {tweet_id, created_at}</div>
        <div class="d-box green"><strong>DELETE /v1/tweets/{id}</strong><br/>Soft-delete &#8212; marks as deleted, removes from timelines</div>
        <div class="d-box blue"><strong>POST /v1/tweets/{id}/like</strong><br/>Idempotent &#8212; returns 200 if already liked</div>
        <div class="d-box blue"><strong>POST /v1/tweets/{id}/retweet</strong><br/>Creates retweet entry, fans out to followers</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Timeline Endpoints</div>
      <div class="d-flow-v">
        <div class="d-box purple"><strong>GET /v1/timeline/home</strong><br/>Params: cursor, limit (default 20)<br/>Returns: {tweets[], next_cursor}</div>
        <div class="d-box purple"><strong>GET /v1/timeline/user/{id}</strong><br/>Params: cursor, limit<br/>Returns: user&#8217;s own tweets + retweets</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Social Graph</div>
      <div class="d-flow-v">
        <div class="d-box amber"><strong>POST /v1/users/{id}/follow</strong><br/>Triggers async fan-out of recent tweets</div>
        <div class="d-box amber"><strong>DELETE /v1/users/{id}/follow</strong><br/>Unfollow &#8212; removes tweets from timeline cache</div>
        <div class="d-box gray"><strong>GET /v1/users/{id}/followers</strong><br/>Paginated &#8212; cursor-based, 100 per page</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-fanout-strategies",
		Title:       "Fan-out Strategies Comparison",
		Description: "Fan-out-on-write vs fan-out-on-read with hybrid approach for celebrities",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fan-out-on-Write (Push)</div>
      <div class="d-flow-v">
        <div class="d-box green">User A posts tweet</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Fan-out service reads A&#8217;s follower list</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Write tweet_id to each follower&#8217;s timeline cache</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">1,000 followers = 1,000 Redis ZADD ops</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Timeline read = single Redis ZREVRANGE &#8212; O(1)</div>
      </div>
      <div class="d-flow-v" style="margin-top:8px">
        <div class="d-box amber">&#10003; Fast reads (&lt; 5ms)</div>
        <div class="d-box red">&#10007; Celebrity problem: 100M writes per tweet</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fan-out-on-Read (Pull)</div>
      <div class="d-flow-v">
        <div class="d-box green">User B opens home timeline</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Fetch B&#8217;s follow list (200 accounts)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Query each followed user&#8217;s recent tweets</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Merge + rank 200 tweet lists</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red">Slow reads: 200 scatter-gather queries &#8776; 100-300ms</div>
      </div>
      <div class="d-flow-v" style="margin-top:8px">
        <div class="d-box amber">&#10003; No write amplification</div>
        <div class="d-box red">&#10007; High read latency at scale</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Hybrid (Best of Both)</div>
      <div class="d-flow-v">
        <div class="d-box green">Normal user (&lt; 10K followers)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Fan-out-on-write &#8212; push to all followers</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo"><strong>Celebrity (&gt; 10K followers)</strong></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Skip fan-out &#8212; store in celebrity tweet cache</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">On read: merge precomputed timeline + celebrity tweets</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Result: fast reads + bounded write cost</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end architecture: CDN, load balancers, services, caches, and storage",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box gray">Client (Mobile / Web)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">CDN (CloudFront)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">ALB</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">API Gateway</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Write Path</div>
        <div class="d-flow-v">
          <div class="d-box green">Tweet Service</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Kafka (tweet-events topic)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-flow">
            <div class="d-box amber">Fan-out Service</div>
            <div class="d-arrow">&#8594;</div>
            <div class="d-box red">Redis Timeline Cache</div>
          </div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Tweet Store (DynamoDB / Vitess)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Read Path</div>
        <div class="d-flow-v">
          <div class="d-box green">Timeline Service</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red">Redis Timeline Cache</div>
          <div class="d-label">Precomputed sorted set per user (tweet_id:score)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Celebrity Tweet Cache</div>
          <div class="d-label">Merge on read for &gt; 10K follower accounts</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Ranking Service</div>
          <div class="d-label">Score = engagement + recency + relevance</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Tweet Store (hydrate tweet objects)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Supporting Services</div>
        <div class="d-flow-v">
          <div class="d-box blue">User Service</div>
          <div class="d-label">Profile, auth, follower counts</div>
          <div class="d-box blue">Social Graph Service</div>
          <div class="d-label">Follow/unfollow, adjacency lists</div>
          <div class="d-box blue">Search Service</div>
          <div class="d-label">Inverted index, hashtag lookup</div>
          <div class="d-box blue">Notification Service</div>
          <div class="d-label">Mentions, likes, retweets</div>
          <div class="d-box blue">Media Service</div>
          <div class="d-label">Image/video upload &#8594; S3 + CDN</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-timeline-generation",
		Title:       "Timeline Generation &#8212; Hybrid Approach",
		Description: "Precomputed Redis timeline merged with on-read celebrity tweets and ranking",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Step 1: Precomputed Timeline (Fan-out-on-Write)</div>
    <div class="d-flow">
      <div class="d-box green">Normal user posts tweet</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Fan-out Service</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">Redis ZADD timeline:{user_id} {score} {tweet_id}</div>
    </div>
    <div class="d-label">Score = tweet timestamp (epoch ms). Sorted set trimmed to 800 entries via ZREMRANGEBYRANK.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Step 2: Celebrity Merge (Fan-out-on-Read)</div>
    <div class="d-flow">
      <div class="d-box indigo">User follows 5 celebrities</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Fetch latest 20 tweets from each celebrity cache</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Merge 100 celebrity tweets with precomputed timeline</div>
    </div>
    <div class="d-label">Celebrity cache: separate Redis sorted set per celebrity. &#8776; 5 &#215; 20 = 100 tweets to merge.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Step 3: Rank &amp; Return</div>
    <div class="d-flow">
      <div class="d-box amber">Ranking model scores merged tweets</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Hydrate top 20 tweet objects from Tweet Store</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Return {tweets[], next_cursor}</div>
    </div>
    <div class="d-label">Total latency: Redis ZREVRANGE (2ms) + celebrity merge (10ms) + ranking (20ms) + hydration (50ms) &#8776; 80ms p50.</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-data-model",
		Title:       "Data Model",
		Description: "Core tables: users, tweets, follows, timelines, likes, retweets, hashtags",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">users</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> user_id (snowflake)</div>
        <div>username (unique)</div>
        <div>display_name</div>
        <div>bio</div>
        <div>follower_count</div>
        <div>following_count</div>
        <div>is_celebrity (bool)</div>
        <div>created_at</div>
        <div><span class="idx idx-hash">IDX</span> username</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header green">tweets</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> tweet_id (snowflake)</div>
        <div><span class="fk">FK</span> user_id</div>
        <div>text (280 chars)</div>
        <div>media_urls[]</div>
        <div>reply_to_id (nullable)</div>
        <div>retweet_of_id (nullable)</div>
        <div>like_count</div>
        <div>retweet_count</div>
        <div>reply_count</div>
        <div>created_at</div>
        <div>is_deleted (bool)</div>
        <div><span class="idx idx-btree">IDX</span> user_id + created_at</div>
        <div><span class="idx idx-btree">IDX</span> reply_to_id</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header purple">follows</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> follower_id + followee_id</div>
        <div><span class="fk">FK</span> follower_id &#8594; users</div>
        <div><span class="fk">FK</span> followee_id &#8594; users</div>
        <div>created_at</div>
        <div><span class="idx idx-btree">IDX</span> followee_id (reverse lookup)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header red">timelines (Redis)</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> timeline:{user_id}</div>
        <div>Sorted set: member=tweet_id, score=timestamp</div>
        <div>Max 800 entries per user</div>
        <div>TTL: 7 days for inactive users</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header amber">likes</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> user_id + tweet_id</div>
        <div>created_at</div>
        <div><span class="idx idx-btree">IDX</span> tweet_id (count query)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header indigo">retweets</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> user_id + tweet_id</div>
        <div>retweet_tweet_id (the new tweet entry)</div>
        <div>created_at</div>
        <div><span class="idx idx-btree">IDX</span> tweet_id</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header gray">hashtags</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> hashtag + tweet_id</div>
        <div>hashtag (lowercase, normalized)</div>
        <div><span class="fk">FK</span> tweet_id &#8594; tweets</div>
        <div>created_at</div>
        <div><span class="idx idx-hash">IDX</span> hashtag</div>
        <div><span class="idx idx-btree">IDX</span> hashtag + created_at</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-celebrity-problem",
		Title:       "The Celebrity Problem",
		Description: "Why fan-out-on-write fails for celebrities and how the hybrid approach solves it",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">The Problem</div>
    <div class="d-flow">
      <div class="d-box indigo"><strong>Celebrity</strong><br/>100M followers</div>
      <div class="d-arrow">&#8594; posts 1 tweet &#8594;</div>
      <div class="d-box red"><strong>Fan-out-on-Write</strong><br/>100M Redis ZADD operations</div>
    </div>
    <div class="d-flow-v" style="margin-top:8px">
      <div class="d-box red">At 100K ZADD/sec per Redis node = 1,000 seconds to complete</div>
      <div class="d-box red">10 shards &#8594; still 100 seconds latency</div>
      <div class="d-box red">Celebrity tweets 5x/day = 500M unnecessary writes/day</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">The Solution: Hybrid Fan-out</div>
    <div class="d-cols">
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box green"><strong>Normal User</strong> (&lt; 10K followers)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">Fan-out-on-write as usual</div>
          <div class="d-label">1,000 followers &#215; 1 ZADD = 1,000 ops &#8776; 10ms</div>
        </div>
      </div>
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box indigo"><strong>Celebrity</strong> (&gt; 10K followers)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Store in celebrity tweet cache only</div>
          <div class="d-label">1 write regardless of follower count</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Read-Time Merge</div>
    <div class="d-flow">
      <div class="d-box red">User&#8217;s precomputed timeline</div>
      <div class="d-box amber">+ merge celebrity tweets (&#8776; 5-10 celebrities followed)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Merged feed in &lt; 15ms extra latency</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-ranking-algorithm",
		Title:       "Feed Ranking Algorithm",
		Description: "Multi-signal ranking: engagement, recency, relevance, and social graph proximity",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Ranking Score Formula</div>
    <div class="d-box purple"><strong>score = w1&#183;engagement + w2&#183;recency + w3&#183;relevance + w4&#183;social_proximity</strong></div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Engagement Score (w1 = 0.3)</div>
        <div class="d-flow-v">
          <div class="d-box green">likes &#215; 1.0</div>
          <div class="d-box green">retweets &#215; 2.0</div>
          <div class="d-box green">replies &#215; 3.0</div>
          <div class="d-box green">click-through &#215; 1.5</div>
          <div class="d-label">Normalized to [0, 1] by percentile rank</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Recency Score (w2 = 0.35)</div>
        <div class="d-flow-v">
          <div class="d-box blue">decay = 1 / (1 + age_hours / 6)</div>
          <div class="d-label">Half-life of 6 hours &#8212; tweets older than 24h score &lt; 0.2</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Relevance Score (w3 = 0.2)</div>
        <div class="d-flow-v">
          <div class="d-box amber">Topic interest from user history</div>
          <div class="d-box amber">Hashtag affinity (past 30 days)</div>
          <div class="d-box amber">Content type preference (images vs text)</div>
          <div class="d-label">ML model: user embedding &#183; tweet embedding</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Social Proximity (w4 = 0.15)</div>
        <div class="d-flow-v">
          <div class="d-box indigo">Mutual follows &#8594; 1.0</div>
          <div class="d-box indigo">Frequent interactions &#8594; 0.8</div>
          <div class="d-box indigo">One-way follow &#8594; 0.4</div>
          <div class="d-box indigo">Friend of friend &#8594; 0.1</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Ranking Pipeline</div>
    <div class="d-flow">
      <div class="d-box gray">Candidate tweets (&#8776; 500)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Light ranker (score formula)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Top 50 &#8594; heavy ML ranker</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Top 20 returned to client</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-search-indexing",
		Title:       "Search &amp; Indexing Pipeline",
		Description: "Real-time inverted index for hashtags and keywords with Kafka-powered ingestion",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Indexing Pipeline (Write Path)</div>
    <div class="d-flow">
      <div class="d-box green">New Tweet</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Kafka (tweet-events)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Tokenizer / NLP Service</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Elasticsearch Cluster</div>
    </div>
    <div class="d-flow" style="margin-top:8px">
      <div class="d-box blue">Tokenizer extracts:</div>
      <div class="d-box gray">#hashtags</div>
      <div class="d-box gray">@mentions</div>
      <div class="d-box gray">keywords (stemmed)</div>
      <div class="d-box gray">URLs (expanded)</div>
      <div class="d-box gray">language</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Inverted Index Structure (Elasticsearch)</div>
    <div class="d-cols">
      <div class="d-col">
        <div class="d-entity">
          <div class="d-entity-header amber">Hashtag Index</div>
          <div class="d-entity-body">
            <div><span class="pk">PK</span> hashtag (lowercase)</div>
            <div>posting_list: [tweet_id, score, timestamp]</div>
            <div>doc_count (for trending)</div>
            <div><span class="idx idx-btree">IDX</span> timestamp (recent first)</div>
          </div>
        </div>
      </div>
      <div class="d-col">
        <div class="d-entity">
          <div class="d-entity-header blue">Full-Text Index</div>
          <div class="d-entity-body">
            <div><span class="pk">PK</span> token (stemmed word)</div>
            <div>posting_list: [tweet_id, tf-idf, position]</div>
            <div>Analyzers: english, cjk, emoji</div>
            <div><span class="idx idx-btree">IDX</span> relevance + recency</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Search Query Path</div>
    <div class="d-flow">
      <div class="d-box green">GET /v1/search?q=#golang</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Elasticsearch query (filtered + scored)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Hydrate tweet objects</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Return ranked results</div>
    </div>
    <div class="d-label">Latency: &#8776; 20-50ms. Index lag: &#8776; 2-5 seconds from tweet post to searchable.</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-scaling",
		Title:       "Scaling Strategy",
		Description: "Sharding tweets by snowflake ID, timeline cache partitioning, and Kafka topic design",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Tweet Store Sharding</div>
      <div class="d-flow-v">
        <div class="d-box green"><strong>Shard key:</strong> tweet_id (snowflake)</div>
        <div class="d-label">Snowflake embeds timestamp &#8594; range scans by time are shard-local</div>
        <div class="d-box blue">16 shards &#215; 3 replicas = 48 nodes</div>
        <div class="d-box blue">Each shard: &#8776; 30M tweets/day</div>
        <div class="d-box purple">User timeline query: scatter to all shards by user_id index</div>
        <div class="d-label">Mitigated by Redis timeline cache (no DB hit for home feed)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Social Graph Sharding</div>
      <div class="d-flow-v">
        <div class="d-box amber"><strong>Shard key:</strong> user_id</div>
        <div class="d-label">All followers/following for a user on same shard</div>
        <div class="d-box amber">8 shards &#8212; follows table is smaller than tweets</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Timeline Cache (Redis Cluster)</div>
      <div class="d-flow-v">
        <div class="d-box red"><strong>Shard key:</strong> hash(user_id) mod 64</div>
        <div class="d-box red">64 Redis masters &#215; 2 replicas = 192 nodes</div>
        <div class="d-box red">Each node: &#8776; 8M user timelines</div>
        <div class="d-label">Memory per node: 800 entries &#215; 16 bytes &#215; 8M = &#8776; 100GB</div>
        <div class="d-box purple">Eviction: LRU for users inactive &gt; 7 days</div>
        <div class="d-label">Cold user returns &#8594; rebuild from DB + fan-out backfill</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Kafka Partitioning</div>
      <div class="d-flow-v">
        <div class="d-box indigo"><strong>tweet-events:</strong> 128 partitions</div>
        <div class="d-label">Partition key: user_id &#8212; all tweets from same user in order</div>
        <div class="d-box indigo"><strong>fanout-tasks:</strong> 256 partitions</div>
        <div class="d-label">Partition key: target_user_id &#8212; ordered delivery per timeline</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-failure-modes",
		Title:       "Failure Modes &amp; Mitigations",
		Description: "Timeline inconsistency, fan-out lag, and celebrity tweet thundering herd scenarios",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">1. Timeline Inconsistency</div>
        <div class="d-flow-v">
          <div class="d-box red"><strong>Failure:</strong> Fan-out partially completes &#8212; some followers see tweet, others don&#8217;t</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber"><strong>Cause:</strong> Fan-out worker crashes mid-batch</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green"><strong>Mitigation:</strong> Kafka consumer offsets &#8212; restart from last committed offset. Idempotent ZADD (re-adding same tweet_id is a no-op).</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">2. Fan-out Lag</div>
        <div class="d-flow-v">
          <div class="d-box red"><strong>Failure:</strong> Viral tweet causes fan-out queue to back up &#8212; 10K+ follower user floods Kafka</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber"><strong>Symptom:</strong> Timeline delivery delay exceeds 30s SLA</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green"><strong>Mitigation:</strong> Priority queue &#8212; celebrity fan-out is already skipped. For high-follower non-celebrities (10K-100K), batch ZADD in groups of 1,000. Auto-scale Kafka consumers on lag metric.</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">3. Celebrity Tweet Thundering Herd</div>
        <div class="d-flow-v">
          <div class="d-box red"><strong>Failure:</strong> Celebrity tweets &#8594; millions of concurrent timeline reads all hit celebrity cache simultaneously</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber"><strong>Symptom:</strong> Redis celebrity cache hot key &#8594; single shard overloaded</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green"><strong>Mitigation:</strong> Replicate celebrity cache across multiple read replicas. Add jitter to client polling intervals. Use Redis Cluster with hash tags to spread celebrity keys across slots.</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">4. Redis Timeline Cache Failure</div>
        <div class="d-flow-v">
          <div class="d-box red"><strong>Failure:</strong> Redis node goes down &#8212; &#8776; 8M user timelines lost</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber"><strong>Symptom:</strong> Cache misses spike &#8594; DB scatter-gather queries surge</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green"><strong>Mitigation:</strong> Redis replica auto-promotes. Circuit breaker on DB: return stale/partial timeline. Background job rebuilds timelines from follows table + tweet store. Rate-limit rebuild to 10K users/sec.</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-trending-topics",
		Title:       "Trending Topics &#8212; Real-Time Detection",
		Description: "Sliding window counting with Count-Min Sketch for approximate top-K trending topics",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Ingestion Pipeline</div>
    <div class="d-flow">
      <div class="d-box green">Tweet Stream</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Hashtag Extractor</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Kafka (hashtag-counts topic)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Trending Service</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Count-Min Sketch (Approximate Counting)</div>
        <div class="d-flow-v">
          <div class="d-box purple"><strong>Structure:</strong> d=5 hash functions &#215; w=10,000 counters</div>
          <div class="d-box purple"><strong>Memory:</strong> 5 &#215; 10K &#215; 4 bytes = 200KB per window</div>
          <div class="d-box blue"><strong>Insert:</strong> hash(hashtag) &#8594; increment d=5 counters</div>
          <div class="d-box blue"><strong>Query:</strong> min of d=5 counters = approximate count</div>
          <div class="d-box amber"><strong>Error:</strong> Over-counts by &#8804; &#949;N with probability 1-&#948;</div>
          <div class="d-label">&#949; = e/w &#8776; 0.03%, &#948; = e^(-d) &#8776; 0.007. Sufficient for top-K ranking.</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Sliding Window (5-Minute Buckets)</div>
        <div class="d-flow-v">
          <div class="d-flow">
            <div class="d-box gray">t-4</div>
            <div class="d-box gray">t-3</div>
            <div class="d-box gray">t-2</div>
            <div class="d-box blue">t-1</div>
            <div class="d-box green">t (current)</div>
          </div>
          <div class="d-label">Window = sum of last 12 buckets (1 hour). Old buckets evicted.</div>
          <div class="d-box amber"><strong>Acceleration detection:</strong><br/>trend_score = count_last_1h / count_prev_1h</div>
          <div class="d-label">Score &gt; 3.0 = &#8220;trending&#8221;. Score &gt; 10.0 = &#8220;viral&#8221;.</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Top-K Extraction</div>
        <div class="d-flow-v">
          <div class="d-box green">Min-heap of size K=50</div>
          <div class="d-label">For each hashtag update: if count &gt; heap.min, replace and sift down</div>
          <div class="d-box green">Publish top 50 trending every 30 seconds</div>
          <div class="d-label">Regionalized: separate heaps per country/city</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})
}

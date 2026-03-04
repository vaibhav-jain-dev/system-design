package diagrams

func registerTwitterFeed(r *Registry) {
	r.Register(&Diagram{
		Slug:        "tf-requirements",
		Title:       "Scale Estimates & Requirements",
		Description: "Functional and non-functional requirements with scale estimates for Twitter feed",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Estimates</div>
      <div class="d-flow-v">
        <div class="d-box blue">500M DAU &#8226; 200M tweets/day</div>
        <div class="d-box blue">Avg user follows 200 accounts</div>
        <div class="d-box blue">Timeline reads: 500M &#215; 10 refreshes = 5B/day</div>
        <div class="d-box purple">Read QPS: 58K &#8226; Peak (3x) = 174K</div>
        <div class="d-box purple">Write QPS: 2.3K &#8226; Peak (5x) = 11.5K</div>
        <div class="d-box amber">Read:Write ratio = 25:1</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green">Post tweets (text, images, links)</div>
        <div class="d-box green">Home timeline: tweets from followed users</div>
        <div class="d-box green">Timeline ranked by relevance + recency</div>
        <div class="d-box blue">Follow/unfollow users</div>
        <div class="d-box blue">Search tweets by keyword</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">Timeline load: &lt; 200ms p99</div>
        <div class="d-box purple">Tweet delivery: &lt; 5s to followers</div>
        <div class="d-box amber">Availability: 99.99%</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-api-design",
		Title:       "API Endpoints",
		Description: "REST API design for tweet operations, timeline retrieval, and follow management",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Tweet Operations</div>
      <div class="d-flow-v">
        <div class="d-box green" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">POST /v1/tweets
{
  "text": "Hello world",
  "media_ids": ["m1", "m2"],
  "reply_to": null
}
&#8594; 201 {"tweet_id": "1234567890"}</div>
        <div class="d-box blue" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">DELETE /v1/tweets/{tweet_id}
&#8594; 204 No Content</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Timeline</div>
      <div class="d-flow-v">
        <div class="d-box amber" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">GET /v1/timeline/home
  ?cursor=eyJ0...&amp;limit=20
&#8594; 200 {
  "tweets": [...],
  "next_cursor": "eyJ0..."
}</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Social Graph</div>
      <div class="d-flow-v">
        <div class="d-box purple" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">POST /v1/users/{id}/follow
DELETE /v1/users/{id}/follow
GET /v1/users/{id}/followers
  ?cursor=...&amp;limit=50</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-fanout-comparison",
		Title:       "Fan-out Strategies: Write vs Read vs Hybrid",
		Description: "Comparison of three fan-out approaches for delivering tweets to follower timelines",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fan-out-on-Write (Push)</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center">User posts tweet</div>
        <div class="d-arrow-down">&#8595; async</div>
        <div class="d-box green" style="text-align:center">Write to all 200 follower timelines</div>
        <div class="d-label">&#10003; Fast reads: O(1) timeline fetch</div>
        <div class="d-label">&#10007; Celebrity: 10M followers = 10M writes</div>
        <div class="d-box blue">Best for: users with &lt; 10K followers</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Fan-out-on-Read (Pull)</div>
      <div class="d-flow-v">
        <div class="d-box amber" style="text-align:center">User opens timeline</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" style="text-align:center">Fetch tweets from all 200 followed users</div>
        <div class="d-label">&#10003; No write amplification</div>
        <div class="d-label">&#10007; Slow reads: 200 lookups + merge + sort</div>
        <div class="d-box blue">Best for: celebrity timelines</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">&#9733; Hybrid (Recommended)</div>
      <div class="d-flow-v">
        <div class="d-box indigo" style="text-align:center">If followers &lt; 10K &#8594; Push</div>
        <div class="d-box indigo" style="text-align:center">If followers &#8805; 10K &#8594; Pull at read</div>
        <div class="d-label">&#10003; Best of both worlds</div>
        <div class="d-label">&#10003; 99% users are push, 1% pull</div>
        <div class="d-box green">Twitter/X uses this exact approach</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end architecture from client through API gateway, tweet service, fan-out, and timeline service",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue" style="text-align:center"><strong>Mobile App</strong></div>
    <div class="d-box blue" style="text-align:center"><strong>Web Client</strong></div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" style="text-align:center"><strong>API Gateway + Load Balancer</strong><br>Rate limiting &#8226; Auth &#8226; Routing</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box green" style="text-align:center"><strong>Tweet Service</strong><br>Create/Delete tweets<br>Store in TweetDB</div>
        <div class="d-arrow-down">&#8595; event</div>
        <div class="d-box amber" style="text-align:center"><strong>Kafka</strong><br>tweet-created topic</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" style="text-align:center"><strong>Fan-out Service</strong><br>Push to follower timelines<br>(if followers &lt; 10K)</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box purple" style="text-align:center"><strong>Timeline Service</strong><br>Read pre-computed timeline<br>+ merge celebrity tweets</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" style="text-align:center"><strong>Redis Cluster</strong><br>Timeline cache<br>ZSET per user</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box indigo" style="text-align:center"><strong>Social Graph Service</strong><br>Follow/Unfollow<br>Follower lists</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray" style="text-align:center"><strong>Graph DB / MySQL</strong><br>Adjacency lists<br>Follower counts</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-tweet-flow",
		Title:       "Tweet Creation &#8594; Fan-out &#8594; Timeline",
		Description: "Step-by-step flow from tweet creation through fan-out workers to follower timeline delivery",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>1. User Posts Tweet</strong><br>POST /v1/tweets &#8594; Tweet Service</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" style="text-align:center"><strong>2. Persist to TweetDB (MySQL)</strong><br>tweets table: id, user_id, text, media, created_at<br>Return 201 to user immediately</div>
  <div class="d-arrow-down">&#8595; async event</div>
  <div class="d-box gray" style="text-align:center"><strong>3. Publish to Kafka</strong><br>Topic: tweet-created &#8226; Key: user_id &#8226; Value: tweet payload</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" style="text-align:center"><strong>4. Fan-out Worker Consumes</strong><br>Fetch follower list from Social Graph (cached)<br>User has 5,000 followers &#8594; 5,000 Redis writes</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red" style="text-align:center"><strong>5. Write to Redis Timelines</strong><br>For each follower: ZADD timeline:{follower_id} {timestamp} {tweet_id}<br>ZREMRANGEBYRANK to cap at 800 entries<br>Pipeline: 50 ZADD per batch &#8226; ~100ms for 5K followers</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" style="text-align:center"><strong>6. Follower Opens Timeline</strong><br>ZREVRANGE timeline:{user_id} 0 19 &#8594; 20 tweet IDs<br>Multi-GET tweet details &#8594; hydrate &#8594; return</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-celebrity-handling",
		Title:       "Celebrity Handling &#8212; Hybrid Fan-out Threshold",
		Description: "How the system handles celebrity accounts with millions of followers using pull-based merge",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Celebrity Posts (&#8805; 10K followers)</div>
        <div class="d-flow-v">
          <div class="d-box red" style="text-align:center"><strong>@celebrity</strong><br>50M followers</div>
          <div class="d-arrow-down">&#8595; tweet</div>
          <div class="d-box amber" style="text-align:center"><strong>NO fan-out!</strong><br>Tweet stored only in TweetDB<br>+ celebrity_tweets Redis list</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Timeline Read (Merge at Read)</div>
        <div class="d-flow-v">
          <div class="d-box blue" style="text-align:center"><strong>User opens timeline</strong></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-flow">
            <div class="d-branch">
              <div class="d-branch-arm">
                <div class="d-box green" style="text-align:center"><strong>Pre-computed</strong><br>ZREVRANGE timeline<br>(non-celebrity tweets)</div>
              </div>
              <div class="d-branch-arm">
                <div class="d-box purple" style="text-align:center"><strong>Pull celebrity</strong><br>User follows 5 celebrities<br>Fetch latest 20 from each</div>
              </div>
            </div>
          </div>
          <div class="d-arrow-down">&#8595; merge + sort</div>
          <div class="d-box green" style="text-align:center"><strong>Merged timeline</strong><br>Top 20 by score &#8226; &lt; 50ms total</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-label">Threshold: 10K followers &#8226; ~0.1% of users are celebrities but cause 99% of fan-out cost if pushed</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-timeline-cache",
		Title:       "Redis Timeline Cache &#8212; Sorted Set Structure",
		Description: "How each user's timeline is stored as a Redis sorted set with tweet IDs scored by timestamp",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Redis Sorted Set per User</div>
      <div class="d-flow-v">
        <div class="d-box red" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">KEY: timeline:{user_id}

ZADD timeline:u123 1704067200 "t9876"
ZADD timeline:u123 1704067190 "t9875"
ZADD timeline:u123 1704067180 "t9874"
...

ZREVRANGE timeline:u123 0 19
&#8594; ["t9876", "t9875", "t9874", ...]

ZCARD timeline:u123
&#8594; 800 (capped)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Memory Calculation</div>
      <div class="d-flow-v">
        <div class="d-box amber">Per entry: ~80 bytes (score + member)</div>
        <div class="d-box amber">Per user: 800 entries &#215; 80B = 64 KB</div>
        <div class="d-box amber">500M users &#215; 64KB = 32 TB total</div>
        <div class="d-box purple">Active users (50M): 3.2 TB</div>
        <div class="d-label">Only cache active users. Inactive: rebuild on demand from DB</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Eviction Strategy</div>
      <div class="d-flow-v">
        <div class="d-box blue">LRU eviction for inactive timelines</div>
        <div class="d-box blue">ZREMRANGEBYRANK cap at 800 per write</div>
        <div class="d-box green">TTL: 7 days for inactive users</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-data-model",
		Title:       "Data Model &#8212; Core Tables",
		Description: "Database schema for tweets, users, follows, and timeline tables",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">tweets (MySQL &#8594; sharded)</div>
      <div class="d-entity-body">
        <div class="pk">tweet_id BIGINT (Snowflake)</div>
        <div class="fk">user_id BIGINT</div>
        <div>text VARCHAR(280)</div>
        <div>media_urls JSON</div>
        <div>reply_to_id BIGINT NULL</div>
        <div>retweet_of_id BIGINT NULL</div>
        <div class="idx idx-btree">created_at TIMESTAMP</div>
        <div>like_count INT DEFAULT 0</div>
        <div>retweet_count INT DEFAULT 0</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header green">users</div>
      <div class="d-entity-body">
        <div class="pk">user_id BIGINT</div>
        <div class="idx idx-hash">username VARCHAR(50) UNIQUE</div>
        <div>display_name VARCHAR(100)</div>
        <div>follower_count INT</div>
        <div>following_count INT</div>
        <div>is_celebrity BOOLEAN</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header purple">follows (Social Graph)</div>
      <div class="d-entity-body">
        <div class="pk">follower_id BIGINT</div>
        <div class="pk">followee_id BIGINT</div>
        <div class="idx idx-btree">created_at TIMESTAMP</div>
        <div>INDEX (followee_id, follower_id)</div>
        <div>INDEX (follower_id, created_at)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header red">timeline_cache (Redis)</div>
      <div class="d-entity-body">
        <div class="pk">KEY: timeline:{user_id}</div>
        <div>TYPE: Sorted Set (ZSET)</div>
        <div>SCORE: tweet timestamp</div>
        <div>MEMBER: tweet_id</div>
        <div>MAX SIZE: 800 entries</div>
        <div>TTL: 7 days (inactive users)</div>
      </div>
    </div>
  </div>
</div>
<div class="d-er-lines">
  <div class="d-er-connector">
    <span class="d-er-from">users</span>
    <span class="d-er-type">1:N</span>
    <span class="d-er-to">tweets (user_id)</span>
  </div>
  <div class="d-er-connector">
    <span class="d-er-from">users</span>
    <span class="d-er-type">M:N</span>
    <span class="d-er-to">follows (follower &#8596; followee)</span>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-trending-detection",
		Title:       "Trending Topic Detection",
		Description: "Sliding window count pipeline using Count-Min Sketch and z-score spike detection",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>Tweet Stream</strong><br>200M tweets/day &#8594; extract hashtags + entities</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" style="text-align:center"><strong>Kafka (hashtag-counts topic)</strong><br>Partitioned by hashtag hash</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Count-Min Sketch (per 5-min window)</div>
        <div class="d-flow-v">
          <div class="d-box amber" style="text-align:center">4 hash functions &#215; 10K counters<br>Memory: 40 KB per window<br>Error rate: &lt; 0.1%</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Spike Detection</div>
        <div class="d-flow-v">
          <div class="d-box purple" style="text-align:center">z-score = (current &#8722; mean) / stddev<br>Trending if z-score &gt; 3.0<br>Compare to same hour last week</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" style="text-align:center"><strong>Trending Topics (Top 30)</strong><br>Ranked by velocity (rate of increase) not raw count<br>Cached in Redis &#8226; Refreshed every 5 min &#8226; Per-region variants</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-scaling",
		Title:       "Sharding Strategy",
		Description: "Sharding approach for tweets, timelines, and social graph across distributed storage",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Tweet Storage (MySQL Vitess)</div>
      <div class="d-flow-v">
        <div class="d-box blue" style="text-align:center"><strong>Shard by user_id</strong><br>All tweets from one user on same shard</div>
        <div class="d-box green" style="text-align:center">16 shards &#215; 3 replicas<br>= 48 MySQL instances</div>
        <div class="d-label">User profile reads: single shard lookup</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Social Graph (MySQL)</div>
      <div class="d-flow-v">
        <div class="d-box purple" style="text-align:center"><strong>Shard by follower_id</strong><br>"Who do I follow?" = single shard</div>
        <div class="d-box purple" style="text-align:center">Reverse index shard by followee_id<br>"Who follows me?" for fan-out</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Timeline Cache (Redis Cluster)</div>
      <div class="d-flow-v">
        <div class="d-box red" style="text-align:center"><strong>Shard by user_id hash</strong><br>Consistent hashing across slots</div>
        <div class="d-flow">
          <div class="d-group" style="flex:1">
            <div class="d-group-title">Active tier</div>
            <div class="d-box amber" style="text-align:center">50M users<br>3.2 TB<br>64 shards</div>
          </div>
          <div class="d-group" style="flex:1">
            <div class="d-group-title">Warm tier</div>
            <div class="d-box gray" style="text-align:center">Rebuild on<br>demand from<br>DB + fan-out</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-real-time-updates",
		Title:       "Real-Time Timeline Updates (SSE/WebSocket)",
		Description: "How new tweets are pushed to connected clients via server-sent events or WebSocket",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Push Flow</div>
        <div class="d-flow-v">
          <div class="d-box green" style="text-align:center"><strong>Fan-out Service</strong><br>Writes to Redis timeline</div>
          <div class="d-arrow-down">&#8595; publish</div>
          <div class="d-box red" style="text-align:center"><strong>Redis Pub/Sub</strong><br>Channel: user:{user_id}:timeline<br>Lightweight notification (tweet_id only)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo" style="text-align:center"><strong>WebSocket Gateway</strong><br>Maintains long-lived connections<br>~5M concurrent connections per node</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue" style="text-align:center"><strong>Client</strong><br>Receives notification &#8594; fetch tweet detail<br>Prepend to timeline UI</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Connection Management</div>
        <div class="d-flow-v">
          <div class="d-box amber">Online users: track in Redis SET</div>
          <div class="d-box amber">Only push to online users</div>
          <div class="d-box amber">Offline: pull on next app open</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Fallback: SSE for Web</div>
        <div class="d-flow-v">
          <div class="d-box purple" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">GET /v1/timeline/stream
Accept: text/event-stream

data: {"tweet_id": "t9876"}
data: {"tweet_id": "t9877"}</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tf-search-indexing",
		Title:       "Tweet Search with Elasticsearch",
		Description: "How tweets are indexed and searched using Elasticsearch with relevance scoring",
		ContentFile: "problems/twitter-feed",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-group">
          <div class="d-group-title">Indexing Pipeline</div>
          <div class="d-flow-v">
            <div class="d-box green" style="text-align:center"><strong>New Tweet</strong></div>
            <div class="d-arrow-down">&#8595;</div>
            <div class="d-box gray" style="text-align:center"><strong>Kafka</strong><br>tweet-created topic</div>
            <div class="d-arrow-down">&#8595;</div>
            <div class="d-box amber" style="text-align:center"><strong>Index Worker</strong><br>Tokenize, stem, extract entities<br>Batch: 1000 docs per bulk request</div>
            <div class="d-arrow-down">&#8595;</div>
            <div class="d-box indigo" style="text-align:center"><strong>Elasticsearch Cluster</strong><br>20 shards &#8226; 2 replicas<br>~50 TB for 30-day window</div>
          </div>
        </div>
      </div>
      <div class="d-branch-arm">
        <div class="d-group">
          <div class="d-group-title">Search Query Flow</div>
          <div class="d-flow-v">
            <div class="d-box blue" style="text-align:center"><strong>User searches "breaking news"</strong></div>
            <div class="d-arrow-down">&#8595;</div>
            <div class="d-box purple" style="text-align:center"><strong>Search Service</strong><br>Query Elasticsearch<br>BM25 + recency boost + engagement</div>
            <div class="d-arrow-down">&#8595;</div>
            <div class="d-box green" style="text-align:center"><strong>Results</strong><br>Ranked by relevance &#215; recency<br>Filter: safe search, blocked users<br>Paginated via search_after</div>
          </div>
        </div>
        <div class="d-group">
          <div class="d-group-title">Index Schema</div>
          <div class="d-flow-v">
            <div class="d-box gray" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">tweet_id: keyword
text: text (analyzed)
user_id: keyword
hashtags: keyword[]
created_at: date
engagement: float</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})
}

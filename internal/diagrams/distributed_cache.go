package diagrams

func registerDistributedCache(r *Registry) {
	r.Register(&Diagram{
		Slug:        "dc-requirements",
		Title:       "Requirements & Scale",
		Description: "Scale targets, NFRs, and scope for distributed cache design.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Targets</div>
      <div class="d-entity">
        <div class="d-entity-header blue">Back-of-Envelope</div>
        <div class="d-entity-body">
          <div class="d-row"><span>Total items</span><span><strong>1 billion</strong></span></div>
          <div class="d-row"><span>Total data size</span><span><strong>1 TB</strong></span></div>
          <div class="d-row"><span>Avg item size</span><span>~1 KB</span></div>
          <div class="d-row"><span>Peak read ops</span><span><strong>500K ops/sec</strong></span></div>
          <div class="d-row"><span>Peak write ops</span><span><strong>50K ops/sec</strong></span></div>
          <div class="d-row"><span>Read:Write ratio</span><span>10:1</span></div>
          <div class="d-row"><span>Target hit rate</span><span>&gt;90%</span></div>
          <div class="d-row"><span>Read latency</span><span><strong>&lt;1ms p99</strong></span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="GET key → value in &lt;1ms. SET key value TTL. DEL key. Basic key-value operations are the core contract.">GET / SET / DEL with optional TTL</div>
        <div class="d-box green" data-tip="Horizontally sharded across N nodes. Adding a node increases both capacity and throughput.">Horizontal sharding across nodes</div>
        <div class="d-box green" data-tip="Each primary node has at least one replica. Reads can be served from replica to reduce primary load.">Primary-replica replication</div>
        <div class="d-box blue" data-tip="When cache is full, old items are evicted to make room. Policy (LRU, LFU) determines which items are removed.">Automatic eviction (LRU default)</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title">Non-Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box indigo" data-tip="99.99% availability = &lt;53 minutes downtime/year. Achieved via multi-node replication and automatic failover.">Availability: 99.99%</div>
        <div class="d-box indigo" data-tip="P99 read &lt;1ms. Achieved by keeping all data in RAM — no disk seeks. Network is the bottleneck, not storage.">Latency: &lt;1ms reads (p99)</div>
        <div class="d-box indigo" data-tip="Add nodes to scale linearly. Consistent hashing minimizes key remapping on topology changes.">Scalability: linear with nodes</div>
        <div class="d-box indigo" data-tip="Hot keys that receive disproportionate traffic must not overload a single node.">Hot-key handling</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-consistent-hashing",
		Title:       "Consistent Hashing Ring",
		Description: "How consistent hashing distributes keys across nodes and minimizes remapping when nodes are added or removed.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Hash Ring Concept</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Hash space is 0 → 2^32. Visualized as a ring where 0 and 2^32 are the same point. Both keys and nodes are hashed onto this ring.">Hash space: 0 → 2<sup>32</sup> (ring)</div>
        <div class="d-box blue" data-tip="MD5 or SHA1 of node hostname/IP → position on ring. Each node occupies a fixed arc of the ring.">Nodes hashed to ring positions</div>
        <div class="d-box blue" data-tip="For key K: hash(K) → position on ring → walk clockwise → first node encountered owns this key.">Keys map to nearest clockwise node</div>
      </div>
      <div class="d-entity" style="margin-top:8px;">
        <div class="d-entity-header indigo">Ring Positions (example)</div>
        <div class="d-entity-body">
          <div class="d-row"><span>Node0</span><span>position 0</span></div>
          <div class="d-row"><span>Node1</span><span>position 90</span></div>
          <div class="d-row"><span>Node2</span><span>position 180</span></div>
          <div class="d-row"><span>Node3</span><span>position 270</span></div>
          <div class="d-row"><span>key "user:123"</span><span>→ pos 45 → Node1</span></div>
          <div class="d-row"><span>key "session:abc"</span><span>→ pos 200 → Node3</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Adding a Node</div>
      <div class="d-flow-v">
        <div class="d-box gray" data-tip="Before: 4 nodes, each owns ~25% of keyspace (positions 0-359 split evenly).">Before: Node0, Node1, Node2, Node3<br><small>Each owns ~25% of keys</small></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Node4 added at position 135. It only takes keys from Node2 (which previously owned positions 90-180). Keys at 90-134 still go to Node1. Keys at 135-179 now go to Node4.">Add Node4 at position 135</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Only keys between 90-134 (previously on Node2) need to move to Node4. ~12.5% of keys remapped.">Only ~12.5% of keys remapped<br><small>(1/N of keyspace affected)</small></div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title">vs. Modulo Hashing</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="With modulo hashing (key % N), adding a 5th node changes the formula from key%4 to key%5. Almost every key maps to a different node.">Modulo: hash(key) % N → node<br><small>Adding node N+1: ~80% of keys remap</small></div>
        <div class="d-box green" data-tip="Consistent hashing: adding 1 node out of N remaps only 1/N of keys. At 10 nodes, only 10% of keys move.">Consistent hashing: only 1/N keys remap<br><small>10 nodes → only 10% of keys move</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Virtual Nodes</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="4 physical nodes with 1 position each → uneven distribution if hashes cluster. Node0 might get 40% and Node3 might get 10%.">Problem: Uneven distribution with few nodes</div>
        <div class="d-box green" data-tip="Each physical node gets 150 virtual nodes (vnodes) spread around the ring. Node0 occupies positions 3, 47, 92, 141, ... (150 positions total).">Solution: 150 vnodes per physical node</div>
        <div class="d-box green" data-tip="With 150 vnodes × 4 nodes = 600 total ring points. Distribution approaches uniform — each node handles ~25% of keys.">4 nodes × 150 vnodes = 600 ring points</div>
        <div class="d-box gray" data-tip="When Node0 fails, its 150 vnodes are removed from ring. Keys redistribute to 150 different successor vnodes — evenly spreading the extra load across all remaining nodes.">Failure: 150 vnodes spread load across<br>all remaining nodes evenly</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-architecture",
		Title:       "Cluster Architecture",
		Description: "End-to-end read and write paths through the distributed cache cluster with client-side sharding.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Clients</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Application server. Uses client library (e.g. Jedis, ioredis) that embeds the consistent hash ring.">App Server A</div>
        <div class="d-box blue" data-tip="App Server B. Each app server has its own copy of the ring — no proxy needed for sharding.">App Server B</div>
        <div class="d-box gray" data-tip="Client library holds the hash ring topology. On topology change (node add/remove), topology is propagated via config service or ZooKeeper watch.">Client lib: embedded hash ring</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Cache Cluster</div>
      <div class="d-flow" style="gap:8px;">
        <div class="d-flow-v" style="flex:1;">
          <div class="d-group">
            <div class="d-group-title">Shard 0</div>
            <div class="d-box green" data-tip="Primary node for shard 0. Handles all writes. Serves reads (or replica can serve reads for lower latency).">Primary 0<br><small>keys 0–25%</small></div>
            <div class="d-arrow-down" style="font-size:10px;">async repl</div>
            <div class="d-box gray" data-tip="Replica for shard 0. Async replication from primary. Can serve reads. Promoted to primary on failure.">Replica 0</div>
          </div>
        </div>
        <div class="d-flow-v" style="flex:1;">
          <div class="d-group">
            <div class="d-group-title">Shard 1</div>
            <div class="d-box green" data-tip="Primary node for shard 1. Independent of Shard 0 — failure of Shard 1 does not affect Shard 0.">Primary 1<br><small>keys 25–50%</small></div>
            <div class="d-arrow-down" style="font-size:10px;">async repl</div>
            <div class="d-box gray" data-tip="Replica for shard 1.">Replica 1</div>
          </div>
        </div>
        <div class="d-flow-v" style="flex:1;">
          <div class="d-group">
            <div class="d-group-title">Shard 2</div>
            <div class="d-box green" data-tip="Primary node for shard 2.">Primary 2<br><small>keys 50–75%</small></div>
            <div class="d-arrow-down" style="font-size:10px;">async repl</div>
            <div class="d-box gray" data-tip="Replica for shard 2.">Replica 2</div>
          </div>
        </div>
        <div class="d-flow-v" style="flex:1;">
          <div class="d-group">
            <div class="d-group-title">Shard 3</div>
            <div class="d-box green" data-tip="Primary node for shard 3.">Primary 3<br><small>keys 75–100%</small></div>
            <div class="d-arrow-down" style="font-size:10px;">async repl</div>
            <div class="d-box gray" data-tip="Replica for shard 3.">Replica 3</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
<div class="d-cols" style="margin-top:8px;">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Read Path</div>
      <div class="d-flow" style="gap:4px;align-items:center;">
        <div class="d-box blue" style="flex:none;padding:6px 10px;">Client</div>
        <div class="d-arrow">→ hash(key)</div>
        <div class="d-box green" style="flex:none;padding:6px 10px;">Primary/Replica</div>
        <div class="d-arrow">→ &lt;1ms</div>
        <div class="d-box blue" style="flex:none;padding:6px 10px;">value</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Write Path</div>
      <div class="d-flow" style="gap:4px;align-items:center;">
        <div class="d-box blue" style="flex:none;padding:6px 10px;">Client</div>
        <div class="d-arrow">→ hash(key)</div>
        <div class="d-box green" style="flex:none;padding:6px 10px;">Primary</div>
        <div class="d-arrow">→ async</div>
        <div class="d-box gray" style="flex:none;padding:6px 10px;">Replica</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-data-partitioning",
		Title:       "Data Partitioning",
		Description: "How 1TB of data is distributed across physical nodes using virtual nodes on the consistent hash ring.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Physical Sharding</div>
      <div class="d-entity">
        <div class="d-entity-header blue">Capacity Planning</div>
        <div class="d-entity-body">
          <div class="d-row"><span>Total data</span><span>1 TB</span></div>
          <div class="d-row"><span>Physical nodes</span><span>8</span></div>
          <div class="d-row"><span>Data per node</span><span><strong>128 GB</strong></span></div>
          <div class="d-row"><span>Instance type</span><span>r6g.4xlarge (128GB RAM)</span></div>
          <div class="d-row"><span>Memory overhead</span><span>~20% (keys, pointers)</span></div>
          <div class="d-row"><span>Usable per node</span><span>~102 GB</span></div>
          <div class="d-row"><span>Eviction starts at</span><span>80% full = 82 GB</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Virtual Node Configuration</div>
      <div class="d-entity">
        <div class="d-entity-header indigo">Hash Ring Topology</div>
        <div class="d-entity-body">
          <div class="d-row"><span>Physical nodes</span><span>8</span></div>
          <div class="d-row"><span>Vnodes per physical</span><span>150</span></div>
          <div class="d-row"><span>Total ring points</span><span>1,200</span></div>
          <div class="d-row"><span>Keys per vnode (avg)</span><span>1,000,000 ÷ 1200 ≈ 833K</span></div>
          <div class="d-row"><span>Keys per physical node</span><span>~125M</span></div>
          <div class="d-row"><span>Distribution variance</span><span>&lt;5% with 150 vnodes</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Rebalancing Cost</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Adding a 9th node: 1/9 of keys move = ~111M keys. At 1KB each = ~111 GB moved. Done gradually in background over hours.">Add node: 1/N keys move<br><small>9th node: ~111 GB transferred</small></div>
        <div class="d-box amber" data-tip="Node failure: 150 vnodes spread across all 7 remaining nodes. Each remaining node takes ~21 extra vnodes = 1/7 of failed node's load.">Node failure: load spread across N-1 nodes<br><small>Each remaining node: +14% load</small></div>
        <div class="d-box gray" data-tip="Rebalancing is done in background with throttling. Cache still serves traffic during rebalance — just with higher miss rate for moved keys.">Rebalancing does not block reads/writes</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-eviction-policies",
		Title:       "Eviction Policies Comparison",
		Description: "LRU, LFU, FIFO, and Random eviction — trade-offs and when to use each.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">LRU — Least Recently Used</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="HashMap for O(1) key lookup. Doubly linked list to track recency order — head = most recent, tail = eviction candidate.">Data structure: HashMap + doubly linked list</div>
        <div class="d-box green" data-tip="GET: move node to head of list. O(1) pointer update. SET new: insert at head, evict tail if full. O(1).">Time: O(1) get and put</div>
        <div class="d-box green" data-tip="Works well for web workloads where recently accessed data is likely to be accessed again (temporal locality).">Best for: temporal locality (web sessions, hot items)</div>
        <div class="d-box gray" data-tip="LRU suffers on large sequential scans — each key is accessed once and then evicted, polluting the cache with single-use items.">Weakness: sequential scan pollution</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">LFU — Least Frequently Used</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Frequency counter per key. Min-heap or frequency doubly-linked list tracks which key has the lowest access count.">Data structure: frequency counter + min-heap</div>
        <div class="d-box blue" data-tip="GET: increment counter. O(log N) for heap reorder. SET: insert with freq=1, evict min-freq key. O(log N).">Time: O(log N) get and put</div>
        <div class="d-box blue" data-tip="Good for CDN-style workloads where popular items are consistently popular over time and should never be evicted.">Best for: stable popularity (CDN, static assets)</div>
        <div class="d-box amber" data-tip="New keys start at frequency 1 — they're immediately eviction candidates even if they'll be popular. Requires aging to handle this.">Weakness: new items unfairly evicted (cold start)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">FIFO &amp; Random</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="FIFO evicts items in insertion order. Ignores access patterns entirely. Evicts a frequently-used old item in favor of a never-used new one.">FIFO: evicts oldest inserted item<br><small>Ignores access patterns — poor hit rate</small></div>
        <div class="d-box red" data-tip="Random eviction is simple and cache-line friendly but has poor hit rate — evicts hot items as often as cold ones.">Random: evicts random item<br><small>Simple but 15–30% worse hit rate than LRU</small></div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title" style="color:#059669;">Recommendation</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="LRU is the default for most web caches. Redis uses a sampled LRU (checks 5 random keys, evicts the least recently used among them) for memory efficiency.">Web cache / API cache → <strong>LRU</strong></div>
        <div class="d-box green" data-tip="CDN edge caches where popular videos/images stay popular for weeks. LFU keeps them resident.">CDN / static content → <strong>LFU</strong></div>
        <div class="d-box gray" data-tip="Redis implements allkeys-lru (evict from all keys) and volatile-lru (evict only keys with TTL). Choose based on whether you set TTLs.">Redis: allkeys-lru or volatile-lru</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-replication",
		Title:       "Replication Strategy",
		Description: "Primary-replica async replication, Sentinel-based failover, and replication lag handling.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Normal Operation</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Client writes go to primary. Primary acknowledges write immediately after writing to its own memory.">Client → Primary (sync write ACK)</div>
        <div class="d-arrow-down">&#8595; async replication</div>
        <div class="d-flow" style="gap:8px;">
          <div class="d-box gray" style="flex:1;" data-tip="Replica 1 receives replication stream from primary. Lag is typically &lt;100ms under normal load.">Replica 1<br><small>lag &lt;100ms</small></div>
          <div class="d-box gray" style="flex:1;" data-tip="Replica 2 for additional redundancy. Can serve reads to reduce primary load by 30-50%.">Replica 2<br><small>lag &lt;100ms</small></div>
        </div>
        <div class="d-box blue" data-tip="Reads can be served by any replica. Client reads stale data only during the replication lag window (&lt;100ms). Acceptable for most use cases.">Read from Replica: eventual consistency<br><small>Stale window: &lt;100ms</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Failover (Sentinel)</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Primary node goes down (crash, network partition, OOM). Sentinel processes detect via heartbeat failure.">Primary fails</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Sentinel requires majority vote (quorum) to declare primary down. With 3 Sentinels, 2 must agree. Prevents split-brain during network partition.">Sentinel quorum detects failure<br><small>Quorum: 2 of 3 Sentinels agree</small></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Sentinel elects the most up-to-date replica as new primary. Selection criteria: replication offset (least lag), then priority config.">Elect most up-to-date replica as primary<br><small>Selection: highest replication offset</small></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Sentinel notifies all clients via pub/sub and reconfigures remaining replicas to replicate from new primary.">Clients notified, replicas reconfigured</div>
        <div class="d-box gray" data-tip="Total failover time: 10s (heartbeat miss) + 5s (quorum vote) + 5s (promotion) ≈ 20-30s. During this time: reads from replicas work, writes fail.">Failover time: ~20–30 seconds total<br><small>Writes blocked during promotion</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Replication Lag</div>
      <div class="d-flow-v">
        <div class="d-box indigo" data-tip="Redis replication uses a replication backlog (circular buffer, default 1MB). If replica falls behind by more than the backlog size, full resync required.">Replication: backlog buffer (1MB default)</div>
        <div class="d-entity" style="margin-top:8px;">
          <div class="d-entity-header blue">Lag by Scenario</div>
          <div class="d-entity-body">
            <div class="d-row"><span>Normal write load</span><span>&lt;100ms</span></div>
            <div class="d-row"><span>Write spike (10× normal)</span><span>100ms–1s</span></div>
            <div class="d-row"><span>Network congestion</span><span>1–5s</span></div>
            <div class="d-row"><span>Partial disconnect</span><span>partial resync</span></div>
            <div class="d-row"><span>Full disconnect &gt;backlog</span><span>full resync</span></div>
          </div>
        </div>
        <div class="d-box amber" data-tip="Monitor replication lag with INFO replication → master_repl_offset vs slave_repl_offset. Alert if gap exceeds 1M bytes or 1 second.">Monitor: master_repl_offset − slave_repl_offset<br><small>Alert threshold: &gt;1M bytes or &gt;1s</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-hot-key-problem",
		Title:       "Hot Key Problem",
		Description: "Why hot keys overload single nodes and three strategies to distribute the load.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title" style="color:#dc2626;">The Problem</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="A viral video, trending topic, or popular product can concentrate 90% of requests on a single cache key. That key lives on exactly one node.">Key &quot;trending:viral-video&quot; → 90% of reads</div>
        <div class="d-arrow-down">&#8595; all map to</div>
        <div class="d-box red" data-tip="Consistent hashing maps this key to Node2. Node2 now handles 90% of total cache traffic. CPU maxed, latency spikes, other keys on Node2 get starved.">Node2: overloaded (CPU 100%, latency 50ms)</div>
        <div class="d-box red" data-tip="Node0, Node1, Node3 are idle while Node2 is overwhelmed. The cluster is imbalanced — scaling out doesn't help because the key can only live on one node.">Node0, Node1, Node3: &lt;5% utilization</div>
        <div class="d-box amber" data-tip="This is a fundamental limitation of consistent hashing: one key → one node. The only fix is to spread the key's load, not the key itself.">Root cause: one key = one node in CH</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title" style="color:#059669;">Solution 1: Local L1 Cache</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Each application server has a local in-process LRU cache (e.g. Caffeine in Java). Holds the 1000 most popular keys in heap memory.">Each app server: LRU cache (1000 items)</div>
        <div class="d-box green" data-tip="Hot key is served from L1 without any network call. 0ms latency, 0 requests to cache cluster.">L1 hit: 0ms, 0 cache cluster requests</div>
        <div class="d-box green" data-tip="L1 miss → check distributed cache (L2). L2 miss → check database. Standard read-through hierarchy.">L1 miss → L2 (distributed cache) → DB</div>
        <div class="d-box amber" data-tip="L1 cache is eventually consistent. On write, distributed cache is updated immediately. L1 becomes stale until TTL expires (typically 1s).">Con: L1 stale for up to TTL (1s)<br><small>Accept for reads, not for financial data</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title" style="color:#059669;">Solutions 2 &amp; 3</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Split key into N shards: trending:viral-video:0, trending:viral-video:1, ...:9. On read, pick random shard. Writes update all 10 shards.">Key splitting: &quot;viral-video:{0-9}&quot;<br><small>Spread 1 key across 10 nodes</small></div>
        <div class="d-box blue" data-tip="10 shards × 1 node each = 10 nodes share the load. Each node handles 1/10 of hot key traffic.">10× throughput for hot key</div>
        <div class="d-box blue" data-tip="Write cost: must update all 10 shards on mutation. Acceptable for read-heavy hot keys (reads >> writes).">Write cost: 10× writes per update</div>
        <div class="d-box purple" style="margin-top:8px;" data-tip="Add read replicas specifically for hot key nodes. Detected by monitoring: if a key gets &gt;X req/sec, auto-promote its node's replica to serve reads too.">Read replicas for hot nodes<br><small>Automated by monitoring: &gt;50K req/sec threshold</small></div>
        <div class="d-box gray" data-tip="Detect hot keys before they cause problems. Redis can report top keys by access count. Set threshold to auto-alert or auto-shard.">Detection: Redis OBJECT FREQ command<br><small>Alert at &gt;10% of total traffic on 1 key</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-consistency-model",
		Title:       "Consistency Models",
		Description: "Strong vs eventual consistency and read-your-own-writes strategies for distributed caches.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">The Replication Lag Window</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Client writes X=100 to primary at T=0. Primary acknowledges immediately.">T=0: Write X=100 to Primary (ACK)</div>
        <div class="d-box amber" data-tip="Replication lag: replica has not yet received the write. Replica still shows X=50 (old value).">T=10ms: Replica still has X=50</div>
        <div class="d-box red" data-tip="Another client reads from replica at T=10ms. Gets X=50 — stale data. This is the eventual consistency window.">T=10ms: Read from Replica → X=50 (stale!)</div>
        <div class="d-box green" data-tip="Replication arrives at T=100ms. Now replica has X=100. Any read after this point is consistent.">T=100ms: Replica updated → X=100</div>
        <div class="d-box gray" data-tip="The inconsistency window is 0–100ms in this example. This is acceptable for most read-heavy workloads (product listings, news feeds).">Inconsistency window: ~0–100ms</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Strong Consistency</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Always route reads to primary. Primary always has the latest data — no replication lag for reads.">Always read from primary</div>
        <div class="d-box green" data-tip="No stale reads possible. Every read reflects the latest write. Required for financial data, inventory counts, auth tokens.">Zero stale reads — latest data guaranteed</div>
        <div class="d-box red" data-tip="Primary handles all reads AND writes. Can't use replicas for read scaling. Primary becomes bottleneck at high read load.">Con: primary is the bottleneck<br><small>Replicas unused for reads</small></div>
        <div class="d-box red" data-tip="Primary is further away geographically → higher latency for users in other regions.">Con: higher latency (1 node path)</div>
        <div class="d-box gray" data-tip="Use strong consistency only when stale reads cause real damage: user's own profile, payment status, inventory, auth sessions.">Use for: financial, auth, inventory</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Eventual Consistency</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Reads go to nearest replica. Reads are fast (nearest node) and primary is freed from read load.">Read from nearest replica</div>
        <div class="d-box blue" data-tip="Primary handles only writes. Replicas handle all reads. At 10:1 read:write ratio, primary CPU freed by ~90%.">Lower latency, primary freed for writes</div>
        <div class="d-box amber" data-tip="Accept that reads may be up to 100ms stale. Fine for: product prices, trending feeds, non-critical counters.">Accept: up to 100ms stale reads</div>
        <div class="d-box gray" data-tip="Use eventual consistency for high-read-volume content that doesn't require precise freshness: social feeds, product catalogs, recommendation results.">Use for: feeds, catalogs, recommendations</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title">Read-Your-Own-Writes</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="After a user writes, route their subsequent reads to the same primary for a short window (e.g. 1 second). Other users can still read from replicas.">Sticky reads: user's writes → primary for 1s</div>
        <div class="d-box purple" data-tip="Implemented via session cookie or user-ID-based routing header. After 1s, replication has caught up and reads can go to replicas again.">Implementation: session-based routing</div>
        <div class="d-box gray" data-tip="Best of both worlds: users see their own writes immediately, while others see eventual consistency. Low overhead since only writers get sticky routing.">Best of both: writers see fresh, others eventual</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-failure-handling",
		Title:       "Failure Handling",
		Description: "Node failure, network partition, and cache stampede scenarios with mitigations.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Node Failure</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Primary node crashes. Health check detects failure within 5–10s.">Primary node crashes</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Consistent hashing remaps the failed node's virtual nodes to the next successor on the ring. Keys temporarily served by neighbor nodes (cache miss rate spikes).">CH ring: failed node's keys → neighbors</div>
        <div class="d-box amber" data-tip="Cache miss rate spikes for ~30s during failover — some requests fall through to database. Pre-warm replica to minimize this window.">Miss rate spikes 20–40% during failover</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Sentinel promotes replica to primary within 20–30s. Clients reconnect. Normal operation resumes.">Replica promoted to primary (~20s)</div>
        <div class="d-box green" data-tip="A new replica is provisioned for the new primary to restore full redundancy.">New replica provisioned for new primary</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Network Partition</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Network split: some Sentinels can see Primary, others cannot. If majority cannot reach Primary, they promote a replica — but Primary is still running and accepting writes. Split-brain!">Network splits — Sentinel quorum split</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Redis uses min-replicas-to-write config: if Primary cannot reach at least 1 replica, it stops accepting writes. This prevents split-brain data divergence.">min-replicas-to-write = 1<br><small>Primary stops writes if isolated</small></div>
        <div class="d-box amber" data-tip="Quorum requirement: majority of Sentinels must agree before promoting a replica. With 3 Sentinels, 2 must agree.">Sentinel quorum: majority must agree<br><small>3 Sentinels → need 2 to agree</small></div>
        <div class="d-box green" data-tip="These two mechanisms together prevent split-brain: at most one Primary accepts writes at any time.">Result: at most one writable primary</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Cache Stampede</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Popular cache key expires (TTL hit). 10,000 concurrent requests all get a cache miss simultaneously. All 10,000 rush to the database to reload the key.">Popular key expires → thundering herd</div>
        <div class="d-box red" data-tip="10,000 simultaneous DB queries for the same key. DB connection pool exhausted, query time spikes, potentially taking down the database.">10,000 DB queries for same key → DB overload</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="On cache miss, acquire distributed lock (Redis SETNX) before fetching from DB. Only the lock holder fetches from DB and populates cache. Others wait.">Mutex lock on cache miss:<br><code>SETNX lock:key 1 EX 5</code></div>
        <div class="d-box green" data-tip="Requests that don't get the lock wait 10ms then retry. By then the lock holder has populated the cache. They get a cache hit.">Waiters retry → cache hit (10ms wait)</div>
        <div class="d-box blue" data-tip="Probabilistic early expiration: before key expires, with probability 1/N refresh it. Keys are refreshed before they expire — no stampede.">Alt: probabilistic early refresh<br><small>Refresh before expiry with probability 1/N</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-write-policies",
		Title:       "Write-Through vs Write-Behind",
		Description: "Three write strategies — write-through, write-behind, and write-around — with latency and consistency trade-offs.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Write-Through</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Every write updates cache AND database synchronously before acknowledging to client. Both stores are always in sync.">Write: App → Cache + DB (sync)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue" data-tip="Client waits for both cache write and DB write to complete before receiving ACK. Total write latency = cache write + DB write.">ACK after both complete<br><small>Latency: cache (0.5ms) + DB (2ms) = 2.5ms</small></div>
        <div class="d-box green" data-tip="Cache and DB are always consistent. A crash immediately after write does not lose data — DB is already updated.">Pro: always consistent, no data loss</div>
        <div class="d-box red" data-tip="Every write goes to DB synchronously. Write latency is dominated by DB latency (2–5ms vs 0.5ms for cache-only). 4–10× slower writes.">Con: write latency dominated by DB<br><small>+2ms per write vs cache-only</small></div>
        <div class="d-box gray" data-tip="Use for financial data, inventory, user profiles — anywhere where cache-DB inconsistency causes real problems.">Use for: financial, inventory, auth</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Write-Behind (Write-Back)</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Write goes to cache immediately. DB is updated asynchronously in the background. Client gets ACK as soon as cache is written.">Write: App → Cache (ACK) → DB async</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Client gets ACK in ~0.5ms (just cache write). DB is updated asynchronously within 50–500ms.">ACK after cache write only<br><small>Latency: ~0.5ms (cache only)</small></div>
        <div class="d-box green" data-tip="Writes are 5× faster than write-through. Good for high-write workloads like analytics event counters, view counts, like counts.">Pro: 5× faster writes than write-through</div>
        <div class="d-box red" data-tip="If cache node crashes before async write reaches DB, the data is permanently lost. Cache data and DB data can diverge for 50–500ms window.">Con: data loss on cache crash<br><small>Loss window: up to 500ms of writes</small></div>
        <div class="d-box gray" data-tip="Use for analytics counters, social metrics (likes, views), recommendation signals — where approximate counts are acceptable and durability is not critical.">Use for: analytics, social metrics, counters</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Write-Around</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="Writes bypass cache entirely — go directly to DB. Cache is only populated on read (read-through). Writes don't pollute cache with infrequently-read data.">Write: App → DB only (bypass cache)</div>
        <div class="d-box purple" data-tip="First read: cache miss → fetch from DB → populate cache. Subsequent reads: cache hit.">Read: cache miss → DB → populate cache</div>
        <div class="d-box green" data-tip="Write-heavy data that is rarely read doesn't pollute the cache. Only cache keys that are actually read get stored.">Pro: cache only holds frequently-read data</div>
        <div class="d-box amber" data-tip="First read after write always goes to DB (cold cache). High read latency for recently written data.">Con: first read after write = cache miss<br><small>Cold read goes to DB every time</small></div>
        <div class="d-box gray" data-tip="Best for bulk imports, batch writes, or write-heavy pipelines where data is rarely read back. ETL pipelines, log ingestion.">Use for: bulk imports, ETL, infrequent reads</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title" style="color:#059669;">Decision Guide</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Financial transactions, user auth data, inventory counts require write-through for consistency guarantees.">Need consistency → Write-Through</div>
        <div class="d-box green" data-tip="Event counters, view counts, social likes — high write rate, some loss acceptable, latency matters → write-behind.">High write rate, loss OK → Write-Behind</div>
        <div class="d-box green" data-tip="Write-once, infrequently-read data (audit logs, archive data) → write-around.">Write-once, rare reads → Write-Around</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-memory-management",
		Title:       "Memory & Cost Analysis",
		Description: "Memory overhead calculation, eviction thresholds, and cost breakdown for a 1TB distributed cache.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Memory Overhead Calculation</div>
      <div class="d-entity">
        <div class="d-entity-header blue">Redis Key Overhead</div>
        <div class="d-entity-body">
          <div class="d-row"><span>Per-key overhead</span><span>~72 bytes</span></div>
          <div class="d-row"><span>Dict entry</span><span>32 bytes</span></div>
          <div class="d-row"><span>RedisObject</span><span>16 bytes</span></div>
          <div class="d-row"><span>LRU pointer</span><span>8 bytes</span></div>
          <div class="d-row"><span>Expiry pointer</span><span>8 bytes</span></div>
          <div class="d-row"><span>Key string</span><span>8 bytes (SDS)</span></div>
        </div>
      </div>
      <div class="d-box amber" style="margin-top:8px;" data-tip="1B keys × 72 bytes = 72 GB of overhead just for key metadata. Plus 1 TB for values. Total RAM needed: 1.072 TB minimum before eviction headroom.">1B keys × 72 bytes = <strong>72 GB key overhead</strong><br><small>On top of 1TB value data</small></div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Node Sizing</div>
      <div class="d-entity">
        <div class="d-entity-header green">r6g.4xlarge (AWS)</div>
        <div class="d-entity-body">
          <div class="d-row"><span>RAM</span><span>128 GB</span></div>
          <div class="d-row"><span>vCPUs</span><span>16</span></div>
          <div class="d-row"><span>Network</span><span>12.5 Gbps</span></div>
          <div class="d-row"><span>Cost</span><span>$1.20/hr (~$876/mo)</span></div>
          <div class="d-row"><span>Usable RAM (80%)</span><span>102 GB</span></div>
          <div class="d-row"><span>Nodes for 1 TB</span><span>10 (with overhead)</span></div>
          <div class="d-row"><span>Total cluster cost</span><span>~$8,760/mo (no replicas)</span></div>
          <div class="d-row"><span>With replicas (2×)</span><span>~$17,520/mo</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Eviction Configuration</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Redis starts evicting when used memory hits maxmemory. Set to 80% of available RAM to leave headroom for memory fragmentation.">maxmemory = 80% RAM per node<br><small>r6g.4xlarge: 102 GB limit</small></div>
        <div class="d-box blue" data-tip="allkeys-lru: evict any key using LRU approximation. Best for general caches where all keys should be evictable.">maxmemory-policy = allkeys-lru</div>
        <div class="d-box blue" data-tip="Redis LRU is sampled — it picks 5 random keys and evicts the least recently used among them. Not exact LRU but close enough. Configurable with maxmemory-samples.">LRU is sampled: maxmemory-samples = 5<br><small>Increase to 10 for better accuracy</small></div>
        <div class="d-box gray" data-tip="Monitor used_memory_rss / used_memory. If ratio > 1.5, memory fragmentation is high. Consider MEMORY PURGE or restart.">Fragmentation ratio: rss/used &lt; 1.5 target</div>
        <div class="d-box gray" data-tip="Alert when used memory hits 75% of maxmemory. At 80% Redis starts evicting — you want warning before eviction starts.">Alert at 75% of maxmemory</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "dc-monitoring",
		Title:       "Monitoring Dashboard",
		Description: "Key metrics, alert thresholds, and dashboard panels for distributed cache operations.",
		ContentFile: "problems/distributed-cache",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Cache Performance Metrics</div>
      <div class="d-entity">
        <div class="d-entity-header blue">Key Metrics</div>
        <div class="d-entity-body">
          <div class="d-row"><span>Hit rate</span><span style="color:#059669;"><strong>&gt;90% target</strong></span></div>
          <div class="d-row"><span>Miss rate</span><span>&lt;10% target</span></div>
          <div class="d-row"><span>Eviction rate</span><span>keys evicted/sec</span></div>
          <div class="d-row"><span>Commands/sec</span><span>ops/sec per node</span></div>
          <div class="d-row"><span>Get latency p99</span><span>&lt;1ms target</span></div>
          <div class="d-row"><span>Set latency p99</span><span>&lt;2ms target</span></div>
          <div class="d-row"><span>Connected clients</span><span>connection count</span></div>
          <div class="d-row"><span>Blocked clients</span><span>should be 0</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Alert Conditions</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Hit rate below 80% means most requests are falling through to the database. Cache is not effective — investigate key eviction rate and TTL configuration.">hit_rate &lt; 80% → PAGE<br><small>Cache not effective, DB overloaded</small></div>
        <div class="d-box red" data-tip="Memory above 90% is dangerous — eviction rate will spike and performance will degrade. Scale out immediately.">memory_used &gt; 90% maxmemory → PAGE<br><small>Scale out: add nodes immediately</small></div>
        <div class="d-box red" data-tip="Replication lag above 5 seconds means replica data is dangerously stale. Investigate network or write throughput.">replication_lag &gt; 5s → PAGE<br><small>Risk: failover to very stale replica</small></div>
        <div class="d-box amber" data-tip="p99 latency above 5ms indicates the cache node is under pressure. Check CPU, network, and memory.">get_latency p99 &gt; 5ms → WARN<br><small>Investigate: CPU, memory, network</small></div>
        <div class="d-box amber" data-tip="Eviction rate above 1000/sec means cache is full and working items are being evicted. Increase memory or add nodes.">eviction_rate &gt; 1000/sec → WARN<br><small>Cache too small, add capacity</small></div>
        <div class="d-box blue" data-tip="Single key receiving high fraction of traffic. Hot key problem — apply sharding or local cache.">single_key_hit_pct &gt; 10% → WARN<br><small>Hot key detected, apply sharding</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Dashboard Panels</div>
      <div class="d-flow-v">
        <div class="d-box indigo" data-tip="Time-series: hit rate and miss rate over time. Dips in hit rate indicate eviction spikes or cold starts.">Hit/Miss rate over time (time series)</div>
        <div class="d-box indigo" data-tip="Per-node memory usage. Should be even across nodes — imbalance indicates hot keys or skewed hash distribution.">Memory used per node (bar chart)</div>
        <div class="d-box indigo" data-tip="Heatmap of p50, p95, p99 latency per endpoint (GET, SET, DEL). Shows latency spikes.">Latency percentiles heatmap</div>
        <div class="d-box indigo" data-tip="Commands/sec per node. Should be relatively even. Large imbalance = hot key or hash ring problem.">Ops/sec per node (balanced?)</div>
        <div class="d-box indigo" data-tip="Replication offset difference between primary and replica over time. Should be near zero.">Replication lag per shard</div>
        <div class="d-box gray" data-tip="Top 10 keys by access frequency. Identify hot keys before they cause problems.">Top 10 hot keys (OBJECT FREQ)</div>
        <div class="d-box gray" data-tip="Number of evictions per second. Near-zero under normal operation. Spikes indicate cache is undersized.">Eviction rate (should be near 0)</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

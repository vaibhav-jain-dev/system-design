package diagrams

func registerAlgorithms(r *Registry) {
	// ---------------------------------------------------------------
	// Base62 Encoding (2 diagrams)
	// ---------------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "algo-base62-process",
		Title:       "Base62 Encoding Process",
		Description: "Step-by-step divide-by-62 conversion of an integer to a Base62 string, showing each remainder mapped to a character.",
		ContentFile: "algorithms/base62-encoding",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v" style="align-items: center;">
  <div class="d-box blue" data-tip="Any positive integer works. 62^7 supports ~3.5 trillion unique short codes."><span class="d-step">1</span> Input: 123456789</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-group">
    <div class="d-group-title">Repeated divide-by-62 <span class="d-metric latency">O(log₆₂ N)</span></div>
    <div class="d-flow-v" style="gap: 0.25rem;">
      <div class="d-flow">
        <div class="d-box gray"><span class="d-step">2</span> 123456789 % 62 = 17</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber" data-tip="Charset: 0-9 (10) + a-z (26) + A-Z (26) = 62 chars. URL-safe, no special characters.">r</div>
      </div>
      <div class="d-flow">
        <div class="d-box gray"><span class="d-step">3</span> 1991239 % 62 = 25</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber">z</div>
      </div>
      <div class="d-flow">
        <div class="d-box gray"><span class="d-step">4</span> 32116 % 62 = 4</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber">e</div>
      </div>
      <div class="d-flow">
        <div class="d-box gray"><span class="d-step">5</span> 517 % 62 = 21</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber">v</div>
      </div>
      <div class="d-flow">
        <div class="d-box gray"><span class="d-step">6</span> 8 % 62 = 8</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber">i</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">↓</div>
  <div class="d-label"><span class="d-step">7</span> Read remainders bottom-to-top</div>
  <div class="d-box green" data-tip="7 chars max = 62^7 = 3.5 trillion unique codes. Enough for 100K URLs/day for 95,000 years.">Result: "ivezr" <span class="d-metric size">5 chars</span></div>
  <div class="d-label" style="margin-top: 0.5rem;">Decode: 8*62^4 + 21*62^3 + 4*62^2 + 25*62 + 17 = 123456789</div>
  <div class="d-caption">Base62 is bijective: every integer maps to exactly one string and vice versa. Encoding and decoding are both O(log₆₂ N) — typically 5-7 iterations for realistic IDs.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-base62-url-shortener",
		Title:       "Base62 in URL Shortener Architecture",
		Description: "Write and read paths of a URL shortener showing how Base62 encoding integrates with the ID generator, database, and cache layers.",
		ContentFile: "algorithms/base62-encoding",
		Type:        TypeHTML,
		HTML: `<div class="d-group" style="margin-bottom: 1.5rem;">
  <div class="d-group-title">Write Path <span class="d-metric latency">~15ms p99</span></div>
  <div class="d-flow-v" style="align-items: center;">
    <div class="d-box blue"><span class="d-step">1</span> Client: POST /shorten</div>
    <div class="d-arrow-down">↓</div>
    <div class="d-flow">
      <div class="d-box purple" data-tip="Stateless — any instance can handle the request. Horizontally scalable behind ALB."><span class="d-step">2</span> API Server</div>
      <div class="d-arrow">→</div>
      <div class="d-box green" data-tip="Counter-based: sequential, simple, single point of failure. KGS: pre-generated key ranges, no coordination needed at write time."><span class="d-step">3</span> ID Generator (counter/KGS)</div>
      <div class="d-arrow">→</div>
      <div class="d-box amber" data-tip="O(log₆₂ N) conversion. 7-char codes support 3.5 trillion URLs. No collision possible — bijective mapping."><span class="d-step">4</span> Base62 Encode <span class="d-metric latency">O(1)</span></div>
    </div>
    <div class="d-arrow-down">↓</div>
    <div class="d-label">short_code = "aB3kX9r" <span class="d-metric size">7 chars</span></div>
    <div class="d-arrow-down">↓</div>
    <div class="d-box indigo" data-tip="Primary key on short_code. DynamoDB: single-digit ms writes. Postgres: append-only B-tree with sequential IDs."><span class="d-step">5</span> Database: INSERT (aB3kX9r, long_url)</div>
  </div>
</div>
<div class="d-group">
  <div class="d-group-title">Read Path <span class="d-metric latency">~5ms p99 (cache hit)</span></div>
  <div class="d-flow-v" style="align-items: center;">
    <div class="d-box blue"><span class="d-step">1</span> Client: GET /aB3kX9r</div>
    <div class="d-arrow-down">↓</div>
    <div class="d-flow">
      <div class="d-box gray" data-tip="CDN cache hit ratio ~90% for popular URLs. TTL 24h. Saves origin from 90% of read traffic."><span class="d-step">2</span> CDN (cached?) <span class="d-metric latency">~2ms</span></div>
      <div class="d-arrow">miss →</div>
      <div class="d-box red" data-tip="Redis read-through cache. Top 20% of URLs serve 80% of reads. ElastiCache r6g.large ~$92/mo."><span class="d-step">3</span> Redis (cached?) <span class="d-metric latency">~5ms</span></div>
      <div class="d-arrow">miss →</div>
      <div class="d-box indigo" data-tip="Only reached on double cache miss — roughly 1-2% of reads."><span class="d-step">4</span> Database (fallback) <span class="d-metric latency">~10ms</span></div>
    </div>
    <div class="d-arrow-down">↓</div>
    <div class="d-box green"><span class="d-step">5</span> HTTP 301 → long URL</div>
  </div>
</div>
<div class="d-caption">Write path is collision-free by design — Base62 bijectively maps unique IDs. Read path optimizes for the 80/20 rule: hot URLs are served from CDN or Redis, keeping DB load minimal.</div>`,
	})

	// ---------------------------------------------------------------
	// Bloom Filter (2 diagrams)
	// ---------------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "algo-bloom-filter-insert-query",
		Title:       "Bloom Filter: Insert and Query",
		Description: "Bit array visualization showing how elements are inserted with k hash functions and how queries produce true positives, true negatives, and false positives.",
		ContentFile: "algorithms/bloom-filter",
		Type:        TypeHTML,
		HTML: `<div class="d-label" style="margin-bottom: 0.5rem;">Bit Array (m = 16 bits), k = 3 hash functions <span class="d-metric size">16 bits</span> <span class="d-metric latency">O(k) per op</span></div>
<div class="d-group" style="margin-bottom: 1rem;">
  <div class="d-group-title"><span class="d-step">1</span> Initial state: all zeros</div>
  <div class="d-bit-array" data-tip="In production, m is sized by: m = -(n * ln p) / (ln 2)^2. For 1B items at 1% FP rate: m = 1.2 GB.">
    <div class="d-bit off">0</div><div class="d-bit off">1</div><div class="d-bit off">2</div><div class="d-bit off">3</div>
    <div class="d-bit off">4</div><div class="d-bit off">5</div><div class="d-bit off">6</div><div class="d-bit off">7</div>
    <div class="d-bit off">8</div><div class="d-bit off">9</div><div class="d-bit off">10</div><div class="d-bit off">11</div>
    <div class="d-bit off">12</div><div class="d-bit off">13</div><div class="d-bit off">14</div><div class="d-bit off">15</div>
  </div>
</div>
<div class="d-group" style="margin-bottom: 1rem;">
  <div class="d-group-title"><span class="d-step">2</span> After INSERT "hello": h1=2, h2=7, h3=13</div>
  <div class="d-bit-array" data-tip="Each hash function independently maps the input to a bit position. Uses murmur3 or xxhash — fast, uniform distribution.">
    <div class="d-bit off">0</div><div class="d-bit off">1</div><div class="d-bit on">2</div><div class="d-bit off">3</div>
    <div class="d-bit off">4</div><div class="d-bit off">5</div><div class="d-bit off">6</div><div class="d-bit on">7</div>
    <div class="d-bit off">8</div><div class="d-bit off">9</div><div class="d-bit off">10</div><div class="d-bit off">11</div>
    <div class="d-bit off">12</div><div class="d-bit on">13</div><div class="d-bit off">14</div><div class="d-bit off">15</div>
  </div>
</div>
<div class="d-group" style="margin-bottom: 1rem;">
  <div class="d-group-title"><span class="d-step">3</span> After INSERT "world": h1=1, h2=7 (shared!), h3=11</div>
  <div class="d-bit-array" data-tip="Bit 7 is shared between 'hello' and 'world'. Shared bits increase false positive probability. Cannot delete — unsetting bit 7 would break 'hello' queries.">
    <div class="d-bit off">0</div><div class="d-bit on">1</div><div class="d-bit on">2</div><div class="d-bit off">3</div>
    <div class="d-bit off">4</div><div class="d-bit off">5</div><div class="d-bit off">6</div><div class="d-bit on">7</div>
    <div class="d-bit off">8</div><div class="d-bit off">9</div><div class="d-bit off">10</div><div class="d-bit on">11</div>
    <div class="d-bit off">12</div><div class="d-bit on">13</div><div class="d-bit off">14</div><div class="d-bit off">15</div>
  </div>
</div>
<div class="d-flow" style="flex-wrap: wrap; gap: 0.5rem; margin-top: 0.5rem;">
  <div class="d-box green" data-tip="All 3 bits set — element may be in set. Bloom filters never produce false negatives."><span class="d-step">4</span> QUERY "hello": bits 2,7,13 all = 1 → POSSIBLY IN SET</div>
  <div class="d-box blue" data-tip="Any single bit = 0 is conclusive proof of absence. This is the key property that makes Bloom filters useful."><span class="d-step">5</span> QUERY "cat": bit[5]=0 → DEFINITELY NOT IN SET</div>
  <div class="d-box red" data-tip="Bits 1,7,11 happen to be set by other insertions. FP rate = (1 - e^(-kn/m))^k ≈ 1% with proper sizing."><span class="d-step">6</span> QUERY "fake": bits 1,7,11 all = 1 → FALSE POSITIVE!</div>
</div>
<div class="d-legend">
  <div class="d-bit on" style="display:inline-block; width:18px; height:18px; vertical-align:middle;"></div> = bit set (1)
  <div class="d-bit off" style="display:inline-block; width:18px; height:18px; vertical-align:middle; margin-left:1rem;"></div> = bit unset (0)
</div>
<div class="d-caption">No false negatives, ever. False positives are tunable: optimal k = (m/n) * ln 2. Cannot delete elements — use counting Bloom filters (4 bits per slot) if deletion is needed.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-bloom-filter-cache-layer",
		Title:       "Bloom Filter as Cache Optimization Layer",
		Description: "Flow diagram showing a Bloom filter sitting in front of a database to skip lookups for keys that are definitely absent.",
		ContentFile: "algorithms/bloom-filter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v" style="align-items: center;">
  <div class="d-box blue"><span class="d-step">1</span> Request: "Does key X exist?"</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-box purple" data-tip="Entire filter lives in memory. At 10 bits/element with k=7 hashes, FP rate is 0.82%. Checking 7 bits takes ~100ns vs ~5ms for a DB query."><span class="d-step">2</span> Bloom Filter <span class="d-metric size">1.2 GB for 1B elements</span> <span class="d-metric latency">O(k) in-memory</span></div>
  <div class="d-arrow-down">↓</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box green" data-tip="If any of the k bits is 0, the element was never inserted. This guarantee is absolute — zero false negatives."><span class="d-step">3a</span> Definitely NOT here</div>
      <div class="d-arrow-down">↓</div>
      <div class="d-box green" data-tip="For URL shorteners, 99% of random short codes don't exist. Bloom filter skips 99% of DB reads — saving ~$50K/mo at scale.">SKIP DB (99% of cases) <span class="d-metric latency">~100ns</span></div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box amber" data-tip="All k bits are 1, but this could be coincidence from other inserts. Must verify against the actual data store."><span class="d-step">3b</span> Maybe here</div>
      <div class="d-arrow-down">↓</div>
      <div class="d-box gray"><span class="d-step">4</span> Database Lookup <span class="d-metric latency">~5ms</span></div>
      <div class="d-arrow-down">↓</div>
      <div class="d-branch">
        <div class="d-branch-arm">
          <div class="d-box green">FOUND (true positive)</div>
        </div>
        <div class="d-branch-arm">
          <div class="d-box red" data-tip="~1% of 'maybe' responses are false positives. Cost: one wasted DB query. Acceptable trade-off for 99% fewer total queries.">NOT FOUND (false positive) <span class="d-metric latency">~1% rate</span></div>
        </div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Bloom filters eliminate 99% of negative lookups at ~100ns each vs ~5ms DB queries. Net effect: 50x throughput improvement for existence checks. Trade-off: 1.2 GB memory for 1B elements.</div>`,
	})

	// ---------------------------------------------------------------
	// Consistent Hashing (3 diagrams)
	// ---------------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "algo-consistent-hash-ring",
		Title:       "Consistent Hash Ring",
		Description: "Circular hash ring with nodes and a key, showing clockwise lookup to determine key ownership and reassignment on node removal.",
		ContentFile: "algorithms/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-ring" style="width: 260px; height: 260px; margin: 0 auto;">
  <div class="d-ring-node blue" style="top: 2%; left: 62%;" data-tip="Position = hash(node_id) % 2^32. Uses SHA-256 or MD5 for uniform distribution across the ring.">Node A (800)</div>
  <div class="d-ring-node green" style="top: 88%; left: 35%;" data-tip="Each node owns the range from its predecessor's position to its own position (clockwise).">Node B (200)</div>
  <div class="d-ring-node purple" style="top: 65%; left: 5%;">Node C (400)</div>
  <div class="d-ring-node amber" style="top: 42%; left: 12%;" data-tip="Key lookup: hash the key, walk clockwise to first node. O(log N) with sorted node list + binary search.">key: user:42 (520)</div>
</div>
<div class="d-flow-v" style="align-items: center; margin-top: 1rem; gap: 0.5rem;">
  <div class="d-label"><span class="d-step">1</span> key "user:42" hashes to position 520. Walk clockwise: next node is A at 800. <span class="d-metric latency">O(log N) lookup</span></div>
  <div class="d-box blue" style="max-width: 300px;">Node A owns "user:42"</div>
  <div class="d-label"><span class="d-step">2</span> If Node A is removed: walk clockwise from 520, wrap around to Node B at 200.</div>
  <div class="d-label"><span class="d-step">3</span> Only keys between C(400) and A(800) are reassigned. Keys owned by B and C do not move.</div>
</div>
<div class="d-caption">On node failure, only K/N keys migrate (K = total keys, N = nodes). With modulo hashing, ALL keys would need remapping. This is why consistent hashing is essential for distributed caches and databases.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-consistent-hash-vnodes",
		Title:       "Virtual Nodes on the Ring",
		Description: "Comparison of load distribution with and without virtual nodes, showing how vnodes reduce standard deviation from 58% to 5%.",
		ContentFile: "algorithms/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-group" style="margin-bottom: 1.5rem;">
  <div class="d-group-title">Without virtual nodes (3 physical nodes) - Unbalanced!</div>
  <div class="d-bit-array" data-tip="With only 3 points on a 2^32 ring, random placement creates huge variance. One node may own 60% of the key space.">
    <div class="d-bit on" style="background: var(--blue, #3b82f6); flex: 6;">A: 60%</div>
    <div class="d-bit on" style="background: var(--green, #22c55e); flex: 1.5;">B: 15%</div>
    <div class="d-bit on" style="background: var(--purple, #a855f7); flex: 2.5;">C: 25%</div>
  </div>
  <div class="d-label">Node A is overloaded! <span class="d-metric latency">Std dev ~58%</span></div>
</div>
<div class="d-group">
  <div class="d-group-title">With 150 virtual nodes per physical node - Balanced <span class="d-metric size">450 ring positions</span></div>
  <div class="d-bit-array" data-tip="Each physical node maps to 150 positions on the ring via hash(node_id + vnode_index). More vnodes = better distribution but higher memory for the ring map (~450 entries in sorted array).">
    <div class="d-bit on" style="background: var(--blue, #3b82f6); flex: 1;"></div>
    <div class="d-bit on" style="background: var(--green, #22c55e); flex: 1;"></div>
    <div class="d-bit on" style="background: var(--purple, #a855f7); flex: 1;"></div>
    <div class="d-bit on" style="background: var(--blue, #3b82f6); flex: 1;"></div>
    <div class="d-bit on" style="background: var(--purple, #a855f7); flex: 1;"></div>
    <div class="d-bit on" style="background: var(--green, #22c55e); flex: 1;"></div>
    <div class="d-bit on" style="background: var(--blue, #3b82f6); flex: 1;"></div>
    <div class="d-bit on" style="background: var(--green, #22c55e); flex: 1;"></div>
    <div class="d-bit on" style="background: var(--purple, #a855f7); flex: 1;"></div>
  </div>
  <div class="d-flow" style="margin-top: 0.5rem; justify-content: center;">
    <div class="d-box blue">A: ~33%</div>
    <div class="d-box green">B: ~33%</div>
    <div class="d-box purple">C: ~33%</div>
  </div>
  <div class="d-label">Standard deviation drops from ~58% to ~5%. <span class="d-metric latency">Std dev ~5%</span></div>
</div>
<div class="d-legend">
  <span style="display:inline-block; width:12px; height:12px; background:var(--blue, #3b82f6); vertical-align:middle;"></span> Node A
  <span style="display:inline-block; width:12px; height:12px; background:var(--green, #22c55e); vertical-align:middle; margin-left:1rem;"></span> Node B
  <span style="display:inline-block; width:12px; height:12px; background:var(--purple, #a855f7); vertical-align:middle; margin-left:1rem;"></span> Node C
</div>
<div class="d-caption">150 vnodes per physical node is the sweet spot (used by Cassandra). Fewer vnodes = worse balance. More vnodes = diminishing returns + larger ring metadata. Memory cost: ~6 bytes per vnode entry in a sorted array.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-consistent-hash-rebalance",
		Title:       "Node Addition: Before and After",
		Description: "Visualization of key redistribution when adding a node to a consistent hash ring, showing only K/N keys move instead of all keys.",
		ContentFile: "algorithms/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-group" style="margin-bottom: 1.5rem;">
  <div class="d-group-title">Before (3 nodes, even distribution with vnodes)</div>
  <div class="d-bit-array">
    <div class="d-bit on" style="background: var(--blue, #3b82f6); flex: 1;">A: 33%</div>
    <div class="d-bit on" style="background: var(--green, #22c55e); flex: 1;">B: 33%</div>
    <div class="d-bit on" style="background: var(--purple, #a855f7); flex: 1;">C: 33%</div>
  </div>
</div>
<div class="d-arrow-down" style="text-align: center;">↓ Add Node D</div>
<div class="d-group" style="margin-top: 0.5rem;">
  <div class="d-group-title">After adding D (4 nodes)</div>
  <div class="d-bit-array">
    <div class="d-bit on" style="background: var(--blue, #3b82f6); flex: 25;">A: 25%</div>
    <div class="d-bit on" style="background: var(--amber, #f59e0b); flex: 8;">D: 8%</div>
    <div class="d-bit on" style="background: var(--green, #22c55e); flex: 25;">B: 25%</div>
    <div class="d-bit on" style="background: var(--purple, #a855f7); flex: 25;">C: 25%</div>
  </div>
  <div class="d-label" style="margin-top: 0.5rem;">Node D takes ~8% of keys from B (its clockwise neighbor). Nodes A and C are unaffected. Total keys moved: ~25% (K/4), not 75% like modulo.</div>
</div>`,
	})

	// ---------------------------------------------------------------
	// Snowflake ID (3 diagrams)
	// ---------------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "algo-snowflake-bit-layout",
		Title:       "Snowflake ID: 64-Bit Layout",
		Description: "Bitfield diagram showing the 64-bit Snowflake ID structure: 1 sign bit, 41 timestamp bits, 10 machine ID bits, and 12 sequence bits.",
		ContentFile: "algorithms/snowflake-id",
		Type:        TypeHTML,
		HTML: `<div class="d-bitfield">
  <div class="d-bitfield-segment gray" data-tip="Reserved for positive values. Snowflake IDs are always positive to avoid issues with signed integer handling in JavaScript and databases.">
    <div class="d-bitfield-bits"><span class="d-metric size">1 bit</span></div>
    <div class="d-bitfield-name">Sign (0)</div>
  </div>
  <div class="d-bitfield-segment blue" data-tip="Milliseconds since custom epoch (e.g., 2015-01-01). 2^41 ms = 69.7 years. Custom epoch extends usable range vs Unix epoch (which would exhaust by 2039).">
    <div class="d-bitfield-bits"><span class="d-metric size">41 bits</span></div>
    <div class="d-bitfield-name">Timestamp (ms since epoch)</div>
  </div>
  <div class="d-bitfield-segment green" data-tip="5 bits datacenter + 5 bits machine = 32 DCs x 32 machines. Or all 10 bits for machines = 1,024 total. Assigned at deploy time via ZooKeeper or config.">
    <div class="d-bitfield-bits"><span class="d-metric size">10 bits</span></div>
    <div class="d-bitfield-name">Machine ID</div>
  </div>
  <div class="d-bitfield-segment purple" data-tip="Resets to 0 each millisecond. If 4,096 IDs exhausted within 1ms, the generator waits (spins) until the next ms. Handles clock skew by refusing to generate if clock goes backward.">
    <div class="d-bitfield-bits"><span class="d-metric size">12 bits</span></div>
    <div class="d-bitfield-name">Sequence</div>
  </div>
</div>
<div class="d-flow" style="margin-top: 1rem; justify-content: center; flex-wrap: wrap; gap: 1rem;">
  <div class="d-box gray" data-tip="Ensures positive int64 in all languages. JavaScript's Number.MAX_SAFE_INTEGER = 2^53, so Snowflake IDs fit.">Sign: always 0</div>
  <div class="d-box blue" data-tip="Custom epoch starting 2015-01-01 gives range until ~2084. Twitter's epoch: 2010-11-04.">~69 years range</div>
  <div class="d-box green" data-tip="1,024 machines generating independently with zero coordination. No single point of failure."><span class="d-metric size">1,024 machines</span></div>
  <div class="d-box purple" data-tip="4,096 per ms per machine = 4,096,000/sec/machine. 1,024 machines = 4 billion IDs/sec globally."><span class="d-metric latency">4,096 IDs/ms</span></div>
</div>
<div class="d-label" style="margin-top: 0.75rem; text-align: center;">Total: 1 + 41 + 10 + 12 = 64 bits. Fits in a single int64. Max: 4,096,000 IDs/sec per machine.</div>
<div class="d-caption">Time-ordered IDs enable B-tree append-only writes (no page splits), range queries by time, and rough chronological sorting without a secondary index. Trade-off: clock skew can cause duplicate IDs if NTP jumps backward.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-snowflake-btree-insert",
		Title:       "B-Tree Insert: Sequential vs Random IDs",
		Description: "Side-by-side comparison of B-tree insert behavior with Snowflake sequential IDs versus UUID v4 random IDs, showing page split differences.",
		ContentFile: "algorithms/snowflake-id",
		Type:        TypeHTML,
		HTML: `<div class="d-cols" style="grid-template-columns: 1fr 1fr; gap: 2rem;">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Snowflake IDs (sequential)</div>
      <div class="d-flow-v" style="align-items: center;">
        <div class="d-label">All inserts append to rightmost leaf</div>
        <div class="d-flow">
          <div class="d-box blue">ID 1</div>
          <div class="d-box blue">ID 2</div>
          <div class="d-box blue">ID 3</div>
          <div class="d-box green">ID 4 ← append</div>
        </div>
        <div class="d-label">No page splits. Buffer pool stays warm. Write amplification = 1x.</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">UUID v4 (random)</div>
      <div class="d-flow-v" style="align-items: center;">
        <div class="d-label">Inserts scatter across all leaf pages</div>
        <div class="d-flow">
          <div class="d-box red">ID x</div>
          <div class="d-box gray">???</div>
          <div class="d-box red">ID y</div>
          <div class="d-box gray">???</div>
        </div>
        <div class="d-label">Constant page splits. Cold buffer pool. Write amplification = 5-10x.</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-snowflake-distributed-arch",
		Title:       "Snowflake ID in Distributed Architecture",
		Description: "Architecture diagram showing multiple app servers generating Snowflake IDs independently with no coordination, writing to sharded databases.",
		ContentFile: "algorithms/snowflake-id",
		Type:        TypeHTML,
		HTML: `<div class="d-flow" style="justify-content: center;">
  <div class="d-box blue">App Server 1<br>machine_id=0<br>Snowflake Gen</div>
  <div class="d-box blue">App Server 2<br>machine_id=1<br>Snowflake Gen</div>
  <div class="d-box blue">App Server 3<br>machine_id=2<br>Snowflake Gen</div>
</div>
<div class="d-flow-v" style="align-items: center; margin-top: 0.5rem;">
  <div class="d-arrow-down">↓</div>
  <div class="d-label">Each server generates IDs independently. No coordination.</div>
  <div class="d-arrow-down">↓</div>
</div>
<div class="d-flow" style="justify-content: center;">
  <div class="d-box green">DB Shard 1<br>(IDs sorted in B-tree)</div>
  <div class="d-box green">DB Shard 2<br>(IDs sorted in B-tree)</div>
</div>
<div class="d-label" style="text-align: center; margin-top: 0.75rem;">No single point of failure. Each server generates IDs independently.</div>`,
	})

	// ---------------------------------------------------------------
	// Token Bucket (3 diagrams)
	// ---------------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "algo-token-bucket-fill-drain",
		Title:       "Token Bucket: Fill and Drain",
		Description: "Visual timeline of a token bucket filling and draining, showing burst consumption, rejection when empty, and gradual refill over time.",
		ContentFile: "algorithms/token-bucket",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-box blue" data-tip="Max burst size = capacity. Choose based on acceptable spike: 10 for API endpoints, 100+ for CDN edge rate limiting.">Capacity: 10 tokens <span class="d-metric size">max burst</span></div>
  <div class="d-box green" data-tip="Sustained throughput ceiling. tokens_added = rate * elapsed_seconds. Lazy evaluation: calculate on each request, don't use timers.">Refill Rate: 2 tokens/sec <span class="d-metric latency">sustained limit</span></div>
</div>
<div class="d-flow-v" style="gap: 0.75rem; margin-top: 1rem;">
  <div class="d-row">
    <div class="d-box green" style="min-width: 160px;"><span class="d-step">1</span> Time 0s: 10/10 (full)</div>
    <div class="d-bit-array">
      <div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div>
      <div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div>
    </div>
  </div>
  <div class="d-label"><span class="d-step">2</span> Burst of 9 requests!</div>
  <div class="d-row">
    <div class="d-box amber" style="min-width: 160px;">Time 0s: 1/10</div>
    <div class="d-bit-array" data-tip="All 9 requests served instantly — this is the burst-friendly behavior that distinguishes token bucket from leaky bucket.">
      <div class="d-bit on"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
    </div>
  </div>
  <div class="d-label"><span class="d-step">3</span> Request!</div>
  <div class="d-row">
    <div class="d-box red" style="min-width: 160px;">Time 0s: 0/10</div>
    <div class="d-bit-array">
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
    </div>
  </div>
  <div class="d-label"><span class="d-step">4</span> Request → REJECTED (429 Too Many Requests)</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-label">...2 seconds pass, +4 tokens refilled (rate x elapsed)...</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-row">
    <div class="d-box green" style="min-width: 160px;"><span class="d-step">5</span> Time 2s: 4/10 (refilled)</div>
    <div class="d-bit-array" data-tip="Refill is lazy: tokens = min(capacity, tokens + rate * (now - last_refill)). No background timer needed.">
      <div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit off"></div>
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
    </div>
  </div>
  <div class="d-label"><span class="d-step">6</span> Request → allowed</div>
  <div class="d-row">
    <div class="d-box blue" style="min-width: 160px;">Time 2s: 3/10</div>
    <div class="d-bit-array">
      <div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit off"></div><div class="d-bit off"></div>
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
    </div>
  </div>
</div>
<div class="d-legend">
  <div class="d-bit on" style="display:inline-block; width:18px; height:18px; vertical-align:middle;"></div> = available token
  <div class="d-bit off" style="display:inline-block; width:18px; height:18px; vertical-align:middle; margin-left:1rem;"></div> = consumed/empty slot
</div>
<div class="d-caption">Burst-friendly: allows up to capacity requests instantly, then throttles to refill_rate. Two parameters per bucket: capacity (burst size) and rate (sustained throughput). State: just 2 values (token count + last refill timestamp).</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-token-vs-leaky-bucket",
		Title:       "Token Bucket vs Leaky Bucket",
		Description: "Side-by-side comparison of token bucket (allows bursts) versus leaky bucket (smooths output to constant drain rate).",
		ContentFile: "algorithms/token-bucket",
		Type:        TypeHTML,
		HTML: `<div class="d-cols" style="grid-template-columns: 1fr 1fr; gap: 2rem;">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Token Bucket (allows bursts)</div>
      <div class="d-flow-v" style="align-items: center;">
        <div class="d-label">Incoming: bursty traffic</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box blue">Tokens: 5</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-label">Outgoing: bursty, passes instantly</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Leaky Bucket (smooths output)</div>
      <div class="d-flow-v" style="align-items: center;">
        <div class="d-label">Incoming: bursty traffic</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box green">Queue: 5</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-label">Outgoing: constant drain rate</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-token-bucket-api-arch",
		Title:       "Token Bucket in API Architecture",
		Description: "Architecture diagram showing token bucket rate limiting as middleware between API server and Redis, with allow/reject branching and response headers.",
		ContentFile: "algorithms/token-bucket",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-box blue"><span class="d-step">1</span> Client</div>
  <div class="d-arrow">→</div>
  <div class="d-box gray" data-tip="ALB distributes across API server fleet. Rate limiting happens per-server with shared Redis state — not at the ALB level."><span class="d-step">2</span> ALB</div>
  <div class="d-arrow">→</div>
  <div class="d-box purple" data-tip="Stateless — rate limit state lives entirely in Redis. Any server can handle any client's request."><span class="d-step">3</span> API Server</div>
</div>
<div class="d-flow-v" style="align-items: center; margin-top: 1rem;">
  <div class="d-arrow-down">↓</div>
  <div class="d-flow">
    <div class="d-box amber" data-tip="Runs before business logic. On Redis failure, fail-open (allow all) is safer than fail-closed (reject all) — a few seconds of unlimited traffic beats a full outage."><span class="d-step">4</span> Rate Limit Middleware</div>
    <div class="d-arrow">→</div>
    <div class="d-box red" data-tip="Lua script executes atomically: read token count + last_refill, calculate new tokens, decrement, write back. Single round trip, no race conditions. ~0.1ms per call."><span class="d-step">5</span> Redis (Lua script) <span class="d-metric latency">~0.1ms</span></div>
  </div>
  <div class="d-label">atomic HMGET + HMSET (single Lua eval, no race conditions)</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box green" data-tip="Token decremented. Request proceeds to business logic. Response includes remaining token count in headers."><span class="d-step">6a</span> ALLOW</div>
      <div class="d-label">continue processing</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box red" data-tip="No tokens available. Return 429 immediately — don't waste compute on business logic. Include Retry-After header so clients can back off intelligently."><span class="d-step">6b</span> REJECT</div>
      <div class="d-label">429 + Retry-After</div>
    </div>
  </div>
</div>
<div class="d-group" style="margin-top: 1rem;">
  <div class="d-group-title">Response Headers (IETF draft-polli-ratelimit)</div>
  <div class="d-flow">
    <div class="d-box gray" data-tip="Bucket capacity. Lets clients know their total allowance per window.">X-RateLimit-Limit: 100</div>
    <div class="d-box gray" data-tip="Tokens left. Clients can throttle proactively before hitting 0.">X-RateLimit-Remaining: 73</div>
    <div class="d-box gray" data-tip="Seconds until tokens refill. Only sent on 429. Smart clients use exponential backoff + jitter.">Retry-After: 8 (on 429)</div>
  </div>
</div>
<div class="d-caption">Redis key per user: rate_limit:{user_id}. Two fields: tokens (float64) and last_refill (timestamp). Total Redis memory: ~100 bytes/user. 10M users = 1 GB. Lua script ensures atomicity without distributed locks.</div>`,
	})
}

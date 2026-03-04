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
  <div class="d-box blue">Input: 123456789</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-group">
    <div class="d-group-title">Repeated divide-by-62</div>
    <div class="d-flow-v" style="gap: 0.25rem;">
      <div class="d-flow">
        <div class="d-box gray">123456789 % 62 = 17</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber">r</div>
      </div>
      <div class="d-flow">
        <div class="d-box gray">1991239 % 62 = 25</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber">z</div>
      </div>
      <div class="d-flow">
        <div class="d-box gray">32116 % 62 = 4</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber">e</div>
      </div>
      <div class="d-flow">
        <div class="d-box gray">517 % 62 = 21</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber">v</div>
      </div>
      <div class="d-flow">
        <div class="d-box gray">8 % 62 = 8</div>
        <div class="d-arrow">→</div>
        <div class="d-box amber">i</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">↓</div>
  <div class="d-label">Read remainders bottom-to-top</div>
  <div class="d-box green">Result: "ivezr"</div>
  <div class="d-label" style="margin-top: 0.5rem;">Decode: 8*62^4 + 21*62^3 + 4*62^2 + 25*62 + 17 = 123456789</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-base62-url-shortener",
		Title:       "Base62 in URL Shortener Architecture",
		Description: "Write and read paths of a URL shortener showing how Base62 encoding integrates with the ID generator, database, and cache layers.",
		ContentFile: "algorithms/base62-encoding",
		Type:        TypeHTML,
		HTML: `<div class="d-group" style="margin-bottom: 1.5rem;">
  <div class="d-group-title">Write Path</div>
  <div class="d-flow-v" style="align-items: center;">
    <div class="d-box blue">Client: POST /shorten</div>
    <div class="d-arrow-down">↓</div>
    <div class="d-flow">
      <div class="d-box purple">API Server</div>
      <div class="d-arrow">→</div>
      <div class="d-box green">ID Generator (counter/KGS)</div>
      <div class="d-arrow">→</div>
      <div class="d-box amber">Base62 Encode</div>
    </div>
    <div class="d-arrow-down">↓</div>
    <div class="d-label">short_code = "aB3kX9r"</div>
    <div class="d-arrow-down">↓</div>
    <div class="d-box indigo">Database: INSERT (aB3kX9r, long_url)</div>
  </div>
</div>
<div class="d-group">
  <div class="d-group-title">Read Path</div>
  <div class="d-flow-v" style="align-items: center;">
    <div class="d-box blue">Client: GET /aB3kX9r</div>
    <div class="d-arrow-down">↓</div>
    <div class="d-flow">
      <div class="d-box gray">CDN (cached?)</div>
      <div class="d-arrow">→</div>
      <div class="d-box red">Redis (cached?)</div>
      <div class="d-arrow">→</div>
      <div class="d-box indigo">Database (fallback)</div>
    </div>
    <div class="d-arrow-down">↓</div>
    <div class="d-box green">HTTP 301 → long URL</div>
  </div>
</div>`,
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
		HTML: `<div class="d-label" style="margin-bottom: 0.5rem;">Bit Array (m = 16 bits), k = 3 hash functions</div>
<div class="d-group" style="margin-bottom: 1rem;">
  <div class="d-group-title">Initial state: all zeros</div>
  <div class="d-bit-array">
    <div class="d-bit off">0</div><div class="d-bit off">1</div><div class="d-bit off">2</div><div class="d-bit off">3</div>
    <div class="d-bit off">4</div><div class="d-bit off">5</div><div class="d-bit off">6</div><div class="d-bit off">7</div>
    <div class="d-bit off">8</div><div class="d-bit off">9</div><div class="d-bit off">10</div><div class="d-bit off">11</div>
    <div class="d-bit off">12</div><div class="d-bit off">13</div><div class="d-bit off">14</div><div class="d-bit off">15</div>
  </div>
</div>
<div class="d-group" style="margin-bottom: 1rem;">
  <div class="d-group-title">After INSERT "hello": h1=2, h2=7, h3=13</div>
  <div class="d-bit-array">
    <div class="d-bit off">0</div><div class="d-bit off">1</div><div class="d-bit on">2</div><div class="d-bit off">3</div>
    <div class="d-bit off">4</div><div class="d-bit off">5</div><div class="d-bit off">6</div><div class="d-bit on">7</div>
    <div class="d-bit off">8</div><div class="d-bit off">9</div><div class="d-bit off">10</div><div class="d-bit off">11</div>
    <div class="d-bit off">12</div><div class="d-bit on">13</div><div class="d-bit off">14</div><div class="d-bit off">15</div>
  </div>
</div>
<div class="d-group" style="margin-bottom: 1rem;">
  <div class="d-group-title">After INSERT "world": h1=1, h2=7 (shared!), h3=11</div>
  <div class="d-bit-array">
    <div class="d-bit off">0</div><div class="d-bit on">1</div><div class="d-bit on">2</div><div class="d-bit off">3</div>
    <div class="d-bit off">4</div><div class="d-bit off">5</div><div class="d-bit off">6</div><div class="d-bit on">7</div>
    <div class="d-bit off">8</div><div class="d-bit off">9</div><div class="d-bit off">10</div><div class="d-bit on">11</div>
    <div class="d-bit off">12</div><div class="d-bit on">13</div><div class="d-bit off">14</div><div class="d-bit off">15</div>
  </div>
</div>
<div class="d-flow" style="flex-wrap: wrap; gap: 0.5rem; margin-top: 0.5rem;">
  <div class="d-box green">QUERY "hello": bits 2,7,13 all = 1 → POSSIBLY IN SET</div>
  <div class="d-box blue">QUERY "cat": bit[5]=0 → DEFINITELY NOT IN SET</div>
  <div class="d-box red">QUERY "fake": bits 1,7,11 all = 1 → FALSE POSITIVE!</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-bloom-filter-cache-layer",
		Title:       "Bloom Filter as Cache Optimization Layer",
		Description: "Flow diagram showing a Bloom filter sitting in front of a database to skip lookups for keys that are definitely absent.",
		ContentFile: "algorithms/bloom-filter",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v" style="align-items: center;">
  <div class="d-box blue">Request: "Does key X exist?"</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-box purple">Bloom Filter (1.2 GB for 1B elements, O(k) in-memory)</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box green">Definitely NOT here</div>
      <div class="d-arrow-down">↓</div>
      <div class="d-box green">SKIP DB (99% of cases)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box amber">Maybe here</div>
      <div class="d-arrow-down">↓</div>
      <div class="d-box gray">Database Lookup</div>
      <div class="d-arrow-down">↓</div>
      <div class="d-branch">
        <div class="d-branch-arm">
          <div class="d-box green">FOUND (true positive)</div>
        </div>
        <div class="d-branch-arm">
          <div class="d-box red">NOT FOUND (false positive)</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
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
  <div class="d-ring-node blue" style="top: 2%; left: 62%;">Node A (800)</div>
  <div class="d-ring-node green" style="top: 88%; left: 35%;">Node B (200)</div>
  <div class="d-ring-node purple" style="top: 65%; left: 5%;">Node C (400)</div>
  <div class="d-ring-node amber" style="top: 42%; left: 12%;">key: user:42 (520)</div>
</div>
<div class="d-flow-v" style="align-items: center; margin-top: 1rem; gap: 0.5rem;">
  <div class="d-label">key "user:42" hashes to position 520. Walk clockwise: next node is A at 800.</div>
  <div class="d-box blue" style="max-width: 300px;">Node A owns "user:42"</div>
  <div class="d-label">If Node A is removed: walk clockwise from 520, wrap around to Node B at 200.</div>
  <div class="d-label">Only keys between C(400) and A(800) are reassigned. Keys owned by B and C do not move.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "algo-consistent-hash-vnodes",
		Title:       "Virtual Nodes on the Ring",
		Description: "Comparison of load distribution with and without virtual nodes, showing how vnodes reduce standard deviation from 58% to 5%.",
		ContentFile: "algorithms/consistent-hashing",
		Type:        TypeHTML,
		HTML: `<div class="d-group" style="margin-bottom: 1.5rem;">
  <div class="d-group-title">Without virtual nodes (3 physical nodes) - Unbalanced!</div>
  <div class="d-bit-array">
    <div class="d-bit on" style="background: var(--blue, #3b82f6); flex: 6;">A: 60%</div>
    <div class="d-bit on" style="background: var(--green, #22c55e); flex: 1.5;">B: 15%</div>
    <div class="d-bit on" style="background: var(--purple, #a855f7); flex: 2.5;">C: 25%</div>
  </div>
  <div class="d-label">Node A is overloaded!</div>
</div>
<div class="d-group">
  <div class="d-group-title">With 150 virtual nodes per physical node - Balanced</div>
  <div class="d-bit-array">
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
  <div class="d-label">Standard deviation drops from ~58% to ~5%.</div>
</div>`,
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
  <div class="d-bitfield-segment gray">
    <div class="d-bitfield-bits">1 bit</div>
    <div class="d-bitfield-name">Sign (0)</div>
  </div>
  <div class="d-bitfield-segment blue">
    <div class="d-bitfield-bits">41 bits</div>
    <div class="d-bitfield-name">Timestamp (ms since epoch)</div>
  </div>
  <div class="d-bitfield-segment green">
    <div class="d-bitfield-bits">10 bits</div>
    <div class="d-bitfield-name">Machine ID</div>
  </div>
  <div class="d-bitfield-segment purple">
    <div class="d-bitfield-bits">12 bits</div>
    <div class="d-bitfield-name">Sequence</div>
  </div>
</div>
<div class="d-flow" style="margin-top: 1rem; justify-content: center; flex-wrap: wrap; gap: 1rem;">
  <div class="d-box gray">Sign: always 0</div>
  <div class="d-box blue">~69 years range</div>
  <div class="d-box green">1,024 machines</div>
  <div class="d-box purple">4,096 IDs/ms</div>
</div>
<div class="d-label" style="margin-top: 0.75rem; text-align: center;">Total: 1 + 41 + 10 + 12 = 64 bits. Fits in a single int64. Max: 4,096,000 IDs/sec per machine.</div>`,
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
  <div class="d-box blue">Capacity: 10 tokens</div>
  <div class="d-box green">Refill Rate: 2 tokens/sec</div>
</div>
<div class="d-flow-v" style="gap: 0.75rem; margin-top: 1rem;">
  <div class="d-row">
    <div class="d-box green" style="min-width: 160px;">Time 0s: 10/10 (full)</div>
    <div class="d-bit-array">
      <div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div>
      <div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div>
    </div>
  </div>
  <div class="d-label">Burst of 9 requests!</div>
  <div class="d-row">
    <div class="d-box amber" style="min-width: 160px;">Time 0s: 1/10</div>
    <div class="d-bit-array">
      <div class="d-bit on"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
    </div>
  </div>
  <div class="d-label">Request!</div>
  <div class="d-row">
    <div class="d-box red" style="min-width: 160px;">Time 0s: 0/10</div>
    <div class="d-bit-array">
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
    </div>
  </div>
  <div class="d-label">Request → REJECTED (empty)</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-label">...2 seconds pass, +4 tokens refilled...</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-row">
    <div class="d-box green" style="min-width: 160px;">Time 2s: 4/10 (refilled)</div>
    <div class="d-bit-array">
      <div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit off"></div>
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
    </div>
  </div>
  <div class="d-label">Request!</div>
  <div class="d-row">
    <div class="d-box blue" style="min-width: 160px;">Time 2s: 3/10</div>
    <div class="d-bit-array">
      <div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit on"></div><div class="d-bit off"></div><div class="d-bit off"></div>
      <div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div><div class="d-bit off"></div>
    </div>
  </div>
</div>
<div class="d-row" style="margin-top: 1rem; gap: 1rem;">
  <div class="d-label">Allows bursts UP TO capacity. Sustained rate cannot exceed refill_rate.</div>
</div>`,
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
  <div class="d-box blue">Client</div>
  <div class="d-arrow">→</div>
  <div class="d-box gray">ALB</div>
  <div class="d-arrow">→</div>
  <div class="d-box purple">API Server</div>
</div>
<div class="d-flow-v" style="align-items: center; margin-top: 1rem;">
  <div class="d-arrow-down">↓</div>
  <div class="d-flow">
    <div class="d-box amber">Rate Limit Middleware</div>
    <div class="d-arrow">→</div>
    <div class="d-box red">Redis (Lua script)</div>
  </div>
  <div class="d-label">atomic HMGET + HMSET</div>
  <div class="d-arrow-down">↓</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box green">ALLOW</div>
      <div class="d-label">continue processing</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box red">REJECT</div>
      <div class="d-label">429 + Retry-After</div>
    </div>
  </div>
</div>
<div class="d-group" style="margin-top: 1rem;">
  <div class="d-group-title">Response Headers</div>
  <div class="d-flow">
    <div class="d-box gray">X-RateLimit-Limit: 100</div>
    <div class="d-box gray">X-RateLimit-Remaining: 73</div>
    <div class="d-box gray">Retry-After: 8 (on 429)</div>
  </div>
</div>`,
	})
}

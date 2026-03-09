package diagrams

func registerIDGenerator(r *Registry) {
	r.Register(&Diagram{
		Slug:        "idgen-requirements",
		Title:       "Requirements & Scale",
		Description: "Scale estimates, NFRs, and scope for distributed ID generator.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="64-bit integer fits in a long/int64. Sortable by time = range scans on primary key work efficiently.">Generate unique 64-bit IDs sortable by time</div>
        <div class="d-box green" data-tip="No coordination required between nodes for ID generation — each worker generates independently.">No central coordinator — fully distributed</div>
        <div class="d-box green" data-tip="IDs issued at time T are always numerically greater than IDs issued at T-1 (within same worker).">Monotonically increasing within a worker</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Non-Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="4096 IDs per ms per worker × 1024 workers = 4.2B IDs/sec system-wide. Single node = 4.09M/sec.">Throughput: 1M+ IDs/sec sustained</div>
        <div class="d-box blue" data-tip="ID generation is pure in-memory — no DB round-trip. Sub-microsecond per ID.">Latency: &lt;1ms per ID (target &lt;1μs)</div>
        <div class="d-box blue" data-tip="Stateless workers. ZooKeeper HA for worker ID registry. Multiple DCs.">Availability: 99.999% (no single SPOF)</div>
        <div class="d-box amber" data-tip="IDs must be unique across all datacenters, all workers, all time — including after restarts.">Uniqueness: globally unique, no collisions ever</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Estimates</div>
      <div class="d-flow-v">
        <div class="d-box indigo" data-tip="Twitter reported ~6K tweets/sec sustained with 150K peak. IDs needed for tweets, likes, follows, DMs.">Target: 1,000,000 IDs/sec peak</div>
        <div class="d-box indigo" data-tip="5 bits for datacenter ID → 2^5 = 32 datacenters maximum. Covers all major cloud regions.">Datacenters: 10 active, 32 max (5 bits)</div>
        <div class="d-box indigo" data-tip="5 bits for worker ID → 32 workers per DC. Each worker handles its own sequence counter.">Workers per DC: 32 (5 bits)</div>
        <div class="d-box indigo" data-tip="12 bits → 2^12 = 4096 IDs per millisecond per worker. At 1ms cadence = 4.09M IDs/sec per worker.">Sequence: 4096 IDs/ms per worker</div>
        <div class="d-box gray" data-tip="41 bits of millisecond timestamp → 2^41 ms = 2,199,023,255,552 ms ÷ 1000 ÷ 60 ÷ 60 ÷ 24 ÷ 365 ≈ 69.7 years.">ID lifetime: ~69 years from epoch</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">What We Are NOT Building</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="UUID v4 is random — not sortable. UUID v7 is time-sorted but 128 bits. We need 64-bit for DB index efficiency.">Not UUID — needs to be 64-bit &amp; sortable</div>
        <div class="d-box red" data-tip="DB auto-increment is a single point of failure and a write bottleneck. Cannot scale across multiple nodes.">Not DB auto-increment — single SPOF</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-bit-layout",
		Title:       "64-Bit Snowflake Layout",
		Description: "Bit partitioning of the 64-bit Snowflake ID into sign, timestamp, datacenter, worker, and sequence fields.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow" style="gap:2px;align-items:stretch;">
    <div class="d-box gray" style="min-width:32px;text-align:center;padding:10px 6px;" data-tip="Always 0. Keeps IDs positive when interpreted as signed 64-bit integers. Java Long is signed, so bit 63 must be 0.">
      <div style="font-size:10px;color:#888;margin-bottom:4px;">bit 63</div>
      <div style="font-weight:700;">0</div>
      <div style="font-size:10px;margin-top:4px;">sign</div>
      <div style="font-size:10px;color:#888;">1 bit</div>
    </div>
    <div class="d-box blue" style="flex:41;text-align:center;padding:10px 6px;" data-tip="Milliseconds since custom epoch (e.g. 2024-01-01 00:00:00 UTC). 41 bits = 2^41 ms = 69.7 years of headroom.">
      <div style="font-size:10px;color:#aac;margin-bottom:4px;">bits 62–22</div>
      <div style="font-weight:700;">timestamp</div>
      <div style="font-size:10px;margin-top:4px;">41 bits</div>
      <div style="font-size:10px;color:#aac;">~69 years</div>
    </div>
    <div class="d-box purple" style="flex:5;text-align:center;padding:10px 6px;" data-tip="5 bits = 32 possible datacenter IDs (0–31). Assigned statically per physical datacenter at deployment time.">
      <div style="font-size:10px;color:#cac;margin-bottom:4px;">bits 21–17</div>
      <div style="font-weight:700;">DC</div>
      <div style="font-size:10px;margin-top:4px;">5 bits</div>
      <div style="font-size:10px;color:#cac;">32 DCs</div>
    </div>
    <div class="d-box green" style="flex:5;text-align:center;padding:10px 6px;" data-tip="5 bits = 32 worker IDs per datacenter. Assigned dynamically via ZooKeeper on startup. Released on shutdown.">
      <div style="font-size:10px;color:#aca;margin-bottom:4px;">bits 16–12</div>
      <div style="font-weight:700;">worker</div>
      <div style="font-size:10px;margin-top:4px;">5 bits</div>
      <div style="font-size:10px;color:#aca;">32/DC</div>
    </div>
    <div class="d-box amber" style="flex:12;text-align:center;padding:10px 6px;" data-tip="12 bits = 4096 unique IDs per millisecond per worker. Resets to 0 each new millisecond. Overflow → wait for next ms.">
      <div style="font-size:10px;color:#ca8;margin-bottom:4px;">bits 11–0</div>
      <div style="font-weight:700;">sequence</div>
      <div style="font-size:10px;margin-top:4px;">12 bits</div>
      <div style="font-size:10px;color:#ca8;">4096/ms</div>
    </div>
  </div>
  <div class="d-flow" style="gap:8px;margin-top:8px;">
    <div class="d-box gray" style="flex:1;" data-tip="Max value: 1. Always 0 so signed int64 stays positive.">sign max: 0</div>
    <div class="d-box blue" style="flex:1;" data-tip="Max value: 2^41 - 1 = 2,199,023,255,551 ms from epoch. With epoch 2024-01-01 → expires ~2093.">ts max: 2<sup>41</sup>−1</div>
    <div class="d-box purple" style="flex:1;" data-tip="Max datacenter ID: 31 (binary 11111). 0–31 gives 32 datacenters.">DC max: 31</div>
    <div class="d-box green" style="flex:1;" data-tip="Max worker ID: 31 per datacenter. 32 workers × 32 DCs = 1024 total workers.">worker max: 31</div>
    <div class="d-box amber" style="flex:1;" data-tip="Max sequence: 4095. After 4095 in the same ms, block until next ms begins.">seq max: 4095</div>
  </div>
  <div class="d-box indigo" data-tip="ID = (timestamp_ms - epoch) &lt;&lt; 22 | datacenter_id &lt;&lt; 17 | worker_id &lt;&lt; 12 | sequence. Bit shifts pack all fields into a single int64.">
    <strong>Compose:</strong> <code>id = (ts − epoch) &lt;&lt; 22 | dc_id &lt;&lt; 17 | worker_id &lt;&lt; 12 | seq</code>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-architecture",
		Title:       "High-Level Architecture",
		Description: "Client to ID Service cluster with ZooKeeper worker ID registry and per-node components.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Clients</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Mobile apps, web frontends calling the ID service via REST or gRPC.">Mobile / Web App</div>
        <div class="d-box blue" data-tip="Backend microservices (tweet service, like service, notification service) calling ID service.">Backend Microservices</div>
        <div class="d-box blue" data-tip="Batch jobs that need large blocks of IDs pre-allocated (e.g. bulk import pipelines).">Batch Jobs</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-flow-v" style="gap:6px;">
      <div class="d-box indigo" style="text-align:center;" data-tip="Layer 7 load balancer. Routes requests to healthy ID service nodes using round-robin. Health checks every 5s.">API Gateway / Load Balancer</div>
      <div class="d-arrow-down">&#8595; round-robin</div>
      <div class="d-group">
        <div class="d-group-title">ID Service Cluster</div>
        <div class="d-flow" style="gap:8px;">
          <div class="d-flow-v" style="flex:1;">
            <div class="d-box green" data-tip="Node 0: Worker ID 3 (assigned by ZooKeeper on startup). Generates IDs fully in-memory — no external calls per ID.">
              <strong>Node 0</strong><br>
              <small>worker_id = 3</small><br>
              <small>seq: 0–4095/ms</small>
            </div>
          </div>
          <div class="d-flow-v" style="flex:1;">
            <div class="d-box green" data-tip="Node 1: Worker ID 7. Independent sequence counter. Never coordinates with Node 0 per request.">
              <strong>Node 1</strong><br>
              <small>worker_id = 7</small><br>
              <small>seq: 0–4095/ms</small>
            </div>
          </div>
          <div class="d-flow-v" style="flex:1;">
            <div class="d-box green" data-tip="Node 2: Worker ID 12. Third independent generator. Combined throughput = 3 × 4.09M = 12.3M IDs/sec.">
              <strong>Node 2</strong><br>
              <small>worker_id = 12</small><br>
              <small>seq: 0–4095/ms</small>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">ZooKeeper Cluster</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="Distributed coordination service. Stores worker ID assignments as ephemeral znodes — auto-deleted when node disconnects.">ZK Leader</div>
        <div class="d-box purple" data-tip="Quorum followers. 3-node ZK ensemble requires 2 nodes to agree. Tolerates 1 failure.">ZK Follower × 2</div>
        <div class="d-box gray" data-tip="/idservice/dc1/workers/3 → node0.example.com. Created on startup, deleted on crash.">Path: /idservice/dc{N}/workers/{id}</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title">Per-Node State (in memory)</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="System clock read via System.currentTimeMillis() or clock_gettime(). No external call — pure CPU instruction.">Local clock (ns precision)</div>
        <div class="d-box blue" data-tip="Atomic int64 incremented per request within same millisecond. Reset to 0 each new ms.">Sequence counter (atomic)</div>
        <div class="d-box blue" data-tip="Assigned once by ZooKeeper. Never changes during node lifetime. Baked into every ID this node generates.">Worker ID (immutable)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-generation-algorithm",
		Title:       "ID Generation Algorithm",
		Description: "Step-by-step ID composition: timestamp acquisition, sequence management, and bit-shift assembly.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-flow-v">
      <div class="d-box blue" data-tip="Read system clock in milliseconds. This is a single CPU instruction — nanosecond cost. No network call.">
        <strong>Step 1: Get timestamp</strong><br>
        <code>now_ms = current_time_ms() − custom_epoch</code><br>
        <small>epoch = 2024-01-01 00:00:00 UTC</small>
      </div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="If same millisecond as last call: increment sequence. If sequence overflows 4095: block until next ms. If new ms: reset sequence to 0.">
        <strong>Step 2: Manage sequence</strong><br>
        <code>if now_ms == last_ms:<br>
&nbsp;&nbsp;seq = (seq + 1) &amp; 0xFFF<br>
&nbsp;&nbsp;if seq == 0: wait_next_ms()<br>
else:<br>
&nbsp;&nbsp;seq = 0<br>
&nbsp;&nbsp;last_ms = now_ms</code>
      </div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box green" data-tip="Bit shifts pack all four fields into a single 64-bit integer. No loops, no allocations. Pure arithmetic — sub-microsecond.">
        <strong>Step 3: Compose ID</strong><br>
        <code>id = (now_ms  &lt;&lt; 22)<br>
&nbsp;&nbsp; | (dc_id    &lt;&lt; 17)<br>
&nbsp;&nbsp; | (worker_id &lt;&lt; 12)<br>
&nbsp;&nbsp; | seq</code>
      </div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box indigo" data-tip="Return the int64 to the caller. Lock is released. Next request starts from Step 1. Total time: &lt;1 microsecond.">
        <strong>Step 4: Return</strong><br>
        <code>return id  // int64</code>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Concurrency Control</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="All three steps (read clock, update sequence, compose) are wrapped in a mutex. Lock held for &lt;100ns — negligible contention.">Mutex per worker node</div>
        <div class="d-box gray" data-tip="Lock is per-node, not per-cluster. No distributed lock needed. Each worker generates IDs independently.">No distributed coordination per ID</div>
        <div class="d-box gray" data-tip="At 4096 IDs/ms, lock is acquired 4096 times per ms = every 244ns. Mutex overhead is ~20ns — acceptable.">Lock acquisition: ~20ns overhead</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title">Example Output</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Binary: 0 [41-bit timestamp] [5-bit DC] [5-bit worker] [12-bit seq]. Decimal looks like a large positive integer.">ID: 7,264,921,035,489,280,100</div>
        <div class="d-box gray" data-tip="Extract timestamp: id &gt;&gt; 22 + epoch. Extract DC: (id &gt;&gt; 17) &amp; 0x1F. Worker: (id &gt;&gt; 12) &amp; 0x1F. Seq: id &amp; 0xFFF.">Decode: ts=2024-06-15 14:32:01.123, dc=2, worker=5, seq=36</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title">wait_next_ms()</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Spin-wait until system clock advances past last_ms. Usually resolves in &lt;1ms. Rare: only when 4096 IDs generated in same ms.">
          <code>while current_time_ms() &lt;= last_ms:<br>
&nbsp;&nbsp;sleep(0)  // yield CPU<br>
return current_time_ms()</code>
        </div>
        <div class="d-box gray" data-tip="Sequence exhaustion means 4096 IDs/ms was hit. At 1M IDs/sec that's 4.096x headroom — only happens at 4.09M+/sec per node.">Triggered only at &gt;4096 IDs/ms per node</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-clock-sync",
		Title:       "Clock Skew Problem & Solution",
		Description: "How clock drift causes duplicate IDs and the strategies used to prevent it.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title" style="color:#dc2626;">Problem: Clock Skew</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Server A has correct wall clock. Generates ID with timestamp 1000ms from epoch.">Server A — clock normal (T=1000ms)</div>
        <div class="d-arrow-down">&#8595; generates</div>
        <div class="d-box blue">ID: <code>1000-dc1-w3-0001</code></div>
        <div class="d-arrow-down" style="margin-top:12px;"></div>
        <div class="d-box red" data-tip="Server B's NTP sync causes its clock to jump backward from T=1005 to T=995. This is clock skew — real hardware clocks drift ±100ms/day.">Server B — clock drifts BACK (T=995ms)</div>
        <div class="d-arrow-down">&#8595; generates</div>
        <div class="d-box red">ID: <code>995-dc1-w7-0001</code> ← timestamp in the past!</div>
        <div class="d-box amber" data-tip="The ID from Server B (timestamp 995) appears before the ID from Server A (timestamp 1000) even though it was generated after. Range queries by ID now return out-of-order results.">Sort by ID: B's ID appears BEFORE A's ID<br><small>Breaking the time-ordering guarantee</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title" style="color:#059669;">Solutions: Clock Protection</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="NTP keeps clocks within 1–50ms of true time in practice. Chrony daemon syncs every 64s. This handles slow drift but not abrupt jumps.">
          <strong>1. NTP sync (baseline)</strong><br>
          Chrony daemon — syncs every 64s<br>
          <small>Handles gradual drift, not jumps</small>
        </div>
        <div class="d-box green" data-tip="On each ID generation: compare current clock to last_ms. If current &lt; last_ms, the clock went backward. Track max drift seen.">
          <strong>2. Detect backward clock</strong><br>
          <code>if now_ms &lt; last_ms:<br>
&nbsp;&nbsp;drift = last_ms − now_ms<br>
&nbsp;&nbsp;handle_clock_skew(drift)</code>
        </div>
        <div class="d-box green" data-tip="If drift &lt; 5ms: spin-wait until clock catches up. Safe and simple. Used by Twitter Snowflake original implementation.">
          <strong>3a. Small drift (&lt;5ms): wait</strong><br>
          Spin until <code>now_ms &gt;= last_ms</code><br>
          <small>P99 wait &lt;5ms — acceptable</small>
        </div>
        <div class="d-box amber" data-tip="If drift &gt; 5ms, something is seriously wrong (NTP misconfiguration, VM migration, leap second). Fail fast — return error to caller.">
          <strong>3b. Large drift (&gt;5ms): refuse</strong><br>
          Return error: <code>CLOCK_SKEW_DETECTED</code><br>
          <small>Alert on-call. Do not generate IDs.</small>
        </div>
        <div class="d-box blue" data-tip="Alternative: use last_ms as the timestamp even when clock goes back. Sequence continues. Caller gets a slightly future-dated ID but never a past-dated one. Used by Sonyflake.">
          <strong>4. Monotonic fallback</strong><br>
          Use <code>max(now_ms, last_ms)</code> as timestamp<br>
          <small>IDs never go backward</small>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-worker-id-assignment",
		Title:       "Worker ID Assignment via ZooKeeper",
		Description: "How each ID service node registers for a unique worker ID at startup and releases it on shutdown.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Startup: Claim Worker ID</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="New ID service node starts up. Needs a unique worker ID before it can generate any IDs.">Node starts up</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Node connects to ZooKeeper ensemble (3 or 5 nodes). ZK guarantees atomic operations — only one node can claim a given worker ID.">Connect to ZooKeeper ensemble</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Create ephemeral sequential znode at /idservice/dc1/workers/. ZK auto-assigns the next available number (0, 1, 2 ...).">Create ephemeral znode:<br><code>/idservice/dc1/workers/&lt;N&gt;</code></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="N becomes the worker_id. Baked into every ID this node generates for its entire lifetime. Never changes while node is running.">Worker ID = N (0–31)<br><small>Immutable for node lifetime</small></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Node can now generate IDs. All 5 bits of worker_id field are filled with N. No further ZK calls needed per ID.">Begin generating IDs</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Shutdown &amp; Failure Recovery</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Graceful shutdown: SIGTERM received. Node stops accepting requests and disconnects from ZooKeeper cleanly.">Node graceful shutdown</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray" data-tip="ZooKeeper ephemeral znodes are automatically deleted when the client session ends. No manual cleanup required.">ZK auto-deletes ephemeral znode</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Worker ID N is now free. Next node to start can claim it. No ID collision risk because time has advanced — same worker_id + new timestamp = different IDs.">Worker ID freed, available for reuse</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title">Node Crash Recovery</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Node crashes (OOM, kernel panic, network partition). ZK session expires after session timeout (default 30s).">Node crashes unexpectedly</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="ZK session timeout = 30s. After 30s of no heartbeat, ZK deletes the ephemeral znode. Worker ID becomes available.">ZK detects after 30s timeout</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="New replacement node starts, claims the same worker_id. Safe because 30 seconds have passed — timestamp bits are different. No duplicate IDs.">Replacement node claims same ID<br><small>30s gap → no timestamp collision</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">ZooKeeper Node Tree</div>
      <div class="d-flow-v">
        <div class="d-box indigo" style="font-family:monospace;font-size:12px;line-height:1.8;" data-tip="ZK path hierarchy. Each ephemeral znode stores the hostname and start time of the claiming node.">
          /idservice/<br>
          &nbsp;&nbsp;dc1/<br>
          &nbsp;&nbsp;&nbsp;&nbsp;workers/<br>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;0 → node-a (active)<br>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;1 → node-b (active)<br>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;2 → node-c (active)<br>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;3 → (free)<br>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;...<br>
          &nbsp;&nbsp;dc2/<br>
          &nbsp;&nbsp;&nbsp;&nbsp;workers/...
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-multi-dc",
		Title:       "Multi-Datacenter Setup",
		Description: "How datacenter bits in the Snowflake ID ensure global uniqueness across multiple datacenters.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Datacenter 1 (DC ID = 0)</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="DC1 has dc_id=0 hardcoded in config. Every ID generated here has bits 21-17 = 00000.">dc_id = 0 (us-east-1)</div>
        <div class="d-flow" style="gap:6px;">
          <div class="d-box green" style="flex:1;" data-tip="Worker 0 in DC1. Generates IDs with dc=0, worker=0.">W0</div>
          <div class="d-box green" style="flex:1;" data-tip="Worker 1 in DC1. Generates IDs with dc=0, worker=1.">W1</div>
          <div class="d-box green" style="flex:1;" data-tip="... up to 32 workers per DC.">...</div>
          <div class="d-box green" style="flex:1;" data-tip="Worker 31 in DC1. Max workers per DC.">W31</div>
        </div>
        <div class="d-box gray" style="font-family:monospace;font-size:11px;" data-tip="Example: ts=1000, dc=0, worker=5, seq=42. DC bits = 00000.">
          ID: <code>ts=1000 | dc=00000 | w=00101 | seq=42</code>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Datacenter 2 (DC ID = 1)</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="DC2 has dc_id=1 hardcoded in config. Every ID generated here has bits 21-17 = 00001. Different from DC1's 00000.">dc_id = 1 (eu-west-1)</div>
        <div class="d-flow" style="gap:6px;">
          <div class="d-box purple" style="flex:1;" data-tip="Worker 0 in DC2. Generates IDs with dc=1, worker=0.">W0</div>
          <div class="d-box purple" style="flex:1;" data-tip="Worker 1 in DC2.">W1</div>
          <div class="d-box purple" style="flex:1;">...</div>
          <div class="d-box purple" style="flex:1;" data-tip="Worker 31 in DC2.">W31</div>
        </div>
        <div class="d-box gray" style="font-family:monospace;font-size:11px;" data-tip="Same timestamp, same worker number, same sequence — but dc=1 vs dc=0 makes the IDs completely different.">
          ID: <code>ts=1000 | dc=00001 | w=00101 | seq=42</code>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Why IDs Are Always Unique</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Even if DC1/W5 and DC2/W5 generate an ID at the exact same millisecond with the same sequence: dc bits differ → IDs differ. Mathematically impossible to collide.">DC bits guarantee cross-DC uniqueness<br><small>Same ts + same worker + same seq → still different IDs</small></div>
        <div class="d-box green" data-tip="Within a DC: worker bits ensure W0 and W1 never collide. Within a worker: sequence bits ensure concurrent calls in same ms don't collide.">Worker bits: within-DC uniqueness</div>
        <div class="d-box green" data-tip="Within a worker: sequence bits ensure up to 4096 unique IDs per millisecond with no coordination.">Sequence bits: within-worker uniqueness</div>
        <div class="d-box indigo" data-tip="Total unique IDs per ms across all DCs and workers: 32 DCs × 32 workers × 4096 seq = 4,294,967,296 IDs/ms = 4.29B/ms = 4.29 trillion IDs/sec.">
          Max throughput: <strong>32 × 32 × 4096 = 4.29B IDs/ms</strong>
        </div>
      </div>
    </div>
  </div>
</div>
<div class="d-box amber" style="margin-top:8px;" data-tip="DC ID is a static config value deployed per datacenter. If you misconfigure two DCs with the same dc_id, you get collisions. This is an operational risk — solved by config management and boot-time validation.">
  <strong>Operational rule:</strong> dc_id is statically configured per datacenter. Must be validated at boot — duplicate dc_ids cause collisions.
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-failure-modes",
		Title:       "Failure Modes & Solutions",
		Description: "How the system handles clock skew, worker crashes, and sequence exhaustion.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title" style="color:#d97706;">Clock Goes Backward</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="NTP adjustment, VM live migration, or leap second causes system clock to go backward.">Trigger: NTP step adjustment or VM migration</div>
        <div class="d-box amber" data-tip="Detect: now_ms &lt; last_ms. Compute drift = last_ms - now_ms.">Detected: now_ms &lt; last_ms</div>
        <div class="d-box green" data-tip="If drift is small (&lt;5ms): spin-wait. At 1M IDs/sec, 5ms pause = 5000 queued requests. Acceptable in practice.">Small drift (&lt;5ms): spin-wait<br><small>Resume once clock catches up</small></div>
        <div class="d-box red" data-tip="Large drift (&gt;5ms) means something is seriously wrong. Return error code. Alert on-call. Do not generate potentially-duplicate IDs.">Large drift (&gt;5ms): return error<br><small>Alert: CLOCK_SKEW_CRITICAL</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title" style="color:#dc2626;">Worker Crash</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Worker process crashes (OOM, segfault, hardware failure). ZooKeeper session kept alive by heartbeat.">Worker process dies</div>
        <div class="d-box amber" data-tip="ZK session expires after 30s of missed heartbeats. Ephemeral znode deleted automatically.">ZK session timeout: 30s</div>
        <div class="d-box amber" data-tip="During 30s window, no IDs generated from this worker. Load balancer stops routing to dead node via health checks (&lt;10s to detect).">~10–30s of reduced capacity</div>
        <div class="d-box green" data-tip="New process starts, claims same or different worker_id from ZK. Safe to reuse because 30s elapsed = timestamp advanced = no collision.">New node claims worker ID<br><small>Sequence gap is normal — not an error</small></div>
        <div class="d-box gray" data-tip="Other workers absorb the load. With 3 workers: 2 workers handle full traffic. At 4.09M/worker, 2 workers = 8.18M/sec capacity.">Other workers absorb traffic</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title" style="color:#7c3aed;">Sequence Exhausted (4096/ms)</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="Sequence counter reaches 4095 within the same millisecond. Next ID request would overflow the 12-bit field.">Sequence hits 4095 in same ms</div>
        <div class="d-box purple" data-tip="wait_next_ms() spins in a tight loop reading the system clock. Usually resolves within &lt;1ms.">Spin-wait for clock to advance</div>
        <div class="d-box purple" data-tip="Once now_ms &gt; last_ms, sequence resets to 0. ID generation resumes. The wait adds &lt;1ms latency to that request.">New ms starts → seq reset to 0</div>
        <div class="d-box green" data-tip="Prevention: add more worker nodes. 4 workers = 16.4M IDs/sec capacity. Sequence exhaustion only occurs above 4.09M IDs/sec per node.">Prevention: scale out workers<br><small>Each added worker = +4.09M IDs/sec</small></div>
        <div class="d-box gray" data-tip="Alternative: pre-generate batches of IDs and serve from a queue. Eliminates per-ID wait at cost of slightly non-monotonic ordering within batch.">Alt: pre-allocate ID batches</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-throughput",
		Title:       "Throughput Analysis",
		Description: "Capacity math for single node, cluster, and full multi-DC deployment, plus comparison with alternatives.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Throughput Math</div>
      <div class="d-entity">
        <div class="d-entity-header blue">Snowflake Throughput</div>
        <div class="d-entity-body">
          <div class="d-row"><span>Single worker</span><span><strong>4,096/ms = 4.09M/sec</strong></span></div>
          <div class="d-row"><span>32 workers/DC</span><span><strong>131M/sec per DC</strong></span></div>
          <div class="d-row"><span>32 DCs</span><span><strong>4.19B/sec system-wide</strong></span></div>
          <div class="d-row"><span>Timestamp lifetime</span><span><strong>69.7 years</strong></span></div>
          <div class="d-row"><span>IDs before epoch expires</span><span><strong>2<sup>63</sup> − 1</strong></span></div>
        </div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title">Latency per ID</div>
      <div class="d-entity">
        <div class="d-entity-header green">Generation Cost</div>
        <div class="d-entity-body">
          <div class="d-row"><span>Clock read</span><span>&lt;50ns</span></div>
          <div class="d-row"><span>Mutex acquire</span><span>~20ns</span></div>
          <div class="d-row"><span>Bit ops (3×)</span><span>&lt;5ns</span></div>
          <div class="d-row"><span>Total per ID</span><span><strong>&lt;100ns</strong></span></div>
          <div class="d-row"><span>Network (gRPC)</span><span>+0.1–1ms</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">ID Scheme Comparison</div>
      <div class="d-entity">
        <div class="d-entity-header indigo">Alternatives</div>
        <div class="d-entity-body">
          <div class="d-row"><span><strong>UUID v4</strong></span><span>128-bit, random</span></div>
          <div class="d-row"><span>Sortable?</span><span style="color:#dc2626;">No (random)</span></div>
          <div class="d-row"><span>DB index size</span><span style="color:#dc2626;">2× larger</span></div>
          <div class="d-row"><span>Coordination?</span><span style="color:#059669;">None needed</span></div>
          <div class="d-row" style="margin-top:8px;"><span><strong>DB auto-increment</strong></span><span>int64, perfect sort</span></div>
          <div class="d-row"><span>Sortable?</span><span style="color:#059669;">Yes</span></div>
          <div class="d-row"><span>Distributed?</span><span style="color:#dc2626;">No — single writer</span></div>
          <div class="d-row"><span>Throughput</span><span style="color:#dc2626;">~50K/sec max</span></div>
          <div class="d-row" style="margin-top:8px;"><span><strong>Snowflake</strong></span><span>64-bit, time-sorted</span></div>
          <div class="d-row"><span>Sortable?</span><span style="color:#059669;">Yes (by time)</span></div>
          <div class="d-row"><span>Distributed?</span><span style="color:#059669;">Yes — per worker</span></div>
          <div class="d-row"><span>Throughput</span><span style="color:#059669;">4.19B/sec</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Timestamp: 41-Bit Lifetime</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="2^41 milliseconds = 2,199,023,255,551 ms ÷ 1000 = 2,199,023,255 seconds ÷ 86400 = 25,451 days ÷ 365 = 69.7 years.">2<sup>41</sup> ms = 69.7 years</div>
        <div class="d-box green" data-tip="Twitter's Snowflake used epoch 2010-11-04. Expires ~2080. You can choose your own epoch — pick something recent to maximize headroom.">Custom epoch: 2024-01-01 → expires 2093</div>
        <div class="d-box amber" data-tip="When 41-bit counter overflows, IDs wrap around to 0 again — timestamps restart. Plan for migration before this happens (60+ years away).">Plan migration before epoch overflow<br><small>~69 years from chosen epoch</small></div>
        <div class="d-box gray" data-tip="Extend timestamp to 42 bits: doubles lifetime to 139 years, reduces sequence to 11 bits (2048/ms). Trade-off based on your expected lifetime.">Want more time? Use 42 bits → 139 years<br><small>Trade: 11-bit seq (2048/ms per worker)</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-api-design",
		Title:       "API Design",
		Description: "REST and gRPC API endpoints for the ID generator service.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">REST API</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Returns a single 64-bit ID. Sub-millisecond generation. Simple to call from any language.">GET /v1/id</div>
        <div class="d-box gray" style="font-family:monospace;font-size:12px;">Response 200:<br>{<br>&nbsp;&nbsp;"id": 7264921035489280100,<br>&nbsp;&nbsp;"id_str": "7264921035489280100"<br>}</div>
        <div class="d-box amber" data-tip="id_str is included because JavaScript's Number type loses precision for int64 > 2^53. Always use id_str in JS clients.">Note: JS clients must use id_str<br><small>JS Number loses precision &gt;2<sup>53</sup></small></div>
        <div class="d-box blue" style="margin-top:8px;" data-tip="Returns a batch of N IDs in one request. Reduces network round-trips for bulk operations. Max 10,000 per call.">POST /v1/ids</div>
        <div class="d-box gray" style="font-family:monospace;font-size:12px;">Body: {"count": 1000}<br><br>Response 200:<br>{<br>&nbsp;&nbsp;"ids": [7264...100, 7264...101, ...],<br>&nbsp;&nbsp;"count": 1000<br>}</div>
        <div class="d-box gray" data-tip="Maximum 10,000 IDs per batch call. Larger batches: call multiple times. IDs within a batch are sequential by sequence counter.">Max batch size: 10,000<br><small>Larger requests: split into multiple calls</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">gRPC API (Low Latency)</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="gRPC adds ~0.05ms vs REST's ~0.5ms for serialization. For 1M calls/sec, that's 450 fewer CPU-seconds/sec of serialization overhead.">Preferred for internal services</div>
        <div class="d-box gray" style="font-family:monospace;font-size:11px;line-height:1.6;">service IDService {<br>&nbsp;&nbsp;rpc GetID(IDRequest) returns (IDResponse);<br>&nbsp;&nbsp;rpc GetIDs(BatchRequest) returns (BatchResponse);<br>}<br><br>message IDResponse {<br>&nbsp;&nbsp;int64 id = 1;<br>&nbsp;&nbsp;string id_str = 2;<br>}</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px;">
      <div class="d-group-title">Error Responses</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Clock drifted backward by more than 5ms. Caller should retry with exponential backoff after 10ms.">503 CLOCK_SKEW_DETECTED<br><small>Retry after 10ms</small></div>
        <div class="d-box amber" data-tip="Sequence exhausted in current ms. Caller should retry immediately — wait is &lt;1ms.">429 SEQUENCE_EXHAUSTED<br><small>Retry immediately (&lt;1ms wait)</small></div>
        <div class="d-box amber" data-tip="count &gt; 10,000 in batch request.">400 BATCH_SIZE_EXCEEDED<br><small>Max 10,000 per request</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-monitoring",
		Title:       "Monitoring & Alerting",
		Description: "Key metrics, alert thresholds, and observability for the ID generator service.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Key Metrics (per worker)</div>
      <div class="d-entity">
        <div class="d-entity-header blue">Operational Metrics</div>
        <div class="d-entity-body">
          <div class="d-row"><span>ids_generated_total</span><span>counter — rate/sec</span></div>
          <div class="d-row"><span>clock_drift_ms</span><span>gauge — milliseconds</span></div>
          <div class="d-row"><span>sequence_utilization_pct</span><span>gauge — 0–100%</span></div>
          <div class="d-row"><span>wait_next_ms_total</span><span>counter — seq exhaustions</span></div>
          <div class="d-row"><span>zk_connection_state</span><span>gauge — 0/1</span></div>
          <div class="d-row"><span>id_generation_latency_ns</span><span>histogram — p50/p99</span></div>
          <div class="d-row"><span>worker_id</span><span>gauge — current assignment</span></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Alert Conditions</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Clock drift above 5ms means IDs may be out of order or generation is halted. Page immediately — this is a P0 incident.">clock_drift_ms &gt; 5ms → PAGE<br><small>Risk: ID generation halted or out-of-order</small></div>
        <div class="d-box red" data-tip="ZK disconnect means no worker ID validation. Node continues generating IDs but cannot renew registration. Risk of ID collision on restart.">zk_connection_state == 0 for &gt;30s → PAGE<br><small>Risk: worker ID conflict on restart</small></div>
        <div class="d-box amber" data-tip="Sequence utilization above 80% means approaching 4096 IDs/ms limit. Add worker nodes before saturation.">sequence_utilization &gt; 80% → WARN<br><small>Action: add worker nodes</small></div>
        <div class="d-box amber" data-tip="wait_next_ms means sequence was exhausted. Should be near-zero under normal load. Frequent occurrences = underprovisioned.">wait_next_ms_rate &gt; 10/sec → WARN<br><small>Action: scale out workers</small></div>
        <div class="d-box blue" data-tip="p99 latency above 1ms means the gRPC/HTTP overhead is too high, or the node is under CPU pressure.">id_generation_latency p99 &gt; 1ms → WARN<br><small>Baseline: p99 &lt; 0.1ms</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Dashboards</div>
      <div class="d-flow-v">
        <div class="d-box indigo" data-tip="Grafana dashboard showing IDs/sec per worker, total cluster throughput, and capacity headroom.">IDs/sec — per worker + cluster total</div>
        <div class="d-box indigo" data-tip="Heatmap showing clock drift over time. Should be flat near 0ms. Spikes indicate NTP events.">Clock drift heatmap</div>
        <div class="d-box indigo" data-tip="Sequence utilization per worker over time. Early warning for saturation.">Sequence utilization trend</div>
        <div class="d-box gray" data-tip="ZooKeeper health: session state, latency to ZK ensemble, number of active worker registrations.">ZooKeeper health panel</div>
        <div class="d-box gray" data-tip="API latency: p50, p95, p99 for GetID and GetIDs endpoints. Broken down by endpoint and response code.">API latency p50/p99 by endpoint</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "idgen-batch-generation",
		Title:       "Batch ID Pre-allocation",
		Description: "Optimization: pre-allocate blocks of IDs in memory to serve without per-ID lock overhead.",
		ContentFile: "problems/id-generator",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Standard: Per-ID Generation</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Each GetID call acquires mutex, reads clock, increments sequence, composes ID, releases mutex.">Request 1: acquire lock → generate → release</div>
        <div class="d-box blue" data-tip="Request 2 must wait for Request 1 to release the lock before it can proceed.">Request 2: wait → acquire → generate → release</div>
        <div class="d-box amber" data-tip="At 4M req/sec, mutex is acquired 4M times/sec. Each acquire = ~20ns → 80ms/sec spent on locking across 4M requests.">4M req/sec = 4M lock acquisitions/sec<br><small>~80ms CPU/sec in locking overhead</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Optimized: Block Pre-allocation</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Allocate a range of 1000 IDs at once under a single lock acquisition. Store [current_id, max_id] range atomically.">Allocate block [1001–2000] under 1 lock</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Serve ID 1001, 1002, ..., 2000 using a simple atomic counter. No lock needed per ID — just atomic increment.">Serve 1001, 1002, ... from atomic counter</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="When counter reaches 2000, allocate next block [2001–3000] under 1 lock. 1 lock per 1000 IDs.">Block exhausted → allocate [2001–3000]</div>
        <div class="d-box indigo" data-tip="1000 IDs per block → 1 lock per 1000 requests. 4M req/sec = 4000 lock acquisitions/sec (vs 4M). 1000x reduction in locking.">1000× fewer lock acquisitions<br><small>4M req/sec = only 4K blocks/sec allocated</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Trade-offs</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Serving IDs from pre-allocated block is a single atomic increment — ~5ns. No mutex contention.">Pro: Near-zero per-ID overhead</div>
        <div class="d-box green" data-tip="Blocks can be distributed across threads without coordination — each thread gets its own block.">Pro: Thread-local ID serving</div>
        <div class="d-box amber" data-tip="If node crashes mid-block (served IDs 1001–1500, then died), IDs 1501–2000 are permanently lost. Gaps in ID sequence — not a bug, normal behavior.">Con: ID gaps on crash<br><small>Lost block = gap in sequence — acceptable</small></div>
        <div class="d-box amber" data-tip="IDs within a block are sequential but block allocation order determines global sort order. IDs from different threads may interleave within the same timestamp.">Con: Weak ordering between threads<br><small>Within-ms ordering not strict</small></div>
        <div class="d-box gray" data-tip="Block size of 1000 is a good default. Larger blocks = fewer allocations but larger gaps on crash. Smaller blocks = more allocations but tighter ordering.">Block size: 1000 (tunable 100–10000)</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

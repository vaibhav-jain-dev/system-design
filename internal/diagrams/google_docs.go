package diagrams

func registerGoogleDocs(r *Registry) {
	r.Register(&Diagram{
		Slug:        "collab-requirements",
		Title:       "Requirements & Scale",
		Description: "Functional requirements, scale targets, and non-functional constraints for collaborative document editing",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green"><span class="d-step">1</span>Real-time collaborative editing — multiple users edit same document simultaneously <div class="d-tag green">&#10003; core</div></div>
        <div class="d-box green"><span class="d-step">2</span>Conflict-free merging — concurrent edits from different cursors converge to identical state <div class="d-tag green">&#10003; core</div></div>
        <div class="d-box green"><span class="d-step">3</span>Persistent document storage — documents survive server restarts and failures <div class="d-tag green">&#10003; core</div></div>
        <div class="d-box green"><span class="d-step">4</span>Presence &amp; cursors — see other editors&#39; cursor positions in real time</div>
        <div class="d-box blue"><span class="d-step">5</span>Offline editing — queue operations locally, sync on reconnect <div class="d-tag blue">P1</div></div>
        <div class="d-box blue"><span class="d-step">6</span>30-day version history — restore any past revision <div class="d-tag blue">P1</div></div>
        <div class="d-box gray"><span class="d-step">7</span>Comments &amp; suggestions — inline annotations <div class="d-tag gray">P2</div></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">1 billion total documents <span class="d-metric size">1B docs</span></div>
        <div class="d-box purple">10 million concurrent editors <span class="d-metric throughput">10M concurrent</span></div>
        <div class="d-box purple">10K operations/sec per popular document <span class="d-metric throughput">10K ops/doc</span></div>
        <div class="d-box purple">500M documents with at least 1 active user <span class="d-metric size">500M active</span></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box amber">Operation propagation: &lt;100ms P99 <span class="d-metric latency">&lt;100ms</span></div>
        <div class="d-box amber">Availability: 99.99% (52 min/year downtime) <span class="d-metric throughput">4 nines</span></div>
        <div class="d-box amber">Convergence: all clients reach identical state after conflicts</div>
        <div class="d-box red">No data loss — every acknowledged operation must persist <div class="d-tag red">critical</div></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end architecture: clients, WebSocket servers, operation log, document store, and version history",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Browser Client<br><small>CRDT/OT engine local</small></div>
    <div class="d-arrow">&#8596; WS</div>
    <div class="d-box blue">Mobile Client<br><small>offline queue (IndexedDB)</small></div>
  </div>
  <div class="d-arrow-down">&#8595; WebSocket (sticky by doc_id)</div>
  <div class="d-group">
    <div class="d-group-title">Load Balancer — consistent hash by doc_id</div>
    <div class="d-flow">
      <div class="d-box green">Document Server A<br><small>docs 0-33%</small></div>
      <div class="d-box green">Document Server B<br><small>docs 33-66%</small></div>
      <div class="d-box green">Document Server C<br><small>docs 66-100%</small></div>
    </div>
  </div>
  <div class="d-flow">
    <div class="d-arrow-down">&#8595; op log (write-through)</div>
    <div class="d-arrow-down">&#8595; snapshot (every 1K ops)</div>
    <div class="d-arrow-down">&#8595; presence (TTL 5s)</div>
  </div>
  <div class="d-flow">
    <div class="d-box purple">Redis<br><small>op log ring buffer<br>presence HSET</small></div>
    <div class="d-box indigo">Firestore / DynamoDB<br><small>operations table<br>document metadata</small></div>
    <div class="d-box amber">S3<br><small>snapshots every 1K ops<br>30-day version history</small></div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-ot-algorithm",
		Title:       "Operational Transformation (OT)",
		Description: "Classic OT conflict: two concurrent edits, transform function adjusts positions before apply",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Concurrent Edit Scenario</div>
    <div class="d-flow">
      <div class="d-box blue">Initial state<br><code>&quot;Hello World&quot;</code><br><small>positions 0&#8211;10</small></div>
      <div class="d-flow-v">
        <div class="d-box green">User A<br>insert(&#39;X&#39;, pos=5)<br><small>&#8594; &quot;HelloX World&quot;</small></div>
        <div class="d-arrow">concurrent</div>
        <div class="d-box amber">User B<br>delete(pos=3)<br><small>&#8594; &quot;HelWorldo&quot; (?)</small></div>
      </div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">&#10007; Without OT — Wrong Result</div>
        <div class="d-flow-v">
          <div class="d-box red">Apply A: insert(&#39;X&#39;, 5) &#8594; &quot;HelloX World&quot;</div>
          <div class="d-box red">Apply B naively: delete(3) &#8594; deletes &#39;l&#39; in wrong position</div>
          <div class="d-box red">Result: &quot;HelXo World&quot; &#8800; what B intended <div class="d-tag red">diverged</div></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">&#10003; With OT — Correct Result</div>
        <div class="d-flow-v">
          <div class="d-box green">Apply A: insert(&#39;X&#39;, 5) &#8594; &quot;HelloX World&quot;</div>
          <div class="d-box green">Transform B vs A: pos 3 &lt; 5 &#8594; no shift needed<br><small>transform(delete(3), insert(5)) &#8594; delete(3)</small></div>
          <div class="d-box green">Apply transformed B: delete(3) &#8594; &quot;HelXo World&quot; <div class="d-tag green">converged</div></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Transformation Function</div>
    <div class="d-flow">
      <div class="d-box indigo">transform(op_b, op_a) &#8594; op_b&#39;<br><small>if op_a inserts before op_b.pos: op_b.pos += op_a.len<br>if op_a deletes before op_b.pos: op_b.pos -= op_a.len<br>if positions equal: tie-break by user_id</small></div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-crdt-comparison",
		Title:       "OT vs CRDT Comparison",
		Description: "Trade-offs between Operational Transformation (server-centric) and CRDTs (serverless convergence)",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">OT — Google Docs Approach</div>
      <div class="d-flow-v">
        <div class="d-box green">&#10003; Mature — 20+ years battle-tested in production</div>
        <div class="d-box green">&#10003; Server serializes all ops — single total order</div>
        <div class="d-box green">&#10003; Lower memory overhead — no tombstones needed</div>
        <div class="d-box red">&#10007; Requires central server for ordering <div class="d-tag red">single point</div></div>
        <div class="d-box red">&#10007; O(N&#178;) transformation complexity — N concurrent ops</div>
        <div class="d-box red">&#10007; Offline editing requires full rebasing on reconnect</div>
        <div class="d-label">Best for: text documents with server always available</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">CRDT — Figma / Notion Approach</div>
      <div class="d-flow-v">
        <div class="d-box green">&#10003; Serverless convergence — ops commutative &amp; associative</div>
        <div class="d-box green">&#10003; Works fully offline — merge on reconnect automatically</div>
        <div class="d-box green">&#10003; No coordination needed — P2P sync possible</div>
        <div class="d-box amber">~ Higher memory — tombstones for deleted items <div class="d-tag amber">trade-off</div></div>
        <div class="d-box amber">~ Complex implementation — Yjs, Automerge libraries needed</div>
        <div class="d-box amber">~ Eventual consistency — brief divergence visible to users</div>
        <div class="d-label">Best for: offline-first, mobile, or multi-master sync</div>
      </div>
    </div>
  </div>
</div>
<div class="d-group" style="margin-top:8px">
  <div class="d-group-title">Decision</div>
  <div class="d-flow">
    <div class="d-box indigo">Use <strong>OT</strong> for server-mediated text docs (Google Docs model) &#8212; simpler client, server owns ordering. Use <strong>CRDT</strong> for offline-first or P2P (Figma, Notion) &#8212; higher complexity, better offline UX.</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-data-model",
		Title:       "Data Model",
		Description: "Core tables: documents, operations, snapshots, and sessions with primary/foreign keys",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header green">documents</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> doc_id (UUID)</div>
        <div>title (string)</div>
        <div>owner_id (user_id FK)</div>
        <div>created_at (timestamp)</div>
        <div>current_version (int)</div>
        <div>content_snapshot_key (S3 key)</div>
      </div>
    </div>
    <div class="d-entity" style="margin-top:8px">
      <div class="d-entity-header amber">sessions</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> session_id (UUID)</div>
        <div><span class="fk">FK</span> doc_id</div>
        <div>user_id</div>
        <div>websocket_id (server-local)</div>
        <div>client_revision (int)</div>
        <div>joined_at (timestamp)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">operations</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> op_id (UUID)</div>
        <div><span class="fk">FK</span> doc_id</div>
        <div>user_id</div>
        <div>op_type (insert | delete | format)</div>
        <div>position (int)</div>
        <div>content (text, nullable)</div>
        <div>revision_num (int) <span class="idx idx-btree">idx</span></div>
        <div>timestamp (int64, ms)</div>
      </div>
    </div>
    <div class="d-entity" style="margin-top:8px">
      <div class="d-entity-header purple">snapshots</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> doc_id</div>
        <div><span class="pk">PK</span> revision_num (composite)</div>
        <div>content_s3_key (string)</div>
        <div>op_count (int, always 1000)</div>
        <div>created_at (timestamp)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-operation-flow",
		Title:       "Operation Processing Flow",
		Description: "Step-by-step sequence from user keystroke to broadcast: OT transform, op log append, and client update",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue"><span class="d-step">1</span>User types &#39;X&#39;<br><small>client creates op:<br>{type:insert, pos:5,<br>char:&#39;X&#39;, client_rev:42}</small></div>
    <div class="d-arrow">&#8594; WS</div>
    <div class="d-box green"><span class="d-step">2</span>Document Server receives<br><small>acquire doc-level lock<br>(Redis SET NX, TTL 100ms)</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber"><span class="d-step">3</span>OT Transform<br><small>compare client_rev=42 vs<br>server_rev=44 &#8594; transform<br>op against 2 missed ops</small></div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box purple"><span class="d-step">4</span>Append to op log<br><small>Redis RPUSH ops:{doc_id}<br>Firestore write (rev=45)<br>release lock</small></div>
    <div class="d-arrow">&#8594; broadcast</div>
    <div class="d-box green"><span class="d-step">5</span>Fan-out to all session clients<br><small>WebSocket send to all<br>connected users on this doc<br>except sender</small></div>
    <div class="d-arrow">&#8594; ack</div>
    <div class="d-box blue"><span class="d-step">6</span>Sender receives ack<br><small>update local client_rev=45<br>apply transformed op<br>to local document</small></div>
  </div>
  <div class="d-group" style="margin-top:8px">
    <div class="d-group-title">Latency Budget</div>
    <div class="d-flow">
      <div class="d-box gray">WS receive: ~1ms</div>
      <div class="d-arrow">+</div>
      <div class="d-box gray">Lock acquire: ~2ms</div>
      <div class="d-arrow">+</div>
      <div class="d-box gray">OT transform: ~0.5ms</div>
      <div class="d-arrow">+</div>
      <div class="d-box gray">Redis write: ~1ms</div>
      <div class="d-arrow">+</div>
      <div class="d-box gray">Broadcast: ~2ms</div>
      <div class="d-arrow">=</div>
      <div class="d-box green">~6ms total <span class="d-metric latency">&lt;10ms P99</span></div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-presence",
		Title:       "Presence & Cursors",
		Description: "Real-time cursor position sharing: Redis HSET per document, TTL-based staleness, 300ms update cadence",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Cursor Update Flow</div>
      <div class="d-flow-v">
        <div class="d-box blue">User moves cursor to position 142<br><small>client throttles: send every 300ms</small></div>
        <div class="d-arrow-down">&#8595; WebSocket</div>
        <div class="d-box green">Document Server receives cursor_update<br><small>{user_id, position:142, color:&#34;#6366F1&#34;}</small></div>
        <div class="d-arrow-down">&#8595; Redis HSET</div>
        <div class="d-box purple">Redis: HSET presence:{doc_id} {user_id} {pos, color, ts}<br><small>EXPIRE presence:{doc_id} 10s (refresh on each update)</small></div>
        <div class="d-arrow-down">&#8595; broadcast</div>
        <div class="d-box amber">Fan-out to all session clients<br><small>all other editors see ghost cursor move</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Stale Cursor Cleanup</div>
      <div class="d-flow-v">
        <div class="d-box amber">User goes idle or disconnects</div>
        <div class="d-box red">Redis field expires after 5s no update<br><small>HDEL presence:{doc_id} {user_id}</small></div>
        <div class="d-box blue">Server broadcasts leave event<br><small>{type:&#34;cursor_leave&#34;, user_id}</small></div>
        <div class="d-box green">Clients remove ghost cursor from UI</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Storage Estimate</div>
      <div class="d-flow-v">
        <div class="d-box gray">10M concurrent editors<br><small>~200 bytes per presence entry<br>= 2 GB total across Redis cluster</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-offline-support",
		Title:       "Offline Editing Support",
		Description: "Client queues ops in IndexedDB while offline; server transforms all queued ops against concurrent ops on reconnect",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">While Offline</div>
        <div class="d-flow-v">
          <div class="d-box amber">Network disconnected<br><small>client detects via WebSocket close event</small></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">User continues editing<br><small>operations applied to local doc immediately</small></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Each op queued to IndexedDB<br><small>[{type:insert,pos:5,char:&#39;A&#39;,base_rev:42},<br> {type:delete,pos:3,base_rev:43}, ...]</small></div>
          <div class="d-label">5 minutes offline = ~10 queued ops (typical)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">On Reconnect</div>
        <div class="d-flow-v">
          <div class="d-box green">WebSocket reconnected<br><small>client sends: {reconnect: true, base_rev: 42, queued_ops: [...]}</small></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Server fetches ops [rev 42 &#8594; current]<br><small>e.g. 15 ops happened while client was offline</small></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Transform each queued op against all missed ops<br><small>10 queued &#215; 15 missed = 150 OT transforms<br>&#8594; O(N&#215;M), completes in &lt;5ms</small></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Send transformed ops + full current state<br><small>client applies, document converges</small></div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-version-history",
		Title:       "Version History & Restore",
		Description: "Checkpoint-based snapshots every 1000 ops in S3; restore by replaying ops from nearest snapshot",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Snapshot Strategy</div>
    <div class="d-flow">
      <div class="d-box blue">ops 1&#8211;999<br><small>stored in op log</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Checkpoint @ rev 1000<br><small>full snapshot &#8594; S3<br>~50KB compressed doc</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">ops 1001&#8211;1999<br><small>stored in op log</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Checkpoint @ rev 2000<br><small>full snapshot &#8594; S3</small></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Restore to Version V Algorithm</div>
    <div class="d-flow">
      <div class="d-box indigo"><span class="d-step">1</span>Find nearest snapshot S &#8804; V<br><small>query snapshots table:<br>WHERE doc_id=X AND rev &#8804; V<br>ORDER BY rev DESC LIMIT 1</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber"><span class="d-step">2</span>Load snapshot from S3<br><small>fetch content_s3_key<br>deserialize document state</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple"><span class="d-step">3</span>Replay ops [S &#8594; V]<br><small>fetch from operations table<br>apply each op in order</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green"><span class="d-step">4</span>Restored document at rev V<br><small>show diff vs current<br>user confirms restore</small></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Example Timeline</div>
    <div class="d-flow">
      <div class="d-box gray">Mon 9am<br>Checkpoint<br>rev=1000</div>
      <div class="d-arrow">&#8594; +20 ops</div>
      <div class="d-box amber">Mon 11am<br>User wants<br>this version</div>
      <div class="d-arrow">&#8594; +30 more ops</div>
      <div class="d-box gray">Mon 2pm<br>Current<br>rev=1050</div>
    </div>
    <div class="d-label">Restore: load snapshot(rev=1000) + replay 20 ops = O(20) work</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-scaling",
		Title:       "Scaling Document Servers",
		Description: "Consistent hash ring routes doc_id to a fixed server; migration on scale-out fetches state from Redis",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Consistent Hash Routing (doc_id &#8594; server)</div>
    <div class="d-flow">
      <div class="d-box blue">Client requests doc_id=&#34;abc123&#34;</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Load Balancer<br><small>hash(&#34;abc123&#34;) % ring<br>always &#8594; Server B</small></div>
      <div class="d-arrow">&#8594; sticky</div>
      <div class="d-box green">Server B<br><small>holds all WebSocket sessions<br>for doc &#34;abc123&#34; in memory</small></div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Steady State — Server Capacity</div>
        <div class="d-flow-v">
          <div class="d-box purple">100K active documents per server<br><small>each doc: ~50 WebSocket connections avg</small></div>
          <div class="d-box purple">10K concurrent editors per server<br><small>50 connections &#215; 200 servers = 10M total</small></div>
          <div class="d-box purple">~4 GB RAM per server<br><small>40KB per active doc &#215; 100K docs</small></div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Scale-Out: Adding a Server</div>
        <div class="d-flow-v">
          <div class="d-box amber"><span class="d-step">1</span>Add Server D to hash ring<br><small>~25% of docs rehash to Server D</small></div>
          <div class="d-box blue"><span class="d-step">2</span>Affected clients get 302 redirect<br><small>reconnect WebSocket to Server D</small></div>
          <div class="d-box green"><span class="d-step">3</span>Server D hydrates from Redis + Firestore<br><small>fetch latest ops + snapshot for migrated docs</small></div>
          <div class="d-box green"><span class="d-step">4</span>Ready to serve in &lt;500ms per doc</small></div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-access-control",
		Title:       "Access Control & Sharing",
		Description: "Permission levels, share-link mechanics, and per-operation ACL checks in WebSocket handler",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Permission Levels</div>
      <div class="d-flow-v">
        <div class="d-box red">Owner — edit + share + delete + transfer <div class="d-tag red">full control</div></div>
        <div class="d-box green">Editor — read + write operations <div class="d-tag green">write</div></div>
        <div class="d-box blue">Commenter — read + add comments only <div class="d-tag blue">comment</div></div>
        <div class="d-box gray">Viewer — read only, no cursor broadcast <div class="d-tag gray">read</div></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Share Link Mechanics</div>
      <div class="d-flow-v">
        <div class="d-box amber">Owner generates link: /doc/abc?token=xyz<br><small>token encodes: doc_id + permission_level + expiry</small></div>
        <div class="d-box blue">Recipient clicks link &#8594; JWT issued with claims<br><small>{doc_id, role:&#34;editor&#34;, exp:+7d}</small></div>
        <div class="d-box green">ACL entry written to documents table<br><small>shared_with: {user_id, role, granted_by}</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Per-Operation ACL Check</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="d-step">1</span>Client connects via WebSocket<br><small>sends JWT in Upgrade request header:<br>Authorization: Bearer &lt;token&gt;</small></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber"><span class="d-step">2</span>Server verifies JWT signature<br><small>extract user_id + doc_id from claims</small></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple"><span class="d-step">3</span>Fetch ACL from Redis cache<br><small>HGET acl:{doc_id} {user_id} &#8594; role<br>TTL 60s; fallback to Firestore</small></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><span class="d-step">4</span>role=editor &#8594; allow write operations<br><small>role=viewer &#8594; allow subscribe only<br>no role &#8594; close WebSocket (403)</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "collab-monitoring",
		Title:       "Monitoring & SLAs",
		Description: "Key metrics dashboard: propagation latency, WebSocket health, divergence detection, and alert thresholds",
		ContentFile: "problems/google-docs",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Key Metrics</div>
      <div class="d-flow-v">
        <div class="d-box purple">Operation propagation latency P50 / P99<br><small>target: P50 &lt;20ms, P99 &lt;100ms</small> <span class="d-metric latency">&lt;100ms P99</span></div>
        <div class="d-box blue">WebSocket active connections<br><small>alert if drops &gt;10% in 60s (server crash indicator)</small></div>
        <div class="d-box amber">Document server memory usage<br><small>alert at 80% — triggers scale-out</small></div>
        <div class="d-box green">Op log size per document<br><small>alert if &gt;5000 ops without snapshot (snapshot lag)</small></div>
        <div class="d-box red">OT divergence detection rate<br><small>target: 0 errors/day — any divergence = critical page <div class="d-tag red">P0</div></small></div>
        <div class="d-box gray">Snapshot creation lag<br><small>alert if &gt;2000 ops since last snapshot</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Alert Thresholds</div>
      <div class="d-flow-v">
        <div class="d-box red">Propagation latency P99 &gt;500ms &#8594; page on-call <div class="d-tag red">immediate</div></div>
        <div class="d-box red">OT error / divergence detected &#8594; critical — document corruption risk <div class="d-tag red">P0 incident</div></div>
        <div class="d-box amber">WebSocket reconnect rate &gt;5%/min &#8594; investigate server stability</div>
        <div class="d-box amber">Op log Redis memory &gt;70% &#8594; add shard</div>
        <div class="d-box blue">Snapshot creation failure &#8594; alert + retry with backoff</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">SLA Summary</div>
      <div class="d-flow-v">
        <div class="d-box green">Availability: 99.99% &#8212; 52 min/year budget</div>
        <div class="d-box green">Propagation: &lt;100ms P99 end-to-end</div>
        <div class="d-box green">Data durability: zero data loss after ack</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

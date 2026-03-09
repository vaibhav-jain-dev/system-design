package diagrams

func registerFileStorage(r *Registry) {
	r.Register(&Diagram{
		Slug:        "fs-requirements",
		Title:       "Requirements & Scale",
		Description: "Scale estimates and NFRs for distributed file storage system.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Estimates</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box blue"><strong>Users:</strong> 1B total, 100M DAU</div>
        <div class="d-box blue"><strong>Storage:</strong> 15GB free/user = 15EB total</div>
        <div class="d-box blue"><strong>Uploads:</strong> 10M files/day</div>
        <div class="d-box blue"><strong>Downloads:</strong> 100M files/day (10:1 read)</div>
        <div class="d-box blue"><strong>File sizes:</strong> 70% &lt;1MB, 20% 1–10MB, 10% &gt;10MB</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">NFRs</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box green"><strong>Durability:</strong> 99.999999999% (11 nines) via S3</div>
        <div class="d-box green"><strong>Availability:</strong> 99.99% upload, 99.999% download</div>
        <div class="d-box amber"><strong>Sync latency:</strong> &lt;30s across devices</div>
        <div class="d-box purple"><strong>Security:</strong> AES-256 encryption at rest + in transit</div>
        <div class="d-box gray"><strong>Versioning:</strong> 30-day history default</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-architecture",
		Title:       "High-Level Architecture",
		Description: "Upload and download paths for distributed file storage.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v" style="gap:12px">
  <div class="d-flow" style="gap:8px;align-items:center">
    <div class="d-box blue">Client<br/><small>mobile/web</small></div>
    <div class="d-arrow">→</div>
    <div class="d-box gray">API Gateway<br/><small>auth + rate limit</small></div>
    <div class="d-arrow">→</div>
    <div class="d-box green">Upload Service<br/><small>pre-signed URLs</small></div>
    <div class="d-arrow">→</div>
    <div class="d-box blue" style="font-weight:700">S3<br/><small>object storage</small></div>
  </div>
  <div style="display:flex;gap:12px;margin-left:160px">
    <div class="d-flow-v" style="gap:6px">
      <div class="d-arrow-down">↓</div>
      <div class="d-box purple">Metadata Service<br/><small>DynamoDB</small></div>
    </div>
    <div class="d-flow-v" style="gap:6px">
      <div class="d-arrow-down">↓</div>
      <div class="d-box amber">Processing Queue<br/><small>SQS → Workers</small></div>
    </div>
  </div>
  <div style="border-top:1px dashed #ccc;padding-top:10px">
    <div class="d-flow" style="gap:8px;align-items:center">
      <div class="d-box blue">Client</div>
      <div class="d-arrow">→</div>
      <div class="d-box green">CloudFront CDN<br/><small>edge cache</small></div>
      <div class="d-arrow">→</div>
      <div class="d-box blue">S3<br/><small>origin</small></div>
    </div>
    <div style="margin-top:4px;color:#64748b;font-size:12px">↑ Download path: CDN-first, 70%+ cache hit rate</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-chunking",
		Title:       "File Chunking & Content-Addressing",
		Description: "How files are split into chunks with SHA256 content-addressing for deduplication.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v" style="gap:10px">
  <div class="d-flow" style="gap:6px;align-items:center">
    <div class="d-box blue" style="min-width:120px">Large File<br/><small>500MB video</small></div>
    <div class="d-arrow">→</div>
    <div class="d-group">
      <div class="d-group-title">Split into 4MB chunks</div>
      <div class="d-flow" style="gap:4px">
        <div class="d-box green" style="font-size:11px">Chunk 1<br/>SHA256=a1b2</div>
        <div class="d-box green" style="font-size:11px">Chunk 2<br/>SHA256=c3d4</div>
        <div class="d-box green" style="font-size:11px">Chunk 3<br/>SHA256=e5f6</div>
        <div class="d-box gray" style="font-size:11px">... 125<br/>chunks</div>
      </div>
    </div>
  </div>
  <div class="d-flow" style="gap:8px;margin-top:6px">
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Chunk Manifest (stored in metadata)</div>
      <div class="d-box gray" style="font-size:11px;font-family:monospace">file_id → [a1b2, c3d4, e5f6, ..., hash125]</div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Content-Addressed Storage (S3)</div>
      <div class="d-box blue" style="font-size:11px">s3://bucket/a1b2 (Chunk 1)</div>
      <div class="d-box blue" style="font-size:11px">s3://bucket/c3d4 (Chunk 2)</div>
      <div class="d-box amber" style="font-size:11px">Same hash = same chunk = stored once!</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-data-model",
		Title:       "Data Model",
		Description: "Core tables for file metadata, chunks, folders, and permissions.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">files</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> file_id UUID</div>
        <div><span class="fk">FK</span> owner_id</div>
        <div>name, mime_type</div>
        <div>size_bytes BIGINT</div>
        <div>chunk_manifest JSONB</div>
        <div>status (uploading/active/deleted)</div>
        <div>parent_folder_id</div>
        <div>created_at, updated_at</div>
      </div>
    </div>
    <div class="d-entity" style="margin-top:8px">
      <div class="d-entity-header purple">chunks</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> chunk_id (SHA256)</div>
        <div>s3_key, size_bytes</div>
        <div>ref_count INT</div>
        <div><small>ref_count=0 → garbage collect</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header green">folders</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> folder_id UUID</div>
        <div><span class="fk">FK</span> owner_id</div>
        <div>name, parent_id (self-ref)</div>
        <div>path (materialized path)</div>
        <div>created_at</div>
      </div>
    </div>
    <div class="d-entity" style="margin-top:8px">
      <div class="d-entity-header amber">permissions</div>
      <div class="d-entity-body">
        <div><span class="fk">FK</span> file_id or folder_id</div>
        <div><span class="fk">FK</span> user_id or link_id</div>
        <div>permission (owner/editor/viewer)</div>
        <div>expires_at (for share links)</div>
        <div>created_at</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-upload-flow",
		Title:       "Upload Flow (Direct-to-S3)",
		Description: "How files are uploaded directly to S3 via pre-signed URLs, bypassing app servers.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v" style="gap:8px">
  <div class="d-flow" style="gap:6px;align-items:center">
    <div class="d-box blue">Client</div>
    <div class="d-arrow">→ POST /files</div>
    <div class="d-box green">Upload Service</div>
  </div>
  <div style="margin-left:120px;color:#64748b;font-size:12px">↓ creates metadata (status=uploading) + generates pre-signed S3 URL</div>
  <div class="d-flow" style="gap:6px;align-items:center;margin-left:40px">
    <div class="d-box green">Upload Service</div>
    <div class="d-arrow">→ returns URL</div>
    <div class="d-box blue">Client</div>
  </div>
  <div style="margin-left:100px;color:#059669;font-size:12px;font-weight:600">↓ Client uploads directly to S3 (bypasses servers — no bandwidth bottleneck!)</div>
  <div class="d-flow" style="gap:6px;align-items:center;margin-left:40px">
    <div class="d-box blue">Client</div>
    <div class="d-arrow">→ PUT chunks</div>
    <div class="d-box blue" style="font-weight:700">S3</div>
    <div class="d-arrow">→ S3 Event</div>
    <div class="d-box amber">SQS</div>
  </div>
  <div style="margin-left:200px;color:#64748b;font-size:12px">↓ SQS triggers workers</div>
  <div class="d-flow" style="gap:8px;margin-left:160px">
    <div class="d-box gray" style="font-size:11px">Virus Scan</div>
    <div class="d-box gray" style="font-size:11px">Thumbnail Gen</div>
    <div class="d-box gray" style="font-size:11px">Search Index</div>
    <div class="d-box green" style="font-size:11px">status=active</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-deduplication",
		Title:       "Content-Based Deduplication",
		Description: "How SHA256 chunk hashing eliminates duplicate storage across users.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">User A uploads vacation.mp4</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box blue">SHA256(file) = abc123</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box green">Chunk abc123 NOT in S3</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box green">Upload to s3://bucket/abc123</div>
        <div class="d-box gray" style="font-size:11px">chunks: {abc123, ref_count=1}</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">User B uploads same file</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box blue">SHA256(file) = abc123</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box amber">Chunk abc123 EXISTS in S3!</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box green">ref_count++ → no S3 upload</div>
        <div class="d-box gray" style="font-size:11px">chunks: {abc123, ref_count=2}</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Result</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box green">1 copy stored, 2 users served</div>
        <div class="d-box green">30% duplicate rate → 30% cost savings</div>
        <div class="d-box gray" style="font-size:11px">On delete: ref_count-- → GC when 0</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-sync-protocol",
		Title:       "Device Sync Protocol",
		Description: "Delta sync protocol using cursors to efficiently sync changes across devices.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v" style="gap:10px">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Delta Sync (Normal)</div>
        <div class="d-flow-v" style="gap:6px">
          <div class="d-box blue">Device stores local_cursor<br/><small>(last sync timestamp)</small></div>
          <div class="d-arrow-down">↓</div>
          <div class="d-box gray" style="font-size:11px">GET /changes?since={cursor}</div>
          <div class="d-arrow-down">↓</div>
          <div class="d-box green">Server returns changed file_ids only</div>
          <div class="d-arrow-down">↓</div>
          <div class="d-box green">Download only changed files</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Conflict (Two devices offline)</div>
        <div class="d-flow-v" style="gap:6px">
          <div class="d-box blue">Device A edits file.doc offline</div>
          <div class="d-box blue">Device B edits same file offline</div>
          <div class="d-arrow-down">↓ both reconnect</div>
          <div class="d-box amber">Conflict detected!</div>
          <div class="d-arrow-down">↓ Dropbox strategy</div>
          <div class="d-box green">file.doc (Device A's version)</div>
          <div class="d-box green">file (Device B's conflicted copy).doc</div>
          <div class="d-box gray" style="font-size:11px">User notified, keeps both</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-versioning",
		Title:       "File Versioning Strategy",
		Description: "Efficient versioning using snapshots and diffs to minimize storage.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v" style="gap:10px">
  <div class="d-flow" style="gap:6px;align-items:center">
    <div class="d-box green" style="font-size:11px">Snapshot<br/>v0 (full)</div>
    <div class="d-arrow">→</div>
    <div class="d-box gray" style="font-size:11px">diff v1</div>
    <div class="d-arrow">→</div>
    <div class="d-box gray" style="font-size:11px">diff v2</div>
    <div class="d-arrow">→</div>
    <div class="d-box gray" style="font-size:11px">... v9</div>
    <div class="d-arrow">→</div>
    <div class="d-box green" style="font-size:11px">Snapshot<br/>v10 (full)</div>
    <div class="d-arrow">→</div>
    <div class="d-box gray" style="font-size:11px">diff v11...</div>
  </div>
  <div class="d-cols" style="margin-top:8px">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Restore to v7</div>
        <div class="d-flow-v" style="gap:4px">
          <div class="d-box blue" style="font-size:11px">Load Snapshot v0</div>
          <div class="d-box gray" style="font-size:11px">Apply diff v1–v7</div>
          <div class="d-box green" style="font-size:11px">Document at v7 ✓</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Storage Savings</div>
        <div class="d-flow-v" style="gap:4px">
          <div class="d-box red" style="font-size:11px">Full copies: 30 × 10MB = 300MB</div>
          <div class="d-box green" style="font-size:11px">Snapshot+diffs: 10MB + 29×50KB = 11.5MB</div>
          <div class="d-box green" style="font-size:11px"><strong>26x storage reduction</strong></div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-permissions",
		Title:       "Sharing & Access Control",
		Description: "Permission levels and share link design for file access control.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Permission Levels</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box blue"><strong>Owner</strong> — full control (edit, share, delete)</div>
        <div class="d-box green"><strong>Editor</strong> — read + write, cannot share</div>
        <div class="d-box amber"><strong>Commenter</strong> — read + add comments</div>
        <div class="d-box gray"><strong>Viewer</strong> — read only</div>
        <div class="d-box purple"><strong>Public link</strong> — anonymous read (time-limited)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Share Link Flow</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box blue">Owner creates share link</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box gray" style="font-size:11px">permissions: {link_id, file_id, viewer, expires=7d}</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box green">Recipient opens link → JWT with permission claim</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box green">Every file access → check JWT + ACL table</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-caching",
		Title:       "Caching Strategy",
		Description: "Three-layer caching for file metadata and content delivery.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v" style="gap:8px">
  <div class="d-flow" style="gap:8px;align-items:center">
    <div class="d-box blue">Client Request</div>
    <div class="d-arrow">→</div>
    <div class="d-box green">Browser Cache<br/><small>Cache-Control: max-age=3600</small></div>
    <div class="d-arrow">→ miss</div>
    <div class="d-box green">CloudFront Edge<br/><small>70%+ hit rate for files</small></div>
    <div class="d-arrow">→ miss</div>
    <div class="d-box green">Redis Metadata<br/><small>TTL=5min for file info</small></div>
    <div class="d-arrow">→ miss</div>
    <div class="d-box blue">DynamoDB<br/><small>source of truth</small></div>
  </div>
  <div class="d-group" style="margin-top:8px">
    <div class="d-group-title">Cache Invalidation on Update</div>
    <div class="d-flow" style="gap:8px">
      <div class="d-box amber">File updated</div>
      <div class="d-arrow">→</div>
      <div class="d-box red" style="font-size:11px">Invalidate CloudFront</div>
      <div class="d-box red" style="font-size:11px">DEL Redis key</div>
      <div class="d-box amber" style="font-size:11px">New ETag for browser</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-failure-handling",
		Title:       "Failure Scenarios",
		Description: "How the system handles upload interruption, S3 failure, and metadata DB failure.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Upload Interrupted</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box amber">Network drops mid-upload</div>
        <div class="d-arrow-down">↓ S3 Multipart</div>
        <div class="d-box green">Resume from last part</div>
        <div class="d-box gray" style="font-size:11px">UploadId saved client-side</div>
        <div class="d-box gray" style="font-size:11px">S3 retains incomplete parts 7 days</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">S3 Regional Failure</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box red">S3 us-east-1 down</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box green">S3 Cross-Region Replication → eu-west-1</div>
        <div class="d-box green">Route53 failover to EU endpoint</div>
        <div class="d-box gray" style="font-size:11px">RPO: near-zero (async replication)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Metadata DB Failure</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box red">DynamoDB unreachable</div>
        <div class="d-arrow-down">↓</div>
        <div class="d-box amber">Serve from Redis cache (read-only)</div>
        <div class="d-box amber">Queue writes to SQS</div>
        <div class="d-box green">Replay queued writes on recovery</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fs-monitoring",
		Title:       "Monitoring & Cost Analysis",
		Description: "Key metrics and cost breakdown for the file storage system.",
		ContentFile: "problems/file-storage",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Key Metrics</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box green">Upload success rate: target &gt;99.9%</div>
        <div class="d-box green">Upload P99 latency: &lt;5s for 10MB</div>
        <div class="d-box green">Download P99: &lt;500ms (CDN-cached)</div>
        <div class="d-box amber">Sync lag: &lt;30s across devices</div>
        <div class="d-box blue">Storage utilization %</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Cost Breakdown</div>
      <div class="d-flow-v" style="gap:6px">
        <div class="d-box blue">S3 Standard: $0.023/GB → 1PB = $23K/mo</div>
        <div class="d-box blue">S3 Intelligent-Tiering auto-moves cold files</div>
        <div class="d-box blue">CloudFront: $0.085/GB egress</div>
        <div class="d-box green">Deduplication saves ~30% storage cost</div>
        <div class="d-box green">Compression saves another 40% on documents</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

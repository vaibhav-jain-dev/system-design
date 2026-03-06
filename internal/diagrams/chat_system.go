package diagrams

func registerChatSystem(r *Registry) {
	r.Register(&Diagram{
		Slug:        "cs-requirements",
		Title:       "Requirements & Scale Estimates",
		Description: "Non-functional requirements and scale targets for a WhatsApp-style chat system",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green">1:1 real-time messaging</div>
        <div class="d-box green">Group chat (up to 256 members)</div>
        <div class="d-box green">Sent/delivered/read receipts</div>
        <div class="d-box green">Media sharing (images, video, docs)</div>
        <div class="d-box blue">Online/offline presence</div>
        <div class="d-box blue">Multi-device sync</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">2B total users, 500M DAU</div>
        <div class="d-box purple">100B messages/day &#8776; 1.15M msgs/sec</div>
        <div class="d-box purple">Delivery latency: &lt; 100ms P99</div>
        <div class="d-box amber">Max group size: 256 members</div>
        <div class="d-box amber">E2E encryption (Signal Protocol)</div>
        <div class="d-box amber">99.99% availability</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Storage Estimates</div>
      <div class="d-flow-v">
        <div class="d-box gray">100B msgs/day &times; 100 bytes avg = 10 TB/day</div>
        <div class="d-box gray">Media: ~50 TB/day (images + video)</div>
        <div class="d-box gray">5-year retention = 18 PB text + 90 PB media</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-api-design",
		Title:       "API Design",
		Description: "Core API endpoints for messaging, conversations, groups, and delivery status",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">REST Endpoints</div>
    <div class="d-flow-v">
      <div class="d-box green">POST /messages &#8212; Send message (text, media_key, conversation_id)</div>
      <div class="d-box green">GET /conversations &#8212; List user conversations (cursor-paginated)</div>
      <div class="d-box blue">POST /groups &#8212; Create group (name, members[], icon)</div>
      <div class="d-box blue">PUT /messages/{id}/status &#8212; Update delivery status (delivered/read)</div>
      <div class="d-box gray">GET /messages?conversation_id=X&amp;cursor=Y &#8212; Fetch message history</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">WebSocket</div>
    <div class="d-flow-v">
      <div class="d-box purple">WS /ws/chat &#8212; Persistent bidirectional connection</div>
      <div class="d-flow">
        <div class="d-box amber">Events: message, typing, presence, receipt</div>
        <div class="d-arrow">&#8594;</div>
        <div class="d-box amber">Auth via JWT token on connect</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-message-flow",
		Title:       "Message Delivery Flow",
		Description: "End-to-end message flow from sender to receiver including offline path",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Online Delivery Path</div>
    <div class="d-flow">
      <div class="d-box green">Sender</div>
      <div class="d-arrow">WebSocket &#8594;</div>
      <div class="d-box blue">Chat Server A</div>
      <div class="d-arrow">publish &#8594;</div>
      <div class="d-box purple">Kafka</div>
      <div class="d-arrow">consume &#8594;</div>
      <div class="d-box blue">Chat Server B</div>
      <div class="d-arrow">WebSocket &#8594;</div>
      <div class="d-box green">Receiver</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Offline Delivery Path</div>
    <div class="d-flow">
      <div class="d-box green">Sender</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Chat Server</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Kafka</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Offline Queue (DynamoDB)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">Push Service (APNs/FCM)</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Persistence (parallel)</div>
    <div class="d-flow">
      <div class="d-box purple">Kafka</div>
      <div class="d-arrow">async &#8594;</div>
      <div class="d-box indigo">DynamoDB (messages table)</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-architecture",
		Title:       "High-Level Architecture",
		Description: "Full system architecture for a WhatsApp-scale chat system",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box gray">CDN (media)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">NLB (Layer 4)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">WebSocket Servers</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box green">WebSocket Servers</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">Chat Service</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple">Kafka (message bus)</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Storage</div>
        <div class="d-flow-v">
          <div class="d-box indigo">DynamoDB (messages)</div>
          <div class="d-box indigo">DynamoDB (conversations)</div>
          <div class="d-box amber">S3 (media files)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Supporting Services</div>
        <div class="d-flow-v">
          <div class="d-box red">Push Service (APNs / FCM)</div>
          <div class="d-box green">Presence Service (Redis)</div>
          <div class="d-box blue">User Service</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-group-chat",
		Title:       "Group Chat Fan-out",
		Description: "Fan-out strategy for delivering messages to group members",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box green">Sender</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">Chat Service</div>
    <div class="d-arrow">lookup members &#8594;</div>
    <div class="d-box purple">Group Service</div>
  </div>
  <div class="d-arrow-down">&#8595; fan-out (max 256 members)</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Online Members</div>
        <div class="d-flow-v">
          <div class="d-box green">Member A &#8592; WebSocket</div>
          <div class="d-box green">Member B &#8592; WebSocket</div>
          <div class="d-box green">Member C &#8592; WebSocket</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Offline Members</div>
        <div class="d-flow-v">
          <div class="d-box amber">Member D &#8594; Offline Queue</div>
          <div class="d-box amber">Member E &#8594; Offline Queue</div>
          <div class="d-box red">Push notification (APNs/FCM)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Why fan-out-on-write?</div>
    <div class="d-flow-v">
      <div class="d-box gray">Max 256 members &#8212; bounded fan-out cost</div>
      <div class="d-box gray">256 writes &times; 1ms = 256ms worst case (parallelized to ~10ms)</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-e2e-encryption",
		Title:       "End-to-End Encryption",
		Description: "Signal Protocol key exchange and message encryption flow",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Key Exchange (X3DH)</div>
    <div class="d-flow">
      <div class="d-box green">Sender</div>
      <div class="d-arrow">fetch prekey bundle &#8594;</div>
      <div class="d-box blue">Key Server</div>
      <div class="d-arrow">&#8592; identity + signed + one-time keys</div>
      <div class="d-box green">Receiver</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; derive shared secret</div>
  <div class="d-group">
    <div class="d-group-title">Double Ratchet</div>
    <div class="d-flow">
      <div class="d-box purple">DH Ratchet</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Chain Key Ratchet</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Message Key</div>
    </div>
    <div class="d-flow-v">
      <div class="d-box gray">New key per message &#8212; forward secrecy</div>
      <div class="d-box gray">Compromised key cannot decrypt past messages</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Message Path</div>
    <div class="d-flow">
      <div class="d-box green">Plaintext</div>
      <div class="d-arrow">encrypt &#8594;</div>
      <div class="d-box amber">Ciphertext</div>
      <div class="d-arrow">transport (opaque) &#8594;</div>
      <div class="d-box amber">Ciphertext</div>
      <div class="d-arrow">decrypt &#8594;</div>
      <div class="d-box green">Plaintext</div>
    </div>
    <div class="d-box red">Server NEVER sees plaintext &#8212; zero-knowledge transport</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-offline-queue",
		Title:       "Offline Message Queue",
		Description: "How messages are queued and synced when the recipient is offline",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Message Arrives &#8212; User Offline</div>
    <div class="d-flow">
      <div class="d-box blue">Chat Service</div>
      <div class="d-arrow">user offline? &#8594;</div>
      <div class="d-box purple">Presence Service (Redis)</div>
    </div>
    <div class="d-arrow-down">&#8595; confirmed offline</div>
    <div class="d-flow">
      <div class="d-box amber">Store in DynamoDB</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">Send push notification (APNs/FCM)</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; user comes online</div>
  <div class="d-group">
    <div class="d-group-title">Sync on Reconnect</div>
    <div class="d-flow">
      <div class="d-box green">Client connects</div>
      <div class="d-arrow">last_sync_ts &#8594;</div>
      <div class="d-box blue">Chat Server</div>
      <div class="d-arrow">query &#8594;</div>
      <div class="d-box indigo">DynamoDB (messages &gt; last_sync_ts)</div>
    </div>
    <div class="d-arrow-down">&#8595;</div>
    <div class="d-flow">
      <div class="d-box green">Client receives batch</div>
      <div class="d-arrow">&#8592; stream messages</div>
      <div class="d-box blue">Chat Server</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-data-model",
		Title:       "Data Model",
		Description: "Core entities and relationships for a chat system",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">users</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> user_id (UUID)</div>
        <div class="d-box gray">phone_number, display_name</div>
        <div class="d-box gray">avatar_url, last_seen_at</div>
        <div class="d-box gray">public_identity_key</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">conversations</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> conversation_id (UUID)</div>
        <div class="d-box gray">type (1:1 | group)</div>
        <div class="d-box gray">group_name, group_icon_url</div>
        <div class="d-box gray">created_at, updated_at</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">participants</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> conversation_id + user_id</div>
        <div class="d-box gray">role (admin | member)</div>
        <div class="d-box gray">joined_at, muted_until</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">messages</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> conversation_id (partition)</div>
        <div class="d-box blue"><span class="pk">SK</span> message_id (Snowflake, sortable)</div>
        <div class="d-box gray">sender_id, encrypted_body</div>
        <div class="d-box gray">media_key, media_url, type</div>
        <div class="d-box gray">created_at</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">read_receipts</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> conversation_id + user_id</div>
        <div class="d-box gray">last_read_message_id</div>
        <div class="d-box gray">updated_at</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">devices</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> user_id + device_id</div>
        <div class="d-box gray">push_token, platform (iOS/Android)</div>
        <div class="d-box gray">prekey_bundle, sync_cursor</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-read-receipts",
		Title:       "Read Receipts & Typing Indicators",
		Description: "Flow for delivery status updates and ephemeral typing signals",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Message Status Lifecycle</div>
    <div class="d-flow">
      <div class="d-box blue">Sent &#10003;</div>
      <div class="d-arrow">server ACK &#8594;</div>
      <div class="d-box green">Delivered &#10003;&#10003;</div>
      <div class="d-arrow">client ACK &#8594;</div>
      <div class="d-box purple">Read (blue &#10003;&#10003;)</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Read Receipt Flow</div>
    <div class="d-flow">
      <div class="d-box green">Receiver opens chat</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Client sends read ACK</div>
      <div class="d-arrow">WebSocket &#8594;</div>
      <div class="d-box blue">Chat Server</div>
      <div class="d-arrow">fan-out &#8594;</div>
      <div class="d-box green">Sender sees blue &#10003;&#10003;</div>
    </div>
    <div class="d-flow">
      <div class="d-box amber">Batch reads: send last_read_message_id (not per-message)</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Typing Indicators</div>
    <div class="d-flow">
      <div class="d-box green">User types</div>
      <div class="d-arrow">ephemeral &#8594;</div>
      <div class="d-box blue">Chat Server</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Other participant</div>
    </div>
    <div class="d-flow-v">
      <div class="d-box gray">Fire-and-forget &#8212; no persistence, no Kafka</div>
      <div class="d-box gray">Throttled to 1 event / 3 seconds per user</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-scaling",
		Title:       "Scaling Strategy",
		Description: "Sharding and partitioning strategy for WebSocket, Kafka, and DynamoDB",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">WebSocket Servers</div>
      <div class="d-flow-v">
        <div class="d-box blue">Shard by user_id (consistent hashing)</div>
        <div class="d-box gray">~100K connections per server</div>
        <div class="d-box gray">500M DAU / 100K = 5,000 WS servers</div>
        <div class="d-box amber">Session registry in Redis: user &#8594; server mapping</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Kafka</div>
      <div class="d-flow-v">
        <div class="d-box purple">Partition by conversation_id</div>
        <div class="d-box gray">Preserves message ordering per conversation</div>
        <div class="d-box gray">~10K partitions for 1M+ msgs/sec throughput</div>
        <div class="d-box amber">Consumer group per region</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">DynamoDB</div>
      <div class="d-flow-v">
        <div class="d-box indigo">Partition key: conversation_id</div>
        <div class="d-box indigo">Sort key: message_id (Snowflake)</div>
        <div class="d-box gray">Hot partition mitigation: write sharding for viral groups</div>
        <div class="d-box gray">On-demand capacity &#8212; auto-scales to 1M+ WCU</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-multi-device",
		Title:       "Multi-Device Sync",
		Description: "How messages are replicated across a user's devices with per-device cursors",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Message Replication</div>
    <div class="d-flow">
      <div class="d-box blue">Incoming message</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Chat Server</div>
      <div class="d-arrow">fan-out to all devices &#8594;</div>
      <div class="d-flow-v">
        <div class="d-box green">Phone</div>
        <div class="d-box green">Tablet</div>
        <div class="d-box green">Desktop</div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Per-Device Sync Cursor</div>
    <div class="d-flow-v">
      <div class="d-flow">
        <div class="d-box green">Phone</div>
        <div class="d-arrow">sync_cursor: msg_1042 &#8594;</div>
        <div class="d-box indigo">devices table</div>
      </div>
      <div class="d-flow">
        <div class="d-box green">Desktop</div>
        <div class="d-arrow">sync_cursor: msg_1038 &#8594;</div>
        <div class="d-box indigo">devices table</div>
      </div>
      <div class="d-box gray">Each device tracks its own cursor &#8212; reconnect fetches only missed messages</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Conflict Resolution</div>
    <div class="d-flow-v">
      <div class="d-box amber">Snowflake IDs provide total ordering &#8212; no conflicts on messages</div>
      <div class="d-box amber">Read receipts: last-writer-wins (latest timestamp)</div>
      <div class="d-box gray">Delete: tombstone propagated to all devices on next sync</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "cs-media-sharing",
		Title:       "Media Message Pipeline",
		Description: "Upload, encrypt, store, and deliver media messages",
		ContentFile: "problems/chat-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Sender Upload</div>
    <div class="d-flow">
      <div class="d-box green">Client generates AES-256 media_key</div>
      <div class="d-arrow">encrypt &#8594;</div>
      <div class="d-box amber">Encrypted media blob</div>
      <div class="d-arrow">upload &#8594;</div>
      <div class="d-box blue">Media Service</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box indigo">S3</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; S3 returns media_url</div>
  <div class="d-group">
    <div class="d-group-title">Thumbnail Generation</div>
    <div class="d-flow">
      <div class="d-box indigo">S3 event</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Lambda (resize + blur)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box indigo">S3 (thumbnail bucket)</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Message Delivery</div>
    <div class="d-flow">
      <div class="d-box green">Sender</div>
      <div class="d-arrow">send message (media_url + media_key) &#8594;</div>
      <div class="d-box blue">Chat Server</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Receiver</div>
    </div>
    <div class="d-box gray">media_key is E2E encrypted &#8212; server cannot decrypt media</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Receiver Download</div>
    <div class="d-flow">
      <div class="d-box green">Receiver</div>
      <div class="d-arrow">download from CDN &#8594;</div>
      <div class="d-box gray">CDN / S3</div>
      <div class="d-arrow">decrypt with media_key &#8594;</div>
      <div class="d-box green">Plaintext media</div>
    </div>
  </div>
</div>`,
	})
}

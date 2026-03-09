package diagrams

func registerNotificationSystem(r *Registry) {
	r.Register(&Diagram{
		Slug:        "notif-requirements",
		Title:       "Requirements & Scale",
		Description: "Scale targets and channel breakdown for a notification system serving 1B users",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Push (FCM/APNs), Email (SendGrid/SES), SMS (Twilio), and in-app notifications. Each channel has different latency guarantees and cost profiles.">Multi-channel delivery: push, email, SMS, in-app</div>
        <div class="d-box green" data-tip="Users can opt out per channel (no email), per category (no marketing), or globally. Opt-outs must be respected within 10 seconds of preference change.">User opt-out &amp; preferences per channel + category</div>
        <div class="d-box green" data-tip="OTP and transaction alerts must arrive in &lt;500ms. Marketing batch can take minutes. Priority tiers drive Kafka topic routing and consumer polling order.">Priority tiers: real-time, high, normal, low (batch)</div>
        <div class="d-box blue" data-tip="Events fire from multiple upstream systems: user actions (like, comment), payment service, marketing campaigns. Notification Service is downstream consumer.">Event-driven: user actions, payments, marketing triggers</div>
        <div class="d-box blue" data-tip="Template library with variable injection (name, amount, order_id). Localized per user locale. Stored in DynamoDB, cached in Redis.">Templated messages with i18n support</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="1B users × 10 notifications/day average = 10B/day ÷ 86400s ≈ 115K notifications/sec peak. Push dominates at 80% of volume.">1B users, 10M notifications/min peak <span class="d-metric throughput">167K/sec</span></div>
        <div class="d-box amber" data-tip="Push: FCM/APNs latency target &lt;1s delivery confirmation. OTP SMS: &lt;3s. Transactional email: &lt;5s. Marketing email: best-effort within 1 hour.">Push &lt;1s delivery, SMS &lt;3s, email &lt;5s (transactional)</div>
        <div class="d-box amber" data-tip="99.9% = 8.7 hours downtime/year. For notifications this means some messages may be delayed but not lost — at-least-once delivery guarantee.">Availability: 99.9% (delay acceptable, no message loss)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Channel Volume Split</div>
      <div class="d-flow-v">
        <div class="d-box green">Push (FCM/APNs) — <strong>8M/min</strong> <span class="d-metric throughput">80%</span></div>
        <div class="d-box blue">Email (SendGrid/SES) — <strong>1M/min</strong> <span class="d-metric throughput">10%</span></div>
        <div class="d-box amber">SMS (Twilio) — <strong>100K/min</strong> <span class="d-metric throughput">1%</span></div>
        <div class="d-box gray">In-app (WebSocket) — <strong>900K/min</strong> <span class="d-metric throughput">9%</span></div>
      </div>
    </div>
  </div>
</div>
<div class="d-flow">
  <div class="d-number"><div class="d-number-value">1B</div><div class="d-number-label">Users</div></div>
  <div class="d-number"><div class="d-number-value">10M</div><div class="d-number-label">Notifs/min</div></div>
  <div class="d-number"><div class="d-number-value">&lt;1s</div><div class="d-number-label">Push Latency</div></div>
  <div class="d-number"><div class="d-number-value">4</div><div class="d-number-label">Channels</div></div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end notification flow from event sources through Kafka to third-party providers",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Event Sources</div>
    <div class="d-flow">
      <div class="d-box blue" data-tip="User posts, likes, follows trigger notifications for followers/mentioned users. High volume, medium priority.">User Actions<br><small>post/like/follow</small></div>
      <div class="d-box amber" data-tip="Payment success, failure, refund. High priority — user expects immediate confirmation. Routes to push + email.">Payment Service<br><small>txn alerts</small></div>
      <div class="d-box purple" data-tip="Scheduled campaigns, A/B tested content. Low priority, bulk volume. Can be delayed minutes without user impact.">Marketing<br><small>campaigns</small></div>
      <div class="d-box red" data-tip="OTP codes, security alerts. Highest priority. Dedicated Kafka partition, consumer polls this first.">Auth Service<br><small>OTP / alerts</small></div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; publish events</div>
  <div class="d-flow">
    <div class="d-col">
      <div class="d-box green" data-tip="Validates payload, looks up user preferences, checks opt-outs, applies rate limits, selects channel, routes to correct Kafka topic. Stateless — scale horizontally."><strong>Notification Service</strong><br><small>validate → route → enqueue</small></div>
    </div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Kafka Topics (by priority)</div>
        <div class="d-flow-v">
          <div class="d-box red">notif.high — OTP, txn alerts</div>
          <div class="d-box amber">notif.normal — social, likes</div>
          <div class="d-box gray">notif.low — marketing batch</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; consume &amp; dispatch</div>
  <div class="d-flow">
    <div class="d-group">
      <div class="d-group-title">Channel Routers (Consumer Groups)</div>
      <div class="d-flow">
        <div class="d-box green" data-tip="Batches 1000 tokens, calls FCM batch API or APNs HTTP/2 multiplexing. Exponential backoff on 500s. Tracks delivery receipts.">Push Router<br><small>FCM / APNs</small></div>
        <div class="d-box blue" data-tip="Calls SendGrid or SES. Handles bounce/complaint webhooks. Unsubscribe one-click link injected per CAN-SPAM.">Email Router<br><small>SendGrid / SES</small></div>
        <div class="d-box amber" data-tip="Twilio REST API. Most expensive channel ($0.0079/SMS). Hard rate limit: 1 SMS/hour/user for non-OTP.">SMS Router<br><small>Twilio</small></div>
        <div class="d-box purple" data-tip="Redis Pub/Sub or WebSocket push for users with open connections. Falls back to storing in DB for next app open.">In-App Router<br><small>WebSocket / Redis</small></div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; deliver</div>
  <div class="d-flow">
    <div class="d-box green">FCM (Android)</div>
    <div class="d-box blue">APNs (iOS)</div>
    <div class="d-box amber">SendGrid</div>
    <div class="d-box red">Twilio</div>
    <div class="d-box purple">User Devices</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-data-model",
		Title:       "Data Model",
		Description: "Core tables for notifications, user preferences, and device tokens",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header green">notifications</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> notification_id UUID</div>
        <div><span class="fk">FK</span> user_id</div>
        <div>channel ENUM(push,email,sms,inapp)</div>
        <div><span class="fk">FK</span> template_id</div>
        <div>status ENUM(pending,sent,delivered,failed)</div>
        <div>priority ENUM(high,normal,low)</div>
        <div>created_at TIMESTAMP</div>
        <div>sent_at TIMESTAMP</div>
        <div>delivered_at TIMESTAMP</div>
        <div>retry_count INT DEFAULT 0</div>
        <div><span class="idx idx-btree">IDX</span> (user_id, created_at DESC)</div>
        <div><span class="idx idx-hash">IDX</span> (status) for DLQ queries</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">user_preferences</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> user_id</div>
        <div>push_enabled BOOL DEFAULT true</div>
        <div>email_enabled BOOL DEFAULT true</div>
        <div>sms_enabled BOOL DEFAULT false</div>
        <div>inapp_enabled BOOL DEFAULT true</div>
        <div>quiet_hours_start TIME <small>(e.g. 22:00)</small></div>
        <div>quiet_hours_end TIME <small>(e.g. 08:00)</small></div>
        <div>timezone VARCHAR(64)</div>
        <div>marketing_opt_out BOOL DEFAULT false</div>
        <div>updated_at TIMESTAMP</div>
      </div>
    </div>
    <div class="d-entity" style="margin-top:8px">
      <div class="d-entity-header amber">device_tokens</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> token_id UUID</div>
        <div><span class="fk">FK</span> user_id</div>
        <div>platform ENUM(ios,android,web)</div>
        <div>token VARCHAR(512) <small>FCM/APNs token</small></div>
        <div>active BOOL DEFAULT true</div>
        <div>updated_at TIMESTAMP</div>
        <div><span class="idx idx-btree">IDX</span> (user_id, platform)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header purple">notification_templates</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> template_id UUID</div>
        <div>channel ENUM(push,email,sms,inapp)</div>
        <div>name VARCHAR(128)</div>
        <div>subject VARCHAR(256) <small>email only</small></div>
        <div>body TEXT <small>Handlebars template</small></div>
        <div>variables JSONB <small>["name","order_id"]</small></div>
        <div>locale VARCHAR(8) DEFAULT 'en'</div>
        <div>version INT</div>
        <div>created_at TIMESTAMP</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-fan-out",
		Title:       "Notification Fan-Out",
		Description: "How a single event fans out to millions of users via Kafka consumer groups",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue" data-tip="A celebrity with 10M followers posts. This triggers one new_post event. The fan-out happens downstream in the notification pipeline.">User with <strong>10M followers</strong> posts</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green" data-tip="Notification Service publishes one event to Kafka. Kafka handles durability and fan-out — no need to write 10M records synchronously.">Notification Service<br><small>1 Kafka event published</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple" data-tip="Kafka topic: notif.normal. Partitioned by user_id for ordering guarantees per user. Replication factor 3 for durability.">Kafka Topic<br><small>new_post_notifications<br>32 partitions</small></div>
  </div>
  <div class="d-arrow-down">&#8595; consumer group reads follower batches</div>
  <div class="d-group">
    <div class="d-group-title">Consumer Group — 10 consumers (auto-scaling to 100 under lag)</div>
    <div class="d-flow">
      <div class="d-box green" data-tip="Each consumer reads from assigned partitions. Fetches 1000 follower user_ids from FollowerDB, looks up device tokens, checks preferences in batch.">Consumer 1<br><small>batch 1000 users</small></div>
      <div class="d-box green" data-tip="Parallel processing. No coordination needed between consumers — each handles its partition independently.">Consumer 2<br><small>batch 1000 users</small></div>
      <div class="d-box gray">Consumer 3–9<br><small>…</small></div>
      <div class="d-box green">Consumer 10<br><small>batch 1000 users</small></div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; batch push to providers</div>
  <div class="d-flow">
    <div class="d-box blue">FCM Batch API<br><small>up to 500 tokens/call</small></div>
    <div class="d-box amber">APNs HTTP/2<br><small>multiplexed streams</small></div>
  </div>
  <div class="d-group" style="margin-top:12px">
    <div class="d-group-title">Fan-Out Timeline</div>
    <div class="d-flow">
      <div class="d-box purple" data-tip="10M followers ÷ 10 consumers = 1M per consumer. At 1000/sec per consumer = 1000s ≈ 17 min. Scale to 100 consumers → ~100 seconds.">10 consumers @ 1K/sec → <strong>~17 min</strong> for 10M</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green" data-tip="Kafka consumer lag triggers auto-scaling. At lag &gt; 50K messages, add 10 more consumers. Scales elastically.">Auto-scale to 100 consumers → <strong>~100 sec</strong></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber" data-tip="For VIP users with 100M+ followers, use a separate high-volume fan-out pipeline with pre-allocated consumer capacity.">VIP tier: pre-scale for celebrities</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-priority-queue",
		Title:       "Priority Queue Design",
		Description: "Kafka topic separation by priority and consumer polling order",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Kafka Topics by Priority</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="OTP codes, transaction alerts, security events. SLA: &lt;500ms delivery. Dedicated consumers always polling. Never starved by lower priority work.">notif.critical<br><small>OTP, txn alerts, security</small><br><span class="d-metric latency">&lt;500ms SLA</span></div>
        <div class="d-box amber" data-tip="Likes, comments, follows, mentions. SLA: &lt;5s. Polled when critical topic is empty or after every 5 critical messages (anti-starvation).">notif.high<br><small>social interactions, mentions</small><br><span class="d-metric latency">&lt;5s SLA</span></div>
        <div class="d-box blue" data-tip="General app notifications, reminders. SLA: &lt;30s. Acceptable to batch with other normal traffic.">notif.normal<br><small>reminders, updates</small><br><span class="d-metric latency">&lt;30s SLA</span></div>
        <div class="d-box gray" data-tip="Marketing campaigns, newsletters. SLA: minutes to hours. Can be held overnight if quiet hours apply. Rate-limited to prevent spam.">notif.low<br><small>marketing, newsletters</small><br><span class="d-metric latency">1–60 min SLA</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Consumer Polling Strategy (weighted)</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Critical consumers always poll first. Dedicated thread pool. Never yield CPU for lower priority work.">1. Always poll notif.critical first</div>
        <div class="d-arrow-down">&#8595; if empty</div>
        <div class="d-box amber" data-tip="Poll high priority. Anti-starvation: after 5 consecutive critical messages, poll high once regardless of critical queue state.">2. Poll notif.high<br><small>Anti-starvation: poll every 5 critical msgs</small></div>
        <div class="d-arrow-down">&#8595; if empty</div>
        <div class="d-box blue">3. Poll notif.normal</div>
        <div class="d-arrow-down">&#8595; if empty</div>
        <div class="d-box gray" data-tip="Low priority only gets CPU when all higher-priority queues are empty. In high-load scenarios, marketing may be delayed by hours — acceptable.">4. Poll notif.low (background)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Latency Comparison</div>
      <div class="d-flow-v">
        <div class="d-box red">OTP → <span class="d-metric latency">&lt;500ms</span><br><small>Twilio SMS direct path</small></div>
        <div class="d-box amber">Like notification → <span class="d-metric latency">&lt;5s</span><br><small>FCM push via high queue</small></div>
        <div class="d-box blue">App reminder → <span class="d-metric latency">&lt;30s</span><br><small>normal queue batch</small></div>
        <div class="d-box gray">Marketing email → <span class="d-metric latency">1–60 min</span><br><small>low queue, quiet hours checked</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-delivery-guarantee",
		Title:       "Delivery Guarantees",
		Description: "At-least-once delivery with exponential backoff retry and dead-letter queue",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box green" data-tip="Consumer dequeues from Kafka. Does NOT commit offset until delivery confirmed. If consumer crashes before ack, message reprocessed from last committed offset.">Dequeue from Kafka<br><small>offset NOT committed yet</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue" data-tip="HTTP call to FCM/APNs/SendGrid/Twilio. Timeout: 3 seconds. If no response in 3s, treat as failure and retry.">Send to Provider<br><small>3s timeout</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber" data-tip="Provider returns 200 + delivery receipt. Commit Kafka offset. Update notification status=delivered in DB. Emit delivery metric.">Provider ACK<br><small>commit Kafka offset</small></div>
  </div>
  <div class="d-arrow-down">&#8595; on failure (timeout or 5xx)</div>
  <div class="d-group">
    <div class="d-group-title">Exponential Backoff Retry</div>
    <div class="d-flow">
      <div class="d-box red" data-tip="First retry: wait 1 second. Use jitter (±20%) to avoid thundering herd when a provider recovers.">Retry 1<br><small>wait 1s (±jitter)</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red" data-tip="Second retry: wait 2 seconds. Still tracking against same notification_id for deduplication.">Retry 2<br><small>wait 2s</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red" data-tip="Third and final retry: wait 4 seconds. After this, the message is routed to the dead-letter queue.">Retry 3<br><small>wait 4s</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box gray" data-tip="DLQ stores failed messages for manual inspection and replay. Alert fires when DLQ depth &gt; 100K. Ops can replay from DLQ after fixing provider issue.">DLQ<br><small>alert + manual replay</small></div>
    </div>
  </div>
  <div class="d-cols" style="margin-top:12px">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Deduplication</div>
        <div class="d-flow-v">
          <div class="d-box purple" data-tip="notification_id is included in every provider API call as a deduplication key. FCM: collapse_key. APNs: apns-collapse-id. SendGrid: x-message-id.">notification_id as idempotency key</div>
          <div class="d-box purple" data-tip="FCM deduplicates within 14 days using collapse_key. APNs deduplicates within 30 minutes. Our DB also tracks sent status to short-circuit duplicate sends.">FCM dedup window: 14 days<br>APNs dedup window: 30 min</div>
          <div class="d-box purple" data-tip="Before sending, check Redis: SET notif:{id} EX 86400 NX. If key exists, skip — already sent. This guards against consumer restarts replaying the same message.">Redis idempotency check: SET NX</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">At-Least-Once Guarantees</div>
        <div class="d-flow-v">
          <div class="d-box green">Kafka offset committed only after provider ACK</div>
          <div class="d-box amber">Consumer crash = re-read from last committed offset</div>
          <div class="d-box amber">Max duplicate risk: 1 extra delivery per provider timeout</div>
          <div class="d-box blue">Acceptable: duplicate push notification &gt; missed OTP</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-rate-limiting",
		Title:       "Rate Limiting per User",
		Description: "Per-user per-channel rate limits using Redis token bucket with quiet hours enforcement",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Per-Channel Limits (enforced in Notification Service)</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="5 push notifications per hour per user. Protects against notification spam. OTP/critical are exempt from this limit.">Push: max <strong>5/hour</strong> per user <small>(OTP exempt)</small></div>
        <div class="d-box blue" data-tip="3 emails per hour. Marketing category is further restricted to 2/day globally. Transactional emails (receipts) exempt.">Email: max <strong>3/hour</strong> per user <small>(transactional exempt)</small></div>
        <div class="d-box amber" data-tip="1 non-OTP SMS per hour. OTP SMS always bypass rate limits — user expects to receive it. Cost driver: each SMS costs ~$0.008.">SMS: max <strong>1/hour</strong> per user <small>(OTP always exempt)</small></div>
        <div class="d-box purple" data-tip="Marketing notifications globally capped at 2/day regardless of channel. Prevents over-messaging across all channels combined.">Marketing category: max <strong>2/day</strong> across all channels</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Redis Token Bucket Implementation</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Key format: ratelimit:{user_id}:{channel}. INCR + EXPIRE via Lua script for atomicity. No race condition between check and increment.">Key: <code>rl:{user_id}:{channel}:{hour}</code></div>
        <div class="d-box amber" data-tip="Lua script: INCR key; if result == 1 then EXPIRE key 3600 end; if result &gt; limit then return 0 end; return 1. Atomic — no separate GET.">Lua atomic: INCR → check limit → set TTL on first use</div>
        <div class="d-box blue" data-tip="If rate limit exceeded, notification is either dropped (marketing) or delayed to next window (transactional). Delay uses a scheduled re-queue back to Kafka.">Exceeded: drop marketing, delay transactional</div>
        <div class="d-box purple" data-tip="Redis key uses sliding hour window: key expires in 3600s from first use. Cheap: 1 Lua call per notification, &lt;1ms P99. 1B users × 4 channels = 4B possible keys (most never set).">Memory: ~100 bytes/key × active users only</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Quiet Hours Check</div>
      <div class="d-flow-v">
        <div class="d-box gray" data-tip="Quiet hours from user_preferences table. Cached in Redis per user_id with 5-min TTL. No DB hit on every notification.">Load user quiet_hours from Redis cache</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Convert notification send time to user's timezone. Check if within quiet window. OTP/critical always bypass quiet hours.">Convert to user timezone, check window</div>
        <div class="d-arrow-down">&#8595; in quiet hours</div>
        <div class="d-box blue" data-tip="Delay push/email/marketing to after quiet_hours_end. Re-enqueue to Kafka with scheduled_at timestamp. SMS OTP still sends — user expects it.">Delay to after quiet_hours_end<br><small>OTP always bypasses</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-push-flow",
		Title:       "Push Notification Flow",
		Description: "Detailed end-to-end flow for a push notification from API call to delivery receipt",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue" data-tip="POST /notifications with event_type, user_id, template_id, variables. Returns 202 Accepted — delivery is async."><span class="d-step">1</span> API Call<br><small>POST /notifications</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green" data-tip="Look up user_preferences from Redis (TTL 5min). Check push_enabled, marketing_opt_out, category opt-out. If opted out, record and return early."><span class="d-step">2</span> Validate Preferences<br><small>Redis cache lookup</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber" data-tip="Fetch quiet_hours_start, quiet_hours_end from user prefs. Convert current time to user timezone. If in quiet hours, reschedule."><span class="d-step">3</span> Check Quiet Hours<br><small>user timezone-aware</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box red" data-tip="Redis token bucket: Lua INCR for rl:{user_id}:push:{hour}. If over limit, drop (marketing) or delay (transactional). OTP exempt."><span class="d-step">4</span> Rate Limit Check<br><small>Redis Lua INCR</small></div>
  </div>
  <div class="d-arrow-down">&#8595; passed all checks</div>
  <div class="d-flow">
    <div class="d-box purple" data-tip="Publish to notif.high or notif.critical Kafka topic depending on priority. Include notification_id, user_id, template rendered payload, device_tokens."><span class="d-step">5</span> Enqueue to Kafka<br><small>with notification_id</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green" data-tip="Consumer fetches 1000 pending notifications. Batch-renders templates with user variables. Groups by FCM project / APNs bundle ID."><span class="d-step">6</span> Consumer Batch<br><small>1000 notifications</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue" data-tip="FCM batch API: up to 500 tokens per HTTP call. APNs HTTP/2: up to 100 concurrent streams. Send with notification_id as collapse_key."><span class="d-step">7</span> FCM / APNs API<br><small>batch HTTP call</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber" data-tip="Provider returns delivery receipt with message_id. Update notification status=sent in DB. Commit Kafka offset."><span class="d-step">8</span> Delivery Receipt<br><small>commit Kafka offset</small></div>
  </div>
  <div class="d-arrow-down">&#8595; async</div>
  <div class="d-flow">
    <div class="d-box purple" data-tip="Update notifications table: status=delivered, delivered_at=now(). Used for analytics and user-facing delivery status."><span class="d-step">9</span> Update DB Status<br><small>status=delivered</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box gray" data-tip="Emit delivery_success or delivery_failure events to analytics pipeline. Track per-channel delivery rate, latency P50/P99, FCM error breakdown."><span class="d-step">10</span> Analytics Event<br><small>delivery metrics</small></div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-opt-out",
		Title:       "Opt-Out & Preferences",
		Description: "Hierarchical opt-out system from global to per-sender with CAN-SPAM compliance",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Opt-Out Hierarchy (evaluated top-down)</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="User marks account as globally opted out. No notifications of ANY type sent (except OTP which is legally required for account access). Stored in user_preferences.global_opt_out=true.">1. Global opt-out → block ALL notifications<br><small>(OTP still sends — account access)</small></div>
        <div class="d-arrow-down">&#8595; else check</div>
        <div class="d-box amber" data-tip="User disables push channel entirely. All push notifications blocked. Email/SMS still work. Stored in push_enabled=false.">2. Channel opt-out → block all [push/email/SMS]</div>
        <div class="d-arrow-down">&#8595; else check</div>
        <div class="d-box blue" data-tip="User opts out of marketing category. Social notifications still work. Stored in per-category flags or a categories_opt_out JSONB array.">3. Category opt-out → block [marketing/social/alerts]</div>
        <div class="d-arrow-down">&#8595; else check</div>
        <div class="d-box gray" data-tip="User blocks a specific sender (e.g., a business account). Stored in per_sender_blocks table: (user_id, sender_id). Low frequency, can afford DB lookup.">4. Per-sender opt-out → block specific sender</div>
        <div class="d-arrow-down">&#8595; all checks passed</div>
        <div class="d-box green">Deliver notification</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">CAN-SPAM / GDPR Compliance</div>
      <div class="d-flow-v">
        <div class="d-box purple" data-tip="One-click unsubscribe link required by CAN-SPAM. POST to /unsubscribe?token={signed_jwt} — no login required. Must be honored within 10 business days (we honor in &lt;1s).">One-click unsubscribe in every marketing email</div>
        <div class="d-box purple" data-tip="List-Unsubscribe header added to all marketing emails. Gmail/Outlook show native unsubscribe button. Format: List-Unsubscribe: &lt;mailto:unsub@domain.com&gt;, &lt;https://...&gt;">List-Unsubscribe header for email clients</div>
        <div class="d-box amber" data-tip="GDPR Right to Erasure: delete all notification history and preferences for user within 30 days of request. Cascade delete from all notification tables.">GDPR: delete notification history on erasure request</div>
        <div class="d-box amber" data-tip="Store opt-out timestamp and source (user action, email link, API). Required for compliance audits. Retained even after account deletion.">Audit log: opt-out timestamp + source stored permanently</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px">
      <div class="d-group-title">Preference Cache Strategy</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Redis key: prefs:{user_id} = serialized user_preferences. TTL 5 minutes. Cache miss triggers DB read + cache fill. On preference update, invalidate immediately.">Redis cache: prefs:{user_id} TTL=5min</div>
        <div class="d-box blue" data-tip="When user updates preferences via settings UI, immediately DELETE Redis key. Next notification will re-fetch from DB. Eventual consistency window: &lt;5 min (only if cached before invalidation).">On update: Redis DEL + DB write (synchronous)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-failure-handling",
		Title:       "Failure Scenarios",
		Description: "Circuit breaker, consumer lag scaling, and token refresh failure handling",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scenario 1: FCM Outage</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Error rate &gt; 50% over 30s window triggers circuit breaker to OPEN state. No more calls to FCM — fail immediately.">FCM error rate &gt;50% → circuit breaker OPEN</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="While circuit is open, all push notifications route to DLQ instead of FCM. Messages preserved — not dropped. DLQ depth triggers alert to on-call.">Route to DLQ — messages preserved, not dropped</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue" data-tip="After 60s, circuit enters HALF-OPEN: send 1 probe request to FCM. If success, close circuit and begin draining DLQ. If fail, reopen for another 60s.">60s probe → HALF-OPEN → drain DLQ on recovery</div>
        <div class="d-box gray">Fallback: send email if push unavailable &gt;5min</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scenario 2: Kafka Consumer Lag</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="CloudWatch metric: kafka_consumer_lag per topic. Alert threshold: lag &gt; 50K messages. Triggers auto-scaling policy on ECS consumer task group.">Kafka consumer lag &gt;50K → CloudWatch alarm</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="ECS auto-scaling adds consumer tasks. Max parallelism = number of Kafka partitions (32). Scale from 10 → 100 consumers in ~2 min.">Auto-scale: 10 → 100 consumers<br><small>max = partition count (32)</small></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue" data-tip="Consumers drain backlog. Lag metric decreases. Scale-in policy: if lag &lt; 5K for 10 min, remove 10 consumer tasks. Gradual scale-in avoids oscillation.">Drain backlog → scale-in when lag &lt;5K</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scenario 3: Expired Device Token</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="FCM returns InvalidRegistration or NotRegistered error when token is stale. APNs returns 410 Gone. These are permanent failures — no point retrying.">FCM 400 InvalidRegistration or APNs 410 Gone</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Mark device_tokens.active=false. Mobile app registers new token on next launch via POST /devices/token. Eventually consistent — token refreshed within days.">Mark device token inactive in DB</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue" data-tip="If user has no active push tokens, fall back to email or in-app notification. Priority: in-app (if user is online) → email → skip.">Fallback: try email or in-app channel</div>
        <div class="d-box gray">App re-registers token on next launch</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-template-engine",
		Title:       "Template Engine",
		Description: "Template storage, variable injection, i18n localization, and rendering pipeline",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-group">
      <div class="d-group-title">Template Storage (DynamoDB)</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="DynamoDB table: notification_templates. PK: (template_id, locale). Supports per-locale variants. Version field for rollback.">PK: (template_id, locale)</div>
        <div class="d-box blue">channel: push | email | sms | inapp</div>
        <div class="d-box blue">body: "Hi &#123;&#123;name&#125;&#125;, your order &#123;&#123;order_id&#125;&#125; is ready!"</div>
        <div class="d-box blue">variables: ["name", "order_id"]</div>
        <div class="d-box gray" data-tip="Templates cached in Redis: tmpl:{template_id}:{locale} TTL=15min. Cache invalidated on template update. Read-heavy, rarely changed.">Redis cache: tmpl:{id}:{locale} TTL=15min</div>
      </div>
    </div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-group">
      <div class="d-group-title">Rendering Pipeline</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Fetch template from Redis (cache hit 99%+). If miss, read from DynamoDB and backfill cache."><span class="d-step">1</span> Fetch template (Redis cache)</div>
        <div class="d-box green" data-tip="Load variable values from notification payload: {name: 'Alice', order_id: 'ORD-123'}. Validate all required variables are present."><span class="d-step">2</span> Load user variables from event payload</div>
        <div class="d-box amber" data-tip="Determine user locale from profile (e.g. 'fr-FR'). If template exists for that locale, use it. Else fall back to 'en' default."><span class="d-step">3</span> Determine locale (user profile → fallback 'en')</div>
        <div class="d-box amber" data-tip="Mustache/Handlebars-style substitution. Replace {{name}} → 'Alice'. Sanitize to prevent injection. Max rendered length: push=256 chars, SMS=160 chars."><span class="d-step">4</span> Render: inject variables + sanitize</div>
        <div class="d-box purple" data-tip="Apply channel-specific formatting. Push: title+body truncation. Email: HTML wrapper with header/footer, unsubscribe link. SMS: 160 char limit."><span class="d-step">5</span> Channel-specific formatting</div>
      </div>
    </div>
  </div>
  <div class="d-group" style="margin-top:12px">
    <div class="d-group-title">Template Examples</div>
    <div class="d-cols">
      <div class="d-col">
        <div class="d-box green"><strong>Push</strong><br>Title: "Order Ready 🛒"<br>Body: "Hi &#123;&#123;name&#125;&#125;, &#123;&#123;order_id&#125;&#125; is ready for pickup!"</div>
      </div>
      <div class="d-col">
        <div class="d-box blue"><strong>Email</strong><br>Subject: "Your order &#123;&#123;order_id&#125;&#125; is ready"<br>Body: Full HTML with header, body, footer, unsubscribe</div>
      </div>
      <div class="d-col">
        <div class="d-box amber"><strong>SMS</strong><br>"&#123;&#123;name&#125;&#125;: Order &#123;&#123;order_id&#125;&#125; ready. Reply STOP to unsubscribe." <small>(≤160 chars)</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "notif-monitoring",
		Title:       "Monitoring & SLAs",
		Description: "Key metrics, alert thresholds, and SLA targets for the notification system",
		ContentFile: "problems/notification-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Key Metrics (per channel)</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="delivery_rate = delivered / (delivered + permanent_failures). Temporary failures (retried) excluded. Target &gt;99% for all channels. Push typically 98-99.5% due to uninstalls.">Delivery rate: <strong>&gt;99%</strong> target per channel</div>
        <div class="d-box green" data-tip="Time from notification_created_at to provider ACK. Measured per channel. P50/P99 tracked. SLA: P99 push &lt;2s, email &lt;10s, SMS &lt;5s.">Delivery latency: P50 / P99 per channel</div>
        <div class="d-box blue" data-tip="FCM error code breakdown: InvalidRegistration (stale token), QuotaExceeded (send rate), InternalServerError (FCM outage). Track each separately.">Provider error rate by error code</div>
        <div class="d-box blue" data-tip="Kafka consumer lag per topic partition. Leading indicator of pipeline health. Spike in lag = consumer scaling event.">Kafka consumer lag per topic</div>
        <div class="d-box amber" data-tip="DLQ depth = notifications that exhausted all retries. Should be near zero. Spike indicates systemic provider issue.">DLQ depth (target: near zero)</div>
        <div class="d-box purple" data-tip="Redis hit rate for template cache and preference cache. Should be &gt;99%. Drop indicates cold start or cache eviction pressure.">Redis cache hit rate (target: &gt;99%)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Alert Thresholds</div>
      <div class="d-flow-v">
        <div class="d-box red" data-tip="Page on-call immediately. Sustained delivery failure = users missing critical notifications. Check circuit breaker state, provider status page, Kafka consumer health.">Delivery rate &lt;95% for 5min → <strong>PAGE</strong></div>
        <div class="d-box amber" data-tip="Warning alert. Investigate but not page-worthy yet. Could be a regional provider degradation or temporary spike.">Delivery rate 95–99% → WARN</div>
        <div class="d-box red" data-tip="DLQ depth growing means messages are accumulating that cannot be delivered. Likely a provider outage or systematic error. Page + investigate.">DLQ depth &gt;100K → <strong>PAGE</strong></div>
        <div class="d-box amber" data-tip="Consumer lag growing means throughput can't keep up with ingest rate. Trigger auto-scaling review. If not recovering, page.">Kafka lag &gt;50K → auto-scale + WARN</div>
        <div class="d-box amber" data-tip="P99 latency spike could indicate provider degradation or consumer bottleneck. Track per channel — push P99 &gt;5s is concerning.">Push P99 latency &gt;5s → WARN</div>
        <div class="d-box gray" data-tip="OTP delivery latency is the most user-visible metric. If an OTP takes &gt;3s, user abandons the flow. Separate SLA from general push.">OTP P99 &gt;3s → <strong>PAGE</strong> (login impact)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">SLA Targets</div>
      <div class="d-flow-v">
        <div class="d-box green">Push (non-OTP): P99 <span class="d-metric latency">&lt;2s</span></div>
        <div class="d-box red">OTP (SMS/push): P99 <span class="d-metric latency">&lt;500ms</span></div>
        <div class="d-box blue">Transactional email: P99 <span class="d-metric latency">&lt;10s</span></div>
        <div class="d-box amber">Marketing email: best-effort <span class="d-metric latency">&lt;1hr</span></div>
        <div class="d-box purple">System availability: <span class="d-metric throughput">99.9%</span></div>
        <div class="d-box gray" data-tip="At-least-once delivery guarantee. A notification may be sent twice (dedup at provider) but never silently lost.">Durability: at-least-once delivery</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

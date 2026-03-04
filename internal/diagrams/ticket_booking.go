package diagrams

func registerTicketBooking(r *Registry) {
	r.Register(&Diagram{
		Slug:        "tb-requirements",
		Title:       "Requirements & Scale Estimates",
		Description: "Non-functional requirements and scale targets for a ticket booking system like BookMyShow",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale</div>
      <div class="d-flow-v">
        <div class="d-box blue">500M registered users</div>
        <div class="d-box blue">10M bookings/day &#8776; 115 bookings/sec avg</div>
        <div class="d-box purple">100K concurrent users for popular shows (Avengers, IPL final)</div>
        <div class="d-box purple">Peak 5x &#8594; 575 bookings/sec burst</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Key Decisions</div>
      <div class="d-flow-v">
        <div class="d-box red">Seat lock duration? 10 min (payment window)</div>
        <div class="d-box red">Overselling tolerance? 0% &#8212; 99.99% no double-booking</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">NFR Targets</div>
      <div class="d-flow-v">
        <div class="d-box green">Booking flow: &lt; 2s end-to-end</div>
        <div class="d-box green">Seat map load: &lt; 500ms p99</div>
        <div class="d-box green">Availability: 99.99% (52 min downtime/yr)</div>
        <div class="d-box amber">Double-booking: 99.99% prevention guarantee</div>
        <div class="d-box amber">Consistency: Strong for seat locks, eventual for seat map views</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-api-design",
		Title:       "API Design",
		Description: "Core API endpoints for browsing shows, locking seats, confirming bookings, and retrieving tickets",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P0 &#8212; Core (Must Have)</div>
      <div class="d-flow-v">
        <div class="d-box green">GET /api/v1/shows?city=&#38;date= &#8594; list shows with availability</div>
        <div class="d-box green">GET /api/v1/shows/{id}/seats &#8594; real-time seat map (available/locked/booked)</div>
        <div class="d-box green">POST /api/v1/bookings/lock &#8594; {show_id, seat_ids[]} &#8594; lock_token, expires_at</div>
        <div class="d-box green">POST /api/v1/bookings/confirm &#8594; {lock_token, payment_id} &#8594; booking_id</div>
        <div class="d-box green">GET /api/v1/bookings/{id}/ticket &#8594; QR code + booking details</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P1 &#8212; Important</div>
      <div class="d-flow-v">
        <div class="d-box blue">POST /api/v1/payments/initiate &#8594; redirect URL or payment intent</div>
        <div class="d-box blue">POST /api/v1/payments/webhook &#8594; payment gateway callback</div>
        <div class="d-box blue">GET /api/v1/users/{id}/bookings &#8594; booking history</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P2 &#8212; Nice to Have</div>
      <div class="d-flow-v">
        <div class="d-box gray">POST /api/v1/reviews/{show_id} &#8594; rating + review</div>
        <div class="d-box gray">GET /api/v1/recommendations &#8594; personalized show suggestions</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-seat-locking",
		Title:       "Seat Locking Flow",
		Description: "Distributed seat locking with Redis SETNX, 10-minute TTL, payment window, and timeout release",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">User selects seats</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">POST /bookings/lock</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple">Redis SETNX lock:show:{id}:seat:{num}</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Happy Path (within 10 min TTL)</div>
        <div class="d-flow-v">
          <div class="d-box green">Lock acquired &#8594; return lock_token</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">User completes payment</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">POST /bookings/confirm</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Write booking to Postgres (SERIALIZABLE)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Generate QR ticket</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Timeout / Failure Path</div>
        <div class="d-flow-v">
          <div class="d-box red">10 min TTL expires</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Redis key auto-deletes</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Seat becomes available again</div>
          <div class="d-label">No manual cleanup needed</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Contention Path</div>
        <div class="d-flow-v">
          <div class="d-box red">SETNX returns 0 (already locked)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Return 409 Conflict &#8212; seat unavailable</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end architecture: CDN, ALB, API services, Redis seat locks, Postgres bookings, payment and QR services",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Client (Web / iOS / Android)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">CloudFront CDN (static assets, seat map cache)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">ALB (Application Load Balancer)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">API Layer</div>
        <div class="d-flow-v">
          <div class="d-box green">Show Service (browse, search)</div>
          <div class="d-box green">Booking Service (lock, confirm, cancel)</div>
          <div class="d-box green">Ticket Service (QR generation, validation)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Data Layer</div>
        <div class="d-flow-v">
          <div class="d-box red">Redis Cluster (seat locks, TTL 10min)</div>
          <div class="d-box indigo">Postgres Primary (bookings, payments)</div>
          <div class="d-box indigo">Postgres Read Replicas (show listings, seat views)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">External Services</div>
        <div class="d-flow-v">
          <div class="d-box amber">Payment Gateway (Stripe / Razorpay)</div>
          <div class="d-box amber">S3 (QR codes, ticket PDFs)</div>
          <div class="d-box amber">SQS (async: notifications, analytics)</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-data-model",
		Title:       "Data Model",
		Description: "Entity relationship diagram: theatres, screens, shows, seats, bookings, payments, tickets",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">theatres</div>
      <div class="d-entity-body">
        <div class="pk">id UUID (PK)</div>
        <div>name VARCHAR(255)</div>
        <div>city VARCHAR(100)</div>
        <div>address TEXT</div>
        <div class="idx idx-btree">city (index for city search)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header blue">screens</div>
      <div class="d-entity-body">
        <div class="pk">id UUID (PK)</div>
        <div class="fk">theatre_id UUID (FK &#8594; theatres)</div>
        <div>name VARCHAR(50)</div>
        <div>total_seats INT</div>
        <div>layout JSONB (row/col config)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header green">shows</div>
      <div class="d-entity-body">
        <div class="pk">id UUID (PK)</div>
        <div class="fk">screen_id UUID (FK &#8594; screens)</div>
        <div>movie_id UUID</div>
        <div>start_time TIMESTAMPTZ</div>
        <div>end_time TIMESTAMPTZ</div>
        <div>status ENUM (scheduled, cancelled, completed)</div>
        <div class="idx idx-btree">screen_id, start_time (composite)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header green">seats</div>
      <div class="d-entity-body">
        <div class="pk">id UUID (PK)</div>
        <div class="fk">screen_id UUID (FK &#8594; screens)</div>
        <div>row CHAR(1)</div>
        <div>number INT</div>
        <div>category ENUM (silver, gold, platinum)</div>
        <div>price_cents INT</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header indigo">bookings</div>
      <div class="d-entity-body">
        <div class="pk">id UUID (PK)</div>
        <div class="fk">user_id UUID (FK)</div>
        <div class="fk">show_id UUID (FK &#8594; shows)</div>
        <div>seat_ids UUID[] (array of seat IDs)</div>
        <div>status ENUM (locked, confirmed, cancelled, expired)</div>
        <div>lock_token UUID (unique)</div>
        <div>locked_at TIMESTAMPTZ</div>
        <div>confirmed_at TIMESTAMPTZ</div>
        <div>total_cents INT</div>
        <div class="idx idx-btree">user_id (booking history)</div>
        <div class="idx idx-btree">show_id, status (availability queries)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header amber">payments</div>
      <div class="d-entity-body">
        <div class="pk">id UUID (PK)</div>
        <div class="fk">booking_id UUID (FK &#8594; bookings)</div>
        <div>gateway_txn_id VARCHAR(255)</div>
        <div>amount_cents INT</div>
        <div>status ENUM (initiated, success, failed, refunded)</div>
        <div>provider ENUM (stripe, razorpay)</div>
        <div>created_at TIMESTAMPTZ</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header purple">tickets</div>
      <div class="d-entity-body">
        <div class="pk">id UUID (PK)</div>
        <div class="fk">booking_id UUID (FK &#8594; bookings)</div>
        <div>qr_token VARCHAR(64) UNIQUE</div>
        <div>qr_image_url TEXT (S3 path)</div>
        <div>scanned_at TIMESTAMPTZ (nullable)</div>
        <div>status ENUM (active, used, cancelled)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-rush-handling",
		Title:       "Rush Scenario (100K&#8594;500 seats)",
		Description: "How 100K concurrent users compete for 500 seats: queue, rate limit, lock acquisition pipeline",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box red">100K users hit &#8220;Book Now&#8221;</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple">CAPTCHA + Device Fingerprint</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber">Virtual Queue (SQS FIFO)</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box blue">Rate limiter: 1K users/sec dequeued</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">Redis SETNX seat lock attempt</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Lock Acquired (~500 users)</div>
        <div class="d-flow-v">
          <div class="d-box green">10 min payment window</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Payment &#8594; Confirm &#8594; Ticket</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Lock Failed (~99.5K users)</div>
        <div class="d-flow-v">
          <div class="d-box red">409 Conflict &#8212; seat taken</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Show waitlist option / other showtimes</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Timeout Recycle</div>
        <div class="d-flow-v">
          <div class="d-box amber">~15% abandon payment (75 seats)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">TTL expires &#8594; seats re-enter pool</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Waitlisted users notified via push</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-payment-flow",
		Title:       "Payment & Failure Handling",
		Description: "Payment lifecycle: lock seat, initiate payment, gateway callback, confirm or rollback",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box green">Seat Locked (lock_token)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">POST /payments/initiate</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber">Payment Gateway (Stripe/Razorpay)</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Success Path</div>
        <div class="d-flow-v">
          <div class="d-box green">Webhook: payment.success</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Verify lock_token still valid</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">BEGIN TX: INSERT booking (confirmed), INSERT payment</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Delete Redis lock (seat now permanently booked)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Async: generate QR, send confirmation email</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Failure Path</div>
        <div class="d-flow-v">
          <div class="d-box red">Webhook: payment.failed</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Delete Redis lock immediately</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Seat released back to pool</div>
          <div class="d-label">User can retry with fresh lock</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Edge Case: Payment succeeds after lock TTL expires</div>
        <div class="d-flow-v">
          <div class="d-box red">Lock expired + payment.success received</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Check if seat re-locked by another user</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red">Auto-refund via gateway API</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-real-time-seats",
		Title:       "Real-Time Seat Availability",
		Description: "WebSocket fan-out for live seat map updates using Redis Pub/Sub",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box green">Seat lock/unlock event</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box red">Redis PUBLISH seat:show:{id}</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box red">Redis Pub/Sub channel: seat:show:{id}</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">WebSocket Server 1</div>
        <div class="d-flow-v">
          <div class="d-box blue">SUBSCRIBE seat:show:{id}</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Client A (seat map open)</div>
          <div class="d-box gray">Client B (seat map open)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">WebSocket Server 2</div>
        <div class="d-flow-v">
          <div class="d-box blue">SUBSCRIBE seat:show:{id}</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray">Client C (seat map open)</div>
          <div class="d-box gray">Client D (seat map open)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Message Payload</div>
        <div class="d-flow-v">
          <div class="d-box purple">{show_id, seat_id, status: "locked"/"available", user_hint: "2 seats left in Row A"}</div>
          <div class="d-label">~50 bytes per update, &lt; 100ms fan-out</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-scaling",
		Title:       "Scaling Strategy",
		Description: "Horizontal scaling: read replicas for browsing, Redis cluster for locks, queue for payment processing",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Read Path (Show Browsing)</div>
      <div class="d-flow-v">
        <div class="d-box purple">CloudFront CDN (show listings, images)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Postgres Read Replicas (3x)</div>
        <div class="d-label">Seat availability cached in Redis (5s TTL)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">ElastiCache Redis (seat map cache)</div>
        <div class="d-label">90% cache hit rate for seat views</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Write Path (Booking)</div>
      <div class="d-flow-v">
        <div class="d-box red">Redis Cluster (6 shards, seat locks)</div>
        <div class="d-label">Shard by show_id &#8594; all seats for one show on same shard</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo">Postgres Primary (single writer)</div>
        <div class="d-label">Partitioned by show_date (monthly)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">SQS FIFO (payment processing queue)</div>
        <div class="d-label">Deduplication by lock_token</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Operational</div>
      <div class="d-flow-v">
        <div class="d-box green">Auto-scaling: ECS tasks scale on CPU/request count</div>
        <div class="d-box blue">Redis: cluster mode, 3 primaries + 3 replicas</div>
        <div class="d-box indigo">Postgres: r6g.2xlarge, 500GB GP3, daily snapshots</div>
        <div class="d-box amber">Cost: ~$3,200/mo (Redis $600 + RDS $1,400 + ECS $800 + misc $400)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-qr-generation",
		Title:       "QR Code Ticket Generation",
		Description: "Async flow: booking confirmed, generate unique token, encode QR, store in S3, push to user",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box green">Booking Confirmed</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber">SQS: ticket-generation queue</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">Ticket Worker (Lambda / ECS)</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Token Generation</div>
        <div class="d-flow-v">
          <div class="d-box purple">Generate: HMAC-SHA256(booking_id + secret + timestamp)</div>
          <div class="d-label">64-char hex token, unique per booking</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Encode QR: token + show_id + seat_info</div>
          <div class="d-label">QR payload ~200 bytes, Version 6 (41&#215;41)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Storage &#38; Delivery</div>
        <div class="d-flow-v">
          <div class="d-box amber">Upload QR PNG to S3</div>
          <div class="d-label">s3://tickets/{booking_id}/qr.png</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Update tickets table: qr_image_url, qr_token</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Push notification + email with ticket link</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Validation at Venue</div>
    <div class="d-flow">
      <div class="d-box gray">Scanner reads QR</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">GET /tickets/validate?token=...</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Verify HMAC + check scanned_at IS NULL</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Mark scanned_at = NOW()</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-anti-scalping",
		Title:       "Anti-Scalping & Bot Prevention",
		Description: "Multi-layer defense: CAPTCHA, rate limiting, device fingerprinting, purchase history, IP reputation",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Incoming booking request</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Layer 1: Edge (CloudFront + WAF)</div>
    <div class="d-flow">
      <div class="d-box purple">AWS WAF rate limit: 100 req/min per IP</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Geo-blocking (restrict to service regions)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Known bot IP blocklist (Spamhaus, AbuseIPDB)</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Layer 2: Application</div>
    <div class="d-flow">
      <div class="d-box amber">CAPTCHA (reCAPTCHA v3, score &gt; 0.7)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Device fingerprint (FingerprintJS)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Session velocity: max 3 lock attempts/min</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Layer 3: Business Logic</div>
    <div class="d-flow">
      <div class="d-box red">Max 6 tickets per user per show</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">Purchase history: flag accounts with &gt; 20 bookings/month</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">Phone verification for high-demand shows</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box green">Legitimate user proceeds to seat selection</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "tb-nfr-estimates",
		Title:       "Non-Functional Requirements",
		Description: "Detailed latency, availability, throughput, and consistency targets",
		ContentFile: "problems/ticket-booking",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Latency Targets</div>
      <div class="d-flow-v">
        <div class="d-box green">Show listing: &lt; 200ms p99 (CDN-cached)</div>
        <div class="d-box green">Seat map load: &lt; 500ms p99 (Redis + read replica)</div>
        <div class="d-box green">Seat lock: &lt; 100ms p99 (Redis SETNX)</div>
        <div class="d-box blue">Booking confirm: &lt; 2s p99 (DB write + QR async)</div>
        <div class="d-box blue">QR generation: &lt; 5s async (SQS + Lambda)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Availability</div>
      <div class="d-flow-v">
        <div class="d-box green">Overall: 99.99% (52 min downtime/yr)</div>
        <div class="d-box amber">Booking path: 99.999% (Redis cluster + Postgres HA)</div>
        <div class="d-box gray">Browse path: 99.9% acceptable (CDN fallback)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Throughput</div>
      <div class="d-flow-v">
        <div class="d-box purple">Read: 50K seat map views/sec (peak rush)</div>
        <div class="d-box purple">Write: 575 bookings/sec (10M/day, 5x peak)</div>
        <div class="d-box purple">Lock ops: 10K SETNX/sec (Redis cluster handles 200K+)</div>
        <div class="d-box amber">Payment: 200 TPS (gateway limit, queued beyond)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Consistency &amp; Durability</div>
      <div class="d-flow-v">
        <div class="d-box red">Seat locks: Strong consistency (Redis single-key atomic)</div>
        <div class="d-box red">Bookings: Serializable isolation (Postgres)</div>
        <div class="d-box amber">Seat map view: Eventual (5s cache TTL acceptable)</div>
        <div class="d-box green">Zero data loss: WAL + synchronous replication for bookings</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

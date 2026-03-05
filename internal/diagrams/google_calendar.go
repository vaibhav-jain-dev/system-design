package diagrams

func registerGoogleCalendar(r *Registry) {
	r.Register(&Diagram{
		Slug:        "gc-requirements",
		Title:       "Scale & Requirements",
		Description: "Scale targets, functional and non-functional requirements for Google Calendar",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale</div>
      <div class="d-flow-v">
        <div class="d-box blue">1B registered users</div>
        <div class="d-box blue">500M monthly active users</div>
        <div class="d-box blue">50B total events stored</div>
        <div class="d-box purple">100K event reads/sec (peak)</div>
        <div class="d-box purple">10K event writes/sec (peak)</div>
        <div class="d-box purple">&#8776; 50 events/user/week average</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P0 &#8212; Core (Must Have)</div>
      <div class="d-flow-v">
        <div class="d-box green">Create / update / delete events</div>
        <div class="d-box green">Recurring events (RRULE)</div>
        <div class="d-box green">Multi-timezone support (IANA)</div>
        <div class="d-box green">Free/busy conflict detection</div>
        <div class="d-box green">Invite attendees + RSVP</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">P1 &#8212; Important</div>
      <div class="d-flow-v">
        <div class="d-box blue">Calendar sharing &amp; permissions (ACL)</div>
        <div class="d-box blue">Multi-channel notifications (email, push, in-app)</div>
        <div class="d-box blue">Incremental sync (mobile/web clients)</div>
        <div class="d-box blue">Day / week / month view rendering</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">P2 &#8212; Nice to Have</div>
      <div class="d-flow-v">
        <div class="d-box gray">Smart scheduling (find optimal slot)</div>
        <div class="d-box gray">Room / resource booking</div>
        <div class="d-box gray">CalDAV / iCal external sync</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">Read latency: &lt; 100ms (p99)</div>
        <div class="d-box purple">Write latency: &lt; 200ms (p99)</div>
        <div class="d-box purple">Availability: 99.99%</div>
        <div class="d-box amber">Timezone correctness: zero DST bugs</div>
        <div class="d-box amber">Sync consistency: eventual (&lt; 2s lag)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-api-design",
		Title:       "API Design",
		Description: "Core REST API endpoints for calendar operations, recurring events, sharing, and notifications",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Event CRUD</div>
      <div class="d-flow-v">
        <div class="d-box green">POST /calendars/{calId}/events &#8212; Create event</div>
        <div class="d-box green">GET /calendars/{calId}/events/{id} &#8212; Get event</div>
        <div class="d-box green">PUT /calendars/{calId}/events/{id} &#8212; Update event</div>
        <div class="d-box green">DELETE /calendars/{calId}/events/{id} &#8212; Delete event</div>
        <div class="d-box blue">GET /calendars/{calId}/events?start=&amp;end= &#8212; List range</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Recurring Events</div>
      <div class="d-flow-v">
        <div class="d-box purple">POST .../events (body: rrule=FREQ=WEEKLY;BYDAY=MO,WE)</div>
        <div class="d-box purple">PUT .../events/{id}/instances/{instanceDate} &#8212; Edit single</div>
        <div class="d-box amber">DELETE .../events/{id}?scope=this|following|all</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Free/Busy</div>
      <div class="d-flow-v">
        <div class="d-box indigo">POST /freeBusy &#8212; Body: {users: [...], timeMin, timeMax}</div>
        <div class="d-box indigo">Response: [{user, busy: [{start, end}, ...]}]</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Calendar Sharing</div>
      <div class="d-flow-v">
        <div class="d-box blue">POST /calendars/{calId}/acl &#8212; Add permission</div>
        <div class="d-box blue">PUT /calendars/{calId}/acl/{ruleId} &#8212; Update role</div>
        <div class="d-box blue">DELETE /calendars/{calId}/acl/{ruleId} &#8212; Revoke</div>
        <div class="d-box gray">Roles: owner | writer | reader | freeBusyOnly</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Notifications</div>
      <div class="d-flow-v">
        <div class="d-box amber">PUT /calendars/{calId}/events/{id}/reminders</div>
        <div class="d-box amber">Body: [{method: email|push|popup, minutes: 15}]</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Sync</div>
      <div class="d-flow-v">
        <div class="d-box green">GET /calendars/{calId}/events?syncToken={token}</div>
        <div class="d-box green">Response: {items: [...], nextSyncToken: "..."}</div>
        <div class="d-box gray">410 Gone &#8594; full sync required (token expired)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-recurring-events",
		Title:       "Recurring Events &#8212; RRULE Expansion",
		Description: "How recurring events are stored, expanded, and how exceptions override individual instances",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Master Event<br/>rrule: FREQ=WEEKLY;BYDAY=MO,WE<br/>dtstart: 2024-01-01T09:00Z<br/>until: 2024-12-31</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">RRULE Engine<br/>Expand on read<br/>(virtual instances)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-group">
      <div class="d-group-title">Generated Instances</div>
      <div class="d-flow-v">
        <div class="d-box gray">Jan 1 (Mon) 9:00</div>
        <div class="d-box gray">Jan 3 (Wed) 9:00</div>
        <div class="d-box gray">Jan 8 (Mon) 9:00</div>
        <div class="d-box gray">... &#8776; 104 instances/year</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Exception Instances (materialized)</div>
        <div class="d-flow-v">
          <div class="d-box amber">Exception: Jan 8 moved to 10:00<br/>recurring_event_id = master.id<br/>original_start = Jan 8 09:00<br/>new_start = Jan 8 10:00</div>
          <div class="d-box red">Exception: Jan 15 cancelled<br/>status = CANCELLED</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Expansion Strategy</div>
        <div class="d-flow-v">
          <div class="d-box purple">Virtual: expand on query &#8212; O(1) storage</div>
          <div class="d-box purple">Materialized: pre-generate N days ahead</div>
          <div class="d-box indigo">Hybrid: virtual + materialize on edit</div>
          <div class="d-box green">Our choice: Hybrid<br/>&#8212; Virtual until user edits an instance<br/>&#8212; Exception row created on modification<br/>&#8212; Query merges master RRULE + exceptions</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end architecture: CDN, ALB, microservices, data stores, and async processing",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box gray">Web / Mobile Clients</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">CDN (CloudFront)<br/>Static assets, calendar UI</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">ALB<br/>TLS termination<br/>path-based routing</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-group">
      <div class="d-group-title">API Gateway</div>
      <div class="d-flow-v">
        <div class="d-box green">Auth / Rate Limit</div>
        <div class="d-box green">Request Routing</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box green">Event Service<br/>CRUD, RRULE expansion<br/>conflict detection</div>
    <div class="d-box green">Notification Service<br/>Reminders, invites<br/>multi-channel delivery</div>
    <div class="d-box green">Sync Service<br/>Delta sync tokens<br/>change tracking</div>
    <div class="d-box green">Sharing Service<br/>ACL management<br/>permission checks</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box indigo">PostgreSQL<br/>Events, calendars,<br/>users, ACLs<br/>(sharded by user_id)</div>
    <div class="d-box purple">Redis<br/>Session cache, free/busy<br/>cache, sync tokens<br/>(ElastiCache cluster)</div>
    <div class="d-box amber">Kafka<br/>Event changes &#8594; notifications<br/>Event changes &#8594; sync fanout<br/>Event changes &#8594; analytics</div>
    <div class="d-box gray">S3<br/>Attachments,<br/>ICS exports,<br/>backups</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-conflict-detection",
		Title:       "Free/Busy Conflict Detection",
		Description: "Interval-based conflict detection with O(log n) queries using sorted intervals and batch availability",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Conflict Detection Flow</div>
      <div class="d-flow-v">
        <div class="d-box blue">New Event: 2pm&#8211;3pm</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Query: SELECT * FROM events<br/>WHERE user_id = ? AND calendar_id = ?<br/>AND start_time &lt; '3pm'<br/>AND end_time &gt; '2pm'<br/>AND status != 'cancelled'</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">B-tree index on (user_id, start_time)<br/>&#8594; O(log n) lookup</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-flow">
          <div class="d-box green">No overlap &#8594; create event</div>
          <div class="d-box amber">Overlap found &#8594; warn user</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Free/Busy Cache (Redis)</div>
      <div class="d-flow-v">
        <div class="d-box purple">Key: freebusy:{user_id}:{date}</div>
        <div class="d-box purple">Value: sorted set of intervals<br/>ZADD score=start_epoch member="start|end"</div>
        <div class="d-box indigo">Batch query: ZRANGEBYSCORE<br/>for all attendees in parallel<br/>&#8594; intersect busy slots</div>
        <div class="d-box green">TTL: 1 hour, invalidate on event change</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Batch Availability (Find a Time)</div>
      <div class="d-flow-v">
        <div class="d-box blue">Input: 5 attendees, next 7 days</div>
        <div class="d-box blue">Merge all busy intervals &#8594; union</div>
        <div class="d-box green">Find gaps &#8805; requested duration</div>
        <div class="d-box green">Return ranked slots (fewest conflicts)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-data-model",
		Title:       "Data Model",
		Description: "Core database tables: users, calendars, events, recurring rules, attendees, notifications, and sharing permissions",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">users</div>
      <div class="d-entity-body">
        <div class="pk">user_id UUID (PK)</div>
        <div>email VARCHAR(255)</div>
        <div>display_name VARCHAR(100)</div>
        <div>default_timezone VARCHAR(50)</div>
        <div>notification_prefs JSONB</div>
        <div class="idx idx-btree">idx_email UNIQUE</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header green">calendars</div>
      <div class="d-entity-body">
        <div class="pk">calendar_id UUID (PK)</div>
        <div class="fk">owner_id UUID (FK &#8594; users)</div>
        <div>title VARCHAR(200)</div>
        <div>color VARCHAR(7)</div>
        <div>timezone VARCHAR(50)</div>
        <div>is_primary BOOLEAN</div>
        <div class="idx idx-btree">idx_owner_id</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header indigo">events</div>
      <div class="d-entity-body">
        <div class="pk">event_id UUID (PK)</div>
        <div class="fk">calendar_id UUID (FK &#8594; calendars)</div>
        <div class="fk">creator_id UUID (FK &#8594; users)</div>
        <div>title VARCHAR(500)</div>
        <div>description TEXT</div>
        <div>start_time TIMESTAMPTZ</div>
        <div>end_time TIMESTAMPTZ</div>
        <div>timezone VARCHAR(50)</div>
        <div>is_all_day BOOLEAN</div>
        <div>status ENUM(confirmed, tentative, cancelled)</div>
        <div class="fk">recurring_rule_id UUID (FK, nullable)</div>
        <div>original_start TIMESTAMPTZ (for exceptions)</div>
        <div>updated_at TIMESTAMPTZ</div>
        <div>etag VARCHAR(64)</div>
        <div class="idx idx-btree">idx_cal_start (calendar_id, start_time)</div>
        <div class="idx idx-btree">idx_cal_end (calendar_id, end_time)</div>
        <div class="idx idx-btree">idx_recurring_rule_id</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header purple">recurring_rules</div>
      <div class="d-entity-body">
        <div class="pk">rule_id UUID (PK)</div>
        <div>rrule VARCHAR(500) (RFC 5545)</div>
        <div>dtstart TIMESTAMPTZ</div>
        <div>dtend TIMESTAMPTZ (series end)</div>
        <div>timezone VARCHAR(50)</div>
        <div>exdates TIMESTAMPTZ[] (excluded dates)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header amber">attendees</div>
      <div class="d-entity-body">
        <div class="pk">attendee_id UUID (PK)</div>
        <div class="fk">event_id UUID (FK &#8594; events)</div>
        <div class="fk">user_id UUID (FK &#8594; users)</div>
        <div>email VARCHAR(255)</div>
        <div>response ENUM(accepted, declined, tentative, needsAction)</div>
        <div>role ENUM(organizer, required, optional)</div>
        <div class="idx idx-btree">idx_event_id</div>
        <div class="idx idx-btree">idx_user_id_event_id UNIQUE</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header red">notifications</div>
      <div class="d-entity-body">
        <div class="pk">notification_id UUID (PK)</div>
        <div class="fk">event_id UUID (FK &#8594; events)</div>
        <div class="fk">user_id UUID (FK &#8594; users)</div>
        <div>method ENUM(email, push, popup)</div>
        <div>trigger_minutes INT (before event)</div>
        <div>scheduled_at TIMESTAMPTZ</div>
        <div>sent BOOLEAN DEFAULT false</div>
        <div class="idx idx-btree">idx_scheduled_unsent (scheduled_at) WHERE NOT sent</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header gray">sharing_permissions</div>
      <div class="d-entity-body">
        <div class="pk">permission_id UUID (PK)</div>
        <div class="fk">calendar_id UUID (FK &#8594; calendars)</div>
        <div class="fk">grantee_id UUID (FK &#8594; users, nullable)</div>
        <div>grantee_email VARCHAR(255)</div>
        <div>role ENUM(owner, writer, reader, freeBusyOnly)</div>
        <div>scope ENUM(user, group, domain, public)</div>
        <div class="idx idx-btree">idx_calendar_grantee UNIQUE</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-notification-system",
		Title:       "Notification System",
		Description: "Multi-channel notification delivery: email via SES, push via FCM/APNs, in-app, and reminder scheduling",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Event Created / Updated</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">Kafka: calendar-events topic</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">Notification Service</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Reminder Scheduler</div>
        <div class="d-flow-v">
          <div class="d-box purple">Scan notifications table<br/>WHERE scheduled_at &lt;= NOW()<br/>AND sent = false</div>
          <div class="d-box purple">Redis sorted set as delay queue<br/>ZADD reminders {fire_time} {notification_id}</div>
          <div class="d-box purple">Worker: ZPOPMIN every 1s<br/>&#8594; dispatch to channel router</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Channel Router</div>
        <div class="d-flow-v">
          <div class="d-flow">
            <div class="d-box amber">Email</div>
            <div class="d-arrow">&#8594;</div>
            <div class="d-box gray">AWS SES<br/>&#8776; $0.10 per 1K emails</div>
          </div>
          <div class="d-flow">
            <div class="d-box amber">Push</div>
            <div class="d-arrow">&#8594;</div>
            <div class="d-box gray">FCM (Android)<br/>APNs (iOS)</div>
          </div>
          <div class="d-flow">
            <div class="d-box amber">In-App</div>
            <div class="d-arrow">&#8594;</div>
            <div class="d-box gray">WebSocket / SSE<br/>to connected clients</div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Notification Types</div>
    <div class="d-flow">
      <div class="d-box green">Event invite</div>
      <div class="d-box green">RSVP update</div>
      <div class="d-box blue">Event modified</div>
      <div class="d-box blue">Event cancelled</div>
      <div class="d-box purple">Reminder (15m, 1h, 1d before)</div>
      <div class="d-box gray">Daily agenda summary</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-timezone-handling",
		Title:       "Timezone Handling",
		Description: "UTC storage, IANA timezone database, DST edge cases, and floating time events",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Storage Strategy</div>
      <div class="d-flow-v">
        <div class="d-box blue">Store all times as TIMESTAMPTZ (UTC)</div>
        <div class="d-box blue">Store original IANA timezone alongside<br/>e.g., "America/New_York"</div>
        <div class="d-box green">Convert to user&#8217;s local TZ on read</div>
        <div class="d-box green">Never store offsets (UTC-5) &#8212; they change with DST</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">DST Edge Cases</div>
      <div class="d-flow-v">
        <div class="d-box amber">Spring forward: 2:30 AM doesn&#8217;t exist<br/>&#8594; Snap to 3:00 AM (next valid time)</div>
        <div class="d-box amber">Fall back: 1:30 AM occurs twice<br/>&#8594; Use wall clock + timezone = unambiguous</div>
        <div class="d-box red">Recurring at 2:30 AM across DST boundary<br/>&#8594; Expand using wall-clock time in IANA TZ<br/>&#8594; Library handles non-existent times</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">IANA Timezone Database</div>
      <div class="d-flow-v">
        <div class="d-box indigo">&#8776; 600 timezone identifiers</div>
        <div class="d-box indigo">Updated 3&#8211;4 times/year (governments change rules)</div>
        <div class="d-box indigo">Ship TZDB with application, not OS-level</div>
        <div class="d-box purple">Use: Go time.LoadLocation() / Java ZoneId.of()</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Floating Time Events</div>
      <div class="d-flow-v">
        <div class="d-box green">All-day events: DATE only, no timezone<br/>e.g., "Birthday on March 15"</div>
        <div class="d-box green">Stored as DATE type, rendered in viewer&#8217;s TZ</div>
        <div class="d-box gray">Traveling user: "Lunch at noon" = local time<br/>&#8594; Store as floating (no TZ conversion)</div>
        <div class="d-box gray">Flag: is_floating BOOLEAN on event</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-sync-protocol",
		Title:       "Incremental Sync Protocol",
		Description: "Sync tokens, delta changes, and conflict resolution for mobile and web clients",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Sync Flow</div>
    <div class="d-flow">
      <div class="d-box blue">Client<br/>syncToken: "abc123"</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">GET /events?syncToken=abc123</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Sync Service</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Query changes since token<br/>SELECT * FROM events<br/>WHERE updated_at &gt; token_timestamp<br/>AND calendar_id IN (accessible)</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Sync Token Implementation</div>
        <div class="d-flow-v">
          <div class="d-box indigo">Token = Base64(user_id + timestamp + page_cursor)</div>
          <div class="d-box indigo">Stored in Redis: sync:{user}:{calendar} &#8594; last_sync_ts</div>
          <div class="d-box amber">Token expiry: 30 days<br/>Expired &#8594; HTTP 410 Gone &#8594; full sync</div>
          <div class="d-box green">Response includes nextSyncToken for next call</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Delta Response</div>
        <div class="d-flow-v">
          <div class="d-box green">Created events: full event object</div>
          <div class="d-box blue">Updated events: full event object (with new etag)</div>
          <div class="d-box red">Deleted events: {id, status: "cancelled"}</div>
          <div class="d-box gray">Page size: max 250 events per response</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Conflict Resolution</div>
        <div class="d-flow-v">
          <div class="d-box purple">Optimistic concurrency: etag on every event</div>
          <div class="d-box purple">PUT with If-Match: {etag}<br/>&#8594; 412 Precondition Failed if stale</div>
          <div class="d-box amber">Server-side: last-writer-wins (by updated_at)</div>
          <div class="d-box amber">Client-side: show conflict UI for concurrent edits</div>
        </div>
      </div>
      <div class="d-group">
        <div class="d-group-title">Push Notifications for Sync</div>
        <div class="d-flow-v">
          <div class="d-box green">Kafka event &#8594; WebSocket push to connected clients</div>
          <div class="d-box green">Mobile: silent push &#8594; triggers background sync</div>
          <div class="d-box gray">Fallback: client polls every 5 minutes</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-scaling",
		Title:       "Scaling Strategy",
		Description: "Sharding by user_id, read replicas for calendar views, and denormalized day-view cache",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Database Sharding</div>
      <div class="d-flow-v">
        <div class="d-box blue">Shard key: user_id (hash-based)</div>
        <div class="d-box blue">All user data co-located on same shard<br/>&#8212; calendars, events, attendees, notifications</div>
        <div class="d-box amber">Cross-user queries (shared calendars)<br/>&#8594; scatter-gather across shards</div>
        <div class="d-box green">&#8776; 16 shards initial, split at 500GB/shard</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Read Replicas</div>
      <div class="d-flow-v">
        <div class="d-box purple">2 read replicas per shard</div>
        <div class="d-box purple">Calendar view reads &#8594; replicas (100K QPS)</div>
        <div class="d-box purple">Writes &#8594; primary only (10K QPS)</div>
        <div class="d-box gray">Replication lag: &lt; 100ms (acceptable for views)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Day-View Cache (Redis)</div>
      <div class="d-flow-v">
        <div class="d-box green">Key: dayview:{user_id}:{date}</div>
        <div class="d-box green">Value: pre-rendered event list for that day</div>
        <div class="d-box green">TTL: 6 hours, invalidate on event change</div>
        <div class="d-box indigo">Cache hit rate: &#8776; 85%<br/>Most users view same day repeatedly</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Hot User Mitigation</div>
      <div class="d-flow-v">
        <div class="d-box amber">CEO calendar: 500+ attendees per event<br/>&#8594; Fan-out writes to attendee shards</div>
        <div class="d-box amber">Rate limit: max 2500 attendees/event</div>
        <div class="d-box red">Organization-wide events (10K+ people)<br/>&#8594; Async notification, no inline fan-out</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Capacity Estimates</div>
      <div class="d-flow-v">
        <div class="d-box gray">50B events &#215; 1KB avg = 50TB raw storage</div>
        <div class="d-box gray">With indexes + replicas &#8776; 200TB total</div>
        <div class="d-box gray">Redis cache: &#8776; 500GB (hot data + free/busy)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-sharing-permissions",
		Title:       "Sharing &amp; Permissions (ACL Model)",
		Description: "Access control list model: owner, editor, viewer roles, organization-wide sharing, and delegation",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Permission Hierarchy</div>
      <div class="d-flow-v">
        <div class="d-box green">Owner<br/>Full control: CRUD events, manage sharing, delete calendar</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Writer (Editor)<br/>Create, edit, delete events. Cannot manage sharing.</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Reader (Viewer)<br/>See event details. Cannot modify.</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray">FreeBusy Only<br/>See busy/free slots. No event details.</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Permission Check Flow</div>
      <div class="d-flow-v">
        <div class="d-box indigo">1. Check user &#8594; calendar direct ACL</div>
        <div class="d-box indigo">2. Check user&#8217;s groups &#8594; calendar group ACL</div>
        <div class="d-box indigo">3. Check domain-wide default (org setting)</div>
        <div class="d-box indigo">4. Check public access flag</div>
        <div class="d-box green">First match wins (most specific scope)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Sharing Scopes</div>
      <div class="d-flow-v">
        <div class="d-box blue">User: share with specific email</div>
        <div class="d-box blue">Group: share with Google Group / team</div>
        <div class="d-box amber">Domain: everyone@company.com &#8594; freeBusy by default</div>
        <div class="d-box red">Public: anyone with link (opt-in, rare)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Delegation</div>
      <div class="d-flow-v">
        <div class="d-box purple">Executive delegation: assistant manages calendar</div>
        <div class="d-box purple">Delegate can: create events, respond to invites, view private events</div>
        <div class="d-box purple">Audit log: all actions tagged with delegate&#8217;s identity</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Caching ACLs</div>
      <div class="d-flow-v">
        <div class="d-box green">Redis: acl:{calendar_id} &#8594; SET of user_id:role</div>
        <div class="d-box green">TTL: 10 minutes, invalidate on ACL change</div>
        <div class="d-box gray">Avoid DB lookup on every event read</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-calendar-view",
		Title:       "Calendar View Rendering",
		Description: "Day, week, and month view query optimization with pre-aggregated time slots",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">View Queries</div>
      <div class="d-flow-v">
        <div class="d-box green">Day View<br/>SELECT * FROM events<br/>WHERE calendar_id IN (...)<br/>AND start_time &lt; day_end<br/>AND end_time &gt; day_start<br/>ORDER BY start_time<br/>&#8594; &#8776; 10&#8211;20 events, &lt; 5ms</div>
        <div class="d-box blue">Week View<br/>Same query, 7-day range<br/>&#8594; &#8776; 50&#8211;100 events, &lt; 15ms</div>
        <div class="d-box purple">Month View<br/>Same query, 30-day range<br/>+ aggregate into day counts<br/>&#8594; &#8776; 200&#8211;400 events, &lt; 50ms</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Recurring Event Expansion</div>
      <div class="d-flow-v">
        <div class="d-box amber">Fetch master events with RRULE</div>
        <div class="d-box amber">Expand RRULE within view window only</div>
        <div class="d-box amber">Merge with exception instances</div>
        <div class="d-box amber">Sort combined list by start_time</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Pre-Aggregated Slots (Month View)</div>
      <div class="d-flow-v">
        <div class="d-box indigo">Key: month:{user_id}:{calendar_id}:{YYYY-MM}</div>
        <div class="d-box indigo">Value: {day1: 3, day2: 0, day3: 5, ...}</div>
        <div class="d-box indigo">Updated async via Kafka consumer</div>
        <div class="d-box green">Month grid renders instantly from cache<br/>Detail loaded on day click</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Multi-Calendar Overlay</div>
      <div class="d-flow-v">
        <div class="d-box blue">User views 5 calendars simultaneously</div>
        <div class="d-box blue">Parallel queries: 1 per calendar</div>
        <div class="d-box blue">Client merges + color-codes by calendar</div>
        <div class="d-box gray">Server-side merge option for mobile (save bandwidth)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Index Strategy</div>
      <div class="d-flow-v">
        <div class="d-box green">Covering index: (calendar_id, start_time) INCLUDE (title, end_time, status)</div>
        <div class="d-box green">Avoids heap lookup for view rendering</div>
        <div class="d-box gray">Partial index: WHERE status != 'cancelled'</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

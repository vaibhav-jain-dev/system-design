package diagrams

func registerGoogleCalendar(r *Registry) {
	r.Register(&Diagram{
		Slug:        "gc-requirements",
		Title:       "Scale Estimates & Requirements",
		Description: "Functional and non-functional requirements with scale estimates for Google Calendar",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Estimates</div>
      <div class="d-flow-v">
        <div class="d-box blue">500M registered users &#8226; 100M DAU</div>
        <div class="d-box blue">Avg 5 events/user/week = 70M events/day</div>
        <div class="d-box blue">Calendar views: 100M &#215; 8 views/day = 800M/day</div>
        <div class="d-box purple">Read QPS: 9.3K &#8226; Peak (3x) = 28K</div>
        <div class="d-box purple">Write QPS: 810 &#8226; Peak (5x) = 4K</div>
        <div class="d-box amber">Storage: 500M users &#215; 200 events avg &#215; 1KB = 100 TB</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green">Create/update/delete events</div>
        <div class="d-box green">Recurring events (daily, weekly, custom RRULE)</div>
        <div class="d-box green">Invite attendees + RSVP (accept/decline/tentative)</div>
        <div class="d-box blue">Reminders (email, push, SMS)</div>
        <div class="d-box blue">Shared calendars with permissions</div>
        <div class="d-box blue">Conflict detection for overlapping events</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">Calendar load: &lt; 300ms p99</div>
        <div class="d-box purple">Sync latency: &lt; 2s cross-device</div>
        <div class="d-box amber">Availability: 99.99%</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-api-design",
		Title:       "API Endpoints",
		Description: "REST API design for calendar and event CRUD, invitations, and sync",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Event Operations</div>
      <div class="d-flow-v">
        <div class="d-box green" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">POST /v1/calendars/{cal_id}/events
{
  "title": "Team Standup",
  "start": "2024-01-15T09:00:00Z",
  "end": "2024-01-15T09:30:00Z",
  "timezone": "America/New_York",
  "recurrence": "RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR",
  "attendees": ["alice@co.com", "bob@co.com"],
  "reminders": [{"method": "push", "minutes": 10}]
}
&#8594; 201 {"event_id": "evt_abc123"}</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Calendar Views</div>
      <div class="d-flow-v">
        <div class="d-box blue" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">GET /v1/calendars/{cal_id}/events
  ?start=2024-01-15T00:00:00Z
  &amp;end=2024-01-22T00:00:00Z
  &amp;expand_recurring=true
&#8594; 200 {"events": [...], "sync_token": "s1"}</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Invitation Response</div>
      <div class="d-flow-v">
        <div class="d-box amber" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">PATCH /v1/events/{event_id}/rsvp
{"status": "accepted"}

GET /v1/calendars/sync
  ?sync_token=s1
&#8594; 200 {"changes": [...], "sync_token": "s2"}</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-recurring-events",
		Title:       "Recurring Events &#8212; RRULE Expansion",
		Description: "How RRULE recurrence rules are stored compactly and expanded at read time",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Storage: Single Row per Series</div>
        <div class="d-flow-v">
          <div class="d-box blue" style="font-family:var(--font-mono);font-size:0.78rem;text-align:left;white-space:pre">event_id: evt_standup
title: "Team Standup"
start: 2024-01-15T09:00:00Z
duration: 30min
rrule: FREQ=WEEKLY;BYDAY=MO,WE,FR
until: 2024-12-31
exceptions: [2024-03-25, 2024-07-04]</div>
          <div class="d-label">One row &#8594; represents 156 occurrences/year</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Read-Time Expansion</div>
        <div class="d-flow-v">
          <div class="d-box amber" style="text-align:center"><strong>Query: Jan 15-22</strong></div>
          <div class="d-arrow-down">&#8595; RRULE engine</div>
          <div class="d-box green" style="text-align:center">Mon Jan 15, 9:00 AM</div>
          <div class="d-box green" style="text-align:center">Wed Jan 17, 9:00 AM</div>
          <div class="d-box green" style="text-align:center">Fri Jan 19, 9:00 AM</div>
          <div class="d-label">Expand only within requested window</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Exception Handling</div>
    <div class="d-flow">
      <div class="d-box red" style="text-align:center"><strong>Delete one instance</strong><br>Add date to exceptions[]</div>
      <div class="d-box purple" style="text-align:center"><strong>Modify one instance</strong><br>Create override row with<br>original_event_id + occurrence_date</div>
      <div class="d-box amber" style="text-align:center"><strong>"This and following"</strong><br>Split: original UNTIL=date<br>New series starts at date</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end architecture from client through API gateway, event service, and supporting services",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue" style="text-align:center"><strong>Web App</strong></div>
    <div class="d-box blue" style="text-align:center"><strong>Mobile App</strong></div>
    <div class="d-box blue" style="text-align:center"><strong>CalDAV Client</strong></div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box gray" style="text-align:center"><strong>API Gateway + Load Balancer</strong><br>Auth &#8226; Rate limiting &#8226; Protocol translation (CalDAV &#8594; REST)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box green" style="text-align:center"><strong>Event Service</strong><br>CRUD events<br>RRULE expansion<br>Conflict detection</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box red" style="text-align:center"><strong>MySQL (Vitess)</strong><br>Events, calendars<br>Sharded by user_id</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box purple" style="text-align:center"><strong>Invitation Service</strong><br>Send/manage invites<br>RSVP tracking</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box gray" style="text-align:center"><strong>Kafka</strong><br>invitation-events<br>reminder-schedule</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box amber" style="text-align:center"><strong>Reminder Service</strong><br>Scheduled delivery<br>Push/Email/SMS</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo" style="text-align:center"><strong>Redis</strong><br>Reminder queue<br>Sync tokens</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-conflict-detection",
		Title:       "Conflict Detection Algorithm",
		Description: "How overlapping events are detected using interval overlap check on sorted event ranges",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Overlap Condition</div>
        <div class="d-flow-v">
          <div class="d-box indigo" style="text-align:center"><strong>Two events overlap when:</strong><br>event_A.start &lt; event_B.end AND event_B.start &lt; event_A.end</div>
          <div class="d-box green" style="text-align:center"><strong>No conflict</strong><br>Meeting A: 9:00-10:00<br>Meeting B: 10:00-11:00<br>10:00 &lt; 11:00 &#10003; BUT 10:00 &lt; 10:00 &#10007;</div>
          <div class="d-box red" style="text-align:center"><strong>Conflict!</strong><br>Meeting A: 9:00-10:30<br>Meeting B: 10:00-11:00<br>9:00 &lt; 11:00 &#10003; AND 10:00 &lt; 10:30 &#10003;</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Query Strategy</div>
        <div class="d-flow-v">
          <div class="d-box blue" style="text-align:center"><strong>1. New event request</strong><br>start=10:00, end=11:00</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber" style="text-align:center"><strong>2. Query existing events</strong><br>WHERE user_id = ? AND end &gt; 10:00 AND start &lt; 11:00<br>Uses composite index (user_id, start, end)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple" style="text-align:center"><strong>3. Expand recurring events</strong><br>Check RRULE occurrences within window<br>Include overrides, exclude exceptions</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" style="text-align:center"><strong>4. Return conflicts</strong><br>List overlapping events &#8226; Client shows warning<br>User can still create (soft conflict)</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-invitation-flow",
		Title:       "Invitation Flow &#8212; Invite &#8594; RSVP &#8594; Sync",
		Description: "End-to-end invitation lifecycle from organizer creating event to attendee RSVP and calendar sync",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>1. Organizer Creates Event</strong><br>POST /events with attendees: [alice, bob, carol]</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" style="text-align:center"><strong>2. Event Service</strong><br>Create event row &#8226; Create invitation rows (status=pending)<br>Publish to Kafka: invitation-events topic</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" style="text-align:center"><strong>3. Invitation Service Consumes</strong><br>For each attendee:<br>&#8226; Internal user &#8594; create shadow event on their calendar<br>&#8226; External user &#8594; send ICS email attachment</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box purple" style="text-align:center"><strong>Push Notification</strong><br>"You're invited to<br>Team Standup"</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box purple" style="text-align:center"><strong>Email</strong><br>ICS attachment<br>Accept/Decline buttons</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; attendee responds</div>
  <div class="d-box indigo" style="text-align:center"><strong>4. RSVP Update</strong><br>PATCH /events/{id}/rsvp {status: "accepted"}<br>Update invitation row &#8594; notify organizer &#8594; sync all attendee views</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" style="text-align:center"><strong>5. Calendar Sync</strong><br>All attendees see updated RSVP status on next sync<br>Organizer sees: 2 accepted, 1 declined</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-reminder-system",
		Title:       "Scheduled Reminder Pipeline",
		Description: "How reminders are scheduled, stored, and delivered via push, email, and SMS channels",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" style="text-align:center"><strong>Event Created/Updated</strong><br>Reminders: [{push, 10min}, {email, 1hr}]</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box green" style="text-align:center"><strong>Reminder Scheduler</strong><br>Compute fire times: event_start &#8722; reminder_offset<br>e.g., 9:00 AM event, 10min reminder &#8594; fire at 8:50 AM</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box red" style="text-align:center"><strong>Redis Sorted Set (Reminder Queue)</strong><br>ZADD reminders {fire_timestamp} {reminder_id}<br>Score = UTC fire time &#8226; Member = reminder payload</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" style="text-align:center"><strong>Reminder Worker (polls every 1s)</strong><br>ZRANGEBYSCORE reminders -inf {now} LIMIT 100<br>Process batch &#8594; ZREM after delivery</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-branch">
      <div class="d-branch-arm">
        <div class="d-box indigo" style="text-align:center"><strong>Push (FCM/APNs)</strong><br>Latency: &lt; 1s<br>90% of reminders</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box purple" style="text-align:center"><strong>Email (SES)</strong><br>Latency: &lt; 30s<br>Calendar attachment</div>
      </div>
      <div class="d-branch-arm">
        <div class="d-box gray" style="text-align:center"><strong>SMS (SNS)</strong><br>Latency: &lt; 5s<br>Premium feature</div>
      </div>
    </div>
  </div>
  <div class="d-label">At-least-once delivery: idempotency key per reminder &#8226; Dedup in last 5 min window</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-data-model",
		Title:       "Data Model &#8212; Core Tables",
		Description: "Database schema for calendars, events, recurrence rules, invitations, and reminders",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">calendars</div>
      <div class="d-entity-body">
        <div class="pk">calendar_id BIGINT</div>
        <div class="fk">owner_id BIGINT</div>
        <div>name VARCHAR(200)</div>
        <div>color VARCHAR(7)</div>
        <div>timezone VARCHAR(50)</div>
        <div class="idx idx-btree">is_primary BOOLEAN</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header green">events</div>
      <div class="d-entity-body">
        <div class="pk">event_id BIGINT (Snowflake)</div>
        <div class="fk">calendar_id BIGINT</div>
        <div class="fk">organizer_id BIGINT</div>
        <div>title VARCHAR(500)</div>
        <div>description TEXT</div>
        <div class="idx idx-btree">start_time TIMESTAMP</div>
        <div class="idx idx-btree">end_time TIMESTAMP</div>
        <div>timezone VARCHAR(50)</div>
        <div>location VARCHAR(500)</div>
        <div>is_recurring BOOLEAN</div>
        <div>status ENUM (confirmed|tentative|cancelled)</div>
        <div class="idx idx-btree">updated_at TIMESTAMP</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header amber">recurrence_rules</div>
      <div class="d-entity-body">
        <div class="pk">rule_id BIGINT</div>
        <div class="fk">event_id BIGINT</div>
        <div>rrule VARCHAR(500)</div>
        <div>until_date DATE NULL</div>
        <div>count INT NULL</div>
        <div>exceptions JSON</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header purple">invitations</div>
      <div class="d-entity-body">
        <div class="pk">invitation_id BIGINT</div>
        <div class="fk">event_id BIGINT</div>
        <div class="fk">attendee_id BIGINT</div>
        <div>email VARCHAR(320)</div>
        <div class="idx idx-hash">status ENUM (pending|accepted|declined|tentative)</div>
        <div>responded_at TIMESTAMP NULL</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header red">reminders</div>
      <div class="d-entity-body">
        <div class="pk">reminder_id BIGINT</div>
        <div class="fk">event_id BIGINT</div>
        <div class="fk">user_id BIGINT</div>
        <div>method ENUM (push|email|sms)</div>
        <div>offset_minutes INT</div>
        <div class="idx idx-btree">fire_at TIMESTAMP</div>
        <div>delivered BOOLEAN</div>
      </div>
    </div>
  </div>
</div>
<div class="d-er-lines">
  <div class="d-er-connector">
    <span class="d-er-from">calendars</span>
    <span class="d-er-type">1:N</span>
    <span class="d-er-to">events</span>
  </div>
  <div class="d-er-connector">
    <span class="d-er-from">events</span>
    <span class="d-er-type">1:1</span>
    <span class="d-er-to">recurrence_rules</span>
  </div>
  <div class="d-er-connector">
    <span class="d-er-from">events</span>
    <span class="d-er-type">1:N</span>
    <span class="d-er-to">invitations</span>
  </div>
  <div class="d-er-connector">
    <span class="d-er-from">events</span>
    <span class="d-er-type">1:N</span>
    <span class="d-er-to">reminders</span>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-timezone-handling",
		Title:       "Timezone Handling &#8212; UTC Storage + Display",
		Description: "How events are stored in UTC and converted to user's local timezone at display time",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Storage: Always UTC</div>
        <div class="d-flow-v">
          <div class="d-box blue" style="text-align:center"><strong>User creates event</strong><br>"Meeting at 3 PM EST"</div>
          <div class="d-arrow-down">&#8595; convert</div>
          <div class="d-box green" style="text-align:center"><strong>Stored as UTC</strong><br>start: 2024-01-15T20:00:00Z<br>end: 2024-01-15T21:00:00Z<br>timezone: America/New_York</div>
          <div class="d-label">IANA timezone ID stored alongside for display</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Display: Convert to Viewer's TZ</div>
        <div class="d-flow-v">
          <div class="d-box amber" style="text-align:center"><strong>Alice (NYC, UTC-5)</strong><br>Sees: 3:00 PM &#8211; 4:00 PM EST</div>
          <div class="d-box purple" style="text-align:center"><strong>Bob (London, UTC+0)</strong><br>Sees: 8:00 PM &#8211; 9:00 PM GMT</div>
          <div class="d-box indigo" style="text-align:center"><strong>Carol (Tokyo, UTC+9)</strong><br>Sees: 5:00 AM &#8211; 6:00 AM JST (+1 day)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Edge Cases</div>
    <div class="d-flow">
      <div class="d-box red" style="text-align:center"><strong>DST Transition</strong><br>Store wall-clock time + IANA TZ<br>"9 AM America/New_York" shifts<br>UTC offset automatically</div>
      <div class="d-box red" style="text-align:center"><strong>All-Day Events</strong><br>Store as DATE not TIMESTAMP<br>No timezone conversion<br>"Jan 15" everywhere</div>
      <div class="d-box red" style="text-align:center"><strong>Floating Time</strong><br>Reminders like "9 AM local"<br>Convert per-user at fire time<br>Not UTC-fixed</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-scaling",
		Title:       "Sharding Strategy &#8212; By User",
		Description: "Database sharding by user_id for events, calendars, and invitations",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box indigo" style="text-align:center"><strong>Shard Key: user_id (calendar owner)</strong><br>All calendars + events for one user on same shard</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 1 (users 0-125M)</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center"><strong>Primary</strong><br>MySQL 8.0<br>~25 TB</div>
        <div class="d-box gray" style="text-align:center">Read Replica 1</div>
        <div class="d-box gray" style="text-align:center">Read Replica 2</div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 2 (users 125-250M)</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center"><strong>Primary</strong><br>MySQL 8.0<br>~25 TB</div>
        <div class="d-box gray" style="text-align:center">Read Replica 1</div>
        <div class="d-box gray" style="text-align:center">Read Replica 2</div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 3 (users 250-375M)</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center"><strong>Primary</strong><br>MySQL 8.0<br>~25 TB</div>
        <div class="d-box gray" style="text-align:center">Read Replica 1</div>
        <div class="d-box gray" style="text-align:center">Read Replica 2</div>
      </div>
    </div>
    <div class="d-group" style="flex:1">
      <div class="d-group-title">Shard 4 (users 375-500M)</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center"><strong>Primary</strong><br>MySQL 8.0<br>~25 TB</div>
        <div class="d-box gray" style="text-align:center">Read Replica 1</div>
        <div class="d-box gray" style="text-align:center">Read Replica 2</div>
      </div>
    </div>
  </div>
  <div class="d-label">Cross-shard query: shared calendar events &#8594; scatter-gather across attendee shards. Mitigated by caching shared calendar in Redis.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-sync-protocol",
		Title:       "Sync Protocol &#8212; Incremental Sync",
		Description: "CalDAV-compatible incremental sync using sync tokens for efficient cross-device synchronization",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Initial Sync (Full)</div>
        <div class="d-flow-v">
          <div class="d-box blue" style="text-align:center"><strong>Client &#8594; Server</strong><br>GET /sync (no token)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" style="text-align:center"><strong>Server returns</strong><br>All events + sync_token = "v1_ts1704067200"</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box gray" style="text-align:center"><strong>Client stores</strong><br>Full event cache + sync_token</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Incremental Sync (Delta)</div>
        <div class="d-flow-v">
          <div class="d-box blue" style="text-align:center"><strong>Client &#8594; Server</strong><br>GET /sync?token=v1_ts1704067200</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber" style="text-align:center"><strong>Server queries</strong><br>WHERE updated_at &gt; token_timestamp<br>Returns only changed/deleted events</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green" style="text-align:center"><strong>Response</strong><br>3 created, 1 updated, 1 deleted<br>New sync_token = "v1_ts1704070800"</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Push Notification for Instant Sync</div>
    <div class="d-flow">
      <div class="d-box purple" style="text-align:center"><strong>Event changed</strong><br>on server</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box indigo" style="text-align:center"><strong>Push via WebSocket</strong><br>"calendar_changed" signal</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue" style="text-align:center"><strong>Client triggers</strong><br>incremental sync</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "gc-shared-calendars",
		Title:       "Shared Calendars &#8212; ACL & Permissions",
		Description: "Access control model for shared calendars with role-based permissions and sharing flows",
		ContentFile: "problems/google-calendar",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Permission Levels</div>
      <div class="d-flow-v">
        <div class="d-box green" style="text-align:center"><strong>Owner</strong><br>Full control &#8226; Delete calendar &#8226; Manage sharing</div>
        <div class="d-box blue" style="text-align:center"><strong>Editor</strong><br>Create/edit/delete events &#8226; Invite others</div>
        <div class="d-box amber" style="text-align:center"><strong>Viewer</strong><br>Read-only &#8226; See event details</div>
        <div class="d-box gray" style="text-align:center"><strong>Free/Busy Only</strong><br>See time slots &#8226; No event details</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">ACL Table</div>
      <div class="d-entity">
        <div class="d-entity-header purple">calendar_acl</div>
        <div class="d-entity-body">
          <div class="pk">acl_id BIGINT</div>
          <div class="fk">calendar_id BIGINT</div>
          <div>grantee_type ENUM (user|group|domain|public)</div>
          <div>grantee_id VARCHAR(320)</div>
          <div>role ENUM (owner|editor|viewer|freebusy)</div>
          <div>created_at TIMESTAMP</div>
        </div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Authorization Check</div>
      <div class="d-flow-v">
        <div class="d-box red" style="text-align:center">On every event read/write:<br>1. Check calendar_acl for user<br>2. Check group memberships<br>3. Check domain-wide sharing<br>Cache ACL in Redis &#8226; TTL 5 min</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

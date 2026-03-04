package diagrams

func registerRideHailing(r *Registry) {
	r.Register(&Diagram{
		Slug:        "rh-requirements",
		Title:       "Requirements & Scale Estimates",
		Description: "Functional requirements and scale targets for a ride-hailing system like Uber/Ola",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green">Rider requests ride with pickup &amp; destination</div>
        <div class="d-box green">Match rider with nearest available driver</div>
        <div class="d-box green">Real-time GPS tracking for rider &amp; driver</div>
        <div class="d-box green">Dynamic (surge) pricing based on demand</div>
        <div class="d-box blue">Fare calculation &amp; payment processing</div>
        <div class="d-box blue">Driver accepts/rejects ride offers</div>
        <div class="d-box gray">Shared rides (ride pooling)</div>
        <div class="d-box gray">Scheduled rides</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">100M registered users</div>
        <div class="d-box purple">20M rides/day (~230 rides/sec)</div>
        <div class="d-box purple">5M active drivers</div>
        <div class="d-box purple">1M concurrent rides at peak</div>
        <div class="d-box amber">Matching latency: &lt; 10 seconds</div>
        <div class="d-box amber">Location update: every 3&#8211;5 seconds</div>
        <div class="d-box amber">ETA accuracy: &#177;2 min for rides &lt; 20 min</div>
        <div class="d-box amber">Availability: 99.99%</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-api-design",
		Title:       "API Design",
		Description: "Core REST API endpoints for ride-hailing: ride requests, driver location, status updates",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Rider APIs</div>
      <div class="d-flow-v">
        <div class="d-box green">POST /rides/request</div>
        <div class="d-box blue" style="font-size:0.85em">{pickup_lat, pickup_lng, dest_lat, dest_lng, ride_type}</div>
        <div class="d-box green">GET /rides/{id}</div>
        <div class="d-box blue" style="font-size:0.85em">Returns: status, driver, ETA, fare_estimate</div>
        <div class="d-box green">PUT /rides/{id}/status</div>
        <div class="d-box blue" style="font-size:0.85em">{action: cancel} &#8212; rider cancellation</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Driver APIs</div>
      <div class="d-flow-v">
        <div class="d-box purple">POST /drivers/location</div>
        <div class="d-box blue" style="font-size:0.85em">{lat, lng, heading, speed} &#8212; every 3-5s</div>
        <div class="d-box purple">GET /drivers/nearby</div>
        <div class="d-box blue" style="font-size:0.85em">?lat=X&amp;lng=Y&amp;radius=5km &#8212; geospatial query</div>
        <div class="d-box purple">PUT /rides/{id}/status</div>
        <div class="d-box blue" style="font-size:0.85em">{action: accept|reject|arrive|start|complete}</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-ride-flow",
		Title:       "Ride Request Flow",
		Description: "End-to-end flow from ride request through matching, pickup, ride, and completion",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-box green">Rider requests ride</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue">Find nearby drivers</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue">Score &amp; rank</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple">Send offer to top driver</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber">Driver accepts</div>
</div>
<div class="d-flow" style="margin-top:1rem">
  <div class="d-box amber">Navigate to pickup</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber">Arrive &amp; confirm</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box green">Start ride</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box green">Complete ride</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box indigo">Process payment</div>
</div>
<div class="d-group" style="margin-top:1rem">
  <div class="d-group-title">Failure Paths</div>
  <div class="d-flow">
    <div class="d-box red">Driver rejects/timeout (15s)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber">Offer next ranked driver</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box red">All reject (3 attempts)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box gray">Expand radius &amp; retry</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-architecture",
		Title:       "High-Level Architecture",
		Description: "Service-oriented architecture for ride-hailing: API gateway, core services, data stores",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box gray">CDN (static assets)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">ALB / NLB</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">API Gateway</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Core Services</div>
    <div class="d-flow">
      <div class="d-box green">Ride Service</div>
      <div class="d-box green">Driver Service</div>
      <div class="d-box purple">Matching Service</div>
      <div class="d-box amber">Pricing Service</div>
      <div class="d-box indigo">Tracking Service</div>
      <div class="d-box blue">Payment Service</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Data Stores &amp; Infrastructure</div>
    <div class="d-flow">
      <div class="d-box purple">PostgreSQL (rides, users)</div>
      <div class="d-box red">Redis GEO (driver locations)</div>
      <div class="d-box amber">Kafka (ride events)</div>
      <div class="d-box blue">S3 (receipts, logs)</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-geospatial",
		Title:       "Geospatial Indexing (S2/H3)",
		Description: "Hexagonal grid cells for spatial indexing: drivers indexed by cell, query neighbors for radius search",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Spatial Index Strategy</div>
      <div class="d-flow-v">
        <div class="d-box blue">City divided into hex cells (H3 resolution 9 &#8776; 0.1 km&#178;)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Each driver mapped to cell ID on location update</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Redis: GEOADD drivers:{city} lng lat driver_id</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Query: GEORADIUS drivers:{city} lng lat 5km</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Neighbor Cell Query</div>
      <div class="d-flow-v">
        <div class="d-box amber">Rider location &#8594; resolve to cell</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Get k-ring neighbors (ring size = radius / cell size)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Fetch all drivers in center + neighbor cells</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Filter by Haversine distance &lt; 5 km</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Result: 10&#8211;50 candidate drivers</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-matching-algorithm",
		Title:       "Driver Matching Algorithm",
		Description: "Weighted scoring algorithm to rank candidate drivers for a ride request",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Scoring Formula</div>
    <div class="d-box blue" style="text-align:center">score = 0.4 &#215; proximity + 0.2 &#215; rating + 0.2 &#215; acceptance_rate + 0.2 &#215; ETA</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Factors (Normalized 0&#8211;1)</div>
        <div class="d-flow-v">
          <div class="d-box green">Proximity (40%): 1 &#8722; (dist / max_radius)</div>
          <div class="d-box blue">Rating (20%): driver_rating / 5.0</div>
          <div class="d-box purple">Acceptance Rate (20%): accepts / offers</div>
          <div class="d-box amber">ETA (20%): 1 &#8722; (eta / max_eta)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Ranking Flow</div>
        <div class="d-flow-v">
          <div class="d-box blue">10&#8211;50 candidates from geo query</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Compute score for each</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Sort descending by score</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Offer to #1, timeout 15s &#8594; #2, &#8230;</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-surge-pricing",
		Title:       "Surge Pricing",
		Description: "Dynamic pricing based on supply-demand ratio per geographic hex cell",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Supply &amp; Demand per Cell</div>
        <div class="d-flow-v">
          <div class="d-box blue">Demand: ride requests in cell / 5 min</div>
          <div class="d-box green">Supply: available drivers in cell</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Ratio = demand / supply</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box red">Surge multiplier = f(ratio)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Multiplier Tiers</div>
        <div class="d-flow-v">
          <div class="d-box green">ratio &lt; 1.0 &#8594; 1.0x (no surge)</div>
          <div class="d-box blue">ratio 1.0&#8211;2.0 &#8594; 1.2x&#8211;1.5x</div>
          <div class="d-box amber">ratio 2.0&#8211;3.0 &#8594; 1.5x&#8211;2.0x</div>
          <div class="d-box red">ratio &gt; 3.0 &#8594; 2.0x&#8211;3.0x (capped)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Fare Calculation</div>
    <div class="d-box purple" style="text-align:center">fare = base_fare + (distance &#215; per_km_rate &#215; surge) + (time &#215; per_min_rate) + booking_fee</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-tracking",
		Title:       "Real-Time Tracking",
		Description: "GPS location flow from driver device through Redis GEO to rider app via WebSocket",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-box green">Driver GPS (3&#8211;5s interval)</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue">Location Service</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box red">Redis GEO</div>
</div>
<div class="d-flow" style="margin-top:1rem">
  <div class="d-box red">Redis GEO</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple">Tracking Service</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber">WebSocket Gateway</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box green">Rider App</div>
</div>
<div class="d-group" style="margin-top:1rem">
  <div class="d-group-title">ETA Updates</div>
  <div class="d-flow">
    <div class="d-box blue">Current driver location</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple">Route engine (OSRM / Google Maps)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber">ETA recalculated every 30s</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">Push to rider via WebSocket</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-data-model",
		Title:       "Data Model",
		Description: "Core entities for ride-hailing: users, drivers, rides, ride events, payments, locations, surge zones",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">users</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> user_id (UUID)</div>
        <div class="d-box gray">name, email, phone, payment_method</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">drivers</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> driver_id (UUID)</div>
        <div class="d-box gray">name, phone, vehicle_type, rating, status</div>
        <div class="d-box green">status: {available, busy, offline}</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">driver_locations</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> driver_id</div>
        <div class="d-box amber">lat, lng, heading, speed, updated_at</div>
        <div class="d-box purple">Redis GEO + TTL 30s</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">rides</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> ride_id (UUID)</div>
        <div class="d-box gray"><span class="fk">FK</span> rider_id, <span class="fk">FK</span> driver_id</div>
        <div class="d-box gray">pickup_lat/lng, dest_lat/lng</div>
        <div class="d-box green">status: {requested, matched, arriving, in_progress, completed, cancelled}</div>
        <div class="d-box amber">fare, surge_multiplier, distance_km, duration_min</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">ride_events</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> event_id</div>
        <div class="d-box gray"><span class="fk">FK</span> ride_id, event_type, timestamp, metadata</div>
        <div class="d-box purple">Append-only log &#8594; Kafka topic</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">payments</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> payment_id</div>
        <div class="d-box gray"><span class="fk">FK</span> ride_id, amount, status, method</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">surge_zones</div>
      <div class="d-flow-v">
        <div class="d-box blue"><span class="pk">PK</span> cell_id (H3)</div>
        <div class="d-box amber">multiplier, demand_count, supply_count, updated_at</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-scaling",
		Title:       "Scaling Strategy",
		Description: "Scaling approach: read replicas, Redis cluster, Kafka partitioning, cell-based architecture",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Data Layer Scaling</div>
      <div class="d-flow-v">
        <div class="d-box blue">PostgreSQL: read replicas for user/ride queries</div>
        <div class="d-box red">Redis Cluster: 6+ nodes for GEO (sharded by city)</div>
        <div class="d-box amber">Kafka: partitioned by city_id for ride events</div>
        <div class="d-box purple">S3: receipts, ride history archive</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Service Scaling</div>
      <div class="d-flow-v">
        <div class="d-box green">Cell-based architecture: shard by city/region</div>
        <div class="d-box green">Each cell is independent: own DB, Redis, services</div>
        <div class="d-box blue">Matching Service: horizontally scaled per city</div>
        <div class="d-box amber">Tracking Service: WebSocket servers behind sticky LB</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Throughput Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">5M drivers &#215; 1 update/4s = 1.25M writes/sec (Redis)</div>
        <div class="d-box purple">230 rides/sec &#215; 6 state changes = 1,380 DB writes/sec</div>
        <div class="d-box purple">1M concurrent WebSocket connections</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-failure-modes",
		Title:       "Failure & Edge Cases",
		Description: "Key failure scenarios: cancellations, driver going offline, payment failures, GPS drift",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Cancellation Handling</div>
      <div class="d-flow-v">
        <div class="d-box red">Driver cancels after accepting</div>
        <div class="d-box gray">&#8594; Re-enter matching with remaining candidates</div>
        <div class="d-box red">Rider cancels after match</div>
        <div class="d-box gray">&#8594; Free driver, apply cancellation fee if &gt; 2 min</div>
        <div class="d-box red">Driver no-show (5 min timeout)</div>
        <div class="d-box gray">&#8594; Auto-cancel, no charge to rider, penalize driver</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Mid-Ride Failures</div>
      <div class="d-flow-v">
        <div class="d-box red">Driver goes offline mid-ride</div>
        <div class="d-box gray">&#8594; No GPS for 60s &#8594; alert ops, contact driver</div>
        <div class="d-box red">Payment fails at completion</div>
        <div class="d-box gray">&#8594; Retry 3x, then mark ride as payment_pending</div>
        <div class="d-box red">GPS drift / tunnel</div>
        <div class="d-box gray">&#8594; Snap to road via map-matching, interpolate</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "rh-shared-rides",
		Title:       "Shared Rides Matching",
		Description: "Algorithm for matching riders heading in the same direction with detour constraints",
		ContentFile: "problems/ride-hailing",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Shared Ride Matching Algorithm</div>
    <div class="d-flow-v">
      <div class="d-box blue">New ride request with pickup P2, destination D2</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box green">Find active rides within pickup radius (&lt; 1 km)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box green">Filter: heading similarity &gt; 0.7 (cosine of direction vectors)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">Calculate detour for each candidate ride</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">detour = route(P1 &#8594; P2 &#8594; D1 &#8594; D2) &#8722; route(P1 &#8594; D1)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box purple">Accept if detour &lt; 40% of original ride &amp;&amp; ETA increase &lt; 10 min</div>
    </div>
  </div>
  <div class="d-cols" style="margin-top:1rem">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Drop-off Order Optimization</div>
        <div class="d-flow-v">
          <div class="d-box blue">Try all permutations of drop-off order</div>
          <div class="d-box green">Pick order that minimizes total travel time</div>
          <div class="d-box amber">Constraint: no rider&#8217;s ETA increases &gt; 50%</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Fare Split</div>
        <div class="d-flow-v">
          <div class="d-box purple">Each rider pays: solo_fare &#215; 0.6&#8211;0.7</div>
          <div class="d-box green">Platform earns: combined fares &#8722; driver payout</div>
          <div class="d-box gray">Driver earns: ~80% of solo fare equivalent</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})
}

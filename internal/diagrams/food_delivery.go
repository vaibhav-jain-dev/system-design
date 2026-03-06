package diagrams

func registerFoodDelivery(r *Registry) {
	r.Register(&Diagram{
		Slug:        "fd-requirements",
		Title:       "Requirements & Scale Estimates",
		Description: "Functional requirements and scale targets for a food delivery platform like Swiggy/Zomato",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green">Browse restaurants &amp; menus by location</div>
        <div class="d-box green">Place order with cart &amp; payment</div>
        <div class="d-box green">Real-time order tracking (GPS)</div>
        <div class="d-box green">Match &amp; assign delivery partners</div>
        <div class="d-box blue">Surge pricing based on demand/supply</div>
        <div class="d-box blue">Ratings &amp; reviews for restaurants/partners</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Non-Functional Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">50M registered users</div>
        <div class="d-box purple">5M orders/day (~58 orders/sec avg)</div>
        <div class="d-box purple">500K active delivery partners</div>
        <div class="d-box amber">Partner matching &lt; 30 seconds</div>
        <div class="d-box amber">100K concurrent in-flight orders</div>
        <div class="d-box amber">GPS update ingestion: 125K writes/sec (500K partners &#215; 1 update/4s)</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Key Decisions</div>
      <div class="d-flow-v">
        <div class="d-box red">Push vs Pull for location updates?</div>
        <div class="d-box red">Fan-out-on-write vs Fan-out-on-read for order status?</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-api-design",
		Title:       "API Design",
		Description: "Core REST API endpoints for the food delivery platform",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Restaurant &amp; Menu</div>
    <div class="d-flow-v">
      <div class="d-box blue">GET /restaurants?lat={lat}&amp;lng={lng}&amp;radius=5km&amp;cuisine=indian</div>
      <div class="d-box blue">GET /restaurants/{id}/menu</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Order Management</div>
    <div class="d-flow-v">
      <div class="d-box green">POST /orders &#8212; { restaurant_id, items[], address, payment_method }</div>
      <div class="d-box green">GET /orders/{id}/track &#8212; SSE stream: status, partner location, ETA</div>
      <div class="d-box green">PUT /orders/{id}/status &#8212; { status: confirmed | picked_up | delivered }</div>
      <div class="d-box amber">POST /orders/{id}/cancel &#8212; { reason } (allowed before pickup only)</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Delivery Partner</div>
    <div class="d-flow-v">
      <div class="d-box purple">GET /partners/nearby?lat={lat}&amp;lng={lng}&amp;radius=3km</div>
      <div class="d-box purple">PUT /partners/{id}/location &#8212; { lat, lng, heading, speed }</div>
      <div class="d-box purple">POST /partners/{id}/accept &#8212; { order_id }</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-order-lifecycle",
		Title:       "Order Lifecycle",
		Description: "State machine for an order from placement through delivery",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-box blue">Order Placed</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box amber">Restaurant Confirms</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple">Partner Assigned</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box indigo">Pickup</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box green">In Transit</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box green">Delivered</div>
</div>
<div class="d-flow" style="margin-top: 1rem;">
  <div class="d-group">
    <div class="d-group-title">Failure Branches</div>
    <div class="d-flow">
      <div class="d-box red">Restaurant Rejects</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">Refund Initiated</div>
    </div>
    <div class="d-flow" style="margin-top: 0.5rem;">
      <div class="d-box red">No Partner Available</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Retry (3x, 30s interval)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">Cancel + Refund</div>
    </div>
    <div class="d-flow" style="margin-top: 0.5rem;">
      <div class="d-box red">Customer Cancels</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Partial refund if food prepared</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end system architecture with core microservices",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box gray">Customer App</div>
    <div class="d-box gray">Partner App</div>
    <div class="d-box gray">Restaurant Dashboard</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box blue">CDN (CloudFront)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box blue">ALB</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box indigo">API Gateway (Auth + Rate Limit)</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Core Services</div>
    <div class="d-flow">
      <div class="d-box green">Order Service</div>
      <div class="d-box green">Restaurant Service</div>
      <div class="d-box green">Partner Service</div>
      <div class="d-box green">Tracking Service</div>
      <div class="d-box green">Payment Service</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Supporting Services</div>
    <div class="d-flow">
      <div class="d-box amber">Matching Engine</div>
      <div class="d-box amber">ETA Service</div>
      <div class="d-box amber">Surge Pricing</div>
      <div class="d-box amber">Notification Service</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Data Layer</div>
    <div class="d-flow">
      <div class="d-box purple">PostgreSQL (Orders, Users)</div>
      <div class="d-box purple">Redis Cluster (Geo, Sessions)</div>
      <div class="d-box purple">Kafka (Order Events)</div>
      <div class="d-box purple">S3 (Images, Menus)</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-partner-matching",
		Title:       "Delivery Partner Matching Algorithm",
		Description: "Scoring and assignment flow for matching orders to delivery partners",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">New Order (restaurant_lat, restaurant_lng)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">GEORADIUS 3km from restaurant in Redis</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Scoring Function (per candidate)</div>
    <div class="d-cols">
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box green">Distance to restaurant (40% weight)</div>
          <div class="d-box green">Acceptance rate (20% weight)</div>
        </div>
      </div>
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box green">Partner rating (20% weight)</div>
          <div class="d-box green">Current load: 0 or 1 orders (20% weight)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber">Rank candidates by composite score</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box purple">Push offer to top partner (30s TTL)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">Partner accepts &#8594; Assign</div>
  </div>
  <div class="d-flow">
    <div class="d-box red">Timeout / Reject</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber">Offer to next partner (up to 3 retries)</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-real-time-tracking",
		Title:       "Real-Time Location Tracking",
		Description: "GPS ingestion and fan-out pipeline for live order tracking",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-flow">
  <div class="d-box gray">Partner App (GPS every 4s)</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box blue">Location Ingestion Service</div>
  <div class="d-arrow">&#8594;</div>
  <div class="d-box purple">Redis GEO (GEOADD partner:{id})</div>
</div>
<div class="d-flow-v" style="margin-top: 1rem;">
  <div class="d-group">
    <div class="d-group-title">Fan-Out to Customer</div>
    <div class="d-flow">
      <div class="d-box purple">Redis Pub/Sub (channel: order:{id})</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">WebSocket Server</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box gray">Customer App (map update)</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Scale Numbers</div>
    <div class="d-flow">
      <div class="d-box amber">500K partners &#215; 0.25 updates/sec = 125K writes/sec</div>
      <div class="d-box amber">100K active orders = 100K Pub/Sub channels</div>
      <div class="d-box amber">Redis GEO: ~40 bytes/entry = ~20MB for all partners</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-eta-prediction",
		Title:       "ETA Prediction Pipeline",
		Description: "ML-based ETA estimation combining multiple signal sources",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Input Signals</div>
    <div class="d-flow">
      <div class="d-box blue">Distance (route API)</div>
      <div class="d-box blue">Traffic conditions</div>
      <div class="d-box blue">Restaurant prep time (historical P50)</div>
      <div class="d-box blue">Partner avg speed</div>
      <div class="d-box blue">Time of day / weather</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box purple">Feature Vector</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box indigo">Gradient Boosted Tree Model</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">ETA &#177; confidence interval</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">ETA Components (shown to customer)</div>
    <div class="d-flow">
      <div class="d-box amber">Prep time: ~15 min</div>
      <div class="d-box amber">Partner to restaurant: ~8 min</div>
      <div class="d-box amber">Restaurant to customer: ~12 min</div>
      <div class="d-box green">Total ETA: 35 min (30&#8211;40 range)</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-surge-pricing",
		Title:       "Surge Pricing Algorithm",
		Description: "Dynamic pricing based on real-time demand/supply ratio per zone",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-group">
      <div class="d-group-title">Demand Signal</div>
      <div class="d-flow-v">
        <div class="d-box blue">Orders placed per zone (5 min window)</div>
        <div class="d-box blue">Search activity in zone</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Supply Signal</div>
      <div class="d-flow-v">
        <div class="d-box green">Available partners per zone</div>
        <div class="d-box green">Partners en route to zone</div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box indigo">Demand / Supply Ratio per Geohash Zone</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box green">Ratio &lt; 1.5 &#8594; 1.0x (no surge)</div>
    <div class="d-box amber">Ratio 1.5&#8211;3.0 &#8594; 1.2x&#8211;1.8x</div>
    <div class="d-box red">Ratio &gt; 3.0 &#8594; 2.0x&#8211;2.5x (capped)</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box purple">Base Fare + (Distance &#215; Rate) + Platform Fee</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple">&#215; Surge Multiplier</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber">Final Delivery Fee</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-data-model",
		Title:       "Data Model",
		Description: "Core database entities and relationships for the food delivery platform",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">restaurants</div>
      <div class="d-entity-body">
        <div><span class="pk">id</span> UUID</div>
        <div>name VARCHAR(255)</div>
        <div>lat DECIMAL(9,6)</div>
        <div>lng DECIMAL(9,6)</div>
        <div>cuisine VARCHAR(50)</div>
        <div>avg_prep_time_min INT</div>
        <div>rating DECIMAL(2,1)</div>
        <div>is_active BOOLEAN</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header blue">menu_items</div>
      <div class="d-entity-body">
        <div><span class="pk">id</span> UUID</div>
        <div><span class="fk">restaurant_id</span> UUID</div>
        <div>name VARCHAR(255)</div>
        <div>price_cents INT</div>
        <div>category VARCHAR(50)</div>
        <div>is_available BOOLEAN</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header green">orders</div>
      <div class="d-entity-body">
        <div><span class="pk">id</span> UUID</div>
        <div><span class="fk">customer_id</span> UUID</div>
        <div><span class="fk">restaurant_id</span> UUID</div>
        <div>status ENUM(placed, confirmed, assigned, picked_up, delivered, cancelled)</div>
        <div>total_cents INT</div>
        <div>delivery_fee_cents INT</div>
        <div>delivery_address_lat DECIMAL(9,6)</div>
        <div>delivery_address_lng DECIMAL(9,6)</div>
        <div>created_at TIMESTAMP</div>
        <div><span class="idx idx-btree">idx_customer_created</span></div>
        <div><span class="idx idx-btree">idx_restaurant_status</span></div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header green">order_items</div>
      <div class="d-entity-body">
        <div><span class="pk">id</span> UUID</div>
        <div><span class="fk">order_id</span> UUID</div>
        <div><span class="fk">menu_item_id</span> UUID</div>
        <div>quantity INT</div>
        <div>unit_price_cents INT</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header purple">delivery_partners</div>
      <div class="d-entity-body">
        <div><span class="pk">id</span> UUID</div>
        <div>name VARCHAR(255)</div>
        <div>phone VARCHAR(15)</div>
        <div>vehicle_type ENUM(bike, scooter, car)</div>
        <div>rating DECIMAL(2,1)</div>
        <div>acceptance_rate DECIMAL(3,2)</div>
        <div>is_online BOOLEAN</div>
        <div>current_lat DECIMAL(9,6)</div>
        <div>current_lng DECIMAL(9,6)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header amber">deliveries</div>
      <div class="d-entity-body">
        <div><span class="pk">id</span> UUID</div>
        <div><span class="fk">order_id</span> UUID</div>
        <div><span class="fk">partner_id</span> UUID</div>
        <div>assigned_at TIMESTAMP</div>
        <div>picked_up_at TIMESTAMP</div>
        <div>delivered_at TIMESTAMP</div>
        <div>distance_km DECIMAL(5,2)</div>
      </div>
    </div>
    <div class="d-entity">
      <div class="d-entity-header red">payments</div>
      <div class="d-entity-body">
        <div><span class="pk">id</span> UUID</div>
        <div><span class="fk">order_id</span> UUID</div>
        <div>amount_cents INT</div>
        <div>method ENUM(card, upi, wallet, cod)</div>
        <div>status ENUM(pending, captured, refunded)</div>
        <div>gateway_txn_id VARCHAR(64)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-scaling",
		Title:       "Scaling Strategy",
		Description: "Scaling approach for each component of the food delivery platform",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Read Path (Menus &amp; Restaurants)</div>
      <div class="d-flow-v">
        <div class="d-box blue">CDN for static menu images</div>
        <div class="d-box blue">PostgreSQL read replicas (3x)</div>
        <div class="d-box blue">Redis cache: restaurant metadata, TTL 5 min</div>
        <div class="d-box blue">ElastiCache for search results, TTL 1 min</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Geo &amp; Location (Hot Path)</div>
      <div class="d-flow-v">
        <div class="d-box purple">Redis Cluster (6 shards) for GEOADD/GEORADIUS</div>
        <div class="d-box purple">Shard by geohash prefix (city-level)</div>
        <div class="d-box purple">125K writes/sec across cluster</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Order Events (Async)</div>
      <div class="d-flow-v">
        <div class="d-box green">Kafka: order-events topic, 32 partitions</div>
        <div class="d-box green">Partition by order_id for ordering guarantee</div>
        <div class="d-box green">Consumers: analytics, notifications, billing</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Notifications</div>
      <div class="d-flow-v">
        <div class="d-box amber">SQS queue per notification type</div>
        <div class="d-box amber">Push: FCM/APNs, SMS: Twilio, Email: SES</div>
        <div class="d-box amber">Deduplication by order_id + event_type</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Database Sharding</div>
      <div class="d-flow-v">
        <div class="d-box indigo">Orders table: shard by city_id (geographic locality)</div>
        <div class="d-box indigo">Partners table: shard by city_id</div>
        <div class="d-box indigo">Restaurants: read-heavy, replicas suffice</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-failure-handling",
		Title:       "Failure & Edge Cases",
		Description: "Handling common failure scenarios in food delivery",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Partner Phone Dies Mid-Delivery</div>
    <div class="d-flow">
      <div class="d-box red">No GPS update for 60s</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Mark partner stale, show last known location</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">After 5 min: SMS partner, alert support</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">After 15 min: reassign to new partner</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Restaurant Delays Preparation</div>
    <div class="d-flow">
      <div class="d-box red">Prep time exceeds 2x estimated</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Update customer ETA + push notification</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Delay partner dispatch to avoid idle wait</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Customer Cancels After Preparation</div>
    <div class="d-flow">
      <div class="d-box red">Cancel request (food already prepared)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Partial refund (delivery fee only)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Compensate restaurant for food cost</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Payment Fails After Order Placed</div>
    <div class="d-flow">
      <div class="d-box red">Payment gateway timeout</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Retry 3x with exponential backoff</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Hold order (do not send to restaurant)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">After 2 min: cancel order, notify customer</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Dual Writes / Consistency</div>
    <div class="d-flow">
      <div class="d-box indigo">Order DB write + Kafka publish &#8594; use Transactional Outbox</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box indigo">CDC (Debezium) tails WAL &#8594; Kafka &#8594; at-least-once delivery</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "fd-batch-delivery",
		Title:       "Batch Delivery Optimization",
		Description: "Algorithm for assigning multiple orders to one partner heading in the same direction",
		ContentFile: "problems/food-delivery",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Batching Eligibility Check</div>
    <div class="d-flow">
      <div class="d-box blue">Partner has 1 active order</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">New order&#39;s restaurant within 500m of current route</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Delivery address within 1.5km of existing drop-off</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Route Optimization</div>
    <div class="d-cols">
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box green">Calculate detour time for pickup</div>
          <div class="d-box green">Calculate detour time for drop-off</div>
          <div class="d-box green">Total added time must be &lt; 10 min</div>
        </div>
      </div>
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box amber">Constraint: existing order ETA cannot increase &gt; 5 min</div>
          <div class="d-box amber">Constraint: max 2 orders per partner</div>
          <div class="d-box amber">Constraint: both restaurants must be ready within 5 min</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box purple">Optimal pickup order (nearest first)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple">Optimal drop-off order (nearest first)</div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box green">Discount delivery fee 30% for batched order</div>
  </div>
</div>`,
	})
}

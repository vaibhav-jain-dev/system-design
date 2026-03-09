package diagrams

func registerLoggingSystem(r *Registry) {
	r.Register(&Diagram{
		Slug:        "log-requirements",
		Title:       "Requirements &amp; Scale",
		Description: "Scale targets, storage math, and retention tiers for a distributed logging system",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">1M servers (fleet-wide)</div>
        <div class="d-box purple">10K log lines/server/sec = <strong>10B logs/sec peak</strong></div>
        <div class="d-box purple">Avg log line: 100 bytes &#8594; <strong>1TB/sec ingestion rate</strong></div>
        <div class="d-box purple">Storage: 100GB/server/day &#215; 1M = <strong>100PB/day</strong></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Query Requirements</div>
      <div class="d-flow-v">
        <div class="d-box blue">Full-text search: &lt; 5s across 7-day hot window</div>
        <div class="d-box blue">Metrics freshness: &lt; 1 min from emission to dashboard</div>
        <div class="d-box blue">Trace lookup by trace_id: &lt; 2s</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Retention Tiers</div>
      <div class="d-flow-v">
        <div class="d-box green"><strong>Hot</strong> &#8212; 0&#8211;7 days &#8212; Elasticsearch SSD, fully indexed</div>
        <div class="d-box amber"><strong>Warm</strong> &#8212; 7&#8211;30 days &#8212; Elasticsearch HDD, read-only, shrunk shards</div>
        <div class="d-box red"><strong>Cold</strong> &#8212; 30 days&#8211;1 year &#8212; S3 Parquet, Athena on-demand query</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Key Constraints</div>
      <div class="d-flow-v">
        <div class="d-box gray">Write path: fire-and-forget, eventual durability OK</div>
        <div class="d-box gray">No log loss on agent crash (offset tracking)</div>
        <div class="d-box gray">Compression 10:1 expected with gzip on structured JSON</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-architecture",
		Title:       "High-Level Architecture",
		Description: "End-to-end logging pipeline: agents, Kafka buffer, processors, and query layer",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Ingestion Layer</div>
    <div class="d-flow">
      <div class="d-box gray">Servers (1M)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Log Agent (Filebeat / Fluentd)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Kafka Cluster (buffer, 30-day retention)</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Log Pipeline</div>
        <div class="d-flow-v">
          <div class="d-box amber">Log Processor (Logstash / Flink)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Elasticsearch Cluster</div>
          <div class="d-label">Hot: SSD, 7 days, full index</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Kibana (query &amp; dashboards)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Metrics Pipeline</div>
        <div class="d-flow-v">
          <div class="d-box amber">Metrics Processor (Flink)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">InfluxDB / Prometheus</div>
          <div class="d-label">Time-series, downsampled over time</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Grafana (metrics dashboards)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Archive Pipeline</div>
        <div class="d-flow-v">
          <div class="d-box amber">S3 Archiver (Flink sink)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">S3 Parquet (cold, 1 year)</div>
          <div class="d-label">Partitioned by service/year/month/day</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box indigo">Athena (SQL on-demand queries)</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-ingestion-pipeline",
		Title:       "Log Ingestion Pipeline",
		Description: "Log agent behavior: tail, parse, enrich, batch, compress, and send to Kafka",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Log Agent (per server &#8212; Filebeat / Fluentd)</div>
    <div class="d-flow">
      <div class="d-box gray">tail -F /var/log/app.log</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Parse JSON / plaintext</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Add metadata</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Batch 1,000 lines</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">gzip compress</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Kafka producer</div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Metadata Added by Agent</div>
        <div class="d-flow-v">
          <div class="d-box green">host &#8212; server hostname</div>
          <div class="d-box green">service &#8212; from /etc/fluentd/service.conf</div>
          <div class="d-box green">environment &#8212; prod / staging / dev</div>
          <div class="d-box green">region &#8212; AWS region / datacenter</div>
          <div class="d-box green">agent_version &#8212; for schema migration</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Kafka Topic Design</div>
        <div class="d-flow-v">
          <div class="d-box purple">Topic: logs-raw</div>
          <div class="d-box purple">Partition key: service_name</div>
          <div class="d-label">All logs from same service ordered on same partition</div>
          <div class="d-box purple">256 partitions &#215; replication factor 3</div>
          <div class="d-box purple">Retention: 30 days (replay buffer)</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Throughput Math</div>
        <div class="d-flow-v">
          <div class="d-box amber">10B logs/sec &#215; 100 bytes = 1TB/sec raw</div>
          <div class="d-box amber">gzip 10:1 &#8594; 100GB/sec on the wire</div>
          <div class="d-box amber">Kafka broker: 1Gbps NIC &#8594; need 800 brokers</div>
          <div class="d-label">Or use 10Gbps NICs &#8594; 80 brokers for ingestion</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-data-model",
		Title:       "Log Data Model",
		Description: "Log schema and Elasticsearch index mapping with field types and analysis strategy",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header blue">Log Record (JSON)</div>
      <div class="d-entity-body">
        <div><span class="pk">PK</span> log_id (Snowflake ID)</div>
        <div>timestamp (ISO 8601, ms precision)</div>
        <div>service (string &#8212; e.g. payment-service)</div>
        <div>level (DEBUG / INFO / WARN / ERROR / FATAL)</div>
        <div>message (free-text string)</div>
        <div>host (server hostname)</div>
        <div>environment (prod / staging / dev)</div>
        <div>trace_id (nullable &#8212; links to distributed trace)</div>
        <div>span_id (nullable &#8212; segment within trace)</div>
        <div>metadata (map&lt;string, string&gt; &#8212; arbitrary fields)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-entity">
      <div class="d-entity-header amber">Elasticsearch Index Mapping</div>
      <div class="d-entity-body">
        <div><span class="idx idx-btree">date</span> timestamp &#8212; range queries, ILM rollover</div>
        <div><span class="idx idx-hash">keyword</span> level &#8212; exact match filter (not analyzed)</div>
        <div><span class="idx idx-hash">keyword</span> service &#8212; exact match, aggregation</div>
        <div><span class="idx idx-hash">keyword</span> host &#8212; exact match</div>
        <div><span class="idx idx-hash">keyword</span> trace_id &#8212; exact match lookup</div>
        <div><span class="idx idx-gin">text</span> message &#8212; full-text, English analyzer</div>
        <div>metadata.* &#8212; dynamic mapping, keyword default</div>
      </div>
    </div>
    <div class="d-group" style="margin-top:8px">
      <div class="d-group-title">keyword vs text</div>
      <div class="d-flow-v">
        <div class="d-box green"><strong>keyword</strong> &#8212; stored as-is, used for exact filter, sort, aggregation</div>
        <div class="d-box blue"><strong>text</strong> &#8212; tokenized &amp; analyzed, used for full-text match (slower)</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-elasticsearch-sharding",
		Title:       "Elasticsearch Shard Strategy",
		Description: "Index-per-day strategy, ILM lifecycle, and shard sizing for log data",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Index Strategy: One Index per Service per Day</div>
    <div class="d-flow">
      <div class="d-box blue">10 services &#215; 30 days</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">300 active indices</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">5 primary + 5 replica shards per index</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Target shard size: 10&#8211;50GB each</div>
    </div>
    <div class="d-label">Naming: logs-{service}-{YYYY.MM.DD}. Allows dropping old indices without reindexing.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">ILM (Index Lifecycle Management) Policy</div>
    <div class="d-flow">
      <div class="d-box green">
        <strong>Hot</strong><br/>
        0&#8211;7 days<br/>
        SSD nodes<br/>
        All shards active<br/>
        Write + read
      </div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">
        <strong>Warm</strong><br/>
        7&#8211;30 days<br/>
        HDD nodes<br/>
        Shrink to 1 shard<br/>
        Read-only
      </div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">
        <strong>Cold</strong><br/>
        30&#8211;90 days<br/>
        S3 searchable snapshots<br/>
        On-demand mount<br/>
        Slow read
      </div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box gray">
        <strong>Delete / Archive</strong><br/>
        &gt;90 days<br/>
        ES index deleted<br/>
        Parquet on S3<br/>
        Athena queries
      </div>
    </div>
  </div>
  <div class="d-group" style="margin-top:8px">
    <div class="d-group-title">Capacity Math</div>
    <div class="d-flow">
      <div class="d-box purple">100PB/day raw &#247; 10 (gzip) = 10PB/day compressed</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">7-day hot tier = 70PB SSD needed</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">30-day warm = 300PB HDD</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-metrics-pipeline",
		Title:       "Metrics Pipeline",
		Description: "Metrics ingestion via Prometheus pull and Statsd push, storage, downsampling, and Grafana",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Push Model (StatsD / Telegraf)</div>
        <div class="d-flow-v">
          <div class="d-box gray">Server emits custom metrics</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">StatsD agent (UDP, fire-and-forget)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Telegraf aggregator</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">InfluxDB</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Pull Model (Prometheus)</div>
        <div class="d-flow-v">
          <div class="d-box gray">Server exposes /metrics endpoint</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box blue">Prometheus scrapes every 15s</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Prometheus TSDB (local)</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Remote write &#8594; Thanos / Cortex (long-term)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Metric Format &amp; Downsampling</div>
    <div class="d-flow">
      <div class="d-box indigo">&#123; name, tags&#123;host, service, env&#125;, value, timestamp &#125;</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">15s resolution &#8212; 15 days</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">1min resolution &#8212; 90 days</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">1hr resolution &#8212; 1 year</div>
    </div>
    <div class="d-label">Downsampling: Flink job rolls up older buckets. Reduces storage 240&#215; from 15s &#8594; 1hr.</div>
  </div>
  <div class="d-group" style="margin-top:8px">
    <div class="d-group-title">Query &amp; Visualization</div>
    <div class="d-flow">
      <div class="d-box indigo">Grafana</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">PromQL (Prometheus queries)</div>
      <div class="d-arrow">+</div>
      <div class="d-box blue">Flux / InfluxQL (InfluxDB queries)</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-alerting",
		Title:       "Alerting System",
		Description: "Alert rule evaluation, deduplication, grouping, routing, and lifecycle states",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Alert Rule Example</div>
    <div class="d-box purple"><strong>IF</strong> error_rate(service=payment, window=5m) <strong>&gt; 5%</strong> <strong>FOR</strong> 5 minutes &#8594; fire alert</div>
    <div class="d-label">Rules evaluated every 30s by Alertmanager against Prometheus metrics.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Alert Manager Processing</div>
        <div class="d-flow-v">
          <div class="d-box blue">Deduplicate &#8212; 100 hosts fire same alert &#8594; 1 notification</div>
          <div class="d-box blue">Group by service &#8212; batch related alerts per service</div>
          <div class="d-box blue">Route by severity &#8212; critical / warning / info</div>
          <div class="d-box blue">Inhibit &#8212; suppress child alerts if parent is firing</div>
          <div class="d-box blue">Silence &#8212; mute during maintenance windows</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Notification Routing</div>
        <div class="d-flow-v">
          <div class="d-box red"><strong>Critical</strong> &#8212; PagerDuty (on-call engineer woken up)</div>
          <div class="d-box amber"><strong>Warning</strong> &#8212; Slack #alerts-warning channel</div>
          <div class="d-box green"><strong>Info</strong> &#8212; Email digest (daily summary)</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Alert Lifecycle</div>
    <div class="d-flow">
      <div class="d-box gray"><strong>Inactive</strong><br/>Rule not firing</div>
      <div class="d-arrow">&#8594; threshold crossed &#8594;</div>
      <div class="d-box amber"><strong>Pending</strong><br/>Rule firing, within grace period<br/>(5 min &#8220;for&#8221; duration)</div>
      <div class="d-arrow">&#8594; grace elapsed &#8594;</div>
      <div class="d-box red"><strong>Firing</strong><br/>Alert active, notifications sent</div>
      <div class="d-arrow">&#8594; metric drops below threshold &#8594;</div>
      <div class="d-box green"><strong>Resolved</strong><br/>Recovery notification sent</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-search-query",
		Title:       "Log Search Architecture",
		Description: "Query parsing, Elasticsearch scatter-gather, and Redis caching for repeated searches",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">User Query &#8594; Elasticsearch DSL</div>
    <div class="d-box gray"><em>Query: &#8220;error level:ERROR service:payment-service last:1h&#8221;</em></div>
    <div class="d-arrow-down">&#8595; parsed &#8595;</div>
    <div class="d-box purple"><strong>bool query:</strong><br/>
      &nbsp;&nbsp;filter: [term(level=ERROR), term(service=payment-service)]<br/>
      &nbsp;&nbsp;filter: [range(timestamp &gt; now-1h)]<br/>
      &nbsp;&nbsp;must: [match(message=error)]
    </div>
    <div class="d-label">filter clauses are cached and do not affect relevance score (faster). must clause applies full-text scoring.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Query Execution (Scatter-Gather)</div>
    <div class="d-flow">
      <div class="d-box blue">Coordinating Node</div>
      <div class="d-arrow">&#8594; broadcast &#8594;</div>
      <div class="d-flow-v">
        <div class="d-box green">Shard 1</div>
        <div class="d-box green">Shard 2</div>
        <div class="d-box green">Shard N</div>
      </div>
      <div class="d-arrow">&#8594; top-K results &#8594;</div>
      <div class="d-box blue">Coordinating Node (merge)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Top 100 hits returned</div>
    </div>
    <div class="d-label">Each shard returns top 100 local hits. Coordinating node merges N&#215;100 results, re-ranks, returns global top 100.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Search Cache (Redis)</div>
    <div class="d-flow">
      <div class="d-box red">Hash(query + time_bucket) &#8594; Redis key</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Cache hit &#8212; return in &lt; 1ms</div>
      <div class="d-box gray">Cache miss &#8212; execute ES query (~1-3s)</div>
    </div>
    <div class="d-label">TTL: 60s for hot queries. Invalidate on new index segment flush. Covers repeated dashboard refreshes.</div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-distributed-tracing",
		Title:       "Distributed Tracing Integration",
		Description: "trace_id propagation across services, Jaeger visualization, and trace-to-log correlation",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Trace Context Propagation</div>
    <div class="d-flow">
      <div class="d-box gray">Request enters API Gateway</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Generate trace_id (Snowflake: globally unique)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box blue">Propagate via HTTP headers: X-Trace-ID, X-Span-ID</div>
    </div>
    <div class="d-flow" style="margin-top:8px">
      <div class="d-box green">Service A<br/>span_id: s1<br/>logs tagged: trace_id:T1, span:s1</div>
      <div class="d-arrow">&#8594; calls &#8594;</div>
      <div class="d-box green">Service B<br/>span_id: s2<br/>logs tagged: trace_id:T1, span:s2</div>
      <div class="d-arrow">&#8594; calls &#8594;</div>
      <div class="d-box green">Service C<br/>span_id: s3<br/>logs tagged: trace_id:T1, span:s3</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Jaeger / Zipkin Visualization</div>
        <div class="d-flow-v">
          <div class="d-box indigo">Waterfall view of all spans</div>
          <div class="d-box blue">&#9632; API Gateway (0&#8211;150ms)</div>
          <div class="d-box green">&nbsp;&nbsp;&#9632; Service A (5&#8211;80ms)</div>
          <div class="d-box amber">&nbsp;&nbsp;&nbsp;&nbsp;&#9632; Redis call (5&#8211;8ms)</div>
          <div class="d-box green">&nbsp;&nbsp;&#9632; Service B (80&#8211;145ms)</div>
          <div class="d-box red">&nbsp;&nbsp;&nbsp;&nbsp;&#9632; DB query (80&#8211;140ms) &#8593; slow</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Log Correlation by trace_id</div>
        <div class="d-flow-v">
          <div class="d-box purple">Kibana: filter trace_id = abc123</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">All log lines from all services for that request</div>
          <div class="d-label">Elasticsearch query: term filter on trace_id keyword field &#8212; O(1) lookup.</div>
          <div class="d-box amber">Link from Jaeger span &#8594; Kibana filtered log view</div>
          <div class="d-label">Deep-link: kibana/logs?trace_id=abc123 &#8212; one click from trace UI to all logs.</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-cold-storage",
		Title:       "Cold Storage &amp; Archival",
		Description: "S3 Parquet archival path, partitioning strategy, and Athena query cost model",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Archival Path (after 90 days)</div>
    <div class="d-flow">
      <div class="d-box red">Elasticsearch index deleted (ILM delete action)</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box amber">Data lives in S3 as Parquet (written by Flink archiver continuously)</div>
    </div>
    <div class="d-label">Flink archiver consumes Kafka and writes Parquet directly &#8212; independent of ES lifecycle.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">S3 Path Structure (Hive Partitioning)</div>
    <div class="d-box indigo">s3://logs-archive/{service}/{year}/{month}/{day}/part-*.parquet</div>
    <div class="d-label">Example: s3://logs-archive/payment-service/2024/01/15/part-00001.parquet</div>
    <div class="d-label">Partition pruning: querying 1 service, 1 day reads only that prefix &#8212; not the full archive.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Athena Query</div>
        <div class="d-box purple">SELECT * FROM logs<br/>WHERE service = &#8216;payment&#8217;<br/>&nbsp;&nbsp;AND date = &#8216;2024-01-15&#8217;<br/>&nbsp;&nbsp;AND level = &#8216;ERROR&#8217;</div>
        <div class="d-label">Athena scans only the matching partition. Parquet columnar format skips non-selected columns.</div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Cost Analysis</div>
        <div class="d-flow-v">
          <div class="d-box amber">Athena: $5 per TB scanned</div>
          <div class="d-box green">1-day partition = 100GB compressed &#8594; $0.50 per query</div>
          <div class="d-box green">Full-month scan = 3TB &#8594; $15 per query</div>
          <div class="d-box red">No partition = full archive = 300PB &#8594; $1.5M per query</div>
          <div class="d-label">Always filter by service + date. Enforce in query UI to prevent runaway costs.</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-failure-handling",
		Title:       "Failure Handling",
		Description: "Recovery strategies for Kafka outage, Elasticsearch indexer failure, and log agent crash",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scenario 1: Kafka Outage</div>
      <div class="d-flow-v">
        <div class="d-box red"><strong>Failure:</strong> Kafka cluster unreachable</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber"><strong>Agent behavior:</strong> buffer to local disk ring buffer (max 1GB per agent)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Ring buffer: oldest entries dropped when full (debug logs first)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><strong>Recovery:</strong> Kafka restores &#8594; agent drains disk buffer in order &#8594; no log loss for &lt; 1h outage</div>
        <div class="d-label">1GB buffer at 100KB/sec net = 10,000 seconds &#8776; 2.7 hours of buffering per agent.</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scenario 2: ES Indexer Failure</div>
      <div class="d-flow-v">
        <div class="d-box red"><strong>Failure:</strong> Logstash / Flink indexer crashes or Elasticsearch is down</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber"><strong>Effect:</strong> Kafka consumer stops committing offsets &#8594; consumer lag grows</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Alert: Kafka consumer lag &gt; 100K messages &#8594; PagerDuty</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><strong>Recovery:</strong> Fix ES &#8594; restart indexer &#8594; consumer resumes from last committed offset &#8594; catches up</div>
        <div class="d-label">Kafka 30-day retention is the replay buffer. Catch-up rate must exceed ingest rate to drain lag.</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scenario 3: Log Agent Crash</div>
      <div class="d-flow-v">
        <div class="d-box red"><strong>Failure:</strong> Filebeat / Fluentd process crashes on a server</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber"><strong>State:</strong> Filebeat stores file offset in ~/.filebeat/registry</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box blue">Registry persists to disk on each checkpoint (every 1s)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green"><strong>Recovery:</strong> Agent restarts, reads registry, continues tailing from last offset &#8212; at most 1s of duplicate logs (idempotent ZADD in ES handles deduplication)</div>
        <div class="d-label">log_id (Snowflake) used as ES document _id &#8212; re-indexing same log is a no-op update.</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "log-monitoring-cost",
		Title:       "Cost Analysis &amp; Optimization",
		Description: "Storage cost breakdown, compression savings, and verbosity controls for log cost management",
		ContentFile: "problems/logging-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Naive Cost (No Optimization)</div>
    <div class="d-flow">
      <div class="d-box red">100PB/day &#215; $0.023/GB = <strong>$2.3M/day</strong> if all stored on S3</div>
      <div class="d-arrow">+</div>
      <div class="d-box red">Hot ES SSD: 7 days &#215; 100PB &#215; $10/GB = <strong>$7B one-time</strong></div>
    </div>
    <div class="d-label">Clearly infeasible. Optimization is not optional &#8212; it is a core design constraint.</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Optimization 1: Compression</div>
        <div class="d-flow-v">
          <div class="d-box green">gzip structured JSON logs: 10:1 ratio</div>
          <div class="d-box green">Real ingest: 100PB &#247; 10 = <strong>10PB/day</strong></div>
          <div class="d-box green">S3 cost: $0.023/GB &#215; 10TB = <strong>$230/day</strong></div>
          <div class="d-label">30-day S3 warm tier: $6,900/month &#8212; very manageable.</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Optimization 2: Log Level Filtering</div>
        <div class="d-flow-v">
          <div class="d-box amber"><strong>Production:</strong> INFO and above only</div>
          <div class="d-box amber">Drop DEBUG/TRACE at agent level (never sent to Kafka)</div>
          <div class="d-box amber">DEBUG typically 60-70% of raw log volume</div>
          <div class="d-label">Combined with compression: effective ratio 25-30:1 vs raw DEBUG-on logs.</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Optimization 3: Dynamic Sampling</div>
        <div class="d-flow-v">
          <div class="d-box blue">DEBUG logs in staging: 100% sampled</div>
          <div class="d-box blue">DEBUG logs in prod (emergency): 0.1% sampled</div>
          <div class="d-box blue">INFO: 10% sampled for high-volume services</div>
          <div class="d-box blue">WARN / ERROR / FATAL: always 100%</div>
          <div class="d-label">Sampling rate configurable per service via feature flag &#8212; no agent restart needed.</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group" style="margin-top:8px">
    <div class="d-group-title">Realistic Cost After Optimization</div>
    <div class="d-flow">
      <div class="d-box green">Hot ES (7 days, compressed): ~700TB SSD &#8212; ~$70K/month</div>
      <div class="d-arrow">+</div>
      <div class="d-box amber">Warm (30 days, HDD): ~3PB &#8212; ~$30K/month</div>
      <div class="d-arrow">+</div>
      <div class="d-box blue">Cold S3 (1 year): ~120PB &#8212; ~$2.8K/month</div>
    </div>
  </div>
</div>`,
	})
}

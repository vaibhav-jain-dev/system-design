package diagrams

func registerPaymentSystem(r *Registry) {
	r.Register(&Diagram{
		Slug:        "pay-requirements",
		Title:       "Requirements & Scale",
		Description: "Functional requirements, scale targets, and compliance constraints for a global payment system",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Functional Requirements</div>
      <div class="d-flow-v">
        <div class="d-box green"><span class="d-step">1</span>Process payments — debit user, credit merchant atomically <div class="d-tag green">&#10003; core</div></div>
        <div class="d-box green"><span class="d-step">2</span>Idempotency — retried payments never double-charge <div class="d-tag green">&#10003; core</div></div>
        <div class="d-box green"><span class="d-step">3</span>Multi-currency — USD, EUR, GBP, 150+ currencies</div>
        <div class="d-box green"><span class="d-step">4</span>Fraud detection — score every transaction &lt;200ms</div>
        <div class="d-box blue"><span class="d-step">5</span>Refunds — full or partial, within 90 days <div class="d-tag blue">P1</div></div>
        <div class="d-box blue"><span class="d-step">6</span>Settlement — batch at EOD, reconcile vs bank statement <div class="d-tag blue">P1</div></div>
        <div class="d-box gray"><span class="d-step">7</span>Subscription billing — recurring charges on schedule <div class="d-tag gray">P2</div></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Scale Targets</div>
      <div class="d-flow-v">
        <div class="d-box purple">10 million transactions per day <span class="d-metric throughput">10M/day</span></div>
        <div class="d-box purple">116 TPS average <span class="d-metric throughput">116 TPS avg</span></div>
        <div class="d-box purple">1,000 TPS peak (flash sales, holidays) <span class="d-metric throughput">1K TPS peak</span></div>
        <div class="d-box purple">Global — multi-region active-active</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Non-Functional &amp; Compliance</div>
      <div class="d-flow-v">
        <div class="d-box red">PCI DSS Level 1 — card data handling <div class="d-tag red">mandatory</div></div>
        <div class="d-box amber">Availability: 99.999% — 5.26 min downtime/year <span class="d-metric throughput">5 nines</span></div>
        <div class="d-box amber">Latency: &lt;2s end-to-end payment processing <span class="d-metric latency">&lt;2s</span></div>
        <div class="d-box amber">Consistency: zero double-charges — ACID ledger</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-architecture",
		Title:       "High-Level Architecture",
		Description: "Request path from client through API gateway, payment service, fraud/auth/ledger, and external card networks",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue">Client<br><small>browser / mobile / merchant SDK</small></div>
    <div class="d-arrow">&#8594; HTTPS</div>
    <div class="d-box green">API Gateway<br><small>JWT auth<br>rate limit: 100 req/s per merchant<br>idempotency key check</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box indigo">Payment Service<br><small>orchestrates saga<br>owns transaction state machine</small></div>
  </div>
  <div class="d-arrow-down">&#8595; fan-out (parallel where possible)</div>
  <div class="d-group">
    <div class="d-group-title">Core Services</div>
    <div class="d-flow">
      <div class="d-box amber">Fraud Detection<br><small>ML score &lt;200ms<br>feature store: Redis</small></div>
      <div class="d-box purple">Authorization Service<br><small>calls Visa/MC/Amex API<br>3DS challenge if needed</small></div>
      <div class="d-box green">Ledger Service<br><small>double-entry accounting<br>PostgreSQL SERIALIZABLE</small></div>
      <div class="d-box blue">Settlement Service<br><small>batch at EOD<br>reconcile vs bank</small></div>
    </div>
  </div>
  <div class="d-flow">
    <div class="d-arrow-down">&#8595; card network APIs</div>
    <div class="d-arrow-down">&#8595; Kafka (async)</div>
  </div>
  <div class="d-flow">
    <div class="d-box red">Visa / Mastercard / Bank<br><small>external auth APIs, ~300ms</small></div>
    <div class="d-box amber">Kafka<br><small>settlement events<br>audit log stream</small></div>
    <div class="d-box gray">Reconciliation Service<br><small>compare ledger vs bank statement<br>flag mismatches &#8594; human review</small></div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-double-entry-ledger",
		Title:       "Double-Entry Ledger",
		Description: "Every transaction debits one account and credits another; sum of all entries is always zero",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Accounting Invariant</div>
    <div class="d-flow">
      <div class="d-box red">&#931; Debits = &#931; Credits<br><small>Net of all ledger entries = $0.00 always.<br>If this breaks, data is corrupted.</small> <div class="d-tag red">invariant</div></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Example: User pays Merchant $100</div>
    <div class="d-cols">
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box amber">Ledger entry 1 of 2<br><strong>DEBIT</strong> — User account<br><small>amount: -$100.00<br>balance_after: $400.00<br>tx_id: tx_9876</small></div>
        </div>
      </div>
      <div class="d-col">
        <div class="d-flow-v">
          <div class="d-box green">Ledger entry 2 of 2<br><strong>CREDIT</strong> — Merchant account<br><small>amount: +$100.00<br>balance_after: $1,100.00<br>tx_id: tx_9876</small></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-entity">
    <div class="d-entity-header blue">ledger_entries</div>
    <div class="d-entity-body">
      <div><span class="pk">PK</span> entry_id (UUID)</div>
      <div><span class="fk">FK</span> tx_id (transaction)</div>
      <div>account_id (user or merchant)</div>
      <div>amount (DECIMAL 18,4 — never FLOAT)</div>
      <div>type (DEBIT | CREDIT)</div>
      <div>balance_after (DECIMAL 18,4)</div>
      <div>created_at (timestamp with tz)</div>
    </div>
  </div>
  <div class="d-group" style="margin-top:8px">
    <div class="d-group-title">PostgreSQL Atomicity</div>
    <div class="d-flow">
      <div class="d-box purple">BEGIN SERIALIZABLE;<br>SELECT balance FROM accounts WHERE id=user FOR UPDATE;<br>INSERT INTO ledger_entries (debit $100);<br>INSERT INTO ledger_entries (credit $100);<br>UPDATE accounts SET balance=balance-100 WHERE id=user;<br>UPDATE accounts SET balance=balance+100 WHERE id=merchant;<br>COMMIT;<br><small>All-or-nothing — partial failure = full rollback</small></div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-payment-flow",
		Title:       "Payment Processing Flow",
		Description: "Detailed sequence: validation, fraud check, pending ledger entry, card network auth, capture or reversal",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-flow">
    <div class="d-box blue"><span class="d-step">1</span>Validate request<br><small>schema check, currency valid<br>idempotency key lookup<br>~5ms</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box amber"><span class="d-step">2</span>Fraud score<br><small>ML model inference<br>feature store: Redis<br>target &lt;200ms</small></div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-box purple"><span class="d-step">3</span>Create PENDING entry<br><small>write to ledger: status=pending<br>reserve funds<br>~10ms (PostgreSQL)</small></div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-box indigo"><span class="d-step">4</span>Call card network<br><small>POST to Visa/MC API<br>auth code or decline<br>~300ms (external)</small></div>
    <div class="d-arrow">&#8594; success</div>
    <div class="d-box green"><span class="d-step">5a</span>Capture<br><small>update ledger: status=completed<br>CREDIT merchant account<br>publish settled event</small></div>
  </div>
  <div class="d-flow">
    <div class="d-box gray">step 4 &#8594;</div>
    <div class="d-arrow">failure / timeout</div>
    <div class="d-box red"><span class="d-step">5b</span>Reverse<br><small>update ledger: status=failed<br>release reserved funds<br>return 402 to client</small></div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Latency Budget (success path)</div>
    <div class="d-flow">
      <div class="d-box gray">Validation: ~5ms</div>
      <div class="d-arrow">+</div>
      <div class="d-box gray">Fraud: ~150ms</div>
      <div class="d-arrow">+</div>
      <div class="d-box gray">Ledger write: ~10ms</div>
      <div class="d-arrow">+</div>
      <div class="d-box gray">Card network: ~300ms</div>
      <div class="d-arrow">+</div>
      <div class="d-box gray">Capture: ~10ms</div>
      <div class="d-arrow">=</div>
      <div class="d-box green">~475ms P50 <span class="d-metric latency">&lt;2s P99</span></div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-idempotency",
		Title:       "Idempotency Design",
		Description: "Idempotency-Key header + Redis SET NX prevents double charges on client retries",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Client-Side Key Format</div>
    <div class="d-flow">
      <div class="d-box blue">Idempotency-Key: {client_id}:{uuid_v4}<br><small>Example: merchant_42:550e8400-e29b-41d4-a716-446655440000<br>Client generates once per payment intent — never reuse across different payments</small></div>
    </div>
  </div>
  <div class="d-cols">
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">First Request</div>
        <div class="d-flow-v">
          <div class="d-box blue">POST /payments<br><small>Idempotency-Key: merchant_42:abc123</small></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box amber">Redis SET NX idempotency:merchant_42:abc123 &#34;processing&#34; EX 86400<br><small>SET if Not eXists — atomic check-and-set</small></div>
          <div class="d-arrow-down">&#8595; key was new &#8594; proceed</div>
          <div class="d-box green">Process payment &#8594; charge card</div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Redis SET idempotency:merchant_42:abc123 {&#34;status&#34;:&#34;success&#34;, &#34;tx_id&#34;:&#34;tx_9876&#34;} EX 86400<br><small>overwrite with final result</small></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box green">Return 201 {tx_id: &#34;tx_9876&#34;, status: &#34;success&#34;}</div>
        </div>
      </div>
    </div>
    <div class="d-col">
      <div class="d-group">
        <div class="d-group-title">Retry (network failure scenario)</div>
        <div class="d-flow-v">
          <div class="d-box blue">POST /payments (retry)<br><small>same Idempotency-Key: merchant_42:abc123</small></div>
          <div class="d-arrow-down">&#8595;</div>
          <div class="d-box purple">Redis GET idempotency:merchant_42:abc123<br><small>&#8594; returns cached result</small></div>
          <div class="d-arrow-down">&#8595; key already exists</div>
          <div class="d-box amber">Return cached response immediately<br><small>201 {tx_id: &#34;tx_9876&#34;, status: &#34;success&#34;}</small></div>
          <div class="d-box green">No second charge issued <div class="d-tag green">&#10003; safe retry</div></div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">In-Flight Race (two concurrent retries)</div>
    <div class="d-flow">
      <div class="d-box indigo">Both retries arrive simultaneously &#8594; only one wins SET NX &#8594; loser polls Redis until result appears &#8594; both return same response<br><small>No locking needed beyond Redis atomic SET NX</small></div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-saga-pattern",
		Title:       "Saga Pattern for Distributed Transactions",
		Description: "Choreography-based saga with compensating transactions on failure; no distributed 2PC",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Happy Path — All Steps Succeed</div>
    <div class="d-flow">
      <div class="d-box green"><span class="d-step">1</span>Authorize<br><small>card network<br>&#8594; auth_code issued</small></div>
      <div class="d-arrow">&#8594; event</div>
      <div class="d-box green"><span class="d-step">2</span>Reserve Funds<br><small>ledger: PENDING<br>&#8594; funds_reserved event</small></div>
      <div class="d-arrow">&#8594; event</div>
      <div class="d-box green"><span class="d-step">3</span>Capture<br><small>bank API<br>&#8594; payment_captured event</small></div>
      <div class="d-arrow">&#8594; event</div>
      <div class="d-box green"><span class="d-step">4</span>Settle<br><small>ledger: COMPLETE<br>notify merchant</small></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Failure at Step 3 — Compensating Transactions (reverse order)</div>
    <div class="d-flow">
      <div class="d-box green"><span class="d-step">1</span>Authorize &#10003;</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green"><span class="d-step">2</span>Reserve &#10003;</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red"><span class="d-step">3</span>Capture FAIL<br><small>bank API timeout</small></div>
    </div>
    <div class="d-flow" style="margin-top:8px">
      <div class="d-box amber">Compensate 2: Release Reserved Funds<br><small>ledger: status=FAILED, release hold</small></div>
      <div class="d-arrow">&#8592; then</div>
      <div class="d-box amber">Compensate 1: Void Authorization<br><small>POST card-network/void auth_code</small></div>
      <div class="d-arrow">&#8592; then</div>
      <div class="d-box red">Return 402 Payment Required<br><small>include error_code: &#34;bank_capture_failed&#34;</small></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Choreography vs Orchestration</div>
    <div class="d-flow">
      <div class="d-box indigo">Each service publishes events to Kafka &#8594; next service listens and reacts. No central coordinator. Failure recovery: compensating event triggers reverse chain. Trade-off: harder to trace than orchestrator pattern but no single point of failure.</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-fraud-detection",
		Title:       "Fraud Detection Pipeline",
		Description: "Real-time ML scoring: feature extraction from Redis, model inference, 3-tier decision threshold",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Feature Extraction (sub-20ms from Redis)</div>
    <div class="d-flow">
      <div class="d-box blue">Transaction<br><small>amount, currency<br>merchant_category<br>device_fingerprint</small></div>
      <div class="d-arrow">+</div>
      <div class="d-box purple">Real-Time Features (Redis)<br><small>velocity: 10 tx in last hour?<br>new merchant for this user?<br>location change since last tx?</small></div>
      <div class="d-arrow">+</div>
      <div class="d-box amber">Historical Features (data warehouse)<br><small>avg tx amount (30d)<br>top merchants (90d)<br>chargeback history</small></div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; ~100ms feature assembly</div>
  <div class="d-group">
    <div class="d-group-title">ML Model Inference</div>
    <div class="d-flow">
      <div class="d-box indigo">Gradient Boosted Trees (XGBoost)<br><small>trained on labeled fraud/non-fraud data (weekly retrain)<br>input: ~50 features &#8594; output: fraud_score 0&#8211;100<br>inference: ~50ms on CPU, ~10ms on GPU</small></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Decision Thresholds</div>
    <div class="d-flow">
      <div class="d-box green">Score 0&#8211;29<br><strong>Approve</strong><br><small>low risk<br>proceed to auth</small> <span class="d-metric throughput">~92% of tx</span></div>
      <div class="d-box amber">Score 30&#8211;69<br><strong>3DS Challenge</strong><br><small>step-up auth<br>OTP to cardholder</small> <span class="d-metric throughput">~7% of tx</span></div>
      <div class="d-box red">Score 70&#8211;100<br><strong>Decline</strong><br><small>return 402<br>flag for review</small> <span class="d-metric throughput">~1% of tx</span></div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-pci-dss",
		Title:       "PCI DSS Compliance",
		Description: "Card data never stored post-auth; tokenization via HSM vault; CVV never persisted",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">What PCI DSS Prohibits</div>
    <div class="d-flow">
      <div class="d-box red">Never store CVV / CVC after authorization<br><small>PCI DSS Req 3.2.1 — no exceptions, even encrypted</small> <div class="d-tag red">prohibited</div></div>
      <div class="d-box red">Never store full PAN unencrypted<br><small>must truncate, hash, or encrypt with HSM key</small> <div class="d-tag red">prohibited</div></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Tokenization Flow</div>
    <div class="d-flow">
      <div class="d-box blue"><span class="d-step">1</span>Card entered by user<br><small>PAN: 4111 1111 1111 1234<br>CVV: 123, Exp: 12/27</small></div>
      <div class="d-arrow">&#8594; HTTPS only</div>
      <div class="d-box red"><span class="d-step">2</span>Vault Service (HSM)<br><small>encrypt PAN with HSM key<br>generate token: tok_abc789<br>discard CVV immediately</small></div>
      <div class="d-arrow">&#8594; token only</div>
      <div class="d-box green"><span class="d-step">3</span>Merchant system stores token<br><small>tok_abc789 stored in merchant DB<br>useless outside vault — no card data exposure</small></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Charging a Saved Card</div>
    <div class="d-flow">
      <div class="d-box blue">Merchant sends: {token: &#34;tok_abc789&#34;, amount: 100}</div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box red">Vault decrypts: PAN=4111...1234<br><small>HSM decryption — key never leaves HSM hardware</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box purple">Send actual PAN to card network<br><small>over TLS to Visa/MC — PAN in transit only</small></div>
      <div class="d-arrow">&#8594;</div>
      <div class="d-box green">Auth response &#8594; discard PAN again</div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Network Segmentation</div>
    <div class="d-flow">
      <div class="d-box indigo">Cardholder Data Environment (CDE) — isolated VPC subnet, no internet access, all traffic logged, quarterly pen-test required. Only Vault Service lives here. All other services use tokens only.</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-consistency",
		Title:       "Consistency & ACID Guarantees",
		Description: "PostgreSQL SERIALIZABLE for balance updates; deadlock prevention; global multi-region eventual consistency",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Single-Region ACID (PostgreSQL)</div>
      <div class="d-flow-v">
        <div class="d-box green"><strong>Atomicity</strong> — debit + credit in one transaction<br><small>partial failure = full rollback</small></div>
        <div class="d-box green"><strong>Consistency</strong> — balance &gt;= 0 constraint enforced<br><small>CHECK (balance &gt;= 0) on accounts table</small></div>
        <div class="d-box green"><strong>Isolation: SERIALIZABLE</strong><br><small>prevents phantom reads during balance check<br>two concurrent charges see each other&#39;s locks</small></div>
        <div class="d-box green"><strong>Durability</strong> — WAL flush before commit ack<br><small>synchronous_commit = on (default)</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Deadlock Prevention</div>
      <div class="d-flow-v">
        <div class="d-box red">Deadlock scenario:<br><small>Tx A: lock account 1 &#8594; lock account 2<br>Tx B: lock account 2 &#8594; lock account 1<br>&#8594; circular wait = deadlock</small></div>
        <div class="d-box green">Fix: always lock accounts in ascending account_id order<br><small>Tx A: lock min(1,2)=1 first, then 2<br>Tx B: also locks 1 first &#8594; waits, no cycle</small></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Multi-Region Consistency</div>
      <div class="d-flow-v">
        <div class="d-box amber">Each region: independent PostgreSQL ledger<br><small>strong consistency within region</small></div>
        <div class="d-box amber">Cross-region: Kafka event stream<br><small>eventual consistency: ~500ms replication lag</small></div>
        <div class="d-box blue">Reconciliation detects cross-region drift<br><small>nightly job compares regional ledgers</small></div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-settlement",
		Title:       "Settlement & Reconciliation",
		Description: "Async batch settlement at EOD, Kafka fan-out, bank file generation, and discrepancy reconciliation",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-group">
    <div class="d-group-title">Settlement Flow (T+1)</div>
    <div class="d-flow">
      <div class="d-box green"><span class="d-step">1</span>Payment captured<br><small>ledger: status=CAPTURED<br>event published to Kafka</small></div>
      <div class="d-arrow">&#8594; async</div>
      <div class="d-box amber"><span class="d-step">2</span>Settlement Service batches<br><small>aggregate by merchant + currency<br>group all EOD captures into batch file</small></div>
      <div class="d-arrow">&#8594; EOD batch</div>
      <div class="d-box purple"><span class="d-step">3</span>Send to bank<br><small>ACH / SWIFT file transfer<br>bank processes T+1 (next business day)</small></div>
      <div class="d-arrow">&#8594; T+1</div>
      <div class="d-box green"><span class="d-step">4</span>Funds in merchant account<br><small>standard: T+1<br>instant settlement: +0.5% fee</small></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Reconciliation (nightly)</div>
    <div class="d-flow">
      <div class="d-box blue">Internal ledger<br><small>all CAPTURED txns for day<br>sum by merchant</small></div>
      <div class="d-arrow">vs</div>
      <div class="d-box indigo">Bank statement<br><small>actual settled amounts<br>received via SFTP file</small></div>
      <div class="d-arrow">&#8594; compare</div>
      <div class="d-box green">Match &#10003;<br><small>mark reconciled<br>no action</small></div>
      <div class="d-arrow">or</div>
      <div class="d-box red">Mismatch &#10007;<br><small>flag for human review<br>SLA: resolved within 24h</small></div>
    </div>
  </div>
  <div class="d-group">
    <div class="d-group-title">Settlement Tiers</div>
    <div class="d-flow">
      <div class="d-box gray">Standard: T+1 &#8212; no extra fee (default for all merchants)</div>
      <div class="d-box amber">Next-day: same day EOD &#8212; +0.25% fee</div>
      <div class="d-box purple">Instant: within 30 min &#8212; +0.5% fee (push to debit card)</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-failure-handling",
		Title:       "Failure Handling",
		Description: "Circuit breaker for bank API timeouts, saga compensation for DB failures, DLQ for consumer lag",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Bank API Timeout</div>
      <div class="d-flow-v">
        <div class="d-box red">Symptom: POST /authorize times out after 3s</div>
        <div class="d-box amber">Circuit Breaker (Hystrix / Resilience4J)<br><small>threshold: 3 failures in 10s &#8594; OPEN circuit<br>open for 30s, then half-open probe</small></div>
        <div class="d-box green">When OPEN: return 503 immediately<br><small>no resource held, clear error to client<br>do NOT retry silently &#8212; double-charge risk</small></div>
        <div class="d-box blue">Fallback: cached auth for low-risk transactions<br><small>fraud score &lt;10 + amount &lt;$50 &#8594; approve offline<br>reconcile on bank recovery</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Database Failure During Saga</div>
      <div class="d-flow-v">
        <div class="d-box red">Symptom: ledger write fails mid-saga (step 3 of 4)</div>
        <div class="d-box amber">Saga compensating transactions fire in reverse:<br><small>void bank authorization &#8594; release reserved funds<br>mark tx as FAILED in idempotency cache</small></div>
        <div class="d-box green">Client receives deterministic 500 + error_id<br><small>retry with same idempotency key is safe</small></div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">Kafka Consumer Lag</div>
      <div class="d-flow-v">
        <div class="d-box amber">Settlement consumer falls behind<br><small>settlement is async — not on critical path</small></div>
        <div class="d-box blue">SLA: reconciliation within 24 hours<br><small>alert if consumer lag &gt;1 hour</small></div>
        <div class="d-box green">Payments still succeed — settlement eventually catches up</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pay-monitoring",
		Title:       "Monitoring & Alerting",
		Description: "Key payment metrics, SLO targets, and alert thresholds with business impact callouts",
		ContentFile: "problems/payment-system",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Key Metrics</div>
      <div class="d-flow-v">
        <div class="d-box purple">Payment success rate<br><small>target &gt;99.9% — below 99% = revenue impact</small> <span class="d-metric throughput">&gt;99.9%</span></div>
        <div class="d-box purple">End-to-end latency P50 / P99<br><small>target: P50 &lt;500ms, P99 &lt;2s</small> <span class="d-metric latency">&lt;2s P99</span></div>
        <div class="d-box amber">Fraud rate<br><small>target &lt;0.1% — above 0.5% triggers model retrain</small> <span class="d-metric throughput">&lt;0.1%</span></div>
        <div class="d-box amber">Chargeback rate<br><small>Visa penalty threshold: &gt;1% of volume</small> <span class="d-metric throughput">&lt;1%</span></div>
        <div class="d-box blue">Revenue at risk (open DLQ)<br><small>payments in dead-letter queue = not yet settled</small></div>
        <div class="d-box blue">Reconciliation mismatch count<br><small>target: 0 per day</small></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Alert Thresholds</div>
      <div class="d-flow-v">
        <div class="d-box red">Success rate &lt;99% for 1 min &#8594; immediate page <div class="d-tag red">P0</div></div>
        <div class="d-box red">Latency P99 &gt;3s &#8594; page on-call <div class="d-tag red">P1</div></div>
        <div class="d-box red">Fraud rate &gt;0.5% &#8594; review + possible traffic block <div class="d-tag red">P1</div></div>
        <div class="d-box amber">Chargeback rate &gt;0.5% &#8594; merchant risk review <div class="d-tag amber">P2</div></div>
        <div class="d-box amber">DLQ depth &gt;0 &#8594; alert + auto-retry <div class="d-tag amber">P2</div></div>
        <div class="d-box blue">Reconciliation mismatch &gt;0 &#8594; human review queue</div>
      </div>
    </div>
    <div class="d-group">
      <div class="d-group-title">SLA Summary</div>
      <div class="d-flow-v">
        <div class="d-box green">Uptime: 99.999% (5.26 min/year budget)</div>
        <div class="d-box green">Latency: &lt;2s P99 payment processing</div>
        <div class="d-box green">Durability: zero data loss after payment ack</div>
        <div class="d-box green">Reconciliation: all settlements within 24h</div>
      </div>
    </div>
  </div>
</div>`,
	})
}

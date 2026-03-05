package diagrams

func registerPatterns(r *Registry) {
	// -------------------------------------------------------
	// Agent Tools
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "pat-react-agent-loop",
		Title:       "ReAct Agent Loop",
		Description: "The ReAct (Reasoning + Acting) agent loop showing the cycle of reason, act, observe, and terminate",
		ContentFile: "patterns/agent-tools",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="The initial user query is parsed and passed to the agent. Complex queries may need decomposition before entering the loop."><span class="d-step">1</span> User Query</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Agent Loop (repeat until done)</div>
    <div class="d-flow-v">
      <div class="d-box green" data-tip="LLM generates a chain-of-thought reasoning trace. This is where the model decides WHICH tool to call and WHY. Cost: 1 LLM call per iteration (~300-800ms)."><span class="d-step">2</span> REASON - Think about what to do next <span class="d-metric latency">300-800ms</span></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="Agent executes exactly ONE tool call per iteration. Tool choice is constrained by the available tool schema. Best practice: keep tool descriptions under 200 tokens each."><span class="d-step">3</span> ACT - Call a tool (e.g. web_search) <span class="d-metric latency">50-5000ms</span></div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box indigo" data-tip="Tool output is appended to the conversation context. Truncate large outputs to avoid context window overflow. Typical limit: keep observation under 4K tokens."><span class="d-step">4</span> OBSERVE - Read tool result</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-branch">
        <div class="d-branch-arm">
          <div class="d-label">Task complete?</div>
          <div class="d-box red" data-tip="Loop continues. Best practice: set a max iteration limit (3-5) to prevent runaway costs and infinite loops.">No &#8594; loop back to REASON</div>
        </div>
        <div class="d-branch-arm">
          <div class="d-label">Task complete?</div>
          <div class="d-box green" data-tip="Agent decides it has enough information to answer. The stop condition can be explicit (tool call 'finish') or implicit (model generates final answer)."><span class="d-status active"></span> Yes &#8595;</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Final answer is formatted and returned. Include source attribution from tool results. Total latency = N iterations x (LLM + tool) time."><span class="d-step">5</span> FINAL ANSWER - Formatted response to user</div>
  <div class="d-caption">Each loop iteration costs one LLM call + one tool execution. Typical agent completes in 2-4 iterations. Set max_iterations to prevent runaway costs.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pat-agent-architecture-patterns",
		Title:       "Agent Architecture Patterns",
		Description: "Comparison of simple agent, router agent, and multi-agent architectures",
		ContentFile: "patterns/agent-tools",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">1. Simple Agent</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Single user interface. All queries go to one agent regardless of complexity.">User</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="One LLM with a ReAct loop and shared tool set. Pros: simple to build and debug. Cons: context window fills fast, can't specialize tool prompts.">Agent Loop</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="All tools share one namespace. Typical limit: 10-20 tools before LLM accuracy degrades on tool selection.">Tools</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">2. Router Agent</div>
      <div class="d-flow-v">
        <div class="d-box blue">User</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Lightweight classifier that routes to specialist agents. One LLM call for routing (~200ms). Pros: each sub-agent has focused context. Cons: misrouting loses a full round-trip.">Router Agent <span class="d-metric latency">~200ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-flow">
          <div class="d-box purple" data-tip="Specialized for code generation. Has code-specific tools (exec, lint, test). Isolated context window.">Code Agent</div>
          <div class="d-box purple" data-tip="Specialized for database queries. Has schema access, SQL execution, result formatting.">SQL Agent</div>
          <div class="d-box purple" data-tip="Specialized for web tasks. Has browser, search, scraping tools. Sandboxed for security.">Web Agent</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">3. Multi-Agent</div>
      <div class="d-flow-v">
        <div class="d-box blue">User</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Orchestrator that decomposes tasks and delegates to workers. Pros: parallel execution, separation of concerns. Cons: complex coordination, higher total cost, harder to debug.">Manager Agent</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-flow">
          <div class="d-box indigo" data-tip="Writes code based on the plan. Passes output to Test Agent for validation.">Coder Agent</div>
          <div class="d-box indigo" data-tip="Runs tests against Coder output. Reports failures back to Manager for re-delegation.">Test Agent</div>
          <div class="d-box indigo" data-tip="Reviews code quality, security, and style. Can request changes from Coder Agent via Manager.">Review Agent</div>
        </div>
      </div>
    </div>
  </div>
</div>
<div class="d-legend">
  <span class="d-box blue" style="display:inline-block;padding:2px 8px;font-size:0.75rem;">User</span> Entry point &nbsp;
  <span class="d-box green" style="display:inline-block;padding:2px 8px;font-size:0.75rem;">Orchestrator</span> Routing/coordination &nbsp;
  <span class="d-box purple" style="display:inline-block;padding:2px 8px;font-size:0.75rem;">Specialist</span> Focused sub-agent &nbsp;
  <span class="d-box indigo" style="display:inline-block;padding:2px 8px;font-size:0.75rem;">Worker</span> Parallel execution
</div>
<div class="d-caption">Simple: best for &lt;10 tools. Router: best for distinct task categories. Multi-Agent: best for complex workflows requiring parallel execution and review cycles.</div>`,
	})

	// -------------------------------------------------------
	// Embeddings & Vector Search
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "pat-embedding-space",
		Title:       "Embedding Space (Simplified to 2D)",
		Description: "2D visualization of how similar concepts cluster together in embedding space",
		ContentFile: "patterns/embeddings-vector-search",
		Type:        TypeHTML,
		HTML: `<div class="d-coord" style="height: 320px;">
  <div class="d-coord-axis-label" style="position: absolute; left: 50%; bottom: 0;">Meaning axis 1</div>
  <div class="d-coord-axis-label" style="position: absolute; left: 0; top: 40%; transform: rotate(-90deg) translateX(-50%);">Meaning axis 2</div>
  <div class="d-coord-point" style="left: 40%; top: 12%; background: var(--blue);" data-tip="'refund policy' — semantic cluster: returns/refunds. Cosine similarity to 'return process': 0.92. Embedding dimension: 1536 (text-embedding-3-small)."></div>
  <div class="d-coord-label" style="left: 28%; top: 5%;">refund policy</div>
  <div class="d-coord-point" style="left: 55%; top: 10%; background: var(--blue);" data-tip="'return process' — nearest neighbor to 'refund policy' at distance 0.08. These would be retrieved together in a top-k search."></div>
  <div class="d-coord-label" style="left: 52%; top: 3%;">return process</div>
  <div class="d-coord-point" style="left: 45%; top: 18%; background: var(--blue);" data-tip="'money back guarantee' — same semantic cluster despite different wording. This is why embeddings beat keyword search: synonyms cluster naturally."></div>
  <div class="d-coord-label" style="left: 32%; top: 22%;">money back guarantee</div>
  <div class="d-coord-point" style="left: 68%; top: 38%; background: var(--green);" data-tip="'product specs' — different cluster from refunds. Cosine distance to 'refund policy': 0.61. Well-separated in vector space."></div>
  <div class="d-coord-label" style="left: 62%; top: 32%;">product specs</div>
  <div class="d-coord-point" style="left: 63%; top: 45%; background: var(--green);" data-tip="'item dimensions' — clusters with 'product specs'. Physical product attributes form their own semantic neighborhood."></div>
  <div class="d-coord-label" style="left: 56%; top: 49%;">item dimensions</div>
  <div class="d-coord-point" style="left: 15%; top: 72%; background: var(--amber);" data-tip="'weather forecast' — completely unrelated to e-commerce. Cosine to 'refund policy': 0.15 (near orthogonal). Would never appear in same retrieval set."></div>
  <div class="d-coord-label" style="left: 8%; top: 66%;">weather forecast</div>
  <div class="d-coord-point" style="left: 22%; top: 78%; background: var(--amber);" data-tip="'temperature today' — clusters with 'weather forecast'. Demonstrates that embeddings capture meaning, not surface form."></div>
  <div class="d-coord-label" style="left: 16%; top: 82%;">temperature today</div>
</div>
<div class="d-label" style="text-align: center; margin-top: 0.5rem;">
  Similar concepts cluster together. cosine("refund policy", "return process") = 0.92 | cosine("refund policy", "weather forecast") = 0.15
</div>
<div class="d-caption">Real embeddings are 768-3072 dimensions projected to 2D here. Cosine similarity &gt; 0.8 = strong match. &lt; 0.3 = unrelated. Clustering quality depends on model choice (e.g. text-embedding-3-large outperforms ada-002).</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pat-hnsw-index-structure",
		Title:       "HNSW Index Structure",
		Description: "Hierarchical Navigable Small World graph with sparse top layers and dense bottom layer",
		ContentFile: "patterns/embeddings-vector-search",
		Type:        TypeHTML,
		HTML: `<div class="d-graph">
  <div class="d-graph-layer" data-tip="Layer 2: Entry point for search. Only ~1% of nodes appear here. Long-range connections enable O(log N) jumps across the graph. Search starts at a random entry node on this layer.">
    <div class="d-graph-layer-label">Layer 2 (sparse) — entry point, long jumps</div>
    <div class="d-graph-nodes">
      <div class="d-graph-node" style="background: var(--purple);" data-tip="Node A promoted to top layer. Probability of promotion: 1/M per layer (M=16 default). Fewer nodes = faster coarse search.">A</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--purple);" data-tip="Node B also at top layer. Edge A-B covers large distance in vector space — enables fast global navigation.">B</div>
    </div>
  </div>
  <div class="d-graph-layer" data-tip="Layer 1: Medium density. More nodes, shorter edges. Search descends here after finding approximate region in Layer 2. Each node connects to M neighbors.">
    <div class="d-graph-layer-label">Layer 1 (medium) — refine neighborhood</div>
    <div class="d-graph-nodes">
      <div class="d-graph-node" style="background: var(--blue);">A</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--blue);">B</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--blue);">C</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--blue);">D</div>
    </div>
  </div>
  <div class="d-graph-layer" data-tip="Layer 0: All nodes present. Dense connections for precise nearest-neighbor search. Final results come from this layer. This is where ef_search controls recall quality.">
    <div class="d-graph-layer-label">Layer 0 (dense) — all nodes, precise results</div>
    <div class="d-graph-nodes">
      <div class="d-graph-node" style="background: var(--green);">A</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--green);">E</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--green);">B</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--green);">F</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--green);">C</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--green);">G</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--green);">D</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--green);">H</div>
    </div>
  </div>
</div>
<div class="d-label" style="text-align: center; margin-top: 0.5rem;">
  Search starts at top layer (few nodes, long jumps). Each layer adds more nodes with shorter connections.<br>
  Tuning: M (connections/node), ef_construction (index quality), ef_search (query recall)
</div>
<div class="d-caption">Search complexity: O(log N) average. For 1M vectors with M=16, ef_search=64: ~1-5ms query latency. Index build time: O(N log N). Memory: ~1.5x raw vectors. Trade-off: higher M = better recall but more memory and slower inserts.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pat-embedding-pipeline",
		Title:       "Production Embedding Pipeline",
		Description: "Indexing and query paths for a production embedding pipeline with document processing and vector storage",
		ContentFile: "patterns/embeddings-vector-search",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Indexing Path</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Source documents can be any format. Use connectors (Unstructured, LlamaIndex) for PDF/HTML extraction. Typical batch: 10K-100K documents."><span class="d-step">1</span> Documents Source (S3, API, DB)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Chunk strategy matters: 256-512 tokens per chunk with 50-token overlap. Too small = lost context. Too large = diluted relevance. Use recursive text splitter for best results."><span class="d-step">2</span> Document Processor (chunk, clean) <span class="d-metric throughput">~500 docs/s</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Batch embedding is 10-50x cheaper than single calls. OpenAI: $0.02/1M tokens. Self-hosted (e5-large): ~$0.002/1M tokens but requires GPU. Batch size: 100-2000 texts per API call."><span class="d-step">3</span> Embedding API (batch) <span class="d-metric latency">~20ms/chunk</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Use UPSERT to handle re-indexing. Store chunk text + metadata (source, page, timestamp) alongside the vector. Create HNSW index after bulk load for faster build."><span class="d-step">4</span> pgvector INSERT (upsert) <span class="d-metric throughput">~5K rows/s</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Query Path</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="User query is typically short (5-20 tokens). Consider query expansion or HyDE (Hypothetical Document Embeddings) for better retrieval."><span class="d-step">A</span> User Query</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="CRITICAL: Must use the EXACT same model as indexing. Mixing models (e.g. ada-002 index + text-embedding-3 query) produces garbage results."><span class="d-step">B</span> Embed Query (same model!) <span class="d-metric latency">~20ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="HNSW index gives approximate nearest neighbors. Set probes=10 for IVFFlat or ef_search=64 for HNSW. Top-k typically 10-20 candidates."><span class="d-step">C</span> pgvector search (HNSW) <span class="d-metric latency">~5ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Return chunks with metadata for citation. Filter by metadata (date, source) in WHERE clause before vector search for efficiency."><span class="d-step">D</span> Results + metadata <span class="d-status active"></span></div>
      </div>
    </div>
  </div>
</div>
<div class="d-label" style="text-align: center; margin-top: 0.5rem;">
  Key: Query embedding model MUST match the document embedding model exactly.
</div>
<div class="d-caption">Total indexing: ~2 hours for 1M documents. Query path end-to-end: ~25ms (embed + search). Cost at scale: embedding dominates indexing cost, pgvector hosting dominates operational cost.</div>`,
	})

	// -------------------------------------------------------
	// Guardrails
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "pat-guardrail-pipeline",
		Title:       "Guardrail Pipeline Architecture",
		Description: "Input and output guardrail pipeline with cheapest-first ordering around the LLM",
		ContentFile: "patterns/guardrails",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="Raw user input before any validation. Could contain prompt injections, PII, or off-topic content. Never pass directly to LLM."><span class="d-step">1</span> User Input</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-group">
      <div class="d-group-title">Input Guardrails (cheapest first)</div>
      <div class="d-flow-v">
        <div class="d-box green" data-tip="Simple string length check. Reject &gt;4K chars. Cost: ~0ms. Blocks prompt-stuffing attacks and accidental file pastes."><span class="d-step">2a</span> Length check <span class="d-metric latency">&lt;1ms</span></div>
        <div class="d-box green" data-tip="Regex + NER model to detect SSN, credit cards, emails, phone numbers. Use Presidio or AWS Comprehend. Redact before LLM sees it."><span class="d-step">2b</span> PII detection <span class="d-metric latency">~5ms</span></div>
        <div class="d-box green" data-tip="Detect prompt injection attempts: 'ignore previous instructions', role-play attacks, encoded payloads. Use a small classifier model (DistilBERT fine-tuned)."><span class="d-step">2c</span> Injection detect <span class="d-metric latency">~15ms</span></div>
        <div class="d-box green" data-tip="Classify if input is on-topic for your application. Use an embedding similarity check against allowed topic clusters. Threshold: cosine &gt; 0.5."><span class="d-step">2d</span> Topic filter <span class="d-metric latency">~25ms</span></div>
      </div>
    </div>
    <div class="d-flow-v">
      <div class="d-label">Fail?</div>
      <div class="d-box red" data-tip="Fail fast: return specific error message for each guard type. Log the failure reason for monitoring. Never expose internal guard details to user.">Return error / ask to rephrase</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; Pass</div>
  <div class="d-box indigo" data-tip="LLM processes validated input. System prompt includes output format constraints. Total input guard latency: ~45ms — negligible vs LLM inference."><span class="d-step">3</span> LLM (your model) <span class="d-metric latency">500-3000ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-group">
      <div class="d-group-title">Output Guardrails (cheapest first)</div>
      <div class="d-flow-v">
        <div class="d-box amber" data-tip="Validate JSON schema, required fields, max length. Use Pydantic or JSON Schema validation. Catches malformed LLM outputs before downstream processing."><span class="d-step">4a</span> Format check <span class="d-metric latency">&lt;1ms</span></div>
        <div class="d-box amber" data-tip="LLM may echo or hallucinate PII. Re-run PII detection on output. Scrub any detected entities with [REDACTED] placeholders."><span class="d-step">4b</span> PII scrubbing <span class="d-metric latency">~5ms</span></div>
        <div class="d-box amber" data-tip="Check for harmful, toxic, or off-brand content. Use a safety classifier (OpenAI Moderation API is free, or Llama Guard for self-hosted)."><span class="d-step">4c</span> Content safety <span class="d-metric latency">~30ms</span></div>
        <div class="d-box amber" data-tip="Most expensive guard: cross-reference claims against retrieved sources. Use NLI model (Natural Language Inference) to check entailment. Critical for RAG systems."><span class="d-step">4d</span> Hallucination check <span class="d-metric latency">~100ms</span></div>
      </div>
    </div>
    <div class="d-flow-v">
      <div class="d-label">Fail?</div>
      <div class="d-box red" data-tip="Strategy depends on failure type: format fail = retry with stricter prompt, safety fail = block entirely, hallucination = retry with higher temperature 0.">Retry / fallback / block</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; Pass</div>
  <div class="d-box purple" data-tip="Fully validated response. Total overhead: ~180ms for all guards combined — adds ~5-10% to typical LLM response time."><span class="d-step">5</span> User Response <span class="d-status active"></span></div>
  <div class="d-caption">Order guards cheapest-first to fail fast. Total guard overhead: ~180ms (~5% of LLM latency). Run input guards in parallel where possible. Hallucination check is most expensive but most impactful for trust.</div>
</div>`,
	})

	// -------------------------------------------------------
	// Prompt Chaining
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "pat-prompt-chain-pipeline",
		Title:       "Prompt Chain Pipeline",
		Description: "Sequential prompt chain with extract, classify, generate, and validate steps plus retry logic",
		ContentFile: "patterns/prompt-chaining",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="Raw input document — could be email, ticket, form submission. Each chain step sees only its predecessor's output, NOT the original input (unless explicitly passed)."><span class="d-step">1</span> Input Document</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-flow-v">
      <div class="d-box green" data-tip="Use a focused prompt: 'Extract these 5 fields from the text.' Structured output (JSON mode) ensures downstream steps can parse reliably. Use a cheaper model (GPT-4o-mini) for extraction."><span class="d-step">2</span> EXTRACT key facts <span class="d-metric latency">~400ms</span></div>
      <div class="d-label">"key: value ..."</div>
    </div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-flow-v">
      <div class="d-box green" data-tip="Classification from extracted facts, not raw text. Simpler input = more accurate classification. Use constrained output (enum of categories) to prevent hallucinated classes."><span class="d-step">3</span> CLASSIFY category <span class="d-metric latency">~200ms</span></div>
      <div class="d-label">"category: X"</div>
    </div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-flow-v">
      <div class="d-box indigo" data-tip="Most expensive step — use the best model here (GPT-4o, Claude). Prompt includes extracted facts + category as context. This is where quality matters most."><span class="d-step">4</span> GENERATE response <span class="d-metric latency">~1500ms</span></div>
      <div class="d-label">"Dear user..."</div>
    </div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-flow-v">
      <div class="d-box amber" data-tip="Quality gate: check tone, factual consistency with extracted facts, format compliance. Use a different model than generator to avoid self-confirmation bias."><span class="d-step">5</span> VALIDATE quality <span class="d-metric latency">~300ms</span></div>
      <div class="d-label">"PASS / FAIL"</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box red" data-tip="Retry includes the validation error as feedback. 'Your response failed because: [reason]. Please regenerate.' Max 2 retries to cap cost. After 2 failures, escalate to human.">FAIL &#8594; Retry Step 3 (max 2x)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green" data-tip="Successful chain completion. Log all intermediate outputs for debugging. Total cost: ~$0.005 per chain execution with GPT-4o-mini for steps 2-3, GPT-4o for step 4."><span class="d-status active"></span> PASS &#8594; Return to User</div>
    </div>
  </div>
  <div class="d-caption">Total chain latency: ~2.4s (sequential). Each step uses focused prompts for higher accuracy than a single monolithic prompt. Cost optimization: use cheaper models for extraction/classification, premium models for generation.</div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pat-gate-pattern-flow",
		Title:       "Gate Pattern Flow",
		Description: "Quality gate pattern with pass/fail branching and retry logic between chain steps",
		ContentFile: "patterns/prompt-chaining",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="Output from the preceding chain step. Should be structured (JSON) to make validation deterministic. Include schema version for forward compatibility.">Step N Output</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber" data-tip="Gate validator can be: (1) LLM-as-judge with rubric, (2) deterministic code checks (regex, schema), or (3) hybrid. Deterministic checks are cheaper and more reliable — use LLM-as-judge only for subjective quality.">GATE Validator <span class="d-metric latency">~50-300ms</span></div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box green" data-tip="Output meets all quality criteria. Pass the validated output (not raw) to the next step. Log the gate decision for audit trail."><span class="d-status active"></span> PASS</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box green" data-tip="Next step receives validated output. Each step in the chain should be idempotent — same input always produces same-quality output.">Step N+1</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box red" data-tip="Gate returns structured failure reason: {passed: false, reason: 'tone too formal', suggestion: 'use casual language'}. This feedback is critical for effective retry.">FAIL</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber" data-tip="Append gate feedback to the retry prompt: 'Previous attempt failed: [reason]. Please fix: [suggestion].' Each retry costs one LLM call. Exponential backoff not needed — LLM retries are deterministic.">Retry Step N with error feedback (max 2-3 retries)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box red" data-tip="After max retries: (1) return best attempt with confidence score, (2) route to human review queue, or (3) use a pre-written fallback template. Never silently fail.">Still FAIL &#8594; Abort chain or use fallback</div>
    </div>
  </div>
  <div class="d-caption">Gates add ~10-20% latency but catch 80%+ of quality issues before they propagate downstream. Use deterministic gates where possible (schema validation, regex) and reserve LLM-as-judge for subjective criteria.</div>
</div>`,
	})

	// -------------------------------------------------------
	// RAG (Retrieval-Augmented Generation)
	// -------------------------------------------------------

	r.Register(&Diagram{
		Slug:        "pat-rag-pipeline",
		Title:       "RAG Pipeline Architecture",
		Description: "Offline indexing and online query phases of a RAG pipeline with chunking, embedding, retrieval, and generation",
		ContentFile: "patterns/rag",
		Type:        TypeHTML,
		HTML: `<div class="d-cols">
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">OFFLINE: Indexing Phase</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="Source documents in any format. Use document loaders (LangChain, LlamaIndex) for PDF/HTML extraction. Track document versions for incremental re-indexing."><span class="d-step">1</span> Documents (PDF, HTML, Markdown)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Chunk size is the #1 tuning knob. Start with 512 tokens, 50-token overlap. Use recursive character splitter to respect paragraph boundaries. Too small = lost context, too large = noise in retrieval."><span class="d-step">2</span> Chunk Documents <span class="d-metric latency">~2ms/doc</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Batch embed for cost efficiency. OpenAI text-embedding-3-small: $0.02/1M tokens, 1536 dims. Self-hosted e5-large: cheaper but requires GPU. Batch size 100-2000 per API call."><span class="d-step">3</span> Embed Chunks <span class="d-metric latency">~20ms/chunk</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Store vector + original text + metadata (source URL, page number, timestamp). Use pgvector, Pinecone, or Weaviate. Build HNSW index after bulk load. Plan for ~4KB per chunk (vector + metadata)."><span class="d-step">4</span> Vector Database (store) <span class="d-metric throughput">~5K rows/s</span></div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">ONLINE: Query Phase</div>
      <div class="d-flow-v">
        <div class="d-box blue" data-tip="User's natural language question. Consider query rewriting: expand abbreviations, resolve pronouns from conversation history."><span class="d-step">A</span> User Query</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Must use the SAME embedding model as indexing. Single query embedding is fast. Consider HyDE: generate a hypothetical answer first, then embed that for better retrieval."><span class="d-step">B</span> Embed Query <span class="d-metric latency">~20ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Retrieve top-k candidates (k=10-20). Use hybrid search: combine dense vector search with BM25 keyword search for best recall. Metadata filters (date range, source) applied here."><span class="d-step">C</span> Vector Search (top-k) <span class="d-metric latency">~5ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber" data-tip="Cross-encoder reranking improves precision by 15-30%. Scores each (query, chunk) pair jointly. Use Cohere Rerank or a cross-encoder model. Worth the latency cost for quality-sensitive applications."><span class="d-step">D</span> Rerank (optional) <span class="d-metric latency">~100ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green" data-tip="Inject top 3-5 chunks into the system prompt. Format: 'Context:\n[chunk1]\n[chunk2]\n\nAnswer the question based on the context above.' Total context: ~2K tokens."><span class="d-step">E</span> Augment Prompt w/ chunks</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo" data-tip="LLM generates answer grounded in retrieved context. Use temperature=0 for factual accuracy. Include 'If the context doesn't contain the answer, say so' to reduce hallucination."><span class="d-step">F</span> LLM Generate <span class="d-metric latency">500-2000ms</span></div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple" data-tip="Return answer with source citations. Link each claim to its source chunk. Users can click through to verify. This builds trust and enables feedback loops."><span class="d-step">G</span> Response + Sources <span class="d-status active"></span></div>
      </div>
    </div>
  </div>
</div>
<div class="d-caption">Indexing is a one-time batch job (~2hrs for 1M docs). Query path end-to-end: ~700ms without rerank, ~800ms with. LLM generation dominates latency. Optimize retrieval quality first, then latency.</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pat-two-stage-retrieval-reranking",
		Title:       "Two-Stage Retrieval + Reranking",
		Description: "Two-stage retrieval pipeline with broad bi-encoder search followed by precise cross-encoder reranking",
		ContentFile: "patterns/rag",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue" data-tip="User query is embedded once and used for both dense retrieval and as input to the cross-encoder. Short queries (5-10 tokens) benefit most from reranking."><span class="d-step">1</span> Query</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box green" data-tip="Bi-encoder: query and documents are embedded independently. Fast because document embeddings are precomputed. Combine dense (semantic) + BM25 (keyword) for hybrid search — typically 10-15% better recall than either alone."><span class="d-step">2</span> Stage 1: Retrieve<br>Bi-encoder search (dense + BM25)</div>
    <div class="d-label">Fast, broad retrieval<br>top-50 candidates <span class="d-metric latency">~10ms</span> <span class="d-metric throughput">1000 QPS</span></div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box amber" data-tip="Cross-encoder: jointly encodes (query, document) pairs for deeper interaction. 10-50x slower than bi-encoder but 15-30% more precise. Only feasible on small candidate sets (50-100). Models: Cohere Rerank, ms-marco-MiniLM."><span class="d-step">3</span> Stage 2: Rerank<br>Cross-encoder (query, doc pairs)</div>
    <div class="d-label">Precise reordering<br>top-5 from 50 <span class="d-metric latency">~100ms</span></div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple" data-tip="Final top-5 chunks are highest quality. Total retrieval: ~110ms. These chunks become the 'context' window for LLM generation. Each chunk includes source attribution for citations."><span class="d-step">4</span> Top-5 chunks injected into LLM prompt <span class="d-status active"></span></div>
  <div class="d-caption">Two-stage retrieval trades ~100ms extra latency for 15-30% precision improvement. Bi-encoder handles scale (millions of docs), cross-encoder handles quality (top-50 candidates). Always measure Recall@5 and NDCG@5 to validate.</div>
</div>`,
	})
}

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
  <div class="d-box blue">User Query</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-group">
    <div class="d-group-title">Agent Loop (repeat until done)</div>
    <div class="d-flow-v">
      <div class="d-box green">REASON - Think about what to do next</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">ACT - Call a tool (e.g. web_search)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box indigo">OBSERVE - Read tool result</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-branch">
        <div class="d-branch-arm">
          <div class="d-label">Task complete?</div>
          <div class="d-box red">No &#8594; loop back to REASON</div>
        </div>
        <div class="d-branch-arm">
          <div class="d-label">Task complete?</div>
          <div class="d-box green">Yes &#8595;</div>
        </div>
      </div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">FINAL ANSWER - Formatted response to user</div>
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
        <div class="d-box blue">User</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Agent Loop</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Tools</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">2. Router Agent</div>
      <div class="d-flow-v">
        <div class="d-box blue">User</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Router Agent</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-flow">
          <div class="d-box purple">Code Agent</div>
          <div class="d-box purple">SQL Agent</div>
          <div class="d-box purple">Web Agent</div>
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
        <div class="d-box green">Manager Agent</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-flow">
          <div class="d-box indigo">Coder Agent</div>
          <div class="d-box indigo">Test Agent</div>
          <div class="d-box indigo">Review Agent</div>
        </div>
      </div>
    </div>
  </div>
</div>`,
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
  <div class="d-coord-point" style="left: 40%; top: 12%; background: var(--blue);"></div>
  <div class="d-coord-label" style="left: 28%; top: 5%;">refund policy</div>
  <div class="d-coord-point" style="left: 55%; top: 10%; background: var(--blue);"></div>
  <div class="d-coord-label" style="left: 52%; top: 3%;">return process</div>
  <div class="d-coord-point" style="left: 45%; top: 18%; background: var(--blue);"></div>
  <div class="d-coord-label" style="left: 32%; top: 22%;">money back guarantee</div>
  <div class="d-coord-point" style="left: 68%; top: 38%; background: var(--green);"></div>
  <div class="d-coord-label" style="left: 62%; top: 32%;">product specs</div>
  <div class="d-coord-point" style="left: 63%; top: 45%; background: var(--green);"></div>
  <div class="d-coord-label" style="left: 56%; top: 49%;">item dimensions</div>
  <div class="d-coord-point" style="left: 15%; top: 72%; background: var(--amber);"></div>
  <div class="d-coord-label" style="left: 8%; top: 66%;">weather forecast</div>
  <div class="d-coord-point" style="left: 22%; top: 78%; background: var(--amber);"></div>
  <div class="d-coord-label" style="left: 16%; top: 82%;">temperature today</div>
</div>
<div class="d-label" style="text-align: center; margin-top: 0.5rem;">
  Similar concepts cluster together. cosine("refund policy", "return process") = 0.92 | cosine("refund policy", "weather forecast") = 0.15
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pat-hnsw-index-structure",
		Title:       "HNSW Index Structure",
		Description: "Hierarchical Navigable Small World graph with sparse top layers and dense bottom layer",
		ContentFile: "patterns/embeddings-vector-search",
		Type:        TypeHTML,
		HTML: `<div class="d-graph">
  <div class="d-graph-layer">
    <div class="d-graph-layer-label">Layer 2 (sparse)</div>
    <div class="d-graph-nodes">
      <div class="d-graph-node" style="background: var(--purple);">A</div>
      <div class="d-graph-edge"></div>
      <div class="d-graph-node" style="background: var(--purple);">B</div>
    </div>
  </div>
  <div class="d-graph-layer">
    <div class="d-graph-layer-label">Layer 1 (medium)</div>
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
  <div class="d-graph-layer">
    <div class="d-graph-layer-label">Layer 0 (dense)</div>
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
</div>`,
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
        <div class="d-box blue">Documents Source (S3, API, DB)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Document Processor (chunk, clean)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Embedding API (batch) - OpenAI / self-host</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">pgvector INSERT (upsert)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">Query Path</div>
      <div class="d-flow-v">
        <div class="d-box blue">User Query</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Embed Query (same model!)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">pgvector search (HNSW)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Results + metadata</div>
      </div>
    </div>
  </div>
</div>
<div class="d-label" style="text-align: center; margin-top: 0.5rem;">
  Key: Query embedding model MUST match the document embedding model exactly.
</div>`,
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
  <div class="d-box blue">User Input</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-group">
      <div class="d-group-title">Input Guardrails (cheapest first)</div>
      <div class="d-flow-v">
        <div class="d-box green">1. Length check</div>
        <div class="d-box green">2. PII detection</div>
        <div class="d-box green">3. Injection detect</div>
        <div class="d-box green">4. Topic filter</div>
      </div>
    </div>
    <div class="d-flow-v">
      <div class="d-label">Fail?</div>
      <div class="d-box red">Return error / ask to rephrase</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; Pass</div>
  <div class="d-box indigo">LLM (your model)</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-group">
      <div class="d-group-title">Output Guardrails (cheapest first)</div>
      <div class="d-flow-v">
        <div class="d-box amber">1. Format check</div>
        <div class="d-box amber">2. PII scrubbing</div>
        <div class="d-box amber">3. Content safety</div>
        <div class="d-box amber">4. Hallucination check</div>
      </div>
    </div>
    <div class="d-flow-v">
      <div class="d-label">Fail?</div>
      <div class="d-box red">Retry / fallback / block</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595; Pass</div>
  <div class="d-box purple">User Response</div>
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
  <div class="d-box blue">Input Document</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-flow">
    <div class="d-flow-v">
      <div class="d-box green">Step 1: EXTRACT key facts</div>
      <div class="d-label">"key: value ..."</div>
    </div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-flow-v">
      <div class="d-box green">Step 2: CLASSIFY category</div>
      <div class="d-label">"category: X"</div>
    </div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-flow-v">
      <div class="d-box indigo">Step 3: GENERATE response</div>
      <div class="d-label">"Dear user..."</div>
    </div>
    <div class="d-arrow">&#8594;</div>
    <div class="d-flow-v">
      <div class="d-box amber">Step 4: VALIDATE quality</div>
      <div class="d-label">"PASS / FAIL"</div>
    </div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box red">FAIL &#8594; Retry Step 3 (max 2x)</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box green">PASS &#8594; Return to User</div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pat-gate-pattern-flow",
		Title:       "Gate Pattern Flow",
		Description: "Quality gate pattern with pass/fail branching and retry logic between chain steps",
		ContentFile: "patterns/prompt-chaining",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Step N Output</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box amber">GATE Validator</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-branch">
    <div class="d-branch-arm">
      <div class="d-box green">PASS</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box green">Step N+1</div>
    </div>
    <div class="d-branch-arm">
      <div class="d-box red">FAIL</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box amber">Retry Step N with error feedback (max 2-3 retries)</div>
      <div class="d-arrow-down">&#8595;</div>
      <div class="d-box red">Still FAIL &#8594; Abort chain or use fallback</div>
    </div>
  </div>
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
        <div class="d-box blue">Documents (PDF, HTML, Markdown)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Chunk Documents</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Embed Chunks</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Vector Database (store)</div>
      </div>
    </div>
  </div>
  <div class="d-col">
    <div class="d-group">
      <div class="d-group-title">ONLINE: Query Phase</div>
      <div class="d-flow-v">
        <div class="d-box blue">User Query</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Embed Query</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Vector Search (top-k)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box amber">Rerank (optional)</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box green">Augment Prompt w/ chunks</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box indigo">LLM Generate</div>
        <div class="d-arrow-down">&#8595;</div>
        <div class="d-box purple">Response + Sources</div>
      </div>
    </div>
  </div>
</div>`,
	})

	r.Register(&Diagram{
		Slug:        "pat-two-stage-retrieval-reranking",
		Title:       "Two-Stage Retrieval + Reranking",
		Description: "Two-stage retrieval pipeline with broad bi-encoder search followed by precise cross-encoder reranking",
		ContentFile: "patterns/rag",
		Type:        TypeHTML,
		HTML: `<div class="d-flow-v">
  <div class="d-box blue">Query</div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box green">Stage 1: Retrieve<br>Bi-encoder search (dense + BM25)</div>
    <div class="d-label">Fast, broad retrieval<br>top-50 candidates, ~10ms</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-row">
    <div class="d-box amber">Stage 2: Rerank<br>Cross-encoder (query, doc pairs)</div>
    <div class="d-label">Precise reordering<br>top-5 from 50, ~100ms</div>
  </div>
  <div class="d-arrow-down">&#8595;</div>
  <div class="d-box purple">Top-5 chunks injected into LLM prompt</div>
</div>`,
	})
}

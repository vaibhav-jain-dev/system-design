# Content Review: Improvements Needed

> **Review Date:** 2026-03-05
> **Scope:** All content files, diagrams, registry, templates, macros, handlers, CSS
> **Total Issues Found:** 202

---

## Table of Contents

1. [Critical Factual Errors (P0)](#1-critical-factual-errors-p0)
2. [Factual Errors (P1)](#2-factual-errors-p1)
3. [Structural Issues](#3-structural-issues)
4. [Missing Content](#4-missing-content)
5. [Quality Issues](#5-quality-issues)
6. [Infrastructure & Code Bugs](#6-infrastructure--code-bugs)
7. [Security Issues](#7-security-issues)
8. [CSS & UI Issues](#8-css--ui-issues)
9. [Documentation vs Implementation Mismatches](#9-documentation-vs-implementation-mismatches)
10. [Improvement Opportunities](#10-improvement-opportunities)

---

## 1. Critical Factual Errors (P0)

These errors would damage interview credibility if repeated. Fix immediately.

### 1.1 Base62 Encoding Diagram Math Is Entirely Wrong

**File:** `internal/diagrams/algorithms.go` lines 14-51 (diagram `algo-base62-process`)
**Also affects:** `content/algorithms/base62-encoding/index.html` line 41

Every single step of the encoding walkthrough for input `123456789` is incorrect:

| Step | Diagram Claims | Correct Value |
|------|---------------|---------------|
| 1 | 123456789 % 62 = 17 → r | 123456789 % 62 = 33 → H |
| 2 | 1991239 % 62 = 25 → z | 1991238 % 62 = 46 → U |
| 3 | 32116 % 62 = 4 → e | 32116 % 62 = 0 → a |
| 4 | 517 % 62 = 21 → v | 518 % 62 = 22 → w |
| 5 | 8 % 62 = 8 → i (only correct step) | 8 % 62 = 8 → i |

The diagram claims the result is `"ivezr"` but the correct encoding is `"iwaHU"` (or `"iwaUH"` depending on alphabet ordering). The decode verification expression also evaluates incorrectly.

**Fix:** Recompute all division steps with correct arithmetic. Verify decode line matches.

---

### 1.2 Read:Write Ratio Contradiction — URL Shortener

**File:** `content/problems/url-shortener/index.html` lines 42-44 vs diagram `url-nfr-estimates`

The text states **100:1** read-to-write ratio, but the diagram and all QPS math use **10:1**. The entire architecture is sized for ~57K peak RPS (derived from 10:1). If it were truly 100:1, peak reads would be 578,500 RPS — a fundamentally different architecture.

**Fix:** Change the text to 10:1 (matches all the math), or recalculate everything for 100:1.

---

### 1.3 Read:Write Ratio Contradiction — Instagram

**File:** `content/problems/instagram/index.html` line 39 vs diagram `ig-nfr-estimates` (instagram.go line 57)

The NFR diagram states **100:1** read-to-write ratio, but the actual numbers (58K reads vs 6K writes) yield approximately **10:1**. Off by an order of magnitude.

**Fix:** Change diagram to say 10:1 to match the math.

---

### 1.4 DynamoDB Partition Limit: "3M reads/sec" Is 1000x Wrong

**File:** `content/problems/url-shortener/index.html` line 146

States a DynamoDB partition has "3,000 RCU capacity (3M reads/sec max)." A partition supports 3,000 RCU = **3,000 reads/sec** (for items ≤4KB), not 3 million. This is a 1,000x error.

**Fix:** Change "3M reads/sec max" to "3,000 reads/sec max (for ≤4KB items)."

---

### 1.5 DynamoDB "Hard Ceiling" Claim Is Outdated

**File:** `content/fundamentals/storage/dynamodb/index.html` line 287

States "Each DynamoDB partition has a hard ceiling of 1,000 WCU and 3,000 RCU regardless of how much total table capacity you provision." Since **2019**, DynamoDB has **adaptive capacity** which allows individual partitions to burst beyond these baseline limits by borrowing unused capacity from other partitions. Calling it a "hard ceiling" is factually incorrect for modern DynamoDB.

**Fix:** Change to "baseline allocation" and mention adaptive capacity allows bursting beyond these limits.

---

### 1.6 NLB "No Security Groups" Is Outdated

**File:** `content/fundamentals/networking/load-balancing/nlb/index.html` line 208

States "NLB does not have security groups by default, use target SGs." Since **August 2023**, NLB supports security groups directly. This advice is now wrong and could embarrass a candidate.

**Fix:** Update to reflect NLB security group support. Recommend using NLB security groups directly.

---

### 1.7 CloudFront Invalidation Limit: 3,000 vs 1,000

**Files:**
- `content/fundamentals/networking/cdn/index.html` line 106: says "3000 free invalidations/mo" — **WRONG**
- `content/fundamentals/networking/cdn/index.html` line 114: says "Free (1000/mo on CF)" — **CORRECT**
- `content/fundamentals/networking/cdn/cloudfront/index.html` line 360: says "3000 free paths per month" — **WRONG**

AWS provides **1,000 free invalidation paths/month**. Two of three references have the wrong number, and they contradict each other within the same file.

**Fix:** Change all instances to 1,000.

---

### 1.8 Bloom Filter False Positive Rate Claims Are Wrong

**File:** `content/algorithms/bloom-filter/index.html`

- **Line 56:** States k=3, 9.6 bits/element yields FP rate of "about 3.5%." Actual: (1 - e^(-3/9.6))^3 ≈ **1.93%**.
- **Line 83:** States "at 2x expected elements, FP rate roughly quadruples." Actual: FP rate goes from ~1% to ~15.7%, a **~16x increase**, not 4x.

**Fix:** Recalculate with correct formula and update both claims.

---

### 1.9 Consistent Hashing Modulo Math Is Wrong

**File:** `content/algorithms/consistent-hashing/index.html` lines 23-26

States: "Only keys where hash(key) is divisible by both 100 and 99 stay. The probability is 1/lcm(99,100) ≈ 0.01%, so ~99.99% move."

Correct: The fraction staying is 99/lcm(99,100) = 99/9900 = **1%**. So ~**99%** move, not 99.99%. The error is in not counting all 99 valid residue classes.

**Fix:** Change to "~1% of keys stay, ~99% move" with correct math.

---

### 1.10 Snowflake Discord Hint Contradicts Its Own Table

**File:** `content/algorithms/snowflake-id/index.html` line 205

States: "They traded 1 bit from worker ID (512 workers instead of 1024) for ~70 extra years."

But the table at line 199 shows Discord's layout as `42 ts + 10 worker + 12 seq`. With 10 bits, Discord supports **1024 workers**, not 512. Discord gained the extra timestamp bit by using all 64 bits (no sign bit), not by reducing worker bits.

**Fix:** Correct the hint to explain that Discord uses the sign bit for timestamp, giving 42 timestamp bits vs Twitter's 41.

---

## 2. Factual Errors (P1)

### 2.1 URL Shortener: Base62 Example Output Inconsistent

**File:** `content/problems/url-shortener/index.html` line 241 and diagram line 310

Shows `base62_encode(123456789) -> "8M0kX"`. With the standard alphabet (`a-z, A-Z, 0-9`), the result should be `"iwaHU"`. The discrepancy suggests a different alphabet ordering. Either fix the output or explicitly document the alphabet.

---

### 2.2 URL Shortener: `zfill(7)` Pads with '0', Not First Alphabet Char

**File:** `content/problems/url-shortener/index.html` line 232

Text says "Left-pad with the first alphabet character." The code uses `.zfill(7)` which pads with the digit `'0'`, not the first alphabet character `'a'`. These are different characters with different Base62 values.

**Fix:** Either change the text to say "left-pad with zeros" or change the code to pad with the correct character.

---

### 2.3 URL Shortener: 43 Bits vs 62^7 Range Unexplained

**File:** Diagram `url-bit-layout` (url_shortener.go lines 330, 365)

43 bits = 8.8T values, but 62^7 = 3.5T. Values between 3.5T and 8.8T would produce **8-character** Base62 strings. This overflow case is not discussed.

**Fix:** Add a note that values must be capped at 62^7, or explain the overflow handling.

---

### 2.4 URL Shortener: MD5 Truncation Inconsistency

**File:** `content/problems/url-shortener/index.html` lines 203-204 vs deepQA line 746

The comparison table describes MD5 truncation in Base62 space (3.5T), but the deepQA describes hex truncation (16^7 = 268M). These are two different truncation methods described inconsistently.

---

### 2.5 URL Shortener: Hystrix Is Deprecated

**File:** `content/problems/url-shortener/index.html` line 386

References "Hystrix on DynamoDB calls." Netflix Hystrix entered maintenance mode in **2018**. The modern replacement is **Resilience4j** (which is mentioned elsewhere in the file at line 405).

**Fix:** Replace Hystrix reference with Resilience4j throughout.

---

### 2.6 Rate Limiter: Sliding Window "0.003%" Accuracy Is Unsubstantiated

**File:** `content/problems/rate-limiter/index.html` lines 125, 137, 439, 451

The claim that sliding window counter achieves "within 0.003% accuracy" is repeated 4 times. This specific number appears fabricated — no published source (Cloudflare or academic) cites this figure. The actual error bound depends heavily on traffic distribution.

**Fix:** Replace with a qualified statement like "typically within a few percent of the true count for uniform traffic" or cite a specific source.

---

### 2.7 Rate Limiter: Cost Contradiction ($580/mo vs $1,164/mo)

**File:** `content/problems/rate-limiter/index.html` lines 403-409 and 428

- Diagram `rl-cost-breakdown` says ElastiCache = ~$550/mo, total ~$580/mo
- Think block says `r6g.large: ~$194/month per node, 6 nodes = $1,164/month`
- Actual AWS pricing: cache.r6g.large ≈ $97/month/node → 6 nodes = ~$582/mo

The $194/node figure is approximately 2x the actual price. The ~$580 total is approximately correct.

**Fix:** Correct the per-node price to ~$97 and update the think block math.

---

### 2.8 Rate Limiter: "2,500 Seconds of Cumulative Latency" Is Misleading

**File:** `content/problems/rate-limiter/index.html` line 45

"At 500K req/sec, even a 5ms Redis call would add 2,500 seconds of cumulative latency per second" — "cumulative latency per second" is not a meaningful metric. Each request adds 5ms individually.

**Fix:** Reframe as "requires 2,500 CPU-seconds of parallel processing capacity per second."

---

### 2.9 Instagram: Storage Variant Multiplier Overstated (4x → 1.24x)

**File:** `content/problems/instagram/index.html` line 28

Claims 4 resize variants at 4x storage (73PB × 4 = ~290PB). But the 4 variants (150px to 1080px) total ~475KB combined, vs ~2MB original. The multiplier is **1.24x** (original + variants), not 4x. Total should be ~90PB, not 290PB.

**Fix:** Recalculate with actual variant sizes.

---

### 2.10 Instagram: S3 Cost "$55K/mo After One Year" Is Off by 30x

**File:** `content/problems/instagram/index.html` line 28

At ~90PB stored after year 1 at $0.023/GB, month-12 storage cost would be ~$2.07M/month. $55K/month implies only ~2.4PB stored (12 days of uploads).

**Fix:** Recalculate with realistic year-end storage volumes.

---

### 2.11 Instagram: CDN Transfer Estimate "~4PB/month" Is Too Low by 10x+

**File:** `content/problems/instagram/index.html` line 27

With 500M DAU × 10 feed loads × 20 images at smallest thumbnail (15KB), CDN edge transfer = ~1.5PB/day = ~45PB/month. The 4PB figure may confuse origin transfer with edge transfer.

---

### 2.12 Instagram: Feed Query Uses OFFSET Despite Arguing Against It

**File:** `content/problems/instagram/index.html` line 190 vs lines 53, 57-58

The API design explicitly argues for cursor-based pagination. But the SQL query uses `LIMIT 20 OFFSET :offset`. Code contradicts the stated best practice.

**Fix:** Change the query to use cursor-based pagination (`WHERE created_at < :cursor`).

---

### 2.13 Instagram: `attribute_not_exists(user_id)` Is Wrong DynamoDB Usage

**File:** `content/problems/instagram/index.html` deepQA line 841

For a table with composite key (post_id, user_id), `attribute_not_exists(user_id)` doesn't work as intended because user_id always exists as a key attribute. The correct pattern is `attribute_not_exists(PK)` where PK is the partition key attribute name.

**Fix:** Change to `attribute_not_exists(post_id)` or explain the correct pattern.

---

### 2.14 Instagram: Fan-Out Math Inconsistency

**File:** `content/problems/instagram/index.html` line 234

States "50K posts/day = 10M Redis writes/day" but the NFR diagram says 100M photos/day at full scale. At 10M DAU with 5% posting, that's 500K posts/day (not 50K). The fan-out WPS is understated by ~10x.

---

### 2.15 Redis: "Single-Threaded" Characterization Is Outdated

**File:** `content/fundamentals/storage/redis/index.html` line 18

States "Memcached is multi-threaded; Redis is single-threaded per shard." Since **Redis 6.0** (2020), Redis supports multi-threaded I/O for network processing via `io-threads` config (command execution remains single-threaded).

**Fix:** Clarify that Redis command execution is single-threaded but network I/O can be multi-threaded since 6.0.

---

### 2.16 CloudFront Functions: ES 5.1 Is Outdated

**File:** `content/fundamentals/networking/cdn/cloudfront/index.html` line 134

States CloudFront Functions runtime is "JavaScript (ES 5.1)." Since 2024, the `cloudfront-js-2.0` runtime supports modern JavaScript (ES 2024+).

**Fix:** Mention both runtimes.

---

### 2.17 ALB/NLB Base Cost: $16.20/mo Should Be $16.43/mo

**Files:** `load-balancing/index.html` lines 248, 258; `alb/index.html` lines 54, 187, 198, 202; `nlb/index.html` line 163

$0.0225/hr × 730 hours = $16.43/mo, not $16.20/mo. Minor but repeated across 7 locations.

---

### 2.18 Base62: Collision Probability Understated by 2x

**File:** `content/algorithms/base62-encoding/index.html` lines 168, 204

States 1B URLs in 3.52T keyspace yields "~0.014% per insert" with "14,000 collisions per 100M inserts." Actual: 1B/3.52T = **0.0284%**, yielding ~28,400 collisions per 100M. Understated by 2x.

---

### 2.19 Consistent Hashing: Standard Deviation Inconsistency

**File:** `content/algorithms/consistent-hashing/index.html` lines 50-57 vs line 175

Table says 150 vnodes → "~5%" std dev. Avoid block says "about 8%." The stated formula `1/sqrt(vnodes)` gives 1/√150 = 8.2%. The table and the formula disagree.

**Fix:** Make the table consistent with the formula, or note the table uses empirical measurements.

---

### 2.20 RAG: BM25 Described as "TF-IDF"

**File:** `content/patterns/rag/index.html` line 312

States BM25 is "TF-IDF term frequency matching." BM25 is a **probabilistic relevance model** that extends TF-IDF with document length normalization and term frequency saturation. Not the same thing.

**Fix:** "BM25 probabilistic keyword matching" or "BM25 (extends TF-IDF with saturation and length normalization)."

---

### 2.21 Embeddings: Euclidean Distance Warning Contradicts Normalized Embeddings

**File:** `content/patterns/embeddings-vector-search/index.html` line 120 vs line 124

Line 120 warns Euclidean distance penalizes longer documents. Line 124 correctly notes OpenAI embeddings are normalized. For normalized vectors, Euclidean and cosine give equivalent rankings — the warning is invalid for the very embeddings being discussed.

---

### 2.22 DynamoDB: Transactions Limit Is "100 Actions" Not "100 Items"

**File:** `content/fundamentals/storage/dynamodb/index.html` line 12

States "ACID transactions across up to 100 items." The limit is **100 actions** per transaction. Multiple actions can target the same item (ConditionCheck + Update = 2 actions, 1 item).

---

### 2.23 Rate Limiter: Sliding Window Lua Script Time Unit Mismatch

**File:** `content/problems/rate-limiter/index.html` lines 278-303

The Lua script uses `now - window * 1000` (converting to ms) but `now` comes from Python's `time.time()` which returns seconds. Mixing seconds and milliseconds.

**Fix:** Use `time.time() * 1000` for millisecond `now`, or keep window in seconds.

---

### 2.24 Rate Limiter: FastAPI Middleware Signature Is Wrong

**File:** `content/problems/rate-limiter/index.html` line 221

`async def rate_limit_middleware(request: Request, call_next)` is not valid Starlette middleware. The correct decorator is `@app.middleware("http")` with signature `(request: Request, call_next: Callable)`.

---

### 2.25 Redis: "Never Run Without Persistence" Is Overly Absolute

**File:** `content/fundamentals/storage/redis/index.html` line 396

The `{{avoid}}` states "Never run Redis without persistence in production, even as a cache." Many production systems intentionally run Redis as a pure cache without persistence — data loss on restart is acceptable, and persistence adds latency.

**Fix:** Qualify: "Consider persistence for any data you can't afford to lose on restart."

---

## 3. Structural Issues

### 3.1 Missing `{{say}}` Openers (System-Wide)

The gold standard requires every phase to open with `{{say}}`. The following phases are missing them:

**Problems:**
| File | Missing `{{say}}` on Phases |
|------|-----------------------------|
| URL Shortener | 12 |
| Rate Limiter | 2, 3, 7, 9, 10, 11, 12 (7 of 12 phases) |
| Instagram | 12 |

**Fundamentals:**
| File | Missing `{{say}}` on Phases |
|------|-----------------------------|
| Load Balancing | 8 |
| ALB | 3, 5 |
| NLB | 2, 5, 6 |
| CDN | 5, 7, 8 |
| CloudFront | 6, 7, 8 |
| Redis | 8 |
| DynamoDB | 8, 9 |

**Algorithms:**
| File | Missing `{{say}}` on Phases |
|------|-----------------------------|
| Base62 Encoding | 4, 6 |
| Bloom Filter | 4, 7 |
| Consistent Hashing | 6, 8 |
| Snowflake ID | 5, 8 |
| Token Bucket | 7 |

**Patterns:**
| File | Missing `{{say}}` on Phases |
|------|-----------------------------|
| RAG | 10 |
| Agent Tools | 9, 10 |
| Prompt Chaining | 7, 9, 10 |
| Guardrails | 6, 7, 8, 9 |
| Embeddings | 9, 10 |

**Total: 48 phases missing `{{say}}` across all files.**

---

### 3.2 Phase Count Violations

**Fundamentals (should be 8 phases):**
| File | Actual | Issue |
|------|--------|-------|
| ALB | 6 | Missing 2 phases |
| NLB | 6 | Missing 2 phases |
| DynamoDB | 9 | 1 extra phase |

**Algorithms (should be 8 phases):**
| File | Actual | Issue |
|------|--------|-------|
| Token Bucket | 7 | Missing 1 phase |

**Patterns (should be 8 phases):**
| File | Actual | Issue |
|------|--------|-------|
| RAG | 10 | 2 extra phases |
| Agent Tools | 10 | 2 extra phases |
| Prompt Chaining | 10 | 2 extra phases |
| Guardrails | 9 | 1 extra phase |
| Embeddings | 10 | 2 extra phases |

---

### 3.3 Missing Closing `{{key}}`

6 of 7 fundamental files, all 5 algorithm files, and all 5 pattern files have `{{key}}` macros mid-file but lack the required single closing `{{key}}` after the final `{{checklist}}`. Only `redis/index.html` correctly ends with `{{checklist}}` followed by `{{key}}`.

**All files with multiple mid-file `{{key}}` blocks (should have exactly 1 at end):**
- Load Balancing (line 101), ALB (line 181), NLB (line 107), CDN (line 99), CloudFront (lines 120, 322), DynamoDB (line 291)
- Base62 (line 28), Bloom Filter (lines 20, 77, 235), Consistent Hashing (line 20), Snowflake (lines 23, 165, 252), Token Bucket (line 22)
- All 5 pattern files

---

### 3.4 Missing `{{think}}` Blocks

The gold standard requires at least 1 `{{think}}` per phase (replacing legacy `{{thought}}`).

**URL Shortener:** Missing in phases 5, 8
**Rate Limiter:** Missing in phases 5, 7, 8
**Instagram:** Missing in phases 4, 6, 8, 11
**Fundamentals:** Multiple phases in each file lack `{{think}}`
**Algorithms:** Generally adequate
**Patterns:** Most phases across all 5 files lack `{{think}}`

---

### 3.5 Missing `{{triggerQs}}`

**Rate Limiter Phase 3:** Zero content — just two diagram slugs with no text, say, hint, think, code, or triggerQs. This phase is essentially empty.

**Instagram:** Missing on phases 8 (Global Scale) and 10 (Analytics)

---

### 3.6 URL Shortener: Two Diagrams Referenced but Not Registered

- `url-db-access-patterns` (content line 132) — not in `url_shortener.go`
- `url-sub-problems` (content line 703) — not in `url_shortener.go`

These will cause runtime rendering errors.

---

### 3.7 Rate Limiter deepQA: Item 1 Missing 3rd Nesting Level

The deepQA's first item has a sub-question without a `dqa-deep` nested inside it. All 5 items should have 3-level nesting per the gold standard.

---

## 4. Missing Content

### 4.1 Instagram: No Video Processing Discussion

**File:** `content/problems/instagram/index.html`

The opening scopes the system to "photo/video sharing" and the API includes Reels. But there is zero discussion of: video transcoding (H.264/H.265/AV1), adaptive bitrate streaming (HLS/DASH), video thumbnail extraction, video CDN delivery, or upload resumability. This is a significant gap for an Instagram design.

---

### 4.2 Instagram: No Reels/Short-Video Architecture

Reels is Instagram's primary growth driver. Even at P2, there should be at least a high-level architectural extension discussion.

---

### 4.3 Instagram: No Feed Ranking ML Details

Line 332 briefly mentions a ranking score but provides no detail on: feature engineering, model architecture, online vs offline scoring, cold-start problem, or A/B testing.

---

### 4.4 Instagram: Stories Buried in Security Phase

Stories are a P1 feature but are discussed inside Phase 9 (Security & Content Moderation) at lines 573-604. They deserve their own phase or should be in Phase 5 (Growing) where architectural additions are discussed.

---

### 4.5 Rate Limiter: Phase 3 Is Completely Empty

Phase 3 "API Design with Rate Headers" consists of only two `{{diagram}}` references. No text, no `{{say}}`, no hints, no code, no explanation of RFC 6585, no mention of the IETF RateLimit header draft standard.

---

### 4.6 Rate Limiter: No `{{qa}}` Macros Anywhere

The gold standard requires `{{qa}}` used "inline throughout." The file has zero `{{qa}}` macros.

---

### 4.7 Rate Limiter: Leaky Bucket Algorithm Not Covered

Content mentions "four main algorithms" but only covers Token Bucket, Sliding Window Counter, Fixed Window, and Sliding Window Log. Leaky Bucket (distinct from Token Bucket — processes at fixed rate, queues excess) is missing.

---

### 4.8 Rate Limiter: No Monitoring/Observability Phase

Rate limiting monitoring (Prometheus metrics, Grafana dashboards, alerting on false positives) is only briefly mentioned in deepQA but never gets dedicated treatment.

---

### 4.9 URL Shortener: No Discussion of 307/308 Redirects

Only compares 301 vs 302. HTTP 307 (preserves method) is used by bit.ly (mentioned in deepQA) but not in main content.

---

### 4.10 DynamoDB: Global Tables Not Discussed in Depth

Multi-region active-active is listed in features table but never explored. No diagram, no conflict resolution discussion. Common interview topic.

---

### 4.11 DynamoDB: DAX Integration Barely Mentioned

Only appears in a hint at line 200. No diagram, no code showing DAX integration pattern.

---

### 4.12 DynamoDB: Streams Barely Mentioned

Only a brief code example at line 214. No discussion of stream processing patterns or consumer group management.

---

### 4.13 Patterns: Severely Under-Diagrammed

All 5 pattern files average only 2 diagrams each (total 10), compared to 19-28 for problem files. Many phases have no visual aids at all.

| Pattern File | Diagrams | Phases Without Diagrams |
|-------------|----------|------------------------|
| RAG | 2 | 1, 3, 4, 5, 7, 8, 9, 10 |
| Agent Tools | 2 | 1, 3, 4, 5, 6, 7, 8, 10 |
| Prompt Chaining | 2 | 2, 3, 4, 6, 7, 8, 9, 10 |
| Guardrails | 1 | 2, 3, 4, 5, 6, 7, 8, 9 |
| Embeddings | 3 | 2, 3, 5, 6, 8, 9, 10 |

---

### 4.14 Patterns: Insufficient `{{hint}}` Density

Most pattern phases lack `{{hint}}` entirely. The gold standard requires 1-2 per phase.

- Agent Tools: Only 1 hint in the entire file (phase 2)
- Guardrails: Only 1 hint in the entire file (phase 2)
- Prompt Chaining: Only 2 hints in the entire file

---

### 4.15 Fundamentals: Missing Diagrams

| File | Total Diagrams | Phases Without Diagrams |
|------|---------------|------------------------|
| CDN | 1 | 1, 3, 4, 5, 6, 7, 8 |
| Load Balancing | 1 | 1, 2, 3, 4, 5, 6, 7 |
| DynamoDB | 3 | 1, 3, 4, 5, 7, 8, 9 |
| Redis | 3 | 1, 7 |

---

### 4.16 Instagram: No Idempotent Fan-Out or Write-Ahead Log

The fan-out code iterates followers in a single pipeline. No discussion of checkpointing, idempotent retry, or durable queue for partial failure recovery.

---

## 5. Quality Issues

### 5.1 Instagram: Phase Timing Totals 77 Minutes

5+5+5+5+8+8+5+5+8+5+8+10 = 77 minutes. A typical system design interview is 45 minutes. No guidance on which phases to prioritize or skip.

---

### 5.2 Instagram: Only 2 `{{avoid}}` and 2 `{{compare}}` Blocks

For a 12-phase problem with many pitfalls (storing images in DB, single auto-increment ID across shards, polling for notifications, fan-out-on-write for celebrities), 2 `{{avoid}}` blocks is sparse.

---

### 5.3 URL Shortener: Cache Hit Rate Math Ambiguous

Line 33 shows cascading hit rates (Browser 30% → CloudFront 60% of remainder → Redis 95% of remainder). The checklist at line 709 presents the same numbers non-cascading, which is ambiguous.

---

### 5.4 URL Shortener: `click_count` Dual Storage Unexplained

`click_count` exists in both the `urls` table and `click_analytics` table. The consistency mechanism between them is never explained.

---

### 5.5 Rate Limiter: Sliding Window Log Lua Uses Non-Unique Members

`ZADD key now member` — if two requests happen at the same millisecond with the same member value, the second overwrites the first (sorted sets have unique members). Need a UUID or counter per request.

---

### 5.6 Rate Limiter: Phase 2 "Capacity Estimation" Is Very Thin

Only a diagram, 2 hints, and a say. No walkthrough of back-of-envelope math with numbers. Missing: DAU assumptions, read/write ratio, storage growth, bandwidth estimation.

---

### 5.7 Instagram: Redis Bitmap Comparison Is Self-Defeating

Line 603 shows the bitmap (62.5MB) is **7.8x larger** than the set approach (8MB) for 1M viewers, undermining the argument for bitmaps. The crossover point (~7.8M viewers) isn't discussed.

---

### 5.8 Fundamentals: Redundant DynamoDB Q&A

`dynamodb/index.html` lines 18 and 27: `{{triggerQs}}` and `{{qa}}` ask nearly identical questions ("When to pick DynamoDB over Postgres?") with nearly identical answers.

---

### 5.9 RAG: GPT-4o Pricing May Be Outdated

Phase 9 cost table uses $5/1M tokens for GPT-4o input. Current pricing is $2.50/1M input, $10/1M output. The calculation only accounts for input tokens.

---

## 6. Infrastructure & Code Bugs

### 6.1 Silent Failure on Unresolved Registry References

**File:** `internal/registry/registry.go` lines 156-174, 193-214

If a problem's `uses` references a non-existent fundamental slug, or an algorithm's `used_in` references a non-existent problem slug, the reference is silently ignored. No warning logged. A YAML typo would silently break cross-linking with no diagnostic output.

**Fix:** Log warnings for unresolved references.

---

### 6.2 Sidebar Section Type Mapping Is Fragile

**File:** `web/templates/base.html` lines 148-158

Section type is determined by DOM index (`sectionTypes[i]`). But Algorithms and Patterns sections are conditionally rendered. If either is empty and omitted, the index mapping shifts and wrong types are assigned.

**Fix:** Use data attributes instead of index-based mapping.

---

### 6.3 `extractKeywords` Splits on Hyphens, Destroying Compound Terms

**File:** `internal/handlers/handlers.go` lines 173-183

Config value "round-robin sticky sessions" yields keywords ["round", "robin", "sticky", "sessions"]. The highlight JS then matches any element containing "round" or "robin" — false-positive highlights.

**Fix:** Keep hyphenated terms intact.

---

### 6.4 Welcome Page Fundamentals Stat Undercounts

**File:** `web/templates/welcome.html` line 26

`{{len .Fundamentals}}` counts only top-level fundamentals (4), not including children (ALB, NLB, CloudFront). Actual count is 7.

**Fix:** Use a recursive count or a pre-computed total.

---

### 6.5 Sidebar "Active" State Not Set on Full Page Load

**File:** `web/templates/sidebar.html`

The `active` class is only added via JavaScript `onclick`. On direct URL access or browser refresh, no sidebar item is highlighted even though `ActiveSlug` is available in template data.

**Fix:** Add template conditional: `{{if eq .ActiveSlug .Slug}}active{{end}}`.

---

### 6.6 Book Mode Wheel Handler Prevents Scrolling in Child Elements

**File:** `web/templates/base.html` lines 559-568

`e.preventDefault()` on all wheel events prevents scrolling inside code blocks, overflow diagrams, and tables.

**Fix:** Check if the event target is within an overflow container before preventing default.

---

### 6.7 No `Content-Type` Header for JSON Responses

**File:** `internal/handlers/handlers.go` lines 316, 321

JSON responses are sent without `Content-Type: application/json`, so they may be served as `text/plain`.

**Fix:** Add `w.Header().Set("Content-Type", "application/json")`.

---

### 6.8 Dead Code

- `doc_card.html` — parsed but never rendered
- `TopLevelCategories()` method — exported but never called
- `FundamentalsByReferenceCount()` method — exported but never called
- `_ = slug` in `AllFundamentals` — should be `for _, f := range`

---

## 7. Security Issues

### 7.1 XSS via Unescaped Error Messages

**File:** `internal/handlers/handlers.go` lines 199, 205

`err.Error()` is concatenated into `template.HTML` without escaping. Low risk (content is embedded) but bad pattern.

---

### 7.2 XSS in `info` Macro via Unescaped Title Attribute

**File:** `internal/macros/macros.go` line 255

The `definition` string is placed into the `title` attribute via `%s` without HTML attribute escaping. A definition containing `"` breaks out of the attribute.

**Fix:** Use `template.HTMLEscapeString(definition)`.

---

### 7.3 XSS via localStorage in "Continue Reading"

**File:** `web/templates/base.html` lines 342-343

`progress.lastPage` from localStorage is injected into `innerHTML` without sanitization.

**Fix:** Use `createElement`/`setAttribute` instead of `innerHTML`.

---

### 7.4 All Macro Functions Embed Content Without Escaping

**File:** `internal/macros/macros.go` lines 156-554

Macros like `say`, `qa`, `checklist` take plain text but embed via `fmt.Sprintf` into `template.HTML`. Currently safe only because content is from trusted embedded files. If content authoring ever opens up, this becomes an XSS vector.

---

### 7.5 Agent Tools: SQL Injection Prevention via String Matching

**File:** `content/patterns/agent-tools/index.html` lines 377-386

The `query_database` function uses `startswith("SELECT")` and keyword scanning. Should note that production code needs parameterized queries and a read-only DB user.

---

### 7.6 Embeddings: SQL Injection via `.format()` for Table Names

**File:** `content/patterns/embeddings-vector-search/index.html` lines 238-257

Uses `.format()` for table names in SQL. Should note this is safe only for controlled, hardcoded table names.

---

## 8. CSS & UI Issues

### 8.1 No CSS Rules for `.doc-card` Classes

`doc_card.html` uses `.doc-card`, `.doc-card-type`, `.doc-card-status`, `.status-missing` — none have CSS rules in `style.css`.

---

### 8.2 `.hint-popup` Clipped by Ancestor `overflow: hidden`

Hint popups are absolutely positioned. If inside a diagram container (`overflow-x: auto`) or book mode (`overflow: hidden`), they get clipped.

**Fix:** Use portal/body-level positioning.

---

### 8.3 Book Mode Forces `overflow: hidden` on All Children

**File:** `web/static/css/style.css` lines 3454-3460

This clips wide tables, long code lines, and diagram containers with no scroll affordance.

---

### 8.4 Pattern Diagrams Use Undocumented CSS Classes

`pat-embedding-space` uses `.d-coord-*` classes and `pat-hnsw-index-structure` uses `.d-graph-*` classes. Neither set is listed in CLAUDE.md or verified to exist in `style.css`.

---

### 8.5 `-webkit-overflow-scrolling: touch` Is Deprecated

**File:** `web/static/css/style.css` line 3212

No effect in modern iOS Safari, generates console warnings.

---

## 9. Documentation vs Implementation Mismatches

### 9.1 CLAUDE.md Says 102 Diagrams; Actual Count Is 197

The server log shows `Loaded 197 diagrams`. The per-file counts in CLAUDE.md are also stale:

| File | CLAUDE.md Claims | Actual |
|------|-----------------|--------|
| `rate_limiter.go` | 19 | 19 |
| `instagram.go` | 24 | 28 |
| `url_shortener.go` | 21 | 26 |
| `algorithms.go` | 13 | 13 |
| `fundamentals.go` | 15 | 15 |
| `patterns.go` | 10 | 10 |
| **Total** | **102** | **111+** |

---

### 9.2 CLAUDE.md Describes Two-Panel Layout; CSS Implements Overlay Sidebar

CLAUDE.md states "Two-panel layout: sidebar + detail area." The actual implementation uses a full-width detail area with the sidebar as a fixed-position overlay toggled by hamburger. There is no desktop breakpoint showing sidebar and detail side-by-side.

---

### 9.3 External CDN Dependency for highlight.js

HTMX and Alpine.js are vendored locally, but highlight.js is loaded from `cdnjs.cloudflare.com`. Inconsistent with the vendoring strategy.

---

## 10. Improvement Opportunities

### 10.1 Add Video Processing Pipeline to Instagram

Even as a section within an existing phase, cover: video transcoding codecs (H.264/H.265/AV1), adaptive bitrate streaming (HLS/DASH), thumbnail extraction, video-specific CDN delivery.

---

### 10.2 Add Leaky Bucket to Rate Limiter

Leaky Bucket is distinct from Token Bucket (processes at fixed rate, queues excess). Include a diagram and comparison.

---

### 10.3 Add `{{compare}}` for Redis vs In-Memory vs DynamoDB in Rate Limiter

Phase 8 discusses Redis config but never formally compares storage options.

---

### 10.4 Add RFC References to Rate Limiter Phase 3

RFC 6585 (429 status code) and IETF draft `draft-ietf-httpapi-ratelimit-headers` would signal domain expertise.

---

### 10.5 Add Capacity Estimation Table to Rate Limiter Phase 2

Structured `{{table}}` walking through back-of-envelope math step by step.

---

### 10.6 Add Feed Ranking ML Details to Instagram

Feature engineering, model architecture (two-tower, wide-and-deep), cold-start problem, A/B testing.

---

### 10.7 Add Media Deduplication to Instagram

Perceptual hashing can deduplicate 20-30% of social media uploads. Good optimization topic.

---

### 10.8 Add Cache Warming Strategy Discussion

Relevant after deployments that change cache key formats, TTL policies, or data structures.

---

### 10.9 Vendor highlight.js Locally

Consistent with HTMX and Alpine.js vendoring strategy. Eliminates external CDN dependency.

---

### 10.10 Add Input Validation to PDF Handlers

`GeneratePDF` and `GenerateStatus` don't validate slug/taskID existence before responding. Add validation for when these are implemented.

---

### 10.11 Add Warning Logs for Unresolved Registry References

Log warnings when `uses`, `used_in`, or `appears_in` reference non-existent slugs. Silent failures make YAML typos invisible.

---

### 10.12 Add Data Attributes to Sidebar Sections

Replace index-based section type mapping with `data-section-type="problem"` attributes. Prevents bugs when sections are conditionally rendered.

---

## Summary by Priority

| Priority | Category | Count |
|----------|----------|-------|
| **P0** | Critical Factual Errors | 10 |
| **P1** | Factual Errors | 15 |
| **P1** | Missing `{{say}}` (48 phases) | 1 (systemic) |
| **P1** | Phase Count Violations | 8 files |
| **P1** | Missing Closing `{{key}}` | 18 files |
| **P2** | Missing Content | 16 |
| **P2** | Quality Issues | 9 |
| **P2** | Infrastructure Bugs | 8 |
| **P2** | Security Issues | 6 |
| **P3** | CSS/UI Issues | 5 |
| **P3** | Documentation Mismatches | 3 |
| **P3** | Improvement Opportunities | 12 |

**Top 5 items to fix first:**
1. Base62 diagram math — every step is wrong (credibility killer)
2. Read:write ratio contradictions in URL Shortener + Instagram — affects all capacity planning
3. DynamoDB "3M reads/sec" — 1000x error
4. CloudFront invalidation limit — wrong in 2 of 3 places, correct in 1
5. 48 phases missing `{{say}}` — systemic gold standard violation

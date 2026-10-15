# 08. Running agents in production · 生产环境中的 Agent

What changes when an agent leaves the demo: validating structured output, controlling cost and latency, and staying reliable under rate limits, overload and restarts. Interviewers ask because most production agent failures are ordinary distributed-systems problems with a model in the middle.

> Agent 离开演示环境之后会遇到什么：校验结构化输出、控制成本和延迟，以及在限流、过载和重启之下保持可靠。面试官会问这些，是因为生产环境中大多数 Agent 故障，其实都是中间夹着一个模型的普通分布式系统问题。

3 questions · generated 2026-10-15T18:25:03+08:00 · [← 目录 / Contents](../README.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Production | Medium | [Supplier data to catalog fields: malformed, truncated, or confidently wrong](#q1) | 供应商数据转商品目录字段：格式错误、输出截断，或者错得理直气壮 |
| 2 | Production | Hard | [A support agent that costs too much and takes 90 seconds to answer](#q2) | 一个太贵、还要 90 秒才能回复的客服 Agent |
| 3 | Production | Hard | [Peak traffic: 429s, overloaded errors, double refunds and restarted tasks](#q3) | 流量高峰：429、过载错误、重复退款和从头再来的任务 |

---

<a id="q1"></a>

## 1. Supplier data to catalog fields: malformed, truncated, or confidently wrong

**供应商数据转商品目录字段：格式错误、输出截断，或者错得理直气壮** · `Production` · `Medium`

### Question · 题目

A pipeline uses an LLM to turn messy supplier product data (titles in mixed languages, free-text specs, HTML descriptions) into ad-platform catalog fields such as `title`, `google_product_category`, `price`, `currency`, `availability` and `gtin`. The prompt just ends with "Return JSON".

> 一条数据管道用 LLM 把杂乱的供应商商品数据（多种语言混杂的标题、自由文本的规格、HTML 描述）转换成广告平台的商品目录字段，比如 `title`、`google_product_category`、`price`、`currency`、`availability` 和 `gtin`。prompt 的结尾只有一句“Return JSON”。

In production, some responses don't parse, some stop in the middle of an object, and some parse fine but are wrong: a price of `1999` instead of `19.99`, `availability: "yes"`, an 11-digit GTIN, a Chinese title in a US feed. The ad platform rejects some of these products and, worse, accepts others. How do you make this pipeline reliable?

> 在生产环境中，有些响应无法解析，有些在对象中间就断掉了，还有一些能正常解析却是错的：价格是 `1999` 而不是 `19.99`，出现了 `availability: "yes"`，GTIN 只有 11 位，美国 feed 里出现了中文标题。广告平台拒收了其中一些商品，更糟的是，另一些被接收了。你如何让这条管道变得可靠？

<details>
<summary><b>Answer · 参考答案</b></summary>

**Three failures, three layers.** These are different problems. Syntax: the output is not valid JSON. Completeness: the output stopped early. Semantics: valid JSON with wrong values. Structured outputs fix the first, checking the stop reason handles the second, and only your own validation catches the third.

> **三种故障，三层防线。** 这是三个不同的问题。语法：输出不是合法的 JSON。完整性：输出提前停止了。语义：JSON 合法，但值是错的。结构化输出解决第一个，检查 stop reason 处理第二个，而第三个只能靠你自己的校验来发现。

**Constrain the format instead of asking for it.** "Return JSON" is a request the model usually follows. Structured outputs make it a constraint: you pass a JSON Schema, and generation is restricted to output that conforms. In the Claude API that is `output_config.format` for the final answer, or `strict: true` on a tool so its input matches its `input_schema`. Use `enum` for fixed sets, `additionalProperties: false`, and list every field in `required`:

> **用约束代替请求。** “Return JSON”只是一个模型通常会遵守的请求。结构化输出（structured outputs）把它变成约束：你传入一个 JSON Schema，生成过程只能产出符合它的内容。在 Claude API 中，最终回答用 `output_config.format`，工具则设置 `strict: true`，让工具输入符合它的 `input_schema`。固定取值的集合用 `enum`，设置 `additionalProperties: false`，并把每个字段都列进 `required`：

```json
"output_config": {
  "format": {
    "type": "json_schema",
    "schema": {
      "type": "object",
      "properties": {
        "title": {"type": "string"},
        "google_product_category": {"type": "string"},
        "price": {"type": "number"},
        "currency": {"type": "string", "enum": ["USD", "EUR", "GBP"]},
        "availability": {"type": "string",
          "enum": ["in_stock", "out_of_stock", "preorder", "backorder"]},
        "gtin": {"anyOf": [{"type": "string"}, {"type": "null"}]}
      },
      "required": ["title", "google_product_category", "price",
                   "currency", "availability", "gtin"],
      "additionalProperties": false
    }
  }
}
```

**What the schema guarantees, and what it doesn't.** You get parseable JSON, every required key, the right types, and enum values from your list. You don't get correct values: `1999` is a perfectly valid number. Some constraints are not enforced during generation at all: in the Claude API, structured outputs don't support numeric ranges (`minimum`, `maximum`) or string lengths (`minLength`, `maxLength`), and the Python and TypeScript SDKs strip them from the schema and check them on the client instead. And the output can still be incomplete.

> **schema 能保证什么，不能保证什么。** 你能得到可解析的 JSON、所有必填字段、正确的类型，以及来自你列表中的 enum 值。你得不到正确的值：`1999` 是一个完全合法的数字。有些约束在生成时根本不会被强制执行：在 Claude API 中，结构化输出不支持数值范围（`minimum`、`maximum`）和字符串长度（`minLength`、`maxLength`），Python 和 TypeScript SDK 会把它们从 schema 中去掉，改为在客户端检查。另外，输出仍然可能是不完整的。

**Check the stop reason before you parse.** If `stop_reason` is `max_tokens`, the output was cut off: never use the partial object, and never run a `tool_use` block from that turn, because its input may be truncated. Raise `max_tokens` (a ceiling, not a target) and send fewer products per request; a small group is far easier to retry than fifty. If it is `refusal`, the output may not match the schema, so route the item to review instead of retrying in a loop.

> **解析之前先检查 stop reason。** 如果 `stop_reason` 是 `max_tokens`，说明输出被截断了：永远不要使用残缺的对象，也永远不要执行这一轮中的 `tool_use` 块，因为它的输入可能是不完整的。调高 `max_tokens`（它是上限，不是目标），并减少每次请求中的商品数量；一小组商品远比五十个容易重试。如果是 `refusal`，输出可能不符合 schema，所以要把这个商品送去人工审核，而不是循环重试。

**Validate semantics after the schema.** This is ordinary backend validation, and it is where the expensive mistakes are caught. Write each error so it can be sent back to the model for a repair:

> **在 schema 之后做语义校验。** 这就是普通的后端校验，代价高昂的错误正是在这里被拦下的。每条错误信息都要写得可以直接发回给模型，用于修复：

```go
func semanticErrors(it CatalogItem, feed Feed) []string {
	var errs []string
	if it.GTIN != nil && !gtinValid(*it.GTIN) { // 8, 12, 13 or 14 digits + check digit
		errs = append(errs, "gtin: invalid length or check digit")
	}
	if it.Currency != feed.Currency {
		errs = append(errs, "currency: want "+feed.Currency)
	}
	if !feed.Taxonomy.Has(it.GoogleProductCategory) {
		errs = append(errs, "google_product_category: not in the taxonomy")
	}
	if lang := detectLanguage(it.Title); lang != feed.Language {
		errs = append(errs, "title: language is "+lang+", want "+feed.Language)
	}
	return errs
}
```

**Don't ask the model for what code can compute.** Parse prices, currencies and GTINs deterministically wherever the supplier data has them, and use the model for the fuzzy parts: titles, descriptions, category mapping, color names. When the model must extract a number from free text, cross-check it against any structured field and a plausible price range. For a large taxonomy, don't put thousands of IDs in an enum; let the model describe the category and map it to an ID in code, or offer a short candidate list.

> **代码能算的，就不要让模型来做。** 只要供应商数据里有价格、币种和 GTIN，就用确定性的代码去解析，把模型用在模糊的部分：标题、描述、类目映射、颜色名称。当模型必须从自由文本中提取数字时，要拿它和任何结构化字段交叉核对，并检查是否落在合理的价格区间内。对于很大的类目体系，不要把几千个 ID 塞进一个 enum；让模型描述类目，再由代码映射成 ID，或者只给它一份简短的候选列表。

**Retry and repair, cheapest first.** When validation fails, escalate step by step:

> **重试与修复，先用最便宜的办法。** 校验失败时，逐级升级处理：

- Deterministic repair: map `"In Stock"` to `in_stock`, trim whitespace, normalize units. No model call needed.
  > 确定性修复：把 `"In Stock"` 映射为 `in_stock`、去掉多余空白、统一单位。不需要调用模型。
- Targeted repair: send the source data, the previous output and the exact validation errors, and ask the model to fix only those fields. Allow one or two attempts.
  > 定向修复：把源数据、上一次的输出和具体的校验错误一起发给模型，让它只修正这些字段。最多尝试一到两次。
- Truncation: retry with a higher `max_tokens` or a smaller batch, not the same request again.
  > 截断：用更高的 `max_tokens` 或更小的批次重试，而不是原样再发一次。
- Give up safely: after the retry budget, send the item to a review queue and keep the last good version live. Publishing nothing beats publishing a wrong price.
  > 安全地放弃：重试预算用完后，把商品放进人工审核队列，线上继续保留上一个正确的版本。什么都不发布，也好过发布一个错误的价格。

Key each job by product ID and a hash of the source data, so reprocessing is idempotent and an unchanged product is never sent to the model again.

> 每个任务以商品 ID 加源数据哈希作为键，这样重新处理是幂等的，没有变化的商品也不会再次发给模型。

**Measure what reaches the platform.** Track failure rates per field and per supplier, repair success, and the ad platform's rejection reasons, which are the ground truth. Sample accepted items for human review, because the dangerous errors are the ones that pass every check, and turn each caught mistake into a test case.

> **衡量真正到达平台的结果。** 按字段和按供应商跟踪失败率、修复成功率，以及广告平台给出的拒收原因，后者才是最终的事实依据。对已被接收的商品抽样做人工审核，因为最危险的错误恰恰是那些通过了所有检查的错误；并把每一个被抓到的错误都变成测试用例。

**Trade-offs.** A new schema has a one-time compilation cost on its first request. Over-strict validators cause more retries and cost, so tune them against real rejections. Small batches resend the shared instructions more often; prompt caching on that prefix keeps this cheap.

> **权衡。** 新的 schema 在第一次请求时会有一次性的编译开销。校验器过于严格会带来更多重试和更高成本，所以要根据真实的拒收数据来调整。小批次会更频繁地重复发送共享的指令；对这段前缀使用 prompt caching，可以让它保持便宜。

</details>

### Follow-ups · 追问

- The taxonomy has thousands of categories. How would you get reliable category IDs without putting them all in the schema?
  > 类目体系有几千个类目。在不把它们全部放进 schema 的情况下，你如何得到可靠的类目 ID？
- Every validator passes, but the ad platform still disapproves 3% of items. How do you find and close that gap?
  > 所有校验都通过了，但广告平台仍然拒收了 3% 的商品。你如何找到并弥补这个差距？
- When would you use strict tool use instead of a structured final answer in this pipeline?
  > 在这条管道里，什么时候你会使用 strict tool use，而不是结构化的最终回答？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| structured outputs | 结构化输出 |
| schema validation | schema 校验 |
| semantic validation | 语义校验 |
| truncated output | 被截断的输出 |
| ground truth | 事实依据（真实标准） |

---

<a id="q2"></a>

## 2. A support agent that costs too much and takes 90 seconds to answer

**一个太贵、还要 90 秒才能回复的客服 Agent** · `Production` · `Hard`

### Question · 题目

Your merchant-support agent answers tickets such as "Why were half of my products disapproved in Google Shopping?" It has a 12,000-token system prompt full of policies, 40 tool definitions, and a loop that averages 14 model calls per ticket. Tools return raw JSON from internal APIs, sometimes 30 KB per call. Every call uses the largest model at a high effort setting.

> 你们的商家客服 Agent 负责回答类似“为什么我一半的商品在 Google Shopping 被拒登了？”这样的工单。它有一个写满各种政策、长达 12,000 token 的 system prompt，40 个工具定义，平均每张工单要调用模型 14 次。工具直接返回内部 API 的原始 JSON，有时一次就有 30 KB。每次调用都使用最大的模型，并且是高 effort 设置。

Finance says the cost per resolved ticket is too high, and the p95 time to answer is 90 seconds. How would you measure the problem, and what would you change, in what order?

> 财务说每张已解决工单的成本太高，而回复时间的 p95 是 90 秒。你会如何衡量这个问题？会改什么？按什么顺序改？

<details>
<summary><b>Answer · 参考答案</b></summary>

**Measure cost per completed task first.** The unit is a resolved ticket, not a request or a token. Cost per resolved ticket is the total model spend on all tickets, including escalated, abandoned and retried ones, divided by the number of tickets resolved without a human. This keeps you honest: a cheaper configuration that resolves fewer tickets can cost more per resolution. Log the `usage` of every call (uncached input, cache writes, cache reads, output) tagged with ticket ID, step, model and effort. For latency, build a waterfall per ticket: time to first token and generation time per model call, time per tool call, and time in queues.

> **先衡量每个完成任务的成本。** 计量单位是一张已解决的工单，而不是一次请求或一个 token。每张已解决工单的成本 = 所有工单（包括转人工、被放弃和重试过的）的模型总花费，除以无需人工就解决的工单数。这个指标能让你保持诚实：一个更便宜、但解决工单更少的配置，摊到每次解决上可能反而更贵。记录每次调用的 `usage`（未缓存输入、缓存写入、缓存读取、输出），并打上工单 ID、步骤、模型和 effort 标签。延迟方面，为每张工单画一张瀑布图：每次模型调用的首 token 时间和生成时间、每次工具调用的耗时，以及在队列中等待的时间。

**Find where the tokens go.** In a tool loop, every call resends the whole prefix: tools, system prompt and the growing history. Here that is 12,000 tokens of policy, 40 tool schemas and every 30 KB tool result so far, sent up to 14 times. Because each turn carries all earlier turns, total input grows roughly with the square of the number of turns. Most of the bill is usually re-sent input, which is exactly what the first levers act on.

> **找出 token 花在了哪里。** 在工具循环中，每次调用都要重新发送整个前缀：工具、system prompt，以及不断增长的历史。在这里，就是 12,000 token 的政策、40 个工具 schema，再加上到目前为止每一个 30 KB 的工具结果，最多要发送 14 次。由于每一轮都带着之前所有轮次的内容，总输入量大致随轮数的平方增长。账单的大头通常就是这些被反复发送的输入，而这正是前几个手段要解决的。

**Lever 1: prompt caching.** The cache matches the exact bytes of the prefix, rendered in the order tools, system, messages. Keep the tools and the policy prompt byte-stable at the front with a cache breakpoint after them, and another on the conversation tail. Cache reads cost about a tenth of normal input on most models; writes with the default 5-minute TTL cost about 1.25x. Remove silent invalidators: a timestamp or ticket ID in the system prompt, non-deterministic JSON key order, per-merchant data mixed into the policies (move it to the first user message), or tools that change order between requests. Then verify in `usage`; deep in a healthy loop, most input is read from the cache (illustrative numbers):

> **手段一：prompt caching。** 缓存按前缀的精确字节进行匹配，渲染顺序是 tools、system、messages。让工具定义和政策 prompt 放在最前面并保持字节级稳定，在它们之后设置一个缓存断点，再在对话末尾设置一个。在大多数模型上，缓存读取的价格约为普通输入的十分之一；使用默认 5 分钟 TTL 的缓存写入约为 1.25 倍。清除那些悄悄让缓存失效的东西：system prompt 里的时间戳或工单 ID、不确定的 JSON key 顺序、混在政策里的商家专属数据（把它移到第一条 user 消息里），或者在不同请求之间顺序会变的工具列表。然后用 `usage` 来验证；在一个健康的循环深处，大部分输入都是从缓存读取的（数字仅作示意）：

```json
"usage": {
  "input_tokens": 1200,
  "cache_creation_input_tokens": 2100,
  "cache_read_input_tokens": 38500,
  "output_tokens": 450
}
```

**Lever 2: trim context and tool outputs.** Return only the fields the model needs, not raw 30 KB responses: project, paginate, and summarize lists in code ("412 products disapproved: 380 missing GTIN, 32 price mismatch"). Give list tools `limit` and `fields` parameters. If many of the 40 tools are rarely used, load them on demand with tool search, or route each ticket type to a smaller tool set. Delete system-prompt prose that only repeats the tool schemas.

> **手段二：精简上下文和工具输出。** 只返回模型需要的字段，而不是 30 KB 的原始响应：在代码里做字段投影、分页，并把列表汇总成摘要（“412 个商品被拒登：380 个缺少 GTIN，32 个价格不一致”）。给列表类工具加上 `limit` 和 `fields` 参数。如果 40 个工具中有很多很少用到，就用 tool search 按需加载，或者按工单类型路由到更小的工具集。删掉 system prompt 中那些只是在重复工具 schema 的文字。

**Lever 3: fewer round trips.** Every model call adds latency and resends the context. Let the model request independent lookups in parallel (feed status, account health and item issues in one response), run them concurrently, and return all the results in one user message. When the model must loop over many items, use programmatic tool calling: it writes a script that calls the tools in a sandbox, and only the summary comes back. And if every feed ticket needs the same three lookups, run them in code before the first model call; where the path is known, a fixed workflow beats an agent.

> **手段三：减少往返次数。** 每一次模型调用都会增加延迟，并重新发送上下文。让模型并行发起互不依赖的查询（在一次响应里同时查询 feed 状态、账户健康度和商品问题），并发执行，然后把所有结果放在一条 user 消息里返回。当模型需要遍历大量条目时，使用 programmatic tool calling：模型写一个脚本，在沙箱里调用工具，只有汇总结果会返回给模型。另外，如果每张 feed 工单都需要同样的三次查询，就在第一次调用模型之前用代码完成；在路径已知的地方，固定的工作流胜过 Agent。

**Lever 4: effort and thinking.** Lower `effort` means less thinking and fewer, more consolidated tool calls, which saves both tokens and time. Sweep the levels on an evaluation set of real tickets and choose one per ticket type. A useful pattern is to run at a lower effort and retry only the failures at a higher one.

> **手段四：effort 与思考。** 更低的 `effort` 意味着更少的思考，以及更少、更集中的工具调用，既省 token 又省时间。在由真实工单组成的评估集上把各个档位跑一遍，再为每种工单类型选定一个档位。一个实用的模式是：先用较低的 effort 运行，只对失败的工单用更高的 effort 重试。

**Lever 5: smaller models for well-scoped subtasks.** Classifying the ticket, summarizing a long error log or extracting fields can go to a smaller, faster model through a subagent, while the main loop stays on one model, because switching the main loop's model throws away its cache. Compare models on the hardest tickets, not the average ones: a cheap model that fails still pays for its tokens, then the retry, then the escalation.

> **手段五：范围明确的子任务交给更小的模型。** 工单分类、总结一段很长的错误日志、提取字段，这些都可以通过子 Agent 交给更小、更快的模型，而主循环始终使用同一个模型，因为切换主循环的模型会丢掉它的缓存。比较模型时要看最难的工单，而不是一般的工单：便宜的模型失败了，照样要为它的 token 付费，然后是重试的费用，再然后是转人工的成本。

**Lever 6: streaming for perceived latency.** Streaming makes nothing cheaper or shorter, but a user who sees "Checking your feed status..." and then an answer appearing within seconds is far happier than one watching a spinner for 90 seconds. For tickets that will take long anyway, switch to an asynchronous reply.

> **手段六：用流式输出改善感知延迟。** 流式输出不会让任务更便宜或更快完成，但一个看到“正在检查你的 feed 状态……”、几秒内就看到答案逐渐出现的用户，远比盯着转圈图标等 90 秒的用户满意。对于注定要花很久的工单，改成异步回复。

**Lever 7: the Batch API for offline work.** Nightly feed-health reports, QA summaries of closed tickets and evaluation runs need no immediate answer. Message batches cost about half the normal price, and results come back asynchronously and in any order, so match them by `custom_id`. Never put an interactive turn on it.

> **手段七：离线工作用 Batch API。** 每晚的 feed 健康报告、已关闭工单的质检摘要、评估跑批，都不需要立即得到结果。Message Batches 的价格约为正常价格的一半，结果以异步方式返回，顺序不定，所以要用 `custom_id` 来对应。千万不要把交互式的对话轮次放到上面。

**Order and trade-offs.** Change one lever at a time, and measure resolution rate, cost per resolved ticket and p95 together on the same evaluation set. Free wins come first (caching, trimming, fewer round trips) because they don't reduce intelligence; trade-offs come last (lower effort, smaller models) because they can raise escalations. Watch the side effects: trimming too hard causes extra lookup turns, a 1-hour cache TTL costs about 2x per write, and a subagent starts with a fresh, uncached prefix.

> **顺序与权衡。** 一次只调整一个手段，并在同一个评估集上同时衡量解决率、每张已解决工单的成本和 p95。先做“免费”的优化（缓存、精简、减少往返），因为它们不会降低智能水平；最后再做有代价的权衡（降低 effort、换小模型），因为它们可能导致更多转人工。还要注意副作用：精简过度会带来额外的查询轮次；1 小时 TTL 的缓存每次写入约为 2 倍价格；子 Agent 会从一个全新的、没有缓存的前缀开始。

</details>

### Follow-ups · 追问

- After you add caching, `cache_read_input_tokens` stays at zero. How do you find out why?
  > 加上缓存之后，`cache_read_input_tokens` 一直是零。你如何查出原因？
- Product wants answers in under 10 seconds for the three most common ticket types. What architecture would you propose?
  > 产品希望最常见的三类工单能在 10 秒内回复。你会提出什么样的架构？
- How would you decide whether a smaller model is good enough for ticket classification?
  > 你如何判断一个更小的模型是否足以胜任工单分类？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| cost per resolved ticket | 每张已解决工单的成本 |
| prompt caching | 提示词缓存 |
| round trip | 往返 |
| time to first token (TTFT) | 首 token 时间 |
| perceived latency | 感知延迟 |

---

<a id="q3"></a>

## 3. Peak traffic: 429s, overloaded errors, double refunds and restarted tasks

**流量高峰：429、过载错误、重复退款和从头再来的任务** · `Production` · `Hard`

### Question · 题目

On Black Friday, your agent platform (hundreds of merchants, each running support and catalog agents) starts failing. Logs show HTTP 429 `rate_limit_error` and 529 `overloaded_error` from the model API, tool calls to the order service timing out, and some customers receiving two refunds after retries.

> 黑色星期五当天，你们的 Agent 平台（数百个商家，每个商家都运行着客服 Agent 和商品目录 Agent）开始出故障。日志显示：模型 API 返回 HTTP 429 `rate_limit_error` 和 529 `overloaded_error`；调用订单服务的工具超时；一些顾客在重试之后收到了两笔退款。

Long catalog-rewrite tasks that take 40 minutes start again from the beginning whenever Kubernetes kills a pod during autoscaling. What is going wrong, and how would you make the platform reliable?

> 那些要跑 40 分钟的商品目录改写任务，每当 Kubernetes 在自动扩缩容时杀掉一个 pod，就会从头开始。哪里出了问题？你会如何让这个平台变得可靠？

<details>
<summary><b>Answer · 参考答案</b></summary>

**The core idea.** An agent run is a long, stateful workflow with a rate-limited, sometimes overloaded dependency in the middle and side effects at the edges. Treat it like any workflow engine: retry only what is safe to retry, make side effects idempotent, persist progress, control concurrency, and decide in advance how to degrade.

> **核心思路。** 一次 Agent 运行，是一个长时间运行、有状态的工作流：中间依赖着一个会限流、有时还会过载的服务，边缘还有各种副作用。像对待任何工作流引擎一样对待它：只重试可以安全重试的东西，让副作用幂等，持久化进度，控制并发，并提前决定如何降级。

**Retry with backoff and jitter, but only what is retryable.** Retry 429 (honor the `retry-after` header), 529, 500 and other 5xx errors, 408, 409, and connection errors. Don't retry 400, 401, 403, 404 or 413: the same request will fail again, so fix the request, the credentials or its size. A `stop_reason` of `refusal` is not a transport error either. Use exponential backoff with full jitter, so thousands of workers don't retry in lockstep, and cap both attempts and total time:

> **用退避加抖动来重试，但只重试可重试的错误。** 可以重试的有：429（遵守 `retry-after` 响应头）、529、500 及其他 5xx 错误、408、409，以及连接错误。不要重试 400、401、403、404 或 413：同样的请求还会再失败，应该去修正请求本身、凭证或请求大小。`stop_reason` 为 `refusal` 也不是传输层错误。使用带完全抖动（full jitter）的指数退避，避免成千上万个 worker 步调一致地同时重试，并且同时限制重试次数和总耗时：

```go
func retryable(status int) bool {
	switch status {
	case 408, 409, 429, 500, 502, 503, 504, 529:
		return true
	}
	return false
}

func backoff(attempt int, retryAfter time.Duration) time.Duration {
	if retryAfter > 0 {
		return retryAfter // the server said when to come back
	}
	ceiling := min(30*time.Second, 500*time.Millisecond<<attempt)
	return time.Duration(rand.Int64N(int64(ceiling))) // full jitter
}
```

**Avoid retry storms.** Retries multiply. The official SDKs already retry twice by default, your agent loop may retry, and the job queue may retry the whole task: three layers of three attempts is 27 calls for one failure, at the exact moment the API is overloaded. Let one layer own retries, set the others to zero, and add a retry budget or circuit breaker that stops retrying when the error rate is high. Note also that retrying a model call is safe because it changes nothing, while retrying a tool call that moves money is not, unless it carries an idempotency key.

> **避免重试风暴。** 重试会层层相乘。官方 SDK 默认已经会重试两次，你的 Agent 循环可能还会重试，任务队列又可能把整个任务重试一遍：三层、每层三次尝试，一次失败就会变成 27 次调用，而且恰好发生在 API 已经过载的时刻。只让一层负责重试，其他层设为零，并加上重试预算或熔断器，在错误率高时停止重试。还要注意：重试模型调用是安全的，因为它不会改变任何东西；而重试一个会动钱的工具调用并不安全，除非它带着幂等键。

**Idempotency keys for side-effecting tools.** The double refunds come from a classic race: the refund timed out on the client but succeeded on the server, and the retry (by your code, or by the model after seeing an error) refunded again. Generate the key once, when the action is first decided, from the business intent (for example `refund:{ticket_id}:{order_id}`), and save it in the task state before calling. Not a random value per attempt, and not the `tool_use` ID, which changes if the turn is regenerated. The order service keeps the key under a unique constraint and returns the original result for a repeat; after a timeout, the harness looks the key up before trying again.

> **有副作用的工具要用幂等键。** 重复退款来自一个经典的竞态：退款在客户端超时了，但在服务端其实已经成功；随后的重试（无论是你的代码发起的，还是看到报错的模型发起的）又退了一次款。幂等键要在第一次决定执行这个动作时，根据业务意图生成一次（例如 `refund:{ticket_id}:{order_id}`），并在调用之前保存到任务状态里。不能每次尝试都用随机值，也不能用 `tool_use` ID，因为这一轮如果重新生成，这个 ID 就会变。订单服务用唯一约束保存幂等键，遇到重复请求时返回原来的结果；超时之后，harness 先按幂等键查询，再决定是否重试。

**Checkpoint agent state so long tasks resume.** The state is the message history (with every `tool_use` and `tool_result` block), a progress cursor, pending actions with their idempotency keys, and counters and budgets. Persist it after every model response and tool result, in a durable store such as Postgres or object storage. When a pod dies, another worker loads the checkpoint and continues from the last completed step; a step that started but was never confirmed is reconciled by its idempotency key, not blindly re-run. For the catalog task, also process products in chunks and save a cursor, so a crash costs one chunk, not 40 minutes. On `SIGTERM`, stop at a step boundary, checkpoint and release the task. A durable workflow engine such as Temporal provides much of this.

> **为 Agent 状态做检查点，让长任务可以续跑。** 状态包括：消息历史（含每一个 `tool_use` 和 `tool_result` 块）、进度游标、待执行的动作及其幂等键，以及各种计数器和预算。在每次模型响应和工具结果之后，把状态持久化到可靠的存储中，比如 Postgres 或对象存储。当 pod 挂掉时，另一个 worker 加载检查点，从最后一个已完成的步骤继续；对于已经开始但从未确认完成的步骤，要通过幂等键去核对，而不是盲目重跑。对于商品目录任务，还要把商品分块处理并保存游标，这样一次崩溃只损失一个分块，而不是 40 分钟。收到 `SIGTERM` 时，在步骤边界停下来，保存检查点并释放任务。像 Temporal 这样的持久化工作流引擎提供了其中的很多能力。

**Queue-based execution with per-tenant limits.** Don't run agent loops inside HTTP handlers. Queue the tasks and let workers lease them with a heartbeat, so a dead worker's task returns to the queue. Limit concurrency globally (sized to your model API's request and token rate limits), per tenant (so one merchant's 50,000-product catalog job cannot starve everyone's support tickets) and per priority (interactive work before batch work). Watch the rate-limit headers on each response and slow down before you hit 429s, not after.

> **基于队列执行，并按租户限流。** 不要在 HTTP handler 里直接跑 Agent 循环。把任务放进队列，worker 以带心跳的租约方式领取任务，这样挂掉的 worker 手里的任务会回到队列。在三个层面限制并发：全局（按模型 API 的请求数和 token 数速率限制来定）、租户（让某个商家 5 万个商品的目录任务无法饿死所有人的客服工单）和优先级（交互式工作先于批处理工作）。关注每个响应中与限流相关的响应头，在撞上 429 之前就主动放慢，而不是撞上之后。

**Graceful degradation.** Decide in advance what happens when capacity is short. Move catalog work off-peak or to the Batch API. If the primary model stays overloaded, route new interactive tasks to a fallback model you have already evaluated, and switch at task boundaries, because a mid-conversation switch loses the prompt cache and can change behavior. Drop optional steps such as extra verification passes, answer common questions from templates, and promise a reply by email. If the order service keeps failing, open a circuit breaker and tell the agent the tool is unavailable, so it hands off to a human instead of looping.

> **优雅降级。** 提前决定容量不足时该怎么做。把商品目录任务挪到低峰时段，或者放到 Batch API 上。如果主模型持续过载，就把新的交互式任务路由到一个事先评估过的备用模型，并且在任务边界上切换，因为在对话中途切换会丢掉 prompt 缓存，还可能改变行为。去掉可选步骤（比如额外的校验轮次），常见问题用模板回答，并承诺通过邮件回复。如果订单服务持续失败，就打开熔断器，并告诉 Agent 这个工具暂时不可用，让它转交给人工，而不是原地打转。

**Timeouts at every layer.** Set a timeout on each model call (stream long outputs and use an idle timeout between events), each tool call, each agent step and the whole task, plus a maximum number of loop iterations. Inner timeouts must be shorter than outer ones, so the outer layer never kills work that was about to succeed. In Go, pass `context.Context` deadlines all the way down, and remember that a zero-value `http.Client` has no timeout at all.

> **每一层都要有超时。** 为每次模型调用（长输出要用流式，并在事件之间设置空闲超时）、每次工具调用、每个 Agent 步骤以及整个任务都设置超时，再加上循环次数的上限。内层超时必须短于外层超时，这样外层永远不会杀掉一个马上就要成功的工作。在 Go 里，把 `context.Context` 的 deadline 一路传下去，并且记住：零值的 `http.Client` 根本没有超时。

**Prove it before the next peak.** Inject faults in staging (kill pods mid-task, return 529s, slow down tools) and check that tasks resume and no action happens twice. Dashboard errors by status code, retry rate, queue depth per tenant, resumed tasks and duplicate side effects (target: zero). Load-test at the expected peak, and request higher rate limits well ahead of time.

> **在下一次高峰之前证明它可靠。** 在预发环境中注入故障（在任务中途杀掉 pod、返回 529、拖慢工具），然后检查任务能否续跑、有没有任何动作被执行两次。在看板上展示按状态码统计的错误、重试率、各租户的队列深度、续跑的任务数，以及重复副作用的次数（目标：零）。按预期峰值做压测，并提前申请更高的速率限制。

</details>

### Follow-ups · 追问

- The refund API has no idempotency support and you cannot change it. How do you still prevent double refunds?
  > 退款 API 不支持幂等，而且你改不了它。你还能如何防止重复退款？
- One merchant has 50 times the traffic of the median merchant. How do you set its concurrency limit fairly?
  > 有一个商家的流量是中位数商家的 50 倍。你如何公平地设置它的并发上限？
- Your fallback model resolves fewer tickets correctly. How do you decide between switching to it and queuing requests for the primary model?
  > 你的备用模型正确解决的工单更少。你如何在“切换到备用模型”和“排队等待主模型”之间做决定？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| exponential backoff with jitter | 带抖动的指数退避 |
| retry storm | 重试风暴 |
| checkpoint | 检查点 |
| graceful degradation | 优雅降级 |
| noisy neighbor | 吵闹邻居（某个租户挤占共享资源） |

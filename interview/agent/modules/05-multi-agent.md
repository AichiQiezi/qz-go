# 05. Planning and multi-agent systems · 规划与多 Agent 系统

How agents plan and track multi-step work, when to split a job across an orchestrator and subagents, and how to design a multi-agent support system with handoffs, escalation and safe side effects. Interviewers ask because long agent tasks usually fail through forgotten requirements, lost context and unclear ownership, not through one bad model call.

> Agent 如何规划和跟踪多步骤工作，什么时候把任务拆给编排者和子 agent，以及如何设计一个带交接、人工升级和安全副作用控制的多 Agent 客服系统。面试官爱问这些，是因为长任务中的 agent 失败，通常源于要求被遗忘、上下文丢失和职责不清，而不是某一次模型调用出错。

3 questions · generated 2026-10-15T18:27:03+08:00 · [← 目录 / Contents](../README.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Multi-agent | Medium | [Keeping a multi-step campaign launch on track with an explicit plan](#q1) | 用显式计划让多步骤的活动上线不跑偏 |
| 2 | Multi-agent | Hard | [Fanning out a 200-merchant ROAS investigation to subagents](#q2) | 把 200 个商家的 ROAS 下滑调查扇出给子 agent |
| 3 | Multi-agent | Senior | [Designing a merchant support system with specialists, handoffs and human escalation](#q3) | 设计商家客服系统：专职 agent、交接与人工升级 |

---

<a id="q1"></a>

## 1. Keeping a multi-step campaign launch on track with an explicit plan

**用显式计划让多步骤的活动上线不跑偏** · `Multi-agent` · `Medium`

### Question · 题目

Your team is building an agent that runs marketing tasks for merchants. A typical request: "Launch our Black Friday campaign: pick our 20 best sellers that have enough stock, write ad copy in English, German and Japanese, and create campaigns on Google Ads and Meta with a total budget of $3,000."

> 你的团队在做一个替商家执行营销任务的 agent。一个典型请求是：“上线我们的黑五活动：挑出库存充足的 20 个畅销商品，用英语、德语和日语写广告文案，并在 Google Ads 和 Meta 上创建广告系列，总预算 3,000 美元。”

The first version is a plain ReAct loop: the model reasons, calls one tool, reads the result and decides the next step. In testing it sometimes forgets the Japanese copy, creates the Meta campaign twice after a timeout, and carries on happily when the product query returns only 8 eligible products. A teammate says: "Just make it write a plan first."

> 第一版是一个朴素的 ReAct 循环：模型思考、调用一个工具、读取结果，再决定下一步。测试中它有时会漏掉日语文案，在一次超时后把 Meta 广告系列创建了两次，而且当商品查询只返回 8 个符合条件的商品时，它照样若无其事地继续往下做。一位同事说：“让它先写个计划就行了。”

Compare plan-then-execute, ReAct-style interleaving and re-planning after each step. How would you make the plan explicit, recover when a step fails, and keep the plan visible and checkable?

> 请比较先计划后执行（plan-then-execute）、ReAct 式的边想边做，以及每一步之后都重新规划这三种方式。你会如何让计划显式化、在某一步失败时如何恢复，以及如何让计划可见、可检查？

<details>
<summary><b>Answer · 参考答案</b></summary>

**The core idea.** Planning means deciding the order of steps toward a goal. ReAct-style interleaving (Yao et al., 2022) alternates reasoning and acting: think, call a tool, observe, repeat. It handles surprises well, but its only memory of the goal is the transcript, so on a long run the original requirements get buried under tool output. Plan-then-execute writes the whole plan first and then runs it: predictable and easy to review, but brittle, because the plan was made before the agent knew that only 8 products qualify. Re-planning after every step revises the remaining steps after each result: the most adaptive option, and the most expensive. The strong production pattern is a hybrid: an explicit plan kept as data, executed step by step, with re-planning triggered by events.

> **核心思路。** 规划就是决定朝目标前进的步骤顺序。ReAct 式的交替执行（Yao et al., 2022）让推理和行动轮流进行：思考、调用工具、观察结果，然后重复。它很能应对意外，但它对目标唯一的“记忆”就是对话记录，所以在长任务里，最初的要求会被大量工具输出淹没。先计划后执行（plan-then-execute）先写出完整计划再执行：可预测、便于审查，但很脆弱，因为计划是在 agent 知道只有 8 个商品符合条件之前就定下的。每一步之后都重新规划，会根据每个结果修改剩余的步骤：这是最灵活的方式，也是最贵的。生产中效果好的是混合式：把显式计划当作数据来维护，逐步执行，并由事件触发重新规划。

**Why the plain loop failed.** Each bug has a different cause:

> **朴素循环为什么会失败。** 每个 bug 的原因都不一样：

- Forgotten Japanese copy: the requirement exists only in the first user message. Thirty tool calls later it is far back in the context, and nothing forces a check.
  > 漏掉日语文案：这个要求只出现在第一条用户消息里。三十次工具调用之后，它已经位于上下文很靠前的位置，没有任何机制强制 agent 去核对。
- Duplicate Meta campaign: a timeout looked like a failure, and the model retried a create call that is not idempotent. A plan does not fix this; idempotency does.
  > Meta 广告系列被重复创建：一次超时看起来像失败，于是模型重试了一个不幂等的创建调用。计划解决不了这个问题，要靠幂等性。
- Carrying on with 8 products: the step had no success criterion, so nobody noticed that a precondition had failed.
  > 只有 8 个商品却继续往下做：这一步没有成功标准，所以没人注意到前置条件已经不成立了。

**Make the plan a first-class object.** Give the agent an `update_plan` tool and keep the plan in your harness, not only in the model's head. Each step has a status, an acceptance check and the outputs it produced, such as the IDs of created objects:

> **让计划成为一等公民。** 给 agent 一个 `update_plan` 工具，并把计划保存在你的 harness 里，而不只是在模型的“脑子”里。每一步都有状态、验收检查，以及它产出的结果，比如已创建对象的 ID：

```json
{
  "goal": "Black Friday: 20 best sellers, copy in en/de/ja, Google Ads + Meta, budget 3000 USD",
  "steps": [
    {"id": "s1", "task": "Select best sellers with stock >= 50", "status": "failed",
     "check": "20 products selected", "result": "only 8 products qualify"},
    {"id": "s2", "task": "Write copy in en, de, ja for each product", "status": "pending",
     "check": "3 languages per product"},
    {"id": "s3", "task": "Create Google Ads campaign (paused)", "status": "pending",
     "check": "campaign_id stored"},
    {"id": "s4", "task": "Create Meta campaign (paused)", "status": "pending",
     "check": "campaign_id stored"}
  ]
}
```

**Let code enforce what code can check.** The harness re-injects the current plan near the end of the context on every turn, where the model attends to it; appending it to the latest message rather than editing the system prompt keeps the prompt cache valid. The harness also validates transitions: a step becomes `done` only when its check passes, and the check runs in code where possible (count the copies per language, ask the ad API whether the campaign exists) instead of trusting the model's claim. Coding agents such as Claude Code keep a visible todo list for the same reason.

> **能用代码检查的，就交给代码强制执行。** harness 每一轮都把当前计划重新注入到上下文末尾附近，那是模型最关注的位置；把它附加到最新的消息上，而不是修改 system prompt，可以让提示缓存保持有效。harness 还会校验状态流转：只有验收检查通过，步骤才能变成 `done`，而且检查尽量在代码里执行（按语言统计文案数量，向广告 API 确认广告系列是否存在），而不是相信模型自己的说法。Claude Code 这类编程 agent 维护一个可见的 todo 列表，也是出于同样的原因。

**Re-plan by failure type.** Classify the failure before reacting:

> **按失败类型重新规划。** 先给失败分类，再做出反应：

- Transient (timeout, HTTP 429, 5xx): retry the same step with backoff. Before retrying a create, check whether the object already exists; if the platform has no idempotency key, put your own request ID in the object's name and search for it.
  > 临时性失败（超时、HTTP 429、5xx）：带退避地重试同一步。重试创建操作之前，先检查对象是否已经存在；如果平台没有幂等键，就把你自己的请求 ID 写进对象名称，并按它搜索。
- Broken precondition (only 8 products): later steps rest on a wrong assumption, so re-plan them. Within what the user asked for, re-plan on your own and record why; if scope, budget or spend would change, ask: "Only 8 products have enough stock. Lower the threshold to 20 units, or launch with these 8?"
  > 前置条件不成立（只有 8 个商品）：后续步骤建立在错误的假设上，所以要重新规划它们。如果仍在用户要求的范围内，就自行重新规划并记录原因；如果范围、预算或花费会改变，就去问商家：“只有 8 个商品库存充足。是把门槛降到 20 件，还是先用这 8 个商品上线？”
- Permanent (Meta ad account disabled): mark the step failed, continue with independent steps such as Google Ads, and report the partial result clearly.
  > 永久性失败（Meta 广告账户被停用）：把这一步标记为失败，继续执行 Google Ads 这类相互独立的步骤，并清楚地报告部分完成的结果。

Re-plan only the remaining steps. Never re-run completed steps that had side effects, and keep their outputs.

> 只重新规划剩余的步骤。绝不要重跑已经完成、带副作用的步骤，并且要保留它们的输出。

**When to use which.** Plain ReAct is fine for short, exploratory tasks such as "why did spend drop yesterday?", where the next step depends on what the data shows. An explicit plan pays off with many steps, requirements that are easy to forget, side effects, steps that can run in parallel, or an approval before money is spent. If the steps never change, do not let the model plan at all: hard-code the workflow and use the model inside the steps. That is the workflow-versus-agent distinction in Anthropic's "Building effective agents".

> **什么时候用哪种。** 对于短小的、探索性的任务，比如“昨天花费为什么下降了？”，朴素的 ReAct 就够了，因为下一步做什么取决于数据显示了什么。当步骤很多、要求容易被遗忘、有副作用、有可以并行的步骤，或者花钱前需要审批时，显式计划才物有所值。如果步骤从不变化，就根本不要让模型做规划：把工作流写死在代码里，只在步骤内部使用模型。这就是 Anthropic《Building effective agents》中 workflow 与 agent 的区别。

**Keep the plan visible and checkable.** Before any step that spends money, show the plan ("I will create 2 paused campaigns with a total budget of 3,000 USD") and wait for approval. Create campaigns paused and enable them after review. Record every plan version in the trace, so you can see when and why it changed. End with a verification step in code against the original goal: 3 languages times 8 products is 24 pieces of copy, both campaigns exist, and the budgets add up to at most 3,000 USD.

> **让计划可见、可检查。** 在任何要花钱的步骤之前，先展示计划（“我将创建 2 个处于暂停状态的广告系列，总预算 3,000 美元”）并等待批准。广告系列先以暂停状态创建，审核后再启用。在 trace 里记录计划的每一个版本，这样你能看清它在何时、因何而改变。最后在代码里用一个验证步骤对照最初的目标：3 种语言乘以 8 个商品就是 24 份文案、两个广告系列都已存在、预算合计不超过 3,000 美元。

**Trade-offs.** Plans cost tokens and latency, a rigid plan can keep the agent on outdated steps, and a model can mark steps done without doing them, which is why the checks belong in code. Start with a lightweight todo list plus code checks, and add heavier re-planning only where real failures call for it.

> **权衡。** 计划会消耗 token、增加延迟；死板的计划会让 agent 执行过时的步骤；模型也可能没做就把步骤标记为完成，这正是检查应该放在代码里的原因。先从轻量的 todo 列表加代码检查开始，只在真实的失败表明确有需要的地方，才引入更重的重新规划。

</details>

### Follow-ups · 追问

- How would you let the merchant edit the plan halfway through without corrupting the agent's state?
  > 如何让商家在执行到一半时修改计划，而不破坏 agent 的状态？
- Where in the prompt would you put the current plan so it stays visible without breaking prompt caching?
  > 你会把当前计划放在提示词的什么位置，既让它保持可见，又不破坏提示缓存？
- How would you test that the agent re-plans correctly when a step fails?
  > 你会如何测试 agent 在某一步失败时能否正确地重新规划？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| plan-then-execute | 先计划后执行 |
| re-planning | 重新规划 |
| acceptance check | 验收检查 |
| precondition | 前置条件 |
| The plan lives in the harness, not only in the model's head. | 计划保存在 harness 里，而不只是在模型的“脑子”里。 |

---

<a id="q2"></a>

## 2. Fanning out a 200-merchant ROAS investigation to subagents

**把 200 个商家的 ROAS 下滑调查扇出给子 agent** · `Multi-agent` · `Hard`

### Question · 题目

An internal analytics agent answers questions like "Why did ROAS drop for these 200 merchants last week?" It is a single agent with tools for ad-platform reports, shop orders, catalog sync logs and tracking-pixel health. A run takes about 40 minutes, often fills its context window (or gets compacted and forgets earlier findings), and the final report covers only the first 60 or so merchants in any depth.

> 一个内部分析 agent 负责回答这类问题：“上周这 200 个商家的 ROAS 为什么下降了？”它是单个 agent，带有查询广告平台报表、店铺订单、商品目录同步日志和追踪像素健康状况的工具。一次运行大约要 40 分钟，经常把上下文窗口塞满（或者被压缩后忘掉之前的发现），最终报告只对前 60 个左右的商家做了深入分析。

Redesign it with an orchestrator and subagents. Explain how context isolation and parallel fan-out help, what a good task spec for a subagent contains, how you merge results and handle failures, and what this costs. When would you keep a single agent instead?

> 请用编排者（orchestrator）加子 agent（subagent）的方式重新设计它。解释上下文隔离和并行扇出（fan-out）为什么有帮助、一份好的子 agent 任务说明应该包含什么、你如何合并结果和处理失败，以及这样做的代价是什么。什么情况下你会继续使用单个 agent？

<details>
<summary><b>Answer · 参考答案</b></summary>

**The core idea.** In the orchestrator-worker pattern, a lead agent (the orchestrator) splits the job into subtasks and starts worker agents (subagents), each with a fresh context and a narrow brief, and it receives only their condensed results. There are two benefits. Context isolation: a subagent reads hundreds of kilobytes of raw reports and returns a kilobyte of findings, so the orchestrator holds conclusions, not raw rows. Parallelism: independent subtasks run at the same time. In backend terms, it is scatter-gather with models as the workers.

> **核心思路。** 在编排者-工作者（orchestrator-worker）模式中，一个主 agent（编排者）把任务拆成子任务，启动若干工作 agent（子 agent），每个子 agent 都有全新的上下文和一份范围很窄的任务说明，编排者只接收它们浓缩后的结果。好处有两个。上下文隔离：子 agent 读几百 KB 的原始报表，只返回 1 KB 的结论，于是编排者手里是结论，而不是原始数据行。并行：相互独立的子任务同时运行。用后端的话说，这就是以模型为工作节点的 scatter-gather。

**Why the single agent struggles.** Every report it fetches stays in its context and is re-sent on every later turn. By merchant 60 the window is full of raw data, quality drops as irrelevant tokens pile up, compaction loses details, and the tool calls run one after another. A bigger window only delays this.

> **单个 agent 为什么吃力。** 它拉取的每份报表都会留在上下文里，并在之后的每一轮被重新发送。到第 60 个商家时，窗口里已经塞满了原始数据；无关 token 越堆越多，质量随之下降，压缩会丢掉细节，而工具调用是一个接一个串行执行的。更大的窗口只能推迟这个问题。

**Decompose before you fan out.** Do not start 200 agents. First run cheap, deterministic triage in SQL: ROAS change per merchant, split into spend versus attributed revenue, by platform and country. Then group merchants by symptom, because many share one cause, such as a pixel release, a catalog sync outage or an attribution change. Fan out by symptom group or hypothesis. Arithmetic belongs in code, not in the model's head.

> **先拆解，再扇出。** 不要一上来就启动 200 个 agent。先用 SQL 做一次便宜的、确定性的初筛：计算每个商家的 ROAS 变化，拆成花费和归因收入两部分，再按平台和国家细分。然后按症状给商家分组，因为很多商家是同一个原因，比如一次像素发布、一次目录同步故障或一次归因规则变化。再按症状分组或按假设来扇出。算术交给代码，而不是让模型心算。

**Write the task spec like an API contract.** A subagent knows only what is in its brief. Vague briefs ("look into tracking") produce duplicated work, gaps and invented causes. A good spec states:

> **像写 API 契约一样写任务说明。** 子 agent 只知道任务说明里写了什么。含糊的说明（“查一下追踪问题”）会导致重复劳动、遗漏和编造出来的原因。一份好的任务说明应写明：

- Goal: the one question to answer, and why it matters.
  > 目标：要回答的那一个问题，以及它为什么重要。
- Inputs: exact merchant IDs, both date ranges with a time zone, and what is already known.
  > 输入：确切的商家 ID、带时区的两个日期区间，以及已经知道的信息。
- Tools and budget: which tools to use, and a cap on tool calls or time.
  > 工具和预算：用哪些工具，以及工具调用次数或时间的上限。
- Output: a schema with evidence and an explicit "inconclusive" option, so it never has to invent a cause.
  > 输出：一个带证据的 schema，并显式提供“无法确定”（inconclusive）选项，这样它永远不必编造原因。
- Boundaries: read-only, only these merchants, and stop once the hypothesis is confirmed or rejected.
  > 边界：只读、只看这些商家，一旦假设被证实或被否定就停止。

```json
{
  "goal": "Decide whether broken conversion tracking explains the ROAS drop for these merchants",
  "merchant_ids": ["m_1042", "m_1077", "m_2210"],
  "period": {"current": "2026-10-05/2026-10-11", "baseline": "2026-09-28/2026-10-04", "tz": "UTC"},
  "known": "Spend is flat; platform-reported conversions fell; shop orders did not fall",
  "tools_allowed": ["get_pixel_events", "get_capi_delivery_stats", "get_orders_daily"],
  "limits": {"max_tool_calls": 30, "read_only": true},
  "output": {
    "per_merchant": [{"merchant_id": "string", "verdict": "tracking_broken | not_tracking | inconclusive",
                      "evidence": "numbers and IDs", "since": "date or null"}],
    "summary": "max 150 words"
  }
}
```

**Run the fan-out like a worker pool.** The orchestrator can request several subagents in one turn, and the harness runs them concurrently with a limit, because the model API and the ad-platform APIs both have rate limits. Each subagent gets a timeout and a token budget and can use a cheaper, faster model, while the orchestrator stays on one model (switching its model would invalidate its prompt cache).

> **像 worker pool 一样运行扇出。** 编排者可以在一轮里请求多个子 agent，harness 带并发上限地同时运行它们，因为模型 API 和广告平台 API 都有限流。每个子 agent 都有超时和 token 预算，并且可以使用更便宜、更快的模型，而编排者始终使用同一个模型（切换它的模型会让它的提示缓存失效）。

```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(8) // respect model and ad-platform rate limits
results := make([]SubResult, len(specs))
for i, spec := range specs {
    g.Go(func() error {
        ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
        defer cancel()
        r, err := runSubagent(ctx, spec)
        if err != nil {
            r = SubResult{SpecID: spec.ID, Status: "failed", Err: err.Error()}
        }
        results[i] = r
        return nil // one failed worker must not cancel its siblings
    })
}
_ = g.Wait()
```

**Merge in code first, then let the model write.** Group the verdicts and count merchants per cause in code. Check coverage: every merchant ID appears in exactly one result, and missing ones are retried or listed as "not analysed", never silently dropped. Treat subagent output as claims: require evidence so a verifier can spot-check it, and settle conflicting verdicts with a targeted follow-up. Keep full outputs outside the context (a file or a table) and pass references. Persist the orchestrator's plan the same way, so compaction cannot lose it.

> **先在代码里合并，再让模型来写。** 在代码里对结论分组，统计每个原因涉及的商家数。检查覆盖率：每个商家 ID 恰好出现在一个结果里，缺失的要么重试，要么列为“未分析”，绝不悄悄丢掉。把子 agent 的输出当作“主张”：要求附上证据，方便校验者抽查；结论相互矛盾时，用一次有针对性的追加调查来解决。把完整输出放在上下文之外（文件或数据表），只传引用。编排者的计划也用同样的方式持久化，这样压缩就不会把它弄丢。

**Handle failure like a job system.** Retry read-only subtasks with the same spec, and never delegate side effects without idempotency. After repeated failure, mark the group inconclusive and report the real coverage instead of guessing. Cap tool calls to stop runaway subagents, and trace each subagent as a child span.

> **像任务调度系统一样处理失败。** 只读子任务可以用同一份任务说明重试；没有幂等性保障，就绝不委派带副作用的操作。多次失败后，把这一组标记为无法确定，并报告真实的覆盖范围，而不是去猜。限制工具调用次数，防止子 agent 失控，并把每个子 agent 记录为一个子 span。

**The costs.** Token multiplication: every subagent pays for its own system prompt, tools and exploration. Anthropic's June 2025 write-up on its multi-agent research system reported that agents use about 4x the tokens of chat, and multi-agent systems about 15x. Duplicated work: overlapping briefs fetch the same reports, so draw clear boundaries and share a cache of tool results keyed by merchant and date. Lost context at handoff: a subagent knows nothing the brief leaves out, and the orchestrator sees only what the subagent reports. Debugging also gets harder, and latency is set by the slowest worker.

> **代价。** token 成倍增加：每个子 agent 都要为自己的 system prompt、工具和探索过程付费。Anthropic 在 2025 年 6 月介绍其多 agent 研究系统的文章中提到，agent 消耗的 token 约为普通聊天的 4 倍，多 agent 系统约为 15 倍。重复劳动：范围重叠的任务说明会拉取同样的报表，所以要划清边界，并共享一个按商家和日期做键的工具结果缓存。交接时丢失上下文：任务说明里没写的，子 agent 一概不知，而编排者只能看到子 agent 汇报的内容。调试也会更难，而且延迟取决于最慢的那个 worker。

**When a single agent is better.** Keep one agent when the steps are tightly coupled (debugging one merchant's tracking, where each finding decides the next check), when the output must be one coherent piece built on shared decisions (Cognition's 2025 post "Don't Build Multi-Agents" argues that parallel workers make conflicting implicit decisions), or when the task is small enough that the overhead dominates. Fan out when the work is breadth-first, the pieces are independent, and each piece needs far more reading than it returns, as with these 200 merchants.

> **什么时候单个 agent 更好。** 当步骤紧密耦合时（排查某一个商家的追踪问题，每个发现都决定下一步查什么），当产出必须是建立在共同决策之上的一个连贯整体时（Cognition 在 2025 年的文章《Don't Build Multi-Agents》认为，并行的 worker 会做出相互冲突的隐含决策），或者当任务小到额外开销占了大头时，就保留单个 agent。当工作是广度优先的、各部分相互独立、且每部分需要阅读的内容远多于它要返回的内容时，就扇出，就像这 200 个商家一样。

</details>

### Follow-ups · 追问

- How would you stop two subagents from fetching the same expensive report?
  > 你会如何避免两个子 agent 去拉取同一份代价很高的报表？
- Should subagents be allowed to start their own subagents? What limits would you set?
  > 应该允许子 agent 再启动自己的子 agent 吗？你会设置哪些限制？
- How would you show that the multi-agent version is actually better than the single agent?
  > 你会如何证明多 agent 版本确实比单 agent 版本更好？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| orchestrator-worker pattern | 编排者-工作者模式 |
| context isolation | 上下文隔离 |
| fan-out / fan-in | 扇出 / 扇入（分发与汇总） |
| task spec (brief) | 任务说明 |
| The subagent only knows what you put in its brief. | 子 agent 只知道你写进任务说明里的内容。 |

---

<a id="q3"></a>

## 3. Designing a merchant support system with specialists, handoffs and human escalation

**设计商家客服系统：专职 agent、交接与人工升级** · `Multi-agent` · `Senior`

### Question · 题目

Your company wants an AI support assistant for merchants, available in the admin dashboard and in chat. Most conversations fall into three areas: subscription billing (invoices, refunds, failed payments), catalog sync errors with Google Merchant Center and Meta, and ad account issues (disapproved ads, suspended accounts, a pixel that stopped firing). Some fixes have side effects: re-triggering a catalog sync, issuing a credit or refund, reconnecting an OAuth integration.

> 你们公司想为商家做一个 AI 客服助手，在管理后台和聊天窗口里都能使用。大多数对话集中在三个方面：订阅计费（发票、退款、扣款失败）、与 Google Merchant Center 和 Meta 之间的商品目录同步错误，以及广告账户问题（广告被拒、账户被封、像素不再触发）。有些修复操作带有副作用：重新触发目录同步、发放抵扣额度或退款、重新连接 OAuth 集成。

Design the system: triage, specialised agents, escalation to humans, shared state, and safety for actions with side effects. Explain what context moves with a conversation during a handoff, when you would use one agent with many tools instead of several agents, how you keep tone and policy consistent, how you evaluate the whole system, and how you would roll it out.

> 请设计这个系统：分诊、专职 agent、升级到人工、共享状态，以及带副作用操作的安全性。解释交接时哪些上下文会随对话一起转移，什么时候你会用一个带很多工具的 agent 而不是多个 agent，如何保持语气和政策的一致，如何评估整个系统，以及你会如何上线。

<details>
<summary><b>Answer · 参考答案</b></summary>

**Start with the shape, not the number of agents.** The default is one agent that owns the conversation, with one shared prompt for tone and policy, tools grouped by domain, and domain instructions loaded only when needed (skills, or tool search with deferred loading). Split off a separate agent only where something really differs: conflicting instructions, different permissions, a different model or cost, a separate owning team, or too many tools to choose from reliably. Billing needs very different permissions from the other two areas, so I would use a front agent that owns the conversation plus three specialists.

> **先定形态，而不是先定 agent 的数量。** 默认做法是由一个 agent 负责整段对话：一份共享的提示词规定语气和政策，工具按领域分组，领域说明只在需要时加载（skills，或带延迟加载的工具搜索）。只有在某些东西确实不同时，才拆出独立的 agent：说明相互冲突、权限不同、模型或成本不同、归属于不同的团队，或者工具多到无法可靠地选择。计费需要的权限与另外两个领域差别很大，所以我会用一个负责对话的前台 agent，再加三个专职 agent。

**Two ways to involve a specialist.** With agent-as-tool, the front agent stays in charge, calls the specialist with a brief and gets a result back, so the merchant hears one voice. With a handoff, the specialist takes over the conversation, which suits a long back-and-forth such as reconnecting OAuth step by step. Prefer agent-as-tool. Frameworks name these differently (the OpenAI Agents SDK calls the second one a handoff), but the design question is the same.

> **引入专职 agent 的两种方式。** 在 agent-as-tool 方式中，前台 agent 始终掌控对话，带着任务说明调用专职 agent 并拿回结果，所以商家听到的始终是同一个“声音”。在交接（handoff）方式中，专职 agent 接管对话，适合需要长时间来回沟通的场景，比如一步步重新连接 OAuth。优先使用 agent-as-tool。不同框架的叫法不同（OpenAI Agents SDK 把第二种叫作 handoff），但设计上要回答的问题是一样的。

**Triage.** A cheap, fast model call with structured output classifies intent, language and urgency, with an explicit "unclear" class that triggers a clarifying question instead of a guess. Hard rules run in code first: chargebacks, legal threats, data-deletion requests and requests for a human go straight to a person.

> **分诊。** 用一次便宜、快速、带结构化输出的模型调用来识别意图、语言和紧急程度，并显式提供一个“不明确”类别，遇到它就提出澄清问题，而不是去猜。硬性规则先用代码执行：拒付（chargeback）、法律威胁、数据删除请求和要求转人工的请求，都直接交给人处理。

**Specialists with least privilege.** The billing agent reads invoices and payment status, and credits are gated. The catalog agent reads sync logs and item-level errors from Google Merchant Center and Meta, and can trigger an idempotent re-sync, rate-limited per merchant. The ad-account agent reads account status, disapproval reasons and pixel health, and never files an appeal without consent. The tool layer injects `merchant_id` from the authenticated session, so the model cannot choose whose data it touches: tenant isolation lives in code, not in the prompt.

> **按最小权限设计专职 agent。** 计费 agent 读取发票和支付状态，抵扣额度需要经过审批关卡。目录 agent 读取同步日志以及 Google Merchant Center 和 Meta 返回的商品级错误，可以触发一次幂等、按商家限流的重新同步。广告账户 agent 读取账户状态、拒登原因和像素健康状况，未经同意绝不提交申诉。工具层从已认证的会话中注入 `merchant_id`，所以模型无法选择要访问哪个商家的数据：租户隔离在代码里，而不在提示词里。

**Shared state lives in a case record, not in the transcript.** Like a ticket, the case holds the merchant ID, verified identity, language, category, collected facts, actions and their results, promises made and pending approvals; agents read and write it through tools. A handoff carries the case, a short summary written for the receiver, the merchant's last message verbatim and the suggested next step. It does not carry the raw transcript (noise, cost, and any injected text rides along), other merchants' data or internal reasoning. A human escalation gets the same packet, so the merchant never repeats themselves.

> **共享状态存放在工单记录里，而不是对话记录里。** 就像一张工单，case 保存商家 ID、已验证的身份、语言、类别、收集到的事实、操作及其结果、已做出的承诺和待审批事项；各个 agent 通过工具读写它。交接时带上这份 case、一段为接收方写的简短摘要、商家最后一条消息的原文和建议的下一步。不带原始对话记录（噪声多、成本高，被注入的文本也会跟着传过去），不带其他商家的数据，也不带内部推理。升级到人工时交接同一个数据包，商家就不用把情况再说一遍。

```json
{
  "case_id": "case_88231",
  "merchant_id": "m_5531",
  "language": "de",
  "category": "catalog_sync",
  "summary": "140 products rejected by GMC since Oct 12 (missing GTIN); merchant says GTINs were added yesterday",
  "facts": {"feed_id": "gmc_feed_17", "rejected_items": 140, "last_sync": "2026-10-14T03:00:00Z"},
  "actions_taken": [{"tool": "get_sync_errors", "result": "140 x missing_gtin"}],
  "promises_made": [],
  "next_step": "Check whether the new GTINs reached our product table, then re-sync",
  "last_user_message": "Ich habe die GTINs gestern eingetragen, warum ist es immer noch rot?"
}
```

**Safety for side effects.** Sort tools by risk. Read-only tools run freely; reversible, low-risk actions such as a re-sync use idempotency keys and rate limits. For money and hard-to-reverse actions, the model only proposes, code checks policy (amount limit, refund window), the merchant confirms with a real UI button rather than a typed "yes" that the model interprets, and everything is audit-logged. Some actions stay out of reach: to change a payment method, the agent sends a link to the billing page. Product titles, ticket text and platform messages are untrusted input, and an injection hidden there must never trigger a credit. Also avoid the lethal trifecta: an agent that sees private data and untrusted content must not also have a free-form way to send data out, such as emailing any address.

> **带副作用操作的安全性。** 按风险给工具分级。只读工具自由运行；可逆的低风险操作（比如重新同步）使用幂等键和限流。对于涉及金钱和难以撤销的操作：模型只能提议，由代码检查政策（金额上限、退款期限），商家在 UI 里点击一个真正的按钮来确认，而不是由模型去理解一句打字输入的“好的”，并且一切都写入审计日志。有些操作 agent 根本碰不到：要更换支付方式，agent 只发送一个指向计费页面的链接。商品标题、工单文本和平台消息都是不可信输入，藏在其中的注入绝不能触发一笔抵扣。还要避免“致命三要素”（lethal trifecta）：一个能看到私有数据和不可信内容的 agent，绝不能同时拥有随意向外发送数据的途径，比如给任意地址发邮件。

**Consistent tone and policy.** All agents share one versioned tone-and-policy block, which is also a stable, cacheable prefix. Facts that change, such as refund windows or plan prices, come from tools or config, so they change in one place. A final output check (code or a small classifier) catches promises the agent cannot make ("we will refund you"), the wrong language and leaked internal data.

> **一致的语气和政策。** 所有 agent 共享同一段带版本号的语气与政策文本，它同时也是一个稳定的、可缓存的前缀。退款期限、套餐价格这类会变化的事实从工具或配置中获取，这样只需在一处修改。最后的输出检查（代码或小分类器）会拦截 agent 无权做出的承诺（“我们会给你退款”）、错误的语言和泄露的内部数据。

**Escalation.** Escalate on policy (amounts over a limit, suspensions, legal topics), on being stuck (a tool failing repeatedly, the turn budget used up), on strong negative sentiment or repeat contact, and always when the merchant asks for a human. Escalation is a tool that opens a ticket with the handoff packet and tells the merchant what happens next.

> **升级到人工。** 以下情况要升级：触及政策（金额超限、账户被封、法律话题）、卡住了（某个工具反复失败、轮次预算用完）、强烈的负面情绪或重复来访，以及商家要求转人工时（一律照办）。升级本身是一个工具：它带着交接数据包创建工单，并告诉商家接下来会发生什么。

**Evaluate in layers.** Triage: a labelled set and a confusion matrix, graded in code. Specialists: sandbox tasks with mocked platform APIs, graded on the end state and on trajectory rules such as "never issue a credit without confirmation". End to end: simulated merchants (an LLM user simulator, as in tau-bench), graded on resolution, policy and tone over repeated trials with pass^k, because merchants need consistency. Online: resolution without a human, escalation rate (and whether humans agreed), repeat contact within a week, satisfaction, reversed actions and cost per resolved case.

> **分层评估。** 分诊：带标注的数据集和混淆矩阵，用代码打分。专职 agent：在沙箱里用模拟的平台 API 跑任务，按最终状态和轨迹规则打分，比如“未经确认绝不发放抵扣额度”。端到端：用模拟商家（LLM 用户模拟器，就像 tau-bench 那样），在多次重复试验中按问题解决、政策合规和语气打分并看 pass^k，因为商家需要稳定一致的表现。线上：无需人工即解决的比例、升级率（以及人工是否认同）、一周内的重复来访率、满意度、被撤销的操作和每个已解决工单的成本。

**Roll out in stages.** Start in draft mode: the agent drafts, human agents send, and you track their edits. Then go live with read-only answers on one topic (say, catalog sync) for some merchants, then low-risk actions, and money actions with approval last. Add feature flags per category and segment, a kill switch and alerts on escalation spikes, and let a human take over any conversation.

> **分阶段上线。** 先用草稿模式：agent 起草，人工客服发送，你统计他们的改动。然后面向部分商家，在一个话题上（比如目录同步）上线只读回答，接着开放低风险操作，需要审批的涉及金钱的操作放在最后。加上按类别和分群的功能开关、紧急关闭开关和升级量激增告警，并让人工可以接管任何对话。

**Trade-offs.** More agents mean clearer ownership and tighter permissions, but also more handoffs, more latency and more lost context. Start with one agent, good tools and policy enforced in code, and split only where permissions or instructions truly diverge.

> **权衡。** agent 越多，归属越清晰、权限越收紧，但交接越多、延迟越高、丢失的上下文也越多。先从一个 agent、好用的工具和在代码里强制执行的政策开始，只在权限或说明确实分化的地方拆分。

</details>

### Follow-ups · 追问

- How would you stop the billing agent from being talked into a credit by text inside a ticket or a product description?
  > 如何防止计费 agent 被工单或商品描述里的文字“说服”去发放抵扣额度？
- A merchant switches from German to English in the middle of a handoff. What could break, and how do you handle it?
  > 商家在交接过程中从德语切换成了英语。可能出什么问题，你会怎么处理？
- What would you put on the dashboard for the first week after launch?
  > 上线后的第一周，你会在监控面板上放哪些指标？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| triage | 分诊（分类分流） |
| handoff | 交接（转交对话控制权） |
| escalation to a human | 升级到人工 |
| side effect | 副作用（对外部产生影响的操作） |
| The case record, not the transcript, is the source of truth. | 以工单记录为准，而不是以对话记录为准。 |

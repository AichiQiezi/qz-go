# 01. The agent loop: from one LLM call to an agent · Agent 循环：从一次 LLM 调用到 Agent

How a single LLM call becomes an agent: when a fixed workflow beats an autonomous agent, exactly what goes over the wire in a tool-use loop, and how to stop a loop from running forever or quitting too early. Interviewers use these questions to check that you can build and debug the core loop yourself, not just call an API.

> 一次 LLM 调用如何变成一个 Agent：什么时候固定的 workflow 比自主 Agent 更合适，工具调用循环中到底传输了什么，以及如何防止循环无休止地运行或过早结束。面试官用这些问题来检验你能否亲手构建并调试这个核心循环，而不只是会调用 API。

3 questions · generated 2026-10-15T18:23:15+08:00 · [← 目录 / Contents](../README.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Agent Loop | Medium | [The PM wants an agent for ad copy and Merchant Center tickets](#q1) | 产品经理想用 Agent 来写广告文案、处理 Merchant Center 工单 |
| 2 | Agent Loop | Medium | [Find the bugs in a hand-written tool-use loop](#q2) | 找出手写工具调用循环中的 bug |
| 3 | Agent Loop | Hard | [A support agent that loops for 40 turns, or quits after one](#q3) | 客服 Agent 要么循环 40 轮，要么一轮就收工 |

---

<a id="q1"></a>

## 1. The PM wants an agent for ad copy and Merchant Center tickets

**产品经理想用 Agent 来写广告文案、处理 Merchant Center 工单** · `Agent Loop` · `Medium`

### Question · 题目

Your product manager wants to "add an AI agent" for two features:

> 产品经理想为两个功能“加一个 AI Agent”：

- Ad copy: every night, generate Google Shopping and Meta ad titles and descriptions for a merchant's 10,000 products, following each platform's length limits and the merchant's brand tone.
  > 广告文案：每天夜里为某个商家的 10,000 个商品生成 Google Shopping 和 Meta 广告的标题与描述，要符合各平台的长度限制和商家的品牌调性。
- Disapproval tickets: answer merchant tickets such as "Why was my product disapproved by Google Merchant Center?" This means looking at the product data, our last feed sync, and the disapproval reason Google returned, and sometimes fixing the feed and resyncing.
  > 拒审工单：回答商家的工单，例如“为什么我的商品被 Google Merchant Center 拒审了？”。这需要查看商品数据、我们最近一次的 feed 同步记录以及 Google 返回的拒审原因，有时还要修正 feed 并重新同步。

For each feature, would you build a fixed workflow or an autonomous agent? Define the difference, name the common workflow patterns, and explain the criteria you use to decide.

> 对每个功能，你会做成固定的 workflow，还是自主的 Agent？请说明两者的区别，列出常见的 workflow 模式，并解释你做决定的标准。

<details>
<summary><b>Answer · 参考答案</b></summary>

**The core idea.** Both options use an LLM with tools; the difference is who controls the flow. In a **workflow**, your code decides the steps: it calls the model at fixed points, in a fixed order, and checks the results. In an **agent**, the model decides: it runs in a loop, picks the next tool call based on what it has seen so far, and decides when it is done. This framing comes from Anthropic's "Building effective agents" (December 2024), and its main advice is to start with the simplest thing that works. Often that is a single well-prompted call with good examples, and you add steps or autonomy only when measurements show you need them.

> **核心思想。** 两种方案用的都是“LLM + 工具”，区别在于由谁来控制流程。在 **workflow** 中，由你的代码决定步骤：在固定的位置、按固定的顺序调用模型，并检查结果。在 **agent** 中，由模型来决定：它在一个循环里运行，根据目前看到的信息选择下一次调用哪个工具，并自己判断什么时候完成。这种划分来自 Anthropic 的《Building effective agents》（2024 年 12 月），文中最核心的建议是：从最简单可行的方案开始。很多时候，这只是一次提示词写得好、带有好示例的调用；只有当度量结果表明确实需要时，才增加步骤或自主性。

**The five workflow patterns.** Each one is just code calling the model:

> **五种 workflow 模式。** 每一种都只是代码在调用模型：

- Prompt chaining: a fixed sequence of calls, each working on the previous output, with code checks (gates) in between. Example: extract key attributes, write the title, then check its length.
  > Prompt chaining（提示链）：按固定顺序进行的一串调用，每一步处理上一步的输出，中间由代码做检查（gate）。例如：先提取关键属性，再写标题，然后检查长度。
- Routing: classify the input first, then send it to a specialised prompt, model or handler. Example: send easy tickets to a cheaper model and billing questions to a billing flow.
  > Routing（路由）：先对输入分类，再把它交给专门的提示词、模型或处理器。例如：简单工单交给更便宜的模型，账单问题走账单流程。
- Parallelization: run independent calls at the same time, either by splitting the work into parts (sectioning) or by running the same task several times and combining the answers (voting).
  > Parallelization（并行化）：同时运行彼此独立的调用，要么把工作拆成几部分（sectioning），要么把同一任务跑多次再合并答案（voting）。
- Orchestrator-workers: a central model breaks a task into subtasks it chooses at run time, hands them to workers, and merges the results. Unlike parallelization, the subtasks are not known in advance.
  > Orchestrator-workers（编排者-执行者）：由一个中心模型在运行时把任务拆成它自己决定的子任务，交给 worker 执行，再合并结果。与并行化不同，子任务事先并不确定。
- Evaluator-optimizer: one call generates, another checks the output against clear criteria and gives feedback, and you repeat until it passes or you hit a retry limit.
  > Evaluator-optimizer（评估者-优化者）：一次调用负责生成，另一次调用按明确的标准检查输出并给出反馈，如此反复，直到通过或达到重试上限。

**What an agent adds, and what it costs.** An agent is the ReAct pattern (Yao et al., 2022) in production: reason, call a tool, read the result, repeat. It handles problems where you cannot predict the steps. The price is more tokens (Anthropic's multi-agent research write-up reported agents using about 4x the tokens of a chat interaction), higher and less predictable latency, different paths on different runs, and errors that compound across steps. It is also harder to test, because there is no fixed path to assert on.

> **Agent 带来了什么，又要付出什么。** Agent 就是投入生产的 ReAct 模式（Yao et al., 2022）：推理、调用工具、读取结果，然后重复。它能处理那些无法预先确定步骤的问题。代价是更多的 token（Anthropic 关于多 Agent 研究系统的文章提到，Agent 消耗的 token 约为普通聊天的 4 倍）、更高且更难预测的延迟、每次运行的路径都可能不同，以及在多个步骤间不断累积的错误。它也更难测试，因为没有一条固定的路径可供断言。

**Decision criteria.** Ask these questions before you build an agent:

> **决策标准。** 在构建 Agent 之前，先问这些问题：

- Can you write the steps down in advance? If so, a workflow is cheaper, faster and easier to test.
  > 能否事先把步骤写下来？如果可以，workflow 更便宜、更快，也更容易测试。
- At your volume, is the task valuable enough to pay for the cost and latency of a loop?
  > 以你的调用量来看，这个任务的价值是否值得承担循环带来的成本和延迟？
- Can the model reliably do the hardest step? Test that on real data before you design around it.
  > 模型能否可靠地完成其中最难的那一步？在围绕它做设计之前，先用真实数据测一测。
- Can errors be caught and recovered from: validators, tests, human review, rollback? An agent takes unscripted paths, so it needs these guardrails even more than a workflow does.
  > 错误能否被发现并恢复：校验器、测试、人工审核、回滚？Agent 会走出脚本之外的路径，所以它比 workflow 更需要这些护栏。
- Does it have to be interactive, or can it run offline in bulk?
  > 它必须是交互式的，还是可以离线批量运行？

**Ad copy for 10,000 products: a workflow.** The steps are the same for every product: build the prompt from catalog data, choose a template by category (routing), generate, then validate in code (length limits, banned words, required attributes). If validation fails, send the error back for one or two retries (evaluator-optimizer with a cap), and send products that still fail to a review queue. Products are independent, so the job parallelizes trivially, and because it runs overnight you can use a batch API: Claude's Message Batches API costs about 50% of the normal price, and results come back in any order, so you match them by `custom_id`. An agent deciding what to do next for each product would add cost and randomness and buy nothing.

> **为 10,000 个商品写广告文案：用 workflow。** 每个商品的步骤都一样：用商品目录数据组装提示词，按品类选择模板（routing），生成文案，然后用代码校验（长度限制、违禁词、必填属性）。如果校验失败，就把错误信息发回去重试一两次（带上限的 evaluator-optimizer），仍然失败的商品进入人工审核队列。商品之间相互独立，所以这个任务天然可以并行；又因为它在夜间运行，可以使用批处理 API：Claude 的 Message Batches API 价格约为正常价格的 50%，结果按任意顺序返回，所以要用 `custom_id` 来对应。如果让 Agent 为每个商品决定下一步做什么，只会增加成本和随机性，没有任何收益。

**Disapproval tickets: a bounded agent behind a router.** Here the path depends on what you find. The cause might be a missing GTIN, a price that differs between the feed and the landing page, a policy issue, or a broken image, and each needs different lookups. Put a router in front: common reasons with known fixes go to deterministic playbooks, and only unclear cases reach the agent. Give the agent read-only tools by default (get the product, get the sync log, get the Merchant Center issue). Put any write action, such as editing the feed or triggering a resync, behind merchant or human approval. Cap turns, tokens and time, and escalate to a human when a limit is reached.

> **拒审工单：放在路由之后、有边界的 Agent。** 这里的排查路径取决于你查到了什么。原因可能是缺少 GTIN、feed 中的价格与落地页不一致、违反政策，或者图片失效，每一种都需要不同的查询。在前面放一个路由：有已知修复方式的常见原因走确定性的处理剧本（playbook），只有不明确的情况才交给 Agent。默认只给 Agent 只读工具（查询商品、查询同步日志、查询 Merchant Center 的问题详情）。任何写操作，例如修改 feed 或触发重新同步，都要经过商家或人工的批准。为轮数、token 和时间设置上限，达到上限时转交人工处理。

**Trade-offs to mention.** Most production systems are hybrids: a workflow skeleton with one agentic step, or an agent whose tools are themselves small workflows. Decide with data, not opinions: replay 50 to 100 historical tickets through each design and compare resolution rate, cost per resolved ticket and latency. Autonomy has to earn its place by beating the simpler design.

> **要提到的权衡。** 大多数生产系统都是混合形态：一个 workflow 骨架中包含一个 agentic 步骤，或者一个 Agent 的工具本身就是小型 workflow。用数据而不是观点来做决定：把 50 到 100 个历史工单分别回放到每种设计中，比较解决率、每个已解决工单的成本以及延迟。自主性必须通过胜过更简单的设计，才能证明自己值得采用。

</details>

### Follow-ups · 追问

- How would you build an evaluation set for the disapproval-ticket agent, and what would you measure?
  > 你会如何为拒审工单 Agent 构建评估集？要度量哪些指标？
- After two retries, 8% of products still fail validation in the batch job. What do you do with them?
  > 在批处理任务中，重试两次之后仍有 8% 的商品校验不通过。你会怎么处理这些商品？
- Where exactly would you put human approval in the ticket flow, and how does the agent wait for it?
  > 在工单流程中，你会把人工审批放在哪一步？Agent 又如何等待审批结果？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| workflow vs. agent | workflow 与 Agent：由代码控制流程，还是由模型控制流程 |
| prompt chaining | 提示链：按固定顺序串联多次 LLM 调用 |
| evaluator-optimizer loop | 评估者-优化者循环：生成、评估、反馈，直到通过 |
| Start with the simplest thing that works. | 从最简单可行的方案开始。 |
| Autonomy has to earn its place. | 自主性必须证明自己值得采用。 |

---

<a id="q2"></a>

## 2. Find the bugs in a hand-written tool-use loop

**找出手写工具调用循环中的 bug** · `Agent Loop` · `Medium`

### Question · 题目

A junior engineer wrote the loop below for a catalog support agent on the Claude Messages API (simplified Go; `client`, `Block` and the helpers are thin wrappers over the HTTP API). It is about to ship.

> 一位初级工程师基于 Claude Messages API，为商品目录客服 Agent 写了下面这个循环（简化的 Go 代码；`client`、`Block` 和各个辅助函数都是对 HTTP API 的薄封装）。它马上就要上线了。

```go
func runAgent(ctx context.Context, userMsg string) (string, error) {
    msgs := []Message{{Role: "user", Content: []Block{Text(userMsg)}}}
    for {
        resp, err := client.CreateMessage(ctx, Request{
            Model: model, MaxTokens: 1024, Tools: tools, Messages: msgs,
        })
        if err != nil {
            return "", err
        }
        calls := toolUses(resp.Content)
        if len(calls) == 0 {
            return resp.Text(), nil // no tools requested: done
        }
        // keep the history small: store only what the model said
        msgs = append(msgs, Message{Role: "assistant", Content: []Block{Text(resp.Text())}})
        for _, c := range calls {
            out, err := runTool(ctx, c.Name, c.Input)
            if err != nil {
                out = "" // tool failed, return nothing
            }
            msgs = append(msgs, Message{Role: "user", Content: []Block{ToolResult(c.ID, out)}})
        }
    }
}
```

First walk through exactly what goes over the wire in one agent turn with tools. Then find the bugs, say what each one would do in production, and fix them.

> 请先讲清楚：带工具的一轮 Agent 交互中，网络上到底传输了什么。然后找出其中的 bug，说明每个 bug 在生产环境中会造成什么后果，并修复它们。

<details>
<summary><b>Answer · 参考答案</b></summary>

**What goes over the wire.** The API is stateless: every request carries the whole conversation, and the client runs the tools. One turn with tools looks like this:

> **网络上传输的内容。** API 是无状态的：每个请求都携带完整的对话，工具由客户端执行。带工具的一轮交互如下：

- The request contains `tools` (name, description, `input_schema`) and `messages`.
  > 请求中包含 `tools`（name、description、`input_schema`）和 `messages`。
- When the model wants tools, the response `content` contains `tool_use` blocks, each with an `id`, a `name` and an `input`, and `stop_reason` is `"tool_use"`. One response can contain several `tool_use` blocks, because parallel tool use is on by default.
  > 当模型想调用工具时，响应的 `content` 中会包含 `tool_use` 块，每个块都有 `id`、`name` 和 `input`，并且 `stop_reason` 为 `"tool_use"`。一个响应可以包含多个 `tool_use` 块，因为并行工具调用默认是开启的。
- The client appends the assistant message exactly as received, runs the tools, and sends one `user` message with a `tool_result` block per call. Each result carries the matching `tool_use_id`, its `content`, and `is_error: true` if the tool failed.
  > 客户端把 assistant 消息原样追加到历史中，执行工具，然后发送一条 `user` 消息，其中每次调用对应一个 `tool_result` 块。每个结果都带有对应的 `tool_use_id` 和 `content`；如果工具执行失败，还要加上 `is_error: true`。
- The loop repeats until `stop_reason` is `"end_turn"`.
  > 循环一直重复，直到 `stop_reason` 为 `"end_turn"`。

Concretely, the history after one parallel turn looks like this:

> 具体来说，一轮并行调用之后，历史记录是这样的：

```json
[
  {"role": "assistant", "content": [
    {"type": "text", "text": "Let me check the product and its last Google sync."},
    {"type": "tool_use", "id": "toolu_01A", "name": "get_product",
     "input": {"product_id": "p_123"}},
    {"type": "tool_use", "id": "toolu_01B", "name": "get_sync_log",
     "input": {"product_id": "p_123", "channel": "google"}}
  ]},
  {"role": "user", "content": [
    {"type": "tool_result", "tool_use_id": "toolu_01A",
     "content": "Linen shirt, price 39.90 USD, gtin: missing"},
    {"type": "tool_result", "tool_use_id": "toolu_01B", "is_error": true,
     "content": "sync-log service timed out after 10s; safe to retry"}
  ]}
]
```

**Bug 1: the tool_use blocks are dropped from history.** To "keep the history small", the loop stores only the text. The next request then contains a `tool_result` whose `tool_use_id` matches no `tool_use` block in the previous assistant message, and the API rejects it with a 400. The model would also lose the record of which calls it made. Fix: append `resp.Content` unchanged. That also keeps any thinking or compaction blocks that the API expects to see again.

> **Bug 1：从历史中丢掉了 tool_use 块。** 为了“让历史保持精简”，这个循环只保存了文本。于是下一个请求中的 `tool_result` 所引用的 `tool_use_id`，在上一条 assistant 消息里找不到任何对应的 `tool_use` 块，API 会以 400 错误拒绝这个请求。模型也会失去它自己发起过哪些调用的记录。修复：原样追加 `resp.Content`。这样也能保留 API 期望再次看到的 thinking 块或 compaction 块。

**Bug 2: parallel results go into separate messages.** When the model asks for two tools, this loop sends two `user` messages in a row, but the API expects all results for one assistant turn in the single user message that follows it. Consecutive same-role messages can be rejected outright, and Anthropic's docs warn that split results teach the model to stop making parallel calls, since the history shows one call per turn. Each independent lookup then costs an extra round trip. Fix: send all results in one message, with the `tool_result` blocks first and any extra text after them.

> **Bug 2：并行调用的结果被放进了不同的消息。** 当模型请求两个工具时，这个循环会连续发送两条 `user` 消息，但 API 期望同一个 assistant 轮次的所有结果都放在紧随其后的那一条 user 消息里。连续出现的同角色消息可能会被直接拒绝；而且 Anthropic 的文档提醒，拆开发送的结果会让模型学会不再发起并行调用，因为历史记录显示的是每轮只调用一次。之后每一次相互独立的查询都要多一次往返。修复：把所有结果放在一条消息里发送，`tool_result` 块放在最前面，其他文本放在它们后面。

**Bug 3: the loop ignores the stop reason.** It treats "no tool calls" as "done" and "has tool calls" as "run them", and both are wrong when the output was cut off. With `stop_reason: "max_tokens"`, an incomplete text answer is returned as final, or the last `tool_use` block has a truncated input that the loop executes anyway, which is dangerous for a write tool. `MaxTokens: 1024` makes this common, especially when thinking counts toward the limit. `refusal` and `pause_turn` (a paused server-side tool loop; send again to resume) need handling too. Fix: switch on `stop_reason`, and set `max_tokens` generously: it is a safety cap the model does not see, not a length control.

> **Bug 3：循环忽略了停止原因。** 它把“没有工具调用”当成“完成”，把“有工具调用”当成“去执行”，而当输出被截断时，这两种判断都是错的。当 `stop_reason: "max_tokens"` 时，不完整的文本答案会被当作最终答案返回，或者最后一个 `tool_use` 块的 input 已被截断，循环却照样执行，对写操作工具来说这很危险。`MaxTokens: 1024` 让这种情况很常见，尤其是当 thinking 也计入这个上限时。`refusal` 和 `pause_turn`（服务端工具循环暂停了；再次发送即可继续）同样需要处理。修复：根据 `stop_reason` 分支处理，并把 `max_tokens` 设得宽裕一些：它是一个模型看不到的安全上限，而不是控制长度的手段。

**Bug 4: tool errors are swallowed.** A timeout becomes an empty string with no error flag, so the model sees a successful call that found nothing. It then tells the merchant the product does not exist, or retries the same call over and over. Fix: return `is_error: true` with a message the model can act on, saying what failed and whether a retry is safe. Keep "not found" as a normal result: it is a real answer, and a failure is not.

> **Bug 4：工具错误被吞掉了。** 一次超时变成了一个没有错误标记的空字符串，于是模型看到的是一次成功但什么都没找到的调用。接着它会告诉商家这个商品不存在，或者一遍又一遍地重试同一个调用。修复：返回 `is_error: true`，并附上模型能据此采取行动的信息：什么失败了、重试是否安全。“未找到”要作为正常结果返回：它是一个真实的答案，而失败不是。

**The fixed loop.**

> **修复后的循环。**

```go
func runAgent(ctx context.Context, userMsg string) (string, error) {
    msgs := []Message{{Role: "user", Content: []Block{Text(userMsg)}}}
    for turn := 0; turn < maxTurns; turn++ {
        resp, err := client.CreateMessage(ctx, Request{
            Model: model, MaxTokens: maxOutputTokens, Tools: tools, Messages: msgs,
        })
        if err != nil {
            return "", err
        }
        // Bug 1: keep the full content, tool_use blocks included.
        msgs = append(msgs, Message{Role: "assistant", Content: resp.Content})

        switch resp.StopReason {
        case "end_turn":
            return resp.Text(), nil
        case "tool_use":
            // run the tools below
        case "pause_turn":
            continue // server-side tool loop paused: send again to resume
        case "max_tokens":
            // Bug 3: output was cut off; a tool input may be incomplete.
            return "", ErrTruncated
        default: // "refusal" or anything new
            return "", fmt.Errorf("agent stopped: %s", resp.StopReason)
        }

        // Bug 2: all results go back in ONE user message.
        var results []Block
        for _, c := range toolUses(resp.Content) {
            out, err := runTool(ctx, c.Name, c.Input)
            if err != nil {
                // Bug 4: report the failure so the model can react to it.
                results = append(results, ToolError(c.ID, describe(err)))
                continue
            }
            results = append(results, ToolResult(c.ID, out))
        }
        msgs = append(msgs, Message{Role: "user", Content: results})
    }
    return "", ErrTooManyTurns
}
```

**Also worth saying.** Add a turn limit and a deadline; the original `for {}` has neither. Run independent read-only tools concurrently (for example with `errgroup`) and side-effecting ones one at a time, and validate inputs before running anything. The official SDKs also ship a beta tool runner that drives this loop, with hooks for approval, error handling and retries; hand-write the loop only when you need control it does not give you.

> **还值得一提。** 加上轮数上限和截止时间，原来的 `for {}` 两者都没有。相互独立的只读工具可以并发执行（例如用 `errgroup`），有副作用的工具则逐个执行；执行任何工具之前都要先校验输入。官方 SDK 还提供了一个 beta 版的 tool runner 来驱动这个循环，并提供审批、错误处理和重试的钩子；只有当你需要它无法提供的控制能力时，才手写这个循环。

</details>

### Follow-ups · 追问

- Which tool calls in this agent would you run concurrently, and how would you keep side-effecting calls safe?
  > 这个 Agent 中哪些工具调用你会并发执行？对有副作用的调用，你会如何保证安全？
- The SDK ships a tool runner. When would you still write the loop yourself?
  > SDK 自带了 tool runner。什么情况下你仍然会自己手写这个循环？
- How would you handle a `max_tokens` stop without losing the work the model has already done?
  > 遇到 `max_tokens` 停止时，你会如何处理，才能不丢掉模型已经完成的工作？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| stop_reason | 停止原因：模型为什么结束这一次响应 |
| tool_use_id | 把 tool_result 与对应的 tool_use 关联起来的 ID |
| parallel tool calls | 并行工具调用：一次响应中请求多个工具 |
| The API is stateless; the client owns the history. | API 是无状态的，对话历史由客户端负责维护。 |
| to swallow an error | 吞掉错误：捕获之后既不上报也不处理 |

---

<a id="q3"></a>

## 3. A support agent that loops for 40 turns, or quits after one

**客服 Agent 要么循环 40 轮，要么一轮就收工** · `Agent Loop` · `Hard`

### Question · 题目

Your merchant-support agent has five tools: `search_help_center`, `get_order`, `get_product`, `get_sync_status` and `create_ticket`. Traces show two failure modes:

> 你们的商家客服 Agent 有五个工具：`search_help_center`、`get_order`、`get_product`、`get_sync_status` 和 `create_ticket`。从 trace 中可以看到两种失败模式：

- Some sessions run for 40 turns, calling `search_help_center` with small variations of the same query, until a gateway timeout kills them.
  > 有些会话会跑 40 轮，用同一个查询的细微变体反复调用 `search_help_center`，直到被网关超时强行终止。
- Other sessions end after one turn with a half answer such as "It looks like there may be a sync issue, let me check that for you." and no tool call.
  > 另一些会话只跑一轮就结束了，只给出半个答案，例如 "It looks like there may be a sync issue, let me check that for you."，却没有任何工具调用。

Explain why agents loop and why they stop too early. Then describe how you would bound and steer this agent before the next release.

> 请解释 Agent 为什么会陷入循环，又为什么会过早停止。然后说明在下一次发布之前，你会如何为这个 Agent 设定边界并加以引导。

<details>
<summary><b>Answer · 参考答案</b></summary>

**The core idea.** An agent loop ends in one of two ways: the model returns a turn with no tool calls (`stop_reason` is `end_turn`), or your harness stops it. So "when to stop" is decided jointly: by the model's judgment, which depends on the goal, the prompt and what the tool results tell it, and by your hard limits. Looping and quitting early are two sides of the same gap: the agent has no clear definition of done and no clear signal of progress.

> **核心思想。** Agent 循环只会以两种方式结束：模型返回一个没有工具调用的轮次（`stop_reason` 为 `end_turn`），或者你的 harness 主动终止它。所以“什么时候停”是由两方共同决定的：一方是模型的判断，它取决于目标、提示词以及工具结果告诉它的信息；另一方是你设定的硬性上限。陷入循环和过早结束是同一个缺口的两面：Agent 既没有清晰的“完成”定义，也没有清晰的进展信号。

**Why agents loop.**

> **Agent 为什么会陷入循环。**

- No progress signal. The search returns the same five articles every time, and nothing tells the model that rewording will not help, so one more query always looks like a reasonable next step.
  > 没有进展信号。搜索每次都返回同样的五篇文章，却没有任何信息告诉模型换个说法也没用，所以“再查一次”看起来总是合理的下一步。
- An open-ended goal. "Help the merchant with their sync problem" has no finish line, so the model keeps gathering information.
  > 目标没有终点。"Help the merchant with their sync problem" 没有明确的终点线，于是模型不停地收集信息。
- Errors that look like empty results. If a timeout or a 403 comes back as an empty list, the model reads it as "nothing found" and tries again with a new query.
  > 错误看起来像空结果。如果超时或 403 被返回成一个空列表，模型会把它理解为“什么都没找到”，然后换个查询再试。
- No exit. If the prompt says "always resolve the issue" and there is no way to escalate or to ask the merchant, another tool call is the only move left.
  > 没有出口。如果提示词要求 "always resolve the issue"，又没有升级处理或询问商家的途径，那么再调用一次工具就是唯一能走的一步。
- Growing context. After many turns, long and repetitive tool output pushes the important parts out of focus, and the model loses track of what it has already tried.
  > 上下文不断膨胀。很多轮之后，冗长而重复的工具输出会把重要信息挤出注意力范围，模型也会忘记自己已经尝试过什么。

**Why agents stop too early.**

> **Agent 为什么会过早停止。**

- Narrating instead of acting. The model writes "let me check" and ends its turn without a tool call. Your loop sees no tool call and treats the text as the final answer.
  > 只说不做。模型写下 "let me check"，然后在没有工具调用的情况下结束了这一轮。你的循环没看到工具调用，就把这段文字当成了最终答案。
- Truncation treated as success. A `max_tokens` stop is returned as if the work were finished.
  > 把截断当成成功。`max_tokens` 导致的停止被当作工作已完成而返回。
- Tool descriptions that do not say when to call them, so the model answers from general knowledge instead of looking things up.
  > 工具描述没有说明什么时候该调用，于是模型凭通用知识作答，而不去实际查询。
- Low effort, or a prompt that rewards brevity, so the model does the minimum.
  > effort 设置过低，或者提示词鼓励简短，于是模型只做最少的工作。

**Bound it: hard limits in the harness.** These are the agent version of a retry policy with a maximum number of attempts:

> **设定边界：在 harness 中设置硬性上限。** 这些相当于 Agent 版本的、带最大尝试次数的重试策略：

- A turn limit per task, set from your traces (for example, comfortably above the 95th percentile of successful sessions).
  > 为每个任务设置轮数上限，并根据 trace 数据来定（例如留出余量，明显高于成功会话的第 95 百分位）。
- A token or cost budget per task, summed from each response's `usage`, plus a wall-clock deadline through `context.Context`.
  > 为每个任务设置 token 或成本预算（从每个响应的 `usage` 累加），再通过 `context.Context` 设置墙钟截止时间。
- Repeated-call detection: key each call by tool name plus normalized input, and refuse exact repeats with a result that explains why (code below).
  > 重复调用检测：用工具名加规范化后的输入作为每次调用的键，拒绝完全相同的重复调用，并在结果中说明原因（代码见下）。
- A graceful exit: when a limit is hit, make one last request with `tool_choice` set to `none`, asking the model to summarize what it found and what is missing, then create a ticket for a human.
  > 优雅退出：触发上限时，最后再发一次 `tool_choice` 为 `none` 的请求，让模型总结已经发现了什么、还缺什么，然后为人工创建工单。

```go
// seen := map[string]int{} is created once per task, before the loop.
key := c.Name + ":" + canonicalJSON(c.Input)
seen[key]++
if seen[key] > 2 {
    results = append(results, ToolError(c.ID,
        "This exact call has already run twice. Do not repeat it: change the "+
            "query substantially, use a different tool, or call create_ticket "+
            "with what you have found so far."))
    continue
}
```

**Steer it: help the model decide.**

> **加以引导：帮助模型做决定。**

- Write explicit completion criteria into the system prompt: you are done when you have found the cause with evidence from a tool and told the merchant the fix, or when you have created a ticket; never end a turn by promising to check something.
  > 在系统提示词中写出明确的完成标准：当你已经用工具返回的证据找到原因并告诉商家修复方法，或者已经创建了工单时，任务才算完成；绝不要以“我去查一下”这样的承诺结束一轮。
- Make tool results informative: say what was searched and how many results came back, separate "no matches" from "the search failed" (with `is_error`), and add hints such as "try the error code returned by get_sync_status".
  > 让工具结果包含有用的信息：说明搜索了什么、返回了多少结果，用 `is_error` 区分“没有匹配”和“搜索失败”，并加上提示，例如 "try the error code returned by get_sync_status"。
- Give it exits: tools to escalate or to ask the merchant a question, with descriptions that say when to use them.
  > 给它出口：提供升级到人工或向商家提问的工具，并在描述中说明何时使用。
- Show progress: after each round of tool results, append a short reminder such as "turn 9 of 15; already searched for X and Y". Newer Claude models accept an appended `system`-role message for this, which keeps the prompt cache intact because the system prompt itself is not edited.
  > 展示进度：每一轮工具结果之后，追加一条简短的提醒，例如 "turn 9 of 15; already searched for X and Y"。较新的 Claude 模型支持为此追加一条 `system` 角色的消息；由于系统提示词本身没有被修改，prompt cache 不会失效。
- Tune effort per route: lower effort means fewer, more consolidated tool calls, which suits simple FAQ questions, while investigations that stop too early may need more.
  > 按路由调整 effort：较低的 effort 意味着更少、更集中的工具调用，适合简单的 FAQ 问题；而过早停止的排查类任务可能需要更高的 effort。
- Catch the promise pattern: if a final answer says "let me check" without a tool call, send one nudge asking the model to call the tool or give its final answer, and allow that only once.
  > 拦截“承诺模式”：如果最终答案里说了 "let me check" 却没有工具调用，就发送一次提醒，要求模型调用工具或给出最终答案，并且只允许提醒一次。

**Measure it.** Log per session: turns, tokens, tool-call counts, repeated-call rate, the final `stop_reason`, budget hits, escalations and the outcome. Turn the bad traces into an evaluation set. Because these failures are intermittent, run each case several times and track pass^k (from tau-bench), the probability that all k runs succeed. It exposes flakiness that a single pass rate hides.

> **度量它。** 为每个会话记录：轮数、token、工具调用次数、重复调用率、最终的 `stop_reason`、触发预算上限的次数、转人工的次数以及结果。把出问题的 trace 变成评估集。由于这些故障是间歇性的，每个用例都要运行多次，并跟踪 pass^k（来自 tau-bench），即 k 次运行全部成功的概率。它能暴露单次通过率所掩盖的不稳定性。

**Trade-offs.** Tight limits mean more escalations and some incomplete answers; loose limits mean higher cost and slow failures. Set them from data and revisit them after model or prompt changes. Treat escalation to a human as a designed outcome, like a dead-letter queue, not as a failure: an agent that hands over a clear summary after 15 turns is better than one that burns 40 turns and times out.

> **权衡。** 上限太紧，意味着更多的转人工和一些不完整的答案；上限太松，意味着更高的成本和更慢的失败。根据数据来设定上限，并在模型或提示词变更后重新评估。把转交人工当作一种设计好的结果，就像死信队列一样，而不是当作失败：一个在 15 轮后交出清晰总结的 Agent，好过一个烧掉 40 轮然后超时的 Agent。

</details>

### Follow-ups · 追问

- Without hard-coding rules for each tool, how would you tell a loop that is making progress from one that is stuck?
  > 在不为每个工具硬编码规则的前提下，你如何区分一个正在取得进展的循环和一个已经卡住的循环？
- Your turn limit is hit in 3% of sessions. How do you decide whether to raise it?
  > 有 3% 的会话触发了轮数上限。你如何决定是否要提高这个上限？
- How would you check that a prompt change fixed early stopping without making the looping worse?
  > 你如何验证一次提示词修改解决了过早停止的问题，同时没有让循环问题变得更严重？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| stop criteria | 停止条件 |
| iteration budget | 迭代预算：每个任务允许的最大轮数 |
| graceful exit | 优雅退出：达到上限时先总结、再交接，而不是直接中断 |
| escalate to a human | 升级给人工处理 |
| It never had a definition of done. | 它从来就没有一个“完成”的定义。 |

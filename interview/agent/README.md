# Agent 面试速成 · Agent Interview Crash Course

面向后端工程师的 AI Agent 面试题，中英双语，按模块从 Agent 循环一路讲到生产环境；之后每天追加一组新题。
每题包含题目、参考答案（默认折叠）、追问和词汇表，中文逐段放在英文下面。

共 13 个模块、49 道题。

- 数据：`interview/agent/data/<模块 id>.json`（与 Agent 面试速成页面的 `modules` 集合结构一致）
- 重新生成 Markdown：在仓库根目录执行 `go run ./interview/render -dir interview/agent -kind modules`

## [01. The agent loop: from one LLM call to an agent · Agent 循环：从一次 LLM 调用到 Agent](modules/01-agent-loop.md)

How a single LLM call becomes an agent: when a fixed workflow beats an autonomous agent, exactly what goes over the wire in a tool-use loop, and how to stop a loop from running forever or quitting too early. Interviewers use these questions to check that you can build and debug the core loop yourself, not just call an API.

> 一次 LLM 调用如何变成一个 Agent：什么时候固定的 workflow 比自主 Agent 更合适，工具调用循环中到底传输了什么，以及如何防止循环无休止地运行或过早结束。面试官用这些问题来检验你能否亲手构建并调试这个核心循环，而不只是会调用 API。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Agent Loop | Medium | [The PM wants an agent for ad copy and Merchant Center tickets](modules/01-agent-loop.md#q1) | 产品经理想用 Agent 来写广告文案、处理 Merchant Center 工单 |
| 2 | Agent Loop | Medium | [Find the bugs in a hand-written tool-use loop](modules/01-agent-loop.md#q2) | 找出手写工具调用循环中的 bug |
| 3 | Agent Loop | Hard | [A support agent that loops for 40 turns, or quits after one](modules/01-agent-loop.md#q3) | 客服 Agent 要么循环 40 轮，要么一轮就收工 |

## [02. Designing tools and MCP servers · 设计工具与 MCP 服务](modules/02-tools-mcp.md)

How to design tools that a model calls correctly, how to expose company systems to AI assistants through the Model Context Protocol (MCP), and what to do when an assistant has too many tools. Interviewers focus on this because many production agent failures trace back to tool design and tool sprawl rather than to the model itself.

> 如何设计能被模型正确调用的工具，如何通过 Model Context Protocol（MCP）把公司的系统开放给 AI 助手，以及当助手接入的工具过多时该怎么办。面试官关注这些，是因为生产环境中很多 Agent 故障的根源在于工具设计和工具泛滥，而不在模型本身。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Tools & MCP | Medium | [The model keeps calling the wrong catalog tool](modules/02-tools-mcp.md#q1) | 模型总是调用错误的商品工具 |
| 2 | Tools & MCP | Hard | [Expose catalog and ad sync to every internal assistant over MCP](modules/02-tools-mcp.md#q2) | 通过 MCP 把商品目录和广告同步能力开放给所有内部助手 |
| 3 | Tools & MCP | Senior | [180 tools from 12 MCP servers made the assistant worse](modules/02-tools-mcp.md#q3) | 接入 12 个 MCP 服务、180 个工具之后，助手反而变差了 |

## [03. Context engineering and memory · 上下文工程与记忆](modules/03-context-memory.md)

What fills an agent's context window, why more tokens can make answers worse and more expensive, and the tools long-running agents rely on: trimming, compaction, external memory, subagents and prompt caching. Interviewers ask because context management is what separates a demo agent from one that runs for hours at a predictable cost.

> Agent 的上下文窗口里装了什么、为什么 token 越多回答可能越差也越贵，以及长时间运行的 Agent 所依赖的手段：裁剪、压缩（compaction）、外部记忆、子 Agent 和提示词缓存。面试官爱问这个，因为上下文管理正是演示级 Agent 和能以可预期的成本连续运行几个小时的 Agent 之间的分水岭。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Context & Memory | Medium | [The window is 1M tokens, so why not send everything?](modules/03-context-memory.md#q1) | 窗口有 100 万 token，为什么不把所有东西都塞进去？ |
| 2 | Context & Memory | Hard | [An overnight audit agent forgets its own decisions](modules/03-context-memory.md#q2) | 通宵运行的审核 Agent 忘了自己做过的决定 |
| 3 | Context & Memory | Hard | [Prompt caching is on, but the cache never hits](modules/03-context-memory.md#q3) | 开了提示词缓存，却从来没命中过 |

## [04. Retrieval and knowledge for agents · Agent 的检索与知识](modules/04-rag.md)

How agents ground answers in your documents and data: the retrieval pipeline, agentic multi-step retrieval and how to evaluate it, and the multi-tenant permission, freshness and deletion problems behind real incidents. Interviewers ask because most production agents depend on retrieval, and many of their wrong answers are really retrieval failures.

> Agent 如何基于你的文档和数据作答：检索流水线、Agent 式多步检索及其评估方法，以及引发真实事故的多租户权限、数据新鲜度和删除问题。面试官爱问这些，因为大多数生产环境中的 Agent 都依赖检索，而它们的很多错误回答其实是检索失败。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | RAG | Medium | [Answering merchants from help-center and ad-policy docs](modules/04-rag.md#q1) | 基于帮助中心和广告政策文档回答商家问题 |
| 2 | RAG | Hard | [Multi-hop questions break single-shot RAG](modules/04-rag.md#q2) | 多跳问题难倒了单次检索的 RAG |
| 3 | RAG | Hard | [Merchant A's answer quoted merchant B's private notes](modules/04-rag.md#q3) | 商家 A 的回答引用了商家 B 的私有笔记 |

## [05. Planning and multi-agent systems · 规划与多 Agent 系统](modules/05-multi-agent.md)

How agents plan and track multi-step work, when to split a job across an orchestrator and subagents, and how to design a multi-agent support system with handoffs, escalation and safe side effects. Interviewers ask because long agent tasks usually fail through forgotten requirements, lost context and unclear ownership, not through one bad model call.

> Agent 如何规划和跟踪多步骤工作，什么时候把任务拆给编排者和子 agent，以及如何设计一个带交接、人工升级和安全副作用控制的多 Agent 客服系统。面试官爱问这些，是因为长任务中的 agent 失败，通常源于要求被遗忘、上下文丢失和职责不清，而不是某一次模型调用出错。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Multi-agent | Medium | [Keeping a multi-step campaign launch on track with an explicit plan](modules/05-multi-agent.md#q1) | 用显式计划让多步骤的活动上线不跑偏 |
| 2 | Multi-agent | Hard | [Fanning out a 200-merchant ROAS investigation to subagents](modules/05-multi-agent.md#q2) | 把 200 个商家的 ROAS 下滑调查扇出给子 agent |
| 3 | Multi-agent | Senior | [Designing a merchant support system with specialists, handoffs and human escalation](modules/05-multi-agent.md#q3) | 设计商家客服系统：专职 agent、交接与人工升级 |

## [06. Evaluating and observing agents · Agent 的评测与可观测性](modules/06-evals.md)

How to trace agent runs, build eval sets and graders, measure reliability, and ship prompt and model changes without surprising production. Interviewers ask because an agent you cannot measure is an agent you cannot improve or ship safely.

> 如何追踪 agent 的运行、构建评测集和评分器、衡量可靠性，以及在不给线上带来意外的前提下发布提示词和模型的改动。面试官爱问这些，是因为一个无法衡量的 agent，既无法改进，也无法安全上线。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Evals & Observability | Medium | [Why some agent runs take three minutes and cost a lot](modules/06-evals.md#q1) | 为什么有些 agent 运行要三分钟、还特别贵 |
| 2 | Evals & Observability | Hard | [Proving a catalog-fixing agent is ready to ship](modules/06-evals.md#q2) | 证明一个修复商品目录的 agent 已经可以上线 |
| 3 | Evals & Observability | Hard | [The eval score went up, and so did production complaints](modules/06-evals.md#q3) | 评测分数涨了，线上投诉也涨了 |

## [07. Safety and security for agents · Agent 的安全](modules/07-safety.md)

How to keep an agent that runs code, reads untrusted content and takes real actions from harming merchants, customers or your company: sandboxing, prompt injection, and guardrails for actions that move money. Interviewers ask because any agent that touches production data needs convincing answers here before it can launch.

> 如何防止一个会运行代码、会读取不可信内容、会执行真实操作的 Agent 伤害商家、顾客或公司：沙箱隔离、提示注入，以及针对“动钱”操作的护栏。面试官会问这些，是因为任何接触生产数据的 Agent，在上线之前都必须在这些问题上给出令人信服的答案。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Safety & Security | Medium | [A data-analysis agent that writes and runs its own scripts](modules/07-safety.md#q1) | 一个自己写脚本、自己跑脚本的数据分析 Agent |
| 2 | Safety & Security | Hard | [Hidden text in a product description hijacks a merchant-ops agent](modules/07-safety.md#q2) | 商品描述里的隐藏文字劫持了商家运营 Agent |
| 3 | Safety & Security | Senior | [Guardrails for an agent that can refund orders and change ad budgets](modules/07-safety.md#q3) | 为一个能退款、能改广告预算的 Agent 设计护栏 |

## [08. Running agents in production · 生产环境中的 Agent](modules/08-production.md)

What changes when an agent leaves the demo: validating structured output, controlling cost and latency, and staying reliable under rate limits, overload and restarts. Interviewers ask because most production agent failures are ordinary distributed-systems problems with a model in the middle.

> Agent 离开演示环境之后会遇到什么：校验结构化输出、控制成本和延迟，以及在限流、过载和重启之下保持可靠。面试官会问这些，是因为生产环境中大多数 Agent 故障，其实都是中间夹着一个模型的普通分布式系统问题。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Production | Medium | [Supplier data to catalog fields: malformed, truncated, or confidently wrong](modules/08-production.md#q1) | 供应商数据转商品目录字段：格式错误、输出截断，或者错得理直气壮 |
| 2 | Production | Hard | [A support agent that costs too much and takes 90 seconds to answer](modules/08-production.md#q2) | 一个太贵、还要 90 秒才能回复的客服 Agent |
| 3 | Production | Hard | [Peak traffic: 429s, overloaded errors, double refunds and restarted tasks](modules/08-production.md#q3) | 流量高峰：429、过载错误、重复退款和从头再来的任务 |

## [Daily set · Oct 16 · 每日新题 · 10 月 16 日](modules/daily-2026-10-16.md)

Five new agent questions: streaming a tool-using chat, long-term memory for merchant preferences, constrained product search, agents that overwrite each other's catalog edits, and choosing between your own loop, an SDK runner, an agent framework and a managed platform.

> 五道新的 Agent 题：带工具调用的聊天如何流式化、商家偏好的长期记忆、带约束条件的商品检索、互相覆盖商品修改的多个 Agent，以及在自研循环、SDK runner、Agent 框架和托管平台之间如何选择。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Production | Medium | [The support chat sits silent for 40 seconds, then dumps the whole answer](modules/daily-2026-10-16.md#q1) | 客服聊天框沉默 40 秒，然后一下子吐出整段回答 |
| 2 | Context & Memory | Medium | [Merchants are tired of repeating their preferences every session](modules/daily-2026-10-16.md#q2) | 商家厌倦了每次会话都重复自己的偏好 |
| 3 | RAG | Hard | ["Dresses under $50, not red" returns red dresses at $79](modules/daily-2026-10-16.md#q3) | “50 美元以下、不要红色的连衣裙”却返回了 79 美元的红裙子 |
| 4 | Multi-agent | Hard | [Two agents keep overwriting each other's catalog edits](modules/daily-2026-10-16.md#q4) | 两个 Agent 不断覆盖彼此对商品目录的修改 |
| 5 | Agent Loop | Senior | [Your own loop, an SDK runner, an agent framework or a managed platform for a nightly feed audit?](modules/daily-2026-10-16.md#q5) | 夜间 feed 审计 Agent：自己写循环、用 SDK runner、用 Agent 框架，还是用托管平台？ |

## [Daily set · Oct 17 · 每日新题 · 10 月 17 日](modules/daily-2026-10-17.md)

Five new agent questions: designing tool errors the model can recover from, keeping merchants' OAuth tokens away from the model, answering reliably from long contracts, finding a success-rate drop with no deploy, and pausing an agent for hours of human approval.

> 五道新的 Agent 题：设计让模型能够恢复的工具错误、让商家的 OAuth token 远离模型、基于长篇合同可靠地回答问题、排查没有发布却下降的成功率，以及让 Agent 暂停数小时等待人工审批。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Tools & MCP | Medium | [The agent gives up after one 429, or retries the same broken call forever](modules/daily-2026-10-17.md#q1) | Agent 遇到一次 429 就放弃，或者无休止地重试同一个错误调用 |
| 2 | Safety & Security | Medium | [Where do the merchants' Google and Meta tokens go when an agent calls their APIs?](modules/daily-2026-10-17.md#q2) | Agent 调用商家的 Google 和 Meta API 时，这些 token 应该放在哪里？ |
| 3 | Context & Memory | Hard | [The agent read a 200-page supplier contract and missed the clause on page 143](modules/daily-2026-10-17.md#q3) | Agent 读了一份 200 页的供应商合同，却漏掉了第 143 页的条款 |
| 4 | Evals & Observability | Hard | [Task success fell from 91% to 84% overnight, and nobody deployed anything](modules/daily-2026-10-17.md#q4) | 任务成功率一夜之间从 91% 掉到 84%，而没有任何人发布过东西 |
| 5 | Agent Loop | Senior | [An agent that waits hours for a merchant's approval, then continues safely](modules/daily-2026-10-17.md#q5) | Agent 等待商家审批几个小时后，如何安全地继续执行 |

## [Daily set · Oct 18 · 每日新题 · 10 月 18 日](modules/daily-2026-10-18.md)

Five new agent questions: a large overnight copy-generation job on a budget, migrating to a new embedding model, stopping an agent from making promises it is not allowed to make, keeping subagent fan-out under control, and designing a public MCP server for merchants' own assistants.

> 五道新的 Agent 题：在预算内完成一个大型的夜间文案生成任务、迁移到新的 embedding 模型、阻止 Agent 做出它无权做出的承诺、控制子 Agent 的扇出规模，以及为商家自己的助手设计一个公开的 MCP server。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Production | Medium | [Ad copy for 50,000 products by tomorrow morning, on a fixed budget](modules/daily-2026-10-18.md#q1) | 明早之前为 5 万个商品生成广告文案，而且预算固定 |
| 2 | RAG | Medium | [Switching to a new embedding model without breaking search](modules/daily-2026-10-18.md#q2) | 更换 embedding 模型，同时不让检索出问题 |
| 3 | Safety & Security | Hard | [The support agent promised a full refund that the policy does not allow](modules/daily-2026-10-18.md#q3) | 客服 Agent 承诺了政策不允许的全额退款 |
| 4 | Multi-agent | Hard | ["Optimize all my campaigns" turned into 300 subagents and a huge bill](modules/daily-2026-10-18.md#q4) | 一句“优化我所有的广告系列”变成了 300 个子 Agent 和一张巨额账单 |
| 5 | Tools & MCP | Senior | [Letting merchants' own AI assistants manage their store through your MCP server](modules/daily-2026-10-18.md#q5) | 让商家自己的 AI 助手通过你们的 MCP server 管理店铺 |

## [Daily set · Oct 19 · 每日新题 · 10 月 19 日](modules/daily-2026-10-19.md)

Five new agent questions: getting parallel tool calls to happen, reading a close prompt comparison correctly, safe and correct text-to-SQL, migrating agents to a newer model, and designing a shared tracing and evaluation platform for many teams.

> 五道新的 Agent 题：让并行工具调用真正发生、正确解读一次势均力敌的提示词对比、安全且正确的文本转 SQL、把 Agent 迁移到更新的模型，以及为多个团队设计共享的追踪与评测平台。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Agent Loop | Medium | [The agent fetches six independent reports one at a time, in six turns](modules/daily-2026-10-19.md#q1) | Agent 把六份互不相关的报表逐个拉取，用了六轮 |
| 2 | Evals & Observability | Medium | [Prompt B beat prompt A 52 to 48 in an LLM-judged comparison. Ship it?](modules/daily-2026-10-19.md#q2) | 在 LLM 评审的对比中，提示词 B 以 52 比 48 胜过提示词 A。要上线吗？ |
| 3 | RAG | Hard | [Letting the agent answer "which of my products..." questions with SQL](modules/daily-2026-10-19.md#q3) | 让 Agent 用 SQL 回答“我的哪些商品……”这类问题 |
| 4 | Production | Hard | [Moving every agent to a newer model without surprising 3,000 merchants](modules/daily-2026-10-19.md#q4) | 把所有 Agent 迁移到更新的模型，同时不让 3000 个商家感到意外 |
| 5 | Evals & Observability | Senior | [Twenty teams build agents. Design the shared tracing and evaluation platform](modules/daily-2026-10-19.md#q5) | 二十个团队都在构建 Agent，请设计共享的追踪与评测平台 |

## [Daily set · Oct 20 · 每日新题 · 10 月 20 日](modules/daily-2026-10-20.md)

Five new agent questions: trimming an overgrown system prompt, making a writer-critic loop converge, changing a tool's schema without breaking live conversations, stopping generated descriptions from inventing product claims, and threat-modelling an agent platform before launch.

> 五道新的 Agent 题：精简过度膨胀的系统提示词、让撰写者与评审者的循环收敛、在不破坏进行中对话的前提下修改工具 schema、阻止生成的描述编造商品宣称，以及在上线前为 Agent 平台做威胁建模。

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Context & Memory | Medium | [A 15,000-token system prompt, and the agent still ignores half the rules](modules/daily-2026-10-20.md#q1) | 系统提示词长达 1.5 万 token，Agent 却仍然无视一半的规则 |
| 2 | Multi-agent | Medium | [The writer and the critic agents argue in circles and never finish](modules/daily-2026-10-20.md#q2) | 撰写 Agent 和评审 Agent 来回争执，永远完不成 |
| 3 | Tools & MCP | Hard | [Renaming one tool parameter broke the agent for a day](modules/daily-2026-10-20.md#q3) | 改了一个工具参数名，Agent 坏了一整天 |
| 4 | RAG | Hard | [Generated product descriptions claim features the product does not have](modules/daily-2026-10-20.md#q4) | 生成的商品描述宣称商品具备它其实没有的功能 |
| 5 | Safety & Security | Senior | [Threat-model the agent platform before it launches to merchants](modules/daily-2026-10-20.md#q5) | 在面向商家上线之前，为 Agent 平台做威胁建模 |

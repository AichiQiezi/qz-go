# 04. Retrieval and knowledge for agents · Agent 的检索与知识

How agents ground answers in your documents and data: the retrieval pipeline, agentic multi-step retrieval and how to evaluate it, and the multi-tenant permission, freshness and deletion problems behind real incidents. Interviewers ask because most production agents depend on retrieval, and many of their wrong answers are really retrieval failures.

> Agent 如何基于你的文档和数据作答：检索流水线、Agent 式多步检索及其评估方法，以及引发真实事故的多租户权限、数据新鲜度和删除问题。面试官爱问这些，因为大多数生产环境中的 Agent 都依赖检索，而它们的很多错误回答其实是检索失败。

3 questions · generated 2026-10-15T18:24:17+08:00 · [← 目录 / Contents](../README.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | RAG | Medium | [Answering merchants from help-center and ad-policy docs](#q1) | 基于帮助中心和广告政策文档回答商家问题 |
| 2 | RAG | Hard | [Multi-hop questions break single-shot RAG](#q2) | 多跳问题难倒了单次检索的 RAG |
| 3 | RAG | Hard | [Merchant A's answer quoted merchant B's private notes](#q3) | 商家 A 的回答引用了商家 B 的私有笔记 |

---

<a id="q1"></a>

## 1. Answering merchants from help-center and ad-policy docs

**基于帮助中心和广告政策文档回答商家问题** · `RAG` · `Medium`

### Question · 题目

Merchants keep asking support questions such as "Why was my product disapproved for misrepresentation on Google?" or "Can I advertise CBD products on TikTok?" You are asked to build an assistant that answers from two sources: your company's help center (about 2,000 articles in several languages) and a curated copy of ad-platform policy pages (Google Merchant Center, Meta and TikTok advertising policies).

> 商家不断提出这样的客服问题：“为什么我的商品在 Google 上因为 misrepresentation（虚假陈述）被拒登？”或者“我能在 TikTok 上投放 CBD 产品的广告吗？”你需要构建一个助手，基于两个来源回答问题：公司的帮助中心（大约 2,000 篇文章，多种语言），以及一份经过整理的广告平台政策页面副本（Google Merchant Center、Meta 和 TikTok 的广告政策）。

Requirements: every answer cites its sources, and when the documents do not cover a question, the assistant says so instead of guessing. Walk through the retrieval pipeline you would build, from raw documents to a cited answer, and explain the trade-offs of each choice.

> 要求：每个回答都要注明来源；当文档没有覆盖某个问题时，助手要明确说出来，而不是去猜。请讲解你会构建的检索流水线，从原始文档一直到带引用的回答，并解释每个选择的权衡。

<details>
<summary><b>Answer · 参考答案</b></summary>

**The core idea.** Retrieval-augmented generation (RAG) means that before the model answers, you search your own documents for the passages most likely to contain the answer, put those passages into the prompt, and instruct the model to answer only from them and to cite them. The model supplies language and reasoning; the documents supply the facts. Most RAG quality problems are really search problems: if the right passage never reaches the prompt, no amount of prompt engineering will fix the answer.

> **核心思想。** 检索增强生成（retrieval-augmented generation，RAG）指的是：在模型回答之前，先在你自己的文档中搜索最有可能包含答案的段落，把这些段落放进提示词，并要求模型只依据它们作答、注明引用。模型提供语言和推理能力，文档提供事实。大多数 RAG 的质量问题其实是搜索问题：如果正确的段落根本没有进入提示词，再怎么做提示词工程也救不了这个回答。

**Ingestion: clean, chunk, enrich.**

> **数据摄入：清洗、分块、补充信息。**

- Clean: convert HTML to text while keeping structure (headings, lists, sensibly flattened tables). Strip navigation, cookie banners and footers, which otherwise match every query.
  > 清洗：把 HTML 转成文本，同时保留结构（标题、列表、合理展平的表格）。去掉导航栏、cookie 提示和页脚，否则它们会匹配上每一个查询。
- Chunk: split documents into passages that are specific enough to match a question but complete enough to make sense alone. Fixed-size windows with some overlap are simple but can cut a rule in half. Structure-aware chunking splits on headings, which suits policy pages: one policy section becomes one chunk.
  > 分块（chunking）：把文档切成段落，既要足够具体以匹配某个问题，又要足够完整以便单独阅读也能看懂。带一定重叠的固定长度窗口很简单，但可能把一条规则拦腰切断。按结构分块会在标题处切分，很适合政策页面：一个政策小节就是一个块。
- Add context to each chunk: prepend the document title and section path ("Google Merchant Center > Misrepresentation > ..."), so a chunk that says "this applies to all products" still makes sense on its own.
  > 给每个块补充上下文：在开头加上文档标题和章节路径（“Google Merchant Center > Misrepresentation > ...”），这样一个写着“这适用于所有商品”的块，单独拿出来也能看懂。
- Attach metadata: `source` (help_center, gmc, meta, tiktok), language, URL, last-updated date and a stable `chunk_id`. Metadata becomes your filters and your citations.
  > 附加元数据：`source`（help_center、gmc、meta、tiktok）、语言、URL、最后更新日期，以及一个稳定的 `chunk_id`。元数据既是你的过滤条件，也是你的引用依据。

**Search: dense plus sparse, then rerank.**

> **搜索：稠密检索加稀疏检索，然后重排序。**

- Dense retrieval: an embedding model turns each chunk and the query into vectors, and you search for nearest neighbours in a vector index (pgvector, OpenSearch or a dedicated vector database). It matches meaning, so "my item got rejected for lying about a discount" can find the misrepresentation policy.
  > 稠密检索（dense retrieval）：embedding 模型把每个块和查询都转换成向量，然后在向量索引（pgvector、OpenSearch 或专门的向量数据库）中搜索最近邻。它匹配的是语义，所以“我的商品因为折扣造假被拒了”也能找到 misrepresentation 政策。
- Sparse retrieval with BM25: classic keyword scoring that matches exact terms, which embeddings tend to blur: policy names, error messages, identifiers such as GTIN or MPN, and brand names.
  > 用 BM25 做稀疏检索：经典的关键词打分，能精确匹配词语，而这恰恰是 embedding 容易模糊掉的东西：政策名称、错误信息、GTIN 或 MPN 这类标识符，以及品牌名。
- Hybrid: run both and merge the ranked lists, for example with reciprocal rank fusion, which scores each chunk by summing 1/(k + rank) across the lists. Support questions mix paraphrase with exact terms, so hybrid usually beats either method alone.
  > 混合检索（hybrid search）：两种都跑，再合并排序结果，例如使用倒数排名融合（reciprocal rank fusion），即把每个块在各个列表中的 1/(k + rank) 相加作为得分。客服问题里既有换种说法的表述，也有精确术语，所以混合检索通常比单用任何一种都好。
- Rerank: take the top 30 to 50 candidates, score each (question, chunk) pair with a cross-encoder reranker or a small LLM, and keep the best 5 to 8. First-stage retrieval buys recall cheaply; the reranker buys precision for some extra latency.
  > 重排序（rerank）：取前 30 到 50 个候选，用 cross-encoder 重排序模型或一个小型 LLM 给每一对（问题，块）打分，保留最好的 5 到 8 个。第一阶段检索以低成本换取召回率；重排序模型以一些额外延迟换取精确率。

**Prompt assembly and citations.** Put the selected chunks into the prompt as clearly delimited documents with ids and metadata, followed by the instructions and the question:

> **组装提示词和引用。** 把选中的块作为界限清晰、带 id 和元数据的文档放进提示词，后面再跟上指令和问题：

```text
<doc id="gmc-misrep-03" source="Google Merchant Center" updated="2026-08-02">
...chunk text...
</doc>
<doc id="hc-1182-2" source="Help Center" updated="2026-09-10">
...chunk text...
</doc>

Answer the merchant's question using only the documents above.
Cite the doc id after each claim, like [gmc-misrep-03].
If the documents do not answer the question, say so and suggest contacting support.
```

Then check the output in code: every cited id must be one you actually supplied, and an answer without citations is retried or replaced with the fallback reply. Some APIs have built-in citation support (the Claude API, for example, can return the exact quoted passages from documents you pass in), which removes the risk of malformed or invented citations. Render each citation as a link to the source URL so merchants and support staff can check it.

> 然后在代码里检查输出：每个被引用的 id 都必须是你确实提供过的；没有引用的回答要重试，或者替换成兜底回复。有些 API 自带引用功能（例如 Claude API 可以返回你传入的文档中被引用的原文段落），这样就消除了引用格式错误或凭空捏造引用的风险。把每条引用渲染成指向来源 URL 的链接，方便商家和客服人员核对。

**Trade-offs to state in the interview.**

> **面试中要讲清楚的权衡。**

- Chunk size: small chunks are precise but lose context; large chunks keep context but dilute the embedding and use more prompt tokens. Structure-aware chunks with a title prefix are a good default; tune the size against an evaluation set.
  > 块的大小：小块精确但丢失上下文；大块保留上下文，但会稀释 embedding，也会占用更多提示词 token。按结构分块并加上标题前缀是一个不错的默认做法；具体大小用评测集来调。
- Embedding model: it must be multilingual here, because merchants ask in Chinese about English policy pages. Switching models means re-embedding the whole corpus, so version the index.
  > Embedding 模型：这里必须支持多语言，因为商家会用中文询问英文的政策页面。换模型意味着要对整个语料库重新生成 embedding，所以要给索引做版本管理。
- Top-k: more chunks raise recall but add cost and distraction. Retrieve wide, rerank, and keep few.
  > Top-k：取更多的块能提高召回率，但会增加成本和干扰。先宽检索、再重排序，最后只保留少数几个。
- Hybrid search and reranking add moving parts and latency. For a small corpus, keep them only if your evaluation shows a gain.
  > 混合检索和重排序会增加组件和延迟。如果语料库很小，只有在评测显示确实有提升时才保留它们。
- Freshness: policy pages change. Re-ingest when they change and show the update date in citations, so an outdated answer is easy to spot.
  > 新鲜度：政策页面会变化。在它们变化时重新摄入，并在引用中显示更新日期，这样过时的回答很容易被发现。

</details>

### Follow-ups · 追问

- How would you build a labelled evaluation set for this assistant?
  > 你会如何为这个助手构建一个带标注的评测集？
- A merchant asks in Chinese about an English policy page. What changes in the pipeline?
  > 商家用中文询问一个英文政策页面的内容。流水线中需要改变什么？
- How do you stop the model from answering from its training data when retrieval finds nothing relevant?
  > 当检索没有找到相关内容时，你如何防止模型凭训练数据作答？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| chunking | 分块：把文档切成适合检索的段落 |
| hybrid search | 混合检索（向量检索加关键词检索） |
| reranker | 重排序模型 |
| grounded answer | 有据可依的回答（只基于检索到的内容作答） |
| If the right passage never reaches the prompt, the model can't use it. | 如果正确的段落根本没有进入提示词，模型就无从使用它。 |

---

<a id="q2"></a>

## 2. Multi-hop questions break single-shot RAG

**多跳问题难倒了单次检索的 RAG** · `RAG` · `Hard`

### Question · 题目

Your support assistant uses classic RAG: embed the merchant's question, fetch the top 5 chunks, generate an answer. It handles "What are Google's GTIN requirements?" well. It fails on questions like this one:

> 你们的客服助手使用经典 RAG：把商家的问题转成 embedding，取前 5 个块，然后生成回答。它能很好地回答“Google 对 GTIN 有什么要求？”这类问题，却答不好下面这种问题：

"Which of my products were disapproved last week for the same policy reason as SKU 123, and what changed in that policy recently?"

> “上周我有哪些商品因为和 SKU 123 相同的政策原因被拒登？这条政策最近有什么变化？”

It retrieves five loosely related policy chunks, invents a list of products, and says nothing concrete about the policy change. Explain why single-shot RAG fails here, how you would redesign it as agentic retrieval, when you would still keep single-shot RAG, and how you would evaluate retrieval separately from generation.

> 它检索出五个关系不大的政策块，编造了一份商品列表，对政策变化也没说出任何具体内容。请解释单次检索的 RAG 为什么在这里失败、你会如何把它重新设计成 Agent 式检索、什么情况下你仍然会保留单次检索的 RAG，以及你会如何把检索和生成分开评估。

<details>
<summary><b>Answer · 参考答案</b></summary>

**Why single-shot RAG fails here.** Classic RAG makes one retrieval decision, up front, from the user's own wording. This question is a chain of dependent lookups: why SKU 123 was disapproved (structured data, not a document); which other products had the same reason last week (a database query); which policy section covers that reason; and what changed in it (a comparison of versions). Each step needs the previous result, and the third step needs words, such as the official policy name, that are not in the question at all. One embedding of the whole question lands somewhere between "disapproval" and "policy change" and pulls generic chunks. And since nothing fetched the real product list, the model invented one.

> **单次检索的 RAG 为什么在这里失败。** 经典 RAG 只做一次检索决策，而且是在一开始、直接用用户的原话去检索。这个问题其实是一串相互依赖的查询：SKU 123 为什么被拒登（这是结构化数据，不是文档）；上周还有哪些商品因同样原因被拒登（一次数据库查询）；哪个政策小节对应这个原因；这一节发生了什么变化（版本之间的比较）。每一步都需要上一步的结果，而第三步需要用到问题中根本没有出现的词，比如官方的政策名称。把整个问题做一次 embedding，得到的向量落在“拒登”和“政策变化”之间的某个位置，取回的都是泛泛的块。又因为没有任何一步真正取回商品列表，模型就编了一份。

**Agentic RAG: retrieval as a tool.** Give the agent retrieval tools and let it decide what to look up next, in a loop, based on what it has found so far (the ReAct pattern applied to search):

> **Agent 式 RAG：把检索当作工具。** 给 Agent 提供检索工具，让它在循环中根据目前已经找到的内容，自己决定下一步查什么（这就是把 ReAct 模式应用到搜索上）：

- `get_product_status(sku)`: a structured lookup in your diagnostics store or the platform API that returns the disapproval reason and date.
  > `get_product_status(sku)`：在你的诊断数据存储或平台 API 中做结构化查询，返回拒登原因和日期。
- `find_disapprovals(reason, from, to)`: a fixed, parameterized query scoped to the current merchant, not free-form SQL written by the model.
  > `find_disapprovals(reason, from, to)`：一个固定的参数化查询，限定在当前商家范围内，而不是由模型自由编写的 SQL。
- `search_policies(query, platform, as_of)`: hybrid search over policy chunks, filtered by platform and effective date.
  > `search_policies(query, platform, as_of)`：对政策块做混合检索，按平台和生效日期过滤。
- `get_policy_history(policy_id)`: version history and a diff, if you store versioned policy documents.
  > `get_policy_history(policy_id)`：如果你存储了带版本的政策文档，就返回版本历史和 diff。

A plausible trace for the question above:

> 针对上面那个问题，一条合理的执行轨迹是：

```text
get_product_status("123")        -> reason "misrepresentation", 2026-10-08
find_disapprovals("misrepresentation", "2026-10-05", "2026-10-11") -> 14 products
search_policies("misrepresentation promotions", "google", "2026-10-15") -> gmc-misrep-03
get_policy_history("gmc-misrep") -> v7 effective 2026-08-02, diff against v6
```

The answer is then built only from what the tools returned, each fact with its source. Cap the loop (a maximum number of tool calls) and give the agent an explicit way to say "not found", so it cannot search forever.

> 最终回答只基于工具返回的内容构建，每条事实都带有来源。给循环设置上限（最大工具调用次数），并给 Agent 一个明确说“没找到”的出口，这样它就不会无休止地搜索下去。

**Query rewriting.** Merchants use their own words; documents use the platform's words. Before searching, have a model (often a cheaper one) rewrite the question into one or more queries: expand abbreviations, translate a Chinese question to match English policy pages, add the official policy name found in an earlier step, and split compound questions into sub-queries. A related technique, HyDE, embeds a hypothetical answer instead of the question, which helps when questions and documents are phrased very differently.

> **查询改写（query rewriting）。** 商家用自己的话提问，文档用的是平台的措辞。搜索之前，先让一个模型（通常是更便宜的模型）把问题改写成一个或多个查询：展开缩写，把中文问题翻译成能匹配英文政策页面的说法，加上前一步查到的官方政策名称，并把复合问题拆成子查询。还有一种相关技术叫 HyDE：它不对问题本身做 embedding，而是对一个假设性的回答做 embedding；当问题和文档的表述方式差异很大时，这很有帮助。

**Structured data plus documents.** Facts about the merchant's own products (status, dates, counts) belong in a database or the platform API, not in a vector index: embedding them is lossy and goes stale, while a query is exact and current. Use documents for "what does the rule say and why", structured tools for "which, how many, when", and let the agent join the two.

> **结构化数据加文档。** 关于商家自己商品的事实（状态、日期、数量）应该放在数据库或平台 API 中，而不是向量索引里：对它们做 embedding 既有损失又会过时，而查询是精确且实时的。文档用来回答“规则怎么说、为什么”，结构化工具用来回答“哪些、多少、什么时候”，再由 Agent 把两者结合起来。

**When single-shot RAG is still better.** Agentic retrieval takes several model calls instead of one, so it is slower, more expensive and less predictable: the same question can take different paths on different runs. For the many single-hop FAQ questions ("what image size does Meta require?"), a tuned pipeline with hybrid search and reranking is faster, cheaper and easier to test. A common production design is a router: classify the question, send simple ones down the fixed pipeline, and only multi-hop or data-dependent ones to the agent.

> **什么情况下单次检索的 RAG 仍然更好。** Agent 式检索需要多次模型调用而不是一次，所以更慢、更贵，也更难预测：同一个问题在不同的运行中可能走不同的路径。对于大量单跳的 FAQ 问题（“Meta 要求的图片尺寸是多少？”），一条调优过、带混合检索和重排序的流水线更快、更便宜，也更容易测试。常见的生产设计是加一个路由器：先对问题分类，简单的问题走固定流水线，只有多跳或依赖数据的问题才交给 Agent。

**Evaluate retrieval separately from generation.** When an answer is wrong, you need to know whether search failed or the model did. Build a labelled set of a few hundred real questions, each tagged with the chunk ids that contain the answer (an LLM can draft the labels; support experts review them). Then measure retrieval on its own:

> **把检索和生成分开评估。** 当回答出错时，你需要知道是搜索失败了还是模型失败了。构建一个带标注的数据集：几百个真实问题，每个问题都标注出包含答案的块 id（可以由 LLM 起草标注，再由客服专家审核）。然后单独衡量检索效果：

- Recall@k: the share of questions with at least one relevant chunk in the top k (or the share of all relevant chunks found). It sets the ceiling for answer quality.
  > Recall@k（前 k 召回率）：前 k 个结果中至少有一个相关块的问题所占的比例（或者，所有相关块中被找到的比例）。它决定了回答质量的上限。
- Precision@k: the share of the top k that is relevant. Low precision means more noise in the prompt.
  > Precision@k（前 k 精确率）：前 k 个结果中相关结果所占的比例。精确率低意味着提示词里的噪声更多。
- MRR (mean reciprocal rank): the average of 1/rank of the first relevant result, which rewards putting the right chunk near the top.
  > MRR（mean reciprocal rank，平均倒数排名）：第一个相关结果排名的倒数（1/rank）的平均值，它奖励把正确的块排在靠前的位置。

For agentic retrieval, also check the trajectory: did it call the structured tool for SKU 123, with the right merchant and dates? Then evaluate generation given the retrieved context: faithfulness (every claim is supported by a cited chunk), correctness against a reference answer, and correct refusals when nothing relevant exists. Run the whole set on every change to chunking, embeddings, rewriting or prompts, like a regression suite.

> 对于 Agent 式检索，还要检查执行轨迹：它有没有为 SKU 123 调用结构化工具，用的商家和日期对不对？然后在给定检索上下文的前提下评估生成效果：忠实度（每个论断都有被引用的块支持）、与参考答案相比的正确性，以及在确实没有相关内容时能否正确拒答。每次修改分块、embedding、查询改写或提示词时，都把整套数据跑一遍，就像跑回归测试一样。

</details>

### Follow-ups · 追问

- How do you stop the agent from searching in circles when the answer does not exist?
  > 当答案根本不存在时，你如何防止 Agent 兜圈子反复搜索？
- How would you build the labelled set without months of manual work?
  > 你如何在不花几个月人工的情况下构建这个带标注的数据集？
- Would you let the agent write its own SQL? What guardrails would that need?
  > 你会让 Agent 自己写 SQL 吗？那需要哪些防护措施？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| multi-hop question | 多跳问题：需要多步、前后依赖的查询才能回答的问题 |
| query rewriting | 查询改写 |
| recall@k | 前 k 召回率 |
| mean reciprocal rank (MRR) | 平均倒数排名 |
| Retrieval quality sets the ceiling for answer quality. | 检索质量决定了回答质量的上限。 |

---

<a id="q3"></a>

## 3. Merchant A's answer quoted merchant B's private notes

**商家 A 的回答引用了商家 B 的私有笔记** · `RAG` · `Hard`

### Question · 题目

Your RAG assistant serves thousands of merchants from one shared vector index. Besides public help-center and policy documents, merchants can upload private notes (supplier terms, internal pricing rules), which are indexed in the same collection. Retrieval works like this: take the top 20 chunks by similarity across the whole index, generate the answer, then remove any citations that point to another merchant's documents before showing it.

> 你们的 RAG 助手用一个共享的向量索引服务几千个商家。除了公开的帮助中心和政策文档，商家还可以上传私有笔记（供应商条款、内部定价规则），这些笔记被索引在同一个集合里。检索流程是这样的：在整个索引中按相似度取前 20 个块，生成回答，然后在展示之前删掉所有指向其他商家文档的引用。

In one week you get two incidents:

> 一周之内出了两起事故：

- Merchant A asked about return rules and received an answer that paraphrased merchant B's private supplier terms, with no visible citation.
  > 商家 A 询问退货规则，得到的回答转述了商家 B 的私有供应商条款，而且看不到任何引用。
- Several answers cited a TikTok policy version that was replaced two months ago.
  > 好几个回答引用了一个两个月前就已被替换的 TikTok 政策版本。

Explain what went wrong. Design retrieval that respects tenant permissions, stays fresh, and supports deletion when a merchant asks for their data to be removed. How would you test for cross-tenant leakage?

> 请解释哪里出了问题。设计一套检索方案，让它遵守租户权限、保持数据新鲜，并且在商家要求删除自己的数据时支持删除。你会如何测试跨租户泄露？

<details>
<summary><b>Answer · 参考答案</b></summary>

**What went wrong.** The tenant check happened in the wrong place. Filtering after generation is like running SQL without `WHERE merchant_id = ?`, sending every row to the client and hiding the extra rows in the UI: the data has already crossed the boundary. Once merchant B's notes were in the prompt, the model could use and paraphrase them; removing the citation only hid the source. The outdated policy is the same bug on another axis: nothing in the query said "only the version in force today".

> **哪里出了问题。** 租户检查放错了位置。在生成之后再过滤，就好比执行一条没有 `WHERE merchant_id = ?` 的 SQL，把所有行都发给客户端，再在 UI 里把多余的行藏起来：数据早已越过了边界。一旦商家 B 的笔记进入了提示词，模型就可以使用并转述它们；删掉引用只是藏起了出处。过时政策是同一个 bug，只是换了一个维度：查询中没有任何条件说明“只要今天正在生效的版本”。

**Rule one: the permission filter is part of the query.** Every retrieval call carries the tenant from the authenticated session, never from the model or the user's text, and the index enforces it:

> **第一条原则：权限过滤是查询的一部分。** 每次检索调用都携带来自已认证会话的租户信息，绝不来自模型或用户的文本，并由索引来强制执行：

- Pre-filter inside the search: `tenant_id IN (current_merchant, "public")` is part of both the vector and the keyword query, so other tenants' chunks are never candidates. Check how your vector store implements filters; filtering after an approximate top-k search can return too few results, or none.
  > 在搜索内部做前置过滤：`tenant_id IN (current_merchant, "public")` 同时是向量查询和关键词查询的一部分，这样其他租户的块永远不会成为候选。要弄清你的向量存储是如何实现过滤的：在近似 top-k 搜索之后再过滤，可能返回过少的结果，甚至一个都没有。
- For stronger isolation, keep each tenant's private data in its own namespace or partition, with a shared one for public documents, and query both. A forgotten filter then returns nothing rather than someone else's data.
  > 如果需要更强的隔离，就把每个租户的私有数据放在各自的命名空间（namespace）或分区里，公开文档放在一个共享的命名空间里，两者都查询。这样即使忘了加过滤条件，返回的也是空结果，而不是别人的数据。
- Enforce it in one place: a retrieval service that cannot be called without a tenant, much like row-level security in Postgres. The agent's search tool takes only the query text; the harness injects the tenant.
  > 在一个地方统一强制执行：一个不带租户就无法调用的检索服务，很像 Postgres 中的行级安全（row-level security）。Agent 的搜索工具只接收查询文本，由编排代码注入租户。
- Audit everything else on the path: a semantic answer cache keyed only by question text leaks just as well, and so can shared conversation memory or evaluation sets built from production logs.
  > 审查链路上的其他所有环节：只以问题文本为键的语义答案缓存同样会泄露数据，共享的对话记忆和基于生产日志构建的评测集也可能泄露。

```go
// The model never chooses the tenant: the handler takes it from the session.
func (h *SearchTool) Handle(ctx context.Context, in SearchInput) (ToolResult, error) {
    tenant := auth.MerchantID(ctx) // set by auth middleware from the verified token
    hits, err := h.index.Search(ctx, Query{
        Text:   in.Query,
        Filter: Filter{TenantIn: []string{tenant, "public"}, EffectiveAt: time.Now()},
        TopK:   20,
    })
    if err != nil {
        return ToolResult{}, err
    }
    return toToolResult(hits), nil
}
```

**Freshness: version documents instead of overwriting them.**

> **新鲜度：给文档做版本管理，而不是直接覆盖。**

- Give every document a stable id and a version, and store `effective_from` and `effective_to` on each chunk. Default queries filter to the version in force now; a "what changed?" question can ask for older versions explicitly.
  > 给每个文档一个稳定的 id 和一个版本号，并在每个块上存储 `effective_from` 和 `effective_to`。默认查询只过滤出当前生效的版本；“有什么变化？”这类问题可以明确地查询旧版本。
- Re-index from change events, not only from a nightly rebuild: change data capture (CDC) on the help-center database, and a crawler that hashes policy pages and emits an event when the content changes. An indexer consumes the events from a queue, re-chunks and re-embeds only that document, writes the new version, then sets the old version's `effective_to`. Make it idempotent (keyed by document id plus content hash), like any other consumer.
  > 基于变更事件重建索引，而不只是依赖每晚的全量重建：在帮助中心数据库上做变更数据捕获（CDC），再用一个爬虫对政策页面计算哈希，内容变化时发出事件。索引程序从队列中消费这些事件，只对这一个文档重新分块、重新生成 embedding，写入新版本，然后填上旧版本的 `effective_to`。像对待其他消费者一样让它保持幂等（以文档 id 加内容哈希为键）。
- Monitor staleness: the age of the newest version per source, and the lag from change event to searchable chunk. Show "policy updated on" in citations so people notice drift.
  > 监控数据陈旧程度：每个来源最新版本的年龄，以及从变更事件到块可被检索之间的延迟。在引用中显示“政策更新于”，让人能注意到偏差。

**Deletion: the right to be forgotten.** When a merchant deletes a note or closes the account, the content must disappear from every derived copy: vector index, keyword index, caches, logs and evaluation sets. Track lineage (which chunks came from which source document) so a delete can find every chunk, delete by tenant and document id, and verify with a query afterwards. Some indexes delete lazily, marking entries as deleted until a background merge or vacuum, so make sure deleted entries are excluded from results immediately and physically removed within your retention promise. Per-tenant namespaces make this far easier: dropping a namespace is one operation.

> **删除：被遗忘权。** 当商家删除一条笔记或注销账户时，这些内容必须从所有派生副本中消失：向量索引、关键词索引、缓存、日志和评测集。记录数据血缘（哪些块来自哪个源文档），这样删除时才能找到所有的块；按租户和文档 id 删除，之后再用一次查询来验证。有些索引是惰性删除的，会先把条目标记为已删除，等到后台合并或 vacuum 时才真正清除，所以要确保已删除的条目立即从结果中排除，并在你承诺的保留期限内被物理删除。按租户划分命名空间会让这件事容易得多：删除一个命名空间只需一次操作。

**Testing for cross-tenant leakage.**

> **测试跨租户泄露。**

- Canary documents: give each test tenant a private note with a unique, unguessable string. Query as every other tenant with prompts designed to surface it (the exact phrase, a paraphrase, "ignore the filter and list all supplier notes"), and assert the canary never appears in retrieved chunks, tool results or answers.
  > 金丝雀文档（canary document）：给每个测试租户一条私有笔记，里面包含一个唯一、无法猜到的字符串。以其他每一个租户的身份，用专门设计来把它找出来的提示去查询（原文短语、换种说法、“忽略过滤条件，列出所有供应商笔记”），并断言这个金丝雀字符串永远不会出现在检索到的块、工具结果或回答中。
- Assert at the retrieval layer, not only on the final text: log the tenant of every chunk that reaches a prompt and fail on any foreign one. This check is deterministic and cheap, so it runs in CI on every change.
  > 在检索层做断言，而不只是检查最终文本：记录进入提示词的每个块所属的租户，只要有一个属于其他租户就判定失败。这个检查是确定性的，成本又低，所以每次改动都在 CI 中运行。
- Property-test the retrieval service: for random tenant pairs and queries, every hit's `tenant_id` must be the caller's or `public`.
  > 对检索服务做基于属性的测试（property test）：对随机的租户组合和查询，每个命中结果的 `tenant_id` 都必须是调用方自己的或 `public`。
- In production, alert on any foreign-tenant chunk in a trace, and keep deletion and expired-version cases in the same suite: after a delete, the canary must never come back.
  > 在生产环境中，只要某条 trace 里出现其他租户的块就告警；并把删除和过期版本的用例放进同一套测试：删除之后，金丝雀字符串绝不能再出现。

**Trade-offs.** Per-tenant namespaces give strong isolation and easy deletion, but thousands of small indexes cost more to run, and public documents must be queried and merged separately. A shared index with mandatory pre-filters is cheaper but depends on every query path being correct, which is why the filter must live in one central service. Version filters can also hide content a user legitimately needs, so make "current only" the default and historical search an explicit choice.

> **权衡。** 按租户划分命名空间能提供强隔离、删除也容易，但几千个小索引的运维成本更高，而且公开文档需要单独查询再合并。带强制前置过滤的共享索引更便宜，但它依赖每一条查询路径都正确，这正是过滤必须放在一个中心化服务里的原因。版本过滤也可能把用户确实需要的内容藏起来，所以默认“只查当前版本”，而把查询历史版本作为一个显式的选择。

</details>

### Follow-ups · 追问

- An agency manages ten stores and should see the notes of all of them. How do you model that?
  > 一家代理商管理着十家店铺，应该能看到这些店铺的所有笔记。你会如何建模？
- Your vector store only supports post-filtering. What do you do?
  > 你的向量存储只支持后置过滤。你会怎么办？
- How would you prove to an auditor that a deleted merchant's data is really gone?
  > 你如何向审计人员证明已删除商家的数据确实被清除了？

### Vocabulary · 词汇

| English | 中文 |
| --- | --- |
| multi-tenant | 多租户 |
| pre-filtering | 前置过滤：在检索过程中就按权限等条件过滤 |
| change data capture (CDC) | 变更数据捕获 |
| right to be forgotten | 被遗忘权 |
| canary document | 金丝雀文档：用于检测数据泄露的唯一标记数据 |

# 每日面试题 · Daily Interview Questions

后端面试题（Go 为主、Java 为辅），中英双语，每天 5 道，覆盖 Go、System Design、Database、Middleware、Java & CS。
每题包含题目、参考答案（默认折叠）、追问和词汇表。

共 34 天、170 道题。

- 数据：`interview/data/YYYY-MM-DD.json`（与 Interview Drill 页面的 `days` 集合结构一致）
- 重新生成 Markdown：在仓库根目录执行 `go run ./interview/render`
- Agent 面试速成：[agent/README.md](agent/README.md)

## [2026-11-14](days/2026-11-14.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [The product API hangs under load while the database is almost idle, and only some endpoints recover on their own](days/2026-11-14.md#q1) | 商品 API 在负载下卡死，而数据库几乎是空闲的，只有部分接口能自己恢复 |
| 2 | System Design | Senior | [One app's bulk sync slows the Admin API for everyone: rate limiting third-party apps on your own platform](days/2026-11-14.md#q2) | 一个应用的批量同步拖慢了所有人的 Admin API：为自己平台上的第三方应用做限流 |
| 3 | Database | Medium | [A SKU prefix search uses the index in staging and scans 30 million rows in production](days/2026-11-14.md#q3) | SKU 前缀搜索在预发环境走索引，在生产环境却扫描了 3000 万行 |
| 4 | Middleware | Medium | [A Redis bitmap of products viewed today: one new product ID made Redis allocate 500 MB and stall](days/2026-11-14.md#q4) | 用 Redis 位图记录今天被浏览过的商品：一个新的商品 ID 让 Redis 分配了 500 MB 内存并卡顿 |
| 5 | Java & CS | Hard | [Campaigns created in late December get end dates a year too late on US servers and in the past on EU servers](days/2026-11-14.md#q5) | 在 12 月底创建的广告活动，在美国服务器上结束日期晚了一年，在欧洲服务器上则变成了过去的日期 |

## [2026-11-13](days/2026-11-13.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [A merchant's CSV import says the "sku" column is missing, and another file fails on a product titled 12" frame](days/2026-11-13.md#q1) | 商家的 CSV 导入提示缺少「sku」列，另一个文件则在一个名为 12" frame 的商品上失败 |
| 2 | System Design | Hard | [Checkout shows shipping rates from three carriers within two seconds, and the price charged must match the price shown](days/2026-11-13.md#q2) | 结账页要在两秒内展示三家承运商的运费，而且最终收取的价格必须与展示的价格一致 |
| 3 | Database | Hard | [Get-or-create for product tags sometimes returns no ID, and in another service fails with a serialization error](days/2026-11-13.md#q3) | 商品标签的「获取或创建」有时拿不到 ID，在另一个服务里则因序列化错误失败 |
| 4 | Middleware | Senior | [Unpaid orders are cancelled by Redis expiry events: some are never cancelled, some are late, and some are cancelled twice](days/2026-11-13.md#q4) | 未支付订单靠 Redis 过期事件来取消：有的永远没被取消，有的取消得很晚，还有的被取消了两次 |
| 5 | Java & CS | Medium | [A token cache hits 99% of the time, yet the ad platform still gets one token request per call and starts returning 429](days/2026-11-13.md#q5) | 令牌缓存的命中率是 99%，广告平台却仍然在每次调用时收到一次令牌请求，并开始返回 429 |

## [2026-11-12](days/2026-11-12.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [A worker selects on an urgent queue and a bulk queue, yet order cancellations wait behind catalog syncs](days/2026-11-12.md#q1) | 一个 worker 同时在紧急队列和批量队列上 select，订单取消却仍然排在商品目录同步后面 |
| 2 | System Design | Medium | [Generating sitemaps and product feeds for 100,000 storefronts every day without hammering the database](days/2026-11-12.md#q2) | 每天为 10 万个店铺生成站点地图和商品 feed，同时不压垮数据库 |
| 3 | Database | Hard | [Checkout fails with "reached maximum value of sequence" on order_items, and the obvious fix would lock the table for an hour](days/2026-11-12.md#q3) | 结账时报错「序列已达到最大值」，出在 order_items 表上，而显而易见的修复会锁表一个小时 |
| 4 | Middleware | Medium | [A flash-sale product's view counter makes one Redis Cluster node hot while the others are idle](days/2026-11-12.md#q4) | 一个限时抢购商品的浏览计数器让 Redis Cluster 的某个节点过热，而其他节点却很空闲 |
| 5 | Java & CS | Senior | [Sizing the thread pool and connection pool for calls to a slow ad platform: Little's law, queues and timeouts](days/2026-11-12.md#q5) | 为调用缓慢广告平台的服务设置线程池和连接池大小：Little 定律、队列与超时 |

## [2026-11-11](days/2026-11-11.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [Merchants' flash sales end hours early: promotion times parsed without a time zone](days/2026-11-11.md#q1) | 商家的限时促销提前几个小时结束：促销时间在解析时没有带时区 |
| 2 | System Design | Senior | [Billing merchants for a subscription plus a share of sales: mid-month plan changes, failed card payments and invoices that must never change](days/2026-11-11.md#q2) | 向商家收取订阅费加销售抽成：月中更换套餐、信用卡扣款失败，以及永远不能更改的发票 |
| 3 | Database | Medium | ["Highest price first" shows products without a price at the top, and the fix makes the query slow](days/2026-11-11.md#q3) | 「价格从高到低」把没有价格的商品排在了最前面，而修复之后查询又变慢了 |
| 4 | Middleware | Hard | [A Lua script that recalculates stock for a large merchant makes Redis answer BUSY to every service, and SCRIPT KILL refuses to stop it](days/2026-11-11.md#q4) | 一个为大商家重新计算库存的 Lua 脚本让 Redis 对所有服务都返回 BUSY，而 SCRIPT KILL 也拒绝终止它 |
| 5 | Java & CS | Hard | [Deploys hang until the pod is killed because a sync worker swallows InterruptedException](days/2026-11-11.md#q5) | 部署一直卡住直到 Pod 被强制杀掉，原因是一个同步 worker 吞掉了 InterruptedException |

## [2026-11-10](days/2026-11-10.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [Large order exports end after exactly 30 seconds with a truncated CSV and HTTP 200](days/2026-11-10.md#q1) | 大额订单导出恰好在 30 秒时结束，得到一个被截断的 CSV，状态码却是 200 |
| 2 | System Design | Hard | [Moving a merchant's whole store from another e-commerce platform: products, customers, orders and SEO redirects, under the source's rate limits](days/2026-11-10.md#q2) | 把商家的整个店铺从另一个电商平台迁移过来：商品、客户、订单和 SEO 重定向，同时受限于源平台的限流 |
| 3 | Database | Medium | [After importing orders with their original IDs, every new order fails with a duplicate key](days/2026-11-10.md#q3) | 导入带原始 ID 的订单之后，每一笔新订单都因重复键而失败 |
| 4 | Middleware | Medium | [Unique visitors per ad campaign: a Redis set per campaign uses gigabytes, and HyperLogLog counts are slightly off](days/2026-11-10.md#q4) | 按广告系列统计独立访客：每个广告系列一个 Redis 集合，占用了好几个 GB，而 HyperLogLog 的计数又略有偏差 |
| 5 | Java & CS | Senior | [Product updates from three services are resolved by timestamp, and a merchant's latest price change is lost because one server's clock is behind](days/2026-11-10.md#q5) | 来自三个服务的商品更新按时间戳决定胜负，结果因为一台服务器的时钟偏慢，商家最新的改价丢失了 |

## [2026-11-09](days/2026-11-09.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [Calls to an ad platform start failing with 401 after the platform moves its API to a new host behind a redirect](days/2026-11-09.md#q1) | 广告平台把 API 迁到一个新主机并加上重定向之后，调用开始返回 401 |
| 2 | System Design | Hard | [Merchants want custom reports over two years of orders and ad spend: designing analytics that never touch the OLTP database](days/2026-11-09.md#q2) | 商家想基于两年的订单和广告花费自定义报表：设计一套完全不碰 OLTP 数据库的分析系统 |
| 3 | Database | Medium | [Two accounts for the same email, and login lookups scan the whole table: case-insensitive emails in Postgres](days/2026-11-09.md#q3) | 同一个邮箱出现了两个账户，登录查询还在扫全表：Postgres 中不区分大小写的邮箱 |
| 4 | Middleware | Hard | [A Redis cleanup job stops after the first empty SCAN page and leaves most expired sessions behind](days/2026-11-09.md#q4) | Redis 清理任务在 SCAN 第一次返回空页后就停了，留下了大部分过期会话 |
| 5 | Java & CS | Senior | [A Bloom filter deduplicates conversion events, and some real purchases are silently dropped](days/2026-11-09.md#q5) | 用布隆过滤器给转化事件去重，结果有些真实的购买被悄悄丢掉了 |

## [2026-11-08](days/2026-11-08.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [Webhook bodies sometimes contain another merchant's order after adding a sync.Pool for buffers](days/2026-11-08.md#q1) | 为缓冲区加上 sync.Pool 之后，webhook 请求体里偶尔出现另一个商家的订单 |
| 2 | System Design | Medium | [Custom domains for 100,000 storefronts: verifying ownership, issuing TLS certificates and renewing them on time](days/2026-11-08.md#q2) | 为 10 万个店铺提供自定义域名：验证所有权、签发 TLS 证书并按时续期 |
| 3 | Database | Medium | [Removing a churned store takes minutes and stalls that store's order writes: ON DELETE CASCADE without an index](days/2026-11-08.md#q3) | 删除一个已流失的店铺要好几分钟，期间该店铺的订单写入都被卡住：没有索引的 ON DELETE CASCADE |
| 4 | Middleware | Hard | [A Redis idempotency key says a conversion was already sent, until a failover makes Meta receive it twice](days/2026-11-08.md#q4) | Redis 幂等键表明一次转化已经发送过，直到一次故障转移让 Meta 收到了两次 |
| 5 | Java & CS | Senior | [A Java service with -Xmx2g is OOMKilled at 3 GB: finding the memory outside the heap](days/2026-11-08.md#q5) | 一个设置了 -Xmx2g 的 Java 服务在 3 GB 时被 OOMKilled：找出堆之外的内存 |

## [2026-11-07](days/2026-11-07.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [An exchange-rate loader wrapped in sync.Once fails once at startup, and the service never recovers](days/2026-11-07.md#q1) | 用 sync.Once 包装的汇率加载器在启动时失败了一次，服务就再也恢复不了 |
| 2 | System Design | Medium | [Ad-spend reports that keep changing for days: importing restated numbers from Meta, Google and TikTok](days/2026-11-07.md#q2) | 过了好几天还在变的广告花费报表：从 Meta、Google 和 TikTok 导入会被修正的数据 |
| 3 | Database | Senior | [Moving tenant isolation into Postgres row-level security: the owner sees every row, and pooled connections leak the tenant](days/2026-11-07.md#q3) | 把租户隔离移到 Postgres 行级安全中：表的所有者能看到所有行，连接池中的连接还会泄漏租户 |
| 4 | Middleware | Hard | [A delayed-job queue on a Redis sorted set runs jobs twice under load and loses them when a worker crashes](days/2026-11-07.md#q4) | 基于 Redis 有序集合的延迟任务队列：高负载时任务执行两次，worker 崩溃时任务丢失 |
| 5 | Java & CS | Hard | [Truncating product titles to an ad platform's limit turns emoji into '?' and some feeds are rejected](days/2026-11-07.md#q5) | 按广告平台的长度限制截断商品标题，表情符号变成了 '?'，部分 feed 还被拒绝 |

## [2026-11-06](days/2026-11-06.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [A config cache guarded by sync.RWMutex: one slow reader and a waiting writer freeze every request](days/2026-11-06.md#q1) | 用 sync.RWMutex 保护的配置缓存：一个慢读者加一个等待中的写者，让所有请求都卡住 |
| 2 | System Design | Hard | [Designing an ad click-tracking redirect service: fast redirects, bot filtering and correct attribution](days/2026-11-06.md#q2) | 设计广告点击跟踪重定向服务：快速重定向、过滤机器人流量以及正确归因 |
| 3 | Database | Senior | [An abandoned CDC replication slot fills the primary's disk overnight](days/2026-11-06.md#q3) | 一个被遗弃的 CDC 复制槽，一夜之间写满了主库的磁盘 |
| 4 | Middleware | Medium | [A Redis MULTI block half-applied a stock reservation: transactions without rollback](days/2026-11-06.md#q4) | 一个 Redis MULTI 块只执行了一半的库存预留：没有回滚的事务 |
| 5 | Java & CS | Medium | [Routing merchants to shards with Math.abs(hashCode()) % n crashes for one key, and the fix moves data](days/2026-11-06.md#q5) | 用 Math.abs(hashCode()) % n 把商家路由到分片，某个 key 会导致崩溃，而修复它又会让数据搬家 |

## [2026-11-05](days/2026-11-05.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [Webhook signatures fail at random: the signed string is built by ranging over a map](days/2026-11-05.md#q1) | webhook 签名随机失败：被签名的字符串是通过遍历 map 拼出来的 |
| 2 | System Design | Senior | [Changing the product search mapping for hundreds of millions of products without downtime](days/2026-11-05.md#q2) | 在不停机的情况下修改数亿商品的搜索映射 |
| 3 | Database | Hard | [Filtering products by a JSONB attribute is slow, even after adding a GIN index](days/2026-11-05.md#q3) | 按 JSONB 属性筛选商品很慢，即使已经加了 GIN 索引 |
| 4 | Middleware | Medium | [Deleting one merchant's cache freezes Redis for half a second: big keys, DEL and UNLINK](days/2026-11-05.md#q4) | 删除一个商家的缓存会让 Redis 卡住半秒：大 key、DEL 与 UNLINK |
| 5 | Java & CS | Hard | [A parallel stream builds the product export list, and on large catalogs products go missing or the job crashes](days/2026-11-05.md#q5) | 用并行流构建商品导出列表，在大型商品目录上会丢失商品，或者任务直接崩溃 |

## [2026-11-04](days/2026-11-04.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Senior | [One malformed webhook takes down the whole webhook service, although net/http recovers handler panics](days/2026-11-04.md#q1) | 一个格式错误的 webhook 让整个 webhook 服务崩溃，尽管 net/http 会恢复 handler 中的 panic |
| 2 | System Design | Medium | [Order alerts to merchants by email, SMS and app push: no duplicates, quiet hours, and per-channel failures](days/2026-11-04.md#q2) | 通过邮件、短信和应用推送向商家发送订单提醒：不重复、免打扰时段，以及单个渠道的故障 |
| 3 | Database | Hard | [A Postgres-backed job table: workers process the same sync job twice, or all wait behind one slow job](days/2026-11-04.md#q3) | 基于 Postgres 的作业表：worker 把同一个同步作业处理了两次，或者全都排在一个慢作业后面等待 |
| 4 | Middleware | Hard | [Merchants cannot page past result 10,000 in product search, and the "export results" button times out](days/2026-11-04.md#q4) | 商家在商品搜索中无法翻到第 10,000 条之后，而“导出结果”按钮会超时 |
| 5 | Java & CS | Medium | [Sorting products by price into a TreeSet silently drops products that share a price](days/2026-11-04.md#q5) | 把商品按价格放进 TreeSet 排序时，价格相同的商品被悄悄丢掉了 |

## [2026-11-03](days/2026-11-03.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [A catalog push reports "done" while products are still being sent, and sometimes misses the error](days/2026-11-03.md#q1) | 商品目录推送已经报告“完成”，而商品还在发送中，有时还会漏掉错误 |
| 2 | System Design | Medium | ["Export all my orders": designing large exports that do not time out or overload the database](days/2026-11-03.md#q2) | “导出我的全部订单”：设计不会超时、也不会压垮数据库的大型导出 |
| 3 | Database | Senior | [Partitioning a huge orders table by month: the unique order number breaks, and queries still scan every partition](days/2026-11-03.md#q3) | 按月对超大订单表分区：唯一订单号约束失效，而查询仍然扫描所有分区 |
| 4 | Middleware | Hard | [After a Kafka broker failover, some conversion events are duplicated and some arrive out of order](days/2026-11-03.md#q4) | Kafka broker 故障切换之后，一些转化事件重复了，一些事件乱序到达 |
| 5 | Java & CS | Medium | [Building a SKU-to-price map crashes on duplicate SKUs and on products without a price](days/2026-11-03.md#q5) | 构建 SKU 到价格的映射时，遇到重复 SKU 和没有价格的商品就会崩溃 |

## [2026-11-02](days/2026-11-02.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [Products imported from a JSON feed inherit the sale price and attributes of the product before them](days/2026-11-02.md#q1) | 从 JSON feed 导入的商品，继承了前一个商品的促销价和属性 |
| 2 | System Design | Hard | [Opening the platform to third-party apps: OAuth install flow, scopes, and revoking access](days/2026-11-02.md#q2) | 向第三方应用开放平台：OAuth 安装流程、权限范围与撤销访问 |
| 3 | Database | Medium | [The conversion-rate report shows 0% for every campaign, and one campaign crashes it with division by zero](days/2026-11-02.md#q3) | 转化率报表中每个广告活动都显示 0%，还有一个活动让报表因除以零而崩溃 |
| 4 | Middleware | Senior | [Order events on Redis Streams: a crashed worker's messages are stuck forever, and trimming loses others](days/2026-11-02.md#q4) | 基于 Redis Streams 的订单事件：崩溃的 worker 手里的消息永远卡住，而裁剪又丢掉了另一些消息 |
| 5 | Java & CS | Hard | [Moving budget between two campaigns occasionally freezes both campaigns forever](days/2026-11-02.md#q5) | 在两个广告活动之间调拨预算，偶尔会让两个活动永远卡住 |

## [2026-11-01](days/2026-11-01.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Senior | [Conversions sent to Meta after the response is written fail with "context canceled"](days/2026-11-01.md#q1) | 在响应写出之后再发往 Meta 的转化事件，以 "context canceled" 失败 |
| 2 | System Design | Hard | [Deleting a shopper's personal data across a dozen services within 30 days](days/2026-11-01.md#q2) | 在 30 天内删除一位顾客分散在十几个服务中的个人数据 |
| 3 | Database | Medium | [A unique constraint should stop duplicate ad connections, but duplicates still appear when one column is NULL](days/2026-11-01.md#q3) | 唯一约束本应阻止重复的广告连接，但当某一列为 NULL 时仍然出现了重复 |
| 4 | Middleware | Medium | [Redis memory keeps growing because keys that were supposed to expire have no TTL](days/2026-11-01.md#q4) | Redis 内存持续增长，因为本该过期的键没有 TTL |
| 5 | Java & CS | Hard | [A product-sync service runs out of memory during a Meta API slowdown, even though it uses a fixed pool of 20 threads](days/2026-11-01.md#q5) | 一个商品同步服务在 Meta API 变慢时内存耗尽，尽管它只用了一个固定 20 个线程的线程池 |

## [2026-10-31](days/2026-10-31.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [A per-merchant rate limiter lets every request through, and the mutex seems to do nothing](days/2026-10-31.md#q1) | 一个按商家限流的限流器放行了所有请求，互斥锁似乎完全不起作用 |
| 2 | System Design | Hard | [Letting merchants A/B test prices without charging shoppers a different price at checkout](days/2026-10-31.md#q2) | 让商家对价格做 A/B 测试，同时不让顾客在结账时被收取另一个价格 |
| 3 | Database | Medium | [An index exists on both columns, yet the merchant's recent-orders page is slow](days/2026-10-31.md#q3) | 两个列上都有索引，商家的近期订单页却依然很慢 |
| 4 | Middleware | Senior | [A Meta API outage stalls Google and TikTok syncs too: designing bulkheads and circuit breakers](days/2026-10-31.md#q4) | Meta API 故障连带拖住了 Google 和 TikTok 的同步：设计舱壁和熔断器 |
| 5 | Java & CS | Medium | [TikTok added a new ad status, and the whole sync job crashes on it](days/2026-10-31.md#q5) | TikTok 新增了一种广告状态，整个同步任务因此崩溃 |

## [2026-10-30](days/2026-10-30.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Senior | [A product cache evicts 99% of its entries, but the service's memory barely drops](days/2026-10-30.md#q1) | 商品缓存淘汰了 99% 的条目，但服务的内存几乎没有下降 |
| 2 | System Design | Medium | [Designing returns and refunds: partial refunds, duplicate clicks and a slow payment provider](days/2026-10-30.md#q2) | 设计退货与退款：部分退款、重复点击和缓慢的支付服务商 |
| 3 | Database | Hard | [A batch product upsert fails with "cannot affect row a second time"](days/2026-10-30.md#q3) | 批量商品 upsert 以“cannot affect row a second time”失败 |
| 4 | Middleware | Medium | [Shoppers in Germany see US prices because the CDN cached a personalized page](days/2026-10-30.md#q4) | 德国顾客看到了美国价格，因为 CDN 缓存了一个个性化页面 |
| 5 | Java & CS | Hard | [The product feed breaks only on servers configured for Germany and Turkey](days/2026-10-30.md#q5) | 商品 feed 只在配置为德国和土耳其的服务器上出错 |

## [2026-10-29](days/2026-10-29.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [Merchants cannot un-pause a campaign or set a budget to zero through the API](days/2026-10-29.md#q1) | 商家无法通过 API 恢复暂停的广告活动，也无法把预算设为零 |
| 2 | System Design | Hard | [Order numbers collide after the order service scales to many instances](days/2026-10-29.md#q2) | 订单服务扩展到多个实例之后，订单号开始冲突 |
| 3 | Database | Hard | [Invoice numbers must have no gaps, but the sequence skips numbers](days/2026-10-29.md#q3) | 发票号必须连续无间隔，但序列会跳号 |
| 4 | Middleware | Senior | [Migrating the order event pipeline from RabbitMQ to Kafka without losing or duplicating events](days/2026-10-29.md#q4) | 把订单事件管道从 RabbitMQ 迁移到 Kafka，同时不丢失、不重复事件 |
| 5 | Java & CS | Medium | [Adding to a list of SKUs throws UnsupportedOperationException after a refactor](days/2026-10-29.md#q5) | 重构之后，向一个 SKU 列表添加元素时抛出 UnsupportedOperationException |

## [2026-10-28](days/2026-10-28.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [A catalog fan-in crashes with "send on closed channel" when the consumer stops early](days/2026-10-28.md#q1) | 商品目录的扇入流程在消费者提前停止时，以“send on closed channel”崩溃 |
| 2 | System Design | Senior | [EU merchants require their data to stay in the EU: turning a single-region SaaS multi-region](days/2026-10-28.md#q2) | 欧盟商家要求数据留在欧盟境内：把单区域 SaaS 改造成多区域 |
| 3 | Database | Medium | [A US merchant's daily sales report puts evening orders on the next day](days/2026-10-28.md#q3) | 美国商家的每日销售报表把晚上的订单算到了第二天 |
| 4 | Middleware | Medium | [A merchant saves a product, but search does not find it for a moment](days/2026-10-28.md#q4) | 商家保存了商品，但搜索在片刻之内找不到它 |
| 5 | Java & CS | Hard | [@Transactional is on the method, but the data is saved without a transaction](days/2026-10-28.md#q5) | 方法上加了 @Transactional，数据却在没有事务的情况下被保存了 |

## [2026-10-27](days/2026-10-27.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [The retry logic never notices rate limits, and an empty error is not nil](days/2026-10-27.md#q1) | 重试逻辑从来识别不出限流，而一个“空”的错误却不等于 nil |
| 2 | System Design | Medium | [Showing prices in 40 currencies without losing money on exchange rates](days/2026-10-27.md#q2) | 用 40 种货币显示价格，而不在汇率上亏钱 |
| 3 | Database | Hard | [A failed CREATE UNIQUE INDEX CONCURRENTLY left an index that does nothing](days/2026-10-27.md#q3) | 一次失败的 CREATE UNIQUE INDEX CONCURRENTLY 留下了一个什么也不做的索引 |
| 4 | Middleware | Senior | [After adding Kafka partitions, product updates are applied out of order](days/2026-10-27.md#q4) | 增加 Kafka 分区之后，商品更新的应用顺序乱了 |
| 5 | Java & CS | Hard | [A product cache throws "Recursive update" for about one key in five](days/2026-10-27.md#q5) | 一个商品缓存对大约五分之一的键抛出“Recursive update” |

## [2026-10-26](days/2026-10-26.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Senior | [Every deploy drops a few hundred in-flight requests, although the service calls Shutdown](days/2026-10-26.md#q1) | 每次发布都会丢掉几百个正在处理的请求，尽管服务调用了 Shutdown |
| 2 | System Design | Hard | [Abandoned-cart emails for 50,000 stores: on time, once, and never after the order](days/2026-10-26.md#q2) | 为 50,000 家店铺发送弃购提醒邮件：准时、只发一次、下单后绝不再发 |
| 3 | Database | Medium | ["Products without ads" returns nothing, although many products have no ads](days/2026-10-26.md#q3) | “没有广告的商品”查询返回空，尽管很多商品都没有广告 |
| 4 | Middleware | Hard | [Conversion events disappear after a consumer restart, although Kafka still has them](days/2026-10-26.md#q4) | 消费者重启后，转化事件不见了，尽管 Kafka 里它们还在 |
| 5 | Java & CS | Medium | [Order dates in a report are sometimes wrong, and parsing occasionally throws odd exceptions](days/2026-10-26.md#q5) | 报表中的订单日期有时是错的，解析偶尔还会抛出奇怪的异常 |

## [2026-10-25](days/2026-10-25.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [Ad IDs from Meta come back slightly wrong after passing through a Go service](days/2026-10-25.md#q1) | Meta 的广告 ID 经过一个 Go 服务之后，回来时变得“差了一点” |
| 2 | System Design | Medium | [Feature flags and gradual rollouts for a multi-tenant SaaS](days/2026-10-25.md#q2) | 多租户 SaaS 中的功能开关与灰度发布 |
| 3 | Database | Senior | [Two admins remove each other at the same time, and a store is left with no owner](days/2026-10-25.md#q3) | 两个管理员同时把对方移除，结果店铺一个所有者都没有了 |
| 4 | Middleware | Medium | [A Redis rate limiter blocks some merchants forever](days/2026-10-25.md#q4) | 一个 Redis 限流器把部分商家永久挡在了门外 |
| 5 | Java & CS | Hard | [Parallel calls with CompletableFuture make the whole service slow, and errors arrive wrapped twice](days/2026-10-25.md#q5) | 用 CompletableFuture 并行调用后整个服务都变慢了，异常还被包了两层 |

## [2026-10-24](days/2026-10-24.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [A service crashes with "fatal error: concurrent map writes", and recover does not help](days/2026-10-24.md#q1) | 服务因“fatal error: concurrent map writes”崩溃，而且 recover 也救不了 |
| 2 | System Design | Senior | [Shared tables, schema per tenant, or database per tenant for an enterprise tier?](days/2026-10-24.md#q2) | 企业版：共享表、每租户一个 schema，还是每租户一个数据库？ |
| 3 | Database | Medium | [Product search with LIKE '%wool%' is slow even though title is indexed](days/2026-10-24.md#q3) | 用 LIKE '%wool%' 搜索商品很慢，尽管 title 已经建了索引 |
| 4 | Middleware | Hard | [Redis ran out of memory: first writes failed, then the job queue vanished](days/2026-10-24.md#q4) | Redis 内存用尽：先是写入失败，后来任务队列消失了 |
| 5 | Java & CS | Medium | [A HashSet of SKUs contains the same SKU twice](days/2026-10-24.md#q5) | 一个 SKU 的 HashSet 里出现了两个相同的 SKU |

## [2026-10-23](days/2026-10-23.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [A job that processes 5,000 feed files fails with "too many open files"](days/2026-10-23.md#q1) | 处理 5000 个 feed 文件的任务报错“too many open files” |
| 2 | System Design | Hard | [Keeping inventory in sync across the shop, Amazon and TikTok Shop without overselling](days/2026-10-23.md#q2) | 在自营店铺、Amazon 和 TikTok Shop 之间同步库存而不超卖 |
| 3 | Database | Hard | [A merchant saves a product and the page still shows the old price](days/2026-10-23.md#q3) | 商家保存了商品，页面却仍然显示旧价格 |
| 4 | Middleware | Medium | [Order events sent over Redis Pub/Sub sometimes never arrive](days/2026-10-23.md#q4) | 通过 Redis Pub/Sub 发送的订单事件有时根本收不到 |
| 5 | Java & CS | Senior | [p99 latency spikes every few minutes in a Java service with a big heap](days/2026-10-23.md#q5) | 一个大堆内存的 Java 服务每隔几分钟 p99 延迟就飙升一次 |

## [2026-10-22](days/2026-10-22.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Senior | [A feed parser keeps only a few IDs, yet memory grows by hundreds of megabytes](days/2026-10-22.md#q1) | feed 解析器只保留了几个 ID，内存却增长了几百 MB |
| 2 | System Design | Medium | [Sending webhooks to thousands of merchant endpoints that are often down](days/2026-10-22.md#q2) | 向成千上万个经常宕机的商家端点发送 webhook |
| 3 | Database | Medium | [The dashboard's order count takes seconds on large merchants](days/2026-10-22.md#q3) | 大商家的后台订单总数要好几秒才能显示出来 |
| 4 | Middleware | Hard | [Streaming responses arrive all at once, and connections drop after a minute](days/2026-10-22.md#q4) | 流式响应一下子全部到达，而且连接在一分钟后断开 |
| 5 | Java & CS | Hard | [One merchant's request runs with another merchant's ID](days/2026-10-22.md#q5) | 某个商家的请求却带着另一个商家的 ID 在运行 |

## [2026-10-21](days/2026-10-21.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [Two identical timestamps compare as different, and a cache never hits](days/2026-10-21.md#q1) | 两个相同的时间戳比较结果却不相等，缓存也永远命中不了 |
| 2 | System Design | Senior | [Sending purchase conversions to Meta and Google from the server, reliably](days/2026-10-21.md#q2) | 从服务端可靠地把购买转化回传给 Meta 和 Google |
| 3 | Database | Hard | [Duplicate platform connections despite a check before insert](days/2026-10-21.md#q3) | 插入之前做了检查，却还是出现了重复的平台连接 |
| 4 | Middleware | Hard | [After moving to Redis Cluster, multi-key commands start failing](days/2026-10-21.md#q4) | 迁移到 Redis Cluster 之后，多键命令开始报错 |
| 5 | Java & CS | Medium | [Removing items from a list in a loop throws, except when it doesn't](days/2026-10-21.md#q5) | 在循环中删除列表元素会抛异常，但有时又不会 |

## [2026-10-20](days/2026-10-20.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [A product-sync worker pool hangs forever when one item fails](days/2026-10-20.md#q1) | 商品同步的 worker 池在某个商品失败后永远卡住 |
| 2 | System Design | Medium | [Product images: uploads, resizing, CDN, and old images that never go away](days/2026-10-20.md#q2) | 商品图片：上传、缩放、CDN，以及怎么也去不掉的旧图 |
| 3 | Database | Hard | [Deep pages of the order list are slow, and rows go missing between pages](days/2026-10-20.md#q3) | 订单列表翻到后面很慢，而且翻页时会漏掉数据 |
| 4 | Middleware | Senior | [Retries at every layer turned a slow payment service into duplicate charges](days/2026-10-20.md#q4) | 每一层都在重试，把一个变慢的支付服务变成了重复扣款 |
| 5 | Java & CS | Medium | [Two equal quantities compare as different, but only above 127](days/2026-10-20.md#q5) | 两个相等的数量比较结果却不相等，而且只在超过 127 时出现 |

## [2026-10-19](days/2026-10-19.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [The API returns null instead of an empty list, and a client crashes](days/2026-10-19.md#q1) | API 返回了 null 而不是空列表，导致客户端崩溃 |
| 2 | System Design | Hard | [An audit log that merchants and auditors can trust](days/2026-10-19.md#q2) | 一份商家和审计人员都能信任的审计日志 |
| 3 | Database | Senior | [Tables keep growing and queries slow down, although rows are being deleted](days/2026-10-19.md#q3) | 表在不断膨胀、查询越来越慢，尽管一直在删除数据 |
| 4 | Middleware | Medium | [Searching for SKU "AB-1234" in Elasticsearch finds nothing](days/2026-10-19.md#q4) | 在 Elasticsearch 里搜索 SKU“AB-1234”什么也搜不到 |
| 5 | Java & CS | Hard | [Money bugs: 1.0 is not equal to 1.00, and split discounts lose a cent](days/2026-10-19.md#q5) | 金额 bug：1.0 不等于 1.00，拆分的折扣少了一分钱 |

## [2026-10-15](days/2026-10-15.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [Every product fails validation after a small refactor](days/2026-10-15.md#q1) | 一次小重构之后，所有商品都校验失败 |
| 2 | System Design | Hard | [A flash sale that oversells and melts the order database](days/2026-10-15.md#q2) | 一场既超卖、又压垮订单库的限时抢购 |
| 3 | Database | Senior | [Orders that vanished after a database failover](days/2026-10-15.md#q3) | 数据库故障切换之后消失的订单 |
| 4 | Middleware | Hard | [One merchant's events back up a whole Kafka partition](days/2026-10-15.md#q4) | 一个商家的事件堵住了整个 Kafka 分区 |
| 5 | Java & CS | Medium | [A scheduled token refresh that silently stopped](days/2026-10-15.md#q5) | 悄无声息停止运行的定时刷新 token 任务 |

## [2026-10-14](days/2026-10-14.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [Log lines that come out scrambled under load](days/2026-10-14.md#q1) | 高负载下日志行变得错乱 |
| 2 | System Design | Senior | [Sharing ad-platform API quotas across many services](days/2026-10-14.md#q2) | 在多个服务之间共享广告平台的 API 配额 |
| 3 | Database | Medium | [Inserts slow down as a table keyed by random UUIDs grows](days/2026-10-14.md#q3) | 以随机 UUID 为主键的表越大，插入越慢 |
| 4 | Middleware | Medium | [A hot cache key expires and the database falls over](days/2026-10-14.md#q4) | 一个热点缓存 key 过期，数据库随之被打垮 |
| 5 | Java & CS | Hard | [A request thread that spins forever waiting for a cache](days/2026-10-14.md#q5) | 一个永远在空转等待缓存的请求线程 |

## [2026-10-13](days/2026-10-13.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [An API client that slowly runs out of ports](days/2026-10-13.md#q1) | 一个慢慢耗尽端口的 API 客户端 |
| 2 | System Design | Medium | [Importing a 2 GB product CSV without timeouts or half-done imports](days/2026-10-13.md#q2) | 导入 2 GB 的商品 CSV，不超时，也不留下半成品 |
| 3 | Database | Hard | [Fast in staging, full table scans in production](days/2026-10-13.md#q3) | 预发环境很快，生产环境却在全表扫描 |
| 4 | Middleware | Hard | [Lost and looping messages in a RabbitMQ pipeline](days/2026-10-13.md#q4) | RabbitMQ 链路中丢失和死循环的消息 |
| 5 | Java & CS | Senior | [Virtual threads made the service slower](days/2026-10-13.md#q5) | 启用虚拟线程后，服务反而变慢了 |

## [2026-10-12](days/2026-10-12.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Senior | [Latency spikes and OOMKills in a containerized Go service](days/2026-10-12.md#q1) | 容器化 Go 服务的延迟尖刺与 OOMKill |
| 2 | System Design | Hard | [Running jobs at 9:00 local time for 200,000 merchants](days/2026-10-12.md#q2) | 为 20 万个商家在当地时间 9:00 运行任务 |
| 3 | Database | Hard | [Coupons redeemed more often than their limit allows](days/2026-10-12.md#q3) | 优惠券被核销的次数超过了上限 |
| 4 | Middleware | Medium | [One Redis command that stalls every client](days/2026-10-12.md#q4) | 一条让所有客户端都卡住的 Redis 命令 |
| 5 | Java & CS | Medium | [Two new order lines go in, only one comes out](days/2026-10-12.md#q5) | 加了两行新的订单明细，只保存下来一行 |

## [2026-10-11](days/2026-10-11.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [A price service that freezes hours after every deploy](days/2026-10-11.md#q1) | 每次发布几小时后就卡死的价格服务 |
| 2 | System Design | Medium | [A webhook endpoint that loses, duplicates and reorders events](days/2026-10-11.md#q2) | 一个会丢失、重复、乱序处理事件的 webhook 接口 |
| 3 | Database | Medium | [An export job that gets slower and still misses rows](days/2026-10-11.md#q3) | 越跑越慢、还会漏数据的导出任务 |
| 4 | Middleware | Senior | [A Redis lock that lets a settlement job run twice](days/2026-10-11.md#q4) | 一个会让结算任务跑两次的 Redis 锁 |
| 5 | Java & CS | Hard | [REQUIRES_NEW and a connection pool that deadlocks under load](days/2026-10-11.md#q5) | REQUIRES_NEW 让连接池在高负载下死锁 |

## [2026-10-10](days/2026-10-10.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Medium | [Batches that silently overwrite each other](days/2026-10-10.md#q1) | 批次之间悄悄互相覆盖 |
| 2 | System Design | Hard | [Refreshing OAuth tokens from many pods without logging merchants out](days/2026-10-10.md#q2) | 多个 pod 并发刷新 OAuth token，又不让商家掉线 |
| 3 | Database | Senior | [Why an "online" ALTER TABLE took the site down](days/2026-10-10.md#q3) | 为什么一个“在线” ALTER TABLE 让整个网站宕机 |
| 4 | Middleware | Hard | [A slow downstream API turns into a Kafka rebalance storm](days/2026-10-10.md#q4) | 下游 API 变慢演变成 Kafka 重平衡风暴 |
| 5 | Java & CS | Medium | [One merchant sees another merchant's report](days/2026-10-10.md#q5) | 一个商家看到了另一个商家的报表 |

## [2026-10-09](days/2026-10-09.md)

| # | Topic | Difficulty | Title | 标题 |
| --- | --- | --- | --- | --- |
| 1 | Go | Hard | [Goroutine leak behind a timed-out downstream call](days/2026-10-09.md#q1) | 下游调用超时后产生的 goroutine 泄漏 |
| 2 | System Design | Senior | [Syncing a product catalog to four ad platforms reliably](days/2026-10-09.md#q2) | 把商品目录可靠地同步到四个广告平台 |
| 3 | Database | Hard | [Deadlocks from locking a row that does not exist](days/2026-10-09.md#q3) | 锁定一行不存在的数据导致的死锁 |
| 4 | Java & CS | Medium | [ConcurrentHashMap and the check-then-act trap](days/2026-10-09.md#q4) | ConcurrentHashMap 与先检查后执行的陷阱 |
| 5 | Middleware | Medium | [How cache-aside still serves stale data](days/2026-10-09.md#q5) | 旁路缓存（cache-aside）为什么仍会读到脏数据 |

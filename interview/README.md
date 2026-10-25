# 每日面试题 · Daily Interview Questions

后端面试题（Go 为主、Java 为辅），中英双语，每天 5 道，覆盖 Go、System Design、Database、Middleware、Java & CS。
每题包含题目、参考答案（默认折叠）、追问和词汇表。

共 15 天、75 道题。

- 数据：`interview/data/YYYY-MM-DD.json`（与 Interview Drill 页面的 `days` 集合结构一致）
- 重新生成 Markdown：在仓库根目录执行 `go run ./interview/render`
- Agent 面试速成：[agent/README.md](agent/README.md)

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

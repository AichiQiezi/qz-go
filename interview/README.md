# 每日面试题 · Daily Interview Questions

后端面试题（Go 为主、Java 为辅），中英双语，每天 5 道，覆盖 Go、System Design、Database、Middleware、Java & CS。
每题包含题目、参考答案（默认折叠）、追问和词汇表。

共 7 天、35 道题。

- 数据：`interview/data/YYYY-MM-DD.json`（与 Interview Drill 页面的 `days` 集合结构一致）
- 重新生成 Markdown：在仓库根目录执行 `go run ./interview/render`

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

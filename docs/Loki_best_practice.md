# Loki 建议

[TOC]

## 1. 尽量使用静态标签

使用静态标签可以在日志时的开销更小。通常日志在发送到 Loki 之前，在注入 label 时，常见的推荐静态标签包含：

- 物理机：kubernetes/hosts
- 应用名：kubernetes/labels/app_kubernetes_io/name
- 组件名：kubernetes/labels/name
- 命名空间：kubernetes/namespace
- 其他 kubernetes/label/\* 的静态标签，如环境、版本等信息

## 2. 谨慎使用动态标签

过多的标签组合会造成大量的流，它会让 Loki 存储大量的索引和小块的对象文件。这些都会显著消耗 Loki 的查询性能。为避免这些问题，在你知道需要之前不要添加标签！`loki的优势在于并行查询`，使用过滤器表达式`( lable = "text", |~ "regex", …)`来查询日志会更有效，并且速度也很快。

**那么，我该什么时候添加标签？**

`chunk_target_size`默认为 1MB，loki 将以 1MB 的压缩后大小来切割日志块，大约等于 5MB 的原始日志文件（根据你配置的压缩级别来决定）。如果在`max_chunk_age`时间内，你的日志流`足以生成一个或者多个压缩块，那么你可以考虑添加标签，将日志流拆得更细一点`。从 Loki 1.4.0 开始，有一个指标可以帮助我们了解日志块刷新的情况

```
sum by (reason) (rate(loki_ingester_chunks_flushed_total{cluster="dev"}[1m]))
```

## 3. 有界的标签值范围

不管怎样，到最后如果你不得不采用动态标签的话，那你也要`注意控制标签的范围和value值的长度`。举个例子，如果你想将 nginx 的访问日志提取一些字段后存储到 loki，

```
{"@timestamp":"2020-09-30T12:16:07+08:00","@source":"172.16.1.1","hostname":"node1","ip":"-","client":"172.16.2.1","request_method":"GET","scheme":"https","domain":"xxx.com","referer":"-","request":"/api/v1/asset/asset?page_size=-1&group=23","args":"page_size=-1&group=23","size":975,"status": 200,"responsetime":0.065,"upstreamtime":"0.064","upstreamaddr":"172.16.3.1:8080","http_user_agent":"python-requests/2.22.0","https":"on"}
```

这里面`@source`代表客户端源地址，由于源地址是公网地址，那么在建立 loki 标签时它的值就是个无界的。 再比如这里面`@request`代表请求 URL。可能存在某些请求参数过长，loki 的标签值也会过大。如果再将两者相乘，那么这个标签的规模是无法接受的。

以上这种情况是比较属于典型无界的动态标签值，在 loki 里面我们用`Cardinality`来表述它，`Cardinality值越高，loki的查询效率越低。`。Loki 社区给出动态标签的范围应尽量`控制在10以内`。

## 4. 客户端应用的动态标签

Loki 的几个客户端（Promtail、Fluentd，Fluent Bit，Docker 插件等）都带有配置标签来创建日志流的方法。我们有时需要在 loki 里面去排查哪些应用使用了动态标签，这时候我们可以用 logcli 工具来辅助我们。在 Loki1.6.0 及更高版本中，`logcli series`命令添加了`--analyze-labels`参数专门用于调试高`cardinality`的标签。例如：

```
$ logcli series --analyze-labels '{app="nginx"}'

Total Streams:  25017
Unique Labels:  8

Label Name  Unique Values  Found In Streams
requestId   24653          24979
logStream   1194           25016
logGroup    140            25016
accountId   13             25016
logger      1              25017
source      1              25016
transport   1              25017
format      1              25017
```

可以看到这里面`requestId`这个标签就有 24653 个值，这是非常不好的。我们`应该将requestId从label里面去删除`,通过这种方式查询

```
{app="nginx"} |= "requestId=1234567"
```

## 5. 配置缓存

### cache_config

**cache_config**是 Loki 的缓存配置区块，当前 Loki 1.6 支持的缓存主要是`in-memory`、`memcached`和`redis`。这三种缓存类型各有自己的场景需求，如果你的 Loki 是 AllinOne 部署话，三选一都可以。如果你的 Loki 是分布式的架构，那么可以选着 redis 作为主要缓存服务

- 内存缓存（in-memory）

```python
# 启动内存缓存
enable_fifocache: <boolean>
# 缺省缓存过期时间
default_validity: <duration>
fifocache:
# 缓存最大内存占用
  max_size_bytes: <int> = 10000
# 字段限制，0无限制
  max_size_items: <int> | default = 0
  validity: <duration>
```

- Memcached

```python
# 启用Memcached后的后台配置
background：
  #回写memcacehd的goroutines数目
  writeback_goroutines: <int> | default = 10]
  writeback_buffer: <int> = 10000
memcached:
  # 过期时间
  expiration: <duration>
  batch_size: <int>
  # 并发限制
  parallelism: <int> | default = 100
memcached_client：
  host: <string>
  service: <string> | default = "memcached"
  # 超时时间
  timeout: <duration> | default = 100ms
  # 最大空闲连接数
  max_idle_conns: <int> | default = 100
  update_interval: <duration> | default = 1m
  consistent_hash: <bool>
```

- Redis

```python
redis
  endpoint: <string>
  timeout: <duration> | default = 100ms
  # 过期时间
  expiration: <duration> | default = 0s
  max_idle_conns: <int> | default = 80
  max_active_conns: <int> | default = 0
  password: <string>
  enable_tls: <boolean> | default = false
```

### 缓存作用域

- 查询结果缓存

`queryrange_config`里面定义了 Loki 查询时关于缓存和切块的配置，缓存的相关配置如下

```python
# 查询缓存开关，默认关闭
cache_results: <boolean> | default = false
results_cache：
# 缓存配置块
  cache: <cache_config>
```

- 日志索引缓存

`index_queries_cache_config`定义 Loki 的索引缓存，大部分情况下可以等同于日志 label 的查询缓存

```python
storage_config:
  #索引缓存有效时间
  index_cache_validity: <duration> | default = 5m
  index_queries_cache_config: <cache_config>
```

> 注意！index_cache_validity 的时间要小于 ingester_config.chunk_idle_period 配置的时间。 大意是日志的入到 Loki 后，缓存的日志索引在原始日志 flush 进存储前都为有效的，以保证查询的缓存索引是正确的。

- 原始日志缓存

`chunk_store_config`定义 Loki 将原始日志写入存储阶段的配置，这里引入缓存其主要目的为增大 Loki 日志接收日志的吞吐量。

```python
chunk_store_config：
  # 日志写入存储前的缓存配置
  chunk_cache_config: <cache_config>
  # 删除重复写入的日志缓存配置
  write_dedupe_cache_config: <cache_config>
```

> 注意！对于日志 chunk 引入缓存，我们务必要将数据持久化，如果采用 in-memory 或者 MemoryCache 存在服务异常掉日志内容的风险，如果采用 redis 则最好把数据持久化打开

- 限制条件

`limits_config`定义 Loki 中提取日志配置全局和按租户限制。

```python
limits_config:
  # 单次查询限制
  max_entries_limit_per_query: <int> | default = 5000 ]
```

### 举个例子

以 redis 作为 Loki 缓存来举个例子，让大家更直观看到关于缓存在全局配置里面的分布。

```
...
frontend:
  compress_responses: true

query_range:
  split_queries_by_interval: 24h
  results_cache:
    cache:
      redis:
        endpoint: redis:6379
        expiration: 1h
  cache_results: true

storage_config:
  index_queries_cache_config:
    redis:
      endpoint: redis:6379
      expiration: 1h

chunk_store_config:
  chunk_cache_config:
    redis:
      endpoint: redis:6379
      expiration: 1h
  write_dedupe_cache_config:
    redis:
      endpoint: redis:6379
      expiration: 1h
...
```

## 6. 日志的时间必须顺序递增

对于一个日志流里面出现时间戳早于该流收到的最新日志，那么这条日志将被删除

```
{job=”syslog”} 00:00:00 i’m a syslog!
{job=”syslog”} 00:00:02 i’m a syslog!
{job=”syslog”} 00:00:01 i’m a syslog! \\这条日志会被删除
```

如果你的服务是分布式跑在多个节点上，而且存在时间差的话，那你只有为这类日志添加新的标签来存储了

```
{job=”syslog”, instance=”host1”} 00:00:00 i’m a syslog!  \\新日志流1
{job=”syslog”, instance=”host1”} 00:00:02 i’m a syslog!
{job=”syslog”, instance=”host2”} 00:00:01 i’m a syslog!  \\新日志流2
{job=”syslog”, instance=”host1”} 00:00:03 i’m a syslog!  \\在日志流1里时间有序
{job=”syslog”, instance=”host2”} 00:00:02 i’m a syslog!  \\在日志流2里时间有序
```

这个没啥好说的，建议日志采集时按照`客户端的时间为每条日志添加时间戳`。如果你的时间戳是从应用日志里面提取出来，并且出现时间乱序的话，那还是请你`先解决应用的问题`

## 7. 使用 chunk_target_size 参数

上文说到`chunk_target_size`可以有效的将日志流压缩到一个合理的空间大小，Loki 中每个日志流都包含一个块。如果我们将日志文件分解成更多的流，内存中存储的块就越多，在被刷新到磁盘之前，理论上来说都有丢日志的风险。那么这个时候就需要组合`max_chunk_age`默认 1h 和`chunk_idle_period`默认 30m，来控制日志刷新的超时时间。

## 8. 使用-print-config-stderr 或-log-config-reverse-order 参数

从 1.6.0 版开始，Loki 和 Promtail 支持这类参数，当启动时，loki 会`把整个配置信息打印到stderr或日志文件中。`，这样我们可以快速看到整个 Loki 配置，便于调试。

当这个参数`-log-config-reverse-order`启用时，我们在 grafna 上查询 loki 时`将以顺序的方式查看日志`，这个可以让我们更加方便一点。

## 9. 使用 query-frontend

query-frontend 可以有效的将日志查询拆分成多个小查询分发给 querier 去并发执行。这件极大的提高 loki 的查询效率，理论上来说你可以`扩容上百个querier去并发处理GB或者TB级别的日志`，不过前提是你的查询客户端能够容得下这些日志。

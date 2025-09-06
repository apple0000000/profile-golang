高并发缓存 Web 服务
这是一个基于 Golang 的高并发缓存 Web 服务演示项目，实现了读写分离和二级缓存架构。

项目概述
本项目展示了一个高并发缓存系统的实现，采用读写分离架构，结合内存缓存和 Redis 作为二级缓存，通过 Kafka 实现数据同步。

架构设计
核心组件
读服务器: 处理数据读取请求，使用内存缓存和 Redis
写服务器: 处理数据写入请求，更新 Redis 并发送 Kafka 消息
Redis: 作为二级缓存，存储全量数据
Kafka: 用于读服务器之间的数据同步

服务链路
客户端 -> [写服务器] -> Redis & Kafka -> [读服务器] -> 内存缓存 -> 返回数据

设计优势
高并发: 通过内存缓存提供低延迟读取
数据一致性: 通过 Kafka 保证各读服务器缓存的一致性
可扩展性: 无状态设计，支持水平扩展
资源优化: 内存只存储热点数据，Redis 存储全量数据

功能特性
读服务器功能
HTTP 接口批量读取数据
二级缓存查询（内存 → Redis）
Kafka 数据同步消费
定期清理过期缓存
写服务器功能
HTTP 接口批量更新数据
同步更新 Redis
异步发送 Kafka 消息通知读服务器


启动服务
启动读服务器：

go run main.go -mode=reader -port=8080
启动写服务器：
go run main.go -mode=writer -port=8081

API 使用说明
写服务器接口
批量写入数据

http
POST /write
Content-Type: application/json
{
  "items": [
    {"key": "key1", "value": "value1"},
    {"key": "key2", "value": "value2"}
  ]
}

读服务器接口
批量读取数据
http
GET /read?keys=key1,key2,key3

响应：
json
{
  "items": [
    {"key": "key1", "value": "value1", "source": "memory"},
    {"key": "key2", "value": "value2", "source": "redis"},
    {"key": "key3", "value": null, "exists": false}
  ]
}



High-Concurrency Caching Web Service
This is a demonstration project written in Go that shows how to build a high-concurrency caching web service with read/write separation and a two-level cache architecture.

Project Overview
The project illustrates a high-concurrency caching system that adopts a read/write split design.
It combines in-memory caching with Redis as the second-level cache and uses Kafka for data synchronization.

Architecture Design
Core Components
Read servers: handle read requests, query both in-memory cache and Redis
Write servers: handle write requests, update Redis and publish Kafka messages
Redis: acts as the second-level cache holding the full data set
Kafka: propagates changes so every read server stays consistent

Request Flow
Client → [Write Server] → Redis & Kafka → [Read Server] → In-Memory Cache → Response

Design Benefits
High concurrency: ultra-low latency reads served from memory
Consistency: Kafka guarantees that all read-server caches stay in sync
Scalability: stateless design allows horizontal scaling
Resource efficiency: memory keeps only hot data; Redis keeps the complete data set

Features
Read-server capabilities
HTTP endpoint for batch reads
Two-level lookup (memory → Redis)
Kafka consumer for invalidation/updates
Periodic eviction of expired entries
Write-server capabilities
HTTP endpoint for batch writes
Synchronous Redis update
Asynchronous Kafka message to notify read servers

Starting the Services
Start a read server:
go run main.go -mode=reader -port=8080

Start a write server:
go run main.go -mode=writer -port=8081

API Usage
Write-server endpoint
Batch write
POST /write
Content-Type: application/json
{
    "items": [
    {"key": "key1", "value": "value1"},
    {"key": "key2", "value": "value2"}
    ]
}

Read-server endpoint
Batch read
GET /read?keys=key1,key2,key3
Response:
{
    "items": [
    {"key": "key1", "value": "value1", "source": "memory"},
    {"key": "key2", "value": "value2", "source": "redis"},
    {"key": "key3", "value": null, "exists": false}
    ]
}
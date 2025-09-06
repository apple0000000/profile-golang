# High-Concurrency Caching Web Service

这是一个基于 Golang 的高并发缓存 Web 服务演示项目，实现了读写分离和二级缓存架构。

## 项目概述

本项目展示了一个高并发缓存系统的实现，采用读写分离架构，结合内存缓存和 Redis 作为二级缓存，通过 Kafka 实现数据同步。

## 架构设计

### 核心组件

- **读服务器**: 处理数据读取请求，使用内存缓存和 Redis
- **写服务器**: 处理数据写入请求，更新 Redis 并发送 Kafka 消息
- **Redis**: 作为二级缓存，存储全量数据
- **Kafka**: 用于读服务器之间的数据同步

### 服务链路

客户端 -> [写服务器] -> Redis & Kafka -> [读服务器] -> 内存缓存 -> 返回数据

### 设计优势

- **高并发**: 通过内存缓存提供低延迟读取
- **数据一致性**: 通过 Kafka 保证各读服务器缓存的一致性
- **可扩展性**: 无状态设计，支持水平扩展
- **资源优化**: 内存只存储热点数据，Redis 存储全量数据

## 功能特性

### 读服务器功能

- HTTP 接口批量读取数据
- 二级缓存查询（内存 → Redis）
- Kafka 数据同步消费
- 定期清理过期缓存

### 写服务器功能

- HTTP 接口批量更新数据
- 同步更新 Redis
- 异步发送 Kafka 消息通知读服务器

## 启动服务

启动读服务器：
go run main.go -mode=reader -port=8080

启动写服务器：
go run main.go -mode=writer -port=8081




This is a demonstration project of a high-concurrency caching web service based on Golang, implementing a read-write separation architecture with a two-level cache.

## Project Overview

This project demonstrates the implementation of a high-concurrency caching system using a read-write separation architecture, combining in-memory cache and Redis as a secondary cache, with data synchronization achieved through Kafka.

## Architecture Design

### Core Components

- **Read Server**: Handles data read requests, using in-memory cache and Redis
- **Write Server**: Handles data write requests, updates Redis, and sends Kafka messages
- **Redis**: Serves as secondary cache, storing complete data
- **Kafka**: Used for data synchronization between read servers

### Service Flow

Client -> [Write Server] -> Redis & Kafka -> [Read Server] -> In-Memory Cache -> Return Data

### Design Advantages

- **High Concurrency**: Provides low-latency reads through in-memory cache
- **Data Consistency**: Ensures consistency across read server caches through Kafka
- **Scalability**: Stateless design supports horizontal scaling
- **Resource Optimization**: Memory stores only hot data, Redis stores complete data

## Feature Set

### Read Server Features

- HTTP interface for batch data reading
- Two-level cache query (Memory → Redis)
- Kafka data synchronization consumption
- Regular cleanup of expired cache

### Write Server Features

- HTTP interface for batch data updates
- Synchronous Redis updates
- Asynchronous Kafka message sending to notify read servers

## Service Startup

Start read server:
go run main.go -mode=reader -port=8080

Start write server:
go run main.go -mode=writer -port=8081
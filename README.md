# Work-Flow

Work-Flow 是一个基于 Go 语言开发的工作流引擎系统，提供了灵活的工作流定义、执行和监控功能。

## 功能特性

- 工作流引擎：支持复杂工作流的定义和执行
- 数据源管理：支持多种数据源（MySQL、Oracle等）的连接和管理
- 异步任务处理：基于 Asynq 的异步任务队列
- 分布式任务调度：使用 DCron 实现分布式任务调度
- 监控和追踪：集成 OpenObserve 实现系统监控和日志追踪
- API 网关：基于 go-zero 框架的 RESTful API 服务
- 规则引擎：支持自定义规则的定义和执行
- 链式处理：支持工作流节点的链式处理

## 技术栈

- 核心框架：go-zero
- 数据库：MySQL、Oracle
- 缓存：Redis
- 异步任务：Asynq
- 分布式调度：DCron
- 监控：OpenObserve
- 工作流引擎：自定义实现
- 规则引擎：自定义实现

## 系统要求

- Go 1.24.0 或更高版本
- Redis 服务器
- MySQL 数据库
- 支持的操作系统：Linux、macOS、Windows


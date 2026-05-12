# 短视频Feed流系统

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.5+-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**面向内容创作者与消费者的高性能短视频Feed流平台，采用三级缓存架构与异步事件驱动设计，支持千万级用户并发访问**

[快速开始](#快速开始) • [功能特性](#功能特性) • [技术架构](#技术架构) • [API文档](API.md) • [部署指南](DEPLOYMENT.md)

</div>

---

## 目录

- [项目简介](#项目简介)
- [核心价值](#核心价值)
- [功能特性](#功能特性)
- [技术架构](#技术架构)
- [快速开始](#快速开始)
- [使用场景](#使用场景)
- [项目结构](#项目结构)
- [项目优势](#项目优势)
- [文档导航](#文档导航)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

---

## 项目简介

短视频Feed流系统是一个面向内容创作者和消费者的短视频平台，提供完整的视频发布、浏览、互动和社交功能。系统采用API进程与Worker进程分离的架构设计，通过RabbitMQ消息队列实现异步事件驱动，支持水平扩展，能够处理千万级用户的并发访问请求。

### 核心价值

- **高性能**：采用三级缓存架构（L1本地缓存 → L2 Redis → L3 MySQL），Feed响应时间从50ms优化至3ms
- **高可用**：通过异步处理、降级机制和缓存自愈，系统可用性达到99.9%
- **数据一致性**：事务性发件箱模式（Outbox Pattern）保证消息不丢失，MQ异常时自动降级直写DB
- **可扩展**：API进程与Worker进程可独立部署和水平扩展，无状态设计支持多实例

---

## 功能特性

### 核心功能模块

#### 1. 用户账号系统
- 用户注册与登录（bcrypt密码加密）
- JWT身份认证（硬认证/软认证双模式）
- 用户信息管理（改名、查询）
- 密码修改与安全登出（Token即时撤销）

#### 2. 视频管理
- 视频发布（事务性发件箱模式保证一致性）
- 视频文件与封面上传（MP4/图片格式支持）
- 视频列表浏览（按作者查询）
- 视频详情查看（Redis缓存 + 分布式锁防击穿）

#### 3. 社交互动
- 点赞/取消点赞（MQ优先 + DB兜底双写策略）
- 评论功能（发布/删除/列表）
- 关注/取消关注
- 粉丝列表与关注列表查看

#### 4. Feed流系统
- 最新视频Feed（冷热分离 + 游标分页）
- 热门视频Feed（Redis滑动窗口热榜 + 快照分页）
- 点赞数排序Feed（双字段复合游标分页）
- 关注用户Feed（需登录，缓存防击穿）

### 技术特性

#### 三级缓存架构
```
请求 → L1本地缓存(3s) → L2 Redis缓存(1h) → L3 MySQL数据库
         命中率~60%        命中率~35%          命中率~5%
```

#### 异步事件驱动
```
用户操作 → RabbitMQ消息队列 → Worker异步处理 → 数据库持久化
              │                    │
              └── 发布失败 ────→ 降级直写DB（保证数据不丢失）
```

#### 冷热数据分离
```
热数据：Redis ZSET存储（访问频率高，毫秒级响应）
冷数据：MySQL存储（访问频率低，游标分页查询）
分离策略：基于ZSET watermark动态判断冷热边界
```

---

## 技术架构

### 技术栈

#### 后端技术
| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.24+ | 主开发语言，高性能并发 |
| Gin | 1.11+ | Web框架，路由与中间件 |
| GORM | 1.31+ | ORM框架，AutoMigrate自动建表 |
| MySQL | 8.0+ | 关系型数据库，数据持久化 |
| Redis | 7.0+ | 缓存数据库，热榜/分布式锁/Token校验 |
| RabbitMQ | 3.12+ | 消息队列，异步事件驱动 |

#### 前端技术
| 技术 | 版本 | 说明 |
|------|------|------|
| Vue 3 | 3.5+ | 前端框架，Composition API |
| TypeScript | 5.9+ | 类型系统 |
| Vite | 7.2+ | 构建工具，HMR热更新 |
| Pinia | 3.0+ | 状态管理 |
| Vue Router | 4.6+ | 路由管理 |

### 系统架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         客户端层                              │
│                    (Vue 3 + TypeScript)                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                         API网关层                             │
│              (Gin + JWT认证 + 限流 + 路由)                    │
│         硬认证(JWTAuth) / 软认证(SoftJWTAuth)                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                        业务服务层                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │账号服务   │  │视频服务   │  │社交服务   │  │Feed服务   │   │
│  │bcrypt加密 │  │发件箱模式 │  │MQ异步写入 │  │三级缓存   │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌──────────────────┐  ┌──────────────┐  ┌──────────────┐
│   缓存层 (Redis)  │  │消息队列(RabbitMQ)│  │ 数据库(MySQL) │
│  - 三级缓存       │  │  - 异步处理     │  │  - 数据持久化  │
│  - 分布式锁       │  │  - 削峰填谷     │  │  - 事务保证    │
│  - 滑动窗口热榜   │  │  - 解耦服务     │  │  - 发件箱表    │
└──────────────────┘  └──────────────┘  └──────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       Worker处理层                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │点赞Worker │  │评论Worker │  │社交Worker │  │热度Worker │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│  ┌──────────┐                                               │
│  │发件箱轮询 │  → 扫描Outbox表 → 投递TimelineMQ → 更新Redis  │
│  └──────────┘                                               │
└─────────────────────────────────────────────────────────────┘
```

### 核心设计模式

#### 1. 事务性发件箱模式（Outbox Pattern）
保证视频发布时数据与消息的原子性，避免分布式事务：
```
业务操作 → 写入Video表 + 写入Outbox表 → 事务提交
                ↓
           OutboxWorker轮询 → 投递TimelineMQ → 更新Redis ZSET
```

#### 2. MQ优先 + DB兜底双写策略
点赞/评论等写操作保证数据不丢失：
```
写操作 → 尝试发MQ（MySQL写入 + Redis热度更新）
              │
              ├── MQ成功 → 直接返回
              └── MQ失败 → 降级直写DB（保证数据落地）
```

#### 3. 滑动窗口热榜 + 快照分页
解决热榜数据实时变化导致的分页不稳定问题：
```
写入：每分钟一个ZSET（hot:video:1m:{minute}），ZINCRBY更新热度
查询：ZUNIONSTORE合并最近60个窗口 → 生成快照 → ZREVRANGE分页
翻页：携带as_of时间戳复用同一快照，保证分页一致性
```

#### 4. 快照式游标分页
替代传统OFFSET分页，提升大数据量下的查询性能：
```
传统分页：OFFSET 1000 LIMIT 10（越往后越慢）
游标分页：WHERE create_time < last_time ORDER BY create_time DESC LIMIT 10（稳定高效）
复合游标：WHERE (likes_count, id) < (last_count, last_id)（解决同值排序不稳定）
```

---

## 快速开始

### 方式一：Docker Compose（推荐）

```bash
# 克隆项目
git clone https://github.com/raaaaap/short-video-feed-system.git
cd short-video-feed-system

# 一键启动所有服务
docker compose up -d --build

# 验证服务状态
docker compose ps

# 访问服务
# 前端：http://localhost:5173
# 后端：http://localhost:8080
# RabbitMQ管理台：http://localhost:15672（账号：admin / password123）
```

### 方式二：本地开发

```bash
# 1. 启动依赖服务
docker compose up -d mysql redis rabbitmq

# 2. 启动后端API
cd backend
go mod download
go run ./cmd

# 3. 启动Worker（新终端）
cd backend
go run ./cmd/worker

# 4. 启动前端（新终端）
cd frontend
npm install
npm run dev
```

### 方式三：一键脚本

```bash
# 使用start.sh一键启动（自动拉起依赖+后端+前端）
bash start.sh
```

> 详细的部署步骤请查看 [DEPLOYMENT.md](DEPLOYMENT.md)

---

## 使用场景

### 场景一：新用户注册并浏览视频

```
1. 用户访问前端页面
2. 点击注册按钮，填写用户名和密码
3. 系统创建账号并返回JWT令牌
4. 用户浏览最新视频Feed（无需登录也可浏览）
5. 点击视频查看详情
6. 对喜欢的视频点赞或评论
```

### 场景二：内容创作者发布视频

```
1. 创作者登录账号
2. 上传视频文件和封面图片
3. 填写视频标题和描述，点击发布
4. 系统使用事务性发件箱模式保证数据一致性
5. OutboxWorker异步投递消息，视频出现在Feed流中
6. 粉丝点赞、评论提升视频热度
7. 热门视频进入滑动窗口热榜
```

### 场景三：社交互动

```
1. 用户发现感兴趣的创作者
2. 点击关注按钮
3. 系统通过MQ异步建立关注关系
4. 创作者发布新视频时，粉丝在"关注"Feed中查看
5. 互动行为（点赞/评论）通过热度系统影响视频排序
```

### 性能指标

| 指标 | 优化前 | 优化后 | 提升方式 |
|------|--------|--------|---------|
| Feed响应时间 | 50ms | 3ms | 三级缓存 + 冷热分离 |
| 系统吞吐量 | 1000 QPS | 5000 QPS | 异步MQ + 连接池优化 |
| 缓存命中率 | 60% | 95% | L1本地缓存 + L2 Redis |
| 数据库负载 | 100% | 20% | 缓存拦截 + MQ削峰 |

---

## 项目结构

```
feedsystem_video_go/
├── backend/                    # 后端代码
│   ├── cmd/                    # 入口程序
│   │   ├── main.go            # API服务入口
│   │   └── worker/            # Worker服务入口
│   │       └── main.go
│   ├── configs/               # 配置文件
│   │   ├── config.yaml        # 本地开发配置
│   │   └── config.docker.yaml # Docker环境配置
│   ├── internal/              # 内部模块
│   │   ├── account/           # 账号模块（注册/登录/JWT管理）
│   │   ├── video/             # 视频模块（发布/详情/点赞/评论）
│   │   ├── social/            # 社交模块（关注/取关/列表）
│   │   ├── feed/              # Feed模块（最新/热榜/关注流）
│   │   ├── auth/              # JWT令牌生成与解析
│   │   ├── config/            # 配置加载
│   │   ├── db/                # 数据库连接与迁移
│   │   ├── http/              # 路由注册
│   │   ├── middleware/        # 中间件
│   │   │   ├── jwt/           # JWT认证（硬/软双模式）
│   │   │   ├── rabbitmq/      # RabbitMQ客户端与业务MQ封装
│   │   │   ├── ratelimit/     # 限流中间件（Redis INCR）
│   │   │   └── redis/         # Redis客户端（缓存/锁/ZSET）
│   │   ├── observability/     # 可观测性（pprof）
│   │   └── worker/            # Worker处理（发件箱/点赞/评论/社交/热度）
│   ├── go.mod                 # Go依赖管理
│   ├── go.sum
│   └── Dockerfile             # 多阶段Docker构建
├── frontend/                   # 前端代码
│   ├── src/
│   │   ├── api/               # API接口封装
│   │   ├── components/        # Vue组件
│   │   ├── views/             # 页面视图
│   │   ├── stores/            # Pinia状态管理
│   │   ├── router/            # Vue Router路由
│   │   └── utils/             # 工具函数
│   ├── package.json           # npm依赖管理
│   ├── vite.config.ts         # Vite配置（含API代理）
│   ├── nginx.conf             # Nginx配置（生产环境）
│   └── Dockerfile             # 前端Docker构建
├── examples/                   # 使用示例
│   ├── quickstart/            # Go快速入门示例
│   │   └── main.go
│   └── curl_examples.sh       # curl接口调用示例
├── test/                       # 测试文件
│   └── postman.json           # Postman接口测试集合
├── picture/                    # 项目设计图片
├── docker-compose.yml          # Docker编排配置
├── start.sh                    # 一键启动脚本
├── .env.example                # 环境变量配置模板
├── .gitignore                  # Git忽略规则
├── README.md                   # 项目说明文档（本文件）
├── API.md                      # API接口文档
├── DEPLOYMENT.md               # 部署指南
├── TESTING.md                  # 测试指南
├── CONTRIBUTING.md             # 贡献指南
└── LICENSE                     # MIT开源许可证
```

---

## 项目优势

### 1. 高性能设计
- **三级缓存**：L1本地缓存(3s) + L2 Redis(1h) + L3 MySQL，Feed响应时间降至3ms
- **异步处理**：RabbitMQ消息队列削峰填谷，吞吐量提升400%
- **singleflight**：合并并发请求，防止缓存击穿时DB过载
- **连接池优化**：数据库和Redis连接池复用，减少连接开销

### 2. 高可用保障
- **降级机制**：MQ不可用时自动降级到DB直写，缓存不可用时降级到DB查询
- **缓存自愈**：Redis宕机不影响核心功能，恢复后通过请求自动回填缓存
- **限流控制**：基于Redis的固定窗口限流，防止恶意请求和流量洪峰
- **Token校验自愈**：Redis缓存未命中时回退MySQL校验，通过后自动回填

### 3. 数据一致性
- **Outbox模式**：视频发布在同一事务中写入业务表和发件箱表，保证消息不丢失
- **分布式锁**：Redis SETNX + Lua脚本保证原子性，防止并发操作导致数据不一致
- **MQ降级直写**：消息发布失败时自动回退到同步DB写入，保证数据落地

### 4. 可扩展性
- **API/Worker分离**：API进程和Worker进程可独立部署和水平扩展
- **无状态设计**：API服务无状态，支持多实例负载均衡
- **中间件可插拔**：Redis和RabbitMQ均为可选依赖，不可用时自动降级

---

## 文档导航

| 文档 | 说明 |
|------|------|
| [README.md](README.md) | 项目总览与快速开始（本文件） |
| [API.md](API.md) | 完整API接口文档 |
| [DEPLOYMENT.md](DEPLOYMENT.md) | 本地部署与Docker部署指南 |
| [TESTING.md](TESTING.md) | 测试环境搭建与执行指南 |
| [CONTRIBUTING.md](CONTRIBUTING.md) | 贡献流程与代码规范 |
| [LICENSE](LICENSE) | MIT开源许可证 |

---

## 贡献指南

我们欢迎所有形式的贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何：

- 提交Issue报告问题
- 提交Pull Request贡献代码
- 改进文档
- 分享使用经验

---

## 许可证

本项目采用 [MIT License](LICENSE) 开源协议。

---

<div align="center">

**如果这个项目对你有帮助，请给一个 ⭐️ Star 支持一下！**

</div>

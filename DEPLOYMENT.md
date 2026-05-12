# 短视频Feed流系统 - 本地部署指南

本文档提供详细的本地开发环境搭建步骤，帮助开发者快速启动和运行项目。

## 目录

- [环境要求](#环境要求)
- [快速启动（Docker Compose）](#快速启动docker-compose)
- [本地开发环境搭建](#本地开发环境搭建)
- [配置说明](#配置说明)
- [常见问题](#常见问题)

---

## 环境要求

### 操作系统
- Windows 10/11
- macOS 10.15+
- Linux (Ubuntu 18.04+, CentOS 7+)

### 必需软件

| 软件 | 最低版本 | 推荐版本 | 说明 |
|------|---------|---------|------|
| Go | 1.21+ | 1.24+ | 后端开发语言 |
| Node.js | 18.0+ | 20.0+ | 前端运行环境 |
| npm | 9.0+ | 10.0+ | 前端包管理器 |
| Docker | 20.10+ | 24.0+ | 容器化部署（推荐） |
| Docker Compose | 2.0+ | 2.20+ | 多容器编排 |
| MySQL | 8.0+ | 8.0+ | 关系型数据库 |
| Redis | 6.0+ | 7.0+ | 缓存数据库 |
| RabbitMQ | 3.10+ | 3.12+ | 消息队列 |

### 版本检查命令

```bash
# 检查Go版本
go version

# 检查Node.js版本
node --version

# 检查npm版本
npm --version

# 检查Docker版本
docker --version

# 检查Docker Compose版本
docker compose version
```

---

## 快速启动（Docker Compose）

### 前置条件
- 已安装 Docker Desktop（Windows/macOS）或 Docker Engine + Docker Compose（Linux）

### 启动步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/raaaaap/short-video-feed-system.git
   cd short-video-feed-system
   ```

2. **一键启动所有服务**
   ```bash
   docker compose up -d --build
   ```

3. **验证服务状态**
   ```bash
   docker compose ps
   ```

   预期输出：
   ```
   NAME                    STATUS    PORTS
   feedsystem-mysql        running   0.0.0.0:3307->3306/tcp
   feedsystem-redis        running   0.0.0.0:6379->6379/tcp
   feedsystem-rabbitmq     running   0.0.0.0:5672->5672/tcp, 0.0.0.0:15672->15672/tcp
   feedsystem-backend      running   0.0.0.0:8080->8080/tcp
   feedsystem-worker       running
   feedsystem-frontend     running   0.0.0.0:5173->80/tcp
   ```

4. **访问服务**
   - 前端界面：http://localhost:5173
   - 后端API：http://localhost:8080
   - RabbitMQ管理台：http://localhost:15672（账号：admin / password123）

### 停止服务

```bash
# 停止所有服务
docker compose down

# 停止并删除数据卷
docker compose down -v
```

### 查看日志

```bash
# 查看所有服务日志
docker compose logs -f

# 查看特定服务日志
docker compose logs -f backend
docker compose logs -f worker
docker compose logs -f frontend
```

---

## 本地开发环境搭建

### 方式一：使用Docker启动依赖服务（推荐）

1. **启动依赖服务**
   ```bash
   docker compose up -d mysql redis rabbitmq
   ```

2. **等待服务就绪**
   ```bash
   # 检查MySQL是否就绪
   docker compose exec mysql mysqladmin ping -h localhost -uroot -p123456

   # 检查Redis是否就绪
   docker compose exec redis redis-cli -a 123456 ping

   # 检查RabbitMQ是否就绪
   docker compose exec rabbitmq rabbitmq-diagnostics -q ping
   ```

3. **启动后端API服务**
   ```bash
   cd backend
   go mod download
   go run ./cmd
   ```

4. **启动Worker进程（新终端）**
   ```bash
   cd backend
   go run ./cmd/worker
   ```

5. **启动前端开发服务器（新终端）**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

### 方式二：完全本地安装

#### 1. 安装MySQL

**Windows:**
1. 下载MySQL安装包：https://dev.mysql.com/downloads/mysql/
2. 运行安装程序，设置root密码为 `123456`
3. 创建数据库：
   ```sql
   CREATE DATABASE feedsystem CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

**macOS (使用Homebrew):**
```bash
brew install mysql
brew services start mysql
mysql -uroot
# 在MySQL中执行
CREATE DATABASE feedsystem CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**Linux (Ubuntu):**
```bash
sudo apt update
sudo apt install mysql-server
sudo mysql_secure_installation
sudo mysql -uroot
# 在MySQL中执行
CREATE DATABASE feedsystem CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 2. 安装Redis

**Windows:**
1. 下载Redis for Windows：https://github.com/microsoftarchive/redis/releases
2. 解压并运行 `redis-server.exe`
3. 设置密码：编辑 `redis.windows.conf`，添加 `requirepass 123456`

**macOS:**
```bash
brew install redis
brew services start redis
# 设置密码
redis-cli CONFIG SET requirepass "123456"
```

**Linux:**
```bash
sudo apt install redis-server
sudo systemctl start redis
sudo systemctl enable redis
# 设置密码
redis-cli CONFIG SET requirepass "123456"
```

#### 3. 安装RabbitMQ

**Windows:**
1. 下载RabbitMQ：https://www.rabbitmq.com/download.html
2. 安装Erlang依赖：https://www.erlang.org/downloads
3. 安装RabbitMQ并启动服务
4. 启用管理插件：
   ```bash
   rabbitmq-plugins enable rabbitmq_management
   ```
5. 创建用户：
   ```bash
   rabbitmqctl add_user admin password123
   rabbitmqctl set_user_tags admin administrator
   rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
   ```

**macOS:**
```bash
brew install rabbitmq
brew services start rabbitmq
# 启用管理插件
rabbitmq-plugins enable rabbitmq_management
# 创建用户
rabbitmqctl add_user admin password123
rabbitmqctl set_user_tags admin administrator
rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
```

**Linux:**
```bash
sudo apt install rabbitmq-server
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server
# 启用管理插件
sudo rabbitmq-plugins enable rabbitmq_management
# 创建用户
sudo rabbitmqctl add_user admin password123
sudo rabbitmqctl set_user_tags admin administrator
sudo rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
```

#### 4. 配置项目

1. **复制环境变量模板**
   ```bash
   cp .env.example .env
   ```

2. **修改配置文件**
   
   编辑 `backend/configs/config.yaml`：
   ```yaml
   server:
     port: 8080

   database:
     host: localhost
     port: 3306
     user: root
     password: 123456
     dbname: feedsystem

   redis:
     host: localhost
     port: 6379
     password: 123456
     db: 0

   rabbitmq:
     host: localhost
     port: 5672
     username: admin
     password: password123
   ```

#### 5. 启动项目

```bash
# 终端1：启动后端API
cd backend
go mod download
go run ./cmd

# 终端2：启动Worker
cd backend
go run ./cmd/worker

# 终端3：启动前端
cd frontend
npm install
npm run dev
```

---

## 配置说明

### 配置文件结构

```
backend/configs/
├── config.yaml          # 本地开发配置
└── config.docker.yaml   # Docker环境配置
```

### 主要配置项

#### 服务器配置
```yaml
server:
  port: 8080              # API服务端口
```

#### 数据库配置
```yaml
database:
  host: localhost         # MySQL主机地址
  port: 3306             # MySQL端口
  user: root             # 数据库用户名
  password: 123456       # 数据库密码
  dbname: feedsystem     # 数据库名称
```

#### Redis配置
```yaml
redis:
  host: localhost         # Redis主机地址
  port: 6379             # Redis端口
  password: 123456       # Redis密码
  db: 0                  # Redis数据库索引
```

#### RabbitMQ配置
```yaml
rabbitmq:
  host: localhost         # RabbitMQ主机地址
  port: 5672             # RabbitMQ端口
  username: admin        # RabbitMQ用户名
  password: password123  # RabbitMQ密码
```

#### 可观测性配置
```yaml
observability:
  pprof:
    enabled: true                    # 是否启用pprof
    api_addr: localhost:6060        # API服务pprof地址
    worker_addr: localhost:6061     # Worker服务pprof地址
```

### 环境变量覆盖

项目支持通过环境变量覆盖配置文件中的值：

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| SERVER_PORT | API服务端口 | 8080 |
| DB_HOST | MySQL主机 | localhost |
| DB_PORT | MySQL端口 | 3306 |
| DB_USER | MySQL用户名 | root |
| DB_PASSWORD | MySQL密码 | 123456 |
| DB_NAME | 数据库名 | feedsystem |
| REDIS_HOST | Redis主机 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| REDIS_PASSWORD | Redis密码 | 123456 |
| RABBITMQ_HOST | RabbitMQ主机 | localhost |
| RABBITMQ_PORT | RabbitMQ端口 | 5672 |
| RABBITMQ_USER | RabbitMQ用户名 | admin |
| RABBITMQ_PASSWORD | RabbitMQ密码 | password123 |

---

## 常见问题

### 1. 端口被占用

**问题：** 启动时提示端口已被占用

**解决方案：**
```bash
# Windows：查找并结束占用端口的进程
netstat -ano | findstr :8080
taskkill /PID <进程ID> /F

# macOS/Linux
lsof -i :8080
kill -9 <PID>
```

或修改配置文件中的端口号。

### 2. MySQL连接失败

**问题：** `Error 1045 (28000): Access denied for user 'root'@'localhost'`

**解决方案：**
1. 检查MySQL服务是否启动
2. 确认用户名和密码正确
3. 检查MySQL用户权限：
   ```sql
   GRANT ALL PRIVILEGES ON feedsystem.* TO 'root'@'localhost';
   FLUSH PRIVILEGES;
   ```

### 3. Redis连接失败

**问题：** `NOAUTH Authentication required`

**解决方案：**
1. 确认Redis密码配置正确
2. 如果Redis未设置密码，在配置文件中将password留空

### 4. RabbitMQ连接失败

**问题：** `Exception (504) Reason: "channel/connection is not open"`

**解决方案：**
1. 检查RabbitMQ服务状态：`rabbitmqctl status`
2. 确认用户权限：
   ```bash
   rabbitmqctl list_users
   rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
   ```

### 5. 前端无法连接后端

**问题：** 前端请求后端API时出现CORS错误

**解决方案：**
1. 确认后端服务已启动
2. 检查 `frontend/vite.config.ts` 中的代理配置：
   ```typescript
   server: {
     proxy: {
       '/api': {
         target: 'http://127.0.0.1:8080',
         changeOrigin: true
       }
     }
   }
   ```

### 6. Docker容器启动失败

**问题：** 容器无法启动或健康检查失败

**解决方案：**
```bash
# 查看容器日志
docker compose logs <服务名>

# 重新构建容器
docker compose down
docker compose up -d --build

# 清理Docker缓存
docker system prune -a
```

### 7. Go依赖下载失败

**问题：** `go: module ...: reading at revision ...: unknown revision`

**解决方案：**
```bash
# 设置Go代理（中国大陆）
go env -w GOPROXY=https://goproxy.cn,direct

# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download
```

### 8. npm依赖安装失败

**问题：** `npm ERR! network request failed`

**解决方案：**
```bash
# 设置npm镜像源（中国大陆）
npm config set registry https://registry.npmmirror.com

# 清理npm缓存
npm cache clean --force

# 删除node_modules重新安装
rm -rf node_modules package-lock.json
npm install
```

---

## 性能优化建议

### 开发环境

1. **Go编译优化**
   ```bash
   # 使用编译缓存加速
   go build -o bin/api ./cmd
   go build -o bin/worker ./cmd/worker
   ```

2. **前端热更新**
   - Vite已内置热更新功能
   - 修改代码后自动刷新浏览器

### 生产环境

1. **数据库优化**
   - 调整MySQL连接池大小
   - 启用查询缓存
   - 定期优化表

2. **Redis优化**
   - 调整最大内存限制
   - 配置淘汰策略
   - 启用持久化

3. **RabbitMQ优化**
   - 调整消费者并发数
   - 配置消息确认机制
   - 启用镜像队列

---

## 下一步

- 查看 [README.md](README.md) 了解项目功能
- 查看 [API.md](API.md) 了解接口文档
- 查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何贡献代码

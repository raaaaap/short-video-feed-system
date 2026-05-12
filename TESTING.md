# 测试指南

本文档详细说明短视频Feed流系统的测试方法，包括环境搭建、测试执行、结果分析和扩展测试。

---

## 目录

- [测试环境搭建](#测试环境搭建)
- [测试方法概览](#测试方法概览)
- [Postman接口测试](#postman接口测试)
- [Go单元测试](#go单元测试)
- [手动功能测试](#手动功能测试)
- [性能测试](#性能测试)
- [测试结果分析](#测试结果分析)
- [添加新测试用例](#添加新测试用例)

---

## 测试环境搭建

### 必需服务

测试前需确保以下服务正常运行：

| 服务 | 用途 | 启动命令 |
|------|------|---------|
| MySQL 8.0+ | 数据存储 | `docker compose up -d mysql` |
| Redis 7.0+ | 缓存/排行榜 | `docker compose up -d redis` |
| RabbitMQ 3.12+ | 消息队列 | `docker compose up -d rabbitmq` |

### 一键启动依赖

```bash
docker compose up -d mysql redis rabbitmq
```

### 验证服务状态

```bash
# 检查MySQL
docker compose exec mysql mysqladmin ping -h localhost -uroot -p123456

# 检查Redis
docker compose exec redis redis-cli -a 123456 ping

# 检查RabbitMQ
docker compose exec rabbitmq rabbitmq-diagnostics -q ping
```

### 启动应用服务

```bash
# 终端1：启动API服务
cd backend && go run ./cmd

# 终端2：启动Worker服务
cd backend && go run ./cmd/worker

# 终端3：启动前端（可选）
cd frontend && npm run dev
```

---

## 测试方法概览

| 测试类型 | 工具 | 覆盖范围 | 位置 |
|---------|------|---------|------|
| 接口测试 | Postman | 全部API端点 | `test/postman.json` |
| 单元测试 | Go testing | 核心模块逻辑 | `backend/internal/**/*_test.go` |
| 手动测试 | 浏览器/curl | 前端交互 | - |
| 性能测试 | pprof | 性能瓶颈分析 | 内置pprof服务 |

---

## Postman接口测试

### 导入测试集合

1. 打开Postman
2. 点击左上角 `Import`
3. 选择项目中的 `test/postman.json`
4. 导入后可见 `FeedSystem API (All-in-One)` 集合

### 推荐执行顺序

按以下顺序执行可确保变量自动传递：

```
1. Account/Register    → 注册测试账号
2. Account/Login       → 获取JWT Token（自动保存到jwt_token变量）
3. Account/Find By Username → 获取accountId（自动保存）
4. Video/Publish       → 发布视频（自动保存videoId）
5. Like/like           → 点赞视频
6. Like/isLiked        → 验证点赞状态
7. Comment/publish     → 发布评论
8. Comment/listAll     → 查看评论列表
9. Social/follow       → 关注用户
10. Feed/listLatest    → 浏览最新Feed
11. Feed/listByPopularity → 浏览热榜
12. Feed/listByFollowing  → 浏览关注Feed
13. Account/logout     → 登出
```

### 使用Collection Runner批量执行

1. 点击集合名称右侧的 `...`
2. 选择 `Run collection`
3. 按上述顺序排列请求
4. 点击 `Run FeedSystem API`
5. 查看测试结果

### 预置变量说明

| 变量 | 默认值 | 说明 |
|------|--------|------|
| host | http://localhost:8080 | API地址 |
| username | 自动生成 | 测试用户名 |
| password | pass123 | 测试密码 |
| jwt_token | 空 | 登录后自动填充 |
| accountId | 1 | 用户ID |
| videoId | 1 | 视频ID |
| feedLimit | 10 | Feed每页数量 |

---

## Go单元测试

### 运行现有测试

```bash
cd backend

# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./internal/observability/...

# 运行测试并显示详细输出
go test -v ./...

# 运行测试并生成覆盖率报告
go test -cover ./...

# 生成详细的覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 当前测试覆盖

| 模块 | 测试文件 | 覆盖内容 |
|------|---------|---------|
| observability | pprof_test.go | pprof服务启停 |

### 测试编写规范

```go
package observability

import (
    "testing"
)

func TestNewPprofServer(t *testing.T) {
    tests := []struct {
        name    string
        enabled bool
        addr    string
        wantErr bool
    }{
        {
            name:    "禁用pprof",
            enabled: false,
            addr:    "",
            wantErr: false,
        },
        {
            name:    "启用pprof",
            enabled: true,
            addr:    "localhost:0",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            srv, err := NewPprofServer("test", tt.enabled, tt.addr)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewPprofServer() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if srv != nil {
                srv.Close()
            }
        })
    }
}
```

---

## 手动功能测试

### 使用curl测试API

#### 1. 注册账号

```bash
curl -X POST http://localhost:8080/account/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test_user","password":"test123"}'
```

预期输出：
```json
{"message":"successfully created"}
```

#### 2. 登录获取Token

```bash
curl -X POST http://localhost:8080/account/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test_user","password":"test123"}'
```

预期输出：
```json
{"token":"eyJhbGciOiJIUzI1NiIs..."}
```

#### 3. 发布视频（需要JWT）

```bash
TOKEN="你的JWT令牌"

curl -X POST http://localhost:8080/video/publish \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title":"测试视频",
    "description":"这是一个测试视频",
    "play_url":"http://example.com/test.mp4",
    "cover_url":"http://example.com/test.jpg"
  }'
```

#### 4. 点赞视频

```bash
curl -X POST http://localhost:8080/like/like \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"video_id":1}'
```

#### 5. 浏览最新Feed

```bash
curl -X POST http://localhost:8080/feed/listLatest \
  -H "Content-Type: application/json" \
  -d '{"limit":10,"latest_time":0}'
```

#### 6. 浏览热榜

```bash
curl -X POST http://localhost:8080/feed/listByPopularity \
  -H "Content-Type: application/json" \
  -d '{"limit":10}'
```

### 前端功能测试

1. 访问 http://localhost:5173
2. 注册新账号
3. 登录系统
4. 发布视频
5. 浏览Feed流
6. 点赞、评论
7. 关注其他用户
8. 查看关注Feed

---

## 性能测试

### 使用pprof分析

系统内置了pprof性能分析服务，默认配置：

| 服务 | 地址 | 说明 |
|------|------|------|
| API pprof | localhost:6060 | API服务性能分析 |
| Worker pprof | localhost:6061 | Worker服务性能分析 |

> 需要在配置文件中启用pprof：`observability.pprof.enabled: true`

### 常用pprof命令

```bash
# 查看CPU性能分析（采集30秒）
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 查看内存分配
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 查看阻塞分析
go tool pprof http://localhost:6060/debug/pprof/block

# Web界面查看（推荐）
# 浏览器访问 http://localhost:6060/debug/pprof/
```

### 使用hey进行HTTP压测

```bash
# 安装hey
go install github.com/rakyll/hey@latest

# 压测Feed接口
hey -n 1000 -c 50 -m POST \
  -H "Content-Type: application/json" \
  -d '{"limit":10,"latest_time":0}' \
  http://localhost:8080/feed/listLatest

# 压测登录接口
hey -n 100 -c 10 -m POST \
  -H "Content-Type: application/json" \
  -d '{"username":"test_user","password":"test123"}' \
  http://localhost:8080/account/login
```

---

## 测试结果分析

### Postman测试结果

在Collection Runner中查看：

| 指标 | 说明 |
|------|------|
| Passed | 测试断言通过的数量 |
| Failed | 测试断言失败的数量 |
| Status Code | HTTP状态码是否为预期值 |
| Response Time | 响应时间，正常应 < 100ms |

### Go测试结果

```bash
# 正常输出
ok      feedsystem_video_go/internal/observability    0.123s

# 失败输出
--- FAIL: TestNewPprofServer (0.00s)
    pprof_test.go:25: unexpected error
FAIL
```

### 性能基准参考

| 接口 | 目标响应时间 | 说明 |
|------|-------------|------|
| /feed/listLatest | < 10ms | 缓存命中时 |
| /feed/listLatest | < 50ms | 缓存未命中时 |
| /video/getDetail | < 5ms | 缓存命中时 |
| /account/login | < 100ms | 含bcrypt校验 |
| /like/like | < 50ms | MQ正常时 |

---

## 添加新测试用例

### 添加Postman测试

1. 在Postman中创建新请求
2. 添加Tests脚本：

```javascript
// 验证状态码
pm.test("状态码为200", function () {
    pm.response.to.have.status(200);
});

// 验证响应格式
pm.test("返回JSON格式", function () {
    pm.response.to.be.json;
});

// 验证响应字段
let json = pm.response.json();
pm.test("包含必要字段", function () {
    pm.expect(json).to.have.property("id");
});

// 自动保存变量
if (json && json.id) {
    pm.collectionVariables.set("newVar", String(json.id));
}
```

3. 导出集合：`File → Export → Collection v2.1`
4. 替换 `test/postman.json`

### 添加Go单元测试

1. 在对应包目录下创建 `_test.go` 文件
2. 编写测试函数：

```go
func TestFunctionName(t *testing.T) {
    got := FunctionName()
    want := expectedValue
    if got != want {
        t.Errorf("FunctionName() = %v, want %v", got, want)
    }
}
```

3. 运行测试验证：

```bash
go test -v ./internal/yourpackage/...
```

### 测试命名规范

| 测试类型 | 命名格式 | 示例 |
|---------|---------|------|
| 基本测试 | TestXxx | TestLogin |
| 分组测试 | TestXxx_Yyy | TestLogin_Success |
| 基准测试 | BenchmarkXxx | BenchmarkListLatest |
| 示例测试 | ExampleXxx | ExampleFeedService |

---

<div align="center">

**完善的测试是高质量代码的保障，欢迎补充更多测试用例！**

</div>

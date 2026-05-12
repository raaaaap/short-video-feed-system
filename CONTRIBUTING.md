# 贡献指南

感谢你对短视频Feed流系统的关注！我们欢迎所有形式的贡献，包括但不限于代码、文档、问题反馈和使用经验分享。

---

## 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发环境搭建](#开发环境搭建)
- [代码风格规范](#代码风格规范)
- [提交信息规范](#提交信息规范)
- [分支管理策略](#分支管理策略)
- [Pull Request流程](#pull-request流程)
- [Issue提交规范](#issue提交规范)
- [版本控制规范](#版本控制规范)
- [沟通渠道](#沟通渠道)

---

## 行为准则

- 尊重所有贡献者，保持友善和建设性的沟通
- 关注问题本身，避免人身攻击
- 欢迎不同观点和经验水平的贡献者
- 以项目利益为先，保持开放和协作的态度

---

## 如何贡献

### 贡献类型

| 类型 | 说明 | 示例 |
|------|------|------|
| Bug修复 | 修复已知问题 | 修复缓存击穿逻辑缺陷 |
| 新功能 | 添加新特性 | 新增视频搜索功能 |
| 性能优化 | 提升系统性能 | 优化数据库查询 |
| 文档改进 | 完善文档内容 | 补充API文档、修正错别字 |
| 测试补充 | 增加测试覆盖 | 添加单元测试、集成测试 |
| 代码重构 | 改善代码质量 | 提取公共方法、优化命名 |

### 贡献流程

```
1. Fork仓库 → 2. 创建分支 → 3. 开发修改 → 4. 提交PR → 5. 代码审查 → 6. 合并
```

---

## 开发环境搭建

### 前置要求

- Go 1.24+
- Node.js 18+
- Docker & Docker Compose

### 搭建步骤

```bash
# 1. Fork并克隆项目
git clone https://github.com/yourusername/feedsystem_video_go.git
cd feedsystem_video_go

# 2. 启动依赖服务
docker compose up -d mysql redis rabbitmq

# 3. 启动后端
cd backend
go mod download
go run ./cmd

# 4. 启动Worker（新终端）
cd backend
go run ./cmd/worker

# 5. 启动前端（新终端）
cd frontend
npm install
npm run dev
```

详细步骤请参考 [DEPLOYMENT.md](DEPLOYMENT.md)

---

## 代码风格规范

### Go代码规范

- 遵循 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `gofmt` 格式化代码
- 使用 `go vet` 检查代码问题
- 包命名使用小写单词，不使用下划线或驼峰
- 导出函数必须添加注释，以函数名开头

```go
// VideoService 视频业务服务
type VideoService struct {
    videoRepo    *VideoRepository
    cache        *rediscache.Client
    popularityMQ *rabbitmq.PopularityMQ
}

// Publish 发布视频，使用事务性发件箱模式保证数据一致性
func (vs *VideoService) Publish(ctx context.Context, ...) error {
    // ...
}
```

### Vue/TypeScript代码规范

- 遵循 [Vue.js 风格指南](https://vuejs.org/style-guide/)
- 使用 TypeScript 严格模式
- 组件命名使用 PascalCase
- 事件命名使用 kebab-case

### 通用规范

- 每个函数/方法只做一件事
- 避免过深的嵌套，尽早返回
- 变量命名要有意义，避免缩写（除常见缩写如ID、URL）
- 错误处理不要忽略，至少记录日志

---

## 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type类型

| 类型 | 说明 |
|------|------|
| feat | 新功能 |
| fix | Bug修复 |
| docs | 文档变更 |
| style | 代码格式（不影响功能） |
| refactor | 代码重构 |
| perf | 性能优化 |
| test | 测试相关 |
| chore | 构建/工具变更 |

### 示例

```
feat(feed): 添加按热度排序Feed流接口

实现Redis滑动窗口热榜+快照分页，支持降级到MySQL游标分页

Closes #12
```

```
fix(like): 修复MQ发布失败时点赞计数不一致的问题

当LikeMQ发布失败时，回退到事务直写DB，确保点赞计数正确更新
```

---

## 分支管理策略

### 分支命名

| 分支 | 命名格式 | 说明 |
|------|---------|------|
| main | `main` | 主分支，稳定版本 |
| develop | `develop` | 开发分支 |
| 功能分支 | `feat/<feature-name>` | 新功能开发 |
| 修复分支 | `fix/<bug-name>` | Bug修复 |
| 热修复 | `hotfix/<issue-name>` | 紧急修复 |
| 文档分支 | `docs/<doc-name>` | 文档更新 |

### 示例

```bash
# 创建功能分支
git checkout -b feat/video-search

# 创建修复分支
git checkout -b fix/like-count-error

# 创建文档分支
git checkout -b docs/api-reference
```

### 合并策略

- 功能分支 → develop：使用 Squash Merge
- develop → main：使用 Merge Commit
- hotfix → main：使用 Merge Commit，同时同步到 develop

---

## Pull Request流程

### 创建PR

1. 确保代码通过本地测试
2. 确保代码风格符合规范
3. 在GitHub上创建Pull Request
4. 填写PR模板中的必要信息

### PR标题格式

```
<type>(<scope>): <description>
```

示例：`feat(feed): 添加视频搜索功能`

### PR描述模板

```markdown
## 变更类型
- [ ] 新功能
- [ ] Bug修复
- [ ] 性能优化
- [ ] 文档更新
- [ ] 代码重构

## 变更说明
<!-- 描述本次PR的主要变更内容 -->

## 关联Issue
<!-- 关联的Issue编号，如 Closes #123 -->

## 测试说明
<!-- 描述如何测试本次变更 -->

## 检查清单
- [ ] 代码通过编译
- [ ] 通过现有测试
- [ ] 添加了必要的测试
- [ ] 更新了相关文档
```

### 代码审查

- 至少需要一位维护者审查通过
- 审查关注点：
  - 代码逻辑正确性
  - 代码风格一致性
  - 是否有潜在的性能问题
  - 是否有安全隐患
  - 是否有更好的实现方式

---

## Issue提交规范

### Bug报告模板

```markdown
## Bug描述
<!-- 简要描述问题 -->

## 复现步骤
1. ...
2. ...
3. ...

## 期望行为
<!-- 描述期望的正确行为 -->

## 实际行为
<!-- 描述实际发生的错误行为 -->

## 环境信息
- 操作系统：
- Go版本：
- Redis版本：
- MySQL版本：
- RabbitMQ版本：

## 附加信息
<!-- 日志、截图等 -->
```

### 功能请求模板

```markdown
## 功能描述
<!-- 描述期望的功能 -->

## 使用场景
<!-- 描述什么场景下需要这个功能 -->

## 建议实现方式
<!-- 如果有实现思路，可以在这里描述 -->

## 替代方案
<!-- 是否有其他替代方案 -->
```

### Issue标签

| 标签 | 说明 |
|------|------|
| bug | Bug报告 |
| enhancement | 功能请求 |
| documentation | 文档相关 |
| good first issue | 适合新贡献者 |
| help wanted | 需要帮助 |
| performance | 性能相关 |
| question | 问题咨询 |

---

## 版本控制规范

### 版本号格式

遵循 [语义化版本](https://semver.org/lang/zh-CN/)：

```
MAJOR.MINOR.PATCH
```

- **MAJOR**：不兼容的API变更
- **MINOR**：向后兼容的功能新增
- **PATCH**：向后兼容的问题修复

### 示例

- `v1.0.0` — 首个正式版本
- `v1.1.0` — 新增视频搜索功能
- `v1.1.1` — 修复搜索结果排序问题

---

## 沟通渠道

| 渠道 | 用途 |
|------|------|
| [GitHub Issues](https://github.com/LeoninCS/feedsystem_video_go/issues) | Bug报告、功能请求 |
| [GitHub Discussions](https://github.com/LeoninCS/feedsystem_video_go/discussions) | 技术讨论、使用问答 |
| Pull Request | 代码贡献、代码审查 |

---

<div align="center">

**感谢你的贡献！每一个PR和Issue都是对项目的宝贵支持。**

</div>

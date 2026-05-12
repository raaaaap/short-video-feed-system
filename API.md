# 短视频Feed流系统 - API接口文档

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

---

## 目录

- [概述](#概述)
- [认证方式](#认证方式)
- [限流策略](#限流策略)
- [通用响应格式](#通用响应格式)
- [账号模块](#账号模块)
- [视频模块](#视频模块)
- [点赞模块](#点赞模块)
- [评论模块](#评论模块)
- [社交模块](#社交模块)
- [Feed流模块](#feed流模块)
- [错误码说明](#错误码说明)

---

## 概述

- **基础URL**：`http://localhost:8080`
- **请求方式**：所有接口均为 `POST`
- **请求格式**：`application/json`（文件上传接口除外）
- **字符编码**：UTF-8

---

## 认证方式

系统采用JWT（JSON Web Token）进行身份认证，提供两种认证模式：

### 硬认证（JWTAuth）

必须携带有效的JWT令牌，否则返回 `401 Unauthorized`。

```
Authorization: Bearer <jwt_token>
```

适用于：登出、改名、视频上传/发布、点赞/取消点赞、评论发布/删除、关注/取关、关注Feed流

### 软认证（SoftJWTAuth）

不强制要求携带令牌。无令牌时以匿名身份访问；有令牌时必须合法，否则返回 `401`。

适用于：最新Feed流、点赞数Feed流、热度Feed流

### 获取令牌

通过 `/account/login` 接口登录后获取：

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 令牌校验流程

```
请求携带Token → 解析JWT Claims → 查询Redis缓存(account:<id>)
  → Redis命中：对比Token是否一致
  → Redis未命中：查询MySQL数据库 → 校验通过后回填Redis
```

---

## 限流策略

基于Redis `INCR` + `EXPIRE` 实现的固定窗口限流，超限返回 `429 Too Many Requests`。

| 接口 | 限流维度 | 限制 |
|------|---------|------|
| `/account/login` | IP | 10次/分钟 |
| `/account/register` | IP | 5次/小时 |
| `/like/like` | 账号 | 30次/分钟 |
| `/like/unlike` | 账号 | 30次/分钟 |
| `/comment/publish` | 账号 | 10次/分钟 |
| `/comment/delete` | 账号 | 10次/分钟 |
| `/social/follow` | 账号 | 20次/分钟 |
| `/social/unfollow` | 账号 | 20次/分钟 |

---

## 通用响应格式

### 成功响应

```json
{
  "message": "操作成功"
}
```

或包含数据的响应：

```json
{
  "id": 1,
  "username": "example_user"
}
```

### 错误响应

```json
{
  "error": "错误描述信息"
}
```

### HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或Token无效 |
| 404 | 资源不存在 |
| 429 | 请求频率超限 |
| 500 | 服务器内部错误 |

---

## 账号模块

### 注册账号

创建新用户账号，密码使用bcrypt加密存储。

```
POST /account/register
```

**限流**：5次/小时/IP

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名，唯一 |
| password | string | 是 | 密码 |

**请求示例**：

```json
{
  "username": "new_user",
  "password": "my_password"
}
```

**成功响应**（200）：

```json
{
  "message": "successfully created"
}
```

**失败响应**（400）：

```json
{
  "error": "username already exists"
}
```

---

### 登录

验证用户名和密码，返回JWT令牌。登录成功后将Token写入数据库和Redis缓存。

```
POST /account/login
```

**限流**：10次/分钟/IP

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例**：

```json
{
  "username": "new_user",
  "password": "my_password"
}
```

**成功响应**（200）：

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**失败响应**（400）：

```json
{
  "error": "invalid username or password"
}
```

---

### 修改密码

验证旧密码后更新为新密码，修改成功后自动登出（清空Token）。

```
POST /account/changePassword
```

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| old_password | string | 是 | 旧密码 |
| new_password | string | 是 | 新密码 |

**请求示例**：

```json
{
  "username": "new_user",
  "old_password": "my_password",
  "new_password": "new_password"
}
```

**成功响应**（200）：

```json
{
  "message": "successfully password changed"
}
```

**失败响应**（400）：

```json
{
  "error": "unsuccessfully password changed"
}
```

---

### 按ID查询用户

```
POST /account/findByID
```

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | uint | 是 | 用户ID |

**请求示例**：

```json
{
  "id": 1
}
```

**成功响应**（200）：

```json
{
  "id": 1,
  "username": "new_user"
}
```

> 注意：响应中不包含密码和Token字段。

---

### 按用户名查询用户

```
POST /account/findByUsername
```

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |

**请求示例**：

```json
{
  "username": "new_user"
}
```

**成功响应**（200）：

```json
{
  "id": 1,
  "username": "new_user"
}
```

---

### 登出 🔒

清除当前用户的Token，使旧令牌立即失效。

```
POST /account/logout
```

**认证**：需要JWT

**请求头**：

```
Authorization: Bearer <jwt_token>
```

**请求体**：无（或空JSON `{}`）

**成功响应**（200）：

```json
{
  "message": "successfully logged out"
}
```

---

### 改名 🔒

修改用户名并生成新的JWT令牌，旧令牌立即失效。

```
POST /account/rename
```

**认证**：需要JWT

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| new_username | string | 是 | 新用户名，唯一 |

**请求示例**：

```json
{
  "new_username": "cool_user"
}
```

**成功响应**（200）：

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**失败响应**（400）：

```json
{
  "error": "username already exists"
}
```

---

## 视频模块

### 上传视频文件 🔒

上传视频文件到服务器，仅支持MP4格式。

```
POST /video/uploadVideo
```

**认证**：需要JWT

**请求格式**：`multipart/form-data`

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 视频文件，仅支持.mp4，最大200MB |

**成功响应**（200）：

```json
{
  "url": "/static/uploads/video_xxx.mp4",
  "play_url": "http://localhost:8080/static/uploads/video_xxx.mp4"
}
```

---

### 上传封面图 🔒

上传视频封面图片。

```
POST /video/uploadCover
```

**认证**：需要JWT

**请求格式**：`multipart/form-data`

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 封面图片，支持.jpg/.jpeg/.png/.webp，最大10MB |

**成功响应**（200）：

```json
{
  "url": "/static/uploads/cover_xxx.jpg",
  "cover_url": "http://localhost:8080/static/uploads/cover_xxx.jpg"
}
```

---

### 发布视频 🔒

发布视频记录到数据库，使用事务性发件箱模式保证数据一致性。

```
POST /video/publish
```

**认证**：需要JWT

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 视频标题 |
| description | string | 否 | 视频描述 |
| play_url | string | 是 | 视频播放地址 |
| cover_url | string | 是 | 封面图片地址 |

**请求示例**：

```json
{
  "title": "我的第一个视频",
  "description": "这是一个很棒的视频",
  "play_url": "http://localhost:8080/static/uploads/video_xxx.mp4",
  "cover_url": "http://localhost:8080/static/uploads/cover_xxx.jpg"
}
```

**成功响应**（200）：

```json
{
  "id": 1,
  "author_id": 1,
  "username": "new_user",
  "title": "我的第一个视频",
  "description": "这是一个很棒的视频",
  "play_url": "http://localhost:8080/static/uploads/video_xxx.mp4",
  "cover_url": "http://localhost:8080/static/uploads/cover_xxx.jpg",
  "create_time": 1700000000,
  "likes_count": 0,
  "popularity": 0
}
```

---

### 获取视频详情

根据ID获取视频详情，支持Redis缓存和分布式锁防击穿。

```
POST /video/getDetail
```

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | uint | 是 | 视频ID |

**请求示例**：

```json
{
  "id": 1
}
```

**成功响应**（200）：

```json
{
  "id": 1,
  "author_id": 1,
  "username": "new_user",
  "title": "我的第一个视频",
  "description": "这是一个很棒的视频",
  "play_url": "http://localhost:8080/static/uploads/video_xxx.mp4",
  "cover_url": "http://localhost:8080/static/uploads/cover_xxx.jpg",
  "create_time": 1700000000,
  "likes_count": 10,
  "popularity": 15
}
```

---

### 按作者查询视频列表

```
POST /video/listByAuthorID
```

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| author_id | uint | 是 | 作者用户ID |

**请求示例**：

```json
{
  "author_id": 1
}
```

**成功响应**（200）：

```json
[
  {
    "id": 1,
    "author_id": 1,
    "username": "new_user",
    "title": "我的第一个视频",
    "description": "这是一个很棒的视频",
    "play_url": "...",
    "cover_url": "...",
    "create_time": 1700000000,
    "likes_count": 10,
    "popularity": 15
  }
]
```

---

## 点赞模块

### 点赞 🔒

对视频进行点赞，采用MQ优先+DB兜底的双写策略。

```
POST /like/like
```

**认证**：需要JWT

**限流**：30次/分钟/账号

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| video_id | uint | 是 | 视频ID |

**请求示例**：

```json
{
  "video_id": 1
}
```

**成功响应**（200）：

```json
{
  "message": "liked"
}
```

**失败响应**（400）：

```json
{
  "error": "already liked"
}
```

---

### 取消点赞 🔒

取消对视频的点赞。

```
POST /like/unlike
```

**认证**：需要JWT

**限流**：30次/分钟/账号

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| video_id | uint | 是 | 视频ID |

**请求示例**：

```json
{
  "video_id": 1
}
```

**成功响应**（200）：

```json
{
  "message": "unliked"
}
```

**失败响应**（400）：

```json
{
  "error": "not liked yet"
}
```

---

### 查询是否已点赞 🔒

查询当前用户是否已对指定视频点赞。

```
POST /like/isLiked
```

**认证**：需要JWT

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| video_id | uint | 是 | 视频ID |

**请求示例**：

```json
{
  "video_id": 1
}
```

**成功响应**（200）：

```json
{
  "is_liked": true
}
```

---

### 获取我点赞过的视频 🔒

获取当前用户点赞过的所有视频列表。

```
POST /like/listMyLikedVideos
```

**认证**：需要JWT

**请求体**：无（或空JSON `{}`）

**成功响应**（200）：

```json
[
  {
    "id": 1,
    "author_id": 1,
    "username": "new_user",
    "title": "我的第一个视频",
    "description": "...",
    "play_url": "...",
    "cover_url": "...",
    "create_time": 1700000000,
    "likes_count": 10,
    "popularity": 15
  }
]
```

---

## 评论模块

### 发布评论 🔒

对视频发布评论，采用MQ优先+DB兜底策略。

```
POST /comment/publish
```

**认证**：需要JWT

**限流**：10次/分钟/账号

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| video_id | uint | 是 | 视频ID |
| content | string | 是 | 评论内容 |

**请求示例**：

```json
{
  "video_id": 1,
  "content": "这个视频太棒了！"
}
```

**成功响应**（200）：

```json
{
  "message": "comment published"
}
```

---

### 获取视频全部评论

获取指定视频的所有评论列表。

```
POST /comment/listAll
```

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| video_id | uint | 是 | 视频ID |

**请求示例**：

```json
{
  "video_id": 1
}
```

**成功响应**（200）：

```json
[
  {
    "id": 1,
    "username": "commenter",
    "video_id": 1,
    "author_id": 2,
    "content": "这个视频太棒了！",
    "created_at": 1700000000
  }
]
```

---

### 删除评论 🔒

删除指定评论，仅评论作者可删除。

```
POST /comment/delete
```

**认证**：需要JWT

**限流**：10次/分钟/账号

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| comment_id | uint | 是 | 评论ID |

**请求示例**：

```json
{
  "comment_id": 1
}
```

**成功响应**（200）：

```json
{
  "message": "comment deleted"
}
```

**失败响应**（400）：

```json
{
  "error": "not the author"
}
```

---

## 社交模块

### 关注 🔒

关注指定用户。

```
POST /social/follow
```

**认证**：需要JWT

**限流**：20次/分钟/账号

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| vlogger_id | uint | 是 | 被关注者ID |

**请求示例**：

```json
{
  "vlogger_id": 2
}
```

**成功响应**（200）：

```json
{
  "message": "followed"
}
```

**失败响应**（400）：

```json
{
  "error": "already followed"
}
```

或

```json
{
  "error": "cannot follow yourself"
}
```

---

### 取消关注 🔒

取消关注指定用户。

```
POST /social/unfollow
```

**认证**：需要JWT

**限流**：20次/分钟/账号

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| vlogger_id | uint | 是 | 被关注者ID |

**请求示例**：

```json
{
  "vlogger_id": 2
}
```

**成功响应**（200）：

```json
{
  "message": "unfollowed"
}
```

---

### 获取粉丝列表 🔒

获取指定用户的粉丝列表。

```
POST /social/getAllFollowers
```

**认证**：需要JWT

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| vlogger_id | uint | 否 | 目标用户ID，为空则查当前登录用户 |

**请求示例**：

```json
{
  "vlogger_id": 2
}
```

**成功响应**（200）：

```json
{
  "followers": [
    {
      "id": 1,
      "username": "follower_user"
    }
  ]
}
```

---

### 获取关注列表 🔒

获取指定用户关注的人列表。

```
POST /social/getAllVloggers
```

**认证**：需要JWT

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| follower_id | uint | 否 | 目标用户ID，为空则查当前登录用户 |

**请求示例**：

```json
{
  "follower_id": 1
}
```

**成功响应**（200）：

```json
{
  "vloggers": [
    {
      "id": 2,
      "username": "followed_user"
    }
  ]
}
```

---

## Feed流模块

Feed流模块是系统的核心，提供四种不同的视频流获取方式，支持三级缓存架构和冷热数据分离。

### Feed视频项统一格式

所有Feed接口返回的视频项格式一致：

```json
{
  "id": 1,
  "author": {
    "id": 1,
    "username": "author_name"
  },
  "title": "视频标题",
  "description": "视频描述",
  "play_url": "http://localhost:8080/static/uploads/video.mp4",
  "cover_url": "http://localhost:8080/static/uploads/cover.jpg",
  "create_time": 1700000000,
  "likes_count": 10,
  "is_liked": false
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 视频ID |
| author | object | 作者信息（id + username） |
| title | string | 视频标题 |
| description | string | 视频描述 |
| play_url | string | 视频播放地址 |
| cover_url | string | 封面图片地址 |
| create_time | int64 | 创建时间（Unix时间戳，秒） |
| likes_count | int64 | 点赞数 |
| is_liked | bool | 当前用户是否已点赞（未登录时为false） |

---

### 最新视频流

获取按时间倒序排列的最新视频，支持冷热分离和游标分页。

```
POST /feed/listLatest
```

**认证**：软认证（SoftJWTAuth）—— 无Token可访问，有Token则解析用户信息

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 是 | 每页数量，建议10-30 |
| latest_time | int64 | 是 | 游标时间戳，首次传0，后续传响应中的next_time |

**请求示例（首次请求）**：

```json
{
  "limit": 10,
  "latest_time": 0
}
```

**请求示例（翻页请求）**：

```json
{
  "limit": 10,
  "latest_time": 1700000000
}
```

**成功响应**（200）：

```json
{
  "video_list": [
    {
      "id": 1,
      "author": {"id": 1, "username": "user1"},
      "title": "视频标题",
      "description": "描述",
      "play_url": "...",
      "cover_url": "...",
      "create_time": 1700000010,
      "likes_count": 5,
      "is_liked": false
    }
  ],
  "next_time": 1700000000,
  "has_more": true
}
```

| 响应字段 | 类型 | 说明 |
|---------|------|------|
| video_list | array | 视频列表 |
| next_time | int64 | 下一页游标时间戳，0表示无更多数据 |
| has_more | bool | 是否还有更多数据 |

**技术实现**：
- **冷热分离**：查询Redis ZSET `feed:global_timeline` 的watermark，热数据走Redis，冷数据走MySQL
- **三级缓存**：L1本地缓存(3s) → L2 Redis(1h) → L3 MySQL
- **singleflight**：防止并发重复查询DB

---

### 按点赞数排序Feed

获取按点赞数倒序排列的视频，使用双字段复合游标分页。

```
POST /feed/listLikesCount
```

**认证**：软认证（SoftJWTAuth）

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 是 | 每页数量 |
| likes_count_before | int64 | 否 | 游标：上一页最后一条的点赞数，首次传0 |
| id_before | uint | 否 | 游标：上一页最后一条的视频ID，首次传0 |

**请求示例（首次请求）**：

```json
{
  "limit": 10,
  "likes_count_before": 0,
  "id_before": 0
}
```

**请求示例（翻页请求）**：

```json
{
  "limit": 10,
  "likes_count_before": 50,
  "id_before": 15
}
```

**成功响应**（200）：

```json
{
  "video_list": [...],
  "next_likes_count_before": 30,
  "next_id_before": 25,
  "has_more": true
}
```

| 响应字段 | 类型 | 说明 |
|---------|------|------|
| video_list | array | 视频列表 |
| next_likes_count_before | int64 | 下一页点赞数游标 |
| next_id_before | uint | 下一页ID游标 |
| has_more | bool | 是否还有更多数据 |

**技术实现**：
- **双字段复合游标**：`likes_count + id` 联合游标，解决点赞数相同时排序不稳定问题
- 纯MySQL查询，保证数据一致性

---

### 按热度排序Feed（热榜）

获取按热度排序的视频，使用Redis滑动窗口+快照分页。

```
POST /feed/listByPopularity
```

**认证**：软认证（SoftJWTAuth）

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 是 | 每页数量 |
| as_of | int64 | 否 | 快照时间戳，首次不传或传0，翻页时传响应中的as_of |
| offset | int | 否 | 快照内偏移量，首次传0，翻页时传响应中的next_offset |
| latest_popularity | int64 | 否 | MySQL降级游标：上一条的热度值 |
| latest_before | int64 | 否 | MySQL降级游标：上一条的时间戳 |
| latest_id_before | uint | 否 | MySQL降级游标：上一条的ID |

**请求示例（首次请求）**：

```json
{
  "limit": 10
}
```

**请求示例（Redis快照翻页）**：

```json
{
  "limit": 10,
  "as_of": 1700000060,
  "offset": 10
}
```

**成功响应**（200）：

```json
{
  "video_list": [...],
  "as_of": 1700000060,
  "next_offset": 20,
  "has_more": true
}
```

| 响应字段 | 类型 | 说明 |
|---------|------|------|
| video_list | array | 视频列表 |
| as_of | int64 | 快照版本号（分钟级时间戳），翻页时需原样传回 |
| next_offset | int | 快照内下一页偏移量 |
| has_more | bool | 是否还有更多数据 |

**技术实现**：
- **滑动窗口热榜**：每分钟一个ZSET（`hot:video:1m:{minute}`），查询时合并最近60个窗口
- **快照分页**：`ZUNIONSTORE` 生成快照，同一翻页会话复用快照，避免"榜单抖动"
- **降级策略**：Redis不可用时自动降级到MySQL三字段游标分页（popularity + create_time + id）

---

### 关注Feed流 🔒

获取当前用户关注的人发布的视频，按时间倒序排列。

```
POST /feed/listByFollowing
```

**认证**：硬认证（JWTAuth）

**请求参数**：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 是 | 每页数量 |
| latest_time | int64 | 否 | 游标时间戳，首次不传或传0 |

**请求示例**：

```json
{
  "limit": 10,
  "latest_time": 0
}
```

**成功响应**（200）：

```json
{
  "video_list": [...],
  "next_time": 1700000000,
  "has_more": true
}
```

**技术实现**：
- Redis缓存 + 分布式锁防击穿
- 缓存key：`feed:listByFollowing:limit={limit}:accountID={id}:before={timestamp}`

---

## 错误码说明

### 通用错误

| 错误信息 | HTTP状态码 | 说明 |
|---------|-----------|------|
| `missing authorization header` | 401 | 缺少Authorization请求头 |
| `invalid or expired token` | 401 | Token无效或已过期 |
| `token has been revoked` | 401 | Token已被撤销（用户已登出/改名/改密码） |
| `Too Many Requests` | 429 | 请求频率超限 |

### 账号模块错误

| 错误信息 | 说明 |
|---------|------|
| `username already exists` | 用户名已被注册 |
| `invalid username or password` | 用户名或密码错误 |
| `unsuccessfully password changed` | 修改密码失败（旧密码错误等） |

### 点赞模块错误

| 错误信息 | 说明 |
|---------|------|
| `already liked` | 已经点赞过该视频 |
| `not liked yet` | 尚未点赞该视频，无法取消 |

### 评论模块错误

| 错误信息 | 说明 |
|---------|------|
| `not the author` | 非评论作者，无法删除 |

### 社交模块错误

| 错误信息 | 说明 |
|---------|------|
| `already followed` | 已经关注过该用户 |
| `not followed yet` | 尚未关注该用户，无法取关 |
| `cannot follow yourself` | 不能关注自己 |
| `account not found` | 目标用户不存在 |

---

## Postman测试集合

项目提供了完整的Postman测试集合，位于 `test/postman.json`。

### 使用方法

1. 打开Postman，点击 `Import`
2. 选择 `test/postman.json` 文件导入
3. 按以下顺序执行接口测试：
   1. `Account/Register` — 注册账号
   2. `Account/Login` — 登录获取Token（自动保存到变量）
   3. `Account/Find By Username` — 查询用户ID（自动保存到变量）
   4. `Video/Publish` — 发布视频（自动保存videoId）
   5. `Like/*` — 点赞相关操作
   6. `Comment/*` — 评论相关操作
   7. `Social/*` — 社交相关操作
   8. `Feed/*` — Feed流操作

### 预置变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| host | http://localhost:8080 | API服务地址 |
| username | 自动生成 | 测试用户名 |
| password | pass123 | 测试密码 |
| jwt_token | 空 | 登录后自动填充 |
| accountId | 1 | 用户ID |
| videoId | 1 | 视频ID |
| feedLimit | 10 | Feed每页数量 |

---

<div align="center">

**如有疑问，请提交 [Issue](https://github.com/LeoninCS/feedsystem_video_go/issues)**

</div>

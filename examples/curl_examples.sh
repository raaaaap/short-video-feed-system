#!/usr/bin/env bash
#
# 短视频Feed流系统 - curl接口调用示例
#
# 使用方法：
#   bash examples/curl_examples.sh
#
# 前置条件：
#   - 后端服务已启动（http://localhost:8080）
#   - 已安装curl和jq
#

set -euo pipefail

BASE_URL="http://localhost:8080"
USERNAME="demo_user_$(date +%s)"
PASSWORD="demo_password"

echo "=== 短视频Feed流系统 - curl接口调用示例 ==="
echo ""

# 步骤1：注册账号
echo "【步骤1】注册新账号"
echo "命令：curl -X POST ${BASE_URL}/account/register -H 'Content-Type: application/json' -d '{\"username\":\"${USERNAME}\",\"password\":\"${PASSWORD}\"}'"
REGISTER_RESULT=$(curl -s -X POST "${BASE_URL}/account/register" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"${USERNAME}\",\"password\":\"${PASSWORD}\"}")
echo "响应：${REGISTER_RESULT}"
echo ""

# 步骤2：登录获取Token
echo "【步骤2】登录获取Token"
echo "命令：curl -X POST ${BASE_URL}/account/login ..."
LOGIN_RESULT=$(curl -s -X POST "${BASE_URL}/account/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"${USERNAME}\",\"password\":\"${PASSWORD}\"}")
TOKEN=$(echo "${LOGIN_RESULT}" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "响应：Token获取${TOKEN:+成功}${TOKEN:-失败}"
echo ""

# 步骤3：查询用户信息
echo "【步骤3】查询用户信息"
FIND_RESULT=$(curl -s -X POST "${BASE_URL}/account/findByUsername" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"${USERNAME}\"}")
ACCOUNT_ID=$(echo "${FIND_RESULT}" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo "响应：用户ID = ${ACCOUNT_ID:-未找到}"
echo ""

# 步骤4：发布视频
echo "【步骤4】发布视频"
PUBLISH_RESULT=$(curl -s -X POST "${BASE_URL}/video/publish" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d "{\"title\":\"示例视频\",\"description\":\"curl测试视频\",\"play_url\":\"http://example.com/demo.mp4\",\"cover_url\":\"http://example.com/demo.jpg\"}")
VIDEO_ID=$(echo "${PUBLISH_RESULT}" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo "响应：视频ID = ${VIDEO_ID:-未创建}"
echo ""

# 步骤5：点赞视频
echo "【步骤5】点赞视频"
LIKE_RESULT=$(curl -s -X POST "${BASE_URL}/like/like" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d "{\"video_id\":${VIDEO_ID:-1}}")
echo "响应：${LIKE_RESULT}"
echo ""

# 步骤6：查询点赞状态
echo "【步骤6】查询是否已点赞"
IS_LIKED_RESULT=$(curl -s -X POST "${BASE_URL}/like/isLiked" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d "{\"video_id\":${VIDEO_ID:-1}}")
echo "响应：${IS_LIKED_RESULT}"
echo ""

# 步骤7：发布评论
echo "【步骤7】发布评论"
COMMENT_RESULT=$(curl -s -X POST "${BASE_URL}/comment/publish" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d "{\"video_id\":${VIDEO_ID:-1},\"content\":\"这个视频太棒了！\"}")
echo "响应：${COMMENT_RESULT}"
echo ""

# 步骤8：查看评论列表
echo "【步骤8】查看评论列表"
COMMENTS_RESULT=$(curl -s -X POST "${BASE_URL}/comment/listAll" \
  -H "Content-Type: application/json" \
  -d "{\"video_id\":${VIDEO_ID:-1}}")
echo "响应：${COMMENTS_RESULT}"
echo ""

# 步骤9：浏览最新Feed
echo "【步骤9】浏览最新Feed流"
FEED_RESULT=$(curl -s -X POST "${BASE_URL}/feed/listLatest" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{"limit":5,"latest_time":0}')
echo "响应：${FEED_RESULT}" | head -c 200
echo "..."
echo ""

# 步骤10：浏览热榜
echo "【步骤10】浏览热度排行"
POPULARITY_RESULT=$(curl -s -X POST "${BASE_URL}/feed/listByPopularity" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{"limit":5}')
echo "响应：${POPULARITY_RESULT}" | head -c 200
echo "..."
echo ""

# 步骤11：取消点赞
echo "【步骤11】取消点赞"
UNLIKE_RESULT=$(curl -s -X POST "${BASE_URL}/like/unlike" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d "{\"video_id\":${VIDEO_ID:-1}}")
echo "响应：${UNLIKE_RESULT}"
echo ""

# 步骤12：登出
echo "【步骤12】登出"
LOGOUT_RESULT=$(curl -s -X POST "${BASE_URL}/account/logout" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{}')
echo "响应：${LOGOUT_RESULT}"
echo ""

echo "=== 所有示例执行完毕 ==="

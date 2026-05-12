package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const baseURL = "http://localhost:8080"

var client = &http.Client{Timeout: 30 * time.Second}

func main() {
	fmt.Println("=== 短视频Feed流系统 - 功能使用示例 ===")
	fmt.Println()

	// 步骤1：注册账号
	fmt.Println("【步骤1】注册新账号")
	username := fmt.Sprintf("demo_user_%d", time.Now().Unix())
	password := "demo_password"
	register(username, password)
	fmt.Println()

	// 步骤2：登录获取Token
	fmt.Println("【步骤2】登录获取Token")
	token := login(username, password)
	fmt.Printf("获取到Token: %s...\n", truncate(token, 30))
	fmt.Println()

	// 步骤3：查询用户信息
	fmt.Println("【步骤3】查询用户信息")
	accountID := findByUsername(username)
	fmt.Printf("用户ID: %d\n", accountID)
	fmt.Println()

	// 步骤4：发布视频
	fmt.Println("【步骤4】发布视频")
	videoID := publishVideo(token, "示例视频标题", "这是一个示例视频的描述")
	fmt.Printf("视频ID: %d\n", videoID)
	fmt.Println()

	// 步骤5：点赞视频
	fmt.Println("【步骤5】点赞视频")
	likeVideo(token, videoID)
	fmt.Println()

	// 步骤6：查询点赞状态
	fmt.Println("【步骤6】查询是否已点赞")
	isLiked := checkLiked(token, videoID)
	fmt.Printf("是否已点赞: %v\n", isLiked)
	fmt.Println()

	// 步骤7：发布评论
	fmt.Println("【步骤7】发布评论")
	publishComment(token, videoID, "这个视频太棒了！")
	fmt.Println()

	// 步骤8：查看评论列表
	fmt.Println("【步骤8】查看评论列表")
	listComments(videoID)
	fmt.Println()

	// 步骤9：浏览最新Feed
	fmt.Println("【步骤9】浏览最新Feed流")
	listLatestFeed(token, 5)
	fmt.Println()

	// 步骤10：浏览热榜
	fmt.Println("【步骤10】浏览热度排行")
	listPopularityFeed(token, 5)
	fmt.Println()

	// 步骤11：取消点赞
	fmt.Println("【步骤11】取消点赞")
	unlikeVideo(token, videoID)
	fmt.Println()

	// 步骤12：登出
	fmt.Println("【步骤12】登出")
	logout(token)
	fmt.Println()

	fmt.Println("=== 所有示例执行完毕 ===")
}

func register(username, password string) {
	body := map[string]string{"username": username, "password": password}
	resp := postJSON("/account/register", body, "")
	result := parseResponse(resp)
	if msg, ok := result["message"]; ok {
		fmt.Printf("注册结果: %v\n", msg)
	} else if err, ok := result["error"]; ok {
		fmt.Printf("注册失败: %v\n", err)
	}
}

func login(username, password string) string {
	body := map[string]string{"username": username, "password": password}
	resp := postJSON("/account/login", body, "")
	result := parseResponse(resp)
	if token, ok := result["token"].(string); ok {
		return token
	}
	fmt.Println("登录失败，无法获取Token")
	return ""
}

func findByUsername(username string) float64 {
	body := map[string]string{"username": username}
	resp := postJSON("/account/findByUsername", body, "")
	result := parseResponse(resp)
	if id, ok := result["id"].(float64); ok {
		return id
	}
	return 0
}

func publishVideo(token, title, description string) float64 {
	body := map[string]string{
		"title":       title,
		"description": description,
		"play_url":    "http://example.com/demo.mp4",
		"cover_url":   "http://example.com/demo.jpg",
	}
	resp := postJSON("/video/publish", body, token)
	result := parseResponse(resp)
	if id, ok := result["id"].(float64); ok {
		return id
	}
	return 0
}

func likeVideo(token string, videoID float64) {
	body := map[string]interface{}{"video_id": videoID}
	resp := postJSON("/like/like", body, token)
	result := parseResponse(resp)
	if msg, ok := result["message"]; ok {
		fmt.Printf("点赞结果: %v\n", msg)
	}
}

func unlikeVideo(token string, videoID float64) {
	body := map[string]interface{}{"video_id": videoID}
	resp := postJSON("/like/unlike", body, token)
	result := parseResponse(resp)
	if msg, ok := result["message"]; ok {
		fmt.Printf("取消点赞结果: %v\n", msg)
	}
}

func checkLiked(token string, videoID float64) bool {
	body := map[string]interface{}{"video_id": videoID}
	resp := postJSON("/like/isLiked", body, token)
	result := parseResponse(resp)
	if liked, ok := result["is_liked"].(bool); ok {
		return liked
	}
	return false
}

func publishComment(token string, videoID float64, content string) {
	body := map[string]interface{}{
		"video_id": videoID,
		"content":  content,
	}
	resp := postJSON("/comment/publish", body, token)
	result := parseResponse(resp)
	if msg, ok := result["message"]; ok {
		fmt.Printf("评论发布结果: %v\n", msg)
	}
}

func listComments(videoID float64) {
	body := map[string]interface{}{"video_id": videoID}
	resp := postJSON("/comment/listAll", body, "")
	var comments []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&comments)
	for i, c := range comments {
		fmt.Printf("  评论%d: [%v] %v\n", i+1, c["username"], c["content"])
	}
}

func listLatestFeed(token string, limit int) {
	body := map[string]interface{}{"limit": limit, "latest_time": 0}
	resp := postJSON("/feed/listLatest", body, token)
	result := parseResponse(resp)
	if videoList, ok := result["video_list"].([]interface{}); ok {
		for i, v := range videoList {
			if video, ok := v.(map[string]interface{}); ok {
				title, _ := video["title"].(string)
				likesCount, _ := video["likes_count"].(float64)
				fmt.Printf("  视频%d: %s (点赞数: %.0f)\n", i+1, title, likesCount)
			}
		}
		if hasMore, ok := result["has_more"].(bool); ok {
			fmt.Printf("  是否有更多: %v\n", hasMore)
		}
	}
}

func listPopularityFeed(token string, limit int) {
	body := map[string]interface{}{"limit": limit}
	resp := postJSON("/feed/listByPopularity", body, token)
	result := parseResponse(resp)
	if videoList, ok := result["video_list"].([]interface{}); ok {
		for i, v := range videoList {
			if video, ok := v.(map[string]interface{}); ok {
				title, _ := video["title"].(string)
				likesCount, _ := video["likes_count"].(float64)
				fmt.Printf("  热榜%d: %s (点赞数: %.0f)\n", i+1, title, likesCount)
			}
		}
	}
}

func logout(token string) {
	body := map[string]interface{}{}
	resp := postJSON("/account/logout", body, token)
	result := parseResponse(resp)
	if msg, ok := result["message"]; ok {
		fmt.Printf("登出结果: %v\n", msg)
	}
}

func uploadVideo(token, filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return ""
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Printf("创建表单文件失败: %v\n", err)
		return ""
	}
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", baseURL+"/video/uploadVideo", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	result := parseResponse(resp)
	if playURL, ok := result["play_url"].(string); ok {
		return playURL
	}
	return ""
}

func postJSON(path string, body interface{}, token string) *http.Response {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		os.Exit(1)
	}
	return resp
}

func parseResponse(resp *http.Response) map[string]interface{} {
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

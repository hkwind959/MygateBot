package api

import (
	"MygateBot/constant"
	"MygateBot/logs"
	"MygateBot/model"
	"MygateBot/utils"
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// GetUserNode 获取用户节点
func GetUserNode(token string, proxy string) model.UserNodeResp {
	maxRetries := 5

	client := utils.NewHttpClient(proxy)
	for retries := 0; retries < maxRetries; retries++ {
		// 复制 constant.Headers
		headers := make(map[string]string)
		for k, v := range constant.Headers {
			headers[k] = v
		}
		// 添加 Authorization 头
		headers["Authorization"] = "Bearer " + token

		resp, err := client.Get(constant.GetNodeUrl, nil, headers, nil)
		if err == nil {
			var userNodeResp model.UserNodeResp
			_ = json.Unmarshal(resp.Body(), &userNodeResp)
			return userNodeResp
		}

		// 记录更详细的日志信息
		logs.I().Error("GetUserNode 获取信息出错",
			zap.Error(err),
			zap.Int("retry", retries+1),
			zap.String("proxy", proxy),
			zap.String("token_prefix", token[:7]+"...")) // 只记录 token 前缀
	}
	return model.UserNodeResp{}
}

// RegisterNode 注册节点
func RegisterNode(token string, proxy string, node string) model.RegisterNodeResp {
	maxRetries := 5
	var retries int
	var uuidStr string

	if node == "" {
		uuidObj, err := uuid.NewRandom()
		if err != nil {
			logs.I().Error("生成UUID失败", zap.Error(err))
		}
		uuidStr = uuidObj.String()
	} else {
		uuidStr = node
	}
	payload := map[string]interface{}{
		"id":             uuidStr,
		"status":         "Good",
		"activationDate": time.Now().UTC().Format(time.RFC3339),
	}
	client := utils.NewHttpClient(proxy)
	for retries < maxRetries {
		// 复制 constant.Headers
		headers := make(map[string]string)
		for k, v := range constant.Headers {
			headers[k] = v
		}
		// 添加 Authorization 头
		headers["Authorization"] = "Bearer " + token
		resp, err := client.Post(constant.RegisterNodeUrl, headers, payload, nil)
		if err == nil {
			var respBody model.RegisterNodeResp
			_ = json.Unmarshal(resp.Body(), &respBody)
			return respBody
		}
		logs.I().Error("注册节点时出错", zap.Error(err))
		retries++
		if retries < maxRetries {
			logs.I().Info("10秒后重试...")
			time.Sleep(10 * time.Second)
		} else {
			logs.I().Error("最大重试次数已超出; 放弃注册。")
			return model.RegisterNodeResp{}
		}
	}
	return model.RegisterNodeResp{}
}

// GetUserInfo 更新刷新token
func GetUserInfo(token string, proxy string) string {
	maxRetries := 5

	client := utils.NewHttpClient(proxy)
	for retries := 0; retries < maxRetries; retries++ {
		// 复制 constant.Headers
		headers := make(map[string]string)
		for k, v := range constant.Headers {
			headers[k] = v
		}
		// 添加 Authorization 头
		headers["Authorization"] = "Bearer " + token

		resp, err := client.Get(constant.GetUserInfoUrl, nil, headers, nil)
		if err == nil {
			return resp.String()
		}
		// 记录更详细的日志信息
		logs.I().Error("GetUserNode 获取信息出错",
			zap.Error(err),
			zap.Int("retry", retries+1),
			zap.String("proxy", proxy),
			zap.String("token_prefix", token[:7]+"...")) // 只记录 token 前缀
	}
	return ""
}

package utils

import (
	"MygateBot/constant"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

func GenerateSignature(data interface{}) (map[string]string, error) {
	// 获取当前 UTC 时间
	now := time.Now().UTC()

	// 构造时间戳（以毫秒为单位）
	timestamp := now.UnixMilli()
	// 将数据序列化为 JSON 字符串
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("无法序列化数据: %v", err)
	}
	// 创建 HMAC 签名

	h := hmac.New(sha256.New, []byte(constant.SecretKey))
	_, err = h.Write(append(dataJSON, fmt.Sprintf("%d", timestamp)...)) // 将数据和时间戳拼接
	if err != nil {
		return nil, fmt.Errorf("无法生成HMAC签名: %v", err)
	}
	// 获取签名的十六进制表示
	signature := fmt.Sprintf("%x", h.Sum(nil))

	// 返回结果
	return map[string]string{
		"timestamp": fmt.Sprintf("%d", timestamp),
		"signature": signature,
	}, nil
}

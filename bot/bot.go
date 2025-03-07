package bot

import (
	"MygateBot/api"
	"MygateBot/constant"
	"MygateBot/logs"
	"MygateBot/model"
	"MygateBot/utils"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Bot struct {
	Token     string
	Proxy     string
	NodeId    string
	lock      sync.Mutex
	writeChan chan interface{}
	conn      *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc
	Remark    string
	Ip        string
}

// NewBot 创建机器人
func NewBot(t model.TokenRequest) *Bot {
	ctx, cancel := context.WithCancel(context.Background())
	return &Bot{
		Token:     t.Token,
		Proxy:     t.Proxy,
		NodeId:    t.NodeId,
		Remark:    t.Remark,
		Ip:        t.Ip,
		writeChan: make(chan interface{}, 100),
		conn:      nil,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// StartBot 启动机器人
func (b *Bot) StartBot() error {
	wsClient, err := utils.NewWebSocketProxyClient(b.Proxy)
	if err != nil {
		logs.I().Error("创建WebSocket客户端失败", zap.Error(err))
		return err
	}
	header := http.Header{
		"Accept-encoding": {"gzip, deflate, br, zstd"},
		"Accept-language": {"en-US,en;q=0.9,id;q=0.8"},
		"Cache-control":   {"no-cache"},
		"Host":            {"api.mygate.network"},
		"Origin":          {"chrome-extension://hajiimgolngmlbglaoheacnejbnnmoco"},
		"Pragma":          {"no-cache"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"},
	}

	signature, err := utils.GenerateSignature(map[string]string{
		"nodeId": b.NodeId,
	})
	if err != nil {
		logs.I().Error("生成签名失败: ", zap.Error(err))
		return err
	}

	wsUrl := constant.GetWssUrl + b.NodeId + "&signature=" + signature["signature"] + "&timestamp=" + signature["timestamp"] + "&version=2&EIO=4&transport=websocket"

	err = wsClient.Connect(wsUrl, header)
	if err != nil {
		maxRetries := 150
		backoff := time.Second * 5    // 初始重试间隔
		maxBackoff := time.Minute * 5 // 最大重试间隔
		for i := 0; i < maxRetries; i++ {
			time.Sleep(backoff)
			backoff *= 2 // 每次重试间隔翻倍
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			// 添加随机因子以减少冲突
			jitter := time.Duration(rand.Int63n(int64(backoff)))
			time.Sleep(jitter)

			err = wsClient.Connect(wsUrl, header)
			if err == nil {
				b.lock.Lock()
				b.conn = wsClient.GetConn()
				b.lock.Unlock()
				break
			}
			select {
			case <-b.ctx.Done():
				return b.ctx.Err()
			default:
			}
			logs.I().Info("重试连接中... ", zap.Int("重试次数", i+1), zap.Duration("重试间隔", backoff))
		}
		if err != nil {
			logs.I().Error("WebSocket客户端连接失败: ", zap.String("User", b.Remark), zap.String("nodeId", b.NodeId), zap.String("Ip", b.Ip))
			return err
		}
		if b.conn != nil {
			b.writeChan <- fmt.Sprintf(`40{"token":"Bearer %s"}`, b.Token)
		}
		logs.I().Info("WebSocket客户端连接成功: ", zap.String("User", b.Remark), zap.String("Ip", b.Ip), zap.String("NodeId", b.NodeId))
		// 发送消息
		go b.writeMessage()
		// 接收消息
		go b.receiveMessage()
	}
	return nil
}

// 发送消息函数
func (b *Bot) writeMessage() {
	for {
		msg, ok := <-b.writeChan
		if !ok {
			return
		}
		b.lock.Lock()
		if b.conn == nil {
			b.lock.Unlock()
			continue
		}
		err := b.conn.WriteMessage(websocket.TextMessage, []byte(msg.(string)))
		b.lock.Unlock()
		if err != nil {
			logs.I().Error("发送消息失败: ", zap.Error(err))
			return
		}
		logs.I().Info("发送消息: ", zap.String("User", b.Remark), zap.String("NodeId", b.NodeId), zap.String("Ip", b.Ip), zap.Any("message", msg))
	}
}

// 接收消息函数
func (b *Bot) receiveMessage() {
	for {
		conn := b.conn
		if conn == nil {
			logs.I().Error("连接未初始化，等待重新连接")
			time.Sleep(time.Second) // 等待一段时间后重试
			continue
		}
		_, message, err := conn.ReadMessage()
		logs.I().Info("收到消息: ", zap.String("User", b.Remark), zap.String("NodeId", b.NodeId), zap.String("Ip", b.Ip), zap.String("message", string(message)))
		if err != nil {
			logs.I().Error("读取消息失败: ", zap.Error(err))
			_ = conn.Close()
			b.conn = nil
			// 重新连接
			b.reconnect()
			return
		}
		s := string(message)
		if s == "2" || s == "41" {
			b.writeChan <- "3"
		}
	}
}

// 重新连接
func (b *Bot) reconnect() {
	for {
		select {
		case <-b.ctx.Done():
			return
		default:
		}
		if err := b.StartBot(); err == nil {
			logs.I().Info("重新连接成功")
			return
		}
		logs.I().Error("重新连接失败 :", zap.String("User", b.Remark), zap.String("IP", b.Ip), zap.String("NodeId", b.NodeId))
		backoff := time.Second * 5    // 初始重试间隔
		maxBackoff := time.Minute * 5 // 最大重试间隔
		for i := 0; i < 150; i++ {
			time.Sleep(backoff)
			backoff *= 2 // 每次重试间隔翻倍
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			// 添加随机因子以减少冲突
			jitter := time.Duration(rand.Int63n(int64(backoff)))
			time.Sleep(jitter)

			if err := b.StartBot(); err == nil {
				logs.I().Info("重新连接成功")
				return
			}
			logs.I().Error("重新连接失败", zap.Int("重试次数", i+1))
		}
		logs.I().Error("达到最大重试次数，无法重新连接")
	}
}

// GetUserInfo 获取用户信息
func GetUserInfo(token string, proxy string) {
	info := api.GetUserInfo(token, proxy)
	logs.I().Info("获取用户信息成功: ", zap.Any("userInfo", info))
}

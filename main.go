package main

import (
	"MygateBot/api"
	"MygateBot/bot"
	"MygateBot/config"
	"MygateBot/logs"
	"MygateBot/model"
	"context"
	"go.uber.org/zap"
	"sync"
	"time"
)

func init() {
}

func main() {
	tokens := config.GetTokens()
	if len(tokens) == 0 {
		logs.I().Info("没有配置Token信息")
		return
	}
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, token := range tokens {
		// 获取用户节点
		nodes := api.GetUserNode(token.Token, token.Proxies[0])
		logs.I().Info("获取用户节点成功 ：", zap.Any("token", token.Remark), zap.Any("nodes", nodes))
		// 注册节点
		for _, proxy := range token.Proxies {
			// 注册节点
			node := api.RegisterNode(token.Token, proxy, "")
			logs.I().Info("注册节点成功 ：", zap.Any("node", node))

			for _, node := range nodes.Data.Item {
				wg.Add(1)
				go func(nodeID string) {
					defer wg.Done()
					b := bot.NewBot(token.Token, proxy, nodeID)
					if err := b.StartBot(); err != nil {
						logs.I().Error("启动机器人失败", zap.Error(err))
					}
					// 监听上下文取消信号
					select {
					case <-ctx.Done():
						logs.I().Info("上下文被取消，停止机器人", zap.String("nodeID", nodeID))
					}
				}(node.ID)
			}
		}
	}
	// 设置时间间隔为 11 分钟
	interval := 11 * time.Minute
	// 创建 Ticker
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	logs.I().Info("定时任务启动...")
	// 启动定时任务
	go func() {
		for {
			select {
			case <-ctx.Done():
				logs.I().Info("定时任务被取消")
				return
			case <-ticker.C:
				logs.I().Info("定时任务触发")
				var taskWg sync.WaitGroup
				for _, token := range tokens {
					taskWg.Add(1)
					go func(ctx context.Context, token model.Token) {
						defer taskWg.Done()
						bot.GetUserInfo(token.Token, token.Proxies[0])
						// 监听上下文取消信号
						select {
						case <-ctx.Done():
							logs.I().Info("上下文被取消，停止获取用户信息", zap.String("token", token.Token))
						default:
						}
					}(ctx, token)
				}
				taskWg.Wait()
			}
		}
	}()

	wg.Wait()
	select {}
}

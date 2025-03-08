package main

import (
	"MygateBot/api"
	"MygateBot/bot"
	"MygateBot/config"
	"MygateBot/logs"
	"MygateBot/model"
	"go.uber.org/zap"
)

func init() {
}

func main() {
	tokens := config.GetTokens()
	if len(tokens) == 0 {
		logs.I().Info("没有配置Token信息")
		return
	}
	//var wg sync.WaitGroup
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	var tokenList []model.TokenRequest
	for _, token := range tokens {
		var tokenModel model.TokenRequest
		// 获取用户节点
		nodes := api.GetUserNode(token.Token, token.Proxies[0])
		logs.I().Info("获取用户节点成功 ：", zap.Any("token", token.Remark), zap.Any("nodes", nodes))
		// 注册节点
		for _, proxy := range token.Proxies {
			checkProxy := api.CheckProxy(proxy)
			logs.I().Info("检查代理 ：", zap.Any("checkProxy", checkProxy["ip"]))
			tokenModel.Token = token.Token
			tokenModel.Remark = token.Remark
			tokenModel.Proxy = proxy
			tokenModel.Ip = checkProxy["ip"]
			// 匹配 IP 并赋值 NodeId
			for _, node := range nodes.Data.Item {
				if node.IP == tokenModel.Ip {
					tokenModel.NodeId = node.ID
					break
				}
			}
			tokenList = append(tokenList, tokenModel)
		}
	}
	logs.I().Info("封装用户节点成功 ：", zap.Any("tokenList", tokenList))

	// 注册 节点，启动机器人
	for _, t := range tokenList {
		// 注册节点
		nodeResp := api.RegisterNode(t.Token, t.Proxy, t.NodeId)
		logs.I().Info("注册节点成功 ：", zap.Any("UserMail", t.Remark), zap.Any("Ip", t.Ip), zap.Any("NodeId", t.NodeId), zap.Any("Resp", nodeResp))
		b := bot.NewBot(t)
		go func() {
			err := b.StartBot()
			if err != nil {
				logs.I().Info("启动机器人失败", zap.Error(err))
			}
			logs.I().Info("启动机器人成功")
		}()
	}
	//}
	// 设置时间间隔为 11 分钟
	//interval := 11 * time.Minute
	//// 创建 Ticker
	//ticker := time.NewTicker(interval)
	//defer ticker.Stop()
	//logs.I().Info("定时任务启动...")
	//// 启动定时任务
	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			logs.I().Info("定时任务被取消")
	//			return
	//		case <-ticker.C:
	//			logs.I().Info("定时任务触发")
	//			var taskWg sync.WaitGroup
	//			for _, token := range tokens {
	//				taskWg.Add(1)
	//				go func(ctx context.Context, token model.Token) {
	//					defer taskWg.Done()
	//					bot.GetUserInfo(token.Token, token.Proxies[0])
	//					// 监听上下文取消信号
	//					select {
	//					case <-ctx.Done():
	//						logs.I().Info("上下文被取消，停止获取用户信息", zap.String("token", token.Token))
	//					default:
	//					}
	//				}(ctx, token)
	//			}
	//			taskWg.Wait()
	//		}
	//	}
	//}()
	select {}
}

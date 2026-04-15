package main

import (
	"bot/config"
	"bot/internal/routes"
	"bot/models"
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

func main() {
	ctx := context.Background()
	log.Println("加载配置")
	key, model, baseURL := config.LoadModelConfig()
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  key,
		Model:   model,
		BaseURL: baseURL,
	})
	if err != nil {
		log.Fatal("创建模型失败", err)
		return
	}
	register := models.NewSessionRegister(chatModel)
	pool := models.NewWorkPool(4, 20, register)
	//agent, _ := adk.NewChatModelAgent(ctx,&adk.ChatModelAgentConfig{})
	go pool.Start()
	log.Println("启动 HTTP 服务器")
	routes.Run(ctx, pool)
}

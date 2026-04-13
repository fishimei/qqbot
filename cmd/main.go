package main

import (
	"bot/config"
	"bot/internal/routes"
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	log.Println("加载配置")
	key, model, baseURL, systemPrompt := config.LoadModelConfig()
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  key,
		Model:   model,
		BaseURL: baseURL,
	})
	//agent, _ := adk.NewChatModelAgent(ctx,&adk.ChatModelAgentConfig{})
	if err != nil {
		log.Fatal(err)
		return
	}
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
	}
	log.Println("启动 HTTP 服务器")
	routes.Run(chatModel, ctx, messages)
}

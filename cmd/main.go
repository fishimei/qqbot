package main

import (
	"bot/config"
	"bot/internal/routes"
	"bot/models"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ctx, cancel := context.WithCancel(ctx)
	register := models.NewSessionRegister(chatModel)
	pool := models.NewWorkPool(4, 20, register)
	//agent, _ := adk.NewChatModelAgent(ctx,&adk.ChatModelAgentConfig{})
	pool.Start(ctx)
	log.Println("启动 HTTP 服务器")
	go routes.Run(pool)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	cancel() // 通知 worker

	// 用 WaitGroup 等 worker 退出
	done := make(chan struct{})
	go func() {
		pool.WaitGroup.Wait() // 等待所有 worker 退出
		close(done)
	}()

	select {
	case <-done:
		log.Println("所有 worker 已退出")
	case <-time.After(30 * time.Second):
		log.Println("worker 退出超时")
	}
}

package routes

import (
	config2 "bot/config"
	"bot/internal/hook"
	"bot/internal/sendMsg"
	"bot/models"
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context, pool *models.WorkPool) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	_, expectedToken, _ := config2.LoadNapcatConfig()
	judgeEnable := config2.LoadJudgeAtConfig()
	r.Use(
		hook.CheckToken(expectedToken),
		hook.ParseMsg(),
		hook.JudgeAt(judgeEnable),
	)
	r.POST("/", sendMsg.SendMsg(ctx, pool))
	port := config2.LoadServerConfig()
	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
		return
	}
}

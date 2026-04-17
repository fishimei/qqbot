package routes

import (
	"bot/config"
	"bot/internal/hook"
	"bot/internal/sendMsg"
	"bot/models"
	"log"

	"github.com/gin-gonic/gin"
)

func Run(pool *models.WorkPool) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	judgeEnable := config.LoadJudgeAtConfig()
	r.Use(
		hook.ParseMsg(),
		hook.JudgeAt(judgeEnable),
	)
	r.POST("/", sendMsg.SendMsg(pool))
	port := config.LoadServerConfig()
	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
		return
	}
}

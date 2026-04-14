package routes

import (
	config2 "bot/config"
	"bot/internal/hook"
	"bot/internal/sendMsg"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

func Run(chatModel *ark.ChatModel, ctx context.Context, messages []*schema.Message) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	apiBaseURL, expectedToken, authToken := config2.LoadNapcatConfig()
	judgeEnable := config2.LoadJudgeAtConfig()
	r.Use(
		hook.PassToken(apiBaseURL, authToken),
		hook.CheckToken(expectedToken),
		hook.ParseMsg(),
		hook.JudgeAt(judgeEnable),
	)
	r.POST("/", sendMsg.SendMsg(chatModel, ctx, messages))
	r.POST("/test", func(c *gin.Context) {
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return
		}
		log.Println(string(all))
		c.JSON(http.StatusOK, gin.H{})
	})
	port := config2.LoadServerConfig()
	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
		return
	}
}

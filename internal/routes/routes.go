package routes

import (
	config2 "bot/config"
	"context"
	"log"
	"net/http"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type Event struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	UserID      int64  `json:"user_id"`
	GroupID     int64  `json:"group_id"`
	RawMessage  string `json:"raw_message"`
	SelfID      int64  `json:"self_id"`
}

type SendMsgReq struct {
	MessageType string `json:"message_type"`
	UserID      int64  `json:"user_id,omitempty"`
	GroupID     int64  `json:"group_id,omitempty"`
	Message     string `json:"message"`
}

func Run(chatModel *ark.ChatModel, ctx context.Context, messages []*schema.Message) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	_, expectedToken := config2.LoadNapcatConfig()
	port := config2.LoadServerConfig()
	r.POST("/", func(c *gin.Context) {
		// 验证 Token
		if expectedToken != "" {
			token := c.GetHeader("Authorization")
			if token != "Bearer "+expectedToken {
				c.String(http.StatusUnauthorized, "invalid token")
				return
			}
		}
		// 解析事件数据
		var event Event
		err := c.BindJSON(&event)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		// 处理消息事件
		if event.PostType != "message" {
			c.String(http.StatusOK, "ignored")
			return
		}
		var reply string
		stream, err := chatModel.Stream(ctx, messages)
		if err != nil {
			return
		}
		defer stream.Close()
		for {
			chunk, err := stream.Recv()
			if err != nil {
				break
			}
			reply += chunk.Content
		}
		c.JSON(http.StatusOK, gin.H{"reply": reply})
	})
	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
		return
	}
}

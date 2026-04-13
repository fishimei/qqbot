package sendMsg

import (
	"bot/internal/utils"
	"bot/models"
	"context"
	"log"
	"net/http"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

func SendMsg(chatModel *ark.ChatModel, ctx context.Context, messages []*schema.Message) func(c *gin.Context) {
	return func(c *gin.Context) {
		event, _ := utils.GetEvent(c)
		var reply string
		reqMessages := make([]*schema.Message, 0, len(messages)+1)
		reqMessages = append(reqMessages, messages...)
		reqMessages = append(reqMessages, StructRequest(event.Message))
		stream, err := chatModel.Stream(ctx, reqMessages)
		if err != nil {
			return
		}
		defer stream.Close()
		for {
			chunk, err := stream.Recv()
			if err != nil {
				log.Println(err)
				break
			}
			reply += chunk.Content
		}
		c.JSON(http.StatusOK, gin.H{"reply": reply})
	}
}

func StructRequest(message []models.MetaMessage) *schema.Message {
	userMessage := &schema.Message{Role: schema.User}
	for _, meta := range message {
		switch meta.Type {
		case "text":
			text, _ := meta.Data["text"].(string)
			userMessage.UserInputMultiContent = append(userMessage.UserInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeText,
				Text: text,
			})
		case "image":
			url, _ := meta.Data["url"].(schema.ImageURLDetail)
			userMessage.UserInputMultiContent = append(userMessage.UserInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					Detail: url,
				},
			})
		default:
			continue
		}
	}
	return userMessage
}

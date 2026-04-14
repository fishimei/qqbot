package sendMsg

import (
	"bot/internal/utils"
	"bot/models"
	"bytes"
	"context"
	"encoding/json"
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
		//每次复制其实挺影响性能的
		reqMessages := make([]*schema.Message, 0, len(messages)+1)
		reqMessages = append(reqMessages, messages...)
		reqMessages = append(reqMessages, StructRequest(event.Message))
		c.JSON(http.StatusOK, gin.H{})
		stream, err := chatModel.Stream(ctx, reqMessages)
		if err != nil {
			log.Println("chatModel stream failed", err)
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
		apiBaseURL, _ := utils.GetApiBaseURL(c)
		authToken, _ := utils.GetAuthToken(c)
		sendMsgToNapcat(event.MessageType, event.GroupID, event.UserID, reply, apiBaseURL, authToken)
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

// storeAssistantReply TODO
func storeAssistantReply(groupID int64, reply string) {

}

func sendMsgToNapcat(msgType string, groupID int64, userID int64, message string, apiBaseURL string, authToken string) {
	reqBody := models.SendMsgReq{
		MessageType: msgType,
		Message:     message,
	}
	if msgType == "private" {
		reqBody.UserID = userID
	} else if msgType == "group" {
		reqBody.GroupID = groupID
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("marshal send msg req failed", err)
		return
	}
	url := apiBaseURL + "/send_msg"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("create send msg req failed", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("send msg to napcat failed", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("send msg to napcat failed, status code:", resp.StatusCode)
		return
	}
}

//根据群号来进行一个会话的在内存中的储存，后续可以考虑持久化到数据库中
//为了保证会话的连续性，还要维护一个全局管道来进行消息的传入和分配

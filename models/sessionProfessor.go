package models

import (
	"bot/config"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

//处理传递来的消息，对接模型，发送最终消息

type SessionProfessor struct {
	session      *Session
	chatModel    *ark.ChatModel
	napcatConfig *NapcatConfig
	msgCh        chan *Event //接收消息通道
	useStatus    atomic.Bool
}

func NewSessionProfessor(sessionID string, chatModel *ark.ChatModel) *SessionProfessor {
	sp := &SessionProfessor{
		session:      NewSession(sessionID, 2, 60),
		msgCh:        make(chan *Event, 20),
		chatModel:    chatModel,
		napcatConfig: NewNapcatConfig(),
	}
	sp.useStatus.Store(false)
	systemPrompts := config.LoadSystemPrompts()
	err := sp.session.SetSystemMessages(schema.SystemMessage(systemPrompts))
	if err != nil {
		log.Println("set system messages is failed ", err)
		return nil
	}
	return sp
}

func (sp *SessionProfessor) Start() {
	if !sp.useStatus.Swap(true) {
		go func() {
			for msg := range sp.msgCh {
				sp.HandleEvent(msg)
			}
		}()
	}
}

func (sp *SessionProfessor) AppendEvent(event *Event) {
	sp.msgCh <- event
}

func (sp *SessionProfessor) HandleEvent(msg *Event) {
	ctx := context.Background()
	sp.session.Append(StructRequest(msg.Message))
	reply := sp.GetReply(ctx, sp.session.GetAll())
	sendMsgToNapcat(sp.session.SessionID, reply, sp.napcatConfig.ApiBaseURL, sp.napcatConfig.AuthToken)
}

func (sp *SessionProfessor) GetReply(ctx context.Context, messages []*schema.Message) string {
	stream, err := sp.chatModel.Stream(ctx, messages)
	if err != nil {
		log.Println("chatModel stream failed", err)
		return ""
	}
	reply := ""
	for {
		chunk, err := stream.Recv()
		if err != nil {
			break
		}
		reply += chunk.Content
	}
	sp.session.Append(&schema.Message{Role: schema.Assistant, AssistantGenMultiContent: []schema.MessageOutputPart{
		{
			Type: schema.ChatMessagePartTypeText,
			Text: reply,
		},
	}})
	return reply
}

func StructRequest(message []MetaMessage) *schema.Message {
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
			url, _ := meta.Data["url"].(string)
			userMessage.UserInputMultiContent = append(userMessage.UserInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{
						URL: &url,
					},
				},
			})
		default:
			continue
		}
	}
	return userMessage
}

func sendMsgToNapcat(sessionID string, message string, apiBaseURL string, authToken string) {
	mataDatas := strings.Split(sessionID, "_")
	reqBody := SendMsgReq{
		MessageType: mataDatas[0],
		Message:     message,
	}
	if mataDatas[0] == "private" {
		reqBody.UserID, _ = strconv.ParseInt(mataDatas[1], 10, 64)
	} else if mataDatas[0] == "group" {
		reqBody.GroupID, _ = strconv.ParseInt(mataDatas[1], 10, 64)
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
		log.Println("send msg to napcat failed, useStatus code:", resp.StatusCode)
		return
	}
}

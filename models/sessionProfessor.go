package models

import (
	"bot/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

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

func (sp *SessionProfessor) GetModeAndID() (string, string) {
	parts := strings.Split(sp.session.SessionID, "_")
	mode := parts[0]
	id := parts[1]
	return mode, id
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

func (sp *SessionProfessor) Start(ctx context.Context) {
	if !sp.useStatus.Swap(true) {
		go func(ctx context.Context) {
			for {
				select {
				case msg := <-sp.msgCh:
					sp.HandleEvent(ctx, msg)
				case <-ctx.Done():
					log.Println("session professor context done, session id:", sp.session.SessionID)
					return
				}
			}
		}(ctx)
	}
}

func (sp *SessionProfessor) AppendEvent(event *Event) {
	sp.msgCh <- event
}

func (sp *SessionProfessor) HandleEvent(ctx context.Context, msg *Event) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	sp.session.Append(sp.StructRequest(msg.Message))
	reply, err := sp.GetReply(ctxTimeout, sp.session.GetAll())
	if err != nil {
		log.Println("get reply failed", err)
		return
	}
	err = sp.sendMsgToNapcat(reply)
	if err != nil {
		log.Println(err)
		return
	}
}

func (sp *SessionProfessor) GetReply(ctx context.Context, messages []*schema.Message) (string, error) {
	log.Println("get reply ing......")
	stream, err := sp.chatModel.Stream(ctx, messages)
	if err != nil {
		return "", err
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
	return reply, nil
}

func (sp *SessionProfessor) StructRequest(message []MetaMessage) *schema.Message {
	userMessage := &schema.Message{Role: schema.User}
	for _, meta := range message {
		switch meta.Type {
		case "text":
			text, _ := meta.Data["text"].(string)
			userMessage.UserInputMultiContent = append(userMessage.UserInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeText,
				Text: text,
			})
		case "face":
			faceText, _ := meta.Data["faceText"].(string)
			userMessage.UserInputMultiContent = append(userMessage.UserInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeText,
				Text: faceText,
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
		case "file":
			fileID, _ := meta.Data["file_id"].(string)
			fileURL, err := sp.GetFileData(fileID)
			if err != nil {
				log.Println("get file data failed", err)
				continue
			}
			userMessage.UserInputMultiContent = append(userMessage.UserInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeFileURL,
				File: &schema.MessageInputFile{
					MessagePartCommon: schema.MessagePartCommon{
						URL: &fileURL,
					},
				},
			})
		default:
			continue
		}
	}
	return userMessage
}

func (sp *SessionProfessor) sendMsgToNapcat(message string) (err error) {
	mode, id := sp.GetModeAndID()
	reqBody := SendMsgReq{
		MessageType: mode,
		Message:     message,
	}
	if mode == "private" {
		reqBody.UserID = id
	} else if mode == "group" {
		reqBody.GroupID = id
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return
	}
	url := sp.napcatConfig.ApiBaseURL + "/send_msg"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if sp.napcatConfig.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+sp.napcatConfig.AuthToken)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("send message to napcat failed, status code: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

func (sp *SessionProfessor) GetFileData(fileID string) (string, error) {
	mode, id := sp.GetModeAndID()
	reqBody := map[string]string{
		"file_id": fileID,
	}
	url := ""
	if mode == "private" {
		reqBody["user_id"] = id
		url = sp.napcatConfig.ApiBaseURL + "/get_private_file_url"
	} else if mode == "group" {
		reqBody["group_id"] = id
		url = sp.napcatConfig.ApiBaseURL + "get_group_file_url"
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("marshal get file req failed", err)
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("create get file req failed", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if sp.napcatConfig.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+sp.napcatConfig.AuthToken)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("get file from napcat failed", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("get file from napcat failed, status code:", resp.StatusCode)
		return "", err
	}
	var result NapcatFile
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("decode get file response failed", err)
		return "", err
	}
	if result.Status != "ok" {
		return "", errors.New("get file from napcat failed, status: " + result.Status)
	}
	return result.Data["url"].(string), nil
}

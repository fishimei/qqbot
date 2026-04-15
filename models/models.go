package models

import (
	"bot/config"
	"strconv"
)

//{"type":"image","data":{"summary":"","file":"2E8289BBC260D9F93ED86E69387ACB93.jpg","sub_type":0,"url":"https://multimedia.nt.qq.com.cn/download?appid=1407&fileid=EhRyld68OCCPnlXppM_GU4DHM3cTzRjl0gog_wooz4-xwIfqkwMyBHByb2RQgL2jAVoQuUwyWBeMqO2e2Xr3kvwRmHoCEhmCAQJneg&rkey=CAMSMNSNtPNodN3RIGV9uivGRSMTpN5fOrtxxt-ORnqrj6fWA2g7jIQVGeyeZfB2KvGZuQ","file_size":"174437"}}]

type Event struct {
	PostType    string        `json:"post_type"`
	MessageType string        `json:"message_type"`
	UserID      int64         `json:"user_id"`
	GroupID     int64         `json:"group_id"`
	Message     []MetaMessage `json:"message"`
	Time        int64         `json:"time"`
	SelfID      int64         `json:"self_id"`
}

func (e *Event) StructSessionID() string {
	switch e.MessageType {
	case "private":
		return "private_" + strconv.FormatInt(e.UserID, 10)
	case "group":
		return "group_" + strconv.FormatInt(e.GroupID, 10)
	default:
		return ""
	}
}

type MetaMessage struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type SendMsgReq struct {
	MessageType string `json:"message_type"`
	UserID      int64  `json:"user_id,omitempty"`
	GroupID     int64  `json:"group_id,omitempty"`
	Message     string `json:"message"`
}

type Session struct {
	SessionID string `json:"session_id"`
	MessagesRingBuffer
}

func NewSession(sessionID string, systemCap, ringCap int) *Session {
	return &Session{
		SessionID:          sessionID,
		MessagesRingBuffer: NewMessagesRingBuffer(systemCap, ringCap),
	}
}

type NapcatConfig struct {
	ApiBaseURL string `json:"api_base_url"`
	AuthToken  string `json:"auth_token"`
}

func NewNapcatConfig() *NapcatConfig {
	apiBaseURL, _, authToken := config.LoadNapcatConfig()
	return &NapcatConfig{
		ApiBaseURL: apiBaseURL,
		AuthToken:  authToken,
	}
}

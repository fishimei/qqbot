package models

import (
	"sync"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

type SessionProfessorRegister struct {
	sync.RWMutex
	sessionMap map[string]*SessionProfessor
	chatModel  *ark.ChatModel
}

func NewSessionRegister(model *ark.ChatModel) *SessionProfessorRegister {
	return &SessionProfessorRegister{
		sessionMap: make(map[string]*SessionProfessor),
		chatModel:  model,
	}
}

func (sr *SessionProfessorRegister) GetSessionProfessor(sessionID string) (*SessionProfessor, bool) {
	sr.RLock()
	defer sr.RUnlock()
	sessionProfessor, exist := sr.sessionMap[sessionID]
	return sessionProfessor, exist
}

func (sr *SessionProfessorRegister) RegisterSessionProfessor(sessionID string, chatModel *ark.ChatModel) *SessionProfessor {
	sr.Lock()
	defer sr.Unlock()
	sp, exist := sr.sessionMap[sessionID]
	if exist {
		return sp
	}
	sessionProfessor := NewSessionProfessor(sessionID, chatModel)
	sr.sessionMap[sessionID] = sessionProfessor
	return sessionProfessor
}

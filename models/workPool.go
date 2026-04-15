package models

import (
	"sync"
)

// WorkPool 是一个工作池结构体，用于处理事件的并发执行。
type WorkPool struct {
	workerCount int
	msgCh       chan *Event
	sync.WaitGroup
	register *SessionProfessorRegister
}

func NewWorkPool(workerCount, poolSize int, register *SessionProfessorRegister) *WorkPool {
	return &WorkPool{
		workerCount: workerCount,
		msgCh:       make(chan *Event, poolSize),
		register:    register,
	}
}

func (wp *WorkPool) AddEvent(event *Event) {
	wp.msgCh <- event
}

func (wp *WorkPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.Add(1)
		go func() {
			defer wp.Done()
			for event := range wp.msgCh {
				wp.handleEvent(event)
			}
		}()
	}
	wp.Wait()
}

func (wp *WorkPool) handleEvent(event *Event) {
	sessionID := event.StructSessionID()
	sessionProfessor, exist := wp.register.GetSessionProfessor(sessionID)
	if !exist {
		sessionProfessor = wp.register.RegisterSessionProfessor(NewSessionProfessor(sessionID, wp.register.chatModel))
	}
	sessionProfessor.AppendEvent(event)
	sessionProfessor.Start()
}

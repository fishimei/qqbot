package models

import (
	"fmt"

	"github.com/cloudwego/eino/schema"
)

type MessagesRingBuffer struct {
	messages  []*schema.Message // 总容量 = systemCap + ringCap
	systemCap int               // 系统消息占用的固定槽位数
	systemNum int               // 当前系统消息数量（必须 <= systemCap）
	ringCap   int               // 环形区域容量
	head      int               // 环形区域的下一个写入位置（相对于环形区域起始索引）
	size      int               // 环形区域中已存储的消息数量
}

func NewMessagesRingBuffer(systemCap, ringCap int) MessagesRingBuffer {
	return MessagesRingBuffer{
		ringCap:   ringCap,
		systemCap: systemCap,
		systemNum: 0,
		messages:  make([]*schema.Message, systemCap+ringCap),
		head:      0,
		size:      0,
	}
}

// SetSystemMessages 设置系统消息（必须在任何普通消息之前调用）
// 索引 0..systemCap-1 被永久占用
func (rb *MessagesRingBuffer) SetSystemMessages(msg *schema.Message) error {
	if rb.systemNum > rb.systemCap {
		return fmt.Errorf("系统消息数量 %d 超过预留容量 %d", rb.systemNum, rb.systemCap)
	}
	rb.messages[rb.systemNum] = msg
	rb.systemNum++
	return nil
}

// Append 追加普通消息（对话消息）
func (rb *MessagesRingBuffer) Append(msg *schema.Message) {
	// 计算环形区域的实际索引
	idx := rb.systemCap + rb.head
	rb.messages[idx] = msg
	rb.head = (rb.head + 1) % rb.ringCap
	if rb.size < rb.ringCap {
		rb.size++
	}
}

// GetAll 获取所有消息（系统消息 + 环形区域内的消息，按时间顺序）
func (rb *MessagesRingBuffer) GetAll() []*schema.Message {
	// 1. 收集系统消息（只收集非 nil 的）
	var result []*schema.Message
	for i := 0; i < rb.systemCap; i++ {
		if rb.messages[i] != nil {
			result = append(result, rb.messages[i])
		}
	}
	// 2. 收集环形区域的消息（按顺序：从最旧到最新）
	if rb.size == 0 {
		return result
	}
	start := rb.head - rb.size
	if start < 0 {
		start += rb.ringCap
	}
	for i := 0; i < rb.size; i++ {
		idx := (start + i) % rb.ringCap
		msg := rb.messages[rb.systemCap+idx]
		if msg != nil {
			result = append(result, msg)
		}
	}
	return result
}

// Size 获取环形区域当前的消息数量
func (rb *MessagesRingBuffer) Size() int {
	return rb.size
}

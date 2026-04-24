package chatServer

import (
	"sync"

	"github.com/cloudwego/eino/schema"
)

// 全局会话存储（内存 map）
var SimpleMemoryMap = &sync.Map{}

// 单个会话的记忆
type SimpleMemory struct {
	mu            sync.Mutex        // 保护并发读写
	ID            string            // 会话ID
	Messages      []*schema.Message // 消息历史
	MaxWindowSize int               // 最大窗口大小，默认是6
}

func NewMemory(id string, max int) error {
	if max <= 0 {
		max = 6
	}
	var simpleMemory *SimpleMemory
	simpleMemory = &SimpleMemory{
		ID:            id,
		Messages:      []*schema.Message{},
		MaxWindowSize: 6,
	}
	SimpleMemoryMap.Store(id, simpleMemory)
	return nil
}

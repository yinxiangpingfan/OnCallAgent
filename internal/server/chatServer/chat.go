package chatServer

import (
	"OnCallAgent/internal/server/ai/agent/chat"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/sirupsen/logrus"
)

type ChatServer interface {
	Chat(ctx context.Context, question string, id string) (string, error)
	ChatSream(ctx context.Context, question string, id string, msgChan *chan string, doneChan *chan struct{}) error
}

type chatServer struct {
	logger *logrus.Logger
	runner compose.Runnable[*chat.UserMessage, *schema.Message]
}

func NewChatServer(log *logrus.Logger, runner compose.Runnable[*chat.UserMessage, *schema.Message]) ChatServer {
	return &chatServer{logger: log, runner: runner}
}

func (c *chatServer) Chat(ctx context.Context, question string, id string) (string, error) {
	// 1. 从 sync.Map 加载会话记忆
	val, ok := SimpleMemoryMap.Load(id)
	if !ok {
		// 会话不存在，自动创建（使用默认窗口大小 6）
		if err := NewMemory(id, 0); err != nil {
			return "", fmt.Errorf("创建会话失败: %w", err)
		}
		val, ok = SimpleMemoryMap.Load(id)
		if !ok {
			return "", fmt.Errorf("创建会话失败")
		}
	}
	memory := val.(*SimpleMemory)

	// 2. 加锁保证并发安全（读取 history 快照 + 写回）
	memory.mu.Lock()
	// 取当前历史的副本传给模型，避免后续写操作互相干扰
	historyCopy := make([]*schema.Message, len(memory.Messages))
	copy(historyCopy, memory.Messages)
	memory.mu.Unlock()

	// 3. 调用大模型
	output, err := c.runner.Invoke(ctx, &chat.UserMessage{
		ID:      id,
		Query:   question,
		History: historyCopy,
	})
	if err != nil {
		c.logger.Errorf("调用大模型失败, id: %s, err: %v", id, err)
		return "", fmt.Errorf("调用大模型失败: %w", err)
	}

	// 4. 构建本轮消息并写回记忆（加锁保证并发安全）
	userMsg := schema.UserMessage(question)
	memory.mu.Lock()
	memory.Messages = append(memory.Messages, userMsg, output)
	// 5. 滑动窗口裁剪，保留最近 MaxWindowSize 条
	if len(memory.Messages) > memory.MaxWindowSize {
		memory.Messages = memory.Messages[len(memory.Messages)-memory.MaxWindowSize:]
	}
	memory.mu.Unlock()

	return output.Content, nil
}

func (c *chatServer) ChatSream(ctx context.Context, question string, id string, msgChan *chan string, doneChan *chan struct{}) error {
	// 1. 从 sync.Map 加载会话记忆
	val, ok := SimpleMemoryMap.Load(id)
	if !ok {
		// 会话不存在，自动创建（使用默认窗口大小 6）
		if err := NewMemory(id, 0); err != nil {
			return fmt.Errorf("创建会话失败: %w", err)
		}
		val, ok = SimpleMemoryMap.Load(id)
		if !ok {
			return fmt.Errorf("创建会话失败")
		}
	}
	memory := val.(*SimpleMemory)

	// 2. 加锁保证并发安全（读取 history 快照 + 写回）
	memory.mu.Lock()
	// 取当前历史的副本传给模型，避免后续写操作互相干扰
	historyCopy := make([]*schema.Message, len(memory.Messages))
	copy(historyCopy, memory.Messages)
	memory.mu.Unlock()

	output, err := c.runner.Stream(ctx, &chat.UserMessage{
		ID:      id,
		Query:   question,
		History: historyCopy,
	})
	if err != nil {
		c.logger.Errorf("调用大模型失败, id: %s, err: %v", id, err)
		return fmt.Errorf("调用大模型失败: %w", err)
	}
	var once sync.Once
	closeOnce := func() {
		once.Do(func() { close(*msgChan) })
	}
	sb := strings.Builder{}
	for {
		// 非阻塞检查 context 是否已取消
		select {
		case <-ctx.Done():
			closeOnce()
			*doneChan <- struct{}{}
			return nil
		default:
		}
		msg, err := output.Recv()
		if errors.Is(err, io.EOF) {
			closeOnce()
			*doneChan <- struct{}{}
			break
		}
		if err != nil {
			c.logger.Errorf("调用大模型失败, id: %s, err: %v", id, err)
			return fmt.Errorf("调用大模型失败: %w", err)
		}
		// 累积完整回复 + 推送 token 到前端
		sb.WriteString(msg.Content)
		*msgChan <- msg.Content
	}

	// 流结束，将本轮对话（用户问题 + AI 完整回复）写回 memory
	userMsg := schema.UserMessage(question)
	assistantMsg := schema.AssistantMessage(sb.String(), nil)
	memory.mu.Lock()
	memory.Messages = append(memory.Messages, userMsg, assistantMsg)
	// 滑动窗口裁剪
	if len(memory.Messages) > memory.MaxWindowSize {
		memory.Messages = memory.Messages[len(memory.Messages)-memory.MaxWindowSize:]
	}
	memory.mu.Unlock()

	return nil
}

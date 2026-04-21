package chat

import (
	"OnCallAgent/pkg/config"

	qdrant_retriever "github.com/cloudwego/eino-ext/components/retriever/qdrant"
)

type chatServer struct {
	retriever *qdrant_retriever.Retriever
	config    *config.Config
}

func NewChatServer(retriever *qdrant_retriever.Retriever, cfg *config.Config) chatServer {
	return chatServer{
		retriever: retriever,
		config:    cfg,
	}
}

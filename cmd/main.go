package main

import (
	indexerr "OnCallAgent/internal/repo/qrdant/indexer"
	initQdrantRepo "OnCallAgent/internal/repo/qrdant/init"
	"OnCallAgent/internal/repo/qrdant/retriever"
	"OnCallAgent/internal/router"
	"OnCallAgent/internal/server/ai/agent/chat"
	knowledgeindex "OnCallAgent/internal/server/ai/agent/knowledge_index"
	"OnCallAgent/internal/server/ai/embeder"
	"OnCallAgent/internal/server/model"
	"OnCallAgent/pkg/config"
	"OnCallAgent/pkg/log"
	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	// 初始化日志记录器
	log := log.InitLogger("info", "log/OnCallAgent.log")
	//初始化配置
	config, err := config.InitConfig("./config/config.json")
	if err != nil {
		panic(err)
	}
	//初始化qdrant
	indexer, err := initQdrantRepo.NewQdrantIndexer(ctx, config)
	if err != nil {
		panic(err)
	}
	//初始化embedding
	embedder, err := embeder.NewEmbedder(ctx, config)
	if err != nil {
		panic(err)
	}
	//初始化retriever
	retriever := retriever.NewRetrieverServer(ctx, indexer, *embedder)
	run, err := retriever.NewRetrieverServer(ctx, "oncallagent", 0.5, 2)
	if err != nil {
		panic(err)
	}
	//初始化chatAgent
	re := chat.NewChatServer(run, config)
	runner, err := re.BuildChatAgent(ctx)
	if err != nil {
		panic(err)
	}
	//新建集合
	indexerr := indexerr.NewQdranIndexerServer(ctx, indexer, *embedder)
	err = indexerr.NewQdrantIndexer(ctx)
	if err != nil {
		panic(err)
	}
	//初始化RAGagent
	knowledgeIndex := knowledgeindex.NewKnowledgeIndex(embeder.NewEmbeddingServer(embedder), indexerr)
	runnerRAG, err := knowledgeIndex.NewGraph(ctx)
	if err != nil {
		panic(err)
	}
	//初始化model
	chatModel, err := model.NewOpenaiModel(ctx, config)
	if err != nil {
		panic(err)
	}
	// 初始化gin
	r := gin.Default()
	router.InitRouter(ctx, r, log, config, runnerRAG, runner, chatModel)
}

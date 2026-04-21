package knowledgeindex

import (
	"OnCallAgent/pkg/tool"
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/qdrant/go-client/qdrant"
)

const (
	FileLoader       = "FileLoader"
	MarkdownSplitter = "MarkdownSplitter"
	QdrantIndexer    = "QdrantIndexer"
)

func (k *knowledgeIndex) NewGraph(ctx context.Context) (r compose.Runnable[document.Source, bool], err error) {
	g := compose.NewGraph[document.Source, bool]()
	lodder, err := k.NewFileLoader(ctx)
	if err != nil {
		return nil, fmt.Errorf("knowledge_index:Graph:生成lodder失败%v", err)
	}
	err = g.AddLoaderNode(FileLoader, lodder)
	if err != nil {
		return nil, fmt.Errorf("knowledge_index:Graph:添加FileLoader节点失败%v", err)
	}
	transformer, err := k.NewSplitMarkdown(ctx)
	if err != nil {
		return nil, fmt.Errorf("knowledge_index:Graph:生成transformer失败%v", err)
	}
	err = g.AddDocumentTransformerNode(MarkdownSplitter, transformer)
	if err != nil {
		return nil, fmt.Errorf("knowledge_index:Graph:添加transformer节点失败%v", err)
	}
	qdantNode := compose.InvokableLambda(k.textToQdrantIndex())
	err = g.AddLambdaNode(QdrantIndexer, qdantNode)
	if err != nil {
		return nil, fmt.Errorf("knowledge_index:Graph:添加QdrantIndexer节点失败%v", err)
	}

	//构建图
	_ = g.AddEdge(compose.START, FileLoader)
	_ = g.AddEdge(QdrantIndexer, compose.END)
	_ = g.AddEdge(FileLoader, MarkdownSplitter)
	_ = g.AddEdge(MarkdownSplitter, QdrantIndexer)
	r, err = g.Compile(ctx, compose.WithGraphName("KnowledgeIndexing"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, fmt.Errorf("knowledge_index:Graph:编译图失败%v", err)
	}

	return r, nil
}

func (k *knowledgeIndex) textToQdrantIndex() func(ctx context.Context, req []*schema.Document) (bool, error) {
	return func(ctx context.Context, req []*schema.Document) (bool, error) {
		points := make([]*qdrant.PointStruct, 0, len(req))
		for _, doc := range req {
			lines := strings.Split(doc.Content, "\n")
			if len(lines) == 0 {
				continue
			}
			title := lines[0]
			// 只处理一级标题
			if !strings.HasPrefix(title, "#") {
				continue
			}
			title = strings.TrimPrefix(title, "# ")
			body := lines[1:]

			// title 权重为 2，其余行各一份
			temp := []string{title, title}
			temp = append(temp, body...)

			t, err := k.embederServer.Embedding(ctx, temp)
			if err != nil {
				return false, fmt.Errorf("knowledge_index:Graph:向量化失败: %w", err)
			}
			res, err := k.embederServer.Average(t)
			if err != nil {
				return false, fmt.Errorf("knowledge_index:Graph:求均值失败: %w", err)
			}
			res = k.embederServer.Normalize(res)

			payload := qdrant.NewValueMap(map[string]any{
				"content": title + "\n" + strings.Join(body, "\n"),
			})
			points = append(points, &qdrant.PointStruct{
				Id: &qdrant.PointId{
					PointIdOptions: &qdrant.PointId_Uuid{Uuid: doc.ID},
				},
				Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Vector: &qdrant.Vector_Dense{Dense: &qdrant.DenseVector{Data: tool.ToFloat32(res)}}}}},
				Payload: payload,
			})
		}
		if len(points) == 0 {
			return true, nil
		}
		if err := k.qdrantServer.AddVector(ctx, &qdrant.UpsertPoints{Points: points}); err != nil {
			return false, fmt.Errorf("knowledge_index:Graph:写入Qdrant失败: %w", err)
		}
		return true, nil
	}
}

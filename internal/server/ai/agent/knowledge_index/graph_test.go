package knowledgeindex

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/qdrant/go-client/qdrant"
)

// --- Mock 实现 ---

type mockEmbeddingServer struct {
	embeddingFn func(ctx context.Context, text []string) ([][]float64, error)
	averageFn   func(embeddings [][]float64) ([]float64, error)
	normalizeFn func(embeddings []float64) []float64
}

func (m *mockEmbeddingServer) Embedding(ctx context.Context, text []string) ([][]float64, error) {
	if m.embeddingFn != nil {
		return m.embeddingFn(ctx, text)
	}
	// 默认：每个 text 返回一个 3 维向量
	result := make([][]float64, len(text))
	for i := range result {
		result[i] = []float64{0.1, 0.2, 0.3}
	}
	return result, nil
}

func (m *mockEmbeddingServer) Average(embeddings [][]float64) ([]float64, error) {
	if m.averageFn != nil {
		return m.averageFn(embeddings)
	}
	if len(embeddings) == 0 {
		return nil, errors.New("empty embeddings")
	}
	dim := len(embeddings[0])
	avg := make([]float64, dim)
	for _, e := range embeddings {
		for i, v := range e {
			avg[i] += v
		}
	}
	for i := range avg {
		avg[i] /= float64(len(embeddings))
	}
	return avg, nil
}

func (m *mockEmbeddingServer) Normalize(embeddings []float64) []float64 {
	if m.normalizeFn != nil {
		return m.normalizeFn(embeddings)
	}
	return embeddings
}

type mockQdrantServer struct {
	addVectorFn func(ctx context.Context, points *qdrant.UpsertPoints) error
	capturedPts *qdrant.UpsertPoints
}

func (m *mockQdrantServer) NewQdrantIndexer(ctx context.Context) error { return nil }

func (m *mockQdrantServer) AddVector(ctx context.Context, points *qdrant.UpsertPoints) error {
	m.capturedPts = points
	if m.addVectorFn != nil {
		return m.addVectorFn(ctx, points)
	}
	return nil
}

// --- 辅助函数 ---

func newTestKnowledgeIndex() *knowledgeIndex {
	return &knowledgeIndex{
		embederServer: &mockEmbeddingServer{},
		qdrantServer:  &mockQdrantServer{},
	}
}

// --- 测试 NewGraph ---

// TestGraph_Compile 验证图可以正常编译
func TestGraph_Compile(t *testing.T) {
	ctx := context.Background()
	k := newTestKnowledgeIndex()

	runnable, err := k.NewGraph(ctx)
	if err != nil {
		t.Fatalf("NewGraph 失败: %v", err)
	}
	if runnable == nil {
		t.Fatal("NewGraph 返回了 nil Runnable")
	}
}

// --- 测试 textToQdrantIndex ---

// TestTextToQdrantIndex_NormalH1 测试包含一级标题的正常文档
func TestTextToQdrantIndex_NormalH1(t *testing.T) {
	ctx := context.Background()
	mockQdrant := &mockQdrantServer{}
	k := &knowledgeIndex{
		embederServer: &mockEmbeddingServer{},
		qdrantServer:  mockQdrant,
	}

	docs := []*schema.Document{
		{
			ID:      "test-uuid-1",
			Content: "# 告警处理\n排查步骤一\n排查步骤二",
		},
	}

	fn := k.textToQdrantIndex()
	ok, err := fn(ctx, docs)
	if err != nil {
		t.Fatalf("textToQdrantIndex 失败: %v", err)
	}
	if !ok {
		t.Error("期望返回 true，实际返回 false")
	}
	if mockQdrant.capturedPts == nil {
		t.Fatal("未调用 AddVector")
	}
}

// TestTextToQdrantIndex_SkipsNoHeading 测试没有任何标题前缀的文档会被跳过
func TestTextToQdrantIndex_SkipsNoHeading(t *testing.T) {
	ctx := context.Background()
	mockQdrant := &mockQdrantServer{}
	k := &knowledgeIndex{
		embederServer: &mockEmbeddingServer{},
		qdrantServer:  mockQdrant,
	}

	// 两个文档均无 '#' 开头的标题行，均应被 continue 跳过
	docs := []*schema.Document{
		{
			ID:      "test-uuid-2",
			Content: "普通文本内容\n没有标题",
		},
		{
			ID:      "test-uuid-3",
			Content: "普通段落\n更多内容",
		},
	}

	fn := k.textToQdrantIndex()
	ok, err := fn(ctx, docs)
	if err != nil {
		t.Fatalf("textToQdrantIndex 失败: %v", err)
	}
	if !ok {
		t.Error("期望返回 true")
	}
	if mockQdrant.capturedPts == nil {
		t.Fatal("未调用 AddVector")
	}
	// 所有文档均被跳过，points 应全为 nil
	for i, pt := range mockQdrant.capturedPts.Points {
		if pt != nil {
			t.Errorf("doc[%d] 无标题前缀，不应生成 PointStruct，但得到: %+v", i, pt)
		}
	}
}

// TestTextToQdrantIndex_AnyHeadingIsProcessed 验证以 '#' 开头的标题（包括多级）都会被处理
func TestTextToQdrantIndex_AnyHeadingIsProcessed(t *testing.T) {
	ctx := context.Background()
	mockQdrant := &mockQdrantServer{}
	k := &knowledgeIndex{
		embederServer: &mockEmbeddingServer{},
		qdrantServer:  mockQdrant,
	}

	docs := []*schema.Document{
		{ID: "uuid-h2", Content: "## 二级标题\n内容"},
		{ID: "uuid-h3", Content: "### 三级标题\n内容"},
	}

	fn := k.textToQdrantIndex()
	ok, err := fn(ctx, docs)
	if err != nil {
		t.Fatalf("textToQdrantIndex 失败: %v", err)
	}
	if !ok {
		t.Error("期望返回 true")
	}
	if mockQdrant.capturedPts == nil {
		t.Fatal("未调用 AddVector")
	}
	// 两个文档均有 '#' 前缀，应当都被处理（points 不为 nil）
	for i, pt := range mockQdrant.capturedPts.Points {
		if pt == nil {
			t.Errorf("doc[%d] 有标题前缀，应生成 PointStruct，但为 nil", i)
		}
	}
}

// TestTextToQdrantIndex_EmptyDocList 测试空文档列表
func TestTextToQdrantIndex_EmptyDocList(t *testing.T) {
	ctx := context.Background()
	mockQdrant := &mockQdrantServer{}
	k := &knowledgeIndex{
		embederServer: &mockEmbeddingServer{},
		qdrantServer:  mockQdrant,
	}

	fn := k.textToQdrantIndex()
	ok, err := fn(ctx, []*schema.Document{})
	if err != nil {
		t.Fatalf("空文档列表不应返回错误, got: %v", err)
	}
	if !ok {
		t.Error("期望返回 true")
	}
}

// TestTextToQdrantIndex_EmbeddingError 测试向量化失败时正确返回错误
func TestTextToQdrantIndex_EmbeddingError(t *testing.T) {
	ctx := context.Background()
	embeddingErr := errors.New("向量化服务不可用")
	k := &knowledgeIndex{
		embederServer: &mockEmbeddingServer{
			embeddingFn: func(_ context.Context, _ []string) ([][]float64, error) {
				return nil, embeddingErr
			},
		},
		qdrantServer: &mockQdrantServer{},
	}

	docs := []*schema.Document{
		{
			ID:      "test-uuid-4",
			Content: "# 故障标题\n故障详情",
		},
	}

	fn := k.textToQdrantIndex()
	ok, err := fn(ctx, docs)
	if err == nil {
		t.Fatal("期望返回错误，但得到 nil")
	}
	if ok {
		t.Error("出错时期望返回 false")
	}
}

// TestTextToQdrantIndex_AverageError 测试求平均失败时正确返回错误
func TestTextToQdrantIndex_AverageError(t *testing.T) {
	ctx := context.Background()
	averageErr := errors.New("平均计算失败")
	k := &knowledgeIndex{
		embederServer: &mockEmbeddingServer{
			averageFn: func(_ [][]float64) ([]float64, error) {
				return nil, averageErr
			},
		},
		qdrantServer: &mockQdrantServer{},
	}

	docs := []*schema.Document{
		{
			ID:      "test-uuid-5",
			Content: "# 故障标题\n故障详情",
		},
	}

	fn := k.textToQdrantIndex()
	ok, err := fn(ctx, docs)
	if err == nil {
		t.Fatal("期望返回错误，但得到 nil")
	}
	if ok {
		t.Error("出错时期望返回 false")
	}
}

// TestTextToQdrantIndex_AddVectorError 测试写入 Qdrant 失败时正确返回错误
func TestTextToQdrantIndex_AddVectorError(t *testing.T) {
	ctx := context.Background()
	qdrantErr := errors.New("Qdrant 不可达")
	k := &knowledgeIndex{
		embederServer: &mockEmbeddingServer{},
		qdrantServer: &mockQdrantServer{
			addVectorFn: func(_ context.Context, _ *qdrant.UpsertPoints) error {
				return qdrantErr
			},
		},
	}

	docs := []*schema.Document{
		{
			ID:      "test-uuid-6",
			Content: "# 告警标题\n告警详情",
		},
	}

	fn := k.textToQdrantIndex()
	ok, err := fn(ctx, docs)
	if err == nil {
		t.Fatal("期望返回错误，但得到 nil")
	}
	if ok {
		t.Error("出错时期望返回 false")
	}
}

// TestTextToQdrantIndex_MultipleH1Docs 测试多个一级标题文档均被处理
func TestTextToQdrantIndex_MultipleH1Docs(t *testing.T) {
	ctx := context.Background()
	mockQdrant := &mockQdrantServer{}
	k := &knowledgeIndex{
		embederServer: &mockEmbeddingServer{},
		qdrantServer:  mockQdrant,
	}

	docs := []*schema.Document{
		{ID: "uuid-a", Content: "# 标题A\n内容A"},
		{ID: "uuid-b", Content: "# 标题B\n内容B"},
		{ID: "uuid-c", Content: "# 标题C\n内容C"},
	}

	fn := k.textToQdrantIndex()
	ok, err := fn(ctx, docs)
	if err != nil {
		t.Fatalf("textToQdrantIndex 失败: %v", err)
	}
	if !ok {
		t.Error("期望返回 true")
	}
	if mockQdrant.capturedPts == nil {
		t.Fatal("未调用 AddVector")
	}
	if len(mockQdrant.capturedPts.Points) != len(docs) {
		t.Errorf("期望 %d 个 point，实际 %d 个", len(docs), len(mockQdrant.capturedPts.Points))
	}
}

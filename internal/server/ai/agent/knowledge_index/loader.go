package knowledgeindex

import (
	"context"

	fileloader "github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino/components/document"

	"github.com/cloudwego/eino/components/document/parser"
)

// 新建文件 loader
func (k *knowledgeIndex) NewFileLoader(ctx context.Context) (document.Loader, error) {
	// 创建扩展解析器，支持多种文件格式
	extParser, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
		FallbackParser: parser.TextParser{}, // 默认使用文本解析器
	})
	if err != nil {
		return nil, err
	}

	// 创建文件加载器
	loader, err := fileloader.NewFileLoader(ctx, &fileloader.FileLoaderConfig{
		UseNameAsID: true, // 使用文件名作为文档 ID
		Parser:      extParser,
	})
	if err != nil {
		return nil, err
	}

	return loader, nil
}

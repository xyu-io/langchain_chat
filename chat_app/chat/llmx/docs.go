package llmx

import (
	"context"
	"langchain_chat/third/langchaingo/documentloaders"
	"langchain_chat/third/langchaingo/schema"
	"langchain_chat/third/langchaingo/textsplitter"
	"os"
)

func DocumentLoader(path string) ([]schema.Document, error) {
	fs, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// 文本加载、分割处理
	docs, err := documentloaders.NewText(fs).LoadAndSplit(context.Background(),
		textsplitter.NewRecursiveCharacter(
			textsplitter.WithChunkSize(1000),
			textsplitter.WithChunkOverlap(20),
		))
	if err != nil {
		return nil, err
	}

	return docs, nil
}

// TextToChunks 函数将文本文件转换为文档块
func TextToChunks(dirFile string, chunkSize, chunkOverlap int) ([]schema.Document, error) {
	file, err := os.Open(dirFile)
	if err != nil {
		return nil, err
	}
	// 创建一个新的文本文档加载器
	docLoaded := documentloaders.NewText(file)
	// 创建一个新的递归字符文本分割器
	split := textsplitter.NewRecursiveCharacter()
	// 设置块大小
	split.ChunkSize = chunkSize
	// 设置块重叠大小
	split.ChunkOverlap = chunkOverlap
	// 加载并分割文档
	docs, err := docLoaded.LoadAndSplit(context.Background(), split)
	if err != nil {
		return nil, err
	}
	return docs, nil
}

func DocToArray(docs []schema.Document) ([]string, error) {
	var contents []string
	for _, doc := range docs {
		contents = append(contents, doc.PageContent)
	}
	return contents, nil
}

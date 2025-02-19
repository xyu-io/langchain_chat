package main

import (
	log "github.com/sirupsen/logrus"
	"langchain_chat/chat_app/chat/llmx"
	"testing"
)

// 本地知识文档上传-需要数据库支撑，目前支持chroma,miluvs
func TestFileUpload(t *testing.T) {
	files := []string{
		"", // exp ./doc_file/chat.txt
	}
	ins, err := llmx.NewLLM(
		llmx.WithModel(DEFAULT_LLM_MODEL),
		llmx.WithEmbedderModel(DEFAULT_EMBEDDING_MODEL),
		llmx.WithLLMUrl(Default_Ollama_Server),
		llmx.UseStore(llmx.WithStoreServer(DEFAULT_MILVUS_URL)),
		llmx.WithContentSize(8*1024))
	if err != nil {
		log.Error(err)
		return
	}

	for _, file := range files {
		err = ins.DocumentLoader(file)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	t.Log("文件上传完成")
}

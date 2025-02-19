package main

import (
	"context"
	"fmt"
	"langchain_chat/chat_app/chat/llmx"
	"langchain_chat/third/langchaingo/llms"
	"testing"
)

// 不需要进行文档知识检索增强
func TestGeneralChat(t *testing.T) {
	llm, err := llmx.NewLLM(
		llmx.WithModel(DEFAULT_LLM_MODEL),
		llmx.WithEmbedderModel(DEFAULT_EMBEDDING_MODEL),
		llmx.WithContentSize(8*1024))
	if err != nil {
		return
	}
	_, err = llm.GetOllamaLLM().Call(
		context.Background(),
		"你好，请简单介绍一下你自己",
		llms.WithStreamingFunc(
			func(_ context.Context, bs []byte) error {
				fmt.Print(string(bs))
				return nil
			}),
	)
	if err != nil {
		t.Error(err)
		return
	}

}

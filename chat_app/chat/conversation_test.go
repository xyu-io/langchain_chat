package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"langchain_chat/base"
	"langchain_chat/chat_app/chat/llmx"
	"langchain_chat/third/langchaingo/chains"
	"langchain_chat/third/langchaingo/memory"
	"testing"
)

func TestConversation(t *testing.T) {
	llm, err := llmx.NewLLM(
		llmx.WithModel(DEFAULT_LLM_MODEL),
		llmx.WithEmbedderModel(DEFAULT_EMBEDDING_MODEL),
		llmx.WithContentSize(8*1024),
		llmx.UseStore(llmx.WithStoreServer(DEFAULT_MILVUS_URL))) // 默认milvus
	if err != nil {
		return
	}
	chain, err := NewChain(WithLLM(llm), WithMemory(memory.NewConversationBuffer())).
		chainRoute(base.ChatSettings{
			ChainType:       base.ConversationChainType,
			Temperature:     0,
			MinResultScore:  800, // 需要结合查询/文档向量来确定，特别是miluvs数据库，不确定可以不填写
			AmountOfResults: 3,   // 检索文档数量
		})
	if err != nil {
		log.Error(err)
		return
	}
	query := "你好，请简单介绍一下你自己"
	answer, err := chains.Call(context.Background(), chain,
		map[string]any{
			"question": query,
			//"agent_scratchpad": "",
		},
		[]chains.ChainCallOption{
			chains.WithStreamingFunc( // 流式输出 - 和前端交互
				func(_ context.Context, bs []byte) error {
					return nil
				}),
		}...,
	)

	if err != nil {
		t.Error(err)
		return
	}
	log.Println(answer)
}

package main

import (
	"context"
	"langchain_chat/base"
	"langchain_chat/chat_app/chat/llmx"
	"langchain_chat/third/langchaingo/chains"
	"langchain_chat/third/langchaingo/memory"
	"log"
	"testing"
)

// sql助手-目前sqlchain基于mysql数据库，用户可以自己拓展
func TestSqlAssistance(t *testing.T) {
	llm, err := llmx.NewLLM(
		llmx.WithModel(DEFAULT_LLM_MODEL),
		llmx.WithEmbedderModel(DEFAULT_EMBEDDING_MODEL),
		llmx.WithContentSize(8*1024))
	if err != nil {
		return
	}

	query := "查询表event中event_id字段值为`100`的记录信息"

	chain, err := NewChain(WithLLM(llm), WithMemory(memory.NewConversationBuffer())).
		chainRoute(base.ChatSettings{
			ChainType:       base.SQLChainType,
			Temperature:     0,
			MinResultScore:  800, // 需要结合查询/文档向量来确定，特别是miluvs数据库，不确定可以不填写
			AmountOfResults: 3,   // 检索文档数量
		})
	if err != nil {
		t.Error(err)
		return
	}
	answer, err := chains.Call(context.Background(), chain,
		map[string]any{
			chains.SQLChainInputKeys[0]: query,
			chains.SQLChainInputKeys[1]: []string{"event"},
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

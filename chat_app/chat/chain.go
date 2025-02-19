package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"langchain_chat/base"
	"langchain_chat/third/langchaingo/agents"
	"langchain_chat/third/langchaingo/chains"
	"langchain_chat/third/langchaingo/llms/ollama"
	"langchain_chat/third/langchaingo/memory"
	"langchain_chat/third/langchaingo/schema"
	"langchain_chat/third/langchaingo/tools"
	"langchain_chat/third/langchaingo/tools/sqldatabase"
	"langchain_chat/third/langchaingo/tools/sqldatabase/mysql"
	"langchain_chat/third/langchaingo/vectorstores"
)

// 模型路由
func (c *Chain) chainRoute(set base.ChatSettings) (chains.Chain, error) {
	if set.ChainType == "" {
		set.ChainType = base.DefaultChainType
	}
	switch set.ChainType {
	case base.SQLChainType:
		return sqlChain(c.llm.GetOllamaLLM(), set)
	case base.RetrievalChainType:
		return retrievalChain(c.llm.GetOllamaLLM(), c.llm.GetStore(), set)
	case base.ConversationChainType:
		return conversationChain(c.llm.GetOllamaLLM(), c.llm.GetStore(), c.memory, set)
	case base.ToolChainType:
		return execChain(c.llm.GetOllamaLLM(), set)
	default:
		return nil, fmt.Errorf("chain type %s not support", set.ChainType)
	}
}

// 模型chain队列编排-多模型顺序处理
func chainQueue(chs ...chains.Chain) (chains.Chain, error) {
	return chains.NewSimpleSequentialChain(chs)

}

// 检索对话模型-知识检索
func retrievalChain(llm *ollama.LLM, emStore vectorstores.VectorStore, set base.ChatSettings) (chains.Chain, error) {
	options := []vectorstores.Option{
		vectorstores.WithScoreThreshold(set.MinResultScore),
	}

	return chains.NewRetrievalQAFromLLM(llm, vectorstores.ToRetriever(emStore, set.AmountOfResults, options...)), nil
}

// 记忆对话模型-支持配置知识检索
func conversationChain(llm *ollama.LLM, emStore vectorstores.VectorStore, memory schema.Memory, set base.ChatSettings) (chains.Chain, error) {
	stuffQAChain := chains.LoadStuffQA(llm)
	questionGeneratorChain := chains.LoadCondenseQuestionGenerator(llm)
	options := []vectorstores.Option{
		vectorstores.WithScoreThreshold(set.MinResultScore),
	}

	return chains.NewConversationalRetrievalQA(stuffQAChain,
		questionGeneratorChain,
		func() schema.Retriever {
			if emStore == nil {
				return nil
			}
			return vectorstores.ToRetriever(emStore, set.AmountOfResults, options...)
		}(),
		memory), nil
}

// 对话理解-工具调用模型
func execChain(llm *ollama.LLM, set base.ChatSettings) (chains.Chain, error) {
	llmTools := []tools.Tool{
		tools.Calculator{},
	}
	chain := agents.NewExecutor(
		agents.NewConversationalAgent(llm, llmTools),
		agents.WithParserErrorHandler(agents.NewParserErrorHandler(func(s string) string {
			log.Warn(s)
			return fmt.Sprintf("Parsing Error. %s", s)
		})),
		agents.WithMaxIterations(set.MaxIterations),
		agents.WithMemory(memory.NewConversationBuffer()))

	return chain, nil
}

// sql数据库处理模型
func sqlChain(llm *ollama.LLM, set base.ChatSettings) (chains.Chain, error) {
	return chains.NewSQLDatabaseChain(llm, set.AmountOfResults, &sqldatabase.SQLDatabase{
		Engine: func() sqldatabase.Engine {
			dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&maxAllowedPacket=%d",
				"root",
				"admin123",
				"127.0.0.1:3306",
				"tdp_plus",
				0,
			)
			engine, err := mysql.NewMySQL(dsn)
			if err != nil {
				return nil
			}
			return engine
		}(),
		SampleRowsNumber: 0,
	}), nil
}

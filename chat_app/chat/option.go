package main

import (
	"langchain_chat/chat_app/chat/llmx"
	"langchain_chat/third/langchaingo/schema"
)

type Option func(*Chain)

type Chain struct {
	llm    *llmx.LLM
	memory schema.Memory
}

func NewChain(opts ...Option) *Chain {
	s := Chain{
		memory: nil,
	}
	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

func WithLLM(llm *llmx.LLM) Option {
	return func(c *Chain) {
		c.llm = llm
	}
}

func WithMemory(memory schema.Memory) Option {
	return func(c *Chain) {
		c.memory = memory
	}
}

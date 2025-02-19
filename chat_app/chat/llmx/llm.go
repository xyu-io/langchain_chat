package llmx

import (
	"context"
	"errors"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"langchain_chat/third/langchaingo/documentloaders"
	"langchain_chat/third/langchaingo/embeddings"
	"langchain_chat/third/langchaingo/llms/ollama"
	"langchain_chat/third/langchaingo/schema"
	"langchain_chat/third/langchaingo/textsplitter"
	"langchain_chat/third/langchaingo/vectorstores"
	"os"
)

type LLM struct {
	model         string
	embedder      *embeddings.EmbedderImpl
	embedderModel string
	url           string
	contentSize   int
	llm           *ollama.LLM

	isUseStore bool
	store      vectorstores.VectorStore
	storeUrl   string
	storeType  StoreType
}

func NewLLM(opts ...Option) (*LLM, error) {
	s, err := applyLLMOptions(opts...)
	if err != nil {
		return nil, err
	}

	llm, err := ollama.New(ollama.WithModel(s.model),
		ollama.WithServerURL(s.url),
		ollama.WithRunnerNumCtx(s.contentSize),
		// ollama.WithHTTPClient(httputil.DebugHTTPClient),
	)
	if err != nil {
		return nil, err
	}

	s.llm = llm

	err = s.initEmbedder()
	if err != nil {
		return nil, err
	}
	if s.isUseStore {
		err = s.initStore()
		if err != nil {
			return nil, err
		}
	}
	return &s, err
}

func (l *LLM) GetOllamaLLM() *ollama.LLM {
	return l.llm
}

func (l *LLM) GetStore() vectorstores.VectorStore {
	return l.store
}

func (l *LLM) initEmbedder() error {
	embeder, err := NewEmbedder(l.embedderModel, l.url, l.contentSize)
	if err != nil {
		return err
	}

	l.embedder = embeder
	return nil
}

func (l *LLM) initStore() error {
	if l.storeType == "" {
		l.storeType = Milvus
	}
	var store vectorstores.VectorStore
	var err error
	switch l.storeType {
	case Chroma:
		if l.storeUrl == "" {
			l.storeUrl = DEFAULT_Chroma_URL
		}
		store, err = NewChromaStore(l.embedder, types.COSINE, l.storeUrl, Default_NameSpace)
		if err != nil {
			return err
		}
	default:
		if l.storeUrl == "" {
			l.storeUrl = DEFAULT_MILVUS_URL
		}
		store, err = NewMilvusStore(l.embedder, entity.L2, l.storeUrl, Default_Collection)
		if err != nil {
			return err
		}
	}

	l.store = store

	return nil
}

func (l *LLM) DocumentLoader(path string) error {
	fs, err := os.Open(path)
	if err != nil {
		return err
	}

	// 文本加载、分割处理
	docs, err := documentloaders.NewText(fs).LoadAndSplit(context.Background(),
		textsplitter.NewRecursiveCharacter(
			textsplitter.WithChunkSize(1000),
			textsplitter.WithChunkOverlap(20),
		))
	if err != nil {
		return err
	}

	err = l.addDocument(docs)
	if err != nil {
		return err
	}
	return nil
}

func (l *LLM) addDocument(docs []schema.Document) error {
	if len(docs) <= 0 {
		return errors.New("no documents")
	}

	_, err := l.store.AddDocuments(context.Background(), docs)
	if err != nil {
		return err
	}

	return nil
}

func newEmbedderWithLLM(llm *ollama.LLM) (*embeddings.EmbedderImpl, error) {
	return embeddings.NewEmbedder(llm)
}

func NewEmbedder(model string, serverURL string, ctxSize int) (*embeddings.EmbedderImpl, error) {
	llm, err := ollama.New(ollama.WithModel(model),
		ollama.WithServerURL(serverURL),
		ollama.WithRunnerNumCtx(ctxSize),
		// ollama.WithHTTPClient(httputil.DebugHTTPClient),
	)
	if err != nil {
		return nil, err
	}

	return embeddings.NewEmbedder(llm)

}

func NewLLMAndEmbedder(model string, serverURL string, ctxSize int) (*ollama.LLM, *embeddings.EmbedderImpl, error) {
	llm, err := NewLLM(
		WithModel(model),
		WithLLMUrl(serverURL),
		WithContentSize(ctxSize),
	)

	embedder, err := newEmbedderWithLLM(llm.llm)
	if err != nil {
		return llm.llm, nil, err
	}

	return llm.llm, embedder, nil
}

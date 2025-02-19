package llmx

const (
	DEFAULT_EMBEDDING_MODEL = "nomic-embed-text:latest"
	DEFAULT_LLM_MODEL       = "deepseek-r1:7b"
	DEFAULT_LLM_SERVER      = "http://127.0.0.1:11434"

	DEFAULT_CONTENT_SIZE = 8 * 1024
	DEFAULT_MILVUS_URL   = "http://192.168.200.222:19530"
	DEFAULT_Chroma_URL   = "http://192.168.200.222:8000"
	Default_Collection   = "milvus_collection"
	Default_NameSpace    = "chroma_space"
)

type StoreType string

const (
	Milvus             StoreType = "milvus"
	Chroma             StoreType = "chroma"
	DEFAULT_STORE_TYPE StoreType = Milvus
)

type Option func(p *LLM)
type Store func(p *LLM)

// WithModel sets the model of ai llm
func WithModel(model string) Option {
	return func(p *LLM) {
		p.model = model
	}
}

func WithEmbedderModel(model string) Option {
	return func(p *LLM) {
		p.embedderModel = model
	}
}

func WithStoreServer(url string) Store {
	return func(p *LLM) {
		p.storeUrl = url
	}
}

func WithStoreType(type_ StoreType) Store {
	return func(p *LLM) {
		p.storeType = type_
	}
}

func WithLLMUrl(url string) Option {
	return func(p *LLM) {
		p.url = url
	}
}

func WithContentSize(size int) Option {
	return func(p *LLM) {
		p.contentSize = size
	}
}

func UseStore(opts ...Store) Option {
	return func(p *LLM) {
		p.isUseStore = true
	}
}

func applyLLMOptions(opts ...Option) (LLM, error) {
	s := LLM{
		model:         DEFAULT_LLM_MODEL,
		url:           DEFAULT_LLM_SERVER,
		contentSize:   DEFAULT_CONTENT_SIZE,
		embedderModel: DEFAULT_EMBEDDING_MODEL,
		storeType:     DEFAULT_STORE_TYPE,
		storeUrl:      DEFAULT_MILVUS_URL,
	}
	for _, opt := range opts {
		opt(&s)
	}

	return s, nil
}

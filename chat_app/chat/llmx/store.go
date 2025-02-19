package llmx

import (
	"context"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	log "github.com/sirupsen/logrus"
	"langchain_chat/third/langchaingo/embeddings"
	"langchain_chat/third/langchaingo/vectorstores"
	"langchain_chat/third/langchaingo/vectorstores/chroma"
	"langchain_chat/third/langchaingo/vectorstores/milvus"
)

const (
	DefaultMilvusURL = "http://192.168.200.222:19530"
	DefaultChromaURL = "http://192.168.200.222:8000"
)

func NewMilvusStore(embedder embeddings.Embedder, metricType entity.MetricType, url, collection string) (vectorstores.VectorStore, error) {
	if metricType == "" {
		metricType = entity.L2
	}
	idx, _ := entity.NewIndexAUTOINDEX(metricType)
	opts := []milvus.Option{
		milvus.WithCollectionName(collection),
		milvus.WithEmbedder(embedder),
		milvus.WithIndex(idx),
	}

	milvusConfig := client.Config{
		Address: func() string {
			if url == "" {
				return DefaultMilvusURL
			}
			return url
		}(),
		Username: "root",
		Password: "Milvus",
	}
	// Create a new milvus vector store.
	return milvus.New(
		context.Background(),
		milvusConfig,
		opts...,
	)
}

func NewChromaStore(
	embedder embeddings.Embedder,
	distanceFunction types.DistanceFunction,
	chromaUrl string,
	nameSpace string,
) (vectorstores.VectorStore, error) {

	if chromaUrl == "" {
		chromaUrl = DefaultChromaURL
	}
	if distanceFunction == "" {
		distanceFunction = types.COSINE
	}

	store, err := chroma.New(
		chroma.WithChromaURL(chromaUrl),
		chroma.WithNameSpace(nameSpace),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction(distanceFunction), // default is cosine l2 ip
	)
	if err != nil {
		log.Error("add doc ", err)
		return nil, err
	}

	return store, nil
}

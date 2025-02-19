package documentloaders

import (
	"context"

	"langchain_chat/third/langchaingo/schema"
	"langchain_chat/third/langchaingo/textsplitter"
)

// Loader is the interface for loading and splitting documents from a source.
type Loader interface {
	// Load loads from a source and returns documents.
	Load(ctx context.Context) ([]schema.Document, error)
	// LoadAndSplit loads from a source and splits the documents using a text splitter.
	LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error)
}

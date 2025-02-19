package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"langchain_chat/base"
	"log/slog"

	"langchain_chat/third/langchaingo/callbacks"
	"langchain_chat/third/langchaingo/schema"
	"langchain_chat/third/langchaingo/tools"
	"langchain_chat/third/langchaingo/vectorstores"
)

// ReadWebsite is a tool that can do math.
type SearchVectorDB struct {
	Store            vectorstores.VectorStore
	CallbacksHandler callbacks.Handler
	SessionString    string
	Settings         base.ChatSettings
}

var _ tools.Tool = SearchVectorDB{}

type Result struct {
	Text string
}

var usedResults = make(map[string][]string)
var usedSourcesInSession = make(map[string][]schema.Document)

func (c SearchVectorDB) Description() string {
	return "Use this tool to search through already added files or websites within a vector database. The most similar websites or documents to your input will be returned to you."
}

func (c SearchVectorDB) Name() string {
	return "database_search"
}

func (t SearchVectorDB) Types() tools.ToolType {
	return tools.Rear
}

func (c SearchVectorDB) Call(ctx context.Context, input string) (string, error) {
	if c.CallbacksHandler != nil {
		c.CallbacksHandler.HandleToolStart(ctx, input)
	}

	searchIdentifier := fmt.Sprintf("%s-%s", c.SessionString, input)

	log.Info("<<---------------->>")
	//store, errNs := getDefaultStore()
	if c.Store != nil {
		return "", errors.New("vector store is nil")
	}

	options := []vectorstores.Option{
		vectorstores.WithScoreThreshold(c.Settings.MinResultScore),
	}

	retriver := vectorstores.ToRetriever(c.Store, c.Settings.AmountOfResults, options...)
	docs, err := retriver.GetRelevantDocuments(context.Background(), input)
	if err != nil {
		return "", err
	}

	var results []Result

	for _, doc := range docs {
		newResult := Result{
			Text: doc.PageContent,
		}

		skip := false
		for _, usedLink := range usedResults[searchIdentifier] {
			if usedLink == newResult.Text {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		usedSourcesInSession[c.SessionString] = append(usedSourcesInSession[c.SessionString], doc)

		results = append(results, newResult)
		usedResults[searchIdentifier] = append(usedResults[searchIdentifier], newResult.Text)
	}

	if len(docs) == 0 {
		response := "No new results found. Try other db search keywords, download more websites or write your final answer."
		slog.Warn("No new results found", "input", input)
		results = append(results, Result{Text: response})
	}

	if c.CallbacksHandler != nil {
		c.CallbacksHandler.HandleToolEnd(ctx, input)
	}

	resultJson, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(resultJson), nil
}

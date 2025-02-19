package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"langchain_chat/base"
	"log/slog"
	"strings"
	"unicode/utf8"

	"langchain_chat/third/langchaingo/llms"
	"langchain_chat/third/langchaingo/schema"
)

type CustomHandler struct {
	OutputChan chan<- base.HttpJsonStreamElement
}

func (l CustomHandler) HandleLLMGenerateContentStart(ctx context.Context, ms []llms.MessageContent) {
	l.LogDebug("Entering LLM with messages:")
	for _, m := range ms {
		var buf strings.Builder
		for _, t := range m.Parts {
			if t, ok := t.(llms.TextContent); ok {
				buf.WriteString(t.Text)
			}
		}
		l.LogDebug(fmt.Sprintf("Role: %s", m.Role))
		l.LogDebug(fmt.Sprintf("Text: %s", buf.String()))
	}
}

func (l CustomHandler) HandleLLMGenerateContentEnd(ctx context.Context, res *llms.ContentResponse) {
	fmt.Println("Exiting LLM with response:")
	for _, c := range res.Choices {
		if c.Content != "" {
			l.LogDebug(fmt.Sprintf("Content: %s", c.Content))
		}
		if c.StopReason != "" {
			l.LogDebug(fmt.Sprintf("StopReason: %s", c.StopReason))
		}
		if len(c.GenerationInfo) > 0 {
			text := ""
			text += fmt.Sprintf("GenerationInfo: ")
			for k, v := range c.GenerationInfo {
				text += fmt.Sprintf("%20s: %v\n", k, v)
			}
			l.LogDebug(text)
		}
		if c.FuncCall != nil {
			l.LogDebug(fmt.Sprintf("FuncCall: %s %s", c.FuncCall.Name, c.FuncCall.Arguments))
		}
	}
}

func (l CustomHandler) LogDebug(text string) {
	l.OutputChan <- base.HttpJsonStreamElement{
		Message: text,
		Stream:  false,
	}
}

func (l CustomHandler) HandleStreamingFunc(ctx context.Context, chunk []byte) {
	l.OutputChan <- base.HttpJsonStreamElement{
		Message:  string(chunk),
		Stream:   true,
		StepType: base.StepHandleStreaming,
	}
}

func (l CustomHandler) HandleText(ctx context.Context, text string) {
	// 监测用户输入文本
	//l.OutputChan <- HttpJsonStreamElement{
	//	Message: text,
	//	Stream:  false,
	//}
}

func (l CustomHandler) HandleLLMStart(ctx context.Context, prompts []string) {
	l.OutputChan <- base.HttpJsonStreamElement{
		Message:  fmt.Sprintf("Entering LLM with prompts: %s", prompts),
		Stream:   false,
		StepType: base.StepHandleLlmStart,
	}
}

func (l CustomHandler) HandleLLMError(_ context.Context, err error) {
	fmt.Println("Exiting LLM with error:", err)
	l.OutputChan <- base.HttpJsonStreamElement{
		Message: err.Error(),
		Stream:  false,
	}
}

func (l CustomHandler) HandleChainStart(ctx context.Context, inputs map[string]any) {
	chainValuesJson, err := json.Marshal(inputs)
	if err != nil {
		fmt.Println("Error marshalling chain values:", err)
	}

	charCount := utf8.RuneCountInString(string(chainValuesJson))
	slog.Info("Entering chain", "tokens", charCount/4)

	var input struct {
		Input string `json:"input"`
	}
	err = json.Unmarshal(chainValuesJson, &input)
	if err == nil {
		l.OutputChan <- base.HttpJsonStreamElement{
			Message:  fmt.Sprintf("【收到用户会话-%d Tokens)】: %s", charCount/4, input.Input),
			Stream:   false,
			StepType: base.StepHandleChainStart,
		}
	}

}

func (l CustomHandler) HandleChainEnd(ctx context.Context, outputs map[string]any) {
	chainValuesJson, err := json.Marshal(outputs)
	if err != nil {
		fmt.Println("Error marshalling chain values:", err)
	}
	var out struct {
		Output string `json:"output"`
	}
	err = json.Unmarshal(chainValuesJson, &out)
	if err == nil {
		l.OutputChan <- base.HttpJsonStreamElement{
			Message:  fmt.Sprintf("【输出】: %s", out.Output),
			Stream:   false,
			StepType: base.StepHandleChainEnd,
		}
	}
}

func (l CustomHandler) HandleChainError(ctx context.Context, err error) {
	message := fmt.Sprintf("Exiting chain with error: %v", err)
	l.OutputChan <- base.HttpJsonStreamElement{
		Message:  message,
		Stream:   false,
		StepType: base.StepHandleChainError,
	}
}

func (l CustomHandler) HandleToolStart(ctx context.Context, input string) {
	l.OutputChan <- base.HttpJsonStreamElement{
		Message:  fmt.Sprintf("Entering tool with input: %s", removeNewLines(input)),
		Stream:   false,
		StepType: base.StepHandleToolStart,
	}
}

func (l CustomHandler) HandleToolEnd(ctx context.Context, output string) {
	l.OutputChan <- base.HttpJsonStreamElement{
		Message:  fmt.Sprintf("Exiting tool with output: %s", removeNewLines(output)),
		Stream:   false,
		StepType: base.StepHandleToolEnd,
	}
}

func (l CustomHandler) HandleToolError(ctx context.Context, err error) {
	l.OutputChan <- base.HttpJsonStreamElement{
		Message: err.Error(),
		Stream:  false,
	}
}

func (l CustomHandler) HandleAgentAction(ctx context.Context, action schema.AgentAction) {
	actionJson, err := json.Marshal(action)
	if err != nil {
		fmt.Println("Error marshalling action:", err)
	}

	l.OutputChan <- base.HttpJsonStreamElement{
		Message:  fmt.Sprintf("【建议】: %s", string(actionJson)),
		Stream:   false,
		StepType: base.StepHandleAgentAction,
	}

}

func (l CustomHandler) HandleAgentFinish(ctx context.Context, finish schema.AgentFinish) {
	finishJson, err := json.Marshal(finish)
	if err != nil {
		fmt.Println("Error marshalling finish:", err)
	}
	var think base.DeepSeekThinkResp
	err = json.Unmarshal(finishJson, &think)
	if err == nil {
		l.OutputChan <- base.HttpJsonStreamElement{
			Message:  fmt.Sprintf("【思考】%s", think.Log),
			Stream:   false,
			StepType: base.StepHandleAgentAction,
		}
	} else {
		l.OutputChan <- base.HttpJsonStreamElement{
			Message:  string(finishJson),
			Stream:   false,
			StepType: base.StepHandleAgentFinish,
		}
	}
}

func (l CustomHandler) HandleRetrieverStart(ctx context.Context, query string) {
	fmt.Println("Entering retriever with query:", removeNewLines(query))
}

func (l CustomHandler) HandleRetrieverEnd(ctx context.Context, query string, documents []schema.Document) {
	// fmt.Println("Exiting retriever with documents for query:", documents, query)
	l.OutputChan <- base.HttpJsonStreamElement{
		Message:  fmt.Sprintf("Exiting retriever with documents for query: %s", query),
		Stream:   false,
		StepType: base.StepHandleRetriverEnd,
	}
}

func (l CustomHandler) HandleVectorFound(ctx context.Context, vectorString string) {
	l.OutputChan <- base.HttpJsonStreamElement{
		Message:  fmt.Sprintf("Found vector %s", vectorString),
		Stream:   false,
		StepType: base.StepHandleVectorFound,
	}
}

func (l CustomHandler) HandleSourceAdded(ctx context.Context, source base.Source) {
	l.OutputChan <- base.HttpJsonStreamElement{
		Message:  "Source added",
		Source:   source,
		Stream:   false,
		StepType: base.StepHandleSourceAdded,
	}
}

func formatChainValues(values map[string]any) string {
	output := ""
	for key, value := range values {
		output += fmt.Sprintf("\"%s\" : \"%s\", ", removeNewLines(key), removeNewLines(value))
	}

	return output
}

func formatAgentAction(action schema.AgentAction) string {
	return fmt.Sprintf("\"%s\" with input \"%s\"", removeNewLines(action.Tool), removeNewLines(action.ToolInput))
}

func removeNewLines(s any) string {
	return strings.ReplaceAll(fmt.Sprint(s), "\n", " ")
}

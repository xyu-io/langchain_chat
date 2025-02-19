package agents

import (
	"context"
	_ "embed"
	"fmt"
	"langchain_chat/third/langchaingo/outputparser"
	"regexp"
	"strings"

	"langchain_chat/third/langchaingo/callbacks"
	"langchain_chat/third/langchaingo/chains"
	"langchain_chat/third/langchaingo/llms"
	"langchain_chat/third/langchaingo/prompts"
	"langchain_chat/third/langchaingo/schema"
	"langchain_chat/third/langchaingo/tools"
)

const (
	_deepSeekFinalAnswerAction = "AI:"
	//_conversationalNoThinkAnswerAction        = "我建议您联网获取时效性较强的信息"
	//_conversationalThinkAnswerActionRegx 	= `<think>(.*?)</think>`
)

// DeepSeekAgent is a struct that represents an agent responsible for DeepSeek
type DeepSeekAgent struct {
	// Chain is the chain used to call with the values. The chain should have an
	// input called "agent_scratchpad" for the agent to put its thoughts in.
	Chain chains.Chain
	// Tools is a list of the tools the agent can use.
	Tools []tools.Tool
	// Output key is the key where the final output is placed.
	OutputKey string
	// CallbacksHandler is the handler for callbacks.
	CallbacksHandler callbacks.Handler
}

var _ Agent = (*DeepSeekAgent)(nil)

func NewDeepSeekAgent(llm llms.Model, tools []tools.Tool, opts ...Option) *DeepSeekAgent {
	options := conversationalDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	chain := chains.NewLLMChain(
		llm,
		options.getDeepSeekPrompt(tools),
		chains.WithCallback(options.callbacksHandler),
	)
	chain.OutputParser = outputparser.NewStructured([]outputparser.ResponseSchema{
		{
			Name:        "action",
			Description: "action to take",
		},
	})
	return &DeepSeekAgent{
		Chain:            chain,
		Tools:            tools,
		OutputKey:        options.outputKey,
		CallbacksHandler: options.callbacksHandler,
	}
}

// Plan decides what action to take or returns the final result of the input.
func (a *DeepSeekAgent) Plan(
	ctx context.Context,
	intermediateSteps []schema.AgentStep,
	inputs map[string]string,
) ([]schema.AgentAction, *schema.AgentFinish, error) {
	fullInputs := make(map[string]any, len(inputs))
	for key, value := range inputs {
		fullInputs[key] = value
	}

	fullInputs["agent_scratchpad"] = constructScratchPad(intermediateSteps)

	var stream func(ctx context.Context, chunk []byte) error

	if a.CallbacksHandler != nil {
		stream = func(ctx context.Context, chunk []byte) error {
			a.CallbacksHandler.HandleStreamingFunc(ctx, chunk)
			return nil
		}
	}

	// 具体调用实现-对话输出/检索输出
	output, err := chains.Predict(
		ctx,
		a.Chain,
		fullInputs,
		chains.WithStopWords([]string{"\nObservation:", "\n\tObservation:"}),
		chains.WithStreamingFunc(stream),
	)
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil

	return a.parseOutput(output)
}

func (a *DeepSeekAgent) GetInputKeys() []string {
	chainInputs := a.Chain.GetInputKeys()

	// Remove inputs given in plan.
	agentInput := make([]string, 0, len(chainInputs))
	for _, v := range chainInputs {
		if v == "agent_scratchpad" {
			continue
		}
		agentInput = append(agentInput, v)
	}

	return agentInput
}

func (a *DeepSeekAgent) GetOutputKeys() []string {
	return []string{a.OutputKey}
}

func (a *DeepSeekAgent) GetTools() []tools.Tool {
	return a.Tools
}

func constructer(steps []schema.AgentStep) string {
	var scratchPad string
	if len(steps) > 0 {
		for _, step := range steps {
			scratchPad += step.Action.Log
			scratchPad += "\nObservation: " + step.Observation
		}
		scratchPad += "\n" + "Thought:"
	}

	return scratchPad
}

func (a *DeepSeekAgent) parseOutput(output string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	if strings.Contains(output, _deepSeekFinalAnswerAction) {
		splits := strings.Split(output, _deepSeekFinalAnswerAction)

		finishAction := &schema.AgentFinish{
			ReturnValues: map[string]any{
				a.OutputKey: splits[len(splits)-1],
			},
			Log: output,
		}

		return nil, finishAction, nil
	}

	r := regexp.MustCompile(`Action: (.*?)[\n]*Action Input: (.*)`)
	matches := r.FindStringSubmatch(output)
	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("%w: %v", ErrUnableToParseOutput, output)
	}

	return []schema.AgentAction{
		{Tool: strings.TrimSpace(matches[1]), ToolInput: strings.TrimSpace(matches[2]), Log: output},
	}, nil, nil
}

func createDeepSeekPrompt(tools []tools.Tool, prefix, instructions, suffix string) prompts.PromptTemplate {
	template := strings.Join([]string{prefix, instructions, suffix}, "\n\n")

	return prompts.PromptTemplate{
		Template:       template,
		TemplateFormat: prompts.TemplateFormatGoTemplate,
		InputVariables: []string{"input", "agent_scratchpad"},
		PartialVariables: map[string]any{
			"tool_names":        toolNames(tools),
			"tool_descriptions": toolDescriptions(tools),
			"history":           "",
		},
	}
}

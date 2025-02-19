package base

type AIChainType string

const (
	//WebChainType          AIChainType = "web"
	SQLChainType          AIChainType = "sql"
	ToolChainType         AIChainType = "tool"
	ConversationChainType AIChainType = "conversation"
	RetrievalChainType    AIChainType = "retrieval"
	DefaultChainType      AIChainType = "conversation"
)

type ChatSettings struct {
	ChainType           AIChainType `json:"chainType"`
	ContextSize         int         `json:"contextSize"`
	MaxIterations       int         `json:"maxIterations"`
	ModelName           string      `json:"modelName"`
	Prompt              []string    `json:"prompt"`
	Session             string      `json:"session"`
	Temperature         float64     `json:"temperature"`
	ToolNames           []string    `json:"toolNames"`
	WebSearchCategories []string    `json:"webSearchCategories"`
	AmountOfResults     int         `json:"amountOfResults"`
	MinResultScore      float32     `json:"minResultScore"`
	AmountOfWebsites    int         `json:"amountOfWebsites"`
	ChunkSize           int         `json:"chunkSize"`
	ChunkOverlap        int         `json:"chunkOverlap"`
	SystemMessage       string      `json:"systemMessage"`
}

type HttpJsonStreamElement struct {
	Message   string   `json:"message"`
	Close     bool     `json:"close"`
	Stream    bool     `json:"stream"`
	StepType  StepType `json:"stepType"`
	Source    Source   `json:"source"`
	Session   string   `json:"session"`
	TimeStamp int64    `json:"timeStamp"`
}

type StepType string

type Source struct {
	Name    string `json:"name"`
	Link    string `json:"link"`
	Summary string `json:"summary"`
	Engine  string `json:"engine"`
	Title   string `json:"title"`
}

type DeepSeekThinkResp struct {
	ReturnValues ReturnInfos `json:"ReturnValues"`
	Log          string      `json:"Log"`
}

type ReturnInfos struct {
	Output string `json:"output"`
}

const (
	StepHandleAgentAction             StepType = "HandleAgentAction"
	StepHandleAgentFinish             StepType = "HandleAgentFinish"
	StepHandleChainEnd                StepType = "HandleChainEnd"
	StepHandleChainError              StepType = "HandleChainError"
	StepHandleChainStart              StepType = "HandleChainStart"
	StepHandleFinalAnswer             StepType = "HandleFinalAnswer"
	StepHandleLLMGenerateContentEnd   StepType = "HandleLLMGenerateContentEnd"
	StepHandleLLMGenerateContentStart StepType = "HandleLLMGenerateContentStart"
	StepHandleLlmEnd                  StepType = "HandleLlmEnd"
	StepHandleLlmError                StepType = "HandleLlmError"
	StepHandleLlmStart                StepType = "HandleLlmStart"
	StepHandleNewSession              StepType = "HandleNewSession"
	StepHandleOllamaStart             StepType = "HandleOllamaStart"
	StepHandleParseError              StepType = "HandleParseError"
	StepHandleRetriverEnd             StepType = "HandleRetriverEnd"
	StepHandleRetriverStart           StepType = "HandleRetriverStart"
	StepHandleSourceAdded             StepType = "HandleSourceAdded"
	StepHandleToolEnd                 StepType = "HandleToolEnd"
	StepHandleToolError               StepType = "HandleToolError"
	StepHandleToolStart               StepType = "HandleToolStart"
	StepHandleVectorFound             StepType = "HandleVectorFound"
	StepHandleFormat                  StepType = "HandleFormat"
	StepHandleStreaming               StepType = "HandleStreaming"
	StepHandleUserMessage             StepType = "HandleUserMessage"
)

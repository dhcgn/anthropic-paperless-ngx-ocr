package types

type Source struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type RequestContent struct {
	Type   string  `json:"type"`
	Source *Source `json:"source,omitempty"`
	Text   string  `json:"text,omitempty"`
}

type Message struct {
	Role    string           `json:"role"`
	Content []RequestContent `json:"content"`
}

type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"input_schema"`
}

type Payload struct {
	Model      string      `json:"model"`
	MaxTokens  int         `json:"max_tokens"`
	Messages   []Message   `json:"messages"`
	Tools      *[]Tool     `json:"tools,omitempty"`
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"`
}

type ToolChoice struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

package anthropictypes

type ResponseContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Usage struct {
	InputTokens              int `json:"input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	OutputTokens             int `json:"output_tokens"`
}

type Response struct {
	ID           string           `json:"id"`
	Type         string           `json:"type"`
	Role         string           `json:"role"`
	Model        string           `json:"model"`
	Content      []RequestContent `json:"content"`
	StopReason   string           `json:"stop_reason"`
	StopSequence *string          `json:"stop_sequence"`
	Usage        Usage            `json:"usage"`
}
package anthropictypes

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

type Payload struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

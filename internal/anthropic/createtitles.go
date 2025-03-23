package anthropic

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"anthropicpaperocr/internal/anthropic/types"
)

//go:embed prompts/createtitles.md
var createTitlesPrompt string

// CreateTitles creates titles from the content using the Anthropic API.
// The generative AI model will generate titles based on the content provided.
// It returns a slice of strings containing the titles and an error if any.
func CreateTitles(content, oldtitle string, apiKeyAnthropic string) ([]string, error) {
	url := "https://api.anthropic.com/v1/messages"

	textContent := createTitlesPrompt
	textContent = strings.Replace(textContent, "{{CONTENT}}", content, 1)
	textContent = strings.Replace(textContent, "{{OLD_TITLE}}", oldtitle, 1)

	payload := types.Payload{
		Model:     "claude-3-7-sonnet-latest",
		MaxTokens: 1024 * 8,
		Messages: []types.Message{
			{
				Role: "user",
				Content: []types.RequestContent{
					{
						Type: "text",
						Text: textContent,
					},
				},
			},
		},
		Tools: &[]types.Tool{
			{
				Name:        "generate_titles",
				Description: "Generate a list of titles from the provided content using well-structured JSON.",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"titles": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
							"description": "List of titles generated from the content.",
						},
					},
					Required: []string{"titles"},
				},
			},
		},
		ToolChoice: &types.ToolChoice{
			Type: "tool",
			Name: "generate_titles",
		},
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Error creating request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKeyAnthropic)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to perform comparison, status code: %d, response body: %s", resp.StatusCode, string(body))
	}

	var response types.Response

	// os.WriteFile(`C:\dev\ai-claude-paperless-ngx-pdf-visual-ocr\debugging\CreateTitles.txt`, body, 0644)

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("Error unmarshaling response: %v", err)
	}

	if len(response.Content) == 0 || response.Content[0].Type != "tool_use" {
		return nil, fmt.Errorf("Unexpected response format")
	}

	titlesInterface, ok := response.Content[0].ToolResult["titles"]
	if !ok {
		return nil, fmt.Errorf("Titles not found in the response")
	}

	titlesSlice, ok := titlesInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Titles are not in the expected format")
	}

	titles := make([]string, len(titlesSlice))
	for i, title := range titlesSlice {
		titles[i], ok = title.(string)
		if !ok {
			return nil, fmt.Errorf("Title is not a string")
		}
	}
	// The response contains the generated titles

	return titles, err
}

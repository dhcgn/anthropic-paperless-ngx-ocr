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

	"github.com/sergi/go-diff/diffmatchpatch"
)

//go:embed prompts/compare.md
var comparePrompt string

func CompareContent(originalContent, newContent, apiKeyAnthropic string) (string, string, error) {
	// Generate diff using diffmatchpatch library
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(originalContent, newContent, false)
	diff := dmp.DiffPrettyText(diffs)

	// Generate AI comparison
	aiComparison, err := generateAIComparison(originalContent, newContent, apiKeyAnthropic)
	if err != nil {
		return "", "", fmt.Errorf("Error generating AI comparison: %v", err)
	}

	return diff, aiComparison, nil
}

func generateAIComparison(originalContent, newContent, apiKeyAnthropic string) (string, error) {
	url := "https://api.anthropic.com/v1/messages"

	textContent := comparePrompt
	textContent = strings.Replace(textContent, "{{TRANSCRIPT_OLD}}", originalContent, 1)
	textContent = strings.Replace(textContent, "{{TRANSCRIPT_NEW}}", newContent, 1)

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
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Error creating request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKeyAnthropic)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to perform comparison, status code: %d, response body: %s", resp.StatusCode, string(body))
	}

	var response types.Response

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("Error unmarshaling response: %v", err)
	}

	if len(response.Content) == 0 || response.Content[0].Type != "text" {
		return "", fmt.Errorf("Unexpected response format")
	}

	return response.Content[0].Text, nil
}

package compare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"anthropicpaperocr/internal/anthropictypes"

	"github.com/sergi/go-diff/diffmatchpatch"
)

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

	textContent := fmt.Sprintf(`
You are tasked with comparing two transcripts of the same document and providing a recommendation on which one is better. Your analysis should be thorough yet concise, focusing on key differences that impact the overall quality of the transcripts.

Here is the first transcript:
<transcript_OLD>
%s
</transcript_OLD>

Here is the second transcript:
<transcript_NEW>
%s
</transcript_NEW>

Please analyze both transcripts, considering the following factors:

1. Accuracy: Compare the transcripts for correctness of words, phrases, and overall content.
2. Completeness: Assess which transcript captures more of the original document's content.
3. Clarity and coherence: Evaluate the readability and flow of each transcript.
4. Any additional factors that you find relevant to determining the quality of the transcripts.

Based on your analysis, provide a recommendation on which transcript is better. Your recommendation should be well-founded and supported by specific examples from the transcripts.

Present your findings and recommendation in the following format:

<analysis>
[Your detailed analysis of the transcripts, comparing them based on the factors mentioned above. Include specific examples to support your points.]
</analysis>

<recommendation>
[Your concise recommendation of which transcript is better, along with a brief justification.]
</recommendation>

Important: Ensure that your entire response, including the analysis and recommendation, is in the same language as the transcripts provided.
`, originalContent, newContent)

	payload := anthropictypes.Payload{
		Model:     "claude-3-7-sonnet-latest",
		MaxTokens: 1024 * 8,
		Messages: []anthropictypes.Message{
			{
				Role: "user",
				Content: []anthropictypes.RequestContent{
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

	var response anthropictypes.Response

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("Error unmarshaling response: %v", err)
	}

	if len(response.Content) == 0 || response.Content[0].Type != "text" {
		return "", fmt.Errorf("Unexpected response format")
	}

	return response.Content[0].Text, nil
}

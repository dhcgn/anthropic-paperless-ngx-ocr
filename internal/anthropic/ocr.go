package anthropic

import (
	"anthropicpaperocr/internal/anthropic/types"
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//go:embed prompts/ocr.md
var ocrPrompt string

func PerformOCR(pdfData []byte, apiKeyAnthropic string) (string, error) {
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfData)

	textContent := ocrPrompt
	textContent = strings.Replace(textContent, "{{PDF_DATA}}", pdfBase64, 1)

	payload := types.Payload{
		Model:     "claude-3-7-sonnet-latest",
		MaxTokens: 64_000,
		Messages: []types.Message{
			{
				Role: "user",
				Content: []types.RequestContent{
					{
						Type: "document",
						Source: &types.Source{
							Type:      "base64",
							MediaType: "application/pdf",
							Data:      pdfBase64,
						},
					},
					{
						Type: "text",
						Text: textContent,
					},
				},
			},
		},
	}

	tempJSONFile := bytes.NewBuffer(nil)
	if err := json.NewEncoder(tempJSONFile).Encode(payload); err != nil {
		return "", err
	}

	// write to to disk for debugging
	// ioutil.WriteFile(`C:\dev\ai-claude-paperless-ngx-pdf-visual-ocr\debugging\temp.json`, tempJSONFile.Bytes(), 0644)

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", tempJSONFile)
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", apiKeyAnthropic)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to perform OCR, status code: %d, corresponding body: %v ", resp.StatusCode, string(body))
	}

	var response types.Response

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Content) == 0 || response.Content[0].Type != "text" {
		return "", fmt.Errorf("unexpected response format")
	}

	return response.Content[0].Text, nil
}

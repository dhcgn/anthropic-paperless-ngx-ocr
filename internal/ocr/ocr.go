package ocr

import (
	"anthropicpaperocr/internal/anthropictypes"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PerformOCR(pdfData []byte, apiKeyAnthropic string) (string, error) {
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfData)

	textContent := `You are an advanced text extraction and transcription system designed to accurately process PDF files. Your task is to transcribe the content of the given PDF file(s), including descriptions of non-text elements.

Please follow these instructions carefully:

1. Read through the entire content of the PDF file(s).

2. Transcribe the text content as accurately as possible, maintaining the original language of the document.

3. For any non-text elements (images, symbols, diagrams, etc.), create descriptive placeholders using square brackets. For example:
   - [Company Logo]
   - [Bar Chart showing quarterly sales figures]
   - [Signature of John Doe, CEO]

4. Insert these placeholders near the relevant text in the document, maintaining the logical flow and layout of the original.

5. If any text is unreadable or unclear, make an educated guess based on the surrounding context. Ensure your guess is plausible and fits with the overall document.

6. Maintain all formatting, paragraph breaks, and section divisions as they appear in the original document.

7. Do not add any commentary, explanations, or additional text beyond the transcription itself.

Before providing the final transcription, ensure all requirements are met. Consider the following steps mentally:

- List out a step-by-step approach for handling different types of content (text, images, tables).
- Identify potential challenges in the transcription process and propose solutions for each.
- Write down 3-5 examples of how you will handle unclear text, showing your thought process.
- Describe your strategy for ensuring placeholders are appropriately placed and described.

It's OK for this section to be quite long.

Once you've completed your mental  planning, provide the final transcription. The transcription should contain only the extracted content, including placeholders, with no additional commentary.

Provide only the final transcription, as your result will be considered as the content of the document and nothing else.`

	payload := anthropictypes.Payload{
		Model:     "claude-3-5-sonnet-20241022",
		MaxTokens: 8192,
		Messages: []anthropictypes.Message{
			{
				Role: "user",
				Content: []anthropictypes.RequestContent{
					{
						Type: "document",
						Source: &anthropictypes.Source{
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
	ioutil.WriteFile("temp.json", tempJSONFile.Bytes(), 0644)

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to perform OCR, status code: %d, corresponding body: %v ", resp.StatusCode, string(body))
	}

	var response anthropictypes.Response

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Content) == 0 || response.Content[0].Type != "text" {
		return "", fmt.Errorf("unexpected response format")
	}

	return response.Content[0].Text, nil
}

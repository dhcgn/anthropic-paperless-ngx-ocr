package paperless

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Document struct {
	Content string `json:"content"`
}

func GetCurrentContent(documentID int, apiKey, url, hostHeader string) (string, error) {
	fullURL := fmt.Sprintf("%s/api/documents/%d/", url, documentID)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Token "+apiKey)
	if hostHeader != "" {
		req.Host = hostHeader
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get document content, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var document Document
	err = json.Unmarshal(body, &document)
	if err != nil {
		return "", err
	}

	return document.Content, nil
}

func SetContent(documentID int, content, apiKey, url, hostHeader string) error {
	fullURL := fmt.Sprintf("%s/api/documents/%d/", url, documentID)
	payload := map[string]string{
		"content": content,
	}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error creating request body: %v", err)
	}

	req, err := http.NewRequest("PATCH", fullURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+apiKey)
	if hostHeader != "" {
		req.Host = hostHeader
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to set document content, status code: %d", resp.StatusCode)
	}

	return nil
}

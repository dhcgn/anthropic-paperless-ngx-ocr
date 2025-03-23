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
	Title   string `json:"title"`
}

func GetFrontendURL(documentID int, url string) string {
	return fmt.Sprintf("%s/document/%d/", url, documentID)
}

func GetCurrentDocument(documentID int, apiKey, url, hostHeader string) (*Document, error) {
	fullURL := fmt.Sprintf("%s/api/documents/%d/", url, documentID)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+apiKey)
	if hostHeader != "" {
		req.Host = hostHeader
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get document content, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var document Document
	err = json.Unmarshal(body, &document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func SetContent(documentID int, content, apiKey, url, hostHeader string) error {
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}
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

func SetTitle(documentID int, title, apiKey, url, hostHeader string) error {
	if title == "" {
		return fmt.Errorf("Title cannot be empty")
	}

	fullURL := fmt.Sprintf("%s/api/documents/%d/", url, documentID)
	payload := map[string]string{
		"title": title,
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

package paperless

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func DownloadPDF(documentID int, apiKey, url, hostHeader string) ([]byte, error) {
	fullURL := fmt.Sprintf("%s/api/documents/%d/download/", url, documentID)
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
		return nil, fmt.Errorf("failed to download PDF, status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

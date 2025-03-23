//go:build integration
// +build integration

package anthropic

import (
	"os"
	"testing"
)

func setapikeyfromlocalfile() error {
	filepath := `C:\dev\ai-claude-paperless-ngx-pdf-visual-ocr\secrets\api_key_anthropic.txt`
	filecontent, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return os.Setenv("ANTHROPIC_API_KEY", string(filecontent))

}

func TestCreateTitlesIntegration(t *testing.T) {
	if err := setapikeyfromlocalfile(); err != nil {
		t.Fatalf("setapikeyfromlocalfile() error = %v", err)
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test; ANTHROPIC_API_KEY not set")
	}

	content := "You are fired from your job at the company XYZ. Reason: Insubordination. You have been warned multiple times about your behavior, but you have not changed. You are hereby terminated from your position effective immediately. You will receive your final paycheck in the next pay period. You are required to return all company property, including your badge and laptop, before you leave the premises. You will be escorted out by security. If you have any questions, please contact HR."
	oldtitle := "Termination Letter"

	titles, err := CreateTitles(content, oldtitle, apiKey)
	if err != nil {
		t.Fatalf("CreateTitles() error = %v", err)
	}
	if len(titles) == 0 {
		t.Errorf("CreateTitles() returned no titles, expected at least one")
	}

	// print to console
	for i, title := range titles {
		t.Logf("Title %d: %s", i+1, title)
	}
}

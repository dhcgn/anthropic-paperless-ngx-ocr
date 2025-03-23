package main

import (
	"anthropicpaperocr/internal/compare"
	"anthropicpaperocr/internal/ocr"
	"anthropicpaperocr/internal/paperless"
	"flag"
	"fmt"

	"github.com/pterm/pterm"
)

const (
	AppName = "AnthropicPaperOCR"
)

var (
	Version = "0.0.2"
)

func main() {
	fmt.Println(AppName, Version)

	// Define flags
	documentID := flag.Int("document_id", 0, "ID of the document")
	apiKeyPaperless := flag.String("api_key_paperless", "", "API key for authentication")
	apiKeyAnthropic := flag.String("api_key_anthropic", "", "API key for authentication")
	hostHeader := flag.String("host_header", "", "Host header for the Paperless instance, if different from the URL")
	url := flag.String("url", "", "The URL for the Paperless instance")

	// Parse flags
	flag.Parse()

	// Validate flags
	if *documentID == 0 || *apiKeyPaperless == "" || *apiKeyAnthropic == "" || *url == "" {
		fmt.Println("All flags -document_id, -api_key_paperless, -api_key_anthropic, and -url are required")
		flag.Usage()
		return
	}

	// Progress bar
	p, _ := pterm.DefaultProgressbar.WithTotal(4).Start()

	// Get current content of document
	p.UpdateTitle("Getting current content of document...")
	currentContent, err := paperless.GetCurrentContent(*documentID, *apiKeyPaperless, *url, *hostHeader)
	if err != nil {
		fmt.Println("Error getting current content:", err)
		return
	}
	p.Increment()

	p.UpdateTitle("Downloading PDF...")
	// Download PDF file in memory
	pdfData, err := paperless.DownloadPDF(*documentID, *apiKeyPaperless, *url, *hostHeader)
	if err != nil {
		fmt.Println("Error downloading PDF:", err)
		return
	}
	p.Increment()

	// Perform OCR
	p.UpdateTitle("Performing OCR...")
	ocrResult, err := ocr.PerformOCR(pdfData, *apiKeyAnthropic)
	if err != nil {
		fmt.Println("Error performing OCR:", err)
		return
	}
	p.Increment()

	// Compare old and new content
	p.UpdateTitle("Comparing content...")
	diff, aiComparison, err := compare.CompareContent(currentContent, ocrResult, *apiKeyAnthropic)
	if err != nil {
		fmt.Println("Error comparing content:", err)
		return
	}
	p.Increment()

	// Display diff and AI comparison
	pterm.DefaultHeader.Println("Diff between original and new content")
	pterm.Info.Println(diff)
	//	fmt.Println(diff)
	pterm.DefaultHeader.Println("AI generated comparison:")
	pterm.Info.Println(aiComparison)
	// fmt.Println(aiComparison)

	p.Stop()

	// Prompt user to set new content
	pterm.DefaultHeader.Println("Next:")
	decision, _ := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Do you want to set the new content in the paperless instance? (yes/no): ")

	if decision == "yes" {
		fmt.Println("Setting the new content in the paperless instance...")
		err := paperless.SetContent(*documentID, ocrResult, *apiKeyPaperless, *url, *hostHeader)
		if err != nil {
			fmt.Println("Error setting new content:", err)
			return
		}
		fmt.Println("The new content has been set in the paperless instance.")
	} else {
		fmt.Println("The new content will not be set in the paperless instance.")
	}

}

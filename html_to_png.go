package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// Path to the HTML file
	htmlFile := "trainmap_table.html"
	outputPng := "trainmap_table_only.png"

	// Read the HTML content from the file (optional, for validation)
	_, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		log.Fatalf("Failed to read HTML file: %v", err)
	}

	// Create a new context for chromedp
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set up options to avoid timeout issues and improve stability
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Define the variable to hold the screenshot data
	var buf []byte

	// Convert the relative file path to an absolute path
	absPath, err := filepath.Abs(htmlFile)
	if err != nil {
		log.Fatalf("Failed to get absolute path of the HTML file: %v", err)
	}

	// Run chromedp tasks to open the local HTML file directly and capture the screenshot
	err = chromedp.Run(ctx, chromedp.Tasks{
		// Set the viewport to the desired size (1200x820)
		// chromedp.EmulateViewport(1200, 820),

		// Navigate to the local HTML file using the file:// protocol
		chromedp.Navigate("file://" + absPath),

		// Wait until the #loader element is visible (to ensure the page is fully loaded)
		chromedp.WaitVisible(`#loader`, chromedp.ByQuery),

		// Apply grayscale filter using JavaScript
		chromedp.Evaluate(`document.body.style.filter = 'grayscale(100%)';`, nil),

		// Capture the #root element as a screenshot
		chromedp.Screenshot(`#root`, &buf, chromedp.ByQuery),
	})

	if err != nil {
		log.Fatalf("Failed to capture table screenshot: %v", err)
	}

	// Write the screenshot to a PNG file
	err = ioutil.WriteFile(outputPng, buf, 0644)
	if err != nil {
		log.Fatalf("Failed to write PNG file: %v", err)
	}

	fmt.Printf("Table screenshot saved as %s\n", outputPng)
}

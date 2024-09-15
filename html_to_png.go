package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// Path to the HTML file
	htmlFile := "trainmap_table.html"
	outputPng := "trainmap_table_only.png"

	// Read the HTML content from the file
	htmlContent, err := ioutil.ReadFile(htmlFile)
	if err != nil {
		log.Fatalf("Failed to read HTML file: %v", err)
	}

	// Start a local HTTP server to serve the HTML content
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlContent)
	})

	// Run the server in the background
	server := &http.Server{Addr: ":8080"}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("HTTP server stopped: %v", err)
		}
	}()
	defer server.Close()

	// Create a new context for chromedp
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set up options to avoid timeout issues and improve stability
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Define the variable to hold the screenshot data
	var buf []byte

	// Run chromedp tasks to navigate to the local server and capture only the table element
	err = chromedp.Run(ctx, chromedp.Tasks{
		// Navigate to the locally served HTML page
		chromedp.Navigate("http://localhost:8080"),
		// Wait until the table is fully loaded
		chromedp.WaitVisible(`table`, chromedp.ByQuery),
		// Capture the table element as a screenshot
		chromedp.Screenshot(`table`, &buf, chromedp.ByQuery),
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

package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

func main() {
	// Path to the HTML file
	htmlFile := "weather.html"
	outputPng := "weather.png"

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

		// Apply grayscale filter using JavaScript
		chromedp.Evaluate(`document.body.style.filter = 'grayscale(100%)';`, nil),

		// Capture the #root element as a screenshot
		chromedp.Screenshot(`#root`, &buf, chromedp.ByQuery),
	})

	// Decode the PNG image from the buffer
	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Fatalf("Failed to decode screenshot: %v", err)
	}

	// Adjust brightness and contrast
	brightnessFactor := -20.0  // Reduce brightness to 70%
  contrastFactor := 40.0    // Increase contrast by 3x

	img = imaging.AdjustBrightness(img, brightnessFactor)
  img = imaging.AdjustContrast(img, contrastFactor)

	// Get the original width and height of the image
	origBounds := img.Bounds()
	origWidth := origBounds.Dx()

	// If the width is greater than 1200px, resize the image
	var resizedImg image.Image
	if origWidth > 1200 {
		// Resize to width 1200px while maintaining the aspect ratio
		resizedImg = resize.Resize(1200, 0, img, resize.Lanczos3) // 0 height keeps the aspect ratio
	} else {
		// No resize needed, use the original image
		resizedImg = img
	}

	// Write the resized image back to a PNG file
	outFile, err := os.Create(outputPng)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode the resized image to PNG and save
	err = png.Encode(outFile, resizedImg)
	if err != nil {
		log.Fatalf("Failed to encode PNG file: %v", err)
	}

	fmt.Printf("Table screenshot saved as %s\n", outputPng)
}

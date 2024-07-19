package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"gopkg.in/yaml.v2"
)

// Helper function to create a dummy image
func createDummyImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), 255, 255})
		}
	}
	return img
}

// Test main functionality
func TestMainFunction(t *testing.T) {
	// Create a dummy image
	img := createDummyImage()

	// Save the dummy image to a buffer
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("Failed to encode dummy image: %v\n", err)
	}

	// Create a temporary input file
	inputFile, err := os.CreateTemp("", "input-*.png")
	if err != nil {
		t.Fatalf("Failed to create temporary input file: %v\n", err)
	}
	defer os.Remove(inputFile.Name())

	// Write the dummy image to the input file
	_, err = buf.WriteTo(inputFile)
	if err != nil {
		t.Fatalf("Failed to write to temporary input file: %v\n", err)
	}

	// Create a temporary output file
	outputFile, err := os.CreateTemp("", "output-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary output file: %v\n", err)
	}
	defer os.Remove(outputFile.Name())

	// Run the main function
	os.Args = []string{"cmd", "-i", inputFile.Name(), "-o", outputFile.Name()}
	main()

	// Check if the output file is created
	if _, err := os.Stat(outputFile.Name()); os.IsNotExist(err) {
		t.Fatalf("Output file was not created\n")
	}

	// Read and decode the output file
	outData, err := os.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v\n", err)
	}

	var palette Palette
	err = yaml.Unmarshal(outData, &palette)
	if err != nil {
		t.Fatalf("Failed to decode output YAML file: %v\n", err)
	}

	// Check the number of colors in the palette
	const expectedColors = 16
	if len(palette.Colors) != expectedColors {
		t.Fatalf("Expected %d colors in the palette, got %d\n", expectedColors, len(palette.Colors))
	}
}

// Test clustering function
func TestClustering(t *testing.T) {
	img := createDummyImage()
	bounds := img.Bounds()

	var points clusters.Observations
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c, ok := colorful.MakeColor(img.At(x, y))
			if !ok {
				continue
			}
			points = append(points, clusters.Coordinates{c.R, c.G, c.B})
		}
	}

	const numClusters = 16
	kmeans := kmeans.New()
	clusterResult, err := kmeans.Partition(points, numClusters)
	if err != nil {
		t.Fatalf("Failed to perform clustering: %v\n", err)
	}

	if len(clusterResult) != numClusters {
		t.Fatalf("Expected %d clusters, got %d\n", numClusters, len(clusterResult))
	}
}

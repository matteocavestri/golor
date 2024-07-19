package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"gopkg.in/yaml.v2"
)

type Palette struct {
	Colors []string `yaml:"colors"`
}

func main() {
	if len(os.Args) != 5 {
		log.Fatalf("Usage: %s -i <input.jpg/png> -o <output.yaml>\n", os.Args[0])
	}

	inputFile := os.Args[2]
	outputFile := os.Args[4]

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v\n", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v\n", err)
	}

	var points clusters.Observations
	bounds := img.Bounds()
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
		log.Fatalf("Failed to perform clustering: %v\n", err)
	}

	var colors []string
	for _, cluster := range clusterResult {
		center := cluster.Center
		c := colorful.Color{R: center[0], G: center[1], B: center[2]}
		colors = append(colors, c.Hex())
	}

	palette := Palette{Colors: colors}

	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v\n", err)
	}
	defer outFile.Close()

	encoder := yaml.NewEncoder(outFile)
	defer encoder.Close()

	err = encoder.Encode(palette)
	if err != nil {
		log.Fatalf("Failed to write to output file: %v\n", err)
	}

	fmt.Printf("Palette written to %s\n", outputFile)
}

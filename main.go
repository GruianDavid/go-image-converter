package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sunshineplan/imgconv"
)

// Usage: go run main.go -threads=6 -input=./original -output=./

var lgHeight int
var mdHeight int
var smHeight int

func main() {
	threads := flag.Int("threads", 6, "Number of concurrent photo processors")
	inputDir := flag.String("input", "./original", "Directory containing photos to process (prefix required like ./)")
	outputDir := flag.String("output", "./", "Directory to save converted webp images (prefix required like ./)")
	lg := flag.Int("lgHeight", 1080, "Generate large size images")
	md := flag.Int("mdHeight", 720, "Generate medium size images")
	sm := flag.Int("smHeight", 480, "Generate small size images")
	flag.Parse()
	lgHeight = *lg
	mdHeight = *md
	smHeight = *sm
	photoChan := getPhotosFromDirectory(*inputDir)

	// Create a semaphore channel
	sem := make(chan struct{}, *threads)
	var wg sync.WaitGroup

	for photoPath := range photoChan {
		sem <- struct{}{} // acquire semaphore
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			fmt.Printf("Processing %s\n", path)
			err := processPhotoToWebP(path, *outputDir)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", path, err)
			}
			<-sem // release semaphore
		}(photoPath)
	}
	wg.Wait()
}

// processPhotoToWebP receives a photo path, converts it to .avif with same dimensions, and saves it in 'avif' directory.
func processPhotoToWebP(photoPath string, output string) error {
	// Prepare output directory
	outDir := output + "/webp"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create webp directory: %w", err)
	}

	// Open the image
	image, err := imgconv.Open(photoPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	// file sizes lg, md, sm
	var imageSizes = []string{"lg", "md", "sm"}
	for _, size := range imageSizes {
		err := convertAndSaveWebP(photoPath, outDir, size, image)
		if err != nil {
			return fmt.Errorf("failed to convert %s to %s: %w", photoPath, size, err)
		}
	}

	return nil
}

func convertAndSaveWebP(photoPath, outDir, size string, image image.Image) error {
	// Prepare output filename
	base := filepath.Base(photoPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	outPath := filepath.Join(outDir, name+"_"+size+".webp")
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Resize image based on size
	switch size {
	case "lg":
		image = imgconv.Resize(image, &imgconv.ResizeOption{Height: lgHeight})
	case "md":
		image = imgconv.Resize(image, &imgconv.ResizeOption{Height: mdHeight})
	case "sm":
		image = imgconv.Resize(image, &imgconv.ResizeOption{Height: smHeight})
	}

	imgconv.Write(outFile, image, &imgconv.FormatOption{
		Format: imgconv.WEBP,
	})

	return nil
}

// getPhotosFromDirectory scans the given directory for photo files and sends their paths to a channel.
func getPhotosFromDirectory(dir string) <-chan string {
	photoChan := make(chan string)
	go func() {
		defer close(photoChan)
		filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // skip files with errors
			}
			if !d.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
					photoChan <- path
				}
			}
			return nil
		})
	}()
	return photoChan
}

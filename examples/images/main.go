package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/free-llms-foundation/retrieval-go"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ret, err := retrieval.New()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	images, err := ret.SearchImagesWithQuery(ctx, "Donald Trump", "y")
	if err != nil {
		log.Fatalf("failed to get images: %v", err)
	}

	for _, img := range images {
		fmt.Printf("URL: %s\n", img.URL)
		fmt.Printf("Thumbnail: %s\n", img.Thumbnail)
		fmt.Printf("PageURL: %s\n", img.PageURL)
		fmt.Printf("Desc: %s\n", img.Desc)
		fmt.Print("-------------------------------------------------------\n\n")
	}
}

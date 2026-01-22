package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/free-llms-foundation/retrieval-go"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	ret, err := retrieval.New()
	if err != nil {
		log.Fatal(err)
	}

	url := "https://go.dev"
	doc, err := ret.ParseContentFromLink(ctx, url)
	if err != nil {
		if errors.Is(err, retrieval.ErrRobotsDenied) {
			log.Println("rodobt.txt denied")
			return
		}

		log.Fatal(err)
	}

	fmt.Printf("Title:%s\n\n", doc.Title)
	fmt.Println("Content:", doc.Content[:min(2000, len(doc.Content))])
	fmt.Println("Byline:", doc.Byline)
	fmt.Println("SiteName:", doc.SiteName)
	fmt.Println("Image:", doc.Image)
	fmt.Println("Favicon:", doc.Favicon)
	fmt.Println("Language:", doc.Language)

}

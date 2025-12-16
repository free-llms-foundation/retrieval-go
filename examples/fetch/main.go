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

	ret := retrieval.New()

	url := "https://go.dev/"
	doc, err := ret.ParseContentFromLink(ctx, url)
	if err != nil {
		if errors.Is(err, retrieval.ErrRobotsDenied) {
			log.Println("rodobt.txt denied")
			return
		}

		log.Fatal(err)
	}

	fmt.Println("Title:", doc.Title)
	fmt.Println("URL:", doc.URL)
	fmt.Println("Content:", doc.Content[:min(1000, len(doc.Content))])
}

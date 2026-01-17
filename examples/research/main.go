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

	pages, err := ret.SearchWithQuery(ctx, "golang")
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range pages {
		fmt.Printf("Link=%s, Title: %s\n\n", p.Link, p.Title)
		doc, err := ret.ParseContentFromLink(ctx, p.Link)
		if err != nil {
			if errors.Is(err, retrieval.ErrRobotsDenied) {
				fmt.Print("Skip. Robots.txt denied\n\n")
				continue
			}
		}
		fmt.Println("Content:", doc.Content[:min(1000, len(doc.Content))])
		fmt.Print("-------------------------------------------------------\n\n")
		time.Sleep(time.Second) // sleep to avoid getting captcha
	}
}

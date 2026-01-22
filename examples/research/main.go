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

	pages, err := ret.SearchWithQuery(ctx, "dollar exchange rate")
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range pages {

		doc, err := ret.ParseContentFromLink(ctx, p.Link)
		if err != nil {
			if errors.Is(err, retrieval.ErrRobotsDenied) {
				fmt.Print("Skip. Robots.txt denied\n\n")
				continue
			}
			if errors.Is(err, retrieval.ErrUnexpectedStatusCode) {
				fmt.Print("Skip. Content is blocked\n\n")
				continue
			}
		}

		fmt.Printf("Title:%s\n\n", doc.Title)
		fmt.Println("Content:", doc.Content[:min(1000, len(doc.Content))])
		fmt.Println("Byline:", doc.Byline)
		fmt.Println("SiteName:", doc.SiteName)
		fmt.Println("Image:", doc.Image)
		fmt.Println("Favicon:", doc.Favicon)
		fmt.Println("Language:", doc.Language)
		fmt.Print("-------------------------------------------------------\n\n")
		time.Sleep(time.Second) // sleep to avoid getting captcha

	}

}

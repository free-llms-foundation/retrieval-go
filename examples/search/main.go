package main

import (
	"context"
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

	pages, err := ret.SearchWithQuery(ctx, "Golang")
	if err != nil {
		log.Fatal(err)
	}

	for i, p := range pages {
		fmt.Printf("%d. %s\n", i+1, p.Title)
		fmt.Printf("%s\n", p.Link)
		fmt.Printf("%s\n", p.Favicon)
		if p.Snippet != "" {
			fmt.Printf("%s\n", p.Snippet)
		}
		fmt.Print("-------------------------------------------------------\n\n")
	}
}

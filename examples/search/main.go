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

	ret := retrieval.New()

	pages, err := ret.SearchWithQuery(ctx, "golang")
	if err != nil {
		log.Fatal(err)
	}

	for i, p := range pages {
		fmt.Printf("%d. %s\n", i+1, p.Title)
		fmt.Printf("%s\n", p.Link)
		if p.Snippet != "" {
			fmt.Printf("%s\n", p.Snippet)
		}
		fmt.Print("-------------------------------------------------------\n\n")
	}
}

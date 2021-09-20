package main

import (
	_ "github.com/sfomuseum/go-sfomuseum-airports/sfomuseum"
)

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-airports"

	"log"
)

func main() {

	lookup_uri := flag.String("lookup-uri", "sfomuseum://", "Valid sources are: sfomuseum://")
	flag.Parse()

	ctx := context.Background()

	lookup, err := airports.NewLookup(ctx, *lookup_uri)

	if err != nil {
		log.Fatal(err)
	}

	for _, code := range flag.Args() {

		results, err := lookup.Find(code)

		if err != nil {
			fmt.Printf("%s *** %s\n", code, err)
			continue
		}

		for _, a := range results {
			fmt.Println(a)
		}
	}

}

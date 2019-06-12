package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-airports"
	"github.com/sfomuseum/go-sfomuseum-airports/sfomuseum"
	"log"
)

func main() {

	source := flag.String("source", "sfomuseum", "Valid sources are: sfomuseum.")
	flag.Parse()

	var lookup airports.Lookup
	var err error

	switch *source {
	case "sfomuseum":
		lookup, err = sfomuseum.NewLookup()
	default:
		err = errors.New("Unknown source")
	}

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

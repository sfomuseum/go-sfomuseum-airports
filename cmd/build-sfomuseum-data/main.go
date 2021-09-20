package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/sfomuseum/go-sfomuseum-airports/sfomuseum"
	"io"
	"log"
	"os"
)

func main() {

	iterator_uri := flag.String("iterator-uri", "repo://?include=properties.sfomuseum:placetype=airport", "...")
	iterator_source := flag.String("iterator-source", "/usr/local/data/sfomuseum-data-whosonfirst", "...")

	target := flag.String("target", "data/sfomuseum.json", "The path to write SFO Museum airline data.")
	stdout := flag.Bool("stdout", false, "Emit SFO Museum aircraft data to SDOUT.")

	flag.Parse()

	ctx := context.Background()

	writers := make([]io.Writer, 0)

	fh, err := os.OpenFile(*target, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Failed to open '%s', %v", *target, err)
	}

	writers = append(writers, fh)

	if *stdout {
		writers = append(writers, os.Stdout)
	}

	wr := io.MultiWriter(writers...)

	lookup, err := sfomuseum.CompileAirportsData(ctx, *iterator_uri, *iterator_source)

	if err != nil {
		log.Fatalf("Failed to compile data, %v", err)
	}

	enc := json.NewEncoder(wr)
	err = enc.Encode(lookup)

	if err != nil {
		log.Fatalf("Failed to marshal results, %w", err)
	}

}

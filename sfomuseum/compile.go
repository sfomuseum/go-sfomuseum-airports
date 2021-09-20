package sfomuseum

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-geojson/feature"
	sfomuseum_props "github.com/sfomuseum/go-sfomuseum-geojson/properties/sfomuseum"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"log"
	"sync"
)

func CompileAirportsData(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]Airport, error) {

	lookup := make([]Airport, 0)
	mu := new(sync.RWMutex)

	iter_cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		path, err := emitter.PathForContext(ctx)

		if err != nil {
			return fmt.Errorf("Failed to derive path from context, %w", err)
		}

		_, uri_args, err := uri.ParseURI(path)

		if err != nil {
			return fmt.Errorf("Failed to parse %s, %w", path, err)
		}

		if uri_args.IsAlternate {
			return nil
		}

		f, err := feature.LoadFeatureFromReader(fh)

		if err != nil {
			return fmt.Errorf("Failed load feature from %s, %w", path, err)
		}

		// TO DO : https://github.com/sfomuseum/go-sfomuseum-airports-tools/issues/1

		pt := sfomuseum_props.Placetype(f)

		if pt != "airport" {
			log.Println("NOT AN AIRPORT", whosonfirst.Id(f), whosonfirst.Name(f), pt)
			return nil
		}

		wof_id := whosonfirst.Id(f)
		name := whosonfirst.Name(f)

		sfom_id := utils.Int64Property(f.Bytes(), []string{"properties.sfomuseum:airport_id"}, -1)

		concordances, err := whosonfirst.Concordances(f)

		if err != nil {
			return err
		}

		a := Airport{
			WOFID:       wof_id,
			SFOMuseumID: int(sfom_id),
			Name:        name,
		}

		iata_code, ok := concordances["iata:code"]

		if ok && iata_code != "" {
			a.IATACode = iata_code
		}

		icao_code, ok := concordances["icao:code"]

		if ok && icao_code != "" {
			a.ICAOCode = icao_code
		}

		mu.Lock()
		lookup = append(lookup, a)
		mu.Unlock()

		return nil
	}

	iter, err := iterator.NewIterator(ctx, iterator_uri, iter_cb)

	if err != nil {
		return nil, fmt.Errorf("Failed to create iterator, %w", err)
	}

	err = iter.IterateURIs(ctx, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to iterate sources, %w", err)
	}

	return lookup, nil
}

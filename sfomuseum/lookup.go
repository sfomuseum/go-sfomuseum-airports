package sfomuseum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-airports"
	"github.com/sfomuseum/go-sfomuseum-airports/data"
	"io"
	"strconv"
	"strings"
	"sync"
)

var lookup_table *sync.Map
var lookup_init sync.Once
var lookup_init_err error

type SFOMuseumLookupFunc func(context.Context)

type SFOMuseumLookup struct {
	airports.Lookup
}

func NewLookup(ctx context.Context, uri string) (airports.Lookup, error) {

	fs := data.FS
	fh, err := fs.Open("sfomuseum.json")

	if err != nil {
		return nil, fmt.Errorf("Failed to load data, %v", err)
	}

	lookup_func := NewLookupFuncWithReader(ctx, fh)
	return NewLookupWithLookupFunc(ctx, lookup_func)
}

// NewLookup will return an `SFOMuseumLookupFunc` function instance that, when invoked, will populate an `airports.Lookup` instance with data stored in `r`.
// `r` will be closed when the `SFOMuseumLookupFunc` function instance is invoked.
// It is assumed that the data in `r` will be formatted in the same way as the procompiled (embedded) data stored in `data/sfomuseum.json`.
func NewLookupFuncWithReader(ctx context.Context, r io.ReadCloser) SFOMuseumLookupFunc {

	lookup_func := func(ctx context.Context) {

		defer r.Close()

		var airport []*Airport

		dec := json.NewDecoder(r)
		err := dec.Decode(&airport)

		if err != nil {
			lookup_init_err = err
			return
		}

		table := new(sync.Map)

		for idx, craft := range airport {

			pointer := fmt.Sprintf("pointer:%d", idx)
			table.Store(pointer, craft)

			str_wofid := strconv.FormatInt(craft.WOFID, 10)

			possible_codes := []string{
				craft.IATACode,
				craft.ICAOCode,
				str_wofid,
			}

			for _, code := range possible_codes {

				if code == "" {
					continue
				}

				pointers := make([]string, 0)
				has_pointer := false

				others, ok := table.Load(code)

				if ok {

					pointers = others.([]string)
				}

				for _, dupe := range pointers {

					if dupe == pointer {
						has_pointer = true
						break
					}
				}

				if has_pointer {
					continue
				}

				pointers = append(pointers, pointer)
				table.Store(code, pointers)
			}

			idx += 1
		}

		lookup_table = table
	}

	return lookup_func
}

// NewLookupWithLookupFunc will return an `airports.Lookup` instance derived by data compiled using `lookup_func`.
func NewLookupWithLookupFunc(ctx context.Context, lookup_func SFOMuseumLookupFunc) (airports.Lookup, error) {

	fn := func() {
		lookup_func(ctx)
	}

	lookup_init.Do(fn)

	if lookup_init_err != nil {
		return nil, lookup_init_err
	}

	l := SFOMuseumLookup{}
	return &l, nil
}

func (l *SFOMuseumLookup) Find(code string) ([]interface{}, error) {

	pointers, ok := lookup_table.Load(code)

	if !ok {
		return nil, errors.New("Not found")
	}

	airport := make([]interface{}, 0)

	for _, p := range pointers.([]string) {

		if !strings.HasPrefix(p, "pointer:") {
			return nil, errors.New("Invalid pointer")
		}

		row, ok := lookup_table.Load(p)

		if !ok {
			return nil, errors.New("Invalid pointer")
		}

		airport = append(airport, row.(*Airport))
	}

	return airport, nil
}

package sfomuseum

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-airports"
	"github.com/sfomuseum/go-sfomuseum-airports/data"
	"io"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var lookup_table *sync.Map
var lookup_idx int64

var lookup_init sync.Once
var lookup_init_err error

type SFOMuseumLookupFunc func(context.Context)

type SFOMuseumLookup struct {
	airports.Lookup
}

func init() {
	ctx := context.Background()
	airports.RegisterLookup(ctx, "sfomuseum", NewLookup)

	lookup_idx = int64(0)
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

	defer r.Close()

	var airports_list []*Airport

	dec := json.NewDecoder(r)
	err := dec.Decode(&airports_list)

	if err != nil {

		lookup_func := func(ctx context.Context) {
			lookup_init_err = err
		}

		return lookup_func
	}

	return NewLookupFuncWithAirports(ctx, airports_list)
}

// NewLookup will return an `SFOMuseumLookupFunc` function instance that, when invoked, will populate an `airports.Lookup` instance with data stored in `airports_list`.
func NewLookupFuncWithAirports(ctx context.Context, airports_list []*Airport) SFOMuseumLookupFunc {

	lookup_func := func(ctx context.Context) {

		table := new(sync.Map)

		for _, data := range airports_list {

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			appendData(ctx, table, data)
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

func NewLookupFromIterator(ctx context.Context, iterator_uri string, iterator_sources ...string) (airports.Lookup, error) {

	airports_list, err := CompileAirportsData(ctx, iterator_uri, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile airport data, %w", err)
	}

	lookup_func := NewLookupFuncWithAirports(ctx, airports_list)
	return NewLookupWithLookupFunc(ctx, lookup_func)
}

func (l *SFOMuseumLookup) Find(ctx context.Context, code string) ([]interface{}, error) {

	pointers, ok := lookup_table.Load(code)

	if !ok {
		return nil, fmt.Errorf("Code '%s' not found", code)
	}

	airport := make([]interface{}, 0)

	for _, p := range pointers.([]string) {

		if !strings.HasPrefix(p, "pointer:") {
			return nil, fmt.Errorf("Invalid pointer, '%s'", p)
		}

		row, ok := lookup_table.Load(p)

		if !ok {
			return nil, fmt.Errorf("Invalid pointer, '%s'", p)
		}

		airport = append(airport, row.(*Airport))
	}

	return airport, nil
}

func (l *SFOMuseumLookup) Append(ctx context.Context, data interface{}) error {
	return appendData(ctx, lookup_table, data.(*Airport))
}

func appendData(ctx context.Context, table *sync.Map, data *Airport) error {

	idx := atomic.AddInt64(&lookup_idx, 1)

	pointer := fmt.Sprintf("pointer:%d", idx)
	table.Store(pointer, data)

	str_wofid := strconv.FormatInt(data.WOFID, 10)
	str_sfomid := strconv.Itoa(data.SFOMuseumID)

	possible_codes := []string{
		data.IATACode,
		data.ICAOCode,
		str_wofid,
		str_sfomid,
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

	return nil
}

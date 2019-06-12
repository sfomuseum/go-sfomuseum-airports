package sfomuseum

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sfomuseum/go-sfomuseum-airports"
	"strconv"
	"strings"
	"sync"
)

type Airport struct {
	WOFID       int64  `json:"wof:id"`
	Name        string `json:"wof:name"`
	SFOMuseumID int    `json:"sfomuseum:airport_id"`
	IATACode    string `json:"iata:code"`
	ICAOCode    string `json:"icao:code"`
}

func (a *Airport) String() string {
	return fmt.Sprintf("%s %s \"%s\" %d", a.IATACode, a.ICAOCode, a.Name, a.WOFID)
}

var lookup_table *sync.Map
var lookup_init sync.Once

type SFOMuseumLookup struct {
	airports.Lookup
}

func NewLookup() (airports.Lookup, error) {

	var lookup_err error

	lookup_func := func() {

		var airport []*Airport

		err := json.Unmarshal([]byte(AirportData), &airport)

		if err != nil {
			lookup_err = err
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

	lookup_init.Do(lookup_func)

	if lookup_err != nil {
		return nil, lookup_err
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

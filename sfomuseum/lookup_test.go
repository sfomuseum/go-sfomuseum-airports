package sfomuseum

import (
	"context"
	"github.com/sfomuseum/go-sfomuseum-airports"
	"testing"
)

func TestSFOMuseumLookup(t *testing.T) {

	wofid_tests := map[string]int64{
		"YUL":  102554351,
		"EGLL": 102556703,
		"162":  1360695653,
		"260":  102525431,
	}

	ctx := context.Background()

	lu, err := airports.NewLookup(ctx, "sfomuseum://")

	if err != nil {
		t.Fatalf("Failed to create lookup, %v", err)
	}

	for code, wofid := range wofid_tests {

		results, err := lu.Find(ctx, code)

		if err != nil {
			t.Fatalf("Unable to find '%s', %v", code, err)
		}

		if len(results) != 1 {
			t.Fatalf("Invalid results for '%s'", code)
		}

		a := results[0].(*Airport)

		if a.WOFID != wofid {
			t.Fatalf("Invalid match for '%s', expected %d but got %d", code, wofid, a.WOFID)
		}
	}
}

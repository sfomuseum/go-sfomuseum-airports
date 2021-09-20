package sfomuseum

import (
	"fmt"
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

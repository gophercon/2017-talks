package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/losinggeneration/geojson"
)

type geojsoner interface {
	GeoJSON() interface{}
}

// A cheater method to hide the fact that there's no data store to get this from.
func GetGeoJSON(g GopherCon) interface{} {
	return geojson.GeoJSON{
		Feature: &geojson.Feature{
			Geometry: &geojson.Geometry{
				Point: &geojson.Point{
					Coordinates: geojson.Position{
						39.742329, -104.9965061,
					},
				},
			},
			Properties: geojson.Properties{
				"id":   g.Id,
				"name": g.Name,
			},
		},
	}
}

// basic wrapper around json.Marshal that will call the model's GeoJSON method or
// return an error if it's not of that type
func geojsonMarshal(v interface{}) ([]byte, error) {
	if g, ok := v.(geojsoner); !ok {
		return nil, errors.New(fmt.Sprintf("cannot convert to geojson: %T", v))
	} else {
		return json.Marshal(g.GeoJSON())
	}
}

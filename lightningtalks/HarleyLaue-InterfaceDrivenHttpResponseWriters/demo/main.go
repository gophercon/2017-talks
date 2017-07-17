package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// common response content types
var (
	jsonAccept    = Accept{"application/json", json.Marshal}
	xmlAccept     = Accept{"application/xml", xml.Marshal}
	geojsonAccept = Accept{"application/vnd.geo+json", geojsonMarshal}
)

type GopherCon struct {
	Id   int64
	Name string
}

func (g GopherCon) GeoJSON() interface{} {
	return GetGeoJSON(g)
}

// Plaintext example uses a string
func TextHandler(w http.ResponseWriter, r *http.Request) {
	if err := Write(w,
		"this is plain text",
		ContentType("text/plain"),
	); err != nil {
		log.Printf("Error: %+v", err)
	}
}

// Example wrapper function around Write
func WriteJSON(w http.ResponseWriter, v interface{}) error {
	return errors.WithStack(Write(w, v, jsonAccept))
}

// JSON example that uses WriteJSON
func JSONHandler(w http.ResponseWriter, r *http.Request) {
	if err := WriteJSON(w, GopherCon{123, "GopherCon 2017"}); err != nil {
		log.Printf("Error: %+v", err)
	}
}

// Mixed handler that can return JSON, XML, & GeoJSON
func MixedHandler(w http.ResponseWriter, r *http.Request) {
	if err := Write(w,
		GopherCon{456, "GopherCon 2017"},
		Acceptable(r, jsonAccept, geojsonAccept, xmlAccept),
	); err != nil {
		log.Printf("Error: %+v", err)
	}
}

func main() {
	http.HandleFunc("/text", TextHandler)
	http.HandleFunc("/json", JSONHandler)
	http.HandleFunc("/mixed", MixedHandler)

	fmt.Println("running on 127.0.0.1:8888")
	http.ListenAndServe("127.0.0.1:8888", nil)
}

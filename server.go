package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type geoFeatures struct {
	Type     string       `json:"type"`
	Features []geoFeature `json:"features"`
}

type geoFeature struct {
	Type     string   `json:"type"`
	Geometry geoPoint `json:"geometry"`

	Properties map[string]interface{} `json:"properties"`
}

type geoPoint struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func townListGeoJSON(towns []*Town) *geoFeatures {
	res := geoFeatures{Type: "FeatureCollection", Features: make([]geoFeature, len(towns))}
	for i, t := range towns {
		res.Features[i] = geoFeature{
			Type:       "Feature",
			Geometry:   geoPoint{Type: "Point", Coordinates: []float64{t.Longitude, t.Latitude}},
			Properties: map[string]interface{}{"name": t.Name, "state": t.State},
		}
	}
	return &res
}

func ListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(townListGeoJSON(searchTowns(req.FormValue("q"))))
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "index.html")
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := indexTowns(f); err != nil {
		log.Fatal(err)
	}
	f.Close()
	log.Println("loaded")

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/list", ListHandler)
	http.ListenAndServe(":8080", nil)
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"googlemaps.github.io/maps"
)

type mapSvc struct {
	client *maps.Client
}

// radius in meters (max 50000)
func (svc *mapSvc) availableLocations(currentLocation *maps.LatLng, radius uint) availableLocationsOutput {
	q := "Fast food restaurants with drive-throughs in Honolulu"
	r := &maps.TextSearchRequest{
		Query:    q,
		Location: currentLocation,
		Radius:   radius,
	}
	resp, err := svc.client.TextSearch(context.Background(), r)
	if err != nil {
		fmt.Println(err)
		return availableLocationsOutput{
			Locations: nil,
			Err:       err,
		}
	}
	// TODO: filter by drive-through hours if available

	return availableLocationsOutput{
		Locations: toLoc(resp),
		Err:       nil,
	}
}

func (svc *mapSvc) attach(mux *http.ServeMux) {
	fmt.Println("POST /api/availableLocations")
	mux.HandleFunc("POST /api/availableLocations", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit: /api/availableLocations")
		var bodyObj struct {
			CurrentLocation struct {
				Lat float64
				Lng float64
			}
			Radius uint
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error reading request body.", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &bodyObj)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Body could not be unmarshalled.", http.StatusInternalServerError)
			return
		}
		// check fields
		if bodyObj.CurrentLocation.Lat == 0.0 || bodyObj.CurrentLocation.Lng == 0.0 {
			http.Error(w, "Missing in body 'currentLocation'", http.StatusBadRequest)
			return
		}
		if bodyObj.Radius == 0 {
			http.Error(w, "Missing in body 'radius' (in meters)", http.StatusBadRequest)
			return
		}
		// call places api
		out := svc.availableLocations(&maps.LatLng{
			Lat: bodyObj.CurrentLocation.Lat,
			Lng: bodyObj.CurrentLocation.Lng,
		}, bodyObj.Radius)
		jsonOut, err := json.Marshal(out)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Data could not be marshalled into json.", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(200)
		w.Write(jsonOut)
	})
}

func newMapSvc() *mapSvc {
	fmt.Println("New Map Service")
	apiKey := os.Getenv("GCP_API_KEY")
	if len(apiKey) == 0 {
		panic("Missing env var: GCP_API_KEY")
	}
	fmt.Println("Found env var GCP_API_KEY: " + apiKey)

	// create maps client
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	return &mapSvc{
		client: c,
	}
}

func toLoc(obj maps.PlacesSearchResponse) []Location {
	ret := make([]Location, 0, 10)
	for i := 0; i < len(obj.Results); i++ {
		place := obj.Results[i]
		ret = append(ret, Location{
			Name:     place.Name,
			Address:  place.FormattedAddress,
			Distance: 0.0,
		})
	}
	return ret
}

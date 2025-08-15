package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type mapSvc struct {
	apiKey string
}

func (svc *mapSvc) getAvailableLocations(currentLocation string) getAvailableLocationsOutput {
	fmt.Println("API: getAvailableLocations")
	fmt.Printf("currentLocation: %s\n", currentLocation)
	// TODO: send post req with
	// -d '{ "textQuery": "THE QUERY" }'
	// -H 'Context-Type: application/json'
	// -H 'X-Goog-Api-Key: '
	// resp, err := http.Post("https://places.googleapis.com/vi/places:searchText")
	// if err != nil {
	// 	return getAvailableLocationsOutput{
	// 		Locations: []Location{},
	// 		Err:       err,
	// 	}
	// }
	// defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return getAvailableLocationsOutput{
	// 		Locations: []Location{},
	// 		Err:       err,
	// 	}
	// }

	// fmt.Println(string(body))

	return getAvailableLocationsOutput{
		Locations: []Location{},
		Err:       nil,
	}
}

func (svc *mapSvc) attach(mux *http.ServeMux) {
	mux.Handle("/api/getAvailableLocations", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cl := r.FormValue("currentLocation")
		if len(cl) == 0 {
			// missing arg
			http.Error(w, "Missing query arg 'currentLocation'", http.StatusBadRequest)
			return
		}
		out := svc.getAvailableLocations(cl)
		jsonOut, err := json.Marshal(out)
		if err != nil {
			http.Error(w, "Data could not be marshalled into json.", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(200)
		w.Write(jsonOut)
	}))
}

func newMapSvc() *mapSvc {
	fmt.Println("New Map Service")
	apiKey := os.Getenv("GCP_API_KEY")
	if len(apiKey) == 0 {
		panic("Missing env var: GCP_API_KEY")
	}
	fmt.Println("Found env var GCP_API_KEY: " + apiKey)
	return &mapSvc{
		apiKey: apiKey,
	}
}

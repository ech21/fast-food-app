package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type mapSvc struct {
	apiKey string
}

func (svc *mapSvc) getAvailableLocations(currentLocation string) getAvailableLocationsOutput {
	fmt.Println("getAvailableLocations")
	fmt.Printf("currentLocation: %s\n", currentLocation)
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
	// read env
	fmt.Println("Getting api key env var")
	apiKey := ""
	return &mapSvc{
		apiKey: apiKey,
	}
}

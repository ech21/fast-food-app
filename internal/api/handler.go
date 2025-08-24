package api

import (
	"github.com/ech21/fast-food-app/internal/types"
	"image"
	"net/http"

	"googlemaps.github.io/maps"
)

type attacher interface {
	attach(mux *http.ServeMux)
}

// Map service --------------------------------------------------------------------

type availableLocationsOutput struct {
	Locations []types.Location `json:"locations"`
	Err       error            `json:"err"`
}

type mapService interface {
	attacher
	// availableLocations takes the current location and radius and returns a list of
	// locations that satisfy the fast food challenge rules and any error.
	availableLocations(currentLocation *maps.LatLng, radius uint) availableLocationsOutput
}

// Receipt service -----------------------------------------------------------------

type receiptService interface {
	attacher
	// ParseReceipt takes an image and returns relevant info and any error.
	ParseReceipt(image image.Image) struct {
		Receipt types.ReceiptData `json:"receipt"`
		Err     error             `json:"err"`
	}
}

// Nutrition service ---------------------------------------------------------------

type nutritionAutocompleteOutput struct {
	Results []types.Item `json:"results"`
	Err     error        `json:"err"`
}

type nutritionService interface {
	attacher
	// nutritionAutocomplete takes a search query and gives a list of potential
	// matches and any error.
	nutritionAutocomplete(q string) nutritionAutocompleteOutput
}

// Joined services object ----------------------------------------------------------

type svc struct {
	Map       mapService
	Receipt   receiptService
	Nutrition nutritionService
}

func Svc() svc {
	return svc{
		Map:       newMapSvc(),
		Nutrition: newNutritionSvc(),
	}
}

func Attach(mux *http.ServeMux) {
	svc := Svc()
	svc.Map.attach(mux)
	svc.Nutrition.attach(mux)
}

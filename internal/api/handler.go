package api

import (
	"image"
	"net/http"

	"googlemaps.github.io/maps"
)

type attacher interface {
	attach(mux *http.ServeMux)
}

// Map service --------------------------------------------------------------------

type availableLocationsOutput struct {
	Locations []Location `json:"locations"`
	Err       error      `json:"err"`
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
		Receipt ReceiptData `json:"receipt"`
		Err     error       `json:"err"`
	}
}

// Nutrition service ---------------------------------------------------------------

type nutritionixSearchInstantResult struct {
	Photo struct {
		Thumb string `json:"thumb"`
	} `json:"photo"`
	Nf_calories          uint   `json:"nf_calories"`
	Serving_unit         string `json:"serving_unit"`
	Serving_qty          uint   `json:"serving_qty"`
	Brand_type           uint   `json:"brand_type"`
	Region               uint   `json:"region"`
	Nix_brand_id         string `json:"nix_brand_id"`
	Brand_name           string `json:"brand_name"`
	Food_name            string `json:"food_name"`
	Brand_name_item_name string `json:"brand_name_item_name"`
	Nix_item_id          string `json:"nix_item_id"`
	Locale               string `json:"locale"`
}

type nutritionAutocompleteOutput struct {
	Results []nutritionixSearchInstantResult `json:"results"`
	Err     error                            `json:"err"`
}

type nutritionInfoOutput struct {
	Info NutritionInfo `json:"info"`
	Err  error         `json:"err"`
}

type nutritionService interface {
	attacher
	// nutritionAutocomplete takes a search query and gives a list of potential
	// matches and any error.
	nutritionAutocomplete(q string) nutritionAutocompleteOutput
	// nutritionInfo gets the nutrition info of a given food item and any error.
	nutritionInfo(item Item) nutritionInfoOutput
}

// Lobby service -------------------------------------------------------------------

type lobbyService interface {
	attacher
	// New creates a new lobby.
	New() *Lobby
	// Join takes a lobby id and a Player object and tries to add the player to
	// the lobby, or have the player rejoin if the player was already in the lobby.
	// It returns the lobby joined and any error while trying to join.
	Join(id string, player *Player) (lobby *Lobby, err error)
	// Close shuts down a lobby and removes all players. It may return an error.
	Close(lobby *Lobby) error
}

// Joined services object ----------------------------------------------------------

type svc struct {
	Map       mapService
	Receipt   receiptService
	Nutrition nutritionService
	Lobby     lobbyService
}

func Svc() svc {
	return svc{
		Map:       newMapSvc(),
		Nutrition: newNutritionSvc(),
	}
}

func AttachHandlers(mux *http.ServeMux) {
	svc := Svc()
	svc.Map.attach(mux)
	svc.Nutrition.attach(mux)
}

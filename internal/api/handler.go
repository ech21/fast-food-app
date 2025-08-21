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

type nutritionAutocompleteOutput struct {
	Results []Item `json:"results"`
	Err     error  `json:"err"`
}

type nutritionService interface {
	attacher
	// nutritionAutocomplete takes a search query and gives a list of potential
	// matches and any error.
	nutritionAutocomplete(q string) nutritionAutocompleteOutput
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

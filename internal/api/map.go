package api

// A Location has details about a fast food restaurant
type Location struct {
	Name    string
	Address string
	// in miles
	Distance float32
}

type MapService interface {
	// GetAvailableLocations takes the current location and returns a list of
	// locations that satisfy the fast food challenge rules and any error.
	GetAvailableLocations(currentLocation string) (locations []Location, err error)
}

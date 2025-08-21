package api

// A Location has details about a fast food restaurant
type Location struct {
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Distance float64 `json:"distance"` // In miles
}

// Details of a food item
type Item struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	LocationName string  `json:"locationName"`
	Price        float64 `json:"price"`
	Calories     uint    `json:"calories"`
}

type ReceiptData struct {
	LocationName string  `json:"locationName"`
	Items        []Item  `json:"items"`
	TotalPrice   float64 `json:"totalPrice"`
}

type Settings struct {
	// in miles, the maximum distance to search for restaurants
	MaxRadius  float64
	MaxPlayers uint
}

type Player struct {
	Name string `json:"name"`
}

type Lobby struct {
	Id      string
	Players []Player
}

package api

// A Location has details about a fast food restaurant
type Location struct {
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Distance float32 `json:"distance"` // In miles
}

// Details of a food item
type Item struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type ReceiptData struct {
	LocationName string  `json:"locationName"`
	Items        []Item  `json:"items"`
	TotalPrice   float32 `json:"totalPrice"`
}

type NutritionInfo struct {
	Name     string
	Calories int
}

type Settings struct {
	// in miles, the maximum distance to search for restaurants
	MaxRadius  float32
	MaxPlayers int
}

type Player struct {
	Id   string
	Name string
}

type Lobby struct {
	Id      string
	Players []Player
}

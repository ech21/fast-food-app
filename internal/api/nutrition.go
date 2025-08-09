package api

type NutritionInfo struct {
	Name     string
	Calories int
}

type NutritionService interface {
	// GetNutritionInfo gets the nutrition info of a given food item.
	GetNutritionInfo(item Item) NutritionInfo
}

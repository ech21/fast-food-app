package api

import (
	"image"
)

// Details of a food item
type Item struct {
	Name  string
	Price string
}

type ReceiptData struct {
	LocationName string
	Items        []Item
	TotalPrice   float32
}

type ReceiptService interface {
	// ParseReceipt takes an image and returns relevant info.
	ParseReceipt(image image.Image) ReceiptData
}

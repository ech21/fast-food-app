package main

import (
	"fmt"
	"net/http"

	"github.com/ech21/fast-food-app/internal/api"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./ui/dist")))

	api.AttachHandlers(http.DefaultServeMux)

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

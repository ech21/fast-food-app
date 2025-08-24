package main

import (
	"fmt"
	"net/http"

	"github.com/ech21/fast-food-app/internal/api"
	"github.com/ech21/fast-food-app/internal/game"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./ui/dist")))

	api.Attach(http.DefaultServeMux)
	game.Attach(http.DefaultServeMux)

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

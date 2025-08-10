package main

import (
	"fmt"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello world!")
}

func main() {
	http.HandleFunc("/", defaultHandler)

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

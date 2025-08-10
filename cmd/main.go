package main

import (
	"fmt"
	"net/http"
)

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Pong!")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./ui/dist")))
	http.HandleFunc("/ping", ping)

	fmt.Println("Listening on port 3000")
	http.ListenAndServe(":3000", nil)
}

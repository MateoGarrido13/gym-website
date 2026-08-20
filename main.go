package main

import (
	"fmt"
	"net/http"
)

func main() {
	staticDir := "./static"
	fileserver := http.FileServer(http.Dir(staticDir))

	http.Handle("/", fileserver)

	port := "8080"
	fmt.Printf("Server is running on port %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}
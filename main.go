package main

import (
	"fmt"
	"net/http"
)
	//MVC HANDLERS, ONNION, SLICE , ARCHIVO DE HANDLERS (POST, GET, PUT, DELETE)

func main() {
	staticDir := "./static"
	fileserver := http.FileServer(http.Dir(staticDir))

	http.Handle("/", fileserver)

	port := "8080"
	fmt.Printf("Servidor iniciado en http://localhost:%s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}

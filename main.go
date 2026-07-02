package main

import (
	"net/http"
)

var Mux = http.NewServeMux()

var Server = http.Server{
	Addr:    ":8080",
	Handler: Mux,
}

func main() {
	Server.ListenAndServe()
}

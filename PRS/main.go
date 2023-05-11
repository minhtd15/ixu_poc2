package main

import (
	"PRS/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/orders", service.OrderHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

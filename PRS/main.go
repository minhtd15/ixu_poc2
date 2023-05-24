package main

import (
	"PRS/controller"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/orders", controller.OrderController).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

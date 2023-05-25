package main

import (
	_ "PRS/client"
	"PRS/controller"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("godror", "system/oracle@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))")
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		controllerDB := controller.NewOrderController(db)
		controllerDB.OrderController(w, r)
	}).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

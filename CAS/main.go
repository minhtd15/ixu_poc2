package main

import (
	"CAS/controller"
	"CAS/service"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	db, err := sql.Open("godror", "system/oracle@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))")
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	fmt.Println("Successfully connected to Oracle")

	cashService := service.NewCashServiceDB(db)
	cashController := controller.NewCashController(cashService)

	r := mux.NewRouter()
	r.HandleFunc("/cash", cashController.CashController).Methods("PUT")
}

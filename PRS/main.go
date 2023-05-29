package main

import (
	"PRS/client"
	_ "PRS/client"
	"PRS/controller"
	"PRS/service"
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

	productService := service.NewProductService(db)
	orderClient := client.NewOrderClient("http://localhost:8081")

	// Tạo một instance của controller và inject ProductService vào
	orderController := controller.NewOrderController(productService, orderClient, db)

	router := mux.NewRouter()
	router.HandleFunc("/order", orderController.OrderController).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

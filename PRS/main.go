package main

import (
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

	// Tạo một instance của controller và inject ProductService vào
	orderController := controller.NewOrderController(productService)

	router := mux.NewRouter()
	router.HandleFunc("/order", orderController.OrderController).Methods("POST")

	// router.HandleFunc("/movies/", CreateMovie).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

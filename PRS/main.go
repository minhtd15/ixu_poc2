package main

import (
	"PRS/client"
	"PRS/controller"
	"PRS/service"
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("godror", "system/oracle@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))")
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	fmt.Println("RabbitMQ Connector")
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ server: %v", err)
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ Instance")

	productService := service.NewProductService(db)
	orderClient := client.NewOrderClient("http://localhost:8081")

	// Tạo một instance của controller và inject ProductService vào
	orderController := controller.NewOrderController(productService, orderClient, db)

	router := mux.NewRouter()
	router.HandleFunc("/order", orderController.OrderController).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

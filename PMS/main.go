package PMS

import (
	"PAS/controller"
	"PMS/service"
	"database/sql"
	"fmt"
	"log"
)

func main() {
	db, err := sql.Open("godror", "system/oracle@(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=TCP)(HOST=localhost)(PORT=1521)))(CONNECT_DATA=(SERVICE_NAME=orclpdb1)))")
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	fmt.Println("Successfully connected to Oracle")

	paymentService := service.NewPaymentHistoryService(db)
	paymentController := controller.NewPaymentHistoryController(paymentService, db)

}

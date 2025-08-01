package main

import (
	"fmt"
	"log"
	"net/http"

	"flashlight-backend/api"
	"flashlight-backend/db"
)

func main() {
	fmt.Println("running main")
	if err := db.InitDB(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	defer sqlDB.Close()

	if err := db.DB.AutoMigrate(&api.Student{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	srv := api.NewServer()

	log.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", srv)
}

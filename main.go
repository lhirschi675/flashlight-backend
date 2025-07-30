package main

import (
	"flashlight-backend/api"
	"fmt"
	"net/http"
)

func main() {
	srv := api.NewServer()

	fmt.Println("Server Listening on :8080")
	http.ListenAndServe(":8080", srv)
}

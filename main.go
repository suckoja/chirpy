package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 1. Create a new ServeMux
	mux := http.NewServeMux()

	// 2. Register a handler function !!NOT YET!!

	// 3. Listen and serve on port 8080 using the custom mux
	fmt.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
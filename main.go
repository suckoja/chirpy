package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Create a file server handler pointer to a root directory
	fileServer := http.FileServer(http.Dir("."))

	// Register the file server to the root path
	mux.Handle("/", fileServer)

	// Listen and serve on port 8080 using the custom mux
	fmt.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
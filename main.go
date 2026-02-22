package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Root: serves index.html from current dir
	mux.Handle("/", http.FileServer(http.Dir(".")))

	// Assets: serves everything in the ./assets folder
	// under the /assets/ URL path.
	// StripPrefix removes "/assets/" from the URL before 
	// looking in the directory.
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	// Listen and serve on port 8080 using the custom mux
	fmt.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
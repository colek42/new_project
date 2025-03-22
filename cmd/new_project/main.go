package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleRoot)
	
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `
<!DOCTYPE html>
<html>
<head>
    <title>New Project</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 40px;
            line-height: 1.6;
        }
        h1 {
            color: #333;
        }
    </style>
</head>
<body>
    <h1>Welcome to New Project!</h1>
    <p>This is a simple Go web application deployed with Cloud Build and Cloud Run.</p>
</body>
</html>
	`)
}
package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	logDir       = "/app/data"
	logFileName  = "unsubmarine.log"
	internalPort = "8080"
)

func main() {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	http.HandleFunc("/unsubscribe", unsubscribeHandler)

	log.Printf("Unsubscribe service listening on internal port %s", internalPort)
	if err := http.ListenAndServe(":"+internalPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {

	// Only GET Requests accepted
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve email from URL

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "<h1>Bad Request</h1><p>Email address was not provided.</p>", http.StatusBadRequest)
		return
	}

	log.Printf("Received unsubscribe request for: %s", email)

	// Log request
	logFilePath := filepath.Join(logDir, logFileName)

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to open log file: %v", err)
		http.Error(w, "<h1>Error</h1><p>Your request could not be processed.</p>", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	timestamp := time.Now().UTC().Format(time.RFC3339)
	logEntry := fmt.Sprintf("%s: Unsubscribe request for %s\n", timestamp, email)

	if _, err := file.WriteString(logEntry); err != nil {
		log.Printf("ERROR: Failed to write to log file: %v", err)
		http.Error(w, "<h1>Error</h1><p>Your request could not be processed.</p>", http.StatusInternalServerError)
		return
	}

	// Serve Confirmation Page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	safeEmail := html.EscapeString(email)

	fmt.Fprintf(w, `
<!DOCTYPE html>  <html lang="en">  <head>  <meta charset="UTF-8">  <meta name="viewport" content="width=device-width, initial-scale=1.0">  <title>Unsubscribe Confirmation</title>  <style>  body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; background-color: #f8f9fa; color: #343a40; margin: 0; }  .container { text-align: center; padding: 2rem; border-radius: 8px; background-color: #ffffff; box-shadow: 0 4px 8px rgba(0,0,0,0.1); }  h1 { color: #007bff; }  p { font-size: 1.1rem; }  .email { font-weight: bold; color: #e83e8c; }  </style>  </head>  <body>  <div class="container">  <h1>Unsubscribe Confirmed</h1>  <p>Your address, <span class="email">%s</span>, has been successfully removed from our mailing list.</p>  <p>Thank you.</p>  </div>  </body>  </html> 
		`, safeEmail)

}

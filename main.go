package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type HealthResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

type ShortenRequest struct {
	LongURL string `json:"longUrl"`
}

type ShortenResponse struct {
	ShortURL string `json:"shortendUrl"`
}

type ShortURLResponse struct {
	LongURL string `json:"longURL"`
}

type ShortURL struct {
	ID       int
	ShortURL string
	LongURL  string
}

func main() {
	// Initialize the logger
	logger := log.New(os.Stdout, "ShortenerAPI ", log.LstdFlags)

	// Connect to the database
	err := connectToDB()
	if err != nil {
		logger.Fatalf("Failed to connect to the database: %v", err)
	}

	// Close the database connection when the application exits
	defer db.Close()

	// Start the HTTP server
	http.HandleFunc("/api/shorten", handleShortenURL)
	http.HandleFunc("/health", handleHealthCheck)

	serverAddr := ":8080"
	logger.Printf("Server listening on %s\n", serverAddr)
	err = http.ListenAndServe(serverAddr, nil)
	if err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}

func handleShortenURL(w http.ResponseWriter, r *http.Request) {
	// Initialize the logger
	logger := log.New(os.Stdout, "ShortenerAPI ", log.LstdFlags)

	if r.Method != http.MethodPost {
		logger.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Println("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var shortenReq ShortenRequest
	err = json.Unmarshal(body, &shortenReq)
	if err != nil {
		logger.Println("Failed to decode request body")
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	longURL := shortenReq.LongURL

	// Step 1: Check if the long URL already exists in the database
	shortURL, err := getShortURLFromDB(longURL)
	if err != nil {

		logger.Println("Failed to query database-s")
		http.Error(w, "Failed to query database-m", http.StatusInternalServerError)
		return
	}

	if shortURL == "" {
		// Step 2: Generate hash-based short URL
		shortURL, err = generateShortURL(longURL)
		if err != nil {
			logger.Println("Failed to generate short URL")
			http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
			return
		}

		// Step 3: Check if the generated short URL already exists in the database
		exists, err := checkShortURLExists(shortURL)
		if err != nil {
			logger.Println("Failed to check short URL existence")
			http.Error(w, "Failed to check short URL existence", http.StatusInternalServerError)
			return
		}

		if exists {
			// Step 4: Append a predefined string to the short URL and regenerate
			shortURL, err = regenerateShortURL(shortURL)
			if err != nil {
				logger.Println("Failed to regenerate short URL")
				http.Error(w, "Failed to regenerate short URL", http.StatusInternalServerError)
				return
			}
		}

		// Step 5: Save the shortened URL in the database
		err = saveShortURLToDB(shortURL, longURL)
		if err != nil {
			logger.Println("Failed to save URL to the database -E")
			http.Error(w, "Failed to save URL to the database_M", http.StatusInternalServerError)
			return
		}
	}

	response := ShortenResponse{
		ShortURL: shortURL,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		logger.Println("Failed to encode JSON response")
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func generateShortURL(longURL string) (string, error) {
	hash := generateRandomString(8)
	return hash, nil
}

func regenerateShortURL(shortURL string) (string, error) {
	// Append a predefined string to the short URL
	return shortURL + "-new", nil
}
func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func checkShortURLExists(shortURL string) (bool, error) {
	query := "SELECT COUNT(*) FROM short_urls WHERE short_url = ?"
	row := db.QueryRow(query, shortURL)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func getShortURLFromDB(longURL string) (string, error) {
	query := "SELECT shortURL FROM short_urls WHERE longURL = ?"
	row := db.QueryRow(query, longURL)

	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return shortURL, nil
}

func saveShortURLToDB(shortURL, longURL string) error {
	query := "INSERT INTO short_urls (shortURL, longURL) VALUES (?, ?)"
	_, err := db.Exec(query, shortURL, longURL)
	if err != nil {
		return err
	}

	return nil
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	healthStatus := "OK"
	healthDesc := "Application is healthy"

	err := db.Ping()
	if err != nil {
		healthStatus = "Error"
		healthDesc = fmt.Sprintf("Failed to connect to the database: %v", err)
	}

	response := HealthResponse{
		Status:      healthStatus,
		Description: healthDesc,
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func connectToDB() error {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(db:3306)/short_urls")
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

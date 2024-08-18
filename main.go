package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:id`
	OriginalURL  string    `json:original_url`
	ShortURL     string    `json:short_url`
	CreationDate time.Time `json:creation_date`
}

var urlDB = make(map[string]URL)

func generateRandomID() string {
	// Generate a random 4-byte array
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func generateShortUrl(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("haser:", hasher)

	data := hasher.Sum(nil)
	fmt.Println("hasher.Sum(nil)", data)

	hash := hex.EncodeToString(data)
	fmt.Println("EncodeToString", hash)
	fmt.Println("Final String", hash[:8])

	return hash[:8]

}

func createUrl(originalUrl string) string {
	shortUrl := generateShortUrl(originalUrl)
	// id := generateRandomID()
	id := shortUrl //Todo: replace this with  generateRandomID()

	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalUrl,
		ShortURL:     shortUrl,
		CreationDate: time.Now(),
	}
	return shortUrl
}

func getUrl(id string) (URL, error) {
	url, ok := urlDB[id]

	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func shortUrlHandler(w http.ResponseWriter, r *http.Request) {

	var data struct {
		URL string `json:url`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result := createUrl(data.URL)
	// fmt.Fprintf(w, "Short URL: %s", result)

	// send resp as json
	response := struct {
		ShortURL string `json:short_url`
	}{
		ShortURL: result,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// if r.Method == "POST" {
	// 	r.ParseForm()
	// 	originalUrl := r.Form.Get("url")
	// 	shortUrl := createUrl(originalUrl)
	// 	fmt.Fprintf(w, "Short URL: %s", shortUrl)
	// } else {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// }
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get Method")
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]

	fmt.Println(id)
	url, err := getUrl(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	OriginalURL := "https://github.com/faizanahmedsy" //	TODO: make this dynamic
	generateShortUrl(OriginalURL)

	// SERVER HANDLERS
	// todo: move to separate file

	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", shortUrlHandler)
	http.HandleFunc("/redirect/", redirectHandler)

	// SERVER START

	fmt.Println("Server is starting......... on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error on starting server", err)
	}

}

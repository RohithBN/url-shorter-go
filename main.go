package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

type Url struct {
	urls map[string]string
}

func (u *Url) shortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		customKey := r.FormValue("customKey")
		if customKey == "" {
			if url == "" {
				http.Error(w, "URL is required", http.StatusBadRequest)
				return
			}
			key := generateShortKey()
			u.urls[key] = url
			shortenedUrl := "http://localhost:8080/" + key


			w.Header().Set("Content-Type", "text/html")
			ResponseHtml := fmt.Sprintf(`
		<h2>URL Shortener</h2>
		<p>Original URL: %s</p>
		<p>Shortened URL: <a href="%s">%s</a></p>
		<form method="post" action="/shorten">
			<input type="text" name="url" placeholder="Enter a URL">
			<input type="text" name="customKey" placeholder="Enter the custom key you would like">
			<input type="submit" value="Shorten">
		</form>`, url, shortenedUrl, shortenedUrl)
			fmt.Fprint(w, ResponseHtml)
		} else {
			if url == "" {
				http.Error(w, "URL is required", http.StatusBadRequest)
			}
			_, ok := u.urls[customKey]
			if ok {
				http.Error(w, "Custom key already exists", http.StatusBadRequest)
				return
			}
			u.urls[customKey] = url
			shortenedUrl := "http://localhost:8080/" + customKey
			w.Header().Set("Content-Type", "text/html")
			ResponseHtml := fmt.Sprintf(`
		<h2>URL Shortener</h2>
		<p>Original URL: %s</p>
		<p>Shortened URL: <a href="%s">%s</a></p>
		<form method="post" action="/shorten">
			<input type="text" name="url" placeholder="Enter a URL">
			<input type="text" name="customKey" placeholder="Enter the custom key you would like">
			<input type="submit" value="Shorten">
		</form>`, url, shortenedUrl, shortenedUrl)
			fmt.Fprint(w, ResponseHtml)
		}
	} else {
		// Handle GET request: Display the form
		w.Header().Set("Content-Type", "text/html")
		ResponseHtml := `
		<h2>URL Shortener</h2>
		<form method="post" action="/shorten">
			<input type="text" name="url" placeholder="Enter a URL">
			<input type="text" name="customKey" placeholder="Enter the custom key you would like">
			<input type="submit" value="Shorten">
			
		</form>`
		fmt.Fprint(w, ResponseHtml)
	}
}

func (u *Url) redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	if key == "" {
		http.NotFound(w, r)
		return
	}

	url, ok := u.urls[key]
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

// generateShortKey generates a random short key
func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(uint64(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func main() {
	shortner := &Url{
		urls: make(map[string]string),
	}

	http.HandleFunc("/shorten", shortner.shortenUrl)

	http.HandleFunc("/", shortner.redirect)

	fmt.Println("URL Shortener is running on :8080")
	http.ListenAndServe(":8080", nil)
}

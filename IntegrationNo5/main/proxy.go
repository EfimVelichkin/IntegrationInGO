package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	backendURLs := []*url.URL{
		{Scheme: "http", Host: "localhost:8080"},
		{Scheme: "http", Host: "localhost:8081"},
	}

	proxy := http.NewServeMux()
	proxy.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		backendURL := backendURLs[0]
		backendURLs = append(backendURLs[1:], backendURLs[0])

		proxy := httputil.NewSingleHostReverseProxy(backendURL)
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8082", proxy))
}

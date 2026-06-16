package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("."))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			http.ServeFile(w, r, "html/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	addr := ":8080"
	log.Println("Starting", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

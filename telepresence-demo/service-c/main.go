package main

import (
	"fmt"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	if caller := r.URL.Query().Get("caller"); caller != "" {
		fmt.Printf("Request from %s\n", caller)
	}
	fmt.Fprintf(w, "Response from service-c! %s", time.Now())
}

func main() {
	fmt.Println("service-c started...")
	http.HandleFunc("/", greet)
	http.ListenAndServe(":80", nil)
}

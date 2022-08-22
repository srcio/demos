package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	if caller := r.URL.Query().Get("caller"); caller != "" {
		fmt.Printf("request from %s\n", caller)
	}

	callServiceC()
	fmt.Fprintf(w, "Response from service-b! %s", time.Now())
}

func callServiceC() {
	url := "http://service-c.default"
	if urlEnv := os.Getenv("SERVICE_C_URL"); urlEnv != "" {
		url = urlEnv
	}
	res, err := http.Get(url + "?caller=service-b")
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}
	fmt.Println(string(body))
}

func main() {
	fmt.Println("service-b started...")
	http.HandleFunc("/", greet)
	http.ListenAndServe(":80", nil)
}

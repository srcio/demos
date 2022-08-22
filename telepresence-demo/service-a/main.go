package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func callServiceB(w http.ResponseWriter, r *http.Request) {
	url := "http://service-b.default"
	if urlEnv := os.Getenv("SERVICE_B_URL"); urlEnv != "" {
		url = urlEnv
	}
	res, err := http.Get(url + "?caller=service-a")
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
	fmt.Fprintf(w, string(body))
}

func main() {
	fmt.Println("service-a started...")
	http.HandleFunc("/", callServiceB)
	http.ListenAndServe(":80", nil)
}

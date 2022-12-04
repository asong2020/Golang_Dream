package main

import (
	"fmt"
	"net/http"
)

func getProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "asong")
}

func main() {
	http.HandleFunc("/profile", getProfile)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("http server failed, err: %v\n", err)
		return
	}
}

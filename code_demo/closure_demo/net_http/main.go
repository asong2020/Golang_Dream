package main

import (
	"fmt"
	"log"
	"net/http"
)

type DecoratorHandler func(http.HandlerFunc) http.HandlerFunc

func MiddlewareHandlerFunc(hp http.HandlerFunc, decors ...DecoratorHandler) http.HandlerFunc {
	for d := range decors {
		dp := decors[len(decors)-1-d]
		hp = dp(hp)
	}
	return hp
}

func VerifyHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" {
			fmt.Fprintf(w,r.URL.Path +" response: Not Logged in")
			return
		}
		h(w,r)
	}
}

func Pong(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w,r.URL.Path +"response: pong")
	return
}


func main()  {
	http.HandleFunc("/api/asong/ping",MiddlewareHandlerFunc(Pong,VerifyHeader))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

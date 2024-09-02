package main

import (
	"fmt"
	"net/http"
	"tool/pkg/web_socket"
)

func main() {

	h := web_socket.NewHub()

	go h.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("index")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		client := web_socket.Join(w, r)
		if client == nil {
			return
		}
	})
	http.ListenAndServe(":8080", nil)

}

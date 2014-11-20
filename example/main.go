package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Some text here2.")
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("e:/home"))))
}

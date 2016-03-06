package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	switch r.Method {
	case "GET":
		get(w, r)
	case "POST":
		post(w, r)
	}
}

var comments = []string{}

func get(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(comments)
}

func post(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comment")
	log.Println("Add comment", comment)
	comments = append(comments, comment)
}

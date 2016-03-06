package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	log.Println("Starting MessageService")

	if err := clear(); err != nil {
		log.Println("Error: issue with clearing:", err)
	}

	http.HandleFunc("/healthCheck", healthHandler)
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

type message struct {
	RoomID   string
	Nickname string
	Text     string
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Figure out how to ping Mongo before returning
	fmt.Fprint(w, "")
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHandler(w, r)
	case "POST":
		postHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("mongo")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	var msg message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println("Error: failed decoding message:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	log.Println("Saving message", msg)

	c := session.DB("Jhat").C("message")
	err = c.Insert(&msg)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Message saved", msg)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("mongo")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	roomID := r.FormValue("roomID")
	if roomID == "" {
		log.Println("Error: no roomID passed to get")
		http.Error(w, "Error: no roomID passed to get", http.StatusBadRequest)
		return
	}

	var messages []message
	if err := session.DB("Jhat").C("message").Find(bson.M{"roomid": roomID}).All(&messages); err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if messages == nil {
		messages = []message{}
	}

	js, err := json.Marshal(messages)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func clear() error {
	session, err := mgo.Dial("mongo")
	if err != nil {
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	return session.DB("Jhat").C("message").DropCollection()
}

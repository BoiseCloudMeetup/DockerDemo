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
	log.Println("Starting UserService")

	if err := clear(); err != nil {
		log.Println("Error: issue with clearing:", err)
	}

	http.HandleFunc("/healthCheck", healthHandler)
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

type user struct {
	RoomID   string
	Nickname string
	Key      string
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
		createHandler(w, r)
	case "DELETE":
		deleteHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Create handler enter")
	session, err := mgo.Dial("mongo")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	var usr user
	json.NewDecoder(r.Body).Decode(&usr)

	c := session.DB("Jhat").C("user")
	err = c.Insert(&usr)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Created:", usr)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("mongo")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	var usr user
	json.NewDecoder(r.Body).Decode(&usr)

	c := session.DB("Jhat").C("user")
	err = c.Remove(&usr)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get handler enter")
	session, err := mgo.Dial("mongo")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	roomID := r.FormValue("roomID")
	if roomID == "" {
		log.Println("Error: roomID wasn't specified in handler")
		http.Error(w, "roomID must be specified", http.StatusBadRequest)
		return
	}
	log.Println("Finding users in room", roomID)

	var users []user
	if err := session.DB("Jhat").C("user").Find(bson.M{"roomid": roomID}).All(&users); err != nil {
		log.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Users:", users)

	js, err := json.Marshal(users)
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

	return session.DB("Jhat").C("user").DropCollection()
}

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/send", sendHandler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var cs = make(connections)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request URL:", r.URL)
	roomID := r.FormValue("roomID")
	if roomID == "" {
		log.Println("Error: roomID wasn't specified in handler")
		http.Error(w, "roomID must be specified", http.StatusBadRequest)
		return
	}

	nickname := r.FormValue("nickname")
	if nickname == "" {
		log.Println("Error: nickname wasn't specified in handler")
		http.Error(w, "nickname must be specified", http.StatusBadRequest)
		return
	}

	u := user{
		RoomID:   roomID,
		Nickname: nickname,
		Key:      getKey(),
	}

	cs.Broadcast("userenter", u.RoomID, u)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error: upgrade failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		cs.Broadcast("userexit", u.RoomID, u)
		return
	}

	if err := createUser(u); err != nil {
		log.Println("Error: couldn't create user at service:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		cs.Broadcast("userexit", u.RoomID, u)
		ws.Close()
		return
	}

	c := connection{ws, make(chan []byte, 1024), u}

	cs.Add(u.RoomID, u.Key, &c)

	go c.writePump()
	c.readPump()
}

func createUser(u user) error {

	b, err := json.Marshal(&u)
	if err != nil {
		return err
	}

	if _, err := http.Post("http://userService/", "application/json", bytes.NewBuffer(b)); err != nil {
		return err
	}

	return nil
}

func deleteUser(u user) error {

	b, err := json.Marshal(&u)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("DELETE", "http://userService/", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	if _, err := http.DefaultClient.Do(request); err != nil {
		return err
	}

	return nil
}

func getKey() string {
	resp, err := http.Get("http://keyService/")
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

type user struct {
	RoomID   string
	Nickname string
	Key      string
}

type message struct {
	RoomID   string
	Nickname string
	Text     string
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering send handler")
	var msg message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil || msg.RoomID == "" || msg.Nickname == "" || msg.Text == "" {
		log.Println("Error: error decoding body", err)
		http.Error(w, "body sucks", http.StatusBadRequest)
		return
	}

	cs.Broadcast("message", msg.RoomID, msg)

	b, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error: error marshaling body", err)
		http.Error(w, "my bad", http.StatusBadRequest)
		return
	}

	_, err = http.Post("http://messageService/", "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("Error: error marshaling body", err)
		http.Error(w, "my bad", http.StatusBadRequest)
		return
	}
}

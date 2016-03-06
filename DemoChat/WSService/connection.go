package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Connections
type connections map[string]map[string]*connection

func (cs connections) Add(id, key string, c *connection) {
	log.Println("Adding connection")
	if _, exists := cs[id]; !exists {
		cs[id] = make(map[string]*connection, 0)
	}

	if _, exists := cs[id][key]; exists {
		log.Fatal("OMG, we shouldn't get the same key twice!!")
	}

	cs[id][key] = c
	log.Println("Added connection", cs[id][key])
}

func (cs connections) Remove(id, key string) {
	if _, exists := cs[id]; exists {
		delete(cs[id], key)
	}
}

func (cs connections) Get(id string) map[string]*connection {
	connectionList, exists := cs[id]
	if !exists {
		return nil
	}
	return connectionList
}

type wsMessage struct {
	Action  string
	Content interface{}
}

func (cs connections) Broadcast(action, roomID string, v interface{}) {

	wsm := wsMessage{
		Action:  action,
		Content: v,
	}

	b, err := json.Marshal(&wsm)
	if err != nil {
		log.Println("Error: couldn't martial struct for", roomID)
	}

	connectionList := cs.Get(roomID)
	for _, c := range connectionList {
		c.send <- b
	}
}

// Connection
type connection struct {
	ws   *websocket.Conn
	send chan []byte
	u    user
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		cs.Remove(c.u.RoomID, c.u.Key)
		deleteUser(c.u)
		cs.Broadcast("userexit", c.u.RoomID, c.u)
		c.ws.Close()
		log.Println("Closing read pump")
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
		log.Println("Closing write pump")
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

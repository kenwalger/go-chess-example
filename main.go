package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type move struct {
	ID   string
	Move string
}

type gameReq struct {
	Username string `json:"username"`
}

type changeEvent struct {
	OperationType string     `bson:"operationType" json:"operationType"`
	FullDocument  gameWithID `bson:"fullDocument" json:"fullDocument"`
}

var clients = make(map[string][]*websocket.Conn)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	eventChan := make(chan changeEvent)

	coll := connectMongo()

	go watchForChanges(coll, eventChan)

	http.Handle("/", http.FileServer(http.Dir("./frontend/dist")))
	http.HandleFunc("/ws", makeHandleWebsockets(eventChan, coll))
	http.HandleFunc("/start", makeHandleStart(coll))

	log.Fatal(http.ListenAndServe("localhost:9001", nil))
}

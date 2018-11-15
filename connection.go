package main

import (
	"log"
	"net/http"
	"sync"	// contains mutex and such
	
	"github.com/gorilla/websocket"
)

type connection struct {
	// Channel receiving or sending only bytes
	send chan []byte
	h *hub
}

type WSHandler struct {
	h *hub
}

func (c *connection) reader(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for {
		_, message, err := wsConn.ReadMessage()	// check what is returned in this function
		if err != nil {
			break
		}
		c.h.broadcast <- message		// send information into channel.. check syntax
	}
}

// writer range over all information in channel.  Need to check related functions.
func (c *connection) writer(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for message := range c.send {
		err := wsConn.WriteMessage(websocket.TextMessage, message)	// Need to check what is returned and types
		if err != nil {
			break
		}
	}
}

// Need to check what this is and where its located.
var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}


func (wsh WSHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Need to check what this function actually does...
	wsConn, err := upgrader.Upgrade(w,req,nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}

	// Need to check what this is actually doing
	c := &connection{send: make(chan []byte), h: wsh.h}

	c.h.addConnection(c)			// adding C to its own hub...did know you could do this....
	defer c.h.removeConnection(c)	// defer removal of actual connection until other connections are closed
	
	var wg sync.WaitGroup			
	wg.Add(2)						// adding a waitgroup to ensure that program does not close before read/write
	go c.writer(&wg, wsConn)		// starting separate writer routine
	go c.reader(&wg, wsConn)		// starting separate reader routine
	wg.Wait()
	wsConn.Close()
}

package other

import (
	"log"
	"sync"
	"time"
)

type Hub struct {
	connectionsMx sync.RWMutex
	connections map[*Connection]struct{}
	broadcast chan []byte
	logMx sync.RWMutex
	log [][]byte
}

func NewHub() *Hub {
	h := &Hub{
		connectionsMx: sync.RWMutex{},
		broadcast: make(chan []byte),
		connections: make(map[*Connection]struct{}),
	}

	go func() {
		for {
			// Check syntax
			msg := <-h.broadcast
			h.connectionsMx.RLock()
			for c := range h.connections {
				select {
				case c.send <- msg:
				case <-time.After(1 *time.Second):		
					log.Printf("TIMEOUT: shutting down this connection %s", c)
					h.removeConnection(c)
				}
			}
			h.connectionsMx.RUnlock()
		}
		}()
		return h
}

// addConnection adds new connection to hub array of connections.
func (h *Hub) addConnection(conn *Connection) {
	// first the connections mutex locks the object so a new connection can be added
	h.connectionsMx.Lock()	
	defer h.connectionsMx.Unlock()  // unlock is defered until the actions are complete
	h.connections[conn] = struct{}{}
}

// removeConnection locks object mutex so specific connection can be removed.
func (h *Hub) removeConnection(conn *Connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	if _, ok := h.connections[conn]; ok {
		delete(h.connections, conn)
		close(conn.send)
	}
}

package socket

import (
  "code.google.com/p/go.net/websocket"
)

type Hub struct {
  // Registered connections.
  connections map[*Connection]bool

  // Inbound messages from the connections.
  broadcast chan string

  // Register requests from the connections.
  register chan *Connection

  // Unregister requests from connections.
  unregister chan *Connection
}

/**
 * Returns a new instance of a hub.
 */
func NewHub() *Hub {
  h := Hub{
    connections: make(map[*Connection]bool),
    broadcast:   make(chan string),
    register:    make(chan *Connection),
    unregister:  make(chan *Connection),
  }

  return &h
}

/**
 * Handles incoming socket connections.
 */
func (hub *Hub) Handler() websocket.Handler {
  return websocket.Handler(func(ws *websocket.Conn) {
    c := &Connection{
      send: make(chan string, 256),
      ws: ws,
      h: hub,
    }
    hub.register <- c
    defer func() { hub.unregister <- c }()
    go c.Writer()
    c.Reader()
  })
}

func (hub *Hub) Run() {
  for {
    select {
    case c := <-hub.register:
      hub.connections[c] = true
    case c := <-hub.unregister:
      delete(hub.connections, c)
      close(c.send)
    case m := <-hub.broadcast:
      for c := range hub.connections {
        select {
        case c.send <- m:
        default:
          delete(hub.connections, c)
          close(c.send)
          go c.ws.Close()
        }
      }
    }
  }
}


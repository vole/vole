package socket

import (
  "code.google.com/p/go.net/websocket"
)

type Connection struct {
  h *Hub
  // The websocket connection.
  ws *websocket.Conn

  // Buffered channel of outbound messages.
  send chan string
}

func (c *Connection) Reader() {
  for {
    var message string
    err := websocket.Message.Receive(c.ws, &message)
    if err != nil {
      break
    }
    c.h.broadcast <- message
  }
  c.ws.Close()
}

func (c *Connection) Writer() {
  for message := range c.send {
    err := websocket.Message.Send(c.ws, message)
    if err != nil {
      break
    }
  }
  c.ws.Close()
}


package easyws

import (
	"net/http"
	"code.google.com/p/go.net/websocket"
)

type Connection struct {
	ws   *websocket.Conn
	send chan string
	h    *Hub
}

type Hub struct {
    connections  map[*Connection]bool
    receiver     chan msginfo
    register     chan *Connection
    unregister   chan *Connection
	onjoin       func(*Connection, *Hub)
}

type msginfo struct {
	conn *Connection
	data string
}

func (c *Connection) Send(message string) {
	c.send <- message
}

func (c *Connection) reader() {
	for {
		var data string
		err := websocket.Message.Receive(c.ws, &data)
		if err != nil {
			break
		}
		message := msginfo{conn: c, data: data}
		c.h.receiver <- message
	}
	c.ws.Close()
}

func (c *Connection) writer() {
	for message := range c.send {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func wsHandler(h *Hub, ws *websocket.Conn) {
	c := &Connection{send: make(chan string, 256), ws: ws, h: h}
	h.onjoin(c, h)
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}

func (h *Hub) run(handle func(string, *Connection, *Hub)) {
    for {
        select {
        case c := <-h.register:
            h.connections[c] = true
        case c := <-h.unregister:
            delete(h.connections, c)
            close(c.send)
        case m := <-h.receiver:
			handle(m.data, m.conn, h)
        }
    }
}

func Socket(path string,
	msgHandle func(string, *Connection, *Hub), 
    joinHandle func(*Connection, *Hub)) {
	h := &Hub{
		receiver:    make(chan msginfo),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
		onjoin:      joinHandle,
	}
	go h.run(msgHandle)
	http.Handle(path, websocket.Handler(func(ws *websocket.Conn){
		wsHandler(h, ws)
	}))
}


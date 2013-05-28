easyws
======

Super-abstracted API for WebSockets in Go. A generalized, reusable form of the WebSocket server implemented by [Gary Burd](http://gary.beagledreams.com/page/go-websocket-chat.html).

### The API

The `easyws` package contains exactly one exported function: 
`easyws.Socket(path string, msgHandle func(string, *Connection, *Hub), joinHandle func(*Connection, *Hub))`

The Socket function creates a new WebSocket at the given path and accepts two handlers: `msgHandle` for when a message is sent to the WebSocket, and `joinHandle` for when a user connects to the WebSocket. Both handlers are called with the corresponding `Connection` and `Hub`. Those datatypes are as follows:

```go
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
```

A `Connection` corresponds to a single connection on the WebSocket, and a `Hub` is a structure which holds data about a WebSocket. Note that a `Connection` struct holds another struct from the `websocket` package, which you can read more about [here](https://code.google.com/p/go/source/browse/websocket/websocket.go?repo=net).

### Example

Luckily, using `easyws` is, as the name implies, easy! Here's a short example:

```go
import (
    "fmt"
    "log"
    "github.com/willcrichton/easyws"
)

func wsOnMessage(msg string, c *easyws.Connection, h *easyws.Hub){
    fmt.Println("Received message: " + msg)
}

func wsOnJoin(c *easyws.Connection, h *easyws.Hub){
    fmt.Println("New user connected")
}

func main(){
    easyws.Socket("/ws", wsOnMessage, wsOnJoin)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

In the above example, once your Go executable is running, you could then connect to your WebSocket in Javascript as follows:

```javascript
ws = new WebSocket("ws://localhost:8080/ws");
ws.send("I AM A MESSAGE");
````

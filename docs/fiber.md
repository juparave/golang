# fiber

## Setup fiber

    $ go get github.com/gofiber/fiber/v2

## Websockets with Fiber

```go
// routes.go
// Websockets sets ws services
func Websockets(app *fiber.App) {
	// ref: https://github.com/gofiber/websocket
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			// get user from session and set userID on locals
			user := getUserFromContext(c)
			c.Locals("userID", user.ID)
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// start listening
	go handlers.SocketListen()

	app.Get("/ws/:id", websocket.New(handlers.Websocket))
}
```

```go
// websockets.go
package handlers

import (
	"log"
	"strconv"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Socket struct {
	Channels   map[string]*websocket.Conn
	Register   chan *websocket.Conn
	Funcpipe   chan *SocketMessage
	Unregister chan *websocket.Conn
}

var socket Socket

type SocketMessage struct {
	Body       string
	Connection *websocket.Conn
	mu         sync.Mutex
}

// Concurrency handling - sending messages
// ref: https://medium.com/swlh/handle-concurrency-in-gorilla-web-sockets-ade4d06acd9c
func (c *SocketMessage) Send(message interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Connection.WriteJSON(message)
}

func CreateSocket() {
	socket.Channels = make(map[string]*websocket.Conn)
	socket.Register = make(chan *websocket.Conn)
	socket.Funcpipe = make(chan *SocketMessage)
	socket.Unregister = make(chan *websocket.Conn)
}

func GetSocket() *Socket {
	return &socket
}

func Websocket(c *websocket.Conn) {
	defer func() {
		// when read message error, server will try to reject user and channel
		socket.Unregister <- c
		// unregister <- c
		c.Close()
	}()
    
  	// c.Locals is added to the *websocket.Conn
	// log.Println(c.Locals("allowed"))     // true
	// log.Println(c.Locals("userID"))      // true
	// log.Println(c.Params("id"))          // 123
	// log.Println(c.Query("v"))            // 1.0
	// log.Println(c.Cookies("session_id")) // ""

	// try to colect user connection
	socket.Register <- c

	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err, mt)
			break
		}

		// send message to funcpipe
		socket.Funcpipe <- &SocketMessage{
			Body:       string(msg),
			Connection: c,
		}
	}
}

// SocketListen infinite loop
func SocketListen() {
	for {
		select {
		case connection := <-socket.Register:
			// append user's first time connection
			socket.Channels[connection.Params("id")] = connection

		case message := <-socket.Funcpipe:
			// handle function for user
			id, _ := strconv.Atoi(message.Connection.Params("id"))
			userID := message.Connection.Locals("userID").(uint)
			if uint(id) != userID {
				// id's don't match, maybe close the connection
				// TODO: websocket close connection if id's dont match
				app.InfoLog.Println("!!!! id != userID")
			}

		case connection := <-socket.Unregister:
			// user has disconected
			socket.Channels[connection.Params("id")] = nil
		}
	}
}
```

```go
...
type Message struct {
	Scope     string `json:"scope"`
	Event     string `json:"event"`
	TaskSid   string `json:"taskSid"`
	WorkerSid string `json:"workerSid"`
}

func TSocketListen() {
	for {
		select {
		case message := <-socket.Twiliopipe:
			app.InfoLog.Println("TSocketListen to event:", string(message.Body))
			// handle message for twilio
			var msg Message
			err := json.Unmarshal([]byte(message.Body), &msg)
			if err != nil {
				fmt.Println("Error decoding JSON:", err)
				return
			}
			if msg.Scope == "call" {
				if msg.Event == "disconnect" {
					// after call is disconnected, if there is a task
					// associated, the task should be marked as `completed` or `canceled`
					UpdateTaskState(msg.TaskSid, "completed", "disconnected")
				}
			}
			if msg.Scope == "twilio" {
				if msg.Event == "getEnvironment" {
					// the client is asking for environment values
					// return workspace Activities values
					app.InfoLog.Println("got socket event getEnvironment")
					err := message.Send(fiber.Map{"activities": app.Twilio.Workspace.Activities})
					if err != nil {
						app.ErrorLog.Println("socket error:", err.Error())
					}
				}
			}
		}
	}
```


## Middlewares

### Authentication

```go
func IsAuthenticated(c *fiber.Ctx) error {
	path := c.Path()

    // if not in `/api` path, let it through
	if !strings.Contains(path, "/api") {
		// not restricted
		return c.Next()
	}

    // get JWT from cookie or header
	jwt := util.GetJWT(c)

	userID, err := util.ParseJwt(jwt)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

    // Get full user's data from the database
	var user models.User
	if err = database.DB.Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
        // if user not found
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// create a string of roles, separated by commas
	roles := ""
	for _, role := range user.Roles {
		roles += role.Name + ","
	}

    // Store the user's data in the context
	c.Locals("user", user)
	c.Locals("roles", roles)

	return c.Next()
}
```

### Authorization

```go
// IsAdmin checks if the user has "admin" role
func IsAdmin(c *fiber.Ctx) error {
	// Retrieve the current user's roles from the request context
	roles := c.Locals("roles")

	// Check if the user has "admin" role
	if util.Contains(strings.Split(roles.(string), ","), "admin") {
		fmt.Println("admin")
		// If the user has "admin" role, call c.Next() to proceed to the next handler
		return c.Next()
	}

	// If the user doesn't have "admin" role, return an error response
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"message": "Access denied. You must have 'admin' role to access this resource.",
	})
}
```

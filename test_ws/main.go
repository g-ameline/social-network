package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	var server *http.ServeMux = http.NewServeMux()

	server.HandleFunc("/", func(responder http.ResponseWriter, request *http.Request) {
		fmt.Println("request hit home")
		page :=
			`<!DOCTYPE html>
	<html lang="en">
	<head>
	    <meta charset="UTF-8">
	    <title>Htmx simple websocket example</title>
	    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
	    <script src="https://unpkg.com/htmx.org@1.9.6/dist/ext/ws.js"></script>
	</head>
	<body>

	<div hx-ext="ws" ws-connect="/ws">
	  Notificaciones:
	    <p/>
	  <div id="notificaciones"></div>
	  </p>
	  <form ws-send hx-posted="javascript:this.reset();">
	    <input type="search" name="country">
	  </form>
	</div>

	</body>
	</html>
	`
		fmt.Fprintf(responder, page)
	})

	server.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request hit /ws")
		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("error upgrade", err)
			return
		}
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Println(err)
				}
				break
			}
			var message map[string]string
			err = json.Unmarshal(data, &message)
			fmt.Println("message", message, err)
		}
	})
	fmt.Println("listen and serve", http.ListenAndServe(":1234", server))
}

// http.HandleFunc("/", index)
// http.ListenAndServe(":8080", nil)

package main

import (
	d "app_server/dom"
	"fmt"

	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

func render_user_chat(data group) {
	c.Green("render user chat occurrence")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	penpal_id := data["penpal_id"].(string)
	penpal_entity := get_user(penpal_id)
	subscribe(socket, user_entity)
	socket_to_last_data[socket] = data

	chat_node := d.New_div().Text("user chat").Id("user_chat").Circle(penpal_id)
	fmt.Println("CHAT ENTITY --------------------")
	fmt.Println(get_user_chat(user_entity, penpal_entity))
	for _, missive := range get_user_chat(user_entity, penpal_entity) {
		chat_node.Add_kid(
			d.New_p(missive).Style("color", "black"),
		)
	}
	chat_node.Bear_kid(message_input_node_with_emojis("user_chat_missive", "missive", func(a_node d.Node) { a_node.Circle(user_id) })).
		Hidden_name_value("user_id", user_id).
		Hidden_name_value("penpal_id", penpal_id)

	packaged_node := package_node_toward_interactivity_frame(chat_node)
	htmx_fragment := packaged_node.Inline()
	err := socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	m.Must(err, "failed during socket's message processing")
}

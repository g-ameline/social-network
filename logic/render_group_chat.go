package main

import (
	d "app_server/dom"
	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

func render_group_chat(data group) {
	c.Green("render group chat occurrence")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	subscribe(socket, group_entity)
	socket_to_last_data[socket] = data

	group_chat_entity := get_group_chat(group_entity.id())

	chat_node := d.New_div().Text("group chat").Id("group_chat").Outline(group_id)
	for _, missive := range group_chat_entity {
		chat_node.Add_kid(
			d.New_p(missive).Style("color", "black"),
		)
	}
	styler := func(a_node d.Node) { a_node.Outline(group_id) }
	chat_node.Bear_kid(message_input_node_with_emojis("group_chat_missive", "missive", styler).
		Text("").
		Hidden_name_value("group_id", group_entity.id()))

	packaged_node := package_node_toward_interactivity_frame(chat_node)
	htmx_fragment := packaged_node.Inline()
	err := socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	m.Must(err, "failed during socket's message processing")
}

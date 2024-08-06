package main

import (
	d "app_server/dom"
	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

func render_group_invite_outsiders(data sa) {
	c.Green("render group invit outsiders occurrence")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	subscribe(socket, group_entity)
	socket_to_last_data[socket] = data

	nodes := d.Nodes{}
	group_info_node := group_info_node(group_entity)
	nodes = append(nodes, group_info_node)

	outsiders_ids := group_entity.outsiders()
	inviting_node := d.New_div().Text("invite outsiders").Id("invite_outsiders_" + group_entity.id())
	for _, outsider_id := range ordered_keys(outsiders_ids) {
		outsider_entity := get_user(outsider_id)
		// review_outsider_node :=
		_ = inviting_node.Bear_kid(message_button_node("group_invite", "outsider_id", outsider_entity.id()).
			Hidden_name_value("group_id", group_entity.id()).
			Circle(outsider_id).
			Text("invite " + outsider_entity.id() + " " + outsider_entity.email()),
		)
	}
	nodes = append(nodes, inviting_node)
	packaged_node := package_node_toward_interactivity_frame(nodes)
	htmx_fragment := packaged_node.Inline()
	err := socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	m.Must(err, "failed during socket's message processing")

}

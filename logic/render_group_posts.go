package main

import (
	d "app_server/dom"
	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

func render_group_posts(data sa) {
	c.Green("render group posts occurrence")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	subscribe(socket, group_entity)
	socket_to_last_data[socket] = data

	nodes := d.Nodes{}
	group_info_node := group_info_node(group_entity)
	nodes = append(nodes, group_info_node)

	styler := func(a_node d.Node) { a_node.Outline(group_id) }
	styler_with_key := func(a_node d.Node, color_key string) { a_node.Outline(color_key) }
	the_posts_node :=
		d.New_div().
			Text("group posts").
			Id("group_posts").
			Margin_V("5").
			Outline(group_id)
	// get psots
	posts_ids := group_entity.posts()
	for post_id := range posts_ids {
		post_entity := get_note(post_id)
		subscribe(socket, post_entity)
		post_node := the_posts_node.Bear_kid(note_node(post_entity, "group_post"))
		tree_node_the_posts(group_entity, post_node, post_entity, styler_with_key)
	}

	// outline_recursively(the_posts_node, group_id)
	nodes = append(nodes, the_posts_node)

	new_group_post_node :=
		one_input_and_one_upload_node(
			group_entity.id(),
			"note_post_group_create",
			ss{"group_id": group_entity.id()},
			styler,
		)
	nodes = append(nodes, new_group_post_node)

	packaged_node := package_node_toward_interactivity_frame(nodes)
	htmx_fragment := packaged_node.Inline()
	err := socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	m.Must(err, "failed during socket's message processing")
}

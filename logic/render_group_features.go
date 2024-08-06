package main

import (
	d "app_server/dom"

	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

func render_group_features(data sa) {
	c.Green("render group features occurrance")
	socket := data["socket"].(*ws.Conn)
	user_id := data["user_id"].(string)
	group_id := data["group_id"].(string)
	unsubscribe(socket)
	socket_to_last_data[socket] = data
	subscribe(socket, get_group(group_id))
	htmx_fragment, err := fragment_group_features(socket, user_id, group_id)
	m.If_nil_must[nvm](err, func() error {
		return socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	})
	m.Must(err, "failed during socket's message processing")
}

func fragment_group_features(socket *ws.Conn, user_id, group_id string) (string, error) {
	c.Gray("user_id", user_id, "groupd_id", group_id)
	group_entity := get_group(group_id)
	user_entity := get_user(user_id)
	group_entity.add_socket(socket)
	group_info_and_features_nodes := group_info_and_features_nodes(group_entity, user_entity)
	packaged_created_node := package_node_toward(group_info_and_features_nodes, "interactivity_frame")
	return packaged_created_node.Inline(), error(nil)
}

func group_info_and_features_nodes(group_entity group, user_entity user) d.Nodes {
	nodes := []d.Node{}
	group_info_node := group_info_node(group_entity)
	nodes = append(nodes, group_info_node)
	switch group_entity.membership(user_entity) {
	case "created":
		nodes = append(nodes,
			review_applications_node(group_entity),
			invite_outsiders_node(group_entity),
			events_node(group_entity),
			group_posts_node(group_entity),
			group_chat_node(group_entity),
		)
	case "joined":
		nodes = append(nodes,
			invite_outsiders_node(group_entity),
			events_node(group_entity),
			group_posts_node(group_entity),
			group_chat_node(group_entity),
		)
	case "invited":
		nodes = append(nodes,
			assent_node(group_entity),
			decline_node(group_entity),
		)
	case "applied":
		nodes = append(nodes,
			unapply_node(group_entity),
		)
	case "applicable":
		nodes = append(nodes,
			apply_node(group_entity),
		)
	default:
		panic("wrong membership: " + group_entity.membership(user_entity))
	}
	return nodes
}

func review_applications_node(group_entity group) d.Node {
	return message_button_node("group_review_applications", "group_id", group_entity.id()).
		Outline(group_entity.id()).
		Text("review applications")
}

func invite_outsiders_node(group_entity group) d.Node {
	return message_button_node("group_invite_outsiders", "group_id", group_entity.id()).
		Outline(group_entity.id()).
		Text("invite outsiders")
}

func events_node(group_entity group) d.Node {
	return message_button_node("group_events", "group_id", group_entity.id()).
		Outline(group_entity.id()).
		Text("events")
}

func group_posts_node(group_entity group) d.Node {
	return message_button_node("note_group_posts", "group_id", group_entity.id()).
		Outline(group_entity.id()).
		Text("posts")
}

func group_chat_node(group_entity group) d.Node {
	return message_button_node("group_chat", "group_id", group_entity.id()).
		Outline(group_entity.id()).
		Text("chat")
}

func apply_node(group_entity group) d.Node {
	group_id := group_entity.id()
	return message_button_node("group_apply", "group_id", group_id)
}
func unapply_node(group_entity group) d.Node {
	group_id := group_entity.id()
	return message_button_node("group_unapply", "group_id", group_id)
}
func assent_node(group_entity group) d.Node {
	group_id := group_entity.id()
	return message_button_node("group_assent", "group_id", group_id)
}
func decline_node(group_entity group) d.Node {
	group_id := group_entity.id()
	return message_button_node("group_decline", "group_id", group_id)
}
func admit_node(applicant_id string) d.Node {
	return message_button_node("group_admit", "user_id", applicant_id)
}
func reject_node(applicant_id string) d.Node {
	return message_button_node("group_reject", "user_id", applicant_id)
}

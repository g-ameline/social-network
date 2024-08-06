package main

import (
	d "app_server/dom"
	"fmt"

	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

// -------------PROFILE
func render_user_list(data sa) {
	c.Green("render user list occurrence")
	socket := data["socket"].(*ws.Conn)
	user_id := data["user_id"].(string)
	unsubscribe(socket)
	socket_to_last_data[socket] = data
	user_entity := get_user(user_id)
	subscribe(socket, get_global_users())
	htmx_fragment, err := fragment_user_list(user_entity)
	m.If_nil_must[nvm](err, func() error {
		return socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	})
	m.Must(err, "failed during socket's message processing")
}
func fragment_user_list(user_entity user) (string, error) {
	packaged_node := package_node_toward_interactivity_frame(user_list_categories_node(user_entity))
	return packaged_node.Inline(), error(nil)
}

func user_list_categories_node(user_entity user) d.Nodes {
	nodes := []d.Node{}
	nodes = append(nodes,
		user_itself_node(user_entity),
		list_followers_node(user_entity),
		list_followees_node(user_entity),
		list_soliciters_node(user_entity),
		list_strangers_node(user_entity),
	)
	return nodes
}

func user_itself_node(user_entity user) d.Node {
	return message_button_node("user_profile", "user_id", user_entity.id()).
		Circle(user_entity.id()).
		Text(fmt.Sprint("user id :", user_entity.id(), " email :", user_entity.email()))
}

func list_followers_node(gazee_entity user) d.Node {
	followers_nodes := []d.Node{}
	for _, follower_id := range ordered_keys(gazee_entity.followers()) {
		follower_e := get_user(follower_id)
		followers_nodes = append(followers_nodes,
			message_button_node("user_profile", "gazee_id", follower_id).
				Text(fmt.Sprint("user id :", follower_e.id(), " email :", follower_e.email())).
				Circle(follower_id))
	}
	user_info_node := accordion_node("followers", followers_nodes...)
	return user_info_node
}

func list_followees_node(gazee_entity user) d.Node {
	followees_nodes := []d.Node{}
	for _, followee_id := range ordered_keys(gazee_entity.followees()) {
		followee_e := get_user(followee_id)
		followees_nodes = append(followees_nodes,
			message_button_node("user_profile", "gazee_id", followee_id).
				Text(fmt.Sprint("user id :", followee_e.id(), " email :", followee_e.email())).
				Circle(followee_id),
		)
	}
	user_info_node := accordion_node("followed users", followees_nodes...)
	return user_info_node
}

func list_soliciters_node(gazee_entity user) d.Node {
	soliciters_nodes := []d.Node{}
	for _, soliciter_id := range ordered_keys(gazee_entity.soliciters()) {
		soliciter_e := get_user(soliciter_id)
		soliciters_nodes = append(soliciters_nodes,
			message_button_node("user_profile", "gazee_id", soliciter_id).
				Text(fmt.Sprint("user id :", soliciter_e.id(), " email :", soliciter_e.email())).
				Circle(soliciter_id),
		)
	}
	user_info_node := accordion_node("want to follow you", soliciters_nodes...)
	return user_info_node
}

func list_strangers_node(user_entity user) d.Node {
	strangers_nodes := []d.Node{}
	all_users_ids := get_users()
	for _, a_user_id := range ordered_keys(all_users_ids) {
		if user_entity.followers()[a_user_id] {
			continue
		}
		if user_entity.followees()[a_user_id] {
			continue
		}
		if user_entity.solicitees()[a_user_id] {
			continue
		}
		if user_entity.id() == a_user_id {
			continue
		}
		stranger_e := get_user(a_user_id)
		strangers_nodes = append(strangers_nodes,
			message_button_node("user_profile", "gazee_id", a_user_id).
				Text(fmt.Sprint("user id :", stranger_e.id(), " email :", stranger_e.email())).
				Circle(a_user_id))
	}
	list_of_stranger_node := accordion_node("others", strangers_nodes...)
	return list_of_stranger_node
}

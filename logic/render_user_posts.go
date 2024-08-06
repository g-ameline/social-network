package main

import (
	d "app_server/dom"
	"fmt"
	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

func render_note_user_posts(data sa) {
	c.Green("render user posts occurrence")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)

	gazer_id := data["user_id"].(string)
	gazer_entity := get_user(gazer_id)
	gazee_id := data["gazee_id"].(string)
	gazee_entity := get_user(gazee_id)
	subscribe(socket, gazee_entity)
	socket_to_last_data[socket] = data

	nodes := d.Nodes{}

	the_posts_node :=
		d.New_div().
			Text("user posts").
			Id("user_posts").
			Margin_V("5").
			Circle(gazee_id)
	// get potss
	// PUBLIC
	public_notes_node := the_posts_node.Bear_kid(d.New_div())
	posts_ids := gazee_entity.posts_public()
	for _, post_id := range ordered_keys(posts_ids) {
		post_entity := get_note(post_id)
		subscribe(socket, post_entity)
		post_node := public_notes_node.Bear_kid(note_node(post_entity, "public_post"))
		tree_node_the_posts(gazee_entity, post_node, post_entity)
	}
	// PRIVATE
	private_notes_node := the_posts_node.Bear_kid(d.New_div())
	posts_ids = gazee_entity.posts_private()
	for _, post_id := range ordered_keys(posts_ids) {
		if gazer_entity.is_follower_of(gazee_entity) || gazer_entity.id() == gazee_entity.id() {
			post_entity := get_note(post_id)
			subscribe(socket, post_entity)
			post_node := private_notes_node.Bear_kid(note_node(post_entity, "private_post"))
			tree_node_the_posts(gazee_entity, post_node, post_entity)
		}
	}
	// EXCLU
	exclusive_notes_node := the_posts_node.Bear_kid(d.New_div()).Text("EXKLU NOTES")
	var posts_ids_addresses_ids ssb = gazee_entity.posts_exclusive()
	fmt.Println("\n", "RENDER EXCLU POST")
	for _, post_id := range ordered_keys(posts_ids_addresses_ids) {
		addresses_ids := posts_ids_addresses_ids[post_id]
		if addresses_ids[gazer_id] || gazee_entity.id() == gazer_id {
			fmt.Println("  checked that watcher can see stuff")
			post_entity := get_note(post_id)
			subscribe(socket, post_entity)
			post_node := exclusive_notes_node.Bear_kid(exclu_post_node(post_entity, gazer_entity))
			tree_node_the_posts(gazee_entity, post_node, post_entity)
		}
	}
	nodes = append(nodes, the_posts_node)

	packaged_node := package_node_toward_interactivity_frame(nodes)
	htmx_fragment := packaged_node.Inline()
	err := socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	m.Must(err, "failed during socket's message processing")
}

func render_note_user_new_post(data sa) {
	c.Green("render new user posts occurrence")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	socket_to_last_data[socket] = data

	new_public_user_post_node :=
		one_input_and_one_upload_node(
			user_entity.id(),
			"note_post_public_create",
			ss{"user_id": user_entity.id()},
		).Margin_V("5")
	new_private_user_post_node :=
		one_input_and_one_upload_node(
			user_entity.id(),
			"note_post_private_create",
			ss{"user_id": user_entity.id()},
		).Margin_V("5")

	new_excluisve_user_post_node :=
		one_input_and_one_upload_node(
			user_entity.id(),
			"note_post_exclusive_create",
			ss{"user_id": user_entity.id()},
		).Margin_V("5")
	nodes := d.Nodes{
		new_public_user_post_node,
		new_private_user_post_node,
		new_excluisve_user_post_node,
	}

	packaged_node := package_node_toward_interactivity_frame(nodes)
	htmx_fragment := packaged_node.Inline()
	err := socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	m.Must(err, "failed during socket's message processing")
}

func exclu_post_node(post_entity note, gazer_entity user) d.Node {
	if post_entity.predecessor() != "" {
		panic("note a post")
	}
	// JUST CREATE NORMAL POST NODE
	post_node := d.New_div()
	// writing_node :=
	post_node.Bear_kid(d.New_node("h3")).
		Id("exclu_post_"+post_entity.id()).
		Class("card").
		Style("font-size", "27px").
		Text_black().
		Text(fmt.Sprintln(post_entity.text())).
		Add_kid(d.New_img(post_entity.picture()))
	// JUST IF POST AUTHOR IS CLIENT OFFER TO SELECT VIEWERS
	fmt.Println("  gazer' s exclusive posts", gazer_entity.posts_exclusive())
	fmt.Println("  post entity", gazer_entity.posts_exclusive())
	fmt.Println("  overlap ?", gazer_entity.posts_exclusive()[post_entity.id()])
	if _, ok := gazer_entity.posts_exclusive()[post_entity.id()]; ok {
		fmt.Println("  gazer is exclu post owner we offer reader selection tool")
		poster_entity := gazer_entity
		// add reader selection option
		followers_nodes := []d.Node{}

		for _, follower_id := range ordered_keys(gazer_entity.followers()) {

			follower_e := get_user(follower_id)
			is_confidant := poster_entity.posts_exclusive()[post_entity.id()][follower_id]
			var occurrence, text string
			if !is_confidant {
				occurrence = "note_divulge"
				text = "allow reading"
			}
			if is_confidant {
				occurrence = "note_hide"
				text = "block reading"
			}
			followers_nodes = append(
				followers_nodes,
				message_button_node(occurrence, "confidant_id", follower_id).
					Hidden_name_value("note_id", post_entity.id()).
					Text(fmt.Sprint(text, " : ", follower_id, " : ", follower_e.email())),
			)
		}
		confidants_selection_node := accordion_node("which follower should see it ?", followers_nodes...)
		post_node.Bear_kid(confidants_selection_node)
	}
	return post_node
}

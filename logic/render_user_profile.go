package main

import (
	d "app_server/dom"

	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

// -------------PROFILE
func render_user_profile(data sa) {
	c.Green("render user profile occurrence")
	socket := data["socket"].(*ws.Conn)
	user_id := data["user_id"].(string)
	unsubscribe(socket)
	socket_to_last_data[socket] = data
	gazer_entity := get_user(user_id)
	gazee_id := func() string {
		if gazee_id, ok := data["gazee_id"]; ok {
			return gazee_id.(string)
		}
		return user_id
	}()
	gazee_entity := get_user(gazee_id)
	subscribe(socket, gazee_entity)
	htmx_fragment, err := fragment_user_profile(gazer_entity, gazee_entity)
	m.If_nil_must[nvm](err, func() error {
		return socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	})
	m.Must(err, "failed during socket's message processing")
}

func fragment_user_profile(gazer_entity, gazee_entity user) (string, error) {
	node_user_profile := package_node_toward_interactivity_frame(user_features_node(gazer_entity, gazee_entity))
	return node_user_profile.Inline(), error(nil)
}

func user_features_node(gazer_entity, gazee_entity user) d.Nodes {
	nodes := []d.Node{}
	// is_self
	if gazer_entity.id() == gazee_entity.id() {
		return []d.Node{
			user_info_disclosed_node(gazer_entity, gazee_entity),
			posts_node(gazee_entity),
			list_followers_node(gazee_entity).Circle(gazee_entity.id()),
			list_followees_node(gazee_entity).Circle(gazee_entity.id()),
			list_soliciters_node(gazee_entity).Circle(gazee_entity.id()),
			privacy_button_node(gazee_entity),
		}
	}
	is_followee := gazer_entity.is_follower_of(gazee_entity)
	is_follower := gazee_entity.is_follower_of(gazer_entity)
	is_public := gazee_entity.privacy() == "public"
	switch true {
	case is_followee || is_public:
		nodes = append(nodes,
			user_info_disclosed_node(gazer_entity, gazee_entity),
			posts_node(gazee_entity),
			list_followers_node(gazee_entity).Circle(gazee_entity.id()),
			list_followees_node(gazee_entity).Circle(gazee_entity.id()),
			acquaintance_node(gazer_entity, gazee_entity),
			user_chat_node(gazee_entity),
		)
	case is_follower || is_public:
		nodes = append(nodes,
			user_info_closed_node(gazer_entity, gazee_entity),
			acquaintance_node(gazer_entity, gazee_entity),
			user_chat_node(gazee_entity),
		)
	default:
		nodes = append(nodes,
			user_info_closed_node(gazer_entity, gazee_entity),
			acquaintance_node(gazer_entity, gazee_entity),
		)
	}
	return nodes
}

func user_info_closed_node(gazer_entity, user_entity user) d.Node {
	info_nodes := []d.Node{
		d.New_div().Text("id : " + user_entity.id()).Circle(user_entity.id()),
		d.New_div().Text("email : " + user_entity.email()).Circle(user_entity.id()),
		picture_node(user_entity.avatar()),
	}
	user_info_node := accordion_node("user info "+gazer_entity.email(), info_nodes...).Circle(user_entity.id())
	return user_info_node
}

func user_info_disclosed_node(gazer_entity, user_entity user) d.Node {
	info_nodes := []d.Node{
		d.New_p("id : " + user_entity.id()).Circle(user_entity.id()),
		d.New_p("email : " + user_entity.email()).Circle(user_entity.id()),
		d.New_p("first name : " + user_entity.first_name()).Circle(user_entity.id()),
		d.New_p("last name : " + user_entity.last_name()).Circle(user_entity.id()),
		d.New_p("date of birth : " + user_entity.birth()).Circle(user_entity.id()),
		d.New_p("nickname : " + user_entity.nickname()).Circle(user_entity.id()),
		d.New_p("about: " + user_entity.about()).Circle(user_entity.id()),
		picture_node(user_entity.avatar()),
	}
	user_info_node := accordion_node("user info "+gazer_entity.email(), info_nodes...).Circle(user_entity.id())
	return user_info_node
}

func privacy_button_node(user_entity user) d.Node {
	switch user_entity.privacy() {
	case "private":
		return message_button_node(
			"user_profile_public",
			"",
			"").
			Text("turn public").
			Circle(user_entity.id())
	case "public":
		return message_button_node(
			"user_profile_private",
			"",
			"").
			Text("turn private").
			Circle(user_entity.id())
	default:
		panic("impossible")
	}
}

func acquaintance_node(gazer_entity, gazee_entity user) d.Node {
	acquaintance_node := d.New_div().Circle(gazee_entity.id())
	is_followee := gazer_entity.followees()[gazee_entity.id()]
	is_follower := gazer_entity.followers()[gazee_entity.id()]
	// is_follower := gazer_entity.followers()[gazee_entity.id()]
	is_solicitee := gazer_entity.solicitees()[gazee_entity.id()]
	is_soliciter := gazer_entity.soliciters()[gazee_entity.id()]
	is_public := gazee_entity.privacy() == "public"
	is_private := gazee_entity.privacy() == "private"

	var gazer_to_gazee_node d.Node
	switch true {
	case is_soliciter:
		gazer_to_gazee_node = message_button_node(
			"user_concede",
			"follower_id",
			gazee_entity.id()).
			Text("accept following request").
			Circle(gazee_entity.id())
	case is_follower:
		gazer_to_gazee_node = message_button_node(
			"user_forshake",
			"follower_id",
			gazee_entity.id()).
			Text("take back following right").
			Circle(gazee_entity.id())
	default:
		gazer_to_gazee_node = d.New_div().
			Text("is not interested in you").
			Circle(gazee_entity.id())
	}
	acquaintance_node.Bear_kid(gazer_to_gazee_node)

	var gazee_to_gazer_node d.Node
	switch true {
	case is_public && !is_followee:
		gazee_to_gazer_node = message_button_node(
			"user_follow",
			"followee_id",
			gazee_entity.id()).
			Text("follow").
			Circle(gazee_entity.id())
	case is_followee:
		gazee_to_gazer_node = message_button_node(
			"user_unfollow",
			"followee_id",
			gazee_entity.id()).
			Text("unfollow").
			Circle(gazee_entity.id())
	case is_solicitee:
		gazee_to_gazer_node = message_button_node(
			"user_unsolicit",
			"followee_id",
			gazee_entity.id()).
			Text("cancel following request").
			Circle(gazee_entity.id())
	case is_private:
		gazee_to_gazer_node = message_button_node(
			"user_solicit",
			"followee_id",
			gazee_entity.id()).
			Text("request following").
			Circle(gazee_entity.id())
	default:
		panic("shall not happen")
	}
	acquaintance_node.Bear_kid(gazee_to_gazer_node)
	return acquaintance_node
}

func posts_node(gazee_entity user) d.Node {
	return message_button_node("note_user_posts", "gazee_id", gazee_entity.id()).
		Text("posts").
		Circle(gazee_entity.id())
}
func user_chat_node(gazee_entity user) d.Node {
	return message_button_node("user_chat", "penpal_id", gazee_entity.id()).
		Text("chat with user").
		Circle(gazee_entity.id())
}

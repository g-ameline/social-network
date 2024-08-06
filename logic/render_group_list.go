package main

import (
	d "app_server/dom"

	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

// -------------CATEGORIES
func render_group_list(data sa) {
	c.Green("render group list occurrance")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	htmx_fragment, err := fragment_group_list()
	m.If_nil_must[nvm](err, func() error {
		return socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	})
	m.Must(err, "failed during socket's message processing")
	// subscribe(socket, get_global_groups())
}
func fragment_group_list() (string, error) {
	node_group_list := package_node_toward_interactivity_frame(group_categories_node())
	return node_group_list.Inline(), error(nil)
}

func group_categories_node() d.Node {
	groups_frame :=
		d.New_div().
			Id("groups_frame").
			Class(d.Container).
			Text("groups_frame")
	memberships := []string{
		"created",
		"joined",
		"invited",
		"applied",
		"applicable",
	}
	for _, membership := range memberships {
		groups_frame.Add_kid(
			message_button_node("group_list_"+membership, "membership", membership).
				Margin_V("10").
				Text(membership + " groups").
				Class(d.Container),
		)
	}
	return groups_frame
}

// -------------CREATED
func render_group_list_created(data sa) {
	c.Green("render group list created")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	user_id := data["user_id"].(string)
	htmx_fragment, err := fragment_group_list_created(socket, user_id)
	m.If_nil_must[nvm](err, func() error {
		return socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	})
	m.Must(err, "failed during socket's message processing")
	subscribe(socket, get_global_groups())
}
func fragment_group_list_created(socket *ws.Conn, user_id string) (string, error) {
	c.Gray("user_id", user_id)
	user_entity := get_user(user_id)
	c.Blue("user entity", len(user_entity))
	created_groups_ids := user_entity.created_groups()
	c.Blue("created groups by user", created_groups_ids)
	created_groups := map[string]group{}
	for _, id := range ordered_keys(created_groups_ids) {
		group_entity := get_group(id)
		created_groups[id] = group_entity
	}
	group_list_created_nodes := list_of_groups_node(created_groups, "created")
	packaged_created_node := package_node_toward(group_list_created_nodes, "interactivity_frame")
	return packaged_created_node.Inline(), error(nil)
}

// -------------JOINED-------------

func render_group_list_joined(data sa) {
	render_group_list_membership(data, "joined")
}
func render_group_list_invited(data sa) {
	render_group_list_membership(data, "invited")
}
func render_group_list_applied(data sa) {
	render_group_list_membership(data, "applied")
}
func render_group_list_applicable(data sa) {
	render_group_list_membership(data, "applicable")
}

func render_group_list_membership(data sa, membership string) {
	c.Green("render group list " + membership)
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	user_id := data["user_id"].(string)
	htmx_fragment, err := fragment_group_list_membership(user_id, membership)
	m.If_nil_must[nvm](err, func() error {
		return socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	})
	m.Must(err, "failed during socket's message processing")
	subscribe(socket, get_global_groups())
}
func fragment_group_list_membership(user_id string, membership string) (string, error) {
	c.Gray("user_id", user_id)
	user_entity := get_user(user_id)
	c.Blue("user entity", user_entity)
	groups_ids := func() sb {
		switch membership {
		case "joined":
			return user_entity.joined_groups()
		case "invited":
			return user_entity.invited_groups()
		case "applied":
			return user_entity.applied_groups()
		case "applicable":
			return user_entity.applicable_groups()
		}
		panic("wrong membership:" + membership)
	}()
	c.Blue("membership groups by user", len(groups_ids))
	groups := map[string]group{}
	for id := range groups_ids {
		group_entity := get_group(id)
		groups[id] = group_entity
	}
	group_list_membership_nodes := list_of_groups_node(groups, membership)
	packaged_membership_node := package_node_toward(group_list_membership_nodes, "interactivity_frame")
	return packaged_membership_node.Inline(), error(nil)
}
func list_of_groups_node(groups_entities map[string]group, membership string) d.Nodes {
	nodes := []d.Node{}
	c.Gray("list of group data", len(groups_entities))
	for _, group_id := range ordered_keys(groups_entities) {
		nodes = append(nodes, group_info_node(get_group(group_id)))
	}
	return nodes
}
func group_info_node(group_entity group) d.Node {
	group_id := group_entity.id()
	title := group_entity.title()
	description := group_entity.description()
	creator_id := group_entity.creator()
	// message_button_node("group_features", "group_id", group_id).

	group_info_button_node := d.New_div().
		Attr("ws-send", "").
		Attr(d.Hx_trigger, "click").
		Id("group_info_"+group_id).
		Attr("role", "button").
		Text("group info").
		Hidden_name_value("occurrence", "group_features").
		Hidden_name_value("group_id", group_id).
		Margin_V("10").
		Outline(group_id)

	group_info_button_node.Bear_kid(d.New_div().Id("group_info_container_" + group_id)).
		Class(d.Grid).
		Add_kid(d.New_div().Text("title<br> " + title).
			Outline(group_id),
		).
		Add_kid(d.New_div().Text("description<br> " + description).
			Outline(group_id),
		).
		Add_kid(d.New_div().
			Text("founder :<br> " + "email " + get_user(creator_id).email() + "<br>id " + get_user(creator_id).id()).
			Outline(group_id),
		).
		Add_kid(d.New_div().Text("group id<br> " + group_id).
			Outline(group_id),
		)

	return group_info_button_node
}

package main

import (
	d "app_server/dom"
	"fmt"
	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

func render_group_events(data sa) {
	c.Green("render group event occurrence")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	subscribe(socket, group_entity)
	socket_to_last_data[socket] = data

	nodes := d.Nodes{}
	group_info_node := group_info_node(group_entity)
	nodes = append(nodes, group_info_node)

	the_events_node := func(group_entity group, user_entity user) d.Node {
		events_node := d.New_div().
			Text("group events").
			Id("group_events").
			Outline(group_id)

		events_ids := group_entity.events()

		going_node := events_node.Bear_kid(d.New_div()).
			Text("events going").Id("group_events_attending")
		for _, event_id := range ordered_keys(events_ids) {
			event_entity := get_event(event_id)
			fmt.Println("    event", event_entity)
			if event_entity.attenders()[user_entity.id()] {
				going_node.Bear_kid(d.New_div()).
					Outline(group_id).
					Margin_V("5").
					Text(fmt.Sprintln("event", event_entity.id(), event_entity.title(), event_entity.date()))
			}
		}
		undecided_node := events_node.Bear_kid(d.New_div()).
			Text("events undecided").Id("group_events_undecided")
		for event_id := range events_ids {
			event_entity := get_event(event_id)
			if !event_entity.absentees()[user_entity.id()] && !event_entity.attenders()[user_entity.id()] {
				undecided_node.Bear_kid(d.New_div()).
					Outline(group_id).
					Class(d.Grid).
					Margin_V("5").
					Text(fmt.Sprintln("event", event_entity.id(), event_entity.title(), event_entity.date())).
					Add_kid(
						message_button_node("event_going", "event_id", event_entity.id()).
							Outline(group_id).
							Text("going"),
					).Add_kid(
					message_button_node("event_not_going", "event_id", event_entity.id()).
						Outline(group_id).
						Text("not going"),
				)
			}
		}
		return events_node.Outline(group_id).Margin_V("5")
	}
	nodes = append(nodes, the_events_node(group_entity, user_entity))
	the_new_event_node := func(group_entity group) d.Node {
		new_event_node := message_form_node("event_create",
			[]sss{
				{"title": ss{"required": ""}},
				{"description": ss{}},
				{"date": ss{
					"required": "",
					"type":     "date",
				}},
			},
			func(a_node d.Node) { a_node.Outline(group_id) },
		)
		new_event_node.
			Hidden_name_value("group_id", group_entity.id()).
			Outline(group_id)
		return new_event_node
	}
	nodes = append(nodes, the_new_event_node(group_entity))
	packaged_node := package_node_toward_interactivity_frame(nodes)
	htmx_fragment := packaged_node.Inline()
	err := socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	m.Must(err, "failed during socket's message processing")

}

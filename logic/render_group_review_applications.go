package main

import (
	d "app_server/dom"
	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

func render_group_review_applications(data sa) {
	c.Green("render group review application occurrence")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	subscribe(socket, group_entity)
	socket_to_last_data[socket] = data

	nodes := d.Nodes{}
	group_info_node := group_info_node(group_entity)
	nodes = append(nodes, group_info_node)

	applicants_ids := group_entity.applicants()
	applications_node := d.New_div().Text("review applications").Id("review_applications")
	for _, applicant_id := range ordered_keys(applicants_ids) {
		applicant_entity := get_user(applicant_id)
		// review_applicant_node :=
		applications_node.Bear_kid(d.New_div()).
			Text(applicant_entity.id() + " " + applicant_entity.email()).
			Circle(applicant_id).
			// Attr("role", "group").
			Class(d.Grid).
			Add_kid(
				message_button_node("group_admit", "applicant_id", applicant_entity.id()).
					Hidden_name_value("group_id", group_entity.id()).
					Circle(applicant_id).
					Text("admit"),
			).Add_kid(
			message_button_node("group_reject", "applicant_id", applicant_entity.id()).
				Hidden_name_value("group_id", group_entity.id()).
				Circle(applicant_id).
				Text("reject"),
		)
	}
	nodes = append(nodes, applications_node)
	packaged_node := package_node_toward_interactivity_frame(nodes)
	htmx_fragment := packaged_node.Inline()
	err := socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	m.Must(err, "failed during socket's message processing")
}

package main

import (
	d "app_server/dom"

	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

func render_group_new(data sa) {
	c.Green("render group features occurrance")
	socket := data["socket"].(*ws.Conn)
	unsubscribe(socket)
	socket_to_last_data[socket] = data
	error_hint, title, description := func() (string, string, string) {
		var error, title, description string
		if error_raw, ok := data["error"]; ok {
			error = error_raw.(string)
		}
		if title_raw, ok := data["title"]; ok {
			title = title_raw.(string)
		}
		if description_raw, ok := data["description"]; ok {
			description = description_raw.(string)
		}
		return error, title, description
	}()
	htmx_fragment, err := fragment_group_new(error_hint, title, description)
	m.If_nil_must[nvm](err, func() error {
		return socket.WriteMessage(ws.TextMessage, []byte(htmx_fragment))
	})
	m.Must(err, "failed during socket's message processing")
}
func fragment_group_new(error_hint, title_value, description_value string) (string, error) {

	new_group_node := message_form_node(
		"group_create",
		[]sss{
			sss{"title": ss{"value": title_value, "required": ""}},
			sss{"description": ss{"value": description_value, "required": ""}},
		},
	)
	m.If_wordly_must[nvm](error_hint, func() { new_group_node.Add_firstborn(d.New_div().Text(error_hint)) })
	packaged_new_group_node := package_node_toward(new_group_node, "interactivity_frame")
	return packaged_new_group_node.Inline(), error(nil)
}

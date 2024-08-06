package main

import (
	d "app_server/dom"
)

func page_logged_inlined(user_id string) string {
	return d.Prefix_doctype(page_logged_node().Inline())
}
func page_logged_node() d.Node {
	page_node, body_node := page_node()
	body_node.Bear_kid(main_frame_logged_node())
	return page_node
}

func page_node() (d.Node, d.Node) {
	page_node :=
		d.New_node(d.Html).
			Add_kid(d.New_node(d.Head).
				Add_kid(d.HTMX_script()).    // htmx stuff
				Add_kid(d.HTMX_ws_script()). // htmx stuff
				Add_kid(d.New_node(d.Link).
					Other(d.Rel, "stylesheet").
					Other(d.Href, "https://cdn.jsdelivr.net/npm/@picocss/pico@next/css/pico.min.css"),
				),
			)
	body_node :=
		page_node.
			Bear_kid(d.New_body()).
			Add_kid(banner_node())

	return page_node, body_node
}
func main_frame_logged_node() d.Node {
	main_frame_node := main_frame_node()
	main_frame_node.Attr(d.Hx_ext, "ws").Attr("ws-connect", slash(interactivity))
	main_frame_node.Add_kid(features_node())
	main_frame_node.Add_kid(notifying_node())
	main_frame_node.Add_kid(interactivity_node())
	return main_frame_node
}
func main_frame_node() d.Node {
	return d.New_div().
		Id("main_frame")
}

func features_node() d.Node {
	return banner_from_nodess(
		d.Nodes{
			message_button_node("user_profile", "", "").Text("your profile"),
			message_button_node("user_list", "", "").Text("browse other users"),
		},
		d.Nodes{
			message_button_node("note_user_new_post", "", "").Text("create new post"),
		},
		d.Nodes{
			message_button_node("group_list", "", "").Text("browse groups"),
			message_button_node("group_new", "", "").Text("create a group"),
		},
	).Id("main_features")
}

func notifying_node() d.Node {
	return banner_from_nodess(
		d.Nodes{
			d.New_div().Text("notifying_frame").
				Id("notifying_frame").
				Class(d.Grid),
		},
	).Id("main_notifying")
}

func interactivity_node() d.Node {
	interactivity_node := d.New_node("section").
		Id("interactivity_frame").
		Text("interactivity_frame").
		Class(d.Container).
		Add_kid(
			message_button_node("user_profile", "", "").
				Attr("hx-trigger", "load"),
		)
	return interactivity_node
}

// <button
//    type="button"
//    hx-trigger="load"
//    nunjucks-template="gistlist"
//    hx-target="#list"
//    hx-swap="innerHTML"
// >Reload</button>

func banner_node() d.Node {
	return banner_from_nodess(
		d.Nodes{
			d.New_button().
				Other("role", d.Button).
				Text("logout").
				Attr(d.Hx_post, slash("user_logout")).
				Attr(d.Hx_target, "#main_frame").
				Attr(d.Hx_swap, d.OuterHTML).
				Attr(d.Hx_trigger, "click"),
		},
		d.Nodes{
			d.New_a().Text("home").Href("./").Other("role", d.Button),
		},
	)
}

func user_profile_node() d.Node {
	return message_button_node("user_profile", "", "")
}
func user_list_node() d.Node {
	return message_button_node("user_list", "", "")
}
func posts_button_node() d.Node {
	return message_button_node("posts", "", "")
}
func chat_button_node() d.Node {
	return message_button_node("chat", "", "")
}
func groups_button_node() d.Node {
	return message_button_node("group_list", "", "")
}
func new_group_button_node() d.Node {
	return message_button_node("group_new", "", "")
}
func new_post_button_node() d.Node {
	return message_button_node("note_user_new_post", "", "")
}

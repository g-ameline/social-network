package main

import (
	d "app_server/dom"
)

func page_credentials_blank_inlined() string {
	return d.Prefix_doctype(page_credentials_node(ss{}, ss{}).Inline())
}

func page_credentials_inlined(login_values, register_values ss) string {
	return d.Prefix_doctype(page_credentials_node(login_values, register_values).Inline())
}

func page_credentials_node(login_values, register_values ss) d.Node {
	page_node, body_node := page_node()
	body_node.Add_kid(main_frame_credentials_node(login_values, register_values))
	return page_node
}

func main_frame_credentials_node(login_values, register_values ss) d.Node {
	main_frame_node := main_frame_node()
	main_frame_node.
		// Style(d.Overflow_y, "auto").
		Class(d.Container).
		Add_kids(credentials_nodes(login_values, register_values))
	return main_frame_node
}

func credentials_nodes(login_values, register_values ss) d.Nodes {
	login_node := d.New_div().
		Class("grid").
		Add_kid(login_node(login_values))
	register_node := d.New_div().
		Class("grid").
		Add_kid(registration_node(register_values))
	nodes := d.Nodes{login_node, register_node}
	return nodes
}

func registration_node_blank() d.Node {
	return registration_node(ss{})
}
func registration_node(values ss) d.Node {
	registration_node := d.New_div().Text("register")
	registration_node.Bear_kid(d.New_node(d.Form)).
		Class("container").
		Id("registration").
		Other(d.Hx_post, slash("user_register")).
		Other(d.Hx_target, "#main_frame").
		Other(d.Hx_swap, d.OuterHTML).
		Other(d.Hx_trigger, "submit").
		Other(d.Hx_encoding, d.Multipart).
		Add_kid(d.New_div().Text(values["error"])).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "email").Text("email")).
			Add_kid(d.New_node(d.Input).Other(d.Type, "text").
				Name("email").
				Type("email").
				Other("required", "").
				Value(values["email"])),
		).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "password").Text("password")).
			Add_kid(d.New_node(d.Input).
				Other(d.Type, "text").
				Type("password").
				Name("password").Other("required", "").
				Value(values["password"])),
		).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "first_name").Text("first name")).
			Add_kid(d.New_node(d.Input).Other(d.Type, "text").Name("first_name").Other("required", "").Value(values["first_name"])),
		).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "last_name").Text("last name")).
			Add_kid(d.New_node(d.Input).Other(d.Type, "text").Name("last_name").Other("required", "").Value(values["last_name"])),
		).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "birth").Text("date of birth")).
			Add_kid(d.New_node(d.Input).Other(d.Type, "text").
				Name("birth").
				Other("required", "").
				Type("date").
				Value(values["date_of_birth"])),
		).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "avatar").Text("avatar")).
			Add_kid(d.New_node(d.Input).Other(d.Type, "file").Name("avatar").Value(values["avatar"])),
		).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "nickname").Text("nickname")).
			Add_kid(d.New_node(d.Input).Other(d.Type, "text").Name("nickname").Value(values["nickname"])),
		).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "about").Text("about")).
			Add_kid(d.New_node(d.Input).Other(d.Type, "text").Name("about").Value(values["about"])),
		).
		Add_kid(d.New_node(d.Button).Other(d.Type, "submit").Text("submit"))
	return registration_node
}
func login_node_blank() d.Node {
	return login_node(ss{})
}
func login_node(values ss) d.Node {
	login_node := d.New_div().Text("login")
	login_node.Bear_kid(d.New_node(d.Form)).
		Id("login").
		Class("container").
		Other(d.Hx_post, slash("user_login")).
		Other(d.Hx_target, "#main_frame").
		Other(d.Hx_swap, d.OuterHTML).
		Other(d.Hx_trigger, "submit").
		Add_kid(d.New_div().Text(values["error"])).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "email").Text("email")).
			Add_kid(d.New_node(d.Input).Other(d.Type, "text").
				Name("email").
				Other("required", "").
				Type("email").
				Value(values["email"])),
		).
		Add_kid(d.New_div().
			Add_kid(d.New_node(d.Label).Other(d.For, "password").Text("password")).
			Add_kid(d.New_node(d.Input).Other(d.Type, "text").
				Type("password").
				Name("password").
				Other("required", "").
				Value(values["password"])),
		).
		Add_kid(d.New_node(d.Button).Other(d.Type, "submit").Text("submit"))
	return login_node
}

func return_home_link_node() d.Node {
	link_node := d.New_button().
		Text("invalid session please return to home page").
		Other("onclick", `location.href='`+url_app_root+`'`).
		Other("type", "button")
	return link_node
}

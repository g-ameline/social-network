package main

import (
	d "app_server/dom"
	m "github.com/g-ameline/maybe"

	"fmt"
)

func message_button_node(occurrence, name, value string) d.Node {
	click_node := d.New_div().
		Id(name).
		Text(occurrence).
		Attr("ws-send", "").
		Attr(d.Hx_trigger, "click").
		Margin_V("10").
		Attr("role", "button")
	click_node.Hidden_name_value("occurrence", occurrence)
	m.If_wordly_must[nvm](name+value, func() {
		click_node.Hidden_name_value(name, value)
	})
	return click_node
}
func message_form_node(occurrence string, names_attributess_valuess []sss, stylers ...func(d.Node)) d.Node {
	var styler func(d.Node)
	switch len(stylers) {
	case 0:
		styler = func(_ d.Node) {}
	case 1:
		styler = stylers[0]
	default:
		panic("")
	}
	form_node := d.New_form().
		Id(occurrence).
		Text(occurrence).
		Attr("ws-send", "").
		Hidden_name_value("occurrence", occurrence)
	for _, name_attributes_values := range names_attributess_valuess {
		for name, attributes_values := range name_attributes_values {
			input_node := form_node.Bear_kid(d.New_input()).
				Id(occurrence+name).
				Name(name).
				Attr(d.Placeholder, name)
			styler(input_node)
			for attribute, value := range attributes_values {
				input_node.Attr(attribute, value)
			}
		}
	}
	submit_node := form_node.Bear_kid(d.New_button()).Type(d.Submit).Text("submit")
	styler(form_node)
	styler(submit_node)
	return form_node
}

func message_input_node(occurrence string, name string, stylers ...func(d.Node)) (d.Node, d.Node) {
	var styler func(d.Node)
	switch len(stylers) {
	case 0:
		styler = func(_ d.Node) {}
	case 1:
		styler = stylers[0]
	default:
		panic("")
	}
	form_node := d.New_form().
		Id(occurrence).
		Text(occurrence).
		Attr("ws-send", "").
		Hidden_name_value("occurrence", occurrence)
	input_node := form_node.Bear_kid(d.New_input()).
		Id(occurrence+"_input").
		Name(name).
		Attr(d.Placeholder, "type your message here").
		Attr("required", "").
		Other(d.Hx_trigger, "keyup[keyCode==13]")

	submit_node := form_node.Bear_kid(d.New_button()).Type(d.Submit).Text("submit")
	styler(form_node)
	styler(input_node)
	styler(submit_node)
	return form_node, input_node
}
func message_input_node_with_emojis(occurrence string, name string, stylers ...func(d.Node)) d.Node {
	var styler func(d.Node)
	switch len(stylers) {
	case 0:
		styler = func(_ d.Node) {}
	case 1:
		styler = stylers[0]
	default:
		panic("")
	}
	form_node, input_node := message_input_node(occurrence, name, styler)
	drawer_node := input_node.Bear_kid(emoji_drawer_node(occurrence+"_input", styler))
	styler(form_node)
	styler(input_node)
	styler(drawer_node)
	return form_node
}
func emoji_drawer_node(input_id string, stylers ...func(d.Node)) d.Node {
	var styler func(d.Node)
	switch len(stylers) {
	case 0:
		styler = func(_ d.Node) {}
	case 1:
		styler = stylers[0]
	default:
		panic("")
	}
	add_emoji_script := func(emoji, target_id string) string {
		return fmt.Sprintf(
			`document.getElementById('%v').value += '%v';`,
			target_id,
			emoji,
		)
	}

	emojis_codes := map[string]string{
		"😀": "&#128512",
		"😁": "&#128513",
		"😂": "&#128514",
		"😃": "&#128515",
		"😄": "&#128516",
		"😅": "&#128517",
	}
	emoji_drawer_id := "emoji_drawer_id"

	emoji_drawer_node :=
		d.New_node("details").
			Id(emoji_drawer_id).
			Add_kid(d.New_node("summary").Text("emojies"))

	for _, emoji_code := range emojis_codes {
		emoji_drawer_node.Bear_kid(d.New_div()).
			Text(emoji_code).
			Attr("onclick", add_emoji_script(emoji_code, input_id))
	}
	styler(emoji_drawer_node)
	return emoji_drawer_node
}
func get_button_node(name string) d.Node {
	click_frame_node := d.New_div().
		Id(name).
		Text(name).
		Class(d.Container).
		Other(d.Hx_get, slash(name)).
		Other(d.Hx_swap, d.None).
		Other(d.Hx_trigger, "click")

	return click_frame_node
}
func add_hidden_inputs_values(the_node d.Node, id string, names_values ss) d.Node {
	for name, value := range names_values {
		the_node.Add_kid(d.New_input().
			Id(id).
			Name(name).
			Value(value).
			Conceal(),
		)
	}

	return the_node
}

func package_node_toward[n d.Node | d.Nodes](some_nodes n, id string) d.Node {
	wraping_node := d.New_div().
		Id(id).
		Attr(d.Hx_swap_OOB, "true :"+id_selector(id))
	if one_node, ok := any(some_nodes).(d.Node); ok {
		wraping_node.Add_kid(one_node)
	}
	if multiple_node, ok := any(some_nodes).(d.Nodes); ok {
		wraping_node.Add_kids(multiple_node)
	}
	return wraping_node
}
func package_node_toward_interactivity_frame[n d.Node | d.Nodes](some_nodes n) d.Node {
	return package_node_toward(some_nodes, "interactivity_frame")
}
func package_node_toward_notifying_frame[n d.Node | d.Nodes](some_nodes n) d.Node {
	return package_node_toward(some_nodes, "notifying_frame")
}

func one_input_and_one_upload_node(top_id, occurrence string, needed_names_values ss, stylers ...func(d.Node)) d.Node {
	var styler func(d.Node)
	switch len(stylers) {
	case 0:
		styler = func(_ d.Node) {}
	case 1:
		styler = stylers[0]
	}
	top_node := d.New_div().Id(top_id).Pack().Text(occurrence + " " + top_id)
	styler(top_node)
	form_node := top_node.Bear_kid(d.New_form()).
		Pack().
		Id(occurrence).
		Text(occurrence).
		Attr("ws-send", "").
		Hidden_name_value("occurrence", occurrence)
	styler(form_node)
	input_id := occurrence + "_text"
	// input_node :=
	input_node := form_node.Bear_kid(d.New_input().
		Id(input_id).
		Pack().
		Name("text").
		Other("required", "").
		Attr(d.Placeholder, "text"),
	)
	styler(input_node)
	picture_id := occurrence + "_picture"
	form_node.Add_kid(d.New_div().Class("grid").
		Id(occurrence + "_picture").Pack().Text("no picture"),
	)
	top_node.Add_kid(
		upload_file_and_replace_target_with_filepath_node(id_selector(picture_id), styler),
	)
	submit_node := form_node.Add_kid(d.New_button().Type(d.Submit).Text("submit"))
	for name, value := range needed_names_values {
		form_node.Hidden_name_value(name, value)
	}
	styler(submit_node)
	return top_node
}

func upload_file_and_replace_target_with_filepath_node(target string, stylers ...func(d.Node)) d.Node {
	var styler func(d.Node)
	switch len(stylers) {
	case 0:
		styler = func(_ d.Node) {}
	case 1:
		styler = stylers[0]
	}
	form_node := d.New_form().
		Id("file_to+"+target).
		Attr(d.Hx_post, ("save_file_return_filepath")).
		Attr(d.Hx_target, target).
		Attr(d.Hx_swap, d.InnerHTML).
		Attr(d.Hx_trigger, "submit").
		Attr(d.Hx_encoding, d.Multipart).
		Class("grid").
		Add_kid(d.New_input().
			Hidden_name_value("dummy_name", "dummy_value").
			Attr(d.Type, "file").Name("file"),
		)
	upload_node := form_node.Bear_kid(d.New_input()).
		Other(d.Type, "submit").
		Value("upload")
	styler(form_node)
	styler(upload_node)
	return form_node
}

func accordion_node(caption string, below_nodes ...d.Node) d.Node {
	details_node := d.New_node("details").Attr("role", d.Button)
	// summary_node :=
	details_node.Bear_kid(d.New_node("summary").Text(caption).Text_black())
	for _, a_node := range below_nodes {
		details_node.Add_kid(a_node)
	}
	return details_node
}

func tree_node_the_posts(main_entity any_entity, parent_node d.Node, parent_note note, stylers ...func(d.Node, string)) {
	var main_styler, comment_styler func(d.Node, string)
	switch len(stylers) {
	case 0:
		main_styler = func(_ d.Node, _ string) {}
		comment_styler = func(_ d.Node, _ string) {}
	case 1:
		main_styler = stylers[0]
		comment_styler = stylers[0]
	case 2:
		main_styler = stylers[0]
		comment_styler = stylers[1]
	default:
		panic("only 2 stylers")
	}
	new_comment_styler := func(a_node d.Node) { comment_styler(a_node, main_entity.id()) }
	switch true {
	case len(parent_note.comments()) == 0:
		new_comment_node :=
			parent_node.Bear_kid(new_comment_node(main_entity, parent_note, parent_node, new_comment_styler))
		main_styler(new_comment_node, main_entity.id())
		return
	default:
		for comment_id := range parent_note.comments() {
			comment_entity := get_note(comment_id)
			comment_node := parent_node.Bear_kid(note_node(comment_entity, "comment"))
			comment_styler(comment_node, comment_entity.id())
			tree_node_the_posts(main_entity, comment_node, comment_entity)
		}
	}
}
func note_node(note_entity note, note_type string, stylers ...func(d.Node)) d.Node {
	var styler func(d.Node)
	switch len(stylers) {
	case 0:
		styler = func(_ d.Node) {}
	case 1:
		styler = stylers[0]
	}
	note_node := d.New_div().
		Id(note_type + "_" + note_entity.id())
	writing_node := note_node.Bear_kid(d.New_node("h4")).
		Class("card").
		Text_black().
		Text(fmt.Sprintln(note_entity.text()))
	note_node.Add_kid(d.New_img(note_entity.picture()))
	styler(note_node)
	styler(writing_node)
	return note_node
}
func new_comment_node(entity any_entity, note_entity note, note_node d.Node, stylers ...func(d.Node)) d.Node {
	return one_input_and_one_upload_node(
		note_entity.id(),
		"note_comment",
		ss{"id": entity.id(), "note_id": note_entity.id()},
		stylers...,
	).Text("write your comment below")
}
func banner_from_nodess(nodess ...d.Nodes) d.Node {
	nav_node := d.New_node(d.Nav)
	for _, nodes := range nodess {
		ul_node := nav_node.Bear_kid(d.New_node(d.Ul))
		for _, node := range nodes {
			ul_node.Bear_kid(d.New_node(d.Li)).Add_kid(node)
		}
	}
	return nav_node
}

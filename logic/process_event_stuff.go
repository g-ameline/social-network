package main

import d "app_server/dom"

func process_event_create(data sa) {
	// update stuff
	// socket := data["socket"].(*ws.Conn)
	user_id := data["user_id"].(string)
	// user_entity := get_user(user_id)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	if user_id != group_entity.creator() && !group_entity.members()[user_id] {
		panic("only group's creator or member can create")
	}
	title := data["title"].(string)
	description := data["description"].(string)
	date := data["date"].(string)
	event_entity := group_entity.create_event(title, description, date)
	refresh_subscribers(group_entity)
	styler := func(a_node d.Node) d.Node {
		return a_node.Outline(group_id)
	}
	creator := get_user(group_entity.creator())
	notify(creator, styler, "a group event ", event_entity.title(), " have been created, group: ", group_entity.title())
	for member_id := range group_entity.members() {
		member_entity := get_user(member_id)
		notify(member_entity, styler, "a group event ", event_entity.title(), " have been created, group: ", group_entity.title())
	}
}

func process_event_going(data sa) {
	user_id := data["user_id"].(string)
	event_id := data["event_id"].(string)
	event_entity := get_event(event_id)
	group_entity := get_group(event_entity.group())
	if user_id != group_entity.creator() && !group_entity.members()[user_id] {
		panic("only group's creator or member can join an event")
	}
	user_entity := get_user(user_id)
	user_entity.attend(event_entity)
	refresh_subscribers(group_entity)
}

func process_event_not_going(data sa) {
	user_id := data["user_id"].(string)
	event_id := data["event_id"].(string)
	event_entity := get_event(event_id)
	group_entity := get_group(event_entity.group())
	if user_id != group_entity.creator() && !group_entity.members()[user_id] {
		panic("only group's creator or member can join an event")
	}
	user_entity := get_user(user_id)
	user_entity.absent(event_entity)
	refresh_subscribers(group_entity)
}

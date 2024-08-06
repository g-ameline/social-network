package main

import (
	"fmt"
	"strings"

	d "app_server/dom"

	// c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	// ws "github.com/gorilla/websocket"
)

func process_group_apply(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	user_entity.apply(group_entity)
	refresh_subscribers(group_entity)
	styler := func(a_node d.Node) d.Node {
		return a_node.Circle(user_id)
	}
	creator := get_user(group_entity.creator())
	notify(creator, styler, "someone, ", user_entity.email(), " ", user_entity.nickname(), " want to join that group ", group_entity.title())
}
func process_group_unapply(data sa) {
	// update stuff
	// socket := data["socket"].(*ws.Conn)
	user_id := data["user_id"].(string)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	get_user(user_id).unapply(group_entity)
	refresh_subscribers(group_entity)
}
func process_group_decline(data sa) {
	// update stuff
	// socket := data["socket"].(*ws.Conn)
	user_id := data["user_id"].(string)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	get_user(user_id).decline(group_entity)
	refresh_subscribers(group_entity)
}
func process_group_invite(data sa) {
	// update stuff
	// socket := data["socket"].(*ws.Conn)
	outsider_id := data["outsider_id"].(string)
	outsider_entity := get_user(outsider_id)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	group_entity.invite(get_user(outsider_id))
	refresh_subscribers(group_entity)
	styler := func(a_node d.Node) d.Node {
		return a_node.Outline(group_entity.id())
	}
	notify(outsider_entity, styler, "you' ve been invited in group ", group_entity.title())
}
func process_group_assent(data sa) {
	// update stuff
	// socket := data["socket"].(*ws.Conn)
	user_id := data["user_id"].(string)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	user_entity := get_user(user_id)
	user_entity.assent(group_entity)
	refresh_subscribers(group_entity)
}

func process_group_admit(data sa) {
	// update stuff
	// socket := data["socket"].(*ws.Conn)
	// user_id := data["user_id"].(string)
	applciant_id := data["applicant_id"].(string)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	applicant_entity := get_user(applciant_id)
	group_entity.admit(applicant_entity)
	refresh_subscribers(group_entity)
}
func process_group_reject(data sa) {
	// update stuff
	// socket := data["socket"].(*ws.Conn)
	// user_id := data["user_id"].(string)
	group_id := data["group_id"].(string)
	applciant_id := data["applicant_id"].(string)
	group_entity := get_group(group_id)
	applicant_entity := get_user(applciant_id)
	group_entity.reject(applicant_entity)
	refresh_subscribers(group_entity)
}
func process_group_create(data sa) {
	// update stuff
	// socket := data["socket"].(*ws.Conn)
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	title := data["title"].(string)
	title = strings.ReplaceAll(title, "'", "''")
	description := data["description"].(string)
	description = strings.ReplaceAll(description, "'", "''")

	groups_entities := get_groups()
	err := error(nil)
	for id := range groups_entities {
		a_group_entity := get_group(id)
		if a_group_entity["title"] == title {
			err = fmt.Errorf("title already taken")
		}
	}
	m.If_error_must[nvm](err, func() {
		data["error"] = err.Error()
		data["title"] = title
		data["description"] = description
		data["occurrence"] = "group_create"
		process_occurrence(data)
	})
	m.If_nil_must[nvm](err, func() {
		user_entity.create_group(title, description)
		refresh_subscribers(get_global_groups())
		data["occurrence"] = "group_list"
		process_occurrence(data)
	})
}

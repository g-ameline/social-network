package main

import "fmt"

func process_note_post_group_create(data sa) {
	user_id := data["user_id"].(string)
	group_id := data["group_id"].(string)
	group_entity := get_group(group_id)
	if user_id != group_entity.creator() && !group_entity.members()[user_id] {
		panic("only group's creator or member can create")
	}
	text := data["text"].(string)
	picture := func() string {
		if picture, ok := data["picture"]; ok {
			return picture.(string)
		}
		return ""
	}()
	user_entity := get_user(user_id)
	user_entity.create_group_post(group_entity, text, picture)
	fmt.Println("\n ICI ----------------", data)
	refresh_subscribers(group_entity)

}
func process_note_post_public_create(data sa) {
	user_id := data["user_id"].(string)
	text := data["text"].(string)
	picture := func() string {
		if picture, ok := data["picture"]; ok {
			return picture.(string)
		}
		return ""
	}()
	user_entity := get_user(user_id)
	user_entity.create_public_post(text, picture)
	refresh_subscribers(user_entity)
	data["occurrence"] = "note_user_posts"
	data["gazee_id"] = user_id
	process_occurrence(data)
}
func process_note_post_private_create(data sa) {
	user_id := data["user_id"].(string)
	text := data["text"].(string)
	picture := func() string {
		if picture, ok := data["picture"]; ok {
			return picture.(string)
		}
		return ""
	}()
	user_entity := get_user(user_id)
	user_entity.create_private_post(text, picture)
	refresh_subscribers(user_entity)
	data["occurrence"] = "note_user_posts"
	data["gazee_id"] = user_id
	process_occurrence(data)
}
func process_note_post_exclusive_create(data sa) {
	user_id := data["user_id"].(string)
	text := data["text"].(string)
	picture := func() string {
		if picture, ok := data["picture"]; ok {
			return picture.(string)
		}
		return ""
	}()
	user_entity := get_user(user_id)
	user_entity.create_exclusive_post(text, picture)
	refresh_subscribers(user_entity)
	data["occurrence"] = "note_user_posts"
	data["gazee_id"] = user_id
	process_occurrence(data)
}
func process_note_comment(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	note_id := data["note_id"].(string)
	note_entity := get_note(note_id)
	text := data["text"].(string)
	var picture_file_path string
	if picture, ok := data["picture"]; ok {
		picture_file_path = picture.(string)
	}
	user_entity.comment(note_entity, text, picture_file_path)
	refresh_subscribers(note_entity.post())
}
func process_note_confide(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	note_id := data["note_id"].(string)
	note_entity := get_note(note_id)
	confidant_id := data["confidant_id"].(string)
	confidant_entity := get_user(confidant_id)
	user_entity.confide(note_entity, confidant_entity)
	refresh_subscribers(note_entity.post())
	refresh_subscribers(user_entity)
}
func process_note_hide(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	note_id := data["note_id"].(string)
	note_entity := get_note(note_id)
	confidant_id := data["confidant_id"].(string)
	confidant_entity := get_user(confidant_id)
	user_entity.distrust(note_entity, confidant_entity)
	refresh_subscribers(note_entity.post())
	refresh_subscribers(user_entity)
}

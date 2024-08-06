package main

func process_group_chat_missive(data sa) {
	user_id := data["user_id"].(string)
	group_id := data["group_id"].(string)
	missive := data["missive"].(string)
	get_group(group_id).add_missive_to_chat(get_user(user_id), missive)
	refresh_subscribers(get_group(group_id))
}
func process_user_chat_missive(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	penpal_id := data["penpal_id"].(string)
	penpal_entity := get_user(penpal_id)
	missive := data["missive"].(string)
	user_entity.add_missive_to_chat(penpal_entity, missive)
	refresh_subscribers(user_entity)
	refresh_subscribers(penpal_entity)
}

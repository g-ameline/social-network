package main

import d "app_server/dom"

func process_user_profile_private(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	user_entity.private()
	refresh_subscribers(user_entity)
}

func process_user_profile_public(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	user_entity.public()
	refresh_subscribers(user_entity)
}

func process_user_follow(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	followee_id := data["followee_id"].(string)
	followee_entity := get_user(followee_id)
	user_entity.follow(followee_entity)
	refresh_subscribers(user_entity, followee_entity)
}

func process_user_unfollow(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	followee_id := data["followee_id"].(string)
	followee_entity := get_user(followee_id)
	user_entity.unfollow(followee_entity)
	refresh_subscribers(user_entity, followee_entity)
}

func process_user_solicit(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	followee_id := data["followee_id"].(string)
	followee_entity := get_user(followee_id)
	if user_entity.followees()[followee_id] {
		panic("state error")
	}
	user_entity.solicit(followee_entity)
	refresh_subscribers(user_entity, followee_entity)
	styler := func(a_node d.Node) d.Node {
		return a_node.Circle(user_entity.id())
	}
	notify(followee_entity, styler, "user ", user_entity.nickname(), " wants to follow you")
}

func process_user_unsolicit(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	followee_id := data["followee_id"].(string)
	followee_entity := get_user(followee_id)
	if user_entity.followees()[followee_id] {
		panic("state error")
	}
	user_entity.unsolicit(followee_entity)
	refresh_subscribers(user_entity, followee_entity)
}

func process_user_concede(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	follower_id := data["follower_id"].(string)
	follower_entity := get_user(follower_id)
	if user_entity.soliciters()[follower_id] {
		user_entity.concede(follower_entity)
		refresh_subscribers(user_entity, follower_entity)
	}
}

func process_user_forshake(data sa) {
	user_id := data["user_id"].(string)
	user_entity := get_user(user_id)
	follower_id := data["follower_id"].(string)
	follower_entity := get_user(follower_id)
	user_entity.forshake(follower_entity)
	refresh_subscribers(user_entity, follower_entity)
}

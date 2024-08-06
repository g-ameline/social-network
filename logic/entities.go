package main

import (
	"fmt"

	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
	http "net/http"
)

type any_entity interface {
	id() string
	sockets() cb
	add_socket(*ws.Conn)
	remove_socket(*ws.Conn)
}

// type se[e any_entity] map[string]e

// type sse[e any_entity] map[string]map[string]e

// type ssse[e any_entity] map[string]map[string]map[string]e

type note map[string]any
type user map[string]any
type group map[string]any
type event map[string]any
type chat []string // []

var socket_to_entities = map[*ws.Conn]map[any]bool{} // [*ws.Conn][entity]true
// or [data][entity]true ..?

var entities sa = sa{
	// "users":         map[string]user{},
	// "groups":        map[string]group{},
	// "events":        map[string]event{}, //[event_id]= event
	// "notes":         map[string]note{},
	"user_chats":  map[string]map[string]chat{}, // [user_id][user_id]chat{}
	"group_chats": map[string]chat{},            // [group_id]chat{}
	"global": sa{
		"users":  user{"sockets": cb{}},
		"groups": group{"sockets": cb{}},
	},
}

// ------------- USERS

func get_users() sa {
	if users, ok := entities["users"]; ok {
		return users.(sa)
	}
	data := ss{"table": "users"}
	ids, err := send_query_to_database_and_parse_result[ss, sb](data_url.get_ids, data)
	m.Must(err, "failed to get users records from database")
	entities["users"] = sa{}
	for id := range ids {
		if _, ok := entities["users"].(sa)[id]; !ok {
			entities["users"].(sa)[id] = true
		}
	}
	return get_users()
}

func get_user(user_id string) user {
	if a_user, ok := get_users()[user_id]; ok {
		if user_entity, ok := a_user.(user); ok {
			return user_entity
		}
	}
	data := ss{
		"table":       "users",
		"key_field_1": "id",
		"key_value_1": user_id,
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.get_record, data)
	})
	m.Must(err, "failed to post data to database server")
	record, err := m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to get missing user's record from database", user_id)
	m.Must(error_from_json(record))
	new_user_entity := user{
		"id":         user_id,
		"email":      record["email"],
		"private":    record["private"],
		"first_name": record["first_name"],
		"last_name":  record["last_name"],
		"birth":      record["birth"],
		"nickname":   record["nickname"],
		"avatar":     record["avatar"],
		"about":      record["about"],
		"sockets":    cb{},
	}
	entities["users"].(sa)[user_id] = new_user_entity
	return get_user(user_id)
}

func get_global_users() user {
	return entities["global"].(sa)["users"].(user)
}

// ----------- GET DATA from user
func (the_user user) email() string {
	if email, ok := the_user["email"].(string); ok {
		return email
	}
	panic("should not happen, field should be there")
}
func (the_user user) first_name() string {
	if first_name, ok := the_user["first_name"].(string); ok {
		return first_name
	}
	panic("should not happen, field should be there")
}

func (the_user user) last_name() string {
	if last_name, ok := the_user["last_name"].(string); ok {
		return last_name
	}
	panic("should not happen, field should be there")
}

func (the_user user) nickname() string {
	if nickname, ok := the_user["nickname"].(string); ok {
		return nickname
	}
	panic("should not happen, field should be there")
}
func (the_user user) birth() string {
	if birth, ok := the_user["birth"].(string); ok {
		return birth
	}
	panic("should not happen, field should be there")
}
func (the_user user) about() string {
	if about, ok := the_user["about"].(string); ok {
		return about
	}
	panic("should not happen, field should be there")
}

func (the_user user) avatar() string {
	if avatar, ok := the_user["avatar"].(string); ok {
		return avatar
	}
	panic("should not happen, field should be there")
}

func (the_user user) privacy() string {
	if private, ok := the_user["private"].(string); ok && private == "1" {
		return "private"
	}
	if private, ok := the_user["private"].(string); ok && private == "0" {
		return "public"
	}
	panic("should not happen, field should be there")
}

func (the_user user) private() user {
	data := sa{
		"table": "users",
		"id":    the_user.id(),
		"field": "private",
		"value": "1",
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.update_record, data)
	})
	m.Must(err, "failed to post data to database server")
	returned_data, err := m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to upsert joining in database", fmt.Sprint(data))
	_, err = m.If_nil_do[nvm](err, func() error {
		error_db, ok := returned_data["error"]
		if ok {
			return fmt.Errorf(error_db)
		}
		return error(nil)
	})
	m.Must(err, "failed to upsert joining in database", fmt.Sprint(data))

	the_user["private"] = "1"

	return the_user
}

func (the_user user) public() user {
	data := sa{
		"table": "users",
		"id":    the_user.id(),
		"field": "private",
		"value": "0",
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.update_record, data)
	})
	m.Must(err, "failed to post data to database server")
	returned_data, err := m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to upsert joining in database", fmt.Sprint(data))

	m.Must(err, "failed to upsert joining in database", fmt.Sprint(data))
	_, err = m.If_nil_do[nvm](err, func() error {
		error_db, ok := returned_data["error"]
		if ok {
			return fmt.Errorf(error_db)
		}
		return error(nil)
	})
	m.Must(err, "failed to upsert joining in database", fmt.Sprint(data))

	the_user["private"] = "0"

	return the_user
}

func (the_user user) followers() sb {
	if followers, ok := the_user["followers"].(sb); ok {
		return followers
	}
	the_user.following()
	return the_user.followers()
}

func (the_user user) followees() sb {
	if followees, ok := the_user["followees"].(sb); ok {
		return followees
	}
	the_user.following()
	return the_user.followees()
}

func (the_user user) soliciters() sb {
	if soliciters, ok := the_user["soliciters"].(sb); ok {
		return soliciters
	}
	the_user.following()
	return the_user.soliciters()
}

func (the_user user) solicitees() sb {
	if solicitees, ok := the_user["solicitees"].(sb); ok {
		return solicitees
	}
	the_user.following()
	return the_user.solicitees()
}

func (the_user user) make_follower_of(other_user user) {
	if the_user.id() == other_user.id() {
		panic("")
	}
	data := sa{
		"table": "followings",
		"record": ss{
			"follower_id": the_user.id(),
			"followee_id": other_user.id(),
			"approval":    "1",
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	returned_data, err := m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to upsert following in database", fmt.Sprint(data))
	m.If_nil_must[nvm](err, func() error {
		if error_db, ok := returned_data["error"]; ok {
			return fmt.Errorf(error_db)
		}
		return error(nil)
	})
	// the_user side
	delete(the_user.solicitees(), other_user.id())
	the_user.followees()[other_user.id()] = true
	// other_user side
	delete(other_user.soliciters(), the_user.id())
	other_user.followers()[the_user.id()] = true
	if !the_user.is_follower_of(other_user) {
		panic("ici")
	}
}
func (the_user user) make_soliciter_of(other_user user) {
	if the_user.id() == other_user.id() {
		panic("")
	}
	data := sa{
		"table": "followings",
		"record": ss{
			"follower_id": the_user.id(),
			"followee_id": other_user.id(),
			"approval":    "0",
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	returned_data, err := m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to upsert following in database", fmt.Sprint(data))
	m.If_nil_must[nvm](err, func() error {
		if error_db, ok := returned_data["error"]; ok {
			return fmt.Errorf(error_db)
		}
		return error(nil)
	})

	// the_user side
	delete(the_user.followees(), other_user.id())
	delete(other_user.followers(), the_user.id())
	the_user.solicitees()[other_user.id()] = true
	// other_user side
	other_user.soliciters()[the_user.id()] = true
	if !the_user.is_soliciter_of(other_user) {
		fmt.Println(`the_user.solicitees()`)
		fmt.Println(the_user.solicitees())
		fmt.Println(`other_user.soliciters()`)
		fmt.Println(other_user.soliciters())
		fmt.Println(`the_user.is_soliciter_of(other_user)`)
		fmt.Println(the_user.is_soliciter_of(other_user))
		panic("failed to properly make one soliciter of other")
	}
}
func (the_user user) is_soliciter_of(other_user user) bool {
	if the_user.solicitees()[other_user.id()] != other_user.soliciters()[the_user.id()] {
		panic("wrong state")
	}
	if the_user.id() == other_user.id() {
		panic("same user")
	}
	return the_user.solicitees()[other_user.id()]
}

func (the_user user) make_unfamilar_of(other_user user) {
	if the_user.id() == other_user.id() {
		panic("")
	}
	data := sa{
		"table": "followings",
		"record": ss{
			"follower_id": the_user.id(),
			"followee_id": other_user.id(),
			"approval":    "0",
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	returned_data, err := m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to upsert following in database", fmt.Sprint(data))
	m.If_nil_must[nvm](err, func() error {
		if error_db, ok := returned_data["error"]; ok {
			return fmt.Errorf(error_db)
		}
		return error(nil)
	})

	// the_user side
	delete(the_user.followees(), other_user.id())
	delete(the_user.solicitees(), other_user.id())
	// other_user side
	delete(other_user.followers(), the_user.id())
	delete(other_user.soliciters(), the_user.id())
	if the_user.is_follower_of(other_user) || the_user.is_soliciter_of(other_user) {
		fmt.Println(`the_user.followees()`)
		fmt.Println(the_user.followees())
		fmt.Println(`other_user.followers()`)
		fmt.Println(other_user.followers())
		panic("koko")
	}
}

func (the_user user) is_follower_of(other_user user) bool {
	is_other_followee_of_user := the_user.followees()[other_user.id()]
	is_user_follower_of_other := other_user.followers()[the_user.id()]
	if is_other_followee_of_user != is_user_follower_of_other {
		fmt.Println(the_user.followees())
		fmt.Println(other_user.followers())
		panic("wrong state")
	}
	if the_user.id() == other_user.id() {
		panic("same user")
	}
	return is_other_followee_of_user
}
func (the_user user) is_unfamiliar_of(other_user user) bool {
	if the_user.id() == other_user.id() {
		panic("same user")
	}
	return !the_user.is_follower_of(other_user) && !the_user.is_soliciter_of(other_user)
}
func (the_user user) following() {
	data := ss{"table": "followings"}
	followings_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_file, data)
	m.Must(err, "failed to get note's followers from database", fmt.Sprint(followings_records))
	followers_ids := sb{}
	followees_ids := sb{}
	soliciters_ids := sb{}
	solicitees_ids := sb{}
	for _, record := range followings_records {

		followee_id := record["followee_id"]
		follower_id := record["follower_id"]
		if record["followee_id"] != the_user.id() && record["follower_id"] != the_user.id() {
			continue
		}
		is_approved := record["approval"] == "1"
		is_solicited := record["approval"] == "0"
		switch is_approved {
		case follower_id == the_user.id():
			followees_ids[followee_id] = true
		case followee_id == the_user.id():
			followers_ids[follower_id] = true
		}
		switch is_solicited {
		case follower_id == the_user.id():
			solicitees_ids[followee_id] = true
		case followee_id == the_user.id():
			soliciters_ids[follower_id] = true
			continue
		}
	}
	the_user["followers"] = followers_ids
	the_user["followees"] = followees_ids
	the_user["soliciters"] = soliciters_ids
	the_user["solicitees"] = solicitees_ids
}
func (the_user user) solicit(stranger user) user {
	if !the_user.is_unfamiliar_of(stranger) {
		fmt.Println(the_user)
		panic("prob")
	}
	the_user.make_soliciter_of(stranger)
	return the_user
}

func (the_user user) unsolicit(solicitee user) user {
	if !the_user.is_soliciter_of(solicitee) {
		fmt.Println(the_user)
		panic("prob")
	}
	the_user.make_unfamilar_of(solicitee)
	return the_user
}

func (the_user user) concede(soliciter user) user {
	if !soliciter.is_soliciter_of(the_user) {
		fmt.Println(the_user)
		panic("prob")
	}
	soliciter.make_follower_of(the_user)
	return the_user
}

func (the_user user) forshake(follower user) user {
	if !follower.is_follower_of(the_user) {
		fmt.Println(the_user)
		panic("prob")
	}
	follower.make_unfamilar_of(the_user)
	return the_user
}

func (the_user user) follow(untracked user) user {
	if the_user.is_follower_of(untracked) {
		fmt.Println(the_user)
		fmt.Println(untracked)
		panic("prob")
	}
	the_user.make_follower_of(untracked)
	if !the_user.is_follower_of(untracked) {
		fmt.Println(the_user)
		fmt.Println(untracked)
		panic("prob")
	}
	return the_user
}

func (the_user user) unfollow(followee user) user {
	if !the_user.is_follower_of(followee) {
		fmt.Println(`the_user`)
		fmt.Println(the_user)
		fmt.Println(the_user.followees())
		fmt.Println(the_user.followees()[followee.id()])
		fmt.Println(`followee`)
		fmt.Println(followee)
		fmt.Println(followee.followers()[the_user.id()])
		panic("prob")
	}
	the_user.make_unfamilar_of(followee)
	return the_user
}

// ---------- UTILS -------------

func unsubscribe(socket *ws.Conn) {
	if previous_linked_entities, ok := socket_to_entities[socket]; ok {
		for entity := range previous_linked_entities {
			remove_socket(entity, socket)
		}
	}
	delete(socket_to_entities, socket)
}

func subscribe[e any_entity](socket *ws.Conn, some_entities ...e) {
	for _, an_entity := range some_entities {
		if _, ok := socket_to_entities[socket]; !ok {
			socket_to_entities[socket] = map[any]bool{}
		}
		// socket_to_entities[socket][an_entity] = true
		socket_to_entities[socket][&an_entity] = true
		add_socket(an_entity, socket)
	}
}
func resubscribe_socket[e any_entity](socket *ws.Conn, some_entities ...e) {
	unsubscribe(socket)
	subscribe(socket, some_entities...)
}
func refresh_subscribers[e any_entity](entities ...e) {
	// get socket
	for _, entity := range entities {
		for socket := range get_sockets_or_panic(entity) {
			data := socket_to_last_data[socket]
			process_occurrence(data)
		}
	}
}

// ------ GENREIC METHODS

func update_something[in any](an_entity any, key string, value in) {
	switch an_entity.(type) {
	case user:
		an_entity.(user)[key] = value
	case note:
		an_entity.(note)[key] = value
	case group:
		an_entity.(group)[key] = value
	case event:
		an_entity.(event)[key] = value
	default:
		panic("wrong type")
	}
}
func remove_something(an_entity any, key string) {
	switch an_entity.(type) {
	case user:
		delete(an_entity.(user), key)
	case note:
		delete(an_entity.(note), key)
	case group:
		delete(an_entity.(group), key)
	case event:
		delete(an_entity.(event), key)
	default:
		panic("wrong type")
	}
}
func add_socket(an_entity any, socket *ws.Conn) {
	switch an_entity.(type) {
	case user:
		an_entity.(user).sockets()[socket] = true
	case note:
		an_entity.(note).sockets()[socket] = true
	case group:
		an_entity.(group).sockets()[socket] = true
	case event:
		an_entity.(event).sockets()[socket] = true
	default:
		panic("wrong type")
	}
}

func (the_entity user) add_socket(socket *ws.Conn)  { add_socket(the_entity, socket) }
func (the_entity group) add_socket(socket *ws.Conn) { add_socket(the_entity, socket) }
func (the_entity event) add_socket(socket *ws.Conn) { add_socket(the_entity, socket) }
func (the_entity note) add_socket(socket *ws.Conn)  { add_socket(the_entity, socket) }

func remove_socket(entity any, socket *ws.Conn) {
	switch entity.(type) {
	case *user:
		delete(entity.(*user).sockets(), socket)
	case *note:
		delete(entity.(*note).sockets(), socket)
	case *group:
		delete(entity.(*group).sockets(), socket)
	case *event:
		delete(entity.(*event).sockets(), socket)
	default:
		panic("wrong type")
	}
}

func (the_entity user) remove_socket(socket *ws.Conn)  { remove_socket(the_entity, socket) }
func (the_entity group) remove_socket(socket *ws.Conn) { remove_socket(the_entity, socket) }
func (the_entity event) remove_socket(socket *ws.Conn) { remove_socket(the_entity, socket) }
func (the_entity note) remove_socket(socket *ws.Conn)  { remove_socket(the_entity, socket) }

func get_something[v any](an_entity any, key string, reaction func(bool)) v {
	switch an_entity.(type) {
	case user:
		value, ok := an_entity.(user)[key].(v)
		reaction(ok)
		return value
	case note:
		value, ok := an_entity.(note)[key].(v)
		reaction(ok)
		return value
	case group:
		value, ok := an_entity.(group)[key].(v)
		reaction(ok)
		return value
	case event:
		value, ok := an_entity.(event)[key].(v)
		reaction(ok)
		return value
	default:
		panic("wrong type")
	}
}

func get_something_or_panic[v any](an_entity any, key string) v {
	return get_something[v](an_entity, key, func(ok bool) {
		if !ok {
			panic(fmt.Sprint(an_entity, key, "not found"))
		}
	})
}
func get_id_or_panic(an_entity any) string {
	return get_something_or_panic[string](an_entity, "id")
}
func (the_entity user) id() string  { return get_id_or_panic(the_entity) }
func (the_entity group) id() string { return get_id_or_panic(the_entity) }
func (the_entity event) id() string { return get_id_or_panic(the_entity) }
func (the_entity note) id() string  { return get_id_or_panic(the_entity) }

// ------ record socket method
func get_sockets_or_panic(an_entity any) cb {
	return get_something_or_panic[cb](an_entity, "sockets")
}
func (the_entity user) sockets() cb  { return get_sockets_or_panic(the_entity) }
func (the_entity group) sockets() cb { return get_sockets_or_panic(the_entity) }
func (the_entity event) sockets() cb { return get_sockets_or_panic(the_entity) }
func (the_entity note) sockets() cb  { return get_sockets_or_panic(the_entity) }

func close_socket(socket *ws.Conn) {
	entities := socket_to_entities[socket]
	for entity := range entities {
		switch entity.(type) {
		case *user:
			delete(entity.(*user).sockets(), socket)
		case *note:
			delete(entity.(*note).sockets(), socket)
		case *group:
			delete(entity.(*group).sockets(), socket)
		case *event:
			delete(entity.(*event).sockets(), socket)
		default:
			panic("wrong type")
		}
	}
	delete(socket_to_entities, socket)
	socket.Close()
}

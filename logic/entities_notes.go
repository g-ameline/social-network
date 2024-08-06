package main

import (
	"fmt"

	m "github.com/g-ameline/maybe"
	// ws "github.com/gorilla/websocket"
	http "net/http"
)

func get_notes() sa {
	if notes, ok := entities["notes"].(sa); ok {
		return notes
	}
	data := ss{"table": "notes"}
	ids, err := send_query_to_database_and_parse_result[ss, sb](data_url.get_ids, data)
	m.Must(err, "failed to get notes records from database")
	m.Must(error_from_json(ids))
	entities["notes"] = sa{}
	for id := range ids {
		if _, ok := entities["notes"].(sa)[id]; !ok {
			entities["notes"].(sa)[id] = true
		}
	}
	return get_notes()
}

func get_note(note_id string) note {
	if a_note, ok := get_notes()[note_id]; ok {
		if note_entity, ok := a_note.(note); ok {
			return note_entity
		}
	}
	data := ss{
		"table":       "notes",
		"key_field_1": "id",
		"key_value_1": note_id,
	}
	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.get_record, data)
	})
	m.Must(err, "failed to post data to database server")
	record, err := m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to get missing note's record from database", note_id)
	m.Must(error_from_json(record), "error with data")

	new_note_entity := note{
		"id":          note_id,
		"author":      record["author_id"],
		"predecessor": record["predecessor_id"],
		"date":        record["date"],
		"text":        record["text"],
		"picture":     record["picture"],
		"sockets":     cb{},
	}
	entities["notes"].(sa)[note_id] = new_note_entity
	return get_note(note_id)
}

func (the_note note) author() string {
	if author_id, ok := the_note["author"].(string); ok {
		return author_id
	}
	panic("should an author field")
}
func (the_note note) predecessor() string {
	if predecessor_id, ok := the_note["predecessor"].(string); ok {
		return predecessor_id
	}
	return ""
}
func (the_note note) post() note {
	var up_note note = the_note
	for up_note.predecessor() != "" {
		up_note = get_note(up_note.predecessor())
	}
	return up_note
}
func (the_note note) date() string {
	if date, ok := the_note["date"].(string); ok {
		return date
	}
	panic("should an datefield")
}
func (the_note note) text() string {
	if text, ok := the_note["text"].(string); ok {
		return text
	}
	panic("should an textfield")
}
func (the_note note) picture() string {
	if picture, ok := the_note["picture"].(string); ok {
		return picture
	}
	panic("should an picturefield")
}
func (the_note note) comments() sb {
	if comments, ok := the_note["comments"].(sb); ok {
		return comments
	}
	data := ss{
		"table":       "notes",
		"key_field_1": "predecessor_id",
		"key_value_1": the_note.id(),
	}
	comments_ids, err := send_query_to_database_and_parse_result[ss, sb](data_url.get_ids, data)
	m.Must(err, "failed to get note's comments from database", fmt.Sprint(comments_ids))
	m.Must(error_from_json(comments_ids))
	the_note["comments"] = comments_ids
	return comments_ids
}
func (the_user user) create_group_post(a_group group, text, picture string) note {
	data := sa{
		"table": "notes",
		"record": ss{
			"text":      text,
			"picture":   picture,
			"author_id": the_user.id(),
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new post in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))
	note_id := returned_data["id"].(string)
	// create adddressing
	data = sa{
		"table": "group_addressings",
		"record": ss{
			"note_id":  note_id,
			"group_id": a_group.id(),
		},
	}
	response, err = m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	_, err = m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new post in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))

	// update entities
	new_note_entity := note{
		"id":       note_id,
		"author":   the_user.id(),
		"text":     text,
		"picture":  picture,
		"comments": sb{},
		"sockets":  cb{},
	}
	get_notes()[note_id] = new_note_entity
	the_user.posts()[note_id] = true
	a_group.posts()[note_id] = true
	return new_note_entity
}

func (the_user user) create_public_post(text, picture string) note {
	data := sa{
		"table": "notes",
		"record": ss{
			"text":      text,
			"picture":   picture,
			"author_id": the_user.id(),
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new post in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))
	note_id := returned_data["id"].(string)

	// create adddressing
	data = sa{
		"table": "public_addressings",
		"record": ss{
			"note_id": note_id,
		},
	}
	response, err = m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	_, err = m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new post in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))

	// update entities
	new_note_entity := note{
		"id":       note_id,
		"author":   the_user.id(),
		"text":     text,
		"picture":  picture,
		"comments": sb{},
		"sockets":  cb{},
	}
	get_notes()[note_id] = new_note_entity
	the_user.posts()[note_id] = true
	the_user.posts_public()[note_id] = true
	return new_note_entity
}

func (the_user user) create_private_post(text, picture string) note {
	data := sa{
		"table": "notes",
		"record": ss{
			"text":      text,
			"picture":   picture,
			"author_id": the_user.id(),
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new post in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))
	note_id := returned_data["id"].(string)

	// create adddressing
	data = sa{
		"table": "private_addressings",
		"record": ss{
			"note_id": note_id,
		},
	}
	response, err = m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	_, err = m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new post in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))

	// update entities
	new_note_entity := note{
		"id":       note_id,
		"author":   the_user.id(),
		"text":     text,
		"picture":  picture,
		"comments": sb{},
		"sockets":  cb{},
	}
	get_notes()[note_id] = new_note_entity
	the_user.posts()[note_id] = true
	the_user.posts_private()[note_id] = true
	return new_note_entity
}

func (the_user user) create_exclusive_post(text, picture string) note {
	data := sa{
		"table": "notes",
		"record": ss{
			"text":      text,
			"picture":   picture,
			"author_id": the_user.id(),
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new post in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))
	note_id := returned_data["id"].(string)

	// create adddressing
	data = sa{
		"table": "exclusive_addressings",
		"record": ss{
			"note_id":      note_id,
			"addressee_id": "0",
			// an addressing with zero allow note to be referenced, and found without attributing a real user
		},
	}
	response, err = m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	_, err = m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new post in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))

	// update entities
	new_note_entity := note{
		"id":       note_id,
		"author":   the_user.id(),
		"text":     text,
		"picture":  picture,
		"comments": sb{},
		"sockets":  cb{},
	}
	get_notes()[note_id] = new_note_entity
	the_user.posts()[note_id] = true
	the_user.posts_exclusive()[note_id] = sb{}
	return get_note(note_id)
}

func (the_author user) confide(a_note note, confidant user) user {
	if !confidant.is_follower_of(the_author) {
		panic("confidant must be follower")
	}
	data := sa{
		"table": "exclusive_addressings",
		"record": ss{
			"note_id":      a_note.id(),
			"addressee_id": confidant.id(),
		},
	}
	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new post in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))

	the_author.posts_exclusive()[a_note.id()][confidant.id()] = true
	return the_author
}
func (the_author user) distrust(a_note note, confidant user) user {
	if !confidant.is_follower_of(the_author) {
		panic("confidant must be follower")
	}
	// fetch addressing id(s)
	data := ss{
		"table":       "exclusive_addressings",
		"key_field_1": "addressee_id",
		"key_value_1": confidant.id(),
		"key_field_2": "note_id",
		"key_value_2": a_note.id(),
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.get_ids, data)
	})
	m.Must(err, "failed to post json data to database server")
	returned_data, err := m.If_nil_do[sb](err, func() (sb, error) {
		return data_from_response_json[sb](response)
	})
	m.Must(err, "failed to insert new addressing  in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))
	ids := returned_data
	if len(ids) != 1 {
		fmt.Println(ids)
		panic("less (or more) than one addressig fetched")
	}
	var id_to_remove string
	for id, _ := range ids {
		id_to_remove = id
	}
	// remove record
	data = ss{
		"table": "exclusive_addressings",
		"id":    id_to_remove,
	}

	response, err = m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.delete_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	m.If_nil_must[any](err, func() (any, error) {
		return data_from_response_json[sb](response)
	})
	// the_author.posts_exclusive()[a_note.id()][confidant.id()] = true
	delete(the_author.posts_exclusive()[a_note.id()], confidant.id())
	return the_author
}
func (the_user user) comment(commented_note note, text, picture string) note {
	data := sa{
		"table": "notes",
		"record": ss{
			"text":           text,
			"picture":        picture,
			"author_id":      the_user.id(),
			"predecessor_id": commented_note.id(),
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.Must(err, "failed to insert new comment in database", fmt.Sprint(returned_data))
	m.Must(error_from_json(returned_data))
	note_id := returned_data["id"].(string)

	new_note_entity := note{
		"id":       note_id,
		"author":   the_user.id(),
		"text":     text,
		"comments": sb{},
		"picture":  picture,
		"sockets":  cb{},
	}
	get_notes()[note_id] = new_note_entity
	the_user.posts()[note_id] = true
	commented_note.comments()[note_id] = true
	return new_note_entity
}

// -------------------- USER POSTS
func (the_user user) posts() sb {
	if posts, ok := the_user["posts"].(sb); ok {
		return posts
	}
	data := ss{
		"table":       "notes",
		"key_field_1": "author_id",
		"key_value_1": the_user.id(),
	}
	notes_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_records, data)
	m.Must(err, "failed to get note's followers from database", fmt.Sprint(notes_records))
	posts_ids := sb{}
	for record_id, record := range notes_records {
		if record["author_id"] != the_user.id() {
			continue
		}
		if record["predecessor"] == "1" {
			continue
		}
		posts_ids[record_id] = true
	}
	the_user["posts"] = posts_ids
	return the_user.posts()
}
func (the_user user) posts_public() sb {
	if posts, ok := the_user["posts_public"].(sb); ok {
		return posts
	}
	data := ss{
		"table": "public_addressings",
	}
	addressings_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_file, data)
	m.Must(err, "failed to get note's  from database", fmt.Sprint(addressings_records))
	m.Must(error_from_json(addressings_records), fmt.Sprint(addressings_records))
	public_posts_ids := sb{}
	for _, record := range addressings_records {
		if the_user.posts()[record["note_id"]] {
			public_posts_ids[record["note_id"]] = true
		}
	}
	the_user["posts_public"] = public_posts_ids
	return the_user.posts_public()
}

func (the_user user) posts_private() sb {
	if posts, ok := the_user["posts_private"].(sb); ok {
		return posts
	}
	data := ss{
		"table": "private_addressings",
	}
	addressings_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_file, data)
	m.Must(err, "failed to get note's  from database", fmt.Sprint(addressings_records))
	m.Must(error_from_json(addressings_records), fmt.Sprint(addressings_records))
	private_posts_ids := sb{}
	for _, record := range addressings_records {
		if the_user.posts()[record["note_id"]] {
			private_posts_ids[record["note_id"]] = true
		}
	}
	the_user["posts_private"] = private_posts_ids
	return the_user.posts_private()
}

func (the_user user) posts_exclusive() ssb {
	if posts_addressees, ok := the_user["posts_exclusive"].(ssb); ok {
		return posts_addressees
	}
	data := ss{
		"table": "exclusive_addressings",
	}
	addressings_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_file, data)
	m.Must(err, "failed to get note's  from database", fmt.Sprint(addressings_records))
	m.Must(error_from_json(addressings_records), fmt.Sprint(addressings_records))
	fmt.Println("\nGETTING ALL EXCLUSIVE POST FROM USER")
	fmt.Println("user_entity", the_user.id())
	posts_ids_addressees_ids := ssb{}
	for _, record := range addressings_records {
		fmt.Println("an adderssing", record)
		a_note := get_note(record["note_id"])
		fmt.Println("notes's author", a_note.author()) //, get_user(a_note.author()))
		if a_note.author() != the_user.id() {
			continue
		}
		addressee_id := record["addressee_id"]
		// post belong to user so thsi i at least an zero x addressing
		// we repertory the note's id in user's exclu notes set
		if _, ok := posts_ids_addressees_ids[a_note.id()]; !ok {
			posts_ids_addressees_ids[a_note.id()] = sb{}
		}
		if addressee_id == "0" {
			continue
			// means there is no confidant associated to it
		}
		fmt.Println("it s a match")
		a_confidant := get_note(addressee_id)
		posts_ids_addressees_ids[a_note.id()][a_confidant.id()] = true
		fmt.Println("found something", posts_ids_addressees_ids)
	}
	the_user["posts_exclusive"] = posts_ids_addressees_ids
	return the_user.posts_exclusive()
}

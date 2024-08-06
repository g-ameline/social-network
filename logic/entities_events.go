package main

import (
	"fmt"

	m "github.com/g-ameline/maybe"
	// ws "github.com/gorilla/websocket"
	http "net/http"
)

// ------------- EVENTS

func get_events() sa {
	if events, ok := entities["events"].(sa); ok {
		return events
	}
	data := ss{"table": "events"}
	ids, err := send_query_to_database_and_parse_result[ss, sb](data_url.get_ids, data)
	m.Must(err, "failed to get events from database", fmt.Sprint(ids))
	m.Must(error_from_json(ids))
	entities["events"] = sa{}
	for id := range ids {
		entities["events"].(sa)[id] = id
	}
	return get_events()
}
func get_event(event_id string) event {
	if an_event, ok := get_events()[event_id]; ok {
		if event_entity, ok := an_event.(event); ok {
			return event_entity
		}
	}
	data := ss{
		"table":       "events",
		"key_field_1": "id",
		"key_value_1": event_id,
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.get_record, data)
	})
	m.Must(err, "failed to post data to database server")
	record, err := m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to get missing event's record from database", event_id)
	m.Must(error_from_json(record))
	new_event_entity := event{
		"id":          event_id,
		"group":       record["group_id"],
		"title":       record["title"],
		"description": record["description"],
		"date":        record["date"],
		"sockets":     cb{},
	}
	entities["events"].(sa)[event_id] = new_event_entity
	return get_event(event_id)
}

func (the_event event) group() string {
	if group_id, ok := the_event["group"].(string); ok {
		return group_id
	}
	panic("event should be there")
}
func (the_event event) title() string {
	if title, ok := the_event["title"].(string); ok {
		return title
	}
	panic("event should be there")
}
func (the_event event) date() string {
	if date, ok := the_event["date"].(string); ok {
		return date
	}
	panic("event should be there")
}
func (the_event event) description() string {
	if description, ok := the_event["description"].(string); ok {
		return description
	}
	panic("event should be there")
}

func (the_event event) attenders() sb {
	if attenders, ok := the_event["attenders"].(sb); ok {
		return attenders
	}
	data := ss{"table": "attendings"}
	attenders_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_file, data)
	m.Must(err, "failed to get note's  from database", fmt.Sprint(attenders_records))
	the_event["attenders"] = sb{}
	for _, record := range attenders_records {
		if record["coming"] == "1" && record["event_id"] == the_event.id() {
			the_event.attenders()[record["attender_id"]] = true
		}
	}
	return the_event["attenders"].(sb)
}

func (the_event event) absentees() sb {
	if absentees, ok := the_event["absentees"].(sb); ok {
		return absentees
	}
	data := ss{"table": "attendings"}
	absentees_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_file, data)
	m.Must(err, "failed to get note's  from database", fmt.Sprint(absentees_records))
	the_event["absentees"] = sb{}
	for _, record := range absentees_records {
		if record["coming"] == "-1" && record["event_id"] == the_event.id() {
			the_event.absentees()[record["attender_id"]] = true
		}
	}
	return the_event["absentees"].(sb)
}

func (the_group group) create_event(title, description, date string) event {
	data := sa{
		"table": "events",
		"record": ss{
			"title":       title,
			"description": description,
			"date":        date,
			"group_id":    the_group.id(),
		},
	}
	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	_, err = m.If_nil_do[nvm](err, func() error {
		error_db, ok := returned_data["error"]
		if ok {
			panic(error_db)
			// return fmt.Errorf(error_db.(string))
		}
		return error(nil)
	})
	event_id, err := m.If_nil_do[string](err, func() string {
		return returned_data["id"].(string)
	})

	new_event_entity := event{
		"id":          event_id,
		"group":       the_group.id(),
		"title":       title,
		"description": description,
		"date":        date,
		"sockets":     cb{},
	}
	get_events()[event_id] = new_event_entity
	the_group["events"] = event_id
	return new_event_entity
}
func (the_user user) attend(the_event event) {
	if the_event.attenders()[the_user.id()] {
		panic("already attending")
	}
	data := sa{
		"table": "attendings",
		"record": ss{
			"coming":      "1",
			"attender_id": the_user.id(),
			"event_id":    the_event.id(),
		},
	}
	upserted_id, err := send_query_to_database_and_parse_result[sa, sss](data_url.upsert_record, data)
	m.Must(err, "failed to upsert attending ", fmt.Sprint(upserted_id))
	delete(the_event.absentees(), the_user.id())
	the_event.attenders()[the_user.id()] = true
}
func (the_user user) absent(the_event event) {
	if the_event.attenders()[the_user.id()] {
		panic("already unattending")
	}
	data := sa{
		"table": "attendings",
		"record": ss{
			"coming":      "-1",
			"attender_id": the_user.id(),
			"event_id":    the_event.id(),
		},
	}
	upserted_id, err := send_query_to_database_and_parse_result[sa, sss](data_url.upsert_record, data)
	m.Must(err, "failed to upsert attending ", fmt.Sprint(upserted_id))
	delete(the_event.attenders(), the_user.id())
	the_event.absentees()[the_user.id()] = true
}

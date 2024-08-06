package main

import (
	"fmt"

	m "github.com/g-ameline/maybe"
	// ws "github.com/gorilla/websocket"
	http "net/http"
)

// ------------- GROUPS

func get_groups() sa {
	if groups, ok := entities["groups"]; ok {
		return groups.(sa)
	}
	data := ss{"table": "groups"}
	ids, err := send_query_to_database_and_parse_result[ss, sb](data_url.get_ids, data)
	m.Must(err, "failed to get groups records from database")
	entities["groups"] = sa{}
	for id := range ids {
		if _, ok := entities["groups"].(sa)[id]; !ok {
			entities["groups"].(sa)[id] = true
		}
	}
	return get_groups()
}

func get_group(group_id string) group {
	if a_group, ok := get_groups()[group_id]; ok {
		if group_entity, ok := a_group.(group); ok {
			return group_entity
		}
	}

	data := ss{
		"table":       "groups",
		"key_field_1": "id",
		"key_value_1": group_id,
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.get_record, data)
	})
	m.Must(err, "failed to post data to database server")
	record, err := m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to get missing group's record from database", group_id)
	m.Must(error_from_json(record), "error with data")
	new_group_entity := group{
		"id":          group_id,
		"title":       record["title"],
		"description": record["description"],
		"creator":     record["creator_id"],
		"sockets":     cb{},
	}
	entities["groups"].(sa)[group_id] = new_group_entity
	entities["groups"].(sa)[group_id].(group).title()

	return get_group(group_id)
}
func get_global_groups() group {
	return entities["global"].(sa)["groups"].(group)
}
func (the_group group) membership(user user) string {
	return membership(user, the_group)
}
func (the_group group) memberships() ssb {
	data := ss{
		"table":       "joinings",
		"key_field_1": "group_id",
		"key_value_1": the_group.id(),
	}
	joinings_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_records, data)
	m.Must(err, "failed to get joinings from database", fmt.Sprint(data))
	members := sb{}
	invitees := sb{}
	applicants := sb{}
	for _, joining_record := range joinings_records {
		if joining_record["group_id"] != the_group.id() {
			continue
		}
		joiner_id := joining_record["joiner_id"]
		appr_creator := joining_record["approval_creator"]
		appr_joiner := joining_record["approval_joiner"]

		if appr_creator == "1" && appr_joiner == "1" {
			members[joiner_id] = true
			continue
		}
		if appr_creator == "1" && appr_joiner == "0" {
			invitees[joiner_id] = true
			continue
		}
		if appr_creator == "0" && appr_joiner == "1" {
			applicants[joiner_id] = true
			continue
		}
	}
	memberships := ssb{}
	the_group["members"] = members
	memberships["members"] = members
	the_group["invitees"] = invitees
	memberships["invitees"] = invitees
	the_group["applicants"] = applicants
	memberships["applicants"] = applicants
	return memberships
}

func (the_group group) members() sb {
	members, ok := the_group["members"].(sb)
	if !ok {
		return the_group.memberships()["members"]
	}
	return members
}
func (the_group group) invitees() sb {
	invitees, ok := the_group["invitees"].(sb)
	if !ok {
		return the_group.memberships()["invitees"]
	}
	return invitees
}
func (the_group group) applicants() sb {
	applicants, ok := the_group["applicants"].(sb)
	if !ok {
		return the_group.memberships()["applicants"]
	}
	return applicants
}

func (the_group group) outsiders() sa {
	outsiders := get_users()
	delete(outsiders, the_group.creator())
	for member_id := range the_group.members() {
		delete(outsiders, member_id)
	}
	for member_id := range the_group.invitees() {
		delete(outsiders, member_id)
	}
	for member_id := range the_group.applicants() {
		delete(outsiders, member_id)
	}
	return outsiders
}

func (the_group group) invite(invitee user) group {
	if the_group.creator() == invitee.id() {
		panic("cant invite itself in a group")
	}
	data := sa{
		"table": "joinings",
		"record": ss{
			"joiner_id":        invitee.id(),
			"group_id":         the_group.id(),
			"approval_creator": "1",
			"approval_joiner":  "0",
		},
	}
	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	_, err = m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to parse joining data from database", fmt.Sprint(data))

	the_group.invitees()[invitee.id()] = true
	return the_group
}
func (the_group group) reject(applicant user) group {
	data := sa{
		"table": "joinings",
		"record": ss{
			"joiner_id":        applicant.id(),
			"creator_id":       the_group.creator(),
			"group_id":         the_group.id(),
			"approval_creator": "1",
			"approval_joiner":  "0",
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	_, err = m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to parse joining data from database", fmt.Sprint(data))

	// _, err := send_query_to_database_and_parse_result[sa, string](data_url.upsert_record, data)
	// m.Must(err, "failed to upsert joining in database", fmt.Sprint(data))

	delete(the_group.applicants(), applicant.id())
	return the_group
}
func (the_group group) admit(applicant user) group {
	data := sa{
		"table": "joinings",
		"record": ss{
			"joiner_id":        applicant.id(),
			"group_id":         the_group.id(),
			"approval_creator": "1",
			"approval_joiner":  "1",
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	_, err = m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to parse joining data from database", fmt.Sprint(data))

	delete(the_group.applicants(), applicant.id())
	delete(applicant.applied_groups(), the_group.id())
	the_group.members()[applicant.id()] = true
	return the_group
}

func (the_creator user) create_group(title, description string) group {
	data := sa{
		"table": "groups",
		"record": ss{
			"title":       title,
			"description": description,
			"creator_id":  the_creator.id(),
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data)
	})
	m.Must(err, "failed to post json data to database server")
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	m.If_nil_must[nvm](err, func() error {
		error_db, ok := returned_data["error"]
		if ok {
			return fmt.Errorf(error_db.(string))
		}
		return error(nil)
	})
	group_id := returned_data["id"].(string)

	new_group_entity := group{
		"id":          group_id,
		"title":       title,
		"description": description,
		"creator":     the_creator.id(),
		"sockets":     cb{},
	}
	get_groups()[group_id] = new_group_entity
	the_creator.created_groups()[group_id] = true

	return get_group(group_id)
}

// ------------ GET DATA from group
func (the_group group) creator() string {
	if creator_id, ok := the_group["creator"].(string); ok {
		return creator_id
	}
	panic("group created without creator id" + fmt.Sprint(the_group))
}

func (the_group group) title() string {
	if title, ok := the_group["title"]; ok {
		if title == "" {
			panic("no title")
		}
		return title.(string)
	}
	panic("should not need to fetch title from group")
}

func (the_group group) description() string {
	if description, ok := the_group["description"].(string); ok {
		return description
	}
	panic("should not need to fetch description from group")
}
func (the_group group) events() sb {
	if events, ok := the_group["events"].(sb); ok {
		return events
	}
	// data := ss{"table": "events"}
	data := ss{
		"table":       "events",
		"key_field_1": "group_id",
		"key_value_1": the_group.id(),
	}
	events_ids, err := send_query_to_database_and_parse_result[ss, sb](data_url.get_ids, data)
	m.Must(err, "failed to get note's followers from database", fmt.Sprint(data))
	// events_ids := sb{}
	// for record_id, record := range events_records {
	// 	if record["group_id"] == the_group.id() {
	// 		events_ids[record_id] = true
	// 	}
	// }
	the_group["events"] = events_ids
	return events_ids
}

func (the_group group) posts() sb {
	if posts, ok := the_group["posts"].(sb); ok {
		return posts
	}
	data := ss{
		"table":       "group_addressings",
		"key_field_1": "group_id",
		"key_value_1": the_group.id(),
	}
	addressings_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_records, data)
	// posts_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_file, data)
	m.Must(err, "failed to get note's  from database", fmt.Sprint(addressings_records))
	m.Must(error_from_json(addressings_records), fmt.Sprint(addressings_records))
	posts_ids := sb{}
	for _, record := range addressings_records {
		if record["group_id"] != the_group.id() {
			continue
		}
		posts_ids[record["note_id"]] = true
	}
	the_group["posts"] = posts_ids
	return posts_ids
}

// ------------- LIST
func (the_user user) joined_groups() sb {
	if joined_groups, ok := the_user["joined_groups"].(sb); ok {
		return joined_groups
	}
	the_user.memberships()
	return the_user.joined_groups()
}

func (the_user user) invited_groups() sb {
	if invited_groups, ok := the_user["invited_groups"].(sb); ok {
		return invited_groups
	}
	the_user.memberships()
	return the_user.invited_groups()
}

func (the_user user) applied_groups() sb {
	if applied_groups, ok := the_user["applied_groups"].(sb); ok {
		return applied_groups
	}
	the_user.memberships()
	return the_user.applied_groups()
}

func (the_user user) applicable_groups() sb {
	if applicable_groups, ok := the_user["applicable_groups"].(sb); ok {
		return applicable_groups
	}
	the_user["applicable_groups"] = sb{}
	for group_id := range get_groups() {
		if the_user.created_groups()[group_id] {
			continue
		}
		if the_user.joined_groups()[group_id] {
			continue
		}
		if the_user.invited_groups()[group_id] {
			continue
		}
		if the_user.applied_groups()[group_id] {
			continue
		}
		the_user["applicable_groups"].(sb)[group_id] = true
	}
	return the_user.applicable_groups()
}
func (the_user user) decline(a_group group) user {
	data := sa{
		"table": "joinings",
		"record": ss{
			"joiner_id":        the_user.id(),
			"creator_id":       a_group.creator(),
			"group_id":         a_group.id(),
			"approval_creator": "1",
			"approval_joiner":  "0",
		},
	}
	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	_, err = m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to parse joining data from database", fmt.Sprint(data))
	a_group.members()[the_user.id()] = true
	delete(a_group.invitees(), the_user.id())
	return the_user
}
func (the_user user) apply(a_group group) user {
	data := sa{
		"table": "joinings",
		"record": ss{
			"joiner_id":        the_user.id(),
			"group_id":         a_group.id(),
			"approval_creator": "0",
			"approval_joiner":  "1",
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	_, err = m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to upsert joining in database", fmt.Sprint(data))

	a_group.applicants()[the_user.id()] = true

	return the_user
}

func (the_user user) unapply(a_group group) user {
	data := sa{
		"table": "joinings",
		"record": ss{
			"joiner_id":        the_user.id(),
			"group_id":         a_group.id(),
			"approval_creator": "0",
			"approval_joiner":  "0",
		},
	}

	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	_, err = m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to upsert joining in database", fmt.Sprint(data))

	m.If_nil_must[ss](err, func() {
		delete(a_group.applicants(), the_user.id())
	})
	return the_user
}

func (the_user user) assent(a_group group) user {
	data := sa{
		"table": "joinings",
		"record": ss{
			"joiner_id":        the_user.id(),
			"creator_id":       a_group.creator(),
			"group_id":         a_group.id(),
			"approval_creator": "1",
			"approval_joiner":  "1",
		},
	}
	response, err := m.If_nil_do[*http.Response](nil, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.upsert_record, data)
	})
	m.Must(err, "failed to post data to database server")
	_, err = m.If_nil_do[ss](err, func() (ss, error) {
		return data_from_response_json[ss](response)
	})
	m.Must(err, "failed to parse joining data from database", fmt.Sprint(data))

	a_group.members()[the_user.id()] = true
	delete(a_group.invitees(), the_user.id())
	delete(the_user.invited_groups(), a_group.id())
	return the_user
}

// ------------- GET DATA from the_user
func (the_user user) created_groups() sb {
	if created_groups, ok := the_user["created_groups"].(sb); ok {
		return created_groups
	}
	data := ss{
		"table":       "groups",
		"key_field_1": "creator_id",
		"key_value_1": the_user.id(),
	}
	created_ids, err := send_query_to_database_and_parse_result[ss, sb](data_url.get_ids, data)
	m.Must(err, "failed to fetch groups ids created from database", fmt.Sprint(data))
	m.Must(error_from_json(created_ids))
	the_user["created_groups"] = created_ids
	return the_user.created_groups()
}
func (the_user user) membership(group_entity group) string {
	return membership(the_user, group_entity)
}

func (the_user user) memberships() user {
	data := ss{"table": "joinings"}
	joinings_records, err := send_query_to_database_and_parse_result[ss, sss](data_url.get_file, data)
	m.Must(err, "failed to get joinings from database", fmt.Sprint(joinings_records))
	joined_ids := sb{}
	invited_ids := sb{}
	applied_ids := sb{}
	for _, joining_record := range joinings_records {
		joiner_id := joining_record["joiner_id"]
		if joiner_id != the_user.id() {
			continue
		}
		group_id := joining_record["group_id"]
		approval_creator := joining_record["approval_creator"]
		approval_joiner := joining_record["approval_joiner"]
		if approval_creator == "1" && approval_joiner == "1" {
			joined_ids[group_id] = true
			continue
		}
		if approval_creator == "1" && approval_joiner == "0" {
			invited_ids[group_id] = true
			continue
		}
		if approval_creator == "0" && approval_joiner == "1" {
			applied_ids[group_id] = true
			continue
		}
		if approval_creator == "-1" || approval_joiner == "-1" {
			continue
		}
	}
	the_user["joined_groups"] = joined_ids
	the_user["invited_groups"] = invited_ids
	the_user["applied_groups"] = applied_ids
	return the_user
}

func membership(user_entity user, group_entity group) string {
	switch true {
	case group_entity.creator() == user_entity.id():
		return "created"
	case group_entity.members()[user_entity.id()]:
		return "joined"
	case group_entity.invitees()[user_entity.id()]:
		return "invited"
	case group_entity.applicants()[user_entity.id()]:
		return "applied"
	default:
		return "applicable"
	}
}

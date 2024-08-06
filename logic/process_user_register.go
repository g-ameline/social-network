package main

import (
	"fmt"
	uuid "github.com/google/uuid"
	"net/http"

	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
)

func process_user_register(data sa) (string, error) {

	c.Magenta("\n handling register demand")
	c.Magenta("DATA FOR REGISTER")
	c.Print_map(data)

	credentials_data, err := register_data(data, files_folder)
	data_for_database := data_for_db_insertion("users", credentials_data)
	c.Yellow("data for db")
	c.Print_map(data_for_database)
	c.Yellow("path to send", data_url.insert_record)
	response, err := m.If_nil_do[*http.Response](err, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.insert_record, data_for_database)
	})
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})
	c.Yellow("returned data from db")
	c.Print_map(returned_data)
	_, err = m.If_nil_do[nvm](err, func() error {
		error_db, ok := returned_data["error"].(string)
		return m.Ok_to_err(ok, error_db)
	})
	_, err = m.If_nil_do[nvm](err, func() {
		c.Magenta("data from db")
		c.Print_map(returned_data)
		data["session"] = uuid.NewString()
		user_id := returned_data["id"].(string)
		data["user_id"] = user_id
		refresh_subscribers(get_user(user_id))
		refresh_subscribers(get_global_users())
	})

	var reply string
	m.If_nil_do[nvm](err,
		func() { reply = main_frame_logged_node().Inline() })

	// if not ok rerender credential stuff with same fields
	m.If_error_do[nvm](err, func() {
		data["error"] = err.Error()
		node := main_frame_credentials_node(ss{}, register_form_data(data))
		reply = node.Inline()
	})
	return reply, err
}

func register_data(data sa, picture_folder string) (row ss, err error) {
	row = ss{}
	var ok bool
	if row["email"], ok = data["email"].(string); !ok {
		return row, fmt.Errorf("need email")
	}
	if row["password"], ok = data["password"].(string); !ok {
		return row, fmt.Errorf("need password")
	}
	row["password"] = Hash_that(data["password"].(string))
	row["first_name"], ok = data["first_name"].(string)
	row["last_name"], ok = data["last_name"].(string)
	row["birth"], ok = data["birth"].(string)
	row["nickname"], ok = data["nickname"].(string)
	row["about"], ok = data["about"].(string)
	row["avatar"], ok = data["filepath"].(string)
	return row, err
}

func register_form_data(data sa) ss {
	register_values := ss{}
	register_values["email"] = data["email"].(string)
	register_values["password"] = data["password"].(string)
	register_values["first_name"] = data["first_name"].(string)
	register_values["last_name"] = data["last_name"].(string)
	register_values["birth"] = data["birth"].(string)
	register_values["nickname"] = data["nickname"].(string)
	register_values["about"] = data["about"].(string)
	if filepath, ok := data["avatar"]; ok {
		register_values["avatar"] = filepath.(string)
	}
	register_values["error"] = data["error"].(string)
	return register_values
}

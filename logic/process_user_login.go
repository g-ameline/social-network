package main

import (
	"fmt"
	"net/http"

	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	uuid "github.com/google/uuid"
)

func process_user_login(data sa) (string, error) {
	c.Magenta("\n handling login demand")
	c.Magenta("DATA FOR LOGIN")
	c.Print_map(data)
	// check if any credential with request (cookies)
	credentials_data, err := login_data(data)
	fmt.Println(credentials_data)

	response, err := m.If_nil_do[*http.Response](err, func() (*http.Response, error) {
		return post_json_data_to_url(data_url.user_login, credentials_data)
	})
	returned_data, err := m.If_nil_do[sa](err, func() (sa, error) {
		return data_from_response_json[sa](response)
	})

	// if data based marked data as erronous update error accordingly
	_, err = m.If_nil_do[nvm](err, func() error {
		error_db, ok := returned_data["error"].(string)
		return m.Ok_to_err(ok, error_db)
	})
	_, err = m.If_nil_do[nvm](err, func() {
		data["session"] = uuid.NewString()
		data["user_id"] = returned_data["user_id"].(string)
	})
	var reply string
	m.If_nil_must[nvm](err,
		func() { reply = main_frame_logged_node().Inline() })

	m.If_error_must[nvm](err,
		func() {
			data["error"] = err.Error()
			node := main_frame_credentials_node(login_form_data(data), ss{})
			reply = node.Inline()
		})
	return reply, err
}
func login_data(data sa) (ss, error) {
	credential := ss{}
	err := error(nil)
	var ok bool
	var email, password string
	if email, ok = data["email"].(string); !ok {
		email = ""
		err = fmt.Errorf("email not supplied")
	}
	if password, ok = data["password"].(string); !ok {
		password = ""
		err = fmt.Errorf("password not supplied")
	}
	credential["email"] = email
	credential["password"] = Hash_that(password)
	return credential, err
}

func login_form_data(data sa) ss {
	login_values := ss{}
	login_values["email"] = data["email"].(string)
	login_values["password"] = data["password"].(string)
	if error, ok := data["error"]; ok {
		login_values["error"] = error.(string)
	}
	return login_values
}

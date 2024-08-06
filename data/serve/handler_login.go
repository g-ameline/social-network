package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
)

func handle_login(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		println("\nHANDLING LOGIN")
		credential, err := get_data_from_request_json[ss](request)
		fmt.Println(credential, err)
		respond_error(responder, err)
		// check with only email
		user_id, err := db.Get_id_two_cond(database_path, "users", "email", credential["email"], "password", credential["password"])
		if err != nil {
			fmt.Println("sending error not find")
			respond_error(responder, fmt.Errorf("could not find user matching credential"))
		}
		response_data := map[string]string{}
		response_data["user_id"] = user_id
		err = api.Respond_json_data(responder, response_data)
		crash(err, "failed to respond")
	}
}
func get_data_from_request_json[data_type any](request *http.Request) (data_type, error) {
	decoder := json.NewDecoder(request.Body)
	var data data_type
	err := decoder.Decode(&data)
	return data, err
}

package main

import (
	"net/http"

	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
)

func handle_registration(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		println("\nHANDLING REGISTRATION")
		data, err := api.Get_data_from_request_json[ss](request)
		_, ok := data["private"]
		if !ok {
			data["private"] = "0"
		}
		respond_error(responder, err, "failed to read request")
		user_id, err := db.Insert_one_row(database_path, "users", data)
		breadcrumb("new user insertion output :", user_id, err)
		respond_error(responder, err)
		response_data := map[string]string{}
		response_data["user_id"] = user_id
		err = api.Respond_json_data(responder, response_data)
		crash(err, "failed to respond")
	}
}

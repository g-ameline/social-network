package main

import (
	"net/http"

	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
)

func handle_user(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		println("\nHANDLING USER")
		data, err := api.Get_data_from_request_json[ss](request)
		user_id := data["user_id"]
		respond_error(responder, err)
		// check with only email
		user_row, err := db.Get_row_one_cond(database_path, "users", "id", user_id)
		user_data := ss{}
		user_data["user_id"] = user_id
		user_data["nickname"] = user_row["nickname"]
		user_data["image"] = user_row["avatar"]
		err = api.Respond_json_data(responder, user_data)
		crash(err, "failed to respond")
	}
}

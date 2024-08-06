package main

import (
	"fmt"
	"net/http"

	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
)

func handle_delete_record(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		println("\nHANDLING DELETION")
		// we only need table name and record's id
		data, err := api.Get_data_from_request_json[ss](request)
		respond_error(responder, err)
		err = db.Delete_rows(database_path, data["table"], "id", data["id"])
		fmt.Println("error", err)
		if err != nil {
			fmt.Println("sending error not find")
			respond_error(responder, fmt.Errorf("could not find user matching credential"))
		}
		response_data := map[string]string{}
		response_data["deleted_id"] = data["id"]
		response_data["table"] = data["table"]
		err = api.Respond_json_data(responder, response_data)
		crash(err, "failed to respond")
	}
}

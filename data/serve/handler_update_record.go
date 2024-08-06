package main

import (
	"fmt"
	"net/http"

	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
)

func handle_update_record(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		data, err := api.Get_data_from_request_json[ss](request)
		fmt.Println("updating a record", data)
		table_name_raw, ok := data["table"]
		if !ok {
			respond_error(responder, fmt.Errorf("need table name key value in data"))
			return
		}
		table_name := table_name_raw
		id, ok := data["id"]
		if !ok {
			respond_error(responder, fmt.Errorf("need id in data"))
			return
		}
		key_field, ok := data["field"]
		if !ok {
			respond_error(responder, fmt.Errorf("need col/key_field in data"))
			return
		}
		new_key_value, ok := data["value"]
		if !ok {
			respond_error(responder, fmt.Errorf("need new_key_value in data"))
			return
		}
		err = db.Update_value(database_path, table_name, id, key_field, new_key_value)
		respond_error(responder, err)
		data_to_send := map[string]any{}
		data_to_send["id"] = id
		fmt.Println("updated a record", data_to_send)
		err = api.Respond_with_json_data(responder, data_to_send)
		crash(err, "failed to respond")
	}
}

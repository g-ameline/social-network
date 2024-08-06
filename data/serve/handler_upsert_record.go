package main

import (
	"fmt"
	"net/http"

	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
)

func handle_upsert_record(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		println("\nrequest for upsert record")
		data, err := get_data_from_request_json[sa](request)
		fmt.Println(data, err)
		table_name_raw, ok := data["table"]
		if !ok {
			respond_error(responder, fmt.Errorf("need table name key value in data"))
			return
		}
		table_name := table_name_raw.(string)

		record_data_raw, ok := data["record"]
		if !ok {
			respond_error(responder, fmt.Errorf("need record dat (column => value) in data"))
			return
		}
		record_data := any_to_ss(record_data_raw)

		fmt.Println("data to add to DB")
		print_dat_map(data["record"].(sa))

		err = db.Upsert_row(database_path, table_name, record_data)
		fmt.Println("error", err)
		respond_error(responder, err)
		data_to_send := map[string]any{}
		err = api.Respond_with_json_data(responder, data_to_send)
		crash(err, "failed to respond")
	}
}

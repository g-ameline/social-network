package main

import (
	"fmt"
	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
	"net/http"
)

func handle_insert_record(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		fmt.Println("\nrequest for inserting record")
		fmt.Println(request.Method)
		fmt.Println(request.Body)
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
		println("table", data["table"].(string))
		print_dat_map(data["record"].(sa))

		record_id, err := db.Insert_one_row(database_path, table_name, record_data)
		respond_error(responder, err)
		data_to_send := map[string]any{}
		data_to_send["id"] = record_id
		fmt.Println("inserted id", record_id)
		err = api.Respond_with_json_data(responder, data_to_send)
		crash(err, "failed to respond")
	}
}
func any_to_ss(a any) ss {
	return sa_to_ss(a.(sa))
}
func sa_to_ss(sa sa) ss {
	new_ss := ss{}
	for k, v := range sa {
		new_ss[k] = v.(string)
	}
	return new_ss
}

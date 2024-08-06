package main

import (
	"fmt"
	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
	"net/http"
)

func handle_get_file(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		println("\nrequest for gettting ids")
		data, err := api.Get_data_from_request_json[sa](request)
		fmt.Println("data for fetching ids")
		print_dat_map(data)
		table_name, ok := data["table"]
		if !ok {
			respond_error(responder, err, "need table key-value")
		}
		rows, err := db.Get_rows(database_path, table_name.(string))
		if err != nil {
			respond_error(responder, err)
		}
		if err == nil {
			fmt.Println("responding with those rows", len(rows))
			err = api.Respond_with_json_data(responder, rows)
		}
		crash(err, "failed to respond")
	}
}

package main

import (
	"fmt"
	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
	"net/http"
)

func handle_get_records(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		println("\nrequest for gettting records")
		data, err := api.Get_data_from_request_json[sa](request)
		fmt.Println("data for fetching all records")
		print_dat_map(data)
		table_name_raw, ok := data["table"]
		if !ok {
			respond_error(responder, fmt.Errorf("need table name key value in data"))
			return
		}
		table_name := table_name_raw.(string)
		key_field_1_raw, ok := data["key_field_1"]
		if !ok {
			respond_error(responder, fmt.Errorf("need record dat (column => value) in data"))
			return
		}
		key_value_1_raw, ok := data["key_value_1"]
		if !ok {
			respond_error(responder, fmt.Errorf("need record dat (column => value) in data"))
			return
		}
		key_field_1, key_value_1 := key_field_1_raw.(string), key_value_1_raw.(string)
		var only_one_condition bool = false
		key_field_2_raw, ok := data["key_field_2"]
		if !ok {
			only_one_condition = true
		}
		key_value_2_raw, ok := data["key_value_2"]
		if !ok {
			only_one_condition = true
		}
		var rows sss
		if only_one_condition == true {
			rows, err = db.Get_rows_one_cond(database_path, table_name, key_field_1, key_value_1)
		}
		if only_one_condition == false {
			key_field_2, key_value_2 := key_field_2_raw.(string), key_value_2_raw.(string)
			rows, err = db.Get_rows_two_cond(database_path, table_name, key_field_1, key_value_1, key_field_2, key_value_2)
		}
		if err != nil {
			fmt.Println("error fetching data", err)
			respond_error(responder, err)
		}
		if err == nil {
			fmt.Println("responding with those data", len(rows), rows)
			err = api.Respond_with_json_data(responder, rows)
			crash(err, "failure here")
		}
	}
}

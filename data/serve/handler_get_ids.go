package main

import (
	"fmt"
	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
	"net/http"
)

func handle_get_ids(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		println("\nrequest for gettting ids")
		data, err := api.Get_data_from_request_json[sa](request)
		fmt.Println("data for fetching ids")
		print_dat_map(data)
		var ids map[string]bool
		table_name_raw, ok := data["table"]
		if !ok {
			respond_error(responder, fmt.Errorf("need record table name in data"))
			return
		}
		table_name := table_name_raw.(string)
		key_field_1_raw, ok := data["key_field_1"]
		if !ok && len(data) == 1 {
			ids, err = db.Get_ids(database_path, table_name_raw.(string))
			fmt.Println("responding with all ids", len(ids))
			err = api.Respond_with_json_data(responder, ids)
			return
		}
		key_value_1_raw, ok := data["key_value_1"]
		if !ok {
			respond_error(responder, fmt.Errorf("need record data (column => value) in data"))
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
		if only_one_condition == true {
			ids, err = db.Get_ids_one_cond(database_path, table_name, key_field_1, key_value_1)
		}
		if only_one_condition == false {
			key_field_2, key_value_2 := key_field_2_raw.(string), key_value_2_raw.(string)
			ids, err = db.Get_ids_two_cond(database_path, table_name, key_field_1, key_value_1, key_field_2, key_value_2)
		}
		if err != nil {
			fmt.Println("error fetching data", err)
			fmt.Println("responding with ids", len(ids))
			err = api.Respond_with_json_data(responder, ss{"error": err.Error()})
		}
		if err == nil {
			fmt.Println("responding with ids", len(ids))
			err = api.Respond_with_json_data(responder, ids)
		}
		crash(err, "failed to respond")
	}
}

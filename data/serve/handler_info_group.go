package main

import (
	// "encoding/json"
	// "fmt"
	"net/http"

	api "github.com/g-ameline/api_helper"
	db "github.com/g-ameline/sql_helper"
)

func handle_info_group(database_path string) http.HandlerFunc {
	return func(responder http.ResponseWriter, request *http.Request) {
		println("HANDLE GROUP INFO")
		var data map[string]string // we need user_id only
		var err error
		data, err = api.Get_data_from_request_json[ss](request)
		user_id := data["user_id"]
		// group_id := data["group_id"]
		crash(err, "failed to get data from request")
		breadcrumb("take all groups")
		groups_rows, err := db.Get_rows(database_path, "groups")
		crash(err)
		// print_data(groups_rows)
		breadcrumb("group list length ", len(groups_rows))
		breadcrumb("take all joinings")
		joinings_rows, err := db.Get_rows(database_path, "joinings")
		crash(err)
		// print_data(joinings_rows)
		breadcrumb("keep all joinings where joiner is user")
		for _, joining_row := range joinings_rows {
			group_id := joining_row["group_id"]
			if joining_row["joiner_id"] != user_id {
				delete(groups_rows, group_id)
				continue
			}
			if joining_row["approval_joiner"] != "1" {
				delete(groups_rows, group_id)
				continue
			}
			if joining_row["approval_creator"] != "1" {
				delete(groups_rows, group_id)
				continue
			}
		}
		breadcrumb("group list length ", len(groups_rows))
		breadcrumb("gather data to send")
		err = api.Respond_with_json_data(responder, groups_rows)
		crash(err, "failed to respond")
	}
}

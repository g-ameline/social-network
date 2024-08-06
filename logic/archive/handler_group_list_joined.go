package main

import (
	c "github.com/g-ameline/colors"
)

func handler_group_list_joined(demand *Demand) (reply sa) {
	// func treat_list_asked_group(demand Demand) (reply sa) {
	Breadcrumb("error state of demand", demand.Error)
	Breadcrumb("\n handling list of joined groups demand", c.RED)
	defer c.Resetting()
	Breadcrumb("DATA FOR JOINED GROUPS")
	Print_dat_map(demand.Data)
	user_id := If_ok_do[string](demand, func() string { return demand.User_id })
	Breadcrumb("take all groups")
	data_for_db := sa{}
	data_for_db["user_id"] = user_id
	data_joined_groups := If_ok_do[sa](demand, func() (sa, error) {
		return send_query_to_database_and_get_result(demand.DB_uri, data_for_db)
	})
	If_ok_do[int](demand, func() {
		reply_ok := sa{}
		reply_ok["what"] = "group_my"
		reply_ok["info"] = data_joined_groups
		reply = reply_ok
	})
	If_nok_do[int](demand, func() {
		reply_nok := sa{}
		reply_nok["what"] = "group my failed"
		reply_nok["info"] = demand.Error.Error()
		reply = reply_nok
	})
	return reply
}

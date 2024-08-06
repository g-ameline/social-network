package main

import (
	d "app_server/dom"
	c "github.com/g-ameline/colors"
)

func handler_group_list_asked(demand *Demand) string {
	// Breadcrumb("error state of demand", demand.Error)
	// Breadcrumb("\n handling list of asked groups demand", c.GREEN)
	// defer c.Resetting()
	// Breadcrumb("DATA FOR ALL ASKED GROUPS")
	// Print_dat_map(demand.Data)
	// user_id := If_ok_do[string](demand, func() string { return demand.User_id })
	// data_user := ss{}
	// data_user["user_id"] = user_id
	// data_asked_groups := If_ok_do[sa](demand, func() (sa, error) {
	// 	return send_query_to_database_and_get_result(demand.DB_uri, data_user)
	// })
	// Breadcrumb("number of asked groups", len(data_asked_groups))

	var reply string
	// If_ok_do[int](demand, func() {
	// 	// for each row create a new group div
	// 	nodes := list_of_groups_node(created_groups_rows)
	// 	reply = d.Inline_nodes(nodes)
	// })
	return reply
}

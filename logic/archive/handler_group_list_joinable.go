package main

import (
	d "app_server/dom"
	c "github.com/g-ameline/colors"
)

func handler_group_list_joinable(demand *Demand) string {
	Breadcrumb("error state of demand", demand.Error)
	Breadcrumb("\n handling list of joinable groups demand", c.RED)
	defer c.Resetting()
	Breadcrumb("DATA FOR JOINABLE GROUPS")
	Print_dat_map(demand.Data)
	user_id := If_ok_do[string](demand, func() string { return demand.User_id })
	data_to_send := ss{}
	data_to_send["user_id"] = user_id
	joinable_groups_rows := If_ok_do[sa](demand, func() (sa, error) {
		return send_query_to_database_and_get_result(demand.DB_uri, data_to_send)
	})
	Breadcrumb("number of joinable groups", len(joinable_groups_rows))
	var reply string
	If_ok_do[int](demand, func() {
		// for each row create a new group div
		nodes := list_of_groups_node(joinable_groups_rows)
		reply = d.Inline_nodes(nodes)
	})

	return reply
}

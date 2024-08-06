package main

import (
	d "app_server/dom"
	c "github.com/g-ameline/colors"
)

func handler_group_list_created(demand *Demand) string {
	Breadcrumb("error state of demand", demand.Error)
	Breadcrumb("\n handling list of joined groups demand", c.RED)
	defer c.Resetting()
	Breadcrumb("demand data", demand.Data)
	Breadcrumb("DATA FOR CREATED GROUPS")
	Print_dat_map(demand.Data)
	data_user := sa{}
	data_user["user_id"] = demand.User_id
	created_groups_rows := If_ok_do[sa](demand, func() (sa, error) {
		return send_query_to_database_and_get_result(demand.DB_uri, data_user)
	})
	Breadcrumb("length created groups", len(created_groups_rows))

	var reply string
	If_ok_do[int](demand, func() {
		// for each row create a new group div
		nodes := list_of_groups_node(created_groups_rows)
		reply = d.Inline_nodes(nodes)
	})
	return reply
}

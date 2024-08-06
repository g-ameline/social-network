package main

import (
	c "github.com/g-ameline/colors"
)

func handler_group_list_invited(demand *Demand) string {
	Breadcrumb("error state of demand", demand.Error)
	Breadcrumb("\n handling list of asked groups demand", c.GREEN)
	defer c.Resetting()
	Breadcrumb("DATA FOR ALL ASKED GROUPS")
	Print_dat_map(demand.Data)
	user_id := If_ok_do[string](demand, func() string { return demand.User_id })
	data_user := ss{}
	data_user["user_id"] = user_id
	data_asked_groups := If_ok_do[sa](demand, func() (sa, error) {
		return send_query_to_database_and_get_result(demand.DB_uri, data_user)
	})
	Breadcrumb("number of groups user is invited", len(data_asked_groups))

	reply := sa{}
	If_ok_do[int](demand, func() {
		reply_ok := sa{}
		reply_ok["what"] = "group_pending"
		reply_ok["info"] = data_asked_groups
		reply = reply_ok
	})
	If_nok_do[int](demand, func() {
		reply_nok := sa{}
		reply_nok["what"] = "group list failed"
		reply_nok["info"] = demand.Error.Error()
		reply = reply_nok
	})
	return reply
}

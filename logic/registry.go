package main

var data_url = struct {
	root          string
	user_register string
	user_login    string
	user_logout   string
	get_id        string
	get_ids       string
	get_record    string
	get_records   string
	get_file      string
	insert_record string
	update_record string
	upsert_record string
	delete_record string
}{
	root:          url_data_root,
	user_register: url_data_root + "user_register",
	user_login:    url_data_root + "user_login",
	user_logout:   url_data_root + "user_logout",
	get_id:        url_data_root + "get_id",
	get_ids:       url_data_root + "get_ids",
	get_record:    url_data_root + "get_record",
	get_records:   url_data_root + "get_records",
	get_file:      url_data_root + "get_file",
	insert_record: url_data_root + "insert_record",
	update_record: url_data_root + "update_record",
	upsert_record: url_data_root + "upsert_record",
	delete_record: url_data_root + "delete_record",
}

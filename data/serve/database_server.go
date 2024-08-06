package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	// api "github.com/g-ameline/api_helper"
	api "github.com/g-ameline/api_helper"
	c "github.com/g-ameline/colors"
	db "github.com/g-ameline/sql_helper"
)

const root_folder = "social-network"
const database_filename = "social-network.db"
const http_protocol = "http://"
const domain = "localhost"
const port = ":5000"
const root = "/"
const main_endpoint = ""
const http_url = http_protocol + domain + port // http://localhost:5000
const home_url = http_url + root               // http://localhost:5000/

const login_path = root + "user_login"

const get_record_path = root + "get_record"
const get_records_path = root + "get_records"
const get_file_path = root + "get_file"
const insert_record_path = root + "insert_record"
const delete_record_path = root + "delete_record"
const update_record_path = root + "update_record"
const upsert_record_path = root + "upsert_record"
const get_ids_path = root + "get_ids"

var verbose = true

func main() {
	breadcrumb("start main")
	flag.Parse()
	arguments := Argufy(0, 1)
	if len(arguments) > 0 && arguments[0] == "quiet" {
		verbose = false
	}
	println("is verbose activated ? ", verbose)
	breadcrumb("check database")
	breadcrumb("search database file")
	path_to_database, err := path_to_other_file_in_repo(root_folder, database_filename)
	crash(err, "failed to find database folder to serve")
	breadcrumb("found database file at", path_to_database)
	breadcrumb("make the server")
	var server *http.ServeMux = http.NewServeMux()
	go func() {
		breadcrumb("activate server at", port)
		crash(http.ListenAndServe(port, server), "failed to start server")
	}()
	breadcrumb("listening at " + home_url)

	breadcrumb("	login request at", http_url+login_path)
	server.HandleFunc(login_path, handle_login(path_to_database))

	breadcrumb("	get file request at", get_file_path)
	server.HandleFunc(get_file_path, handle_get_file(path_to_database))
	breadcrumb("	delete stuff request at", delete_record_path)
	server.HandleFunc(delete_record_path, handle_delete_record(path_to_database))
	breadcrumb("	get stuff request at", get_record_path)
	server.HandleFunc(get_record_path, handle_get_record(path_to_database))
	breadcrumb("	get stuff request at", get_records_path)
	server.HandleFunc(get_records_path, handle_get_records(path_to_database))
	breadcrumb("	insert record request at", insert_record_path)
	server.HandleFunc(insert_record_path, handle_insert_record(path_to_database))
	breadcrumb("	update stuff request at", update_record_path)
	server.HandleFunc(update_record_path, handle_update_record(path_to_database))
	breadcrumb("	upsert stuff request at", upsert_record_path)
	server.HandleFunc(upsert_record_path, handle_upsert_record(path_to_database))
	breadcrumb("	get ids request at", get_ids_path)
	server.HandleFunc(get_ids_path, handle_get_ids(path_to_database))

	server.HandleFunc("/", func(responder http.ResponseWriter, request *http.Request) {
		c.Magenta("unhandled path ", request.URL.Path)
		http.Error(responder, "unhandled path", http.StatusNotFound)

	})
	select {} // to prevent main function to return while listener/handlers are still running
}

func path_to_other_file_in_repo(repo_root_folder, target_folder string) (string, error) {
	var find_the_root_of_the_repo func(string) string
	find_the_root_of_the_repo = func(path string) string {
		absolute_path, _ := filepath.Abs(path)
		if _, err := os.Stat(absolute_path + "/data"); err == nil {
			working_directory, _ := os.Getwd()
			path_to_root, _ := filepath.Rel(working_directory, absolute_path)
			return path_to_root
		}
		if filepath.Base(absolute_path) == repo_root_folder {
			working_directory, _ := os.Getwd()
			path_to_root, _ := filepath.Rel(working_directory, absolute_path)
			return path_to_root
		}
		if filepath.Base(absolute_path) == repo_root_folder+".git" {
			working_directory, _ := os.Getwd()
			path_to_root, _ := filepath.Rel(working_directory, absolute_path)
			return path_to_root
		}
		if filepath.Base(absolute_path) == "back" {
			working_directory, _ := os.Getwd()
			path_to_root, _ := filepath.Rel(working_directory, absolute_path)
			return path_to_root
		}
		if absolute_path == "/" {
			return ""
		}
		return find_the_root_of_the_repo(path + "../")
	}
	relative_path_to_root := find_the_root_of_the_repo("./")
	if relative_path_to_root == "" {
		return "", fmt.Errorf("could not find root repo")
		// relative_path_to_root = "../../.."
	}
	var result_path string
	// will walk the directory from the top root and if it finds the desired folder (suffix pattern)
	// will copy its value to result_path
	filepath.Walk(relative_path_to_root, func(walked_path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(walked_path, target_folder) {
			result_path = walked_path
			return nil
		}
		return nil
	})
	if result_path == "" {
		return "", fmt.Errorf("could not find root repo")
	}
	return result_path, nil
}

type sb = map[string]bool
type sa = map[string]any
type ss = map[string]string
type ssa = map[string]map[string]any
type sss = map[string]map[string]string
type s = string
type b = bool

func breadcrumb(hints ...any) {
	if verbose {
		for _, h := range hints {
			fmt.Print(h, " ")
		}
		fmt.Print("\n")
	}
}

func update_user_activity(db_path, id string) error {
	now_milliseconds := time.Now().UnixMilli()
	now_string := strconv.FormatInt(now_milliseconds, 10)
	breadcrumb("updating user id", id, "last activity at ", now_string)
	return db.Update_value(db_path, "users", id, "last_activity", now_string)
}

func print_dat_map(damap sa) {
	c.Print_map(damap)
}
func print_keys(damap sa) {
	c.Print_map(damap)
}

func from_sa_to_row(not_row map[string]any) map[string]string {
	row := map[string]string{}
	for k, v := range not_row {
		row[k] = v.(s)
	}
	return row
}

func make_null_if_empty(list_as_map map[string]string, key string) map[string]string {
	if _, ok := list_as_map["session"]; ok {
		if list_as_map["session"] == "" {
			list_as_map["sesssion"] = db.NULL
		}
	}
	return list_as_map
}

func Argufy(min int, max int) []string {
	arguments := flag.Args() // get argument and not flags
	switch {                 // we try to cath any wrong "inputing"
	case len(arguments) < min:
		fmt.Println("we need " + strconv.Itoa(min) + " argument")
		os.Exit(1)
	case len(arguments) > max:
		fmt.Println("we only accept " + strconv.Itoa(max) + " argument")
		os.Exit(1)
	}
	return arguments
}

func crash(err error, messages ...string) {
	if err != nil {
		var message string
		for _, mes := range messages {
			message += mes
		}
		panic(message + " " + err.Error())
	}
}

func respond_error(responder http.ResponseWriter, err error, messages ...string) {
	if err != nil {
		data := map[string]string{}
		for _, m := range messages {
			data["error"] += m + " "
		}
		data["error"] += err.Error()
		fmt.Println("data as error data", data)
		api.Respond_json_data(responder, data)
	}
}

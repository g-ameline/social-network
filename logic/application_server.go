package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
)

const root_folder = "social-network"
const name_website = "social-network"
const http_protocol = "http://"
const ws_protocol = "ws://"
const domain = "localhost"
const app_port = ":8000"
const data_port = ":5000"
const root = "/"

const url_app = http_protocol + domain + app_port     // http://localhost:8000
const url_app_root = url_app + root                   // http://localhost:8000/
const url_ws = ws_protocol + domain + app_port + root // ws://localhost:8000

const url_data = http_protocol + domain + data_port // http://localhost:5000
const url_data_root = url_data + root               // http://localhost:5000/

const db_root_uri = http_protocol + domain + data_port + root

const name_of_avatar_picture = "avatar"
const files_folder = "files"

const interactivity = "interactivity"
const interactivity_path = root + "interactivity"

const everyone = "everyone"

const session_max_period = 3000000                // 5 minutes ?
var session_to_last_activity = map[string]int64{} // [session]time_in_ms
var session_to_id = map[string]string{}           // [session]user_id
var id_to_display = map[string]*ws.Conn{}

func main() {

	c.Gray("start main")
	c.Gray("make the server")
	var server *http.ServeMux = http.NewServeMux()
	go func() {
		c.Gray("activate server at", app_port)
		m.Must(http.ListenAndServe(app_port, server), "failed to start server")
	}()
	c.Gray("listening at " + url_app_root)
	server.HandleFunc(slash(files_folder+"/"), handler_files)
	server.HandleFunc(slash(interactivity), handler_socket)
	server.HandleFunc(root, handler_home)
	handle_ajax(server, slash("user_login"), process_user_login)
	handle_ajax(server, slash("user_logout"), process_user_logout)
	handle_ajax(server, slash("user_register"), process_user_register)
	handle_ajax(server, slash("save_file_return_filepath"), process_file_upload)
	select {} // to prevent main function to return while listener/handlers are still running
}

func path_to_other_file_in_repo(repo_root_folder, target_folder string) (string, error) {
	var find_the_root_of_the_repo func(string) string
	find_the_root_of_the_repo = func(path string) string {
		absolute_path, _ := filepath.Abs(path)
		if filepath.Base(absolute_path) == repo_root_folder {
			wd, _ := os.Getwd()
			path_to_root, _ := filepath.Rel(wd, absolute_path)
			return path_to_root
		}
		if filepath.Base(absolute_path) == repo_root_folder+".git" {
			wd, _ := os.Getwd()
			path_to_root, _ := filepath.Rel(wd, absolute_path)
			return path_to_root
		}
		c.Gray("updated, session binders")
		c.Gray(session_to_last_activity)
		c.Gray(session_to_id)
		if filepath.Base(absolute_path) == "back" {
			wd, _ := os.Getwd()
			path_to_root, _ := filepath.Rel(wd, absolute_path)
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

func handle_ajax(server *http.ServeMux, path string, process_data func(sa) (string, error)) {
	server.HandleFunc(path, func(responder http.ResponseWriter, request *http.Request) {
		c.Gray("\nreceived ajax request at ", path)
		data := sa{}
		{ // if the credential path
			data_form, err := data_from_request(request)
			m.If_nil_do[nvm](err,
				func() {
					for k, v := range data_form {
						data[k] = v
					}
				})
			filepath, err := m.If_nil_do[string](err,
				func() (string, error) { return filepath_from_file_from_request(request) })
			m.If_nil_do[nvm](err,
				func() { data["filepath"] = filepath })
		}
		fmt.Println("after parsing for files", data)
		{ // if the already logged path
			session, err := session_from_cookie_from_request(request)
			user_id, err := m.If_nil_do[string](err,
				func() (string, error) { return session_valid(session) })
			m.If_nil_do[nvm](err,
				func() { data["session"] = session; data["user_id"] = user_id })
		}
		// --------------------- PROCESSING START
		htmx_fragment, err := process_data(data)
		m.Warn(err)
		// --------------------- PROCESSING OVER
		c.Gray("data after processing")
		c.Print_map(data)
		{
			session, ok := data["session"].(string)
			m.If_nil_try[nvm](m.Nok_to_err(ok, "no generated session in data"), func() {
				c.Gray("adding cookie to response", session)
				set_cookie_into_response(responder, name_website, session)
			})
			m.If_nok_do[nvm](ok, func() { c.Gray("no session was found in data") })
		}
		fmt.Fprintf(responder, htmx_fragment)
		// post_processing
		{
			session, is_session := data["session"].(string)
			user_id, is_user_id := data["user_id"].(string)
			err := m.Nok_to_err(is_session && is_user_id, "no session or user_id in data")
			m.If_nil_must[nvm](err, func() {
				update_ephemeris(session)
				update_almanac(session, user_id)
			})
		}
	})
}

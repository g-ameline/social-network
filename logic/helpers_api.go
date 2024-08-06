package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	// c "github.com/g-ameline/colors"
	mb "github.com/g-ameline/maybe"
)

const content_multipart = "multipart/form-data"
const content_form = "application/x-www-form-urlencoded"
const content_json = "application/json"

func content_type_from_request(request *http.Request) string {
	return cast_content_type(request.Header.Get("Content-type"))
}
func cast_content_type(input string) string {
	content_types := []string{content_multipart, content_form, content_json}
	for _, content_type := range content_types {
		if strings.Contains(input, content_type) {
			return content_type
		}
	}
	return "unknown"
}

func From_raw[M any](raw_message []byte) (M, error) {
	var message_once_parsed M
	err := json.Unmarshal(raw_message, &message_once_parsed)
	mb.Warn(err, "failed to parse/assert raw_message")
	return message_once_parsed, err
}
func To_raw[M any](to_parse M) ([]byte, error) {
	message_once_parsed, err := json.Marshal(to_parse)
	mb.Warn(err, "failed to marshall raw_message")
	return message_once_parsed, err
}

func Respond_demand(w http.ResponseWriter, reply ss) error {
	raw_reply, err := To_raw(reply)
	mb.Warn(err, "failed to parse reply into ")
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(raw_reply)
	mb.Warn(err, "failed to write back reply")
	return err
}
func Update_users_activities(session string) {
	if session != "" {
		session_to_last_activity[session] = time.Now().UnixMilli()
	}
}

func Update_users_binder(session, user_id string) {
	if session != "" {
		session_to_id[session] = user_id
	}
}

func filepath_from_file_from_request(request *http.Request) (filepath string, err error) {
	content_type := content_type_from_request(request)
	if content_type != content_multipart {
		return "", fmt.Errorf("content is not multipart")
	}
	err = request.ParseMultipartForm(32 << 10)
	if err != nil {
		return "", fmt.Errorf("failed to parse multipart form")
	}
	// c.Gray("files inside after parsing", len(request.MultipartForm.File), request.MultipartForm.File)
	if len(request.MultipartForm.File) != 1 {
		return "", fmt.Errorf("less or more than one file into multipart form")
	}
	if len(request.MultipartForm.File) == 1 {
		for _, v := range request.MultipartForm.File {
			file, err := v[0].Open()
			if err != nil {
				return "", fmt.Errorf("failed to open the file from multipart form")
			}
			return save_file_and_return_path(file, files_folder)
		}
	}
	return filepath, fmt.Errorf("impossible state")
}

func data_from_request(request *http.Request) (data map[string]any, err error) {
	data = sa{}
	switch content_type_from_request(request) {
	case content_json:
		data, err = data_from_request_json[map[string]any](request)
	case content_form:
		data, err = data_from_request_form(request)
	case content_multipart:
		data, err = data_from_request_multipart(request)
	}
	return data, err
}
func data_from_request_multipart(request *http.Request) (data map[string]any, err error) {
	data = map[string]any{}
	err = request.ParseMultipartForm(32 << 10)
	if err != nil {
		return data, err
	}
	// fmt.Println("form of multipart", request.MultipartForm)
	if len(request.MultipartForm.Value) <= 0 {
		return data, fmt.Errorf("no data into multipart form")
	}
	for k, v := range request.MultipartForm.Value {
		switch len(v) {
		case 1:
			data[k] = v[0]
		case 0:
			continue
		default:
			data[k] = v
		}
	}
	return data, err
}
func data_from_request_json[data_type any](request *http.Request) (data_type, error) {
	decoder := json.NewDecoder(request.Body)
	var data data_type
	err := decoder.Decode(&data)
	return data, err
}
func data_from_response_json[data_type any](request *http.Response) (data_type, error) {
	decoder := json.NewDecoder(request.Body)
	var data data_type
	err := decoder.Decode(&data)
	return data, err
}
func data_from_request_form(request *http.Request) (map[string]any, error) {
	err := request.ParseForm()
	data := sa{}
	if err != nil {
		return data, err
	}
	if len(request.PostForm) <= 0 {
		return data, fmt.Errorf("no data into post form")
	}
	for k, v := range request.PostForm {
		// fmt.Println("into form", k, v, len(v))
		switch len(v) {
		case 1:
			data[k] = v[0]
		case 0:
			continue
		default:
			data[k] = v
		}
	}
	return data, err
}
func save_file_and_return_path(file multipart.File, folder string) (string, error) {
	defer file.Close()
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return "", err
	}
	path := folder + fmt.Sprintf("/%d", time.Now().UnixNano())
	destination, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer destination.Close()
	_, err = io.Copy(destination, file)
	if err != nil {
		return "", err
	}
	return path, err
}
func send_query_to_database_and_parse_result[In any, Out any](database_url string, data_to_send In) (responded_data Out, err error) {
	response, err := post_json_data_to_url[In](database_url, data_to_send)
	if err != nil {
		return responded_data, fmt.Errorf("%v; %w", "bad response from data server", err)
	}
	responded_data, err = data_from_response_json[Out](response)
	// c.Print_map(responded_data)
	if err != nil {
		return responded_data, fmt.Errorf("%v; %w", "error when parsing response body into json", err)
	}
	return responded_data, error_from_json(responded_data)
}

func error_from_json(data_from_json any) error {
	detyped_map := any(data_from_json)
	switch detyped_map.(type) {
	case ss:
		if err_msg, ok := detyped_map.(ss)["error"]; ok {
			return fmt.Errorf(err_msg)
		}
	case sa:
		if err_msg, ok := detyped_map.(sa)["error"]; ok {
			return fmt.Errorf(err_msg.(string))
		}
	default:
		return error(nil)
	}
	return error(nil)
}

func post_json_data_to_url[data_type any](a_url string, data data_type) (response *http.Response, err error) {
	data_json, err := json.Marshal(data)
	if err != nil {
		return response, err
	}

	// //--------------debug
	// c.Yellow("if we try to read a buffer made of the json")
	// a_json, _ := json.Marshal(data)
	// b := bytes.NewBuffer(a_json)
	// // fmt.Fprintf(b, "world!")
	// b.WriteTo(os.Stdout)
	// //--------------debug

	data_reader := bytes.NewBuffer(data_json)
	response, err = http.Post(a_url, content_json, data_reader)
	if err != nil {
		return response, err
	}
	// defer func() {
	// 	err := response.Body.Close()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	if response.StatusCode != http.StatusOK {
		return response, fmt.Errorf("Error: unexpected status code: %d", response.StatusCode)
	}
	return response, err
}
func session_from_cookie_from_request(request *http.Request) (string, error) {
	cookie_struct, err := request.Cookie(name_website)
	switch true {
	case err != nil:
		return "", err
	case cookie_struct == nil:
		return "", err
	default:
		return cookie_struct.Value, err
	}
}
func set_cookie_into_response(responder http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(responder, &cookie)
}
func unset_session_cookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Value:    "",
		Name:     "social-network",
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, &cookie)
}

package main

import (
	"cmp"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	ws "github.com/gorilla/websocket"

	d "app_server/dom"
	m "github.com/g-ameline/maybe"
)

func slash(suffix string) string {
	return "/" + suffix
}

func id_selector(suffix string) string {
	return "#" + suffix
}

func From_sa_to_row(not_row map[string]any) map[string]string {
	row := map[string]string{}
	for k, v := range not_row {
		row[k] = v.(string)
	}
	return row
}

func Hash_it(something string) (hashed string) {
	hash_thing := sha256.New()
	hash_thing.Write([]byte(something))
	return string(hash_thing.Sum(nil))
}

func Hash_that(something string) (hashed string) {
	hash_thing := sha256.New()
	hash_thing.Write([]byte(something))
	bytes := hash_thing.Sum(nil)
	var numbers_as_string string
	for _, b := range bytes {
		numbers_as_string += strconv.Itoa(int(b))
	}
	return numbers_as_string
}

func from_map_to_slice(data_as_map sa) []any {
	var data_as_slice []any
	for _, value := range data_as_map {
		data_as_slice = append(data_as_slice, value)
	}
	return data_as_slice
}

func id_from_session(session string) (string, error) {
	user_id, ok := session_to_id[session]
	_, err := m.If_nok_do[nvm](ok, fmt.Errorf("session is unregistered"))
	_, err = m.If_nok_do[nvm](is_session_valid(session), fmt.Errorf("session is expired"))
	return user_id, err
}
func notify(user user, styler func(d.Node) d.Node, thingamajig ...any) {
	socket := user_to_socket(user)
	if socket == nil {
		return
	}
	fmt.Println("\nID_TO_DISPLAY", id_to_display)
	text := fmt.Sprint(thingamajig...)
	notification_node := notification_node(text)
	notification_node = styler(notification_node)
	packaged_node := package_node_toward_notifying_frame(notification_node)
	message := packaged_node.Inline()
	socket.WriteMessage(ws.TextMessage, []byte(message))
}
func notification_node(text string) d.Node {
	return d.New_div().
		Text(text).
		On_load(text)
}
func user_to_socket(a_user user) *ws.Conn {
	return id_to_display[a_user.id()]
}

func is_session_valid(session string) bool {
	last_activity, ok := session_to_last_activity[session]
	if !ok {
		return false
	}
	is_valid := session_max_period > time.Now().UnixMilli()-last_activity
	if !is_valid {
		unregister_session(session)
	}
	return is_valid
}

func session_valid(session string) (string, error) {
	fmt.Println("session", session)
	fmt.Println("almanac", session_to_last_activity)
	fmt.Println("ephemeris", session_to_id)
	user_id, ok := session_to_id[session]
	if !ok {
		return "", fmt.Errorf("no user_id matching session in almanac")
	}
	last_activity, ok := session_to_last_activity[session]
	if !ok {
		return "", fmt.Errorf("session not registered in ephemeris")
	}
	is_valid := time.Now().UnixMilli()-last_activity < session_max_period
	if !is_valid {
		fmt.Println("SESSION NOT VALID ANYMORE")
		fmt.Println(time.Now().UnixMilli() - last_activity)
		fmt.Println(session_max_period)
		unregister_session(session)
		return "", fmt.Errorf("session expired")
	}
	return user_id, nil
}
func update_ephemeris(session string) {
	session_to_last_activity[session] = time.Now().UnixMilli()
}

func update_almanac(session, user_id string) {
	session_to_id[session] = user_id
}

func unregister_session(session string) {
	fmt.Println("removing from binders", session)
	delete(session_to_last_activity, session)
	delete(session_to_id, session)
}

func upgrade_request(responder http.ResponseWriter, request *http.Request) (*ws.Conn, error) {
	upgrader := ws.Upgrader{ // needed by gorilla websocket package to transform ws request into websocket connection
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(request *http.Request) bool { return true },
	}
	return upgrader.Upgrade(responder, request, nil)
}
func send_message(socket *ws.Conn, message string) error {
	if message == "" {
		return nil
	}
	return socket.WriteMessage(ws.TextMessage, []byte(message))
}

func data_for_db_insertion(table_name string, record_data ss) sa {
	data_for_database := sa{}
	data_for_database["table"] = table_name
	data_for_database["record"] = record_data
	return data_for_database
}
func raw_to_data[M any](raw_message []byte) (M, error) {
	var message_once_parsed M
	err := json.Unmarshal(raw_message, &message_once_parsed)
	m.Warn(err, "failed to parse/assert raw_message")
	return message_once_parsed, err
}
func data_to_raw[M any](to_parse M) ([]byte, error) {
	message_once_parsed, err := json.Marshal(to_parse)
	m.Warn(err, "failed to marshall data")
	return message_once_parsed, err
}
func data_to_string[M any](to_parse M) (string, error) {
	message_once_parsed, err := json.Marshal(to_parse)
	m.Warn(err, "failed to marshall data")
	return string(message_once_parsed), err
}

func ordered_keys[k cmp.Ordered, v any](a_map map[k]v) []k {
	keys_to_order := make([]k, len(a_map))
	i := 0
	for key := range a_map {
		keys_to_order[i] = key
		i++
	}
	slices.Sort(keys_to_order)
	return keys_to_order
}

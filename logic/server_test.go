package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	api "github.com/g-ameline/api_helper"
	c "github.com/g-ameline/colors"
)

const web_port = ":3000"
const web_http_url = http_protocol + domain + web_port // http://localhost:3000
var testing_http_url = http_url

const pres_register_path = root + "register"
const pres_login_path = root + "login"
const pres_logout_path = root + "logout"
const pres_user_search_path = root + "user_search"

const pres_asked_groups_path = root + "group_pending"
const pres_created_groups_path = root + "group_created"
const pres_joinable_groups_path = root + "group_list"
const pres_joined_groups_path = root + "group_my"

const pres_new_group_path = root + "group_create"
const pres_new_joining_path = root + "group_request"

func Test_server(t *testing.T) {
	fmt.Println("STARTING")
	t.Log("Say hie")
	test_group_stuff()
	// poster_public := map[string]any{}
	// poster_private := map[string]any{}
	// poster_exclusive := map[string]any{}

	// follower := map[string]any{}
	// solower := map[string]any{}

}
func test_group_stuff() {
	joiner := fresh_user()
	invitee := fresh_user()
	creator := fresh_user()
	fmt.Println("\ncreating group")
	// group data
	creator["name"] = "new_group_" + random_syllables(9, 13)
	creator["description"] = "description_" + random_syllables(9, 13)

	fmt.Println("\ngetting list of created groups, should work")
	send_data(creator, created_groups_path, c.Blue_light, c.BLACK)

	output := send_data(creator, new_group_path, c.Black, c.RED_LIGHT)
	created_group_id := output["group_id"]

	fmt.Println("\ngetting list of created groups, should work")
	send_data(creator, created_groups_path, c.Black, c.GRAY_LIGHT)

	fmt.Println("\ngetting list of groups where invitee is invited")
	send_data(invitee, invited_groups_path, c.Gray_light, c.BLACK)

	fmt.Println("\ncreator  , invite invitee in created group")
	creator["group_id"] = created_group_id
	print_data(invitee)
	creator["joiner_id"] = invitee["user_id"]
	invitee_joining_data := send_data(creator, creator_invite_user_path, c.Black, c.YELLOW_LIGHT)
	println(invitee_joining_data, joiner)

	fmt.Println("\ngetting list of groups where invitee is invited")
	send_data(invitee, invited_groups_path, c.Gray_light, c.BLACK)

	fmt.Println("\n invitee , accept creator demand  ")
	invitee["joining_id"] = invitee_joining_data["joining_id"]
	send_data(invitee, user_accept_group_path, c.Black, c.CYAN_LIGHT)

	fmt.Println("\ngetting list of joined groups, should work")
	send_data(invitee, joined_groups_path, c.White, c.BLACK)

	fmt.Println("\ngetting list of asked groups, should work")
	send_data(joiner, asked_groups_path, c.Cyan_light, c.BLACK)

	fmt.Println("\njoiner  , ask creator to join group")
	joiner["group_id"] = created_group_id
	joiner_joining_data := send_data(joiner, user_ask_creator_path, c.Black, c.BLUE_LIGHT)
	println(joiner_joining_data)
	fmt.Println("\ngetting list of asked groups, should work")
	send_data(joiner, asked_groups_path, c.Cyan_light, c.BLACK)

	// fmt.Println("\n creator , accept joiner in ")
	// creator["joiner_id"] = joiner["user_id"]
	// creator["joining_id"] = joiner_joining_data["joining_id"]
	// send_data(creator, creator_accept_user_path, c.Black, c.MAGENTA_LIGHT)
	// fmt.Println("\n invitee , leave group  ")
	// send_data(invitee, user_leave_group_path, c.Black, c.WHITE)

	// fmt.Println("\ngetting list of joined groups, should work")
	// send_data(invitee, joined_groups_path, c.White, c.BLACK)

	// const new_group_path = root + "group_create"
	// const user_ask_creator_path = root + "group_request"
	// const creator_invite_user_path = root + "group_invit"
	// const user_accept_group_path = root + "TODO"
	// const creator_refuse_user_path = root + "TODO"
	// const creator_accept_user_path = root + "group_approve_owner"
	// const user_refuse_group_path = root + "group_approve_user"
	// const user_leave_group_path = root + "group_leave"

	// fmt.Println("\ngetting list of joinable groups, should work")
	// joinable_groups := send_data(joiner, joinable_groups_path, c.Black, c.WHITE)

	// fmt.Println("group to join", joinable_groups[group_id])
	// ok_client_data["group_id"] = group_id
	// output = send_data(ok_client_data, new_joining_path, c.Black, c.BLUE_LIGHT)
}

func send_data(sending_data sa, sending_path, color_fb, color_bg string) sa {
	fmt.Println(color_fb, color_bg, "sending this data")
	print_data(sending_data)
	// sending_uri := http_url + sending_path
	sending_uri := testing_http_url + sending_path
	fmt.Println("to", sending_uri)
	response, err := api.Post_json_data_to_url(sending_uri, sending_data)
	received_data := sa{}
	if err != nil {
		fmt.Println("error when sending to app server", err)
		return received_data
	}
	fmt.Println("we could send data to app server", response.StatusCode)
	cookie, err := api.Get_cookie_data_from_response(response, "social-network")
	if err != nil {
		panic(err)
	}
	received_data, err = api.Get_data_from_response_json[map[string]any](response)
	if err != nil {
		fmt.Println("error when parsing data from response", err)
		return received_data
	}
	if cookie != nil {
		fmt.Println("cookie from request", cookie)
		received_data["cookie"] = cookie
	}
	fmt.Println("data from app server")
	print_data(received_data)
	c.Resetting()
	return received_data
}

func print_data(data any) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(data)
		return
	}
	if len(b) < 300 {
		fmt.Println(string(b))
		return
	}
	fmt.Println(string(b[:300]) + "\n...")
}

// upgrade websocket:
// time.Sleep(1 * time.Second)
// fmt.Println("ws url", ws_url)
// fmt.Println("ws path", ws_path)
// connection, response, err := ws.DefaultDialer.Dial(ws_url, nil)
// time.Sleep(1 * time.Second)
// connection.WriteMessage(ws.TextMessage, []byte("PATATE-MAN EST LA"))
// connection.SetReadDeadline(time.Now().Add(2 * time.Second))
// mt, ms, err := connection.ReadMessage()
// if err != nil {
// 	t.Error(err)
// }
// fmt.Println("from websocket conenction", mt, string(ms), err)

func random_figure(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = "0123456789"
	letters = letters[:n]
	return string(letters[rand.Intn(len(letters))])
}

func random_string(strings ...string) string {
	rand.Seed(time.Now().UnixNano())
	return string(strings[rand.Intn(len(strings))])
}

func randomstring(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = "nadwqgeiuo39482 dfwq2387 klrqwjbvg&"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func random_vowels() string {
	rand.Seed(time.Now().UnixNano())
	var pool = "aeiou"
	lengths := []int{1, 1, 1, 1, 1, 2, 3}
	length := lengths[rand.Intn(len(lengths))]
	var vowels []byte
	for i := 0; i < length; i++ {
		vowels = append(vowels, pool[rand.Intn(len(pool))])
	}
	return string(vowels)
}

func random_consonants() string {
	rand.Seed(time.Now().UnixNano())
	var pool = "bkjptdlrvzw"
	lengths := []int{1, 1, 1, 1, 2, 2, 3}
	length := []int{1, 1, 1, 1, 1, 1, 1, 1, 2, 2}[rand.Intn(len(lengths))]
	consonant := pool[rand.Intn(len(pool))]
	var consonants []byte
	for i := 0; i < length; i++ {
		consonants = append(consonants, consonant)
	}
	return string(consonants)
}

func heads_or_tails() bool {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		return true
	}
	return false
}

func random_syllables(min, max int) string {
	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(max-min) + min
	zero_or_one := rand.Intn(2)
	var heads_or_tails [2]func() string
	heads_or_tails[zero_or_one] = random_consonants
	heads_or_tails[(zero_or_one+1)%2] = random_vowels
	var word []byte
	for i := 0; i < length; i++ {
		word = append(word, []byte(heads_or_tails[i%2]())...)
	}
	return string(word)
}

func random_sentence(length int) string {
	rand.Seed(time.Now().UnixNano())
	var phrase string
	space := " "
	for i := 0; i < length; i++ {
		phrase += space + random_syllables(2, 5)
	}
	return string(phrase)
}
func random_key_from_row[T any](row map[string]T) string {
	for key := range row {
		return key
	}
	panic("list empty I guess")
}

func n_different_random_keys_from_row[T any](n int, row map[string]T) []string {
	if len(row) == 0 {
		panic("empty list dude")
	}
	keys := []string{}
	count := 0
	for key := range row {
		if count == n {
			break
		}
		keys = append(keys, key)
	}
	return keys
}

func random_value_from_row(row map[string]string) string {
	for _, value := range row {
		return value
	}
	return "yolo"
}

func random_date() string {
	rand.Seed(time.Now().UnixNano())
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2099, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	date := time.Unix(sec, 0)
	date_text := date.Format(time.DateTime)
	return date_text
	// date := time.Time.GoString(time.Unix(sec, 0))
	// return date.Format(time.DateTime)
}
func test_credentials() {
	fmt.Println("\ntry to register with minimal random credentials, should be ok")
	ok_client_data := map[string]any{}
	ok_client_data["email"] = random_syllables(9, 12) + "@" + random_syllables(9, 12) + "." + random_syllables(9, 12)
	ok_client_data["password"] = "a_password"
	send_data(ok_client_data, register_path, c.Yellow_light, c.CYAN)

	fmt.Println("\nre try with same credentials, should failed")
	send_data(ok_client_data, register_path, c.Magenta_light, c.CYAN)

	// send login demand
	fmt.Println("\ntry to login with previous credentials, should work")
	logged_in_data := send_data(ok_client_data, login_path, c.Black, c.RED)
	ok_client_data["session"] = logged_in_data["user_id"]
	ok_client_data["session"] = logged_in_data["uuid"]

	fmt.Println("\ntry to login with random credentials, should fail")
	bad_client_data := map[string]any{}
	bad_client_data["email"] = random_syllables(7, 8) + "@e.mail"
	bad_client_data["password"] = random_syllables(6, 8)
	send_data(bad_client_data, login_path, c.Green_light, c.RED)

	// logout
	fmt.Println("\ntry to logout with previous credentials, should work", http_url+logout_path)
	fmt.Println("data", ok_client_data)
	send_data(ok_client_data, logout_path, c.Black, c.RED)

	fmt.Println("\ntry to logout with wrong credentials, should failed", http_url+logout_path)
	bad_client_data["email"] = random_syllables(7, 8) + "@e.mail"
	bad_client_data["session"] = ""
	send_data(bad_client_data, logout_path, c.Black, c.RED)

	// send login demand
	fmt.Println("\ntry to login with previous credentials, should work")
	logged_in_data = send_data(ok_client_data, login_path, c.Black, c.RED)
	ok_client_data["session"] = logged_in_data["user_id"]
	ok_client_data["session"] = logged_in_data["uuid"]

	fmt.Println("\ntry to search for users, itself, should work")
	ok_client_data["user"] = ok_client_data["email"]
	send_data(ok_client_data, user_search_path, c.Black, c.RED)

}

func fresh_user() map[string]any {
	new_user := map[string]any{}
	new_user["email"] = random_syllables(9, 12) + "@" + random_syllables(9, 12) + "." + random_syllables(9, 12)
	new_user["password"] = random_syllables(9, 12)
	registration_data := send_data(new_user, register_path, c.Yellow_light, c.CYAN)
	if _, ok := registration_data["error"]; ok {
		panic("registration failed")
	}
	if user_id, ok := registration_data["user_id"]; ok {
		new_user["user_id"] = user_id
	}
	login_data := send_data(new_user, login_path, c.Yellow_light, c.CYAN)
	if _, ok := login_data["error"]; ok {
		panic("login failed")
	}
	if session, ok := login_data["uuid"]; ok {
		new_user["session"] = session
	}
	return new_user
}

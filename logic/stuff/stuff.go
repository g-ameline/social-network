package stuff

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	c "github.com/g-ameline/colors"
	mb "github.com/g-ameline/maybe"
	db "github.com/g-ameline/sql_helper"
	ws "github.com/gorilla/websocket"
)

type sb = map[string]bool
type sa = map[string]any
type ss = map[string]string
type ssa = map[string]map[string]any
type sss = map[string]map[string]string
type s = string
type b = bool
type Demand struct {
	// backtrack data
	Request    *http.Request
	Responder  http.ResponseWriter
	Connection *ws.Conn
	Data       sa
	Method     string
	What       string
	// constants
	Files          string
	Period_session int64
	// processing data
	Error  error
	DB_uri string
	// intra user data
	User_id string
	Session string
	// inter user data
	Ephemeris map[string]int64    // [session]last activity in ms
	Almanac   map[string]string   // [session]user id
	Hub       map[string]*ws.Conn // [session]user id
}

var verbose = true

func Fresh_demand() Demand {
	new_demand := *new(Demand)
	new_demand.Data = sa{}
	new_demand.Ephemeris = map[string]int64{}
	new_demand.Almanac = map[string]string{}
	return new_demand
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

func Breadcrumb(hints ...any) {
	if verbose {
		for _, h := range hints {
			fmt.Print(h, " ")
		}
		fmt.Print("\n")
	}
}

func Print_dat_map(damap sa) {
	for k, v := range damap {
		c.Print(k, c.Magenta)
		print(" : ")
		fmt.Printf("%T = ", v)
		c.Println(v, c.Cyan)
	}
}
func Print_keys(damap sa) {
	for k, v := range damap {
		c.Print(k, c.Magenta)
		print(c.Cyan)
		fmt.Printf(": %T \n", v)
	}
}

func From_sa_to_row(not_row map[string]any) map[string]string {
	row := map[string]string{}
	for k, v := range not_row {
		row[k] = v.(s)
	}
	return row
}

func Make_null_if_empty(list_as_map map[string]string, key string) map[string]string {
	if _, ok := list_as_map["session"]; ok {
		if list_as_map["session"] == "" {
			list_as_map["sesssion"] = db.NULL
		}
	}
	return list_as_map
}

func If_ok_do[T any](demand *Demand, something any) T {
	if demand.Error != nil {
		Breadcrumb("..skip.. " + demand.Error.Error())
		return *new(T)
	}
	switch something.(type) {
	case error:
		err := something.(error)
		demand.Error = err
		return something.(T)
	case func() error:
		err := something.(func() error)()
		demand.Error = err
		return *new(T)
	case T:
		return something.(T)
	case func() T:
		f := something.(func() T)
		return f()
	case func() (T, error):
		var res T
		res, demand.Error = something.(func() (T, error))()
		return res
	case func():
		something.(func())()
		return *new(T)
	}
	fmt.Printf("Underlying Type: %T\n", something)
	demand.Error = fmt.Errorf("badly fail at func Type assertion when binding")
	panic("failed into assertion")
}

func If_nok_do[T any](demand *Demand, something any) T {
	if demand.Error == nil {
		Breadcrumb("..skip..", demand.Error)
		return *new(T)
	}
	switch something.(type) {
	case error:
		err := something.(error)
		demand.Error = err
		return something.(T)
	case func() error:
		err := something.(func() error)()
		demand.Error = err
		return *new(T)
	case T:
		return something.(T)
	case func() T:
		f := something.(func() T)
		return f()
	case func() (T, error):
		var res T
		res, demand.Error = something.(func() (T, error))()
		return res
	case func():
		something.(func())()
		return *new(T)
	}
	fmt.Printf("Underlying Type: %T\n", something)
	demand.Error = fmt.Errorf("badly fail at func Type assertion when binding")
	panic("failed into assertion")
}

func If_ok_try[T any](demand *Demand, something any) T {
	// fmt.Printf("inside ifokdo output Type: %T\n", *new(T))
	if demand.Error != nil {
		Breadcrumb("..skip.. " + demand.Error.Error())
		return *new(T)
	}
	switch something.(type) {
	case error:
		return something.(T)
	case func() error:
		something.(func() error)()
		return *new(T)
	case T:
		return something.(T)
	case func() T:
		f := something.(func() T)
		return f()
	case func() (T, error):
		res, err := something.(func() (T, error))()
		if err != nil {
			res = *new(T)
		}
		return res
	case func():
		something.(func())()
		return *new(T)
	}
	fmt.Printf("Underlying Type: %T\n", something)
	panic("failed into assertion")
}
func If_nok_try[T any](demand *Demand, something any) T {
	// fmt.Printf("inside ifokdo output Type: %T\n", *new(T))
	if demand.Error == nil {
		Breadcrumb("..skip.. " + demand.Error.Error())
		return *new(T)
	}
	switch something.(type) {
	case error:
		return something.(T)
	case func() error:
		err := something.(func() error)()
		return err.(T)
	case T:
		return something.(T)
	case func() T:
		f := something.(func() T)
		return f()
	case func() (T, error):
		var res T
		res, _ = something.(func() (T, error))()
		return res
	case func():
		something.(func())()
		return *new(T)
	}
	fmt.Printf("Underlying Type: %T\n", something)
	panic("failed into assertion")
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

func Respond_demand(w http.ResponseWriter, reply ss) error {
	raw_reply, err := To_raw(reply)
	mb.Warn(err, "failed to parse reply into ")
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(raw_reply)
	mb.Warn(err, "failed to write back reply")
	return err
}
func Update_users_activities(session_to_last_activity map[string]int64, session string) {
	if session != "" {
		session_to_last_activity[session] = time.Now().UnixMilli()
	}
}

func Update_users_binder(session_to_user_id map[string]string, session, user_id string) {
	if session != "" {
		session_to_user_id[session] = user_id
	}
}

func Check_session_valid(demand *Demand) (user_id string, err error) {
	var session_to_last_activity map[string]int64 = demand.Ephemeris
	var session_period int64 = demand.Period_session
	var session string = demand.Session
	if session == "" {
		return user_id, fmt.Errorf("there is no sesssion in demand")
	}
	now := time.Now().UnixMilli()
	last_activity, ok := session_to_last_activity[session]
	if !ok {
		return "", fmt.Errorf("could not find session id")
	}
	if session_period < now-last_activity {
		Breadcrumb("user session expired")
		return "", fmt.Errorf("session expired")
	}
	Breadcrumb("user session still valid")
	user_id, ok = demand.Almanac[session]
	if !ok {
		return "", fmt.Errorf("could not find matching user_id in server")
	}
	return user_id, nil
}
func Save_file_and_return_path(file multipart.File, filename, folder string) (string, error) {
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

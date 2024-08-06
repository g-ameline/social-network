package main

import (
	"database/sql"
	"strings"

	// "errors"
	"fmt"
	"log"
	"math/rand"

	col "github.com/g-ameline/colors"
	db "github.com/g-ameline/sql_helper"

	// "strconv"
	"time"
)

const path_to_db = "../social-network.db"

func main() {
	check_table_exist(path_to_db, "database")
	version_the_database(path_to_db, "0.1", "0.2")
	check_table_exist(path_to_db, "users")
	populate_users(path_to_db)
	populate_followings(path_to_db)
	populate_groups(path_to_db)
	populate_joinings(path_to_db)
	populate_events(path_to_db)
	populate_attendings(path_to_db)
	populate_public_posts(path_to_db)
	populate_exclu_posts(path_to_db)
	populate_private_posts(path_to_db)
	populate_group_posts(path_to_db)
}

func version_the_database(path_to_database, version_before, version_after string) {
	rows, err := db.Get_all_rows_sorted(path_to_database, "database", "id")
	fatal_error(err)
	// check find old version
	if (version_before != "") && (rows[len(rows)-1]["version"] != version_before) {
		log.Fatalln("wrong version")
	}
	// insert new version
	new_database_version_row := map[string]string{}
	new_database_version_row["version"] = version_after
	_, err = db.Insert_row(path_to_database, "database", new_database_version_row)
	fatal_error(err)
}

var database *sql.DB

const database_driver = "sqlite3"

func check_table_exist(path_to_database, table_name string) {
	res, err := db.Get_all_rows_from_table(path_to_db, table_name, "id")
	fatal_error(err)
	fmt.Println(res)
}

func populate_users(Path_to_db string) {
	fmt.Println("inserting users")
	number_of_users := 30
	// create users
	for i := 0; i < number_of_users; i++ {
		fmt.Println("creating a user row")
		new_user_row := map[string]string{} //:= new_row("user")
		new_user_row["email"] = random_syllables(1, 5) + "@" + random_syllables(1, 3) + "." + randomstring(2)
		new_user_row["password"] = "parool"
		new_user_row["first_name"] = random_syllables(0, 3)
		new_user_row["last_name"] = random_syllables(0, 3)
		new_user_row["birth"] = random_date()
		new_user_row["nickname"] = random_syllables(0, 4)
		new_user_row["about"] = random_sentence(3)
		new_user_row["private"] = [2]string{"0", "1"}[rand.Intn(2)]
		fmt.Println("new user row to insert", new_user_row)
		_, err := db.Insert_row(Path_to_db, "users", new_user_row)
		talk("error after inserting row user", err)
	}
	fmt.Println("user inserted")
	fmt.Println("")
}

func populate_followings(Path_to_db string) {
	fmt.Println("inserting followings")
	// create users
	users_ids, err := db.Get_ids(Path_to_db, "users")
	// number_of_followings := len(users_ids) * len(users_ids)
	fmt.Println(users_ids)
	talk("failed to get users", err)
	for follower_id := range users_ids {
		for followee_id := range users_ids {
			if rand.Intn(8) == 0 {
				continue
			}
			fmt.Print("creating a following row ")
			new_following_row := map[string]string{} //:= new_row("user")
			new_following_row["follower_id"] = follower_id
			new_following_row["followee_id"] = followee_id
			new_following_row["approval"] = random_string("-1", "0", "1")
			fmt.Println("new following row to insert", new_following_row)
			_, err := db.Insert_row(Path_to_db, "followings", new_following_row)
			talk("error after inserting row user", err)
			new_following_row = nil
		}
	}
	fmt.Println("following inserted")
	fmt.Println("")
}
func populate_groups(Path_to_db string) {
	fmt.Println("inserting groups")
	users_ids, err := db.Get_ids(Path_to_db, "users")
	fmt.Println(users_ids["1"], len(users_ids))
	talk("failed to get users", err)
	for user_id, _ := range users_ids {
		if rand.Intn(4) == 0 {
			continue
		}
		fmt.Print("creating a group row ")
		new_group_row := map[string]string{} //:= new_row("user")
		new_group_row["title"] = random_syllables(2, 6)
		new_group_row["creator_id"] = user_id
		new_group_row["description"] = random_sentence(9)
		fmt.Println("new group row to insert", new_group_row)
		_, err := db.Insert_row(Path_to_db, "groups", new_group_row)
		talk("error after inserting row user", err)
		new_group_row = nil
	}
	fmt.Println("group inserted")
	fmt.Println("")
}

func populate_joinings(path_to_database string) {
	fmt.Println("inserting joinings")
	// create users
	users_ids, err := db.Get_ids(path_to_database, "users")
	groups_ids, err := db.Get_ids(path_to_database, "groups")
	talk("failed to get users", err)
	for user_id := range users_ids {
		for group_id := range groups_ids {
			if rand.Intn(6) == 0 {
				continue
			}
			fmt.Println("creating a joining row")
			new_joining_row := map[string]string{} //:= new_row("user")
			new_joining_row["joiner_id"] = user_id
			new_joining_row["group_id"] = group_id
			new_joining_row["approval_creator"] = random_string("-1", "0", "1", "1", "1", "1", "1", "1")
			new_joining_row["approval_joiner"] = random_string("-1", "0", "1", "1", "1", "1", "1", "1", "1")
			fmt.Println("new joining row to insert", new_joining_row)
			_, err := db.Insert_row(path_to_database, "joinings", new_joining_row)
			talk("error after inserting row user", err)
		}
	}
	fmt.Println("joining inserted")
	fmt.Println("")
}

func populate_events(path_to_database string) {
	fmt.Println("inserting events")
	number_of_events := 15
	// create users
	groups_ids, err := db.Get_ids(path_to_database, "groups")
	talk("failed to get users", err)
	for i := 0; i < number_of_events; i++ {
		a_group_id := random_key_from_rows[bool](groups_ids)
		fmt.Println("creating a event row")
		new_event_row := map[string]string{} //:= new_row("user")
		new_event_row["group_id"] = a_group_id
		new_event_row["title"] = random_syllables(2, 5)
		new_event_row["description"] = random_sentence(6)
		new_event_row["date"] = random_date()
		fmt.Println("new event row to insert", new_event_row)
		_, err := db.Insert_row(path_to_database, "events", new_event_row)
		talk("error after inserting row user", err)
		new_event_row = nil
	}
	fmt.Println("event inserted")
	fmt.Println("")
}

func populate_attendings(path_to_database string) {
	fmt.Println("inserting attendings")
	groups_ids, err := db.Get_ids(path_to_database, "groups")
	fmt.Println("fetched groupds", groups_ids)
	talk("faield getting groups;srows", err)
	// events_ids, err := db.Get_ids(path_to_database, "events")
	// for each group
	for group_id := range groups_ids {
		// get all members
		joining_rows, err := db.Get_all_rows_from_table(path_to_database, "joinings", "id")
		talk("failed to fetch members (by joinings)", err)
		for _, joining_row := range joining_rows {
			if joining_row["group_id"] != group_id {
				continue
			}
			if joining_row["approval_creator"] != "1" {
				continue
			}
			if joining_row["approval_joiner"] != "1" {
				continue
			}
			member_id := joining_row["Joiner_id"]
			// create a attending row
			new_attending_row := map[string]string{}
			new_attending_row["attender_id"] = member_id
			new_attending_row["group_id"] = group_id
			new_attending_row["coming"] = random_string("0", "1")
			fmt.Println("new attending row to insert", new_attending_row)
			_, err := db.Insert_row(path_to_database, "attendings", new_attending_row)
			talk("error after inserting row attending", err)
		}
		fmt.Println("")
	}
	fmt.Println("attending inserted")
	fmt.Println("")
}
func populate_public_posts(path_to_database string) {
	fmt.Print(col.All["WHITE"], col.All["black"])
	fmt.Println("PUBLIC NOTE")
	fmt.Println("inserting notes")
	users_ids, err := db.Get_ids(path_to_database, "users")
	talk("failed to get users", err)
	for addresser_id := range users_ids {
		new_note_row := map[string]string{} //:= new_row("user")
		new_note_row["text"] = random_sentence(5)
		new_note_row["author_id"] = addresser_id
		note_id, err := db.Insert_row(path_to_database, "notes", new_note_row)
		talk("error after inserting row note 2", err)
		fmt.Println("public note inserted")
		// addressing
		new_addressing_row := map[string]string{}
		new_addressing_row["note_id"] = note_id
		_, err = db.Insert_row(path_to_database, "public_addressings", new_addressing_row)
		talk("error after inserting row addressing 1", err)
		fmt.Println("public addressing inserted")
		fmt.Println("")
	}
}

// PRIVATE

func populate_private_posts(path_to_database string) {
	fmt.Print(col.All["WHITE"], col.All["black"])
	fmt.Println("PRIVATE NOTE")
	fmt.Println("inserting notes")
	users_ids, err := db.Get_ids(path_to_database, "users")
	talk("failed to get users", err)
	for addresser_id := range users_ids {
		new_note_row := map[string]string{} //:= new_row("user")
		new_note_row["text"] = random_sentence(5)
		new_note_row["author_id"] = addresser_id
		note_id, err := db.Insert_row(path_to_database, "notes", new_note_row)
		talk("error after inserting row note 2", err)
		// fmt.Println("private note inserted")
		// addressing
		new_addressing_row := map[string]string{}
		new_addressing_row["note_id"] = note_id
		_, err = db.Insert_row(path_to_database, "private_addressings", new_addressing_row)
		talk("error after inserting row addressing 1", err)
		// fmt.Println("private addressing inserted")
		fmt.Println("")
	}
}

func populate_exclu_posts(path_to_database string) {
	fmt.Print(col.All["YELLOW"], col.All["black"])
	fmt.Println("EXCLU NOTE")
	fmt.Println("inserting notes")
	users_ids, err := db.Get_ids(path_to_database, "users")
	talk("failed to get users", err)
	followings_rows, err := db.Get_all_rows_from_table(path_to_database, "followings", "id")
	talk("failed to get followings", err)

	for addresser_id := range users_ids {
		new_note_row := map[string]string{} //:= new_row("user")
		new_note_row["text"] = random_sentence(5)
		new_note_row["author_id"] = addresser_id
		note_id, err := db.Insert_row(path_to_database, "notes", new_note_row)
		talk("error after inserting row note 2", err)
		fmt.Println("exclu note inserted", note_id)
		followers := map[string]bool{}
		for _, following_row := range followings_rows {
			if following_row["followee_id"] != addresser_id {
				continue
			}
			if following_row["approval"] != "1" {
				continue
			}
			follower_id := following_row["follower_id"]
			followers[follower_id] = true
		}
		fmt.Println("that many followers", len(followers))

		some_followers_ids := n_different_random_keys_from_rows[bool](3, followers)
		for _, a_follower_id := range some_followers_ids {
			new_addressing_row := map[string]string{}
			new_addressing_row["addressee_id"] = a_follower_id
			new_addressing_row["note_id"] = note_id
			_, err := db.Insert_row(path_to_database, "exclusive_addressings", new_addressing_row)
			talk("error after inserting row addressing 1", err)
			fmt.Println("    exclu addressing inserted")
		}
		fmt.Println("")
	}
	exclu_addressings_rows, err := db.Get_all_rows_from_table(path_to_database, "exclusive_addressings", "id")
	fmt.Println("that many exclu addressings ", len(exclu_addressings_rows))
}

func populate_group_posts(path_to_database string) {
	fmt.Print(col.All["WHITE"], col.All["black"])
	fmt.Println("Group NOTE")
	fmt.Println("inserting notes")
	joinings_rows, err := db.Get_all_rows_from_table(path_to_database, "joinings", "id")
	talk("failed to get joinings", err)
	// for each user find another one that he follows or is followed by
	fmt.Println("number of joinings rows", len(joinings_rows))
	for _, joining_row := range joinings_rows {
		if joining_row["approval_creator"] != "1" && joining_row["approval_joiner"] != "1" {
			continue
		}
		addresser_id := joining_row["joinee_id"]
		group_id := joining_row["group_id"]
		// fmt.Println("group", addressee_id, "writer", addresser_id)
		new_note_row := map[string]string{} //:= new_row("user")
		new_note_row["text"] = random_sentence(5)
		new_note_row["author_id"] = addresser_id
		note_id, err := db.Insert_row(path_to_database, "notes", new_note_row)
		talk("error after inserting row note 2", err)
		// fmt.Print("exclu note inserted ")
		new_addressing_row := map[string]string{}
		new_addressing_row["group_id"] = group_id
		new_addressing_row["note_id"] = note_id
		_, err = db.Insert_row(path_to_database, "group_addressings", new_addressing_row)
		talk("error after inserting row addressing 1", err)
		// fmt.Println("Group addressing inserted")
		// fmt.Println("")
	}
}

func fatal_error(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func talk(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
	}
}
func capitalize(word string) string {
	name := strings.ToUpper(word[:1]) + word[1:]
	return name
}
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
func random_key_from_rows[T any](rows map[string]T) string {
	for key := range rows {
		return key
	}
	log.Fatal("list empty I guess")
	panic(9)
}

func n_different_random_keys_from_rows[T any](n int, rows map[string]T) []string {
	if len(rows) == 0 {
		log.Fatalln("empty list dude")
	}
	keys := []string{}
	count := 0
	for key := range rows {
		if count == n {
			break
		}
		keys = append(keys, key)
		count++
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

package main

import (
	"database/sql"
	// "errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	db_f "rtf_db/database_functions"
	// "strconv"
	"time"
)

func main() {
	Path_to_db := "./rtf.db"
	// delete database
	Delete_db(Path_to_db)
	// create database
	Initiate_DB(Path_to_db)
	Open_DB(Path_to_db)
	// create tables
	Create_tables(Path_to_db, List_tables_fields, List_tables_constraints)
	// create some users
	Creates_users(Path_to_db, list_users)
	// create some categories
	Create_categories(Path_to_db, list_categories)
	// create some errands
	Create_first_public_errands(Path_to_db)
	Create_following_public_errands(Path_to_db)
	Create_first_private_errands(Path_to_db)
	Create_reply_private_errands(Path_to_db)
}

func Delete_db(Path_to_db string) {
	fmt.Println("deletting existing DB")
	err := os.Remove(Path_to_db)
	if err != nil {
		fmt.Println(" error during file deletion", err)
		// os.Exit()
	}
	fmt.Println("")
}

func Initiate_DB(Path_to_db string) {
	fmt.Println("initiate DB")
	db_f.Initiate_database(Path_to_db)
	fmt.Println("")
}
func Open_DB(Path_to_db string) {
	fmt.Println("opening DB")
	var err error
	database, err := sql.Open("sqlite3", Path_to_db)
	if err != nil {
		fmt.Println("error opening/creating database")
	}
	defer database.Close() // good practice
	fmt.Println("")
}

func Create_tables(Path_to_db string, list_tables_fields map[string]map[string]string, List_tables_constraints map[string]map[string]string) {
	// create tables
	fmt.Println("creating tables")
	for table_name, types_by_fields := range List_tables_fields {
		db_f.Create_table(Path_to_db, table_name, types_by_fields, List_tables_constraints[table_name])
		fmt.Println("")
	}
	fmt.Println("")
}

var dummy_user_entry_1 map[string]string = map[string]string{
	"Email":         "bop1@gmail.ee",
	"Nickname":      "bob1",
	"Password":      "secret1",
	"session":       "1",
	"First_name":    "",
	"Last_name":     "",
	"Gender":        "",
	"Age":           "",
	"Last_activity": "", // date
}

// "Id":            "INTEGER PRIMARY KEY AUTOINCREMENT",
// "Email":         "TEXT NOT NULL UNIQUE",
// "Nickname":      "TEXT NOT NULL UNIQUE",
// "First_name":    "TEXT",
// "Gender":        "TEXT",
// "Age":           "INTEGER",
// "Last_name":     "TEXT",
// "Password":      "TEXT CHECK (length(password) > 3 )",
// "Session":       "TEXT UNIQUE",
// "Last_activity": "TEXT", // date

var dummy_user_entry_2 map[string]string = map[string]string{
	"Email":         "bop2@gmail.ee",
	"Nickname":      "bob2",
	"Password":      "secret2",
	"session":       "2",
	"First_name":    "",
	"Last_name":     "",
	"Gender":        "",
	"Age":           "",
	"Last_activity": "", // date
}
var dummy_user_entry_3 map[string]string = map[string]string{
	"Email":         "bop3@gmail.ee",
	"Nickname":      "bob3",
	"Password":      "secret3",
	"session":       "3",
	"First_name":    "",
	"Last_name":     "",
	"Gender":        "",
	"Age":           "",
	"Last_activity": "", // date
}
var dummy_user_entry_4 map[string]string = map[string]string{
	"Email":         "bop4@gmail.ee",
	"Nickname":      "bob4",
	"Password":      "secret4",
	"session":       "4",
	"First_name":    "",
	"Last_name":     "",
	"Gender":        "",
	"Age":           "",
	"Last_activity": "", // date
}
var list_users []map[string]string = []map[string]string{
	dummy_user_entry_1,
	dummy_user_entry_2,
	dummy_user_entry_3,
	dummy_user_entry_4,
}

func Creates_users(Path_to_db string, list_users []map[string]string) {
	fmt.Println("inserting users")
	// create users
	for _, user_data := range list_users {
		db_f.Insert_row(Path_to_db, "users", user_data)
	}
	fmt.Println("")
}

func Create_categories(Path_to_db string, list_categories map[string]bool) {
	fmt.Println("creating categories")
	for categ_name := range list_categories {
		fmt.Println("category name : ", categ_name)
		categ_row := New_row("Category")
		fmt.Println("virgin row: ", categ_row)
		categ_row["Category"] = categ_name
		fmt.Println("inserting row filled", categ_row)
		db_f.Insert_row(Path_to_db, "categories", categ_row)
	}
	fmt.Println("")
}

var list_categories map[string]bool = map[string]bool{
	"I like cat":  true,
	"I like porc": true,
	"whatever":    true,
}

func Create_first_public_errands(Path_to_db string) {
	fmt.Println("creating comments")
	fmt.Println("getting all categories")
	categories, err := db_f.Get_all_rows_from_table(Path_to_db, "Categories", "Id")
	fatal_error(err)
	fmt.Println("all categories", categories)
	users, err := db_f.Get_all_rows_from_table(Path_to_db, "Users", "Id")
	fatal_error(err)
	for category_id := range categories {
		for user_id := range users {
			an_errand := New_row("Errand")
			an_errand["Category_id"] = category_id
			an_errand["Sender_id"] = user_id
			an_errand["Date"] = rdm_date()
			an_errand["Text"] = randomstring(9)
			fmt.Println("\ntrying to insert that new errand into Errands", an_errand)
			res, err := db_f.Insert_row(Path_to_db, "Errands", an_errand)
			fmt.Println("result of insertion", res)
			fatal_error(err)
		}
	}
	fmt.Println("")
}

func Create_following_public_errands(Path_to_db string) {
	fmt.Println("creating reply pub errands")
	fmt.Println("getting all errands")
	errands, err := db_f.Get_all_rows_from_table(Path_to_db, "Errands", "Id")
	fatal_error(err)
	fmt.Println("getting all users")
	users, err := db_f.Get_all_rows_from_table(Path_to_db, "Users", "Id")
	fatal_error(err)
	for errand_id, errand_values := range errands {
		for user_id := range users {
			new_errand := New_row("Errand")
			new_errand["Category_id"] = errand_values["Category_id"]
			new_errand["previous_id"] = errand_id
			new_errand["opening_id"] = errand_id
			new_errand["Sender_id"] = user_id
			new_errand["Date"] = rdm_date()
			new_errand["Text"] = randomstring(9)
			_, err := db_f.Insert_row(Path_to_db, "Errands", new_errand)
			fatal_error(err)
		}
	}
	fmt.Println("")
}

func Create_first_private_errands(Path_to_db string) {
	fmt.Println("creating 1st  private erands")
	users, err := db_f.Get_all_rows_from_table(Path_to_db, "Users", "Id")
	fatal_error(err)
	for sender_id := range users {
		an_errand := New_row("Errand")
		an_errand["Sender_id"] = sender_id
		for sendee_id := range users {
			if sendee_id == sender_id {
				continue
			}
			an_errand["Date"] = rdm_date()
			an_errand["Text"] = randomstring(9)
			an_errand["Sendee_id"] = sendee_id
			fmt.Println("inserting new errands as 1st private", an_errand)
			_, err := db_f.Insert_row(Path_to_db, "Errands", an_errand)
			fatal_error(err)
		}
	}
	fmt.Println("")
}

func Create_reply_private_errands(Path_to_db string) {
	fmt.Println("creating reply privte errands")
	errands, err := db_f.Get_all_rows_from_table(Path_to_db, "Errands", "Id")
	fatal_error(err)
	for prev_errand_id, prev_errand_values := range errands {
		if prev_errand_values["Sendee"] == "" {
			continue
		}
		an_errand := New_row("Errand")
		an_errand["Sender_id"] = prev_errand_values["Sendee"]
		an_errand["Sendee_id"] = prev_errand_id
		an_errand["Opening_id"] = prev_errand_id
		an_errand["Date"] = rdm_date()
		an_errand["Text"] = randomstring(9)
		db_f.Insert_row(Path_to_db, "Errands", an_errand)
	}
	fmt.Println("")
}

func fatal_error(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func random_figure(n int) string {
	var letters = "123456789"
	letters = letters[:n]
	return string(letters[rand.Intn(len(letters))])
}

func randomstring(n int) string {
	var letters = "nadwqgeiuo39482 dfwq2387 klrqwjbvg&"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func rdm_date() string {
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

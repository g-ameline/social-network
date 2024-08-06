package main

import (
	"database/sql"

	// "errors"
	"fmt"
	db "github.com/g-ameline/sql_helper"
	"log"
	"os"
	// "strconv"
)

const path_to_db = "../social-network.db"

func main() {
	delete_database_if_exist(path_to_db)
	db.Create_database(path_to_db)
	// check all right
	fatal_error(db.Try_open_close_database(path_to_db))
	// create tables
	create_all_tables(path_to_db, list_tables_fields, list_tables_constraints)
	check_table_exist(path_to_db, "database")
	check_table_exist(path_to_db, "users")
	version_the_database(path_to_db, "", "0.1")
}
func delete_database_if_exist(path_to_database string) {
	fmt.Println("delete table if exist")
	os.Remove(path_to_database)
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

func create_all_tables(Path_to_db string, list_tables_fields map[string]map[string]string, List_tables_constraints map[string]map[string]string) {
	// create tables
	fmt.Println("creating tables")
	for table_name, types_by_fields := range list_tables_fields {
		fmt.Println("table", table_name)
		db.Create_table(Path_to_db, table_name, types_by_fields, List_tables_constraints[table_name])
		fmt.Println("")
	}
	fmt.Println("")
}

var database *sql.DB

const database_driver = "sqlite3"

func check_table_exist(path_to_database, table_name string) {
	res, err := db.Get_all_rows_from_table(path_to_db, table_name, "id")
	fatal_error(err)
	fmt.Println(res)
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

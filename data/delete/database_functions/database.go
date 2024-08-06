package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
)

/* some terminology :
if we take the user entity we have
a table users
with column (fields) like username or email
with also rows (records) like id : 2 , username : jojo , email : john@gmail.com, etc.
id is the key and "jojo" is the value matching a column and a row

note:
unless we need to do operation on queried values, it will be handle as string
*/

var database *sql.DB

var database_driver = "sqlite3"

const comma = " , "

const verbose = true

func Initiate_database(path_to_database string) {
	// check if database not exist
	if _, err := os.Stat(path_to_database); errors.Is(err, os.ErrNotExist) {
		breadcrumb(verbose, "no DB , creating DB")
		Create_database(path_to_database) // if not then we create it
	}
	//try to open DB
	var err error
	database, err = sql.Open(database_driver, path_to_database)
	if err != nil {
		log.Fatalln("error opening/creating database")
	}
	defer database.Close() // good practice
}

func Create_database(path_to_db string) {
	file, err := os.Create(path_to_db)
	defer file.Close()
	if err != nil {
		fmt.Println("could not create databse ", err)
	}
	file.Close()
}

func Create_table(path_to_database string, table_name string, types_by_fields map[string]string, constraints map[string]string) error {
	database, err := sql.Open(database_driver, path_to_database)
	if err != nil {
		log.Fatalln("error opening/creating database")
	}
	defer database.Close() // in case of
	var statement string
	statement = statement_create_table(table_name, types_by_fields, constraints)
	query, err := database.Prepare(statement)
	if_wrong(err, "error creating table "+table_name)
	_, err = query.Exec()
	return err
}

func statement_create_table(table_name string, types_by_fields map[string]string, constraints map[string]string) string {
	var statement string
	statement += "CREATE TABLE IF NOT EXISTS " + table_name
	statement += " ("
	for field, data_type := range types_by_fields {
		statement += field + " " + data_type + comma
	}
	// add constraints to statement
	for constraint_name, link := range constraints {
		statement += constraint_name + " " + link + comma
	}
	statement = statement[:len(statement)-len(comma)] // remove last comma
	statement += ") "

	breadcrumb(verbose, "create table statement", statement)
	return statement
}

func Insert_row(path_to_database string, table_name string, values_by_fields map[string]string) (string, error) {
	database, err := sql.Open(database_driver, path_to_database)
	defer database.Close() // good practice
	single_quote_text_values(values_by_fields)
	statement_new_row := statement_insert_row(table_name, values_by_fields)
	breadcrumb(verbose, "statement:", statement_new_row)
	result, err := database.Exec(statement_new_row)
	id_int64, err := sql.Result.LastInsertId(result)
	id_string := strconv.FormatInt(id_int64, 10)
	return id_string, err
}

func statement_insert_row(table_name string, values_by_fields map[string]string) string {
	var statement string
	var fields_part, values_part string
	statement += "INSERT INTO " + table_name
	// column first
	fields_part += "( "
	values_part += "( "
	for field, value := range values_by_fields {
		fields_part += field + comma
		values_part += value + comma
	}
	if len(fields_part) > len(comma) && len(values_part) > len(comma) {
		fields_part = fields_part[:len(fields_part)-len(comma)] // remove last comma
		values_part = values_part[:len(values_part)-len(comma)] // remove last comma
	}
	fields_part += ") "
	values_part += ") "
	// then values
	statement += fields_part
	statement += "VALUES "
	statement += values_part
	return statement
}

func Delete_records(path_to_database string, table_name, field, value string) error {
	database, err := sql.Open(database_driver, path_to_database)
	defer database.Close() // good practice
	value = single_quote_text(value)
	single_quote_text(value)
	statement := statement_delete_rows(table_name, field, value)
	breadcrumb(verbose, "deletion statement:", statement)
	// f := func() (sql.Result, error) { return m_database.Value.Exec(statement) }
	// m_result := m.Ligate[*sql.DB, sql.Result](m_database, f)
	_, err = database.Exec(statement)
	return err
}

func statement_delete_rows(table_name, field, value string) string {
	var statement string
	statement += "DELETE FROM " + table_name + " WHERE " + field + " = " + value
	return statement
}

func Get_one_row(path_to_database string, table_name, field_key, value_key string) (map[string]string, error) {
	database, err := sql.Open(database_driver, path_to_database)
	defer database.Close() // good practice
	var query string
	value_key = single_quote_text(value_key)
	query = query_get_one_row(table_name, field_key, value_key)
	breadcrumb(verbose, "query:", query)
	rows, err := database.Query(query)
	defer rows.Close()
	breadcrumb(verbose, "rows from  query :", rows)
	// var fields []string
	fields, err3 := rows.Columns()
	if err3 != nil { // catch here
		return nil, err
	}
	values := make([]string, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range values {
		pointers_v[i] = &values[i]
	}
	var row_as_map map[string]string
	if rows.Next() {
		err := rows.Scan(pointers_v...)
		breadcrumb(verbose, "fields", fields)
		breadcrumb(verbose, "values", values)
		row_as_map, err = zip_map(fields, values)
		return row_as_map, err
	}
	return row_as_map, errors.New("error during row scanning")
}

func Get_all_rows_from_table(path_to_database string, table_name, sorting_field string) (map[string]map[string]string, error) {
	table_as_map := make(map[string]map[string]string)
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error opening/creating database")
	defer database.Close() // good practice
	var query string
	query = query_get_all_rows(table_name, sorting_field)
	breadcrumb(verbose, "query for getting all rows", query)
	rows, err := database.Query(query)
	if_wrong(err, "error when fetching all rows")
	defer rows.Close()
	var fields []string
	fields, err = rows.Columns()
	if_wrong(err, "issue when fetching columns names")
	breadcrumb(verbose, "fields from table", fields)
	for rows.Next() {
		var row_as_map map[string]string
		values := make([]string, len(fields))
		pointers_v := make([]any, len(fields))
		for i := range values {
			pointers_v[i] = &values[i]
		}
		err = rows.Scan(pointers_v...)
		if_wrong(err, "error during scanning of a row"+" "+table_name+" "+sorting_field)
		breadcrumb(verbose, "fields", fields)
		breadcrumb(verbose, "values", values)
		row_as_map, _ = zip_map(fields, values)
		table_as_map[row_as_map[sorting_field]] = row_as_map
	}
	return table_as_map, err
}

func Get_all_rows_sorted(path_to_database string, table_name, sorting_field string) ([]map[string]string, error) {
	var table_as_slice []map[string]string
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error opening/creating database")
	defer database.Close() // good practice
	var query string
	query = query_get_all_rows(table_name, sorting_field)
	breadcrumb(verbose, "query for getting all rows", query)
	rows, err := database.Query(query)
	if_wrong(err, "error when fetching all rows")
	defer rows.Close()
	var fields []string
	fields, err = rows.Columns()
	if_wrong(err, "issue when fetching columns names")
	breadcrumb(verbose, "fields from table", fields)
	for rows.Next() {
		var row_as_map map[string]string
		values := make([]string, len(fields))
		pointers_v := make([]any, len(fields))
		for i := range values {
			pointers_v[i] = &values[i]
		}
		err = rows.Scan(pointers_v...)
		if_wrong(err, "error during scanning of a row"+" "+table_name+" "+sorting_field)
		breadcrumb(verbose, "fields", fields)
		breadcrumb(verbose, "values", values)
		row_as_map, _ = zip_map(fields, values)
		table_as_slice = append(table_as_slice, row_as_map)
	}
	return table_as_slice, err
}

func query_get_all_rows(table_name, sorting_field string) string {
	var query string
	query += "SELECT * FROM " + table_name + " ORDER BY " + sorting_field
	return query
}

func query_get_row_two_cond(table_name, field, value, other_field, other_value string) string {
	var query string
	query += "SELECT * FROM " + table_name + " WHERE " + field + "=" + value + " AND " + other_field + "=" + other_value
	return query
}

func query_get_one_row(table_name, field, value string) string {
	var query string
	query += "SELECT * FROM " + table_name + " WHERE " + field + " = " + value
	return query
}

func Get_id(path_to_database string, table_name, field_key, value_key string) (string, error) {
	var row_as_map map[string]string
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error opening database")
	defer database.Close() // good practice
	var query string
	value_key = single_quote_text(value_key)
	query = query_get_one_row(table_name, field_key, value_key)
	rows, err := database.Query(query)
	if_wrong(err, "error while fetching row")
	defer rows.Close()
	var fields []string
	fields, err = rows.Columns()
	if_wrong(err, "error while reading row")
	values := make([]string, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range values {
		pointers_v[i] = &values[i]
	}
	rows.Next()
	err = rows.Scan(pointers_v...)
	if_wrong(err, "error during scanning of single row to get Id"+" "+table_name+" "+field_key+" "+value_key)
	row_as_map, _ = zip_map(fields, values)
	return row_as_map["Id"], err
}

func Get_id_two_cond(path_to_database string, table_name, field_key, value_key, other_field, other_value string) (string, error) {
	var row_as_map map[string]string
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error opening database")
	defer database.Close() // good practice
	var query string
	value_key = single_quote_text(value_key)
	query = query_get_row_two_cond(table_name, field_key, value_key, other_field, other_value)
	rows, err := database.Query(query)
	if_wrong(err, "error while fetching row")
	defer rows.Close()
	var fields []string
	fields, err = rows.Columns()
	if_wrong(err, "error while reading row")
	values := make([]string, len(fields))
	pointers_v := make([]any, len(fields))
	for i := range values {
		pointers_v[i] = &values[i]
	}
	rows.Next()
	err = rows.Scan(pointers_v...)
	if_wrong(err, "error during scanning of single row to get Id"+" "+table_name+" "+field_key+" "+value_key)
	row_as_map, _ = zip_map(fields, values)
	return row_as_map["Id"], err
}

func query_ids_from_table(table_name string) string {
	var query string
	query += "SELECT Id FROM " + table_name
	return query
}

func Check_if_record(path_to_database string, table, field_1, value_1, field_2, value_2 string) (bool, error) { // for likes or dislikes
	query := "SELECT 1 FROM " + table + " WHERE " + field_1 + "=" + value_1 + " AND " + field_2 + "=" + value_2
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error accessing database")
	defer database.Close() // good practice
	rows, err := database.Query(query)
	if_wrong(err, "error while querying all row/record")
	defer rows.Close()
	// fmt.Println("res", rows)
	// fmt.Println("row.next", rows.Next())
	return rows.Next(), err
}

func Is_in_database(path_to_database string, table, field_1, value_1, field_2, value_2 string) (bool, error) {
	query := "SELECT 1 FROM " + table + " WHERE " + field_1 + "=" + value_1 + " AND " + field_2 + "=" + value_2
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error accessing database")
	defer database.Close() // good practice
	rows, err := database.Query(query)
	if_wrong(err, "error while querying all row/record")
	defer rows.Close()
	// fmt.Println("res", rows)
	// fmt.Println("row.next", rows.Next())
	return rows.Next(), err
}

func statement_rows(table_name, field, value string) string {
	var statement string
	statement += "DELETE FROM " + table_name + " WHERE " + field + " = " + value
	return statement
}

func Count_all_rows(path_to_database string, table_name string) (int, error) {
	database, err := sql.Open(database_driver, path_to_database)
	if_wrong(err, "error accessing database")
	defer database.Close() // good practice
	query := query_ids_from_table(table_name)
	breadcrumb(verbose, "counting statement:", query)
	rows, err := database.Query(query)
	if_wrong(err, "error while querying all row/record")
	defer rows.Close()
	var counter int
	for rows.Next() {
		counter++
	}
	return counter, err
}

func zip_map(keys_slice []string, values_slice []string) (map[string]string, error) {
	if len(keys_slice) != len(values_slice) {
		return nil, fmt.Errorf("different length of slices when zipping it")
	}
	if len(keys_slice) == 0 {
		return nil, fmt.Errorf("zero length slice of slices when zipping it")
	}
	keys_values := make(map[string]string)
	for i := 0; i < len(keys_slice); i++ {
		keys_values[keys_slice[i]] = values_slice[i]
	}
	return keys_values, nil
}

func breadcrumb(v bool, helpers ...any) {
	if v {
		for _, h := range helpers {
			fmt.Print(h, " ")
		}
		fmt.Print("\n")
	}
}

func if_wrong(err error, message string) {
	if err != nil {
		println(message, err.Error())
	}
}
func is_wrong(err error, message string) bool {
	if err != nil {
		println(message, err.Error())
		return true
	}
	return false
}

func single_quote_text_values(values_by_fields map[string]string) {
	for field, value := range values_by_fields {
		if value != `''` {
			values_by_fields[field] = single_quote_text(value)
		}
	}
}

func single_quote_text(value string) string {
	_, err := strconv.Atoi(value) // if can be inferred to an int then it is an int
	if err != nil {
		return "'" + value + "'"
	} // what about boolean ?

	return value
}

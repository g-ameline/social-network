package main

import (
	// "database/sql"
	// "strings"

	// "errors"
	// "fmt"
	// col "github.com/g-ameline/colors"
	"fmt"

	db "github.com/g-ameline/sql_helper"
	// "log"
	// "math/rand"
	// "strconv"
	// "time"
)

const path_to_db = "../social-network.db"

func main() {
	fmt.Println(db.Try_open_close_database(path_to_db + "2"))
}

package main

import "log"

func New_row(entity string) map[string]string {
	var new_row map[string]string = map[string]string{}
	switch entity {
	case "User":
		for k := range Users_types_by_fields {
			new_row[k] = ""
		}
	case "Category":
		for k := range Categories_types_by_fields {
			new_row[k] = ""
		}
		println("new row", new_row)
	case "Errand":
		for k := range Errands_types_by_fields {
			new_row[k] = ""
		}
	default:
		log.Fatalln("entity (table name) not defined", entity)
	}
	delete(new_row, "Id")
	return new_row
}

var List_tables_fields map[string]map[string]string = map[string]map[string]string{
	"Users":      Users_types_by_fields,
	"Categories": Categories_types_by_fields,
	"Errands":    Errands_types_by_fields,
}

var List_tables_constraints map[string]map[string]string = map[string]map[string]string{
	"Users":      nil,
	"Categories": nil,
	"Errands":    Errands_constraints,
}

var Users_types_by_fields map[string]string = map[string]string{
	"Id":            "INTEGER PRIMARY KEY AUTOINCREMENT",
	"Email":         "TEXT NOT NULL UNIQUE",
	"Nickname":      "TEXT NOT NULL UNIQUE",
	"First_name":    "TEXT",
	"Gender":        "TEXT",
	"Age":           "INTEGER",
	"Last_name":     "TEXT",
	"Password":      "TEXT CHECK (length(password) > 3 )",
	"Session":       "TEXT UNIQUE",
	"Last_activity": "TEXT", // date
}

var Categories_types_by_fields map[string]string = map[string]string{
	"Id":       "INTEGER PRIMARY KEY AUTOINCREMENT",
	"Category": "TEXT NOT NULL UNIQUE",
}
var Errands_types_by_fields map[string]string = map[string]string{
	"Id":          "INTEGER PRIMARY KEY AUTOINCREMENT",
	"Category_id": "INT",
	"Opening_id":  "INTEGER",
	"Previous_id": "INTEGER",
	"Sender_id":   "INTEGER",
	"Sendee_id":   "INTEGER",
	"Date":        "TEXT",
	"Text":        "TEXT NOT NULL",
}

var Errands_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_Sender":   "FOREIGN KEY (Sender_id) REFERENCES Users(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_Sendee":   "FOREIGN KEY (Sendee_id) REFERENCES Users(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_Category": "FOREIGN KEY (Category_id) REFERENCES Categories(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_Opening":  "FOREIGN KEY (Opening_id) REFERENCES Errands(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_Previous": "FOREIGN KEY (Previous_id) REFERENCES Errands(Id) ON DELETE SET NULL",
}

package main

var list_tables_fields map[string]map[string]string = map[string]map[string]string{
	"database":              database_types_by_fields,
	"users":                 users_types_by_fields,
	"followings":            followings_types_by_fields,
	"groups":                groups_types_by_fields,
	"joinings":              joinings_types_by_fields,
	"events":                events_types_by_fields,
	"attendings":            attendings_types_by_fields,
	"public_addressings":    public_addressings_types_by_fields,
	"private_addressings":   private_addressings_types_by_fields,
	"exclusive_addressings": exclusive_addressings_types_by_fields,
	"group_addressings":     group_addressings_types_by_fields,
	"notes":                 notes_types_by_fields,
}

var list_tables_constraints map[string]map[string]string = map[string]map[string]string{
	"database":              nil,
	"users":                 nil,
	"groups":                groups_constraints,
	"followings":            followings_constraints,
	"joinings":              joinings_constraints,
	"events":                events_constraints,
	"attendings":            attendings_constraints,
	"public_addressings":    public_addressings_constraints,
	"private_addressings":   private_addressings_constraints,
	"exclusive_addressings": exclusive_addressings_constraints,
	"group_addressings":     group_addressings_constraints,
	"notes":                 notes_constraints,
}

var database_types_by_fields map[string]string = map[string]string{
	"id":      "INTEGER PRIMARY KEY AUTOINCREMENT",
	"version": "TEXT NOT NULL UNIQUE",
	"date":    "DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP",
}

var users_types_by_fields map[string]string = map[string]string{
	"id":         "INTEGER PRIMARY KEY AUTOINCREMENT",
	"email":      "TEXT NOT NULL UNIQUE",
	"password":   "BLOB NOT NULL", // CHECK" (length(password) > 3 )",
	"first_name": "TEXT NOT NULL",
	"last_name":  "TEXT NOT NULL",
	"birth":      "TEXT NOT NULL",
	"nickname":   "TEXT",
	"avatar":     "TEXT UNIQUE",
	"about":      "TEXT",
	"private":    "INT DEFAULT 0", // 0 = public or 1 = private
}

var followings_types_by_fields map[string]string = map[string]string{
	"id":          "INTEGER PRIMARY KEY AUTOINCREMENT",
	"follower_id": "INT NOT NULL",
	"followee_id": "INT NOT NULL",
	"approval":    "INT NOT NULL", //-1 refused | 0 pending | +1 accepted
	"pair_id":     "TEXT UNIQUE GENERATED ALWAYS AS (follower_id || '-' || followee_id) STORED",
}

var followings_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_follower":       "FOREIGN KEY (follower_id) REFERENCES users(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_followee":       "FOREIGN KEY (followee_id) REFERENCES users(Id) ON DELETE SET NULL",
	"CONSTRAINT no_self_following": "CHECK  (followee_id <> follower_id)",
}

var groups_types_by_fields map[string]string = map[string]string{
	"id":          "INTEGER PRIMARY KEY AUTOINCREMENT",
	"creator_id":  "INT NOT NULL",
	"title":       "TEXT NOT NULL UNIQUE",
	"description": "TEXT NOT NULL",
}
var groups_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_creator": "FOREIGN KEY (creator_id) REFERENCES users(Id) ON DELETE SET NULL",
}

var joinings_types_by_fields map[string]string = map[string]string{
	"id":               "INTEGER PRIMARY KEY AUTOINCREMENT",
	"joiner_id":        "INT NOT NULL",
	"group_id":         "INT NOT NULL",
	"approval_creator": "INT NOT NULL DEFAULT 0", // -1 refused | 0 pending | +1 approved
	"approval_joiner":  "INT NOT NULL DEFAULT 0", // -1 refused | 0 pending | +1 approved
	"pair_id":          "INT UNIQUE GENERATED ALWAYS AS (group_id || '-' || joiner_id) STORED",
}
var joinings_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_joiner": "FOREIGN KEY (joiner_id) REFERENCES users(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_group":  "FOREIGN KEY (group_id) REFERENCES groups(Id) ON DELETE SET NULL",
}

var events_types_by_fields map[string]string = map[string]string{
	"id":          "INTEGER PRIMARY KEY AUTOINCREMENT",
	"group_id":    "INT NOT NULL",
	"title":       "TEXT NOT NULL",
	"description": "TEXT",
	"date":        "TEXT",
}
var events_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_group": "FOREIGN KEY (group_id) REFERENCES groups(Id) ON DELETE SET NULL",
}

var attendings_types_by_fields map[string]string = map[string]string{
	"id":          "INTEGER PRIMARY KEY AUTOINCREMENT",
	"attender_id": "INT NOT NULL",
	"event_id":    "INT NOT NULL",
	"coming":      "INT NOT NULL", // -1 no | 0 nsp | +1 yes
}
var attendings_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_attender": "FOREIGN KEY (attender_id) REFERENCES users(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_event":    "FOREIGN KEY (event_id) REFERENCES events(Id) ON DELETE SET NULL",
}

var notes_types_by_fields map[string]string = map[string]string{
	"id":             "INTEGER PRIMARY KEY AUTOINCREMENT",
	"author_id":      "INT NOT NULL",
	"predecessor_id": "INT",
	"date":           "DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP",
	"text":           "TEXT NOT NULL",
	"picture":        "TEXT",
}

var notes_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_author":      "FOREIGN KEY (author_id) REFERENCES users(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_predecessor": "FOREIGN KEY (predecessor_id) REFERENCES notes(Id) ON DELETE SET NULL",
}
var public_addressings_types_by_fields map[string]string = map[string]string{
	"id":      "INTEGER PRIMARY KEY AUTOINCREMENT",
	"note_id": "INT UNIQUE NOT NULL",
}
var public_addressings_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_note": "FOREIGN KEY (note_id) REFERENCES notes(Id) ON DELETE SET NULL",
}
var private_addressings_types_by_fields map[string]string = map[string]string{
	"id":      "INTEGER PRIMARY KEY AUTOINCREMENT",
	"note_id": "INT UNIQUE NOT NULL",
}
var private_addressings_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_note": "FOREIGN KEY (note_id) REFERENCES notes(Id) ON DELETE SET NULL",
}
var exclusive_addressings_types_by_fields map[string]string = map[string]string{
	"id":           "INTEGER PRIMARY KEY AUTOINCREMENT",
	"note_id":      "INT NOT NULL",
	"addressee_id": "INT NOT NULL",
}
var exclusive_addressings_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_note":      "FOREIGN KEY (note_id) REFERENCES notes(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_addressee": "FOREIGN KEY (addressee_id) REFERENCES users(Id) ON DELETE SET NULL",
	"CONSTRAINT uniqueness":   "UNIQUE (addressee_id, note_id)",
}
var group_addressings_types_by_fields map[string]string = map[string]string{
	"id":       "INTEGER PRIMARY KEY AUTOINCREMENT",
	"note_id":  "INT UNIQUE NOT NULL",
	"group_id": "INT NOT NULL",
}

var group_addressings_constraints map[string]string = map[string]string{
	"CONSTRAINT fk_note":  "FOREIGN KEY (note_id) REFERENCES notes(Id) ON DELETE SET NULL",
	"CONSTRAINT fk_group": "FOREIGN KEY (group_id) REFERENCES groups(Id) ON DELETE SET NULL",
}

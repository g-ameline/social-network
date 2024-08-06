package main

import (
	"fmt"
	col "github.com/g-ameline/colors"
	db "github.com/g-ameline/sql_helper"
	"testing"
)

func Test_server(t *testing.T) {
	list_all_groups(path_to_db)
	// list_all_public_notes(path_to_db)
	// list_all_private_notes(path_to_db)
	// list_all_exclu_notes(path_to_db)
	// list_all_chat_notes(path_to_db)
}
func list_all_groups(path_to_database string) {
	fmt.Print(col.All["BLUE"], col.All["black"])
	fmt.Println("groups")
	// users_rows, err := db.Get_all_rows_from_table(path_to_database, "users", "id")
	// talk("failed to get users", err)
	groups_rows, err := db.Get_all_rows_from_table(path_to_database, "groups", "id")
	talk("failed to get addressings", err)
	for _, group_row := range groups_rows {
		fmt.Println(group_row)
	}
}
func list_all_public_notes(path_to_database string) {
	fmt.Print(col.All["BLUE"], col.All["black"])
	fmt.Println("Public NOTE")
	// users_rows, err := db.Get_all_rows_from_table(path_to_database, "users", "id")
	// talk("failed to get users", err)
	public_addressings_rows, err := db.Get_all_rows_from_table(path_to_database, "public_addressings", "id")
	talk("failed to get addressings", err)
	notes_rows, err := db.Get_all_rows_from_table(path_to_database, "notes", "id")
	talk("failed to get notes", err)
	for _, public_addressing_row := range public_addressings_rows {
		public_note_id := public_addressing_row["note_id"]
		notes_rows[public_note_id] = notes_rows[public_note_id]
	}
	fmt.Println("public notes", len(notes_rows))
}
func list_all_private_notes(path_to_database string) {
	fmt.Print(col.All["WHITE"], col.All["blue"])
	fmt.Println("private NOTE")
	users_rows, err := db.Get_all_rows_from_table(path_to_database, "users", "id")
	talk("failed to get users", err)
	private_addressings_rows, err := db.Get_all_rows_from_table(path_to_database, "private_addressings", "id")
	talk("failed to get addressings", err)
	followings_rows, err := db.Get_all_rows_from_table(path_to_database, "followings", "id")
	talk("failed to get followings", err)
	notes_rows, err := db.Get_all_rows_from_table(path_to_database, "notes", "id")
	talk("failed to get notes", err)
	for user_id, _ := range users_rows {
		acquaintances_ids := map[string]bool{}
		for _, following_row := range followings_rows {
			followee_id := following_row["followee_id"]
			follower_id := following_row["follower_id"]
			if following_row["approval"] != "1" {
				continue
			}
			if followee_id == user_id {
				acquaintances_ids[follower_id] = true
			}
			if follower_id == user_id {
				acquaintances_ids[followee_id] = true
			}
		}
		// get all addressing that are private and addresser is an acquaintance
		private_notes_rows := map[string]map[string]string{}
		for _, private_addressing_row := range private_addressings_rows {
			private_note_id := private_addressing_row["note_id"]
			private_notes_rows[private_note_id] = notes_rows[private_note_id]
		}
		private_notes_rows_from_acquaintances := map[string]map[string]string{}
		for note_id, private_note_row := range private_notes_rows {
			author_id := private_note_row["author_id"]
			if !acquaintances_ids[author_id] {
				continue
			}
			private_notes_rows_from_acquaintances[note_id] = private_note_row
		}
		fmt.Println("private notes visible by user", user_id, len(private_notes_rows_from_acquaintances))
	}
}
func list_all_exclu_notes(path_to_database string) {
	fmt.Print(col.All["BLUE"], col.All["white"])
	fmt.Println("exclu NOTE")
	users_rows, err := db.Get_all_rows_from_table(path_to_database, "users", "id")
	talk("failed to get users", err)
	exclu_addressings_rows, err := db.Get_all_rows_from_table(path_to_database, "exclusive_addressings", "id")
	talk("failed to get addressings", err)
	notes_rows, err := db.Get_all_rows_from_table(path_to_database, "notes", "id")
	talk("failed to get notes", err)
	fmt.Println("that many exclu addressings ", len(exclu_addressings_rows))
	for user_id, _ := range users_rows {
		exclu_notes := map[string]map[string]string{}
		for _, exclu_addressing_row := range exclu_addressings_rows {
			addressee_id := exclu_addressing_row["addressee_id"]
			if addressee_id != user_id {
				continue
			}
			exclu_note_id := exclu_addressing_row["note_id"]
			exclu_note_row := notes_rows[exclu_note_id]
			exclu_notes[exclu_note_id] = exclu_note_row
		}
		fmt.Println("that many exclu note destiend to the user", user_id, len(exclu_notes))
	}
}
func list_all_chat_notes(path_to_database string) {
	fmt.Print(col.All["BLACK"], col.All["green"])
	fmt.Println("chat NOTE")
	users_rows, err := db.Get_all_rows_from_table(path_to_database, "users", "id")
	talk("failed to get users", err)
	chat_addressings_rows, err := db.Get_all_rows_from_table(path_to_database, "chat_addressings", "id")
	talk("failed to get addressings", err)
	notes_rows, err := db.Get_all_rows_from_table(path_to_database, "notes", "id")
	talk("failed to get notes", err)
	fmt.Println("that many chat addressings ", len(chat_addressings_rows))
	for user_id, _ := range users_rows {
		chat_notes := map[string]map[string]string{}
		for _, chat_addressing_row := range chat_addressings_rows {
			addresser_id := chat_addressing_row["addresser_id"]
			addressee_id := chat_addressing_row["addressee_id"]
			if addressee_id != user_id && addresser_id != user_id {
				continue
			}
			chat_note_id := chat_addressing_row["note_id"]
			chat_note_row := notes_rows[chat_note_id]
			chat_notes[chat_note_id] = chat_note_row
		}
		fmt.Println("that many chat note destiend to the user", user_id, len(chat_notes))
	}
}

// func list_all_group_notes(path_to_database string) {
// 	fmt.Print(col.All["GREEN"], col.All["white"])
// 	fmt.Println("group NOTE")
// 	users_rows, err := db.Get_all_rows_from_table(path_to_database, "users", "id")
// 	talk("failed to get users", err)
// 	joinings_rows, err := db.Get_all_rows_from_table(path_to_database, "joinings", "id")
// 	talk("failed to get joinings", err)
// 	group_addressings_rows, err := db.Get_all_rows_from_table(path_to_database, "group_addressings", "id")
// 	talk("failed to get addressings", err)
// 	notes_rows, err := db.Get_all_rows_from_table(path_to_database, "notes", "id")
// 	talk("failed to get notes", err)
// 	fmt.Println("that many group addressings ", len(group_addressings_rows))
// 	for _, joining_row := range joinings_rows {
// 	if joining_row["approval_creator"] != "1" && joining_row["approval_joiner"] != "1" {
// 		continue
// 	}
// 	addresser_id := joining_row["joinee_id"]
// 	group_id := joining_row["group_id"]

// 	for user_id, _ := range users_rows {
// 		group_notes := map[string]map[string]string{}
// 		for _, group_addressing_row := range group_addressings_rows {
// 			addresser_id := group_addressing_row["addresser_id"]
// 			addressee_id := group_addressing_row["addressee_id"]
// 			if addressee_id != user_id && addresser_id != user_id {
// 				continue
// 			}
// 			group_note_id := group_addressing_row["note_id"]
// 			group_note_row := notes_rows[group_note_id]
// 			group_notes[group_note_id] = group_note_row
// 		}
// 		fmt.Println("that many group note destiend to the user", user_id, len(group_notes))
// 	}
// }

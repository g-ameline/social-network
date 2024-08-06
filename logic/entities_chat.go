package main

import (
	"fmt"
	// m "github.com/g-ameline/maybe"
	// ws "github.com/gorilla/websocket"
	// http "net/http"
)

// --------- CHATS
func get_user_chat(user_entity_A, user_entity_B user) chat {
	user_id_A := user_entity_A.id()
	user_id_B := user_entity_B.id()
	_, ok_A := entities["user_chats"].(map[string]map[string]chat)[user_id_A]
	_, ok_B := entities["user_chats"].(map[string]map[string]chat)[user_id_B]
	if !ok_A && !ok_B {
		println("adding  map[string]chat to fist is entry")
		entities["user_chats"].(map[string]map[string]chat)[user_id_A] = map[string]chat{}
		entities["user_chats"].(map[string]map[string]chat)[user_id_B] = map[string]chat{}
	}
	chat_A_B, ok_A_B := entities["user_chats"].(map[string]map[string]chat)[user_id_A][user_id_B]
	chat_B_A, ok_B_A := entities["user_chats"].(map[string]map[string]chat)[user_id_B][user_id_A]
	if ok_A_B && ok_B_A { //&& &chat_A_B == &chat_B_A {
		fmt.Println("found it", chat_A_B, chat_B_A)
		return chat_A_B
	}
	// if ok_A_B && ok_B_A {
	// 	fmt.Println("prob it", &chat_A_B, &chat_B_A)
	// 	fmt.Println("equal ?", &chat_A_B == &chat_B_A)
	// 	panic("issue")
	// }
	if ok_A_B != ok_B_A {
		panic("asymetry")
	}
	println("adding new chat{} to each entries")
	new_user_chat := chat{} // list of notes[date or author or ocntent]
	entities["user_chats"].(map[string]map[string]chat)[user_id_A][user_id_B] = new_user_chat
	entities["user_chats"].(map[string]map[string]chat)[user_id_B][user_id_A] = new_user_chat
	return new_user_chat
	// return get_user_chat(user_entity_B, user_entity_A)
}

func get_group_chat(group_id string) chat {
	if chat, ok := entities["group_chats"].(map[string]chat)[group_id]; ok {
		return chat
	}
	entities["group_chats"].(map[string]chat)[group_id] = chat{}
	return entities["group_chats"].(map[string]chat)[group_id]
}
func (the_chat *chat) add(author user, text string) *chat {
	missive := fmt.Sprint(author.email(), " : ", text)
	*the_chat = append(*the_chat, missive)
	return the_chat
}

func (the_group group) add_missive_to_chat(author user, text string) chat {
	missive := fmt.Sprint(author.email(), " : ", text)
	the_chat := entities["group_chats"].(map[string]chat)[the_group.id()]
	the_chat = append(the_chat, missive)
	entities["group_chats"].(map[string]chat)[the_group.id()] = the_chat
	return the_chat
}

func (the_user user) add_missive_to_chat(a_sendee user, text string) chat {
	missive := fmt.Sprint(the_user.email(), " : ", text)
	the_chat := get_user_chat(the_user, a_sendee)
	the_chat = append(the_chat, missive)
	entities["user_chats"].(map[string]map[string]chat)[the_user.id()][a_sendee.id()] = the_chat
	entities["user_chats"].(map[string]map[string]chat)[a_sendee.id()][the_user.id()] = the_chat
	return the_chat
}

func (the_group group) chat() chat {
	if chat, ok := entities["group_chats"].(map[string]chat)[the_group.id()]; ok {
		return chat
	}
	group_chat_entity := chat{}
	entities["group_chats"].(map[string]chat)[the_group.id()] = group_chat_entity
	return group_chat_entity
}

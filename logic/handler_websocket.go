package main

import (
	"fmt"
	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	ws "github.com/gorilla/websocket"
	"net/http"
)

var socket_to_last_data = csa{}

func handler_socket(responder http.ResponseWriter, request *http.Request) {
	c.Green("treating that request in ws upgrader handler", request.URL.Path)
	session, err := session_from_cookie_from_request(request)
	user_id, err := m.If_nil_do[string](err,
		func() (string, error) { return session_valid(session) })
	m.If_nil_must[nvm](err, func() { update_almanac(session, user_id) })
	m.If_error_do[nvm](err, func() { fmt.Fprintf(responder, return_home_link_node().Inline()) })
	socket, err := m.If_nil_do[*ws.Conn](err,
		func() (*ws.Conn, error) { return upgrade_request(responder, request) })
	defer socket.Close()
	m.If_nil_must[nvm](err, func() {
		socket.SetReadLimit(999999)
		id_to_display[user_id] = socket
		fmt.Println("\nID_TO_DISPLAY", id_to_display)
	})
	// m.If_nil_must[nvm](err, func() {
	defer func() {
		fmt.Println("\nID_TO_DISPLAY before delete", id_to_display)
		if id_to_display[user_id] == socket {
			delete(id_to_display, user_id)
		}
		fmt.Println("\nID_TO_DISPLAY after delete", id_to_display)
	}()
	// })
	m.Must(err, "error during websocket handling")
	for {
		fmt.Println("websocket loop")
		c.Cyan("   processing socket messages")
		c.Cyan("      user_id", user_id)
		_, received_raw_message, err := socket.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			println("break")
			fmt.Println(err)
			break
		}
		data, err := m.If_nil_do[sa](err,
			func() (sa, error) { return From_raw[sa](received_raw_message) })
		data["user_id"] = user_id
		data["socket"] = socket
		c.Cyan("      data from message")
		c.Print_map(data)
		process_occurrence(data)
		update_ephemeris(session)
	}

	select {}
}

func process_occurrence(data sa) {
	switch data["occurrence"].(string) {
	case "dummy_occurrence":
		dummy_occurrence_processor(data)
		// --------------- USER OCCURRENCE  ---------------
	case "user_list":
		render_user_list(data)
	case "user_profile":
		render_user_profile(data)
	case "user_profile_private":
		process_user_profile_private(data)
	case "user_profile_public":
		process_user_profile_public(data)
	case "user_chat":
		render_user_chat(data)
	case "user_chat_missive":
		process_user_chat_missive(data)
	// --------------- USER OCCURRENCE features ---------------
	case "user_follow":
		process_user_follow(data)
	case "user_unfollow":
		process_user_unfollow(data)
	case "user_solicit":
		process_user_solicit(data)
	case "user_unsolicit":
		process_user_unsolicit(data)
	case "user_concede":
		process_user_concede(data)
	case "user_forshake":
		process_user_forshake(data)
		// --------------- user---------------
	case "note_user_posts":
		render_note_user_posts(data)
	case "note_user_new_post":
		render_note_user_new_post(data)
	case "note_comment":
		process_note_comment(data)
	case "note_post_public_create":
		process_note_post_public_create(data)
	case "note_post_private_create":
		process_note_post_private_create(data)
	case "note_post_exclusive_create":
		process_note_post_exclusive_create(data)
	case "note_divulge":
		process_note_confide(data)
	case "note_hide":
		process_note_hide(data)
		// --------------- group ---------------
	case "note_post_group_create":
		process_note_post_group_create(data)
	case "note_group_posts":
		render_group_posts(data)
		// --------------- GROUP OCCURRENCES listing ---------------
	case "group_list":
		render_group_list(data)
	case "group_list_created":
		render_group_list_created(data)
	case "group_list_joined":
		render_group_list_joined(data)
	case "group_list_invited":
		render_group_list_invited(data)
	case "group_list_applied":
		render_group_list_applied(data)
	case "group_list_applicable":
		render_group_list_applicable(data)
		// --------------- GROUP OCCURRENCES create ---------------
	case "group_new":
		render_group_new(data)
	case "group_create":
		process_group_create(data)
		// --------------- GROUP OCCURRENCES FEATURES render ---------------
	case "group_features":
		render_group_features(data)
		// --------------- GROUP OCCURRENCES sub features ---------------
	case "group_review_applications":
		render_group_review_applications(data)
	case "group_invite_outsiders":
		render_group_invite_outsiders(data)
	case "group_events":
		render_group_events(data)
	case "group_chat":
		render_group_chat(data)
		// --------------- GROUP OCCURRENCES manipulation ---------------
	case "group_apply":
		process_group_apply(data)
	case "group_unapply":
		process_group_unapply(data)
	case "group_admit":
		process_group_admit(data)
	case "group_reject":
		process_group_reject(data)
	case "group_invite":
		process_group_invite(data)
	case "group_assent":
		process_group_assent(data)
	case "group_decline":
		process_group_decline(data)
	case "group_chat_missive":
		process_group_chat_missive(data)
		// --------------- EVENTS OCCURRENCE---------------
	case "event_create":
		process_event_create(data)
	case "event_going":
		process_event_going(data)
	case "event_not_going":
		process_event_not_going(data)
	default:
		c.Red(data["occurrence"].(string))
		panic(fmt.Errorf("there is no processor for that occurrence message"))
	}
}

func dummy_occurrence_processor(data sa) {
	c.Magenta("dummy occurence processing")
	socket := data["socket"].(*ws.Conn)
	c.Print_map(data)
	dummy_content, err := To_raw(data)
	m.Must(err, "voila")
	dummy_htmx_fragment := append([]byte(`<div>`), dummy_content...)
	dummy_htmx_fragment = append(dummy_htmx_fragment, []byte(`<\div>`)...)
	socket.WriteMessage(ws.TextMessage, []byte(dummy_htmx_fragment))
}

package main

import (
	"fmt"
	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
	"net/http"
	_ "net/http/httputil"
)

func handler_home(responder http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		c.Yellow("received request that did match any other paths", request.URL.Path)
		c.Yellow("sending processing it as generic occurrence")
		session, err := session_from_cookie_from_request(request)
		data := sa{}
		data["user_id"], err = m.If_nil_do[string](err, func() (string, error) {
			c.Gray("checking if session is valid")
			return session_valid(session)
		})
		data["occurrence"] = request.URL.Path
		process_occurrence(data)
		responder.WriteHeader(204)
		responder.Write([]byte{})
		return
	}
	// check credentials/cookies in request
	var respond_with_credentials_home_page func(http.ResponseWriter) = func(responder http.ResponseWriter) {
		_, err := fmt.Fprintf(responder, page_credentials_blank_inlined())
		m.Warn(err, "failed to return credentials page")
		return
	}
	var respond_with_logged_home_page func(http.ResponseWriter, string) = func(responder http.ResponseWriter, user_id string) {
		_, err := fmt.Fprintf(responder, page_logged_inlined(user_id))
		m.Warn(err, "failed to return credentials page")
		return
	}
	session, err := session_from_cookie_from_request(request)
	user_id, err := m.If_nil_do[string](err, func() (string, error) {
		c.Gray("checking if session is valid")
		return session_valid(session)
	})
	m.If_error_must[nvm](err, func() {
		respond_with_credentials_home_page(responder)
		return
	})
	m.If_nil_must[nvm](err, func() { update_ephemeris(session) })
	m.If_nil_must[nvm](err, func() {
		respond_with_logged_home_page(responder, user_id)
	})

	return
}

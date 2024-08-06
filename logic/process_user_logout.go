package main

import (
	c "github.com/g-ameline/colors"
	m "github.com/g-ameline/maybe"
)

func process_user_logout(data sa) (string, error) {
	c.Blue("error state of demand")
	c.Blue("\n handling logout demand")

	session, ok := data["session"].(string)
	_, err := m.If_nil_do[string](m.Nok_to_err(ok, "no session in data"),
		func() (string, error) { return session_valid(session) })
	m.If_nil_try[nvm](err, func() {
		delete(data, "session")
		delete(data, "user_id")
		unregister_session(session)
	})
	m.Warn(err, "log out did not happen as expected")
	// in both case return main_frame_credentials_blank dom
	reply := main_frame_credentials_node(ss{}, ss{}).Inline()
	return reply, err
}

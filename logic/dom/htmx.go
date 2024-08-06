package dom

const (
	// Request Modifiers
	Hx_get        = "hx-get"
	Hx_post       = "hx-post"
	Hx_encoding   = "hx-encoding"
	Hx_put        = "hx-put"
	Hx_delete     = "hx-delete"
	Hx_patch      = "hx-patch"
	Hx_include    = "hx-include"
	Hx_disinherit = "hx-disinherit"

	// Request Headers and Content-Type
	Hx_headers = "hx-headers"
	Hx_content = "hx-content"

	// Request Parameters
	Hx_params = "hx-params"
	Hx_values = "hx-values"

	// Request Timeout and Retries
	Hx_timeout       = "hx-timeout"
	Hx_retry         = "hx-retry"
	Hx_retry_timeout = "hx-retry-timeout"

	// Response Processing
	Hx_target   = "hx-target"
	Hx_swap_OOB = "hx-swap-oob"
	Hx_select   = "hx-select"
	Hx_ext      = "hx-ext"
	Hx_vals     = "hx-vals"
	Hx_swap     = "hx-swap"

	//swap values
	InnerHTML   = "innerHTML"
	OuterHTML   = "outerHTML"
	Beforebegin = "beforebegin"
	Afterbegin  = "afterbegin"
	Beforeend   = "beforeend"
	Afterend    = "afterend"
	Delete      = "delete"
	None        = "none"

	// Events
	Hx_trigger            = "hx-trigger"
	Hx_confirm            = "hx-confirm"
	Hx_on                 = "hx-on"
	Hx_triggering_element = "hx-triggering-element"
	Hx_triggering_event   = "hx-triggering-event"

	// Indicators
	Hx_indicator = "hx-indicator"

	// History
	Hx_pushURL      = "hx-push-url"
	Hx_history_elt  = "hx-history-elt"
	Hx_history_attr = "hx-history-attr"

	// Error Handling
	Hx_boost = "hx-boost"
	Hx_error = "hx-error"

	// Caching
	Hx_cache = "hx-cache"

	// Triggers
	Click  = "click"
	Submit = "submit"
	Load   = "load"

	// Classes
	Container = "container"
	Card      = "card"

	// Request Modifiers
	Hx_ws      = "hx-ws"
	Hx_connect = "hx-connect"
	Hx_send    = "hx-send"
)

func HTMX_script() Node {
	script_node := New_node(Script).
		// Other("src", "https://unpkg.com/htmx.org@1.9.9").
		Other("src", "https://unpkg.com/htmx.org")
		// Other("crossorigin", "anonymous")
	return script_node
}

func HTMX_ws_script() Node {
	script_node := New_node(Script).
		Other("src", "https://unpkg.com/htmx.org/dist/ext/ws.js")
		// Other("crossorigin", "anonymous")
	return script_node
}

func (the_node Node) Hidden_name_value(name, value string) Node {
	// var include_value string = "["
	the_node.Other(Hx_include, "this").
		Add_kid(New_input().Name(name).Value(value).Conceal())
	return the_node
}

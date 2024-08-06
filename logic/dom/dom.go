package dom

import (
	"fmt"
	"strings"
)

// tags
const (
	Html   = "html"
	Body   = "body"
	Head   = "head"
	Meta   = "meta"
	Tag    = "tag"
	Div    = "div"
	Img    = "img"
	Button = "button"
	Form   = "form"
	Input  = "input"
	Label  = "label"
	Script = "script"
	Link   = "link"
	Aside  = "aside"
	Nav    = "nav"
	Ul     = "ul"
	Li     = "li"
)

// attributes
const (
	Id          = "id"
	Name        = "name"
	Text        = "text"
	Style       = "style"
	Class       = "classlist"
	Data        = "dataset"
	Kids        = "kids"
	Value       = "value"
	Placeholder = "placeholder"
	Src         = "src"
)
const (
	Rel  = "rel"
	Href = "href"
	Type = "type"
	For  = "for"
)

const (
	Multipart = "multipart/form-data"
	// Json= "hx-enc"
)

type Styles map[string]string
type Classes map[string]bool
type Dataset map[string]string
type Nodes []Node

type Node map[string]any

func New_node(node_tag string) Node {
	if node_tag == "" {
		node_tag = Div
	}
	a_node := make(Node)
	a_node[Tag] = node_tag
	return a_node
}

func New_div() Node {
	return New_node(Div)
}

func New_link() Node {
	return New_node(Link)
}

func New_a() Node {
	return New_node("a")
}

func New_p(text string) Node {
	return New_node("p").Text(text)
}

func New_form() Node {
	return New_node(Form)
}

func New_input() Node {
	return New_node(Input)
}

func New_script() Node {
	return New_node(Script)
}

func New_button() Node {
	return New_node(Button)
}

func New_body() Node {
	return New_node(Body)
}

func New_img(src string) Node {
	return New_node(Img).Attr(Src, src)
}

func (the_node Node) Name(identifier string) Node {
	the_node[Name] = identifier
	return the_node
}
func (the_node Node) Value(value string) Node {
	the_node[Value] = value
	return the_node
}
func (the_node Node) Type(type_value string) Node {
	the_node[Type] = type_value
	return the_node
}
func (the_node Node) Id(identifier string) Node {
	the_node[Id] = identifier
	return the_node
}
func (the_node Node) Conceal() Node {
	the_node.Style("display", "none")
	return the_node
}
func (the_node Node) Get_id() string {
	id, ok := the_node[Id]
	crash(ok, "no id")
	return id.(string)
}
func (the_node Node) Text(content string) Node {
	the_node[Text] = content
	return the_node
}
func (the_node Node) Get_text() string {
	content, ok := the_node[Text]
	crash(ok, "no text")
	return content.(string)
}
func (the_node Node) Style(style_key, style_value string) Node {
	_, ok := the_node[Style]
	if !ok {
		the_node[Style] = Styles{}
	}
	the_node[Style].(Styles)[style_key] = style_value
	return the_node
}
func (the_node Node) Get_styles() Styles {
	styles, ok := the_node[Style].(Styles)
	crash(ok, "not any style")
	return styles
}
func (the_node Node) Get_style(key string) string {
	styles := the_node.Get_styles()
	style, ok := styles[key]
	crash(ok, "no style")
	return style
}
func (the_node Node) Class(class string) Node {
	_, ok := the_node[Class]
	if !ok {
		the_node[Class] = Classes{}
	}
	the_node[Class].(Classes)[class] = true
	return the_node
}
func (the_node Node) Get_classes() Classes {
	classes, ok := the_node[Class].(Classes)
	crash(ok, "not any class")
	return classes
}
func (the_node Node) Href(url string) Node {
	the_node[Href] = url
	return the_node
}
func (the_node Node) Other(key, value string) Node {
	the_node[key] = value
	return the_node
}
func (the_node Node) Attr(key, value string) Node {
	the_node[key] = value
	return the_node
}
func (the_node Node) On_load(script string) Node {
	the_node["onload"] = script
	return the_node
}
func (the_node Node) Get_other(key string) string {
	other, ok := the_node[key]
	crash(ok, "no "+key)
	return other.(string)
}
func (the_node Node) Bear_kid(kid_node Node) Node {
	_, ok := the_node[Kids]
	if !ok {
		the_node[Kids] = []*Node{}
	}
	the_node[Kids] = append(the_node[Kids].([]*Node), &kid_node)
	return kid_node
}
func (the_node Node) Bear_firstborn(kid_node Node) Node {
	_, ok := the_node[Kids]
	if !ok {
		the_node[Kids] = []*Node{}
	}
	the_node[Kids] = append([]*Node{&kid_node}, the_node[Kids].([]*Node)...)
	return kid_node
}
func (the_node Node) Add_kid(kid_node Node) Node {
	the_node.Bear_kid(kid_node)
	return the_node
}
func (the_node Node) Add_kids(kids_nodes Nodes) Node {
	for _, kid_node := range kids_nodes {
		the_node.Add_kid(kid_node)
	}
	return the_node
}
func (the_node Node) Add_firstborn(kid_node Node) Node {
	the_node.Bear_firstborn(kid_node)
	return the_node
}
func (the_node Node) Get_kids() Nodes {
	kids_raw, ok := the_node[Kids]
	if !ok {
		return Nodes{}
	}
	kids_ps, ok := kids_raw.([]*Node)
	var kids []Node
	for _, kid := range kids_ps {
		kids = append(kids, *kid)
	}
	return kids
}

func inline_id(the_node Node) string {
	id, ok := the_node[Id].(string)
	if !ok {
		return ""
	}
	inline := "id = \"" + id + "\""
	return inline
}
func inline_styles(styles Styles) string {
	inline := ""
	for k, v := range styles {
		inline += k + ": " + v + "; "
	}
	return inline
}
func inline_classes(classes Classes) string {
	inline := ""
	for class := range classes {
		inline += class + " "
	}
	return strings.TrimSpace(inline)
}

func (the_node Node) Inline() string {
	inline := "<" + the_node[Tag].(string) + " "
	inline += the_node.inline_attributes()
	inline += ">\n"
	inline += the_node.inline_text()
	inline += the_node.inline_kids() + "\n"
	inline += "</" + the_node[Tag].(string) + ">"
	return inline
}
func (the_node Node) inline_attributes() string {
	inline := ""
	for key, value := range the_node {
		switch value.(type) {
		case Classes:
			inline += "class=" + quote(inline_classes(value.(Classes))) + "\n"
		case Styles:
			inline += "style=" + quote(inline_styles(value.(Styles))) + "\n"
		case string:
			if key == Text {
				break
			}
			if key == Tag {
				break
			}
			inline += key + " =" + quote(value.(string)) + "\n"
		case Nodes:
			break
		}
	}
	return inline
}
func (the_node Node) inline_text() string {
	text, ok := the_node[Text]
	if !ok {
		return ""
	}
	return text.(string)
}
func (the_node Node) inline_kids() string {
	kids, ok := the_node[Kids]
	if !ok {
		return ""
	}
	inline := ""
	for _, kid := range kids.([]*Node) {
		inline += (*kid).Inline()
	}
	return inline
}

func Inline_nodes(nodes Nodes) string {
	inline := ""
	for _, node := range nodes {
		inline += (node).Inline()
	}
	return inline
}
func (the_nodes Nodes) Inline() string {
	inline := ""
	for _, node := range the_nodes {
		inline += node.Inline()
	}
	return inline
}

func Prefix_doctype(inlined_node string) string {
	return "<!DOCTYPE html>" + inlined_node
}

func crash(false_or_error any, message string) {
	switch false_or_error.(type) {
	case error:
		if false_or_error != nil {
			panic(message)
		}
	case bool:
		if false_or_error == false {
			panic(message)
		}
	}
}

func quote(text string) string {
	return fmt.Sprintf("%q", text)
}

func (the_node Node) Pack() Node {
	the_node.
		Style(Border_style, Solid).
		Style(Padding, "10px").
		Style(Box_sizing, Border_box)
	return the_node
}
func (the_node Node) Order(i string) Node {
	return the_node.Style("order", i)
}
func (the_node Node) Margin(m string) Node {
	return the_node.Style(Margin, m+"px")
}
func (the_node Node) Margin_V(m string) Node {
	return the_node.
		Style(Margin_top, m+"px").
		Style(Margin_bottom, m+"px")
}

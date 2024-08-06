package dom

import "testing"
import "fmt"

func Test_a(t *testing.T) {
	// create node
	div := New_div()
	div.Id("1").
		Bear_kid(New_input().Id("2")).
		Bear_kid(New_button().Id("3")).Bear_kid(New_form().Id("4"))

	fmt.Println(div.Inline())
	outline_recursively(div, "8")

	fmt.Println(div.Inline())
}
func outline_recursively(node Node, key_color string) {
	node.Attr("AAAAAAAAAAAA", "BBBBBBBBB")
	for _, kid_node := range node.Get_kids() {
		outline_recursively(kid_node, key_color)
	}
}

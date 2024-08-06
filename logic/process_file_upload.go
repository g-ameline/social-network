package main

import (
	// "fmt"
	// uuid "github.com/google/uuid"
	d "app_server/dom"

	c "github.com/g-ameline/colors"
	// m "github.com/g-ameline/maybe"
)

func process_file_upload(data sa) (string, error) {
	c.Magenta("\n handling file upload and filepath return")
	c.Magenta("DATA FOR FILE UPLOAD")
	c.Print_map(data)
	filepath, ok := data["filepath"].(string)
	if !ok {
		panic("no file path")
	}
	picture_node := picture_node(filepath).Hidden_name_value("picture", filepath)
	return picture_node.Inline(), error(nil)
}

func picture_node(filepath string) d.Node {
	return d.New_img(filepath).Attr(d.Height, "50 !important")
}

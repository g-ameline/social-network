package main

import (
	"fmt"
	"net/http"
	_ "net/http/httputil"
	"os"
)

func handler_files(responder http.ResponseWriter, request *http.Request) {
	fmt.Println("request path", request.URL.Path)
	fmt.Println("serving file ?", request.URL.Path)
	// will serve files only if path match exactly a **file** in the file folder
	println(root + files_folder)
	println(root + files_folder + "/")
	if request.URL.Path == root+files_folder || request.URL.Path == root+files_folder+"/" {
		http.Error(responder, "can't access that folder", 404)
		return
	}
	path := "." + request.URL.Path
	_, err := os.Stat(path)
	if err != nil {
		http.Error(responder, "wrong file (path)", 404)
	}
	http.ServeFile(responder, request, path)
	return
}

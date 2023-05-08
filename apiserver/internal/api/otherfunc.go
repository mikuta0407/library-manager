package api

import (
	"fmt"
	"net/http"
	"strings"
)

func getRouteParams(r *http.Request) []string {
	splited := strings.Split(r.RequestURI, "/")
	var params []string
	for i := 0; i < len(splited); i++ {
		if len(splited[i]) != 0 {
			params = append(params, splited[i])
		}
	}
	return params
}

var libraryMode string

func judgeMode(params []string) error {
	if params[1] == "cd" {
		libraryMode = "cd"
	} else if params[1] == "book" {
		libraryMode = "book"
	} else {
		return fmt.Errorf("not cd or book")
	}
	return nil
}

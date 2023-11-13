package api

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func getRouteParams(r *http.Request, limit int) ([]string, error) {
	splited := strings.Split(r.RequestURI, "/")
	var params []string
	for i := 0; i < len(splited); i++ {
		if len(splited[i]) != 0 {
			params = append(params, splited[i])
		}
	}

	log.Println(params)
	if len(params) != limit {
		log.Println("neko", params)
		return nil, errors.New("Param length error")
	}
	return params, nil
}

package api

import (
	"fmt"
	"log"
	"net/http"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	params := getRouteParams(r)
	fmt.Println(params)
	if err := judgeMode(params); err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}
}

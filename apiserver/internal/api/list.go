package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mikuta0407/library-manager/internal/database"
)

func List(w http.ResponseWriter, r *http.Request) {
	params := getRouteParams(r)
	fmt.Println(params)
	if err := judgeMode(params); err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}
	database.GetList(libraryMode)
}

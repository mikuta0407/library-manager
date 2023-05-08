package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mikuta0407/library-manager/internal/database"
	"github.com/mikuta0407/library-manager/internal/models"
)

func List(w http.ResponseWriter, r *http.Request) {
	params := getRouteParams(r)
	fmt.Println(params)
	if err := judgeMode(params); err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}
	var items []models.Item
	items, err := database.GetList(libraryMode)

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}

	log.Println(items)
	fmt.Fprintf(w, "%v", items)

}

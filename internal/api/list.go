package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mikuta0407/library-manager/internal/database"
	"github.com/mikuta0407/library-manager/internal/models"
)

func List(w http.ResponseWriter, r *http.Request) {

	// POSTだけを受け入れる
	switch r.Method {
	case "GET":
	default:
		returnErrorMessage(w, http.StatusMethodNotAllowed, errors.New("Use GET Method"))
		return
	}

	params, err := getRouteParams(r, 3) // /api /list /category
	if err != nil {
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	var items models.ItemArray

	log.Println("Param:", params[2])
	items, err = database.GetList(params[2])

	if err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	log.Println(items)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")       // Content-Type指定
	if err := json.NewEncoder(w).Encode(items); err != nil { // JSON生成、応答
		log.Printf("json encode Error: %s", err)
	}
	return

}

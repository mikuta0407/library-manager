package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mikuta0407/library-manager/internal/database"
	"github.com/mikuta0407/library-manager/internal/models"
)

func Detail(w http.ResponseWriter, r *http.Request) {
	// /list/(book|cd)/{id}

	switch r.Method {
	case "GET":
	default:
		returnErrorMessage(w, http.StatusMethodNotAllowed, errors.New("Use GET Method"))
	}

	params, err := getRouteParams(r, 3)
	if err != nil {
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	fmt.Println(params)
	if err := judgeMode(params); err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(params[2])
	if err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusBadRequest, err)
	}
	var item models.Item
	item, err = database.GetDetail(libraryMode, id)

	if err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	log.Println(item)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")      // Content-Type指定
	if err := json.NewEncoder(w).Encode(item); err != nil { // JSON生成、応答
		log.Printf("json encode Error: %s", err)
	}
	return

}

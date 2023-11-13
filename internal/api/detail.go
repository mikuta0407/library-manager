package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/mikuta0407/library-manager/internal/database"
	"github.com/mikuta0407/library-manager/internal/models"
)

func Detail(w http.ResponseWriter, r *http.Request) {
	// /detail/{id}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	//必要なメソッドを許可する
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	switch r.Method {
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
		return
	case "GET":
	default:
		returnErrorMessage(w, http.StatusMethodNotAllowed, errors.New("Use GET Method"))
		return
	}

	params, err := getRouteParams(r, 3) // /api /detail /{id}
	if err != nil {
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(params[2])
	if err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusBadRequest, err)
	}
	var item models.Item
	item, err = database.GetDetail(id)

	if err != nil {
		log.Println(err)
		if err.Error() == "sql: no rows in result set" {
			returnErrorMessage(w, http.StatusNotFound, errors.New("No record"))
		} else {
			returnErrorMessage(w, http.StatusInternalServerError, err)
		}
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

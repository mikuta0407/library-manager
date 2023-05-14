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

func Delete(w http.ResponseWriter, r *http.Request) {
	// /delete/(book|cd)/{id}

	// DELETEだけを受け入れる
	// POSTだけを受け入れる
	switch r.Method {
	case "DELETE":
	default:
		returnErrorMessage(w, http.StatusMethodNotAllowed, errors.New("Use DELETE Method"))
		return
	}

	// パラメータ数確認
	params, err := getRouteParams(r, 4)
	if err != nil {
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	// book/cd判定
	fmt.Println(params)
	if err := judgeMode(params); err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	// 最低限バリデーション
	// id
	var id int
	id, err = strconv.Atoi(params[3])
	if err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusBadRequest, errors.New("id is not numeric"))
		return
	}

	err = database.DeleteItem(libraryMode, id)
	if err != nil {
		log.Println(err)
		if err.Error() == "No record" {
			returnErrorMessage(w, http.StatusNotFound, err)
		} else {
			returnErrorMessage(w, http.StatusInternalServerError, err)
		}
		return
	}

	// IDの返答
	log.Println("Delete OK!")
	resMsg := models.SuccessResponseMessage{
		Message: "Success",
		Id:      strconv.FormatInt(int64(id), 10),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")       // Content-Type指定
	if err = json.NewEncoder(w).Encode(resMsg); err != nil { // JSON生成、応答
		log.Printf("json encode Error: %s", err)
	}
	return
}

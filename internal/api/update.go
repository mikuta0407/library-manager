package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/mikuta0407/library-manager/internal/database"
	"github.com/mikuta0407/library-manager/internal/models"
)

func Update(w http.ResponseWriter, r *http.Request) {
	// /update/(book|cd)/

	// PUTだけを受け入れる
	switch r.Method {
	case "PUT":
	default:
		returnErrorMessage(w, http.StatusMethodNotAllowed, errors.New("Use PUT Method"))
		return
	}

	// MIMEタイプ確認
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("MIME Type Error:", r.Header.Get("Content-Type"))
		returnErrorMessage(w, http.StatusUnsupportedMediaType, errors.New("Use application/json"))
		return
	}

	// パラメータ数確認
	params, err := getRouteParams(r, 3) // /api /update /id
	if err != nil {
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	// 容量確認 (実BodyData)
	const maxDataSize int = 5242880                                             // 5MBに制限
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, int64(maxDataSize))) //maxDataSize分だけ読む、のをReadAllで読んでdataへ

	if err != nil {
		if err.Error() == "http: request body too large" { // 読んだバイト数が(指定した)最大バイト数を超えた場合蹴る
			log.Println(err)
			returnErrorMessage(w, http.StatusRequestEntityTooLarge, errors.New("Body should be less than 5242880 bytes"))
			return
		} else {
			log.Println(err)
			returnErrorMessage(w, http.StatusInternalServerError, err)
			return
		}
	}
	log.Println("Request Size: OK", len(data))

	// POSTされてきたJSONから構造体へデコード
	var item models.Item
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&item)
	if err != nil {
		log.Println("JSON Decode Error:", err)
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	log.Println(item)

	// 最低限バリデーション

	// id
	item.Id, err = strconv.Atoi(params[2])
	if err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusBadRequest, errors.New("id is not numeric"))
		return
	}
	// タイトル
	if item.Title == "" {
		log.Println(err)
		returnErrorMessage(w, http.StatusBadRequest, errors.New("no title specified"))
		return
	}

	// UPDATE実行
	err = database.UpdateItem(item)
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
	log.Println("Update OK!")
	resMsg := models.SuccessResponseMessage{
		Message: "Success",
		Id:      strconv.FormatInt(int64(item.Id), 10),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")       // Content-Type指定
	if err = json.NewEncoder(w).Encode(resMsg); err != nil { // JSON生成、応答
		log.Printf("json encode Error: %s", err)
	}
	return

}

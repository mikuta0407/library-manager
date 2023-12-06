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

func Create(w http.ResponseWriter, r *http.Request) {
	// /api/create

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	//必要なメソッドを許可する
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	log.Println(r.Method)

	// POSTだけを受け入れる
	switch r.Method {
	case "OPTIONS":
		log.Println("OPTIONS!")
		w.WriteHeader(http.StatusOK)
		return
	case "POST":
	default:
		returnErrorMessage(w, http.StatusMethodNotAllowed, errors.New("Use POST Method"))
		return
	}

	// MIMEタイプ確認
	if r.Header.Get("Content-Type") != "application/json" {
		log.Println("MIME Type Error:", r.Header.Get("Content-Type"))
		returnErrorMessage(w, http.StatusUnsupportedMediaType, errors.New("Use application/json"))
		return
	}

	// パラメータ数確認
	_, err := getRouteParams(r, 2) // /api /create
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
	if item.Title == "" {
		log.Println(err)
		returnErrorMessage(w, http.StatusBadRequest, errors.New("no title specified"))
		return
	}
	if item.Category == "" {
		log.Println(err)
		returnErrorMessage(w, http.StatusBadRequest, errors.New("no category specified"))
		return
	}

	if item.Image == "" {
		item.Image = "No Image"
	}

	id, err := database.CreateItem(item)
	if err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	// IDの返答
	log.Println("Regist OK!")
	resMsg := models.SuccessResponseMessage{
		Message: "Success",
		Id:      strconv.FormatInt(id, 10),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")       // Content-Type指定
	if err = json.NewEncoder(w).Encode(resMsg); err != nil { // JSON生成、応答
		log.Printf("json encode Error: %s", err)
	}
	return

}

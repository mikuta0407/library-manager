package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/mikuta0407/library-manager/internal/database"
	"github.com/mikuta0407/library-manager/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// /api/login
	var userLoginInfo models.UserExternalLogin
	var userResponse models.UserExternalResponse
	var userConfidential models.UserInternal

	// POSTだけを受け入れる
	switch r.Method {
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
	_, err := getRouteParams(r, 2) // /api /login
	if err != nil {
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	// 容量確認 (実BodyData)
	const maxDataSize int = 1024                                                // 1KBに制限
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
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&userLoginInfo)
	if err != nil {
		log.Println("JSON Decode Error:", err)
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	// ユーザー情報取得
	userConfidential, err = database.GetUserByName(userLoginInfo.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			returnErrorMessage(w, http.StatusUnauthorized, err)
			log.Println("Login failed(No Users):", userLoginInfo.UserName)
		} else {
			log.Println(err)
			returnErrorMessage(w, http.StatusInternalServerError, errors.New("Internal Server Error"))
			log.Println("Login failed(DB Error):", userLoginInfo.UserName)
		}
		return
	}

	// ハッシュ一致
	if err := bcrypt.CompareHashAndPassword([]byte(userConfidential.PasswordHash), []byte(userLoginInfo.Password)); err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusUnauthorized, errors.New("Login failed"))
		log.Println("Login failed(Password mismatch):", userLoginInfo.UserName)
		return
	}
	log.Println("Login Success:", userLoginInfo.UserName)

	userResponse.UserName = userConfidential.UserName
	userResponse.UUID = userConfidential.UUID
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")             // Content-Type指定
	if err = json.NewEncoder(w).Encode(userResponse); err != nil { // JSON生成、応答
		log.Printf("json encode Error: %s", err)
	}
	return

}

func Regist(w http.ResponseWriter, r *http.Request) {
	// /api/regist
	var userLoginInfo models.UserExternalLogin
	var userConfidential models.UserInternal

	// POSTだけを受け入れる
	switch r.Method {
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
	_, err := getRouteParams(r, 2) // /api /regist
	if err != nil {
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	// 容量確認 (実BodyData)
	const maxDataSize int = 1024                                                // 1KBに制限
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
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&userLoginInfo)
	if err != nil {
		log.Println("JSON Decode Error:", err)
		returnErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	// ユーザー情報取得(あるかないか)
	userConfidential, err = database.GetUserByName(userLoginInfo.UserName)
	var hash []byte
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		returnErrorMessage(w, http.StatusInternalServerError, err)
		return
	} else if err == sql.ErrNoRows {
		hash, err = bcrypt.GenerateFromPassword([]byte(userLoginInfo.Password), 12)
		if err != nil {
			log.Println(err)
			returnErrorMessage(w, http.StatusInternalServerError, err)
			return
		}

	} else {
		returnErrorMessage(w, http.StatusInternalServerError, errors.New("User is already exists"))
		return
	}

	uuidTmp, err := uuid.NewV4()

	userConfidential.UserName = userLoginInfo.UserName
	userConfidential.PasswordHash = string(hash)
	userConfidential.UUID = uuidTmp.String()

	if err := database.AddUser(userConfidential); err != nil {
		log.Println(err)
		returnErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	log.Println("Regist OK!")
	resMsg := models.SuccessResponseMessage{
		Message: "Success",
		Id:      userConfidential.UUID,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")       // Content-Type指定
	if err = json.NewEncoder(w).Encode(resMsg); err != nil { // JSON生成、応答
		log.Printf("json encode Error: %s", err)
	}
	return

}

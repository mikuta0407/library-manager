package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mikuta0407/library-manager/internal/models"
)

func returnErrorMessage(w http.ResponseWriter, code int, err error) {

	log.Println("Code:", code)
	// メッセージ生成
	var errMsg models.ErrMessage
	switch code {
	case http.StatusBadRequest:
		log.Println("Bad Request")
		errMsg.ErrMessage = "Bad Request"
	case http.StatusNotFound:
		log.Println("Not Found")
		errMsg.ErrMessage = "Not Found"
	case http.StatusMethodNotAllowed:
		log.Println("Method Not Allowed")
		errMsg.ErrMessage = "Method Not Allowed"
	case http.StatusConflict:
		log.Println("Conflict")
		errMsg.ErrMessage = "Conflict"
	case http.StatusRequestEntityTooLarge:
		log.Println("Payload Too Large")
		errMsg.ErrMessage = "Payload Too Large"
	case http.StatusRequestURITooLong:
		log.Println("URI Too Long")
		errMsg.ErrMessage = "URI Too Long"
	case http.StatusUnsupportedMediaType:
		log.Println("Unsupported Media Type")
		errMsg.ErrMessage = "Unsupported Media Type"
	case http.StatusInternalServerError:
		log.Println("Internal Server Error")
		errMsg.ErrMessage = "Internal Server Error"
	default:
		log.Println("Undefined code: ", code)
		code = http.StatusInternalServerError // 返答は500とする
		errMsg.ErrMessage = "Internal Server Error"
	}

	errMsg.ErrDetail = err.Error()

	w.WriteHeader(code)

	// メッセージ応答
	w.Header().Set("Content-Type", "application/json")        // Content-Type指定
	if err := json.NewEncoder(w).Encode(errMsg); err != nil { // JSON生成、応答
		log.Printf("json encode Error: %s", err)
	}
	return
}

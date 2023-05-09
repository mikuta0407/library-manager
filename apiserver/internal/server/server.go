package server

import (
	"log"
	"net/http"

	"github.com/mikuta0407/library-manager/internal/api"
	"github.com/mikuta0407/library-manager/internal/database"
)

func HandleRequests() {
	if err := database.ConnectDB("./library.db"); err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/list/", api.List)     // 一覧
	http.HandleFunc("/detail/", api.Detail) // 詳細
	http.HandleFunc("/search/", api.Search) // 検索
	http.HandleFunc("/create/", api.Create) // レコード作成

	//http.HandleFunc("/update/", api.Update) // レコード編集
	//http.HandleFunc("/delete/", api.Delete) // レコード削除
	log.Fatal(http.ListenAndServe(":8081", nil))
	defer database.DisconnectDB()
}

func main() {
	HandleRequests()
}

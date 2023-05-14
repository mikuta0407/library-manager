package server

import (
	"log"
	"net/http"

	"github.com/mikuta0407/library-manager/internal/api"
	"github.com/mikuta0407/library-manager/internal/database"
)

func HandleRequests(httpHost string, httpPort string, sqliteDBPath string) {
	if err := database.ConnectDB(sqliteDBPath); err != nil {
		log.Fatalln(err)
	}

	// APIサーバー
	http.HandleFunc("/api/list/", api.List)     // 一覧
	http.HandleFunc("/api/detail/", api.Detail) // 詳細
	http.HandleFunc("/api/search/", api.Search) // 検索
	http.HandleFunc("/api/create/", api.Create) // レコード作成
	http.HandleFunc("/api/update/", api.Update) // レコード編集
	http.HandleFunc("/api/delete/", api.Delete) // レコード削除
	log.Fatal(http.ListenAndServe(httpHost+":"+httpPort, nil))
	defer database.DisconnectDB()
}

package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")

	params := getRouteParams(r)
	fmt.Println(params)

}

func HandleRequests() {
	http.HandleFunc("/list/", homePage)   // 一覧
	http.HandleFunc("/detail/", homePage) // 詳細
	http.HandleFunc("/create/", homePage) // レコード作成
	http.HandleFunc("/update/", homePage) // レコード編集
	http.HandleFunc("/detail/", homePage) // レコード削除
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	HandleRequests()
}

func getRouteParams(r *http.Request) []string {
	splited := strings.Split(r.RequestURI, "/")
	var params []string
	for i := 0; i < len(splited); i++ {
		if len(splited[i]) != 0 {
			params = append(params, splited[i])
		}
	}
	return params
}

package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikuta0407/library-manager/internal/models"
)

var db *sql.DB

// DB接続
func ConnectDB(filename string) error {
	var err error

	// 接続開始
	db, err = sql.Open("sqlite3", filename)
	if err != nil {
		return err
	}

	// Ping確認
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB Connected")

	return nil
}

// DB切断
func DisconnectDB() error {
	if err := db.Close(); err != nil {
		log.Fatalln(err)
	}
	return nil
}

// book/cd全一覧取得
func GetList(libraryMode string) (models.ItemArray, error) {

	var items models.ItemArray

	// SELECT *する
	rows, err := db.Query("SELECT * FROM " + libraryMode)
	if err != nil {
		return items, err
	}

	defer rows.Close()

	// rowごとに一旦突っ込んでappendでスライスに追加
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Id, &item.Title, &item.Author, &item.Code, &item.Place, &item.Note, &item.Image); err != nil {
			return items, err
		}
		items.ItemList = append(items.ItemList, item)
	}
	return items, nil
}

// 1つだけ取得(IDをもとに)
func GetDetail(libraryMode string, id int) (models.Item, error) {

	var item models.Item

	// プリペアドステートメント作成
	var prepStmt string
	if libraryMode == "book" {
		prepStmt = "SELECT * FROM book WHERE id = $1"
	} else if libraryMode == "cd" {
		prepStmt = "SELECT * FROM cd WHERE id = $1"
	}

	// 実行
	prep, err := db.Prepare(prepStmt)
	if err != nil {
		return item, err
	}
	defer prep.Close()

	// item構造体に突っ込んで返却
	err = prep.QueryRow(id).Scan(&item.Id, &item.Title, &item.Author, &item.Code, &item.Place, &item.Note, &item.Image)
	if err != nil {
		return item, err
	}

	return item, nil
}

// DBにINSERTする
func CreateItem(libraryMode string, item models.Item) (int64, error) {

	var insertId int64 // insert idを入れる

	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	// クエリ準備
	var prepStmt string
	if libraryMode == "book" {
		prepStmt = "INSERT INTO book (title, author, code, place, note, image) values ($1, $2, $3, $4, $5, $6)"
	} else if libraryMode == "cd" {
		prepStmt = "INSERT INTO cd (title, artist, code, place, note, image) values ($1, $2, $3, $4, $5, $6)"
	}

	// INSRT実行
	res, err := tx.Exec(prepStmt, item.Title, item.Author, item.Code, item.Place, item.Note, item.Image)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// ID取得
	insertId, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// コミット
	tx.Commit()

	return insertId, nil
}

func SearchItem(libraryMode string, item models.Item) (models.ItemArray, error) {

	// 使うやつ宣言

	var items models.ItemArray // 返すアイテム一覧

	var searchResultId []int64 //IDリスト
	var res []int64            // 一時スライス
	var err error              // errorよう

	/* 一旦それぞれの項目で検索し、該当のIDを取得する (重複削除前提でとりあえず全部スライスに入れる) */
	// タイトル検索
	if item.Title != "" {
		res, err = searchColumnId(libraryMode, "title", item.Title)
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// 著者・アーティスト名検索
	if item.Author != "" {
		if libraryMode == "book" {
			res, err = searchColumnId(libraryMode, "author", item.Author)
		} else if libraryMode == "cd" {
			res, err = searchColumnId(libraryMode, "artist", item.Author)
		}
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// コード検索
	if item.Code != "" {
		res, err = searchColumnId(libraryMode, "code", item.Code)
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// 場所検索
	if item.Place != "" {
		res, err = searchColumnId(libraryMode, "place", item.Place)
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// 備考欄検索
	if item.Note != "" {
		res, err = searchColumnId(libraryMode, "place", item.Note)
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// 重複削除
	m := make(map[int64]bool)
	uniqSearchResultId := []int64{}

	for _, ele := range searchResultId {
		if !m[ele] {
			m[ele] = true
			uniqSearchResultId = append(uniqSearchResultId, ele)
		}
	}

	// 得られたIDを元にSELECT
	prepStmt := "SELECT * FROM " + libraryMode + " WHERE id in ("
	for i, id := range uniqSearchResultId {
		prepStmt = prepStmt + strconv.FormatInt(id, 10)
		if i != len(uniqSearchResultId)-1 {
			prepStmt = prepStmt + ", "
		}
	}
	prepStmt = prepStmt + ")"

	log.Println(prepStmt)

	rows, err := db.Query(prepStmt)
	if err != nil {
		return items, err
	}

	defer rows.Close()

	// GetListと同じようにitemスライスに入れる
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Id, &item.Title, &item.Author, &item.Code, &item.Place, &item.Note, &item.Image); err != nil {
			return items, err
		}
		items.ItemList = append(items.ItemList, item)
	}
	return items, nil

}

// SearchItemで項目ごとに検索するための関数。IDのスライスを返却する
func searchColumnId(libraryMode string, columnName string, value string) ([]int64, error) {
	var res []int64

	prepStmt := "SELECT id FROM " + libraryMode + " WHERE " + columnName + " LIKE '%'||?||'%'"
	fmt.Println(prepStmt)
	rows, err := db.Query(prepStmt, value)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return res, err
		}

		res = append(res, id)
	}

	return res, nil

}

// UPDATEする
func UpdateItem(libraryMode string, item models.Item) (int64, error) {

	itemId := int64(item.Id)

	// 存在するのか確認
	if _, err := GetDetail(libraryMode, int(itemId)); err != nil {
		return itemId, errors.New("No record")
	}

	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	// UPDATEクエリ準備
	var prepStmt string
	if libraryMode == "book" {
		prepStmt = "UPDATE book SET title = ?, author = ?, code = ?, place = ?, note = ?, image = ? WHERE id = ?"
	} else if libraryMode == "cd" {
		prepStmt = "UPDATE cd SET title = ?, artist = ?, code = ?, place = ?, note = ?, image = ? WHERE id = ?"
	}

	// Update実行
	_, err = tx.Exec(prepStmt, item.Title, item.Author, item.Code, item.Place, item.Note, item.Image, item.Id)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// コミット
	tx.Commit()

	return itemId, nil
}

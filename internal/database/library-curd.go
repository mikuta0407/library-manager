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

// カテゴリ内全一覧取得
func GetList(category string) (models.ItemArray, error) {

	var items models.ItemArray

	// SELECTする
	var rows *sql.Rows
	var err error

	if category == "all" {
		rows, err = db.Query("SELECT id, title, category, author, code, purchase, place, note, image FROM library")
		if err != nil {
			return items, err
		}
	} else {
		rows, err = db.Query("SELECT id, title, category, author, code, purchase, place, note, image FROM library WHERE category = $1", category)
		if err != nil {
			return items, err
		}
	}

	defer rows.Close()

	// rowごとに一旦突っ込んでappendでスライスに追加
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Id, &item.Title, &item.Category, &item.Author, &item.Code, &item.Purchase, &item.Place, &item.Note, &item.Image); err != nil {
			return items, err
		}
		items.ItemList = append(items.ItemList, item)
	}
	return items, nil
}

// 1つだけ取得(IDをもとに)
func GetDetail(id int) (models.Item, error) {

	var item models.Item

	// プリペアドステートメント作成
	prepStmt := "SELECT id, title, category, author, code, purchase, place, note, image FROM library WHERE id = $1"

	// 実行
	prep, err := db.Prepare(prepStmt)
	if err != nil {
		return item, err
	}
	defer prep.Close()

	// item構造体に突っ込んで返却
	err = prep.QueryRow(id).Scan(&item.Id, &item.Title, &item.Category, &item.Author, &item.Code, &item.Purchase, &item.Place, &item.Note, &item.Image)
	if err != nil {
		return item, err
	}

	return item, nil
}

// DBにINSERTする
func CreateItem(item models.Item) (int64, error) {

	var insertId int64 // insert idを入れる

	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	// クエリ準備
	prepStmt := "INSERT INTO library (title, category, author, code, purchase, place, note, image) values ($1, $2, $3, $4, $5, $6, $7, $8)"

	// INSRT実行
	res, err := tx.Exec(prepStmt, item.Title, item.Category, item.Author, item.Code, item.Purchase, item.Place, item.Note, item.Image)
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

func SearchItem(item models.Item) (models.ItemArray, error) {

	// 使うやつ宣言

	var items models.ItemArray // 返すアイテム一覧

	var searchResultId []int64 //IDリスト
	var res []int64            // 一時スライス
	var err error              // errorよう

	/* 一旦それぞれの項目で検索し、該当のIDを取得する (重複削除前提でとりあえず全部スライスに入れる) */
	// タイトル検索
	if item.Title != "" {
		res, err = searchColumnId("title", item.Title)
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// 著者・アーティスト名検索
	if item.Author != "" {
		res, err = searchColumnId("author", item.Author)
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// コード検索
	if item.Code != "" {
		res, err = searchColumnId("code", item.Code)
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// 購入場所検索
	if item.Purchase != "" {
		res, err = searchColumnId("purchase", item.Purchase)
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// 場所検索
	if item.Place != "" {
		res, err = searchColumnId("place", item.Place)
		if err != nil {
			return items, err
		}
		searchResultId = append(searchResultId, res...)
	}

	// 備考欄検索
	if item.Note != "" {
		res, err = searchColumnId("node", item.Note)
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
	prepStmt := "SELECT id, title, category, author, code, purchase, place, note, image FROM library WHERE id in ("
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
		if err := rows.Scan(&item.Id, &item.Title, &item.Author, &item.Code, &item.Purchase, &item.Place, &item.Note, &item.Image); err != nil {
			return items, err
		}
		items.ItemList = append(items.ItemList, item)
	}
	return items, nil

}

// SearchItemで項目ごとに検索するための関数。IDのスライスを返却する
func searchColumnId(columnName string, value string) ([]int64, error) {
	var res []int64

	prepStmt := "SELECT id FROM library WHERE " + columnName + " LIKE '%'||?||'%'"
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
func UpdateItem(item models.Item) error {

	itemId := int64(item.Id)

	// 存在するのか確認
	if _, err := GetDetail(int(itemId)); err != nil {
		return errors.New("No record")
	}

	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// UPDATEクエリ準備
	prepStmt := "UPDATE library SET title = ?, category = ?,author = ?, code = ?, purchase = ?, place = ?, note = ?, image = ? WHERE id = ?"

	// Update実行
	_, err = tx.Exec(prepStmt, item.Title, item.Category, item.Author, item.Code, item.Purchase, item.Place, item.Note, item.Image, item.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// コミット
	tx.Commit()

	return nil
}

// DELETEする
func DeleteItem(itemId int) error {

	// 存在するのか確認
	if _, err := GetDetail(itemId); err != nil {
		return errors.New("No record")
	}

	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// DELETEクエリ準備
	prepStmt := "DELETE FROM library WHERE id = ?"

	// Delete実行
	_, err = tx.Exec(prepStmt, itemId)
	if err != nil {
		tx.Rollback()
		return err
	}

	// コミット
	tx.Commit()

	return nil
}

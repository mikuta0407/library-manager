package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikuta0407/library-manager/internal/models"
)

var db *sql.DB

func ConnectDB(filename string) error {
	var err error
	db, err = sql.Open("sqlite3", filename) //接続開始（example.sqlに保存する）ConnectDB
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB Connected")

	return nil
}

func DisconnectDB() error {
	if err := db.Close(); err != nil {
		log.Fatalln(err)
	}
	return nil
}

func GetList(libraryMode string) (models.ItemArray, error) {

	var items models.ItemArray

	rows, err := db.Query("SELECT * FROM " + libraryMode)
	if err != nil {
		return items, err
	}

	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Id, &item.Title, &item.Author, &item.Code, &item.Place, &item.Note, &item.Image); err != nil {
			return items, err
		}
		items.ItemList = append(items.ItemList, item)
	}
	return items, nil
}

func GetDetail(libraryMode string, id int) (models.Item, error) {

	var item models.Item

	var prepStmt string
	if libraryMode == "book" {
		prepStmt = "SELECT * FROM book WHERE id = $1"
	} else if libraryMode == "cd" {
		prepStmt = "SELECT * FROM cd WHERE id = $1"
	}
	prep, err := db.Prepare(prepStmt)
	if err != nil {
		return item, err
	}
	defer prep.Close()
	err = prep.QueryRow(id).Scan(&item.Id, &item.Title, &item.Author, &item.Code, &item.Place, &item.Note, &item.Image)
	if err != nil {
		return item, err
	}

	return item, nil
}

func CreateItem(libraryMode string) {

}

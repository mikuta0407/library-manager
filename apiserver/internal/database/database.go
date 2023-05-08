package database

import (
	"database/sql"
	"fmt"
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

func GetList(libraryMode string) {
	rows, err := db.Query(
		`SELECT * FROM book`,
	)
	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Id, &item.Title, &item.Author, &item.Code, &item.Place, &item.Note, &item.Image); err != nil {
			log.Fatal("rows.Scan()", err)
			return
		}
		fmt.Println(item)
	}
}

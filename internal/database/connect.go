package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
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
		return err
	}
	log.Println("DB Connected")

	// TABLE作成(ない場合に)
	createLibraryTable := `
	CREATE TABLE IF NOT EXISTS "library" (
		"id"	INTEGER,
		"title"	TEXT NOT NULL,
		"category" TEXT,
		"author"	TEXT,
		"code"	TEXT,
		"purchase" TEXT,
		"place"	TEXT,
		"note"	TEXT,
		"image"	TEXT,
		"user" TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	)`

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER,
		"username" VARCHAR(20),
		"password_hash" TEXT NOT NULL,
		"uuid" TEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	)`

	_, err = db.Exec(createLibraryTable)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// DB切断
func DisconnectDB() error {
	if err := db.Close(); err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

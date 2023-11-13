package database

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikuta0407/library-manager/internal/models"
)

func GetUserByName(username string) (models.UserInternal, error) {
	var user models.UserInternal

	// プリペアドステートメント作成
	prepStmt := "SELECT * FROM users WHERE username = $1"

	// 実行
	prep, err := db.Prepare(prepStmt)
	if err != nil {
		return user, err
	}
	defer prep.Close()

	err = prep.QueryRow(username).Scan(&user.Id, &user.UserName, &user.PasswordHash, &user.UUID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func GetUserByUUID(uuid string) (int64, error) {
	var userId int64

	// プリペアドステートメント作成
	prepStmt := "SELECT id FROM users WHERE uuid = $1"

	// 実行
	prep, err := db.Prepare(prepStmt)
	if err != nil {
		return 0, err
	}
	defer prep.Close()

	err = prep.QueryRow(uuid).Scan(&userId)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return userId, nil
}

func AddUser(user models.UserInternal) error {

	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// クエリ準備
	prepStmt := "INSERT INTO users (username, password_hash, uuid) values ($1, $2, $3)"

	// INSRT実行
	_, err = tx.Exec(prepStmt, user.UserName, user.PasswordHash, user.UUID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// コミット
	tx.Commit()

	return nil
}

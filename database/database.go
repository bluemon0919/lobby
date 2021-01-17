package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tenntenn/sqlite"
)

// EntitySQL データを管理する
type EntitySQL struct {
	filename string // TODO:ファイルの拡張子はシステム側でつけてもいいかも
	db       *sql.DB
}

// Item is SQL item
type Item struct {
	UserName   string
	NumOfGames int
	NumOfWins  int
}

// NewSQL creates Entity
func NewSQL(filename string) *EntitySQL {
	db, err := sql.Open(sqlite.DriverName, filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	const sql = `CREATE TABLE IF NOT EXISTS item (
		key   INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		UserName TEXT NOT NULL,
		NumOfGames INTEGER NOT NULL,
		NumOfWins INTEGER NOT NULL);`

	if _, err = db.Exec(sql); err != nil {
		return nil
	}

	return &EntitySQL{
		filename: filename,
		db:       db,
	}
}

// Add adds item
func (e *EntitySQL) Add(item *Item) error {
	switch {
	case item.UserName == "":
		return errors.New("UserName is empty")
	case item.NumOfGames < 0:
		return errors.New("NumOfGames is out of range")
	case item.NumOfWins < 0:
		return errors.New("NumOfWins is out of range")
	}
	const sql = "INSERT INTO item(UserName, NumOfGames, NumOfWins) values (?,?,?)"
	_, err := e.db.Exec(sql, item.UserName, item.NumOfGames, item.NumOfWins)
	return err
}

// Delete delete item from key
func (e *EntitySQL) Delete(key int) error {
	sql := "DELETE FROM item WHERE key = ?"
	_, err := e.db.Exec(sql, key)
	return err
}

// Update Entityの指定のキーを入力ステータスでアップデートする
func (e *EntitySQL) Update(key, numOfGames, numOfWins int) error {
	sql := "UPDATE item SET numOfGames = ?, numOfWins = ? WHERE key = ?"
	_, err := e.db.Exec(sql, numOfGames, numOfWins, key)
	return err
}

// Get Entityからアイテムを取得する
func (e *EntitySQL) Get(userName string) (int, Item, error) {
	var item Item
	const sql = "SELECT * FROM item WHERE UserName = ?"
	rows, err := e.db.Query(sql, userName)
	if err != nil {
		return 0, item, err
	}

	key := 0
	for rows.Next() {
		if err := rows.Scan(
			&key,
			&item.UserName,
			&item.NumOfGames,
			&item.NumOfWins); err != nil {
			return 0, item, err
		}
	}
	return key, item, nil
}

// IsEmpty judge Item empty
func (i Item) IsEmpty() bool {
	if i.UserName == "" {
		return true
	}
	return false
}

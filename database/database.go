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

// Opponent is Opponent list
type Opponent struct {
	UserName     string
	OpponentName string // 対戦相手の名前

	// 対戦順を管理する。1〜5で１が最新、５が最古となるように設定する
	// databaseの利用側で番号管理すること
	Num int
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

	const sqlOpponent = `CREATE TABLE IF NOT EXISTS opponent (
		key   INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		UserName TEXT NOT NULL,
		OpponentName TEXT NOT NULL,
		Num INTEGER NOT NULL);`

	if _, err = db.Exec(sql); err != nil {
		return nil
	}
	if _, err = db.Exec(sqlOpponent); err != nil {
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
	defer rows.Close()

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

// IsEmpty judge Opponent empty
func (opp Opponent) IsEmpty() bool {
	if opp.UserName == "" || opp.OpponentName == "" {
		return true
	}
	return false
}

// AddOpponent adds opponent
func (e *EntitySQL) AddOpponent(opp *Opponent) error {
	switch {
	case opp.UserName == "":
		return errors.New("UserName is empty")
	case opp.OpponentName == "":
		return errors.New("OpponentName is empty")
	default:
		const sql = "INSERT INTO opponent(UserName, OpponentName, Num) values (?,?,?)"
		_, err := e.db.Exec(sql, opp.UserName, opp.OpponentName, opp.Num)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteOpponent delete opponent from key
func (e *EntitySQL) DeleteOpponent(key int) error {
	sql := "DELETE FROM opponent WHERE key = ?"
	_, err := e.db.Exec(sql, key)
	return err
}

// UpdateOpponent Entityの指定のキーを入力ステータスでアップデートする
func (e *EntitySQL) UpdateOpponent(key int, UserName, OpponentName string, Num int) error {
	sql := "UPDATE opponent SET UserName = ?, OpponentName = ?, Num = ? WHERE key = ?"
	_, err := e.db.Exec(sql, UserName, OpponentName, Num, key)
	return err
}

// GetOpponent Entityからアイテムを取得する
func (e *EntitySQL) GetOpponent(userName string) ([]int, []Opponent, error) {
	var opps []Opponent
	const sql = "SELECT * FROM opponent WHERE UserName = ?"
	rows, err := e.db.Query(sql, userName)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	keys := []int{}
	for rows.Next() {
		key := 0
		var opp Opponent
		if err := rows.Scan(
			&key,
			&opp.UserName,
			&opp.OpponentName,
			&opp.Num); err != nil {
			return nil, nil, err
		}
		keys = append(keys, key)
		opps = append(opps, opp)
	}
	return keys, opps, nil
}

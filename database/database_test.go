package database

import (
	"log"
	"os"
	"testing"
)

const TestEntitySQLTestFileName = "test.db"

func TestNewSQL(t *testing.T) {

	expect := TestEntitySQLTestFileName
	// ファイルを削除する
	if _, err := os.Stat(expect); !os.IsNotExist(err) {
		if err = os.Remove(expect); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewSQL(expect)
	if ent.filename != expect {
		log.Fatalf("file name does not match. %s\n", ent.filename)
	}
	if ent.db == nil {
		log.Fatal("db is nil")
	}

	// ファイルを削除する
	if _, err := os.Stat(expect); !os.IsNotExist(err) {
		if err = os.Remove(expect); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}

func TestSQLAdd(t *testing.T) {

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewSQL(TestEntitySQLTestFileName)

	item := Item{
		UserName:   "hoge",
		NumOfGames: 10,
		NumOfWins:  5,
	}
	if err := ent.Add(&item); err != nil {
		log.Fatal("registration failed")
	}

	item = Item{
		NumOfGames: 10,
		NumOfWins:  5,
	}
	if err := ent.Add(&item); err == nil {
		log.Fatal("register even if the UserName is empty")
	}

	item = Item{
		UserName:   "fuga",
		NumOfGames: -1,
		NumOfWins:  9,
	}
	if err := ent.Add(&item); err == nil {
		log.Fatal("register even if the NumOfGames is out of range")
	}

	item = Item{
		UserName:   "hoge",
		NumOfGames: 10,
		NumOfWins:  -1,
	}
	if err := ent.Add(&item); err == nil {
		log.Fatal("register even if the NumOfWins is out of range")
	}

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}

func TestSQLDelete(t *testing.T) {

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewSQL(TestEntitySQLTestFileName)

	item := Item{
		UserName:   "hoge",
		NumOfGames: 10,
		NumOfWins:  5,
	}

	if err := ent.Add(&item); err != nil {
		log.Fatal("registration failed")
	}

	if err := ent.Delete(1); err != nil {
		log.Fatal("delete failed")
	}

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}

func TestSQLUpdate(t *testing.T) {

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewSQL(TestEntitySQLTestFileName)

	item := Item{
		UserName:   "hoge",
		NumOfGames: 10,
		NumOfWins:  5,
	}

	if err := ent.Add(&item); err != nil {
		log.Fatal("registration failed")
	}

	if err := ent.Update(1, 11, 6); err != nil {
		log.Fatal("update failed")
	}

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}

func TestSQLGet(t *testing.T) {

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}

	ent := NewSQL(TestEntitySQLTestFileName)

	item := Item{
		UserName:   "hoge",
		NumOfGames: 10,
		NumOfWins:  5,
	}

	if err := ent.Add(&item); err != nil {
		log.Fatal("registration failed")
	}

	_, item, err := ent.Get("hoge")
	if err != nil {
		log.Fatal("update failed.", err)
	}
	switch {
	case item.UserName != "hoge":
		log.Fatal("UserName does not match")
	case item.NumOfGames != 10:
		log.Fatal("NumOfGames does not match")
	case item.NumOfWins != 5:
		log.Fatal("NumOfWins does not match")
	}

	// ファイルを削除する
	if _, err := os.Stat(TestEntitySQLTestFileName); !os.IsNotExist(err) {
		if err = os.Remove(TestEntitySQLTestFileName); err != nil {
			log.Fatal("failed to delete the file.")
		}
	}
}

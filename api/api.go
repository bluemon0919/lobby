package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bluemon0919/lobby/database"
	"github.com/bluemon0919/lobby/sessions"
	"github.com/bluemon0919/lobby/websocket"
	"github.com/gin-gonic/gin"
)

// WebAPI is webAPI information
type WebAPI struct {
	sql *database.EntitySQL
}

// NewWebAPI create new WebAPI
func NewWebAPI(sql *database.EntitySQL) *WebAPI {
	return &WebAPI{
		sql: sql,
	}
}

// PlayersGET gets player information
func (w *WebAPI) PlayersGET(c *gin.Context) {
	sessionID := c.Param("SessionID")
	hubManger := websocket.NewManager()
	hub, _ := hubManger.Get(sessionID)
	usersSessionIDs := hubManger.Users(hub) // users sessionID list

	datas := make(map[string]interface{})
	for _, id := range usersSessionIDs {
		player := sessions.GetPlayer(id)
		datas[id] = player.Name
	}

	// sessionID : PlayerNameのリスト(datas)を渡す
	c.JSON(http.StatusOK, datas)
}

// ResultNotify notices game result
// PlayerName(string)とResult(string)を受け取ってSQLに書き込む
// Resuleは"1"なら勝ち、"0"なら負けを表す
func (w *WebAPI) ResultNotify(c *gin.Context) {
	playerName := c.Param("PlayerName")
	result := c.Param("Result")

	_, item, err := w.sql.Get(playerName)
	if err != nil {
		log.Fatal(err)
	}
	if item.IsEmpty() {
		item := database.Item{
			UserName:   playerName,
			NumOfGames: 0,
			NumOfWins:  0,
		}
		w.sql.Add(&item)
	}

	key, item, err := w.sql.Get(playerName)
	if err != nil {
		c.String(http.StatusInternalServerError, "database update error")
		return
	}
	tmpNumOfGames := item.NumOfGames
	tmpNumOfWins := item.NumOfWins

	item.NumOfGames++
	if result == "1" {
		item.NumOfWins++
	}

	fmt.Print("PlayerName:", playerName)
	fmt.Println(" Result:", result)
	fmt.Println(" NumOfGames:", tmpNumOfGames, "->", item.NumOfGames)
	fmt.Println(" NumOfWins:", tmpNumOfWins, "->", item.NumOfWins)

	err = w.sql.Update(key, item.NumOfGames, item.NumOfWins)
	if err != nil {
		c.String(http.StatusInternalServerError, "database update error")
		return
	}
	c.String(http.StatusOK, "result set")
}

// GetHistory get player's games history.
// PlayerIDを基にプレイヤーのゲーム履歴を取得する
func (w *WebAPI) GetHistory(c *gin.Context) {
	playerID := c.Param("SessionID")
	player := sessions.GetPlayer(playerID)
	fmt.Println("id:", playerID, "name:", player)

	_, item, err := w.sql.Get(player.Name)
	if err != nil {
		c.String(http.StatusInternalServerError, "database get error")
		return
	}
	tmpNumOfGames := item.NumOfGames
	tmpNumOfWins := item.NumOfWins

	// key,value形式のJsonデータを作る
	datas := make(map[string]interface{})
	datas["player"] = player.Name
	datas["games"] = tmpNumOfGames
	datas["wins"] = tmpNumOfWins

	c.JSON(http.StatusOK, datas)
}

package api

import (
	"net/http"

	"github.com/bluemon0919/lobby/sessions"
	"github.com/bluemon0919/lobby/websocket"
	"github.com/gin-gonic/gin"
)

type jsonData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func PlayersGET(c *gin.Context) {
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

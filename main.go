package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bluemon0919/lobby/api"
	"github.com/bluemon0919/lobby/sessions"
	"github.com/bluemon0919/lobby/websocket"
	"github.com/gin-gonic/gin"
)

const cookieName = "gameid"

var addr = flag.String("addr", ":8080", "http service address")

func serveLobby(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/room.html")
}

func serveLoginHandler(c *gin.Context) {
	account := c.PostForm("account")
	if "" == account {
		http.Error(c.Writer, "account not set", http.StatusNotFound)
		return
	}

	// セッションを開始
	manager := sessions.NewManager()
	session, err := manager.Start(c.Writer, c.Request, cookieName)
	if err != nil {
		http.Error(c.Writer, "session start faild", http.StatusMethodNotAllowed)
		return
	}
	session.Set("account", account)
	if err := session.Save(); err != nil {
		http.Error(c.Writer, "session save faild", http.StatusMethodNotAllowed)
		return
	}
	player := sessions.NewPlayer(account)
	sessions.SetPlayer(session.ID, *player)

	hubmaneger := websocket.NewManager()
	hub, err := hubmaneger.Get(session.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(c.Writer, "websocket hub get faild", http.StatusMethodNotAllowed)
		return
	}

	if hubmaneger.Count(hub) >= 2 {
		us := hubmaneger.Users(hub)
		u1, u := us[0], us[0]
		u2 := us[1]
		hub.Boardcast([]byte("&u1=" + u1 + "&u2=" + u2 + "&u=" + u))
		// PlayRoomを返す
		http.Redirect(c.Writer, c.Request, "/play?u1="+u2+"&u2="+u1+"&u="+u, http.StatusFound)
	} else {
		// LobbyRoomを返す
		http.Redirect(c.Writer, c.Request, "/lobby", http.StatusMovedPermanently)
	}
}

func serveWebsocket(c *gin.Context) {
	manager := sessions.NewManager()
	session, err := manager.Get(c.Request, cookieName)
	if err != nil {
		return
	}
	hubManger := websocket.NewManager()
	hub, _ := hubManger.Get(session.ID)
	websocket.ServeWs(hub, c.Writer, c.Request)
}

func huboutHandler(c *gin.Context) {
	manager := sessions.NewManager()
	session, err := manager.Get(c.Request, cookieName)
	if err != nil {
		fmt.Println(err)
		return
	}
	hubManger := websocket.NewManager()
	hubManger.Destroy(session.ID)
}

func main() {
	flag.Parse()
	router := gin.Default()
	v1 := router.Group("")
	{
		v1.GET("/", func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, "static/front.html")
		})
		v1.GET("/lobby", func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, "static/room.html")
		})
		v1.GET("/play", func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, "static/connect4.html")
		})
		v1.GET("/ws", serveWebsocket)
		v1.POST("/login", serveLoginHandler)
		v1.POST("/hubout", huboutHandler)
		v1.GET("/players/:SessionID", api.PlayersGET)
	}
	router.Run(":8080")

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

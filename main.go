// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bluemon0919/lobby/sessions"
	"github.com/bluemon0919/lobby/websocket"
	"github.com/rs/xid"
)

const cookieName = "gameid"

var addr = flag.String("addr", ":8080", "http service address")

var count = 0
var userList []string
var flg = false

type User struct {
	name       string
	id         string
	roomNumber int
}

// Room manages game room user information.
type Room struct {
	users []User
	count int
}

func newRoom() *Room {
	return &Room{}
}

func serveFront(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/front.html")
}

func serveLobby(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/room.html")
}

func serveLoginHandler(room *Room, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if "" == r.FormValue("account") {
		http.Error(w, "account not set", http.StatusNotFound)
		return
	}

	// セッションを開始
	manager := sessions.NewManager()
	session, err := manager.Start(w, r, cookieName)
	if err != nil {
		http.Error(w, "session start faild", http.StatusMethodNotAllowed)
		return
	}
	session.Set("account", r.FormValue("account"))
	if err := session.Save(); err != nil {
		http.Error(w, "session save faild", http.StatusMethodNotAllowed)
		return
	}

	// ここから先を別のハンドラにする
	// セッションはつないでいるので、session.Get(sessionID)でデータは取れるはず
	hubmaneger := websocket.NewManager()
	hub, _ := hubmaneger.Get(session.ID)

	// 待機用のページを返す
	if hubmaneger.Count(hub) >= 2 {
		us := hubmaneger.Users(hub)
		u1, u := us[0], us[0]
		u2 := us[1]
		hub.Boardcast([]byte("&u1=" + u1 + "&u2=" + u2 + "&u=" + u))
		http.Redirect(w, r, "/play?u1="+u2+"&u2="+u1+"&u="+u, http.StatusFound)
	}
	http.Redirect(w, r, "/lobby", http.StatusMovedPermanently)
}

func (r *Room) id() string {
	guid := xid.New()
	return guid.String()
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	manager := sessions.NewManager()
	session, err := manager.Get(r, cookieName)
	if err != nil {
		return
	}
	hubManger := websocket.NewManager()
	hub, _ := hubManger.Get(session.ID)
	websocket.ServeWs(hub, w, r)
}

func servePlay(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/connect4.html")
}

func main() {
	flag.Parse()
	room := newRoom()
	http.HandleFunc("/", serveFront)
	http.HandleFunc("/lobby", serveLobby)
	http.HandleFunc("/play", servePlay)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		serveLoginHandler(room, w, r)
	})
	http.HandleFunc("/ws", serveWebsocket)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

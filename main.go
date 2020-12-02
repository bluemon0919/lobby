// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
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

type RoomManager struct {
	rooms []OneRoom
}

// OneRoom は各ゲームルームの情報です
type OneRoom struct {
	// ルームナンバー
	id string
	// 入室中のユーザーリスト
	users []User
	// 通信ハブ
	hub *websocket.Hub
}

func newRoom() *Room {
	return &Room{}
}

func serveFront(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "static/front.html")
}

func serveLobby(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/lobby" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "static/room.html")
}

func serveLoginHandler(room *Room, w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/login" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	fmt.Println("Method:", r.Method)
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("account=", r.FormValue("account"))
	if "" == r.FormValue("account") {
		http.Error(w, "account not set", http.StatusNotFound)
		return
	}

	// セッションを開始
	manager := sessions.NewManager()
	session, err := manager.Get(r, sessions.DefaultCookieName)
	if err != nil {
		session, err = manager.New(w, r, sessions.DefaultCookieName)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "session get faild", http.StatusMethodNotAllowed)
			return
		}
	}
	session.Set("account", r.FormValue("account"))
	//session.Set("hub", hub)
	if err := session.Save(); err != nil {
		http.Error(w, "session save faild", http.StatusMethodNotAllowed)
		return
	}

	hubmaneger := websocket.NewManager()
	hub, _ := hubmaneger.Get(session.ID)

	// 待機用のページを返す
	if hubmaneger.Count(hub) >= 2 {
		us := hubmaneger.Users(hub)
		fmt.Println("us:", us)
		u1, u := us[0], us[0]
		u2 := us[1]
		hub.Boardcast([]byte("&u1=" + u1 + "&u2=" + u2 + "&u=" + u))
		fmt.Println("casted")
		http.Redirect(w, r, "/play?u1="+u2+"&u2="+u1+"&u="+u, http.StatusFound)
	}
	http.Redirect(w, r, "/lobby", http.StatusMovedPermanently)
}

func serverRoomTop(w http.ResponseWriter, r *http.Request) {
	// セッションを取得
	manager := sessions.NewManager()
	session, err := manager.Get(r, sessions.DefaultCookieName)
	if err != nil {
		http.Error(w, "session get faild", http.StatusMethodNotAllowed)
		return
	}

	if _, ok := session.Values["account"]; ok {
		fmt.Println("登録済み")
	}
}

func (r *Room) id() string {
	guid := xid.New()
	return guid.String()
}

func (r *Room) roomNumber() int {
	return 1 // 固定
}

// Number returns the number of people in the room
func (r *Room) Number(roomNum int) (int, error) {
	num := 0
	for _, u := range r.users {
		if roomNum == u.roomNumber {
			num++
		}
	}
	return num, nil
}

func (r *Room) Register(name string) (string, error) {
	for _, u := range r.users {
		if u.name == name {
			return u.id, fmt.Errorf("already registered")
		}
	}

	id := r.id()
	r.users = append(r.users, User{
		name:       name,
		id:         id,
		roomNumber: r.roomNumber(),
	})
	log.Println("add id:", name)
	return id, nil
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	manager := sessions.NewManager()
	session, err := manager.Get(r, sessions.DefaultCookieName)
	if err != nil {
		return
	}
	hubManger := websocket.NewManager()
	hub, _ := hubManger.Get(session.ID)
	websocket.ServeWs(hub, w, r)
}

func servePlay(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/play" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
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

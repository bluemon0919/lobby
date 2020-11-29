// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bluemon0919/websocket/connect4/sessions"
	"github.com/rs/xid"
)

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
	hub *Hub
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

	http.ServeFile(w, r, "front.html")
}

func serveRoom(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/room" {
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
	//router := gin.Default()
	manager := sessions.NewManager()
	//router.Use(sessions.StartDefaultSession(manager))
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

	/*
		hubmaneger := NewManager()
		hub := hubmaneger.Get(session.ID)
		hubmaneger.Save(session.ID, hub)
	*/
	/*
		var err error
		if name := r.FormValue("account"); 0 != len(name) {
			/// UserNameをサーバーに登録
			/// UesrName,id,RoomNumberをセットで管理する
			_, err = room.Register(name)
			if err != nil {
				http.Error(w, "Not found ID", http.StatusNotFound)
				return
			}
		} else {
			http.Error(w, "Not found ID", http.StatusNotFound)
			return
		}
	*/

	// 待機用のページを返す
	//http.ServeFile(w, r, "static/room.html?u1="+r.FormValue("account"))
	if len(hub.clients) >= 2 {
		http.Redirect(w, r, "/connect4", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/room", http.StatusMovedPermanently)
	}

	/*
		// メンバーが揃ったら別スレッドからWebSocket通信で通知する
		roomNumber := 1
		n, err := room.Number(roomNumber)
		if err != nil {
			http.Error(w, "room number not found", http.StatusNotFound)
			return
		}
		if n >= 2 {
			u := room.users[0].name
			u1 := room.users[0].name
			u2 := room.users[1].name
			go func() {
				time.Sleep(time.Second * 2)
				hub.broadcast <- []byte("u1=" + u1 + "&u2=" + u2 + "&u=" + u)
			}()
		}
	*/

	/*
		if len(room.userList) >= 2 {
			// もう一人に通知する
			if !flg {
				hub.broadcast <- []byte("揃った")
				flg = true
			}
			u := room.userList[0]
			u2 := room.userList[1]
			if u2 == r.FormValue("id") {
				u2 = room.userList[0]
			}
			http.Redirect(w, r, "/connect4?u1="+r.FormValue("id")+"&u2="+u2+"&u="+u, http.StatusFound)
		} else {
			http.ServeFile(w, r, "static/room.html")
		}
		//tpl := template.Must(template.ParseFiles("static/front.html"))
		//tpl.Execute(w, nil)
	*/
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

func serveConnet4(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/connect4" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/connect4.html")
}

var hub *Hub

func main() {
	flag.Parse()
	hub = newHub()
	go hub.run()
	room := newRoom()
	http.HandleFunc("/", serveFront)
	http.HandleFunc("/room", serveRoom)
	http.HandleFunc("/connect4", serveConnet4)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		serveLoginHandler(room, w, r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

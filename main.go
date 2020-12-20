package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bluemon0919/lobby/sessions"
	"github.com/bluemon0919/lobby/websocket"
)

const cookieName = "gameid"

var addr = flag.String("addr", ":8080", "http service address")

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

func serveLoginHandler(w http.ResponseWriter, r *http.Request) {
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

	hubmaneger := websocket.NewManager()
	hub, err := hubmaneger.Get(session.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "websocket hub get faild", http.StatusMethodNotAllowed)
		return
	}

	if hubmaneger.Count(hub) >= 2 {
		us := hubmaneger.Users(hub)
		u1, u := us[0], us[0]
		u2 := us[1]
		hub.Boardcast([]byte("&u1=" + u1 + "&u2=" + u2 + "&u=" + u))
		// PlayRoomを返す
		http.Redirect(w, r, "/play?u1="+u2+"&u2="+u1+"&u="+u, http.StatusFound)
	} else {
		// LobbyRoomを返す
		http.Redirect(w, r, "/lobby", http.StatusMovedPermanently)
	}
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

func huboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	manager := sessions.NewManager()
	session, err := manager.Get(r, cookieName)
	if err != nil {
		fmt.Println(err)
		return
	}
	hubManger := websocket.NewManager()
	hubManger.Destroy(session.ID)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", serveFront)
	http.HandleFunc("/lobby", serveLobby)
	http.HandleFunc("/play", servePlay)
	http.HandleFunc("/login", serveLoginHandler)
	http.HandleFunc("/ws", serveWebsocket)
	http.HandleFunc("/hubout", huboutHandler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

var count = 0
var userList []string
var flg = false

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
	if len(r.FormValue("id")) == 0 {
		http.Error(w, "Not found ID", http.StatusNotFound)
		return
	}
	found := false
	for _, u := range userList {
		if u == r.FormValue("id") {
			found = true
		}
	}
	if !found {
		userList = append(userList, r.FormValue("id"))
		fmt.Println("add id:", r.FormValue("id"))
	}
	if len(userList) >= 2 {
		// もう一人に通知する
		if !flg {
			hub.broadcast <- []byte("揃った")
			flg = true
		}
		u := userList[0]
		u2 := userList[1]
		if u2 == r.FormValue("id") {
			u2 = userList[0]
		}
		http.Redirect(w, r, "/connect4?u1="+r.FormValue("id")+"&u2="+u2+"&u="+u, http.StatusFound)
	} else {
		http.ServeFile(w, r, "static/room.html")
	}
	//tpl := template.Must(template.ParseFiles("static/front.html"))
	//tpl.Execute(w, nil)
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
	http.HandleFunc("/", serveFront)
	http.HandleFunc("/room", serveRoom)
	http.HandleFunc("/connect4", serveConnet4)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

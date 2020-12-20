package websocket

import (
	"fmt"
	"net/http"
)

// Manager manages websocket hubs
type Manager struct {
	database map[string]*Hub
	pool     []*Hub
	count    map[*Hub]int
	users    map[*Hub][]string
}

// PoolMax プールに生成するHubの数
const PoolMax = 2

var m Manager

func init() {
	m.database = make(map[string]*Hub, PoolMax*2)
	m.count = make(map[*Hub]int, PoolMax)
	m.users = make(map[*Hub][]string, PoolMax)
	for i := 0; i < PoolMax; i++ {
		hub := newHub()
		m.pool = append(m.pool, hub)
		go hub.run()
	}
}

// NewManager websocket hub managerを取得する
func NewManager() *Manager {
	return &m
}

// Get websocket hubを取得する
func (m *Manager) Get(key string) (*Hub, error) {
	if hub, ok := m.database[key]; ok {
		return hub, nil
	}
	if len(m.pool) <= 0 {
		return nil, fmt.Errorf("no hub")
	}

	hub := m.pool[0] // 先頭のHubを割り当てる
	m.database[key] = hub
	m.count[hub]++
	if m.count[hub] >= PoolMax {
		m.pool = m.pool[1:] // 先頭のHubを除外する
	}
	m.users[hub] = append(m.users[hub], key)
	return hub, nil
}

// Count hubの割当人数を取得する
func (m *Manager) Count(hub *Hub) int {
	return m.count[hub]
}

// Users ユーザ名を取得する
func (m *Manager) Users(hub *Hub) []string {
	return m.users[hub]
}

// Destroy キーを削除する
func (m *Manager) Destroy(key string) {
	if hub, exist := m.database[key]; exist {
		delete(m.database, key)
		users := m.Users(hub)
		fmt.Println(users)
		for u := 0; u < len(users); u++ {
			if users[u] == key {
				m.users[hub] = remove(m.users[hub], u)
				users = m.users[hub]
			}
		}
		if 0 == len(users) {
			// ペアが解放されたらhubをpoolに戻す
			delete(m.count, hub)
			m.pool = append(m.pool, hub)
		} else {
			// ペアが解放されるのを待つ
		}
	}
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// Boardcast wrap boradcast
func (h *Hub) Boardcast(bytes []byte) {
	h.broadcast <- bytes
}

// ServeWs wrap serveWs
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	serveWs(hub, w, r)
}

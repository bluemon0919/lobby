package websocket

import "fmt"

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

func NewManager() *Manager {
	return &m
}

func (m *Manager) Get(key string) (*Hub, error) {
	if hub, ok := m.database[key]; ok {
		return hub, nil
	}
	if len(m.pool) <= 0 {
		return nil, fmt.Errorf("no hub")
	}

	hub := m.pool[0] // 先頭のHubを返す
	m.database[key] = hub
	m.count[hub]++
	if m.count[hub] >= 2 {
		m.pool = m.pool[1:] // 先頭のHubを除外する
	}
	m.users[hub] = append(m.users[hub], key)
	return hub, nil
}

func (m *Manager) Count(hub *Hub) int {
	return m.count[hub]
}

func (m *Manager) Users(hub *Hub) []string {
	return m.users[hub]
}

func (m *Manager) Destroy(key string) {
	delete(m.database, key)
}

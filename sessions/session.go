package sessions

import "net/http"

// Session セッション情報
type Session struct {
	cookieName string
	ID         string
	manager    *Manager
	Values     map[string]interface{}
	request    *http.Request
	writer     http.ResponseWriter
}

// NewSession 新しいセッションを生成する
func NewSession(manager *Manager, cookieName string) *Session {
	return &Session{
		cookieName: cookieName,
		manager:    manager,
		Values:     map[string]interface{}{},
	}
}

// Save セッションを保存する
func (s *Session) Save() error {
	return s.manager.Save(s)
}

// CookieName Cookie名を取得する
func (s *Session) CookieName() string {
	return s.cookieName
}

// Get セッション変数を取得する
func (s *Session) Get(key string) (interface{}, bool) {
	ret, exists := s.Values[key]
	return ret, exists
}

// Set セッション変数に設定する
func (s *Session) Set(key string, val interface{}) {
	s.Values[key] = val
}

// Delete セッション変数を削除する
func (s *Session) Delete(key string) {
	delete(s.Values, key)
}

// Terminate セッションを終了する
func (s *Session) Terminate() {
	s.manager.Destroy(s.ID)
}

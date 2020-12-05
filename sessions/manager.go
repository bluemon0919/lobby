package sessions

import (
	"errors"
	"net/http"

	"github.com/rs/xid"
)

// Manager manages session
type Manager struct {
	database map[string]interface{}
}

var m Manager

func init() {
	m.database = map[string]interface{}{}
}

func NewManager() *Manager {
	return &m
}

// 新しいセッションIDを発行する
func (m *Manager) sessionID() string {
	guid := xid.New()
	return guid.String()
}

// セッションIDが存在するか確認する
func (m *Manager) Exists(sessionID string) bool {
	_, ok := m.database[sessionID]
	return ok
}

func (m *Manager) Start(w http.ResponseWriter, r *http.Request, cookieName string) (*Session, error) {
	session, err := m.Get(r, cookieName)
	if err != nil {
		session, err = m.New(w, r, cookieName)
		if err != nil {
			http.Error(w, "session get faild", http.StatusMethodNotAllowed)
			return nil, err
		}
	}
	return session, nil
}

func (m *Manager) Save(r *http.Request, w http.ResponseWriter, session *Session) error {
	m.database[session.ID] = session

	c := &http.Cookie{
		Name:  session.CookieName(),
		Value: session.ID,
		Path:  "/",
	}
	http.SetCookie(w, c)
	//ttp.SetCookie(session.writer, c)
	return nil
}

//func (m *Manager) New(r *http.Request, cookieName string) (*Session, error) {
func (m *Manager) New(w http.ResponseWriter, r *http.Request, cookieName string) (*Session, error) {
	cookie, err := r.Cookie(cookieName)
	if err == nil && m.Exists(cookie.Value) {
		return nil, errors.New("sessionIDはすでに発行されています")
	}
	session := NewSession(m, cookieName)
	session.ID = m.sessionID()
	session.request = r
	session.writer = w
	return session, nil
}

// セッションを取得する
func (m *Manager) Get(r *http.Request, cookieName string) (*Session, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		// リクエストからcookie情報を取得できない場合
		return nil, err
	}

	sessionID := cookie.Value
	// cookie情報からセッション情報を取得
	buffer, exists := m.database[sessionID]
	if !exists {
		return nil, errors.New("無効なセッションIDです")
	}

	session := buffer.(*Session)
	session.request = r // ここでフロントの情報をサーバーに格納しているっぽい⭐️
	return session, nil
}

// セッションを破棄する
func (m *Manager) Destroy(sessionID string) {
	delete(m.database, sessionID)
}

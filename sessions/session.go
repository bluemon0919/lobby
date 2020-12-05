package sessions

import (
	"net/http"
)

type Session struct {
	cookieName string
	ID         string
	manager    *Manager
	request    *http.Request
	writer     http.ResponseWriter
	Values     map[string]interface{}
}

func NewSession(manager *Manager, cookieName string) *Session {
	return &Session{
		cookieName: cookieName,
		manager:    manager,
		Values:     map[string]interface{}{},
	}
}

func (s *Session) Save() error {
	return s.manager.Save(s.request, s.writer, s)
}

func (s *Session) CookieName() string {
	return s.cookieName
}

func (s *Session) Get(key string) (interface{}, bool) {
	ret, exists := s.Values[key]
	return ret, exists
}

func (s *Session) Set(key string, val interface{}) {
	s.Values[key] = val
}

func (s *Session) Delete(key string) {
	delete(s.Values, key)
}

func (s *Session) Terminate() {
	s.manager.Destroy(s.ID)
}

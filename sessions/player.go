package sessions

// Player is player informaion
type Player struct {
	Name string
}

var p map[string]Player

func init() {
	p = make(map[string]Player)
}

// NewPlayer create new player information
func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
	}
}

// SetPlayer sets player information
func SetPlayer(sessionID string, player Player) {
	p[sessionID] = player
}

// GetPlayer gets player information
func GetPlayer(sessionID string) Player {
	return p[sessionID]
}

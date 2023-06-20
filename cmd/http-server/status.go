package main

type PlayerStatus struct {
	X  float32 `json:"x"`
	Y  float32 `json:"y"`
	VX float32 `json:"vx"`
	VY float32 `json:"vy"`
}

func DefaultPlayerStatus() PlayerStatus {
	return PlayerStatus{X: 0, Y: 0, VX: 0, VY: 0}
}

type Statuses struct {
	Players map[string]PlayerStatus `json:"players"`
}

func DefaultStatuses() Statuses {
	return Statuses{Players: make(map[string]PlayerStatus)}
}

// Set status if not exists
func (s *Statuses) SetPlayerStatusIfNotExists(userID string) {
	if _, ok := s.Players[userID]; !ok {
		s.Players[userID] = DefaultPlayerStatus()
	}
}

package main

type Room struct {
	Exist    bool                    `json:"exist"`
	Host     User                    `json:"host"`
	Players  []User                  `json:"players"`
	Actions  map[string]*ActionQueue `json:"actions"`
	Statuses map[string]Status       `json:"statuses"`
}

// Get default room
func defaultRoom() Room {
	return Room{Exist: false, Host: User{}, Players: []User{}, Actions: make(map[string]*ActionQueue), Statuses: make(map[string]Status)}
}

// Add user to room
func (room *Room) AddUser(user User) {
	if len(room.Players) == 0 {
		room.Exist = true
		room.Host = user
	}
	room.Players = append(room.Players, user)
	room.Actions[user.UserID] = DefaultActionQueue()
	room.Statuses[user.UserID] = Status{}
}

// Remove user from room
func (room *Room) RemoveUser(user User) {
	for i, player := range room.Players {
		if player.UserID == user.UserID {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			break
		}
	}

	// delete nomatter exist or not
	delete(room.Actions, user.UserID)
	delete(room.Statuses, user.UserID)

	if len(room.Players) == 0 {
		room.Exist = false
	}
}

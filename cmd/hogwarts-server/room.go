package main

type Room struct {
	Actions  map[string]*ActionQueue `json:"actions"`
	Statuses Statuses                `json:"statuses"`
}

// Get default room
func defaultRoom() Room {
	return Room{Actions: make(map[string]*ActionQueue), Statuses: DefaultStatuses()}
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const (
	MAX_ROOMS = 1000
)

var (
	rooms      [MAX_ROOMS]Room
	mu_rooms   [MAX_ROOMS]sync.Mutex
	n_rooms    = 0
	mu_n_rooms sync.Mutex
)

func init() {
	for i := 0; i < MAX_ROOMS; i++ {
		rooms[i] = defaultRoom()
	}

	// Create room 0 with host "1" and player "1", "2"
	actions := make(map[string]*ActionQueue)
	statuses := DefaultStatuses()
	actions["1"] = DefaultActionQueue()
	statuses.Players["1"] = DefaultPlayerStatus()
	actions["2"] = DefaultActionQueue()
	statuses.Players["2"] = DefaultPlayerStatus()
	room := Room{Actions: actions, Statuses: statuses}

	mu_rooms[0].Lock()
	defer mu_rooms[0].Unlock()
	rooms[0] = room
	n_rooms = 1
}

func main() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/createroom", createRoom)
	http.HandleFunc("/action", action)
	http.HandleFunc("/update", update)
	http.HandleFunc("/getactions", getActions)
	http.HandleFunc("/getstatuses", getStatuses)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var res struct {
		Message string `json:"message"`
	}

	res.Message = "pong"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		UserID string `json:"userID"`
	}
	var res struct {
		RoomID int `json:"roomID"`
	}

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	mu_n_rooms.Lock()
	roomID := n_rooms
	n_rooms = (n_rooms + 1) % MAX_ROOMS
	mu_n_rooms.Unlock()

	mu_rooms[n_rooms].Lock()
	defer mu_rooms[n_rooms].Unlock()

	// Create the room
	rooms[roomID] = defaultRoom()

	// Return the room
	res.RoomID = roomID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func action(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		UserID string `json:"userID"`
		RoomID int    `json:"roomID"`
		Action Action `json:"action"`
	}
	var res struct {
	}

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	mu_rooms[req.RoomID].Lock()
	defer mu_rooms[req.RoomID].Unlock()

	// Check if action queue exists
	if _, ok := rooms[req.RoomID].Actions[req.UserID]; !ok {
		rooms[req.RoomID].Actions[req.UserID] = DefaultActionQueue()
	}
	// Update the action
	rooms[req.RoomID].Actions[req.UserID].Push(req.Action)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		UserID   string   `json:"userID"`
		RoomID   int      `json:"roomID"`
		Statuses Statuses `json:"statuses"`
	}
	var res struct {
	}

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	mu_rooms[req.RoomID].Lock()
	defer mu_rooms[req.RoomID].Unlock()

	// Update all player statuses
	rooms[req.RoomID].Statuses = req.Statuses

	for userID := range req.Statuses.Players {
		// if no player action is set, set it to default
		if _, ok := rooms[req.RoomID].Actions[userID]; !ok {
			rooms[req.RoomID].Actions[userID] = DefaultActionQueue()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func getActions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		RoomID int `json:"roomID"`
	}
	var res map[string]Action = make(map[string]Action)

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	mu_rooms[req.RoomID].Lock()
	defer mu_rooms[req.RoomID].Unlock()

	for userID, action := range rooms[req.RoomID].Actions {
		res[userID] = action.Pop()
	}

	// Return all player actions
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func getStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		RoomID int `json:"roomID"`
	}
	var res Statuses

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	mu_rooms[req.RoomID].Lock()
	defer mu_rooms[req.RoomID].Unlock()

	res = rooms[req.RoomID].Statuses

	// Return all player statuses
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

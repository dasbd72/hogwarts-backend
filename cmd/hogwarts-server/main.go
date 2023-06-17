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
	actions := make(map[string]Action)
	statuses := make(map[string]Status)
	actions["1"] = defaultAction()
	statuses["1"] = defaultStatus()
	actions["2"] = defaultAction()
	statuses["2"] = defaultStatus()
	room := Room{Exist: true, Host: User{UserID: "1"}, Players: []User{{UserID: "1"}, {UserID: "2"}}, Actions: actions, Statuses: statuses}

	mu_rooms[0].Lock()
	defer mu_rooms[0].Unlock()
	rooms[0] = room
	n_rooms = 1
}

func main() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/createroom", createRoom)
	http.HandleFunc("/joinroom", joinRoom)
	http.HandleFunc("/leaveroom", leaveRoom)
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
	n_rooms++
	mu_n_rooms.Unlock()

	mu_rooms[n_rooms].Lock()
	defer mu_rooms[n_rooms].Unlock()

	// Create the room
	rooms[roomID] = defaultRoom()
	rooms[roomID].AddUser(User{UserID: req.UserID})

	// Return the room
	res.RoomID = roomID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		UserID string `json:"userID"`
		RoomID int    `json:"roomID"`
	}
	var res struct {
		RoomID int `json:"roomID"`
	}

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	mu_rooms[req.RoomID].Lock()
	defer mu_rooms[req.RoomID].Unlock()

	// Add the user to the room
	rooms[req.RoomID].AddUser(User{UserID: req.UserID})

	// Return the room
	res.RoomID = req.RoomID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func leaveRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		UserID string `json:"userID"`
		RoomID int    `json:"roomID"`
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

	// Remove the user from the room
	rooms[req.RoomID].RemoveUser(User{UserID: req.UserID})

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

	// // Check if the room exists
	// if !rooms[req.RoomID].Exist {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// Update the action
	rooms[req.RoomID].Actions[req.UserID] = req.Action
	// if no player status is set, set it to default
	if _, ok := rooms[req.RoomID].Statuses[req.UserID]; !ok {
		rooms[req.RoomID].Statuses[req.UserID] = defaultStatus()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		UserID   string            `json:"userID"`
		RoomID   int               `json:"roomID"`
		Statuses map[string]Status `json:"statuses"`
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

	// // Check if the room exists
	// if !rooms[req.RoomID].Exist {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// // Check if the user is the host
	// if req.UserID != rooms[req.RoomID].Host.UserID {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// Update all player statuses
	for userID, status := range req.Statuses {
		rooms[req.RoomID].Statuses[userID] = status
		// if no player action is set, set it to default
		if _, ok := rooms[req.RoomID].Actions[userID]; !ok {
			rooms[req.RoomID].Actions[userID] = defaultAction()
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
	var res map[string]Action

	// Parse the request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	mu_rooms[req.RoomID].Lock()
	defer mu_rooms[req.RoomID].Unlock()

	res = rooms[req.RoomID].Actions

	// Return all player actions
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func getStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		RoomID int `json:"roomID"`
	}
	var res map[string]Status

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

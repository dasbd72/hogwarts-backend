package main

type Action struct {
	Left   bool `json:"left_move"`
	Right  bool `json:"right_move"`
	Jump   bool `json:"jump"`
	Ground bool `json:"on_ground"`
	Attack bool `json:"attack"`
}

// Get default action
func defaultAction() Action {
	return Action{Left: false, Right: false, Jump: false, Ground: false, Attack: false}
}

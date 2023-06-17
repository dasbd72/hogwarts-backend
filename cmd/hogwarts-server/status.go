package main

type Status struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func DefaultStatus() Status {
	return Status{X: 0, Y: 0}
}

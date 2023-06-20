package main

const (
	MAX_QUEUE int = 2
)

type Action struct {
	Left   bool `json:"left_move"`
	Right  bool `json:"right_move"`
	Jump   bool `json:"jump"`
	Ground bool `json:"on_ground"`
	Attack bool `json:"attack"`
}

// Get default action
func DefaultAction() Action {
	return Action{Left: false, Right: false, Jump: false, Ground: false, Attack: false}
}

// Or two actions
func (a *Action) Or(b Action) {
	a.Left = a.Left || b.Left
	a.Right = a.Right || b.Right
	a.Jump = a.Jump || b.Jump
	a.Ground = a.Ground || b.Ground
	a.Attack = a.Attack || b.Attack
}

type ActionQueue struct {
	Actions [MAX_QUEUE]Action `json:"actions"`
	Head    int               `json:"head"`
	Tail    int               `json:"tail"`
	Temp    Action            `json:"temp"`
}

// Get default action queue
func DefaultActionQueue() *ActionQueue {
	return &ActionQueue{Actions: [MAX_QUEUE]Action{}, Head: 0, Tail: 0}
}

// Push action to queue
func (aq *ActionQueue) Push(action Action) {
	if (aq.Head+1)%MAX_QUEUE == aq.Tail {
		// queue is full
		temp := aq.Pop()
		aq.Temp.Or(temp)
	}
	aq.Actions[aq.Head] = action
	aq.Head = (aq.Head + 1) % MAX_QUEUE
}

// Pop action from queue
func (aq *ActionQueue) Pop() Action {
	if aq.Head == aq.Tail {
		// queue is empty
		return DefaultAction()
	}
	action := aq.Actions[aq.Tail]
	aq.Tail = (aq.Tail + 1) % MAX_QUEUE

	action.Or(aq.Temp)
	aq.Temp = DefaultAction()
	return action
}

// Size of queue
func (aq *ActionQueue) Size() int {
	return (aq.Head - aq.Tail + MAX_QUEUE) % MAX_QUEUE
}

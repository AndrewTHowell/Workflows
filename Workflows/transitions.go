package Workflows

type Transition interface {
	StartState() State
	Input() Input
	EndState() State
}

func NewTransition(startState State, input Input, endState State) Transition {
	return &transition{
		startState: startState,
		input: input,
		endState: endState,
	}
}

type transition struct {
	startState State
	input Input
	endState State
}

func (t *transition) StartState() State {
	return t.startState
}

func (t *transition) Input() Input {
	return t.input
}

func (t *transition) EndState() State {
	return t.endState
}

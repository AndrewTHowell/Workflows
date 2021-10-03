package main

import (
	"fmt"

	"Workflows/Workflows"
)

// TODO: state events, draw FSM diagrams

func main() {
	alphabet := Workflows.NewAlphabet('a', 'b')

	stateA := Workflows.NewState("hello there")
	stateB := Workflows.NewState("general kenobi!")
	states := []Workflows.State{stateA, stateB}

	transitions := []Workflows.Transition{
		Workflows.NewTransition(stateA, 'a', stateA),
		Workflows.NewTransition(stateA, 'a', stateB),
	}

	fsm, err := Workflows.NewFSM(alphabet, states, stateA, []Workflows.State{stateB}, transitions)
	if err != nil {
		panic(err)
	}

	final, err := fsm.Inputs('a')
	if err != nil {
		panic(err)
	}

	fmt.Println(final)
}

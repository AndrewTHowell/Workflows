package main

import (
	"fmt"

	"Workflows/Workflows"
)

// TODO: state events, non-deterministic FSM (each state can take same inputs to different states),

func main() {
	alphabet := Workflows.NewAlphabet('a', 'b')

	stateA := Workflows.NewState("hello there")
	stateB := Workflows.NewState("general kenobi!")
	states := []Workflows.State{stateA, stateB}

	transitions := []Workflows.Transition{
		Workflows.NewTransition(stateA, 'b', stateB),
		Workflows.NewTransition(stateA, 'a', stateA),
	}

	fsm, err := Workflows.NewFSM(alphabet, states, stateA, []Workflows.State{stateB}, transitions)
	if err != nil {
		panic(err)
	}

	final, err := fsm.Inputs('a', 'b')
	if err != nil {
		panic(err)
	}

	fmt.Println(final)
}

package main

import (
	"fmt"

	"Workflows/Workflows"
)

// TODO: state events, non-deterministic FSM (each state can take same inputs to different states)

func main() {
	alphabet := Workflows.NewAlphabet('a', 'b', 'c')

	stateA := Workflows.NewState("hello there")
	stateB := Workflows.NewState("general kenobi!")

	fsm, err := Workflows.NewFSM(alphabet, []Workflows.State{stateA, stateB}, stateA, []Workflows.State{stateB})
	if err != nil {
		panic(err)
	}
	err = fsm.AddTransition(stateA, 'b', stateB)
	if err != nil {
		panic(err)
	}
	err = fsm.AddTransition(stateA, 'a', stateA)
	if err != nil {
		panic(err)
	}

	final, err := fsm.Inputs('a', 'b')
	if err != nil {
		panic(err)
	}

	fmt.Println(final)
}

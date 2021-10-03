package main

import (
	"fmt"

	"Workflows/Workflows"
)

// TODO: raw FSM diagrams, global context, state context

func main() {
	alphabet := Workflows.NewAlphabet('a', 'b')

	stateA := Workflows.NewState("hello there", func() { fmt.Println("entering ht") }, func() { fmt.Println("exiting ht") })
	stateB := Workflows.NewState("general kenobi!", func() { fmt.Println("entering gk") }, func() { fmt.Println("exiting gk") })
	states := []Workflows.State{stateA, stateB}

	transitions := []Workflows.Transition{
		Workflows.NewTransition(stateA, 'a', stateB),
		Workflows.NewTransition(stateB, 'b', stateB),
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

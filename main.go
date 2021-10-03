package main

import (
	"fmt"

	"Workflows/Workflows"
)

// TODO: global context, state context

func main() {
	config := &Workflows.Config{
		Alphabet: []Workflows.Input{'a', 'b'},
		States: []Workflows.StateConfig{
			{
				ID: "hello there",
				StartState: true,
				EntryEvent: func() { fmt.Println("entering ht") },
				ExitEvent: func() { fmt.Println("exiting ht") },
			},
			{
				ID: "general kenobi!",
				FinalState: true,
				EntryEvent: func() { fmt.Println("entering gk") },
				ExitEvent: func() { fmt.Println("exiting gk") },
			},
		},
		Transitions: []Workflows.TransitionConfig{
			{
				StartStateID: "hello there",
				Input: 'a',
				EndStateID: "general kenobi!",
			},
			{
				StartStateID: "general kenobi!",
				Input: 'b',
				EndStateID: "general kenobi!",
			},
		},
	}
	fsm, err := Workflows.NewFSM(config)
	if err != nil {
		panic(err)
	}

	final, err := fsm.Inputs('a', 'b')
	if err != nil {
		panic(err)
	}

	fmt.Println(final)
}

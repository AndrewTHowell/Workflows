package main

import (
	"fmt"

	"Workflows/Workflows"
)

// TODO: global context, state context, better error messages, refactor fsm.go to reduce complexity, all validation in parse...()

func main() {
	config := &Workflows.Config{
		Alphabet: []Workflows.Input{'a', 'b'},
		States: []Workflows.StateConfig{
			{
				ID:         "A",
				Name:       "State A",
				StartState: true,
				EntryEvent: func() { fmt.Println("entering A") },
				ExitEvent:  func() { fmt.Println("exiting A") },
			},
			{
				ID:         "B",
				Name:       "State B",
				FinalState: true,
				EntryEvent: func() { fmt.Println("entering B") },
				ExitEvent:  func() { fmt.Println("exiting B") },
			},
		},
		Transitions: []Workflows.TransitionConfig{
			{
				StartStateID: "A",
				Input:        'a',
				EndStateID:   "B",
			},
			{
				StartStateID: "B",
				Input:        'b',
				EndStateID:   "B",
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

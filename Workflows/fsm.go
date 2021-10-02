package Workflows

import (
	"fmt"

	"github.com/pkg/errors"
)

func NewFSM(alphabet Alphabet, states []State, startState State, finalStates []State) (*fsm, error) {
	fsmStates := make(map[State]struct{}, len(states))
	for _, state := range states {
		fsmStates[state] = struct{}{}
	}
	if _, ok := fsmStates[startState]; !ok {
		return nil, errors.New("'startState' must be a subset of 'states'")
	}
	fmsFinalStates := make(map[State]struct{}, len(finalStates))
	for _, finalState := range finalStates {
		if _, ok := fsmStates[finalState]; !ok {
			return nil, errors.New("'finalStates' must be a subset of 'states'")
		}
		fmsFinalStates[finalState] = struct{}{}
	}

	return &fsm{
		alphabet:     alphabet,
		states:       fsmStates,
		currentState: startState,
		finalStates:  fmsFinalStates,
		transitions:  map[State]map[Input]State{},
	}, nil

}

type fsm struct {
	alphabet     Alphabet
	currentState State
	states       map[State]struct{}
	finalStates  map[State]struct{}
	transitions  map[State]map[Input]State
}

func (fsm *fsm) AddTransition(startState State, input Input, endState State) error {
	if err := fsm.validateTransition(startState, input, endState); err != nil {
		return errors.Wrap(err, "invalid transition")
	}

	if _, ok := fsm.transitions[startState]; !ok {
		fsm.transitions[startState] = map[Input]State{}
	}

	fsm.transitions[startState][input] = endState
	return nil
}

func (fsm *fsm) validateTransition(startState State, input Input, endState State) error {
	if _, valid := fsm.states[startState]; !valid {
		return errors.New("startState not within states")
	}
	if !fsm.alphabet.Valid(input) {
		return errors.New("input not within alphabet")
	}
	if _, valid := fsm.states[endState]; !valid {
		return errors.New("endState not within states")
	}
	return nil
}

func (fsm *fsm) Inputs(inputs ...Input) (bool, error) {
	for i, input := range inputs {
		if _, err := fsm.Input(input); err != nil {
			return false, errors.Wrap(err, fmt.Sprintf("failure at input %d", i))
		}
	}
	return fsm.IsInFinalState(), nil
}

func (fsm *fsm) Input(input Input) (State, error) {
	if !fsm.alphabet.Valid(input) {
		return nil, errors.New("invalid input: not within alphabet")
	}

	nextState, valid := fsm.transitions[fsm.currentState][input]
	if !valid {
		return nil, errors.New("invalid input: invalid transition for current state")
	}

	fsm.currentState = nextState
	return fsm.currentState, nil
}

func (fsm *fsm) IsInFinalState() bool {
	_, finalState := fsm.finalStates[fsm.currentState]
	return finalState
}

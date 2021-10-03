package Workflows

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

type FSM interface {
	Inputs(inputs ...Input) (bool, error)
	Input(input Input) (State, error)
	IsInFinalState() bool
}

func NewFSM(alphabet Alphabet, states []State, startState State, finalStates []State, transitions []Transition) (FSM, error) {
	rand.Seed(time.Now().Unix())

	fsm := &fsm{
		alphabet: alphabet,
	}
	fsm.addStates(states)
	if err := fsm.addFinalStates(finalStates); err != nil {
		return nil, err
	}
	if err := fsm.addTransitions(transitions); err != nil {
		return nil, err
	}

	if _, ok := fsm.states[startState]; !ok {
		return nil, errors.New("'startState' must be a subset of 'states'")
	}
	fsm.setNewState(startState)

	return fsm, nil
}

type fsm struct {
	alphabet     Alphabet
	currentState State
	states       map[State]struct{}
	finalStates  map[State]struct{}
	transitions  map[State]map[Input][]State
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

	nextStates, valid := fsm.transitions[fsm.currentState][input]
	if !valid {
		return nil, errors.New("invalid input: invalid transition for current state")
	}
	nextStateIndex := rand.Intn(len(nextStates))
	nextState := nextStates[nextStateIndex]

	fsm.setNewState(nextState)
	return fsm.currentState, nil
}

func (fsm *fsm) setNewState(newState State) {
	// Only run exit event if current state is non-empty. Avoids panic on setting first state
	if fsm.currentState != nil {
		fsm.currentState.RunExitEvent()
	}
	fsm.currentState = newState
	fsm.currentState.RunEntryEvent()
	fmt.Println("> Current State: ", fsm.currentState)
}

func (fsm *fsm) IsInFinalState() bool {
	_, finalState := fsm.finalStates[fsm.currentState]
	return finalState
}

func (fsm *fsm) addStates(states []State) {
	fsm.states = make(map[State]struct{}, len(states))
	for _, state := range states {
		fsm.states[state] = struct{}{}
	}
}

func (fsm *fsm) addFinalStates(finalStates []State) error {
	fsm.finalStates = make(map[State]struct{}, len(finalStates))
	for _, finalState := range finalStates {
		if _, ok := fsm.states[finalState]; !ok {
			return errors.New("'finalStates' must be a subset of 'states'")
		}
		fsm.finalStates[finalState] = struct{}{}
	}
	return nil
}

func (fsm *fsm) addTransitions(transitions []Transition) error {
	fsm.transitions = map[State]map[Input][]State{}
	for _, transition := range transitions {
		err := fsm.addTransition(transition)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fsm *fsm) addTransition(transition Transition) error {
	if err := fsm.validateTransition(transition); err != nil {
		return errors.Wrap(err, "invalid transition")
	}

	if _, ok := fsm.transitions[transition.StartState()]; !ok {
		fsm.transitions[transition.StartState()] = map[Input][]State{}
	}

	endStates := fsm.transitions[transition.StartState()][transition.Input()]
	endStates = append(endStates, transition.EndState())
	fsm.transitions[transition.StartState()][transition.Input()] = endStates
	return nil
}

func (fsm *fsm) validateTransition(transition Transition) error {
	if _, valid := fsm.states[transition.StartState()]; !valid {
		return errors.New("startState not within states")
	}
	if !fsm.alphabet.Valid(transition.Input()) {
		return errors.New("input not within alphabet")
	}
	if _, valid := fsm.states[transition.EndState()]; !valid {
		return errors.New("endState not within states")
	}
	return nil
}

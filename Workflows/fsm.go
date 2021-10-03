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

func NewFSM(fsmConfig *Config) (FSM, error) {
	rand.Seed(time.Now().Unix())

	config, err := parseConfig(fsmConfig)
	if err != nil {
		return nil, err
	}

	return newFSM(config)
}

func newFSM(config *config) (FSM, error) {
	fsm := &fsm{
		alphabet: config.alphabet,
	}
	fsm.addStates(config.states)
	if err := fsm.addFinalStates(config.finalStates); err != nil {
		return nil, err
	}
	if err := fsm.addTransitions(config.transitions); err != nil {
		return nil, err
	}
	if err := fsm.setStartState(config.startState); err != nil {
		return nil, err
	}

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

	fsm.currentState.RunExitEvent()
	fsm.setNewState(nextState)
	return fsm.currentState, nil
}

func (fsm *fsm) setStartState(newState State) error {
	if _, ok := fsm.states[newState]; !ok {
		return errors.New("'startState' must be a subset of 'states'")
	}
	fsm.setNewState(newState)
	return nil
}

func (fsm *fsm) setNewState(newState State) {
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

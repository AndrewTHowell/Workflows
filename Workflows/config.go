package Workflows

import (
	"errors"
	"fmt"
)

type Config struct {
	Alphabet    []Input
	States      []StateConfig
	Transitions []TransitionConfig
}

type StateConfig struct {
	ID         string
	Name       string
	StartState bool
	FinalState bool
	EntryEvent func()
	ExitEvent  func()
}

type TransitionConfig struct {
	StartStateID string
	Input        Input
	EndStateID   string
}

type config struct {
	alphabet    Alphabet
	startState  State
	states      []State
	finalStates []State
	transitions []Transition
}

func parseConfig(fsmConfig *Config) (*config, error) {
	config := &config{
		NewAlphabet(fsmConfig.Alphabet...),
		nil,
		[]State{},
		[]State{},
		[]Transition{},
	}

	validStateIDs := config.parseStateConfigs(fsmConfig.States)

	err := config.parseTransitionConfigs(fsmConfig.Transitions, validStateIDs)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *config) parseStateConfigs(stateConfigs []StateConfig) map[string]State {
	validStateIDs := make(map[string]State, len(stateConfigs))
	for _, stateConfig := range stateConfigs {
		id, state := c.parseStateConfig(stateConfig)
		validStateIDs[id] = state
	}
	return validStateIDs
}

func (c *config) parseStateConfig(stateConfig StateConfig) (string, State) {
	state := NewState(stateConfig.ID, stateConfig.Name, stateConfig.EntryEvent, stateConfig.ExitEvent)
	if stateConfig.StartState {
		c.startState = state
	}
	if stateConfig.FinalState {
		c.finalStates = append(c.finalStates, state)
	}
	c.states = append(c.states, state)
	return stateConfig.ID, state
}

func (c *config) parseTransitionConfigs(transitionConfigs []TransitionConfig, validStateIDs map[string]State) error {
	for _, transitionConfig := range transitionConfigs {
		if err := c.parseTransitionConfig(transitionConfig, validStateIDs); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) parseTransitionConfig(transitionConfig TransitionConfig, validStateIDs map[string]State) error {
	startState, ok := validStateIDs[transitionConfig.StartStateID]
	if !ok {
		return errors.New(fmt.Sprintf("invalid transition: StartStateID <%v> does not exist", transitionConfig.StartStateID))
	}
	endState, ok := validStateIDs[transitionConfig.EndStateID]
	if !ok {
		return errors.New(fmt.Sprintf("invalid transition: EndStateID <%v> does not exist", transitionConfig.EndStateID))
	}
	if !c.alphabet.Valid(transitionConfig.Input) {
		return errors.New(fmt.Sprintf("invalid transition: Input <%v> not in Alphabet", transitionConfig.Input))
	}
	c.transitions = append(c.transitions, NewTransition(startState, transitionConfig.Input, endState))
	return nil
}

package Workflows

type State interface {
	RunEntryEvent()
	RunExitEvent()
	String() string
}

func NewState(id string, entryEvent, exitEvent func()) State {
	return &state{
		id:         id,
		entryEvent: entryEvent,
		exitEvent:  exitEvent,
	}
}

type state struct {
	id         string
	entryEvent func()
	exitEvent  func()
}

func (s *state) RunEntryEvent() {
	s.entryEvent()
}

func (s *state) RunExitEvent() {
	s.exitEvent()
}

func (s *state) String() string {
	return s.id
}

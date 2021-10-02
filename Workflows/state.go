package Workflows

type State interface {
	String() string
}

func NewState(id string) State {
	return &state{
		id: id,
	}
}

type state struct {
	id string
}

func (s *state) String() string {
	return s.id
}

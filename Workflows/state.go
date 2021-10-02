package Workflows

func NewState(id string) State {
	return &state{
		id: id,
	}
}

type State interface {
	String() string
}

type state struct {
	id string
}

func (s *state) String() string {
	return s.id
}

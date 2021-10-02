package Workflows

type Alphabet interface {
	Valid(Input) bool
}

type Input rune

func NewAlphabet(inputs ...Input) Alphabet {
	validInputs := make(map[Input]struct{}, len(inputs))
	for _, input := range inputs {
		validInputs[input] = struct{}{}
	}
	return &alphabet{
		validInputs: validInputs,
	}
}

type alphabet struct {
	validInputs map[Input]struct{}
}

func (a *alphabet) Valid(input Input) bool {
	_, ok := a.validInputs[input]
	return ok
}

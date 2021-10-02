package Workflows

func NewAlphabet(inputs ...Input) Alphabet {
	validInputs := make(map[Input]struct{}, len(inputs))
	for _, input := range inputs {
		validInputs[input] = struct{}{}
	}
	return &alphabet{
		validInputs: validInputs,
	}
}

type Alphabet interface {
	Valid(Input) bool
}

type alphabet struct {
	validInputs map[Input]struct{}
}

type Input rune

func (a *alphabet) Valid(input Input) bool {
	_, ok := a.validInputs[input]
	return ok
}

package turing

func MakeMachineFromJson(jsonMachine *JsonMachine) (*Machine, error) {
	var parse parseHelper

	// Alphabet
	if err := parse.alphabetFromString(jsonMachine.Alphabet); err != nil {
		return nil, err
	}

	// Lambda
	if err := parse.lambdaFromString(jsonMachine.Lambda); err != nil {
		return nil, err
	}

	// Start State
	if err := parse.setStartStateName(jsonMachine.Start); err != nil {
		return nil, err
	}

	// States
	if err := parse.statesFromJsonStates(jsonMachine.States); err != nil {
		return nil, err
	}

	// Final States
	if err := parse.finalStatesFromStrings(jsonMachine.Final); err != nil {
		return nil, err
	}

	return parse.makeMachine(), nil
}

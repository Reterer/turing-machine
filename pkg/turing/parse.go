package turing

import "errors"

type parseHelper struct {
	alphContainsAllRunes bool
	alphabet             map[rune]bool
	lambda               rune

	statesMap map[string]int
	states    []State

	start string
	final []string
}

// check and parse alphabet from string
func (ph *parseHelper) alphabetFromString(alphabet string) error {
	ph.alphabet = make(map[rune]bool)
	ph.alphContainsAllRunes = len(alphabet) == 0

	if ph.alphContainsAllRunes {
		return nil
	}

	for _, ch := range alphabet {
		if _, ok := ph.alphabet[ch]; ok {
			return errors.New("the alphabet contains repeated characters: '" + string(ch) + "'")
		}
		ph.alphabet[ch] = true
	}

	return nil
}

// check and parse char
func (ph *parseHelper) parseCharString(char string) (rune, error) {
	runes := []rune(char)
	if len(runes) != 1 {
		return 0, errors.New("char has several characters: '" + char + "'")
	}
	if !ph.alphContainsAllRunes {
		if _, ok := ph.alphabet[runes[0]]; !ok {
			return 0, errors.New("char '" + char + "' is not in the alphabet")
		}
	}
	return runes[0], nil
}

// check and parse lambda from string
func (ph *parseHelper) lambdaFromString(lambda string) error {
	// TODO SIMPLIFY WITH parseCharString
	lambdaRunes := []rune(lambda)
	if len(lambdaRunes) != 1 {
		return errors.New("lambda has several characters: " + lambda)
	}
	ph.lambda = lambdaRunes[0]

	if !ph.alphContainsAllRunes {
		if _, ok := ph.alphabet[ph.lambda]; !ok {
			return errors.New("lambda is not in the alphabet")
		}
	}

	return nil
}

// check state name without state map
func (ph *parseHelper) checkStateNameWithoutStateMap(name string) error {
	if len(name) == 0 {
		return errors.New("length of the state is zero")
	}
	return nil
}

// check state name
func (ph *parseHelper) checkStateName(name string) error {
	if err := ph.checkStateNameWithoutStateMap(name); err != nil {
		return err
	}
	if _, ok := ph.statesMap[name]; !ok {
		return errors.New("the state: '" + name + "' has no definition")
	}
	return nil
}

// check and set start state name
func (ph *parseHelper) setStartStateName(start string) error {
	if err := ph.checkStateNameWithoutStateMap(start); err != nil {
		return err
	}
	ph.start = start
	return nil
}

// check and parse states map
func (ph *parseHelper) statesMapFromJsonStates(states []JsonState) error {
	ph.statesMap = make(map[string]int)
	for i, jState := range states {
		if _, ok := ph.statesMap[jState.Name]; ok {
			return errors.New("state '" + jState.Name + "' has several definitions")
		}
		if err := ph.checkStateNameWithoutStateMap(jState.Name); err != nil {
			return err
		}

		ph.statesMap[jState.Name] = i
	}
	return nil
}

// check and parse transitions
func (ph *parseHelper) stateTransitionFromJsonStates(states []JsonState) error {
	ph.states = make([]State, len(ph.statesMap))
	for _, jState := range states {
		currState := &ph.states[ph.statesMap[jState.Name]]
		currState.commands = make(map[rune]Command)
		for _, jCommand := range jState.Transitions {
			var char, newChar rune
			var nextState string
			var move Move
			var err error

			// char
			char, err = ph.parseCharString(jCommand.Char)
			if err != nil {
				// TODO CUSTOM ERROR
				return err
			}
			if _, ok := currState.commands[char]; ok {
				return errors.New("multiple identical transitions. from: '" + jState.Name + "' ch: '" + jCommand.Char + "'")
			}

			// new char
			newChar, err = ph.parseCharString(jCommand.NewChar)
			if err != nil {
				return err
			}

			// move
			switch jCommand.Move {
			case "l":
				move = mLeft
			case "r":
				move = mRight
			case "n":
				move = mNone
			default:
				return errors.New("unccorect movement: '" + jCommand.Move + "' from: '" + jState.Name + "'")
			}

			// next state
			if err = ph.checkStateName(jCommand.NextState); err != nil {
				return err
			}
			nextState = jCommand.NextState

			// add transition
			currState.commands[char] = Command{
				move:      move,
				char:      newChar,
				nextState: ph.statesMap[nextState],
			}
		}
	}
	return nil
}

// check and parse states
func (ph *parseHelper) statesFromJsonStates(states []JsonState) error {
	// parse states map
	if err := ph.statesMapFromJsonStates(states); err != nil {
		return err
	}

	// swap start state idx with zero idx state
	ph.statesMap[ph.start], ph.statesMap[states[0].Name] = ph.statesMap[states[0].Name], ph.statesMap[ph.start]

	// parse states
	if err := ph.stateTransitionFromJsonStates(states); err != nil {
		return err
	}

	return nil
}

// check and parse final states
func (ph *parseHelper) finalStatesFromStrings(final []string) error {
	ph.final = make([]string, len(final))
	copy(ph.final, final)

	for _, f := range ph.final {
		if err := ph.checkStateName(f); err != nil {
			return err
		}
		ph.states[ph.statesMap[f]].isFinite = true
	}

	return nil
}

// make machine
func (ph *parseHelper) makeMachine() *Machine {
	return &Machine{
		lambda: ph.lambda,
		states: ph.states,
	}
}

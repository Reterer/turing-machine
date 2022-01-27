package turing

import (
	"encoding/json"
	"errors"
)

type JsonMachine struct {
	Alphabet string      `json:"alphabet"`
	Lambda   string      `json:"lambda"`
	Start    string      `json:"start"`
	Final    []string    `json:"final"`
	States   []JsonState `json:"states"`
}

type JsonState struct {
	Name        string        `json:"name"`
	Transitions []JsonCommand `json:"transitions"`
}

type JsonCommand struct {
	Char      string `json:"ch"`
	NewChar   string `json:"nch"`
	Move      string `json:"mv"`
	NextState string `json:"nst"`
}

func LoadJsonMachineFromFile(filename string) (*JsonMachine, error) {

	return LoadJsonMachineFromJson([]byte(`
	{
		"alphabet" : "abcd*",
		"lambda" : "*",
		"start" : "replacer",
		"final" : ["final"],
		"states" : [
			{
				"name" : "replacer",
				"transitions" : [
					{ "ch" : "a", "nch" : "*", "mv" : "r", "nst" : "replacer" },
					{ "ch" : "b", "nch" : "*", "mv" : "r", "nst" : "replacer" },
					{ "ch" : "c", "nch" : "*", "mv" : "r", "nst" : "replacer" },
					{ "ch" : "d", "nch" : "*", "mv" : "r", "nst" : "replacer" },
					{ "ch" : "*", "nch" : "*", "mv" : "n", "nst" : "final" }
				]
			},
			{
				"name" : "final",
				"transitions" : []
			}
		]
	}
	`))
}

func LoadJsonMachineFromJson(jsonData []byte) (*JsonMachine, error) {
	var jsonMachine JsonMachine
	err := json.Unmarshal(jsonData, &jsonMachine)
	return &jsonMachine, err
}

func MakeMachineFromJson(jsonMachine *JsonMachine) (*Machine, error) {
	var machine Machine
	// TODO ADD VALIDATE NAMES
	// todo add call validateJsonMachine
	// Alphabet
	alphabet := make(map[rune]bool)
	for _, ch := range jsonMachine.Alphabet {
		alphabet[ch] = true
	}

	// Lambda
	lambda := []rune(jsonMachine.Lambda)
	if _, ok := alphabet[lambda[0]]; !ok {
		return nil, errors.New("lambda has several runes")
	}
	machine.lambda = lambda[0]

	// Start State
	startState := jsonMachine.Start

	// Final States
	finalStates := make(map[string]bool)
	for _, name := range jsonMachine.Final {
		if _, ok := finalStates[name]; ok {
			return nil, errors.New("several final states has same name")
		} else {
			finalStates[name] = true
		}
	}

	// States
	// get idx for names
	states := make(map[string]int)
	statesCount := 0
	zeroStateName := ""
	for i := 0; i < len(jsonMachine.States); i++ {
		if _, ok := states[jsonMachine.States[i].Name]; ok {
			return nil, errors.New("several states has same name: " + jsonMachine.States[i].Name)
		} else {
			if statesCount == 0 {
				zeroStateName = jsonMachine.States[i].Name
			}
			states[jsonMachine.States[i].Name] = statesCount
			statesCount++
		}
	}
	if _, ok := states[startState]; !ok {
		return nil, errors.New("the initial state is not defined: " + startState)
	}
	for final := range finalStates {
		if _, ok := states[final]; !ok {
			return nil, errors.New("the final state is not defined: " + final)
		}
	}
	// Swap zero state with start state
	states[startState], states[zeroStateName] = states[zeroStateName], states[startState]

	// load states into machine
	machine.states = make([]State, statesCount)
	for _, jState := range jsonMachine.States {
		state := &machine.states[states[jState.Name]]
		state.commands = make(map[rune]Command)
		state.isFinite = false

		// this state is finite
		if _, ok := finalStates[jState.Name]; ok {
			state.isFinite = true
		}

		// add transitions
		for _, jCommand := range jState.Transitions {
			char := []rune(jCommand.Char)
			newChar := []rune(jCommand.NewChar)
			// TODO ADD VALIDATION

			if _, ok := state.commands[char[0]]; ok {
				return nil, errors.New("multiple identical transitions. from: " + jState.Name + " ch: " + jCommand.Char)
			} else {
				var move Move
				switch jCommand.Move {
				case "l":
					move = mLeft
				case "r":
					move = mRight
				case "n":
					move = mNone
				default:
					return nil, errors.New("unccorect movement: " + jCommand.Move + " from: " + jState.Name)
				}

				nextState, ok := states[jCommand.NextState]
				if !ok {
					return nil, errors.New("transition to an undefined state. from: " + jState.Name + " to: " + jCommand.NextState)
				}

				state.commands[char[0]] = Command{
					move:      move,
					char:      newChar[0],
					nextState: nextState,
				}
			}
		}
	}

	return &machine, nil
}

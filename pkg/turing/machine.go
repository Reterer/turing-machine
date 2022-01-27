package turing

import (
	"errors"
)

type Move uint8

const (
	mLeft Move = iota
	mRight
	mNone
)

type Command struct {
	move      Move
	char      rune
	nextState int
}

type State struct {
	commands map[rune]Command
	isFinite bool
}

type Machine struct {
	lambda rune

	tape   []rune
	states []State

	currState    int
	currTapeCell int
	isRunning    bool
}

func MakeMachine(lambda rune, tape []rune, states []State, startState int, startTapeCell int) *Machine {
	return &Machine{
		lambda:       lambda,
		tape:         tape,
		states:       states,
		currState:    startState,
		currTapeCell: startState,
		isRunning:    true,
	}
}

func (m *Machine) Iterate() error {
	if !m.isRunning {
		return errors.New("Machine was stopped")
	}

	currState := &m.states[m.currState]
	if currState.isFinite {
		m.isRunning = false
		return nil
	}
	if m.currTapeCell >= len(m.tape) {
		m.tape = append(m.tape, m.lambda)
	}

	currChar := &m.tape[m.currTapeCell]
	if cmd, ok := currState.commands[*currChar]; ok {
		// Обработка возможных ошибок
		if cmd.move == mLeft && m.currTapeCell == 0 {
			return errors.New("Index out of bounds < 0")
		}

		// Изменение символа на ленте
		*currChar = cmd.char

		// Движение ячейки на ленте
		switch cmd.move {
		case mLeft:
			m.currTapeCell--
		case mRight:
			m.currTapeCell++
		}

		// Смена состояния на следующее
		m.currState = cmd.nextState
	} else {
		return errors.New("Unkown char on tape")
	}

	return nil
}

func (m *Machine) IsRunning() bool {
	return m.isRunning
}

func (m *Machine) GetTape() []rune {
	return m.tape
}

func MakeTestStates() []State {
	return []State{
		State{
			commands: map[rune]Command{
				' ': Command{
					move:      mRight,
					char:      '+',
					nextState: 1},
			},
			isFinite: false,
		},
		State{
			commands: map[rune]Command{
				' ': Command{
					move:      mRight,
					char:      '+',
					nextState: 2},
			},
			isFinite: false,
		},
		State{
			isFinite: true,
		},
	}
}

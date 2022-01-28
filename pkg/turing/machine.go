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

func (m *Machine) Iterate() error {
	if !m.isRunning {
		return errors.New("machine was stopped")
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
			return errors.New("index out of bounds < 0")
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
		return errors.New("unkown char on tape")
	}

	return nil
}

func (m *Machine) IsRunning() bool {
	return m.isRunning
}

func (m *Machine) GetTapePosition() int {
	return m.currTapeCell
}

func (m *Machine) GetTape() []rune {
	return m.tape
}

func (m *Machine) SetTape(tape []rune, initPosition int) error {
	if m.isRunning {
		return errors.New("machine is running")
	}
	if len(tape) < initPosition {
		return errors.New("the length of the tape is less than the initial position")
	}

	m.currTapeCell = initPosition
	m.tape = tape
	return nil
}

func (m *Machine) TurnOn() {
	m.isRunning = true
}

/*
func (m *Machine) Run() ||
func (m *Machine) Run(breakpoints []int)
*/

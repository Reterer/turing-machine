package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Reterer/turing-machine/pkg/turing"
	arg "github.com/alexflint/go-arg"
)

type Args struct {
	InTape        string `arg:"positional,required" help:"input tape"`
	OutTape       string `arg:"positional,required" help:"output tape"`
	TuringMachine string `arg:"positional,required" help:"turing machine josn file"`
	Debug         bool   `help:"enable debug mode"`
}

type TTuringMachine int

const (
	TJson TTuringMachine = iota
)

func main() {
	var args Args
	arg.MustParse(&args)

	inputTape, err := loadTape(args.InTape)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	machine, err := loadTuringMachine(args.TuringMachine, TJson)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	outputTape := make([]rune, 0)
	if args.Debug {
		outputTape, err = debugMode(inputTape, machine)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		outputTape, err = defaultMode(inputTape, machine)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	err = saveTape(args.OutTape, outputTape)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func loadTape(input string) ([]rune, error) {
	inFile, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	inBytes, err := ioutil.ReadAll(inFile)
	if err != nil {
		return nil, err
	}
	return []rune(string(inBytes)), nil
}

func saveTape(filename string, tape []rune) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data := []byte(string(tape))
	bytesWrited, err := file.Write(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("was writed %d/%d bytes", bytesWrited, len(data))

	return nil
}

func loadTuringMachine(filename string, typeTm TTuringMachine) (*turing.Machine, error) {
	var machine *turing.Machine
	switch typeTm {
	case TJson:
		jsonMachine, err := turing.LoadJsonMachineFromFile(filename)
		if err != nil {
			return nil, err
		}
		machine, err = turing.MakeMachineFromJson(jsonMachine)
		if err != nil {
			return nil, err
		}
	default:
	}

	return machine, nil
}

func defaultMode(inputTape []rune, machine *turing.Machine) ([]rune, error) {
	if err := machine.SetTape(inputTape, 0); err != nil {
		return nil, err
	}
	machine.TurnOn()

	for machine.IsRunning() {
		err := machine.Iterate()
		if err != nil {
			return nil, err
		}
		t := machine.GetTape()
		p := machine.GetTapePosition()
		fmt.Println(string(t))
		for i := 0; i < p; i++ {
			fmt.Print(" ")
		}
		fmt.Println("^")
	}
	return machine.GetTape(), nil
}

func debugMode(inputTape []rune, machine *turing.Machine) ([]rune, error) {
	return nil, errors.New("not implemented")
}

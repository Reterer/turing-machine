package main

/*
	1. Загрузить ленту
	2. Загрузить файл
	3. Интепретировать
*/
import (
	"fmt"

	"github.com/Reterer/turing-machine/pkg/turing"
)

func main() {
	jsonMachine, err := turing.LoadJsonMachineFromFile("examples/json/replacer.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	machine, err := turing.MakeMachineFromJson(jsonMachine)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var str string
	fmt.Println("Введите ленту:")
	fmt.Scan(&str)

	machine.SetTape([]rune(str), 0)
	machine.TurnOn()
	for machine.IsRunning() {
		machine.Iterate()
		fmt.Println(string(machine.GetTape()))
		pos := machine.GetTapePosition()
		for i := 0; i < pos; i++ {
			fmt.Print(" ")
		}
		fmt.Println("^")
	}
}

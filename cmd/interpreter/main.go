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
	states := turing.MakeTestStates()

	machine := turing.MakeMachine(rune(' '), []rune(""), states, 0, 0)
	for machine.IsRunning() {
		machine.Iterate()
	}
	fmt.Println(string(machine.GetTape()))
}

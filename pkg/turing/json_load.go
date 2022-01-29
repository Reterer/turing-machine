package turing

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return LoadJsonMachineFromJson(data)
}

func LoadJsonMachineFromJson(jsonData []byte) (*JsonMachine, error) {
	var jsonMachine JsonMachine
	err := json.Unmarshal(jsonData, &jsonMachine)
	return &jsonMachine, err
}

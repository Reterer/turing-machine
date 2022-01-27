# Turing-Machine
[![Go Report Card](https://goreportcard.com/badge/github.com/Reterer/turing-machine)](https://goreportcard.com/report/github.com/Reterer/turing-machine)

Implementation of a Turing machine

## Todo list
- [x] turing machine interpreter
- [ ] loading a json description of a turing machine
- [ ] debug mode
- [ ] support for a simpler description of the Turing machine
- [ ] learn how to use go tests

## How to run interpreter
```
go run cmd\interpreter\main.go <in tape> <out tape> <turing machine> [-debug]
```
Where
* in tape - is a text file that displays the input tape
* out tape - is a text file where the resulting tape will be saved
* turing machine - is a json file describing a turing machine
* [debug] - flag for enabling debug mode
  * n - next instruction

## how to wrtie turing machine code
example of a valid file:
```json
{
    // set the alphabet of the turing machine
    // an empty string ("") means that any utf character can be used 
    "alphabet" : "abcd*",
    // space character
    // the symbol must belong to the alphabet
    "lambda" : "*",
    // initial state
    "start" : "replacer",
    // final states
    // there may be several
    "final" : ["final"],
    // description of states
    "states" : [
        // description of state
        {
            // name of the state
            "name" : "replacer",
            // transitions between states
            "transitions" : [
                /*
                    ch - read character
                    nch - the character that will be written to this cell
                    mv - tape movement (l - left, r - right, n - nothing)
                    nst - next state
                */
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
```
this program replaces characters with spaces.




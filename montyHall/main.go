package main

// Monty-Hall simulates the outcome of running the Monty-Hall problem
// It uses a Monte-Carlo simulation

import (
	"fmt"
	"math/rand/v2"
)

// Created: Mon May  3 19:31:12 2021

const (
	dfltTrialCount = 1e6
	dfltDoorCount  = 3
)

// prog contains the parameters of the program
type prog struct {
	trials     int
	doorCount  int
	changeDoor bool
	wins       int
}

// newProg returns a new Prog instance with the default values set
func newProg() *prog {
	return &prog{
		trials:    dfltTrialCount,
		doorCount: dfltDoorCount,
	}
}

func main() {
	prog := newProg()
	ps := makeParamSet(prog)
	ps.Parse()

	for range prog.trials {
		prog.runTrial()
	}

	fmt.Printf("doors: %6d  change door? %5t  win %%age: %.9f\n",
		prog.doorCount,
		prog.changeDoor,
		100*float64(prog.wins)/float64(prog.trials))
}

// runTrial runs a single trial incrementing the win count if the player
// strategy would have won
//
//nolint:gosec
func (prog *prog) runTrial() {
	prizeDoor := rand.IntN(prog.doorCount)
	chosenDoor := rand.IntN(prog.doorCount)

	if chosenDoor == prizeDoor { // if we had chosen the right door ...
		if !prog.changeDoor { // ... and we don't change door
			prog.wins++
		}
	} else if prog.changeDoor { // ... we chose the wrong door but we change
		prog.wins++
	}
}

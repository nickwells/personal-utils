package main

// Monty-Hall simulates the outcome of running the Monty-Hall problem
// It uses a Monte-Carlo simulation

import (
	"fmt"
	"math/rand"
)

// Created: Mon May  3 19:31:12 2021

const (
	dfltTrialCount = 1e6
	dfltDoorCount  = 3
)

type Prog struct {
	trials     int64
	doorCount  int64
	changeDoor bool
	wins       int
}

// NewProg returns a new Prog instance with the default values set
func NewProg() *Prog {
	return &Prog{
		trials:    dfltTrialCount,
		doorCount: dfltDoorCount,
	}
}

func main() {
	prog := NewProg()
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
func (prog *Prog) runTrial() {
	prizeDoor := rand.Intn(int(prog.doorCount))
	chosenDoor := rand.Intn(int(prog.doorCount))

	if chosenDoor == prizeDoor { // if we had chosen the right door ...
		if !prog.changeDoor { // ... and we don't change door
			prog.wins++
		}
	} else if prog.changeDoor { // ... we chose the wrong door but we change
		prog.wins++
	}
}

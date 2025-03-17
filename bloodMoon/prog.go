package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/nickwells/tempus.mod/tempus"
	"github.com/nickwells/verbose.mod/verbose"
)

// Prog holds program parameters and status
type Prog struct {
	exitStatus int
	stack      *verbose.Stack
	year       int64
}

// NewProg returns a new Prog instance with the default values set
func NewProg() *Prog {
	return &Prog{
		stack: &verbose.Stack{},
	}
}

// SetExitStatus sets the exit status to the new value. It will not do this
// if the exit status has already been set to a non-zero value.
func (prog *Prog) SetExitStatus(es int) {
	if prog.exitStatus == 0 {
		prog.exitStatus = es
	}
}

// ForceExitStatus sets the exit status to the new value. It will do this
// regardless of the existing exit status value.
func (prog *Prog) ForceExitStatus(es int) {
	prog.exitStatus = es
}

// Run is the starting point for the program, it should be called from main()
// after the command-line parameters have been parsed. Use the setExitStatus
// method to record the exit status and then main can exit with that status.
//
//nolint:gosec
func (prog *Prog) Run() {
	days := []int64{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	month := rand.Int64N(int64(len(days)))
	dom := rand.Int64N(days[month])
	h := rand.Int64N(tempus.HoursPerDay)
	m := rand.Int64N(tempus.MinutesPerHour)
	s := rand.Int64N(tempus.SecondsPerMinute)

	fmt.Printf("%4d%02d%02d%02d%02d%02d\n", prog.year, month+1, dom+1, h, m, s)
}

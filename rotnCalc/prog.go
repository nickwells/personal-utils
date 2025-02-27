package main

import (
	"fmt"
	"math"

	"github.com/nickwells/verbose.mod/verbose"
)

const defaultOutputPrecision = 2

// Prog holds program parameters and status
type Prog struct {
	exitStatus int
	stack      *verbose.Stack

	acc    float64
	accSet int

	rpm    float64
	rpmSet int

	radius    float64
	radiusSet int

	precision int
}

// NewProg returns a new Prog instance with the default values set
func NewProg() *Prog {
	return &Prog{
		stack:     &verbose.Stack{},
		precision: defaultOutputPrecision,
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
func (prog *Prog) Run() {
	const (
		twoPi   = math.Pi * 2
		rps2rpm = 60 / twoPi
		rpm2rps = twoPi / 60
	)

	if prog.rpmSet == 1 && prog.accSet == 1 {
		omega := prog.rpm * rpm2rps
		fmt.Printf("Radius: %.*f m\n",
			prog.precision, prog.acc/(omega*omega))

		return
	}

	if prog.rpmSet == 1 && prog.radiusSet == 1 {
		omega := prog.rpm * rpm2rps
		fmt.Printf("Acceleration: %.*f m/s^2\n",
			prog.precision, omega*omega*prog.radius)

		return
	}

	if prog.radiusSet == 1 && prog.accSet == 1 {
		rpm := math.Sqrt(prog.acc/prog.radius) * rps2rpm
		fmt.Printf("RPM: %.*f revolutions per minute\n",
			prog.precision, rpm)

		return
	}
}

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

	acc    float64 // acceleration
	accSet int

	rpm    float64 // revolutions per minute
	rpmSet int
	omega  float64 // radians per second
	omega2 float64 // omega squared

	radius    float64
	radiusSet int

	precision int

	showCircSpeed bool
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
	)

	switch {
	case prog.rpmSet == 1 && prog.accSet == 1:
		prog.radius = prog.acc / prog.omega2
		fmt.Printf("Radius: %.*f m\n", prog.precision, prog.radius)
	case prog.rpmSet == 1 && prog.radiusSet == 1:
		prog.acc = prog.omega2 * prog.radius
		fmt.Printf("Acceleration: %.*f m/s^2\n", prog.precision, prog.acc)
	case prog.radiusSet == 1 && prog.accSet == 1:
		prog.rpm = math.Sqrt(prog.acc/prog.radius) * rps2rpm
		prog.setOmega()
		fmt.Printf("RPM: %.*f revolutions per minute\n",
			prog.precision, prog.rpm)
	}

	if prog.showCircSpeed {
		fmt.Printf("circumferential speed: %.*f m/s\n",
			prog.precision, prog.omega*prog.radius)
	}
}

// setOmega will set the omega and omega2 values for prog
func (prog *Prog) setOmega() {
	const rpm2RadPerSec = math.Pi * 2 / 60

	prog.omega = prog.rpm * rpm2RadPerSec
	prog.omega2 = prog.omega * prog.omega
}

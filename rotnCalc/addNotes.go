package main

import (
	"github.com/nickwells/param.mod/v6/param"
)

const (
	noteBaseName = "rotnCalc - "

	noteNameAlgo = noteBaseName + "Algorithm"
)

// addNotes adds the notes for this program.
func addNotes(_ *Prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.AddNote(noteNameAlgo, "The formula used to calculate the values is:"+
			"\n\n"+
			"acceleration = radius * radians per second squared"+
			"\n\n"+
			"Note that you can convert from rpm to radians per second by"+
			" multiplying by 2*Pi and dividing by 60.")

		return nil
	}
}

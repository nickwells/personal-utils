package main

import (
	"github.com/nickwells/param.mod/v6/param"
)

const (
	noteBaseName = "penney - "

	noteNameRuns = noteBaseName + "runs"
)

// addNotes adds the notes for this program.
func addNotes(_ *Prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.AddNote(noteNameRuns,
			"A run is a sequence of wins by the same player. It is worth"+
				" taking note of the maximum runs for each choice as this"+
				" gives you an idea of the maximum loss you can expect"+
				" before the errect of the long-term odds takes over.")

		return nil
	}
}

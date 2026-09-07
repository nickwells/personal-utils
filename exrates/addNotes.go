package main

import (
	"github.com/nickwells/param.mod/v7/param"
)

const (
	noteBaseName = "exrates - "

	noteNameCaching = noteBaseName + "caching"
)

// addNotes adds the notes for this program.
func addNotes(_ *prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.AddNote(noteNameCaching,
			"the exchange rates are cached the first time they are used"+
				" and only refreshed when new rates are expected in the website.")

		return nil
	}
}

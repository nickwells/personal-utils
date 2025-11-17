package main

import (
	"github.com/nickwells/param.mod/v6/param"
)

const (
	noteBaseName = "cgt - "

	noteNameTrades = noteBaseName + "trades"
)

// addNotes adds the notes for this program.
func addNotes(_ *prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.AddNote(noteNameTrades, "This program takes a file of trade"+
			" data and calculates the capital gains on the trtades. The"+
			" trade file can be populated from the Fundsmith spreadsheet"+
			" by selecting the relevant columns and pasting them into a"+
			" text file. Note that the output is formatted with commas"+
			" between thousands so you will need to edit the file to"+
			" remove these. Also the Fundsmith spreadsheet has the number"+
			" of shares sold as a negative value (it is reducing the"+
			" total holding) but for tax return purposes the value should"+
			" be positive. The Total Capital Gain figure is the one you"+
			" need to provide on your tax return.")

		return nil
	}
}

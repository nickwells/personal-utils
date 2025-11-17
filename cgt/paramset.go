package main

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/verbose.mod/verbose"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// makeParamSet generates the param set ready for parsing
func makeParamSet(prog *prog) *param.PSet {
	return paramset.NewOrPanic(
		verbose.AddParams,
		verbose.AddTimingParams(prog.stack),
		versionparams.AddParams,

		addParams(prog),
		addNotes(prog),

		param.SetProgramDescription(
			"This program takes a file of trade data (number of units and"+
				" price), representing sales, an aggregate purchase price"+
				" and a tax rate and will report the value of each trade,"+
				" the associated capital gain and the associated tax. It"+
				" will then report the totals of these amounts."),
	)
}

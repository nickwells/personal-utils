package main

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/verbose.mod/verbose"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// makeParamSet generates the param set ready for parsing
func makeParamSet(prog *Prog) *param.PSet {
	return paramset.NewOrPanic(
		verbose.AddParams,
		verbose.AddTimingParams(prog.stack),
		versionparams.AddParams,

		addParams(prog),

		param.SetProgramDescription(
			"This program will generate a randon=m birth stamp for the"+
				" given year. the birth stamp is a concept in the Blood"+
				" Moon universe where each citizen of the moon is"+
				" allocated a fourteen digit number recording the date"+
				" and time when they were first registered."),
	)
}

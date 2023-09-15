package main

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// makeParamSet generates the param set ready for parsing
func makeParamSet(prog *Prog) *param.PSet {
	return paramset.NewOrPanic(
		addParams(prog),
		versionparams.AddParams,
		param.SetProgramDescription(
			"this will run a monte-carlo simulation"+
				" demonstrating the Monte Hall problem"),
	)
}

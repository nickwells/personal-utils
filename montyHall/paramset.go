package main

import (
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramset"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// makeParamSet generates the param set ready for parsing
func makeParamSet(prog *prog) *param.PSet {
	return paramset.New(
		addParams(prog),
		versionparams.AddParams,
		param.SetProgramDescription(
			"this will run a monte-carlo simulation"+
				" demonstrating the Monte Hall problem"),
	)
}

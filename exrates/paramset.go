package main

import (
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramset"
	"github.com/nickwells/verbose.mod/verbose"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// makeParamSet generates the param set ready for parsing
func makeParamSet(prog *prog) *param.PSet {
	return paramset.New(
		verbose.AddParams,
		verbose.AddTimingParams(prog.stack),
		versionparams.AddParams,

		addParams(prog),
		addNotes(prog),

		param.SetProgramDescription(
			"This program will use the latest"+
				" foreign exchange rate data from the"+
				" UK inland revenue site and convert the"+
				" given amount between the"+
				" two supplied currencies."),
	)
}

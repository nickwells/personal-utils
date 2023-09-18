package main

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// makeParamSet generates the param set ready for parsing
func makeParamSet(prog *Prog) *param.PSet {
	return paramset.NewOrPanic(
		versionparams.AddParams,

		addParams(prog),

		param.SetProgramDescription(
			"This will read Go files given after '--' and will"+
				" look for lines with invalid Deprecation comments."+
				" It will report any such line it finds."),
	)
}

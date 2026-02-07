package main

import (
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/paramset"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// makeParamSet generates the param set ready for parsing
func makeParamSet() *param.PSet {
	return paramset.New(
		versionparams.AddParams,

		param.SetTrailingParamsName("Go-file"),

		param.SetProgramDescription(
			"This will read Go files and will"+
				" look for lines with invalid Deprecation comments."+
				" It will report any such line it finds."),
	)
}

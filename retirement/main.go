package main

// retirement

import (
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/personal-utils/retirement/model"
	"github.com/nickwells/versionparams.mod/versionparams"
)

// main
func main() {
	m := model.New()
	ps := paramset.NewOrDie(
		versionparams.AddParams,
		model.AddParams(m),
		SetConfigFile,
		param.SetProgramDescription(
			"this will simulate various scenarios for retirement"+
				" allowing you to explore the effect of changes"+
				" in your portfolio, inflation etc"))
	ps.Parse()

	m.Report(m.CalcValues())

	m.ReportModelMetrics()
}

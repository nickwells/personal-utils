package main

import (
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

func addParams(prog *prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.Add("trials",
			psetter.Int[int]{Value: &prog.trials},
			"the number of trials to perform",
		)

		ps.Add("doors",
			psetter.Int[int]{Value: &prog.doorCount},
			"the number of doors",
		)

		ps.Add("change",
			psetter.Bool{Value: &prog.changeDoor},
			"change door",
		)

		return nil
	}
}

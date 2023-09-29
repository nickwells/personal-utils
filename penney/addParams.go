package main

import (
	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

// addParams adds the parameters for this program
func addParams(prog *Prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.Add("trials",
			psetter.Int[int]{
				Value: &prog.trials,
				Checks: []check.ValCk[int]{
					check.ValGT[int](200),
				},
			},
			"how many coin tosses should be performed")
		return nil
	}
}

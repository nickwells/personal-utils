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
		ps.Add("coins",
			psetter.Int[int]{
				Value: &prog.coinCount,
				Checks: []check.ValCk[int]{
					check.ValGT[int](2),
				},
			},
			"how many coins should be chosen")
		ps.Add("try-all",
			psetter.Bool{
				Value: &prog.tryAll,
			},
			"player 2 tries all the possible alternatives")
		ps.Add("show-win-count",
			psetter.Bool{
				Value: &prog.showWinCount,
			},
			"show the number of wins")
		ps.Add("show-rough-results",
			psetter.Bool{
				Value: &prog.showRoughly,
			},
			"show the proportion of wins as the nearest 'neat' figure"+
				" within 1% of the actual figure (that's within 1/100 of"+
				" the percentage value)",
			param.AltNames("show-roughly"))
		return nil
	}
}

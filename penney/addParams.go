package main

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

// addParams adds the parameters for this program
func addParams(prog *Prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		const (
			minTrials = 200
			minCoins  = 3
		)

		ps.Add("trials",
			psetter.Int[int]{
				Value: &prog.trials,
				Checks: []check.ValCk[int]{
					check.ValGE(minTrials),
				},
			},
			"how many coin tosses should be performed")

		ps.Add("coins",
			psetter.Int[int]{
				Value: &prog.coinCount,
				Checks: []check.ValCk[int]{
					check.ValGE(minCoins),
				},
			},
			"how many coins should be chosen")

		copyCountParam := ps.Add("copy-count",
			psetter.Int[int]{
				Value: &prog.copyCount,
				Checks: []check.ValCk[int]{
					check.ValGE(0),
				},
			},
			"how many coins should be copied when generating player"+
				" number two's coin sequence")

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

		ps.Add("show-run-info",
			psetter.Bool{
				Value: &prog.showRunInfo,
			},
			"show the run information where a run is a"+
				" sequence of wins by the same player")

		ps.Add("show-rough-results",
			psetter.Bool{
				Value: &prog.showRoughly,
			},
			"show the proportion of wins as the nearest 'neat' figure"+
				" within 1% of the actual figure (that's within 1/100 of"+
				" the percentage value)",
			param.AltNames("show-roughly"))

		ps.AddFinalCheck(func() error {
			if prog.copyCount >= prog.coinCount {
				return fmt.Errorf("The copy count (%d) must be less than"+
					" the number of coins (%d)",
					prog.copyCount, prog.coinCount)
			}
			return nil
		})

		ps.AddFinalCheck(func() error {
			if !copyCountParam.HasBeenSet() {
				prog.copyCount = prog.coinCount - 1
			}
			return nil
		})

		return nil
	}
}

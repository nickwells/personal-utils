package main

import (
	"errors"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/param.mod/v6/paction"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

const (
	paramNameAcc        = "acc"
	paramNameRpm        = "rpm"
	paramNameRadius     = "radius"
	paramNamePrecisionn = "precision"

	numParamsReqd = 2
)

// addParams adds the parameters for this program
func addParams(prog *Prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.Add(paramNameAcc,
			psetter.Float[float64]{
				Value: &prog.acc,
				Checks: []check.ValCk[float64]{
					check.ValGT(0.0),
				},
			},
			"the acceleration",
			param.PostAction(paction.SetVal(&prog.accSet, 1)),
			param.AltNames("a"),
		)

		ps.Add(paramNameRpm,
			psetter.Float[float64]{
				Value: &prog.rpm,
				Checks: []check.ValCk[float64]{
					check.ValGT(0.0),
				},
			},
			"the revolutions per minute",
			param.PostAction(paction.SetVal(&prog.rpmSet, 1)),
		)

		ps.Add(paramNameRadius,
			psetter.Float[float64]{
				Value: &prog.radius,
				Checks: []check.ValCk[float64]{
					check.ValGT(0.0),
				},
			},
			"the radius",
			param.PostAction(paction.SetVal(&prog.radiusSet, 1)),
			param.AltNames("r"),
		)

		ps.AddFinalCheck(func() error {
			if (prog.accSet + prog.rpmSet + prog.radiusSet) != numParamsReqd {
				return errors.New("Two and only two of " +
					english.Join(
						[]string{paramNameAcc, paramNameRpm, paramNameRadius},
						", ", " and ") +
					" must be set")
			}

			return nil
		})

		ps.Add(paramNamePrecisionn,
			psetter.Int[int]{
				Value: &prog.precision,
				Checks: []check.ValCk[int]{
					check.ValGE(0),
				},
			},
			"the precision with which to print the result",
			param.AltNames("p", "prec"),
		)

		// ps.AddGroup("group-name", "description")
		// ps.AddExample("example", "description")
		// ps.AddNote("headline", "text")
		// ps.AddReference("name", "description")
		return nil
	}
}

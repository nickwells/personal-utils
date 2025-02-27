package main

import (
	"fmt"

	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

const (
	paramNameYear = "year"
)

const (
	minYear = 2087
	maxYear = 2151
)

// addParams adds the parameters for this program
func addParams(prog *Prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.Add(paramNameYear, psetter.Int[int64]{Value: &prog.year},
			"The year is the year of birth or arrival at the moon",
			param.Attrs(param.MustBeSet))

		ps.AddFinalCheck(func() error {
			if prog.year < minYear {
				return fmt.Errorf("The earliest year is %d", minYear)
			}

			if prog.year > maxYear {
				return fmt.Errorf("The latest year is %d", maxYear)
			}

			return nil
		})
		// ps.AddGroup("group-name", "description")
		// ps.AddExample("example", "description")
		// ps.AddNote("headline", "text")
		// ps.AddReference("name", "description")
		return nil
	}
}

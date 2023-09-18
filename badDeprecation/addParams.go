package main

import (
	"github.com/nickwells/param.mod/v6/param"
)

// addParams will add parameters to the passed ParamSet
func addParams(prog *Prog) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		err := ps.SetNamedRemHandler(prog, "Go-file")
		if err != nil {
			return err
		}

		return nil
	}
}

package main

import (
	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

const (
	paramNameTradeFile     = "trades-file"
	paramNamePurchasePrice = "purchase-price"
	paramNameCGTRate       = "cgt-rate"
)

// addParams adds the parameters for this program
func addParams(prog *prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.Add(paramNameTradeFile,
			psetter.Pathname{
				Value:       &prog.filename,
				Expectation: filecheck.FileExists(),
			},
			"the name of the file containing the trades",
			param.Attrs(param.MustBeSet),
			param.AltNames("trade-file", "file"),
		)

		ps.Add(paramNamePurchasePrice,
			psetter.Float[float64]{
				Value: &prog.purchasePx,
				Checks: []check.ValCk[float64]{
					check.ValGT[float64](0),
				},
			},
			"the weighted average of the purchase price",
			param.Attrs(param.MustBeSet),
			param.AltNames("px"),
		)

		const (
			minTaxRate = 0
			maxTaxRate = 100
		)
		ps.Add(paramNameCGTRate,
			psetter.Float[float64]{
				Value: &prog.cgtRate,
				Checks: []check.ValCk[float64]{
					check.ValBetween[float64](minTaxRate, maxTaxRate),
				},
			},
			"the current capital gains tax rate (as a percentage)",
			param.Attrs(param.MustBeSet),
		)

		return nil
	}
}

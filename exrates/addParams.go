package main

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v7/paction"
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/psetter"
)

const (
	paramNameCacheFile = "cache-file"
	paramNameCacheDir  = "cache-dir"
	paramNameFrom      = "from-currency"
	paramNameTo        = "to-currency"
	paramNameAmount    = "amount"
	paramNameAsOf      = "as-of"
)

// addParams adds the parameters for this program
func addParams(prog *prog) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		ps.Add(paramNameCacheFile,
			psetter.Pathname{
				Value:         &prog.cacheFile,
				Expectation:   filecheck.FileExists(),
				ForceAbsolute: true,
			},
			"use this given filename (which must exist)"+
				" to supply the currency rates",
			param.PostAction(
				paction.SetVal(&prog.useGivenFile, true)))
		ps.Add(paramNameCacheDir,
			psetter.Pathname{
				Value:       &prog.cacheDir,
				Expectation: filecheck.DirExists(),
			},
			"use this directory in place of the default cache directory")
		fromParam := ps.Add(paramNameFrom,
			psetter.String[CurrencyCode]{
				Value: &prog.from,
				Checks: []check.ValCk[CurrencyCode]{
					check.ValCk[CurrencyCode](CheckCurrencyCode),
				},
			},
			"the currency to convert from",
			param.AltNames("from", "from-ccy"),
		)
		toParam := ps.Add(paramNameTo,
			psetter.String[CurrencyCode]{
				Value: &prog.to,
				Checks: []check.ValCk[CurrencyCode]{
					check.ValCk[CurrencyCode](CheckCurrencyCode),
				},
			},
			"the currency to convert to",
			param.AltNames("to", "to-ccy"),
		)
		ps.Add(paramNameAmount,
			psetter.Float[float64]{
				Value: &prog.amount,
				Checks: []check.Float64{
					check.ValGT[float64](0.0),
				},
			},
			"set the amount to be converted",
			param.AltNames("value", "val", "v", "amt"))

		ps.AddFinalCheck(func() error {
			if !fromParam.HasBeenSet() &&
				!toParam.HasBeenSet() {
				return fmt.Errorf(
					"at least one of the parameters: %q and %q must be set",
					paramNameFrom,
					paramNameTo)
			}

			return nil
		})

		return nil
	}
}

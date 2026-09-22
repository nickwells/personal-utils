package main

import (
	"fmt"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/groupsetter.mod/groupsetter"
	"github.com/nickwells/locale.mod/locale"
	"github.com/nickwells/param.mod/v7/paction"
	"github.com/nickwells/param.mod/v7/param"
	"github.com/nickwells/param.mod/v7/psetter"
	"github.com/nickwells/timesetter.mod/timesetter"
)

const (
	paramNameCacheFile = "cache-file"
	paramNameCacheDir  = "cache-dir"
	paramNameFrom      = "from-currency"
	paramNameTo        = "to-currency"
	paramNameAmount    = "amount"
	paramNameAsOf      = "as-of"
)

// makeAsOfSetter creates a groupsetter.Single for the asOf date.
func makeAsOfSetter(prog *prog) *groupsetter.Single[asOf] {
	now := time.Now()

	const (
		paramNameMonth = "month"
		paramNameYear  = "year"
	)

	const minYear = 2021

	gSetter := groupsetter.NewSingle(&prog.asOf)

	gSetter.Sep = "/"

	gSetter.AddByPosParam(paramNameMonth,
		timesetter.MonthSetter{
			Value: &gSetter.InterimVal.m,
			Time:  &locale.TimeEnglish,
		},
		"set the value of the month.",
	)
	gSetter.AddByPosParam(paramNameYear,
		psetter.Int[int]{
			Value: &gSetter.InterimVal.y,
			Checks: []check.ValCk[int]{
				check.ValGE(minYear),
				check.ValLE(now.Year()),
			},
		},
		"set the value of the year.",
	)

	return gSetter
}

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
			param.AltNames("value", "val", "v", "amt"),
		)

		ps.Add(paramNameAsOf,
			makeAsOfSetter(prog),
			"the month and year for which to give the exchange rates",
		)

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

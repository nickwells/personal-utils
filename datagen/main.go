//nolint:mnd
package main

// datagen

import (
	"fmt"
	"strings"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/datagen.mod/datagen"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

// Created: Sat Aug 20 12:19:22 2022

var startTime = time.Date(
	2022, time.August, 25,
	13, 0, 0, 0,
	time.FixedZone("UTC", 0))

// prog holds the parameters and current status of the program
type prog struct {
	// the number of records to generate
	count int64
}

// newProg creates and returns a new instance of a prog structure with the
// count set to a default of 1
func newProg() *prog {
	return &prog{
		count: 1,
	}
}

func main() {
	prog := newProg()
	ps := makeParamSet(prog)

	ps.Parse()

	sortCodeGen := datagen.NewGen(datagen.GenSetValue("'12-34-56"))
	acctNumGen := datagen.NewGen(datagen.GenSetValue("09876543"))

	transTypeGen := datagen.NewWStringGen(datagen.Random,
		datagen.WeightedString{Str: "BGC", Weight: 1},
		datagen.WeightedString{Str: "SO", Weight: 8},
		datagen.WeightedString{Str: "DD", Weight: 20},
		datagen.WeightedString{Str: "DEB", Weight: 97},
	)

	debitDescGen := datagen.NewWStringGen(datagen.Random,
		datagen.WeightedString{Str: "LOCAL CAFE", Weight: 3},
		datagen.WeightedString{Str: "LOCAL SUPERMARKET", Weight: 5},
		datagen.WeightedString{Str: "POSH SUPERMARKET", Weight: 1},
		datagen.WeightedString{Str: "CHEAP SUPERMARKET", Weight: 3},
	)
	creditDescGen := datagen.NewGen(datagen.GenSetValue("EMPLOYER"))

	transIsBGC := datagen.NewValCk(check.ValEQ("BGC"), transTypeGen)
	transIsNotBGC := datagen.NewValCk(
		check.Not(check.ValEQ("BGC"), ""),
		transTypeGen)

	gb := datagen.Countries["GB"]
	gbp := gb.Ccy()
	nf := datagen.NewNumFmt(
		datagen.NumFmtSetSepCount(),
		datagen.NumFmtSetZeroVal(""),
	)
	toStr := gbp.MoneyMkStrFunc(nf)
	moneySM := datagen.NewMoneyStringMaker(func(m datagen.Money) string {
		return toStr(m.Amt)
	})
	setMoneySM := datagen.GenSetStringMaker(moneySM)

	zeroMoney := datagen.Money{Ccy: gbp, Amt: 0}
	zeroGen := datagen.NewGen(setMoneySM, datagen.GenSetValue(zeroMoney))
	salary := datagen.Money{Ccy: gbp, Amt: 275635}
	salaryGen := datagen.NewGen(setMoneySM, datagen.GenSetValue(salary))
	randMoneyGen := datagen.NewGen(setMoneySM,
		datagen.GenSetValue(datagen.Money{Ccy: gbp, Amt: 2340}),
		datagen.GenSetValSetter(
			datagen.NewMoneyValSetter(
				datagen.NewNormValSetter[int64](150, 50000, 6000, 3000))))

	debitGen := datagen.NewSwitchGen(zeroGen,
		datagen.NewCase(transIsNotBGC, randMoneyGen))
	creditGen := datagen.NewSwitchGen(zeroGen,
		datagen.NewCase(transIsBGC, salaryGen))
	aggregator := func(v *datagen.Money,
		vals ...datagen.TypedGenerator[datagen.Money],
	) {
		v.Amt = v.Amt - vals[0].Value().Amt + vals[1].Value().Amt
	}
	balanceGen := datagen.NewGen(setMoneySM,
		datagen.GenSetValue(datagen.Money{Ccy: gbp, Amt: 132045}),
		datagen.GenSetValSetter(
			datagen.NewComputedValSetter(
				aggregator, debitGen, creditGen)))

	r := datagen.NewRecord("",
		datagen.NewField("Transaction Date",
			datagen.NewTimeGen(
				datagen.TimeGenSetLayout("02/01/2006"),
				datagen.TimeGenSetInitialTime(startTime))),
		datagen.NewField("Transaction Type", transTypeGen),
		datagen.NewField("Sort Code", sortCodeGen),
		datagen.NewField("Account Number", acctNumGen),
		datagen.NewField("Transaction Description",
			datagen.NewSwitchGen(
				debitDescGen,
				datagen.NewCase(transIsBGC, creditDescGen))),
		datagen.NewField("Debit Amount", debitGen),
		datagen.NewField("Credit Amount", creditGen),
		datagen.NewField("Balance", balanceGen),
	)

	fmt.Println(strings.Join(r.GenerateTitles(), ","))

	for range prog.count {
		fmt.Println(strings.Join(r.Generate(), ","))
		r.Next()
	}
}

// addParams will add parameters to the passed ParamSet
func addParams(prog *prog) func(ps *param.PSet) error {
	return func(ps *param.PSet) error {
		ps.Add("count",
			psetter.Int[int64]{
				ValueReqMandatory: psetter.ValueReqMandatory{},
				Value:             &prog.count,
				Checks:            []check.Int64{check.ValGT[int64](0)},
			},
			"how many records to generate")

		return nil
	}
}

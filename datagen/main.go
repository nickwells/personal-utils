//nolint:mnd
package main

// datagen

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/datagen.mod/datagen"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/psetter"
)

// Created: Sat Aug 20 12:19:22 2022

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

// intermittentIterator replaces the Next function on a Generator with one
// that is only called intermittently.
type intermittentIterator[T any] struct {
	next int
	min  int
	max  int
}

// Next ...
func (ti *intermittentIterator[T]) Next(tg datagen.TypedGenerator[T]) {
	if ti.next == 0 {
		ti.next = ti.min + rand.IntN(ti.max-ti.min) //nolint:gosec
	}

	ti.next--
	if ti.next == 0 {
		tg.Next()
	}
}

// makeTimeGen creates and returns a generator for the date
func makeTimeGen() datagen.TypedGenerator[time.Time] {
	interIter := &intermittentIterator[time.Time]{
		min: 3,
		max: 8,
	}
	startTime := time.Date(2025, time.June, 1, 12, 0, 0, 0,
		time.FixedZone("UTC", 0))

	timeGen, err := datagen.NewGenWrapper[time.Time](datagen.NewTimeGen(
		datagen.TimeGenSetLayout("Mon 02/01/2006"),
		datagen.TimeGenSetInitialTime(startTime),
		datagen.TimeGenSetIntervalF(
			datagen.TimeGenConstIntervalF(time.Hour*24)),
	),
		datagen.GenWrapperSetIterator[time.Time](interIter),
	)
	if err != nil {
		fmt.Println("Cannot build the GenWrapper", err)
		os.Exit(1)
	}

	return timeGen
}

// sinceLast satisfies the datagen.Passer interface. Its Passes method checks
// that the supplied TypedVal[time.Time] is more than the given time since
// the last time it passed.
type sinceLast struct {
	val      datagen.TypedVal[time.Time]
	last     time.Time
	interval time.Duration
}

// Passes returns true if the time falls on the first Monday of the month and
// the embedded TypedVal[time.Time] is more that 'interval' since 'last'. It
// assumes that the time values returned by val.Value() are increasing.
func (sl *sinceLast) Passes() bool {
	if sl.val.Value().Sub(sl.last) > sl.interval {
		sl.last = sl.val.Value()

		return true
	}

	return false
}

// makeSalaryCheck returns a Passer that will pass if the record should be a
// Salary record. It is based on the date - if the date is the first Monday
// of the month and this is the first time we've seen it in the last 2 days.
func makeSalaryCheck(timeVal datagen.TypedVal[time.Time]) datagen.Passer {
	return datagen.MakeAndPasser(
		datagen.MakePasser(
			check.TimeIsNthWeekdayOfMonth(1, time.Monday),
			timeVal),
		&sinceLast{
			interval: time.Hour * 24 * 2, // two days
			val:      timeVal,
		},
	)
}

// salaryGenImpl has a SetVal method thereby satisfying the ValSetter
// interface
type salaryGenImpl struct {
	p datagen.Passer
}

// SetVal ...
func (sgi salaryGenImpl) SetVal(v *bool) {
	*v = sgi.p.Passes()
}

func main() {
	prog := newProg()
	ps := makeParamSet(prog)

	ps.Parse()

	timeGen := makeTimeGen()
	isSalary := datagen.NewGen(
		datagen.GenSetValue(false),
		datagen.GenSetValSetter[bool](
			salaryGenImpl{
				p: makeSalaryCheck(timeGen),
			}))
	salaryCheck := datagen.MakeBoolPasser(isSalary)
	sortCodeGen := datagen.NewGen(datagen.GenSetValue("12-34-56"))
	acctNumGen := datagen.NewGen(datagen.GenSetValue("09876543"))

	nonSalaryTransTypeGen := datagen.NewWStringGen(datagen.Random,
		datagen.WeightedString{Str: "DEB", Weight: 10},
		datagen.WeightedString{Str: "DD", Weight: 1},
	)
	transTypeGen := datagen.NewSwitchGen(
		nonSalaryTransTypeGen, // default
		datagen.NewCase(salaryCheck, datagen.NewGen(datagen.GenSetValue("BGC"))))
	directDebitCheck := datagen.MakePasser(check.ValEQ("DD"), transTypeGen)
	debitDescGen := datagen.NewWStringGen(datagen.Random,
		datagen.WeightedString{Str: "LOCAL CAFE", Weight: 3},
		datagen.WeightedString{Str: "LOCAL SUPERMARKET", Weight: 5},
		datagen.WeightedString{Str: "POSH SUPERMARKET", Weight: 1},
		datagen.WeightedString{Str: "CHEAP SUPERMARKET", Weight: 3},
		datagen.WeightedString{Str: "PUB", Weight: 4},
		datagen.WeightedString{Str: "GASTROPUB", Weight: 2},
		datagen.WeightedString{Str: "FANCY RESTAURANT", Weight: 1},
		datagen.WeightedString{Str: "BOOKSHOP", Weight: 1},
	)
	directDebitDescGen := datagen.NewWStringGen(datagen.Random,
		datagen.WeightedString{Str: "COUNCIL", Weight: 3},
		datagen.WeightedString{Str: "PHONE COMPANY", Weight: 5},
		datagen.WeightedString{Str: "GAS AND ELEC", Weight: 1},
		datagen.WeightedString{Str: "WATER COMPANY", Weight: 3},
	)
	transDescGen := datagen.NewSwitchGen(
		debitDescGen, // default
		datagen.NewCase(directDebitCheck, directDebitDescGen),
		datagen.NewCase(salaryCheck, datagen.NewGen(datagen.GenSetValue("EMPLOYER"))),
	)

	transIsNotBGC := datagen.MakePasser(
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
		datagen.NewCase(salaryCheck, salaryGen))
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

	r, err := datagen.NewRecord("",
		datagen.NewField("Transaction Date", timeGen),
		datagen.NewHiddenField(isSalary),
		datagen.NewField("Transaction Type", transTypeGen),
		datagen.NewField("Sort Code", sortCodeGen),
		datagen.NewField("Account Number", acctNumGen),
		datagen.NewField("Transaction Description", transDescGen),
		datagen.NewField("Debit Amount", debitGen),
		datagen.NewField("Credit Amount", creditGen),
		datagen.NewField("Balance", balanceGen),
	)
	if err != nil {
		fmt.Println("cannot build the report:", err)
		os.Exit(1)
	}

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

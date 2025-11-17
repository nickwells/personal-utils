package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/nickwells/col.mod/v5/col"
	"github.com/nickwells/col.mod/v5/colfmt"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/mathutil.mod/v2/mathutil"
	"github.com/nickwells/verbose.mod/verbose"
)

// prog holds program parameters and status
type prog struct {
	exitStatus int
	stack      *verbose.Stack
	// parameters
	filename   string
	purchasePx float64
	cgtRate    float64
}

// newProg returns a new Prog instance with the default values set
func newProg() *prog {
	return &prog{
		stack: &verbose.Stack{},
	}
}

// setExitStatus sets the exit status to the new value. It will not do this
// if the exit status has already been set to a non-zero value.
func (prog *prog) setExitStatus(es int) {
	if prog.exitStatus == 0 {
		prog.exitStatus = es
	}
}

// run is the starting point for the program, it should be called from main()
// after the command-line parameters have been parsed. Use the setExitStatus
// method to record the exit status and then main can exit with that status.
func (prog *prog) run() {
	trades, err := prog.getTrades()
	if err != nil {
		prog.setExitStatus(1)
		fmt.Println("Couldn't get the trades: ", err)

		return
	}

	const (
		colWidth = 10
		colPrec  = 2
	)

	rpt := col.StdRpt(
		col.New(&colfmt.Float{W: colWidth, Prec: colPrec}, "Trade", "Amount"),
		col.New(&colfmt.Float{W: colWidth, Prec: colPrec}, "Trade", "Price"),
		col.New(&colfmt.Float{W: colWidth, Prec: colPrec}, "Total", "Value"),
		col.New(&colfmt.Float{W: colWidth, Prec: colPrec}, "Capital Gain"),
		col.New(&colfmt.Float{W: colWidth, Prec: colPrec}, "CGT"),
	)

	totVal := 0.0
	totCG := 0.0
	totTax := 0.0

	const penceToPounds = 100

	for _, t := range trades {
		tradeVal := t.num * t.px / penceToPounds
		totVal += tradeVal
		gain := t.num * (t.px - prog.purchasePx) / penceToPounds
		totCG += gain
		cgt := gain * mathutil.FromPercent(prog.cgtRate)
		totTax += cgt

		if err := rpt.PrintRow(t.num, t.px, tradeVal, gain, cgt); err != nil {
			prog.setExitStatus(1)
			fmt.Println("Unexpected error found while printing a row: ", err)

			return
		}
	}

	const skipCols = 2
	if err = rpt.PrintFooterVals(skipCols, totVal, totCG, totTax); err != nil {
		prog.setExitStatus(1)
		fmt.Println("Unexpected error found while printing the report footer: ",
			err)

		return
	}
}

// parseTradePart attempts to parse 'text' into a float64 and returns an
// error reporting the name and location if the parsing fails.
func (prog *prog) parseTradePart(loc *location.L, name, text string,
) (float64, error) {
	v, err := strconv.ParseFloat(text, 64)
	if err != nil {
		prog.setExitStatus(1)

		return 0.0, fmt.Errorf("%s: couldn't parse the %s: %s", loc, name, err)
	}

	return v, err
}

// getTrades reads the trades from the file
func (prog *prog) getTrades() ([]trade, error) {
	const expectedFieldCount = 2

	trades := []trade{}

	f, err := os.Open(prog.filename)
	if err != nil {
		prog.setExitStatus(1)

		return trades, fmt.Errorf(
			"cannot open file: %q: %s", prog.filename, err)
	}

	tradeScanner := bufio.NewScanner(f)

	re := regexp.MustCompile(`\s`)

	loc := location.New(prog.filename)
	for tradeScanner.Scan() {
		loc.Incr()

		parts := re.Split(tradeScanner.Text(), -1)
		if len(parts) != expectedFieldCount {
			prog.setExitStatus(1)

			return trades, fmt.Errorf(
				"%s: cannot split the line into two parts", loc)
		}

		var t trade

		t.num, err = prog.parseTradePart(loc, "trade amount", parts[0])
		if err != nil {
			return trades, err
		}

		t.px, err = prog.parseTradePart(loc, "trade price", parts[1])
		if err != nil {
			return trades, err
		}

		trades = append(trades, t)
	}

	if err := tradeScanner.Err(); err != nil {
		prog.setExitStatus(1)

		return trades, fmt.Errorf("error reading %q : %v", prog.filename, err)
	}

	return trades, nil
}

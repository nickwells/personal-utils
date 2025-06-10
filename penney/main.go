package main

import (
	"fmt"
	"os"

	"github.com/nickwells/col.mod/v5/col"
	"github.com/nickwells/col.mod/v5/colfmt"
	"github.com/nickwells/mathutil.mod/v2/mathutil"
)

// Created: Fri Sep 29 00:27:06 2023

const (
	choiceWidth = 7
	choicePrec  = 2
)

func main() {
	prog := NewProg()
	ps := makeParamSet(prog)
	ps.Parse()

	allPossibleChoices := prog.makeAllPossibleChoices()

	if prog.tryAll {
		cols := []*col.Col{}
		for _, choice := range allPossibleChoices {
			cols = append(cols,
				col.New(&colfmt.Float{W: choiceWidth, Prec: choicePrec},
					prog.uintToStr(choice)))
		}

		rpt := col.NewReportOrPanic(col.NewHeaderOrPanic(),
			os.Stdout,
			col.New(&colfmt.String{}, "P1/P2").SetSep(": "),
			cols...)

		for i := range len(allPossibleChoices) {
			p2Choice := allPossibleChoices[i]
			p2Choices := make([]uint, len(allPossibleChoices))

			for idx := range len(allPossibleChoices) {
				p2Choices[idx] = p2Choice
			}

			p1 := newPlayer("P1", allPossibleChoices)
			p2 := newPlayer("P2", p2Choices)

			prog.play(p1, p2)

			vals := []any{
				prog.uintToStr(p2Choice),
			}

			for idx, p1R := range p1.r {
				p2R := p2.r[idx]

				if p1.choices[idx] == p2Choice {
					vals = append(vals, col.Skip{})
				} else {
					vals = append(vals, float64(p2R.myWins)/float64(p1R.myWins))
				}
			}

			err := rpt.PrintRow(vals...)
			if err != nil {
				fmt.Println("error printing report: ", err)
			}
		}
	} else {
		p1 := newPlayer("P1", allPossibleChoices)
		p2 := newPlayer("P2", prog.makeOtherChoices(p1.choices))
		prog.play(p1, p2)

		prog.reportResults(p1, p2)
	}
}

// reportResults reports the results for the two players
//
//nolint:mnd
func (prog Prog) reportResults(p1, p2 *player) {
	perChoiceCols := 4

	if prog.showWinCount {
		perChoiceCols++
	}

	if prog.showRunInfo {
		perChoiceCols += 2
	}

	cols := make([]*col.Col, 0, prog.choiceCount()*perChoiceCols)

	for range len(p1.choices) {
		cols = append(cols,
			col.New(
				&colfmt.String{
					W: uint(prog.coinCount), //nolint:gosec
				},
				"chc"))

		if prog.showWinCount {
			maxWinWidth := mathutil.Digits(prog.trials)
			cols = append(cols,
				col.New(
					&colfmt.Int{
						W: uint(maxWinWidth), //nolint:gosec
					},
					"wins"))
		}

		pctCol := col.New(
			&colfmt.Percent{W: 7, Prec: 2, SuppressPct: true}, "%age")
		if prog.showRunInfo {
			cols = append(cols, pctCol)
			cols = append(cols, col.New(&colfmt.Int{W: 3}, "max", "run"))
			cols = append(cols, col.New(
				&colfmt.Float{W: 5, Prec: 1},
				"avg", "run").SetSep(" | "))
		} else {
			cols = append(cols, pctCol.SetSep(" | "))
		}
	}

	hdr := col.NewHeaderOrPanic()
	rpt := col.NewReportOrPanic(hdr,
		os.Stdout,
		col.New(&colfmt.String{}, "player").SetSep(": "),
		cols...)

	p1.reportResults(rpt, prog)
	p2.reportResults(rpt, prog)

	if prog.showExcess {
		prog.reportExcess(rpt, p1, p2)
	}
}

// reportExcess shows the difference between the percentages for the two
// Players
func (prog Prog) reportExcess(rpt *col.Report, p1, p2 *player) {
	vals := []any{p2.ID + "-" + p1.ID}

	for i, p1R := range p1.r {
		p2R := p2.r[i]

		vals = append(vals, col.Skip{})

		if prog.showWinCount {
			vals = append(vals, col.Skip{})
		}

		vals = append(vals, p2R.percVal(prog)-p1R.percVal(prog))

		if prog.showRunInfo {
			vals = append(vals, col.Skip{}, col.Skip{})
		}
	}

	err := rpt.PrintRow(vals...)
	if err != nil {
		fmt.Println("error printing report: ", err)
	}
}

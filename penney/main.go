package main

import (
	"os"

	"github.com/nickwells/col.mod/v4/col"
	"github.com/nickwells/col.mod/v4/colfmt"
	"github.com/nickwells/mathutil.mod/v2/mathutil"
)

// Created: Fri Sep 29 00:27:06 2023

func main() {
	prog := NewProg()
	ps := makeParamSet(prog)
	ps.Parse()

	p1 := NewPlayer("P1", prog.makeAllPossibleChoices())

	if prog.tryAll {
		allOtherChoices := prog.makeAllOtherChoices(p1.choices)

		for i, c := range p1.choices {
			dupChoice := make([]uint, len(p1.choices)-1)
			for idx := range dupChoice {
				dupChoice[idx] = c
			}

			p1rpt := NewPlayer(p1.ID, dupChoice)
			p2 := NewPlayer("P2", allOtherChoices[i])

			prog.play(p1rpt, p2)

			prog.reportResults(p1rpt, p2)
		}
	} else {
		p2 := NewPlayer("P2", prog.makeOtherChoices(p1.choices))
		prog.play(p1, p2)

		prog.reportResults(p1, p2)
	}
}

// reportResults reports the results for the two players
//
//nolint:mnd
func (prog Prog) reportResults(p1, p2 *Player) {
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
			col.New(colfmt.String{W: uint(prog.coinCount)}, "chc")) //nolint:gosec

		if prog.showWinCount {
			maxWinWidth := mathutil.Digits(prog.trials)
			cols = append(cols,
				col.New(colfmt.Int{W: uint(maxWinWidth)}, "wins")) //nolint:gosec
		}

		pctCol := col.New(&colfmt.Percent{W: 7, Prec: 2}, "%age")
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
		col.New(colfmt.String{}, "player").SetSep(": "),
		cols...)

	p1.reportResults(rpt, prog)
	p2.reportResults(rpt, prog)
}

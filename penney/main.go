package main

import (
	"os"

	"github.com/nickwells/col.mod/v3/col"
	"github.com/nickwells/col.mod/v3/col/colfmt"
	"github.com/nickwells/mathutil.mod/v2/mathutil"
)

// Created: Fri Sep 29 00:27:06 2023

func main() {
	prog := NewProg()
	ps := makeParamSet(prog)
	ps.Parse()

	p1 := NewPlayer("P1", prog.makeAllPossibleChoices())
	p2 := NewPlayer("P2", prog.makeOtherChoices(p1.choices))
	prog.play(p1, p2)

	cols := make([]*col.Col, 0, prog.choiceCount()*3)
	for i := 0; i < prog.choiceCount(); i++ {
		cols = append(cols, col.New(colfmt.String{W: prog.coinCount}, "chc"))
		cols = append(cols,
			col.New(colfmt.Int{W: mathutil.Digits[int](prog.trials)}, "wins"))
		cols = append(cols,
			col.New(&colfmt.Percent{W: 6, Prec: 2}, "%age"))
		cols = append(cols,
			col.New(&colfmt.Int{W: 3}, "max", "run"))
		cols = append(cols,
			col.New(&colfmt.Float{W: 5, Prec: 1}, "avg", "run"))
	}
	hdr := col.NewHeaderOrPanic()
	rpt := col.NewReport(hdr,
		os.Stdout,
		col.New(colfmt.String{}, "player"),
		cols...)
	p1.reportResults(rpt, *prog)
	p2.reportResults(rpt, *prog)
}

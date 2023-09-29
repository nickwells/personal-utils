package main

import (
	"fmt"

	"github.com/nickwells/col.mod/v3/col"
)

type Player struct {
	ID string

	choices []uint
	r       []Results
}

// NewPlayer returns a new Player
func NewPlayer(id string, choices []uint) *Player {
	if id == "" {
		panic("bad id - must not be empty")
	}

	p := &Player{
		ID:      id,
		choices: choices,
	}

	p.r = make([]Results, len(choices))
	for i, r := range p.r {
		r.ID = id
		p.r[i] = r
	}

	return p
}

// reportResults reports each of the results
func (p Player) reportResults(rpt *col.Report, prog Prog) {
	vals := make([]any, 0, 1+prog.choiceCount()*3)
	vals = append(vals, p.ID)
	for i, r := range p.r {
		r.notify(0, "")
		vals = append(vals,
			prog.uintToStr(p.choices[i]),
			r.myWins,
			float64(r.myWins)/float64(r.totalWins),
			r.maxRunLength,
			float64(r.totalRunLength)/float64(r.runCount))
	}
	err := rpt.PrintRow(vals...)
	if err != nil {
		fmt.Println("error printing report: ", err)
	}
}

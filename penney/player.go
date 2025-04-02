package main

import (
	"fmt"

	"github.com/nickwells/col.mod/v4/col"
)

// player represents a player of the game
type player struct {
	ID string

	choices []uint
	r       []results
}

// newPlayer returns a new Player
func newPlayer(id string, choices []uint) *player {
	if id == "" {
		panic("bad id - must not be empty")
	}

	p := &player{
		ID:      id,
		choices: choices,
	}

	p.r = make([]results, len(choices))
	for i, r := range p.r {
		r.ID = id
		p.r[i] = r
	}

	return p
}

// reportResults reports each of the results
func (p player) reportResults(rpt *col.Report, prog Prog) {
	vals := make([]any, 0, 1+prog.choiceCount()*3)

	vals = append(vals, p.ID)

	for i, r := range p.r {
		r.notify(0, "")

		vals = append(vals, prog.uintToStr(p.choices[i]))

		if prog.showWinCount {
			vals = append(vals, r.myWins)
		}

		vals = append(vals, r.percVal(prog))

		if prog.showRunInfo {
			vals = append(vals,
				r.maxRunLength,
				float64(r.totalRunLength)/float64(r.runCount))
		}
	}

	err := rpt.PrintRow(vals...)
	if err != nil {
		fmt.Println("error printing report: ", err)
	}
}

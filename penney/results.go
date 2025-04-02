package main

import "github.com/nickwells/mathutil.mod/v2/mathutil"

// results captures the results of the game
type results struct {
	ID           string
	lastWinnerID string

	totalWins  int
	myWins     int
	totalFlips int
	maxFlips   int

	currentRunLength int
	maxRunLength     int
	totalRunLength   int
	runCount         int
}

// notify checks to see if the ID matches the winner and adds to the winning
// results if so. If not, it checks to see if the ID matches the last winner
// and if so it updates the runLengths.
func (r *results) notify(flips int, winnerID string) {
	if flips > 0 {
		r.totalWins++
	}

	if winnerID == r.ID {
		r.myWins++
		r.totalFlips += flips

		if flips > r.maxFlips {
			r.maxFlips = flips
		}

		r.currentRunLength++
	} else {
		if r.lastWinnerID == r.ID {
			r.runCount++
			r.totalRunLength += r.currentRunLength

			if r.currentRunLength > r.maxRunLength {
				r.maxRunLength = r.currentRunLength
			}
		}

		r.currentRunLength = 0
	}

	r.lastWinnerID = winnerID
}

// percVal returns the players wins as a proportion of the total wins. If the
// program showRoughly flag is set then the result is passed to the
// mathutil.Roughly func before being returned.
func (r *results) percVal(prog Prog) float64 {
	percVal := float64(r.myWins) / float64(r.totalWins)

	if prog.showRoughly {
		percVal = mathutil.Roughly(percVal, 1)
	}

	return percVal
}

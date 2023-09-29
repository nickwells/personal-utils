package main

type Results struct {
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
func (r *Results) notify(flips int, winnerID string) {
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

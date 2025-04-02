package main

import (
	"math/rand/v2"
)

// Prog holds program parameters and status
type Prog struct {
	trials             int
	coinCount          int
	copyCount          int // must be < coinCount
	leadingPickupShift int // must be < coinCount
	tryAll             bool
	showWinCount       bool
	showRunInfo        bool
	showRoughly        bool
	showExcess         bool
}

// NewProg returns a new Prog instance with the default values set
func NewProg() *Prog {
	const (
		dfltTrials    = 1000
		dfltCoinCount = 3
		dfltCopyCount = 2
	)

	return &Prog{
		trials:             dfltTrials,
		coinCount:          dfltCoinCount,
		copyCount:          dfltCopyCount,
		leadingPickupShift: 1,
	}
}

// uintToStr converts a uint, less than the limit into a string of
// H's and T's
func (prog Prog) uintToStr(v uint) string {
	coins := [2]string{"H", "T"}
	val := ""

	for i := prog.coinCount - 1; i >= 0; i-- {
		cIdx := (v & (1 << i)) >> i
		val += coins[cIdx]
	}

	return val
}

// choiceCount returns the number of choices corresponding to the coinCount
func (prog Prog) choiceCount() int {
	return 1 << prog.coinCount
}

// makeAllPossibleChoices creates the full set of choices of values
func (prog Prog) makeAllPossibleChoices() []uint {
	limit := prog.choiceCount()
	choices := make([]uint, 0, limit)

	for i := range limit {
		choices = append(choices, uint(i)) //nolint:gosec
	}

	return choices
}

// makeOtherChoices constructs the other (winning) choices given the choices
// of the first player
func (prog Prog) makeOtherChoices(choices []uint) []uint {
	otherChoices := make([]uint, 0, len(choices))
	shift := prog.coinCount - prog.copyCount
	shiftMask := makeShiftMask(prog.coinCount, shift, prog.leadingPickupShift)

	for _, c := range choices {
		oc := c >> shift
		leadingBits := ((^c) & shiftMask) << uint(prog.leadingPickupShift) //nolint:gosec
		oc |= leadingBits
		otherChoices = append(otherChoices, oc)
	}

	return otherChoices
}

// makeBitMask returns a bit-mask of length bitCount
func makeBitMask(bitCount int) uint {
	var bm uint
	for range bitCount {
		bm = (bm << 1) | 1
	}

	return bm
}

// makeShiftMask returns a bit-mask covering just the unshifted bits in the
// value
func makeShiftMask(maxLen, shift, leadingPickupShift int) uint {
	var bm uint

	for i := range maxLen {
		var nextBit uint = 1

		if i >= shift {
			nextBit = 0
		}

		bm = (bm << 1) | nextBit
	}

	return bm >> uint(leadingPickupShift) //nolint:gosec
}

// play runs the trials collecting the results in the players results fields
func (prog Prog) play(p1, p2 *player) {
	flips := make([]int, len(p1.choices))

	var match uint

	mask := makeBitMask(prog.coinCount)

	for range prog.trials {
		toss := uint(rand.IntN(2)) //nolint:gosec,mnd
		match <<= 1
		match |= toss
		match &= mask

		for c, fc := range flips {
			fc++

			if fc >= int(prog.coinCount) {
				switch match {
				case p1.choices[c]:
					p1.r[c].notify(fc, p1.ID)
					p2.r[c].notify(fc, p1.ID)
					fc = 0
				case p2.choices[c]:
					p1.r[c].notify(fc, p2.ID)
					p2.r[c].notify(fc, p2.ID)
					fc = 0
				}
			}

			flips[c] = fc
		}
	}
}

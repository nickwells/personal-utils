package main

import (
	"math/rand"
)

// Prog holds program parameters and status
type Prog struct {
	trials             int
	coinCount          int
	copyCount          int // must be < coinCount
	leadingPickupShift int // must be < coinCount
	tryAll             bool
	showWinCount       bool
}

// NewProg returns a new Prog instance with the default values set
func NewProg() *Prog {
	return &Prog{
		trials:             1000,
		coinCount:          3,
		copyCount:          2,
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
	for i := 0; i < limit; i++ {
		choices = append(choices, uint(i))
	}
	return choices
}

// makeOtherChoices constructs the other (winning) choices given the choices
// of the first player
func (prog Prog) makeOtherChoices(choices []uint) []uint {
	otherChoices := make([]uint, 0, len(choices))
	shift := prog.coinCount - prog.copyCount
	shiftMask := prog.makeShiftMask(shift)
	for _, c := range choices {
		oc := c >> shift
		leadingBits := ((^c) & shiftMask) << uint(prog.leadingPickupShift)
		oc |= leadingBits
		otherChoices = append(otherChoices, oc)
	}
	return otherChoices
}

// makeBitMask returns a bit-mask covering all the bits in the value
func (prog Prog) makeBitMask() uint {
	var bm uint
	for i := 0; i < prog.coinCount; i++ {
		bm = (bm << 1) | 1
	}
	return bm
}

// makeShiftMask returns a bit-mask covering just the unshifted bits in the
// value
func (prog Prog) makeShiftMask(shift int) uint {
	var bm uint
	for i := 0; i < prog.coinCount; i++ {
		var nextBit uint = 1
		if i >= shift {
			nextBit = 0
		}
		bm = (bm << 1) | nextBit
	}
	return bm >> uint(prog.leadingPickupShift)
}

// play runs the trials collecting the results in the players results fields
func (prog Prog) play(p1, p2 *Player) {
	flips := make([]int, len(p1.choices))
	var match uint
	mask := prog.makeBitMask()

	for i := 0; i < prog.trials; i++ {
		toss := uint(rand.Intn(2))
		match <<= 1
		match |= toss
		match &= mask

		for c, fc := range flips {
			fc++
			if fc >= int(prog.coinCount) {
				if match == p1.choices[c] {
					p1.r[c].notify(fc, p1.ID)
					p2.r[c].notify(fc, p1.ID)
					fc = 0
				} else if match == p2.choices[c] {
					p1.r[c].notify(fc, p2.ID)
					p2.r[c].notify(fc, p2.ID)
					fc = 0
				}
			}
			flips[c] = fc
		}
	}
}

package main

import (
	"fmt"
)

// CurrencyCode represents an ISO 4217 currency code. This is a three letter
// code where the first two letters are the two letters of the ISO 3166-1
// alpha-2 country code and the last letter is derived (usually) from the
// initial letter of the currency name. Note that the unofficial codes for
// minor currency units are not supported.
type CurrencyCode string

// IsValid returns a non-nil error if the Currency code is not valid. A code
// must be 3 uppercase ASCII letters.
func (cc CurrencyCode) IsValid() error {
	const maxCodeLen = 3
	if len(cc) != maxCodeLen {
		return fmt.Errorf("%s the code length (%d) must equal %d",
			cc.errorIntro(), len(cc), maxCodeLen)
	}

	for i, r := range cc {
		if r < 'A' || r > 'Z' {
			return fmt.Errorf(
				"%s character %d (%q) must be a capital Latin letter",
				cc.errorIntro(), i, r)
		}
	}

	return nil
}

// errorIntro returns a common string for the start of a CurrencyCode error.
func (cc CurrencyCode) errorIntro() string {
	return fmt.Sprintf("bad CurrencyCode: %q:", cc)
}

// CheckCurrencyCode satisfies the check.ValCk type
func CheckCurrencyCode(cc CurrencyCode) error {
	return cc.IsValid()
}

// CurrencyPair records a symbol for a pair of currencies
type CurrencyPair struct {
	baseCcy    CurrencyCode
	counterCcy CurrencyCode
}

// MakeCurrencyPair returns a populated CurrencyPair
func MakeCurrencyPair(baseCcy, counterCcy CurrencyCode) CurrencyPair {
	return CurrencyPair{
		baseCcy:    baseCcy,
		counterCcy: counterCcy,
	}
}

// String returns a formatted currency pair string
func (cp CurrencyPair) String() string {
	return string(cp.baseCcy + cp.counterCcy)
}

// Swap generates the opposite pair
func (cp CurrencyPair) Swap() CurrencyPair {
	return CurrencyPair{
		baseCcy:    cp.counterCcy,
		counterCcy: cp.baseCcy,
	}
}

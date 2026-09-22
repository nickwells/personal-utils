package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/nickwells/tempus.mod/v2/tempus"
	"github.com/nickwells/verbose.mod/verbose"
	"github.com/nickwells/xdg.mod/xdg"
)

const (
	baseCcy = "GBP"

	cacheDirPerms = 0o750
)

// asOf records the Month and Year for which to report the exchange rate
type asOf struct {
	m time.Month
	y int
}

// initAsOf constructs the default asOf structure. It uses the current month
// and year unless the date is after the second last Thursday of the month in
// which case it advances the date by one month (possibly advancing the year
// if the date is in December).
func initAsOf() asOf {
	target := time.Now()

	domPenultimateThu, _ := tempus.NthWeekdayOfMonthYear(
		-2, time.Thursday,
		target.Month(), target.Year())
	if domPenultimateThu < target.Day() {
		target = tempus.AddMonth(1, target)
	}

	return asOf{
		y: target.Year(),
		m: target.Month(),
	}
}

// prog holds program parameters and status
type prog struct {
	exitStatus int
	stack      *verbose.Stack

	// parameters
	cacheFile    string
	useGivenFile bool

	cacheDir string

	from CurrencyCode
	to   CurrencyCode

	asOf asOf

	amount float64
	// program data
	rates     map[CurrencyPair]float64
	countries map[string]CurrencyCode
	ccyNames  map[CurrencyCode]string
}

// newProg returns a new Prog instance with the default values set
func newProg() *prog {
	return &prog{
		stack: &verbose.Stack{},

		cacheDir: mkCacheDirPath(),

		from: baseCcy,
		to:   baseCcy,

		asOf: initAsOf(),

		amount: 1,

		rates: map[CurrencyPair]float64{},
		countries: map[string]CurrencyCode{
			"Great Britain": baseCcy,
		},
		ccyNames: map[CurrencyCode]string{
			baseCcy: "Pound",
		},
	}
}

// mkCacheDirPath returns the name of the default cache directory
func mkCacheDirPath() string {
	ccyCachePath := []string{
		"github.com",
		"nickwells",
		"personal-utils",
		"exrates",
		"ccyRates",
	}

	return filepath.Join(append([]string{xdg.CacheHome()}, ccyCachePath...)...)
}

// run is the starting point for the program, it should be called from main()
// after the command-line parameters have been parsed. Use the setExitStatus
// method to record the exit status and then main can exit with that status.
func (prog *prog) run() {
	var err error

	defer func() {
		if err != nil {
			fmt.Println(err)

			prog.exitStatus = 1
		}
	}()

	cacheFilePath := prog.cacheFile
	if !prog.useGivenFile {
		cacheFilePath = filepath.Join(prog.cacheDir, prog.makeCcyFileName())
	}

	err = prog.populateRates(cacheFilePath)
	if err != nil {
		if prog.useGivenFile {
			prog.exitStatus = 1
			return
		}

		if errors.Is(err, os.ErrNotExist) {
			err = prog.fillRatesCacheFile()
		}

		if err != nil {
			return
		}

		err = prog.populateRates(cacheFilePath)
		if err != nil {
			return
		}
	}

	ccyPair := MakeCurrencyPair(prog.from, prog.to)

	rate, ok := prog.rates[ccyPair]
	if !ok {
		err = fmt.Errorf("currency pair %q is unknown", ccyPair)
		return
	}

	fmt.Printf("%s: %.2f %s = %.4f %s\n",
		ccyPair,
		prog.amount, prog.ccyNames[prog.from],
		prog.amount*rate, prog.ccyNames[prog.to])
}

// populateRates reads the contents of the given file and uses that to
// populate the currency details. It returns any errors detected
func (prog *prog) populateRates(filename string) error {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		return err
	}

	defer f.Close()

	reader := csv.NewReader(f)

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("read failure from %q: %s", filename, err)
	}

	const (
		idxCountryName = iota
		idxCcyName
		idxCcyCode
		idxRate
	)
	for _, rec := range records {
		ccyCode := CurrencyCode(rec[idxCcyCode])
		ccyName := rec[idxCcyName]
		countryName := rec[idxCountryName]
		rateStr := rec[idxRate]

		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			return fmt.Errorf(
				"parse failure from %q: Country: %q: rate: %s : %s",
				filename, countryName, rateStr, err)
		}

		prog.rates[MakeCurrencyPair("GBP", ccyCode)] = rate
		prog.rates[MakeCurrencyPair(ccyCode, "GBP")] = 1 / rate
		prog.countries[countryName] = CurrencyCode(ccyCode)
		prog.ccyNames[ccyCode] = ccyName
	}

	for ccy1 := range prog.ccyNames {
		if ccy1 == "GBP" {
			continue
		}

		fromGBP := prog.rates[MakeCurrencyPair("GBP", ccy1)]
		for ccy2 := range prog.ccyNames {
			if ccy2 == "GBP" {
				continue
			}

			if ccy1 == ccy2 {
				continue
			}

			toGBP := prog.rates[MakeCurrencyPair(ccy2, "GBP")]
			prog.rates[MakeCurrencyPair(ccy2, ccy1)] = toGBP * fromGBP
		}
	}

	return nil
}

// fillRatesCacheFile retrieves the currency rates from the target site and
// constructs a matrix of exchange rates between each pair of currencies
func (prog *prog) fillRatesCacheFile() error {
	const (
		siteURL = "https://www.trade-tariff.service.gov.uk"
		baseURL = siteURL + "/uk/api/exchange_rates/files"
	)

	fileName := prog.makeCcyFileName()
	url := baseURL + "/" + fileName

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return fmt.Errorf("download failure: %s: %s", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failure: %s: status: %s", url, resp.Status)
	}

	reader := csv.NewReader(resp.Body)

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("cannot read the response: %s", err)
	}

	err = os.MkdirAll(prog.cacheDir, cacheDirPerms)
	if err != nil {
		return fmt.Errorf("cannot create the cache directory; %q: %s",
			prog.cacheDir, err)
	}

	cacheFilePath := filepath.Join(prog.cacheDir, fileName)

	w, err := os.Create(cacheFilePath) //nolint:gosec
	if err != nil {
		return fmt.Errorf("cannot create the cache file; %q: %s",
			cacheFilePath, err)
	}

	csvWriter := csv.NewWriter(w)
	allRecords := [][]string{}

	const (
		idxCountryName = iota
		idxCcyName
		idxCcyCode
		idxRate
	)
	for _, rec := range records[1:] {
		allRecords = append(allRecords,
			[]string{
				rec[idxCountryName],
				rec[idxCcyName],
				rec[idxCcyCode],
				rec[idxRate],
			})
	}

	return csvWriter.WriteAll(allRecords)
}

// makeCcyFileName constructs the name of the latest currency file from the
// program asOf field which is initialised to take the current date possibly
// with some adjustment for the second-last-Thursday-of-the-month rule; see
// the initAsOf func for details.
func (prog *prog) makeCcyFileName() string {
	return fmt.Sprintf("monthly_csv_%d-%d.csv", prog.asOf.y, prog.asOf.m)
}

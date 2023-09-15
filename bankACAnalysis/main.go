package main

// bankACAnalysis

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nickwells/col.mod/v3/col"
	"github.com/nickwells/col.mod/v3/col/colfmt"
)

// Created: Sun May 12 16:39:24 2019

// Xactn represents a single transaction
type Xactn struct {
	lineNum   int
	date      time.Time
	xaType    string
	desc      string
	debitAmt  float64
	creditAmt float64
	balance   float64
}

// Summary represents a summary of the account transactions
type Summary struct {
	name       string
	count      int
	firstDate  time.Time
	lastDate   time.Time
	debitAmt   float64
	creditAmt  float64
	parent     *Summary
	depth      int
	components map[string]*Summary
}

const (
	catAll     = "all"
	catUnknown = "unknown"
	catCash    = "cash"
	catCheque  = "cheque"

	editTypeSearch  = "search"
	editTypeReplace = "replace"
)

const xactnMapDesc = "map of transaction types"

const tabWidth = 4

// Edit represents a substitution to be made to a transaction description
type Edit struct {
	search      string
	searchRE    *regexp.Regexp
	replacement string
}

type Summaries struct {
	parentOf     map[string]string
	summaries    map[string]*Summary
	edits        []Edit
	maxDepth     int
	maxNameWidth int
}

type reportStyle int

const (
	showLeafEntries reportStyle = iota
	summaryReport
)

// openFileOrDie will try to open the given file and will return the open
// file if successful and will print an error message and exit of not.
func openFileOrDie(fileName, desc string) *os.File {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Couldn't open the %s file: %s", desc, err)
		os.Exit(1)
	}
	return f
}

// populateParents constructs the parent tree of transactions from the
// transaction map file
func (s *Summaries) populateParents(prog *Prog) {
	s.parentOf[catAll] = catAll
	err := s.addParent(catAll, catUnknown)
	if err != nil {
		fmt.Printf("Cannot initialise the %s: %s\n", xactnMapDesc, err)
		os.Exit(1)
	}
	err = s.addParent(catAll, catCash)
	if err != nil {
		fmt.Printf("Cannot initialise the %s: %s\n", xactnMapDesc, err)
		os.Exit(1)
	}
	err = s.addParent(catAll, catCheque)
	if err != nil {
		fmt.Printf("Cannot initialise the %s: %s\n", xactnMapDesc, err)
		os.Exit(1)
	}

	mf := openFileOrDie(prog.xactMapFileName, xactnMapDesc)
	defer mf.Close()

	mScanner := bufio.NewScanner(mf)
	lineNum := 0

	for mScanner.Scan() {
		lineNum++
		line := mScanner.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		err = s.addParent(parts[0], parts[1])
		if err != nil {
			fmt.Printf("%s:%d: Bad entry in the %s: %s\n",
				prog.xactMapFileName, lineNum, xactnMapDesc, err)
		}
	}
}

// populateEdits constructs the slice of editing rules to be performed on
// transaction descriptions
func (s *Summaries) populateEdits(prog *Prog) {
	ef := openFileOrDie(prog.editFileName, "transaction edits")
	defer ef.Close()

	eScanner := bufio.NewScanner(ef)
	lineNum := 0
	prevType := ""
	var searchRE *regexp.Regexp
	var searchStr string
	var errFound bool
	var err error
	const errIntro = "Bad transaction edits entry"

	for eScanner.Scan() {
		lineNum++
		line := eScanner.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("%s:%d: %s: Missing '=' : %s\n",
				prog.editFileName, lineNum, errIntro, line)
			errFound = true
			continue
		}
		entryType := parts[0]
		switch entryType {
		case editTypeSearch:
			if prevType == editTypeSearch {
				fmt.Printf(
					"%s:%d: %s: %q entry missing for previous search\n",
					prog.editFileName, lineNum, errIntro, editTypeReplace)
			}
			errFound = false
			searchStr = parts[1]
			searchRE, err = regexp.Compile(searchStr)
			if err != nil {
				fmt.Printf("%s:%d: %s: Couldn't compile the regexp: %s\n",
					prog.editFileName, lineNum, errIntro, err)
				errFound = true
			}
		case editTypeReplace:
			if !errFound {
				s.edits = append(s.edits, Edit{
					search:      searchStr,
					searchRE:    searchRE,
					replacement: parts[1],
				})
			}
		default:
			fmt.Printf("%s:%d: %s: Bad type: %s\n",
				prog.editFileName, lineNum, errIntro, entryType)
			errFound = true
		}
		prevType = entryType
	}
}

// initSummaries returns an initialised Summaries structure
func (prog *Prog) initSummaries() *Summaries {
	s := Summaries{
		parentOf:  make(map[string]string),
		summaries: make(map[string]*Summary),
	}

	s.summaries[catAll] = &Summary{
		name:       catAll,
		components: make(map[string]*Summary),
	}
	s.populateParents(prog)

	s.populateEdits(prog)

	return &s
}

// addParent adds the parent/child relationship so that a given summary can
// find its parent. It is an error if the parent does not already exist.
func (s *Summaries) addParent(parent, child string) error {
	if _, ok := s.parentOf[parent]; !ok {
		return fmt.Errorf("%q (parent of %q) doesn't exist",
			parent, child)
	}

	if oldParent, ok := s.parentOf[child]; ok {
		if oldParent != parent {
			return fmt.Errorf("%q already has a parent: %q != %q",
				child, parent, oldParent)
		}
		return nil
	}

	pSum, ok := s.summaries[parent]
	if !ok {
		return fmt.Errorf("%q (parent of %q) has no summary record"+
			" - check the transaction map file",
			parent, child)
	}
	cSum := &Summary{
		name:       child,
		parent:     pSum,
		depth:      pSum.depth + 1,
		components: make(map[string]*Summary),
	}
	s.summaries[child] = cSum
	if cSum.depth > s.maxDepth {
		s.maxDepth = cSum.depth
	}
	if len(cSum.name) > s.maxNameWidth {
		s.maxNameWidth = len(cSum.name)
	}

	pSum.components[child] = cSum
	s.parentOf[child] = parent
	return nil
}

// summarise will summarise the transaction working its way up to the top of
// the tree of Summary records
func (s *Summaries) summarise(xa Xactn) {
	summ, ok := s.summaries[xa.desc]
	if !ok {
		fmt.Println("Couldn't find the summary record for :", xa)
		return
	}
	summ.add(xa)
}

// add will add the values to the summary record and move on to the parent
// (if there is one)
func (s *Summary) add(xa Xactn) {
	if s.count == 0 {
		s.firstDate = xa.date
		s.lastDate = xa.date
	} else {
		if xa.date.After(s.lastDate) {
			s.lastDate = xa.date
		} else if s.firstDate.After(xa.date) {
			s.firstDate = xa.date
		}
	}
	s.count++
	s.debitAmt += xa.debitAmt
	s.creditAmt += xa.creditAmt
	if s.parent != nil {
		s.parent.add(xa)
	}
}

// Prog holds the parameters and current status of the program
type Prog struct {
	// the name of a file containing  transactions
	acFileName string
	// all the transaction files, including any given after the terminal
	// parameter
	files []string

	// the name of the file containing the mappings between transaction names
	// and categories
	xactMapFileName string

	// the name of the file containing the replacements to make to
	// transaction names
	editFileName string

	// don't suppress printing of summary records for which there are no
	// transactions
	showZeros bool

	// Skip the first (header) line in the file of transactions.
	skipFirstLine bool

	style         reportStyle
	minimalAmount float64
	showCats      []string
}

func NewProg() *Prog {
	return &Prog{
		skipFirstLine: true,
		style:         showLeafEntries,
		showCats:      []string{catAll},
	}
}

func main() {
	prog := NewProg()
	ps := makeParamSet(prog)

	ps.Parse()

	prog.files = ps.Remainder()
	if prog.acFileName != "" {
		prog.files = append(prog.files, prog.acFileName)
	}
	summaries := prog.getAccountData()

	sep := ""
	for _, cat := range prog.showCats {
		fmt.Print(sep)
		sep = "\n"
		summaries.report(prog, cat)
	}
}

// getAccountData opens each file in turn and reads from it to populate the
// summaries
func (prog *Prog) getAccountData() *Summaries {
	prog.checkFiles()

	s := prog.initSummaries()

	for _, name := range prog.files {
		f := openFileOrDie(name, "bank account")
		r := csv.NewReader(f)
		s.populateSummaries(prog, name, r)
		f.Close()
	}
	return s
}

// checkFiles checks the slice of files and if a duplicate is found it will
// report an error and exit
func (prog *Prog) checkFiles() {
	if len(prog.files) == 0 {
		fmt.Println("Some account files must be given")
		os.Exit(1)
	}

	m := map[string]bool{}
	var dupFound int
	for _, f := range prog.files {
		if m[f] {
			fmt.Println("File name", f,
				"appears more than once in the list of files")
			dupFound++
		}
		m[f] = true
	}
	if dupFound > 0 {
		os.Exit(1)
	}
}

// populateSummaries fills in the summaries from the lines read from the
// io.Reader
func (s *Summaries) populateSummaries(prog *Prog, name string, r *csv.Reader) {
	lineNum := 0
	for {
		parts, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error found while reading:", name)
			fmt.Println(err)
			os.Exit(1)
		}

		lineNum++
		if prog.skipFirstLine && lineNum == 1 {
			continue // ignore the first line of headings
		}
		xa, err := s.mkXactn(lineNum, parts)
		if err != nil {
			fmt.Println(err)
			continue
		}
		s.createNewMapEntries(name, lineNum, xa)
		s.summarise(xa)
	}
}

// createNewMapEntries will create new parent/child map entries for the
// transaction if it is not already known or if it is a cheque or cashpoint
// withdrawal
func (s *Summaries) createNewMapEntries(fileName string, lineNum int, xa Xactn) {
	if xa.xaType == "CHQ" {
		err := s.addParent(catCheque, xa.desc)
		if err != nil {
			fmt.Printf(
				"%s:%d: Can't add the cheque to the %s: %s\n",
				fileName, lineNum, xactnMapDesc, err)
		}
	} else if xa.xaType == "CPT" {
		err := s.addParent(catCash, xa.desc)
		if err != nil {
			fmt.Printf(
				"%s:%d: Can't add the cashpoint withdrawal to the %s: %s\n",
				fileName, lineNum, xactnMapDesc, err)
		}
	} else {
		if _, ok := s.parentOf[xa.desc]; ok {
			return
		}

		err := s.addParent(catUnknown, xa.desc)
		if err != nil {
			fmt.Printf(
				"%s:%d: Can't add the unknown entry to the %s: %s\n",
				fileName, lineNum, xactnMapDesc, err)
		}
	}
}

// normalise converts the string into a 'normal' form - this involves editing
// it to replace multiple alternative spellings into a single variant. It
// returns after all the edits have been applied
func (s *Summaries) normalise(str string) string {
	for _, ed := range s.edits {
		str = ed.searchRE.ReplaceAllLiteralString(str, ed.replacement)
	}
	return str
}

// report will report the summaries
func (s *Summaries) report(prog *Prog, cat string) {
	summ, ok := s.summaries[cat]
	if !ok {
		fmt.Printf("*** category: %q is not recognised\n", cat)
		return
	}
	floatCol := colfmt.Float{
		W:    10,
		Prec: 2,
		Zeroes: &colfmt.FloatZeroHandler{
			Handle:  true,
			Replace: "",
		},
	}
	pctCol := colfmt.Percent{
		W: 5,
		Zeroes: &colfmt.FloatZeroHandler{
			Handle:  true,
			Replace: "",
		},
	}

	rpt := col.StdRpt(
		col.New(colfmt.String{W: tabWidth*s.maxDepth + s.maxNameWidth},
			"Transaction Type"),
		col.New(colfmt.Int{W: 5}, "Count"),
		col.New(&colfmt.Time{Format: "2006-Jan-02"},
			"Date of", "First", "Transaction"),
		col.New(&colfmt.Time{Format: "2006-Jan-02"},
			"Date of", "Last", "Transaction"),
		col.New(&floatCol, "Debit", "Amount"),
		col.New(&pctCol, "%age"),
		col.New(&floatCol, "Credit", "Amount"),
		col.New(&pctCol, "%age"),
		col.New(&floatCol, "Nett", "Amount"),
	)

	summ.report(prog, rpt, summ.debitAmt, summ.creditAmt, 0)
}

// calcPct calculates the amount as a proportion of the total, if the total
// is zero, the proportion is zero regardless of the amount
func calcPct(amt, tot float64) float64 {
	if tot == 0 {
		return 0
	}
	return amt / tot
}

func (s *Summary) report(
	prog *Prog,
	rpt *col.Report,
	totDebit, totCredit float64,
	indent int,
) {
	if prog.style == summaryReport && len(s.components) == 0 {
		return
	}
	if !prog.showZeros && s.count == 0 {
		return
	}
	if s.creditAmt+s.debitAmt < prog.minimalAmount {
		return
	}

	err := rpt.PrintRow(
		strings.Repeat(" ", tabWidth*indent)+s.name,
		s.count,
		s.firstDate, s.lastDate,
		s.debitAmt, calcPct(s.debitAmt, totDebit),
		s.creditAmt, calcPct(s.creditAmt, totCredit),
		s.creditAmt-s.debitAmt)
	if err != nil {
		fmt.Println("Couldn't print the row:", err)
	}

	compList := []*Summary{}
	for _, c := range s.components {
		compList = append(compList, c)
	}
	sort.Slice(compList, func(i, j int) bool {
		return (compList[i].debitAmt + compList[i].creditAmt) >
			(compList[j].debitAmt + compList[j].creditAmt)
	})
	for _, c := range compList {
		c.report(prog, rpt, totDebit, totCredit, indent+1)
	}
}

// parseNum returns 0.0 if the string is empty, otherwise it will parse the
// number as a float
func parseNum(s, name string) (float64, error) {
	if s == "" {
		return 0.0, nil
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, fmt.Errorf("Couldn't parse the %s: %s", name, err)
	}
	return n, nil
}

// mkXactn converts the slice of strings into an transaction record
func (s *Summaries) mkXactn(lineNum int, parts []string) (Xactn, error) {
	date, err := time.Parse("02/01/2006", parts[0])
	if err != nil {
		return Xactn{}, fmt.Errorf("Couldn't parse the date: %s", err)
	}

	da, err := parseNum(parts[5], "debit amount")
	if err != nil {
		return Xactn{}, err
	}

	ca, err := parseNum(parts[6], "debit amount")
	if err != nil {
		return Xactn{}, err
	}

	bal, err := parseNum(parts[7], "balance amount")
	if err != nil {
		return Xactn{}, err
	}

	desc := parts[4]
	if _, ok := s.parentOf[desc]; !ok {
		desc = s.normalise(desc)
	}

	return Xactn{
		lineNum:   lineNum,
		date:      date,
		xaType:    parts[1],
		desc:      desc,
		debitAmt:  da,
		creditAmt: ca,
		balance:   bal,
	}, nil
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"

	check "github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v6/param"
)

// Prog holds the parameters and current status of the program
type Prog struct {
	files    []string
	provisos filecheck.Provisos
}

// HandleRemainder checks that each trailing argument is a non-empty Go file
// and adds them to the program file list. It records an error if any
// parameter is not a Go file.
func (p *Prog) HandleRemainder(ps *param.PSet, _ *location.L) {
	for _, fileName := range ps.Remainder() {
		if err := p.provisos.StatusCheck(fileName); err != nil {
			ps.AddErr("bad file", err)
			continue
		}

		p.files = append(p.files, fileName)
	}

	if len(p.files) == 0 {
		ps.AddErr("no files",
			errors.New("at least one Go file must be supplied"))
	}
}

// NewProg returns a new Prog value with any initial values set
func NewProg() *Prog {
	provisos := filecheck.FileNonEmpty()
	provisos.Checks = append(provisos.Checks,
		check.FileInfoName(check.StringHasSuffix[string](".go")))

	return &Prog{
		provisos: provisos,
	}
}

func main() {
	prog := NewProg()
	ps := makeParamSet(prog)

	ps.Parse()

	prog.files = ps.Remainder()

	var lineNum int

	lastLine := ""

	var (
		depRE       = regexp.MustCompile("^// Deprecated:")
		blankLineRE = regexp.MustCompile("^// *$")
	)

	for _, filename := range prog.files {
		f, err := os.Open(filename) //nolint:gosec
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening: %q : %v\n", filename, err)
			continue
		}

		lineNum = 0

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			if depRE.MatchString(line) {
				if !blankLineRE.MatchString(lastLine) {
					fmt.Println(filename, ":", lineNum, ":", line)
				}
			}

			lastLine = line
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %q : %v\n", filename, err)
		}

		f.Close()
	}
}

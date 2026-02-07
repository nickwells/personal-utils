package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	check "github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/filecheck.mod/filecheck"
)

// Prog holds the parameters and current status of the program
type Prog struct {
	files    []string
	provisos filecheck.Provisos
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
	ps := makeParamSet()

	ps.Parse()

	for _, fileName := range ps.TrailingParams() {
		if err := prog.provisos.StatusCheck(fileName); err != nil {
			fmt.Fprint(os.Stderr, "bad file:", err)
			continue
		}

		prog.files = append(prog.files, fileName)
	}

	if len(prog.files) == 0 {
		fmt.Fprint(os.Stderr, "at least one Go file must be supplied")
		os.Exit(1)
	}

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

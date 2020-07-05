package main

import (
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/location.mod/location"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// addParams will add parameters to the passed ParamSet
func addParams(ps *param.PSet) error {
	ps.Add("ac-file",
		psetter.Pathname{
			Value:       &acFileName,
			Expectation: filecheck.FileExists(),
		},
		"the name of the file containing the bank account transactions."+
			" This can also be given as a list of files after a "+
			ps.TerminalParam()+" parameter."+
			"\n"+
			" The file is expected to contain lines of comma-separated"+
			" values with the values as follows:\n\n"+
			"transaction date in the form DD/MM/YYYY\n"+
			"transaction type\n"+
			"sort-code\n"+
			"account number\n"+
			"transaction description\n"+
			"debit amount\n"+
			"credit amount\n"+
			"balance",
	)

	ps.Add("map-file",
		psetter.Pathname{
			Value:       &xactMapFileName,
			Expectation: filecheck.FileExists(),
		},
		"the name of the file containing the transaction name map.\n\n"+
			"Each non-blank line in the file should contain a word"+
			" representing the 'parent' group of transactions"+
			" followed by a space and the rest of the line which"+
			" represents the 'child' group of transactions.\n\n"+
			"There is an initial group called '"+catAll+"' with a child,"+
			" called '"+catUnknown+"' and the entries in this file are"+
			" intended"+
			" to construct the tree of transaction groups. Any"+
			" transaction description which is not found in this map"+
			" will automatically be placed in the 'unknown' group so you"+
			" can find the transactions you haven't classified by"+
			" looking in that group. In"+
			" order to create a new group you make an entry in this"+
			" file with parent set to 'all' and child set to the new"+
			" group name. Then each transaction that you want to put in"+
			" that group should have an entry with parent set to the"+
			" group name and the child set to the transaction"+
			" description. Groups can be nested to an arbitrary depth.",
		param.Attrs(param.MustBeSet))

	ps.Add("edit-file",
		psetter.Pathname{
			Value:       &editFileName,
			Expectation: filecheck.FileExists(),
		},
		"the name of the file containing the transaction name"+
			" replacements. Transaction descriptions that are not mapped"+
			" will be edited according to the rules in this file.\n\n"+
			"Each editing rule is given by a pair of lines,"+
			" the first must start with '"+editTypeSearch+"='"+
			" and the second must start with '"+editTypeReplace+"='."+
			" The first line value should be a valid regular expression",
		param.Attrs(param.MustBeSet))

	ps.Add("show-zeroes", psetter.Bool{Value: &showZeros},
		"don't suppress entries which have no transactions")

	ps.Add("dont-skip-line1",
		psetter.Bool{
			Value:  &skipFirstLine,
			Invert: true,
		},
		"don't ignore the first line of the transactions file")

	ps.Add("summary", psetter.Nil{},
		"show a summary report with no leaf transactions",
		param.PostAction(func(_ location.L,
			_ *param.ByName,
			_ []string) error {
			style = summaryReport
			return nil
		}))

	ps.Add("show-categories",
		psetter.StrList{
			Value: &showCats,
			Checks: []check.StringSlice{
				check.StringSliceLenGT(0),
			},
		},
		"show the report only for the listed categories",
		param.AltName("show-cats"),
		param.AltName("cats"),
	)

	ps.Add("minimal-amount", psetter.Float64{Value: &minimalAmount},
		"don't show summaries where the total transactions are less than this")

	// allow trailing arguments
	err := ps.SetNamedRemHandler(param.NullRemHandler{}, "bank-AC files")
	if err != nil {
		return err
	}

	return nil
}

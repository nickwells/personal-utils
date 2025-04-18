package main

// Code generated by mkparamfilefunc; DO NOT EDIT.
// with parameters set at:
//	[command line]: Argument:2: "-funcs" "personalOnly"
import (
	"path/filepath"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/xdg.mod/xdg"
)

/*
SetConfigFile adds a config file to the set which the param parser will process
before checking the command line parameters.
*/
func SetConfigFile(ps *param.PSet) error {
	baseDir := xdg.ConfigHome()

	ps.AddConfigFileStrict(
		filepath.Join(baseDir,
			"github.com",
			"nickwells",
			"personal-utils",
			"retirement",
			"common.cfg"),
		filecheck.Optional)

	return nil
}

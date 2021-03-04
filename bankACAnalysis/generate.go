// +build generate

package main

//go:generate mkparamfilefunc -funcs personalOnly
//go:generate mkdoc -build-arg -tags -build-arg version_no_check

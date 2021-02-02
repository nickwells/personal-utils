<!-- Code generated by mkbadge; DO NOT EDIT. START -->
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-green?logo=go)](https://pkg.go.dev/mod/github.com/nickwells/personal-utils)
[![Go Report Card](https://goreportcard.com/badge/github.com/nickwells/personal-utils)](https://goreportcard.com/report/github.com/nickwells/personal-utils)
![GitHub License](https://img.shields.io/github/license/nickwells/personal-utils)
<!-- Code generated by mkbadge; DO NOT EDIT. END -->
# personal-utils
this holds utility programs of personal interest not expected to be of
interest to a technical audience.

All these tools use the standard param package to handle command-line flags
and so they support the standard '-help' parameter which will print out a
comprehensive usage message.

The tools all use the version package and so the executables must be built
with the ldflags set. See the `goBuildLdflags` script
in `github.com/nickwells/version.mod/_sh`. This will set the necessary
build flags provided you are in the package directory when you call it. This
is so as to allow the git tools to find the necessary values.

```
go build -ldflags="$(goBuildLdflags)"
```


## retirement

[See here](retirement/_retirement.DOC.md)

## bankACAnalysis

[See here](bankACAnalysis/_bankACAnalysis.DOC.md)

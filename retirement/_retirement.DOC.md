<!-- Created by mkdoc DO NOT EDIT. -->

# retirement

this will simulate various scenarios for retirement allowing you to explore the
effect of changes in your portfolio, inflation etc



<!-- This file is inserted into markdown files generated by mkdoc -->
<!-- if the program being documented depends on this module       -->
<!-- ============================================================ -->
<!-- See github.com/nickwells/utilities/mkdoc                     -->
## Parameters

This uses the `param` package and so it has access to the help parameters
which give a comprehensive message describing the usage of the program and
the parameters you can give. The `-help` parameter on its own will print the
standard parameters that the program can accept but you can also give
parameters to show both more or less help, in more or less detail. Other
standard parameters allow you to explore where parameters have been set and
where they can be set. The description of the `-help` parameter is a good
place to start to explore the help available.

The intention of the `param` package is to provide complete documentation
for the program from the command line.


<!-- This file is inserted into markdown files generated by mkdoc -->
<!-- if the program being documented depends on this module       -->
<!-- ============================================================ -->
<!-- See github.com/nickwells/utilities/mkdoc                     -->
## Version information

This uses the `version` package which provides version information for the
program. It is important that the program is built with either the arguments

`-tags version_no_check`

which prevents the startup check that the version information has been set,
or else with the ldflags set; the script

`github.com/nickwells/version.mod/_sh/goBuildLdflags`

will print a string setting the ldflags appropriately. It can be passed as
the value of the `-ldflags` parameter to `go build` or `go install`.


<!-- This file is inserted into markdown files generated by mkdoc -->
<!-- if the program being documented depends on this module       -->
<!-- ============================================================ -->
<!-- See github.com/nickwells/utilities/mkdoc                     -->
## Version parameters

This offers version-querying parameters that allow the user to discover the
program version.

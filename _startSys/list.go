package main

import (
	"fmt"
	"golem/paramSetter"
	"golem/params"
	"strings"
	"sysStarter/config"
)

var listStyle = "full"
var reportSections = []string{}

func init() {
	{
		p := params.New("style",
			paramSetter.EnumSetter{
				Value: &listStyle,
				AllowedVals: paramSetter.AValMap{
					"full":     "a full description of each item",
					"expanded": "an expanded description - each item on a single line",
					"short":    "a brief description"}},
			"specify how to list the results")
		p.SetGroupName(paramGroupName)
	}
	{
		p := params.New("only-show",
			paramSetter.EnumListSetter{
				Value: &reportSections,
				AllowedVals: paramSetter.AValMap{
					"servers": "show server details",
					"classes": "show class details",
					"hosts":   "show host details"}},
			"specify which section of the report to show")
		p.SetGroupName(paramGroupName)
	}
}

func list(sc *config.SysConfig) {
	if len(reportSections) == 0 {
		reportSections = []string{"servers", "classes", "hosts"}
	}

	sectionSep := ""
	for _, section := range reportSections {
		fmt.Print(sectionSep)
		if listStyle == "full" {
			fmt.Println(section)
			fmt.Println(strings.Repeat("=", len(section)))

			sectionSep = "============================================\n\n"
		} else if listStyle == "short" {
			fmt.Println(section)
			sectionSep = "\n"
		} else {
			sectionSep = "\n"
		}

		if section == "servers" {
			printAllServers(sc)
		} else if section == "hosts" {
			printAllHosts(sc)
		} else if section == "classes" {
			printAllClasses(sc)
		}
	}
}

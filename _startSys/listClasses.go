package main

import (
	"fmt"
	"math"
	"sort"
	"sysStarter/config"
)

type byServerDetails []*config.ServerDetails

func (s byServerDetails) Len() int {
	return len(s)
}

func (s byServerDetails) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byServerDetails) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func printClass(sc *config.SysConfig, className config.ClassName) {
	svrsPrompt := "server:"

	switch listStyle {
	case "full":
		strLen := sc.MaxClassNameLen()
		strLen = int(math.Max(float64(strLen), float64(len(svrsPrompt))))

		intro := "    "
		if class, ok := sc.FindClass(className); ok {
			fmt.Printf("%s%-*.*s %s\n", intro, strLen, strLen, className, class.Desc)
			if len(class.Servers) > 0 {
				fmt.Printf("%s%*.*s", intro, strLen, strLen, svrsPrompt)
				indent2 := fmt.Sprintf("%s%*.*s", intro, strLen, strLen, "")
				indent := ""
				classSvrs := class.Servers
				sort.Sort(byServerDetails(classSvrs))
				for _, svr := range classSvrs {
					fmt.Printf(" %s%s\n", indent, svr.Name)
					indent = indent2
				}
			}
		}
		fmt.Print("\n")

	case "expanded":
		classPrompt := "class:"
		if class, ok := sc.FindClass(className); ok {
			fmt.Println(classPrompt, className, "description:", class.Desc)
			if len(class.Servers) > 0 {
				classSvrs := class.Servers
				sort.Sort(byServerDetails(classSvrs))
				for _, svr := range classSvrs {
					fmt.Println(classPrompt, className, svrsPrompt, svr.Name)
				}
			}
		}

	case "short":
		fmt.Println("    ", className)
	}
}

type byClassName []config.ClassName

func (s byClassName) Len() int {
	return len(s)
}

func (s byClassName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byClassName) Less(i, j int) bool {
	return s[i] < s[j]
}

func printAllClasses(sc *config.SysConfig) {
	classNames := make([]config.ClassName, 0, len(classes))
	for name := range classes {
		classNames = append(classNames, name)
	}
	sort.Sort(byClassName(classNames))
	for _, name := range classNames {
		printClass(sc, name)
	}
}

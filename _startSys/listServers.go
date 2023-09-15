package main

import (
	"fmt"
	"math"
	"sort"
	"sysStarter/config"
)

func printServerName(intro string, strLen int, svr config.ServerDetails, svrOk bool) {
	fmt.Printf("%s%-*.*s",
		intro, strLen, strLen, svr.Name)
	if !svrOk {
		fmt.Println(" server not found")
		return
	}

	fmt.Printf(" hosts: %6d", len(svr.Hosts))
	var totInstances int
	for _, instances := range svr.Hosts {
		totInstances += instances
	}
	fmt.Printf(" instances: %6d\n", totInstances)
}

func printServerHosts(sc *config.SysConfig,
	intro string,
	strLen int,
	svr config.ServerDetails,
	hostsPrompt string) {
	if len(svr.Hosts) > 0 {
		indent2 := printPrompt(intro, strLen, hostsPrompt)
		indent := ""
		for hostName, val := range svr.Hosts {
			fmt.Printf(" %s%-*.*s %4d instance",
				indent, sc.MaxHostNameLen(), sc.MaxHostNameLen(), hostName, val)
			if val > 1 {
				fmt.Print("s")
			}
			fmt.Print(":\n")
			indent = indent2
			if hd, ok := sc.FindHost(hostName); ok {
				if hd.Servers != nil {
					for _, cmd := range hd.Servers[svr.Name] {
						fmt.Printf("   %s  %s\n", indent, cmd.Command)
					}
				}
			}
		}
	}
}

func printPrompt(intro string, strLen int, prompt string) string {
	fmt.Printf("%s%*.*s", intro, strLen, strLen, prompt)
	indent := fmt.Sprintf("%s%*.*s", intro, strLen, strLen, "")
	return indent
}

func printServerClasses(intro string, strLen int, classes []config.ClassName, prompt string) {
	if len(classes) > 0 {
		indent2 := printPrompt(intro, strLen, prompt)
		indent := ""
		for _, className := range classes {
			fmt.Printf(" %s%s\n", indent, className)
			indent = indent2
		}
	}
}

func printServerFull(sc *config.SysConfig, intro string, strLen int, svr config.ServerDetails,
	classesPrompt, hostsPrompt, needsPrompt, neededByPrompt string) {
	printServerClasses(intro, strLen, svr.Classes, classesPrompt)

	printServerHosts(sc, intro, strLen, svr, hostsPrompt)

	if len(svr.Needs) > 0 {
		indent2 := printPrompt(intro, strLen, needsPrompt)
		indent := ""
		for needs := range svr.Needs {
			fmt.Printf(" %s%s\n", indent, needs)
			indent = indent2
		}
	}
	if len(svr.NeededBy) > 0 {
		indent2 := printPrompt(intro, strLen, neededByPrompt)
		indent := ""
		for neededBy := range svr.NeededBy {
			fmt.Printf(" %s%s\n", indent, neededBy)
			indent = indent2
		}
	}
	fmt.Print("\n")
}

func printServerExpanded(sc *config.SysConfig, svr config.ServerDetails,
	svrPrompt, classesPrompt, hostsPrompt, needsPrompt, neededByPrompt string) {

	for _, className := range svr.Classes {
		fmt.Println(svrPrompt, svr.Name, classesPrompt, className)
	}

	for hostName := range svr.Hosts {
		if hd, ok := sc.FindHost(hostName); ok {
			if hd.Servers != nil {
				for _, cmd := range hd.Servers[svr.Name] {
					fmt.Println(svrPrompt, svr.Name,
						hostsPrompt, hostName,
						"command:", cmd.Command)
				}
			}
		}
	}

	for needs := range svr.Needs {
		fmt.Println(svrPrompt, svr.Name, needsPrompt, needs)
	}

	for neededBy := range svr.NeededBy {
		fmt.Println(svrPrompt, svr.Name, neededByPrompt, neededBy)
	}

}

func printServer(sc *config.SysConfig, svrName config.ServerName) {
	classesPrompt := "class:"
	hostsPrompt := "host:"
	needsPrompt := "needs:"
	neededByPrompt := "needed-by:"

	strLen := sc.MaxServerNameLen()
	strLen = int(math.Max(float64(strLen), float64(len(classesPrompt))))
	strLen = int(math.Max(float64(strLen), float64(len(hostsPrompt))))
	strLen = int(math.Max(float64(strLen), float64(len(needsPrompt))))
	strLen = int(math.Max(float64(strLen), float64(len(neededByPrompt))))

	intro := "    "

	svr, svrOk := sc.FindServer(svrName)

	switch listStyle {
	case "full":
		printServerName(intro, strLen, svr, svrOk)
		if !svrOk {
			return
		}
		printServerFull(sc, intro, strLen, svr,
			classesPrompt, hostsPrompt, needsPrompt, neededByPrompt)

	case "expanded":
		svrPrompt := "server:"
		if !svrOk {
			fmt.Println(svrPrompt, svrName, "NOT-FOUND")
			return
		}
		printServerExpanded(sc, svr,
			svrPrompt, classesPrompt, hostsPrompt, needsPrompt, neededByPrompt)

	case "short":
		printServerName(intro, strLen, svr, svrOk)
	}
}

type byServerName []config.ServerName

func (s byServerName) Len() int {
	return len(s)
}

func (s byServerName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byServerName) Less(i, j int) bool {
	return s[i] < s[j]
}

func printAllServers(sc *config.SysConfig) {
	serverNames := make([]config.ServerName, 0, len(servers))
	for name := range servers {
		serverNames = append(serverNames, name)
	}
	sort.Sort(byServerName(serverNames))
	for _, name := range serverNames {
		printServer(sc, name)
	}
}

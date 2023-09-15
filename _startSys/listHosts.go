package main

import (
	"fmt"
	"math"
	"sort"
	"sysStarter/config"
)

func printFullHostDetails(sc *config.SysConfig,
	hostName config.HostName,
	attrsPrompt, svrsPrompt string) {
	strLen := sc.MaxHostNameLen()
	strLen = int(math.Max(float64(strLen), float64(len(attrsPrompt))))
	strLen = int(math.Max(float64(strLen), float64(len(svrsPrompt))))

	if host, ok := sc.FindHost(hostName); ok {
		intro := "    "
		fmt.Printf("%s%-*.*s", intro, strLen, strLen, hostName)
		fmt.Printf("  %-*.*s", sc.MaxDCNameLen(), sc.MaxDCNameLen(), host.Datacentre)
		fmt.Printf("  %-*.*s\n", sc.MaxOSNameLen(), sc.MaxOSNameLen(), host.OS)
		if len(host.Attrs) > 0 {
			fmt.Printf("%s%*.*s", intro, strLen, strLen, attrsPrompt)
			indent2 := fmt.Sprintf("%s%*.*s", intro, strLen, strLen, "")
			indent := ""
			for attrName := range host.Attrs {
				for _, attrVal := range host.Attrs[attrName] {
					fmt.Printf(" %s%s = %s\n", indent, attrName, attrVal)
					indent = indent2
				}
			}
		}

		if len(host.Servers) > 0 {
			fmt.Printf("%s%*.*s", intro, strLen, strLen, svrsPrompt)
			indent2 := fmt.Sprintf("%s%*.*s", intro, strLen, strLen, "")
			indent := ""
			for svrName := range host.Servers {
				fmt.Printf(" %s%s\n", indent, svrName)
				indent = indent2
				for _, cmd := range host.Servers[svrName] {
					fmt.Printf("%s    %s\n", indent, cmd.Command)
				}
			}
		}
	}
	fmt.Print("\n")
}

func printExpandedHostDetails(sc *config.SysConfig,
	hostName config.HostName,
	attrsPrompt, svrsPrompt string) {
	hostPrompt := "host:"
	if host, ok := sc.FindHost(hostName); ok {
		fmt.Println(hostPrompt, hostName, "datacentre:", host.Datacentre)
		fmt.Println(hostPrompt, hostName, "OS:", host.OS)

		for attrName := range host.Attrs {
			for _, attrVal := range host.Attrs[attrName] {
				fmt.Println(hostPrompt, hostName, attrsPrompt, attrName, "val:", attrVal)
			}
		}

		for svrName := range host.Servers {
			for _, cmd := range host.Servers[svrName] {
				fmt.Println(hostPrompt, hostName, svrsPrompt, svrName, "command:", cmd.Command)
			}
		}
	}
}

func printHost(sc *config.SysConfig, hostName config.HostName) {
	attrsPrompt := "attr:"
	svrsPrompt := "server:"

	switch listStyle {
	case "full":
		printFullHostDetails(sc, hostName, attrsPrompt, svrsPrompt)

	case "expanded":
		printExpandedHostDetails(sc, hostName, attrsPrompt, svrsPrompt)

	case "short":
		fmt.Println("    ", hostName)
	}
}

type byHostName []config.HostName

func (s byHostName) Len() int {
	return len(s)
}

func (s byHostName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byHostName) Less(i, j int) bool {
	return s[i] < s[j]
}

func printAllHosts(sc *config.SysConfig) {
	hostNames := make([]config.HostName, 0, len(hosts))
	for name := range hosts {
		hostNames = append(hostNames, name)
	}
	sort.Sort(byHostName(hostNames))
	for _, name := range hostNames {
		printHost(sc, name)
	}
}

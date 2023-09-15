package main

import (
	"fmt"
	"golem/appl"
	"golem/paramSetter"
	"golem/params"
	"os"

	"sysStarter/config"
)

var action = "start"

var hostList = []string{}
var hostExclusionList = []string{}
var serverList = []string{}
var serverExclusionList = []string{}
var classList = []string{}
var classExclusionList = []string{}
var addReqs = false
var addDeps = false

var paramGroupName = "params.SysStarter"

func init() {
	params.SetGroupDescription(paramGroupName, `SysStarter specific parameters

These alter the behaviour of the sysStarter program.`)

	{
		p := params.New("action",
			paramSetter.EnumSetter{
				Value: &action,
				AllowedVals: paramSetter.AValMap{
					"list":  "list the configuration details",
					"start": "start the system"}},
			"specify what action to take")
		p.SetGroupName(paramGroupName)
	}

	{
		p := params.New("host",
			paramSetter.StrListSetter{Value: &hostList},
			"only servers on these hosts will be used")
		p.SetGroupName(paramGroupName)
	}

	{
		p := params.New("not-hst",
			paramSetter.StrListSetter{Value: &hostExclusionList},
			"only servers not on these hosts will be used")
		p.SetGroupName(paramGroupName)
	}

	{
		p := params.New("server",
			paramSetter.StrListSetter{Value: &serverList},
			"only these servers will be used")
		p.AddAltName("svr")
		p.SetGroupName(paramGroupName)
	}

	{
		p := params.New("not-server",
			paramSetter.StrListSetter{Value: &serverExclusionList},
			"these servers will not be used")
		p.AddAltName("not-svr")
		p.SetGroupName(paramGroupName)
	}

	{
		p := params.New("class",
			paramSetter.StrListSetter{Value: &classList},
			"only servers in these classes will be used")
		p.SetGroupName(paramGroupName)
	}

	{
		p := params.New("not-class",
			paramSetter.StrListSetter{Value: &classExclusionList},
			"only servers not a member of these classes will be used")
		p.SetGroupName(paramGroupName)
	}

	{
		p := params.New("add-required-svrs",
			paramSetter.BoolSetter{Value: &addReqs},
			"any servers required by those specified will be added to the set")
		p.AddAltName("add-reqs")
		p.SetGroupName(paramGroupName)
	}

	{
		p := params.New("add-dependencies",
			paramSetter.BoolSetter{Value: &addDeps},
			"any servers that are dependent on those specified will be added to the set")
		p.AddAltName("add-deps")
		p.SetGroupName(paramGroupName)
	}
}

func reportParamErrs(paramErrs []string) {
	if errCount := len(paramErrs); errCount > 0 {
		fmt.Println("\n", params.ProgName())
		fmt.Println(errCount, " error(s) detected with the parameters:")
		sep := ""
		for _, err := range paramErrs {
			fmt.Println(sep, err)
			sep = "\n"
		}
		os.Exit(1)
	}
}

var servers = make(map[config.ServerName]bool)
var serversToExclude = make(map[config.ServerName]bool)

var hosts = make(map[config.HostName]bool)
var hostsToExclude = make(map[config.HostName]bool)

var classes = make(map[config.ClassName]bool)
var classesToExclude = make(map[config.ClassName]bool)

func populateClassesToExclude(sc *config.SysConfig) []string {
	errs := make([]string, 0, 0)

	if len(classExclusionList) > 0 {
		if len(classList) > 0 {
			errs = append(errs, "you have specified a class list and a class exclusion list - you should only supply one or the other")
		} else {
			for _, name := range classExclusionList {
				if _, ok := sc.FindClass(config.ClassName(name)); ok {
					classesToExclude[config.ClassName(name)] = true
				} else {
					errs = append(errs, "Class: "+name+" from the class exclusion list is not valid (not found)")
				}
			}
		}
	}
	return errs
}

func populateHostsToExclude(sc *config.SysConfig) []string {
	errs := make([]string, 0, 0)

	if len(hostExclusionList) > 0 {
		if len(hostList) > 0 {
			errs = append(errs, "you have specified a host list and a host exclusion list - you should only supply one or the other")
		} else {
			for _, name := range hostExclusionList {
				if _, ok := sc.FindHost(config.HostName(name)); ok {
					hostsToExclude[config.HostName(name)] = true
				} else {
					errs = append(errs, "Host: "+name+" from the host exclusion list is not valid (not found)")
				}
			}
		}
	}
	return errs
}

func populateServersToExclude(sc *config.SysConfig) []string {
	errs := make([]string, 0, 0)

	if len(serverExclusionList) > 0 {
		if len(serverList) > 0 {
			errs = append(errs, "you have specified a server list and a server exclusion list - you should only supply one or the other")
		} else {
			for _, name := range serverExclusionList {
				if _, ok := sc.FindServer(config.ServerName(name)); ok {
					serversToExclude[config.ServerName(name)] = true
				} else {
					errs = append(errs, "Server: "+name+" from the server exclusion list is not valid (not found)")
				}
			}
		}
	}
	return errs
}

func processParams(sc *config.SysConfig) {
	errs := make([]string, 0, 0)

	errs = append(errs, populateClassesToExclude(sc)...)
	errs = append(errs, populateHostsToExclude(sc)...)
	errs = append(errs, populateServersToExclude(sc)...)

	if len(classList) == 0 {
		for _, name := range sc.AllClasses() {
			if !classesToExclude[name] {
				classes[name] = true
			}
		}
	} else {
		for _, name := range classList {
			if class, ok := sc.FindClass(config.ClassName(name)); ok {
				classes[config.ClassName(name)] = true
				for _, svr := range class.Servers {
					if !serversToExclude[svr.Name] {
						servers[svr.Name] = true
					}
				}
			} else {
				errs = append(errs, "Class: "+name+" from the class list is not valid (not found)")
			}
		}
	}

	if len(hostList) == 0 {
		for _, name := range sc.AllHosts() {
			if !hostsToExclude[name] {
				hosts[name] = true
			}
		}
	} else {
		for _, name := range hostList {
			if _, ok := sc.FindHost(config.HostName(name)); ok {
				hosts[config.HostName(name)] = true
			} else {
				errs = append(errs, "Host: "+name+" from the host list is not valid (not found)")
			}
		}
	}

	if len(serverList) == 0 && len(classList) == 0 {
		for _, svr := range sc.AllServers() {
			exclude := serversToExclude[svr.Name]
			for _, className := range svr.Classes {
				if classesToExclude[className] {
					exclude = true
					break
				}
			}
			if !exclude {
				servers[svr.Name] = true
			}
		}
	}
	for _, name := range serverList {
		if _, ok := sc.FindServer(config.ServerName(name)); ok {
			servers[config.ServerName(name)] = true
		}
	}

	reportParamErrs(errs)
}

func main() {
	params.SetEnvPrefix("SYS_STARTER_")
	_, err := appl.New()
	if err != nil {
		fmt.Printf("Error starting application: %s\n", err)
		os.Exit(1)
	}

	sc := config.GetConfig()

	processParams(sc)

	if action == "start" {
		startSystem(sc)
	} else if action == "list" {
		list(sc)
	} else if action == "check" {
		checkConfig(sc)
	}
}

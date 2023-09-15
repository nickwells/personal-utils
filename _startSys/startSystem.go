package main

import (
	"fmt"
	"golem/paramSetter"
	"golem/params"
	"sort"
	"strings"
	"sysStarter/config"
	"time"

	"golang.org/x/crypto/ssh"
)

type cmdID uint

type cmdStatus uint

const (
	cmdUnknown cmdStatus = iota
	cmdOK
	cmdCouldntRun
	cmdFail

	cmdStatusMax
)

type hostStatus uint

const (
	hostUnknown hostStatus = iota
	hostOK
	hostConnFail
	hostSessionFail

	hostStatusMax
)

type identifiedCmd struct {
	id       cmdID
	hostName config.HostName
	svrName  config.ServerName
	cmd      string
}

func (idCmd identifiedCmd) String() string {
	return fmt.Sprintf("[%d] server: %s host: '%s' command: '%s'",
		idCmd.id,
		idCmd.svrName,
		idCmd.hostName,
		idCmd.cmd)
}

type identifiedCmdStatus struct {
	idCmd        identifiedCmd
	status       cmdStatus
	statusOfHost hostStatus
	desc         string
}

type hostStatusSummary struct {
	statusCount  [hostStatusMax]uint
	latestStatus hostStatus
	total        uint
}

type svrStatusSummary struct {
	statusCount  [cmdStatusMax]uint
	latestStatus cmdStatus
	total        uint
}

type sysStatus struct {
	hosts           map[config.HostName]hostStatusSummary
	servers         map[config.ServerName]svrStatusSummary
	commands        []identifiedCmdStatus
	responseCount   cmdID
	cmdFailureCount uint
}

func (ss sysStatus) String() string {
	return fmt.Sprintf("system status summary: responses: %d failures: %d hosts: %d servers: %d",
		ss.responseCount, ss.cmdFailureCount, len(ss.hosts), len(ss.servers))
}

type channelSet struct {
	cmdChan    chan identifiedCmd
	statusChan chan identifiedCmdStatus
}

var hostChannel = make(map[config.HostName]channelSet)

var startServers = true

func init() {
	p := params.New("dont-start",
		paramSetter.BoolSetterNot{Value: &startServers},
		"Don't start the servers")
	p.SetGroupName(paramGroupName)
}

func printError(msgs ...interface{}) {
	fmt.Println("*****************")
	fmt.Println("ERROR:", msgs)
	fmt.Println("*****************")
}

func allNeededServersLaunched(needs map[config.ServerName]bool, launched map[config.ServerName]bool) bool {
	for sn := range needs {
		if !launched[sn] {
			return false
		}
	}
	return true
}

func launchableServers(sc *config.SysConfig, servers []*config.ServerDetails, launched map[config.ServerName]bool) (launchableServers []*config.ServerDetails) {
	for _, svr := range servers {
		if launched[svr.Name] {
			continue
		}

		if !allNeededServersLaunched(svr.Needs, launched) {
			continue
		}

		launchableServers = append(launchableServers, svr)
	}
	return launchableServers
}

func hostLauncher(host config.HostDetails, cs channelSet) {
	var launcherStatus = hostUnknown
	var launcherStatusDetails string
	var session *ssh.Session

	portNum := "22"
	vals, ok := host.Attrs["SSHPort"]
	if ok && len(vals) > 0 {
		portNum = vals[0]
	}

	hostDetails := fmt.Sprint("Host: ", host.Name, " Port: ", portNum)
	params.Verbose("Open ssh connection to: ", hostDetails, "\n")
	conn, err := makeClientConn(string(host.Name) + ":" + portNum)
	if err != nil {
		launcherStatusDetails = fmt.Sprint("couldn't connect to: ",
			hostDetails, " Problem:", err)
		launcherStatus = hostConnFail
	} else {
		launcherStatus = hostOK
	}

	for {
		idCmd, ok := <-cs.cmdChan
		if !ok {
			break
		}

		var description string
		cmdStatus := cmdUnknown

		if launcherStatus != hostOK {
			cmdStatus = cmdCouldntRun
			description = launcherStatusDetails
		} else {
			session, err = conn.NewSession()
			if err != nil {
				launcherStatusDetails = fmt.Sprint("couldn't make the session.",
					hostDetails, "Problem:", err)
				launcherStatus = hostSessionFail
				cmdStatus = cmdCouldntRun
				description = launcherStatusDetails

				printError(launcherStatusDetails)
			} else {
				defer session.Close()
				params.Verbose("starting server on host: ", host.Name,
					" with: ", idCmd.cmd, "\n")

				if err = session.Run(idCmd.cmd + " &"); err == nil {
					cmdStatus = cmdOK
				} else {
					printError("Couldn't run command:", idCmd.cmd,
						"on", hostDetails, "Problem:", err)
					cmdStatus = cmdFail
					description = err.Error()
				}
			}
		}

		cs.statusChan <- identifiedCmdStatus{
			idCmd:        idCmd,
			status:       cmdStatus,
			statusOfHost: launcherStatus,
			desc:         description}
	}
}

func launchServer(host config.HostDetails, idCmd identifiedCmd, statusChan chan identifiedCmdStatus) {
	if params.VerboseMode() || !startServers {
		fmt.Print("\t\t\t")
		if !startServers {
			fmt.Print("would have started ")
		} else {
			fmt.Print("starting ")
		}
		fmt.Println(idCmd)
	}
	if startServers {
		cs, ok := hostChannel[host.Name]
		if !ok {
			cs.cmdChan = make(chan identifiedCmd, host.ServerCount)
			cs.statusChan = statusChan
			hostChannel[host.Name] = cs
			go hostLauncher(host, cs)
		}
		cs.cmdChan <- idCmd
	}
}

func monitorStatus(sc *config.SysConfig,
	statusChan chan identifiedCmdStatus,
	statusReqChan chan chan sysStatus) {
	var systemStatus = sysStatus{
		hosts:    make(map[config.HostName]hostStatusSummary),
		servers:  make(map[config.ServerName]svrStatusSummary),
		commands: make([]identifiedCmdStatus, sc.TotalServerCount()),
	}
	for {
		select {
		case responseChan := <-statusReqChan:
			responseChan <- systemStatus
		case cmdStatus := <-statusChan:
			systemStatus.responseCount++
			if cmdStatus.status != cmdOK {
				systemStatus.cmdFailureCount++
			}
			hName := cmdStatus.idCmd.hostName
			sName := cmdStatus.idCmd.svrName

			hs := systemStatus.hosts[hName]
			hs.latestStatus = cmdStatus.statusOfHost
			hs.statusCount[cmdStatus.statusOfHost]++
			hs.total++
			systemStatus.hosts[hName] = hs

			ss := systemStatus.servers[sName]
			ss.latestStatus = cmdStatus.status
			ss.statusCount[cmdStatus.status]++
			ss.total++
			systemStatus.servers[sName] = ss

			systemStatus.commands[cmdStatus.idCmd.id] = cmdStatus
		}
	}
}

func (ss sysStatus) reportBadHosts() {
	hostIntro := fmt.Sprintf("               host failures: ")
	altHostIntro := strings.Repeat(" ", len(hostIntro))
	var badHostCount int
	for hostName, hostStatus := range ss.hosts {
		if hostStatus.statusCount[hostOK] != hostStatus.total {
			badHostCount++

			fmt.Print(hostIntro)
			hostIntro = altHostIntro
			fmt.Print(hostName)

			if hostStatus.statusCount[hostConnFail] > 0 {
				fmt.Print(": Connection failed")
			} else if hostStatus.statusCount[hostSessionFail] > 0 {
				fmt.Print(": Session failed")
			} else {
				fmt.Print(" *** Unknown error ***")
			}
			fmt.Println()
		}
	}
	if badHostCount > 0 {
		fmt.Printf("%stotal failures: %6d\n\n",
			altHostIntro, badHostCount)
	}
}

// (ss sysStatus) maxBadSvrNameLen returns the maximum name length of any server that has failed
func (ss sysStatus) maxBadSvrNameLen() int {
	maxLen := 0
	for name, svrStatus := range ss.servers {
		failures := svrStatus.total - svrStatus.statusCount[cmdOK]
		if failures != 0 {
			if maxLen < len(name) {
				maxLen = len(name)
			}
		}
	}
	return maxLen
}

// serverKeys returns a sorted slice containing the keys to the servers map
func (ss sysStatus) serverKeys() []config.ServerName {
	serverKeys := make([]config.ServerName, 0, len(ss.servers))

	for name := range ss.servers {
		serverKeys = append(serverKeys, name)
	}

	sort.Slice(serverKeys, func(i, j int) bool {
		return serverKeys[i] < serverKeys[j]
	})
	return serverKeys
}

func (ss sysStatus) reportBadServers() {
	svrIntro := "               svr failure: "
	altSvrIntro := strings.Repeat(" ", len(svrIntro))

	var badSvrCount int

	nameLen := ss.maxBadSvrNameLen()
	blankName := strings.Repeat(" ", nameLen)

	serverKeys := ss.serverKeys()

	for _, svrName := range serverKeys {
		svrStatus := ss.servers[svrName]
		failures := svrStatus.total - svrStatus.statusCount[cmdOK]
		if failures != 0 {
			badSvrCount += int(failures)

			formattedName := fmt.Sprintf("%-*.*s", nameLen, nameLen, svrName)

			if svrStatus.statusCount[cmdCouldntRun] > 0 {
				fmt.Printf("%s%s : Couldn't run command (%6d times)",
					svrIntro, formattedName, svrStatus.statusCount[cmdCouldntRun])
				svrIntro = altSvrIntro
				formattedName = blankName
			}
			if svrStatus.statusCount[cmdFail] > 0 {
				fmt.Printf("%s%s : Command failed (%6d times)",
					svrIntro, formattedName, svrStatus.statusCount[cmdFail])
				svrIntro = altSvrIntro
				formattedName = blankName
			}
			fmt.Println()
		}
	}
	if badSvrCount > 0 {
		fmt.Printf("%stotal failures: %6d\n\n",
			altSvrIntro, badSvrCount)
	}
}

func (ss sysStatus) report() {
	fmt.Printf("System Status: responses: %d failures: %d\n",
		ss.responseCount, ss.cmdFailureCount)

	if ss.cmdFailureCount > 0 {
		ss.reportBadHosts()
		ss.reportBadServers()
	}
}

func reportStatus(responsesExpected cmdID, waitSecs uint, statusReqChan chan chan sysStatus) {
	var statusRespChan = make(chan sysStatus)
	defer close(statusRespChan)

	statusReqChan <- statusRespChan

	var attempts uint
Loop:
	for ss := range statusRespChan {
		if ss.responseCount < responsesExpected {
			if attempts > waitSecs {
				fmt.Printf("Timed Out: some responses have not been received: expected: %d\n",
					responsesExpected)
				ss.report()

				break Loop
			}

			time.Sleep(time.Second)
			statusReqChan <- statusRespChan
			attempts++
		} else {
			ss.report()
			break Loop
		}
	}
}

func startSystem(sc *config.SysConfig) {
	servers := sc.AllServers()
	launched := make(map[config.ServerName]bool)

	var id cmdID
	var pendingServers = make([]identifiedCmd, 0, sc.TotalServerCount())
	var statusChan = make(chan identifiedCmdStatus, sc.TotalServerCount())
	var statusReqChan = make(chan chan sysStatus)

	if startServers {
		go monitorStatus(sc, statusChan, statusReqChan)
	}

	var round int
	for ls := launchableServers(sc, servers, launched); len(ls) > 0; ls = launchableServers(sc, servers, launched) {
		params.Verbose("launching servers (round ", round, ")\n")
		for _, svr := range ls {
			launched[svr.Name] = true

			for hostName := range svr.Hosts {
				host, hostExists := sc.FindHost(hostName)

				if !hostExists {
					printError("'"+hostName+"'", "is not found in the set of hosts")
					continue
				}

				svrCommands, ok := host.Servers[svr.Name]
				if !ok {
					printError("'"+svr.Name+"'",
						"is not found in the set of servers running on: ",
						"'"+hostName+"'")
					continue
				}
				for _, command := range svrCommands {
					idCmd := identifiedCmd{
						id:       id,
						hostName: host.Name,
						svrName:  svr.Name,
						cmd:      command.Command}

					launchServer(host, idCmd, statusChan)
					pendingServers = append(pendingServers, idCmd)
					id++
				}
			}
		}
		if startServers {
			reportStatus(id, 30, statusReqChan)
		}

		round++
	}
}

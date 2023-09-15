package config

import (
	"database/sql"
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type configDetails struct {
	getFunc     func(*SysConfig, *sql.DB)
	goodSqlInit func(mock sqlmock.Sqlmock)
	badSqlInit  func(mock sqlmock.Sqlmock)
}

var getConfigFuncs []configDetails

func init() {
	getConfigFuncs = []configDetails{
		{
			getFunc:     (*SysConfig).getClasses,
			goodSqlInit: initDBExpectations_classes,
			badSqlInit:  initDBExpectations_badSQL_classes,
		},
		{
			getFunc:     (*SysConfig).getDatacentres,
			goodSqlInit: initDBExpectations_datacentres,
			badSqlInit:  initDBExpectations_badSQL_datacentres,
		},
		{
			getFunc:     (*SysConfig).getDependencyDetails,
			goodSqlInit: initDBExpectations_dependencyDetails,
			badSqlInit:  initDBExpectations_badSQL_dependencyDetails,
		},
		{
			getFunc:     (*SysConfig).getHostAttrs,
			goodSqlInit: initDBExpectations_hostAttrs,
			badSqlInit:  initDBExpectations_badSQL_hostAttrs,
		},
		{
			getFunc:     (*SysConfig).getHostAttrTemplates,
			goodSqlInit: initDBExpectations_hostAttrTemplates,
			badSqlInit:  initDBExpectations_badSQL_hostAttrTemplates,
		},
		{
			getFunc:     (*SysConfig).getHostServerMap,
			goodSqlInit: initDBExpectations_hostServerMap,
			badSqlInit:  initDBExpectations_badSQL_hostServerMap,
		},
		{
			getFunc:     (*SysConfig).getHosts,
			goodSqlInit: initDBExpectations_hosts,
			badSqlInit:  initDBExpectations_badSQL_hosts,
		},
		{
			getFunc:     (*SysConfig).getLocations,
			goodSqlInit: initDBExpectations_locations,
			badSqlInit:  initDBExpectations_badSQL_locations,
		},
		{
			getFunc:     (*SysConfig).getOperatingSystems,
			goodSqlInit: initDBExpectations_operatingSystems,
			badSqlInit:  initDBExpectations_badSQL_operatingSystems,
		},
		{
			getFunc:     (*SysConfig).getServerClassMap,
			goodSqlInit: initDBExpectations_serverClassMap,
			badSqlInit:  initDBExpectations_badSQL_serverClassMap,
		},
		{
			getFunc:     (*SysConfig).getServerHostMustHave,
			goodSqlInit: initDBExpectations_serverHostMustHave,
			badSqlInit:  initDBExpectations_badSQL_serverHostMustHave,
		},
		{
			getFunc:     (*SysConfig).getServers,
			goodSqlInit: initDBExpectations_servers,
			badSqlInit:  initDBExpectations_badSQL_servers,
		},
	}
}

func reportUnexpectedErrors(t *testing.T, sc *SysConfig, errCntBefore, errCntAfter, expectedErrCnt int, action string) {
	errorsFound := errCntAfter - errCntBefore
	if errorsFound != expectedErrCnt {
		t.Error(errorsFound, " errors were found while ", action, ": expected: ", expectedErrCnt)
		if errCntAfter > errCntBefore {
			for _, err := range sc.errs[errCntBefore:] {
				fmt.Println("\t", err)
			}
		}
	}
}

func reportUnexpectedValueCount(t *testing.T, valueCount, expectedValueCount int, action string) {
	if valueCount != expectedValueCount {
		t.Error("unexpected number of values when", action,
			"- expected: ", expectedValueCount,
			". There were: ", valueCount)
	}
}

func (sc *SysConfig) initLocations(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var locs = []struct {
		name LocationName
		desc string
	}{
		{LocationName("Balham"), "a suburb of London - South West"},
		{LocationName("Croydon"), "a suburb of London - South East"},
		{LocationName("Brooklyn"), "a suburb of New York - East"},
		{LocationName("Location-withLongName"), "nowhere"},
	}

	for _, l := range locs {
		sc.populateLocation(l.name, l.desc)
	}

	testName := "initialising locations"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.LocationCount(), len(locs), testName)
	if expMaxLen := len(locs[len(locs)-1].name); sc.maxLocationNameLen != expMaxLen {
		t.Error("while", testName,
			"the max name length was:", sc.maxLocationNameLen,
			"it was expected to be:", expMaxLen)
	}

	errCntBefore = sc.ErrCount()
	locCountBefore := sc.LocationCount()
	testName = "initialising locations - with empty name"
	sc.populateLocation("", "description") // empty name - expected to fail

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.LocationCount(), locCountBefore, testName)

	errCntBefore = sc.ErrCount()
	locCountBefore = sc.LocationCount()
	testName = "initialising locations - with duplicate"
	sc.populateLocation("location1", "description")
	sc.populateLocation("location1", "description") // duplicate name - expected to fail

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.LocationCount(), locCountBefore+1, testName)
}

func (sc *SysConfig) initOperatingSystems(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var os = []OSName{
		OSName("Linux"),
		OSName("OSX"),
		OSName("Android"),
		OSName("OS-withLongName"),
	}

	for _, o := range os {
		sc.populateOperatingSystem(o)
	}

	testName := "initialising operating systems"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.OperatingSystemCount(), len(os), testName)
	if expMaxLen := len(os[len(os)-1]); sc.maxOSNameLen != expMaxLen {
		t.Error("while", testName,
			"the max name length was:", sc.maxOSNameLen,
			"it was expected to be:", expMaxLen)
	}

	errCntBefore = sc.ErrCount()
	osCountBefore := sc.OperatingSystemCount()
	testName = "initialising operating systems - with duplicate"
	sc.populateOperatingSystem("PrimeOS")
	sc.populateOperatingSystem("PrimeOS") // duplicate name - expected to fail

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.OperatingSystemCount(), osCountBefore+1, testName)

	errCntBefore = sc.ErrCount()
	osCountBefore = sc.OperatingSystemCount()
	testName = "initialising operating systems - with blank name"
	sc.populateOperatingSystem("") // blank name - expected to fail

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.OperatingSystemCount(), osCountBefore, testName)
}

func (sc *SysConfig) initDatacentres(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var dcs = []struct {
		name DCName
		loc  LocationName
	}{
		{name: DCName("DC-UK1"), loc: LocationName("Balham")},
		{name: DCName("DC-UK2"), loc: LocationName("Balham")},
		{name: DCName("DC-UK3"), loc: LocationName("Balham")},
		{name: DCName("DC-UK4"), loc: LocationName("Croydon")},
		{name: DCName("DC-NY1"), loc: LocationName("Brooklyn")},
		{name: DCName("DC-withLongName"), loc: LocationName("Brooklyn")},
	}

	for _, dc := range dcs {
		sc.populateDatacentre(dc.name, dc.loc)
	}

	testName := "initialising datacentres"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.DatacentreCount(), len(dcs), testName)
	if expMaxLen := len(dcs[len(dcs)-1].name); sc.maxDCNameLen != expMaxLen {
		t.Error("while", testName,
			"the max name length was:", sc.maxDCNameLen,
			"it was expected to be:", expMaxLen)
	}

	errCntBefore = sc.ErrCount()
	dcCntBefore := sc.DatacentreCount()
	testName = "initialising datacentres (with duplicate)"
	sc.populateDatacentre("DC-Test1", "Balham")
	sc.populateDatacentre("DC-Test1", "Balham") // duplicate - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.DatacentreCount(), dcCntBefore+1, testName)

	errCntBefore = sc.ErrCount()
	dcCntBefore = sc.DatacentreCount()
	testName = "initialising datacentres (with invalid location)"
	sc.populateDatacentre("DC-Test2", "Nonesuch") // invalid location - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.DatacentreCount(), dcCntBefore, testName)

	errCntBefore = sc.ErrCount()
	dcCntBefore = sc.DatacentreCount()
	testName = "initialising datacentres (with empty name)"
	sc.populateDatacentre("", "Balham") // empty name - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.DatacentreCount(), dcCntBefore, testName)
}

func (sc *SysConfig) initHosts(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var hosts = []struct {
		name HostName
		dc   DCName
		os   OSName
	}{
		{name: HostName("ldnprd1"), dc: DCName("DC-UK1"), os: OSName("Linux")},
		{name: HostName("ldnprd2"), dc: DCName("DC-UK1"), os: OSName("Linux")},
		{name: HostName("ldnprd3"), dc: DCName("DC-UK1"), os: OSName("Linux")},
		{name: HostName("ldnprd-withLongName"), dc: DCName("DC-UK1"), os: OSName("Linux")},
	}

	for _, h := range hosts {
		sc.populateHost(h.name, h.dc, h.os)
	}

	testName := "initialising hosts"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.HostCount(), len(hosts), testName)
	if expMaxLen := len(hosts[len(hosts)-1].name); sc.maxHostNameLen != expMaxLen {
		t.Error("while", testName,
			"the max name length was:", sc.maxHostNameLen,
			"it was expected to be:", expMaxLen)
	}

	errCntBefore = sc.ErrCount()
	hostCntBefore := sc.HostCount()
	testName = "initialising hosts (with duplicate)"
	sc.populateHost("ldnprd11", "DC-Test1", "Linux")
	sc.populateHost("ldnprd11", "DC-Test1", "Linux") // duplicate - expected to fail

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostCount(), hostCntBefore+1, testName)

	errCntBefore = sc.ErrCount()
	hostCntBefore = sc.HostCount()
	testName = "initialising hosts (with blank name)"
	sc.populateHost("", "DC-Test1", "Linux")

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostCount(), hostCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hostCntBefore = sc.HostCount()
	testName = "initialising hosts (with invalid datacentre)"
	sc.populateHost("ldnprd12", "Nonesuch", "Linux")

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostCount(), hostCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hostCntBefore = sc.HostCount()
	testName = "initialising hosts (with invalid OS)"
	sc.populateHost("ldnprd12", "DC-Test1", "Nonesuch")

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostCount(), hostCntBefore, testName)
}

func (sc *SysConfig) initServers(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var servers = []struct {
		name      ServerName
		ms        int
		keepAlive bool
	}{
		{name: ServerName("TestSvr1"), ms: 0, keepAlive: true},
		{name: ServerName("TestSvr2"), ms: 1, keepAlive: true},
		{name: ServerName("TestSvr3"), ms: 2, keepAlive: true},
		{name: ServerName("TestSvrWithLongName"), ms: 2, keepAlive: true},
	}

	for _, s := range servers {
		sc.populateServer(s.name, s.ms, s.keepAlive)
	}
	testName := "initialising servers"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.ServerCount(), len(servers), testName)
	if expMaxLen := len(servers[len(servers)-1].name); sc.maxServerNameLen != expMaxLen {
		t.Error("while", testName,
			"the max name length was:", sc.maxServerNameLen,
			"it was expected to be:", expMaxLen)
	}

	errCntBefore = sc.ErrCount()
	serverCntBefore := sc.ServerCount()
	testName = "initialising servers (with duplicate)"
	sc.populateServer("TestSvr4", 0, false)
	sc.populateServer("TestSvr4", 1, true) // duplicate - expected to fail

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.ServerCount(), serverCntBefore+1, testName)

	errCntBefore = sc.ErrCount()
	serverCntBefore = sc.ServerCount()
	testName = "initialising servers (with blank name)"
	sc.populateServer("", 1, true) // blank name - expected to fail

	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.ServerCount(), serverCntBefore, testName)
}

func (sc SysConfig) countHostServerCommands(hostName HostName) int {
	var count int
	if host, ok := sc.hosts[hostName]; ok {
		for _, cmd := range host.Servers {
			count += len(cmd)
		}
	}
	return count
}

func (sc SysConfig) countHostServers() int {
	var count int
	for _, h := range sc.hosts {
		count += len(h.Servers)
	}
	return count
}

func (sc SysConfig) countServerHosts() int {
	var count int
	for _, s := range sc.servers {
		for _, h := range s.Hosts {
			count += h
		}
	}
	return count
}

func (sc *SysConfig) initHostServerMap(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var hs = []struct {
		hostName HostName
		svrName  ServerName
		cmd      string
	}{
		{hostName: HostName("ldnprd1"), svrName: ServerName("TestSvr1"), cmd: "runTestSvr1"},
		{hostName: HostName("ldnprd1"), svrName: ServerName("TestSvr1"), cmd: "runTestSvr1 backup"},
		{hostName: HostName("ldnprd1"), svrName: ServerName("TestSvr2"), cmd: "runTestSvr2"},
		{hostName: HostName("ldnprd2"), svrName: ServerName("TestSvr1"), cmd: "runTestSvr1"},
	}

	for _, hsc := range hs {
		sc.populateHostServerMap(hsc.hostName, hsc.svrName, hsc.cmd)
	}
	testName := "initialising host-server map"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.countHostServers(), 3, testName+" (host/server)")
	reportUnexpectedValueCount(t, sc.countServerHosts(), len(hs), testName+" (server/host)")
	reportUnexpectedValueCount(t, sc.countHostServerCommands("ldnprd1"), 3, testName+" (ServerCount - ldnprd1)")
	reportUnexpectedValueCount(t, sc.countHostServerCommands("ldnprd2"), 1, testName+" (ServerCount - ldnprd2)")
	reportUnexpectedValueCount(t, sc.countHostServerCommands("NoneSuch"), 0, testName+" (ServerCount - NoneSuch)")
	reportUnexpectedValueCount(t, int(sc.totalServerCount), 4, testName+" (ServerCount - total)")

	hostSvrCntBefore := sc.countHostServers()
	svrHostCntBefore := sc.countServerHosts()

	errCntBefore = sc.ErrCount()
	testName = "initialising host-server map (with empty command)"
	sc.populateHostServerMap("ldnprd2", "TestSvr1", "") // empty command - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)

	errCntBefore = sc.ErrCount()
	testName = "initialising host-server map (with unknown host and server)"
	sc.populateHostServerMap("Nonesuch", "Nonesuch", "command") // unknown host and server - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)

	errCntBefore = sc.ErrCount()
	testName = "initialising host-server map (with unknown server)"
	sc.populateHostServerMap("ldnprd1", "Nonesuch", "command") // unknown server - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)

	errCntBefore = sc.ErrCount()
	testName = "initialising host-server map (with unknown host)"
	sc.populateHostServerMap("Nonesuch", "TestSvr1", "command") // unknown host - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)

	testName = "initialising host-server map with errors"
	reportUnexpectedValueCount(t, sc.countHostServers(), hostSvrCntBefore, testName+" (host/server)")
	reportUnexpectedValueCount(t, sc.countServerHosts(), svrHostCntBefore, testName+" (server/host)")
}

func (sc *SysConfig) initHostAttrTemplates(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var hats = []struct {
		name                           AttrName
		hasDflt                        bool
		dfltVal                        string
		minCnt, maxCnt                 int
		desc, units, valRegex, valDesc string
	}{
		{name: AttrName("attr0"), hasDflt: false, dfltVal: "", minCnt: 2, maxCnt: -1, desc: "zero or more", units: "", valRegex: ".*", valDesc: "any string, at least 2"},
		{name: AttrName("attr1"), hasDflt: false, dfltVal: "", minCnt: 0, maxCnt: -1, desc: "zero or more", units: "", valRegex: ".*", valDesc: "any string"},
		{name: AttrName("attr2"), hasDflt: false, dfltVal: "", minCnt: 0, maxCnt: -1, desc: "zero or more", units: "", valRegex: `^\d+$`, valDesc: "any number"},
		{name: AttrName("attr3"), hasDflt: false, dfltVal: "", minCnt: 1, maxCnt: 1, desc: "mandatory and no default", units: "", valRegex: "^ABC.*$", valDesc: "any string starting with ABC"},
		{name: AttrName("attr4"), hasDflt: true, dfltVal: "dflt3", minCnt: 1, maxCnt: 1, desc: "mandatory with default", units: "", valRegex: ".*", valDesc: "any string"},
		{name: AttrName("attr5"), hasDflt: true, dfltVal: "dflt5", minCnt: 2, maxCnt: -1, desc: "zero or more", units: "", valRegex: ".*", valDesc: "any string, at least 2"},
		{name: AttrName("attr6"), hasDflt: true, dfltVal: "dflt6", minCnt: 1, maxCnt: 1, desc: "zero or more", units: "", valRegex: ".*", valDesc: "any string, at least 2"},
	}

	for _, hat := range hats {
		sc.populateHostAttrTemplate(hat.name,
			hat.hasDflt, hat.dfltVal,
			hat.minCnt, hat.maxCnt,
			hat.desc, hat.units,
			hat.valRegex, hat.valDesc)
	}

	testName := "initialising host attribute templates"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.HostAttrTemplateCount(), len(hats), testName)

	errCntBefore = sc.ErrCount()
	hatCntBefore := sc.HostAttrTemplateCount()
	testName = "initialising host attribute templates (with duplicate)"
	sc.populateHostAttrTemplate("attr99", false, "", 0, -1, "", "", ".*", "")
	sc.populateHostAttrTemplate("attr99", true, "99", 0, -1, "", "", ".*", "") // duplicate - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostAttrTemplateCount(), hatCntBefore+1, testName)

	errCntBefore = sc.ErrCount()
	hatCntBefore = sc.HostAttrTemplateCount()
	testName = "initialising host attribute templates (with blank name)"
	sc.populateHostAttrTemplate("", false, "", 0, -1, "", "", ".*", "") // blank name - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostAttrTemplateCount(), hatCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hatCntBefore = sc.HostAttrTemplateCount()
	testName = "initialising host attribute templates (with bad regex)"
	sc.populateHostAttrTemplate("attr-bad", false, "", 0, -1, "", "", "*", "") // bad regex - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostAttrTemplateCount(), hatCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hatCntBefore = sc.HostAttrTemplateCount()
	testName = "initialising host attribute templates (with bad default)"
	sc.populateHostAttrTemplate("attr-bad", true, "", 0, -1, "", "", `\d+`, "") // bad default - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostAttrTemplateCount(), hatCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hatCntBefore = sc.HostAttrTemplateCount()
	testName = "initialising host attribute templates (with min>max)"
	sc.populateHostAttrTemplate("attr-bad", true, "", 3, 2, "", "", `.*`, "") // min>man - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostAttrTemplateCount(), hatCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hatCntBefore = sc.HostAttrTemplateCount()
	testName = "initialising host attribute templates (with min==0 and has dflt)"
	sc.populateHostAttrTemplate("attr-bad", true, "kk", 0, -1, "", "", `.*`, "") // min==0 && has dflt - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.HostAttrTemplateCount(), hatCntBefore, testName)
}

func (sc SysConfig) countHostAttrs() int {
	var count int
	for _, h := range sc.hosts {
		count += len(h.Attrs)
	}
	return count
}

func (sc *SysConfig) initHostAttrs(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var has = []struct {
		host HostName
		attr AttrName
		val  string
	}{
		{host: HostName("ldnprd1"), attr: AttrName("attr0"), val: "val0.0"},
		{host: HostName("ldnprd1"), attr: AttrName("attr0"), val: "val0.1"},
		{host: HostName("ldnprd1"), attr: AttrName("attr1"), val: "val1"},
		{host: HostName("ldnprd1"), attr: AttrName("attr2"), val: "99"},
		{host: HostName("ldnprd1"), attr: AttrName("attr3"), val: "ABC99"},
		{host: HostName("ldnprd1"), attr: AttrName("attr5"), val: "val5"},
		{host: HostName("ldnprd1"), attr: AttrName("attr6"), val: "val6"},
		{host: HostName("ldnprd2"), attr: AttrName("attr0"), val: "val0.0"},
		{host: HostName("ldnprd2"), attr: AttrName("attr5"), val: "val5"},
		{host: HostName("ldnprd2"), attr: AttrName("attr6"), val: "val6"},
		{host: HostName("ldnprd3"), attr: AttrName("attr0"), val: "val0.0"},
		{host: HostName("ldnprd3"), attr: AttrName("attr0"), val: "val0.1"},
		{host: HostName("ldnprd3"), attr: AttrName("attr1"), val: "val1"},
		{host: HostName("ldnprd3"), attr: AttrName("attr2"), val: "99"},
		{host: HostName("ldnprd3"), attr: AttrName("attr3"), val: "ABC99"},
		{host: HostName("ldnprd3"), attr: AttrName("attr5"), val: "val5"},
		{host: HostName("ldnprd3"), attr: AttrName("attr6"), val: "val6"},
		{host: HostName("ldnprd3"), attr: AttrName("attr6"), val: "val6"},
	}

	for _, ha := range has {
		sc.populateHostAttr(ha.host, ha.attr, ha.val)
	}

	testName := "initialising host attributes"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.countHostAttrs(), len(has)-3, testName)

	errCntBefore = sc.ErrCount()
	hostAttrCntBefore := sc.countHostAttrs()
	testName = "initialising host attributes (invalid host)"
	sc.populateHostAttr("Nonesuch", "attr1", "any") // bad host - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countHostAttrs(), hostAttrCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hostAttrCntBefore = sc.countHostAttrs()
	testName = "initialising host attributes (invalid attr)"
	sc.populateHostAttr("ldnprd1", "Nonesuch", "any") // bad attr - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countHostAttrs(), hostAttrCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hostAttrCntBefore = sc.countHostAttrs()
	testName = "initialising host attributes (invalid value)"
	sc.populateHostAttr("ldnprd1", "attr2", "not a number") // bad value - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countHostAttrs(), hostAttrCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hostAttrCntBefore = sc.countHostAttrs()
	testName = "checking host attributes (good host all mandatory attrs present)"
	sc.checkAttrs(sc.hosts["ldnprd1"])
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.countHostAttrs(), hostAttrCntBefore+1, testName)

	errCntBefore = sc.ErrCount()
	hostAttrCntBefore = sc.countHostAttrs()
	testName = "checking host attributes (non-existant host)"
	sc.checkAttrs(sc.hosts["Nonesuch"])
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.countHostAttrs(), hostAttrCntBefore, testName)

	errCntBefore = sc.ErrCount()
	hostAttrCntBefore = sc.countHostAttrs()
	testName = "checking host attributes (valid host missing mandatory attrs)"
	sc.checkAttrs(sc.hosts["ldnprd2"])
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 2, testName)
	reportUnexpectedValueCount(t, sc.countHostAttrs(), hostAttrCntBefore+1, testName)

	errCntBefore = sc.ErrCount()
	hostAttrCntBefore = sc.countHostAttrs()
	testName = "checking host attributes (valid host too many attrs)"
	sc.checkAttrs(sc.hosts["ldnprd3"])
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countHostAttrs(), hostAttrCntBefore+1, testName)
}

func (sc *SysConfig) initClasses(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var classes = []struct {
		name ClassName
		desc string
	}{
		{name: ClassName("Class1"), desc: "class desc"},
		{name: ClassName("Class2"), desc: "class desc"},
		{name: ClassName("Class3"), desc: "class desc"},
		{name: ClassName("Class-withLongName"), desc: "class desc"},
	}

	for _, c := range classes {
		sc.populateClasses(c.name, c.desc)
	}

	testName := "initialising classes"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.ClassCount(), len(classes), testName)
	if expMaxLen := len(classes[len(classes)-1].name); sc.maxClassNameLen != expMaxLen {
		t.Error("while", testName,
			"the max name length was:", sc.maxClassNameLen,
			"it was expected to be:", expMaxLen)
	}

	errCntBefore = sc.ErrCount()
	classCntBefore := sc.ClassCount()
	testName = "initialising classes (empty class name)"
	sc.populateClasses("", "desc") // empty name - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.ClassCount(), classCntBefore, testName)

	errCntBefore = sc.ErrCount()
	classCntBefore = sc.ClassCount()
	testName = "initialising classes (duplicates)"
	sc.populateClasses("Class99", "desc")
	sc.populateClasses("Class99", "desc") // duplicate - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.ClassCount(), classCntBefore+1, testName)
}

func (sc SysConfig) countClassesByServer() int {
	var count int
	for _, s := range sc.servers {
		count += len(s.Classes)
	}
	return count
}

func (sc SysConfig) countServersByClass() int {
	var count int
	for _, c := range sc.classes {
		count += len(c.Servers)
	}
	return count
}

func (sc *SysConfig) initServerClassMap(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var svrClassMappings = []struct {
		svrName   ServerName
		className ClassName
	}{
		{svrName: ServerName("TestSvr1"), className: ClassName("Class1")},
		{svrName: ServerName("TestSvr2"), className: ClassName("Class1")},
		{svrName: ServerName("TestSvr1"), className: ClassName("Class2")},
		{svrName: ServerName("TestSvr1"), className: ClassName("Class3")},
	}

	for _, scm := range svrClassMappings {
		sc.populateServerClassMap(scm.svrName, scm.className)
	}

	testName := "initialising server/class map"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.countClassesByServer(), len(svrClassMappings), testName+" (classes/server)")
	reportUnexpectedValueCount(t, sc.countServersByClass(), len(svrClassMappings), testName+" (server/class)")

	errCntBefore = sc.ErrCount()
	cbsBefore := sc.countClassesByServer()
	sbcBefore := sc.countServersByClass()
	testName = "initialising server/class map with unknown server"
	sc.populateServerClassMap("Nonesuch", "Class1") // bad server - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countClassesByServer(), cbsBefore, testName+" (classes/server)")
	reportUnexpectedValueCount(t, sc.countServersByClass(), sbcBefore, testName+" (server/class)")

	errCntBefore = sc.ErrCount()
	cbsBefore = sc.countClassesByServer()
	sbcBefore = sc.countServersByClass()
	testName = "initialising server/class map with unknown class"
	sc.populateServerClassMap("TestSvr1", "Nonesuch") // bad class - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countClassesByServer(), cbsBefore, testName+" (classes/server)")
	reportUnexpectedValueCount(t, sc.countServersByClass(), sbcBefore, testName+" (server/class)")

	errCntBefore = sc.ErrCount()
	cbsBefore = sc.countClassesByServer()
	sbcBefore = sc.countServersByClass()
	testName = "initialising server/class map with unknown class and server"
	sc.populateServerClassMap("Nonesuch", "Nonesuch") // bad class and server - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countClassesByServer(), cbsBefore, testName+" (classes/server)")
	reportUnexpectedValueCount(t, sc.countServersByClass(), sbcBefore, testName+" (server/class)")
}

func (sc SysConfig) countSvrHostMustHave() int {
	var count int
	for _, svr := range sc.servers {
		count += len(svr.HostMustHave)
	}
	return count
}

func (sc *SysConfig) initServerHostMustHave(t *testing.T) {
	errCntBefore := sc.ErrCount()
	var svrHostMustHave = []struct {
		svrName  ServerName
		attrName AttrName
		val      string
	}{
		{svrName: ServerName("TestSvr1"), attrName: AttrName("attr1"), val: "val1"},
		{svrName: ServerName("TestSvr1"), attrName: AttrName("attr2"), val: "99"},
		{svrName: ServerName("TestSvr1"), attrName: AttrName("attr3"), val: "ABC99"},
		{svrName: ServerName("TestSvr2"), attrName: AttrName("attr6"), val: "ABC99"},
	}

	for _, scm := range svrHostMustHave {
		sc.populateServerHostMustHave(scm.svrName, scm.attrName, scm.val)
	}

	testName := "initialising server host-must-have values"
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.countSvrHostMustHave(), len(svrHostMustHave), testName)

	errCntBefore = sc.ErrCount()
	shmhCountBefore := sc.countSvrHostMustHave()
	testName = "initialising server host-must-have values - bad server"
	sc.populateServerHostMustHave("Nonesuch", "attr2", "99") // bad server - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countSvrHostMustHave(), shmhCountBefore, testName)

	errCntBefore = sc.ErrCount()
	shmhCountBefore = sc.countSvrHostMustHave()
	testName = "initialising server host-must-have values - bad attr"
	sc.populateServerHostMustHave("TestSvr2", "Nonesuch", "99") // bad attr - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countSvrHostMustHave(), shmhCountBefore, testName)

	errCntBefore = sc.ErrCount()
	shmhCountBefore = sc.countSvrHostMustHave()
	testName = "initialising server host-must-have values - bad server and attr"
	sc.populateServerHostMustHave("Nonesuch", "Nonesuch", "99") // bad server and attr - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countSvrHostMustHave(), shmhCountBefore, testName)

	errCntBefore = sc.ErrCount()
	shmhCountBefore = sc.countSvrHostMustHave()
	testName = "initialising server host-must-have values - bad value"
	sc.populateServerHostMustHave("TestSvr2", "attr2", "a99") // bad value - expected to fail
	sc.populateServerHostMustHave("TestSvr2", "attr2", "99a") // bad value - expected to fail
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 2, testName)
	reportUnexpectedValueCount(t, sc.countSvrHostMustHave(), shmhCountBefore, testName)

	errCntBefore = sc.ErrCount()
	testName = "confirm host meets server needs"
	sc.confirmHostsMeetSvrNeeds()
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 4, testName)
}

func (sc *SysConfig) countSvrNeeds() int {
	var count int
	for _, svr := range sc.servers {
		count += len(svr.Needs)
	}
	return count
}

func (sc *SysConfig) countSvrNeededBy() int {
	var count int
	for _, svr := range sc.servers {
		count += len(svr.NeededBy)
	}
	return count
}

func (sc *SysConfig) initDependencyDetails(t *testing.T) {
	errCntBefore := sc.ErrCount()
	testName := "dependencyDetails - constructServerSet with invalid type name"
	if _, ok := sc.constructServerSet("Nonesuch", "value", "expected invalid type name"); ok {
		t.Error(testName, "- didn't fail")
	} else {
		reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	}

	errCntBefore = sc.ErrCount()
	testName = "dependencyDetails - constructServerSet with invalid class name"
	if _, ok := sc.constructServerSet("class", "Nonesuch", "expected invalid class name"); ok {
		t.Error(testName, "- didn't fail")
	} else {
		reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	}

	errCntBefore = sc.ErrCount()
	testName = "dependencyDetails - constructServerSet with invalid server name"
	if _, ok := sc.constructServerSet("server", "Nonesuch", "expected invalid server name"); ok {
		t.Error(testName, "- didn't fail")
	} else {
		reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	}

	errCntBefore = sc.ErrCount()
	testName = "dependencyDetails - constructServerSet - server"
	svrSet, ok := sc.constructServerSet("server", "TestSvr1", "expected valid")
	if !ok {
		t.Error(testName, "- unexpected failure")
	}
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, len(svrSet), 1, testName)

	errCntBefore = sc.ErrCount()
	testName = "dependencyDetails - constructServerSet class"
	svrSet, ok = sc.constructServerSet("class", "Class1", "expected valid")
	if !ok {
		t.Error(testName, "- unexpected failure")
	}
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, len(svrSet), 2, testName)

	errCntBefore = sc.ErrCount()
	testName = "dependencyDetails - initialise"
	snCount := sc.countSvrNeeds()
	sc.populateServerSet("class", "Class1", "server", "TestSvr3")
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.countSvrNeeds(), sc.countSvrNeededBy(), testName+" Needs should == NeededBy")
	reportUnexpectedValueCount(t, sc.countSvrNeeds(), snCount+2, testName+" dependency count")

	errCntBefore = sc.ErrCount()
	testName = "dependencyDetails - initialise (self dependency)"
	snCount = sc.countSvrNeeds()
	sc.populateServerSet("server", "TestSvr3", "server", "TestSvr3")
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)
	reportUnexpectedValueCount(t, sc.countSvrNeededBy(), snCount, testName+" dependency count")

	errCntBefore = sc.ErrCount()
	testName = "dependencyDetails - detect cycles - none present"
	sc.detectCycles()
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 0, testName)

	errCntBefore = sc.ErrCount()
	testName = "dependencyDetails - detect cycles - 1 present"
	sc.populateServerSet("server", "TestSvr1", "server", "TestSvr2")
	sc.populateServerSet("server", "TestSvr2", "server", "TestSvr1")
	sc.detectCycles()
	reportUnexpectedErrors(t, sc, errCntBefore, sc.ErrCount(), 1, testName)
	reportUnexpectedValueCount(t, sc.countSvrNeeds(), sc.countSvrNeededBy(), testName+" Needs should == NeededBy")
}

func (sc *SysConfig) populate(t *testing.T) {
	sc.initLocations(t)
	sc.initOperatingSystems(t)
	sc.initDatacentres(t)
	sc.initHosts(t)
	sc.initServers(t)
	sc.initHostServerMap(t)
	sc.initHostAttrTemplates(t)
	sc.initHostAttrs(t)
	sc.initClasses(t)
	sc.initServerClassMap(t)
	sc.initServerHostMustHave(t)
	sc.initDependencyDetails(t)
}

func TestGetConfig(t *testing.T) {
	abortOnError = false
	sc := NewSysConfig()
	sc.populate(t)
}

func initDBExpectations_locations(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"loc_name", "loc_desc"})
	rs.AddRow("Balham", "London - South West")
	mock.ExpectQuery(locationsQuery).WillReturnRows(rs)
}

func initDBExpectations_operatingSystems(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"os_name"})
	rs.AddRow("Linux")
	mock.ExpectQuery(operatingSystemsQuery).WillReturnRows(rs)
}
func initDBExpectations_datacentres(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"dc_name", "loc_name"})
	rs.AddRow("LDN1", "Balham")
	mock.ExpectQuery(datacentresQuery).WillReturnRows(rs)
}
func initDBExpectations_hosts(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"host_name", "dc_name", "os_name"})
	rs.AddRow("ldnprd1", "LDN1", "Linux")
	mock.ExpectQuery(hostsQuery).WillReturnRows(rs)
}
func initDBExpectations_servers(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"svr_name", "startup_millisecs", "keep_alive"})
	rs.AddRow("TestSvr1", 10, true)
	mock.ExpectQuery(serversQuery).WillReturnRows(rs)
}
func initDBExpectations_hostServerMap(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"host_name", "svr_name", "command"})
	rs.AddRow("ldnprd1", "TestSvr1", "run_TestSvr1")
	mock.ExpectQuery(hostServerMapQuery).WillReturnRows(rs)
}
func initDBExpectations_hostAttrTemplates(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"host_attr",
		"default_attr_val",
		"attr_min_count",
		"attr_max_count",
		"attr_description",
		"attr_units",
		"value_regex",
		"value_desc",
	})
	rs.AddRow("attr1", nil, 0, 1, "description", "", ".*", "")
	mock.ExpectQuery(hostAttrTemplatesQuery).WillReturnRows(rs)
}
func initDBExpectations_hostAttrs(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"host_name", "host_attr", "attr_val"})
	rs.AddRow("ldnprd1", "attr1", "val1")
	mock.ExpectQuery(hostAttrsQuery).WillReturnRows(rs)
}
func initDBExpectations_classes(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"class_name", "class_desc"})
	rs.AddRow("Class1", "class desc")
	rs.AddRow("Class2", "class desc")
	mock.ExpectQuery(classesQuery).WillReturnRows(rs)
}
func initDBExpectations_serverClassMap(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"server_name", "class_name"})
	rs.AddRow("TestSvr1", "Class1")
	mock.ExpectQuery(serverClassMapQuery).WillReturnRows(rs)
}
func initDBExpectations_serverHostMustHave(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"svr_name", "host_attr", "attr_val"})
	rs.AddRow("TestSvr1", "attr1", "val1")
	mock.ExpectQuery(serverHostMustHaveQuery).WillReturnRows(rs)
}
func initDBExpectations_dependencyDetails(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"from_type", "from_val", "to_type", "to_val"})
	rs.AddRow("server", "TestSvr1", "class", "Class2")
	mock.ExpectQuery(dependencyDetailsQuery).WillReturnRows(rs)
}

func initDBExpectations(mock sqlmock.Sqlmock) {
	initDBExpectations_locations(mock)
	initDBExpectations_operatingSystems(mock)
	initDBExpectations_datacentres(mock)
	initDBExpectations_hosts(mock)
	initDBExpectations_servers(mock)
	initDBExpectations_hostServerMap(mock)
	initDBExpectations_hostAttrTemplates(mock)
	initDBExpectations_hostAttrs(mock)
	initDBExpectations_classes(mock)
	initDBExpectations_serverClassMap(mock)
	initDBExpectations_serverHostMustHave(mock)
	initDBExpectations_dependencyDetails(mock)
}

func initDBExpectations_badSQL_locations(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"loc_name"})
	rs.AddRow("too few")
	rs.RowError(0, fmt.Errorf("locations - row error"))
	mock.ExpectQuery(locationsQuery).WillReturnRows(rs)
}

func initDBExpectations_badSQL_operatingSystems(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"os_name"})
	rs.RowError(0, fmt.Errorf("operatingSystems - row error"))
	mock.ExpectQuery(operatingSystemsQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_datacentres(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"dc_name", "loc_name"})
	rs.RowError(0, fmt.Errorf("datacentres - row error"))
	mock.ExpectQuery(datacentresQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_hosts(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"host_name", "dc_name", "os_name"})
	rs.RowError(0, fmt.Errorf("hosts - row error"))
	mock.ExpectQuery(hostsQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_servers(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"svr_name", "startup_millisecs", "keep_alive"})
	rs.RowError(0, fmt.Errorf("servers - row error"))
	mock.ExpectQuery(serversQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_hostServerMap(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"host_name", "svr_name", "command"})
	rs.RowError(0, fmt.Errorf("hostServerMap - row error"))
	mock.ExpectQuery(hostServerMapQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_hostAttrTemplates(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"host_attr",
		"default_attr_val",
		"attr_min_count",
		"attr_max_count",
		"attr_description",
		"attr_units",
		"value_regex",
		"value_desc",
	})
	rs.RowError(0, fmt.Errorf("hostAttrTemplates - row error"))
	mock.ExpectQuery(hostAttrTemplatesQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_hostAttrs(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"host_name", "host_attr", "attr_val"})
	rs.RowError(0, fmt.Errorf("hostAttrs - row error"))
	mock.ExpectQuery(hostAttrsQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_classes(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"class_name", "class_desc"})
	rs.RowError(0, fmt.Errorf("classes - row error"))
	mock.ExpectQuery(classesQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_serverClassMap(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"server_name", "class_name"})
	rs.RowError(0, fmt.Errorf("serverClassMap - row error"))
	mock.ExpectQuery(serverClassMapQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_serverHostMustHave(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"svr_name", "host_attr", "attr_val"})
	rs.RowError(0, fmt.Errorf("serverHostMustHave - row error"))
	mock.ExpectQuery(serverHostMustHaveQuery).WillReturnRows(rs)
}
func initDBExpectations_badSQL_dependencyDetails(mock sqlmock.Sqlmock) {
	rs := sqlmock.NewRows([]string{"from_type", "from_val", "to_type", "to_val"})
	rs.RowError(0, fmt.Errorf("dependencyDetails - row error"))
	mock.ExpectQuery(dependencyDetailsQuery).WillReturnRows(rs)
}

func (sc *SysConfig) testPopulateFromDB(t *testing.T, db *sql.DB, mock sqlmock.Sqlmock) {
	initDBExpectations(mock)
	sc.getFromDB(db)
	if ecnt := sc.ErrCount(); ecnt != 0 {
		sc.ReportErrors()
		t.Error("there were ", ecnt, "errors when constructing sysStarter/config from the db")
	}
}

func TestDBGetConfig(t *testing.T) {
	abortOnError = false
	sc := NewSysConfig()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("Couldn't open mock db: ", err.Error())
		return
	}
	defer db.Close()
	sc.testPopulateFromDB(t, db, mock)
}

func functionName(f func(*SysConfig, *sql.DB)) string {
	funcName := "Unknown"
	if rtFunc := runtime.FuncForPC(reflect.ValueOf(f).Pointer()); rtFunc != nil {
		funcName = rtFunc.Name()
	}
	return funcName
}

func TestDBNoSQL(t *testing.T) {
	abortOnError = false

	db, _, err := sqlmock.New()
	if err != nil {
		t.Error("While testing with no SQL - couldn't open mock db: ", err.Error())
		return
	}
	defer db.Close()

	for _, f := range getConfigFuncs {
		sc := NewSysConfig()
		f.getFunc(sc, db)
		if sc.ErrCount() != 1 && !sc.fatal {
			t.Error("calling ",
				functionName(f.getFunc),
				" function with no sql didn't cause a fatal error")
		}
	}
}

// func TestDBBadSQL(t *testing.T) {
// 	abortOnError = false

// 	for _, f := range getConfigFuncs {
// 		db, mock, err := sqlmock.New()
// 		if err != nil {
// 			t.Error("While testing with Bad SQL - couldn't open mock db: ", err.Error())
// 			return
// 		}
// 		defer db.Close()
// 		f.badSqlInit(mock)
// 		sc := NewSysConfig()
// 		f.getFunc(sc, db)
// 		if sc.ErrCount() != 1 && !sc.fatal {
// 			t.Error("calling ",
// 				functionName(f.getFunc),
// 				" function with bad sql didn't cause a fatal error")
// 		}
// 	}
// }

func TestFatalFlagStopsEarly(t *testing.T) {
	abortOnError = false

	db, _, err := sqlmock.New()
	if err != nil {
		t.Error("Couldn't open mock db: ", err.Error())
		return
	}
	defer db.Close()

	for _, f := range getConfigFuncs {
		sc := NewSysConfig()
		sc.fatal = true
		f.getFunc(sc, db)
		if sc.ErrCount() > 0 {
			t.Error("calling ",
				functionName(f.getFunc),
				" with the fatal flag set should have caused an early return but the routine was entered")
		}
	}
}

func TestSysConfigErrors(t *testing.T) {
	sc := NewSysConfig()
	sc.addError("test")
	if sc.ErrCount() != 1 {
		t.Error("calling addError didn't add an error")
	}

	if sc.fatal {
		t.Error("the fatal flag was set before any fatal errors were added")
	}
	errCntBefore := sc.ErrCount()
	sc.addFatalError("test")
	if sc.ErrCount() != errCntBefore+1 {
		t.Error("calling addFatalError didn't add an error")
	}
	if !sc.fatal {
		t.Error("calling addFatalError didn't set the fatal flag")
	}
}

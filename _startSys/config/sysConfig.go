package config

import "regexp"

type DCName string

type DatacentreDetails struct {
	Name  DCName
	Loc   LocationName
	Hosts []*HostDetails
}

type LocationName string

type LocationDetails struct {
	Name        LocationName
	Description string
	DCs         []*DatacentreDetails
}

type AttrName string

type HostAttrTemplate struct {
	Name                           AttrName
	HasDefault                     bool
	DefaultVal                     string
	Min, Max                       int
	Desc, Units, ValRegex, ValDesc string
	Re                             *regexp.Regexp
}

type HostName string
type AttrValMap map[AttrName][]string

type ServerCommandDetails struct {
	Name    ServerName
	Command string
}

type HostDetails struct {
	Name        HostName
	Datacentre  DCName
	OS          OSName
	Attrs       AttrValMap
	Servers     map[ServerName][]ServerCommandDetails
	ServerCount uint
}

type ClassName string
type ClassDetails struct {
	Name    ClassName
	Desc    string
	Servers []*ServerDetails
}

type ServerName string

type ServerDetails struct {
	Name         ServerName
	StartupMsecs int
	KeepAlive    bool
	Classes      []ClassName
	HostMustHave AttrValMap
	Hosts        map[HostName]int
	Needs        map[ServerName]bool
	NeededBy     map[ServerName]bool
}

type OSName string

type SysConfig struct {
	errs  []string
	fatal bool

	operatingSystems  map[OSName]bool
	datacentres       map[DCName]*DatacentreDetails
	locations         map[LocationName]*LocationDetails
	hostAttrTemplates map[AttrName]HostAttrTemplate
	hosts             map[HostName]*HostDetails
	classes           map[ClassName]*ClassDetails
	servers           map[ServerName]*ServerDetails

	maxServerNameLen,
	maxHostNameLen,
	maxClassNameLen,
	maxOSNameLen,
	maxDCNameLen,
	maxLocationNameLen int

	totalServerCount uint
}

func (sc *SysConfig) TotalServerCount() uint { return sc.totalServerCount }

func (sc *SysConfig) FindServer(name ServerName) (svr ServerDetails, ok bool) {
	sp, ok := sc.servers[name]
	if ok {
		return *sp, ok
	} else {
		return ServerDetails{}, ok
	}
}

func (sc *SysConfig) AllServerNames() []ServerName {
	rval := make([]ServerName, 0, len(sc.servers))
	for name, _ := range sc.servers {
		rval = append(rval, name)
	}
	return rval
}

func (sc *SysConfig) AllServers() []*ServerDetails {
	rval := make([]*ServerDetails, 0, len(sc.servers))
	for _, svr := range sc.servers {
		rval = append(rval, svr)
	}
	return rval
}

func (sc *SysConfig) FindHost(name HostName) (host HostDetails, ok bool) {
	hp, ok := sc.hosts[name]
	if ok {
		return *hp, ok
	} else {
		return HostDetails{}, ok
	}
}

func (sc *SysConfig) AllHosts() []HostName {
	rval := make([]HostName, 0, len(sc.hosts))
	for name, _ := range sc.hosts {
		rval = append(rval, name)
	}
	return rval
}

func (sc *SysConfig) FindClass(name ClassName) (svr ClassDetails, ok bool) {
	sp, ok := sc.classes[name]
	if ok {
		return *sp, ok
	} else {
		return ClassDetails{}, ok
	}
}

func (sc *SysConfig) AllClasses() []ClassName {
	rval := make([]ClassName, 0, len(sc.classes))
	for name, _ := range sc.classes {
		rval = append(rval, name)
	}
	return rval
}

func (sc *SysConfig) HostCount() int             { return len(sc.hosts) }
func (sc *SysConfig) ClassCount() int            { return len(sc.classes) }
func (sc *SysConfig) ServerCount() int           { return len(sc.servers) }
func (sc *SysConfig) LocationCount() int         { return len(sc.locations) }
func (sc *SysConfig) DatacentreCount() int       { return len(sc.datacentres) }
func (sc *SysConfig) OperatingSystemCount() int  { return len(sc.operatingSystems) }
func (sc *SysConfig) HostAttrTemplateCount() int { return len(sc.hostAttrTemplates) }
func (sc *SysConfig) ErrCount() int              { return len(sc.errs) }

func (sc *SysConfig) MaxServerNameLen() int   { return sc.maxServerNameLen }
func (sc *SysConfig) MaxHostNameLen() int     { return sc.maxHostNameLen }
func (sc *SysConfig) MaxClassNameLen() int    { return sc.maxClassNameLen }
func (sc *SysConfig) MaxOSNameLen() int       { return sc.maxOSNameLen }
func (sc *SysConfig) MaxDCNameLen() int       { return sc.maxDCNameLen }
func (sc *SysConfig) MaxLocationNameLen() int { return sc.maxLocationNameLen }

func NewSysConfig() *SysConfig {
	return &SysConfig{
		operatingSystems:  make(map[OSName]bool),
		datacentres:       make(map[DCName]*DatacentreDetails),
		locations:         make(map[LocationName]*LocationDetails),
		hostAttrTemplates: make(map[AttrName]HostAttrTemplate),
		hosts:             make(map[HostName]*HostDetails),
		classes:           make(map[ClassName]*ClassDetails),
		servers:           make(map[ServerName]*ServerDetails)}
}

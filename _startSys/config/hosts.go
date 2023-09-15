package config

import (
	"database/sql"
)

func (sc *SysConfig) populateHost(name HostName, dc DCName, os OSName) {
	if name == "" {
		sc.addError("Host names must not be blank")
		return
	}
	if host, hostExists := sc.hosts[name]; hostExists {
		sc.addError("Host: '" + string(name) +
			"' is already in the set of hosts (with datacentre: '" +
			string(host.Datacentre) +
			"' and OS: '" +
			string(host.OS) +
			"'). Host names must be unique")
		return
	}

	dcd, dcExists := sc.datacentres[dc]
	if !dcExists {
		sc.addError("Host: '" + string(name) +
			"' has an invalid datacentre: '" + string(dc) +
			"': it is not found in the set of Datacentres")
		return
	}

	if _, ok := sc.operatingSystems[os]; !ok {
		errStr := "Host: '" + string(name) +
			"' has an invalid OS: '" + string(os) +
			"': it is not found in the set of operating systems"
		sc.addError(errStr)
		return
	}

	if nameLen := len(name); nameLen > sc.maxHostNameLen {
		sc.maxHostNameLen = nameLen
	}
	details := &HostDetails{
		Name:       name,
		Datacentre: dc,
		OS:         os,
		Attrs:      make(AttrValMap)}

	sc.hosts[name] = details
	dcd.Hosts = append(dcd.Hosts, details)
}

const hostsQuery = "select host_name, dc_name, os_name from hosts"

func (sc *SysConfig) getHosts(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(hostsQuery)
	if err != nil {
		sc.addFatalError("getHosts failed: SQL error: " + err.Error())
	} else {
		defer rows.Close()
		var (
			name HostName
			dc   DCName
			os   OSName
		)
		for rows.Next() {
			err := rows.Scan((*string)(&name), (*string)(&dc), (*string)(&os))
			if err != nil {
				sc.addFatalError(
					"getHosts failed: couldn't scan the columns: " +
						err.Error())
				break
			}
			sc.populateHost(name, dc, os)
		}
	}
}

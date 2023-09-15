package config

import "database/sql"

func (sc *SysConfig) populateHostServerMap(
	hostName HostName,
	svrName ServerName,
	command string) {
	host, hostOk := sc.hosts[hostName]
	svr, svrOk := sc.servers[svrName]

	if command == "" {
		sc.addError("host-to-server mapping error: host: " + string(hostName) +
			" running server: " + string(svrName) +
			" (with command: " + command +
			") - the command must not be null")
		return
	}

	if !hostOk && !svrOk {
		sc.addError("host-to-server mapping error: host: " + string(hostName) +
			" running server: " + string(svrName) +
			" (with command: " + command +
			") - both host and server are unknown")
		return
	}

	if !hostOk {
		sc.addError("host-to-server mapping error: host: " + string(hostName) +
			" running server: " + string(svrName) +
			" (with command: " + command +
			") - the host is unknown")
		return
	}

	if !svrOk {
		sc.addError("host-to-server mapping error: host: " + string(hostName) +
			" running server: " + string(svrName) +
			" (with command: " + command +
			") - the server is unknown")
		return
	}

	if host.Servers == nil {
		host.Servers = make(map[ServerName][]ServerCommandDetails)
	}
	host.Servers[svrName] = append(host.Servers[svrName],
		ServerCommandDetails{
			Name:    svrName,
			Command: command})
	host.ServerCount++
	svr.Hosts[hostName]++
	sc.totalServerCount++
}

const hostServerMapQuery = "select host_name, svr_name, command from host_server_map"

func (sc *SysConfig) getHostServerMap(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(hostServerMapQuery)
	if err != nil {
		sc.addFatalError("getHostServerMap failed: SQL error: " + err.Error())
	} else {
		defer rows.Close()
		var (
			hostName HostName
			svrName  ServerName
			command  string
		)
		for rows.Next() {
			err := rows.Scan(
				(*string)(&hostName),
				(*string)(&svrName),
				&command)
			if err != nil {
				sc.addFatalError(
					"getHostServerMap failed: couldn't scan the columns: " +
						err.Error())
				break
			}
			sc.populateHostServerMap(hostName, svrName, command)
		}
	}
}

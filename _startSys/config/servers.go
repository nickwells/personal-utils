package config

import "database/sql"

func (sc *SysConfig) populateServer(name ServerName, startupMsecs int, keepAlive bool) {
	if name == "" {
		sc.addError("Server names must not be blank")
		return
	}

	if _, svrExists := sc.servers[name]; svrExists {
		sc.addError("server: " + string(name) +
			" is already in the set of servers." +
			" Server names must be unique")
		return
	}

	if nameLen := len(name); nameLen > sc.maxServerNameLen {
		sc.maxServerNameLen = nameLen
	}
	sc.servers[name] = &ServerDetails{
		Name:         name,
		StartupMsecs: startupMsecs,
		KeepAlive:    keepAlive,
		HostMustHave: make(AttrValMap),
		Hosts:        make(map[HostName]int),
		Needs:        make(map[ServerName]bool),
		NeededBy:     make(map[ServerName]bool)}
}

const serversQuery = "select svr_name, startup_millisecs, keep_alive from servers"

func (sc *SysConfig) getServers(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(serversQuery)
	if err != nil {
		sc.addFatalError("getServers failed: SQL error: " + err.Error())
	} else {
		defer rows.Close()
		var (
			name         ServerName
			startupMsecs int
			keepAlive    bool
		)
		for rows.Next() {
			err := rows.Scan(
				(*string)(&name),
				(*int)(&startupMsecs),
				&keepAlive)
			if err != nil {
				sc.addFatalError(
					"getServers failed: couldn't scan the columns: " +
						err.Error())
				break
			}

			sc.populateServer(name, startupMsecs, keepAlive)
		}
	}
}

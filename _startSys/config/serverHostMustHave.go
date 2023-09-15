package config

import "database/sql"

func (sc *SysConfig) populateServerHostMustHave(svrName ServerName, attr AttrName, val string) {
	svr, svrOk := sc.servers[svrName]
	tmplt, attrOk := sc.hostAttrTemplates[attr]

	if !svrOk && !attrOk {
		sc.addError("serverHostMustHave error: invalid server: '" + string(svrName) +
			"' and attr: '" + string(attr) +
			"' ( = '" + val + "')")
		return
	}

	if !svrOk {
		sc.addError("serverHostMustHave error: server: '" + string(svrName) +
			"' with attr: '" + string(attr) +
			"' = '" + val +
			"' - the server name is invalid")
		return
	}

	if !attrOk {
		sc.addError("serverHostMustHave error: server: '" + string(svrName) +
			"' with attr: '" + string(attr) +
			"' = '" + val +
			"' - the attr name is invalid")
		return
	}

	if !tmplt.Re.MatchString(val) {
		sc.addError("serverHostMustHave error: server: '" + string(svrName) +
			"' has an invalid attr: '" + string(attr) +
			"': value: '" + val +
			"' doesn't match regexp: " + tmplt.ValRegex)
		return
	}

	svr.HostMustHave[attr] = append(svr.HostMustHave[attr], val)
}

func (sc *SysConfig) confirmHostsMeetSvrNeeds() {
	for _, svr := range sc.servers {
		for hostName, _ := range svr.Hosts {
			host := sc.hosts[hostName]
			var problemCnt int
			for attrName, svrAttrVals := range svr.HostMustHave {
				if hostAttrVals, ok := host.Attrs[attrName]; ok {
					for _, svrVal := range svrAttrVals {
						found := false
						for _, hostVal := range hostAttrVals {
							if hostVal == svrVal {
								found = true
								break
							}
							if !found {
								problemCnt++
								err := "host '" + string(hostName) +
									"' doesn't satisfy server '" + string(svr.Name) +
									"' - Host attribute: '" + string(attrName) +
									"' should be: '" + svrVal +
									"', is: "
								sep := ""
								for _, hostVal := range hostAttrVals {
									err = err + sep + "'" + hostVal + "'"
									sep = ", "
								}
								sc.addError(err)
							}
						}
					}
				} else {
					problemCnt++
					sc.addError("host '" + string(hostName) +
						"' doesn't satisfy server '" + string(svr.Name) +
						"' - '" + string(attrName) +
						"' is not defined for the host")
				}
			}
		}
	}
}

const serverHostMustHaveQuery = "select svr_name, host_attr, attr_val from server_host_must_have;"

func (sc *SysConfig) getServerHostMustHave(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(serverHostMustHaveQuery)
	if err != nil {
		sc.addFatalError(
			"getServerHostMustHave failed: SQL error: " + err.Error())
	} else {
		defer rows.Close()
		var (
			svrName ServerName
			attr    AttrName
			val     string
		)
		for rows.Next() {
			err := rows.Scan((*string)(&svrName), (*string)(&attr), &val)
			if err != nil {
				sc.addFatalError(
					"getServerHostMustHave failed: couldn't scan the columns: " +
						err.Error())
				break
			}
			sc.populateServerHostMustHave(svrName, attr, val)
		}
		sc.confirmHostsMeetSvrNeeds()
	}
}

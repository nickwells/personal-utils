package config

import (
	"database/sql"
)

func (sc *SysConfig) populateOperatingSystem(name OSName) {
	if name == "" {
		sc.addError("Operating system names must not be blank")
		return
	}

	if _, osExists := sc.operatingSystems[name]; osExists {
		sc.addError("Operating system: " + string(name) +
			" is already in the set of operating systems." +
			" Operating system names must be unique")
		return
	}

	if nameLen := len(name); nameLen > sc.maxOSNameLen {
		sc.maxOSNameLen = nameLen
	}
	sc.operatingSystems[name] = true
}

const operatingSystemsQuery = "select os_name from operating_systems"

func (sc *SysConfig) getOperatingSystems(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(operatingSystemsQuery)
	if err != nil {
		sc.addFatalError("getOperatingSystems failed: SQL error: " +
			err.Error())
	} else {
		defer rows.Close()
		var name OSName
		for rows.Next() {
			err := rows.Scan((*string)(&name))
			if err != nil {
				sc.addFatalError(
					"getOperatingSystems failed: couldn't scan the columns: " +
						err.Error())
				break
			}
			sc.populateOperatingSystem(name)
		}
	}
}

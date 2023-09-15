package config

import (
	"database/sql"
)

func (sc *SysConfig) populateLocation(name LocationName, desc string) {
	if name == "" {
		sc.addError("Location names must not be blank")
		return
	}
	if _, locExists := sc.locations[name]; locExists {
		sc.addError("location: " + string(name) +
			" is already in the set of locations." +
			" Location names must be unique")
		return
	}

	if nameLen := len(name); nameLen > sc.maxLocationNameLen {
		sc.maxLocationNameLen = nameLen
	}
	sc.locations[name] = &LocationDetails{
		Name:        name,
		Description: desc}
}

const locationsQuery = "select loc_name, loc_desc from locations"

func (sc *SysConfig) getLocations(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(locationsQuery)
	if err != nil {
		sc.addFatalError("getLocations failed: SQL error: " +
			err.Error())
	} else {
		defer rows.Close()
		var (
			name LocationName
			desc string
		)
		for rows.Next() {
			err := rows.Scan((*string)(&name), &desc)
			if err != nil {
				sc.addFatalError(
					"getLocations failed: couldn't scan the columns: " +
						err.Error())
				break
			}
			sc.populateLocation(name, desc)
		}
	}
}

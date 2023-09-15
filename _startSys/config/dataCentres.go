package config

import "database/sql"

func (sc *SysConfig) populateDatacentre(name DCName, loc LocationName) {
	dcDetails := &DatacentreDetails{
		Name: name,
		Loc:  loc}

	if name == "" {
		sc.addError("Datacentre names must not be blank")
		return
	}

	if dc, dcExists := sc.datacentres[name]; dcExists {
		sc.addError("datacentre: " +
			string(name) +
			" is already in the set of datacentres (with location: " +
			string(dc.Loc) +
			"). Datacentre names must be unique")
		return
	}

	locDetails, locExists := sc.locations[loc]
	if !locExists {
		errStr := "datacentre: " + string(name) +
			" has an invalid location: '" + string(loc) +
			"': it is not found in the set of Locations"
		sc.addError(errStr)
		return
	}

	if nameLen := len(name); nameLen > sc.maxDCNameLen {
		sc.maxDCNameLen = nameLen
	}
	sc.datacentres[name] = dcDetails
	locDetails.DCs = append(locDetails.DCs, dcDetails)
}

const datacentresQuery = "select dc_name, loc_name from datacentres"

func (sc *SysConfig) getDatacentres(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(datacentresQuery)
	if err != nil {
		sc.addFatalError("getDatacentres failed: SQL error: " + err.Error())
	} else {
		defer rows.Close()
		var (
			name DCName
			loc  LocationName
		)
		for rows.Next() {
			err := rows.Scan((*string)(&name), (*string)(&loc))
			if err != nil {
				sc.addFatalError(
					"getDatacentres failed: couldn't scan the columns: " +
						err.Error())
				break
			}
			sc.populateDatacentre(name, loc)
		}
	}
}

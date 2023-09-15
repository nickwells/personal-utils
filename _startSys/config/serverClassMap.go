package config

import "database/sql"

func (sc *SysConfig) populateServerClassMap(svrName ServerName, className ClassName) {
	svr, svrExists := sc.servers[svrName]
	class, classExists := sc.classes[className]

	if !svrExists && !classExists {
		sc.addError("serverClassMap error: neither class: '" +
			string(className) +
			"' nor server: '" + string(svrName) + "' exist")
		return
	}

	if !svrExists {
		sc.addError("serverClassMap error: server: '" + string(svrName) +
			"' doesn't exist")
		return
	}

	if !classExists {
		sc.addError("serverClassMap error: class: '" + string(className) +
			"' doesn't exist")
		return
	}

	svr.Classes = append(svr.Classes, className)
	class.Servers = append(class.Servers, svr)
}

const serverClassMapQuery = "select svr_name, class_name from server_class_map"

func (sc *SysConfig) getServerClassMap(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(serverClassMapQuery)
	if err != nil {
		sc.addFatalError("getServerClassMap failed: SQL error: " +
			err.Error())
	} else {
		defer rows.Close()
		var (
			svrName   ServerName
			className ClassName
		)
		for rows.Next() {
			err := rows.Scan((*string)(&svrName), (*string)(&className))
			if err != nil {
				sc.addFatalError(
					"getServerClassMap failed: couldn't scan the columns: " +
						err.Error())
				break
			}

			sc.populateServerClassMap(svrName, className)
		}
	}
}

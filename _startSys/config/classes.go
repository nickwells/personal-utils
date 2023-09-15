package config

import (
	"database/sql"
)

func (sc *SysConfig) populateClasses(name ClassName, desc string) {
	if name == "" {
		sc.addError("Class names must not be blank")
		return
	}
	if c, classExists := sc.classes[name]; classExists {
		sc.addError("Class: '" + string(name) +
			"' is already in the set of classes (with description: '" +
			c.Desc +
			"'). Class names must be unique")
		return
	}

	if nameLen := len(name); nameLen > sc.maxClassNameLen {
		sc.maxClassNameLen = nameLen
	}
	sc.classes[name] = &ClassDetails{
		Name: name,
		Desc: desc}
}

const classesQuery = "select class_name, class_desc from classes"

func (sc *SysConfig) getClasses(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(classesQuery)
	if err != nil {
		sc.addFatalError("getClasses failed: SQL error: " +
			err.Error())
	} else {
		defer rows.Close()
		var (
			name ClassName
			desc string
		)
		for rows.Next() {
			err := rows.Scan((*string)(&name), &desc)
			if err != nil {
				sc.addFatalError(
					"getClasses failed: couldn't scan the columns: " +
						err.Error())
				break
			}
			sc.populateClasses(name, desc)
		}
	}
}

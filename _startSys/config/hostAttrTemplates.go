package config

import (
	"database/sql"
	"fmt"
	"regexp"
)

func (sc *SysConfig) populateHostAttrTemplate(name AttrName,
	hasDflt bool,
	dfltVal string,
	minCnt, maxCnt int,
	desc, units, valRegex, valDesc string) {
	if _, attrExists := sc.hostAttrTemplates[name]; attrExists {
		sc.addError("HostAttrTemplate: attr: '" + string(name) +
			"' is already in the set of host attribute templates." +
			" Attribute names must be unique")
		return
	}

	if name == "" {
		sc.addError("HostAttrTemplate: the name must not be blank")
		return
	}

	if maxCnt > 0 && minCnt > maxCnt {
		sc.addError(
			fmt.Sprintf("HostAttrTemplate: attr: '%s' has invalid min (%d) and max (%d) counts: if a max count is >0 it must be >= min count",
				name, minCnt, maxCnt))
		return
	}

	re, err := regexp.Compile(valRegex)
	if err != nil {
		sc.addError("HostAttrTemplate: attr: '" + string(name) +
			"' - the regular expression: '" + valRegex +
			"' did not compile: " + err.Error())
		return
	}

	if hasDflt {
		if !re.MatchString(dfltVal) {
			sc.addError("HostAttrTemplate: attr: '" + string(name) +
				"' - the default value: '" + dfltVal +
				"' does not match the regular expression: '" + valRegex + "'")
			return
		}
	}

	if minCnt < 0 {
		minCnt = 0
	}

	if minCnt == 0 && hasDflt {
		sc.addError("HostAttrTemplate: attr: '" + string(name) +
			"' - is optional and has a default value: '" + dfltVal +
			"'")
		return
	}

	sc.hostAttrTemplates[name] = HostAttrTemplate{
		Name:       name,
		HasDefault: hasDflt,
		DefaultVal: dfltVal,
		Min:        minCnt,
		Max:        maxCnt,
		Desc:       desc,
		Units:      units,
		ValRegex:   valRegex,
		ValDesc:    valDesc,
		Re:         re}
}

const hostAttrTemplatesQuery = "select host_attr, default_attr_val, attr_min_count, attr_max_count, attr_description, attr_units, value_regex, value_desc from host_attr_templates"

func (sc *SysConfig) getHostAttrTemplates(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(hostAttrTemplatesQuery)
	if err != nil {
		sc.addFatalError("getHostAttrTemplates failed: SQL error: " +
			err.Error())
	} else {
		defer rows.Close()
		var (
			name                           AttrName
			dfltValNullable                sql.NullString
			dfltVal                        string
			minCnt, maxCnt                 int
			desc, units, valRegex, valDesc string
		)
		for rows.Next() {
			err := rows.Scan((*string)(&name),
				&dfltValNullable,
				&minCnt,
				&maxCnt,
				&desc,
				&units,
				&valRegex,
				&valDesc)
			if err != nil {
				sc.addFatalError(
					"getHostAttrTemplates failed: couldn't scan the columns: " +
						err.Error())
				break
			}

			hasDflt := dfltValNullable.Valid
			if hasDflt {
				dfltVal = dfltValNullable.String
			} else {
				dfltVal = ""
			}
			sc.populateHostAttrTemplate(name, hasDflt, dfltVal, minCnt, maxCnt, desc, units, valRegex, valDesc)

		}
	}
}

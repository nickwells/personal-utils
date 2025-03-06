package config

import (
	"database/sql"
	"fmt"
)

func (sc *SysConfig) checkAttrs(hd *HostDetails) {
	if hd == nil {
		return
	}

	hostAttrs := hd.Attrs
	hasAttr := make(map[AttrName]bool)
	for attr, vals := range hostAttrs {
		if tmplt, ok := sc.hostAttrTemplates[attr]; !ok {
			sc.addError("hostAttr: host: '" + string(hd.Name) +
				"' attr: '" + string(attr) +
				"' - the attr name is not in the set of attribute templates")
		} else {
			valCount := len(vals)
			hasAttr[attr] = true
			if valCount < tmplt.Min {
				if tmplt.HasDefault {
					for i := valCount; i < tmplt.Min; i++ {
						hostAttrs[attr] =
							append(hostAttrs[attr], tmplt.DefaultVal)
					}
				} else {
					sc.addError(
						fmt.Sprintf(
							"hostAttr: host: '%s' attr: '%s': has too few values (%d) and no default, there should be at least: %d",
							hd.Name, attr, valCount, tmplt.Min))
				}
			} else if tmplt.Max > 0 && valCount > tmplt.Max {
				sc.addError(fmt.Sprintf(
					"hostAttr: host: '%s' attr: '%s': has too many values: %d, should be at most: %d",
					hd.Name, attr, valCount, tmplt.Max))
			}
		}
	}
	for attr, tmplt := range sc.hostAttrTemplates {
		if tmplt.Min > 0 && !hasAttr[attr] {
			if tmplt.HasDefault {
				for range tmplt.Min {
					hostAttrs[attr] = append(hostAttrs[attr], tmplt.DefaultVal)
				}
			} else {
				sc.addError(
					fmt.Sprintf("hostAttr: host: '%s' attr: '%s' - the attr is mandatory, missing and has no default value",
						hd.Name, attr))
			}
		}
	}
}

func (sc *SysConfig) populateHostAttr(host HostName, attr AttrName, val string) {
	hd, ok := sc.hosts[host]
	if !ok {
		sc.addError("hostAttr: host: '" + string(host) +
			"' with attr: '" + string(attr) +
			"' and val: '" + val +
			"' is invalid: the host is not found in the set of hosts")
		return
	}

	tmplt, ok := sc.hostAttrTemplates[attr]
	if !ok {
		sc.addError("hostAttr: attr: '" + string(attr) +
			"' with val: '" + val +
			"' being added to host: '" + string(host) +
			"' is invalid: the attribute is not found in the set of attribute templates")
		return
	}

	if !tmplt.Re.MatchString(val) {
		sc.addError("hostAttr: value: '" + val +
			"' for attr: '" + string(attr) +
			"' being added to host: '" + string(host) +
			"' doesn't match regexp: '" + tmplt.ValRegex + "'")
		return
	}

	hd.Attrs[attr] = append(hd.Attrs[attr], val)
}

const hostAttrsQuery = "select host_name, host_attr, attr_val from host_attrs"

func (sc *SysConfig) getHostAttrs(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(hostAttrsQuery)
	if err != nil {
		sc.addFatalError("getHostAttrs failed: SQL error: " + err.Error())
	} else {
		defer rows.Close()
		var (
			host HostName
			attr AttrName
			val  string
		)
		for rows.Next() {
			err := rows.Scan((*string)(&host), (*string)(&attr), &val)
			if err != nil {
				sc.addFatalError(
					"getHostAttrs failed: couldn't scan the columns: " +
						err.Error())
				break
			}
			sc.populateHostAttr(host, attr, val)
		}
		for _, hostDetails := range sc.hosts {
			sc.checkAttrs(hostDetails)
		}
	}
}

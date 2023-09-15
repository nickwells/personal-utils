package config

import (
	"database/sql"
	"fmt"
	"sort"
)

type serverDependencies struct {
	Name     ServerName
	Needs    map[ServerName]bool
	NeededBy map[ServerName]bool
}
type sdMap map[ServerName]serverDependencies

// copyDependencies will create a copy of the dependency attributes for all
// the servers
func (sc *SysConfig) copyDependencies() sdMap {
	m := make(sdMap)

	for _, s := range sc.servers {
		m[s.Name] = serverDependencies{
			Name:     s.Name,
			Needs:    make(map[ServerName]bool, len(s.Needs)),
			NeededBy: make(map[ServerName]bool, len(s.NeededBy)),
		}
		for k, v := range s.Needs {
			m[s.Name].Needs[k] = v
		}
		for k, v := range s.NeededBy {
			m[s.Name].NeededBy[k] = v
		}
	}
	return m
}

type ByName []ServerName

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i] < a[j] }

func getFirstSvr(svrs map[ServerName]bool) ServerName {
	var svrList []ServerName
	for sn, _ := range svrs {
		svrList = append(svrList, sn)
	}
	sort.Sort(ByName(svrList))
	return svrList[0]
}

// detectCycles will find and report any cycles in the dependency graph
//
// It does this by repeatedly pruning the graph. It removes any entries that
// either need no other entries or else are needed by no other entries as
// they can never be part of a cycle. Whatever is left after all the pruning
// has been done must be part of a cycle
func (sc *SysConfig) detectCycles() {
	m := sc.copyDependencies()

	for {
		svrsToRemove := make(map[ServerName]bool)
		for _, svr := range m {
			if len(svr.Needs) == 0 {
				svrsToRemove[svr.Name] = true
			} else if len(svr.NeededBy) == 0 {
				svrsToRemove[svr.Name] = true
			}
		}
		if len(svrsToRemove) == 0 {
			break
		}
		for svrName := range svrsToRemove {
			for neededBy := range m[svrName].NeededBy {
				delete(m[neededBy].Needs, svrName)
			}
			for needs := range m[svrName].Needs {
				delete(m[needs].NeededBy, svrName)
			}
			delete(m, svrName)
		}
	}
	if len(m) != 0 {
		sc.reportLoops(m)
	}
}

func (sc *SysConfig) reportLoops(m sdMap) {
	var start ServerName
	for sn, _ := range m {
		start = sn
		break
	}
	loopMsg := string(start)
	sn := getFirstSvr(m[start].Needs)
	for {
		loopMsg += " -> " + string(sn)
		if sn == start {
			break
		}
		sn = getFirstSvr(m[sn].Needs)
	}
	sc.addError("A dependency loop has been detected: " + loopMsg)
}

func (sc *SysConfig) constructServerSet(typeName, val, recDesc string) (map[ServerName]bool, bool) {
	svrSet := map[ServerName]bool{}
	var ok bool

	switch typeName {
	case "class":
		var cd *ClassDetails
		if cd, ok = sc.classes[ClassName(val)]; !ok {
			sc.addError("dependency-detail error: unknown class: '" + val +
				"' (" + recDesc + ")")
		} else {
			for _, svr := range cd.Servers {
				svrSet[svr.Name] = true
			}
		}
	case "server":
		if _, ok = sc.servers[ServerName(val)]; !ok {
			sc.addError("dependency-detail error: unknown server: '" + val +
				"' (" + recDesc + ")")
		} else {
			svrSet[ServerName(val)] = true
		}
	default:
		sc.addError("dependency-detail error: invalid type name: '" + typeName +
			"' (" + recDesc + ")")
	}

	return svrSet, ok
}

// printSvrSet will take a server set and if it is ok it will print out the
// names of the constituent servers
func printSvrSet(name string, ss map[ServerName]bool, ok bool) {
	fmt.Printf("%s: ok? %v:\n", name, ok)
	if ok {
		for svrName, _ := range ss {
			fmt.Printf("\t%v\n", svrName)
		}
	}
}

func (sc *SysConfig) populateServerSet(fromType, fromVal, toType, toVal string) {
	recDesc := fromType + "[" + fromVal + "] needs " + toType + "[" + toVal + "]"

	fromSvrSet, fromSetOk := sc.constructServerSet(fromType, fromVal, recDesc)
	toSvrSet, toSetOk := sc.constructServerSet(toType, toVal, recDesc)

	if fromSetOk && toSetOk {
		for svr, _ := range fromSvrSet {
			for needs, _ := range toSvrSet {
				if svr != needs {
					sc.servers[svr].Needs[needs] = true
					sc.servers[needs].NeededBy[svr] = true
				}
			}
		}
	}
}

const dependencyDetailsQuery = "select from_type, from_val, to_type, to_val from dependency_details"

func (sc *SysConfig) getDependencyDetails(db *sql.DB) {
	if sc.fatal {
		return
	}

	rows, err := db.Query(dependencyDetailsQuery)
	if err != nil {
		sc.addFatalError("getDependencyDetails failed: SQL error: " +
			err.Error())
	} else {
		defer rows.Close()
		var (
			fromType string
			fromVal  string
			toType   string
			toVal    string
		)

		for rows.Next() {
			err := rows.Scan(&fromType, &fromVal, &toType, &toVal)
			if err != nil {
				sc.addFatalError(
					"getDependencyDetails failed: couldn't scan the columns: " +
						err.Error())
				break
			}
			sc.populateServerSet(fromType, fromVal, toType, toVal)
		}
	}
	sc.detectCycles()
}

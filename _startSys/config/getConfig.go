package config

import (
	"database/sql"
	"golem/paramSetter"
	"golem/params"

	_ "github.com/lib/pq"
)

var abortOnError bool = true
var reportErrors bool = true

func init() {
	var paramGroupName = "params.SysStarter.config"
	params.SetGroupDescription(paramGroupName, `SysStarter configuration parameters

These are concerned with the configuration details for starting the system.`)

	{
		p := params.New("dont-exit-on-config-error",
			paramSetter.BoolSetterNot{Value: &abortOnError},
			`don't exit if there have been any errors found when constructing the configuration details.
The default behaviour is to abort`)
		p.SetGroupName(paramGroupName)
	}
	{
		p := params.New("dont-report-errors",
			paramSetter.BoolSetterNot{Value: &reportErrors},
			`don't report any errors found when constructing the configuration details.
The default behaviour is to print them`)
		p.SetGroupName(paramGroupName)
	}
}

const (
	dbName = "sys_starter_config"
)

func (sc *SysConfig) getFromDB(db *sql.DB) {
	sc.getLocations(db)
	sc.getOperatingSystems(db)
	sc.getDatacentres(db)
	sc.getHosts(db)
	sc.getServers(db)
	sc.getHostServerMap(db)
	sc.getHostAttrTemplates(db)
	sc.getHostAttrs(db)
	sc.getClasses(db)
	sc.getServerClassMap(db)
	sc.getServerHostMustHave(db)
	sc.getDependencyDetails(db)
}

func GetConfig() *SysConfig {
	var sc *SysConfig = NewSysConfig()

	db, err := sql.Open("postgres", "password=secret dbname="+dbName)

	if err != nil {
		sc.addFatalError("failed to open the database connection: " + err.Error())
	} else {
		defer db.Close()

		sc.getFromDB(db)
	}

	if reportErrors {
		sc.ReportErrors()
	}

	return sc
}

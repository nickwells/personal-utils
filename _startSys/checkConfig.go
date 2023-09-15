package main

import (
	"fmt"
	"sysStarter/config"
)

func checkConfig(sc *config.SysConfig) {
	for _, svr := range sc.AllServers() {
		if len(svr.Hosts) == 0 {
			fmt.Println("server:", svr.Name, "is not running on any hosts")
		}
	}
}

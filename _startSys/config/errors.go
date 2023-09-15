package config

import (
	"fmt"
	"os"
)

func (sc *SysConfig) addError(msg string) {
	sc.errs = append(sc.errs, msg)
}
func (sc *SysConfig) addFatalError(msg string) {
	sc.addError(msg)
	sc.fatal = true
}

func (sc *SysConfig) ReportErrors() {
	if errCnt := len(sc.errs); errCnt > 0 {
		fmt.Println(errCnt, " system configuration error(s) detected:")
		for _, err := range sc.errs {
			fmt.Println(err)
		}
		if abortOnError {
			fmt.Println("Aborting")
			os.Exit(1)
		}
	}
}

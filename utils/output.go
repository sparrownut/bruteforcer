package utils

import "bruteforcer/Global"

func DBGLOG(str string) {
	if Global.DBG {
		println(str)
	}
}

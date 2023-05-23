package conf

import "log"

func LogConf() {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
}

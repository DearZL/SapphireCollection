package conf

import "log"

func DefaultConf() {
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
}
func init() {
	DefaultConf()
}

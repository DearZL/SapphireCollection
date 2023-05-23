package main

import (
	"P/router"
	"github.com/spf13/viper"
	"log"
)

func main() {
	r := router.NewRouter()
	//err = r.RunTLS(":9090", "./cert.pem", "./private.key")
	err := r.Run(viper.GetString("app.host") + ":" + viper.GetString("app.port"))
	if err != nil {
		log.Println(err.Error())
	}
}

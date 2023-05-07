package main

import (
	"P/model"
	"P/resp"
	"P/router"
	"encoding/gob"
	"log"
)

func main() {
	gob.Register(&model.User{})
	gob.Register(&resp.User{})
	r := router.NewRouter()
	//err = r.RunTLS(":9090", "./cert.pem", "./private.key")
	err := r.Run(":9090")
	if err != nil {
		log.Println(err.Error())
	}
}

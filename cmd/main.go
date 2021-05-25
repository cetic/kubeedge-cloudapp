package main

import (
	"CloudApp/config"
	"CloudApp/internal/core"
	log "github.com/sirupsen/logrus"
)


func main(){
	config.Parse()
	log.Infof("%+v",config.Conf)
	log.Debug("%+v",config.Conf)
	for _, dev := range config.Conf.Devices {
		d:= core.Device{ID: dev}
		go d.Listen()
	}
	select{}
}
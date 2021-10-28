package config

import (
	"CloudApp/internal/api"
	"flag"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var Conf Config

type Api struct {
	Url string `yaml:"url"`
}

type Log struct {
	LogLevels     string `yaml:"level"`
	LogFormatters string `yaml:"formatter"`
}

type Triggering struct {
	Condition string     `yaml:"condition"`
	Action    api.Update `yaml:"action"`
}

// config means all configurations used by the applications
type Config struct {
	Triggering []Triggering `yaml:"triggering"`
	Devices    []string     `yaml:"devices"`
	Api        Api          `yaml:"api"`
	Log        Log          `yaml:"log"`
	Polling    int          `yaml:"polling"`
}

func Parse() {
	configFile := flag.String("c", "../config/config.yaml", "config file")
	flag.Parse()
	args := flag.Args()
	myself := os.Args[0]
	if len(args) != 0 {
		log.Errorf("Wrong number of argument : %s [-c configfile] \n", myself)
		os.Exit(1)
	}
	yamlFile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Errorf("Unmarshall err   #%v ", err)
	}
	lvl, err := log.ParseLevel(Conf.Log.LogLevels)
	if err != nil {
		log.Errorf("Can Parse Log Level option   #%v ", err)
	}
	log.SetLevel(lvl)
	if Conf.Log.LogFormatters == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

package core

import (
	"CloudApp/config"
	"CloudApp/internal/api"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Device struct {
	ID          string
	Value       api.Response
	lastTrigger string
}

func (d *Device) UpdateValue() error {
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	log.Debug(config.Conf.Api.Url + "/api/v1/device/" + d.ID)
	request, err := http.NewRequest("GET", config.Conf.Api.Url+"/api/v1/device/"+d.ID, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Origin", "EnedisApi")
	request.Header.Set("X-Message-Id", "unset")

	rsp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, e := ioutil.ReadAll(rsp.Body)
	log.Debug(string(b))
	err = json.Unmarshal(b, &d.Value)
	if err != nil {
		return err
	}
	if e != nil {
		log.Error(e)
	}
	return nil
}

func (d *Device) SetAction(input api.Update) error {
	jsonaction, err := json.Marshal(input)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	request, err := http.NewRequest("POST", config.Conf.Api.Url+"/api/v1/device/"+d.ID, bytes.NewBuffer(jsonaction))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Origin", "gac")
	request.Header.Set("X-Message-Id", "unset")

	rsp, err := client.Do(request)
	if err != nil {
		return err
	}
	log.Info(rsp.Body)
	return nil
}

func (d *Device) Listen() {
	for {
		log.Debugf("Old Trigger : %s", d.lastTrigger)
		e := d.UpdateValue()
		if e != nil {
			log.Error(e)
		}
		log.Debugf("Updated Trigger : %s", d.Value.Trigger)
		for _, trig := range config.Conf.Triggering {
			if trig.Condition == d.Value.Trigger && d.lastTrigger != d.Value.Trigger {
				log.Debug("Trigger Find")
				d.lastTrigger = d.Value.Trigger
				e = d.SetAction(trig.Action)
				if e != nil {
					log.Error(e)
				}
			}
		}
		time.Sleep(time.Millisecond * time.Duration(config.Conf.Polling))
	}
}

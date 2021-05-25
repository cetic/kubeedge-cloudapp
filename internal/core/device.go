package core

import (
	"CloudApp/config"
	"CloudApp/internal/api"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type Device struct {
	ID string
	Value api.Response
}

func (d *Device) UpdateValue() error {
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	log.Debug(config.Conf.Api.Url+"/api/v1/device/"+d.ID)
	request, err := http.NewRequest("GET",config.Conf.Api.Url+"/api/v1/device/"+d.ID, nil)
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
	err = json.Unmarshal(b,&d.Value)
	if err != nil {
		return err
	}
	if e != nil {
		log.Error(e)
	}
	return nil
}


func (d *Device) SetAction(input api.Update) error {
	jsonMetering, err := json.Marshal(input)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	request, err := http.NewRequest("POST", config.Conf.Api.Url+"/api/v1/device/"+d.ID, bytes.NewBuffer(jsonMetering))
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
	if rsp.StatusCode != http.StatusNoContent {
		defer rsp.Body.Close()
		b, e := ioutil.ReadAll(rsp.Body)
		if e != nil {
			log.Error(e)
		}
		return fmt.Errorf("unable to post the metering %s: %s", rsp.Status, string(b))
	}
	log.Info(rsp.Body)
	return nil
}

func (d *Device) Listen() {
	for {
		oldTrigger := d.Value.Trigger
		log.Debugf("Old Trigger : %s",oldTrigger)
		e := d.UpdateValue()
		if e != nil {
			log.Error(e)
		}
		log.Debugf("Updated Trigger : %s",d.Value.Trigger)
		for _, trig := range config.Conf.Triggering {
			if trig.Condition == d.Value.Trigger && oldTrigger == d.Value.Trigger {
				log.Debug("Trigger Find")
				e = d.SetAction(trig.Action)
				if e != nil {
					log.Error(e)
				}
			}
		}
		time.Sleep(time.Millisecond * time.Duration(config.Conf.Polling))
	}
}
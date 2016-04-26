package tray_sensu_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/cvtienhoven/tray_plugin"
	"io/ioutil"
	"net/http"
)

/*
type Results []struct {
	Client string
	Check  struct {
		Status int
		Name   string
	}
}

type Clients []struct {
  Name string
  Subscriptions []string
}
*/
type Events []struct {
	Client struct {
		Tags          []string
		Name          string
		Subscriptions []string
	}
	Check struct {
		Name   string
		Tags   []string
		Status int
	}
}

var url string
var subscription string = ""
var tag string = ""

func setConfig(config tray_plugin.Config) {
	for _, element := range config {
		fmt.Println("Config: ", element.Key, "=", element.Value)
		if element.Key == "url" {
			url = element.Value
		} else if element.Key == "subscription" {
			subscription = element.Value
		} else if element.Key == "tag" {
			tag = element.Value
		}
	}
}

func GetStatus(config tray_plugin.Config) int {
	setConfig(config)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while executing request")
	}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error while reading response body")
	}

	defer res.Body.Close()

	var events = new(Events)
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Error parsing results: ", err)
	}

	var status int = 0
	for _, element := range *events {
		if status == 2 {
			break
		}
		//current event matches subscription on client level
		if subscription != "" && element.Client.Subscriptions != nil {
			if contains(element.Client.Subscriptions, subscription) && element.Check.Status > status && element.Check.Status < 3 {
				status = element.Check.Status
				fmt.Println(element.Client.Name, " matches subscription ", subscription, ": ", element.Client.Subscriptions)
				continue
			}
		}
		//current event matches on tag level
		if tag != "" && element.Check.Tags != nil {
			if contains(element.Check.Tags, tag) && element.Check.Status > status && element.Check.Status < 3 {
				status = element.Check.Status
				fmt.Println(element.Check.Name, " matches tag ", tag, ": ", element.Check.Tags)
				continue
			}
		}

		if tag == "" && subscription == "" && element.Check.Status > status && element.Check.Status < 3 {
			status = element.Check.Status
			continue
		}
	}

	fmt.Println(status)
	return status
}

func contains(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

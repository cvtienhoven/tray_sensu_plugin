package tray_sensu_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/cvtienhoven/tray_plugin"
	"io/ioutil"
	"net/http"
)

type Results []struct {
	Client string
	Check  struct {
		Status int
		Name   string
	}
}

var url string

func setConfig(config tray_plugin.Config) {
	for _, element := range config {
		if element.Key == "url" {
			url = element.Value
		}
	}
}

func GetStatus(config tray_plugin.Config) int {
	setConfig(config)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error")
	}
	body, err := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	var results = new(Results)
	err = json.Unmarshal(body, &results)
	if err != nil {
		fmt.Println("Error parsing results: ", err)
	}

	var status int = 0
	for _, element := range *results {
		if element.Check.Status > 0 {
			fmt.Println(element.Client, " - ", element.Check.Name, ": ", element.Check.Status)
			if element.Check.Status < 3 && element.Check.Status > status {
				status = element.Check.Status
			}
			if status == 2 {
				break
			}
		}
	}

	fmt.Println(status)
	return status
}

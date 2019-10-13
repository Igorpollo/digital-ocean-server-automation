package main

import (
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/digitalocean/godo"
	"os"

)

type WriteServer struct {
	ID string
	IP string
}


func Writeyml(droplet *godo.Droplet) {
	configData, _ := ioutil.ReadFile("./config.yml")
	config := Configs{}
	err := yaml.Unmarshal([]byte(configData), &config)
	if err != nil {
		fmt.Println(err)
	}


	ip, _ := droplet.PublicIPv4()
	config.CreatedByIP[ip] = CreatedServer{
		ID: droplet.ID,
		IP: ip,
	}

	d, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println(err)
	}

	// write to file
	f, err := os.Create("/tmp/dat2")
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("config.yml", d, 0644)
	if err != nil {
		fmt.Println(err)
	}

	f.Close()

}


func deleteIPYML(ip string) {
	configData, _ := ioutil.ReadFile("./config.yml")
	config := Configs{}
	err := yaml.Unmarshal([]byte(configData), &config)
	if err != nil {
		fmt.Println(err)
	}

	delete(config.CreatedByIP, ip)

	d, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println(err)
	}

	// write to file
	f, err := os.Create("/tmp/dat2")
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("config.yml", d, 0644)
	if err != nil {
		fmt.Println(err)
	}

	f.Close()	
}
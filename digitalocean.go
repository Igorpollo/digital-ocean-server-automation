package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/digitalocean/godo"
	log "github.com/igorpollo/go-custom-log"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
	"errors"
)

func createDigitalOcean(config Server, serverName string) error {
	tokenSource := &TokenSource{
		AccessToken: pat,
	}
	log.Info("Creating droplet...")
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)
	var keysSSH []godo.DropletCreateSSHKey
	keysSSH = append(keysSSH, godo.DropletCreateSSHKey{
		Fingerprint: os.Getenv("FINGERPRINT"),
	})
	createRequest := &godo.DropletCreateRequest{
		Name:    serverName,
		Region:  config.Region,
		Size:    config.Config,
		SSHKeys: keysSSH,
		Image: godo.DropletCreateImage{
			Slug: config.Image,
		},
	}

	ctx := context.TODO()

	newDroplet, _, err := client.Droplets.Create(ctx, createRequest)
	var drop = newDroplet
	if err != nil {
		fmt.Printf("Something bad happened: %s\n\n", err)
		return err
	}

	for drop.Status != "active" {
		drop, _, err = client.Droplets.Get(ctx, newDroplet.ID)
		time.Sleep(3 * time.Second)
	}
	go Writeyml(drop)
	log.Success("====== DROPLET CRIADO =======")
	time.Sleep(5 * time.Second)
	log.Info("CONECTANDO....")
	configSSH := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			hai(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	ip, err := drop.PublicIPv4()
	var conn *ssh.Client
	conn, err = ssh.Dial("tcp", ip+":22", configSSH)
	if err != nil {
		time.Sleep(5 * time.Second)
		conn, err = ssh.Dial("tcp", ip+":22", configSSH)
		if err != nil {
			return err
		}
	}
	log.Success("Conectado!")
	log.Info("Running commands...")
	for i:=0; i <  len(config.Commands); i++ {
		runCommand(config.Commands[i], conn)
	}
	log.Success("Runned all commands")
	conn.Close()
	return nil
}


func DeleteByIP(ip string) error {
	ctx := context.TODO()
	config := readConfig()
	if _, ok := config.CreatedByIP[ip]; !ok {
		return errors.New("Can't find server")
	}
	server := config.CreatedByIP[ip]
	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)

	_, err := client.Droplets.Delete(ctx, server.ID)
	if err != nil {
		return err
	}
	go deleteIPYML(ip)
	return nil
}
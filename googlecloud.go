package main

import (
	"context"
	"fmt"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"golang.org/x/crypto/ssh"
	log "github.com/igorpollo/go-custom-log"
	"io/ioutil"
	//"path/filepath"
	//"os"
	//"time"


)

func hai2(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)

	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		fmt.Println(err)
	}
	return ssh.PublicKeys(signer)
}


func createGoogleCloud(config Server, serverName string) error {
	log.Info("Creating Google Cloud Instance...")

	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile("./pollinho-a59f50e2734e.json"))
	if err != nil {
		fmt.Println(err)
	}
	// oi, err := computeService.Instances.Get("pollinho","us-central1-a","903903683863491545").Do()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	var networks []*compute.NetworkInterface
	var disks []*compute.AttachedDisk
	var acess []*compute.AccessConfig

	networks = append(networks, &compute.NetworkInterface{
		Name: "mynet",
		
	})
	disks = append(disks, &compute.AttachedDisk{
		Boot: true,
		AutoDelete: true,
		Mode: "READ_WRITE",
		Interface: "SCSI",
		InitializeParams: &compute.AttachedDiskInitializeParams{
			SourceImage: "projects/debian-cloud/global/images/family/debian-9",
			DiskName: "disk22-tesaate1",
		},
	})
	acess = append(acess, &compute.AccessConfig{
		SetPublicPtr: true,
		Type: "ONE_TO_ONE_NAT",
		Name: "My external IP",
		NetworkTier: "PREMIUM",
	})
	var instance compute.Instance
	instance.Zone = "us-central1-a"
	inst := compute.Instance{
		MachineType: "zones/us-central1-a/machineTypes/f1-micro",
		Zone: config.Region,
		Name: serverName,
		NetworkInterfaces: networks,
		Disks: disks,
	}
	oi, err := computeService.Instances.Insert("pollinho","us-central1-a", &inst).Do()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(oi.Id)

	// oi := computeService.Instances.Insert()
	return nil
}
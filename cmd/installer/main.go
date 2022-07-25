package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// usage prints a usage message and exits.
func usage(name string) {
	fmt.Printf("Usage:\n")
	fmt.Printf("    %s -i  install service\n", name)
	fmt.Printf("    %s -r  remove service\n", name)
	fmt.Printf("    %s -u  update service\n", name)
	os.Exit(1)
}

func main() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting executable: %v", err)
	}
	dir := filepath.Dir(executable)

	name, args := os.Args[0], os.Args[1:]
	if len(args) != 1 {
		usage(name)
	}
	switch args[0] {
	case "-i":
		install(dir)
	case "-r":
		remove()
	case "-u":
		update(dir)
	default:
		usage(name)
	}
}

// install creates a service to run service.exe, then starts the service.
func install(dir string) {
	m, err := mgr.Connect()
	if err != nil {
		log.Fatalf("error connecting to service control manager: %v", err)
	}

	servicePath := filepath.Join(dir, "service.exe")
	config := mgr.Config{
		DisplayName: "Sample Windows service written in Go",
		StartType:   mgr.StartAutomatic,
	}
	s, err := m.CreateService("sample-service", servicePath, config)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}
	defer s.Close()

	err = s.Start()
	if err != nil {
		log.Fatalf("error starting service: %v", err)
	}
}

// remove deletes and stops the service.
func remove() {
	m, err := mgr.Connect()
	if err != nil {
		log.Fatalf("error connecting to service control manager: %v", err)
	}

	s, err := m.OpenService("sample-service")
	if err != nil {
		log.Fatalf("error opening service: %v", err)
	}
	defer s.Close()

	err = s.Delete()
	if err != nil {
		log.Fatalf("error marking service for deletion: %v", err)
	}

	_, err = s.Control(svc.Stop)
	if err != nil {
		log.Fatalf("error requesting service to stop: %v", err)
	}
}

// update updates the service to version 2.
func update(dir string) {
	log.Printf("updating service to version 2...")

	m, err := mgr.Connect()
	if err != nil {
		log.Fatalf("error connecting to service control manager: %v", err)
	}

	service, err := m.OpenService("sample-service")
	if err != nil {
		log.Fatalf("error accessing service: %v", err)
	}

	config, err := service.Config()
	if err != nil {
		log.Fatalf("error getting service config: %v", err)
	}
	config.BinaryPathName = filepath.Join(dir, "service-2.exe")
	err = service.UpdateConfig(config)
	if err != nil {
		log.Fatalf("error updating config: %v", err)
	}

	log.Print("requesting service to stop")
	status, err := service.Control(svc.Stop)
	if err != nil {
		log.Fatalf("error requesting service to stop: %v", err)
	}
	log.Printf("sent stop; service state is %v", status.State)

	for i := 0; i <= 12; i++ {
		log.Printf("querying service status (attempt %d)", i)
		status, err = service.Query()
		if err != nil {
			log.Fatalf("error querying service: %v", err)
		}
		log.Printf("service state: %v", status.State)

		if status.State == svc.Stopped {
			log.Println("service state is 'stopped'")
			break
		}

		time.Sleep(10 * time.Second)
	}

	log.Println("starting service")
	err = service.Start()
	if err != nil {
		log.Fatalf("error starting service: %v", err)
	}
	log.Print("service updated to version 2")
}

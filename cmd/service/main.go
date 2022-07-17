package main

// Sample code for a Windows service.

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/svc"
)

const (
	exitCodeDirNotFound     = 1
	exitCodeErrorOpeningLog = 2
)

// main runs as a Windows service.
func main() {
	svc.Run("sample-service", new(handler))
}

type handler struct{}

// Execute implements the Windows service.
func (h *handler) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	s <- svc.Status{State: svc.StartPending}

	executable, err := os.Executable()
	if err != nil {
		return true, exitCodeDirNotFound
	}
	dir := filepath.Dir(executable)
	logPath := filepath.Join(dir, "service.log")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return true, exitCodeErrorOpeningLog
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	s <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}
	loop(r, s, logger)
	logger.Print("service: shutting down")
	return true, 0
}

// loop runs the work function every 30s and reacts to stop/shutdown requests.
func loop(r <-chan svc.ChangeRequest, s chan<- svc.Status, logger *log.Logger) {
	tick := time.Tick(30 * time.Second)
	logger.Print("service: up and running")
	for {
		select {
		case <-tick:
			work(logger)
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				s <- c.CurrentStatus
			case svc.Stop:
				logger.Print("service: got stop signal, exiting")
				return
			case svc.Shutdown:
				logger.Print("service: got shutdown signal, exiting")
				return
			}
		}
	}
}

// work is meant to do the actual work of the service.
func work(logger *log.Logger) {
	logger.Print("service says: hello")
}

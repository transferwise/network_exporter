//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log/level"
)

func reloadSignal() {

	// Signal handling
	hup := make(chan os.Signal, 1)
	signal.Notify(hup, syscall.SIGHUP)
	susr := make(chan os.Signal, 1)
	signal.Notify(susr, syscall.SIGUSR1)
	go func() {
		for {
			select {
			case <-hup:
				level.Debug(logger).Log("msg", "Signal: HUP")
				level.Info(logger).Log("msg", "ReLoading config")
				if err := sc.ReloadConfig(logger, *configFile); err != nil {
					level.Error(logger).Log("msg", "Reloading config skipped", "err", err)
					continue
				} else {
					monitorPING.DelTargets()
					monitorPING.CheckActiveTargets()
					monitorPING.AddTargets()
					monitorMTR.DelTargets()
					monitorMTR.CheckActiveTargets()
					monitorMTR.AddTargets()
					monitorTCP.DelTargets()
					monitorTCP.CheckActiveTargets()
					monitorTCP.AddTargets()
					monitorHTTPGet.DelTargets()
					monitorHTTPGet.AddTargets()
				}
			case <-susr:
				level.Debug(logger).Log("msg", "Signal: USR1")
				fmt.Printf("PING: %+v\n", monitorPING)
				fmt.Printf("MTR: %+v\n", monitorMTR)
				fmt.Printf("TCP: %+v\n", monitorTCP)
				fmt.Printf("HTTPGet: %+v\n", monitorHTTPGet)
			}
		}
	}()
}

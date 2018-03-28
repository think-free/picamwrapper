package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/fulldump/goconfig"
	"github.com/think-free/other/picamwrapper/config"
	"github.com/think-free/other/picamwrapper/httpserver"
	"github.com/think-free/other/picamwrapper/mqtt"
	"github.com/think-free/other/picamwrapper/status"
)

func main() {

	// Getting configuration
	cfg := &config.Config{DisableMqtt: false}
	goconfig.Read(cfg)

	// Getting internal object
	internal := &config.Internal{}
	internal.AutoMode = make(chan bool)
	internal.Chwritest = make(chan bool)
	internal.Chwriteauto = make(chan bool)
	internal.Chgetauto = make(chan bool)
	internal.Chgetst = make(chan bool)

	// Create the picam process
	pccmd := exec.Command("bash", "-c", cfg.PiCamPath+"/picam "+cfg.PiCamParams)
	pccmd.Stdout = os.Stdout
	pccmd.Stderr = os.Stderr

	// Handle ctrl+c

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for {
			select {
			case _ = <-c:
				pccmd.Process.Kill()
				fmt.Println("\nClosing application")
				os.Exit(0)
			}
		}
	}()

	// Creating application routines
	mq := mqtt.New(cfg, internal)
	go mq.Run()

	st := status.New(cfg, internal)
	go st.Run()

	ht := httpserver.New(cfg, internal)
	go ht.Run()

	// Starting picam
	pccmd.Start()
	pccmd.Wait()
}

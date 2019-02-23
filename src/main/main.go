// Package main implements the CLI application :
// Write an application (CLI) that creates a batch of 100 unique DevEUIs and registers
// them with the LoRaWAN api.
package main

import (
	"../data"
	"../multiThreadApiCall"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

func main() {
	// creating and initializing the waiting group
	var wg sync.WaitGroup
	wg.Add(1)

	// setting the max num of routines (writers routines 7, launcher routine 1, main routine 1, reader routine 1)
	runtime.GOMAXPROCS(data.MAXWRITERS + 3)

	// initialize the channel for the responses
	readerChannel := make(chan data.ResponseAPInewDevEUI)

	// initialize the channel used for the replay of errored calls (!200)
	redoChan := make(chan int, data.N)
	for i := 0; i < data.N; i++ {
		redoChan <- 1
	}

	// generate the suffix of the DevEUI
	suffix, _ := multiThreadApiCall.RandomHex(11)

	// handle the signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// writers launcher
	go multiThreadApiCall.RoutinesLauncher(suffix, readerChannel, redoChan, sigs)

	// reader routine
	go multiThreadApiCall.ResponsesReader(readerChannel, redoChan, sigs, &wg)

	wg.Wait()
}

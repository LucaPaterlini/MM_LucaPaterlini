package coreCli

import (
	"../data"
	"../multiThreadApiCall"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

// Create100NewIds create and launch the go routines that are going to launch the functions that start
// the writers function and start the readers go routine, keeping a final loop to handle the print
// of the reader go routine

func Create100NewIds(debug bool) string{
	stdout := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)

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
	go multiThreadApiCall.ResponsesReader(readerChannel, redoChan,  &wg,stdout)

	//go Offload the answers
	var v []string
	for item := range stdout {
		v =append(v, item)
		if debug {fmt.Println(item)}
	}
	urlsJson, _ := json.MarshalIndent(v,"","  ")
	wg.Wait()
	return string(urlsJson)
}
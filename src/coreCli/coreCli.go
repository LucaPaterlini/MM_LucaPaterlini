package coreCli

import (
	"../data"
	"../multiThreadApiCall"
	"encoding/json"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

func Create100NewIds() string{
	stdout := make(chan string)
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
	go multiThreadApiCall.ResponsesReader(readerChannel, redoChan,  &wg,stdout)

	//go Offload(stdout,&wg)
	var v []string
	for item := range stdout {
		v =append(v, item)
		//fmt.Println(item)
	}
	urlsJson, _ := json.MarshalIndent(v,"","  ")
	wg.Wait()
	return string(urlsJson)
}

//func Offload(std chan string, wg *sync.WaitGroup){
//	var v []string
//	for item := range std {
//		v =append(v, item)
//		//fmt.Println(item)
//	}
//	urlsJson, _ := json.MarshalIndent(v,"","  ")
//	fmt.Println(string(urlsJson))
//	wg.Done()
//}

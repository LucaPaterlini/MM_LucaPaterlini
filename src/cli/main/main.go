// Package main inplements the CLI application :
// Write an application (CLI) that creates a batch of 100 unique DevEUIs and registers
// them with the LoRaWAN api.
package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

// postCallJson execute a Post call to the default url+path of LoRaWAN
func postCallJson(DevEUI string,responseChannel chan<- responseAPInewDevEUI, semaphoreW <-chan int){

	// response initialization
	responseItem := responseAPInewDevEUI{0,nil,DevEUI}

	// prepare the header and the body of the call
	url := urlLoRaWAN+pathRegisterNewDevEUI
	var jsonBody = []byte(`{"deveui":"`+DevEUI+`"}`)
	req,err := http.NewRequest("POST", url,  bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if err!=nil {
		responseItem.err = err
		responseChannel <- responseItem
		return
	}
	// execute the call
	client := &http.Client{}
	res,err := client.Do(req)
	if err!=nil {
		responseItem.err = err
		responseChannel <- responseItem
		return
	}
	// just considering the statuscode
	responseItem.sCode = res.StatusCode
	responseChannel <- responseItem
	fmt.Println("starting: ",DevEUI)
	<-semaphoreW
	return
}

// random returns an the suffix of the DevEUI prefix of 11 characters
func randomHex(n int) (returnString string, err error) {
	rand.Seed(time.Now().Unix())
	bytes := make([]byte, n/2+1)
	if _, err = rand.Read(bytes); err != nil {
		return
	}
	returnString = hex.EncodeToString(bytes)[:n]
	return
}

// routinesLauncher start the go routines with the limit of 1

func routinesLauncher(suffix string, readerChannel chan<- responseAPInewDevEUI,
	redoChan <-chan int, term <-chan os.Signal) {
	// create a semaphore
	semaphoreW := make(chan int,MAXWRITERS)
	// counter loops
	i := 0

	L:
	// reply all the wrong calls and breaks in case the redo channel is closed
	// break if the channel is closed
	for range redoChan {

		// assembling the DevEUI and changing the case to uppercase
		DevEUI := strings.ToUpper(suffix + fmt.Sprintf("%05x", i))
		// inc the writers counter
		i++
		select{
			case <-term:
				fmt.Println("stop pls")
				close(readerChannel)
				break L
			default :
				semaphoreW <- 1
				go postCallJson(DevEUI, readerChannel,semaphoreW)
			}
	}
}


// responsesReader read from the readerChannel and write on redochan in case of errored calls
func responsesReader(readerChannel <-chan responseAPInewDevEUI, redoChan chan<- int,term chan<-os.Signal, wg *sync.WaitGroup){
	// the counter for the output is general inside the function itself
	i:=0
	// dec the waiting group

	defer wg.Done()
	for item := range  readerChannel {

		// Debug : && item.DevEUI[len(item.DevEUI)-1:]!="7"  append to the following condition to test
		// in case there are no !200 responses from the endpoint
		if item.err == nil && item.sCode == 200 && item.DevEUI[len(item.DevEUI)-1:]!="7" {
			fmt.Printf("DevEUI #% 3d: %s\n", i, item.DevEUI)
			i++
		}else {
			redoChan<-1
		}
		if i>=N { close(redoChan);return  }
	}

}


func main() {
	// creating and initializing the waiting group
	var wg sync.WaitGroup
	wg.Add(1)

	// setting the max num of routines (writers routines 7, launcher routine 1, main routine 1, reader routine 1)
	runtime.GOMAXPROCS(MAXWRITERS+3)

	// initialize the channel for the responses
	readerChannel := make(chan responseAPInewDevEUI)

	// initialize the channel used for the replay of errored calls (!200)
	redoChan := make(chan int,N)
	for i:=0;i<N;i++{redoChan<-1}

	// generate the suffix of the DevEUI
	suffix, _ := randomHex(11)

	// handle the signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM,syscall.SIGKILL)

	// writers launcher
	go routinesLauncher(suffix,readerChannel,redoChan,sigs)

	// reader routine
	go responsesReader(readerChannel,redoChan,sigs,&wg)

	wg.Wait()

}

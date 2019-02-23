package multiThreadApiCall

import (
	"../data"
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// postCallJson execute a Post call to the default url+path of LoRaWAN
func PostCallJson(DevEUI string, responseChannel chan<- data.ResponseAPInewDevEUI, semaphoreW <-chan int) {

	// response initialization
	responseItem := data.ResponseAPInewDevEUI{0, nil, DevEUI}

	// prepare the header and the body of the call
	url := data.UrlLoRaWAN + data.PathRegisterNewDevEUI
	var jsonBody = []byte(`{"deveui":"` + DevEUI + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		responseItem.Err = err
		responseChannel <- responseItem
		return
	}
	// execute the call
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		responseItem.Err = err
		responseChannel <- responseItem
		return
	}
	// just considering the statuscode
	responseItem.SCode = res.StatusCode
	responseChannel <- responseItem
	<-semaphoreW
	return
}

// random returns an the suffix of the DevEUI prefix of 11 characters
func RandomHex(n int) (returnString string, err error) {
	rand.Seed(time.Now().Unix())
	bytes := make([]byte, n/2+1)
	if _, err = rand.Read(bytes); err != nil {
		return
	}
	returnString = hex.EncodeToString(bytes)[:n]
	return
}

// routinesLauncher start the go routines with the limit of 1

func RoutinesLauncher(suffix string, readerChannel chan<- data.ResponseAPInewDevEUI,
	redoChan <-chan int, term <-chan os.Signal) {
	// create a semaphore
	semaphoreW := make(chan int, data.MAXWRITERS)
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
		select {
		case <-term:
			fmt.Println("stop pls")
			close(readerChannel)
			break L
		default:
			semaphoreW <- 1
			go PostCallJson(DevEUI, readerChannel, semaphoreW)
		}
	}
}

// responsesReader read from the readerChannel and write on redochan in case of errored calls
func ResponsesReader(readerChannel <-chan data.ResponseAPInewDevEUI, redoChan chan<- int, wg *sync.WaitGroup) {
	// the counter for the output is general inside the function itself
	i := 0
	// dec the waiting group

	defer wg.Done()
	for item := range readerChannel {

		// Debug : && item.DevEUI[len(item.DevEUI)-1:]!="7"  append to the following condition to test
		// in case there are no !200 responses from the endpoint
		if item.Err == nil && item.SCode == 200 && item.DevEUI[len(item.DevEUI)-1:] != "7" {
			fmt.Printf("DevEUI #% 3d: %s\n", i, item.DevEUI)
			i++
		} else {
			redoChan <- 1
		}
		if i >= data.N {
			close(redoChan)
			return
		}
	}

}

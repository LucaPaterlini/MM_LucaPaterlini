package multiThreadApiCall

import (
	"../data"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
)

func TestPostCallJson(t *testing.T) {
	responseChannel := make(chan data.ResponseAPInewDevEUI)
	semaphoreW := make(chan int)
	defer close(responseChannel)
	defer close(semaphoreW)

	go func() { semaphoreW <- 1 }()
	go PostCallJson("85058B9D35200069", responseChannel, semaphoreW)

	c := <-responseChannel

	if c.Err != nil {
		t.Error(c.Err.Error())
	}
}

func TestRandomHex(t *testing.T) {
	s,err:=RandomHex(11)
	fmt.Println(s,err)
	if err != nil || len(s) != 11  {
		t.Error(err.Error())
	}
}

func TestRoutinesLauncher(t *testing.T) {
	suffix:="AAAAAAAAA"
	responseChannel := make(chan data.ResponseAPInewDevEUI)
	redoChan := make(chan int)
	defer close(responseChannel)
	defer close(redoChan)

	go func() { redoChan <- 1 }()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// writers launcher
	go RoutinesLauncher(suffix,responseChannel,redoChan, sigs)

	ass :=<-responseChannel
	s := fmt.Sprint(ass)
	if s != "{200 <nil> AAAAAAAAA00000}"{
		t.Error(ass.Err.Error())
	}

}

func TestResponsesReader(t *testing.T) {

	responseChannel := make(chan data.ResponseAPInewDevEUI)
	var redoChan  chan int
	var wg sync.WaitGroup
	wg.Add(1)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	stdout := make(chan string)

	go func (){
		responseChannel<-data.ResponseAPInewDevEUI{200 ,nil ,"AAAAAAAAA000000" }
		close(responseChannel)
		}()
	go ResponsesReader(responseChannel, redoChan, &wg,stdout)

	wg.Wait()

}
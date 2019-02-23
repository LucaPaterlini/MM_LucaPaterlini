package main

import (
	"fmt"
	"math/rand"
	"time"
)

import gocache "github.com/pmylund/go-cache"

var cache = gocache.New(1*time.Minute, 3*time.Minute)


func f ()int {
	time.Sleep(1*time.Second)
	return rand.Int()
}

func main(){
	//key := "1"
	//for i:=0;i<2;i+=1 {
	//	_ = cache.Add(key, f(), gocache.DefaultExpiration)
	//	response, _ := cache.Get(key)
	//	fmt.Println(response)
	//}
	var i []string
	i =  append(i,"ciao")
	i = append(i," cane")
	fmt.Println(i)

}

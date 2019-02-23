// Package main implements the CLI application :
// Write an application (CLI) that creates a batch of 100 unique DevEUIs and registers
// them with the LoRaWAN api.
package main

import (
	"../coreCli"
	"fmt"
)

func main() {
	fmt.Println(coreCli.Create100NewIds())
}

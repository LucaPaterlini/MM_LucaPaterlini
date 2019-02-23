// Package main implements the CLI application :
// Write an application (CLI) that creates a batch of 100 unique DevEUIs and registers
// them with the LoRaWAN api.
package main

import (
	"../coreCli"
)

func main() {
	coreCli.Create100NewIds(true)
}

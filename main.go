package main

import (
	"github.com/cmcpasserby/unity-loader/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

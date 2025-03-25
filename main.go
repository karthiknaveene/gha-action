package main

import (
	"gha-action/cmd"
	"log"
)

func main() {

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/dev2choiz/sasspiler/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

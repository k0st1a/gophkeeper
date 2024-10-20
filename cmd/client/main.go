package main

import (
	"fmt"

	"github.com/k0st1a/gophkeeper/internal/application/client"
	"github.com/rs/zerolog/log"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n"+
		"Build date: %s\n"+
		"Build commit: %s\n",
		buildVersion, buildDate, buildCommit)

	err := client.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

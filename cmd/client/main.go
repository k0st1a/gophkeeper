package main

import (
	"github.com/k0st1a/gophkeeper/internal/application/client"
	"github.com/rs/zerolog/log"
)

func main() {
	err := client.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

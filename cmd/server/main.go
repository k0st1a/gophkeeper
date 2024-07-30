package main

import (
	"github.com/k0st1a/gophkeeper/internal/application/server"
	"github.com/rs/zerolog/log"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

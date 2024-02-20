package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	debugFlag := flag.Bool("debug", false, "Flag for debug level with console log outputs")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Debug().Msg("End Init Func")
}

func main() {
	log.Debug().Msg("Start Main Func")
}

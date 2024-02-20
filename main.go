package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	debugFlag := flag.Bool("debug", false, "Flag for debug level with console log outputs")

	flag.Parse()

	logFileHandle := &lumberjack.Logger{
		Filename: "./logs/log",
		MaxSize:  100,
		MaxAge:   31,
		Compress: true,
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.
		With().Caller().Logger().
		With().Timestamp().Logger()

	log.Logger = log.Output(logFileHandle)
	if *debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFileHandle)
		log.Logger = log.Output(multiWriter)
	}

	log.Debug().Msg("End Init Func")
}

func main() {
	log.Debug().Msg("Start Main Func")
}

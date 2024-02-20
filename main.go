package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	mymiddleware "github.com/hmcalister/WebAuthnSSO/middleware"
	"github.com/hmcalister/WebAuthnSSO/routes"
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

	router := chi.NewRouter()
	router.Use(mymiddleware.ZerologLogger)
	router.Use(mymiddleware.RecoverWithInternalServerError)

	router.Get("/api/v1/heartbeat", routes.Heartbeat)
	router.Get("/api/v1/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("panic function called")
	})

	http.ListenAndServe(":3000", router)
}

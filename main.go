package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/hmcalister/AuthSSO/database"
	mymiddleware "github.com/hmcalister/AuthSSO/middleware"
	"github.com/hmcalister/AuthSSO/routes/api/apiv1"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	databaseManager *database.DatabaseManager
)

func init() {
	var err error

	debugFlag := flag.Bool("debug", false, "Flag for debug level with console log outputs")
	databaseFilePath := flag.String("databaseFilePath", "database.sqlite", "The path to the database file on disk.")
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

	apiV1Router := apiv1.NewApiRouter()
	router.Mount("/api/v1", apiV1Router.Router())

	http.ListenAndServe(":3000", router)
}

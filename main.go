package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	authenticationmaster "github.com/hmcalister/AuthSSO/authenticationMaster"
	"github.com/hmcalister/AuthSSO/database"
	commonMiddleware "github.com/hmcalister/GoChi-CommonMiddleware"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	databaseManager *database.DatabaseManager
	port            *int
	secretKey       []byte
)

func init() {
	var err error

	port = flag.Int("port", 6585, "The port to use for the HTTP server.")
	debugFlag := flag.Bool("debug", false, "Flag for debug level with console log outputs.")
	databaseFilePath := flag.String("databaseFilePath", "database.sqlite", "The path to the database file on disk.")
	secretKeyFile := flag.String("secretKeyFile", "key.secret", "The path to the file containing the secret key for JWTAuth.")
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

	log.Debug().
		Str("DatabaseFilePath", *databaseFilePath).
		Msg("Creating Database")
	databaseManager, err = database.NewDatabase(*databaseFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error during creation of database manager")
	}

	secretKey, err = os.ReadFile(*secretKeyFile)
	if err != nil {
		log.Fatal().Str("FilePath", *secretKeyFile).Err(err).Msg("Could not open secret file for JWTAuth")
	}

	log.Debug().Msg("End Init Func")
}

func main() {
	defer databaseManager.CloseDatabase()

	log.Debug().Msg("Start Main Func")

	router := chi.NewRouter()
	router.Use(commonMiddleware.ZerologLogger)
	router.Use(commonMiddleware.RecoverWithInternalServerError)

	authMaster := authenticationmaster.NewAuthenticationMaster(databaseManager, secretKey)
	router.Post("/api/register", authMaster.Register)
	router.Post("/api/login", authMaster.Login)
	router.Get("/api/authenticate", authMaster.AuthenticateRequest)

	// --------------------------------------------------------------------------------
	// Example authenticated route
	// --------------------------------------------------------------------------------
	// router.Group(func(r chi.Router) {
	// 	r.Use(jwtauth.Verifier(tokenAuth))
	// 	r.Use(jwtauth.Authenticator(tokenAuth))
	// 	r.Get("/private", ...)
	// })

	targetBindAddress := fmt.Sprintf("localhost:%v", *port)
	log.Info().Msgf("Starting server on %v", targetBindAddress)
	err := http.ListenAndServe(targetBindAddress, router)
	if err != nil {
		log.Fatal().Err(err).Msg("Error during http listen and serve")
	}
}

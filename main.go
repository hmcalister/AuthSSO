package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	authenticationmaster "github.com/hmcalister/AuthSSO/authenticationMaster"
	"github.com/hmcalister/AuthSSO/database"
	commonMiddleware "github.com/hmcalister/GoChi-CommonMiddleware"
	"github.com/phsym/console-slog"

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

	var slogHandler slog.Handler
	if *debugFlag {
		slogHandler = console.NewHandler(os.Stdout, &console.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	} else {
		slogHandler = slog.NewJSONHandler(logFileHandle, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	slog.SetDefault(slog.New(
		slogHandler,
	))

	slog.Debug("Creating Database", "DatabaseFilePath", *databaseFilePath)
	databaseManager, err = database.NewDatabase(*databaseFilePath)
	if err != nil {
		slog.Error("Error during creation of database manager", "Error", err)
		os.Exit(1)
	}

	secretKey, err = os.ReadFile(*secretKeyFile)
	if err != nil {
		slog.Error("Could not open secret file for JWTAuth", "FilePath", *secretKeyFile, "Error", err)
		os.Exit(1)
	}
}

func main() {
	defer databaseManager.CloseDatabase()

	slog.Debug("Start Main Func")

	router := chi.NewRouter()
	router.Use(commonMiddleware.SlogLogger)
	router.Use(commonMiddleware.RecoverWithInternalServerError)

	authMaster := authenticationmaster.NewAuthenticationMaster(databaseManager, secretKey)
	router.Post("/api/register", authMaster.Register)
	router.Post("/api/login", authMaster.Login)
	router.Get("/api/authenticate", authMaster.AuthenticateRequest)

	targetBindAddress := fmt.Sprintf("localhost:%v", *port)
	slog.Info("Starting server", "Address", targetBindAddress)
	err := http.ListenAndServe(targetBindAddress, router)
	if err != nil {
		slog.Error("Error during http listen and serve", "Error", err)
		os.Exit(1)
	}
}

package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Application struct
type application struct {
	logger *slog.Logger
}

func main() {
	// Command line flag "addr" to set the port at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse the command line flag, needs to be called before before using the addr variable or default value will be used
	flag.Parse()

	// Structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// Initialize application struct containing the dependencies
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", "addr", *addr)
	// Call app.routes() and pass it to http.ListenAndServe
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

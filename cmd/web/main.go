package main

import (
	"database/sql"
	"flag"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.macaronj/internal/models"
)

// Application struct
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	// Command line flag "addr" to set the port at runtime
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Command line flag for MSQL DSN string
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	// Parse the command line flag, needs to be called before before using the addr variable or default value will be used
	flag.Parse()

	// Structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// defer to close connection pool before main() exits
	defer db.Close()

	// Initialize template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	// Initialize application struct containing the dependencies
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	logger.Info("starting server", "addr", *addr)
	// Call app.routes() and pass it to http.ListenAndServe
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// openDB() wraps sql.Open() and returns sql.DB connection for given string
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

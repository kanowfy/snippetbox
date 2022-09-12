package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kanowfy/snippetbox/pkg/models/mysql"
)

type application struct {
	infoLog       *log.Logger
	errLog        *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":8080", "Network Address")
	dsn := flag.String("dsn", "web:a@/snippetbox?parseTime=true", "MySQL Data Source Name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Println(err.Error())
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errLog.Println(err.Error())
	}

	app := &application{
		infoLog:       infoLog,
		errLog:        errLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Handler:      app.routes(),
		Addr:         *addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		ErrorLog:     app.errLog,
	}

	app.infoLog.Printf("Listening on port %s", *addr)
	app.errLog.Fatal(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
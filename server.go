package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/unrolled/render"

	"log"
)

var format = render.New()

func initDatabase(db *sqlx.DB) error {
	// Check if persons table exists
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='persons')")
	if err != nil {
		return err
	}

	if !exists {
		// Read and execute SQL dump
		file, err := os.Open("dump.sql")
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var currentStatement strings.Builder

		for scanner.Scan() {
			line := scanner.Text()

			// Skip comments and empty lines
			if strings.HasPrefix(strings.TrimSpace(line), "--") || strings.TrimSpace(line) == "" {
				continue
			}

			currentStatement.WriteString(line)
			currentStatement.WriteString(" ")

			// If the line ends with a semicolon, execute the statement
			if strings.HasSuffix(strings.TrimSpace(line), ";") {
				stmt := currentStatement.String()
				_, err := db.Exec(stmt)
				if err != nil {
					return err
				}
				currentStatement.Reset()
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
		log.Println("Database initialized successfully!")
	}
	return nil
}

func main() {
	// config
	Config.LoadFromFile("./config.yml")

	// db connection
	db, err := sqlx.Connect("sqlite3", Config.DataSourceName())
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database if needed
	if err := initDatabase(db); err != nil {
		log.Fatal(err)
	}

	// http router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	r.Get("/api/data/{table}", func(w http.ResponseWriter, r *http.Request) {
		table := chi.URLParam(r, "table")
		err := getDataFromDB(w, db, table, nil)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}
	})

	r.Post("/api/data/{table}", func(w http.ResponseWriter, r *http.Request) {
		table := chi.URLParam(r, "table")

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}

		err = getDataFromDB(w, db, table, body)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}
	})

	r.Get("/api/data/{table}/{field}/suggest", func(w http.ResponseWriter, r *http.Request) {
		table := chi.URLParam(r, "table")
		field := chi.URLParam(r, "field")

		err := getSuggestFromDB(w, db, table, field)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}
	})

	log.Printf("Start server at %s", Config.Server.Port)
	http.ListenAndServe(Config.Server.Port, r)
}

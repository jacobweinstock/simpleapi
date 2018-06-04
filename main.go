package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var logger, err = zap.NewProduction()

func main() {
	defer logger.Sync()
	logger.Info("starting server...")
	connStr := "postgres://myuser:mypassword@localhost:5432/simpleapi?sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open connections: %s", err)
	}

	db := &database{conn: conn}

	http.Handle("/developers", withLogging(&handler{db: db}))
	http.HandleFunc("/healthz", healthz)

	http.ListenAndServe(":8080", nil)
}

func withLogging(h http.Handler) http.Handler {
	defer logger.Sync()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("before")
		h.ServeHTTP(w, r)
		logger.Info("after")
	})
}

type database struct {
	conn *sql.DB
}

type handler struct {
	db *database
}

func (d *database) AllDevelopers() []developer {
	rows, err := d.conn.Query("SELECT id, name FROM developers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var devs []developer
	for rows.Next() {
		var dev developer
		if err := rows.Scan(&dev.ID, &dev.Name); err != nil {
			log.Fatal(err)
		}
		devs = append(devs, dev)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return devs
}

func pingdb() bool {
	connStr := "postgres://myuser:mypassword@localhost:5432/simpleapi?sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	logger.Info("before ping")
	err = conn.Ping()
	defer conn.Close()
	result := true
	if err != nil {
		result = false
		logger.Info("TEST")
		// log.Fatalf("failed to open connections: %s", err)

	}
	logger.Info("err", zap.Error(err))
	logger.Info("result", zap.Bool("result", result))
	return result
}

func healthz(w http.ResponseWriter, r *http.Request) {
	defer logger.Sync()
	w.Header().Add("content-type", "application/json;charset=utf-8")
	result := pingdb()
	w.Write([]byte(fmt.Sprintf(`{"healthy": %v}`, result)))
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json;charset=utf-8")
	ds := h.db.AllDevelopers()
	json.NewEncoder(w).Encode(developerList{Developers: ds})
}

type developer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type developerList struct {
	Developers []developer `json:"developers"`
}

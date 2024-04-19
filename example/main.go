package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	codec "github.com/andyfusniak/scs-json-codec"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Initialize a new session manager and configure the session lifetime.
	sessionManager := scs.New()

	// Open a SQLite3 database.
	db, err := sql.Open("sqlite3", "sessions.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialize a new session manager and configure it to use sqlite3store as the session store.
	sessionManager = scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Codec = codec.JSONCodec{}

	mux := http.NewServeMux()
	mux.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		// Store a new key and value in the session data.
		sessionManager.Put(r.Context(), "message", "Hello from a session!")
	})

	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		// Use the GetString helper to retrieve the string value associated with a
		// key. The zero value is returned if the key does not exist.
		msg := sessionManager.GetString(r.Context(), "message")
		_, _ = io.WriteString(w, msg)
	})

	// Wrap your handlers with the LoadAndSave() middleware.
	if err := http.ListenAndServe(":4000", sessionManager.LoadAndSave(mux)); err != nil {
		return err
	}

	return nil
}

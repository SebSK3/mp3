package main

import (
	"net/http"

	"github.com/sebsk/mp3/assets"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /list", app.list)
	mux.HandleFunc("GET /status/{id}", app.status)
	// mux.HandleFunc("GET /list/{page}", app.listPage)
	mux.HandleFunc("POST /download", app.download)

	return app.logAccess(app.recoverPanic(app.securityHeaders(mux)))
}

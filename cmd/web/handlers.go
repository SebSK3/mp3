package main

import (
	"net/http"
	"strconv"

	"github.com/sebsk/mp3/internal/models"
	"github.com/sebsk/mp3/internal/request"
	"github.com/sebsk/mp3/internal/response"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	err := response.Page(w, http.StatusOK, data, "pages/home.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.serverError(w, r, err)
	}
	song, err := app.db.Progress(id)
	if err != nil {
		app.serverError(w, r, err)
	}
	err = response.JSON(w, http.StatusOK, map[string]int64{"progress": song.Progress})
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) list(w http.ResponseWriter, r *http.Request) {
	var err error

	data := app.newTemplateData(r)
	data["Mp3"], err = app.db.GetList(10, 0)
	if err != nil {
		app.serverError(w, r, err)
	}

	err = response.Page(w, http.StatusOK, data, "pages/list.tmpl")
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) download(w http.ResponseWriter, r *http.Request) {
	var song internal.Song
	err := request.DecodePostForm(r, &song)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	app.logger.Debug("Url", "url", song.Url)
	song.Id, err = app.db.NewSong(song)
	app.backgroundTask(r, func() error {
		return executeDownload(app, &song)
	})

	err = response.JSON(w, http.StatusOK, map[string]int64{"id": song.Id})
	if err != nil {
		app.serverError(w, r, err)
	}
}

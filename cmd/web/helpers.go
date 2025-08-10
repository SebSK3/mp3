package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
	"go.senan.xyz/taglib"

	"github.com/sebsk/mp3/internal/env"
	"github.com/sebsk/mp3/internal/models"
	"github.com/sebsk/mp3/internal/version"
)

func (app *application) newTemplateData(r *http.Request) map[string]any {
	data := map[string]any{
		"Version": version.Get(),
	}

	return data
}

func (app *application) backgroundTask(r *http.Request, fn func() error) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			pv := recover()
			if pv != nil {
				app.reportServerError(r, fmt.Errorf("%v", pv))
			}
		}()

		err := fn()
		if err != nil {
			app.reportServerError(r, err)
		}
	}()
}

func reportProgress(app *application, song *internal.Song, progress ytdlp.ProgressUpdate) {
	app.logger.Debug(fmt.Sprintf("%s: %s", progress.Filename, progress.PercentString()))
	song.Progress = int64(progress.Percent())
	app.db.UpdateSong(song)
}

func songFilepath(song *internal.Song) string {
	return filepath.Join(env.GetString("DOWNLOAD_DIR", "./"), song.Artist, "/", song.Album+"/")
}

func executeDownload(app *application, song *internal.Song) error {
	dl := ytdlp.New().Print("title")
	result, err := dl.Run(context.TODO(), song.Url)
	if err != nil {
		return err
	}
	song.Fullname = result.Stdout

	if song.Mix {
		song.Artist = "Mix"
		// TODO: better way of handling mixes?
		// 	dl := ytdlp.New().Print("channel")
		// 	result, err := dl.Run(context.TODO(), song.Url)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	song.Artist = result.Stdout
		// 	app.logger.Debug("[MIX] Artist:", "artist", song.Artist)
	}
	song.Filename = song.Fullname + ".opus"
	app.logger.Debug("Filename:", "filename", song.Filename)

	extractDataFromFullname(app, song)

	// Download
	dl = ytdlp.New().
		FormatSort("acodec").
		ExtractAudio().
		AudioFormat("opus").
		EmbedThumbnail().
		ConvertThumbnails("jpg").
		ProgressFunc(100*time.Millisecond, func(progress ytdlp.ProgressUpdate) {
			reportProgress(app, song, progress)
		}).
		Output(filepath.Join(songFilepath(song), "%(title)s.%(ext)s"))

	_, err = dl.Run(context.TODO(), song.Url)
	if err != nil {
		return err
	}
	setTags(app, song)
	return nil
}

func extractDataFromFullname(app *application, song *internal.Song) {
	if song.Artist == "" && song.Title == "" {
		// Try to extract from title
		re := regexp.MustCompile(`(.*)(\s-\s)([\w ',]*)(\()?`)
		matches := re.FindStringSubmatch(song.Fullname)
		if len(matches) >= 3 {
			if song.Artist == "" {
				song.Artist = extractArtist(strings.TrimSpace(matches[1]))
			}
			song.Title = strings.TrimSpace(matches[3])
		}
	}
	if song.Album == "" {
		song.Album = "Single"
	}
	app.logger.Debug("Final data:", "artist", song.Artist, "title", song.Title, "album", song.Album)
}

// Workaround for golang's regex greedy matching
func extractArtist(artist string) string {
	re := regexp.MustCompile(`(?i)(\w+.*)(ft\.?|feat\.?|&)`)
	matches := re.FindStringSubmatch(artist)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	} else {
		return artist
	}
}

func setTags(app *application, song *internal.Song) {
	app.logger.Info("Setting tags:", "title", song.Title, "artist", song.Artist, "album", song.Album)

	_ = taglib.WriteTags(filepath.Join(songFilepath(song), song.Filename), map[string][]string{
		taglib.Title:  {song.Title},
		taglib.Artist: {song.Artist},
		taglib.Album:  {song.Album},
	}, taglib.Clear)
	app.db.UpdateSong(song)
}

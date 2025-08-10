package database

import (
	"context"
	"time"

	"github.com/sebsk/mp3/internal/models"

	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

const defaultTimeout = 3 * time.Second

type DB struct {
	*sqlx.DB
}

func New(dsn string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "sqlite", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)
	schema := `
	CREATE TABLE IF NOT EXISTS songs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		artist TEXT,
		album TEXT,
		mix BOOLEAN,
		url TEXT,
		fullname TEXT,
		filename TEXT,
		progress INTEGER
	);`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) GetList(limit int, offset int) ([]internal.Song, error) {
	query := `
		SELECT id, title, artist, album, mix, url, fullname, filename, progress
		FROM songs
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`

	var songs []internal.Song
	err := db.Select(&songs, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func (db *DB) Progress(id int) (*internal.Song, error) {
	query := `
		SELECT id, title, artist, album, mix, url, fullname, filename, progress
		FROM songs
		WHERE ID = ?
	`

	var song internal.Song
	err := db.Get(&song, query, id)

	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (db *DB) NewSong(song internal.Song) (int64, error) {
	query := `
		INSERT INTO
			songs (title, artist, album, mix, url, fullname, filename, progress)
			VALUES (:title, :artist, :album, :mix, :url, :fullname, :filename, :progress)
	`
	res, err := db.NamedExec(query, map[string]any{
		"title":    song.Title,
		"artist":   song.Artist,
		"album":    song.Album,
		"mix":      song.Mix,
		"url":      song.Url,
		"fullname": song.Fullname,
		"filename": song.Filename,
		"progress": 0,
	})
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}

// func (db *DB) UpdateSongProgress(song *internal.Song, progress int64) error {
// 	_, err := db.Exec("UPDATE songs SET progress=$1 WHERE id = $2", progress, song.Id)
// 	return err
// }

func (db *DB) UpdateSong(song *internal.Song) error {
	query := `
		UPDATE songs
			SET title = :title,
			artist = :artist,
			album = :album,
			mix = :mix,
			url = :url,
			fullname = :fullname,
			filename = :filename,
			progress = :progress			
		WHERE id = :id
	`
	_, err := db.NamedExec(query, map[string]any{
		"id":       song.Id,
		"title":    song.Title,
		"artist":   song.Artist,
		"album":    song.Album,
		"mix":      song.Mix,
		"url":      song.Url,
		"fullname": song.Fullname,
		"filename": song.Filename,
		"progress": song.Progress,
	})
	return err
}

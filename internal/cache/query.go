package cache

import (
	"fmt"
	playerStruct "lrcsnc/internal/player/struct"
	"strings"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

const (
	_DB_TABLE_NAME = "lrcsnc_cache"
)

func createTableQuery() (string, []any) {
	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(_DB_TABLE_NAME).IfNotExists()
	ctb.Define("title", "TEXT", "NOT NULL")
	ctb.Define("artists", "TEXT", "NOT NULL")
	ctb.Define("album", "TEXT", "NOT NULL")
	ctb.Define("duration", "REAL", "NOT NULL")
	ctb.Define("lyrics", "TEXT")
	ctb.Define("state", "TINYINT", "NOT NULL")
	ctb.Define("updated_at", "DATETIME", "NOT NULL")

	return ctb.BuildWithFlavor(sqlbuilder.SQLite)
}

func createIndexQuery() (string, []any) {
	return fmt.Sprintf(`CREATE INDEX IF NOT EXISTS lrcsnc_cache_idx ON %v (title, artists, album, duration)`, _DB_TABLE_NAME), []any{}
}

func selectQuery(song *playerStruct.Song) (string, []any) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("lyrics", "state", "updated_at").From(_DB_TABLE_NAME)
	sb.Where(sb.EQ("title", song.Metadata.Title))
	sb.Where(sb.EQ("artists", strings.Join(song.Metadata.Artists, ", ")))
	sb.Where(sb.EQ("album", song.Metadata.Album))
	sb.Where(sb.EQ("duration", song.Metadata.Duration))
	sb.Limit(1)

	return sb.BuildWithFlavor(sqlbuilder.SQLite)
}

func insertQuery(song *playerStruct.Song) (string, []any) {
	ib := sqlbuilder.NewInsertBuilder()
	ib.ReplaceInto(_DB_TABLE_NAME)
	ib.Cols("title", "artists", "album", "duration", "lyrics", "state", "updated_at")
	ib.Values(song.Metadata.Title, strings.Join(song.Metadata.Artists, ", "), song.Metadata.Album, song.Metadata.Duration, song.LyricsData.Lyrics.ToJSON(), byte(song.LyricsData.LyricsState), time.Now())

	return ib.BuildWithFlavor(sqlbuilder.SQLite)
}

func removeQuery(song *playerStruct.Song) (string, []any) {
	db := sqlbuilder.NewDeleteBuilder()
	db.DeleteFrom(_DB_TABLE_NAME)
	db.Where(db.EQ("title", song.Metadata.Title))
	db.Where(db.EQ("artists", strings.Join(song.Metadata.Artists, ", ")))
	db.Where(db.EQ("album", song.Metadata.Album))
	db.Where(db.EQ("duration", song.Metadata.Duration))

	return db.BuildWithFlavor(sqlbuilder.SQLite)
}

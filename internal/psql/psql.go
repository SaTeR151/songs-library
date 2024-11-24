package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jmoiron/sqlx"
	"github.com/sater-151/song-library/internal/models"
	logger "github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("not found")

func init() {
	logger.SetFormatter(&logger.TextFormatter{FullTimestamp: true})
	lvl, ok := os.LookupEnv("LOG_LEVEL")

	if !ok {
		lvl = "debug"
	}

	ll, err := logger.ParseLevel(lvl)
	if err != nil {
		ll = logger.DebugLevel
	}

	logger.SetLevel(ll)
}

type DBStruct struct {
	db *sql.DB
}

func (db *DBStruct) Close() {
	db.db.Close()
}

func Open(config models.DBConfig) (*DBStruct, error) {
	var err error
	connInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Pass,
		config.Dbname,
		config.Port,
		config.Sslmode)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	DB := &DBStruct{db: db}
	return DB, nil
}

func (db *DBStruct) SetDatestyle() error {
	_, err := db.db.Exec("set datestyle = 'German, DMY'")
	if err != nil {
		return err
	}
	return nil
}

func (db *DBStruct) Delete(id string) error {
	res, err := db.db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return ErrNotFound
	}

	return nil
}

func (db *DBStruct) Insert(song models.Song) (string, error) {
	logger.Debug("data addition")
	var id string
	db.db.QueryRow("INSERT INTO songs (song, name_group, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		song.Song, song.Group, song.ReleaseDate, song.Lyric, song.Link).Scan(&id)

	logger.Debug("data added")
	return id, nil
}

func (db *DBStruct) SelectSong(cfg models.SelectConfig) ([]models.Song, error) {
	logger.Debug("sql string configuration")
	var listSong []models.Song
	var sep bool
	row := fmt.Sprintf("SELECT * FROM %s", cfg.Table)
	if cfg.Where {
		row += " WHERE"
	}
	if cfg.Group != "" {
		row += fmt.Sprintf(" name_group = %s", "'"+cfg.Group+"'")
		sep = true
	}
	if cfg.Song != "" {
		if sep {
			row += " and"
		}
		row += fmt.Sprintf(" song = %s", "'"+cfg.Song+"'")
	}
	if cfg.Date != "" {
		if sep {
			row += " and"
		}
		t := strings.Split(cfg.Date, ".")
		t[0], t[2] = t[2], t[0]
		cfg.Date = strings.Join(t, ".")
		row += fmt.Sprintf(" release_date = %s", "'"+cfg.Date+"'")
	}
	if cfg.Lyric != "" {
		if sep {
			row += " and"
		}
		row += fmt.Sprintf(" text like %s", "'%"+cfg.Lyric+"%'")
	}
	if cfg.Id != "" {
		if sep {
			row += " and"
		}
		row += fmt.Sprintf(" id = %s", cfg.Id)
	}
	if cfg.Sort != "" {
		if cfg.Sort == "group" {
			cfg.Sort = "name_group"
		}
		row += fmt.Sprintf(" ORDER BY %s %s", cfg.Sort, cfg.TypeSort)
	}
	if cfg.Limit != "" {
		row += fmt.Sprintf(" LIMIT %s", cfg.Limit)
	}
	if cfg.Offset != "" {
		row += fmt.Sprintf(" OFFSET %s", cfg.Offset)
	}
	logger.Debug("database query")
	res, err := db.db.Query(row)
	if err != nil {
		return listSong, err
	}
	logger.Debug("parseing the received data")
	for res.Next() {
		song := models.Song{}
		err = res.Scan(&song.Id, &song.Song, &song.Group, &song.ReleaseDate, &song.Lyric, &song.Link)
		t := strings.Split(song.ReleaseDate, "T")
		t = strings.Split(t[0], "-")
		t[0], t[2] = t[2], t[0]
		song.ReleaseDate = strings.Join(t, ".")
		if err != nil {
			return listSong, err
		}
		listSong = append(listSong, song)
	}
	return listSong, nil
}

func (db *DBStruct) SelectSongInfo(cfg models.SelectConfig) ([]models.SongInfo, error) {
	logger.Debug("database query is executed ")
	var songInfo []models.SongInfo
	res, err := db.db.Query("SELECT release_date, text, link FROM songs WHERE song = $1 and name_group = $2", cfg.Song, cfg.Group)
	if err != nil {
		return songInfo, err
	}
	for res.Next() {
		song := models.SongInfo{}
		err = res.Scan(&song.ReleaseDate, &song.Text, &song.Link)
		t := strings.Split(song.ReleaseDate, "T")
		t = strings.Split(t[0], "-")
		t[0], t[2] = t[2], t[0]
		song.ReleaseDate = strings.Join(t, ".")
		if err != nil {
			return songInfo, err
		}
		songInfo = append(songInfo, song)
	}
	logger.Debug("request fulfilled")
	return songInfo, nil
}

func (db *DBStruct) UpdateSongs(song models.Song) (string, error) {
	logger.Debug("data update")
	res, err := db.db.Exec("UPDATE songs SET song = $1, name_group=$2, release_date=$3, text=$4, link=$5 WHERE id = $6",
		song.Song, song.Group, song.ReleaseDate, song.Lyric, song.Link, song.Id)
	if err != nil {
		return "", err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return "", err
	}
	if rows == 0 {
		return "", fmt.Errorf("data update error")
	}
	logger.Debug("data has been updated")
	return song.Id, nil
}

func (db *DBStruct) Migration() error {
	driver, err := postgres.WithInstance(db.db, &postgres.Config{})
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return err
	}
	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

package service

import (
	"os"

	"github.com/sater-151/song-library/internal/models"
	"github.com/sater-151/song-library/internal/pkg/deezer"
	"github.com/sater-151/song-library/internal/pkg/lyric_api"
	"github.com/sater-151/song-library/internal/psql"
	logger "github.com/sirupsen/logrus"
)

type ServiceStruct struct {
	Db     *psql.DBStruct
	Client *deezer.DeezerClient
	Liryc  *lyric_api.LyricStruct
}

func New(db *psql.DBStruct) *ServiceStruct {
	service := &ServiceStruct{Db: db}
	service.Client = deezer.New()
	service.Liryc = lyric_api.New()
	return service
}

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

func (s *ServiceStruct) GetSongDetail(postSong models.PostSongJSON) (models.Song, error) {
	var song models.Song
	logger.Debug("getting song id")
	id, err := s.Client.GetID(postSong)
	if err != nil {
		return song, err
	}
	logger.Debug("getting song info")
	song, err = s.Client.GetSongInfo(id)
	if err != nil {
		return song, err
	}
	logger.Debug("getting song lyric")
	song.Lyric, err = s.Liryc.GetLyric(postSong)
	if err != nil {
		return song, err
	}
	song.Song = postSong.Song
	song.Group = postSong.Group

	return song, nil
}

func (s *ServiceStruct) AddSong(postSong models.PostSongJSON) ([]models.Song, error) {
	var selectConfig models.SelectConfig
	var song []models.Song
	songInsert, err := s.GetSongDetail(postSong)
	if err != nil {
		return song, err
	}
	id, err := s.Db.Insert(songInsert)
	if err != nil {
		return song, err
	}
	selectConfig.Table = "songs"
	selectConfig.Id = id
	selectConfig.Where = true
	song, err = s.Db.SelectSong(selectConfig)
	if err != nil {
		return song, err
	}
	return song, nil
}

func (s *ServiceStruct) UpdateSong(postSong models.PostSongJSON, id string) ([]models.Song, error) {
	var selectConfig models.SelectConfig
	var song []models.Song
	songUpdate, err := s.GetSongDetail(postSong)
	if err != nil {
		return song, err
	}
	songUpdate.Id = id
	_, err = s.Db.UpdateSongs(songUpdate)
	if err != nil {
		return song, err
	}
	selectConfig.Table = "songs"
	selectConfig.Id = id
	selectConfig.Where = true
	song, err = s.Db.SelectSong(selectConfig)
	if err != nil {
		return song, err
	}
	return song, nil
}

package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	_ "github.com/jackc/pgx"
	"github.com/sater-151/song-library/internal/models"
	"github.com/sater-151/song-library/internal/pkg/deezer"
	"github.com/sater-151/song-library/internal/psql"
	"github.com/sater-151/song-library/internal/service"
	logger "github.com/sirupsen/logrus"
)

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

func ErrorHandler(res http.ResponseWriter, err error, status int) {
	var errJS models.ErrorJSON
	errJS.Err = err.Error()
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(errJS)
}

// @Summary		GetSongInfo
// @Tags Get
// @Description	get song info
// @Accept			json
// @Produce			json
// @Param			song	query	string	true	"song name"
// @Param			group	query	string	true	"group name"
// @Succes			200 {object} models.SongInfo
// @Failuere		400 {string} string error
// @Failuere		404 {string} string error
// @Failuere		500 {string} string error
// @Router			/info [get]
func GetSongInfo(DB *psql.DBStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("GetSongInfo request received")
		res.Header().Set("Content-type", "application/json; charset=UTF-8")
		song := req.FormValue("song")
		if song == "" {
			logger.Error("song required")
			http.Error(res, "name song required", http.StatusBadRequest)
			return
		}
		group := req.FormValue("group")
		if group == "" {
			logger.Error("group required")
			http.Error(res, "name group required", http.StatusBadRequest)
			return
		}
		var selectConfig models.SelectConfig
		selectConfig.Table = "songs"
		selectConfig.Group = group
		selectConfig.Song = song
		logger.Debug("selectConfig is configured")
		getSongJSON, err := DB.SelectSongInfo(selectConfig)
		if err != nil {
			logger.Error(err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(getSongJSON) == 0 {
			http.Error(res, "song not found", http.StatusNotFound)
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(getSongJSON[0])
		logger.Info("GetSongInfo processed")
	}
}

//		@Summary		GetListSong
//		@Tags Get
//		@Description	get list song
//		@Produce		json
//	 	@Param 			limit 	query 	string 	false 	"limit songs in body"
//		@Param			offset	query	string	false	"Offsetting the song list"
//		@Param			song	query	string	false	"song name"
//		@Param			group	query	string	false	"group name"
//		@Param			release_date	query	string	false	"song release date "
//		@Param			lyric	query	string	false	"song lyrics"
//		@Param			sort	query	string	false	"field sorting"
//		@Param			type_sort	query	string	false	"sorting type"
//		@Succes			200 {object} models.Song
//		@Failuere		404 {string} string error
//		@Failuere		500 {string} string error
//		@Router			/songs [get]
func GetListSong(DB *psql.DBStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("GetListSong request received")
		res.Header().Set("Content-type", "application/json; charset=UTF-8")
		var selectConfig models.SelectConfig
		selectConfig.Table = "songs"

		limit := req.FormValue("limit")
		offset := req.FormValue("offset")
		selectConfig.Limit = limit
		selectConfig.Offset = offset

		song := req.FormValue("song")
		if song != "" {
			selectConfig.Song = song
			selectConfig.Where = true
		}
		group := req.FormValue("group")
		if group != "" {
			selectConfig.Group = group
			selectConfig.Where = true
		}
		releaseDate := req.FormValue("release_date")
		if releaseDate != "" {
			selectConfig.Date = releaseDate
			selectConfig.Where = true
		}
		lyric := req.FormValue("lyric")
		if lyric != "" {
			selectConfig.Lyric = lyric
			selectConfig.Where = true
		}
		sort := req.FormValue("sort")
		if sort != "" {
			selectConfig.Sort = sort
		}
		typeSort := req.FormValue("type_sort")
		if typeSort != "" {
			selectConfig.TypeSort = typeSort
		}

		SongsJSON, err := DB.SelectSong(selectConfig)

		if err != nil {
			logger.Error(err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(SongsJSON) == 0 {
			http.Error(res, "songs not found", http.StatusNotFound)
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(SongsJSON)
		logger.Info("GetListSong processed")
	}
}

//		@Summary		GetSong
//		@Tags Get
//		@Description	get song
//		@Produce		json
//	 	@Param 			id 	query 	string 	true 	"song id"
//		@Succes			200 {object} models.Song
//		@Failuere		400 {string} string error
//		@Failuere		404 {string} string error
//		@Failuere		500 {string} string error
//		@Router			/song [get]
func GetSong(DB *psql.DBStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("GetSong request received")
		res.Header().Set("Content-type", "application/json; charset=UTF-8")
		id := req.FormValue("id")
		if id == "" {
			logger.Error("id required")
			http.Error(res, "id required", http.StatusBadRequest)
			return
		}
		var selectConfig models.SelectConfig
		selectConfig.Table = "songs"
		selectConfig.Id = id
		selectConfig.Where = true
		listSong, err := DB.SelectSong(selectConfig)
		if err != nil {
			logger.Error(err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(listSong) == 0 {
			http.Error(res, "song not found", http.StatusNotFound)
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(listSong[0])
		logger.Info("GetSong processed")
	}
}

// @Summary		PutSong
// @Tags Put
// @Description	Song update
// @Accept			json
// @Produce			json
// @Param			id	query	string	true	"song id"
// @Succes			200 {object} models.SongInfo
// @Failuere		400 {string} string error
// @Failuere		404 {string} string error
// @Failuere		500 {string} string error
// @Router			/song [put]
func PutSong(s *service.ServiceStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("PutSong request received")
		id := req.FormValue("id")
		if id == "" {
			logger.Error("id required")
			http.Error(res, "id required", http.StatusBadRequest)
			return
		}
		var song models.PostSongJSON
		var buf bytes.Buffer
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			logger.Error(err.Error())
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(buf.Bytes(), &song)
		if err != nil {
			logger.Error(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		songs, err := s.UpdateSong(song, id)
		if err != nil {
			logger.Error(err.Error())
			if err == deezer.ErrBadGateway {
				http.Error(res, err.Error(), http.StatusBadGateway)
				return
			}
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(songs) == 0 {
			http.Error(res, "saving error", http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(songs[0])
		logger.Info("PutSong processed")
	}
}

// @Summary		DeleteSong
// @Tags Delete
// @Description	Delete song
// @Param			id	query	string	true	"song id"
// @Succes			204
// @Failuere		400 {string} string error
// @Failuere		404 {string} string error
// @Failuere		500 {string} string error
// @Router			/song [delete]
func DeleteSong(db *psql.DBStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("DeleteSong request received")
		id := req.FormValue("id")
		if id == "" {
			logger.Error("id required")
			http.Error(res, "id required", http.StatusBadRequest)
			return
		}
		err := db.Delete(id)
		if err != nil {
			if err == psql.ErrNotFound {
				http.Error(res, err.Error(), http.StatusNotFound)
			}
			logger.Error(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusNoContent)
		logger.Info("DeleteSong processed")
	}
}

// @Summary		PostSong
// @Tags Post
// @Description	Adding song
// @Accept			json
// @Produce			json
// @Succes			201
// @Failuere		400 {string} string error
// @Failuere		404 {string} string error
// @Failuere		500 {string} string error
// @Failuere		502 {string} string error
// @Router			/song [post]
func PostSong(s *service.ServiceStruct) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger.Info("PostSong request received")
		var song models.PostSongJSON
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			logger.Error(err)
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &song)
		if err != nil {
			logger.Error(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		songs, err := s.AddSong(song)
		if err != nil {
			logger.Error(err.Error())
			if err == deezer.ErrBadGateway {
				http.Error(res, err.Error(), http.StatusBadGateway)
				return
			}
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(songs) == 0 {
			http.Error(res, "saving error", http.StatusNotFound)
			return
		}
		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(songs[0])
		logger.Info("PostSong processed")
	}
}

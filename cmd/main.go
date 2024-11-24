package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/sater-151/song-library/docs"
	cfg "github.com/sater-151/song-library/internal/config"
	"github.com/sater-151/song-library/internal/handlers"
	"github.com/sater-151/song-library/internal/psql"
	"github.com/sater-151/song-library/internal/service"
	logger "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Song Library API
//	@version		1.0
//	@description	API server for Song Library Application

//	@host		localhost:8080
//	@BasePath	/

// Logger Configuration
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

func main() {
	logger.Info(".env download")
	err := godotenv.Load()
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("configuration generation")
	serverConfig := cfg.GetServerConfig()
	dbConfig := cfg.GetDBConfig()

	logger.Info("database connection")
	db, err := psql.Open(dbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Database connect error: %s\n", err.Error()))
		return
	}
	defer db.Close()
	logger.Info("database connected")
	err = db.SetDatestyle()
	if err != nil {
		logger.Error(fmt.Sprintf("Database connect error: %s\n", err.Error()))
		return
	}

	logger.Info("start migration")
	err = db.Migration()

	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("migration done")

	service := service.New(db)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	r.Get("/info", handlers.GetSongInfo(db))
	r.Get("/song", handlers.GetSong(db))
	r.Get("/songs", handlers.GetListSong(db))
	r.Delete("/song", handlers.DeleteSong(db))
	r.Post("/song", handlers.PostSong(service))
	r.Put("/song", handlers.PutSong(service))

	logger.Info(fmt.Sprintf("server start at port: %s\n", serverConfig.Port))
	if err := http.ListenAndServe(":"+serverConfig.Port, r); err != nil {
		logger.Error(fmt.Sprintf("Server startup error: %s\n", err.Error()))
		return
	}
}

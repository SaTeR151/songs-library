package config

import (
	"os"

	"github.com/sater-151/song-library/internal/models"
)

func GetServerConfig() models.ServerConfig {
	var config models.ServerConfig
	config.Port = os.Getenv("SERVER_PORT")
	return config
}

func GetDBConfig() models.DBConfig {
	var dbConfig models.DBConfig
	dbConfig.User = os.Getenv("POSTGRES_USER")
	dbConfig.Pass = os.Getenv("POSTGRES_PASSWORD")
	dbConfig.Dbname = os.Getenv("POSTGRES_DB")
	dbConfig.Sslmode = os.Getenv("SSLMODE")
	dbConfig.Port = os.Getenv("DB_PORT")
	dbConfig.Host = os.Getenv("DB_HOST")
	return dbConfig
}

package models

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	User    string
	Pass    string
	Dbname  string
	Sslmode string
	Port    string
	Host    string
}

type PostSongJSON struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type ErrorJSON struct {
	Err string `json:"error"`
}

type SongInfo struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SelectConfig struct {
	Id       string
	Limit    string
	Offset   string
	Sort     string
	TypeSort string
	Table    string
	Group    string
	Song     string
	Lyric    string
	Date     string
	Where    bool
}

type InsertSongDB struct {
	Group       string
	Name        string
	ReleaseDate string
	Text        string
	Link        string
}

type Song struct {
	Id          string `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Lyric       string `json:"text"`
	Link        string `json:"link"`
}

package lyric_api

import (
	lyrics "github.com/rhnvrm/lyric-api-go"
	"github.com/sater-151/song-library/internal/models"
)

type LyricStruct struct {
	Lyric lyrics.Lyric
}

func New() *LyricStruct {
	l := lyrics.New(lyrics.WithAllProviders())
	lyr := &LyricStruct{Lyric: l}
	return lyr
}

func (l *LyricStruct) GetLyric(song models.PostSongJSON) (string, error) {
	lyric, err := l.Lyric.Search(song.Song, song.Group)
	if err != nil {
		if err.Error() == "Not Found" {
			return "", nil
		} else {
			return "", err
		}
	}
	return lyric, nil
}

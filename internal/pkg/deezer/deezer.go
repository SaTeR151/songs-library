package deezer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sater-151/song-library/internal/models"
)

var ErrBadRequest = errors.New("incorrect request")
var ErrBadGateway = errors.New("bad gateway")

type DeezerClient struct {
	Link string
}

func New() *DeezerClient {
	dc := &DeezerClient{Link: "https://deezerdevs-deezer.p.rapidapi.com"}
	return dc
}

func (dc *DeezerClient) GetID(postSong models.PostSongJSON) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	group, err := url.Parse(postSong.Group)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	song, err := url.Parse(postSong.Song)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/search?q=artist:\"%s\"track:\"%s\"", dc.Link, group, song), nil)

	req.Header.Add("x-rapidapi-key", "dd86916a26mshdaf373ce683e073p180b95jsnad25f54223fd")
	req.Header.Add("x-rapidapi-host", "deezerdevs-deezer.p.rapidapi.com")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusBadRequest {
			return "", ErrBadRequest
		} else {
			return "", fmt.Errorf("external API status code: %v", res.StatusCode)
		}
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return "", err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		return "", err
	}
	if m["error"] != nil {
		msg := m["error"].(map[string]interface{})
		if msg["message"].(string) == "Quota limit exceeded" {
			return "", ErrBadGateway
		}
	}
	if m["message"] != nil {
		return "", errors.New(m["message"].(string))
	}
	ma := m["data"].([]interface{})
	if len(ma) == 0 {
		return "", fmt.Errorf("song not found")
	}

	ab := ma[0].(map[string]interface{})

	id := strconv.Itoa(int(ab["id"].(float64)))
	return id, nil
}

func (dc *DeezerClient) GetSongInfo(id string) (models.Song, error) {
	var song models.Song
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/track/%s", dc.Link, id), nil)
	if err != nil {
		return song, err
	}

	req.Header.Add("x-rapidapi-key", "dd86916a26mshdaf373ce683e073p180b95jsnad25f54223fd")
	req.Header.Add("x-rapidapi-host", "deezerdevs-deezer.p.rapidapi.com")

	res, err := client.Do(req)
	if err != nil {
		return song, err
	}
	defer res.Body.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return song, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		return song, err
	}
	if m["error"] != nil {
		msg := m["error"].(map[string]interface{})
		if msg["message"].(string) == "Quota limit exceeded" {
			return song, ErrBadGateway
		}
	}
	if m["message"] != nil {
		fmt.Println(1)
		return song, errors.New(m["message"].(string))
	}
	song.Link = m["link"].(string)
	song.ReleaseDate = m["release_date"].(string)
	return song, nil
}

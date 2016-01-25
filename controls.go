package main

import (
	"fmt"
	"github.com/HarmonyProject/songster/musicservice"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func play() {
	if playlistSize() == 0 {
		seed()
	}
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		refresh()
	}
}

func enqueue(s musicservice.Song, agent string) {
	addToPlaylist(s, agent)
	UpdateSongdetails(s)
}

func CurrentlyPlaying() musicservice.Song {
	s := getSong(firstSongId())
	return s
}

func refresh() {
	s := CurrentlyPlaying()
	if s.Seek < s.Length {
		updateSeek(s.Id)
		//fmt.Printf("\r%d/%d - %s  ", s.Seek, s.Length, s.Name)
	} else {
		remove(s)
		refresh()
	}
}

func Skip() {
	s := CurrentlyPlaying()
	remove(s)
	refresh()
}

func seed() {
	seedQuery := "tum se hi"
	searchResults := musicservice.Search(seedQuery)
	seedSong := searchResults[0]
	clearPlaylist()
	enqueue(seedSong, "system")
}

func GetPlaylist() []musicservice.Song {
	var playlist []musicservice.Song
	ids := currentPlaylistIds()
	for _, id := range ids {
		playlist = append(playlist, getSong(id))
	}
	return playlist
}

func getLastSong() musicservice.Song {
	s := getSong(lastSongId())
	return s
}

func remove(s musicservice.Song) {
	removeFromPlaylist(s.Id)
}

func getVideoid(youtubeLink string) string {
	u, err := url.Parse(youtubeLink)
	if err != nil {
		fmt.Println("unable to parse URL")
	}
	videoid := u.Query().Get("v")
	return videoid
}

func getQueryResults(query string) []musicservice.Song {
	var songs []musicservice.Song
	if strings.Contains(query, "www.youtube.com/watch?v=") {
		song := musicservice.CreateSong(getVideoid(query))
		if song.Length != -1 {
			songs = append(songs, song)
		}
	} else {
		songs = musicservice.Search(query)
	}
	return songs
}

func UserAdd(query string, user string) bool {
	if strings.Contains(query, "www.youtube.com/watch?v=") {
		song := musicservice.CreateSong(getVideoid(query))
		if song.Length == -1 || user == "" {
			return false
		} else {
			enqueue(song, user)
			return true
		}
	}
	searchResults := musicservice.Search(query)
	if len(searchResults) == 0 {
		return false
	}
	enqueue(searchResults[0], user)
	return true
}

func autoAdd() {
	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		c := CurrentlyPlaying()
		timeRemaining := c.Length - c.Seek
		if playlistSize() == 1 && timeRemaining < 30 {
			newSong := musicservice.Recommend(getLastSong())
			enqueue(newSong, "system")
		}
	}
}

func UpdateLibrary(form url.Values) bool {
	status := false
	var song musicservice.LibSong
	song.Videoid = form.Get("songvideoid")
	song.Artist = form.Get("songartist")
	song.Track = form.Get("songtrack")
	song.Rating, _ = strconv.Atoi(form.Get("songrating"))
	if form.Get("songfav") == "0" {
		song.Fav = false
	} else {
		song.Fav = true
	}

	var user musicservice.User
	user.Name = form.Get("username")
	user.Id = form.Get("userid")

	operation := form.Get("operation")

	if operation == "add" {
		status = addToLibrary(song, user)
	} else if operation == "remove" {
		status = removeFromLibrary(song, user)
	}

	return status
}

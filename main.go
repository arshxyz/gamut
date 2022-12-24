package main

import (
	"context"
	"fmt"
	"image/color"
	"log"

	"github.com/arshxyz/gamut/authenticate"
	"github.com/arshxyz/gamut/colorestimation"
	"github.com/arshxyz/gamut/utils"
	"github.com/zmb3/spotify/v2"
)

type coloredAlbum struct {
	Name        string
	ColorString string
	colorVal    color.RGBA
}

type songInfo struct {
	Name string
	ID   spotify.ID
}

// Maps playlist ID to color string
type PlaylistIDMap map[string]spotify.ID

// Maps colour string to playlist name
var ColouredPlaylists = map[string]string{
	"red":         "❤️❤️ Liked ❤️❤️",
	"green":       "💚💚 Liked 💚💚",
	"blue":        "💙💙 Liked 💙💙",
	"yellow":      "💛💛 Liked 💛💛",
	"pink/purple": "💜💜 Liked 💜💜",
	"bw":          "🤍🖤 Liked 🤍🖤",
	"orange":      "🧡🧡 Liked 🧡🧡",
}

// Finds promiment colour for all albums of liked songs
// and adds all liked songs in their corresponding playlists
func classify(client *spotify.Client, playlistID PlaylistIDMap) {
	count := 1
	tracks, err := client.CurrentUsersTracks(context.Background())
	total := tracks.Total
	if err != nil {
		log.Fatalln(err)
	}
	// Maps Album ID to Name and promiment colour
	albumColor := make(map[spotify.ID]coloredAlbum)
	// Maps Album ID to liked tracks in the Album
	albumSongs := make(map[spotify.ID][]songInfo)

	for page := 0; ; page++ {
		if len(tracks.Tracks) == 0 {
			break
		}
		tracks, err = client.CurrentUsersTracks(context.Background(), spotify.Offset(20*page))
		if err != nil {
			break
		}
		// TODO: Fix spaghetti code
		for _, v := range tracks.Tracks {
			var clstring string
			var colorVal color.RGBA
			album := v.Album
			// Only process album cover if not in map
			if _, ok := albumColor[album.ID]; !ok {
				clstring, colorVal = utils.GetAlbumColour(v.Album)
				albumColor[album.ID] = coloredAlbum{
					Name:        album.Name,
					ColorString: clstring,
					colorVal:    colorVal,
				}
			}
			// Add song to Album ID
			albumSongs[album.ID] = append(albumSongs[album.ID], songInfo{Name: v.Name, ID: v.ID})
			fmt.Printf("\033[2K\r%s", fmt.Sprint("Processing track: ", count, " of ", total))
			count++
		}

	}
	// Start count again for adding tracks to albums after processing
	count = 1
	for albumID, songs := range albumSongs {
		closest := colorestimation.FindClosest(albumColor[albumID].colorVal)
		songIDs := make([]spotify.ID, 0, len(songs))
		for _, song := range songs {
			songIDs = append(songIDs, song.ID)
			fmt.Printf("\033[2K%s\r", fmt.Sprint("Adding track: ", count, " of ", total))
			count++
		}

		_, err := client.AddTracksToPlaylist(context.Background(), playlistID[closest], songIDs...)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// Create playlists for all colours
func createPlaylists(client *spotify.Client, playlistID PlaylistIDMap) {
	user, _ := client.CurrentUser(context.Background())
	fmt.Println("Logged in as", user.DisplayName)
	userID := user.ID
	count := 1
	for pcolor, pname := range ColouredPlaylists {
		fmt.Printf("\033[2KCreating Playlist %d of %d: %s\r", count, len(ColouredPlaylists), pname)
		playlist, err := client.CreatePlaylistForUser(context.Background(), userID, pname, "Generated with github.com/arshxyz/gamut", false, false)
		if err != nil {
			log.Fatalln(err)
		}
		playlistID[pcolor] = playlist.ID
		count++
	}
	fmt.Println("\033[2KCreated Playlists!\r")
}

func main() {
	client, _ := authenticate.Authenticate()
	playlistID := make(PlaylistIDMap)
	createPlaylists(client, playlistID)
	classify(client, playlistID)
}

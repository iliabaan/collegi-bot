package command

import (
	"context"
	"fmt"
	"github.com/ndrewnee/go-yamusic/yamusic"
	"github.com/zmb3/spotify/v2"
	"os"
	"strconv"
)

func AddSpotify(title string, client spotify.Client) spotify.FullTrack {
	// Search for the track
	ctx := context.Background()
	searchResult, err := client.Search(ctx, title, spotify.SearchTypeTrack)
	if err != nil {
		panic(err)
	}

	// Get the first track
	track := searchResult.Tracks.Tracks[0]

	fmt.Println("Found track:", track.Name)
	// Add the track to the playlist
	playlistID := os.Getenv("SPOTIFY_PLAYLIST_ID")
	_, err = client.AddTracksToPlaylist(ctx, spotify.ID(playlistID), track.ID)

	if err != nil {
		panic(err)
	}

	return track
}

func AddYandex(ctx context.Context, title string, client *yamusic.Client) {
	opts := &yamusic.SearchOptions{
		Page:      0,
		NoCorrect: false,
	}

	// find track to add
	searchResult, _, err := client.Search().Tracks(ctx, title, opts)

	if err != nil {
		panic(err)
	}

	// get playlist where to add
	playlistResp, _ := getPlaylist(client, converterToInt(os.Getenv("YANDEX_PLAYLIST_ID")))
	playlist := playlistResp.Result

	track := searchResult.Result.Tracks.Results[0]

	// create tracks to add
	tracksToAdd := make([]yamusic.PlaylistsTrack, 1)
	tracksToAdd[0] = yamusic.PlaylistsTrack{
		ID:      track.ID,
		AlbumID: track.Albums[0].ID,
	}

	pATO := &yamusic.PlaylistsAddTracksOptions{
		At: 0,
	}

	// add track to playlist
	resp, _, err := client.Playlists().AddTracks(ctx, converterToInt(os.Getenv("YANDEX_PLAYLIST_ID")), playlist.Revision, tracksToAdd, pATO)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Added track: %v", resp)
}

func getPlaylist(client *yamusic.Client, playlistId int) (*yamusic.PlaylistsGetResp, error) {
	ctx := context.Background()
	playlist, _, err := client.Playlists().GetByUserIDAndKind(ctx, os.Getenv("YANDEX_ID"), playlistId)
	if err != nil {
		panic(err)
	}
	return playlist, nil
}

func converterToInt(data string) int {
	integer, _ := strconv.Atoi(data)
	return integer
}

package authentication

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
	"net/http"
	"os"
)

var (
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func Spotify() (*spotify.Client, error) {
	var auth = spotifyauth.New(
		spotifyauth.WithClientID(os.Getenv("SPOTIFY_ID")),
		spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_SECRET")),
		spotifyauth.WithRedirectURL(os.Getenv("APP_URL")+"/callback"),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPrivate, spotifyauth.ScopePlaylistModifyPublic))

	// start an HTTP server
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		completeAuth(w, r, auth)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go func() {
		err := http.ListenAndServe(":"+os.Getenv("EXPOSED_PORT"), nil)
		if err != nil {
			log.Println(err)
		}
	}()

	url := auth.AuthURL(state)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Println(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	return client, nil
}

func completeAuth(w http.ResponseWriter, r *http.Request, auth *spotifyauth.Authenticator) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Printf("State mismatch: %s != %s\n", st, state)
		return
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	_, err = fmt.Fprintf(w, "Login Completed!")
	if err != nil {
		return
	}
	ch <- client
}

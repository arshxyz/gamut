package authenticate

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arshxyz/gamut/utils"
	"github.com/dghubble/sling"
	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
	"github.com/zmb3/spotify/v2"
	auth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type RefreshRequest struct {
	GrantType    string `url:"grant_type,omitempty"`
	RefreshToken string `url:"refresh_token,omitempty"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

type RefreshError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

var allscopes = []string{auth.ScopeUserReadCurrentlyPlaying, auth.ScopeUserReadPlaybackState, auth.ScopeUserModifyPlaybackState, auth.ScopeUserReadRecentlyPlayed, auth.ScopeStreaming, auth.ScopeUserLibraryRead, auth.ScopeUserLibraryModify, auth.ScopePlaylistModifyPrivate}

var (
	SpotifyAuth = func(clid, cls string) *auth.Authenticator {
		return auth.New(
			auth.WithRedirectURL(redirectURI),
			auth.WithScopes(allscopes...),
			auth.WithClientID(clid),
			auth.WithClientSecret(cls))
	}

	ch          = make(chan *oauth2.Token)
	state       = "spotistate"
	redirectURI = "http://localhost:8888/callback"
)

func Authenticate() (client *spotify.Client, err error) {
	utils.InitViper()
	var clid, cls string
	var tok *oauth2.Token
	var returningUser bool
	// Check if Secrets are present in config file
	sc, err := utils.GetSecrets()
	if err != nil {
		log.Println(err.Error())
		returningUser = false
	} else {
		returningUser = true
	}
	// Secrets not found, prompt user to enter ClientID and ClientSecret
	if !returningUser {
		clidPrompt := promptui.Prompt{
			Label: "Client ID",
		}
		clid, err = clidPrompt.Run()
		if err != nil {
			log.Fatalln(err)
		}
		clsPrompt := promptui.Prompt{
			Label: "Client Secret",
		}
		cls, err = clsPrompt.Run()
		if err != nil {
			log.Fatalln(err)
		}
		tok = getFirstToken(clid, cls)

		err = utils.WriteSecrets(utils.Secrets{ClientID: clid, ClientSecret: cls, RefreshToken: tok.RefreshToken})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("\033[2K\r%s\n", "Login Successful!")
	} else {
		// Get ClientID and ClientSecret from config
		clid = sc.ClientID
		cls = sc.ClientSecret
		rt := sc.RefreshToken

		tok, err = RefreshToken(rt, clid, cls)
		if err != nil {
			// Refresh token has expired, fetch new token using auth code flow
			fmt.Println("Oauth Token needs to be refreshed, please login again!")
			tok = getFirstToken(clid, cls)
			fmt.Printf("\033[2K\r%s\n", "Login Successful!")
			utils.WriteSecrets(utils.Secrets{ClientID: clid, ClientSecret: cls, RefreshToken: tok.RefreshToken})

		}

	}
	// Start a client
	client = spotify.New(SpotifyAuth(clid, cls).Client(context.Background(), tok), spotify.WithRetry(true))
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Logged in as", user.DisplayName)
	return client, nil
}

// Get a Refresh Token from ClientID and ClientSecret
func getFirstToken(clid, cls string) (t *oauth2.Token) {
	http.HandleFunc("/callback", completeAuth(clid, cls))
	// start a webserver listening on the redirectURI port
	go func() {
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Print("Please log in to Spotify, opening browser...")
	time.Sleep(3 * time.Second)
	// Setting handler functions for callback. This must match the callback path in redirectURI

	url := SpotifyAuth(clid, cls).AuthURL(state)
	browser.OpenURL(url)
	// Trigger on auth complete
	tok := <-ch
	return tok
}

// Get token from /callback
func completeAuth(clid, cls string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, err := SpotifyAuth(clid, cls).Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatal(err)
		}
		// Check for state
		if st := r.FormValue("state"); st != state {
			http.NotFound(w, r)
			log.Fatalf("State mismatch: %s != %s\n", st, state)
		}
		// Notify through browser
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "Login Completed! Close this window")
		ch <- tok
	}
}

// Refresh Token if the previous one has expired. In practice I've only
// had this fire when scopes were changed.
// While Spotify docs meantion a timeout it is probably not enforced
func RefreshToken(rtoken, clientId, clientSecret string) (token *oauth2.Token, err error) {
	authstring := base64.URLEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	body := &RefreshRequest{
		GrantType:    "refresh_token",
		RefreshToken: rtoken,
	}
	resp := &RefreshResponse{}
	respErr := &RefreshError{}

	slingClient := sling.New().
		Base("https://accounts.spotify.com/api/").
		Set("Authorization", "Basic "+authstring)

	req, err := slingClient.Post("token").BodyForm(body).Receive(resp, respErr)
	if err != nil {
		return token, err
	}

	if req.StatusCode != 200 {
		return token, fmt.Errorf("token refresh error: %s", req.Status)
	}
	return &oauth2.Token{
		AccessToken: resp.AccessToken,
	}, nil
}

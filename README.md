## Gamut
Gamut helps you organize your Spotify library by colour.
It splits your "liked" tracks into different playlists by colour.

![Gamut Intro UltraCompressed](https://user-images.githubusercontent.com/23417273/210441169-90ce5492-c7f8-401c-87b9-c1d963335a79.gif)

### Installation
- Make sure you have the [latest version of Go](https://go.dev/dl/)
- Install `gamut` using the `go get` command
```console
go get github.com/arshxyz/gamut
```
- Head over to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/) and create a new project. Note the Client ID and Secret key. Set the callback URL for your project to `http://localhost:8888/callback` by going to the "Edit Settings" page and entering this URL in the "Redirect URIs" field.
- Finally run `gamut` in the terminal. If this is your first time using it, it should automatically prompt for your Client ID and Secret.

## TODO
- [x] Remove unused OAuth scopes
- [x] Use prettier spinners and prompts
- [x] Show color info for each track during processing/adding
- [x] Finetune colour matching
- [ ] Write a proper README

## Maybe
- [ ] Add alternative names for playlists without emojis for Windows
- [ ] Error handling can be better
- [ ] Add paging for AddTrackToPlaylist (handles edge case where 100+ tracks have been liked from the same album)

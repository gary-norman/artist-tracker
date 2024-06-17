package api

import (
	"fmt"
	"strings"
)

// ANSI escape codes for coloring text
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Cyan   = "\033[36m"
	Blue   = "\033[34m"
	Yellow = "\033[33m"
)

// PrintArtistDetails prints detailed information about the artist with colors
func PrintArtistDetails(artist *Artist) {
	fmt.Println()
	fmt.Println(Bold + Yellow + "*********Searching result************" + Reset)
	fmt.Printf(Bold+Cyan+"Artist Name: "+Reset+Blue+"%v\n"+Reset, artist.Name)
	fmt.Printf(Bold+Cyan+"ID: "+Reset+Blue+"%v\n"+Reset, artist.Id)
	fmt.Printf(Bold+Cyan+"Image: "+Reset+Blue+"%v\n"+Reset, artist.Image)
	fmt.Printf(Bold+Cyan+"Members: "+Reset+Blue+"%v\n"+Reset, strings.Join(artist.Members, ", "))
	fmt.Printf(Bold+Cyan+"Creation Date: "+Reset+Blue+"%v\n"+Reset, artist.CreationDate)
	fmt.Printf(Bold+Cyan+"First Album: "+Reset+Blue+"%v\n"+Reset, artist.FirstAlbum)
	fmt.Printf(Bold+Cyan+"Album Name: "+Reset+Blue+"%v\n"+Reset, artist.TadbAlbum.Album)
	fmt.Printf(Bold+Cyan+"Album Image Link: "+Reset+Blue+"%v\n"+Reset, artist.TadbAlbum.AlbumThumb)
	fmt.Printf(Bold+Cyan+"Album Year Released: "+Reset+Blue+"%v\n"+Reset, artist.TadbAlbum.YearReleased)
	fmt.Printf(Bold+Cyan+"Album Description: "+Reset+Blue+"%v\n"+Reset, artist.TadbAlbum.DescriptionEN)
	fmt.Printf(Bold+Cyan+"Concert Dates and Locations: "+Reset+Blue+"%v\n"+Reset, artist.DatesLocations)

	for location, dates := range artist.DatesLocations {
		fmt.Printf(Bold+Cyan+"Location: "+Reset+Blue+"%v, "+Bold+Cyan+"Dates: "+Reset+Blue+"%v\n"+Reset, location, strings.Join(dates, ", "))
	}

	fmt.Printf(Bold+Cyan+"Label: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.Label)
	fmt.Printf(Bold+Cyan+"Genre: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.Genre)
	fmt.Printf(Bold+Cyan+"Website: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.Website)
	fmt.Printf(Bold+Cyan+"Biography: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.BiographyEn)
	fmt.Printf(Bold+Cyan+"Artist Thumb: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistThumb)
	fmt.Printf(Bold+Cyan+"Artist Logo: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistLogo)
	fmt.Printf(Bold+Cyan+"Artist Banner: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistBanner)
	fmt.Printf(Bold+Cyan+"MusicBrainz ID: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.MusicBrainzID)
}

// PrintSuggestionsDetails prints detailed information about each suggestion with colors
func PrintSuggestionsDetails(suggestions []Suggestion) {
	fmt.Println()
	fmt.Println(Bold + Yellow + "********* Suggestions ************" + Reset)
	fmt.Printf("%sThere are %s%d%s result(s) found:\n%s", Blue, Yellow, len(suggestions), Blue, Reset)

	for index, suggestion := range suggestions {
		fmt.Printf(Bold+Cyan+"Type: "+Reset+Blue+"%s\n"+Reset, suggestion.Category)
		fmt.Printf(Bold+Cyan+"Match: "+Reset+Blue+"%s\n"+Reset, suggestion.MatchItem)

		// Print artist details if available
		if suggestion.Artist != nil {
			fmt.Println()
			fmt.Printf(Bold+Yellow+"Matching artist result No.%v info:\n"+Reset, index+1)
			PrintArtistDetails(suggestion.Artist)
		}

		fmt.Println("---------------------------------")
	}
}

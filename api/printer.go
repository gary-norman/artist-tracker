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
	fmt.Println(Bold + Cyan + "*********Searching result************" + Reset)
	fmt.Printf(Bold+Cyan+"Artist Name: "+Reset+Blue+"%v\n"+Reset, artist.Name)
	fmt.Printf(Bold+Cyan+"ID: "+Reset+Blue+"%v\n"+Reset, artist.Id)
	fmt.Printf(Bold+Cyan+"Image: "+Reset+Blue+"%v\n"+Reset, artist.Image)
	fmt.Printf(Bold+Cyan+"MemberList: "+Reset+Blue+"%v\n"+Reset, strings.Join(artist.MemberList, ", "))
	// Member's pictures
	fmt.Println(Bold + Cyan + "------Members details------" + Reset)
	for member, picLink := range artist.Members {
		fmt.Printf(Bold+Cyan+"member: "+Reset+Blue+"%v, "+Bold+Cyan+"Picture Link: "+Reset+Blue+"%v\n"+Reset, member, picLink)
	}
	fmt.Printf(Bold+Cyan+"Creation Date: "+Reset+Blue+"%v\n"+Reset, artist.CreationDate)
	fmt.Printf(Bold+Cyan+"First Album: "+Reset+Blue+"%v\n"+Reset, artist.FirstAlbum)

	// TadbAlbum details
	fmt.Println(Bold + Cyan + "------First Album Details------" + Reset)
	fmt.Printf(Bold+Cyan+"First Album Name: "+Reset+Blue+"%v\n"+Reset, artist.FirstAlbumStruct.Album)
	fmt.Printf(Bold+Cyan+"First Album Image Link: "+Reset+Blue+"%v\n"+Reset, artist.FirstAlbumStruct.AlbumThumb)
	fmt.Printf(Bold+Cyan+"First Album Year Released: "+Reset+Blue+"%v\n"+Reset, artist.FirstAlbumStruct.YearReleased)
	fmt.Printf(Bold+Cyan+"Fisrt Album Description: "+Reset+Blue+"%v\n"+Reset, artist.FirstAlbumStruct.DescriptionEN)
	fmt.Printf(Bold+Cyan+"First MusicBrainz Album ID: "+Reset+Blue+"%v\n"+Reset, artist.FirstAlbumStruct.MusicBrainzAlbumID)

	// TadbAlbum details
	fmt.Println(Bold + Cyan + "------All Album Details------" + Reset)
	for i := range artist.AllAlbums.Album {
		fmt.Printf(Bold+Cyan+"Album Name: "+Reset+Blue+"%v\n"+Reset, artist.AllAlbums.Album[i].Album)
		fmt.Printf(Bold+Cyan+"Album Image Link: "+Reset+Blue+"%v\n"+Reset, artist.AllAlbums.Album[i].AlbumThumb)
		fmt.Printf(Bold+Cyan+"Album Year Released: "+Reset+Blue+"%v\n"+Reset, artist.AllAlbums.Album[i].YearReleased)
		fmt.Printf(Bold+Cyan+"Album Genre: "+Reset+Blue+"%v\n"+Reset, artist.AllAlbums.Album[i].Genre)
		fmt.Printf(Bold+Cyan+"Album Description: "+Reset+Blue+"%v\n"+Reset, artist.AllAlbums.Album[i].DescriptionEN)
		fmt.Printf(Bold+Cyan+"MusicBrainz Album ID: "+Reset+Blue+"%v\n"+Reset, artist.AllAlbums.Album[i].MusicBrainzAlbumID)
	}
	fmt.Printf(Bold+Cyan+" --- %v albums in total:---\n\n"+Reset, len(artist.AllAlbums.Album))

	// Dates and Locations
	fmt.Println(Bold + Cyan + "------Concert Dates and Locations------" + Reset)
	for location, dates := range artist.DatesLocations {
		fmt.Printf(Bold+Cyan+"Location: "+Reset+Blue+"%v, "+Bold+Cyan+"Dates: "+Reset+Blue+"%v\n"+Reset, location, strings.Join(dates, ", "))
	}

	// TourDetails (Concerts)
	fmt.Println(Bold + Cyan + "------ Concert Details ------" + Reset)
	for _, concert := range artist.TourDetails.Data {
		fmt.Printf(Bold+Cyan+"Concert ID: "+Reset+Blue+"%v\n"+Reset, concert.ConcertId)
		//fmt.Printf(Bold+Cyan+"Description: "+Reset+Blue+"%v\n"+Reset, concert.Description)
		fmt.Printf(Bold+Cyan+"Start Date: "+Reset+Blue+"%v\n"+Reset, concert.StartDate)
		//fmt.Printf(Bold+Cyan+"End Date: "+Reset+Blue+"%v\n"+Reset, concert.EndDate)
		//fmt.Printf(Bold+Cyan+"Image: "+Reset+Blue+"%v\n"+Reset, concert.Image)
		fmt.Printf(Bold+Cyan+"Location: "+Reset+Blue+"%v \n"+Reset, concert.Location)
	}

	// TheAudioDbArtist details
	fmt.Println(Bold + Cyan + "------Artist Details------" + Reset)
	fmt.Printf(Bold+Cyan+"Label: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.Label)
	fmt.Printf(Bold+Cyan+"Genre: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.Genre)
	fmt.Printf(Bold+Cyan+"Website: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.Website)
	fmt.Printf(Bold+Cyan+"Biography: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.BiographyEn)
	fmt.Printf(Bold+Cyan+"Artist Thumb: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistThumb)
	fmt.Printf(Bold+Cyan+"Artist Logo: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistLogo)
	fmt.Printf(Bold+Cyan+"Artist Cutout: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistCutout)
	fmt.Printf(Bold+Cyan+"Artist Clearart: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistClearart)
	fmt.Printf(Bold+Cyan+"Artist Widethumb: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistWidethumb)
	fmt.Printf(Bold+Cyan+"Artist Fanart: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistFanart)
	fmt.Printf(Bold+Cyan+"Artist Fanart2: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistFanart2)
	fmt.Printf(Bold+Cyan+"Artist Fanart3: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistFanart3)
	fmt.Printf(Bold+Cyan+"Artist Fanart4: "+Reset+Blue+"%v\n"+Reset, artist.TheAudioDbArtist.ArtistFanart4)
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

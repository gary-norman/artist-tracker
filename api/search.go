package api

import (
	"fmt"
	"strings"
)

// Custom String method for Artist struct to format output
func (a Artist) String() string {
	result := fmt.Sprintf("Id: %d\nImage: %s\nName: %s\nMembers: %v\nCreationDate: %d\nFirstAlbum: %s\n"+
		"Album: %+v\nTheAudioDbArtist:\n  IdArtist: %s\n  Label: %s\n  Genre: %s\n  Website: %s\n  BiographyEn: %s\n  ArtistThumb: %s\n"+
		"  ArtistLogo: %s\n  ArtistCutout: %s\n  ArtistClearart: %s\n  ArtistWidethumb: %s\n  ArtistFanart: %s\n  ArtistFanart2: %s\n"+
		"  ArtistFanart3: %s\n  ArtistFanart4: %s\n  ArtistBanner: %s\n  MusicBrainzID: %s\n",
		a.Id, a.Image, a.Name, a.Members, a.CreationDate, a.FirstAlbum, a.SpotifyAlbum, a.IdArtist, a.Label, a.Genre,
		a.Website, a.BiographyEn, a.ArtistThumb, a.ArtistLogo, a.ArtistCutout, a.ArtistClearart, a.ArtistWidethumb,
		a.ArtistFanart, a.ArtistFanart2, a.ArtistFanart3, a.ArtistFanart4, a.ArtistBanner, a.MusicBrainzID)

	result += "DatesLocations:\n"
	for location, dates := range a.DatesLocations {
		result += fmt.Sprintf("  %s: %v\n", location, dates)
	}
	return result
}

// SearchArtist function searches for an artist by name and returns the artist details
func SearchArtist(artists []Artist, name string) (*Artist, error) {
	for _, artist := range artists {
		if strings.ToLower(artist.Name) == strings.ToLower(name) {
			return &artist, nil
		}
	}
	return nil, fmt.Errorf("artist not found")
}

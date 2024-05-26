package api

import (
	"fmt"
	"strings"
)

// Custom String method for Artist struct to format output
func (a Artist) String() string {
	result := fmt.Sprintf("Id: %d\nImage: %s\nName: %s\nMembers: %v\nCreationDate: %d\nFirstAlbum: %s\n"+
		"TadbAlbum:\n  IdAlbum: %s\n  Album: %s\n  YearReleased: %s\n  AlbumThumb: %s\n  DecriptionEN: %s\n  "+
		"MusicBrainzAlbumID: %s\nTheAudioDbArtist:\n  IdArtist: %s\n  Label: %s\n  Genre: %s\n  Website: %s\n"+
		"  BiographyEn: %s\n  ArtistThumb: %s\n  ArtistLogo: %s\n  ArtistCutout: %s\n  ArtistClearart: %s\n"+
		"  ArtistWidethumb: %s\n  ArtistFanart: %s\n  ArtistFanart2: %s\n  ArtistFanart3: %s\n  ArtistFanart4: %s\n"+
		"  ArtistBanner: %s\n  MusicBrainzID: %s\nTourDetails:\n ConcertId: %v\n",
		a.Id, a.Image, a.Name, a.Members, a.CreationDate, a.FirstAlbum, a.IdAlbum, a.Album, a.YearReleased, a.AlbumThumb,
		a.DescriptionEN, a.MusicBrainzAlbumID, a.IdArtist, a.Label, a.Genre,
		a.Website, a.BiographyEn, a.ArtistThumb, a.ArtistLogo, a.ArtistCutout, a.ArtistClearart, a.ArtistWidethumb,
		a.ArtistFanart, a.ArtistFanart2, a.ArtistFanart3, a.ArtistFanart4, a.ArtistBanner, a.MusicBrainzID, a.TourDetails)

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
			result := &artist
			return result, nil
		}
	}
	return &Artist{}, fmt.Errorf("artist not found")
}

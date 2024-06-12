package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// Custom String method for Artist struct to format output
func (a Artist) String() string {
	result := fmt.Sprintf("Id: %d\nImage: %s\nName: %s\nMembers: %v\nCreationDate: %d\nFirstAlbum: %s\n"+
		"TadbAlbum:\n  IdAlbum: %s\n  Album: %s\n  YearReleased: %s\n  AlbumThumb: %s\n  DecriptionEN: %s\n  "+
		"MusicBrainzAlbumID: %s\nTheAudioDbArtist:\n  IdArtist: %s\n  Label: %s\n  Genre: %s\n  Website: %s\n"+
		"  BiographyEn: %s\n  ArtistThumb: %s\n  ArtistLogo: %s\n  ArtistCutout: %s\n  ArtistClearart: %s\n"+
		"  ArtistWidethumb: %s\n  ArtistFanart: %s\n  ArtistFanart2: %s\n  ArtistFanart3: %s\n  ArtistFanart4: %s\n"+
		"  ArtistBanner: %s\n  MusicBrainzID: %s\n",
		a.Id, a.Image, a.Name, a.Members, a.CreationDate, a.FirstAlbum, a.IdAlbum, a.Album, a.YearReleased, a.AlbumThumb,
		a.DescriptionEN, a.MusicBrainzAlbumID, a.IdArtist, a.Label, a.Genre,
		a.Website, a.BiographyEn, a.ArtistThumb, a.ArtistLogo, a.ArtistCutout, a.ArtistClearart, a.ArtistWidethumb,
		a.ArtistFanart, a.ArtistFanart2, a.ArtistFanart3, a.ArtistFanart4, a.ArtistBanner, a.MusicBrainzID)

	result += "DatesLocations:\n"
	for location, dates := range a.DatesLocations {
		result += fmt.Sprintf("  %s: %v\n", location, dates)
	}
	return result
}

// String method for Address to provide a custom string representation.
func (a Address) String() string {
	return fmt.Sprintf("StreetAddress: %s\n     AddressLocality: %s\n     AddressRegion: %s\n     PostalCode: %s\n     AddressCountry: %s",
		a.StreetAddress, a.AddressLocality, a.AddressRegion, a.PostalCode, a.AddressCountry)
}

// String method for Geo to provide a custom string representation.
func (g Geo) String() string {
	return fmt.Sprintf(" Geo {\n     Type: %s\n     Latitude: %f\n     Longitude: %f\n    }", g.Type, g.Latitude, g.Longitude)
}

// String method for Location to provide a custom string representation.
func (l Location) String() string {
	return fmt.Sprintf("   Location {\n    Name: %s\n    Address {\n     %s\n   %s", l.Name, l.Address, l.Geo)
}

// String method for ConcertData to provide a custom string representation.
func (d ConcertData) String() string {
	return fmt.Sprintf("Data {\n   Concert ID: %s\n   Description: %s\n   Start Date: %s\n   End Date: %s\n   Image: %s\n%s\n}",
		d.ConcertId, d.Description, d.StartDate, d.EndDate, d.Image, d.Location)
}

// String method for TourDetails to provide a custom string representation.
func (td TourDetails) String() string {
	var sb strings.Builder
	for _, data := range td.Data {
		sb.WriteString(data.String() + "\n")
	}
	return sb.String()
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

// display live suggestion when typing
func SuggestHandler(w http.ResponseWriter, r *http.Request, artists []Artist, tpl *template.Template) {
	searchQuery := strings.ToLower(r.FormValue("query"))
	var suggestions []string

	for _, artist := range artists {
		// Check artist name
		if strings.Contains(strings.ToLower(artist.Name), searchQuery) {
			suggestions = append(suggestions, fmt.Sprintf("Artist: %s", artist.Name))
		}

		// Check artist members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), searchQuery) {
				suggestions = append(suggestions, fmt.Sprintf("Member: %s", member))
			}
		}

		// Check artist locations
		for locations, _ := range artist.DatesLocations {
			if strings.Contains(strings.ToLower(locations), searchQuery) {
				suggestions = append(suggestions, fmt.Sprintf("Location: %s", locations))
			}
		}

		// Check artist first album
		if strings.Contains(strings.ToLower(artist.FirstAlbum), searchQuery) {
			suggestions = append(suggestions, fmt.Sprintf("Album: %s", artist.FirstAlbum))
		}
	}

	jsonData, err := json.Marshal(suggestions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

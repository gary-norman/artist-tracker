package api

import (
	"fmt"
	"strings"
)

// Custom String method for Artist struct to format output
func (a Artist) String() string {
	result := fmt.Sprintf("Id: %d\nImage: %s\nName: %s\nMembers: %v\nCreationDate: %d\nFirstAlbum: %s\n",
		//result := fmt.Sprintf("%d\n%s\n%s\n%v\n%d\n%s\n",
		a.Id, a.Image, a.Name, a.Members, a.CreationDate, a.FirstAlbum)

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

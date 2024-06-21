package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type WikiQuery struct {
	Pages map[string]WikiPage `json:"pages"`
}

type WikiPage struct {
	WikiThumbnail WikiThumbnail `json:"thumbnail"`
}

type WikiThumbnail struct {
	Source string `json:"source"`
}

// only for debug, delete later
type WikiResponse struct {
	WikiQuery struct {
		Pages map[string]struct {
			WikiThumbnail struct {
				Source string `json:"source"`
			} `json:"thumbnail"`
		} `json:"pages"`
	} `json:"query"`
}

func FetchAllArtistsImages(artists []Artist) {
	for i := range artists {
		WikiImageFetcher(&artists[i])
	}
}

// get individul arttist's image
func WikiImageFetcher(artist *Artist) {
	artist.Members = make(map[string]string)
	for _, member := range artist.MemberList {
		encodedMember := url.QueryEscape(member)
		queryURL := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&titles=%s&prop=pageimages&format=json&pithumbsize=500", encodedMember)
		resp, err := http.Get(queryURL)
		if err != nil {
			fmt.Println("here")
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		// debug only, delete later
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: received status code %d for member %s\n", resp.StatusCode, member)
			body, _ := io.ReadAll(resp.Body)
			fmt.Println("Response body:", string(body))
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			// return
			continue
		}

		var result WikiResponse
		if err = json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			fmt.Println("Error:", err)
			//return
			continue
		}

		// Add the image URL to the map
		for _, page := range result.WikiQuery.Pages {
			if page.WikiThumbnail.Source != "" {
				artist.Members[member] = page.WikiThumbnail.Source
			} else {
				artist.Members[member] = "/icons/artist_placeholder_08.png"
			}
			//fmt.Println("Main Image URL:", page.WikiThumbnail.Source)
		}

		// Add the image URL to the struct
		for _, page := range result.WikiQuery.Pages {
			memberName := strings.Replace(member, "_", " ", -1)
			if page.WikiThumbnail.Source != "" {
				artist.MemberStruct = append(artist.MemberStruct, Member{memberName, page.WikiThumbnail.Source})
			} else {
				artist.MemberStruct = append(artist.MemberStruct, Member{memberName, "/icons/artist_placeholder_08.png"})
			}
			//fmt.Println("Main Image URL:", page.WikiThumbnail.Source)
		}
	}
	/* fmt.Printf("Struct Artist: %v\nStruct Image:%v\n", artist.MemberStruct[0].MemberName, artist.MemberStruct[0].MemberImage)
	fmt.Printf("=== member's data of aritst:%v === \n", artist.Name)
	for memberName, imgLink := range artist.Members {
		fmt.Printf("Member: %v, imgLink: %v\n", memberName, imgLink)
	} */
}

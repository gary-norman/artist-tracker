package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type WikiResponse struct {
	WikiQuery WikiQuery `json:"query"`
}

func WikiImageFetcher(artist *Artist) {
	artist.Members = make(map[string]string)
	for _, member := range artist.MemberList {
		member = strings.Replace(member, " ", "_", -1)
		queryURL := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&titles=%s&prop=pageimages&format=json&pithumbsize=500", member)
		resp, err := http.Get(queryURL)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var result WikiResponse
		if err = json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Add the image URL to the map
		for _, page := range result.WikiQuery.Pages {
			if page.WikiThumbnail.Source != "" {
				artist.Members[strings.Replace(member, "_", " ", -1)] = page.WikiThumbnail.Source
			} else {
				artist.Members[strings.Replace(member, "_", " ", -1)] = "/icons/artist_placeholder_08.png"
			}
			//fmt.Println("Main Image URL:", page.WikiThumbnail.Source)
		}

		// Add the image URL to the struct
		for _, page := range result.WikiQuery.Pages {
			if page.WikiThumbnail.Source != "" {
				artist.MemberStruct = append(artist.MemberStruct, Member{strings.Replace(member, "_", " ", -1), page.WikiThumbnail.Source})
			} else {
				artist.MemberStruct = append(artist.MemberStruct, Member{strings.Replace(member, "_", " ", -1), "/icons/artist_placeholder_08.png"})
			}
			//fmt.Println("Main Image URL:", page.WikiThumbnail.Source)
		}
	}
	fmt.Printf("Struct Artist: %v\nStruct Image:%v\n", artist.MemberStruct[0].MemberName, artist.MemberStruct[0].MemberImage)
	fmt.Printf("=== member's data of aritst:%v === \n", artist.Name)
	for memberName, imgLink := range artist.Members {
		fmt.Printf("Member: %v, imgLink: %v\n", memberName, imgLink)
	}
}

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
		err := WikiImageFetcher(&artists[i])
		if err != nil {
			log.Fatal(err)
		}
	}
}

// WikiImageFetcher get individual member's image
func WikiImageFetcher(artist *Artist) error {
	notFound := "/icons/artist_placeholder.svg"
	//fmt.Printf("Getting member images for %v\n", artist.Name)
	artist.Members = make(map[string]string)
	searchCounter := 1
	for _, member := range artist.MemberList {
		var queryURL string
		var encodedMember string
		if member == "Roger Taylor" {
			encodedMember = "Roger_Taylor_(Queen_drummer)"
		} else {
			encodedMember = url.QueryEscape(member)
		}
		if artist.Name == "Mamonas Assassinas" {
			queryURL = fmt.Sprintf("https://pt.wikipedia.org/w/api.php?action=query&titles=%s&prop=pageimages&format=json&pithumbsize=500", encodedMember)
		} else {
			queryURL = fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&titles=%s&prop=pageimages&format=json&pithumbsize=500", encodedMember)
		}
		err := func() error {
			resp, err := http.Get(queryURL)
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}
			defer func(Body io.ReadCloser) {
				err = Body.Close()
				if err != nil {
					fmt.Println("Error closing connection:", err)
				}
			}(resp.Body)

			// debug only, delete later
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error: received status code %d for member %s\n", resp.StatusCode, member)
				body, _ := io.ReadAll(resp.Body)
				fmt.Println("Response body:", string(body))
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}
			var result WikiResponse
			if err = json.Unmarshal(body, &result); err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				fmt.Println("Error:", err)
				return err
			}

			// Add the image URL to the map
			for _, page := range result.WikiQuery.Pages {
				if page.WikiThumbnail.Source != "" {
					artist.Members[member] = page.WikiThumbnail.Source
					//fmt.Printf("%v - %v: success!\n", searchCounter, member)
				} else {
					artist.Members[member] = notFound
					//fmt.Printf("%v - %v: no member image found\n", searchCounter, member)
				}
				searchCounter += 1
			}

			// Add the image URL to the struct
			/* 	for _, page := range result.WikiQuery.Pages {
				if page.WikiThumbnail.Source != "" {
					artist.MemberStruct = append(artist.MemberStruct, Member{member, page.WikiThumbnail.Source})
				} else {
					artist.MemberStruct = append(artist.MemberStruct, Member{member, notFound})
				}
			} */
			return nil
		}()
		if err != nil {
			return err
		}
	}
	/* 	fmt.Println("*****************")
	   	for member, picLink := range artist.Members {
	   		fmt.Printf(Bold+Cyan+"member: "+Reset+Blue+"%v, "+Bold+Cyan+"Picture Link: "+Reset+Blue+"%v\n"+Reset, member, picLink)
	   	}*/
	return nil
}

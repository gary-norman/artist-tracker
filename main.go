package main

import (
	"artist-tracker/api"
	"fmt"
	"strconv"
	"time"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func main() {
	// Define the date layouts
	const layoutUK = "02-01-2006"
	const layoutUS = "2006-01-02"

	// setup config file to parse the key and token
	api.ConfigSetup()

	err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Artist", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithRGB("-", pterm.NewRGB(255, 215, 0)),
		putils.LettersFromStringWithStyle("Tracker", pterm.FgLightMagenta.ToStyle())).
		Render()
	if err != nil {
		fmt.Printf("Could not print BigText: %v", err)
		return
	} // Render the big text to the terminal
	spinnerInfo, _ := pterm.DefaultSpinner.Start("Fetching artist information")
	start := time.Now()
	// Populate the artist struct with API data
	Artists := api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
	t := time.Now()
	timetaken := t.Sub(start).Milliseconds()
	if len(Artists) == 0 {
		spinnerInfo.Fail()
	} else {
		spinnerInfo.Success("Fetched artist information in " + strconv.FormatInt(timetaken, 10) + "ms")
	}
	update := time.Now()
	// call the functions to populate extra information from TADB
	api.UpdateArtistInfo(Artists)
	for i := range Artists {
		api.FindFirstAlbum(&Artists[i])
	}
	t = time.Now()
	timetaken = t.Sub(update).Milliseconds()
	spinnerInfo.Success("Updated artist information in " + strconv.FormatInt(timetaken, 10) + "ms\n")
	api.CorrectMisnamedMembers(Artists)
	api.FetchAllArtistsImages(Artists)
	timetaken = t.Sub(start).Milliseconds()
	pterm.Info.Println("All tasks completed successfully in " + pterm.Green(strconv.FormatInt(timetaken, 10)+"ms"))
	//pterm.Println(pterm.Cyan(Artists[i]))
	//pterm.Println(pterm.Cyan("TourDetails {"))
	//pterm.Println(pterm.Cyan(Artists[i].TourDetails))
	//pterm.Println(pterm.Cyan("}"))

	// debug print, to see better all the information of an artist
	/* 	artistsResult, _ := api.SearchArtist(Artists, "Queen")
	   	api.PrintArtistDetails(artistsResult) */

	// test DateConvert
	/*fmt.Printf("UK > US: %v\n", api.DateConvert("11-09-2001", layoutUK, layoutUS))
	fmt.Printf("US > UK: %v\n", api.DateConvert("2001-09-11", layoutUS, layoutUK))

	fmt.Printf("Queen first event date: %v\nQueen first event date (converted): %v\n",
		Artists[0].FirstAlbum, api.DateConvert(Artists[0].FirstAlbum, layoutUK, layoutUS))*/
	for i := range Artists {
		fmt.Printf("%v: %v\n", i, Artists[i].Name)
	}
	api.HandleRequests(Artists, api.GetTemplate())
}

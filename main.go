package main

import (
	"artist-tracker/api"
	"fmt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"strconv"
	"time"
)

func main() {
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
	t = time.Now()
	timetaken = t.Sub(update).Milliseconds()
	spinnerInfo.Success("Updated artist information in " + strconv.FormatInt(timetaken, 10) + "ms\n")
	tour := time.Now()
	// Fetch and update tour information
	//i := 11
	//api.GetTourInfo(Artists, Artists[i].Name, i)
	pterm.DefaultBasicText.Println("Tour info" + pterm.Green(" only ") + "available for the following artists:")
	//tourLogger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	for i := 0; i < 12; i++ {
		api.UnmarshallTourInfo(Artists, i)
		if len(Artists[i].Data) > 0 {
			pterm.DefaultBasicText.Println(Artists[i].Name + ": " + pterm.Green("success"))
		}
		//pterm.Println(pterm.Cyan(Artists[i]))
		//pterm.Println(pterm.Cyan("TourDetails {"))
		//pterm.Println(pterm.Cyan(Artists[i].TourDetails))
		//pterm.Println(pterm.Cyan("}"))
	}
	t = time.Now()
	timetaken = t.Sub(tour).Milliseconds()
	spinnerInfo.Success("Updated tour information in " + strconv.FormatInt(timetaken, 10) + "ms\n\n")
	//err2 := pterm.DefaultBigText.WithLetters(
	//	putils.LettersFromStringWithStyle("Artist", pterm.FgCyan.ToStyle()),
	//	putils.LettersFromStringWithRGB("-", pterm.NewRGB(255, 215, 0)),
	//	putils.LettersFromStringWithStyle("Tracker", pterm.FgLightMagenta.ToStyle())).
	//	Render()
	//if err2 != nil {
	//	fmt.Printf("Could not print BigText: %v", err2)
	//	return
	//} // Render the big text to the terminal
	timetaken = t.Sub(start).Milliseconds()
	pterm.Info.Println("All tasks completed successfully in " + pterm.Green(strconv.FormatInt(timetaken, 10)+"ms"))
	//fmt.Println(api.GetTourInfo(Artists, "queen"))
	//for i := 1; i < 10; i++ { // TODO run this on monday
	//i := 3
	//api.GetTourInfo(Artists, Artists[i].Name, i)
	//}

	api.HandleRequests(Artists, api.GetTemplate())
}

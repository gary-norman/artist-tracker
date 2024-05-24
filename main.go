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
	// Call AllJsonToStruct and print the result
	Artists := api.AllJsonToStruct("https://groupietrackers.herokuapp.com/api/artists")
	t := time.Now()
	timetaken := t.Sub(start).Milliseconds()
	if len(Artists) == 0 {
		spinnerInfo.Fail()
	} else {
		spinnerInfo.Success("Fetched artist information in " + strconv.FormatInt(timetaken, 10) + "ms")
	}
	api.UpdateArtistInfo(Artists)
	pterm.Println(pterm.Cyan(Artists[23]))
	err2 := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Artist", pterm.FgCyan.ToStyle()),
		putils.LettersFromStringWithRGB("-", pterm.NewRGB(255, 215, 0)),
		putils.LettersFromStringWithStyle("Tracker", pterm.FgLightMagenta.ToStyle())).
		Render()
	if err2 != nil {
		fmt.Printf("Could not print BigText: %v", err2)
		return
	} // Render the big text to the terminal
	t = time.Now()
	timetaken = t.Sub(start).Milliseconds()
	pterm.Info.Println("All tasks completed successfully in " + pterm.Green(strconv.FormatInt(timetaken, 10)+"ms"))
	api.HandleRequests(Artists, api.GetTemplate())

}

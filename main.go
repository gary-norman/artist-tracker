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
	MapBoxHtmlValues := make([][]string, 52)
	MapBoxHtmlValues[0] = []string{"clwunn3x6016c01qx2kio2sfj", Artists[0].Name + "-tourdates"}
	MapBoxHtmlValues[2] = []string{"clwxi2kg0017h01pca2qs5ay1", Artists[2].Name + "-tourdates"}
	MapBoxHtmlValues[4] = []string{"clwxj7j4h01gn01qrg1ba9esp", Artists[4].Name + "-tourdates"}
	MapBoxHtmlValues[5] = []string{"clwxq6lkk01h501qr1hyxargt", Artists[5].Name + "-tourdates"}
	MapBoxHtmlValues[6] = []string{"clwxqtxg9018101pc5860cdo8", Artists[5].Name + "-tourdates"}
	MapBoxHtmlValues[8] = []string{"clwxr03zk01h801qr41qs1kym", Artists[8].Name + "-tourdates"}
	MapBoxHtmlValues[11] = []string{"clwxr4pik01f301nye5wh4cxl", Artists[11].Name + "-tourdates"}
	MapBoxHtmlValues[13] = []string{"clwxrbg4g01f401nygj1d46xd", Artists[13].Name + "-tourdates"}
	MapBoxHtmlValues[14] = []string{"clwxrf5aa01am01qx2sj8d83f", Artists[14].Name + "-tourdates"}
	MapBoxHtmlValues[15] = []string{"clwxrkjf501f501nygntrglk7", Artists[15].Name + "-tourdates"}
	MapBoxHtmlValues[16] = []string{"", Artists[16].Name + "-tourdates"}
	MapBoxHtmlValues[18] = []string{"", Artists[18].Name + "-tourdates"}
	MapBoxHtmlValues[19] = []string{"", Artists[19].Name + "-tourdates"}
	MapBoxHtmlValues[20] = []string{"", Artists[20].Name + "-tourdates"}
	MapBoxHtmlValues[21] = []string{"", Artists[21].Name + "-tourdates"}
	MapBoxHtmlValues[23] = []string{"", Artists[23].Name + "-tourdates"}
	MapBoxHtmlValues[26] = []string{"", Artists[26].Name + "-tourdates"}
	MapBoxHtmlValues[32] = []string{"", Artists[32].Name + "-tourdates"}
	MapBoxHtmlValues[34] = []string{"", Artists[34].Name + "-tourdates"}
	MapBoxHtmlValues[35] = []string{"", Artists[35].Name + "-tourdates"}
	MapBoxHtmlValues[40] = []string{"", Artists[40].Name + "-tourdates"}
	MapBoxHtmlValues[42] = []string{"", Artists[42].Name + "-tourdates"}
	MapBoxHtmlValues[48] = []string{"", Artists[48].Name + "-tourdates"}
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
	//for i := 49; i < 52; i++ {
	//	api.GetTourInfo(Artists, Artists[i].Name, i)
	//	duration := time.Second
	//	time.Sleep(duration)
	//}
	api.UpdateACDC(Artists)
	// Fetch and update tour information
	//api.GetTourInfo(Artists, Artists[i].Name, i)
	//var i int
	//for i = 0; i < 52; i++ {
	//	api.UnmarshallTourInfo(Artists, i)
	//	if len(Artists[i].Data) > 0 {
	//		pterm.DefaultBasicText.Println(Artists[i].Name + ": " + pterm.Green("success"))
	//	}
	//	api.RapidToMapbox(i)
	//}
	fmt.Println("Artists with geojson data:")
	for i := 0; i < 52; i++ {
		api.GeojsonCheck(i, Artists[i].Name)
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
	//i = 48
	pterm.Info.Println("All tasks completed successfully in " + pterm.Green(strconv.FormatInt(timetaken, 10)+"ms"))
	//pterm.Println(pterm.Cyan(Artists[i]))
	//pterm.Println(pterm.Cyan("TourDetails {"))
	//pterm.Println(pterm.Cyan(Artists[i].TourDetails))
	//pterm.Println(pterm.Cyan("}"))
	//fmt.Println(api.GetTourInfo(Artists, "queen"))
	//for i := 1; i < 10; i++ { // TODO run this on monday
	//i := 3
	//api.GetTourInfo(Artists, Artists[i].Name, i)
	//}
	// upload datasets to mapbox
	//indices := []int{13, 14, 15, 16, 18, 19, 20, 21, 23, 26, 32, 34, 35, 40, 42, 48}
	//for _, i := range indices {
	//	api.MapboxDataset(i, Artists[i].Name)
	//}
	api.HandleRequests(Artists, api.GetTemplate())
}

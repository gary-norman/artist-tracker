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
	//MapBoxHtmlValues := make(map[string][]string, 52)
	//MapBoxHtmlValues[Artists[0].Name] = []string{"clwunn3x6016c01qx2kio2sfj", strings.Replace(Artists[0].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[1].Name] = []string{"clwunn3x6016c01qx2kio2sfj", strings.Replace(Artists[1].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[2].Name] = []string{"clwxi2kg0017h01pca2qs5ay1", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[3].Name] = []string{"clx3cl1uv01j801qs41ma8oal", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[4].Name] = []string{"clwxj7j4h01gn01qrg1ba9esp", strings.Replace(Artists[4].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[5].Name] = []string{"clwxq6lkk01h501qr1hyxargt", strings.Replace(Artists[5].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[6].Name] = []string{"clwxqtxg9018101pc5860cdo8", strings.Replace(Artists[5].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[7].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[8].Name] = []string{"clwxr03zk01h801qr41qs1kym", strings.Replace(Artists[8].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[9].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[10].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[11].Name] = []string{"clwxr4pik01f301nye5wh4cxl", strings.Replace(Artists[11].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[12].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[13].Name] = []string{"clwxrbg4g01f401nygj1d46xd", strings.Replace(Artists[13].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[14].Name] = []string{"clwxrf5aa01am01qx2sj8d83f", strings.Replace(Artists[14].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[15].Name] = []string{"clwxrkjf501f501nygntrglk7", strings.Replace(Artists[15].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[16].Name] = []string{"clwz2oy7p00wx01poh4qlc1o8", strings.Replace(Artists[16].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[17].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[18].Name] = []string{"clx09zgzt01b101qsdabn9piw", strings.Replace(Artists[18].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[19].Name] = []string{"clx1vt09001mq01pn10pj71gu", strings.Replace(Artists[19].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[20].Name] = []string{"clx1w2mm401ev01qs2fuaeg8l", strings.Replace(Artists[20].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[21].Name] = []string{"clx1w8l2u01rv01qs9s24b3ym", strings.Replace(Artists[21].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[22].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[23].Name] = []string{"clx393obm01ma01qxbu3c79d2", strings.Replace(Artists[23].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[24].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[25].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[26].Name] = []string{"clx39u1kv01q801pn3ixs2ijh", strings.Replace(Artists[26].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[27].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[28].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[29].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[30].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[31].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[32].Name] = []string{"clx39uy7u01q901pnf0z31wnp", strings.Replace(Artists[32].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[33].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[34].Name] = []string{"clx39v50t01ip01qsa3fvhedt", strings.Replace(Artists[34].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[35].Name] = []string{"clx39v9aj01k801pc71obdf80", strings.Replace(Artists[35].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[36].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[37].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[38].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[39].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[40].Name] = []string{"clx39vj6p01v301qsforzd597", strings.Replace(Artists[40].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[41].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[42].Name] = []string{"clx39vfpd01v201qs6fo4ai8a", strings.Replace(Artists[42].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[43].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[44].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[45].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[46].Name] = []string{"clx09zgzt01b101qsdabn9piw", strings.Replace(Artists[46].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[47].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[48].Name] = []string{"clx3av9id007l01qs4pcm2u3u", strings.Replace(Artists[48].Name, " ", "-", -1) + "-tourdates"}
	//MapBoxHtmlValues[Artists[49].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[50].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
	//MapBoxHtmlValues[Artists[51].Name] = []string{"", strings.Replace(Artists[2].Name, " ", "-", -1) + "-tourdates-std"}
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
	api.UpdateACDC(Artists)
	timetaken = t.Sub(start).Milliseconds()
	i := 48
	pterm.Info.Println("All tasks completed successfully in " + pterm.Green(strconv.FormatInt(timetaken, 10)+"ms"))
	pterm.Println(pterm.Cyan(Artists[i]))
	pterm.Println(pterm.Cyan("TourDetails {"))
	pterm.Println(pterm.Cyan(Artists[i].TourDetails))
	pterm.Println(pterm.Cyan("}"))
	//fmt.Println(api.GetTourInfo(Artists, "queen", 0))
	//for i, artist := range Artists {
	//	fmt.Printf("%d: %v\n", i, artist.Name)
	//for location, dates := range artist.DatesLocations {
	//	fmt.Printf("Location: %v\n", location)
	//	for _, date := range dates {
	//		fmt.Printf("Date: %v\n", date)
	//	}
	//}
	//}
	//api.MapboxReverseLookup(1, Artists[1])

	// debug print, to see better all the informations of an aritsts
	artistsResult, _ := api.SearchArtist(Artists, "Led Zeppelin")
	api.PrintArtistDetails(artistsResult)

	api.HandleRequests(Artists, api.GetTemplate())
}

package api

import (
	"fmt"
	"time"
)

// DateConvert converts from US to GB dates or vice versa
// eg newdate := DateConvert(2019-08-16, LayoutUS, LayoutUK) (newdate = 16-08-2019)
func DateConvert(original, from, to string) (output string) {
	date, err := time.Parse(from, original)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	return date.Format(to)
}

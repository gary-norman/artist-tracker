package api

import (
	"fmt"
	"math"
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

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

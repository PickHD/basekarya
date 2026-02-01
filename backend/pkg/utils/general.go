package utils

import "fmt"

func FormatNumber(n float64) string {
	s := fmt.Sprintf("%.0f", n)
	nStr := ""
	for i, j := len(s)-1, 0; i >= 0; i, j = i-1, j+1 {
		if j > 0 && j%3 == 0 {
			nStr = "." + nStr
		}
		nStr = string(s[i]) + nStr
	}
	return nStr
}

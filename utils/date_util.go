package utils

import (
	"fmt"
	"time"
)

func GetYMD() string {
	currentTime := time.Now()
	year := currentTime.Year()
	month := currentTime.Month()
	day := currentTime.Day()
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

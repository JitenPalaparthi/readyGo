package global

import "time"

// Yet to de developeds

// GetCurrentDateTimeInStr is to get the current date and time in string format
func GetCurrentDateTimeInStr() string {
	return time.Now().String()
}

// GetDefaultStr is to get the default value
func GetDefaultStr(str string) string {
	return str
}

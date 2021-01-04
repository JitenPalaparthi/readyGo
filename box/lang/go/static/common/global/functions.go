package global

import "time"

// Yet to de developeds

// GetCurrentDateTimeInStr is to get the current date and time in string format
func GetCurrentDateTimeInStr() string {
	return time.Now().String()
}

// GetCurrentDateTimeInTime is to get the current date and time in string format
func GetCurrentDateTimeInTime() time.Time {
	return time.Now()
}

// GetDefaultStr is to get the default value which is of string type
func GetDefaultStr(str string) string {
	return str
}

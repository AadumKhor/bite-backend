package utils

import "time"

// Global variable to be used elsewhere
var LocIndia *time.Location

func InitTimeZone(timeZone string) error {
	locIndia, err := time.LoadLocation(timeZone)
	if err != nil {
		return err
	}

	// set value to global variable via the function
	LocIndia = locIndia
	return nil
}
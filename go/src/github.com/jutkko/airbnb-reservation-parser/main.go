package main

import (
	"flag"
	"fmt"
	"github.com/jutkko/airbnb-reservation-parser/listing"
	"log"
	"time"
)

func main() {
	fromDatePtr := flag.String("from-date", "", "the starting date for the book rate calculation, e.g., 2019-01-20")
	toDatePtr := flag.String("to-date", "", "the ending date for the book rate calculation, e.g., 2019-01-21")
	flag.Parse()

	var fromDate, toDate time.Time
	var err error
	currentTime := time.Now()

	if *fromDatePtr == "" {
		y, m, d := currentTime.Date()
		// Defaults to begin of the day
		fromDate = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		if err != nil {
			log.Fatal("failed to get today's date ", err)
		}
	} else {
		fromDate, err = time.Parse(listing.LayoutISO, *fromDatePtr)
		if err != nil {
			log.Fatal("failed to parse from date ", err)
		}
	}

	if *toDatePtr == "" {
		y, m, _ := currentTime.Date()
		beginningOfMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
		// Defaults to end of the month, which is the beginning of the next month
		toDate = beginningOfMonth.AddDate(0, 1, 0)
	} else {
		toDate, err = time.Parse(listing.LayoutISO, *toDatePtr)

		// Defaults to the next day, since the calculation is not inclusive
		toDate = toDate.AddDate(0, 0, 1)
		if err != nil {
			log.Fatal("failed to parse to date ", err)
		}
	}

	myListing, err := listing.ProcessData("data/reservations.csv")
	if err != nil {
		log.Fatal("failed to process date ", err)
	}

	bookRate, bookPrice := myListing.GetBookRateAndPrice(fromDate, toDate)
	fmt.Printf("book rate: %.2f price %.2f\n", bookRate, bookPrice)
}

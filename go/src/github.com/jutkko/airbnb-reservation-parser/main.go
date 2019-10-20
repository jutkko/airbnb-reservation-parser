package main

import (
	"fmt"
	"github.com/jutkko/airbnb-reservation-parser/listing"
	"log"
)

func main() {
	listing, err := listing.ProcessData("data/reservations.csv")
	if err != nil {
		log.Fatal("failed to process data", err)
	}

	for _, reservation := range listing.Reservations {
		fmt.Printf("Start: %s\n", reservation.StartDate)
		fmt.Printf("End: %s\n", reservation.EndDate)
	}
}

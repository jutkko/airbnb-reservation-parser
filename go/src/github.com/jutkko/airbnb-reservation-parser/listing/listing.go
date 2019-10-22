package listing

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	LayoutISO = "2006-01-02"
	Cancelled = "已取消"
	Confirmed = "已确认"
)

type Date struct {
	time.Time
}

// Convert the CSV string as internal date
func (d*Date) UnmarshalCSV(csv string) (err error) {
	d.Time, err = time.Parse(LayoutISO, csv)
	return err
}

type Price struct {
	float64
}

// Convert the CSV string as internal date
func (p *Price) UnmarshalCSV(csv string) (err error) {
	wholePrice := strings.Split(csv, "€")

	p.float64, err = strconv.ParseFloat(wholePrice[1], 64)
	return err
}

type Reservation struct {
	ReservationCode string `csv:"reservation_code"`
	Status string `csv:"status"`
	Name string `csv:"name"`
	PhoneNumber string `csv:"phone_number"`
	Adults int `csv:"adults"`
	Children int `csv:"children"`
	Infants int `csv:"infants"`
	StartDate Date `csv:"start_date"`
	EndDate Date `csv:"end_date"`
	Nights int `csv:"nights"`
	ConfirmationDate Date `csv:"confirmation_date"`
	Flat string `csv:"flat"`
	Price Price `csv:"price"`
}

type Listing struct {
	Reservations []*Reservation
}

func ProcessData(filename string) (*Listing, error) {
	rawReservationFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer rawReservationFile.Close()

	var rawReservations []*Reservation
	if err := gocsv.UnmarshalFile(rawReservationFile, &rawReservations); err != nil {
		return nil, err
	}

	var reservations []*Reservation
	for _, reservation := range rawReservations {
		if reservation.Status == Confirmed {
			reservations = append(reservations, reservation)
		}
	}

	// Sort the reservations ascending by start date
	sort.Slice(reservations, func(i, j int) bool {
		return reservations[i].StartDate.Time.Before(reservations[j].StartDate.Time)
	})

	return &Listing{
		Reservations: reservations,
	}, nil
}

// From has to be strictly earlier than to
func (l *Listing) GetBookRate(from, to time.Time) float64 {
	var bookedNightsInRange int
	for _, reservation := range l.Reservations {
		if to.Before(reservation.StartDate.Time) || from.After(reservation.EndDate.Time) {
			continue
		}

		fmt.Printf("for %s's booking before %d\n", reservation.Name, bookedNightsInRange)
		if (reservation.StartDate.Time.After(from) || reservation.StartDate.Time.Equal(from)) && (reservation.EndDate.Time.Before(to) || reservation.EndDate.Time.Equal(to)) {
			// (from to)          *****
			// (reserve start end) ***
			bookedNightsInRange += reservation.Nights
		} else if reservation.StartDate.Time.After(from) && reservation.EndDate.After(to) {
			// (from to)          *****
			// (reserve start end)  *****
			bookedNightsInRange += int(to.Sub(reservation.StartDate.Time).Hours()/24)
		} else if reservation.StartDate.Time.Before(from) && reservation.EndDate.Before(to) {
			// (from to)            *****
			// (reserve start end) *****
			bookedNightsInRange += int(reservation.EndDate.Sub(from).Hours()/24)
		} else if reservation.StartDate.Time.Before(from) && reservation.EndDate.After(to) {
			// (from to)            ***
			// (reserve start end) *****
			bookedNightsInRange += int(to.Sub(from).Hours()/24)
		} else if reservation.StartDate.Time.Equal(from) && reservation.EndDate.After(to) {
			// (from to)           ***
			// (reserve start end) *****
			bookedNightsInRange += int(to.Sub(reservation.StartDate.Time).Hours()/24)
		} else if reservation.StartDate.Time.Before(from) && reservation.EndDate.Equal(to) {
			// (from to)             ***
			// (reserve start end) *****
			bookedNightsInRange += int(reservation.EndDate.Time.Sub(from).Hours()/24)
		}

		fmt.Printf("for %s's booking after %d\n", reservation.Name, bookedNightsInRange)
	}

	fmt.Printf("booked nights: %d\n", bookedNightsInRange)
	return float64(bookedNightsInRange) / (to.Sub(from).Hours()/24)
}

package listing

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testFilename = "../../../../../../data/test.csv"

func TestProcessData(t *testing.T) {
	_, err := ProcessData(testFilename)
	assert.NoError(t, err, "process should not fail")
}

func TestProcessDataIgnoreCancelled(t *testing.T) {
	testListing, err := ProcessData(testFilename)
	assert.NoError(t, err, "process should not fail")

	for _, reservation := range testListing.Reservations {
		assert.Equal(t, reservation.Status, Confirmed, "all reservations should be confirmed")
	}
}

func TestProcessDataSorted(t *testing.T) {
	testListing, err := ProcessData(testFilename)
	assert.NoError(t, err, "process should not fail")

	var prev *Reservation
	prev = testListing.Reservations[0]
	for _, reservation := range testListing.Reservations[1:] {
		assert.True(t, prev.StartDate.Time.Before(reservation.StartDate.Time), "all reservations should be confirmed")
		prev = reservation
	}
}

func TestListingGetBookRate(t *testing.T) {
	testListing, err := ProcessData(testFilename)
	assert.NoError(t, err, "process should not fail")

	from, err := time.Parse(LayoutISO, "2019-12-30")
	assert.NoError(t, err, "parse should not fail")

	to, err := time.Parse(LayoutISO, "2019-12-31")
	assert.NoError(t, err, "parse should not fail")

	assert.Equal(t, float64(1), testListing.GetBookRate(from, to), "book rate should match")
}

func TestListingGetBookRate1(t *testing.T) {
	testListing, err := ProcessData(testFilename)
	assert.NoError(t, err, "process should not fail")

	from, err := time.Parse(LayoutISO, "2019-09-13")
	assert.NoError(t, err, "parse should not fail")

	to, err := time.Parse(LayoutISO, "2019-09-14")
	assert.NoError(t, err, "parse should not fail")

	assert.Equal(t, float64(0), testListing.GetBookRate(from, to), "book rate should match")
}

func TestListingGetBookRate2(t *testing.T) {
	testListing, err := ProcessData(testFilename)
	assert.NoError(t, err, "process should not fail")

	from, err := time.Parse(LayoutISO, "2019-09-13")
	assert.NoError(t, err, "parse should not fail")

	to, err := time.Parse(LayoutISO, "2019-09-15")
	assert.NoError(t, err, "parse should not fail")

	assert.Equal(t, 0.5, testListing.GetBookRate(from, to), "book rate should match")
}

func TestListingGetBookRate3(t *testing.T) {
	testListing, err := ProcessData(testFilename)
	assert.NoError(t, err, "process should not fail")

	from, err := time.Parse(LayoutISO, "2019-09-13")
	assert.NoError(t, err, "parse should not fail")

	to, err := time.Parse(LayoutISO, "2019-09-18")
	assert.NoError(t, err, "parse should not fail")

	assert.Equal(t, 0.4, testListing.GetBookRate(from, to), "book rate should match")
}

func TestListingGetBookRate4(t *testing.T) {
	testListing, err := ProcessData(testFilename)
	assert.NoError(t, err, "process should not fail")

	from, err := time.Parse(LayoutISO, "2019-09-13")
	assert.NoError(t, err, "parse should not fail")

	to, err := time.Parse(LayoutISO, "2019-09-19")
	assert.NoError(t, err, "parse should not fail")

	assert.Equal(t, 0.5, testListing.GetBookRate(from, to), "book rate should match")
}

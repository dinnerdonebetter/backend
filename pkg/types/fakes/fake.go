package fakes

import (
	"fmt"
	"math"
	"time"

	"github.com/dinnerdonebetter/backend/internal/identifiers"

	fake "github.com/brianvoe/gofakeit/v5"
)

func init() {
	fake.Seed(time.Now().UnixNano())
}

const (
	exampleQuantity = 3
)

// BuildFakeSQLQuery builds a fake SQL query and arg pair.
func BuildFakeSQLQuery() (query string, args []any) {
	s := fmt.Sprintf("%s %s WHERE things = ? AND stuff = ?",
		fake.RandomString([]string{"SELECT * FROM", "INSERT INTO", "UPDATE"}),
		fake.Password(true, true, true, false, false, 32),
	)

	return s, []any{"things", "stuff"}
}

// BuildFakeID builds a fake ID.
func BuildFakeID() string {
	return identifiers.New()
}

func BuildFakeNumber() float64 {
	return math.Round(float64((fake.Number(101, math.MaxInt8-1) * 100) / 100))
}

// BuildFakeTime builds a fake time.
func BuildFakeTime() time.Time {
	return fake.Date().Add(0).Truncate(time.Second).UTC()
}

// buildUniqueString builds a fake string.
func buildUniqueString() string {
	return fake.LoremIpsumSentence(7)
}

// BuildFakePassword builds a fake password.
func BuildFakePassword() string {
	return fake.Password(true, true, true, true, false, 32)
}

package main

import (
	"testing"

	"github.com/cristalhq/builq"
	"github.com/stretchr/testify/assert"
)

func Test_applyToEach(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleInput := []string{
			"things",
			"and",
			"stuff",
		}

		callCount := 0
		exampleFunc := func(x string) string {
			callCount += 1
			return x
		}

		expected := []string{
			"things",
			"and",
			"stuff",
		}
		actual := applyToEach(exampleInput, exampleFunc)

		assert.Equal(t, callCount, len(exampleInput))
		assert.Equal(t, expected, actual)
	})
}

func Test_buildRawQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var whatever builq.Builder

		builder := whatever.Addf("SELECT * FROM things")

		expected := "SELECT * FROM things\n"
		actual := buildRawQuery(builder)

		assert.Equal(t, expected, actual)
	})
}

func Test_filterForInsert(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exception := "whatever"
		exampleColumns := []string{
			"things",
			"and",
			"stuff",
			createdAtColumn,
			lastUpdatedAtColumn,
			archivedAtColumn,
			exception,
		}

		expected := []string{
			"things",
			"and",
			"stuff",
		}
		actual := filterForInsert(exampleColumns, exception)

		assert.Equal(t, expected, actual)
	})
}

func Test_formatQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		example := `SELECT stuff
FROM things
				WHERE id = 1
`

		expected := "SELECT stuff FROM things WHERE id = 1;"
		actual := formatQuery(example)

		assert.Equal(t, expected, actual)
	})
}

func Test_fullColumnName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "things.stuff"
		actual := fullColumnName("things", "stuff")

		assert.Equal(t, expected, actual)
	})
}

func Test_mergeColumns(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"webhooks.id",
			"webhooks.name",
			"webhooks.content_type",
			"webhooks.url",
			"webhooks.method",
			"webhook_trigger_events.id",
			"webhook_trigger_events.trigger_event",
			"webhook_trigger_events.belongs_to_webhook",
			"webhook_trigger_events.created_at",
			"webhook_trigger_events.archived_at",
			"webhooks.created_at",
			"webhooks.last_updated_at",
			"webhooks.archived_at",
			"webhooks.belongs_to_household",
		}

		actual := mergeColumns(
			applyToEach(webhooksColumns, func(s string) string {
				return fullColumnName(webhooksTableName, s)
			}),
			applyToEach(webhookTriggerEventsColumns, func(s string) string {
				return fullColumnName(webhookTriggerEventsTableName, s)
			}),
			5,
		)

		assert.Equal(t, expected, actual)
	})
}

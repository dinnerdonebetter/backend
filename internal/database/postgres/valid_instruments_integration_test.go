package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func createValidInstrumentForTest(t *testing.T, ctx context.Context, exampleValidInstrument *types.ValidInstrument, dbc *Querier) *types.ValidInstrument {
	t.Helper()

	// create
	if exampleValidInstrument == nil {
		exampleValidInstrument = fakes.BuildFakeValidInstrument()
	}
	dbInput := converters.ConvertValidInstrumentToValidInstrumentDatabaseCreationInput(exampleValidInstrument)

	created, err := dbc.CreateValidInstrument(ctx, dbInput)
	exampleValidInstrument.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidInstrument, created)

	validInstrument, err := dbc.GetValidInstrument(ctx, created.ID)
	exampleValidInstrument.CreatedAt = validInstrument.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validInstrument, exampleValidInstrument)

	return created
}

func TestQuerier_Integration_ValidInstruments(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidInstrument := fakes.BuildFakeValidInstrument()
	createdValidInstruments := []*types.ValidInstrument{}

	// create
	createdValidInstruments = append(createdValidInstruments, createValidInstrumentForTest(t, ctx, exampleValidInstrument, dbc))

	// update
	updatedValidInstrument := fakes.BuildFakeValidInstrument()
	updatedValidInstrument.ID = createdValidInstruments[0].ID
	assert.NoError(t, dbc.UpdateValidInstrument(ctx, updatedValidInstrument))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidInstrument()
		input.Name = fmt.Sprintf("%s %d", updatedValidInstrument.Name, i)
		createdValidInstruments = append(createdValidInstruments, createValidInstrumentForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validInstruments, err := dbc.GetValidInstruments(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validInstruments.Data)
	assert.Equal(t, len(createdValidInstruments), len(validInstruments.Data))

	// fetch as list of IDs
	validInstrumentIDs := []string{}
	for _, validInstrument := range createdValidInstruments {
		validInstrumentIDs = append(validInstrumentIDs, validInstrument.ID)
	}

	byIDs, err := dbc.GetValidInstrumentsWithIDs(ctx, validInstrumentIDs)
	assert.NoError(t, err)
	assert.Equal(t, validInstruments.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidInstruments(ctx, updatedValidInstrument.Name)
	assert.NoError(t, err)
	assert.Equal(t, validInstruments.Data, byName)

	// delete
	for _, validInstrument := range createdValidInstruments {
		assert.NoError(t, dbc.ArchiveValidInstrument(ctx, validInstrument.ID))

		var exists bool
		exists, err = dbc.ValidInstrumentExists(ctx, validInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidInstrument
		y, err = dbc.GetValidInstrument(ctx, validInstrument.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

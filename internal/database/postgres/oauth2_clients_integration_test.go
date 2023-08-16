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

func createOAuth2ClientForTest(t *testing.T, ctx context.Context, exampleOAuth2Client *types.OAuth2Client, dbc *Querier) *types.OAuth2Client {
	t.Helper()

	// create
	if exampleOAuth2Client == nil {
		exampleOAuth2Client = fakes.BuildFakeOAuth2Client()
	}
	dbInput := converters.ConvertOAuth2ClientToOAuth2ClientDatabaseCreationInput(exampleOAuth2Client)

	created, err := dbc.CreateOAuth2Client(ctx, dbInput)
	exampleOAuth2Client.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleOAuth2Client, created)

	oauth2Client, err := dbc.GetOAuth2ClientByDatabaseID(ctx, created.ID)
	exampleOAuth2Client.CreatedAt = oauth2Client.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, oauth2Client, exampleOAuth2Client)

	return created
}

func TestQuerier_Integration_OAuth2Clients(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleOAuth2Client := fakes.BuildFakeOAuth2Client()
	createdOAuth2Clients := []*types.OAuth2Client{}

	// create
	createdOAuth2Clients = append(createdOAuth2Clients, createOAuth2ClientForTest(t, ctx, exampleOAuth2Client, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeOAuth2Client()
		input.Name = fmt.Sprintf("%s %d", exampleOAuth2Client.Name, i)
		createdOAuth2Clients = append(createdOAuth2Clients, createOAuth2ClientForTest(t, ctx, input, dbc))
	}

	// fetch as list
	oauth2Clients, err := dbc.GetOAuth2Clients(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, oauth2Clients.Data)
	assert.Equal(t, len(createdOAuth2Clients), len(oauth2Clients.Data))

	// delete
	for _, oauth2Client := range createdOAuth2Clients {
		assert.NoError(t, dbc.ArchiveOAuth2Client(ctx, oauth2Client.ID))

		var y *types.OAuth2Client
		y, err = dbc.GetOAuth2ClientByClientID(ctx, oauth2Client.ClientID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
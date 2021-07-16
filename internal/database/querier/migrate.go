package querier

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// Migrate is a simple wrapper around the core querier Migrate call.
func (q *SQLQuerier) Migrate(ctx context.Context, maxAttempts uint8, testUserConfig *types.TestUserCreationConfig) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	q.logger.Info("migrating db")

	if !q.IsReady(ctx, maxAttempts) {
		return database.ErrDatabaseNotReady
	}

	q.migrateOnce.Do(q.sqlQueryBuilder.BuildMigrationFunc(q.db))

	if testUserConfig != nil {
		q.logger.Debug("creating test user")

		testUserExistenceQuery, testUserExistenceArgs := q.sqlQueryBuilder.BuildGetUserByUsernameQuery(ctx, testUserConfig.Username)
		userRow := q.getOneRow(ctx, q.db, "user", testUserExistenceQuery, testUserExistenceArgs...)

		_, _, _, err := q.scanUser(ctx, userRow, false)
		if err != nil {
			testUserCreationQuery, testUserCreationArgs := q.sqlQueryBuilder.BuildTestUserCreationQuery(ctx, testUserConfig)

			// these structs will be fleshed out by createUser
			user := &types.User{Username: testUserConfig.Username}
			account := &types.Account{}

			if err = q.createUser(ctx, user, account, testUserCreationQuery, testUserCreationArgs); err != nil {
				observability.AcknowledgeError(err, q.logger, span, "creating test user")
			}

			q.logger.WithValue(keys.UsernameKey, testUserConfig.Username).Debug("created test user and account")
		}
	}

	return nil
}

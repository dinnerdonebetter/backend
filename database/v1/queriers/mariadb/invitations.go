package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	invitationsTableName = "invitations"
)

var (
	invitationsTableColumns = []string{
		"id",
		"code",
		"consumed",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanInvitation takes a database Scanner (i.e. *sql.Row) and scans the result into an Invitation struct
func scanInvitation(scan database.Scanner) (*models.Invitation, error) {
	x := &models.Invitation{}

	if err := scan.Scan(
		&x.ID,
		&x.Code,
		&x.Consumed,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanInvitations takes a logger and some database rows and turns them into a slice of invitations
func scanInvitations(logger logging.Logger, rows *sql.Rows) ([]models.Invitation, error) {
	var list []models.Invitation

	for rows.Next() {
		x, err := scanInvitation(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildGetInvitationQuery constructs a SQL query for fetching an invitation with a given ID belong to a user with a given ID.
func (m *MariaDB) buildGetInvitationQuery(invitationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(invitationsTableColumns...).
		From(invitationsTableName).
		Where(squirrel.Eq{
			"id":         invitationID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetInvitation fetches an invitation from the mariadb database
func (m *MariaDB) GetInvitation(ctx context.Context, invitationID, userID uint64) (*models.Invitation, error) {
	query, args := m.buildGetInvitationQuery(invitationID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanInvitation(row)
}

// buildGetInvitationCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of invitations belonging to a given user that meet a given query
func (m *MariaDB) buildGetInvitationCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(CountQuery).
		From(invitationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetInvitationCount will fetch the count of invitations from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetInvitationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetInvitationCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allInvitationsCountQueryBuilder sync.Once
	allInvitationsCountQuery        string
)

// buildGetAllInvitationsCountQuery returns a query that fetches the total number of invitations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllInvitationsCountQuery() string {
	allInvitationsCountQueryBuilder.Do(func() {
		var err error
		allInvitationsCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(invitationsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allInvitationsCountQuery
}

// GetAllInvitationsCount will fetch the count of invitations from the database
func (m *MariaDB) GetAllInvitationsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllInvitationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetInvitationsQuery builds a SQL query selecting invitations that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetInvitationsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(invitationsTableColumns...).
		From(invitationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	m.logQueryBuildingError(err)

	return query, args
}

// GetInvitations fetches a list of invitations from the database that meet a particular filter
func (m *MariaDB) GetInvitations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.InvitationList, error) {
	query, args := m.buildGetInvitationsQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for invitations")
	}

	list, err := scanInvitations(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetInvitationCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching invitation count: %w", err)
	}

	x := &models.InvitationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Invitations: list,
	}

	return x, nil
}

// GetAllInvitationsForUser fetches every invitation belonging to a user
func (m *MariaDB) GetAllInvitationsForUser(ctx context.Context, userID uint64) ([]models.Invitation, error) {
	query, args := m.buildGetInvitationsQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching invitations for user")
	}

	list, err := scanInvitations(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateInvitationQuery takes an invitation and returns a creation query for that invitation and the relevant arguments.
func (m *MariaDB) buildCreateInvitationQuery(input *models.Invitation) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(invitationsTableName).
		Columns(
			"code",
			"consumed",
			"belongs_to",
			"created_on",
		).
		Values(
			input.Code,
			input.Consumed,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildInvitationCreationTimeQuery takes an invitation and returns a creation query for that invitation and the relevant arguments
func (m *MariaDB) buildInvitationCreationTimeQuery(invitationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(invitationsTableName).
		Where(squirrel.Eq{"id": invitationID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateInvitation creates an invitation in the database
func (m *MariaDB) CreateInvitation(ctx context.Context, input *models.InvitationCreationInput) (*models.Invitation, error) {
	x := &models.Invitation{
		Code:      input.Code,
		Consumed:  input.Consumed,
		BelongsTo: input.BelongsTo,
	}

	query, args := m.buildCreateInvitationQuery(x)

	// create the invitation
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing invitation creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildInvitationCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateInvitationQuery takes an invitation and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdateInvitationQuery(input *models.Invitation) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(invitationsTableName).
		Set("code", input.Code).
		Set("consumed", input.Consumed).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateInvitation updates a particular invitation. Note that UpdateInvitation expects the provided input to have a valid ID.
func (m *MariaDB) UpdateInvitation(ctx context.Context, input *models.Invitation) error {
	query, args := m.buildUpdateInvitationQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveInvitationQuery returns a SQL query which marks a given invitation belonging to a given user as archived.
func (m *MariaDB) buildArchiveInvitationQuery(invitationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(invitationsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          invitationID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveInvitation marks an invitation as archived in the database
func (m *MariaDB) ArchiveInvitation(ctx context.Context, invitationID, userID uint64) error {
	query, args := m.buildArchiveInvitationQuery(invitationID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
